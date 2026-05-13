//ff:func feature=endpoint type=helper control=sequence
//ff:what Converts a file path to an HTTP route path following Next.js conventions
package endpoint

import (
	"path/filepath"
	"strings"
)

// deriveRoutePath derives the HTTP route path from a file's relative path.
// e.g., "app/api/users/route.ts" -> "/api/users"
// e.g., "pages/api/users/index.ts" -> "/api/users"
// e.g., "pages/api/users/[id].ts" -> "/api/users/[id]"
// e.g., "src/app/api/users/route.ts" -> "/api/users"
func deriveRoutePath(relPath string) string {
	p := filepath.ToSlash(relPath)
	p = strings.TrimPrefix(p, "src/")

	dir := filepath.Dir(p)
	base := filepath.Base(p)

	if base == "route.ts" || base == "route.js" {
		dir = strings.TrimPrefix(filepath.ToSlash(dir), "app")
		if dir == "" {
			dir = "/"
		}
		return dir
	}

	dir = filepath.ToSlash(dir)
	dir = strings.TrimPrefix(dir, "pages")

	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	if name == "index" {
		if dir == "" {
			return "/"
		}
		return dir
	}

	return dir + "/" + name
}
