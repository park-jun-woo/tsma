package coverage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestParseCoverageFinalJSONValid(t *testing.T) {
	dir := t.TempDir()
	data := map[string]coverageFinalEntry{
		"src/handler.ts": {
			StatementMap: map[string]coverageRange{
				"0": {Start: coveragePosition{Line: 1}, End: coveragePosition{Line: 1}},
			},
			S:         map[string]int{"0": 1},
			BranchMap: map[string]coverageBranch{},
			B:         map[string][]int{},
			FnMap:     map[string]coverageFunction{},
			F:         map[string]int{},
		},
	}

	raw, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "coverage-final.json"), raw, 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := parseCoverageFinalJSON(dir)
	if err != nil {
		t.Fatalf("parseCoverageFinalJSON: %v", err)
	}
	if _, ok := result["src/handler.ts"]; !ok {
		t.Error("expected key 'src/handler.ts'")
	}
}

func TestParseCoverageFinalJSONMissingFile(t *testing.T) {
	dir := t.TempDir()
	_, err := parseCoverageFinalJSON(dir)
	if err == nil {
		t.Fatal("expected error when file does not exist")
	}
}

func TestParseCoverageFinalJSONInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "coverage-final.json"), []byte("bad"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := parseCoverageFinalJSON(dir)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
