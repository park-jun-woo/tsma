//ff:func feature=chain type=implementation control=iteration dimension=1
//ff:what Reads a Python file and classifies imports as external or internal
package chain

import (
	"bufio"
	"os"
	"strings"
)

// collectImports reads a Python file and classifies imports.
// Returns a map of identifier -> "external" or "internal".
func collectImports(filePath string) map[string]string {
	imports := make(map[string]string)

	f, err := os.Open(filePath)
	if err != nil {
		return imports
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		classifyImportLine(trimmed, imports)
	}

	return imports
}
