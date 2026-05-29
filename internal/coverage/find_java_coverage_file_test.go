package coverage

import "testing"

func TestFindJavaCoverageFile(t *testing.T) {
	cov := &jacocoCoverage{Files: []jacocoFile{
		{Path: "com/example/Calculator.java"},
		{Path: "com/example/Other.java"},
	}}

	// Project-relative file ends with the package-qualified jacoco path.
	f := findJavaCoverageFile(cov, "src/main/java/com/example/Calculator.java", "/proj")
	if f == nil || f.Path != "com/example/Calculator.java" {
		t.Fatalf("got %+v, want Calculator", f)
	}

	if findJavaCoverageFile(cov, "src/main/java/com/example/Missing.java", "/proj") != nil {
		t.Error("expected nil for unmatched file")
	}
}
