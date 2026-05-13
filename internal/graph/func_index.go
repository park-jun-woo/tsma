//ff:type feature=graph type=model
//ff:what Maps function names and qualified names to indices for fast lookup
package graph

// funcIndex provides fast lookup for functions by name and qualified name.
type funcIndex struct {
	byQualified map[string]int   // qualified_name -> index
	byName      map[string][]int // name -> []index
}
