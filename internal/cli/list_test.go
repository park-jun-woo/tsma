package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestRunList_getProjectRootError(t *testing.T) {
	// Force getProjectRoot() to fail by removing the current working directory
	// out from under the process, so os.Getwd() returns an error. This covers
	// the early `return err` branch in runList.
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer os.Chdir(orig)

	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	// Remove the cwd; subsequent os.Getwd() should fail.
	if err := os.Remove(dir); err != nil {
		t.Skipf("could not remove cwd to simulate getwd failure: %v", err)
	}

	if _, gErr := os.Getwd(); gErr == nil {
		t.Skip("os.Getwd did not fail after removing cwd on this platform")
	}

	listPage = 1
	if err := runList(nil, nil); err == nil {
		t.Fatal("expected error when getProjectRoot fails")
	}
}

func TestRunList_noSession(t *testing.T) {
	dir := t.TempDir()

	// Change to temp dir so getProjectRoot returns it
	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	err := runList(nil, nil)
	if err == nil {
		t.Fatal("expected error when no session exists")
	}
}

func TestRunList_withSession(t *testing.T) {
	dir := t.TempDir()

	// Create a valid session file
	sess := &model.Session{
		Project: dir,
		Lang:    "go",
		Functions: []model.Function{
			{Name: "A", Status: model.StatusPass, CoveragePct: 100},
			{Name: "B", Status: model.StatusTodo},
		},
		Summary: model.Summary{Total: 2, Pass: 1, Todo: 1},
	}
	sessDir := filepath.Join(dir, ".tsma")
	os.MkdirAll(sessDir, 0o755)
	os.MkdirAll(filepath.Join(sessDir, "tests"), 0o755)
	data, _ := json.MarshalIndent(sess, "", "  ")
	os.WriteFile(filepath.Join(sessDir, "session.json"), data, 0o644)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	// Reset page to valid value
	listPage = 1

	output := captureStdout(func() {
		err := runList(nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if output == "" {
		t.Error("expected non-empty output")
	}
}

func TestRunList_corruptSession(t *testing.T) {
	dir := t.TempDir()

	sessDir := filepath.Join(dir, ".tsma")
	os.MkdirAll(sessDir, 0o755)
	// Write an invalid JSON so Exists passes but Load fails.
	os.WriteFile(filepath.Join(sessDir, "session.json"), []byte("{not json"), 0o644)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	listPage = 1
	err := runList(nil, nil)
	if err == nil {
		t.Fatal("expected error loading corrupt session")
	}
}

func TestRunList_pageBelowOne(t *testing.T) {
	dir := t.TempDir()

	sess := &model.Session{
		Project: dir,
		Lang:    "go",
		Functions: []model.Function{
			{Name: "A", Status: model.StatusTodo},
		},
		Summary: model.Summary{Total: 1, Todo: 1},
	}
	sessDir := filepath.Join(dir, ".tsma")
	os.MkdirAll(sessDir, 0o755)
	os.MkdirAll(filepath.Join(sessDir, "tests"), 0o755)
	data, _ := json.MarshalIndent(sess, "", "  ")
	os.WriteFile(filepath.Join(sessDir, "session.json"), data, 0o644)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	// page < 1 should be normalized to 1 (covers the listPage < 1 branch).
	listPage = 0
	output := captureStdout(func() {
		if err := runList(nil, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	if output == "" {
		t.Error("expected non-empty output")
	}
	if listPage != 1 {
		t.Errorf("expected listPage normalized to 1, got %d", listPage)
	}
}

func TestRunList_multiplePages(t *testing.T) {
	dir := t.TempDir()

	// More than one page (pageSize is 20) so the totalPages > 1 footer prints.
	fns := make([]model.Function, 0, 25)
	for i := 0; i < 25; i++ {
		fns = append(fns, model.Function{Name: "fn", Status: model.StatusTodo})
	}
	sess := &model.Session{
		Project:   dir,
		Lang:      "go",
		Functions: fns,
		Summary:   model.Summary{Total: 25, Todo: 25},
	}
	sessDir := filepath.Join(dir, ".tsma")
	os.MkdirAll(sessDir, 0o755)
	os.MkdirAll(filepath.Join(sessDir, "tests"), 0o755)
	data, _ := json.MarshalIndent(sess, "", "  ")
	os.WriteFile(filepath.Join(sessDir, "session.json"), data, 0o644)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	listPage = 1
	output := captureStdout(func() {
		if err := runList(nil, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(output, "Page 1/2") {
		t.Errorf("expected pagination footer 'Page 1/2', got: %q", output)
	}
}

func TestRunList_pageOutOfRange(t *testing.T) {
	dir := t.TempDir()

	sess := &model.Session{
		Project: dir,
		Lang:    "go",
		Functions: []model.Function{
			{Name: "A", Status: model.StatusTodo},
		},
		Summary: model.Summary{Total: 1, Todo: 1},
	}
	sessDir := filepath.Join(dir, ".tsma")
	os.MkdirAll(sessDir, 0o755)
	os.MkdirAll(filepath.Join(sessDir, "tests"), 0o755)
	data, _ := json.MarshalIndent(sess, "", "  ")
	os.WriteFile(filepath.Join(sessDir, "session.json"), data, 0o644)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	listPage = 999

	output := captureStdout(func() {
		err := runList(nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if output == "" {
		t.Error("expected non-empty output")
	}
}
