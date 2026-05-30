package match

import "testing"

func TestBuildPkgSourceReceiversSameNameMultiple(t *testing.T) {
	root, pkgDir := writePkg(t, map[string]string{
		"go_file.go": `package gen
type GoFile struct{}
func (f *GoFile) GetFuncs() int { return 1 }
`,
		"cs_file.go": `package gen
type CSharpFile struct{}
func (f *CSharpFile) GetFuncs() int { return 2 }
`,
	})

	r, err := BuildPkgSourceReceivers(root, pkgDir)
	if err != nil {
		t.Fatal(err)
	}
	if !r.isSameNameMultiple("GetFuncs") {
		t.Errorf("GetFuncs should be same-name-multiple (GoFile + CSharpFile)")
	}
}

func TestBuildPkgSourceReceiversSameNameSingle(t *testing.T) {
	root, pkgDir := writePkg(t, map[string]string{
		"go_file.go": `package gen
type GoFile struct{}
func (f *GoFile) GetFuncs() int { return 1 }
func (f *GoFile) GetPath() string { return "" }
`,
	})

	r, err := BuildPkgSourceReceivers(root, pkgDir)
	if err != nil {
		t.Fatal(err)
	}
	if r.isSameNameMultiple("GetFuncs") {
		t.Errorf("GetFuncs declared on one receiver only -> single")
	}
	// Pointer and value receivers normalize to the same type name.
	if r.isSameNameMultiple("GetPath") {
		t.Errorf("GetPath single")
	}
}

func TestBuildPkgSourceReceiversFreeFuncMethodCoexist(t *testing.T) {
	// A free function and a method sharing a name -> set {"", "GoFile"} -> multiple.
	root, pkgDir := writePkg(t, map[string]string{
		"x.go": `package gen
type GoFile struct{}
func (f *GoFile) Parse() int { return 1 }
func Parse() int { return 2 }
`,
	})

	r, err := BuildPkgSourceReceivers(root, pkgDir)
	if err != nil {
		t.Fatal(err)
	}
	if !r.isSameNameMultiple("Parse") {
		t.Errorf("Parse as free func + method should be same-name-multiple")
	}
}

func TestBuildPkgSourceReceiversPointerValueNormalized(t *testing.T) {
	// *GoFile and GoFile receivers must collapse to one distinguisher.
	root, pkgDir := writePkg(t, map[string]string{
		"x.go": `package gen
type GoFile struct{}
func (f *GoFile) A() {}
func (f GoFile) B() {}
func (f *GoFile) A2() {}
`,
	})

	r, err := BuildPkgSourceReceivers(root, pkgDir)
	if err != nil {
		t.Fatal(err)
	}
	// Same method name on *T and T (synthesize): declare A on both forms.
	root2, pkgDir2 := writePkg(t, map[string]string{
		"x.go": `package gen
type GoFile struct{}
func (f *GoFile) A() {}
`,
		"y.go": `package gen
func (f GoFile) A() {}
`,
	})
	r2, err := BuildPkgSourceReceivers(root2, pkgDir2)
	if err != nil {
		t.Fatal(err)
	}
	if r2.isSameNameMultiple("A") {
		t.Errorf("A on *GoFile and GoFile must normalize to one receiver (single)")
	}
	_ = r
}

func TestBuildPkgSourceReceiversSkipsTestFiles(t *testing.T) {
	root, pkgDir := writePkg(t, map[string]string{
		"x.go": `package gen
type GoFile struct{}
func (f *GoFile) GetFuncs() int { return 1 }
`,
		// A method on a different receiver, but in a _test.go: must be ignored,
		// so GetFuncs stays single.
		"x_test.go": `package gen
type CSharpFile struct{}
func (f *CSharpFile) GetFuncs() int { return 2 }
`,
	})

	r, err := BuildPkgSourceReceivers(root, pkgDir)
	if err != nil {
		t.Fatal(err)
	}
	if r.isSameNameMultiple("GetFuncs") {
		t.Errorf("_test.go declarations must not count toward source receivers")
	}
}

func TestBuildPkgSourceReceiversSkipsUnparseableFile(t *testing.T) {
	// A syntactically broken .go file must be skipped (not abort the build), so
	// the valid file's declarations are still recorded. This is the resilience
	// branch: one bad file does not poison the whole package map.
	root, pkgDir := writePkg(t, map[string]string{
		"good.go": `package gen
type GoFile struct{}
func (f *GoFile) GetFuncs() int { return 1 }
`,
		// Not valid Go: unterminated func body / garbage tokens.
		"broken.go": `package gen
func (f *CSharpFile) GetFuncs( {{{ this is not valid go
`,
	})

	r, err := BuildPkgSourceReceivers(root, pkgDir)
	if err != nil {
		t.Fatalf("unparseable file must not abort the build: %v", err)
	}
	// The good file's method is recorded; the broken file's bogus second
	// receiver for GetFuncs was skipped, so it stays single.
	if _, ok := r.byName["GetFuncs"]; !ok {
		t.Errorf("GetFuncs from the valid file must still be recorded")
	}
	if r.isSameNameMultiple("GetFuncs") {
		t.Errorf("broken file must be skipped, leaving GetFuncs single")
	}
}

func TestBuildPkgSourceReceiversReadDirError(t *testing.T) {
	root := t.TempDir()
	if _, err := BuildPkgSourceReceivers(root, "does/not/exist"); err == nil {
		t.Fatal("expected error for missing dir")
	}
}

func TestIsSameNameMultipleNilSafe(t *testing.T) {
	var r *PkgSourceReceivers
	if r.isSameNameMultiple("anything") {
		t.Errorf("nil receiver map must report single (false)")
	}
}
