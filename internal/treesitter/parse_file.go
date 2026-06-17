//ff:func feature=index type=helper control=iteration dimension=1
//ff:what ParseFile: convenience over Run+ParseXML for a single source file — returns its parse-tree root Node by matching the source name (or the sole source). Shared by match (D2 call-site extraction) and smell (D4 node detectors); errors propagate so callers graceful-fallback.
package treesitter

import "fmt"

// ParseFile runs tree-sitter on a single file and returns its root Node. It
// returns an error when the CLI fails or the file produced no parse tree.
func ParseFile(command, grammarDir, absPath string) (*Node, error) {
	out, err := Run(command, grammarDir, []string{absPath})
	if err != nil {
		return nil, err
	}
	sources, err := ParseXML(out)
	if err != nil {
		return nil, err
	}
	for _, s := range sources {
		if s.Name == absPath && s.Root != nil {
			return s.Root, nil
		}
	}
	if len(sources) == 1 && sources[0].Root != nil {
		return sources[0].Root, nil
	}
	return nil, fmt.Errorf("treesitter: no parse tree for %s", absPath)
}
