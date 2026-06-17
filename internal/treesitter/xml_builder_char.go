//ff:func feature=index type=helper control=sequence
//ff:what (xmlTreeBuilder).handleChar: attaches trimmed character data to the current node's Text, but only the first non-empty run — so leaf identifiers keep their text while interior nodes stay empty.
package treesitter

import (
	"encoding/xml"
	"strings"
)

// handleChar records the first non-empty character data on the current node.
func (b *xmlTreeBuilder) handleChar(t xml.CharData) {
	if len(b.stack) == 0 {
		return
	}
	text := strings.TrimSpace(string(t))
	if text == "" {
		return
	}
	top := b.stack[len(b.stack)-1]
	if top.Text == "" {
		top.Text = text
	}
}
