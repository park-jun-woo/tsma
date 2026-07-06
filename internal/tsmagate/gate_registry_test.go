//ff:func feature=gate type=test
//ff:what gateRegistry 단위테스트: `검사`.<규칙> 15심볼이 전부 등록됐는지(중복 등록 에러로 존재 확인), gate.md가 이 레지스트리로 tangeul.Load 전량 검사(구문·미등록 심볼)를 통과하는지, 로드된 문서의 "제출 통과" 판정이 PASS(전체 커버)·FAIL(테스트 실패, RootCause=tests-must-pass — RuleMeta.ID 승계)로 흐르는지 고정한다.

package tsmagate

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/park-jun-woo/reins/pkg/tangeul"
	"github.com/park-jun-woo/tsma/internal/coverage"
)

// TestGateRegistry pins gateRegistry: all 15 `검사`.<X> symbols are registered
// (probed via the duplicate-registration error), gate.md whole-checks against
// the registry via tangeul.Load, and the loaded doc's "제출 통과" case judges a
// fully-covered measurement PASS and a failed-tests measurement FAIL with
// RootCause "tests-must-pass" (RuleMeta.ID inheritance through RulePred).
func TestGateRegistry(t *testing.T) {
	reg := gateRegistry()
	if reg == nil {
		t.Fatal("gateRegistry() = nil")
	}

	symbols := []string{
		"TestsMustPass", "BranchCoverageBelow100",
		"GoUnsafe", "GoReflectDynamic", "GoLinkname",
		"TsAsAny", "TsReflect", "TsOwnProperty",
		"JavaReflect", "JavaSetAccessible",
		"CsReflect", "CsReflectInfo",
		"RsUnsafe", "RsTransmute", "RsPtr",
	}
	for _, name := range symbols {
		// RegisterPred rejects a duplicate alias.name — an error here proves
		// the symbol is already registered.
		if err := reg.RegisterPred("검사", name, tangeul.Pred{}); err == nil {
			t.Errorf("검사.%s is not registered (duplicate probe did not error)", name)
		}
	}

	doc, err := tangeul.Load(gateDoc, "internal/tsmagate/gate.md", gateRegistry())
	if err != nil {
		t.Fatalf("tangeul.Load(gate.md) whole-check failed: %v", err)
	}

	pass, err := doc.GateVerdict("제출 통과", gate.Context{
		Submission: &measurement{Report: &coverage.Report{AllCovered: true, TotalPct: 100}},
	}, nil, nil)
	if err != nil {
		t.Fatalf("GateVerdict(covered) error: %v", err)
	}
	if pass.Outcome != quest.OutPass {
		t.Errorf("covered measurement: Outcome = %v, want OutPass (%+v)", pass.Outcome, pass)
	}

	fail, err := doc.GateVerdict("제출 통과", gate.Context{
		Submission: &measurement{TestFailed: true, FailOutput: "boom", FuncName: "F"},
	}, nil, nil)
	if err != nil {
		t.Fatalf("GateVerdict(test failed) error: %v", err)
	}
	if fail.Outcome != quest.OutFail {
		t.Errorf("failed measurement: Outcome = %v, want OutFail (%+v)", fail.Outcome, fail)
	}
	if fail.RootCause != "tests-must-pass" {
		t.Errorf("RootCause = %q, want tests-must-pass (RuleMeta.ID inheritance)", fail.RootCause)
	}
}
