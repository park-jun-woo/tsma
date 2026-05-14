package cli

import "testing"

func TestMeasureResult_constants(t *testing.T) {
	// Verify outcome constants have expected values
	tests := []struct {
		name string
		val  string
		want string
	}{
		{"outcomeTestFail", outcomeTestFail, "test_fail"},
		{"outcomePass", outcomePass, "pass"},
		{"outcomeDone", outcomeDone, "done"},
		{"outcomeRetry", outcomeRetry, "retry"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.val != tt.want {
				t.Errorf("expected %q, got %q", tt.want, tt.val)
			}
		})
	}
}

func TestMeasureResult_struct(t *testing.T) {
	// Verify measureResult struct can be created and fields accessed
	r := measureResult{
		outcome:     outcomePass,
		mtime:       "2024-01-01T00:00:00Z",
		coveragePct: 85.5,
		attempt:     2,
		failOutput:  "",
		uncovered:   nil,
	}
	if r.outcome != outcomePass {
		t.Errorf("expected outcome=%q, got %q", outcomePass, r.outcome)
	}
	if r.coveragePct != 85.5 {
		t.Errorf("expected coveragePct=85.5, got %f", r.coveragePct)
	}
	if r.attempt != 2 {
		t.Errorf("expected attempt=2, got %d", r.attempt)
	}
}
