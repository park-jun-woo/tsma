//ff:func feature=chain type=helper control=iteration dimension=1
//ff:what Parses a Python from-import line and classifies identifiers as external or internal
package chain

import "strings"

// parseFromImport parses a "from module import name" line.
func parseFromImport(trimmed string, imports map[string]string) {
	parts := strings.SplitN(trimmed, "import", 2)
	if len(parts) != 2 {
		return
	}
	modPart := strings.TrimSpace(strings.TrimPrefix(parts[0], "from"))
	names := parts[1]

	classification := "external"
	if strings.HasPrefix(modPart, ".") {
		classification = "internal"
	}

	for _, name := range strings.Split(names, ",") {
		name = strings.TrimSpace(name)
		fields := strings.Fields(name)
		if len(fields) >= 3 && fields[1] == "as" {
			imports[fields[2]] = classification
		} else if len(fields) > 0 {
			imports[fields[0]] = classification
		}
	}
}
