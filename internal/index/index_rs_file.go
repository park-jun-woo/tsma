//ff:func feature=index type=implementation control=iteration dimension=1
//ff:what Parses a single Rust file and extracts fn/impl-method declarations via brace-scope tracking
package index

import (
	"bufio"
	"os"

	"github.com/park-jun-woo/tsma/internal/model"
)

// indexRsFile parses a single Rust file and extracts function declarations.
// It tracks brace depth to associate fns with their enclosing impl/mod scopes,
// and skips functions inside #[cfg(test)] modules.
func indexRsFile(relPath, absPath string) []model.Function {
	f, err := os.Open(absPath)
	if err != nil {
		return nil
	}
	defer f.Close()

	st := &rsParseState{relDir: pkgDirOf(relPath), relPath: relPath}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		processRsLine(st, scanner.Text())
	}

	closePrevTSEndLineAtEOF(st.functions, st.lastNonEmptyLine)

	return st.functions
}
