//ff:func feature=coverage type=implementation control=sequence lang=csharp
//ff:what Runs dotnet test with the XPlat coverage collector to produce a Cobertura report
package coverage

import (
	"fmt"
	"os/exec"
)

// runCsCoverage runs `dotnet test` with the given coverage-collection arguments
// in projectRoot, producing a Cobertura XML report under the results directory.
// A working .NET SDK + coverlet.collector is required (E2E only; not exercised
// in sandbox environments).
func runCsCoverage(bin, projectRoot string, args []string) error {
	cmd := exec.Command(bin, args...)
	cmd.Dir = projectRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("dotnet coverage run: %s\n%s", err, string(output))
	}
	return nil
}
