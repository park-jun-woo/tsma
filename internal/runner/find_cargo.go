//ff:func feature=runner type=helper control=sequence
//ff:what Locates the cargo binary on PATH, returning an error if absent
package runner

import (
	"fmt"
	"os/exec"
)

// findCargo returns the path to the cargo binary, or an error if cargo is not
// installed (so the runner can report a clear toolchain-missing message).
func findCargo() (string, error) {
	path, err := exec.LookPath("cargo")
	if err != nil {
		return "", fmt.Errorf("cargo not found on PATH: install the Rust toolchain")
	}
	return path, nil
}
