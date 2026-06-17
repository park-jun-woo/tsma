//ff:type feature=index type=model lang=python
//ff:what pyAstFunc: the JSON DTO unmarshaled from pyAstDefScript's stdout — one entry per Python function the `ast` walk found. parsePyAst converts it into a model.Function. Receiver is the owning class for a method, "" for a top-level or nested function; Decorators is captured for fidelity (not yet surfaced on model.Function).
package index

// pyAstFunc is one function as emitted by the embedded ast dump script.
type pyAstFunc struct {
	Name       string   `json:"name"`
	StartLine  int      `json:"start_line"`
	EndLine    int      `json:"end_line"`
	Col        int      `json:"col"`
	Receiver   string   `json:"receiver"`
	Decorators []string `json:"decorators"`
}
