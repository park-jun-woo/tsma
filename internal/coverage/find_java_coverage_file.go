//ff:func feature=coverage type=helper control=iteration dimension=1
//ff:what Finds the JaCoCo file entry whose package-qualified path matches the target source file
package coverage

// findJavaCoverageFile finds the flattened JaCoCo file whose package-qualified
// path matches the target relative path. JaCoCo paths are package-relative
// (e.g. com/example/Foo.java) while the target is project-relative (e.g.
// src/main/java/com/example/Foo.java), so a suffix match is used.
func findJavaCoverageFile(cov *jacocoCoverage, file, projectRoot string) *jacocoFile {
	for i := range cov.Files {
		if matchesPyPath(file, cov.Files[i].Path, projectRoot) {
			return &cov.Files[i]
		}
	}
	return nil
}
