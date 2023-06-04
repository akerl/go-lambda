package mux

import (
	"crypto/subtle"
	"regexp"

	"github.com/akerl/go-lambda/apigw/events"
)

// Route is a receiver that works based on a regex path
type Route struct {
	Path *regexp.Regexp
	SimpleReceiver
}

// Check tests if the path matches for the route
func (r *Route) Check(req events.Request) bool {
	match := r.Path.FindStringSubmatch(req.Path)
	if len(match) == 0 {
		return false
	}
	return r.SimpleReceiver.Check(req)
}

// Handle runs the handle func with path regexp injected
func (r *Route) Handle(req events.Request) (events.Response, error) {
	match := r.Path.FindStringSubmatch(req.Path)
	for i, name := range r.Path.SubexpNames() {
		if name != "" {
			req.PathParameters[name] = match[i]
		}
	}
	return r.SimpleReceiver.Handle(req)
}

// NewRoute is a helper to convert a regexp and handlefunc into a Route Receiver
func NewRoute(path *regexp.Regexp, handler HandleFunc) *Route {
	return &Route{Path: path, SimpleReceiver: SimpleReceiver{HandleFunc: handler}}
}

// NewRouteWithAuth is a helper to conver a regexp, handler, and auth func into a Route Receiver
func NewRouteWithAuth(path *regexp.Regexp, handler HandleFunc, auth HandleFunc) *Route {
	return &Route{
		Path:           path,
		SimpleReceiver: SimpleReceiver{HandleFunc: handler, AuthFunc: auth},
	}
}

// NewRouteWithBasicAuth is a helper to create a route protected by HTTP basic auth
func NewRouteWithBasicAuth(path *regexp.Regexp, handler HandleFunc, users map[string]string) *Route {
	return NewRouteWithAuth(path, handler, basicAuthFunc(users))
}

func basicAuthFunc(users map[string]string) HandleFunc {
	return func(req events.Request) (events.Response, error) {
		user, pass, ok := parseUserName(req.Headers["Authorization"])
		if !ok || c.Users[user] == "" || subtle.ConstantTimeCompare([]byte(c.Users[user]), []byte(pass)) != 1 {
			return events.Response{
				StatusCode: 401,
				Body:       "Unauthorized",
				Headers: map[string]string{
					"WWW-Authenticate": "Basic realm=\"Please authenticate\"",
				},
			}, nil
		}
		return events.Response{}, nil
	}
}
