package search

import (
	"math"
	"sort"
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
)

type Config struct {
	Engine          string          `json:"engine"`
	Algo            string          `json:"algo"`
	MatchMode       string          `json:"matchMode"`
	CaseSensitivity string          `json:"caseSensitivity"`
	Normalize       bool            `json:"normalize"`
	Tiebreak        []TiebreakEntry `json:"tiebreak"`
}

type TiebreakEntry struct {
	Key    string
	Bucket string
}

type ItemMeta struct {
	UseCount int
	LastUsed time.Time
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
		tiebreak = []TiebreakEntry{{Key: TiebreakScore}, {Key: TiebreakLength}, {Key: TiebreakIndex}}
	}
	useFrecency := metaLookup != nil && hasKey(tiebreak, TiebreakFrecency)
	matchFn := algo.FuzzyMatchV2
	if cfg.Algo == AlgoV1 {
		matchFn = algo.FuzzyMatchV1
	}
	if cfg.MatchMode == MatchModeExact {
		matchFn = algo.ExactMatchNaive
	}

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
				res, pos := matchFn(caseSensitive, cfg.Normalize, true, &text, pattern, true, slab)
				if res.Start < 0 {
					matched = false
					break
				}
				totalScore += res.Score
				if pos != nil {
					matchedPositions = append(matchedPositions, *pos...)
					for _, p := range *pos {
						if p < begin {
							begin = p
						}
						if p > end {
							end = p
						}
					}
				} else {
					for p := res.Start; p < res.End; p++ {
						matchedPositions = append(matchedPositions, p)
					}
					if res.Start < begin {
						begin = res.Start
					}
					if res.End-1 > end {
						end = res.End - 1
					}
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
	case "log2":
		if val <= 1 {
			return 0
		}
		return math.Floor(math.Log2(val))
	default:
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
	if m.UseCount == 0 {
		return 0
	}
	if m.LastUsed.IsZero() {
		return float64(m.UseCount)
	}
	age := max(now.Sub(m.LastUsed), 0)
	return float64(m.UseCount) * math.Exp(-float64(age)/float64(frecencyHalflife))
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
