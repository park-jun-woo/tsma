//ff:func feature=endpoint type=helper control=iteration dimension=1
//ff:what Extracts method names from a Python list literal like "GET", "POST"
package endpoint

import (
	"regexp"
	"strings"
)

// parseMethodsList extracts method names from a Python list literal.
// Input examples: `"GET", "POST"` or `'GET', 'POST'`
func parseMethodsList(raw string) []string {
	var methods []string
	re := regexp.MustCompile(`["'](\w+)["']`)
	for _, m := range re.FindAllStringSubmatch(raw, -1) {
		methods = append(methods, strings.ToUpper(m[1]))
	}
	return methods
}
