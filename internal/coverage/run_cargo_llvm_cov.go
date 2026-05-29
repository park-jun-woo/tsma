//ff:func feature=coverage type=implementation control=sequence
//ff:what Runs cargo llvm-cov and writes the JSON coverage export to the given path
package coverage

import (
	"fmt"
	"os/exec"
)

// runCargoLLVMCov runs `cargo llvm-cov --json` in projectRoot and writes the
// JSON export to coverJSONPath. A working cargo + cargo-llvm-cov toolchain is
// required (E2E only; not exercised in sandbox environments).
func runCargoLLVMCov(cargo, projectRoot, coverJSONPath string) error {
	cmd := exec.Command(cargo, "llvm-cov", "--json", "--output-path", coverJSONPath)
	cmd.Dir = projectRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cargo llvm-cov: %s\n%s", err, string(output))
	}
	return nil
}
