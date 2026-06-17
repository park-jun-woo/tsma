package smell

import "testing"

// TestScanJavaNoTreeSitter covers the tree-sitter-unavailable branch: ScanJava
// returns (nil, err) so the caller ignores the file.
func TestScanJavaNoTreeSitter(t *testing.T) {
	t.Setenv("TSMA_TREE_SITTER", "/nonexistent/abs/tree-sitter")
	findings, err := ScanJava("../../testdata/java/smell/ReflectionTest.java")
	if err == nil {
		t.Error("expected error when tree-sitter unavailable")
	}
	if findings != nil {
		t.Errorf("expected nil findings, got %+v", findings)
	}
}

// TestScanJavaParseError covers the ParseFile-failure branch with a real CLI
// present but a nonexistent source path.
func TestScanJavaParseError(t *testing.T) {
	if !locateJavaSmell(t) {
		t.Skip("tree-sitter CLI + java grammar not available")
	}
	findings, err := ScanJava("../../testdata/java/smell/does_not_exist.java")
	if err == nil {
		t.Error("expected parse error for a missing file")
	}
	if findings != nil {
		t.Errorf("expected nil findings, got %+v", findings)
	}
}
