package coverage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindCoberturaReport(t *testing.T) {
	dir := t.TempDir()
	guidDir := filepath.Join(dir, "abc-123-guid")
	if err := os.MkdirAll(guidDir, 0o755); err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(guidDir, "coverage.cobertura.xml")
	if err := os.WriteFile(want, []byte("<coverage/>"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := findCoberturaReport(dir)
	if err != nil {
		t.Fatalf("findCoberturaReport: %v", err)
	}
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFindCoberturaReportMissing(t *testing.T) {
	dir := t.TempDir()
	if _, err := findCoberturaReport(dir); err == nil {
		t.Error("expected error when no cobertura report exists")
	}
}

// TestFindCoberturaReportSkipsAfterFound covers the `found != ""` SkipDir
// branch: once a report is located, later sibling directories are skipped, so
// the first lexical match wins.
func TestFindCoberturaReportSkipsAfterFound(t *testing.T) {
	dir := t.TempDir()
	// Two guid dirs; "aaa" sorts before "bbb" so its report is found first.
	first := filepath.Join(dir, "aaa-guid")
	second := filepath.Join(dir, "bbb-guid")
	if err := os.MkdirAll(first, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(second, 0o755); err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(first, "coverage.cobertura.xml")
	if err := os.WriteFile(want, []byte("<coverage/>"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(second, "coverage.cobertura.xml"), []byte("<coverage/>"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := findCoberturaReport(dir)
	if err != nil {
		t.Fatalf("findCoberturaReport: %v", err)
	}
	if got != want {
		t.Errorf("got %q, want first match %q", got, want)
	}
}

// TestFindCoberturaReportWalkError covers the walk-error branch (err != nil in
// the WalkFunc): an unreadable subdirectory yields a walk error that is
// swallowed, and the report from a readable sibling is still located.
func TestFindCoberturaReportWalkError(t *testing.T) {
	dir := t.TempDir()
	// "zzz" sorts last so the unreadable dir is visited after the report is
	// already found is not what we want here; instead make the unreadable dir
	// sort first so its error path (err != nil) is exercised before the match.
	bad := filepath.Join(dir, "000-bad")
	good := filepath.Join(dir, "111-good")
	if err := os.MkdirAll(bad, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(good, 0o755); err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(good, "coverage.cobertura.xml")
	if err := os.WriteFile(want, []byte("<coverage/>"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Remove read/exec so walking into bad yields a permission error.
	if err := os.Chmod(bad, 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(bad, 0o755) })

	got, err := findCoberturaReport(dir)
	if err != nil {
		t.Fatalf("findCoberturaReport: %v", err)
	}
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
