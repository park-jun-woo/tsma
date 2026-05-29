package index

import "testing"

func TestBuildJavaQualifiedName(t *testing.T) {
	cases := []struct {
		pkg    string
		scopes []javaScope
		name   string
		want   string
	}{
		{"com.example", []javaScope{{typeName: "Calculator"}}, "add", "com.example.Calculator.add"},
		{"com.example", []javaScope{{typeName: "Outer"}, {typeName: "Inner"}}, "run", "com.example.Outer.Inner.run"},
		{"", []javaScope{{typeName: "Foo"}}, "bar", "Foo.bar"},
		{"p", nil, "topLevel", "p.topLevel"},
		{"", nil, "loose", "loose"},
	}
	for _, c := range cases {
		if got := buildJavaQualifiedName(c.pkg, c.scopes, c.name); got != c.want {
			t.Errorf("buildJavaQualifiedName(%q,%v,%q) = %q, want %q", c.pkg, c.scopes, c.name, got, c.want)
		}
	}
}
