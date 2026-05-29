//ff:func feature=index type=helper control=sequence
//ff:what Defines regex patterns for matching Java package, class, and method declarations
package index

import "regexp"

var (
	// javaPackagePattern captures the package name from a package declaration:
	//   package com.example.app;
	javaPackagePattern = regexp.MustCompile(`^package\s+([\w.]+)\s*;`)

	// javaTypePattern matches class/interface/enum/record declarations and
	// captures the type name. Modifiers (public, abstract, final, static, ...)
	// preceding the keyword are tolerated best-effort.
	//   public class Foo {        final class Bar implements X {
	//   interface I {             enum E {            public record R(...) {
	javaTypePattern = regexp.MustCompile(
		`\b(?:class|interface|enum|record)\s+([A-Za-z_]\w*)`,
	)

	// javaMethodPattern matches a method declaration of the form
	//   [modifiers] returnType name(...) {
	// capturing the method name. The return-type/modifier prefix is required so
	// bare method-call statements (`foo();`) are not mistaken for declarations.
	// Generics on the return type (<T>) and array brackets are tolerated
	// best-effort.
	javaMethodPattern = regexp.MustCompile(
		`^[A-Za-z_<>\[\],.\s@]+\s+([A-Za-z_]\w*)\s*\([^;]*\)\s*(?:throws\s+[\w.,\s]+)?\{`,
	)

	// javaConstructorPattern matches a constructor declaration: an uppercase-led
	// name immediately followed by a parameter list and an opening brace, with
	// optional access modifiers.
	//   public Foo(...) {        Foo(...) {
	javaConstructorPattern = regexp.MustCompile(
		`^(?:(?:public|protected|private)\s+)?([A-Z]\w*)\s*\([^;]*\)\s*(?:throws\s+[\w.,\s]+)?\{`,
	)
)
