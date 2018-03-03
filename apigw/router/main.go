package router

import (
	"regexp"

	"github.com/akerl/go-lambda/apigw/events"

	"github.com/aws/aws-lambda-go/lambda"
)

// Route defines a path for handling requests
type Route struct {
	Path    *regexp.Regexp
	Handler events.Handler
}

// NewRoute converts a string and handler into a Route
func NewRoute(pathString string, handler events.Handler) (Route, error) {
	path, err := regexp.Compile(pathString)
	if err != nil {
		return Route{}, err
	}
	return Route{Path: path, Handler: handler}, nil
}

// Handle proxies to the route handler
func (r *Route) Handle(req events.Request) (events.Response, error) {
	return r.Handler(req)
}

// Router defines a path-based request handler
type Router struct {
	Routes []Route
}

// Handle handles an incoming request
func (r *Router) Handle(req events.Request) (events.Response, error) {
	for _, route := range r.Routes {
		if route.Path.MatchString(req.Path) {
			return route.Handle(req)
		}
	}
	return events.Fail("no handler found")
}

// Start runs the API GW Lambda
func (r *Router) Start() {
	lambda.Start(r.Handle)
}
