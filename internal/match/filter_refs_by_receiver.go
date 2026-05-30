//ff:func feature=match type=helper control=iteration dimension=1 lang=go
//ff:what Filters a name bucket's test refs by the function's receiver, conservatively
package match

import "github.com/park-jun-woo/tsma/internal/model"

// filterRefsByReceiver applies the receiver-aware attribution policy to the
// test refs found in fn's name bucket. For a free function (fn.Receiver == "")
// it returns the bucket unchanged — receiver is irrelevant, behavior is
// unchanged. For a method it keeps only the refs that belong to fn's type:
//
//   - a ref whose call-site receiver equals fn.Receiver is always kept;
//   - a ref with an unknown receiver ("") is kept only when fn.Name is
//     same-name-single in the package (no other declaration shares the name),
//     which preserves the pre-existing correct behavior for the common case;
//     when fn.Name is same-name-multiple the unknown ref is dropped to avoid
//     mis-attributing another type's test to fn.
//
// srcReceivers may be nil (package sources unparsed); isSameNameMultiple then
// reports false, so unknown refs are kept — the conservative non-regressing
// choice. Returning an empty slice means receiver-aware matching attributed
// nothing, and the caller falls back to file-name matching.
func filterRefsByReceiver(refs []testRef, srcReceivers *PkgSourceReceivers, fn *model.Function) []testRef {
	if fn.Receiver == "" {
		return refs
	}
	multiple := srcReceivers.isSameNameMultiple(fn.Name)
	var kept []testRef
	for _, r := range refs {
		if keepRefForReceiver(r, fn.Receiver, multiple) {
			kept = append(kept, r)
		}
	}
	return kept
}
