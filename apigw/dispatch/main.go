package dispatch

import (
	"github.com/akerl/go-lambda/apigw/events"

	"github.com/aws/aws-lambda-go/lambda"
)

// Receiver describes a wrapper for a user-provided function
type Receiver interface {
	Check(events.Request) bool
	Handle(events.Request) (events.Response, error)
	Auth(events.Request) (bool, string)
}

// Dispatcher defines a dynamic handler
type Dispatcher struct {
	Receivers []Receiver
}

// Handle handles an incoming request
func (d *Dispatcher) Handle(req events.Request) (events.Response, error) {
	for _, h := range d.Receivers {
		if h.Check(req) {
			if ok, msg := h.Auth(req); !ok {
				return events.Respond(403, msg)
			}
			return h.Handle(req)
		}
	}
	return events.Fail("no handler found")
}

// Start runs the API GW Lambda
func (d *Dispatcher) Start() {
	lambda.Start(d.Handle)
}

type checkFunc func(events.Request) bool
type handleFunc func(events.Request) (events.Response, error)
type authFunc func(events.Request) (bool, string)

type shim struct {
	CheckFunc  checkFunc
	HandleFunc handleFunc
	AuthFunc   authFunc
}

// Check runs the check func
func (s *shim) Check(req events.Request) bool {
	return s.CheckFunc(req)
}

// Handle runs the handle func
func (s *shim) Handle(req events.Request) (events.Response, error) {
	return s.HandleFunc(req)
}

// Auth runs the auth func
func (s *shim) Auth(req events.Request) (bool, string) {
	return s.AuthFunc(req)
}

// NewReceiver creates a receiver from individual functions
func NewReceiver(cf checkFunc, hf handleFunc, af authFunc) Receiver {
	return &shim{
		CheckFunc:  cf,
		HandleFunc: hf,
		AuthFunc:   af,
	}
}
