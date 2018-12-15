package events

import (
	"bytes"
	"net/http"
	"net/url"
	"os"

	"github.com/aws/aws-lambda-go/events"
)

// Response aliases the API GW Response type
type Response events.APIGatewayProxyResponse

// Request aliases the API GW Request type
type Request events.APIGatewayProxyRequest

// BodyAsParams attempts to parse the request body as URL parameters
func (r *Request) BodyAsParams() (map[string]string, error) {
	result := make(map[string]string)
	vals, err := url.ParseQuery(r.Body)
	if err != nil {
		return result, err
	}
	for key := range vals {
		result[key] = vals.Get(key)
	}
	return result, nil
}

// ToHTTP returns the API Gateway request as an HTTP Request object
func (r *Request) ToHTTP() (*http.Request, error) {
	url := "https://" + r.Headers["Host"] + r.Path
	body := bytes.NewBufferString(r.Body)
	hr, err := http.NewRequest(r.HTTPMethod, url, body)
	if err != nil {
		return hr, err
	}
	for k, v := range r.Headers {
		hr.Header.Add(k, v)
	}
	return hr, nil
}

// Fail returns a message with an HTTP 500
func Fail(msg string) (Response, error) {
	return Respond(500, msg)
}

// Succeed returns a message with an HTTP 200
func Succeed(msg string) (Response, error) {
	return Respond(200, msg)
}

// Reject returns a message with an HTTP 403
func Reject(msg string) (Response, error) {
	return Respond(403, msg)
}

// Redirect returns a redirect to a new URL
func Redirect(url string, code int) (Response, error) {
	return Response{
		StatusCode: code,
		Headers: map[string]string{
			"Location": url,
		},
	}, nil
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
