package config

import (
	"encoding/json"
	"testing"
)

func TestDuplicatesAllowed(t *testing.T) {
	cases := []struct {
		name         string
		json         string
		text, images bool
	}{
		{"default", `{}`, false, false},
		{"images inherit", `{"allowDuplicates": true}`, true, true},
		{"images deduped", `{"allowDuplicates": true, "allowImageDuplicates": false}`, true, false},
		{"only images duplicated", `{"allowDuplicates": false, "allowImageDuplicates": true}`, false, true},
	}
	for _, tc := range cases {
		c := defaultConfig()
		if err := json.Unmarshal([]byte(tc.json), &c); err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		if got := c.DuplicatesAllowed(false); got != tc.text {
			t.Errorf("%s: text = %v, want %v", tc.name, got, tc.text)
		}
		if got := c.DuplicatesAllowed(true); got != tc.images {
			t.Errorf("%s: images = %v, want %v", tc.name, got, tc.images)
		}
	}
}
