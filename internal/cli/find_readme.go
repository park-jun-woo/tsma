//ff:func feature=cli type=helper control=sequence
//ff:what Resolves the README location from the working directory, else the GitHub URL
package cli

import "os"

// findReadme locates README.md by searching upward from the current working
// directory. When the working directory is unavailable or no local README.md
// exists, it returns the GitHub URL instead.
func findReadme() string {
	dir, err := os.Getwd()
	if err != nil {
		return readmeGitHubURL
	}
	return findReadmeFrom(dir)
}
