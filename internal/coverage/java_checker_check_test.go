package coverage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// fakeMaven installs a fake `mvn` executable on PATH whose script body is the
// supplied shell snippet, and returns a project dir containing a pom.xml.
func fakeMaven(t *testing.T, script string) string {
	t.Helper()
	binDir := t.TempDir()
	bin := filepath.Join(binDir, "mvn")
	if err := os.WriteFile(bin, []byte("#!/bin/sh\n"+script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	proj := t.TempDir()
	if err := os.WriteFile(filepath.Join(proj, "pom.xml"), []byte("<project/>"), 0o644); err != nil {
		t.Fatal(err)
	}
	return proj
}

// TestJavaCheckerCheckToolNotFound covers the findJavaTool error branch
// (line 21): a Maven project (pom.xml) with no mvnw wrapper and no `mvn` on
// PATH, so the build-tool binary cannot be located.
func TestJavaCheckerCheckToolNotFound(t *testing.T) {
	proj := t.TempDir()
	if err := os.WriteFile(filepath.Join(proj, "pom.xml"), []byte("<project/>"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Point PATH at an empty directory so exec.LookPath("mvn") fails.
	t.Setenv("PATH", t.TempDir())

	c := &JavaChecker{}
	_, err := c.Check(proj, mkMatch("src/test/java/com/example/CalculatorTest.java"),
		&model.Function{File: "src/main/java/com/example/Calculator.java", StartLine: 5, EndLine: 5})
	if err == nil {
		t.Fatal("expected error when mvn is not found on PATH")
	}
	if !strings.Contains(err.Error(), "not found on PATH") {
		t.Errorf("expected 'not found on PATH', got: %v", err)
	}
}

func TestJavaCheckerCheckNoBuildTool(t *testing.T) {
	c := &JavaChecker{}
	_, err := c.Check(t.TempDir(), mkMatch("src/test/java/p/FooTest.java"),
		&model.Function{File: "src/main/java/p/Foo.java", StartLine: 1, EndLine: 5})
	if err == nil {
		t.Fatal("expected error when no build tool marker present")
	}
	if !strings.Contains(err.Error(), "build tool") {
		t.Errorf("error should mention build tool: %v", err)
	}
}

func TestJavaCheckerCheckRunFails(t *testing.T) {
	proj := fakeMaven(t, "exit 1\n")
	c := &JavaChecker{}
	_, err := c.Check(proj, mkMatch("src/test/java/com/example/CalculatorTest.java"),
		&model.Function{File: "src/main/java/com/example/Calculator.java", StartLine: 5, EndLine: 5})
	if err == nil {
		t.Fatal("expected error when build tool run fails")
	}
	if !strings.Contains(err.Error(), "java coverage failed") {
		t.Errorf("expected 'java coverage failed', got: %v", err)
	}
}

func TestJavaCheckerCheckParseFails(t *testing.T) {
	// mvn succeeds but writes no report; parseJacoco then fails to read it.
	proj := fakeMaven(t, "exit 0\n")
	c := &JavaChecker{}
	_, err := c.Check(proj, mkMatch("src/test/java/com/example/CalculatorTest.java"),
		&model.Function{File: "src/main/java/com/example/Calculator.java", StartLine: 5, EndLine: 5})
	if err == nil {
		t.Fatal("expected parse error when report is missing")
	}
	if !strings.Contains(err.Error(), "parse jacoco xml") {
		t.Errorf("expected 'parse jacoco xml', got: %v", err)
	}
}

func TestJavaCheckerCheckSuccess(t *testing.T) {
	fixtureAbs, err := filepath.Abs(jacocoFixture)
	if err != nil {
		t.Fatal(err)
	}
	// mvn writes the fixture to the conventional Maven report path.
	proj := fakeMaven(t, `
mkdir -p target/site/jacoco
cat "`+fixtureAbs+`" > target/site/jacoco/jacoco.xml
exit 0
`)

	c := &JavaChecker{}
	report, err := c.Check(proj, mkMatch("src/test/java/com/example/CalculatorTest.java"),
		&model.Function{File: "src/main/java/com/example/Calculator.java", Name: "classify", StartLine: 9, EndLine: 12})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report == nil {
		t.Fatal("expected non-nil report")
	}
	if report.AllCovered {
		t.Error("expected classify() to be partially covered")
	}
}

func TestJavaCheckerImplementsChecker(t *testing.T) {
	var _ Checker = &JavaChecker{}
}
