package config

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/savedra1/clipse/config"
)

func Test(_ *testing.T) {}

func TestTiebreakListUnmarshalMixed(t *testing.T) {
	raw := []byte(`["score","length",{"key":"frecency","bucket":"log2"},"index"]`)
	var got config.TiebreakList
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	want := config.TiebreakList{
		{Key: "score"},
		{Key: "length"},
		{Key: "frecency", Bucket: "log2"},
		{Key: "index"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestTiebreakListMarshalPreservesShape(t *testing.T) {
	in := config.TiebreakList{
		{Key: "score"},
		{Key: "frecency", Bucket: "log2"},
		{Key: "index"},
	}
	out, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	want := `["score",{"key":"frecency","bucket":"log2"},"index"]`
	if string(out) != want {
		t.Errorf("got %s, want %s", string(out), want)
	}
}
