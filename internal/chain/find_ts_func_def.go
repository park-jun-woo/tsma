//ff:func feature=chain type=implementation control=sequence
//ff:what Walks project TS/JS files searching for a named function definition via regex
package chain

import (
	"os"
	"path/filepath"
)

// findTSFuncDef searches the project for a function definition matching the given name.
// Returns nil if not found.
func findTSFuncDef(projectRoot, funcName, currentFile string) *tsFuncDef {
	pattern := tsFuncDefSearchRegex(funcName)

	// First, search in the current file (most likely location).
	absCurrentFile := filepath.Join(projectRoot, currentFile)
	if def := searchFileForFunc(absCurrentFile, currentFile, funcName, pattern); def != nil {
		return def
	}

	// Walk the project to find the definition elsewhere.
	var result *tsFuncDef
	filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil || result != nil {
			return nil
		}
		if info.IsDir() {
			return skipTSDir(path)
		}
		if !isTSOrJSSourceFile(path) {
			return nil
		}

		relPath, _ := filepath.Rel(projectRoot, path)
		if relPath == currentFile {
			return nil // already searched
		}

		if def := searchFileForFunc(path, relPath, funcName, pattern); def != nil {
			result = def
		}
		return nil
	})

	return result
}
