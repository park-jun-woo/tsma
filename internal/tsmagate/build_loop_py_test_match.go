//ff:func feature=gate type=helper control=sequence level=error lang=python
//ff:what buildLoopPyTestMatch: the Python isolation analogue of buildLoopTSTestMatch. Writes the generated test to .tsma/test/gen_<slug>.py (gitignored scratch) with a sys.path header injected so it imports the real module from there, and returns a TestMatch pointing at that backing file (no overlay — Python test location is free, so the runner runs the backing directly). The source tree is untouched during measurement; the original (un-injected) src is promoted to the canonical test_<file>.py only on PASS. A write failure propagates so Prepare surfaces TestFailed.
package tsmagate

import (
	"path/filepath"

	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/park-jun-woo/tsma/internal/match"
)

// buildLoopPyTestMatch writes the sys.path-injected generated test to a backing
// file under .tsma/test and returns the TestMatch plus the backing's
// root-relative path (for cleanup). src is the sanitized, NOT-yet-injected
// source; the canonical promotion uses that original.
func buildLoopPyTestMatch(p funcPayload, it *quest.Item, src string) (match.TestMatch, string, error) {
	backingRel := filepath.Join(".tsma", "test", "gen_"+pyBackingSlug(it.Key)+".py")

	sourceDirAbs := filepath.Join(p.Root, filepath.Dir(p.Fn.File))
	injected := injectPySyspath(src, sourceDirAbs)

	if err := writeTestFile(p.Root, backingRel, injected); err != nil {
		return match.TestMatch{}, "", err
	}
	tm := match.TestMatch{Files: []string{backingRel}, TestFuncs: nil}
	return tm, backingRel, nil
}
