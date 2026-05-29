//ff:func feature=detect type=helper control=iteration dimension=1 lang=csharp
//ff:what Returns true if the project root carries a C#/.NET marker file
package detect

import (
	"os"
	"path/filepath"
)

// detectCSharp reports whether projectRoot contains a C#/.NET marker:
//   - any *.csproj or *.sln file (matched via glob, since their names vary), or
//   - a Directory.Build.props file (fixed name).
//
// Glob matching is used for the project/solution files because, unlike the
// fixed marker names of other languages, their base names are arbitrary.
func detectCSharp(projectRoot string) bool {
	for _, pattern := range []string{"*.csproj", "*.sln"} {
		matches, err := filepath.Glob(filepath.Join(projectRoot, pattern))
		if err == nil && len(matches) > 0 {
			return true
		}
	}
	if _, err := os.Stat(filepath.Join(projectRoot, "Directory.Build.props")); err == nil {
		return true
	}
	return false
}
