package runner

import "testing"

func TestBuildTestArgsVitestFramework(t *testing.T) {
	args := buildTestArgs(frameworkVitest, "src/handler.test.ts")
	expected := []string{"vitest", "run", "src/handler.test.ts", "--reporter=verbose"}
	assertArgsEqual(t, args, expected)
}

func TestBuildTestArgsJestFramework(t *testing.T) {
	args := buildTestArgs(frameworkJest, "src/handler.test.ts")
	expected := []string{"jest", "src/handler.test.ts", "--verbose"}
	assertArgsEqual(t, args, expected)
}

func TestBuildTestArgsMochaFramework(t *testing.T) {
	args := buildTestArgs(frameworkMocha, "src/handler.test.ts")
	expected := []string{"mocha", "src/handler.test.ts"}
	assertArgsEqual(t, args, expected)
}

func TestBuildTestArgsUnknownFramework(t *testing.T) {
	args := buildTestArgs("unknown", "test.ts")
	expected := []string{"vitest", "run", "test.ts", "--reporter=verbose"}
	assertArgsEqual(t, args, expected)
}

func assertArgsEqual(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %v (len %d), want %v (len %d)", got, len(got), want, len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("arg[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
