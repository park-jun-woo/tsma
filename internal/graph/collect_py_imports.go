//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what Extracts import name to module path mappings from Python import statements
package graph

import "strings"

// collectPyImports extracts import name -> module path mappings from Python imports.
func collectPyImports(lines []string) map[string]string {
	imports := make(map[string]string)
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		classifyPyImport(trimmed, imports)
	}
	return imports
}
