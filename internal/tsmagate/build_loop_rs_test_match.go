//ff:func feature=gate type=helper control=sequence level=error lang=rust
//ff:what buildLoopRsTestMatch: the Rust D5 isolation step (confirmed strategy ③ — in-file mod injection + terminal rollback). Unlike Go/TS, Rust unit tests live INSIDE the source file, so there is no overlay/scratch-only path: it first copies the ORIGINAL source to a gitignored backing under .tsma/test (the rollback anchor + crash recovery), then writes the source file itself with the generated #[cfg(test)] mod injected (injectRsTestMod). Returns a TestMatch pointing at the source file (the in-file test target; RsRunner runs plain `cargo test`, RsChecker measures llvm-cov) plus the backing's root-relative path for finalizeRsBacking. A read/write failure propagates so prepareLoopRs surfaces TestFailed.
package tsmagate

import (
	"os"
	"path/filepath"

	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/park-jun-woo/tsma/internal/match"
)

// buildLoopRsTestMatch backs up the original source, injects the generated
// #[cfg(test)] mod into the source file, and returns the in-file TestMatch plus
// the backing's root-relative path. Precondition (plan §5.2c): a git-clean source
// tree — the backing guarantees recovery, and finalizeRsBacking restores on any
// non-terminal result.
func buildLoopRsTestMatch(p funcPayload, it *quest.Item, genSrc string) (match.TestMatch, string, error) {
	absSrc := filepath.Join(p.Root, p.Fn.File)
	original, err := os.ReadFile(absSrc)
	if err != nil {
		return match.TestMatch{}, "", err
	}
	backingRel := filepath.Join(".tsma", "test", "orig-"+pyBackingSlug(it.Key)+".rs")
	if err := writeTestFile(p.Root, backingRel, string(original)); err != nil {
		return match.TestMatch{}, "", err
	}
	injected := injectRsTestMod(string(original), genSrc)
	if err := writeTestFile(p.Root, p.Fn.File, injected); err != nil {
		return match.TestMatch{}, "", err
	}
	return match.TestMatch{Files: []string{p.Fn.File}, TestFuncs: nil}, backingRel, nil
}
