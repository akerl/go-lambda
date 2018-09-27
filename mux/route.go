package mux

import (
	"regexp"

	"github.com/akerl/go-lambda/apigw/events"
)

// Route is a receiver that works based on a regex path
type Route struct {
	Path *regexp.Regexp
	receiverStruct
}

// Check tests if the path matches for the route
func (r *Route) Check(req events.Request) bool {
	match := r.Path.FindStringSubmatch(req.Path)
	if len(match) == 0 {
		return false
	}
	return r.receiverStruct.Check(req)
}

// Handle runs the handle func with path regexp injected
func (r *Route) Handle(req events.Request) (events.Response, error) {
	match := r.Path.FindStringSubmatch(req.Path)
	for i, name := range r.Path.SubexpNames() {
		if name != "" {
			req.PathParameters[name] = match[i]
		}
	}
	return r.receiverStruct.Handle(req)
}
