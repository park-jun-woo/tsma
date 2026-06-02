package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// TestRunStatus_getProjectRootError covers the getProjectRoot error branch
// (line 21) by removing the cwd so os.Getwd() fails.
func TestRunStatus_getProjectRootError(t *testing.T) {
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer os.Chdir(orig)

	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	if err := os.Remove(dir); err != nil {
		t.Skipf("could not remove cwd: %v", err)
	}
	if _, gErr := os.Getwd(); gErr == nil {
		t.Skip("os.Getwd did not fail after removing cwd on this platform")
	}

	if err := runStatus(nil, nil); err == nil {
		t.Fatal("expected error when getProjectRoot fails")
	}
}

func TestRunStatus_noSession(t *testing.T) {
	dir := t.TempDir()

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	err := runStatus(nil, nil)
	if err == nil {
		t.Fatal("expected error when no session exists")
	}
}

func TestRunStatus_withSession(t *testing.T) {
	dir := t.TempDir()
	// Real source so the on-load reconcile (Phase012) matches by QualifiedName
	// ("pkg.<Name>") and preserves each function's status instead of dropping it.
	writeGoFunc(t, dir, "A")
	writeGoFunc(t, dir, "B")
	writeGoFunc(t, dir, "C")

	sess := &model.Session{
		Project: dir,
		Lang:    "go",
		Functions: []model.Function{
			{QualifiedName: "pkg.A", Name: "A", File: "pkg/a.go", Status: model.StatusPass, CoveragePct: 100},
			{QualifiedName: "pkg.B", Name: "B", File: "pkg/b.go", Status: model.StatusDone, CoveragePct: 80},
			{QualifiedName: "pkg.C", Name: "C", File: "pkg/c.go", Status: model.StatusTodo},
		},
		Summary: model.Summary{Total: 3, Pass: 1, Done: 1, Todo: 1},
	}
	sessDir := filepath.Join(dir, ".tsma")
	os.MkdirAll(sessDir, 0o755)
	os.MkdirAll(filepath.Join(sessDir, "tests"), 0o755)
	data, _ := json.MarshalIndent(sess, "", "  ")
	os.WriteFile(filepath.Join(sessDir, "session.json"), data, 0o644)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	output := captureStdout(func() {
		err := runStatus(nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(output, "3 functions") {
		t.Errorf("expected '3 functions' in output, got %q", output)
	}
	if !strings.Contains(output, "PASS") {
		t.Error("expected PASS in output")
	}
	if !strings.Contains(output, "DONE") {
		t.Error("expected DONE in output")
	}
	if !strings.Contains(output, "TODO") {
		t.Error("expected TODO in output")
	}
	if !strings.Contains(output, "Coverage average") {
		t.Error("expected coverage average in output")
	}
}

func TestRunStatus_corruptSession(t *testing.T) {
	dir := t.TempDir()

	sessDir := filepath.Join(dir, ".tsma")
	os.MkdirAll(sessDir, 0o755)
	// Exists true, Load fails -> error from the load branch.
	os.WriteFile(filepath.Join(sessDir, "session.json"), []byte("{broken"), 0o644)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	err := runStatus(nil, nil)
	if err == nil {
		t.Fatal("expected error loading corrupt session")
	}
}

func TestRunStatus_emptyFunctions(t *testing.T) {
	dir := t.TempDir()

	sess := &model.Session{
		Project:   dir,
		Lang:      "go",
		Functions: []model.Function{},
		Summary:   model.Summary{Total: 0},
	}
	sessDir := filepath.Join(dir, ".tsma")
	os.MkdirAll(sessDir, 0o755)
	os.MkdirAll(filepath.Join(sessDir, "tests"), 0o755)
	data, _ := json.MarshalIndent(sess, "", "  ")
	os.WriteFile(filepath.Join(sessDir, "session.json"), data, 0o644)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	output := captureStdout(func() {
		err := runStatus(nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(output, "No functions found") {
		t.Errorf("expected 'No functions found' in output, got %q", output)
	}
}
