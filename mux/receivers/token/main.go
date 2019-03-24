package slack

import (
	"net/http"

	"github.com/akerl/go-lambda/apigw/events"

	"golang.org/x/crypto/bcrypt"
)

const (
	headerName = "Authorization"
)

// Handler checks for Slack requests and returns Slack messages
type Handler struct {
	HandleFunc   func(events.Request) (events.Response, error)
	ErrorFunc    func(events.Request, error) (events.Response, error)
	BcryptTokens map[string]string
}

// Check validates the Slack body parameter exists
func (h *Handler) Check(req events.Request) bool {
	return req.Headers[headerName] != ""
}

// Auth checks if the auth token is valid
func (h *Handler) Auth(req events.Request) (events.Response, error) {
	if len(h.BcryptTokens) == 0 {
		return events.Reject("no tokens provided")
	}

	httpReq := http.Request{}
	httpReq.Header.Add(headerName, req.Headers[headerName])
	user, password, ok := httpReq.BasicAuth()
	if !ok {
		return events.Reject("invalid token")
	}

	expectedHash, ok := h.BcryptTokens[user]
	if !ok {
		return events.Reject("invalid token")
	}

	err := bcrypt.CompareHashAndPassword([]byte(expectedHash), []byte(password))
	if err == nil {
		return events.Response{}, nil
	}

	return events.Reject("invalid token")
}

// Handle processes the message
func (h *Handler) Handle(req events.Request) (events.Response, error) {
	if h.HandleFunc == nil {
		return events.Fail("No handle function defined")
	}
	return h.HandleFunc(req)
}

// Error processes errors during handler processing
func (h *Handler) Error(req events.Request, e error) (events.Response, error) {
	if h.ErrorFunc == nil {
		return events.Fail("No error function defined")
	}
	return h.ErrorFunc(req, e)
}
