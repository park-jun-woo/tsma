//ff:func feature=coverage type=helper control=sequence
//ff:what Locates the cargo binary on PATH, returning an error if absent
package coverage

import (
	"fmt"
	"os/exec"
)

// findCargo returns the path to the cargo binary, or an error if cargo is not
// installed (so the checker can report a clear toolchain-missing message).
func findCargo() (string, error) {
	path, err := exec.LookPath("cargo")
	if err != nil {
		return "", fmt.Errorf("cargo not found on PATH: install the Rust toolchain and cargo-llvm-cov")
	}
	return path, nil
}
