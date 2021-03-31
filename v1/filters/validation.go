package filters

import (
	"context"
	"fmt"
	"github.com/4thel00z/libhttp"

	"github.com/4thel00z/libservice/v1"
)

const (
	ValidationResult = "validation_result"
)

func Validation(app v1.App) libhttp.Filter {
	return func(req libhttp.Request, svc libhttp.Service) libhttp.Response {
		pattern := app.Router.Pattern(req)
		routes := app.Routes()
		route, ok := routes[pattern]
		if !ok {
			return svc(req)
		}

		if route.Validator == nil {
			return svc(req)
		}

		val, err := (*route.Validator)(req)

		if err != nil {
			msg := err.Error()
			return req.Response(v1.GenericResponse{
				Message: fmt.Sprintf("[%s] %s validation error", pattern, route.Method),
				Error:   &msg,
			})
		}

		req.Context = context.WithValue(req.Context, ValidationResult, val)
		return svc(req)

	}
}
