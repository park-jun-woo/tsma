//ff:func feature=match type=helper control=sequence lang=go
//ff:what Looks up a function's test references by name then filters them by receiver
package match

import "github.com/park-jun-woo/tsma/internal/model"

// MatchFuncByName returns the test references attributed to a function within
// its package, using a prebuilt content-aware index keyed by bare name plus the
// package's source-receiver map for receiver-aware filtering. It looks the
// function up by its bare model.Function.Name (consistent with the indexer's
// bare-name keys), then, when the function is a method (fn.Receiver != ""),
// filters the name bucket so only refs belonging to fn's type are returned
// (filterRefsByReceiver): exact-receiver refs always, unknown-receiver refs only
// when the name is same-name-single in srcReceivers. A free function returns the
// whole bucket unchanged. found is false when the name is absent or the
// receiver filter leaves no refs, so the caller falls back to file-name
// matching. srcReceivers may be nil (sources unparsed), in which case filtering
// is non-regressing (unknown refs kept).
func MatchFuncByName(idx *PkgTestIndex, srcReceivers *PkgSourceReceivers, fn *model.Function) ([]testRef, bool) {
	if idx == nil || fn == nil {
		return nil, false
	}
	refs, ok := idx.Lookup(fn.Name)
	if !ok {
		return nil, false
	}
	filtered := filterRefsByReceiver(refs, srcReceivers, fn)
	if len(filtered) == 0 {
		return nil, false
	}
	return filtered, true
}
