//ff:func feature=coverage type=helper control=iteration dimension=2
//ff:what Finds the matching file entry across all data blocks in llvm-cov output
package coverage

// findRsCoverageFile finds the file entry whose filename matches the target
// relative path, scanning every data block in the llvm-cov export.
func findRsCoverageFile(cov *llvmCovJSON, file, projectRoot string) *llvmCovFile {
	for di := range cov.Data {
		for fi := range cov.Data[di].Files {
			f := &cov.Data[di].Files[fi]
			if matchesPyPath(f.Filename, file, projectRoot) {
				return f
			}
		}
	}
	return nil
}
