//ff:func feature=chain type=helper control=sequence
//ff:what Dispatches a single import line to the appropriate parser based on prefix
package chain

import "strings"

// classifyImportLine parses a single import line and adds entries to the imports map.
func classifyImportLine(trimmed string, imports map[string]string) {
	if pyRelativeImportRe.MatchString(trimmed) {
		parseRelativeImport(trimmed, imports)
		return
	}

	if strings.HasPrefix(trimmed, "import ") {
		parseAbsoluteImport(trimmed, imports)
		return
	}

	if strings.HasPrefix(trimmed, "from ") {
		parseFromImport(trimmed, imports)
	}
}
