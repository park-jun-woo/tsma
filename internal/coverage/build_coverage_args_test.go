package coverage

import "testing"

func TestBuildCoverageArgsForVitest(t *testing.T) {
	args := buildCoverageArgs(coverVitest, "src/app.test.ts", "/tmp/cov")
	if len(args) == 0 {
		t.Fatal("expected non-empty args")
	}
	if args[0] != "vitest" {
		t.Errorf("args[0] = %q, want 'vitest'", args[0])
	}
	if args[1] != "run" {
		t.Errorf("args[1] = %q, want 'run'", args[1])
	}
	if args[2] != "src/app.test.ts" {
		t.Errorf("args[2] = %q, want test file", args[2])
	}
}

func TestBuildCoverageArgsForJest(t *testing.T) {
	args := buildCoverageArgs(coverJest, "src/app.test.ts", "/tmp/cov")
	if len(args) == 0 {
		t.Fatal("expected non-empty args")
	}
	if args[0] != "jest" {
		t.Errorf("args[0] = %q, want 'jest'", args[0])
	}
	if args[1] != "src/app.test.ts" {
		t.Errorf("args[1] = %q, want test file", args[1])
	}
}

func TestBuildCoverageArgsDefaultFramework(t *testing.T) {
	args := buildCoverageArgs("unknown", "test.ts", "/out")
	if args[0] != "vitest" {
		t.Errorf("default framework args[0] = %q, want 'vitest'", args[0])
	}
}

func TestBuildCoverageArgsCoverDir(t *testing.T) {
	args := buildCoverageArgs(coverVitest, "test.ts", "/custom/dir")
	found := false
	for _, a := range args {
		if a == "--coverage.reportsDirectory=/custom/dir" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected coverDir in vitest args")
	}

	args = buildCoverageArgs(coverJest, "test.ts", "/custom/dir")
	found = false
	for _, a := range args {
		if a == "--coverageDirectory=/custom/dir" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected coverDir in jest args")
	}
}
