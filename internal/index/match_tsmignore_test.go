package index

import "testing"

func TestMatchTsmIgnoreEmptyPatterns(t *testing.T) {
	if MatchTsmIgnore("src/main.go", "main.go", false, nil) {
		t.Error("expected false for nil patterns")
	}
	if MatchTsmIgnore("src/main.go", "main.go", false, []string{}) {
		t.Error("expected false for empty patterns")
	}
}

func TestMatchTsmIgnoreMatch(t *testing.T) {
	patterns := []string{"vendor/", "*.log"}
	if !MatchTsmIgnore("vendor", "vendor", true, patterns) {
		t.Error("expected vendor directory to match")
	}
	if !MatchTsmIgnore("src/debug.log", "debug.log", false, patterns) {
		t.Error("expected debug.log to match *.log")
	}
}

func TestMatchTsmIgnoreNoMatch(t *testing.T) {
	patterns := []string{"vendor/", "*.log"}
	if MatchTsmIgnore("src/main.go", "main.go", false, patterns) {
		t.Error("expected main.go to not match any pattern")
	}
	if MatchTsmIgnore("src", "src", true, patterns) {
		t.Error("expected src directory to not match any pattern")
	}
}
