//ff:func feature=gate type=test lang=typescript
//ff:what TS 루프 잔여 분기 단위테스트: prepareLoopTS의 backing 쓰기 실패
// (.tsma가 파일이라 MkdirAll 실패 → TestFailed+FailOutput, py 버전과 동형)와
// rewriteTSImports의 filepath.Rel 실패 폴백(상대 backing 경로 → 원 지정자 유지)을
// 결정적으로 덮는다.
package tsmagate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/park-jun-woo/tsma/internal/model"
)

// TestPrepareLoopTS_BuildFailureIsTestFailed mirrors the Python variant: with
// .tsma as a regular file the backing write inside buildLoopTSTestMatch fails,
// which prepareLoopTS must surface as TestFailed with the error output.
func TestPrepareLoopTS_BuildFailureIsTestFailed(t *testing.T) {
	emptyPath(t) // keep tidyTSSource an identity (no npx lookup side effects)
	root := t.TempDir()
	// .tsma as a file → backing write fails inside buildLoopTSTestMatch.
	if err := os.WriteFile(filepath.Join(root, ".tsma"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	p := funcPayload{Lang: "typescript", Root: root, Fn: model.Function{File: "src/math.ts", QualifiedName: "src.add"}}
	it := &quest.Item{Key: "src.add"}
	m := prepareLoopTS(it, p, []byte("test('x', () => {});\n"))
	if !m.TestFailed || m.FailOutput == "" {
		t.Fatalf("expected TestFailed with output on a write failure, got %+v", m)
	}
}

// TestRewriteTSImports_RelErrorKeepsSpecifier covers the filepath.Rel failure
// fallback: a relative backing dir cannot be made a base for the absolute
// target, so the original import specifier is kept verbatim.
func TestRewriteTSImports_RelErrorKeepsSpecifier(t *testing.T) {
	src := `import { f } from "./mod";`
	if got := rewriteTSImports(src, "/abs/src", "relative-backing"); got != src {
		t.Errorf("Rel failure must keep the specifier, got %q", got)
	}
}
