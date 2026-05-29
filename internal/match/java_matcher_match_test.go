package match

import (
	"os"
	"path/filepath"
	"testing"
)

func writeJavaTest(t *testing.T, root, rel, content string) {
	t.Helper()
	abs := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(abs, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestJavaMatcherStandardLayoutTest(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join("src", "main", "java", "p", "Foo.java")
	test := filepath.Join("src", "test", "java", "p", "FooTest.java")
	writeJavaTest(t, dir, src, "package p;\npublic class Foo {}\n")
	writeJavaTest(t, dir, test, "package p;\npublic class FooTest {}\n")

	m := &JavaMatcher{}
	got, found := m.Match(dir, src)
	if !found {
		t.Fatal("expected match for FooTest.java")
	}
	if got != test {
		t.Errorf("testFile = %q, want %q", got, test)
	}
}

func TestJavaMatcherTestsSuffix(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join("src", "main", "java", "p", "Bar.java")
	test := filepath.Join("src", "test", "java", "p", "BarTests.java")
	writeJavaTest(t, dir, src, "package p;\npublic class Bar {}\n")
	writeJavaTest(t, dir, test, "package p;\npublic class BarTests {}\n")

	m := &JavaMatcher{}
	got, found := m.Match(dir, src)
	if !found {
		t.Fatal("expected match for BarTests.java")
	}
	if got != test {
		t.Errorf("testFile = %q, want %q", got, test)
	}
}

func TestJavaMatcherNoTest(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join("src", "main", "java", "p", "Lonely.java")
	writeJavaTest(t, dir, src, "package p;\npublic class Lonely {}\n")

	m := &JavaMatcher{}
	if _, found := m.Match(dir, src); found {
		t.Error("expected no match when no test file exists")
	}
}

func TestJavaMatcherImplementsMatcher(t *testing.T) {
	var _ Matcher = &JavaMatcher{}
}
