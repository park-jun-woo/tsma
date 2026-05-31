package cli

import "testing"

// TestAttemptOutcome_threshold checks the retry/done boundary at maxAttempts,
// including max=3 (1,2 retry / 3 done), max=1 (immediate done), and max=5.
func TestAttemptOutcome_threshold(t *testing.T) {
	cases := []struct {
		name        string
		attempt     int
		maxAttempts int
		want        string
	}{
		{"max3_attempt1_retry", 1, 3, outcomeRetry},
		{"max3_attempt2_retry", 2, 3, outcomeRetry},
		{"max3_attempt3_done", 3, 3, outcomeDone},
		{"max3_attempt4_done", 4, 3, outcomeDone},
		{"max1_attempt1_done", 1, 1, outcomeDone},
		{"max5_attempt4_retry", 4, 5, outcomeRetry},
		{"max5_attempt5_done", 5, 5, outcomeDone},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := attemptOutcome(c.attempt, c.maxAttempts)
			if got != c.want {
				t.Fatalf("attemptOutcome(%d, %d) = %v, want %v",
					c.attempt, c.maxAttempts, got, c.want)
			}
		})
	}
}
