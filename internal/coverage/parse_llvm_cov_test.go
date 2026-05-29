package coverage

import (
	"os"
	"path/filepath"
	"testing"
)

const llvmCovFixture = "../../testdata/rust/coverage/llvm-cov.json"

func TestParseLLVMCov(t *testing.T) {
	cov, err := parseLLVMCov(llvmCovFixture)
	if err != nil {
		t.Fatalf("parseLLVMCov: %v", err)
	}
	if len(cov.Data) != 1 {
		t.Fatalf("data blocks = %d, want 1", len(cov.Data))
	}
	files := cov.Data[0].Files
	if len(files) != 1 {
		t.Fatalf("files = %d, want 1", len(files))
	}
	f := files[0]
	if filepath.Base(f.Filename) != "lib.rs" {
		t.Errorf("filename = %q, want .../lib.rs", f.Filename)
	}
	if len(f.Segments) != 9 {
		t.Errorf("segments = %d, want 9", len(f.Segments))
	}
	if len(f.Branches) != 2 {
		t.Errorf("branches = %d, want 2", len(f.Branches))
	}

	// Spot-check a decoded segment: line 3, count 2, region entry, has count.
	var seg3 *llvmSegment
	for i := range f.Segments {
		if f.Segments[i].Line == 3 && f.Segments[i].IsRegionEntry {
			seg3 = &f.Segments[i]
			break
		}
	}
	if seg3 == nil {
		t.Fatal("missing region-entry segment at line 3")
	}
	if seg3.Count != 2 || !seg3.HasCount {
		t.Errorf("segment line 3: count=%d hasCount=%v, want 2/true", seg3.Count, seg3.HasCount)
	}

	// Spot-check a branch: line 9 with false side covered, true side not.
	var br9 *llvmBranch
	for i := range f.Branches {
		if f.Branches[i].LineStart == 9 {
			br9 = &f.Branches[i]
			break
		}
	}
	if br9 == nil {
		t.Fatal("missing branch at line 9")
	}
	if br9.ExecCount != 0 || br9.FalseExecCount != 1 {
		t.Errorf("branch line 9: exec=%d false=%d, want 0/1", br9.ExecCount, br9.FalseExecCount)
	}
}

func TestParseLLVMCovMissingFile(t *testing.T) {
	if _, err := parseLLVMCov("/nonexistent/llvm-cov.json"); err == nil {
		t.Error("expected error for missing file")
	}
}

func TestParseLLVMCovInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("{not valid json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := parseLLVMCov(path); err == nil {
		t.Error("expected unmarshal error for invalid JSON")
	}
}
