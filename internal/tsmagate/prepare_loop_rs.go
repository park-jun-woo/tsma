//ff:func feature=gate type=helper control=sequence lang=rust
//ff:what prepareLoopRs: the Rust loop's measurement pipeline (D5, confirmed strategy ③ — in-file #[cfg(test)] mod injection + terminal rollback; the prepareLoopGo mirror for Rust's in-file convention). Sanitizes+tidies the generated raw, checks for truncation via brace balance (C3: TestFailed+Truncated, nothing written), rejects a block declaring no #[test] fn (C1 — nothing would run), injects the mod into the REAL source file (backing the original first), scans the injected source for escape-hatch smells, and measures via the Rust runner/checker (plain `cargo test` + llvm-cov; vacuous-pass guard in measureLoop downgrades a 0%-coverage pass). finalizeRsBacking is deferred so the source is rolled back on every non-terminal result AND on a panic (panic-safe restore), and kept only at a terminal pass. cargo absence makes measureLoop a clean TestFailed (the documented e2e skip-gate); the injection/rollback/backing logic is pure and unit-tested.
package tsmagate

import "github.com/park-jun-woo/reins/pkg/quest"

// prepareLoopRs runs the Rust loop's in-file measurement. A truncated or
// test-less generation is rejected before any disk write; otherwise the generated
// mod is injected into the source file, measured, and rolled back (deferred,
// panic-safe) unless a terminal pass materializes it.
func prepareLoopRs(it *quest.Item, p funcPayload, raw []byte) *measurement {
	m := &measurement{FuncName: p.Fn.QualifiedName}
	src := sanitizeRsSource(string(raw))
	funcs, ok := parseRsTestFuncs(src)
	if !ok {
		m.TestFailed, m.Truncated = true, true
		return m
	}
	if len(funcs) == 0 {
		m.TestFailed = true
		m.FailExpected = "at least one #[test] fn that exercises " + p.Fn.QualifiedName
		m.FailOutput = "generated block declares no #[test] function, so nothing runs"
		return m
	}
	tm, backingRel, err := buildLoopRsTestMatch(p, it, src)
	if err != nil {
		m.TestFailed, m.FailOutput = true, err.Error()
		return m
	}
	defer finalizeRsBacking(p, it, m, backingRel)
	m.Smells = scanRsSmells(p.Root, []string{p.Fn.File})
	measureLoop(m, p, tm)
	return m
}
