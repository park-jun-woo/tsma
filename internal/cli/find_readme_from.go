//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Walks up from a directory to locate README.md, falling back to the GitHub URL
package cli

import (
	"os"
	"path/filepath"
)

// readmeGitHubURL points to the README on the project's GitHub repository.
// It is used when no local README.md can be found.
const readmeGitHubURL = "https://github.com/park-jun-woo/tsma/blob/main/README.md"

// findReadmeFrom searches for README.md starting at startDir and walking up to
// the filesystem root. It returns the absolute path of the first README.md
// found, or the GitHub URL when none exists along the way.
func findReadmeFrom(startDir string) string {
	dir := startDir
	for {
		candidate := filepath.Join(dir, "README.md")
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return readmeGitHubURL
}
