//ff:func feature=match type=helper control=iteration dimension=1 lang=python
//ff:what ingestPyTestFile: ast-parses one Python test file and records, for every referenced name, a back-reference to that file in the index. The Python analogue of ingestTSTestFile; files that fail to parse contribute nothing so one bad file never aborts the package index.
package match

// ingestPyTestFile appends rel to every referenced name's bucket in idx.
func ingestPyTestFile(idx *PyPkgTestIndex, python, absPath, rel string) {
	for _, name := range collectPyCalledNames(python, absPath) {
		idx.refs[name] = append(idx.refs[name], rel)
	}
}
