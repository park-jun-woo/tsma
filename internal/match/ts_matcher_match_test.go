package match

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTSMatcherMatchTestTS(t *testing.T) {
	dir := t.TempDir()

	srcDir := filepath.Join(dir, "src")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "handler.ts"), []byte("export function handler() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "handler.test.ts"), []byte("describe('handler', () => {});\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	m := &TSMatcher{}
	testFile, found := m.Match(dir, "src/handler.ts")
	if !found {
		t.Fatal("expected to find .test.ts file")
	}
	want := filepath.Join("src", "handler.test.ts")
	if testFile != want {
		t.Errorf("testFile = %q, want %q", testFile, want)
	}
}

func TestTSMatcherMatchSpecTS(t *testing.T) {
	dir := t.TempDir()

	srcDir := filepath.Join(dir, "src")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "service.ts"), []byte("export function createUser() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "service.spec.ts"), []byte("describe('createUser', () => {});\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	m := &TSMatcher{}
	testFile, found := m.Match(dir, "src/service.ts")
	if !found {
		t.Fatal("expected to find .spec.ts file")
	}
	want := filepath.Join("src", "service.spec.ts")
	if testFile != want {
		t.Errorf("testFile = %q, want %q", testFile, want)
	}
}

func TestTSMatcherMatchInTestsDir(t *testing.T) {
	dir := t.TempDir()

	srcDir := filepath.Join(dir, "src")
	testsDir := filepath.Join(dir, "src", "__tests__")
	if err := os.MkdirAll(testsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "handler.ts"), []byte("export function handler() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(testsDir, "handler.test.ts"), []byte("describe('handler', () => {});\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	m := &TSMatcher{}
	testFile, found := m.Match(dir, "src/handler.ts")
	if !found {
		t.Fatal("expected to find test file in __tests__/ dir")
	}
	want := filepath.Join("src", "__tests__", "handler.test.ts")
	if testFile != want {
		t.Errorf("testFile = %q, want %q", testFile, want)
	}
}

func TestTSMatcherMatchNotFound(t *testing.T) {
	dir := t.TempDir()

	srcDir := filepath.Join(dir, "src")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "handler.ts"), []byte("export function handler() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	m := &TSMatcher{}
	_, found := m.Match(dir, "src/handler.ts")
	if found {
		t.Error("expected no match when no test files exist")
	}
}

func TestTSMatcherMatchJSFile(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "handler.js"), []byte("function handler() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "handler.test.js"), []byte("describe('handler', () => {});\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	m := &TSMatcher{}
	testFile, found := m.Match(dir, "handler.js")
	if !found {
		t.Fatal("expected to find .test.js file")
	}
	if testFile != "handler.test.js" {
		t.Errorf("testFile = %q, want %q", testFile, "handler.test.js")
	}
}

func TestTSMatcherMatchTSXFile(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "component.tsx"), []byte("export function Component() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "component.test.ts"), []byte("describe('Component', () => {});\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	m := &TSMatcher{}
	testFile, found := m.Match(dir, "component.tsx")
	if !found {
		t.Fatal("expected to find test file for .tsx source")
	}
	if testFile != "component.test.ts" {
		t.Errorf("testFile = %q, want %q", testFile, "component.test.ts")
	}
}
