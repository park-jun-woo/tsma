package index

import "testing"

func TestIsRsImplLine(t *testing.T) {
	if !isRsImplLine("impl Foo {", "impl Foo {") {
		t.Error("expected impl-with-brace to match")
	}
	if !isRsImplLine("impl<T> Trait for Foo<T> {", "impl<T> Trait for Foo<T> {") {
		t.Error("expected generic trait impl to match")
	}
	// impl spanning lines (brace on next line) is not treated as an open here.
	if isRsImplLine("impl Foo", "impl Foo") {
		t.Error("impl without brace should not match")
	}
	if isRsImplLine("fn foo() {", "fn foo() {") {
		t.Error("fn line should not be an impl line")
	}
}
