//ff:func feature=index type=helper control=iteration dimension=1
//ff:what attrValue: returns the value of a named XML attribute on a start element, or "" — used by the builder to read field/srow/scol/erow/ecol.
package treesitter

import "encoding/xml"

// attrValue returns the value of the named attribute, or "".
func attrValue(e xml.StartElement, name string) string {
	for _, a := range e.Attr {
		if a.Name.Local == name {
			return a.Value
		}
	}
	return ""
}
