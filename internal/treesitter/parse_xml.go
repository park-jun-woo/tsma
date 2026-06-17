//ff:func feature=index type=helper control=iteration dimension=1
//ff:what ParseXML: streams `tree-sitter parse --xml` output into a []Source tree via an xmlTreeBuilder. The marshalling half of the shared pipeline — unit-testable with canned XML (no CLI), and the place graceful-fallback errors originate. The loop stays flat: each token is handed to the builder.
package treesitter

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
)

// ParseXML parses the `<sources>` XML emitted by `tree-sitter parse --xml` into
// one Source per file. It returns an error on malformed XML so callers can fall
// back to a line-based path.
func ParseXML(data []byte) ([]Source, error) {
	dec := xml.NewDecoder(bytes.NewReader(data))
	b := &xmlTreeBuilder{}
	for {
		tok, err := dec.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		b.consume(tok)
	}
	return b.sources, nil
}
