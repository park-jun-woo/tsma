package index

import "testing"

func TestMatchCsMethod(t *testing.T) {
	cases := []struct {
		line     string
		wantName string
		wantOK   bool
	}{
		{"public int Add(int a, int b) {", "Add", true},
		{"private static string Classify(int n) {", "Classify", true},
		{"public Foo(int x) {", "Foo", true},
		// Modifier-less constructor: csMethodPattern needs a return-type prefix,
		// so this falls through to the csConstructorPattern branch (lines 23-26).
		{"Bar() {", "Bar", true},
		{"if (x > 0) {", "", false},
		{"foreach (var i in items) {", "", false},
		{"public class Foo {", "", false},
		{"Console.WriteLine(x);", "", false},
	}
	for _, tc := range cases {
		name, ok := matchCsMethod(tc.line)
		if ok != tc.wantOK || name != tc.wantName {
			t.Errorf("matchCsMethod(%q) = (%q, %v), want (%q, %v)",
				tc.line, name, ok, tc.wantName, tc.wantOK)
		}
	}
}
