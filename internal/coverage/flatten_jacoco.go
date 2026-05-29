//ff:func feature=coverage type=helper control=iteration dimension=1
//ff:what Flattens a JaCoCo report's package/sourcefile nesting into path-qualified file entries
package coverage

// flattenJacoco collapses the package/sourcefile nesting of a JaCoCo report
// into a flat list of files, prefixing each source-file name with its package
// path so the result carries a full relative path (e.g. com/example/Foo.java).
func flattenJacoco(report *jacocoReport) *jacocoCoverage {
	cov := &jacocoCoverage{}
	for _, pkg := range report.Packages {
		cov.Files = append(cov.Files, flattenJacocoPackage(pkg)...)
	}
	return cov
}
