package index

import "testing"

func TestMatchJavaMethod(t *testing.T) {
	cases := []struct {
		line     string
		wantName string
		wantOK   bool
	}{
		{"public int add(int a, int b) {", "add", true},
		{"private static void main(String[] args) {", "main", true},
		{"public List<String> names() {", "names", true},
		{"int helper(int x) throws IOException {", "helper", true},
		{"public Foo(int n) {", "Foo", true},         // constructor
		{"Foo() {", "Foo", true},                     // package-private constructor
		{"if (x > 0) {", "", false},                  // control statement
		{"for (int i = 0; i < n; i++) {", "", false}, // control statement
		{"} else {", "", false},
		{"public class Foo {", "", false},     // type declaration
		{"interface Bar {", "", false},        // type declaration
		{"foo();", "", false},                 // bare call, no body open
		{"return compute(x);", "", false},     // statement
		{"public int field = compute();", "", false},
	}
	for _, c := range cases {
		name, ok := matchJavaMethod(c.line)
		if ok != c.wantOK || name != c.wantName {
			t.Errorf("matchJavaMethod(%q) = (%q,%v), want (%q,%v)", c.line, name, ok, c.wantName, c.wantOK)
		}
	}
}
