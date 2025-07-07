package slack

import (
	"encoding/json"

	"github.com/akerl/go-lambda/apigw/events"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// Handler checks for Slack requests and returns Slack messages
type Handler struct {
	HandleFunc    func(slackevents.EventsAPIInnerEvent) error
	ErrorFunc     func(events.Request, error) (events.Response, error)
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

	body, err := req.DecodedBody()
	if err != nil {
		return events.Reject("failed to decode body")
	}

	byteBody := []byte(body)

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

	event, err := slackevents.ParseEvent(
		json.RawMessage(req.Body),
		slackevents.OptionNoVerifyToken(),
	)
	if err != nil {
		return events.Fail("failed to parse event")
	}

	if event.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(req.Body), &r)
		if err != nil {
			return events.Fail("failed to unmarshal challenge")
		}
		return events.Response{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type": "text",
			},
			Body: r.Challenge,
		}, nil
	}

	innerEvent := event.InnerEvent

	err = h.HandleFunc(innerEvent)
	if err != nil {
		return events.Fail("failed to process event")
	}
	return events.Response{}, nil
}

func (h *Handler) Error(req events.Request, e error) (events.Response, error) {
	if h.ErrorFunc == nil {
		return events.Fail("No error function defined")
	}
	return h.ErrorFunc(req, e)
}
