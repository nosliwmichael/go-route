package goroute_test

import (
	"go-route/goroute"
	grmiddleware "go-route/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ServeMuxTestCase struct {
	name       string
	req        *http.Request
	wantStatus int
}

var subroutes = []goroute.Route{
	{RequestMethod: http.MethodGet, Path: "/user", HandlerFunc: okHandler},
	{RequestMethod: http.MethodGet, Path: "/account", HandlerFunc: okHandler},
}

func createRequests(method string, path string) *http.Request {
	req, _ := http.NewRequest(method, path, nil)
	return req
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func buildRoutes() *goroute.Mux {
	rootMux := goroute.NewMux("")
	v1Mux := goroute.NewMux("/v1")
	v2Mux := goroute.NewMux("/v2")

	rootMux.AddRoutes([]goroute.Route{
		{Path: "/v1/", Handler: v1Mux},
		{Path: "/v2/", Handler: v2Mux},
	})

	v1Mux.AddSubRoutes("/api/", subroutes)

	v2Mux.AddRoute(goroute.Route{
		Path:        "/api/",
		HandlerFunc: okHandler,
		SubRoutes:   subroutes,
		Middleware: []grmiddleware.Middleware{
			grmiddleware.NewRequestMethodCheck(http.MethodGet),
		},
	})

	return rootMux
}

func TestMux(t *testing.T) {
	testCases := []ServeMuxTestCase{
		{"GET-v1-api-user", createRequests(http.MethodGet, "/v1/api/user"), http.StatusOK},
		{"GET-v1-api-account", createRequests(http.MethodGet, "/v1/api/account"), http.StatusOK},
		{"POST-v1-api-user", createRequests(http.MethodPost, "/v1/api/user"), http.StatusMethodNotAllowed},
		{"POST-v1-api-account", createRequests(http.MethodPost, "/v1/api/account"), http.StatusMethodNotAllowed},

		{"GET-v2-api", createRequests(http.MethodGet, "/v2/api"), http.StatusOK},
		{"GET-v2-api-user", createRequests(http.MethodGet, "/v2/api/user"), http.StatusOK},
		{"GET-v2-api-account", createRequests(http.MethodGet, "/v2/api/account"), http.StatusOK},
		{"POST-v2-api", createRequests(http.MethodPost, "/v2/api"), http.StatusMethodNotAllowed},
		{"POST-v2-api-user", createRequests(http.MethodPost, "/v2/api/user"), http.StatusMethodNotAllowed},
		{"POST-v2-api-account", createRequests(http.MethodPost, "/v2/api/account"), http.StatusMethodNotAllowed},
	}
	for _, tc := range testCases {
		testMuxRequest(t, buildRoutes(), tc)
	}
}

func testMuxRequest(t *testing.T, mux http.Handler, tc ServeMuxTestCase) {
	t.Run(tc.name, func(t *testing.T) {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, tc.req)
		if tc.wantStatus != w.Code {
			t.Errorf("Unexpected status: want [%d] got [%d]", tc.wantStatus, w.Code)
		}
	})
}
