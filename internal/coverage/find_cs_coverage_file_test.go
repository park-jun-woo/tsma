package coverage

import "testing"

func TestFindCsCoverageFile(t *testing.T) {
	cov := &csCoverage{Files: []csFile{
		{Path: "Calculator.cs"},
		{Path: "App/Other.cs"},
	}}

	// Suffix match: the cobertura class filename is a suffix of the longer
	// project-relative path.
	f := findCsCoverageFile(cov, "src/App/Calculator.cs", "/proj")
	if f == nil || f.Path != "Calculator.cs" {
		t.Fatalf("got %+v, want Calculator via suffix", f)
	}

	// Exact relative path match.
	f = findCsCoverageFile(cov, "App/Other.cs", "/proj")
	if f == nil || f.Path != "App/Other.cs" {
		t.Fatalf("got %+v, want Other", f)
	}

	if findCsCoverageFile(cov, "App/Missing.cs", "/proj") != nil {
		t.Error("expected nil for unmatched file")
	}
}
