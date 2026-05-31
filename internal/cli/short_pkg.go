//ff:func feature=cli type=helper control=sequence
//ff:what Renders a package directory relative to the project root for progress logging
package cli

import "path/filepath"

// shortPkg renders a package directory relative to the project root for the
// "Measuring coverage: pkg i/N (..)" progress logs, falling back to the absolute
// path if it cannot be made relative.
func shortPkg(root, pkgDir string) string {
	if rel, err := filepath.Rel(root, pkgDir); err == nil {
		return rel
	}
	return pkgDir
}
