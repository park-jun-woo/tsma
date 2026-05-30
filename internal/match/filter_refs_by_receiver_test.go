package match

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// srcRecvWith builds a PkgSourceReceivers directly from name->receivers pairs.
func srcRecvWith(t *testing.T, decls map[string][]string) *PkgSourceReceivers {
	t.Helper()
	r := &PkgSourceReceivers{byName: make(map[string]map[string]struct{})}
	for name, recvs := range decls {
		for _, rec := range recvs {
			addNameReceiver(r, name, rec)
		}
	}
	return r
}

func TestFilterRefsByReceiverFreeFuncUnchanged(t *testing.T) {
	refs := []testRef{
		{TestFunc: "TestA", Receiver: ""},
		{TestFunc: "TestB", Receiver: "Whatever"},
	}
	fn := &model.Function{Name: "Parse", Receiver: ""}
	got := filterRefsByReceiver(refs, nil, fn)
	if len(got) != 2 {
		t.Fatalf("free function must return the whole bucket unchanged, got %v", got)
	}
}

func TestFilterRefsByReceiverExactMatch(t *testing.T) {
	refs := []testRef{
		{TestFunc: "TestGo", Receiver: "GoFile"},
		{TestFunc: "TestCS", Receiver: "CSharpFile"},
	}
	src := srcRecvWith(t, map[string][]string{"GetFuncs": {"GoFile", "CSharpFile"}})
	fn := &model.Function{Name: "GetFuncs", Receiver: "CSharpFile"}
	got := filterRefsByReceiver(refs, src, fn)
	if len(got) != 1 || got[0].TestFunc != "TestCS" {
		t.Fatalf("got %v, want only TestCS (CSharpFile)", got)
	}
}

func TestFilterRefsByReceiverUnknownSingleKept(t *testing.T) {
	// receiver unknown ref + same-name-single -> kept (regression guard).
	refs := []testRef{{TestFunc: "TestGo", Receiver: ""}}
	src := srcRecvWith(t, map[string][]string{"GetFuncs": {"GoFile"}})
	fn := &model.Function{Name: "GetFuncs", Receiver: "GoFile"}
	got := filterRefsByReceiver(refs, src, fn)
	if len(got) != 1 || got[0].TestFunc != "TestGo" {
		t.Fatalf("unknown ref + single name must be kept, got %v", got)
	}
}

func TestFilterRefsByReceiverUnknownMultipleDropped(t *testing.T) {
	// receiver unknown ref + same-name-multiple -> dropped (mis-attribution guard).
	refs := []testRef{{TestFunc: "TestGo", Receiver: ""}}
	src := srcRecvWith(t, map[string][]string{"GetFuncs": {"GoFile", "CSharpFile"}})
	fn := &model.Function{Name: "GetFuncs", Receiver: "CSharpFile"}
	got := filterRefsByReceiver(refs, src, fn)
	if len(got) != 0 {
		t.Fatalf("unknown ref + multiple name must be dropped, got %v", got)
	}
}

func TestFilterRefsByReceiverNilSrcNonRegressing(t *testing.T) {
	// nil source map -> isSameNameMultiple false -> unknown refs kept.
	refs := []testRef{{TestFunc: "TestGo", Receiver: ""}}
	fn := &model.Function{Name: "GetFuncs", Receiver: "GoFile"}
	got := filterRefsByReceiver(refs, nil, fn)
	if len(got) != 1 {
		t.Fatalf("nil src must keep unknown refs (non-regressing), got %v", got)
	}
}
