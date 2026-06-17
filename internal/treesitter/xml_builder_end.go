//ff:func feature=index type=helper control=sequence dimension=1
//ff:what (xmlTreeBuilder).handleEnd: closes an element — clears the <sources> guard, finalizes a <source> into Sources (with its captured root), or pops the node stack. Flattened from ParseXML (depth ≤2).
package treesitter

import "encoding/xml"

// handleEnd processes an XML end element.
func (b *xmlTreeBuilder) handleEnd(t xml.EndElement) {
	name := t.Name.Local
	if name == "sources" {
		b.inSources = false
		return
	}
	if name == "source" {
		b.sources = append(b.sources, Source{Name: b.curSourceName, Root: b.curRoot})
		b.stack = nil
		b.curRoot = nil
		return
	}
	if len(b.stack) > 0 {
		b.stack = b.stack[:len(b.stack)-1]
	}
}
