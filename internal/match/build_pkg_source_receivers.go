//ff:func feature=match type=implementation control=iteration dimension=1 lang=go
//ff:what Builds a name->receiver-set map from a package's non-test Go sources
package match

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// BuildPkgSourceReceivers parses every non-test .go file in the package
// directory pkgDir (relative to projectRoot) and builds a PkgSourceReceivers
// mapping each declared identifier name to the set of distinguishers it appears
// with: a method's bare receiver type (pointer/value/generic normalized via
// srcReceiver), or "" for a free function. It is the source-side counterpart of
// BuildPkgTestIndex — same per-directory, build-once lifetime — and lets
// receiver-aware lookup decide whether a name is same-name-single or
// same-name-multiple without depending on the batch caller's fns slice (so the
// single-func MatchFunc path resolves it identically). Unparseable files are
// skipped so one bad file does not abort the map. pkgDir "" or "." means the
// root. A read error on the directory is returned; a nil map is never returned
// alongside a nil error.
func BuildPkgSourceReceivers(projectRoot, pkgDir string) (*PkgSourceReceivers, error) {
	r := &PkgSourceReceivers{byName: make(map[string]map[string]struct{})}
	absDir := filepath.Join(projectRoot, pkgDir)
	entries, err := os.ReadDir(absDir)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		abs := filepath.Join(absDir, name)
		fset := token.NewFileSet()
		f, perr := parser.ParseFile(fset, abs, nil, parser.AllErrors)
		if perr != nil {
			continue
		}
		recordSourceDecls(r, f)
	}
	return r, nil
}
