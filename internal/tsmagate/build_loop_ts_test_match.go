//ff:func feature=gate type=helper control=sequence level=error lang=typescript
//ff:what buildLoopTSTestMatch: the TS isolation analogue of buildLoopTestMatch. Writes the generated test to .tsma/test/gen-<item>.test.ts (gitignored scratch) with its relative imports rewritten to reach the real module from there, and returns a TestMatch pointing at that backing file (no overlay — TS test location is free, so the runner runs the backing directly). The source tree is untouched during measurement; the original (un-rewritten) src is promoted to the canonical path only on PASS. A write failure propagates so Prepare surfaces TestFailed.
package tsmagate

import (
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/park-jun-woo/tsma/internal/match"
)

// buildLoopTSTestMatch writes the import-rewritten generated test to a backing
// file under .tsma/test and returns the TestMatch plus the backing's
// root-relative path (for the smell scan and cleanup). src is the sanitized,
// NOT-yet-rewritten source; the canonical promotion uses that original.
func buildLoopTSTestMatch(p funcPayload, it *quest.Item, src string) (match.TestMatch, string, error) {
	slug := strings.NewReplacer("/", "_", ".", "_", "*", "_", "(", "_", ")", "_", " ", "_", "[", "_", "]", "_").Replace(it.Key)
	backingRel := filepath.Join(".tsma", "test", "gen-"+slug+".test.ts")

	sourceDirAbs := filepath.Join(p.Root, filepath.Dir(p.Fn.File))
	backingDirAbs := filepath.Join(p.Root, ".tsma", "test")
	rewritten := rewriteTSImports(src, sourceDirAbs, backingDirAbs)

	if err := writeTestFile(p.Root, backingRel, rewritten); err != nil {
		return match.TestMatch{}, "", err
	}
	tm := match.TestMatch{Files: []string{backingRel}, TestFuncs: nil}
	return tm, backingRel, nil
}
