package match

import (
	"os"
	"path/filepath"
	"testing"
)

func writeCsTest(t *testing.T, root, rel, content string) {
	t.Helper()
	abs := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(abs, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestCsMatcherTestProjectTests(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join("App", "Services", "Foo.cs")
	test := filepath.Join("App.Tests", "Services", "FooTests.cs")
	writeCsTest(t, dir, src, "namespace App;\npublic class Foo {}\n")
	writeCsTest(t, dir, test, "namespace App.Tests;\npublic class FooTests {}\n")

	m := &CsMatcher{}
	got, found := m.Match(dir, src)
	if !found {
		t.Fatal("expected match for FooTests.cs in App.Tests")
	}
	if got != test {
		t.Errorf("testFile = %q, want %q", got, test)
	}
}

func TestCsMatcherSameDirTest(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join("src", "Bar.cs")
	test := filepath.Join("src", "BarTest.cs")
	writeCsTest(t, dir, src, "public class Bar {}\n")
	writeCsTest(t, dir, test, "public class BarTest {}\n")

	m := &CsMatcher{}
	got, found := m.Match(dir, src)
	if !found {
		t.Fatal("expected match for BarTest.cs in same dir")
	}
	if got != test {
		t.Errorf("testFile = %q, want %q", got, test)
	}
}

func TestCsMatcherPrefersTestsSuffix(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join("App", "Baz.cs")
	writeCsTest(t, dir, src, "public class Baz {}\n")
	// Both suffixes exist in the same dir; "Tests" should win.
	writeCsTest(t, dir, filepath.Join("App", "BazTests.cs"), "public class BazTests {}\n")
	writeCsTest(t, dir, filepath.Join("App", "BazTest.cs"), "public class BazTest {}\n")

	m := &CsMatcher{}
	got, found := m.Match(dir, src)
	if !found {
		t.Fatal("expected a match")
	}
	if got != filepath.Join("App", "BazTests.cs") {
		t.Errorf("testFile = %q, want App/BazTests.cs", got)
	}
}

func TestCsMatcherNoTest(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join("App", "Lonely.cs")
	writeCsTest(t, dir, src, "public class Lonely {}\n")

	m := &CsMatcher{}
	if _, found := m.Match(dir, src); found {
		t.Error("expected no match when no test file exists")
	}
}

func TestCsMatcherImplementsMatcher(t *testing.T) {
	var _ Matcher = &CsMatcher{}
}
