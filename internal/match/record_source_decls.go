//ff:func feature=match type=helper control=iteration dimension=1 lang=go
//ff:what Records every func/method declaration of a parsed file into the receiver map
package match

import "go/ast"

// recordSourceDecls walks the top-level declarations of a parsed source file and
// records each function/method declaration into r: the declared name mapped to
// its distinguisher (the bare receiver type for a method, or "" for a free
// function). Non-func declarations are ignored. It is the per-file inner step
// of BuildPkgSourceReceivers, extracted so the builder's directory loop stays
// shallow.
func recordSourceDecls(r *PkgSourceReceivers, f *ast.File) {
	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		addNameReceiver(r, fd.Name.Name, funcDeclReceiver(fd))
	}
}
