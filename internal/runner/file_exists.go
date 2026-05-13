//ff:func feature=runner type=helper control=sequence
//ff:what Checks if a file exists at the given path
package runner

import "os"

// fileExists checks if a file exists at the given path.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
