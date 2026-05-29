package coverage

import "testing"

func TestFlattenJacocoPackage(t *testing.T) {
	pkg := jacocoPackage{
		Name: "com/example",
		SourceFiles: []jacocoSourceFile{
			{Name: "Foo.java", Lines: []jacocoLine{{Nr: 1, Ci: 1}}},
			{Name: "Bar.java"},
		},
	}
	files := flattenJacocoPackage(pkg)
	if len(files) != 2 {
		t.Fatalf("files = %d, want 2", len(files))
	}
	if files[0].Path != "com/example/Foo.java" || files[1].Path != "com/example/Bar.java" {
		t.Errorf("paths = %q, %q", files[0].Path, files[1].Path)
	}
}

func TestFlattenJacocoPackageNoName(t *testing.T) {
	pkg := jacocoPackage{SourceFiles: []jacocoSourceFile{{Name: "Root.java"}}}
	files := flattenJacocoPackage(pkg)
	if len(files) != 1 || files[0].Path != "Root.java" {
		t.Errorf("got %+v, want Root.java with no prefix", files)
	}
}
