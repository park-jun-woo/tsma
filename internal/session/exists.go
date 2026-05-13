//ff:func feature=session type=helper control=sequence
//ff:what Checks whether a session file exists on disk
package session

import (
	"os"
	"path/filepath"
)

// Exists checks whether a session file exists.
func Exists(projectRoot string) bool {
	p := filepath.Join(Dir(projectRoot), sessionFile)
	_, err := os.Stat(p)
	return err == nil
}
