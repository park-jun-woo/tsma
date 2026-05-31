//ff:func feature=cli type=util control=iteration dimension=1
//ff:what Appends items not already present to a slice, preserving order
package cli

// appendUnique appends the items not already present in dst, preserving order.
// Used to build a package's deduplicated union of attributed test functions.
func appendUnique(dst, items []string) []string {
	seen := map[string]struct{}{}
	for _, d := range dst {
		seen[d] = struct{}{}
	}
	for _, it := range items {
		if _, ok := seen[it]; ok {
			continue
		}
		seen[it] = struct{}{}
		dst = append(dst, it)
	}
	return dst
}
