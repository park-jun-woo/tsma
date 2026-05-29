//ff:func feature=coverage type=helper control=iteration dimension=2 lang=csharp
//ff:what Flattens a Cobertura report's package/class nesting into path-qualified file entries
package coverage

// flattenCobertura collapses the package/class nesting of a Cobertura report
// into a flat list of files keyed by class filename. Classes that share a
// filename (C# partial classes) are merged so all of a file's lines appear in a
// single entry, preserving the order of first appearance.
func flattenCobertura(report *coberturaReport) *csCoverage {
	cov := &csCoverage{}
	index := map[string]int{}
	for _, pkg := range report.Packages {
		for _, cls := range pkg.Classes {
			if cls.Filename == "" {
				continue
			}
			if i, ok := index[cls.Filename]; ok {
				cov.Files[i].Lines = append(cov.Files[i].Lines, cls.Lines...)
				continue
			}
			index[cls.Filename] = len(cov.Files)
			cov.Files = append(cov.Files, csFile{Path: cls.Filename, Lines: cls.Lines})
		}
	}
	return cov
}
