package match

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// TestAttributeFunc_contentAwareWins verifies attributeFunc keeps the precise
// content-aware match when a test references the function by name, even though a
// conventional test file is also present.
func TestAttributeFunc_contentAwareWins(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/m\n\ngo 1.22\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "widget.go"), []byte("package m\n\nfunc Widget() int { return 1 }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "direct_test.go"),
		[]byte("package m\n\nimport \"testing\"\n\nfunc TestDirect(t *testing.T) { _ = Widget() }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "widget_test.go"),
		[]byte("package m\n\nimport \"testing\"\n\nfunc TestIndirect(t *testing.T) {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	idx, err := BuildPkgTestIndex(root, ".")
	if err != nil {
		t.Fatal(err)
	}
	srcRecv, _ := BuildPkgSourceReceivers(root, ".")
	fn := &model.Function{Name: "Widget", File: "widget.go"}
	tm, ok := attributeFunc(root, idx, srcRecv, fn)
	if !ok {
		t.Fatal("expected Widget to be attributed content-aware")
	}
	if len(tm.Files) != 1 || filepath.Base(tm.Files[0]) != "direct_test.go" {
		t.Fatalf("Files = %v, want [direct_test.go] (content-aware, not fallback)", tm.Files)
	}
	if len(tm.TestFuncs) != 1 || tm.TestFuncs[0] != "TestDirect" {
		t.Fatalf("TestFuncs = %v, want [TestDirect]", tm.TestFuncs)
	}
}

// TestAttributeFunc_fallbackWhenUnreferenced verifies attributeFunc falls back
// to the conventional <base>_test.go when content-aware finds no reference.
func TestAttributeFunc_fallbackWhenUnreferenced(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/m\n\ngo 1.22\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "agent.go"), []byte("package m\n\nfunc agentCmd() int { return 1 }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "agent_test.go"),
		[]byte("package m\n\nimport \"testing\"\n\nfunc TestAgent(t *testing.T) { runCmd(t, \"agent\") }\n\nfunc runCmd(t *testing.T, name string) {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	idx, err := BuildPkgTestIndex(root, ".")
	if err != nil {
		t.Fatal(err)
	}
	srcRecv, _ := BuildPkgSourceReceivers(root, ".")
	fn := &model.Function{Name: "agentCmd", File: "agent.go"}
	tm, ok := attributeFunc(root, idx, srcRecv, fn)
	if !ok {
		t.Fatal("expected agentCmd to fall back to agent_test.go")
	}
	if len(tm.Files) != 1 || filepath.Base(tm.Files[0]) != "agent_test.go" {
		t.Fatalf("Files = %v, want [agent_test.go]", tm.Files)
	}
	if tm.TestFuncs != nil {
		t.Fatalf("TestFuncs = %v, want nil (runner resolves)", tm.TestFuncs)
	}
}

// TestAttributeFunc_noMatch verifies attributeFunc reports found false when
// neither content-aware nor the file-name fallback yields a match.
func TestAttributeFunc_noMatch(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/m\n\ngo 1.22\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "lonely.go"), []byte("package m\n\nfunc Lonely() int { return 1 }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// scenario_test.go is not the conventional name for lonely.go and does not
	// reference Lonely, so neither path matches.
	if err := os.WriteFile(filepath.Join(root, "scenario_test.go"),
		[]byte("package m\n\nimport \"testing\"\n\nfunc TestScenario(t *testing.T) {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	idx, err := BuildPkgTestIndex(root, ".")
	if err != nil {
		t.Fatal(err)
	}
	srcRecv, _ := BuildPkgSourceReceivers(root, ".")
	fn := &model.Function{Name: "Lonely", File: "lonely.go"}
	if tm, ok := attributeFunc(root, idx, srcRecv, fn); ok {
		t.Fatalf("expected Lonely to be unmatched, got %v", tm)
	}
}

// TestAttributeFunc_nilIndexFallback verifies that with a nil index (package
// could not be indexed) attributeFunc still uses the file-name fallback.
func TestAttributeFunc_nilIndexFallback(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "thing.go"), []byte("package m\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "thing_test.go"), []byte("package m\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	fn := &model.Function{Name: "Thing", File: "thing.go"}
	tm, ok := attributeFunc(root, nil, nil, fn)
	if !ok {
		t.Fatal("expected nil-index path to use file-name fallback")
	}
	if len(tm.Files) != 1 || filepath.Base(tm.Files[0]) != "thing_test.go" {
		t.Fatalf("Files = %v, want [thing_test.go]", tm.Files)
	}
}
