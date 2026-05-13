//ff:func feature=match type=helper control=iteration dimension=1
//ff:what Checks if a filename matches one of the expected test suffixes for a source name
package match

import "path/filepath"

// matchTSByFilename checks whether entryName matches srcName + any test suffix.
func matchTSByFilename(projectRoot, dir, entryName, srcName string) (bool, string) {
	for _, suffix := range tsTestSuffixes {
		if entryName != srcName+suffix {
			continue
		}
		absPath := filepath.Join(dir, entryName)
		rel, err := filepath.Rel(projectRoot, absPath)
		if err != nil {
			return true, absPath
		}
		return true, rel
	}
	return false, ""
}
