package apigw

import (
	"github.com/aws/aws-lambda-go/lambda"
)

// Start runs the API GW Lambda
func Start(r Router) {
	lambda.Start(r.Handler())
}
