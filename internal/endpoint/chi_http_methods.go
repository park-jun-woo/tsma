package endpoint

// chiHTTPMethods lists method names used in Chi's router API.
// Chi uses lowercase-initial method names (Get, Post, etc.).
var chiHTTPMethods = map[string]string{
	"Get":        "GET",
	"Post":       "POST",
	"Put":        "PUT",
	"Patch":      "PATCH",
	"Delete":     "DELETE",
	"Head":       "HEAD",
	"Options":    "OPTIONS",
	"Connect":    "CONNECT",
	"Trace":      "TRACE",
	"Handle":     "",
	"HandleFunc": "",
}
