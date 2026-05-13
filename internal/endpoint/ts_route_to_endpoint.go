//ff:func feature=endpoint type=helper control=sequence
//ff:what Converts a tsRouteRegistration to a model.Endpoint
package endpoint

import (
	"strings"

	"github.com/park-jun-woo/tsma/internal/model"
)

func tsRouteToEndpoint(r tsRouteRegistration) model.Endpoint {
	return model.Endpoint{
		Name:   r.handler,
		Method: strings.ToUpper(r.method),
		Path:   r.path,
		Handler: model.FuncLocation{
			File:      r.file,
			StartLine: r.startLine,
			EndLine:   r.endLine,
		},
		Status: model.StatusTodo,
	}
}
