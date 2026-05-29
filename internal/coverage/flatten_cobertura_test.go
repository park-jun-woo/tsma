package coverage

import "testing"

func TestFlattenCoberturaMergesPartialClasses(t *testing.T) {
	report := &coberturaReport{
		Packages: []coberturaPackage{
			{Classes: []coberturaClass{
				{Filename: "App/Foo.cs", Lines: []coberturaLine{{Number: 1, Hits: 1}}},
				{Filename: "App/Foo.cs", Lines: []coberturaLine{{Number: 9, Hits: 0}}},
				{Filename: "App/Bar.cs", Lines: []coberturaLine{{Number: 3, Hits: 2}}},
			}},
		},
	}
	cov := flattenCobertura(report)
	if len(cov.Files) != 2 {
		t.Fatalf("files = %d, want 2 (Foo merged)", len(cov.Files))
	}
	if cov.Files[0].Path != "App/Foo.cs" || len(cov.Files[0].Lines) != 2 {
		t.Errorf("Foo entry = %+v, want 2 lines", cov.Files[0])
	}
}

func TestFlattenCoberturaSkipsBlankFilename(t *testing.T) {
	report := &coberturaReport{
		Packages: []coberturaPackage{
			{Classes: []coberturaClass{{Filename: "", Lines: []coberturaLine{{Number: 1}}}}},
		},
	}
	if cov := flattenCobertura(report); len(cov.Files) != 0 {
		t.Errorf("expected blank-filename class skipped, got %+v", cov.Files)
	}
}
