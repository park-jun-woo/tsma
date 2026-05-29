package coverage

import "testing"

func TestFlattenJacoco(t *testing.T) {
	report := &jacocoReport{
		Packages: []jacocoPackage{
			{
				Name: "com/example",
				SourceFiles: []jacocoSourceFile{
					{Name: "Foo.java", Lines: []jacocoLine{{Nr: 1, Ci: 2}}},
					{Name: "Bar.java"},
				},
			},
			{
				Name:        "",
				SourceFiles: []jacocoSourceFile{{Name: "Root.java"}},
			},
		},
	}
	cov := flattenJacoco(report)
	if len(cov.Files) != 3 {
		t.Fatalf("files = %d, want 3", len(cov.Files))
	}
	if cov.Files[0].Path != "com/example/Foo.java" {
		t.Errorf("files[0].Path = %q", cov.Files[0].Path)
	}
	if cov.Files[2].Path != "Root.java" {
		t.Errorf("files[2].Path = %q, want Root.java (no package prefix)", cov.Files[2].Path)
	}
}
