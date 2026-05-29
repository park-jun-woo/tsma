package index

import "testing"

func TestRsFnPattern(t *testing.T) {
	matches := map[string]string{
		"fn foo() {":                     "foo",
		"pub fn bar(x: i32) -> i32 {":    "bar",
		"pub(crate) fn baz() {":          "baz",
		"pub async fn fetch() {":         "fetch",
		"const fn c() -> u8 {":           "c",
		"unsafe fn danger() {":           "danger",
		"pub fn generic<T>(t: T) -> T {": "generic",
		`pub extern "C" fn ext() {`:      "ext",
	}
	for line, want := range matches {
		m := rsFnPattern.FindStringSubmatch(line)
		if m == nil {
			t.Errorf("rsFnPattern did not match %q", line)
			continue
		}
		if m[1] != want {
			t.Errorf("rsFnPattern(%q) = %q, want %q", line, m[1], want)
		}
	}

	nonMatches := []string{"let x = 1;", "struct Foo {", "// fn comment"}
	for _, line := range nonMatches {
		if rsFnPattern.MatchString(line) {
			t.Errorf("rsFnPattern should not match %q", line)
		}
	}
}

func TestRsImplPattern(t *testing.T) {
	matches := map[string]string{
		"impl Foo {":                   "Foo",
		"impl<T> Foo<T> {":             "Foo",
		"impl Trait for Foo {":         "Foo",
		"impl<T> Display for Bar<T> {": "Bar",
	}
	for line, want := range matches {
		m := rsImplPattern.FindStringSubmatch(line)
		if m == nil {
			t.Errorf("rsImplPattern did not match %q", line)
			continue
		}
		if m[1] != want {
			t.Errorf("rsImplPattern(%q) = %q, want %q", line, m[1], want)
		}
	}
}

func TestRsModPattern(t *testing.T) {
	m := rsModPattern.FindStringSubmatch("pub mod tests {")
	if m == nil || m[1] != "tests" {
		t.Errorf("rsModPattern failed on inline mod: %v", m)
	}
	if rsModPattern.MatchString("mod external;") {
		t.Error("rsModPattern should not match a file-module declaration")
	}
}
