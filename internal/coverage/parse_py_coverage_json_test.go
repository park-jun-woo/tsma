package coverage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestParsePyCoverageJSONValid(t *testing.T) {
	dir := t.TempDir()
	covData := pyCoverageJSON{
		Files: map[string]pyCoverageFile{
			"handler.py": {
				ExecutedLines: []int{1, 2, 3},
				MissingLines:  []int{4},
			},
		},
	}

	raw, err := json.Marshal(covData)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "coverage.json")
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := parsePyCoverageJSON(path)
	if err != nil {
		t.Fatalf("parsePyCoverageJSON: %v", err)
	}
	if len(result.Files) != 1 {
		t.Errorf("got %d files, want 1", len(result.Files))
	}
}

func TestParsePyCoverageJSONMissingFile(t *testing.T) {
	_, err := parsePyCoverageJSON("/nonexistent/coverage.json")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestParsePyCoverageJSONBadJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "coverage.json")
	if err := os.WriteFile(path, []byte("{bad json}"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := parsePyCoverageJSON(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
