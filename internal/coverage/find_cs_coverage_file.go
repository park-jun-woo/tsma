//ff:func feature=coverage type=helper control=iteration dimension=1 lang=csharp
//ff:what Finds the Cobertura file entry whose source path matches the target source file
package coverage

// findCsCoverageFile finds the flattened Cobertura file whose source-relative
// path matches the target project-relative path. Cobertura class filenames are
// relative to the report's <source> root, so a suffix match (via matchesPyPath)
// is used to tolerate differing path prefixes.
func findCsCoverageFile(cov *csCoverage, file, projectRoot string) *csFile {
	for i := range cov.Files {
		if matchesPyPath(file, cov.Files[i].Path, projectRoot) {
			return &cov.Files[i]
		}
	}
	return nil
}
