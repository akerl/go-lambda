package text

import (
	"github.com/akerl/go-lambda/apigw/events"
)

// Handler does basic text responses
type Handler struct {
	Func     func(events.Request) (string, error)
	AuthFunc func(events.Request) (events.Response, error)
}

// Check is always true for text handlers
func (h *Handler) Check(req events.Request) bool {
	return true
}

// Handle calls the func with the provided request
func (h *Handler) Handle(req events.Request) (events.Response, error) {
	resp, err := h.Func(req)
	if err != nil {
		return events.Fail(err.Error())
	}
	return events.Succeed(resp)
}

// Auth runs an auth check if provided
func (h *Handler) Auth(req events.Request) (events.Response, error) {
	if h.AuthFunc == nil {
		return events.Succeed("")
	}
	return h.AuthFunc(req)
}
