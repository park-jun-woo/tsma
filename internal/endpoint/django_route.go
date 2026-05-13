//ff:type feature=endpoint type=model
//ff:what Holds a parsed URL pattern from Django urls.py
package endpoint

// djangoRoute holds a parsed URL pattern from urls.py.
type djangoRoute struct {
	path     string
	viewName string
	methods  []string
}
