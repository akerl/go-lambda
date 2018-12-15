package slack

import (
	"encoding/json"

	"github.com/akerl/go-lambda/apigw/events"
	"github.com/nlopes/slack"
)

// Handler checks for Slack requests and returns Slack messages
type Handler struct {
	HandleFunc    func(events.Request) (*slack.Msg, error)
	ErrorFunc     func(events.Request, error) (*slack.Msg, error)
	SigningTokens []string
}

// Check validates the Slack body parameter exists
func (h *Handler) Check(req events.Request) bool {
	bodyParams, err := req.BodyAsParams()
	return err == nil && bodyParams["trigger_id"] != ""
}

// Auth checks if the auth token is valid
func (h *Handler) Auth(req events.Request) (events.Response, error) {
	if len(h.SigningTokens) == 0 {
		return events.Reject("no signing tokens provided")
	}

	byteBody := []byte(req.Body)

	for _, i := range h.SigningTokens {
		sv, err := slack.NewSecretsVerifier(req.MultiValueHeaders, i)
		if err != nil {
			return events.Reject("failed to create secret verifier")
		}
		if _, err := sv.Write(byteBody); err != nil {
			return events.Reject("failed to parse body")
		}
		if err := sv.Ensure(); err == nil {
			return events.Response{}, nil
		}
	}

	return events.Reject("invalid signature")
}

// Handle processes the message
func (h *Handler) Handle(req events.Request) (events.Response, error) {
	if h.HandleFunc == nil {
		return events.Fail("No handle function defined")
	}
	resp, err := h.HandleFunc(req)
	if err != nil {
		return events.Response{}, err
	}
	return processRequest(resp)
}

func (h *Handler) Error(req events.Request, err error) (events.Response, error) {
	if h.ErrorFunc == nil {
		return events.Fail("No error function defined")
	}
	resp, err := h.ErrorFunc(req, err)
	if err != nil {
		return events.Fail("Error function failed")
	}
	return processRequest(resp)
}

func processRequest(resp *slack.Msg) (events.Response, error) {
	jsonMsg, err := json.Marshal(resp)
	if err != nil {
		return events.Fail("failed to serialize response")
	}
	return events.Succeed(string(jsonMsg))
}
