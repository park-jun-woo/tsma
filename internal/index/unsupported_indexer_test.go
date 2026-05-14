package index

import "testing"

func TestUnsupportedIndexerReturnsError(t *testing.T) {
	u := &UnsupportedIndexer{Lang: "rust"}
	funcs, err := u.Index("/some/project")
	if err == nil {
		t.Fatal("expected error from UnsupportedIndexer.Index")
	}
	if funcs != nil {
		t.Errorf("expected nil functions, got %v", funcs)
	}
	want := "indexing not implemented for: rust"
	if err.Error() != want {
		t.Errorf("error = %q, want %q", err.Error(), want)
	}
}

func TestUnsupportedIndexerDifferentLangs(t *testing.T) {
	langs := []string{"java", "ruby", "c++", "php"}
	for _, lang := range langs {
		u := &UnsupportedIndexer{Lang: lang}
		_, err := u.Index(".")
		if err == nil {
			t.Errorf("expected error for lang %q", lang)
		}
	}
}
