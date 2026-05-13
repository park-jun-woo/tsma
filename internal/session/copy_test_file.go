//ff:func feature=session type=implementation control=sequence
//ff:what Copies a test file into .tsma/tests/ and returns the relative path
package session

import (
	"fmt"
	"os"
	"path/filepath"
)

// CopyTestFile copies a test file into .tsma/tests/ and returns the
// destination path relative to the project root.
func CopyTestFile(projectRoot, srcPath string) (string, error) {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return "", fmt.Errorf("read test file: %w", err)
	}
	base := filepath.Base(srcPath)
	dstDir := filepath.Join(Dir(projectRoot), testDir)
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return "", fmt.Errorf("create tests dir: %w", err)
	}
	dst := filepath.Join(dstDir, base)
	if err := os.WriteFile(dst, data, 0o644); err != nil {
		return "", fmt.Errorf("write test copy: %w", err)
	}
	rel, _ := filepath.Rel(projectRoot, dst)
	return rel, nil
}
