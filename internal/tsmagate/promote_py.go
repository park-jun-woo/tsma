//ff:func feature=gate type=helper control=sequence level=error lang=python
//ff:what promotePy: commits a passing Python loop test to its canonical path (match.CanonicalTestPath → test_<file>.py beside the source). Like promoteTS, it accumulates the ORIGINAL sanitized src (with its original module imports intact), not the sys.path-injected backing — the canonical file lives beside the source so pytest's prepend import mode resolves the original imports without the .tsma/test header. BUG-002: promoteMerged accumulates by function (QualifiedName) marker block (`#` comments), so multiple functions in one .py file no longer overwrite each other. A write failure surfaces as TestFailed (a PASS that cannot persist is fed back as FAIL, never silent).
package tsmagate

import "github.com/park-jun-woo/tsma/internal/match"

// promotePy accumulates the original sanitized src into the canonical test path
// (BUG-002: per-function marker block, not a whole-file overwrite). The canonical
// sits in the source directory, so pytest resolves the original imports without
// the .tsma/test sys.path injection.
func promotePy(p funcPayload, m *measurement, src string) {
	canonical := match.CanonicalTestPath(p.Lang, p.Fn.File)
	if canonical == "" {
		m.TestFailed = true
		m.FailOutput = "cannot derive canonical test path for " + p.Fn.File
		return
	}
	if err := promoteMerged(p.Root, canonical, src, p.Fn.QualifiedName, p.Lang); err != nil {
		m.TestFailed = true
		m.FailOutput = err.Error()
	}
}
