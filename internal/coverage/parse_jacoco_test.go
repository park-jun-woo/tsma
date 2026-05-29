package coverage

import (
	"os"
	"path/filepath"
	"testing"
)

const jacocoFixture = "../../testdata/java/coverage/jacoco.xml"

func TestParseJacoco(t *testing.T) {
	cov, err := parseJacoco(jacocoFixture)
	if err != nil {
		t.Fatalf("parseJacoco: %v", err)
	}
	if len(cov.Files) != 1 {
		t.Fatalf("files = %d, want 1", len(cov.Files))
	}
	f := cov.Files[0]
	if f.Path != "com/example/Calculator.java" {
		t.Errorf("path = %q, want com/example/Calculator.java", f.Path)
	}
	if len(f.Lines) != 4 {
		t.Fatalf("lines = %d, want 4", len(f.Lines))
	}

	// Spot-check the branch line (nr=9): mb=1, cb=1, ci=2.
	var l9 *jacocoLine
	for i := range f.Lines {
		if f.Lines[i].Nr == 9 {
			l9 = &f.Lines[i]
		}
	}
	if l9 == nil {
		t.Fatal("missing line nr=9")
	}
	if l9.Cb != 1 || l9.Mb != 1 || l9.Ci != 2 {
		t.Errorf("line 9 = %+v, want cb=1 mb=1 ci=2", *l9)
	}
}

func TestParseJacocoMissingFile(t *testing.T) {
	if _, err := parseJacoco("/nonexistent/jacoco.xml"); err == nil {
		t.Error("expected error for missing file")
	}
}

func TestParseJacocoInvalidXML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.xml")
	if err := os.WriteFile(path, []byte("<report><not-closed"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := parseJacoco(path); err == nil {
		t.Error("expected decode error for invalid XML")
	}
}
