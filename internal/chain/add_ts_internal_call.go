//ff:func feature=chain type=implementation control=sequence
//ff:what Adds an internal TS function call to the chain and recurses
package chain

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/tsma/internal/model"
)

// addTSInternalCall adds an internal TS function call to the chain and recurses.
func addTSInternalCall(def *tsFuncDef, displayName, projectRoot string, visited map[string]bool, entries *[]model.ChainEntry, depth int) {
	key := fmt.Sprintf("%s:%d", def.file, def.startLine)
	if visited[key] {
		return
	}
	visited[key] = true

	*entries = append(*entries, model.ChainEntry{
		Func:      displayName,
		File:      def.file,
		StartLine: def.startLine,
		EndLine:   def.endLine,
	})

	absDefFile := filepath.Join(projectRoot, def.file)
	defData, err := os.ReadFile(absDefFile)
	if err != nil {
		return
	}
	defLines := strings.Split(string(defData), "\n")
	traceTSFunc(projectRoot, def.file, defLines, def.startLine, def.endLine, visited, entries, depth+1)
}
