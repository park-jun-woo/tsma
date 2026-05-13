//ff:func feature=session type=helper control=sequence
//ff:what Returns the .tsma directory path for the given project root
package session

import "path/filepath"

// Dir returns the .tsma directory path for the given project root.
func Dir(projectRoot string) string {
	return filepath.Join(projectRoot, dirName)
}
