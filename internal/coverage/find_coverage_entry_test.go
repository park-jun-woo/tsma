package coverage

import "testing"

func TestFindCoverageEntryNilMap(t *testing.T) {
	entry := findCoverageEntry("src/handler.ts", nil, "/project")
	if entry != nil {
		t.Error("expected nil for nil map")
	}
}

func TestFindCoverageEntryEmptyMap(t *testing.T) {
	data := map[string]coverageFinalEntry{}
	entry := findCoverageEntry("src/handler.ts", data, "/project")
	if entry != nil {
		t.Error("expected nil for empty map")
	}
}

func TestFindCoverageEntryEmptyProjectRoot(t *testing.T) {
	data := map[string]coverageFinalEntry{
		"src/handler.ts": {S: map[string]int{"0": 1}},
	}
	entry := findCoverageEntry("src/handler.ts", data, "")
	if entry == nil {
		t.Fatal("expected to find entry with direct match even with empty projectRoot")
	}
}

func TestFindCoverageEntryReturnsNonNilForMatch(t *testing.T) {
	data := map[string]coverageFinalEntry{
		"src/handler.ts": {S: map[string]int{"0": 5}},
	}
	entry := findCoverageEntry("src/handler.ts", data, "/project")
	if entry == nil {
		t.Fatal("expected non-nil entry")
	}
	if entry.S["0"] != 5 {
		t.Errorf("S[\"0\"] = %d, want 5", entry.S["0"])
	}
}
