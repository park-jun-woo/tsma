package index

import "testing"

func TestBuildRsQualifiedName(t *testing.T) {
	cases := []struct {
		name   string
		pkgDir string
		scopes []rsScope
		fn     string
		want   string
	}{
		{"top level no dir", "", nil, "main", "main"},
		{"with dir", "src", nil, "run", "src.run"},
		{"impl method", "src", []rsScope{{receiver: "Foo"}}, "bar", "src.Foo.bar"},
		{"module + impl", "src", []rsScope{{module: "util"}, {receiver: "Foo"}}, "bar", "src.util::Foo.bar"},
		{"nested modules", "", []rsScope{{module: "a"}, {module: "b"}}, "f", "a::b.f"},
	}
	for _, c := range cases {
		if got := buildRsQualifiedName(c.pkgDir, c.scopes, c.fn); got != c.want {
			t.Errorf("%s: buildRsQualifiedName = %q, want %q", c.name, got, c.want)
		}
	}
}
