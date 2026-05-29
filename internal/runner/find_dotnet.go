//ff:func feature=runner type=helper control=sequence lang=csharp
//ff:what Locates the dotnet binary on PATH, returning an error if absent
package runner

import (
	"fmt"
	"os/exec"
)

// findDotnet returns the path to the dotnet binary, or an error if the .NET SDK
// is not installed (so the runner can report a clear toolchain-missing message).
func findDotnet() (string, error) {
	path, err := exec.LookPath("dotnet")
	if err != nil {
		return "", fmt.Errorf("dotnet not found on PATH: install the .NET SDK")
	}
	return path, nil
}
