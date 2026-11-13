package repository

import (
	"strings"
	"testing"
)

func TestTruncateKeepsUTF8Valid(t *testing.T) {
	value := "ООО «Весовая Тест»"
	truncated := truncate(value, 6)
	if !strings.HasPrefix(value, truncated) {
		t.Fatalf("truncated value %q must be a prefix of %q", truncated, value)
	}
	if len([]rune(truncated)) != 6 {
		t.Fatalf("expected 6 runes, got %d", len([]rune(truncated)))
	}
}
