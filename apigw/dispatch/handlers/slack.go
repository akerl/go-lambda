package handlers

import (
	"encoding/json"

	"github.com/akerl/go-lambda/apigw/events"

	"github.com/nlopes/slack"
)

// SlackHandler checks for Slack requests and returns Slack messages
type SlackHandler struct {
	Func func(events.Request) (*slack.Msg, error)
}

// Check validates the Slack body parameter exists
func (h *SlackHandler) Check(req events.Request) bool {
	bodyParams, _ := req.BodyAsParams()
	if bodyParams["trigger_id"] == "" {
		return false
	}
	return true
}

// Run checks the auth token and processes the message
func (h *SlackHandler) Run(req events.Request) (events.Response, error) {
	bodyParams, _ := req.BodyAsParams()
	actualToken := bodyParams["token"]

	params := events.Params{Request: &req}
	expectedToken := params.Lookup("slack_token")

	if expectedToken == "" {
		return events.Fail("no slack_token provided")
	} else if expectedToken != "skip" && expectedToken != actualToken {
		return events.Fail("invalid slack token")
	}

	resp, err := h.Func(req)

	jsonMsg, err := json.Marshal(resp)
	if err != nil {
		return events.Fail("failed to serialize response")
	}
	return events.Succeed(string(jsonMsg))
}
