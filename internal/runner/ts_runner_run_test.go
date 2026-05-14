package runner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestTSRunnerDetectsFrameworkAndBuildsArgs(t *testing.T) {
	dir := t.TempDir()

	// Write package.json with jest
	pkg := struct {
		DevDependencies map[string]string `json:"devDependencies"`
	}{DevDependencies: map[string]string{"jest": "^29.0.0"}}
	data, err := json.Marshal(pkg)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "package.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}

	framework := detectTSTestFramework(dir)
	if framework != frameworkJest {
		t.Fatalf("framework = %q, want %q", framework, frameworkJest)
	}

	args := buildTestArgs(framework, "src/handler.test.ts")
	if len(args) != 3 {
		t.Fatalf("args = %v, want len 3", args)
	}
	if args[0] != "jest" {
		t.Errorf("args[0] = %q, want \"jest\"", args[0])
	}
	if args[1] != "src/handler.test.ts" {
		t.Errorf("args[1] = %q, want \"src/handler.test.ts\"", args[1])
	}
	if args[2] != "--verbose" {
		t.Errorf("args[2] = %q, want \"--verbose\"", args[2])
	}
}

func TestTSRunnerFallbackFramework(t *testing.T) {
	dir := t.TempDir()
	// No package.json -> vitest fallback
	framework := detectTSTestFramework(dir)
	if framework != frameworkVitest {
		t.Errorf("framework = %q, want %q (fallback)", framework, frameworkVitest)
	}

	args := buildTestArgs(framework, "test.ts")
	if args[0] != "vitest" {
		t.Errorf("args[0] = %q, want \"vitest\"", args[0])
	}
}
