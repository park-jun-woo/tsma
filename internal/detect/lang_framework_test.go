package detect

import "testing"

func TestLangFrameworkStructFields(t *testing.T) {
	lf := LangFramework{Lang: "go"}
	if lf.Lang != "go" {
		t.Errorf("Lang = %q, want %q", lf.Lang, "go")
	}
}

func TestLangFrameworkZeroValue(t *testing.T) {
	var lf LangFramework
	if lf.Lang != "" {
		t.Errorf("zero Lang = %q, want empty", lf.Lang)
	}
}

func TestLangFrameworkDifferentLanguages(t *testing.T) {
	cases := []string{"go", "typescript", "python", ""}
	for _, lang := range cases {
		lf := LangFramework{Lang: lang}
		if lf.Lang != lang {
			t.Errorf("Lang = %q, want %q", lf.Lang, lang)
		}
	}
}
