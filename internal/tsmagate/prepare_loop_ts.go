//ff:func feature=gate type=helper control=sequence lang=typescript
//ff:what prepareLoopTS: the TS loop's non-invasive measurement pipeline (D5, the TS analogue of prepareLoopGo). Sanitizes+tidies the generated raw, writes an import-rewritten backing test under .tsma/test (source tree untouched), scans the backing for smells, runs+measures via the TS runner/checker, and only at a terminal pass promotes the ORIGINAL src to the canonical foo.test.ts. TS test location is free, so isolation needs no overlay — just the scratch dir + import rewrite. Returns the measurement for the gate rules.
package tsmagate

import "github.com/park-jun-woo/reins/pkg/quest"

// prepareLoopTS runs the TypeScript loop's isolated measurement. A build/write
// failure or a non-pass run is surfaced as TestFailed via measureLoop; the
// backing scratch is always cleaned up, and the canonical test reaches disk only
// on a terminal pass (finalizeTSBacking → promoteTS).
func prepareLoopTS(it *quest.Item, p funcPayload, raw []byte) *measurement {
	m := &measurement{FuncName: p.Fn.QualifiedName}
	src := sanitizeTSSource(string(raw))
	tm, backingRel, err := buildLoopTSTestMatch(p, it, src)
	if err != nil {
		m.TestFailed, m.FailOutput = true, err.Error()
		return m
	}
	m.Smells = scanTSSmells(p.Root, []string{backingRel})
	measureLoop(m, p, tm)
	finalizeTSBacking(p, it, m, src, backingRel)
	return m
}
