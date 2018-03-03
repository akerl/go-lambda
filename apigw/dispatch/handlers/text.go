package handlers

import (
	"github.com/akerl/go-lambda/apigw/events"
)

// TextHandler does basic text responses
type TextHandler struct {
	Func func(events.Request) (string, error)
}

// Check is always true for text handlers
func (h *TextHandler) Check(req events.Request) bool {
	return true
}

// Run calls the func with the provided request
func (h *TextHandler) Run(req events.Request) (events.Response, error) {
	resp, err := h.Func(req)
	if err != nil {
		return events.Fail(err.Error())
	}
	return events.Succeed(resp)
}
