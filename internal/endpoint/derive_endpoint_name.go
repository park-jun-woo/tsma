//ff:func feature=endpoint type=helper control=iteration dimension=1
//ff:what Converts a file path to a PascalCase endpoint name
package endpoint

import (
	"path/filepath"
	"strings"
)

// deriveEndpointName derives a PascalCase endpoint name from a file path.
// e.g., "app/api/users/route.ts" -> "ApiUsers"
// e.g., "pages/api/users/[id].ts" -> "ApiUsersId"
func deriveEndpointName(relPath string) string {
	p := filepath.ToSlash(relPath)
	p = strings.TrimPrefix(p, "src/")

	ext := filepath.Ext(p)
	p = strings.TrimSuffix(p, ext)

	p = strings.TrimPrefix(p, "app/")
	p = strings.TrimPrefix(p, "pages/")

	p = strings.TrimSuffix(p, "/route")
	p = strings.TrimSuffix(p, "/index")

	segments := strings.Split(p, "/")
	var parts []string
	for _, seg := range segments {
		if seg == "" {
			continue
		}
		seg = strings.NewReplacer("[", "", "]", "", "...", "").Replace(seg)
		if seg == "" {
			continue
		}
		parts = append(parts, capitalize(seg))
	}

	if len(parts) == 0 {
		return "Root"
	}
	return strings.Join(parts, "")
}
