package apigw

import (
	"encoding/json"

	"github.com/nlopes/slack"
)

// SlackHandler checks for Slack requests and returns Slack messages
type SlackHandler struct {
	Func       HandlerFunc
	bodyParams map[string]string
}

// Check validates the Slack body parameter exists
func (h *SlackHandler) Check(req Request) bool {
	bodyParams, _ := req.BodyAsParams()
	if bodyParams["trigger_id"] == "" {
		return false
	}
	return true
}

// Run checks the auth token and processes the message
func (h *SlackHandler) Run(req Request, params Params) (Response, error) {
	bodyParams, _ := req.BodyAsParams()
	actualToken := bodyParams["token"]

	expectedToken := params.Lookup("slack_token")

	if expectedToken == "" {
		return Fail("no slack_token provided")
	} else if expectedToken != "skip" && expectedToken != actualToken {
		return Fail("invalid slack token")
	}

	resp, err := h.Func(req, params)

	var msg slack.Msg
	switch val := resp.(type) {
	case string:
		msg = slack.Msg{
			Text:         val,
			ResponseType: "in_channel",
		}
	case *slack.Msg:
		msg = *val
	case slack.Msg:
		msg = val
	default:
		return Fail("handler returned unexpected type")
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return Fail("failed to serialize response")
	}
	return Succeed(string(jsonMsg))
}
