package slack

import (
	"encoding/json"
	"fmt"

	"github.com/akerl/go-lambda/apigw/events"

	"github.com/nlopes/slack"
)

// Handler checks for Slack requests and returns Slack messages
type Handler struct {
	Func        func(events.Request) (*slack.Msg, error)
	SlackTokens []string
}

// Check validates the Slack body parameter exists
func (h *Handler) Check(req events.Request) bool {
	bodyParams, _ := req.BodyAsParams()
	if bodyParams["trigger_id"] == "" {
		return false
	}
	return true
}

// Handle processes the message
func (h *Handler) Handle(req events.Request) (events.Response, error) {
	resp, err := h.Func(req)

	jsonMsg, err := json.Marshal(resp)
	if err != nil {
		return events.Fail("failed to serialize response")
	}
	return events.Succeed(string(jsonMsg))
}

// Auth checks if the auth token is valid
func (h *Handler) Auth(req events.Request) (events.Response, error) {
	bodyParams, _ := req.BodyAsParams()
	actualToken := bodyParams["token"]

	params := events.Params{Request: &req}
	expectedToken := params.Lookup("slack_token")

	if expectedToken == "" {
		return events.Response{
			StatusCode: 403,
			Body:       "no slack_token provided",
		}, nil
	} else if expectedToken == "skip" {
		return events.Response{}, nil
	}

	for _, i := range h.SlackTokens {
		if i == actualToken {
			return events.Response{}, nil
		}
	}

	return events.Response{
		StatusCode: 403,
		Body:       "invalid slack_token",
	}, nil
}
