//ff:func feature=runner type=implementation control=sequence
//ff:what Executes a Python test file using pytest or unittest
package runner

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/tsma/internal/match"
)

// pyFallbackHint is prepended to a failing unittest-fallback run's output so the
// failure is not mistaken for "tsma doesn't support Python". It is added only on
// failure (D3); successful fallback runs pass quietly.
const pyFallbackHint = "[tsma] pytest 미감지로 unittest 폴백 실행; src-layout/pytest-style 프로젝트면 import 실패가 정상적 결과 — pyproject에 pytest 선언/pytest.ini 확인\n"

// Run executes the given Python test file against the project.
func (r *PyRunner) Run(projectRoot string, m match.TestMatch) (*Result, error) {
	testFile := testFileFromMatch(m)
	absTest := filepath.Join(projectRoot, testFile)

	usePytest := detectPytest(projectRoot)
	python := findPython()

	var cmd *exec.Cmd
	if usePytest {
		cmd = exec.Command(python, "-m", "pytest", absTest, "-v")
	} else {
		rel, _ := filepath.Rel(projectRoot, absTest)
		modulePath := strings.TrimSuffix(rel, ".py")
		modulePath = strings.ReplaceAll(modulePath, string(filepath.Separator), ".")
		cmd = exec.Command(python, "-m", "unittest", modulePath, "-v")
	}
	cmd.Dir = projectRoot

	output, err := cmd.CombinedOutput()
	result := &Result{
		Output: string(output),
		Pass:   err == nil,
	}
	// D3 — diagnose a failed unittest fallback. The runner can't reliably tell
	// an import-stage failure from a genuine assertion failure, so we attach the
	// hint whenever the fallback path failed, and stay quiet on success.
	if !usePytest && !result.Pass {
		result.Output = pyFallbackHint + result.Output
	}
	return result, nil
}
