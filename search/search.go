package search

import (
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/charmbracelet/bubbles/list"
	"github.com/junegunn/fzf/src/algo"
	"github.com/junegunn/fzf/src/util"
)

func init() {
	algo.Init("default")
}

const (
	EngineDefault = "default"
	EngineFzf     = "fzf"

	AlgoV1 = "v1"
	AlgoV2 = "v2"

	MatchModeFuzzy = "fuzzy"
	MatchModeExact = "exact"

	CaseSmart   = "smart"
	CaseRespect = "respect"
	CaseIgnore  = "ignore"

	TiebreakScore    = "score"
	TiebreakLength   = "length"
	TiebreakIndex    = "index"
	TiebreakFrecency = "frecency"
	TiebreakBegin    = "begin"
	TiebreakEnd      = "end"

	frecencyHalflife = 24 * time.Hour
	slab16Size       = 100 * 1024
	slab32Size       = 2048

	typoEditPenalty  = 4
	typoMaxTargetLen = 256

	scoreContiguousBonus = 1 << 13
)

type Config struct {
	Engine          string          `json:"engine"`
	Algo            string          `json:"algo"`
	MatchMode       string          `json:"matchMode"`
	CaseSensitivity string          `json:"caseSensitivity"`
	Normalize       bool            `json:"normalize"`
	TypoTolerance   bool            `json:"typoTolerance"`
	MaxScatter      int             `json:"maxScatter"`
	Tiebreak        []TiebreakEntry `json:"tiebreak"`
}

type TiebreakEntry struct {
	Key    string
	Bucket string
}

type ItemMeta struct {
	UseCount int
	LastUsed time.Time
	Recorded time.Time
}

type MetaLookup func(target string) ItemMeta

func Filter(cfg Config, metaLookup MetaLookup) func(term string, targets []string) []list.Rank {
	if cfg.Engine != EngineFzf {
		return list.DefaultFilter
	}
	return fzfFilter(cfg, metaLookup)
}

type rankWithScore struct {
	rank     list.Rank
	score    int
	length   int
	frecency float64
	begin    int
	end      int
}

func fzfFilter(cfg Config, metaLookup MetaLookup) func(string, []string) []list.Rank {
	tiebreak := cfg.Tiebreak
	if len(tiebreak) == 0 {
		tiebreak = []TiebreakEntry{{Key: TiebreakScore, Bucket: "32"}, {Key: TiebreakFrecency}, {Key: TiebreakLength}, {Key: TiebreakIndex}}
	}
	useFrecency := metaLookup != nil && hasKey(tiebreak, TiebreakFrecency)
	matchFn := algo.FuzzyMatchV2
	if cfg.Algo == AlgoV1 {
		matchFn = algo.FuzzyMatchV1
	}
	if cfg.MatchMode == MatchModeExact {
		matchFn = algo.ExactMatchNaive
	}
	typo := cfg.TypoTolerance && cfg.MatchMode != MatchModeExact

	return func(term string, targets []string) []list.Rank {
		slab := util.MakeSlab(slab16Size, slab32Size)
		term = strings.TrimSpace(term)
		if term == "" {
			out := make([]list.Rank, len(targets))
			for i := range targets {
				out[i] = list.Rank{Index: i}
			}
			return out
		}

		tokens := strings.Fields(term)
		now := time.Now()

		results := make([]rankWithScore, 0, len(targets))
		for idx, target := range targets {
			totalScore := 0
			var matchedPositions []int
			begin, end := math.MaxInt, -1
			matched := true
			for _, tok := range tokens {
				caseSensitive := isCaseSensitive(cfg.CaseSensitivity, tok)
				text := util.ToChars([]byte(target))
				patternStr := tok
				if !caseSensitive {
					patternStr = strings.ToLower(patternStr)
				}
				pattern := []rune(patternStr)
				if cfg.Normalize {
					pattern = algo.NormalizeRunes(pattern)
				}

				tokScore := math.MinInt
				var tokPositions []int
				tokBegin, tokEnd := math.MaxInt, -1

				ptext := util.ToChars([]byte(string(pattern)))
				perfect, _ := matchFn(caseSensitive, cfg.Normalize, true, &ptext, pattern, true, slab)

				res, pos := matchFn(caseSensitive, cfg.Normalize, true, &text, pattern, true, slab)
				if res.Start >= 0 {
					tokScore = res.Score
					if pos != nil {
						tokPositions = append(tokPositions, *pos...)
						for _, p := range *pos {
							if p < tokBegin {
								tokBegin = p
							}
							if p > tokEnd {
								tokEnd = p
							}
						}
					} else {
						for p := res.Start; p < res.End; p++ {
							tokPositions = append(tokPositions, p)
						}
						tokBegin, tokEnd = res.Start, res.End-1
					}
					if cfg.MaxScatter > 0 && tokEnd-tokBegin+1-len(pattern) > cfg.MaxScatter {
						tokScore, tokBegin, tokEnd, tokPositions = math.MinInt, math.MaxInt, -1, nil
					}
				}

				if typo {
					if ts, tb, te, tp, ok := typoTokenMatch(pattern, target, cfg.Normalize, perfect.Score); ok && ts > tokScore {
						tokScore, tokBegin, tokEnd, tokPositions = ts, tb, te, tp
					}
				}

				if cb, ce, ok := contiguousMatch(pattern, target, cfg.Normalize, caseSensitive); ok {
					tokScore = scoreContiguousBonus
					tokBegin, tokEnd = cb, ce
					tokPositions = tokPositions[:0]
					for p := cb; p <= ce; p++ {
						tokPositions = append(tokPositions, p)
					}
				}

				if tokScore == math.MinInt {
					matched = false
					break
				}
				totalScore += tokScore
				matchedPositions = append(matchedPositions, tokPositions...)
				if tokBegin < begin {
					begin = tokBegin
				}
				if tokEnd > end {
					end = tokEnd
				}
			}
			if !matched {
				continue
			}
			fr := 0.0
			if useFrecency {
				m := metaLookup(target)
				fr = frecencyScore(m, now)
			}
			results = append(results, rankWithScore{
				rank:     list.Rank{Index: idx, MatchedIndexes: matchedPositions},
				score:    totalScore,
				length:   len(target),
				frecency: fr,
				begin:    begin,
				end:      end,
			})
		}

		sort.SliceStable(results, func(i, j int) bool {
			return lessByTiebreak(results[i], results[j], tiebreak)
		})

		out := make([]list.Rank, len(results))
		for i, r := range results {
			out[i] = r.rank
		}
		return out
	}
}

func lessByTiebreak(a, b rankWithScore, tiebreak []TiebreakEntry) bool {
	for _, tb := range tiebreak {
		switch tb.Key {
		case TiebreakScore:
			av, bv := bucketize(float64(a.score), tb.Bucket), bucketize(float64(b.score), tb.Bucket)
			if av != bv {
				return av > bv
			}
		case TiebreakLength:
			av, bv := bucketize(float64(a.length), tb.Bucket), bucketize(float64(b.length), tb.Bucket)
			if av != bv {
				return av < bv
			}
		case TiebreakIndex:
			if a.rank.Index != b.rank.Index {
				return a.rank.Index < b.rank.Index
			}
		case TiebreakFrecency:
			av, bv := bucketize(a.frecency, tb.Bucket), bucketize(b.frecency, tb.Bucket)
			if av != bv {
				return av > bv
			}
		case TiebreakBegin:
			av, bv := bucketize(float64(a.begin), tb.Bucket), bucketize(float64(b.begin), tb.Bucket)
			if av != bv {
				return av < bv
			}
		case TiebreakEnd:
			av, bv := bucketize(float64(a.end), tb.Bucket), bucketize(float64(b.end), tb.Bucket)
			if av != bv {
				return av > bv
			}
		}
	}
	return false
}

func bucketize(val float64, strategy string) float64 {
	switch strategy {
	case "", "raw":
		return val
	case "log2":
		if val <= 1 {
			return 0
		}
		return math.Floor(math.Log2(val))
	default:
		// A numeric strategy is a linear quantization width: values are
		// collapsed into floor(val/width) buckets, so anything within one
		// width ties and a later tiebreak (e.g. frecency) decides. fzf's
		// per-character score unit is 16, so a width of ~16 gives "gentle"
		// near-tie breaking. Unparsable / non-positive widths fall back to raw.
		if w, err := strconv.Atoi(strategy); err == nil && w > 0 {
			return math.Floor(val / float64(w))
		}
		return val
	}
}

func hasKey(entries []TiebreakEntry, key string) bool {
	for _, e := range entries {
		if e.Key == key {
			return true
		}
	}
	return false
}

func frecencyScore(m ItemMeta, now time.Time) float64 {
	last := m.Recorded
	if m.LastUsed.After(last) {
		last = m.LastUsed
	}
	if last.IsZero() {
		return 0
	}
	age := max(now.Sub(last), 0)
	return float64(1+m.UseCount) * math.Exp(-float64(age)/float64(frecencyHalflife))
}

func typoMaxEdits(n int) int {
	switch {
	case n < 3:
		return 0
	case n <= 4:
		return 1
	default:
		return 2
	}
}

func typoTokenMatch(pattern []rune, target string, normalize bool, perfectScore int) (score, begin, end int, positions []int, ok bool) {
	maxEd := typoMaxEdits(len(pattern))
	if maxEd == 0 || len(target) > typoMaxTargetLen {
		return 0, 0, 0, nil, false
	}
	runes := []rune(target)
	bestD := maxEd + 1
	bestStart, bestLen := -1, 0
	for i := 0; i < len(runes); {
		if !isWordRune(runes[i]) {
			i++
			continue
		}
		j := i
		for j < len(runes) && isWordRune(runes[j]) {
			j++
		}
		word := make([]rune, 0, j-i)
		for k := i; k < j; k++ {
			word = append(word, unicode.ToLower(runes[k]))
		}
		if normalize {
			word = algo.NormalizeRunes(word)
		}
		if absInt(len(word)-len(pattern)) <= maxEd {
			if d := damerauLevenshtein(pattern, word, maxEd); d < bestD {
				bestD, bestStart, bestLen = d, i, j-i
			}
		}
		if len(word) > len(pattern) {
			if d := damerauLevenshtein(pattern, word[:len(pattern)], maxEd); d < bestD {
				bestD, bestStart, bestLen = d, i, len(pattern)
			}
		}
		i = j
	}
	if bestStart < 0 || bestD > maxEd {
		return 0, 0, 0, nil, false
	}
	positions = make([]int, bestLen)
	for k := range positions {
		positions[k] = bestStart + k
	}
	score = perfectScore - typoEditPenalty*bestD
	return score, bestStart, bestStart + bestLen - 1, positions, true
}

func damerauLevenshtein(a, b []rune, maxD int) int {
	la, lb := len(a), len(b)
	if absInt(la-lb) > maxD {
		return maxD + 1
	}
	prev2 := make([]int, lb+1)
	prev := make([]int, lb+1)
	curr := make([]int, lb+1)
	for j := 0; j <= lb; j++ {
		prev[j] = j
	}
	for i := 1; i <= la; i++ {
		curr[0] = i
		rowMin := curr[0]
		for j := 1; j <= lb; j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			m := min(prev[j]+1, curr[j-1]+1)
			m = min(m, prev[j-1]+cost)
			if i > 1 && j > 1 && a[i-1] == b[j-2] && a[i-2] == b[j-1] {
				m = min(m, prev2[j-2]+1)
			}
			curr[j] = m
			if m < rowMin {
				rowMin = m
			}
		}
		if rowMin > maxD {
			return maxD + 1
		}
		prev2, prev, curr = prev, curr, prev2
	}
	return prev[lb]
}

func contiguousMatch(pattern []rune, target string, normalize, caseSensitive bool) (begin, end int, ok bool) {
	if len(pattern) == 0 {
		return 0, 0, false
	}
	hay := []rune(target)
	if len(hay) < len(pattern) {
		return 0, 0, false
	}
	for i := range hay {
		if !caseSensitive {
			hay[i] = unicode.ToLower(hay[i])
		}
	}
	if normalize {
		hay = algo.NormalizeRunes(hay)
	}
	for i := 0; i+len(pattern) <= len(hay); i++ {
		match := true
		for j := range pattern {
			if hay[i+j] != pattern[j] {
				match = false
				break
			}
		}
		if match {
			return i, i + len(pattern) - 1, true
		}
	}
	return 0, 0, false
}

func isWordRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func isCaseSensitive(mode, pattern string) bool {
	switch mode {
	case CaseRespect:
		return true
	case CaseIgnore:
		return false
	default: // smart
		for _, r := range pattern {
			if unicode.IsUpper(r) {
				return true
			}
		}
		return false
	}
}
