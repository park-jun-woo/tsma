package endpoint

// echoHTTPMethods lists HTTP method names used in Echo's router API.
var echoHTTPMethods = map[string]bool{
	"GET":     true,
	"POST":    true,
	"PUT":     true,
	"PATCH":   true,
	"DELETE":  true,
	"HEAD":    true,
	"OPTIONS": true,
	"Add":     true,
}
