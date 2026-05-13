//ff:type feature=endpoint type=model
//ff:what Holds parsed route info from a TypeScript or JavaScript file
package endpoint

// tsRouteRegistration holds parsed route info from a TypeScript/JavaScript file.
type tsRouteRegistration struct {
	method    string
	path      string
	handler   string
	file      string
	startLine int
	endLine   int
}
