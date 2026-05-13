//ff:func feature=coverage type=helper control=iteration dimension=1
//ff:what Finds the matching file entry in Python coverage data
package coverage

// findPyCoverageFile finds the matching file entry in Python coverage data.
func findPyCoverageFile(covData *pyCoverageJSON, file, projectRoot string) *pyCoverageFile {
	for covPath, cov := range covData.Files {
		if matchesPyPath(covPath, file, projectRoot) {
			c := cov
			return &c
		}
	}
	return nil
}
