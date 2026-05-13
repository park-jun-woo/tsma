package session

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()

	sess := &model.Session{
		Project:   dir,
		Lang:      "go",
		Framework: "gin",
		Created:   time.Now(),
		Endpoints: []model.Endpoint{
			{
				Name:   "Login",
				Method: "POST",
				Path:   "/api/login",
				Status: model.StatusTodo,
			},
			{
				Name:   "Signup",
				Method: "POST",
				Path:   "/api/signup",
				Status: model.StatusDone,
			},
		},
	}

	if err := Save(dir, sess); err != nil {
		t.Fatalf("Save: %v", err)
	}

	if !Exists(dir) {
		t.Fatal("Exists returned false after Save")
	}

	loaded, err := Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if loaded.Lang != "go" {
		t.Errorf("Lang = %q, want go", loaded.Lang)
	}
	if len(loaded.Endpoints) != 2 {
		t.Errorf("Endpoints count = %d, want 2", len(loaded.Endpoints))
	}
	if loaded.Summary.Total != 2 {
		t.Errorf("Summary.Total = %d, want 2", loaded.Summary.Total)
	}
	if loaded.Summary.Done != 1 {
		t.Errorf("Summary.Done = %d, want 1", loaded.Summary.Done)
	}
}

func TestDelete(t *testing.T) {
	dir := t.TempDir()

	sess := &model.Session{
		Project: dir,
		Lang:    "go",
	}
	if err := Save(dir, sess); err != nil {
		t.Fatalf("Save: %v", err)
	}

	if err := Delete(dir); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	if Exists(dir) {
		t.Error("Exists returned true after Delete")
	}
}

func TestCopyTestFile(t *testing.T) {
	dir := t.TempDir()

	// Create .tsma dir.
	os.MkdirAll(filepath.Join(dir, ".tsma", "tests"), 0o755)

	// Create a test file.
	srcPath := filepath.Join(dir, "my_test.go")
	if err := os.WriteFile(srcPath, []byte("package foo"), 0o644); err != nil {
		t.Fatal(err)
	}

	rel, err := CopyTestFile(dir, srcPath)
	if err != nil {
		t.Fatalf("CopyTestFile: %v", err)
	}

	if rel != filepath.Join(".tsma", "tests", "my_test.go") {
		t.Errorf("relative path = %q, unexpected", rel)
	}

	dst := filepath.Join(dir, rel)
	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("read copy: %v", err)
	}
	if string(data) != "package foo" {
		t.Errorf("content = %q, want 'package foo'", string(data))
	}
}
