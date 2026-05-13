//ff:func feature=chain type=implementation control=sequence
//ff:what Walks the project tree and indexes all Python function definitions
package chain

import (
	"os"
	"path/filepath"
	"strings"
)

// indexPyFunctions walks the project and indexes all Python function definitions.
func indexPyFunctions(projectRoot string) (map[string]*pyFuncInfo, error) {
	funcs := make(map[string]*pyFuncInfo)

	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if pyTracerSkipDirs[filepath.Base(path)] {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".py") {
			return nil
		}

		relPath, _ := filepath.Rel(projectRoot, path)
		parsedFuncs := parsePyFunctions(path, relPath)
		for k, v := range parsedFuncs {
			funcs[k] = v
		}
		return nil
	})

	return funcs, err
}
