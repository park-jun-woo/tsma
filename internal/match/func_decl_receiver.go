//ff:func feature=match type=helper control=selection lang=go
//ff:what Returns a func declaration's bare receiver type, or "" for a free function
package match

import "go/ast"

// funcDeclReceiver returns the bare receiver type name of a method declaration
// (pointer/value/generic normalized via srcReceiver), or "" when fd is a free
// function with no receiver. The "" value is also the distinguisher used for
// free functions in PkgSourceReceivers, so a free function and a method sharing
// a name still produce a multi-element set.
func funcDeclReceiver(fd *ast.FuncDecl) string {
	switch {
	case fd.Recv == nil || len(fd.Recv.List) == 0:
		return ""
	default:
		return srcReceiver(fd.Recv.List[0].Type)
	}
}
