package cli

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeWithMtime(t *testing.T, root, rel string, mt time.Time) {
	t.Helper()
	abs := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(abs, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(abs, mt, mt); err != nil {
		t.Fatal(err)
	}
}

// TestCombinedTestMtime_returnsLatest verifies the max mtime across the set is
// returned, formatted as RFC3339.
func TestCombinedTestMtime_returnsLatest(t *testing.T) {
	root := t.TempDir()
	older := time.Date(2022, 5, 1, 0, 0, 0, 0, time.UTC)
	newer := time.Date(2024, 6, 2, 12, 0, 0, 0, time.UTC)
	writeWithMtime(t, root, "a_test.go", older)
	writeWithMtime(t, root, "b_test.go", newer)

	got := combinedTestMtime(root, []string{"a_test.go", "b_test.go"})
	parsed, err := time.Parse(time.RFC3339, got)
	if err != nil {
		t.Fatalf("combinedTestMtime returned unparseable %q: %v", got, err)
	}
	if !parsed.Equal(newer) {
		t.Errorf("combinedTestMtime = %v, want %v", parsed, newer)
	}
}

// TestCombinedTestMtime_missingFilesIgnored verifies missing files contribute
// nothing and the max is taken from existing files only.
func TestCombinedTestMtime_missingFilesIgnored(t *testing.T) {
	root := t.TempDir()
	mt := time.Date(2023, 3, 3, 3, 3, 3, 0, time.UTC)
	writeWithMtime(t, root, "present_test.go", mt)

	got := combinedTestMtime(root, []string{"missing_test.go", "present_test.go"})
	parsed, err := time.Parse(time.RFC3339, got)
	if err != nil {
		t.Fatalf("combinedTestMtime returned unparseable %q: %v", got, err)
	}
	if !parsed.Equal(mt) {
		t.Errorf("combinedTestMtime = %v, want %v", parsed, mt)
	}
}

// TestCombinedTestMtime_emptySet verifies an empty set yields "".
func TestCombinedTestMtime_emptySet(t *testing.T) {
	if got := combinedTestMtime(t.TempDir(), nil); got != "" {
		t.Errorf("combinedTestMtime = %q, want empty", got)
	}
}

// TestCombinedTestMtime_allMissing verifies an all-missing set yields "".
func TestCombinedTestMtime_allMissing(t *testing.T) {
	got := combinedTestMtime(t.TempDir(), []string{"nope_test.go", "gone_test.go"})
	if got != "" {
		t.Errorf("combinedTestMtime = %q, want empty", got)
	}
}
