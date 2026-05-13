//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what Parses comma-separated imported names from a TS/JS import statement
package graph

import "strings"

// classifyTSImportNames parses imported names from an import line.
func classifyTSImportNames(namesStr, modulePath string, imports map[string]string) {
	names := strings.Split(namesStr, ",")
	for _, n := range names {
		n = strings.TrimSpace(n)
		if n == "" {
			continue
		}
		parts := strings.Fields(n)
		if len(parts) == 3 && parts[1] == "as" {
			imports[parts[2]] = modulePath
		} else {
			imports[parts[0]] = modulePath
		}
	}
}
