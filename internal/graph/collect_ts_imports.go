//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what Extracts import name to module path mappings from TS/JS import statements
package graph

// collectTSImports extracts import name -> module path mappings.
func collectTSImports(lines []string) map[string]string {
	imports := make(map[string]string)
	for _, line := range lines {
		m := tsNamedImportPattern.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		classifyTSImportNames(m[1], m[2], imports)
	}
	return imports
}
