package filters

import (
	"github.com/4thel00z/libhttp"
	"github.com/4thel00z/libservice/v1"
)

func Auth(app v1.App) libhttp.Filter {
	return func(req libhttp.Request, svc libhttp.Service) libhttp.Response {
		pattern := app.Router.Pattern(req)
		routes := app.Routes()
		route, ok := routes[pattern]
		if !ok {
			return svc(req)
		}

		if route.TokenValidator == nil {
			return svc(req)
		}

		validator := route.TokenValidator

		if validator == nil {
			return svc(req)
		}

		return validator(req, svc)

	}
}
