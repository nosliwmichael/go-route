# go-route
A small library designed to build complex routing behavior in Go. It makes use of Go's standard library by using ServeMux to facilitate route matching. It does introduce new behaviors or functionality but rather allows the user to build out routes in a more structured way.

# Examples

## Single layer routing
Most of the time, your routing needs are quite simple and all you need is a way to define a few paths off of your context root.\
You could add them one by one:
```golang
mux := goroute.NewMux(contextRoot)
mux.AddRoute(goroute.Route{RequestMethod: http.MethodGet, Path: "/user", HandlerFunc: handler.GetUser})
mux.AddRoute(goroute.Route{RequestMethod: http.MethodGet, Path: "/account", HandlerFunc: handler.GetAccount})
```
Or you could add them all at once:
```golang
mux := goroute.NewMux(contextRoot)
mux.AddRoutes([]goroute.Route{
    {RequestMethod: http.MethodGet, Path: "/user", HandlerFunc: handler.GetUser},
    {RequestMethod: http.MethodGet, Path: "/account", HandlerFunc: handler.GetAccount},
})
```
## Multi-layer (Nested) routing
As your application grows in complexity, a need may arise for creating a hierarchy of routes. This can be a bit cumbersome with the standard library as it requires either using nested ServeMux instances or lots of redundant string building. Both of those approaches are used within the goroute package but they're conveniently taken care behind the scenes.
```golang
mux := goroute.NewMux(contextRoot)
mux.AddSubRoutes("/v1/", []goroute.Route{
    {RequestMethod: http.MethodGet, Path: "/user", HandlerFunc: handler.GetUser},
    {RequestMethod: http.MethodGet, Path: "/account", HandlerFunc: handler.GetAccount},
})
mux.AddSubRoutes("/v2/", []goroute.Route{
    {RequestMethod: http.MethodGet, Path: "/user", HandlerFunc: handler.GetUserV2},
    {RequestMethod: http.MethodGet, Path: "/account", HandlerFunc: handler.GetAccountV2},
})
```
Or you could add them all at once:
```golang
mux := goroute.NewServeMux("")
mux.AddRoute(
    goroute.Route{Path: contextRoot, SubRoutes: []goroute.Route{
        {Path: "/v1/", SubRoutes: []goroute.Route{
            {RequestMethod: http.MethodGet, Path: "/user", HandlerFunc: handler.GetUser},
            {RequestMethod: http.MethodGet, Path: "/account", HandlerFunc: handler.GetAccount},
        }},
        {Path: "/v2/", SubRoutes: []goroute.Route{
            {RequestMethod: http.MethodGet, Path: "/user", HandlerFunc: handler.GetUserV2},
            {RequestMethod: http.MethodGet, Path: "/account", HandlerFunc: handler.GetAccountV2},
        }},
    },
})
```
## Middleware
If you're not familiar with middleware, they can be thought of as request interceptors. They are a chain of request handlers that do a little bit of processing before handing the request off to the next handler. There are several ways to add them to your routes.\
You can add them directly on a route. In which case it will be applied only to that route.
```golang
mux := goroute.NewMux(contextRoot)
mux.AddRoute(goroute.Route{RequestMethod: http.MethodGet, Path: "/user", HandlerFunc: handler.GetUser, Middleware: []grmiddleware.Middleware{grmiddleware.ContentTypeJSON()}})
mux.AddRoute(goroute.Route{RequestMethod: http.MethodGet, Path: "/account", HandlerFunc: handler.GetAccount})
```
Or you can apply it to multiple routes being added at once:
```golang
mux := goroute.NewMux(contextRoot)
mux.AddRoutes([]goroute.Route{
    {RequestMethod: http.MethodGet, Path: "/user", HandlerFunc: handler.GetUser},
    {RequestMethod: http.MethodGet, Path: "/account", HandlerFunc: handler.GetAccount},
}, grmiddleware.ContentTypeJSON())
```
Or you can apply it to all routes belonging to a mux: (Note: This means the middleware is applied BEFORE the mux has done any routing)
```golang
mux := goroute.NewMux(contextRoot)
mux.AddRoute(goroute.Route{RequestMethod: http.MethodGet, Path: "/user", HandlerFunc: handler.GetUser})
mux.AddRoute(goroute.Route{RequestMethod: http.MethodGet, Path: "/account", HandlerFunc: handler.GetAccount})
mux.AddRootMiddleware(grmiddleware.ContentTypeJSON())
```
And here's what it would look like to apply middleware at all the different levels:
```golang
rootMux := goroute.NewServeMux(contextRoot)
v1Mux := rootMux.AddSubRoutes("/v1/", []goroute.Route{
    {RequestMethod: http.MethodGet, Path: "/user", HandlerFunc: handler.GetUser},
    {RequestMethod: http.MethodGet, Path: "/account", HandlerFunc: handler.GetAccount},
}, grmiddleware.ContentTypeJSON())
v2Mux := rootMux.AddSubRoutes("/v1/", []goroute.Route{
    {RequestMethod: http.MethodGet, Path: "/user", HandlerFunc: handler.GetUser, Middleware: []grmiddleware.Middleware{grmiddleware.NewRequestMethodCheck(http.MethodGet)}},
    {RequestMethod: http.MethodGet, Path: "/account", HandlerFunc: handler.GetAccount},
})
v2Mux.AddRootMiddleware(middleware.ContentTypeJSON())
```