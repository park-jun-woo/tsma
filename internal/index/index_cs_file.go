//ff:func feature=index type=implementation control=iteration dimension=1 lang=csharp
//ff:what Parses a single C# file and extracts method declarations via brace-scope tracking
package index

import (
	"bufio"
	"os"

	"github.com/park-jun-woo/tsma/internal/model"
)

// indexCsFile parses a single C# file and extracts method/constructor
// declarations. It tracks brace depth to associate methods with their enclosing
// namespace/type scopes and records the file-scoped namespace for qualified
// names.
func indexCsFile(relPath, absPath string) []model.Function {
	f, err := os.Open(absPath)
	if err != nil {
		return nil
	}
	defer f.Close()

	st := &csParseState{relDir: pkgDirOf(relPath), relPath: relPath}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		processCsLine(st, scanner.Text())
	}

	closePrevTSEndLineAtEOF(st.functions, st.lastNonEmptyLine)

	return st.functions
}
