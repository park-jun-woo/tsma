//ff:func feature=gate type=helper control=sequence lang=typescript
//ff:what finalizeTSBacking: the TS analogue of finalizeBacking. On a terminal pass (shouldMaterialize) it promotes the ORIGINAL src to the canonical test path; then, on every path, it removes the .tsma/test backing scratch. Reuses cleanupBacking (the extra overlay.json remove is a harmless no-op for TS).
package tsmagate

import "github.com/park-jun-woo/reins/pkg/quest"

// finalizeTSBacking promotes the original src to disk on a terminal pass and
// always sweeps the .tsma/test backing afterward.
func finalizeTSBacking(p funcPayload, it *quest.Item, m *measurement, src, backingRel string) {
	if shouldMaterialize(m, it) {
		promoteTS(p, m, src)
	}
	cleanupBacking(p.Root, backingRel)
}
