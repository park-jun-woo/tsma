//ff:type feature=index type=model lang=csharp
//ff:what Holds an enclosing C# scope (namespace/class/struct/interface/record/enum) tied to a brace depth
package index

// csScope records an enclosing namespace or type context opened at a given brace
// depth, used to build nested qualified names.
type csScope struct {
	depth    int    // brace depth at which this scope's body was opened
	typeName string // declared namespace or type name
}
