package mux

import (
	"github.com/akerl/go-lambda/apigw/events"
)

// SimpleReceiver defines a basic Receiver composed of 4 functions
type SimpleReceiver struct {
	CheckFunc  CheckFunc
	AuthFunc   HandleFunc
	HandleFunc HandleFunc
	ErrorFunc  ErrorFunc
}

// Check determines if the receiver can handle the given events.Request
func (rs *SimpleReceiver) Check(req events.Request) bool {
	if rs.CheckFunc == nil {
		rs.CheckFunc = NoCheck
	}
	return rs.CheckFunc(req)
}

// Auth determines if the events.Request is authorized to proceed
func (rs *SimpleReceiver) Auth(req events.Request) (events.Response, error) {
	if rs.AuthFunc == nil {
		rs.AuthFunc = NoAuth
	}
	return rs.AuthFunc(req)
}

// Handle responds to the events.Request
func (rs *SimpleReceiver) Handle(req events.Request) (events.Response, error) {
	if rs.HandleFunc == nil {
		return events.Fail("no handler found")
	}
	return rs.HandleFunc(req)
}

// Error generates an error events.Response of the events.Request could not be handled
func (rs *SimpleReceiver) Error(req events.Request, err error) (events.Response, error) {
	if rs.ErrorFunc == nil {
		rs.ErrorFunc = NoError
	}
	return rs.ErrorFunc(req, err)
}
