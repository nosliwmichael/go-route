package grmiddleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func ChainMiddleware(handler http.Handler, mw ...Middleware) http.Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		handler = mw[i](handler)
	}
	return handler
}
