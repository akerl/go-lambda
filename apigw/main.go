package apigw

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Start runs the API GW Lambda
func Start(l Lambda) {
	lambda.Start(l.BuildHandler())
}

// Request aliases the API GW Request type
type Request events.APIGatewayProxyRequest

// Response aliases the API GW Response type
type Response events.APIGatewayProxyResponse

// Handler describes the signature for an API GW request handler
type Handler func(Request, Params) (string, error)

// Lambda defines a set of handlers
type Lambda struct {
	Handlers map[string]Handler
	Defaults map[string]string
}

// BuildHandler returns a function that routes to the appropriate handler
func (l *Lambda) BuildHandler() func(req Request) (Response, error) {
	return func(req Request) (Response, error) {
		if req.QueryStringParameters["trigger_id"] != "" {
			actualToken := req.QueryStringParameters["token"]
			params := Params{
				Request: &req,
			}
			expectedToken := params.Lookup("slack_token")
			if expectedToken == "" {
				return Fail("no slack_token provided")
			} else if expectedToken != actualToken {
				return Fail("invalid slack token")
			}
			return l.run("slack", req)
		}
		return l.run("default", req)
	}
}

func (l *Lambda) run(name string, req Request) (Response, error) {
	params := Params{
		Request:  &req,
		Defaults: l.Defaults,
	}
	h := l.Handlers[name]
	if h == nil {
		h = l.Handlers["default"]
	}
	if h == nil {
		return Fail(fmt.Sprintf("no valid handler found: %s", name))
	}
	body, err := h(req, params)
	if err != nil {
		return Fail(err.Error())
	}
	return Succeed(body)
}

// Fail returns a message with an HTTP 500
func Fail(msg string) (Response, error) {
	return Respond(500, msg)
}

// Succeed returns a message with an HTTP 200
func Succeed(msg string) (Response, error) {
	return Respond(200, msg)
}

// Respond builds a response with a given HTTP code and text message
func Respond(code int, msg string) (Response, error) {
	return Response{
		Body:       msg,
		StatusCode: code,
	}, nil
}

// Params allows looking up parameters from Stage Variables, Path Parameters,
// and the Lambda environment variables
type Params struct {
	Request  *Request
	Defaults map[string]string
}

// Lookup returns a value for the parameter, if it's set in an available field
func (p *Params) Lookup(name string) string {
	options := []string{
		p.Request.StageVariables[name],
		p.Request.PathParameters[name],
		p.Defaults[name],
		os.Getenv(name),
	}
	for _, i := range options {
		if i != "" {
			return i
		}
	}
	return ""
}
