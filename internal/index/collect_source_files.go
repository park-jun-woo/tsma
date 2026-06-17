//ff:func feature=index type=helper control=sequence
//ff:what collectSourceFiles: walks a project tree honoring .tsmignore + the language's skipDir/isSource filters and returns the indexable files (rel+abs). Language-neutral — the same walk the line-based Index methods do, factored out so the tree-sitter pipeline can batch every file into one CLI invocation.
package index

import (
	"os"
	"path/filepath"
)

// collectSourceFiles returns every indexable source file under projectRoot,
// applying the same .tsmignore + skipDir + isSource filtering the per-language
// line-based indexers use.
func collectSourceFiles(projectRoot string, isSource func(string) bool, skipDir func(string) error) ([]sourceFile, error) {
	ignorePatterns := ParseTsmIgnore(filepath.Join(projectRoot, ".tsmignore"))
	var files []sourceFile

	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		relPath, _ := filepath.Rel(projectRoot, path)
		if MatchTsmIgnore(relPath, info.Name(), info.IsDir(), ignorePatterns) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.IsDir() {
			return skipDir(path)
		}
		if !isSource(path) {
			return nil
		}
		files = append(files, sourceFile{rel: relPath, abs: path})
		return nil
	})

	return files, err
}
