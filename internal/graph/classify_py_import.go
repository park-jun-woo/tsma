//ff:func feature=graph type=helper control=sequence
//ff:what Classifies a single Python import line and adds entries to the import map
package graph

// classifyPyImport classifies a single import line.
func classifyPyImport(trimmed string, imports map[string]string) {
	if m := pyImportAbsPattern.FindStringSubmatch(trimmed); m != nil {
		module := m[1]
		alias := m[2]
		if alias != "" {
			imports[alias] = module
		} else {
			imports[module] = module
		}
		return
	}

	m := pyFromImportPattern.FindStringSubmatch(trimmed)
	if m == nil {
		return
	}

	classifyPyFromImportNames(m[1], m[2], imports)
}
