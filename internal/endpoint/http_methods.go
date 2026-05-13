package endpoint

// httpMethods lists HTTP method names used in Gin's router API.
var httpMethods = map[string]bool{
	"GET":     true,
	"POST":    true,
	"PUT":     true,
	"PATCH":   true,
	"DELETE":  true,
	"HEAD":    true,
	"OPTIONS": true,
	"Any":     true,
	"Handle":  true,
}
