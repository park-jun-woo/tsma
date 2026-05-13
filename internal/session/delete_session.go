//ff:func feature=session type=implementation control=sequence
//ff:what Removes the entire .tsma directory
package session

import (
	"errors"
	"os"
)

// Delete removes the entire .tsma directory.
func Delete(projectRoot string) error {
	dir := Dir(projectRoot)
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return os.RemoveAll(dir)
}
