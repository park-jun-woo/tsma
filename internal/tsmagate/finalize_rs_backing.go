//ff:func feature=gate type=helper control=sequence lang=rust
//ff:what finalizeRsBacking: the Rust D5 terminal handler (strategy ③ rollback), the inverse of finalizeBacking/finalizeTSBacking. Because the generated test was injected into the REAL source file, a terminal pass (shouldMaterialize) KEEPS the injected source as the committed in-file test (nothing more to write), whereas every non-terminal result (retry, failure, vacuous pass, or panic via defer) RESTORES the original from the backing so the source tree returns to its pre-measurement state. On both paths the backing scratch under .tsma/test is swept (cleanupBacking). A restore error is recorded on the measurement (never silent) so a stuck mutation surfaces as a failure rather than corrupting the tree quietly.
package tsmagate

import "github.com/park-jun-woo/reins/pkg/quest"

// finalizeRsBacking keeps the injected source on a terminal pass and otherwise
// rolls it back from the backing, always sweeping the .tsma/test scratch.
func finalizeRsBacking(p funcPayload, it *quest.Item, m *measurement, backingRel string) {
	if !shouldMaterialize(m, it) {
		if err := restoreRsSource(p, backingRel); err != nil {
			m.TestFailed = true
			m.FailOutput = "failed to restore source after in-file test injection: " + err.Error()
		}
	}
	cleanupBacking(p.Root, backingRel)
}
