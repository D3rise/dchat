package server

import (
	"net/http"

	"go.uber.org/fx"
)

type Route interface {
	Pattern() string
	http.Handler
}

// AsRoute annotates the given constructor to state that
// it provides a route to the "routes" group.
func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Route)),
		fx.ResultTags(`group:"routes"`),
	)
}
