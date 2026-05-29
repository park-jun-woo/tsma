//ff:func feature=coverage type=helper control=iteration dimension=1
//ff:what Flattens one JaCoCo package's source files into path-qualified file entries
package coverage

import "strings"

// flattenJacocoPackage converts a single JaCoCo package's source files into
// flat file entries, prefixing each source-file name with the package path.
func flattenJacocoPackage(pkg jacocoPackage) []jacocoFile {
	prefix := ""
	if pkg.Name != "" {
		prefix = strings.TrimSuffix(pkg.Name, "/") + "/"
	}
	files := make([]jacocoFile, 0, len(pkg.SourceFiles))
	for _, sf := range pkg.SourceFiles {
		files = append(files, jacocoFile{Path: prefix + sf.Name, Lines: sf.Lines})
	}
	return files
}
