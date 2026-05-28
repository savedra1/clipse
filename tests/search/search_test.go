package search_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/list"

	"github.com/savedra1/clipse/search"
)

func TestDefaultEngineMatchesListDefaultFilter(t *testing.T) {
	targets := []string{
		"git commit -m fix",
		"go test ./...",
		"git checkout main",
	}
	terms := []string{"git", "go", "gx", "CHECKOUT"}

	cfg := search.Config{Engine: search.EngineDefault}
	filter := search.Filter(cfg, nil)

	for _, term := range terms {
		got := filter(term, targets)
		want := list.DefaultFilter(term, targets)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("term %q: default engine diverged from list.DefaultFilter\n got=%+v\nwant=%+v", term, got, want)
		}
	}
}

func TestFzfRanksWordBoundaryAbove(t *testing.T) {
	targets := []string{
		"git commit",
		"go compile output",
		"git checkout origin",
	}
	cfg := search.Config{
		Engine:    search.EngineFzf,
		Algo:      search.AlgoV2,
		Normalize: true,
		Tiebreak:  []search.TiebreakEntry{{Key: search.TiebreakScore}, {Key: search.TiebreakLength}, {Key: search.TiebreakIndex}},
	}
	filter := search.Filter(cfg, nil)

	ranks := filter("gco", targets)
	if len(ranks) == 0 {
		t.Fatal("expected at least one match for 'gco'")
	}
	top := targets[ranks[0].Index]
	if top != "git checkout origin" && top != "go compile output" {
		t.Errorf("expected a word-boundary match at top, got %q", top)
	}
	for _, r := range ranks {
		if targets[r.Index] == "git commit" && targets[ranks[0].Index] == "git commit" {
			t.Errorf("unexpected: 'git commit' ranked top for pattern 'gco'")
		}
	}
}

func TestFzfMultiTermAnd(t *testing.T) {
	targets := []string{
		"git commit",
		"git checkout main",
		"go test",
	}
	cfg := search.Config{Engine: search.EngineFzf, Algo: search.AlgoV2, Normalize: true}
	filter := search.Filter(cfg, nil)

	ranks := filter("git ch", targets)
	if len(ranks) != 1 || targets[ranks[0].Index] != "git checkout main" {
		t.Errorf("expected only 'git checkout main', got %v", ranks)
	}
}

func TestFzfSmartCase(t *testing.T) {
	targets := []string{"Hello World", "hello there"}
	cfg := search.Config{Engine: search.EngineFzf, Algo: search.AlgoV2, CaseSensitivity: search.CaseSmart, Normalize: true}
	filter := search.Filter(cfg, nil)

	if ranks := filter("hello", targets); len(ranks) != 2 {
		t.Errorf("smart case lowercase: expected 2 matches, got %d", len(ranks))
	}
	ranks := filter("Hello", targets)
	if len(ranks) != 1 || targets[ranks[0].Index] != "Hello World" {
		t.Errorf("smart case mixed: expected only 'Hello World', got %+v", ranks)
	}
}

func TestFzfNormalize(t *testing.T) {
	targets := []string{"café au lait", "tea"}
	cfg := search.Config{Engine: search.EngineFzf, Algo: search.AlgoV2, Normalize: true}
	filter := search.Filter(cfg, nil)

	ranks := filter("cafe", targets)
	if len(ranks) == 0 {
		t.Errorf("normalize=true: expected 'cafe' to match 'café au lait'")
	}
}

func TestFrecencyTiebreak(t *testing.T) {
	targets := []string{"foo bar", "foo baz"}
	now := time.Now()
	meta := map[string]search.ItemMeta{
		"foo bar": {UseCount: 1, LastUsed: now.Add(-48 * time.Hour)},
		"foo baz": {UseCount: 10, LastUsed: now.Add(-1 * time.Hour)},
	}
	lookup := func(t string) search.ItemMeta { return meta[t] }

	cfg := search.Config{
		Engine:    search.EngineFzf,
		Algo:      search.AlgoV2,
		Normalize: true,
		Tiebreak:  []search.TiebreakEntry{{Key: search.TiebreakFrecency}, {Key: search.TiebreakIndex}},
	}
	filter := search.Filter(cfg, lookup)

	ranks := filter("foo", targets)
	if len(ranks) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(ranks))
	}
	if targets[ranks[0].Index] != "foo baz" {
		t.Errorf("frecency tiebreak: expected 'foo baz' first, got %q", targets[ranks[0].Index])
	}
}

func TestFrecencyDisabledWhenLookupNil(t *testing.T) {
	targets := []string{"foo bar", "foo baz"}
	cfg := search.Config{
		Engine:    search.EngineFzf,
		Algo:      search.AlgoV2,
		Normalize: true,
		Tiebreak:  []search.TiebreakEntry{{Key: search.TiebreakFrecency}, {Key: search.TiebreakIndex}},
	}
	filter := search.Filter(cfg, nil)
	ranks := filter("foo", targets)
	if len(ranks) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(ranks))
	}
	if ranks[0].Index != 0 {
		t.Errorf("nil lookup: expected index 0 first, got %d", ranks[0].Index)
	}
}

func TestFrecencyBucketLog2LetsLaterTiebreakDecide(t *testing.T) {
	targets := []string{"foo bar", "foo baz"}
	now := time.Now()
	meta := map[string]search.ItemMeta{
		"foo bar": {UseCount: 100, LastUsed: now},
		"foo baz": {UseCount: 105, LastUsed: now},
	}
	lookup := func(t string) search.ItemMeta { return meta[t] }

	cfg := search.Config{
		Engine:    search.EngineFzf,
		Algo:      search.AlgoV2,
		Normalize: true,
		Tiebreak: []search.TiebreakEntry{
			{Key: search.TiebreakScore},
			{Key: search.TiebreakFrecency, Bucket: "log2"},
			{Key: search.TiebreakIndex},
		},
	}
	ranks := search.Filter(cfg, lookup)("foo", targets)
	if len(ranks) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(ranks))
	}
	if targets[ranks[0].Index] != "foo bar" {
		t.Errorf("log2 bucket should tie frecencies 100 vs 105, letting index decide (foo bar first); got %q", targets[ranks[0].Index])
	}

	cfg.Tiebreak[1].Bucket = ""
	ranks = search.Filter(cfg, lookup)("foo", targets)
	if targets[ranks[0].Index] != "foo baz" {
		t.Errorf("unbucketed: expected 'foo baz' to win on frecency, got %q", targets[ranks[0].Index])
	}
}

func TestBeginTiebreak(t *testing.T) {
	targets := []string{"world hello", "hello world"}
	cfg := search.Config{
		Engine:    search.EngineFzf,
		Algo:      search.AlgoV2,
		Normalize: true,
		Tiebreak:  []search.TiebreakEntry{{Key: search.TiebreakBegin}, {Key: search.TiebreakIndex}},
	}
	ranks := search.Filter(cfg, nil)("hello", targets)
	if len(ranks) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(ranks))
	}
	if targets[ranks[0].Index] != "hello world" {
		t.Errorf("begin tiebreak: expected 'hello world' first (match at position 0), got %q", targets[ranks[0].Index])
	}
}

func TestEndTiebreak(t *testing.T) {
	targets := []string{"hello world", "world hello"}
	cfg := search.Config{
		Engine:    search.EngineFzf,
		Algo:      search.AlgoV2,
		Normalize: true,
		Tiebreak:  []search.TiebreakEntry{{Key: search.TiebreakEnd}, {Key: search.TiebreakIndex}},
	}
	ranks := search.Filter(cfg, nil)("hello", targets)
	if len(ranks) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(ranks))
	}
	if targets[ranks[0].Index] != "world hello" {
		t.Errorf("end tiebreak: expected 'world hello' first (match closer to tail), got %q", targets[ranks[0].Index])
	}
}

func TestBeginBucketLog2LetsLaterTiebreakDecide(t *testing.T) {
	// 'x' is at byte 2 in "aax" and byte 3 in "aaax" — both log2-bucket to 1.
	targets := []string{"aaax", "aax"}
	cfg := search.Config{
		Engine:    search.EngineFzf,
		Algo:      search.AlgoV2,
		Normalize: true,
		Tiebreak: []search.TiebreakEntry{
			{Key: search.TiebreakBegin, Bucket: "log2"},
			{Key: search.TiebreakIndex},
		},
	}
	ranks := search.Filter(cfg, nil)("x", targets)
	if len(ranks) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(ranks))
	}
	if targets[ranks[0].Index] != "aaax" {
		t.Errorf("log2 bucket should tie begins 2 vs 3; index should decide (aaax first), got %q", targets[ranks[0].Index])
	}

	cfg.Tiebreak[0].Bucket = ""
	ranks = search.Filter(cfg, nil)("x", targets)
	if targets[ranks[0].Index] != "aax" {
		t.Errorf("unbucketed: expected 'aax' (begin 2 < 3) first, got %q", targets[ranks[0].Index])
	}
}

func TestFzfEmptyTerm(t *testing.T) {
	targets := []string{"a", "b", "c"}
	cfg := search.Config{Engine: search.EngineFzf, Algo: search.AlgoV2}
	filter := search.Filter(cfg, nil)
	ranks := filter("", targets)
	if len(ranks) != 3 {
		t.Errorf("empty term should pass all items, got %d", len(ranks))
	}
}

// TestScoreBucketAbsorbsBoundaryNoiseSoLengthDecides guards against a fine score
// bucket letting fzf's boundary-bonus noise override length. fzf scores a match
// preceded by whitespace higher than the same match preceded by '[', so with a
// narrow bucket a long log line that merely contains "git" can outrank a short,
// focused "[git]" entry — and length, a later tiebreak, cannot rescue it across
// buckets. A width-32 score bucket collapses that boundary-bonus spread into one
// bucket so length sorts clean matches short-first, while genuinely scattered
// matches still fall into a lower bucket.
func TestScoreBucketAbsorbsBoundaryNoiseSoLengthDecides(t *testing.T) {
	longLog := "deploy notes: cut a release candidate\n  remember to checkout main before tagging\n" +
		"  then run the pipeline ............................................"
	targets := []string{
		longLog,                // clean " main" but long -> max score (whitespace boundary)
		"[main] release notes", // clean but '[' boundary -> slightly lower score, short
		"main branch",          // clean prefix, shortest
		"magician",             // scattered m..a..i..n, shortest overall -> lower score bucket
	}
	cfg := search.Config{
		Engine:    search.EngineFzf,
		Algo:      search.AlgoV2,
		Normalize: true,
		Tiebreak: []search.TiebreakEntry{
			{Key: search.TiebreakScore, Bucket: "32"},
			{Key: search.TiebreakLength},
			{Key: search.TiebreakIndex},
		},
	}
	ranks := search.Filter(cfg, nil)("main", targets)
	order := make([]string, len(ranks))
	for i, r := range ranks {
		order[i] = targets[r.Index]
	}
	// Short clean matches must come before the long clean match, even though the
	// long one has a higher raw fzf score (whitespace boundary beats '[').
	posLong, posBracket, posPrefix := indexOf(order, longLog), indexOf(order, "[main] release notes"), indexOf(order, "main branch")
	if posPrefix >= posBracket || posBracket >= posLong {
		t.Errorf("expected short clean matches before long log; got order:\n  %q", order)
	}
	// "magician" is the shortest target but a scattered match: a coarse bucket
	// (e.g. log2) would float it to the top, but width-32 keeps it below the
	// clean matches.
	if posScatter := indexOf(order, "magician"); posScatter < posLong {
		t.Errorf("scattered match should rank below clean matches; got order:\n  %q", order)
	}
}

func indexOf(s []string, v string) int {
	for i, x := range s {
		if x == v {
			return i
		}
	}
	return -1
}
