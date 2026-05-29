package runner

import (
	"reflect"
	"testing"
)

// TestAppendUniqueFuncs_dedupAcrossCalls verifies names are appended in
// first-seen order and duplicates (within src and against the seen set) are
// skipped.
func TestAppendUniqueFuncs_dedupAcrossCalls(t *testing.T) {
	seen := make(map[string]struct{})
	var dst []string

	dst = appendUniqueFuncs(seen, dst, []string{"TestA", "TestB", "TestA"})
	if want := []string{"TestA", "TestB"}; !reflect.DeepEqual(dst, want) {
		t.Fatalf("after first call dst = %v, want %v", dst, want)
	}

	dst = appendUniqueFuncs(seen, dst, []string{"TestB", "TestC"})
	if want := []string{"TestA", "TestB", "TestC"}; !reflect.DeepEqual(dst, want) {
		t.Fatalf("after second call dst = %v, want %v", dst, want)
	}

	if _, ok := seen["TestC"]; !ok {
		t.Error("seen set should record TestC")
	}
}

// TestAppendUniqueFuncs_emptySrc verifies an empty src leaves dst unchanged.
func TestAppendUniqueFuncs_emptySrc(t *testing.T) {
	seen := map[string]struct{}{"X": {}}
	dst := []string{"X"}

	got := appendUniqueFuncs(seen, dst, nil)
	if want := []string{"X"}; !reflect.DeepEqual(got, want) {
		t.Errorf("dst = %v, want %v", got, want)
	}
}

// TestAppendUniqueFuncs_allDuplicates verifies that when every src name is
// already seen, dst is returned unchanged.
func TestAppendUniqueFuncs_allDuplicates(t *testing.T) {
	seen := map[string]struct{}{"A": {}, "B": {}}
	dst := []string{"A", "B"}

	got := appendUniqueFuncs(seen, dst, []string{"A", "B"})
	if want := []string{"A", "B"}; !reflect.DeepEqual(got, want) {
		t.Errorf("dst = %v, want %v", got, want)
	}
}
