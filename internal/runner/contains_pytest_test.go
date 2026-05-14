package runner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestContainsPytestMatch(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "requirements.txt")
	if err := os.WriteFile(path, []byte("pytest==7.0.0\nflask==2.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if !containsPytest(path, "pytest") {
		t.Error("expected true for file containing 'pytest'")
	}
}

func TestContainsPytestNoMatch(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "requirements.txt")
	if err := os.WriteFile(path, []byte("flask==2.0\nrequests==2.28\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if containsPytest(path, "pytest") {
		t.Error("expected false for file not containing 'pytest'")
	}
}

func TestContainsPytestCaseInsensitive(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.txt")
	if err := os.WriteFile(path, []byte("PYTEST==7.0.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if !containsPytest(path, "pytest") {
		t.Error("expected true for case-insensitive match")
	}
}

func TestContainsPytestNonexistentFile(t *testing.T) {
	if containsPytest("/nonexistent/path/file.txt", "pytest") {
		t.Error("expected false for nonexistent file")
	}
}

func TestContainsPytestEmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.txt")
	if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}

	if containsPytest(path, "pytest") {
		t.Error("expected false for empty file")
	}
}
