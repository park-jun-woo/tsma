//ff:func feature=index type=helper control=sequence
//ff:what Defines regex patterns for matching Rust fn, impl, and mod declarations
package index

import "regexp"

var (
	// rsFnPattern matches free function and method declarations:
	//   fn foo(...)            pub fn foo(...)            pub(crate) fn foo(...)
	//   pub async fn foo(...)  const fn foo(...)          unsafe fn foo(...)
	// Generics (<T>) after the name are tolerated via best-effort.
	rsFnPattern = regexp.MustCompile(
		`^(?:pub(?:\([^)]*\))?\s+)?(?:default\s+)?(?:const\s+)?(?:async\s+)?(?:unsafe\s+)?(?:extern\s+(?:"[^"]*"\s+)?)?fn\s+(\w+)`,
	)
	// rsImplPattern matches impl blocks, capturing the implementing type.
	//   impl Foo {                impl<T> Foo<T> {
	//   impl Trait for Foo {      impl<T> Trait<T> for Foo<T> {
	rsImplPattern = regexp.MustCompile(
		`^impl(?:\s*<[^>]*>)?\s+(?:.+\s+for\s+)?([A-Za-z_]\w*)`,
	)
	// rsModPattern matches inline module declarations: mod name {  /  pub mod name {
	rsModPattern = regexp.MustCompile(`^(?:pub(?:\([^)]*\))?\s+)?mod\s+(\w+)\s*\{`)
)
