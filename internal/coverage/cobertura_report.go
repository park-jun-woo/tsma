//ff:type feature=coverage type=model lang=csharp
//ff:what Represents the top-level <coverage> element of a Cobertura XML export
package coverage

import "encoding/xml"

// coberturaReport is the root <coverage> element of a Cobertura XML export
// (the format produced by coverlet via `dotnet test --collect:"XPlat Code
// Coverage"`).
type coberturaReport struct {
	XMLName  xml.Name           `xml:"coverage"`
	Packages []coberturaPackage `xml:"packages>package"`
}
