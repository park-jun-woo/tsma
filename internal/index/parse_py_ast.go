//ff:func feature=index type=helper control=iteration dimension=1 lang=python
//ff:what parsePyAst: converts the ast dump JSON into the same []model.Function the line-based indexPyFile yields, so the matcher and coverage stages are unchanged. Builds the qualified name from relDir + receiver (buildQualifiedName, shared with the line path) and reuses the line path's exported heuristic (leading uppercase, no leading underscore). EndLine comes straight from ast end_lineno — the precise range that raises D3 branch accuracy.
package index

import (
	"encoding/json"

	"github.com/park-jun-woo/tsma/internal/model"
)

// parsePyAst unmarshals the ast dump for relPath and returns the functions.
func parsePyAst(data []byte, relPath string) ([]model.Function, error) {
	var raw []pyAstFunc
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	relDir := pkgDirOf(relPath)
	var out []model.Function
	for _, f := range raw {
		out = append(out, pyAstFuncToModel(f, relDir, relPath))
	}
	return out, nil
}
