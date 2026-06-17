//ff:func feature=index type=helper control=sequence dimension=1
//ff:what (xmlTreeBuilder).handleStart: opens an element — toggles the <sources> guard, starts a new <source>, or creates a Node (attaching it under the stack top or as the source root via attach) and pushes it. Flattened from ParseXML so no branch exceeds depth 2.
package treesitter

import "encoding/xml"

// handleStart processes an XML start element.
func (b *xmlTreeBuilder) handleStart(t xml.StartElement) {
	name := t.Name.Local
	if name == "sources" {
		b.inSources = true
		return
	}
	if name == "source" {
		b.curSourceName = attrValue(t, "name")
		b.stack = nil
		b.curRoot = nil
		return
	}
	if !b.inSources {
		return
	}
	node := &Node{
		Type:  name,
		Field: attrValue(t, "field"),
		SRow:  attrInt(t, "srow"),
		SCol:  attrInt(t, "scol"),
		ERow:  attrInt(t, "erow"),
		ECol:  attrInt(t, "ecol"),
	}
	b.attach(node)
	b.stack = append(b.stack, node)
}
