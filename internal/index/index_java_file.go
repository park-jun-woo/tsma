//ff:func feature=index type=implementation control=iteration dimension=1
//ff:what Parses a single Java file and extracts method declarations via brace-scope tracking
package index

import (
	"bufio"
	"os"

	"github.com/park-jun-woo/tsma/internal/model"
)

// indexJavaFile parses a single Java file and extracts method/constructor
// declarations. It tracks brace depth to associate methods with their enclosing
// class/interface scopes and records the package name for qualified names.
func indexJavaFile(relPath, absPath string) []model.Function {
	f, err := os.Open(absPath)
	if err != nil {
		return nil
	}
	defer f.Close()

	st := &javaParseState{relDir: pkgDirOf(relPath), relPath: relPath}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		processJavaLine(st, scanner.Text())
	}

	closePrevTSEndLineAtEOF(st.functions, st.lastNonEmptyLine)

	return st.functions
}
