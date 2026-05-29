package coverage

import "testing"

func TestFindRsCoverageFile(t *testing.T) {
	cov, err := parseLLVMCov(llvmCovFixture)
	if err != nil {
		t.Fatal(err)
	}

	f := findRsCoverageFile(cov, "src/lib.rs", "/home/dev/demo")
	if f == nil {
		t.Fatal("expected to find src/lib.rs entry")
	}
	if len(f.Segments) != 9 {
		t.Errorf("segments = %d, want 9", len(f.Segments))
	}

	if findRsCoverageFile(cov, "src/other.rs", "/home/dev/demo") != nil {
		t.Error("expected nil for unknown file")
	}
}
