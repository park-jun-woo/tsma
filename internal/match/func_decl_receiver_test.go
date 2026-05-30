package match

import "testing"

func TestFuncDeclReceiver(t *testing.T) {
	method := parseFuncDecl(t, `package p
func (f *GoFile) M() {}
`, "M")
	if got := funcDeclReceiver(method); got != "GoFile" {
		t.Errorf("method receiver = %q, want GoFile", got)
	}

	free := parseFuncDecl(t, `package p
func Free() {}
`, "Free")
	if got := funcDeclReceiver(free); got != "" {
		t.Errorf("free func receiver = %q, want \"\"", got)
	}
}
