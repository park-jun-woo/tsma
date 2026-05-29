//ff:func feature=index type=helper control=sequence lang=csharp
//ff:what Defines regex patterns for matching C# namespace, type, and method declarations
package index

import "regexp"

var (
	// csNamespacePattern captures a namespace name from either a block-scoped or
	// file-scoped namespace declaration:
	//   namespace Com.Example.App {        namespace Com.Example.App;
	csNamespacePattern = regexp.MustCompile(`^namespace\s+([\w.]+)\s*[{;]?`)

	// csTypePattern matches class/struct/interface/record/enum declarations and
	// captures the type name. Modifiers (public, internal, sealed, abstract,
	// static, partial, ...) preceding the keyword are tolerated best-effort.
	//   public class Foo {     internal sealed class Bar : Base {
	//   public record R(...) { interface I {    struct S {    enum E {
	csTypePattern = regexp.MustCompile(
		`\b(?:class|struct|interface|record|enum)\s+([A-Za-z_]\w*)`,
	)

	// csMethodPattern matches a method declaration of the form
	//   [modifiers] returnType Name(...) [where ...] [{]
	// capturing the method name. The return-type/modifier prefix is required so
	// bare method-call statements (`Foo();`) are not mistaken for declarations.
	// The trailing brace is optional to support C#'s common Allman style where
	// the body brace is placed on the next line; the line must not end in a
	// statement terminator. Generics on the return type (<T>) and array brackets
	// are tolerated best-effort.
	csMethodPattern = regexp.MustCompile(
		`^[A-Za-z_<>\[\],.\s@]+\s+([A-Za-z_]\w*)\s*(?:<[\w\s,]+>)?\s*\([^;{]*\)\s*(?:where\s+[\w.,:\s<>()]+)?(?:\{|$)`,
	)

	// csConstructorPattern matches a constructor declaration: an uppercase-led
	// name immediately followed by a parameter list, with optional access
	// modifiers and an optional base/this initialiser. The trailing brace is
	// optional (Allman style); otherwise the line must end after the signature.
	//   public Foo(...) {   Foo(...)   internal Foo() : base() {
	csConstructorPattern = regexp.MustCompile(
		`^(?:(?:public|protected|private|internal|static)\s+)*([A-Z]\w*)\s*\([^;{]*\)\s*(?::\s*(?:base|this)\s*\([^;{]*\)\s*)?(?:\{|$)`,
	)
)
