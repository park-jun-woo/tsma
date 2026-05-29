//ff:func feature=coverage type=helper control=sequence lang=csharp
//ff:what Locates the dotnet binary on PATH, returning an error if absent
package coverage

import (
	"fmt"
	"os/exec"
)

// findDotnet returns the path to the dotnet binary, or an error if the .NET SDK
// is not installed (so the checker can report a clear toolchain-missing
// message).
func findDotnet() (string, error) {
	path, err := exec.LookPath("dotnet")
	if err != nil {
		return "", fmt.Errorf("dotnet not found on PATH: install the .NET SDK and coverlet.collector")
	}
	return path, nil
}
