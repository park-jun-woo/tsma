//ff:func feature=cli type=helper control=sequence
//ff:what Returns the current working directory as the project root
package cli

import (
	"fmt"
	"os"
)

// getProjectRoot returns the current working directory as the project root.
func getProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}
	return dir, nil
}
