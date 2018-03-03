package dispatch

import (
	"github.com/akerl/go-lambda/apigw/events"

	"github.com/aws/aws-lambda-go/lambda"
)

// Receiver describes a wrapper for a user-provided function
type Receiver interface {
	Check(events.Request) bool
	Handle(events.Request) (events.Response, error)
}

// Dispatcher defines a dynamic handler
type Dispatcher struct {
	Receivers []Receiver
}

// Handle handles an incoming request
func (d *Dispatcher) Handle(req events.Request) (events.Response, error) {
	for _, h := range d.Receivers {
		if h.Check(req) {
			return h.Handle(req)
		}
	}
	return events.Fail("no handler found")
}

// Start runs the API GW Lambda
func (d *Dispatcher) Start() {
	lambda.Start(d.Handle)
}
