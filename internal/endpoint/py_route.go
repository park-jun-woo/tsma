//ff:type feature=endpoint type=model
//ff:what Holds parsed route info from FastAPI decorators
package endpoint

// pyRoute holds parsed route info from FastAPI decorators.
type pyRoute struct {
	method    string
	path      string
	handler   string
	file      string
	startLine int
	endLine   int
}
