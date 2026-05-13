//ff:func feature=detect type=implementation control=sequence
//ff:what Reads package.json and detects the TypeScript web framework
package detect

import (
	"os"
	"path/filepath"
)

func detectTSFramework(projectRoot string) string {
	data, err := os.ReadFile(filepath.Join(projectRoot, "package.json"))
	if err != nil {
		return "unknown"
	}
	content := string(data)
	if containsImport(content, "\"next\"") || containsImport(content, "'next'") {
		return "nextjs"
	}
	if containsImport(content, "\"express\"") || containsImport(content, "'express'") {
		return "express"
	}
	return "unknown"
}
