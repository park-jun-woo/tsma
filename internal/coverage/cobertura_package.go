//ff:type feature=coverage type=model lang=csharp
//ff:what Represents a <package> element of a Cobertura XML export holding its classes
package coverage

// coberturaPackage is a <package> element. Its classes carry per-file line
// coverage; class filenames are paths relative to the configured <source> root.
type coberturaPackage struct {
	Name    string           `xml:"name,attr"`
	Classes []coberturaClass `xml:"classes>class"`
}
