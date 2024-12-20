package grmiddleware

import (
	grheader "go-route/header"
	"net/http"
)

func ContentTypeJSON() Middleware {
	return NewContentTypeMW(grheader.ApplicationJSON)
}

func NewContentTypeMW(value string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(grheader.ContentTypeHeader, value)
		})
	}
}
