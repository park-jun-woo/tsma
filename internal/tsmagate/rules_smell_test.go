//ff:func feature=gate type=test
//ff:what TS-REFL 룰 통합테스트(§8): measurement를 직접 구성해 d.Rules()를 gate.Evaluate로 평가, (a) 100%+smell→REVIEW, (b) 100%+clean→PASS, (c) 부분커버+smell→FAIL(커버리지 우선, smell 묻힘) verdict를 검증한다. 룰 카탈로그가 Fail 2 + Review 3 = 5개임도 확인.

package tsmagate

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/park-jun-woo/tsma/internal/coverage"
	"github.com/park-jun-woo/tsma/internal/smell"
)

func evalVerdict(m *measurement) quest.Verdict {
	d := New()
	return gate.Evaluate(d.Rules(), gate.Context{Submission: m})
}

func TestRules_CleanFullCoverage_Pass(t *testing.T) {
	m := &measurement{
		FuncName: "pkg.Fn",
		Report:   &coverage.Report{AllCovered: true, TotalPct: 100},
	}
	if got := evalVerdict(m).Outcome; got != quest.OutPass {
		t.Fatalf("Outcome = %s, want PASS", got)
	}
}

func TestRules_FullCoverageWithSmell_Review(t *testing.T) {
	m := &measurement{
		FuncName: "pkg.Fn",
		Report:   &coverage.Report{AllCovered: true, TotalPct: 100},
		Smells: []smell.Finding{
			{Rule: "TS-REFL-001", File: "x_test.go", Line: 7, Note: `import "unsafe"`},
		},
	}
	v := evalVerdict(m)
	if v.Outcome != quest.OutReview {
		t.Fatalf("Outcome = %s, want REVIEW", v.Outcome)
	}
	if v.RootCause != "TS-REFL-001" {
		t.Fatalf("RootCause = %q, want TS-REFL-001", v.RootCause)
	}
}

func TestRules_PartialCoverageWithSmell_Fail(t *testing.T) {
	m := &measurement{
		FuncName: "pkg.Fn",
		Report: &coverage.Report{
			AllCovered: false,
			TotalPct:   50,
			Uncovered:  []coverage.UncoveredBranch{{File: "x.go", Line: 3}},
		},
		Smells: []smell.Finding{
			{Rule: "TS-REFL-002", File: "x_test.go", Line: 9, Note: "reflect dynamic .MethodByName"},
		},
	}
	v := evalVerdict(m)
	if v.Outcome != quest.OutFail {
		t.Fatalf("Outcome = %s, want FAIL (coverage wins over smell)", v.Outcome)
	}
	if v.RootCause != "branch-coverage-below-100" {
		t.Fatalf("RootCause = %q, want branch-coverage-below-100", v.RootCause)
	}
}

func TestRules_FullCoverageWithJavaSmell_Review(t *testing.T) {
	m := &measurement{
		FuncName: "p.Calculator.add",
		Report:   &coverage.Report{AllCovered: true, TotalPct: 100},
		Smells: []smell.Finding{
			{Rule: "TS-REFL-JV-002", File: "CalculatorTest.java", Line: 13, Note: "setAccessible(true)"},
		},
	}
	v := evalVerdict(m)
	if v.Outcome != quest.OutReview {
		t.Fatalf("Outcome = %s, want REVIEW", v.Outcome)
	}
	if v.RootCause != "TS-REFL-JV-002" {
		t.Fatalf("RootCause = %q, want TS-REFL-JV-002", v.RootCause)
	}
}

func TestRules_Catalog_TenRules(t *testing.T) {
	// 2 Fail (G-001, G-002/G-004) + 3 Go-reflect Review + 3 TS-reflect Review
	// (005a) + 2 Java-reflect Review (005c) = 10.
	rules := New().Rules()
	if len(rules) != 10 {
		t.Fatalf("len(Rules) = %d, want 10", len(rules))
	}
	var fail, review int
	for _, r := range rules {
		switch r.Meta.Level {
		case gate.LevelFail:
			fail++
		case gate.LevelReview:
			review++
		}
	}
	if fail != 2 || review != 8 {
		t.Fatalf("levels: fail=%d review=%d, want fail=2 review=8", fail, review)
	}
}
