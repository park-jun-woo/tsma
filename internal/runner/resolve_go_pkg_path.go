//ff:func feature=runner type=helper control=sequence
//ff:what Resolves the Go package path from the test file location
package runner

import (
	"fmt"
	"path/filepath"
)

// resolveGoPkgPath resolves the Go package path from the test file location.
func resolveGoPkgPath(projectRoot, absTest string) (string, error) {
	pkgDir := filepath.Dir(absTest)
	relPkg, err := filepath.Rel(projectRoot, pkgDir)
	if err != nil {
		return "", fmt.Errorf("resolve relative package: %w", err)
	}
	return "./" + filepath.ToSlash(relPkg), nil
}
