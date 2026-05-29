//ff:type feature=index type=model
//ff:what Holds an enclosing Java type scope (class/interface/enum) tied to a brace depth
package index

// javaScope records an enclosing class/interface/enum/record context opened at
// a given brace depth, used to build nested qualified names.
type javaScope struct {
	depth    int    // brace depth at which this type's body was opened
	typeName string // declared type name
}
