package cli

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// TestEffectiveMaxAttempts covers effectiveMaxAttempts directly: an unset (or
// sub-1) MaxAttempts falls back to defaultMaxAttempts, while any value >= 1 is
// returned verbatim.
func TestEffectiveMaxAttempts(t *testing.T) {
	cases := []struct {
		name        string
		maxAttempts int
		want        int
	}{
		{"unset falls back to default", 0, defaultMaxAttempts},
		{"negative falls back to default", -5, defaultMaxAttempts},
		{"one is kept", 1, 1},
		{"explicit value is kept", 7, 7},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sess := &model.Session{MaxAttempts: tc.maxAttempts}
			if got := effectiveMaxAttempts(sess); got != tc.want {
				t.Fatalf("effectiveMaxAttempts(MaxAttempts=%d) = %d, want %d", tc.maxAttempts, got, tc.want)
			}
		})
	}
}
