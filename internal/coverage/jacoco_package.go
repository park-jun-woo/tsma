//ff:type feature=coverage type=model
//ff:what Represents a <package> element of a JaCoCo XML export holding its source files
package coverage

// jacocoPackage is a <package> element. Its Name is a slash-separated package
// path (e.g. "com/example") prepended to source-file names to form the full
// relative path.
type jacocoPackage struct {
	Name        string             `xml:"name,attr"`
	SourceFiles []jacocoSourceFile `xml:"sourcefile"`
}
