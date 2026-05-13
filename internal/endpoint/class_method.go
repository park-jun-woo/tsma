//ff:type feature=endpoint type=model
//ff:what Holds info about an HTTP method in a Django class-based view
package endpoint

// classMethod holds info about an HTTP method in a class-based view.
type classMethod struct {
	name      string
	startLine int
	endLine   int
}
