package search

import (
	"math"
	"slices"
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

	CaseSmart   = "smart"
	CaseRespect = "respect"
	CaseIgnore  = "ignore"

	TiebreakScore    = "score"
	TiebreakLength   = "length"
	TiebreakIndex    = "index"
	TiebreakFrecency = "frecency"

	frecencyHalflife = 24 * time.Hour
	slab16Size       = 100 * 1024
	slab32Size       = 2048
)

type Config struct {
	Engine          string   `json:"engine"`
	Algo            string   `json:"algo"`
	CaseSensitivity string   `json:"caseSensitivity"`
	Normalize       bool     `json:"normalize"`
	Tiebreak        []string `json:"tiebreak"`
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
}

func fzfFilter(cfg Config, metaLookup MetaLookup) func(string, []string) []list.Rank {
	useFrecency := metaLookup != nil && containsStr(cfg.Tiebreak, TiebreakFrecency)
	tiebreak := cfg.Tiebreak
	if len(tiebreak) == 0 {
		tiebreak = []string{TiebreakScore, TiebreakLength, TiebreakIndex}
	}
	matchFn := algo.FuzzyMatchV2
	if cfg.Algo == AlgoV1 {
		matchFn = algo.FuzzyMatchV1
	}
	slab := util.MakeSlab(slab16Size, slab32Size)

	return func(term string, targets []string) []list.Rank {
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
				} else {
					for p := res.Start; p < res.End; p++ {
						matchedPositions = append(matchedPositions, p)
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

func lessByTiebreak(a, b rankWithScore, tiebreak []string) bool {
	for _, tb := range tiebreak {
		switch tb {
		case TiebreakScore:
			if a.score != b.score {
				return a.score > b.score
			}
		case TiebreakLength:
			if a.length != b.length {
				return a.length < b.length
			}
		case TiebreakIndex:
			if a.rank.Index != b.rank.Index {
				return a.rank.Index < b.rank.Index
			}
		case TiebreakFrecency:
			if a.frecency != b.frecency {
				return a.frecency > b.frecency
			}
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

func containsStr(ss []string, s string) bool {
	return slices.Contains(ss, s)
}
