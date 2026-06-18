//ff:func feature=detect type=helper control=sequence lang=python
//ff:what Reports whether a TOML section header is a PEP 621 dependency table
package detect

import "strings"

// isDepTableHeader reports whether a trimmed TOML section header line opens a
// PEP 621 dependency table whose entries may name pytest:
// [project.optional-dependencies] / [dependency-groups] and their subtables.
func isDepTableHeader(line string) bool {
	return line == "[project.optional-dependencies]" ||
		line == "[dependency-groups]" ||
		strings.HasPrefix(line, "[project.optional-dependencies.") ||
		strings.HasPrefix(line, "[dependency-groups.")
}
