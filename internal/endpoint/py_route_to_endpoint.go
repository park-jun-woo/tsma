//ff:func feature=endpoint type=helper control=sequence
//ff:what Converts a pyRoute to a model.Endpoint
package endpoint

import "github.com/park-jun-woo/tsma/internal/model"

func pyRouteToEndpoint(r pyRoute) model.Endpoint {
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
