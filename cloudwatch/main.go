package cloudwatch

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

// SNSEvent proxies to events object
type SNSEvent events.SNSEvent

// SNSEventRecord proxies to events object
type SNSEventRecord events.SNSEventRecord

// SNSEntity proxies to events object
type SNSEntity events.SNSEntity

// AlarmMessage describes a CloudWatch Alarm SNS event
type AlarmMessage struct {
	Subject          string       `json:"-"`
	AlarmName        string       `json:"AlarmName"`
	AlarmDescription string       `json:"AlarmDescription"`
	AWSAccountID     string       `json:"AWSAccountId"`
	OldStateValue    string       `json:"OldStateValue"`
	NewStateValue    string       `json:"NewStateValue"`
	NewStateReason   string       `json:"NewStateReason"`
	StateChangeTime  time.Time    `json:"StateChangeTime"`
	Region           string       `json:"Region"`
	Trigger          AlarmTrigger `json:"Trigger"`
}

// AlarmTrigger describes the SNS Alarm trigger
type AlarmTrigger struct {
	MetricName         string            `json:"MetricName"`
	Namespace          string            `json:"Namespace"`
	Statistic          string            `json:"Statistic"`
	Unit               string            `json:"Unit"`
	Dimensions         map[string]string `json:"Dimensions"`
	Period             int               `json:"Period"`
	EvaluationPeriods  int               `json:"EvaluationPeriods"`
	ComparisonOperator string            `json:"ComparisonOperator"`
	Threshold          float32           `json:"Threshold"`
}

// DecodedMessage decodes the message body of an SNS record
func (r *SNSEventRecord) DecodedMessage() (AlarmMessage, error) {
	var m AlarmMessage
	m.Subject = r.SNS.Subject
	err := json.Unmarshal([]byte(r.SNS.Message), &m)
	return m, err
}
