package index

import "testing"

func TestJavaPackagePattern(t *testing.T) {
	m := javaPackagePattern.FindStringSubmatch("package com.example.app;")
	if m == nil || m[1] != "com.example.app" {
		t.Errorf("javaPackagePattern failed: %v", m)
	}
	if javaPackagePattern.MatchString("import com.example.X;") {
		t.Error("javaPackagePattern should not match an import")
	}
}

func TestJavaTypePattern(t *testing.T) {
	matches := map[string]string{
		"public class Foo {":               "Foo",
		"final class Bar implements X {":    "Bar",
		"interface I {":                     "I",
		"enum Color {":                      "Color",
		"public record Point(int x, int y){": "Point",
	}
	for line, want := range matches {
		m := javaTypePattern.FindStringSubmatch(line)
		if m == nil || m[1] != want {
			t.Errorf("javaTypePattern(%q) = %v, want %q", line, m, want)
		}
	}
	if javaTypePattern.MatchString("int x = 1;") {
		t.Error("javaTypePattern should not match a field declaration")
	}
}

func TestJavaMethodPattern(t *testing.T) {
	m := javaMethodPattern.FindStringSubmatch("public int add(int a, int b) {")
	if m == nil || m[1] != "add" {
		t.Errorf("javaMethodPattern failed on method: %v", m)
	}
}

func TestJavaConstructorPattern(t *testing.T) {
	m := javaConstructorPattern.FindStringSubmatch("public Widget(int n) {")
	if m == nil || m[1] != "Widget" {
		t.Errorf("javaConstructorPattern failed: %v", m)
	}
}
