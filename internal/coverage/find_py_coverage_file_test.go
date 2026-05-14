package coverage

import "testing"

func TestFindPyCoverageFileMatch(t *testing.T) {
	covData := &pyCoverageJSON{
		Files: map[string]pyCoverageFile{
			"handler.py": {ExecutedLines: []int{1, 2, 3}},
		},
	}

	result := findPyCoverageFile(covData, "handler.py", "/project")
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if len(result.ExecutedLines) != 3 {
		t.Errorf("ExecutedLines = %d, want 3", len(result.ExecutedLines))
	}
}

func TestFindPyCoverageFileNoMatch(t *testing.T) {
	covData := &pyCoverageJSON{
		Files: map[string]pyCoverageFile{
			"other.py": {ExecutedLines: []int{1}},
		},
	}

	result := findPyCoverageFile(covData, "handler.py", "/project")
	if result != nil {
		t.Error("expected nil for non-matching file")
	}
}

func TestFindPyCoverageFileEmpty(t *testing.T) {
	covData := &pyCoverageJSON{
		Files: map[string]pyCoverageFile{},
	}

	result := findPyCoverageFile(covData, "handler.py", "/project")
	if result != nil {
		t.Error("expected nil for empty files map")
	}
}

func TestFindPyCoverageFileSuffixMatch(t *testing.T) {
	covData := &pyCoverageJSON{
		Files: map[string]pyCoverageFile{
			"/some/path/handler.py": {ExecutedLines: []int{1}},
		},
	}

	result := findPyCoverageFile(covData, "handler.py", "/project")
	if result == nil {
		t.Fatal("expected non-nil result for suffix match")
	}
}

func TestFindPyCoverageFileReturnsCopy(t *testing.T) {
	covData := &pyCoverageJSON{
		Files: map[string]pyCoverageFile{
			"handler.py": {ExecutedLines: []int{1, 2, 3}},
		},
	}

	result := findPyCoverageFile(covData, "handler.py", "/project")
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	result.ExecutedLines = append(result.ExecutedLines, 99)
	orig := covData.Files["handler.py"]
	if len(orig.ExecutedLines) != 3 {
		t.Error("modifying returned value should not affect original")
	}
}
