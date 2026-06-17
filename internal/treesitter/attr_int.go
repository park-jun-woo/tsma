//ff:func feature=index type=helper control=sequence
//ff:what attrInt: returns a named XML attribute parsed as an int, or 0 — used for the 0-based row/col positions.
package treesitter

import (
	"encoding/xml"
	"strconv"
)

// attrInt returns the named attribute parsed as an int, or 0.
func attrInt(e xml.StartElement, name string) int {
	v, _ := strconv.Atoi(attrValue(e, name))
	return v
}
