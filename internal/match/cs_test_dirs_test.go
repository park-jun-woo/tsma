package match

import "testing"

func TestCsTestDirsNested(t *testing.T) {
	got := csTestDirs("App/Services")
	want := []string{"App/Services", "App.Tests/Services"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("dir[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestCsTestDirsTopLevelProject(t *testing.T) {
	got := csTestDirs("App")
	want := []string{"App", "App.Tests"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("dir[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestCsTestDirsRoot(t *testing.T) {
	got := csTestDirs(".")
	if len(got) != 1 || got[0] != "." {
		t.Errorf("got %v, want [.]", got)
	}
}
