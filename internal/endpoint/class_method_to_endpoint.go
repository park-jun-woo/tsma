//ff:func feature=endpoint type=helper control=sequence
//ff:what Converts a Django class method to a model.Endpoint
package endpoint

import (
	"strings"

	"github.com/park-jun-woo/tsma/internal/model"
)

func classMethodToEndpoint(cm classMethod, relPath string, route djangoRoute) model.Endpoint {
	return model.Endpoint{
		Name:   route.viewName + "." + cm.name,
		Method: strings.ToUpper(cm.name),
		Path:   route.path,
		Handler: model.FuncLocation{
			File:      relPath,
			StartLine: cm.startLine,
			EndLine:   cm.endLine,
		},
		Status: model.StatusTodo,
	}
}
