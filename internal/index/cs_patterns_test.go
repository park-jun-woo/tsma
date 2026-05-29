package index

import "testing"

func TestCsNamespacePattern(t *testing.T) {
	cases := map[string]string{
		"namespace Com.Example.App {": "Com.Example.App",
		"namespace Com.Example.App;":  "Com.Example.App",
		"namespace Foo":               "Foo",
	}
	for line, want := range cases {
		m := csNamespacePattern.FindStringSubmatch(line)
		if m == nil {
			t.Errorf("no namespace match for %q", line)
			continue
		}
		if m[1] != want {
			t.Errorf("namespace(%q) = %q, want %q", line, m[1], want)
		}
	}
}

func TestCsTypePattern(t *testing.T) {
	cases := map[string]string{
		"public class Foo {":                 "Foo",
		"internal sealed class Bar : Base {": "Bar",
		"public record R(int X) {":           "R",
		"interface IThing {":                 "IThing",
		"public struct Point {":              "Point",
		"enum Color {":                       "Color",
	}
	for line, want := range cases {
		m := csTypePattern.FindStringSubmatch(line)
		if m == nil {
			t.Errorf("no type match for %q", line)
			continue
		}
		if m[1] != want {
			t.Errorf("type(%q) = %q, want %q", line, m[1], want)
		}
	}
}

func TestCsMethodPattern(t *testing.T) {
	cases := map[string]string{
		"public int Add(int a, int b) {":                "Add",
		"private static string Classify(int n) {":       "Classify",
		"public List<int> Items() {":                    "Items",
		"public T Get<T>(string key) where T : new() {": "Get",
	}
	for line, want := range cases {
		m := csMethodPattern.FindStringSubmatch(line)
		if m == nil {
			t.Errorf("no method match for %q", line)
			continue
		}
		if m[1] != want {
			t.Errorf("method(%q) = %q, want %q", line, m[1], want)
		}
	}
}

func TestCsConstructorPattern(t *testing.T) {
	cases := map[string]string{
		"public Foo(int x) {":       "Foo",
		"Bar() {":                   "Bar",
		"internal Baz() : base() {": "Baz",
	}
	for line, want := range cases {
		m := csConstructorPattern.FindStringSubmatch(line)
		if m == nil {
			t.Errorf("no constructor match for %q", line)
			continue
		}
		if m[1] != want {
			t.Errorf("constructor(%q) = %q, want %q", line, m[1], want)
		}
	}
}
