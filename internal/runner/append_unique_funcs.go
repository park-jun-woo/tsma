//ff:func feature=runner type=util control=iteration dimension=1 lang=go
//ff:what Appends names not yet in the seen set to dst, preserving order
package runner

// appendUniqueFuncs appends each name from src that is not already in seen to
// dst, recording it in seen, and returns the extended slice. It preserves
// first-seen order and is used to build the deduplicated union of test function
// names across multiple files.
func appendUniqueFuncs(seen map[string]struct{}, dst, src []string) []string {
	for _, name := range src {
		if _, dup := seen[name]; dup {
			continue
		}
		seen[name] = struct{}{}
		dst = append(dst, name)
	}
	return dst
}
