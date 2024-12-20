package grmiddleware

import "net/http"

func NewRequestMethodCheck(requestMethod string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if requestMethod != r.Method {
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			}
			next.ServeHTTP(w, r)
		})
	}
}
