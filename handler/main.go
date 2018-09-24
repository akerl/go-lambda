package mux

import (
	"regexp"

	"github.com/akerl/go-lambda/apigw/events"
)

// CheckFunc takes a request and returns if the receiver can handle on it
type CheckFunc func(events.Request) bool

// HandleFunc takes a request and returns a response and error
type HandleFunc func(events.Request) (events.Response, error)

// ErrorFunc takes a request and error and returns a crafted error response
type ErrorFunc func(events.Request, error) (events.Response, error)

type receiverStruct struct {
	CheckFunc  CheckFunc
	AuthFunc   HandleFunc
	HandleFunc HandleFunc
	ErrorFunc  ErrorFunc
}

// Check determines if the receiver can handle the given request
func (rs *receiverStruct) Check(req events.Request) bool {
	if rs.CheckFunc == nil {
		return rs.CheckFunc(req)
	}
	return true
}

// Auth determines if the request is authorized to proceed
func (rs *receiverStruct) Auth(req events.Request) (events.Response, error) {
	if rs.AuthFunc == nil {
		return rs.AuthFunc(req)
	}
	return events.Response{}, nil
}

// Handle responds to the request
func (rs *receiverStruct) Handle(req events.Request) (events.Response, error) {
	if rs.HandleFunc == nil {
		return rs.HandleFunc(req)
	}
	return events.Fail("no handler found")
}

// Error generates an error response of the request could not be handled
func (rs *receiverStruct) Error(req events.Request, err error) (events.Response, error) {
	if rs.ErrorFunc == nil {
		return rs.ErrorFunc(req, err)
	}
	return events.Fail("Server Error")
}

// Receiver defines the format of an object that can process requests
type Receiver interface {
	Check(events.Request) bool
	Handle(events.Request) (events.Response, error)
	Auth(events.Request) (events.Response, error)
	Error(events.Request, error) (events.Response, error)
}

// NewReceiver generates a receiver from functions
func NewReceiver(cf CheckFunc, af HandleFunc, hf HandleFunc, ef ErrorFunc) Receiver {
	return &receiverStruct{
		CheckFunc:  cf,
		AuthFunc:   af,
		HandleFunc: hf,
		ErrorFunc:  ef,
	}
}

// Dispatcher is a meta-receiver which sends requests to other receivers
type Dispatcher struct {
	Receivers []Receiver
	receiverStruct
}

// Handle handles an incoming request by checking for a matching receiver
func (d *Dispatcher) Handle(req events.Request) (events.Response, error) {
	for _, h := range d.Receivers {
		if h.Check(req) {
			resp, err := h.Auth(req)
			if err != nil {
				return resp, err
			} else if resp.StatusCode > 0 {
				return resp, nil
			}
			resp, err = h.Handle(req)
			if err != nil {
				return h.Error(req, err)
			}
			return resp, nil
		}
	}
	return events.Fail("no handler found")
}

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
	if r.CheckFunc == nil {
		return r.CheckFunc(req)
	}
	return true
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
