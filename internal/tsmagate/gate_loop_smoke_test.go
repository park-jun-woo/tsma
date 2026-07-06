//ff:func feature=gate type=test
//ff:what E2E loop 스모크(gate.md 경로): cli.Options에 Gate: GateOptions()를 명시해 그래프 판정 경로로 loop를 구동한다(기존 loop 테스트의 loopOpts는 Gate 미설정 = 평탄 경로라 이 경로는 여기서만 탄다). 부분커버 테스트→FAIL(branch-coverage-below-100 RootCause 코칭)→재시도→완전커버→PASS 수렴이 gate.md 경로에서 성립함을 stub 백엔드로 결정론적으로 고정한다.

package tsmagate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/cli"
	"github.com/park-jun-woo/reins/pkg/llm"
)

// loopOptsGate returns the production Loop config plus the gate.md graph gate,
// with the stub backend injected. Unlike loopOpts (flat path), this pins the
// GateVerdict path so the E2E smoke exercises the migrated judgment.
func loopOptsGate(backend llm.Backend) cli.Options {
	lo := LoopOptions()
	lo.LLM = backend
	return cli.Options{Loop: lo, Gate: GateOptions()}
}

// TestLoop_GateDoc_PartialThenFullConverges is the Phase003 E2E smoke: on the
// gate.md graph path, a partial test FAILs (coverage below 100), the loop feeds
// the branch-coverage RootCause coaching back, and the retry's full test PASSes.
func TestLoop_GateDoc_PartialThenFullConverges(t *testing.T) {
	root := classifyModule(t)
	chdirTo(t, root)
	session, out := sessionPaths(t)
	calls := 0
	backend := llm.CallFunc(func(system, user string) (string, error) {
		calls++
		if calls == 1 {
			return classifyPartialTest, nil // passes but coverage < 100 → FAIL
		}
		return classifyFullTest, nil // retry covers the missing branch → PASS
	})
	opts := loopOptsGate(backend)

	if _, err := runTsma(t, opts, session, out, "scan", root); err != nil {
		t.Fatalf("scan: %v", err)
	}
	got, err := runTsma(t, opts, session, out, "loop")
	if err != nil {
		t.Fatalf("loop: %v", err)
	}
	if calls != 2 {
		t.Fatalf("backend called %d times, want 2 (FAIL then PASS on the gate.md path)", calls)
	}
	if !strings.Contains(got, "-> FAIL") || !strings.Contains(got, "-> PASS") {
		t.Fatalf("loop output = %q, want a FAIL then a PASS", got)
	}
}
