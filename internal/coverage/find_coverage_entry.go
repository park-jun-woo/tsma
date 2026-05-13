//ff:func feature=coverage type=helper control=iteration dimension=1
//ff:what Matches a relative file path against coverage-final.json keys using suffix comparison
package coverage

import "path/filepath"

// findCoverageEntry finds the coverage entry matching the given relative file path.
func findCoverageEntry(relFile string, coverageData map[string]coverageFinalEntry, projectRoot string) *coverageFinalEntry {
	normalizedRel := filepath.ToSlash(relFile)

	for key, entry := range coverageData {
		if matchesCoverageKey(key, normalizedRel, relFile, projectRoot) {
			e := entry
			return &e
		}
	}

	return nil
}
