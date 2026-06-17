//ff:func feature=gate type=test
//ff:what 게이트 룰 헬퍼 단위테스트: asMeasurement(측정/비측정 submission)·uncoveredLocations(빈/정상/cap초과)·firstFinding(매칭/무매칭)·loc·firstLines(절단/무절단)의 모든 분기를 이름으로 직접 호출해 덮는다.

package tsmagate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/tsma/internal/coverage"
	"github.com/park-jun-woo/tsma/internal/smell"
)

func TestAsMeasurement_OK(t *testing.T) {
	want := &measurement{FuncName: "pkg.Fn"}
	got, ok := asMeasurement(gate.Context{Submission: want})
	if !ok {
		t.Fatal("ok = false, want true for a *measurement submission")
	}
	if got != want {
		t.Errorf("got %p, want %p", got, want)
	}
}

func TestAsMeasurement_NotMeasurement(t *testing.T) {
	// A non-measurement submission must report ok=false (defensive path).
	if _, ok := asMeasurement(gate.Context{Submission: "not a measurement"}); ok {
		t.Fatal("ok = true, want false for a non-measurement submission")
	}
}

func TestUncoveredLocations_NotMeasurement(t *testing.T) {
	if got := uncoveredLocations(gate.Context{Submission: 42}); got != "" {
		t.Errorf("got %q, want empty for a non-measurement submission", got)
	}
}

func TestUncoveredLocations_NilReport(t *testing.T) {
	if got := uncoveredLocations(gate.Context{Submission: &measurement{}}); got != "" {
		t.Errorf("got %q, want empty when Report is nil", got)
	}
}

func TestUncoveredLocations_JoinsLocations(t *testing.T) {
	m := &measurement{Report: &coverage.Report{
		Uncovered: []coverage.UncoveredBranch{
			{File: "a.go", Line: 3},
			{File: "b.go", Line: 7},
		},
	}}
	got := uncoveredLocations(gate.Context{Submission: m})
	if got != "a.go:3, b.go:7" {
		t.Errorf("got %q, want \"a.go:3, b.go:7\"", got)
	}
}

func TestUncoveredLocations_CapsAtTen(t *testing.T) {
	// 11 uncovered branches: only the first 10 are listed, then an ellipsis.
	var ubs []coverage.UncoveredBranch
	for i := 1; i <= 11; i++ {
		ubs = append(ubs, coverage.UncoveredBranch{File: "x.go", Line: i})
	}
	m := &measurement{Report: &coverage.Report{Uncovered: ubs}}
	got := uncoveredLocations(gate.Context{Submission: m})
	parts := strings.Split(got, ", ")
	if len(parts) != 11 {
		t.Fatalf("parts = %d, want 11 (10 locs + ellipsis): %q", len(parts), got)
	}
	if parts[10] != "…" {
		t.Errorf("last part = %q, want ellipsis", parts[10])
	}
	if parts[0] != "x.go:1" {
		t.Errorf("first part = %q, want x.go:1", parts[0])
	}
}

func TestFirstFinding_Found(t *testing.T) {
	m := &measurement{Smells: []smell.Finding{
		{Rule: "TS-REFL-001", File: "a_test.go", Line: 2},
		{Rule: "TS-REFL-002", File: "b_test.go", Line: 5},
	}}
	f := firstFinding(m, "TS-REFL-002")
	if f == nil {
		t.Fatal("firstFinding = nil, want the TS-REFL-002 finding")
	}
	if f.File != "b_test.go" || f.Line != 5 {
		t.Errorf("got %+v, want b_test.go:5", f)
	}
}

func TestFirstFinding_NotFound(t *testing.T) {
	m := &measurement{Smells: []smell.Finding{
		{Rule: "TS-REFL-001", File: "a_test.go", Line: 2},
	}}
	if f := firstFinding(m, "TS-REFL-003"); f != nil {
		t.Errorf("firstFinding = %+v, want nil when no rule matches", f)
	}
}

func TestLoc(t *testing.T) {
	got := loc(&smell.Finding{File: "x_test.go", Line: 9})
	if got != "x_test.go:9" {
		t.Errorf("got %q, want x_test.go:9", got)
	}
}

func TestFirstLines_Truncates(t *testing.T) {
	s := "l1\nl2\nl3\nl4"
	got := firstLines(s, 2)
	if got != "l1\nl2\n…" {
		t.Errorf("got %q, want truncation with ellipsis", got)
	}
}

func TestFirstLines_NoTruncation(t *testing.T) {
	s := "l1\nl2"
	got := firstLines(s, 5)
	if got != s {
		t.Errorf("got %q, want unchanged %q", got, s)
	}
}
