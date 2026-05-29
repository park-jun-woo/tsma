package index

import "testing"

func TestIsJavaControlKeyword(t *testing.T) {
	keywords := []string{"if", "for", "while", "switch", "catch", "synchronized", "try", "else", "new", "case"}
	for _, k := range keywords {
		if !isJavaControlKeyword(k) {
			t.Errorf("isJavaControlKeyword(%q) = false, want true", k)
		}
	}
	notKeywords := []string{"compute", "main", "Foo", "getValue"}
	for _, k := range notKeywords {
		if isJavaControlKeyword(k) {
			t.Errorf("isJavaControlKeyword(%q) = true, want false", k)
		}
	}
}
