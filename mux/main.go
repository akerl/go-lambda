package mux

import (
	"github.com/akerl/go-lambda/apigw/events"

	"github.com/aws/aws-lambda-go/lambda"
)

// Receiver defines the format of an object that can process events.Requests
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

// Start launches a lambda handler from the Receiver
func Start(r Receiver) {
	lambda.Start(r.Handle)
}

// CheckFunc takes a events.Request and returns if the receiver can handle on it
type CheckFunc func(events.Request) bool

// HandleFunc takes a events.Request and returns a events.Response and error
type HandleFunc func(events.Request) (events.Response, error)

// ErrorFunc takes a events.Request and error and returns a crafted error events.Response
type ErrorFunc func(events.Request, error) (events.Response, error)

type receiverStruct struct {
	CheckFunc  CheckFunc
	AuthFunc   HandleFunc
	HandleFunc HandleFunc
	ErrorFunc  ErrorFunc
}

// Check determines if the receiver can handle the given events.Request
func (rs *receiverStruct) Check(req events.Request) bool {
	if rs.CheckFunc != nil {
		return rs.CheckFunc(req)
	}
	return true
}

// Auth determines if the events.Request is authorized to proceed
func (rs *receiverStruct) Auth(req events.Request) (events.Response, error) {
	if rs.AuthFunc != nil {
		return rs.AuthFunc(req)
	}
	return events.Response{}, nil
}

// Handle responds to the events.Request
func (rs *receiverStruct) Handle(req events.Request) (events.Response, error) {
	if rs.HandleFunc != nil {
		return rs.HandleFunc(req)
	}
	return events.Fail("no handler found")
}

// Error generates an error events.Response of the events.Request could not be handled
func (rs *receiverStruct) Error(req events.Request, err error) (events.Response, error) {
	if rs.ErrorFunc != nil {
		return rs.ErrorFunc(req, err)
	}
	return events.Fail("Server Error")
}

// NoCheck is an empty CheckFunc
func NoCheck(_ events.Request) bool {
	return true
}

// NoAuth is an empty AuthFunc
func NoAuth(_ events.Request) (events.Response, error) {
	return events.Response{}, nil
}

// NoError is an empty ErrorFunc
func NoError(_ events.Request, err error) (events.Response, error) {
	return events.Response{}, err
}
