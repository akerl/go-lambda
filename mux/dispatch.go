package mux

import (
	"github.com/akerl/go-lambda/apigw/events"
)

// Dispatcher is a meta-receiver which sends requests to other receivers
type Dispatcher struct {
	Receivers []Receiver
	SimpleReceiver
}

// Handle handles an incoming request by checking for a matching receiver
func (d *Dispatcher) Handle(req events.Request) (events.Response, error) {
	for _, h := range d.Receivers {
		if h.Check(req) {
			resp, err := h.Auth(req)
			if err != nil {
				return h.Error(req, err)
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

// NewDispatcher is a helper to build a dispatcher from a list of receivers
func NewDispatcher(receivers ...Receiver) *Dispatcher {
	return &Dispatcher{Receivers: receivers}
}
