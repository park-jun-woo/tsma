//ff:func feature=match type=helper control=iteration dimension=1
//ff:what Strips TS/JS file extensions to get the base name for matching
package match

import "strings"

// stripTSExtension removes .ts, .js, .tsx, .jsx extensions from a filename.
func stripTSExtension(name string) string {
	for _, ext := range []string{".tsx", ".jsx", ".ts", ".js"} {
		if strings.HasSuffix(name, ext) {
			return strings.TrimSuffix(name, ext)
		}
	}
	return name
}
