//ff:type feature=endpoint type=model
//ff:what Holds parsed route info shared by Go framework detectors
package endpoint

// routeRegistration holds parsed route info.
type routeRegistration struct {
	method    string
	path      string
	handler   string
	file      string
	startLine int
	endLine   int
}
