package runner

import "testing"

func TestAnchorRunPattern_empty(t *testing.T) {
	if got := AnchorRunPattern(nil); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

func TestAnchorRunPattern_single(t *testing.T) {
	if got := AnchorRunPattern([]string{"TestFoo"}); got != "^TestFoo$" {
		t.Errorf("got %q, want ^TestFoo$", got)
	}
}

func TestAnchorRunPattern_multiple(t *testing.T) {
	got := AnchorRunPattern([]string{"TestA", "TestB", "TestC"})
	want := "^(TestA|TestB|TestC)$"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
