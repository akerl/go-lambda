package slack

import (
	"encoding/json"

	"github.com/akerl/go-lambda/apigw/events"

	"github.com/nlopes/slack"
)

type slackFunc func(events.Request) (*slack.Msg, error)

// Handler checks for Slack requests and returns Slack messages
type Handler struct {
	HandleFunc  slackFunc
	ErrorFunc   slackFunc
	SlackTokens []string
}

// Check validates the Slack body parameter exists
func (h *Handler) Check(req events.Request) bool {
	bodyParams, err := req.BodyAsParams()
	return err == nil && bodyParams["trigger_id"] != ""
}

// Handle processes the message
func (h *Handler) Handle(req events.Request) (events.Response, error) {
	return processRequest(req, h.HandleFunc)
}

// Auth checks if the auth token is valid
func (h *Handler) Auth(req events.Request) (events.Response, error) {
	bodyParams, err := req.BodyAsParams()
	if err != nil {
		return events.Fail("failed to process params")
	}
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

func (h *Handler) Error(req events.Request) (events.Response, error) {
	return processRequest(req, h.ErrorFunc)
}

func processRequest(req events.Request, f slackFunc) (events.Response, error) {
	if f == nil {
		return events.Fail("No handle function defined")
	}
	resp, err := f(req)
	if err != nil {
		return events.Fail("failed to process request")
	}

	jsonMsg, err := json.Marshal(resp)
	if err != nil {
		return events.Fail("failed to serialize response")
	}
	return events.Succeed(string(jsonMsg))
}
