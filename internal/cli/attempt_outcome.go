//ff:func feature=cli type=helper control=sequence
//ff:what Decides retry vs auto-DONE for a still-partial function from its presentation count
package cli

// attemptOutcome decides whether a still-partial function should be retried or
// auto-accepted as DONE, based on how many times it has been presented. Reaching
// maxAttempts yields outcomeDone (best-effort accept); below it yields
// outcomeRetry (keep TODO). Shared by the measure and resurface paths so the
// threshold lives in exactly one place.
func attemptOutcome(attempt, maxAttempts int) string {
	if attempt >= maxAttempts {
		return outcomeDone
	}
	return outcomeRetry
}
