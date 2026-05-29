//ff:type feature=coverage type=model
//ff:what Represents the top-level <report> element of a JaCoCo XML coverage export
package coverage

import "encoding/xml"

// jacocoReport is the root <report> element of a JaCoCo XML export.
type jacocoReport struct {
	XMLName  xml.Name        `xml:"report"`
	Packages []jacocoPackage `xml:"package"`
}
