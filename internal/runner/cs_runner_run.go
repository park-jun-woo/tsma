//ff:func feature=runner type=implementation control=sequence lang=csharp
//ff:what Executes C# tests for the given test file using dotnet test with a class filter
package runner

import (
	"os/exec"

	"github.com/park-jun-woo/tsma/internal/match"
)

// E2E note: running this path requires a working .NET SDK (`dotnet`), which is
// not available in CI/sandbox environments. The argument builder is split out
// (buildCsTestArgs) for environment-independent unit testing.

// Run executes the given C# test file against the project using `dotnet test`,
// filtering to the single test class derived from the file name.
func (r *CsRunner) Run(projectRoot string, m match.TestMatch) (*Result, error) {
	testFile := testFileFromMatch(m)
	dotnet, err := findDotnet()
	if err != nil {
		return nil, err
	}

	args := buildCsTestArgs(csTestClass(testFile))
	cmd := exec.Command(dotnet, args...)
	cmd.Dir = projectRoot

	output, err := cmd.CombinedOutput()
	result := &Result{
		Output: string(output),
		Pass:   err == nil,
	}
	return result, nil
}
