//ff:func feature=chain type=helper control=iteration dimension=1
//ff:what Parses a Python absolute import line and adds identifiers as external
package chain

import "strings"

// parseAbsoluteImport parses an "import module" line.
func parseAbsoluteImport(trimmed string, imports map[string]string) {
	rest := trimmed[len("import "):]
	for _, part := range strings.Split(rest, ",") {
		part = strings.TrimSpace(part)
		fields := strings.Fields(part)
		if len(fields) >= 3 && fields[1] == "as" {
			imports[fields[2]] = "external"
		} else if len(fields) > 0 {
			root := strings.Split(fields[0], ".")[0]
			imports[root] = "external"
		}
	}
}
