//ff:func feature=index type=helper control=selection dimension=1
//ff:what (xmlTreeBuilder).consume: dispatches one tree-sitter XML token to the start/end/chardata handler. Keeps ParseXML's loop body a single call so the switch lives here, not nested inside the loop (depth ≤2).
package treesitter

import "encoding/xml"

// consume dispatches one XML token to the matching handler.
func (b *xmlTreeBuilder) consume(tok xml.Token) {
	switch t := tok.(type) {
	case xml.StartElement:
		b.handleStart(t)
	case xml.EndElement:
		b.handleEnd(t)
	case xml.CharData:
		b.handleChar(t)
	}
}
