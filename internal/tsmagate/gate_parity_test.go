//ff:func feature=gate type=test
//ff:what gate.md 그래프 판정(tangeul GateVerdict)과 legacy 평탄 판정(gate.Evaluate(Rules))의 동등성 골든 테스트: 픽스처 매트릭스 F1~F8 + smell 13종 파라미터라이즈로 Outcome·RootCause·Facts 집합(FAIL+smell 동시 발화 시 smell Fact 보존)·EscalateOn 래치가 두 경로에서 동일함을 고정한다. Feedback은 구조(RootCause 키)만 보고 문자열 스냅샷은 하지 않는다. 불일치 시 교정은 gate.md 위상뿐 — rule_*.go는 불변.

package tsmagate

import (
	"reflect"
	"sort"
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/park-jun-woo/reins/pkg/tangeul"
	"github.com/park-jun-woo/tsma/internal/coverage"
	"github.com/park-jun-woo/tsma/internal/smell"
)

// parityDoc loads gate.md through the production registry once per test. A load
// error here would also be the whole-check regression (unregistered symbol), so
// the parity suite doubles as a load smoke.
func parityDoc(t *testing.T) *tangeul.Doc {
	t.Helper()
	doc, err := tangeul.Load(gateDoc, "internal/tsmagate/gate.md", gateRegistry())
	if err != nil {
		t.Fatalf("tangeul.Load(gate.md): %v", err)
	}
	return doc
}

// runBoth judges one measurement on both paths: legacy flat gate.Evaluate over
// the Rules() catalog, and the gate.md graph via GateVerdict. tsma has no ground
// needs (all tier-0), so snap/provider are nil.
func runBoth(t *testing.T, doc *tangeul.Doc, m *measurement) (legacy, graph quest.Verdict) {
	t.Helper()
	ctx := gate.Context{Submission: m}
	legacy = gate.Evaluate(New().Rules(), ctx)
	g, err := doc.GateVerdict("제출 통과", ctx, nil, nil)
	if err != nil {
		t.Fatalf("GateVerdict error: %v", err)
	}
	return legacy, g
}

// factKey serializes a Fact for order-independent set comparison (axis ③).
func factKey(f quest.Fact) string {
	return f.Rule + "\x00" + f.Where + "\x00" + f.Expected + "\x00" + f.Actual
}

// factsEqualSet reports whether two Fact slices are equal as sets (a Fact carries
// its own Rule ID, so order is not load-bearing — the plan compares Facts 집합).
func factsEqualSet(a, b []quest.Fact) bool {
	if len(a) != len(b) {
		return false
	}
	sa := append([]quest.Fact(nil), a...)
	sb := append([]quest.Fact(nil), b...)
	sort.Slice(sa, func(i, j int) bool { return factKey(sa[i]) < factKey(sa[j]) })
	sort.Slice(sb, func(i, j int) bool { return factKey(sb[i]) < factKey(sb[j]) })
	return reflect.DeepEqual(sa, sb)
}

// hasFactRule reports whether any Fact was stamped with the given rule ID.
func hasFactRule(facts []quest.Fact, rule string) bool {
	for _, f := range facts {
		if f.Rule == rule {
			return true
		}
	}
	return false
}

// assertEquiv pins the three compared axes between the two paths: Outcome (①),
// RootCause (②), and the Fact set (③). Feedback (④) is deliberately not string-
// compared — the graph path adds a walkthrough; equivalence is that both select
// the same RootCause key (already asserted), which is what RuleSystem indexes.
func assertEquiv(t *testing.T, name string, legacy, graph quest.Verdict) {
	t.Helper()
	if graph.Outcome != legacy.Outcome {
		t.Fatalf("%s: Outcome graph=%s legacy=%s", name, graph.Outcome, legacy.Outcome)
	}
	if graph.RootCause != legacy.RootCause {
		t.Fatalf("%s: RootCause graph=%q legacy=%q", name, graph.RootCause, legacy.RootCause)
	}
	if !factsEqualSet(legacy.Facts, graph.Facts) {
		t.Fatalf("%s: Facts set mismatch:\n legacy=%+v\n graph=%+v", name, legacy.Facts, graph.Facts)
	}
}

// assertGolden asserts the legacy path's current judgment (the golden) and then
// that the graph path is equivalent on all axes.
func assertGolden(t *testing.T, doc *tangeul.Doc, name string, m *measurement, wantOutcome quest.Outcome, wantRoot string) {
	t.Helper()
	legacy, graph := runBoth(t, doc, m)
	if legacy.Outcome != wantOutcome {
		t.Fatalf("%s: legacy Outcome = %s, want %s", name, legacy.Outcome, wantOutcome)
	}
	if legacy.RootCause != wantRoot {
		t.Fatalf("%s: legacy RootCause = %q, want %q", name, legacy.RootCause, wantRoot)
	}
	assertEquiv(t, name, legacy, graph)
}

// TestGateParity_Matrix is the F1~F8 golden parity table. Each row judges one
// measurement on both paths and asserts identical Outcome/RootCause/Facts.
func TestGateParity_Matrix(t *testing.T) {
	doc := parityDoc(t)

	// F1: tests failed → FAIL, tests-must-pass.
	assertGolden(t, doc, "F1/test-failed",
		&measurement{TestFailed: true, FailOutput: "boom", FuncName: "F"},
		quest.OutFail, "tests-must-pass")

	// F2: tests pass but coverage below 100 → FAIL, branch-coverage-below-100.
	f2 := &measurement{
		FuncName: "pkg.Fn",
		Report: &coverage.Report{AllCovered: false, TotalPct: 50,
			Uncovered: []coverage.UncoveredBranch{{File: "x.go", Line: 3}}},
	}
	assertGolden(t, doc, "F2/coverage-below", f2, quest.OutFail, "branch-coverage-below-100")

	// F3: 100% covered + one smell → REVIEW, the smell's ID.
	assertGolden(t, doc, "F3/covered-smell",
		&measurement{FuncName: "pkg.Fn",
			Report: &coverage.Report{AllCovered: true, TotalPct: 100},
			Smells: []smell.Finding{{Rule: "TS-REFL-001", File: "x_test.go", Line: 7, Note: `import "unsafe"`}}},
		quest.OutReview, "TS-REFL-001")

	// F4: fully clean → PASS, no RootCause, no Facts.
	assertGolden(t, doc, "F4/clean",
		&measurement{FuncName: "pkg.Fn", Report: &coverage.Report{AllCovered: true, TotalPct: 100}},
		quest.OutPass, "")

	// F5: test-failed AND a smell fire together → FAIL, tests-must-pass; the smell
	// Fact must survive (axis ③ crux — legacy records it, the graph must not
	// exclude it).
	f5 := &measurement{TestFailed: true, FailOutput: "boom", FuncName: "F",
		Smells: []smell.Finding{{Rule: "TS-REFL-001", File: "x_test.go", Line: 7, Note: `import "unsafe"`}}}
	legacy5, graph5 := runBoth(t, doc, f5)
	if legacy5.Outcome != quest.OutFail || legacy5.RootCause != "tests-must-pass" {
		t.Fatalf("F5: legacy = %s/%q, want FAIL/tests-must-pass", legacy5.Outcome, legacy5.RootCause)
	}
	assertEquiv(t, "F5/fail+smell", legacy5, graph5)
	if !hasFactRule(legacy5.Facts, "TS-REFL-001") {
		t.Fatalf("F5: legacy dropped the smell Fact: %+v", legacy5.Facts)
	}
	if !hasFactRule(graph5.Facts, "TS-REFL-001") {
		t.Fatalf("F5: graph dropped the smell Fact (over-exclusion): %+v", graph5.Facts)
	}

	// F6: coverage-below AND a smell fire together → FAIL, branch-coverage-below-100;
	// the smell Fact must survive.
	f6 := &measurement{FuncName: "pkg.Fn",
		Report: &coverage.Report{AllCovered: false, TotalPct: 50,
			Uncovered: []coverage.UncoveredBranch{{File: "x.go", Line: 3}}},
		Smells: []smell.Finding{{Rule: "TS-REFL-002", File: "x_test.go", Line: 9, Note: "reflect .MethodByName"}}}
	legacy6, graph6 := runBoth(t, doc, f6)
	if legacy6.Outcome != quest.OutFail || legacy6.RootCause != "branch-coverage-below-100" {
		t.Fatalf("F6: legacy = %s/%q, want FAIL/branch-coverage-below-100", legacy6.Outcome, legacy6.RootCause)
	}
	assertEquiv(t, "F6/coverage+smell", legacy6, graph6)
	if !hasFactRule(graph6.Facts, "TS-REFL-002") {
		t.Fatalf("F6: graph dropped the smell Fact (over-exclusion): %+v", graph6.Facts)
	}

	// F7: truncated submission. No hardcoded expectation — record the legacy
	// judgment as the golden and assert the graph matches it (equivalence, not
	// redesign).
	f7 := &measurement{TestFailed: true, Truncated: true, FuncName: "F"}
	legacy7, graph7 := runBoth(t, doc, f7)
	t.Logf("F7 golden (legacy): Outcome=%s RootCause=%q Facts=%d", legacy7.Outcome, legacy7.RootCause, len(legacy7.Facts))
	assertEquiv(t, "F7/truncated", legacy7, graph7)

	// F8: F2 reused — EscalateOn latch (axis ⑤). Both paths root-cause
	// branch-coverage-below-100, and that ID is the loop's EscalateOn signal.
	legacy8, graph8 := runBoth(t, doc, f2)
	if legacy8.RootCause != "branch-coverage-below-100" || graph8.RootCause != "branch-coverage-below-100" {
		t.Fatalf("F8: RootCause legacy=%q graph=%q, want branch-coverage-below-100 on both", legacy8.RootCause, graph8.RootCause)
	}
	escalated := false
	for _, id := range LoopOptions().EscalateOn {
		if id == "branch-coverage-below-100" {
			escalated = true
		}
	}
	if !escalated {
		t.Fatalf("F8: EscalateOn %v does not latch branch-coverage-below-100", LoopOptions().EscalateOn)
	}
}

// TestGateParity_Smells parameterizes F3 over all 13 TS-REFL rule IDs: each fires
// a single smell on a fully-covered submission and must yield REVIEW with that ID
// as RootCause on both paths.
func TestGateParity_Smells(t *testing.T) {
	doc := parityDoc(t)
	ids := []string{
		"TS-REFL-001", "TS-REFL-002", "TS-REFL-003",
		"TS-REFL-TS-001", "TS-REFL-TS-002", "TS-REFL-TS-003",
		"TS-REFL-JV-001", "TS-REFL-JV-002",
		"TS-REFL-CS-001", "TS-REFL-CS-002",
		"TS-REFL-RS-001", "TS-REFL-RS-002", "TS-REFL-RS-003",
	}
	for _, id := range ids {
		id := id
		t.Run(id, func(t *testing.T) {
			m := &measurement{FuncName: "pkg.Fn",
				Report: &coverage.Report{AllCovered: true, TotalPct: 100},
				Smells: []smell.Finding{{Rule: id, File: "t_test.go", Line: 1, Note: "n"}}}
			assertGolden(t, doc, "smell/"+id, m, quest.OutReview, id)
		})
	}
}
