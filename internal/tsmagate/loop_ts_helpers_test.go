//ff:func feature=gate type=test
//ff:what Phase005a D5 unit tests for the TypeScript loop helpers: sanitizeTSSource (fence unwrap), rewriteTSImports (relative specifier re-pointing), buildLoopTSTestMatch (isolated backing under .tsma/test, source tree untouched), and finalizeTSBacking/promoteTS (canonical written with ORIGINAL imports only on a terminal pass; backing always swept). The jest-backed measureLoop is exercised only when jest is present, so these drive the deterministic, tooling-free pieces.
package tsmagate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/park-jun-woo/tsma/internal/coverage"
	"github.com/park-jun-woo/tsma/internal/model"
)

func TestSanitizeTSSource_UnwrapsFence(t *testing.T) {
	raw := "Here is the test:\n```typescript\nimport { add } from \"./math\";\ntest(\"x\", () => {});\n```\nDone."
	got := sanitizeTSSource(raw)
	if strings.Contains(got, "```") || strings.Contains(got, "Here is") || strings.Contains(got, "Done.") {
		t.Errorf("fence/prose not stripped: %q", got)
	}
	if !strings.Contains(got, `import { add } from "./math";`) {
		t.Errorf("test body lost: %q", got)
	}
}

func TestRewriteTSImports_RepointsRelative(t *testing.T) {
	src := `import { add } from "./math";
import { z } from "../util/z";
import { jest } from "jest";
const m = await import("./dyn");
`
	got := rewriteTSImports(src, "/proj/src", "/proj/.tsma/test")
	if !strings.Contains(got, `from "../../src/math"`) {
		t.Errorf("./math not repointed: %q", got)
	}
	if !strings.Contains(got, `from "../../util/z"`) {
		t.Errorf("../util/z not repointed: %q", got)
	}
	if !strings.Contains(got, `from "jest"`) {
		t.Errorf("bare specifier must be untouched: %q", got)
	}
	if !strings.Contains(got, `import("../../src/dyn")`) {
		t.Errorf("dynamic import not repointed: %q", got)
	}
}

// writeTSProj creates a temp project with the given files and returns its root.
func writeTSProj(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for rel, content := range files {
		abs := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(abs, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestBuildLoopTSTestMatch_IsolatesAndRewrites(t *testing.T) {
	root := writeTSProj(t, map[string]string{
		"src/math.ts": "export function add(a:number,b:number){return a+b;}\n",
	})
	p := funcPayload{Lang: "typescript", Root: root, Fn: model.Function{File: "src/math.ts", Name: "add"}}
	it := &quest.Item{Key: "src.add"}
	src := "import { add } from \"./math\";\ntest(\"add\", () => { expect(add(1,2)).toBe(3); });\n"

	tm, backingRel, err := buildLoopTSTestMatch(p, it, src)
	if err != nil {
		t.Fatalf("buildLoopTSTestMatch: %v", err)
	}
	if len(tm.Files) != 1 || tm.Files[0] != backingRel {
		t.Fatalf("TestMatch.Files = %v, want [%s]", tm.Files, backingRel)
	}
	if !strings.HasPrefix(backingRel, filepath.Join(".tsma", "test")) || !strings.HasSuffix(backingRel, ".test.ts") {
		t.Errorf("backing path not under .tsma/test: %s", backingRel)
	}
	data, err := os.ReadFile(filepath.Join(root, backingRel))
	if err != nil {
		t.Fatalf("backing not written: %v", err)
	}
	if !strings.Contains(string(data), `from "../../src/math"`) {
		t.Errorf("backing imports not rewritten to reach the real module: %s", data)
	}
	// The source tree must be untouched — no foo.test.ts beside the source.
	if _, err := os.Stat(filepath.Join(root, "src", "math.test.ts")); !os.IsNotExist(err) {
		t.Errorf("source tree must stay clean during measurement, stat err = %v", err)
	}
}

func TestFinalizeTSBacking_PromotesOnPassSweepsAlways(t *testing.T) {
	root := writeTSProj(t, map[string]string{
		"src/math.ts": "export function add(a:number,b:number){return a+b;}\n",
	})
	p := funcPayload{Lang: "typescript", Root: root, Fn: model.Function{File: "src/math.ts", Name: "add"}}
	it := &quest.Item{Key: "src.add", Tries: 0}
	originalSrc := "import { add } from \"./math\";\ntest(\"add\", () => {});\n"
	backingRel := filepath.Join(".tsma", "test", "gen-src_add.test.ts")
	if err := writeTestFile(root, backingRel, "rewritten backing"); err != nil {
		t.Fatal(err)
	}

	// Passing, fully-covered measurement → promote original to canonical.
	m := &measurement{Report: &coverage.Report{AllCovered: true}}
	finalizeTSBacking(p, it, m, originalSrc, backingRel)

	canonical := filepath.Join(root, "src", "math.test.ts")
	data, err := os.ReadFile(canonical)
	if err != nil {
		t.Fatalf("canonical not promoted: %v", err)
	}
	if !strings.Contains(string(data), `from "./math"`) {
		t.Errorf("canonical must carry the ORIGINAL relative import, got: %s", data)
	}
	if _, err := os.Stat(filepath.Join(root, backingRel)); !os.IsNotExist(err) {
		t.Errorf("backing must be swept, stat err = %v", err)
	}
}

func TestFinalizeTSBacking_FailDoesNotTouchSource(t *testing.T) {
	root := writeTSProj(t, map[string]string{
		"src/math.ts": "export function add(a:number,b:number){return a+b;}\n",
	})
	p := funcPayload{Lang: "typescript", Root: root, Fn: model.Function{File: "src/math.ts", Name: "add"}}
	it := &quest.Item{Key: "src.add", Tries: 0}
	backingRel := filepath.Join(".tsma", "test", "gen-src_add.test.ts")
	if err := writeTestFile(root, backingRel, "rewritten backing"); err != nil {
		t.Fatal(err)
	}

	// Failed measurement → no canonical, backing still swept (non-invasive).
	m := &measurement{TestFailed: true}
	finalizeTSBacking(p, it, m, "x", backingRel)

	if _, err := os.Stat(filepath.Join(root, "src", "math.test.ts")); !os.IsNotExist(err) {
		t.Errorf("a failed loop must not write into the source tree, stat err = %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, backingRel)); !os.IsNotExist(err) {
		t.Errorf("backing must be swept even on failure, stat err = %v", err)
	}
}
