//ff:func feature=endpoint type=helper control=sequence
//ff:what Strips module prefixes from Django view references like views.my_handler
package endpoint

import "strings"

// extractViewName strips module prefixes, e.g. "views.my_handler" -> "my_handler".
// Also handles ".as_view()" suffix from class-based views.
func extractViewName(raw string) string {
	raw = strings.TrimSuffix(raw, ".as_view")
	parts := strings.Split(raw, ".")
	return parts[len(parts)-1]
}
