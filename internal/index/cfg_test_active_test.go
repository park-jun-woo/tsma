package index

import "testing"

func TestCfgTestActive(t *testing.T) {
	if !cfgTestActive(nil, true) {
		t.Error("pending #[cfg(test)] should be active")
	}
	if cfgTestActive(nil, false) {
		t.Error("no scopes and not pending should be inactive")
	}
	if !cfgTestActive([]rsScope{{module: "tests", cfgTest: true}}, false) {
		t.Error("enclosing cfgTest scope should be active")
	}
	if cfgTestActive([]rsScope{{receiver: "Foo"}}, false) {
		t.Error("plain impl scope should be inactive")
	}
}
