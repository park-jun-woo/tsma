//ff:func feature=graph type=helper control=sequence
//ff:what Adds a bidirectional edge between caller and callee functions
package graph

import "github.com/park-jun-woo/tsma/internal/model"

// addEdge adds a caller->callee edge and callee->caller reverse edge.
func addEdge(functions []model.Function, callerIdx, calleeIdx int, ambiguous bool) {
	callerQN := functions[callerIdx].QualifiedName
	calleeQN := functions[calleeIdx].QualifiedName

	// Add callee to caller's callees.
	functions[callerIdx].Callees = append(functions[callerIdx].Callees, model.Edge{
		Target:    calleeQN,
		Ambiguous: ambiguous,
	})

	// Add caller to callee's callers.
	functions[calleeIdx].Callers = append(functions[calleeIdx].Callers, model.Edge{
		Target:    callerQN,
		Ambiguous: ambiguous,
	})
}
