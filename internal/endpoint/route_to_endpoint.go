//ff:func feature=endpoint type=helper control=sequence
//ff:what Converts a routeRegistration to a model.Endpoint
package endpoint

import "github.com/park-jun-woo/tsma/internal/model"

func routeToEndpoint(r routeRegistration) model.Endpoint {
	return model.Endpoint{
		Name:   r.handler,
		Method: r.method,
		Path:   r.path,
		Handler: model.FuncLocation{
			File:      r.file,
			StartLine: r.startLine,
			EndLine:   r.endLine,
		},
		Status: model.StatusTodo,
	}
}
