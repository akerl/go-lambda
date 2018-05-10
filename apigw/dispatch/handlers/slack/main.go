package slack

import (
	"encoding/json"

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

	if len(h.SlackTokens) == 0 {
		return events.Response{
			StatusCode: 403,
			Body:       "no slacktokens provided",
		}, nil
	}
	for _, i := range h.SlackTokens {
		if i == "skip" || i == actualToken {
			return events.Response{}, nil
		}
	}

	return events.Response{
		StatusCode: 403,
		Body:       "invalid slack_token",
	}, nil
}
