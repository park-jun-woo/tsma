//ff:func feature=gate type=helper control=sequence level=error lang=typescript
//ff:what promoteTS: commits a passing TS loop test to its canonical path (match.CanonicalTestPath → foo.test.ts beside the source). Unlike the Go overlay promote, it writes the ORIGINAL sanitized src (with its original relative imports intact), not the import-rewritten backing — the canonical file lives beside the source so the original `./mod` imports are correct there. A write failure surfaces as TestFailed (a PASS that cannot persist is fed back as FAIL, never silent).
package tsmagate

import "github.com/park-jun-woo/tsma/internal/match"

// promoteTS writes the original sanitized src to the canonical test path. The
// canonical sits in the source directory, so the original relative imports
// resolve without the .tsma/test rewrite.
func promoteTS(p funcPayload, m *measurement, src string) {
	canonical := match.CanonicalTestPath(p.Lang, p.Fn.File)
	if canonical == "" {
		m.TestFailed = true
		m.FailOutput = "cannot derive canonical test path for " + p.Fn.File
		return
	}
	if err := writeTestFile(p.Root, canonical, src); err != nil {
		m.TestFailed = true
		m.FailOutput = err.Error()
	}
}
