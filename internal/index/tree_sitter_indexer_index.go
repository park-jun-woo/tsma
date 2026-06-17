//ff:func feature=index type=implementation control=sequence dimension=1
//ff:what (TreeSitterIndexer).Index: the precise indexing path. When tree-sitter is unavailable it delegates wholesale to the line-based fallback (zero regression); otherwise it batch-parses every source file and runs the language extractor, per-file-falling-back to the line-based path for any file tree-sitter failed to parse. Output is the same []model.Function the fallback yields, so matcher/coverage are unchanged.
package index

import (
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// Index collects function declarations using tree-sitter, falling back to the
// line-based indexer when the CLI/grammar is absent or a file fails to parse.
func (t *TreeSitterIndexer) Index(projectRoot string) ([]model.Function, error) {
	if !t.available() {
		return t.fallback.Index(projectRoot)
	}

	files, err := collectSourceFiles(projectRoot, t.isSource, t.skipDir)
	if err != nil {
		return t.fallback.Index(projectRoot)
	}
	if len(files) == 0 {
		return nil, nil
	}

	abs := make([]string, len(files))
	for i, f := range files {
		abs[i] = f.abs
	}

	out, runErr := treesitter.Run(t.command, t.grammarDir, abs)
	if runErr != nil || len(out) == 0 {
		return t.fallbackFiles(files), nil
	}

	sources, perr := treesitter.ParseXML(out)
	if perr != nil {
		return t.fallbackFiles(files), nil
	}

	byName := make(map[string]*treesitter.Node, len(sources))
	for i := range sources {
		byName[sources[i].Name] = sources[i].Root
	}

	var functions []model.Function
	for _, f := range files {
		root := byName[f.abs]
		if root == nil {
			functions = append(functions, t.fileFallback(f.rel, f.abs)...)
			continue
		}
		functions = append(functions, t.extract(root, f.rel, pkgDirOf(f.rel))...)
	}
	return functions, nil
}
