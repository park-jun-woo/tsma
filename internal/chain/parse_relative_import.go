//ff:func feature=chain type=helper control=iteration dimension=1
//ff:what Parses a Python relative import line and adds identifiers as internal
package chain

import "strings"

// parseRelativeImport parses a "from .foo import bar" line.
func parseRelativeImport(trimmed string, imports map[string]string) {
	idx := strings.Index(trimmed, "import")
	if idx < 0 {
		return
	}
	names := trimmed[idx+len("import"):]
	for _, name := range strings.Split(names, ",") {
		name = strings.TrimSpace(name)
		fields := strings.Fields(name)
		if len(fields) >= 3 && fields[1] == "as" {
			imports[fields[2]] = "internal"
		} else if len(fields) > 0 {
			imports[fields[0]] = "internal"
		}
	}
}
