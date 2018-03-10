package cloudwatch

import (
	"encoding/json"
	"time"
)

// SNSEvent is a set of EventRecords
type SNSEvent struct {
	Records []SNSEventRecord `json:"Records"`
}

// SNSEventRecord is an SNS event
type SNSEventRecord struct {
	EventVersion         string    `json:"EventVersion"`
	EventSubscriptionArn string    `json:"EventSubscriptionArn"`
	EventSource          string    `json:"EventSource"`
	SNS                  SNSEntity `json:"Sns"`
}

// SNSEntity is the message details for the SNS event
type SNSEntity struct {
	Signature         string                 `json:"Signature"`
	MessageID         string                 `json:"MessageId"`
	Type              string                 `json:"Type"`
	TopicArn          string                 `json:"TopicArn"`
	MessageAttributes map[string]interface{} `json:"MessageAttributes"`
	SignatureVersion  string                 `json:"SignatureVersion"`
	Timestamp         string                 `json:"Timestamp"`
	SigningCertURL    string                 `json:"SigningCertUrl"`
	Message           string                 `json:"Message"`
	UnsubscribeURL    string                 `json:"UnsubscribeUrl"`
	Subject           string                 `json:"Subject"`
}

// AlarmMessage describes a CloudWatch Alarm SNS event
type AlarmMessage struct {
	Subject          string       `json:"-"`
	AlarmName        string       `json:"AlarmName"`
	AlarmDescription string       `json:"AlarmDescription"`
	AWSAccountID     string       `json:"AWSAccountId"`
	OldStateValue    string       `json:"OldStateValue"`
	NewStateValue    string       `json:"NewStateValue"`
	NewStateReason   string       `json:"NewStateReason"`
	StateChangeTime  string       `json:"StateChangeTime"`
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
