package coverage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// fakeDotnet installs a fake `dotnet` executable on PATH whose script body is
// the supplied shell snippet, and returns a fresh project dir.
func fakeDotnet(t *testing.T, script string) string {
	t.Helper()
	binDir := t.TempDir()
	bin := filepath.Join(binDir, "dotnet")
	if err := os.WriteFile(bin, []byte("#!/bin/sh\n"+script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	return t.TempDir()
}

// TestCsCheckerCheckToolNotFound covers the findDotnet error branch: dotnet is
// absent from PATH.
func TestCsCheckerCheckToolNotFound(t *testing.T) {
	t.Setenv("PATH", t.TempDir())

	c := &CsChecker{}
	_, err := c.Check(t.TempDir(), mkMatch("App.Tests/CalculatorTests.cs"),
		&model.Function{File: "App/Calculator.cs", StartLine: 5, EndLine: 5})
	if err == nil {
		t.Fatal("expected error when dotnet is not found on PATH")
	}
	if !strings.Contains(err.Error(), "dotnet not found") {
		t.Errorf("expected 'dotnet not found', got: %v", err)
	}
}

func TestCsCheckerCheckRunFails(t *testing.T) {
	proj := fakeDotnet(t, "exit 1\n")
	c := &CsChecker{}
	_, err := c.Check(proj, mkMatch("App.Tests/CalculatorTests.cs"),
		&model.Function{File: "App/Calculator.cs", StartLine: 5, EndLine: 5})
	if err == nil {
		t.Fatal("expected error when dotnet run fails")
	}
	if !strings.Contains(err.Error(), "csharp coverage failed") {
		t.Errorf("expected 'csharp coverage failed', got: %v", err)
	}
}

func TestCsCheckerCheckReportMissing(t *testing.T) {
	// dotnet succeeds but writes no report; locating it then fails.
	proj := fakeDotnet(t, "exit 0\n")
	c := &CsChecker{}
	_, err := c.Check(proj, mkMatch("App.Tests/CalculatorTests.cs"),
		&model.Function{File: "App/Calculator.cs", StartLine: 5, EndLine: 5})
	if err == nil {
		t.Fatal("expected error when cobertura report is missing")
	}
	if !strings.Contains(err.Error(), "locate cobertura xml") {
		t.Errorf("expected 'locate cobertura xml', got: %v", err)
	}
}

func TestCsCheckerCheckSuccess(t *testing.T) {
	fixtureAbs, err := filepath.Abs(coberturaFixture)
	if err != nil {
		t.Fatal(err)
	}
	// dotnet writes the fixture under the conventional coverlet output layout.
	proj := fakeDotnet(t, `
mkdir -p .tsma/coverage/guid-1
cat "`+fixtureAbs+`" > .tsma/coverage/guid-1/coverage.cobertura.xml
exit 0
`)

	c := &CsChecker{}
	report, err := c.Check(proj, mkMatch("App.Tests/CalculatorTests.cs"),
		&model.Function{File: "App/Calculator.cs", Name: "Classify", StartLine: 9, EndLine: 12})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report == nil {
		t.Fatal("expected non-nil report")
	}
	if report.AllCovered {
		t.Error("expected Classify() to be partially covered")
	}
}

// TestCsCheckerCheckMkdirFails covers the coverage-dir preparation error branch:
// dotnet is on PATH, but the project's .tsma path already exists as a regular
// file, so clearing/creating .tsma/coverage underneath it fails.
func TestCsCheckerCheckMkdirFails(t *testing.T) {
	binDir := t.TempDir()
	bin := filepath.Join(binDir, "dotnet")
	if err := os.WriteFile(bin, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	proj := t.TempDir()
	// .tsma exists as a FILE, so MkdirAll(.tsma/coverage) must fail.
	if err := os.WriteFile(filepath.Join(proj, ".tsma"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	c := &CsChecker{}
	_, err := c.Check(proj, mkMatch("App.Tests/CalculatorTests.cs"),
		&model.Function{File: "App/Calculator.cs", StartLine: 5, EndLine: 5})
	if err == nil {
		t.Fatal("expected error when .tsma coverage dir cannot be created")
	}
	if !strings.Contains(err.Error(), ".tsma coverage dir") {
		t.Errorf("expected '.tsma coverage dir' error, got: %v", err)
	}
}

// TestCsCheckerCheckParseFails covers the parseCobertura error branch: dotnet
// writes a report file that is found but contains malformed XML.
func TestCsCheckerCheckParseFails(t *testing.T) {
	proj := fakeDotnet(t, `
mkdir -p .tsma/coverage/guid-1
printf 'this is not valid xml <<<' > .tsma/coverage/guid-1/coverage.cobertura.xml
exit 0
`)

	c := &CsChecker{}
	_, err := c.Check(proj, mkMatch("App.Tests/CalculatorTests.cs"),
		&model.Function{File: "App/Calculator.cs", StartLine: 5, EndLine: 5})
	if err == nil {
		t.Fatal("expected error when cobertura xml is malformed")
	}
	if !strings.Contains(err.Error(), "parse cobertura xml") {
		t.Errorf("expected 'parse cobertura xml', got: %v", err)
	}
}

func TestCsCheckerImplementsChecker(t *testing.T) {
	var _ Checker = &CsChecker{}
}
