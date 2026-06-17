//ff:func feature=gate type=helper control=sequence lang=python
//ff:what finalizePyBacking: the Python analogue of finalizeTSBacking. On a terminal pass (shouldMaterialize) it promotes the ORIGINAL src to the canonical test path; then, on every path, it removes the .tsma/test backing scratch. Reuses cleanupBacking (the extra overlay.json remove is a harmless no-op for Python).
package tsmagate

import "github.com/park-jun-woo/reins/pkg/quest"

// finalizePyBacking promotes the original src to disk on a terminal pass and
// always sweeps the .tsma/test backing afterward.
func finalizePyBacking(p funcPayload, it *quest.Item, m *measurement, src, backingRel string) {
	if shouldMaterialize(m, it) {
		promotePy(p, m, src)
	}
	cleanupBacking(p.Root, backingRel)
}
