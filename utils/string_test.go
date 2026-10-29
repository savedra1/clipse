package utils

import (
	"testing"
	"unicode/utf8"
)

func TestShortenMultibyte(t *testing.T) {
	// 44 runes, 83 bytes: a byte-based cut lands mid-rune.
	s := "ускорить анимацию перехода по рабочим столам"

	if got := Shorten(s, 72); got != s {
		t.Errorf("Shorten(%q, 72) = %q, want unchanged", s, got)
	}

	got := Shorten(s, 20)
	if !utf8.ValidString(got) {
		t.Errorf("Shorten(%q, 20) = %q, not valid UTF-8", s, got)
	}
	if want := "ускорить анимацию..."; got != want {
		t.Errorf("Shorten(%q, 20) = %q, want %q", s, got, want)
	}
}
