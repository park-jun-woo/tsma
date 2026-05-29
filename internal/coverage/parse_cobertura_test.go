package coverage

import (
	"os"
	"path/filepath"
	"testing"
)

const coberturaFixture = "../../testdata/csharp/coverage/cobertura.xml"

func TestParseCobertura(t *testing.T) {
	cov, err := parseCobertura(coberturaFixture)
	if err != nil {
		t.Fatalf("parseCobertura: %v", err)
	}
	if len(cov.Files) != 1 {
		t.Fatalf("files = %d, want 1", len(cov.Files))
	}
	f := cov.Files[0]
	if f.Path != "App/Calculator.cs" {
		t.Errorf("path = %q, want App/Calculator.cs", f.Path)
	}
	if len(f.Lines) != 4 {
		t.Fatalf("lines = %d, want 4", len(f.Lines))
	}

	// Spot-check the branch line (number=9): branch=true, 50% (1/2), hits=2.
	var l9 *coberturaLine
	for i := range f.Lines {
		if f.Lines[i].Number == 9 {
			l9 = &f.Lines[i]
		}
	}
	if l9 == nil {
		t.Fatal("missing line number=9")
	}
	if l9.Branch != "true" || l9.Hits != 2 || l9.ConditionCoverage != "50% (1/2)" {
		t.Errorf("line 9 = %+v, want branch=true hits=2 cc=50%% (1/2)", *l9)
	}
}

func TestParseCoberturaMissingFile(t *testing.T) {
	if _, err := parseCobertura("/nonexistent/cobertura.xml"); err == nil {
		t.Error("expected error for missing file")
	}
}

func TestParseCoberturaInvalidXML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.xml")
	if err := os.WriteFile(path, []byte("<coverage><not-closed"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := parseCobertura(path); err == nil {
		t.Error("expected decode error for invalid XML")
	}
}
