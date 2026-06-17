//ff:func feature=gate type=test
//ff:what Phase003 C2/C3 overlay 측정 테스트: (1) 비침습 — 신규 함수가 디스크 테스트 없이 overlay로 측정되고(MatchFunc 우회) 비종결 시 패키지 디렉터리가 측정 전후 불변, (2) 신규함수 PASS materialize — 종결 통과 시 정명 경로로만 확정(zzz 가상파일 잔존 0), (3) 오염 방지 — 깨진 생성테스트가 형제 함수 측정을 오염 안 함, (4) C3 — 잘린 소스는 truncated Fact + 파일 미기록.
package tsmagate

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/park-jun-woo/tsma/internal/model"
)

// pkgDirEntries returns the sorted file names directly under root/pkg, the
// snapshot used to prove an overlay measurement never touches the source tree.
func pkgDirEntries(t *testing.T, root string) []string {
	t.Helper()
	ents, err := os.ReadDir(filepath.Join(root, "pkg"))
	if err != nil {
		t.Fatalf("readdir pkg: %v", err)
	}
	var names []string
	for _, e := range ents {
		names = append(names, e.Name())
	}
	sort.Strings(names)
	return names
}

// loopSession returns a session with MetaLoop set, driving Prepare's Go overlay
// branch.
func loopSession() *quest.Session {
	s := quest.New()
	s.SetMeta(quest.MetaLoop, true)
	return s
}

// TestPrepare_OverlayMeasuresNewFuncWithoutTouchingTree proves a brand-new
// function with no test on disk is still measured (the disk re-match is bypassed),
// and that a non-terminal partial result leaves the package directory byte-for-byte
// unchanged (overlay only — nothing is written to the source tree).
func TestPrepare_OverlayMeasuresNewFuncWithoutTouchingTree(t *testing.T) {
	root := writeGoPkg(t, map[string]string{
		"go.mod":          "module ovl\n\ngo 1.22\n",
		"pkg/classify.go": classifySrc,
	})
	chdirTo(t, root)
	before := pkgDirEntries(t, root)

	fn := model.Function{QualifiedName: "pkg.Classify", Name: "Classify", File: filepath.Join("pkg", "classify.go"), StartLine: 3, EndLine: 8}
	it := itemWithPayload(t, "go", root, fn) // Tries == 0 (not the final try)
	ctx, _, err := New().Prepare(loopSession(), it, []byte(classifyPartialTest))
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	m, _ := asMeasurement(ctx)
	if m.TestFailed {
		t.Fatalf("overlay measurement failed unexpectedly: %s", m.FailOutput)
	}
	if m.Report == nil {
		t.Fatal("a brand-new function must still be measured via overlay (disk re-match bypassed)")
	}
	if m.Report.AllCovered {
		t.Fatal("partial test must measure < 100% coverage")
	}
	if after := pkgDirEntries(t, root); strings.Join(after, ",") != strings.Join(before, ",") {
		t.Fatalf("source tree changed during overlay measurement: before=%v after=%v", before, after)
	}
}

// TestPrepare_OverlayNewFuncPassMaterializes proves a full covering test for a
// brand-new function passes via overlay and is materialized to the canonical path
// only — no virtual zzz_ file is ever left behind in the package directory.
func TestPrepare_OverlayNewFuncPassMaterializes(t *testing.T) {
	root := writeGoPkg(t, map[string]string{
		"go.mod":          "module ovl\n\ngo 1.22\n",
		"pkg/classify.go": classifySrc,
	})
	chdirTo(t, root)
	fn := model.Function{QualifiedName: "pkg.Classify", Name: "Classify", File: filepath.Join("pkg", "classify.go"), StartLine: 3, EndLine: 8}
	ctx, _, err := New().Prepare(loopSession(), itemWithPayload(t, "go", root, fn), []byte(classifyFullTest))
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	m, _ := asMeasurement(ctx)
	if m.TestFailed || m.Report == nil || !m.Report.AllCovered {
		t.Fatalf("expected a fully-covered overlay PASS, got TestFailed=%v report=%+v", m.TestFailed, m.Report)
	}
	if _, err := os.Stat(filepath.Join(root, "pkg", "classify_test.go")); err != nil {
		t.Fatalf("a PASS must be materialized to the canonical test path: %v", err)
	}
	for _, name := range pkgDirEntries(t, root) {
		if strings.Contains(name, "zzz_tsma_gen") {
			t.Fatalf("the virtual overlay file must never reach the source tree, found %q", name)
		}
	}
}

// TestPrepare_OverlayBrokenSiblingDoesNotContaminate proves a broken generated
// test for one function in a multi-function package never reaches disk, so it
// cannot poison the measurement of a sibling function in the same package.
func TestPrepare_OverlayBrokenSiblingDoesNotContaminate(t *testing.T) {
	root := writeGoPkg(t, map[string]string{
		"go.mod":       "module ovl\n\ngo 1.22\n",
		"pkg/alpha.go": "package pkg\n\nfunc Alpha() int { return 1 }\n",
		"pkg/beta.go":  "package pkg\n\nfunc Beta() int { return 2 }\n",
	})
	chdirTo(t, root)

	// Beta gets a parseable but non-compiling test (references an undefined symbol):
	// it fails to build, but is overlay-only so it never lands in pkg/.
	betaFn := model.Function{QualifiedName: "pkg.Beta", Name: "Beta", File: filepath.Join("pkg", "beta.go"), StartLine: 3, EndLine: 3}
	broken := "package pkg\n\nimport \"testing\"\n\nfunc TestBeta(t *testing.T) { _ = Nonexistent() }\n"
	bctx, _, err := New().Prepare(loopSession(), itemWithPayload(t, "go", root, betaFn), []byte(broken))
	if err != nil {
		t.Fatalf("Prepare(Beta): %v", err)
	}
	bm, _ := asMeasurement(bctx)
	if !bm.TestFailed {
		t.Fatal("a non-compiling generated test must be TestFailed")
	}
	if _, statErr := os.Stat(filepath.Join(root, "pkg", "beta_test.go")); statErr == nil {
		t.Fatal("the broken test must not have been written to the source tree")
	}

	// Alpha now measures cleanly — Beta's broken test never contaminated the package.
	alphaFn := model.Function{QualifiedName: "pkg.Alpha", Name: "Alpha", File: filepath.Join("pkg", "alpha.go"), StartLine: 3, EndLine: 3}
	good := "package pkg\n\nimport \"testing\"\n\nfunc TestAlpha(t *testing.T) { if Alpha() != 1 { t.Fatal(\"x\") } }\n"
	actx, _, err := New().Prepare(loopSession(), itemWithPayload(t, "go", root, alphaFn), []byte(good))
	if err != nil {
		t.Fatalf("Prepare(Alpha): %v", err)
	}
	am, _ := asMeasurement(actx)
	if am.TestFailed {
		t.Fatalf("sibling Alpha measurement was contaminated: %s", am.FailOutput)
	}
	if am.Report == nil || !am.Report.AllCovered {
		t.Fatalf("expected Alpha to pass at full coverage, got %+v", am.Report)
	}
}

// TestPrepare_OverlayTruncatedFeedback proves a truncated (unparseable) generated
// source skips measurement entirely: it is flagged truncated, writes nothing to
// the source tree, and tests-must-pass emits the dedicated "emit the ENTIRE file"
// Fact rather than a raw compiler error.
func TestPrepare_OverlayTruncatedFeedback(t *testing.T) {
	root := writeGoPkg(t, map[string]string{
		"go.mod":          "module ovl\n\ngo 1.22\n",
		"pkg/classify.go": classifySrc,
	})
	chdirTo(t, root)
	before := pkgDirEntries(t, root)
	fn := model.Function{QualifiedName: "pkg.Classify", Name: "Classify", File: filepath.Join("pkg", "classify.go"), StartLine: 3, EndLine: 8}

	truncated := "package pkg\n\nimport \"testing\"\n\nfunc TestClassify(t *testing.T) {\n\tif Classify(1) != \"pos\" {" // cut off: no closing braces
	ctx, _, err := New().Prepare(loopSession(), itemWithPayload(t, "go", root, fn), []byte(truncated))
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	m, _ := asMeasurement(ctx)
	if !m.TestFailed || !m.Truncated {
		t.Fatalf("truncated source must set TestFailed+Truncated, got TestFailed=%v Truncated=%v", m.TestFailed, m.Truncated)
	}
	if after := pkgDirEntries(t, root); strings.Join(after, ",") != strings.Join(before, ",") {
		t.Fatalf("truncated source must write nothing to the tree: before=%v after=%v", before, after)
	}
	fired, fact := testsMustPass.Check(gate.Context{Submission: m})
	if !fired {
		t.Fatal("tests-must-pass must fire on a truncated measurement")
	}
	if fact.Expected != "complete compilable file" || !strings.Contains(fact.Actual, "truncated") {
		t.Fatalf("expected the dedicated truncated Fact, got %+v", fact)
	}
}

// tsmaTestResidue returns the names directly under root/.tsma/test (gen backings
// and overlay JSON) — the scratch C2 must leave empty after a measurement.
func tsmaTestResidue(t *testing.T, root string) []string {
	t.Helper()
	ents, err := os.ReadDir(filepath.Join(root, ".tsma", "test"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		t.Fatalf("readdir .tsma/test: %v", err)
	}
	var names []string
	for _, e := range ents {
		names = append(names, e.Name())
	}
	sort.Strings(names)
	return names
}

// TestPrepare_OverlayVacuousZeroCoverNotMaterialized (C1-a) proves a generated
// test that compiles, runs, and exits 0 but never calls the target (0% coverage)
// is downgraded to TestFailed with the dedicated "exercises"/"covers 0%" Fact, and
// is NOT materialized even on the final try (the *통과DONE* path requires cover>0).
func TestPrepare_OverlayVacuousZeroCoverNotMaterialized(t *testing.T) {
	root := writeGoPkg(t, map[string]string{
		"go.mod":          "module ovl\n\ngo 1.22\n",
		"pkg/classify.go": classifySrc,
	})
	chdirTo(t, root)
	fn := model.Function{QualifiedName: "pkg.Classify", Name: "Classify", File: filepath.Join("pkg", "classify.go"), StartLine: 3, EndLine: 8}
	it := itemWithPayload(t, "go", root, fn)
	it.Tries = quest.MaxTries - 1 // final try: would lock *통과DONE* if accepted

	noCall := "package pkg\n\nimport \"testing\"\n\nfunc TestClassify(t *testing.T) {\n\tif 1+1 != 2 {\n\t\tt.Fatal(\"x\")\n\t}\n}\n"
	ctx, _, err := New().Prepare(loopSession(), it, []byte(noCall))
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	m, _ := asMeasurement(ctx)
	if !m.TestFailed {
		t.Fatal("a zero-coverage pass must be downgraded to TestFailed")
	}
	if m.Report != nil {
		t.Fatalf("a downgraded vacuous pass must carry no Report, got %+v", m.Report)
	}
	fired, fact := testsMustPass.Check(gate.Context{Submission: m})
	if !fired {
		t.Fatal("tests-must-pass must fire on a vacuous (0%) pass")
	}
	if !strings.Contains(fact.Expected, "exercises") || !strings.Contains(fact.Actual, "covers 0%") {
		t.Fatalf("expected the dedicated 0%% Fact, got %+v", fact)
	}
	if _, statErr := os.Stat(filepath.Join(root, "pkg", "classify_test.go")); !os.IsNotExist(statErr) {
		t.Fatalf("a vacuous pass must not materialize even on the final try, stat err = %v", statErr)
	}
}

// TestPrepare_OverlayMalformedTestNameRejected (C1-b) proves a malformed test name
// (`TestpyIndent` — lowercase after "Test", which go test silently ignores) is
// rejected before measurement with a clear name-specific Fact.
func TestPrepare_OverlayMalformedTestNameRejected(t *testing.T) {
	root := writeGoPkg(t, map[string]string{
		"go.mod":          "module ovl\n\ngo 1.22\n",
		"pkg/classify.go": classifySrc,
	})
	chdirTo(t, root)
	fn := model.Function{QualifiedName: "pkg.Classify", Name: "Classify", File: filepath.Join("pkg", "classify.go"), StartLine: 3, EndLine: 8}

	malformed := "package pkg\n\nimport \"testing\"\n\nfunc TestpyIndent(t *testing.T) {\n\t_ = Classify(1)\n}\n"
	ctx, _, err := New().Prepare(loopSession(), itemWithPayload(t, "go", root, fn), []byte(malformed))
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	m, _ := asMeasurement(ctx)
	if !m.TestFailed || m.Report != nil {
		t.Fatalf("a malformed test name must be rejected pre-measure (TestFailed, no Report), got %+v", m)
	}
	fired, fact := testsMustPass.Check(gate.Context{Submission: m})
	if !fired || !strings.Contains(fact.Actual, "TestpyIndent") || !strings.Contains(fact.Actual, "malformed") {
		t.Fatalf("expected a name-specific malformed Fact, got fired=%v %+v", fired, fact)
	}
	// Pre-measure reject writes no backing, so .tsma/test stays empty.
	if r := tsmaTestResidue(t, root); len(r) != 0 {
		t.Fatalf("a pre-measure reject must leave no .tsma/test residue, got %v", r)
	}
}

// TestPrepare_OverlayNoResidueAfterMeasure (C2) proves every overlay measurement
// path leaves .tsma/test empty: the materialize path (full cover, final try) and
// the retry path (partial cover, mid-loop) both sweep the backing and overlay JSON.
func TestPrepare_OverlayNoResidueAfterMeasure(t *testing.T) {
	cases := []struct {
		name  string
		tries int
		raw   string
	}{
		{"materialize", quest.MaxTries - 1, classifyFullTest},
		{"retry", 0, classifyPartialTest},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			root := writeGoPkg(t, map[string]string{
				"go.mod":          "module ovl\n\ngo 1.22\n",
				"pkg/classify.go": classifySrc,
			})
			chdirTo(t, root)
			fn := model.Function{QualifiedName: "pkg.Classify", Name: "Classify", File: filepath.Join("pkg", "classify.go"), StartLine: 3, EndLine: 8}
			it := itemWithPayload(t, "go", root, fn)
			it.Tries = c.tries
			if _, _, err := New().Prepare(loopSession(), it, []byte(c.raw)); err != nil {
				t.Fatalf("Prepare: %v", err)
			}
			if r := tsmaTestResidue(t, root); len(r) != 0 {
				t.Fatalf("measurement must leave no .tsma/test residue, got %v", r)
			}
		})
	}
}
