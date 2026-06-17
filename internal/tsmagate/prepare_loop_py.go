//ff:func feature=gate type=helper control=sequence lang=python
//ff:what prepareLoopPy: the Python loop's non-invasive measurement pipeline (D5, the Python analogue of prepareLoopTS). Sanitizes+tidies the generated raw, writes a sys.path-injected backing test under .tsma/test (source tree untouched), runs+measures via the Python runner/checker, and only at a terminal pass promotes the ORIGINAL src to the canonical test_<file>.py. Python test location is free, so isolation needs no overlay — just the scratch dir + sys.path injection. No smell scan: Python smell is held (parent §7-3 lock, Phase005b §4). Returns the measurement for the gate rules.
package tsmagate

import "github.com/park-jun-woo/reins/pkg/quest"

// prepareLoopPy runs the Python loop's isolated measurement. A build/write
// failure or a non-pass run is surfaced as TestFailed via measureLoop; the
// backing scratch is always cleaned up, and the canonical test reaches disk only
// on a terminal pass (finalizePyBacking → promotePy).
func prepareLoopPy(it *quest.Item, p funcPayload, raw []byte) *measurement {
	m := &measurement{FuncName: p.Fn.QualifiedName}
	src := sanitizePySource(string(raw))
	tm, backingRel, err := buildLoopPyTestMatch(p, it, src)
	if err != nil {
		m.TestFailed, m.FailOutput = true, err.Error()
		return m
	}
	measureLoop(m, p, tm)
	finalizePyBacking(p, it, m, src, backingRel)
	return m
}
