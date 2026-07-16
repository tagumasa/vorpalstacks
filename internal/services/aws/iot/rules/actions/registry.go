package actions

import (
	"context"
	"fmt"
)

type ActionPayload struct {
	Topic      string
	Raw        map[string]interface{}
	JSONBytes  []byte
	JSONString string
}

type ActionHandler func(d *Dispatcher, ctx context.Context, config *ActionConfig, p *ActionPayload) error

var actionRegistry = map[string]ActionHandler{
	"lambda":           (*Dispatcher).dispatchLambda,
	"sqs":              (*Dispatcher).dispatchSQS,
	"sns":              (*Dispatcher).dispatchSNS,
	"dynamoDB":         (*Dispatcher).dispatchDynamoDB,
	"dynamoDBv2":       (*Dispatcher).dispatchDynamoDBv2,
	"s3":               (*Dispatcher).dispatchS3,
	"kinesis":          (*Dispatcher).dispatchKinesis,
	"cloudwatchMetric": (*Dispatcher).dispatchCloudWatch,
	"cloudwatchAlarm":  (*Dispatcher).dispatchCloudWatchAlarm,
	"cloudwatchLogs":   (*Dispatcher).dispatchCloudWatchLogs,
	"republish":        (*Dispatcher).dispatchRepublish,
	"stepFunctions":    (*Dispatcher).dispatchStepFunctions,
	"iotEvents":        (*Dispatcher).dispatchIoTEvents,
	"http":             (*Dispatcher).dispatchHTTP,
	"firehose": func(_ *Dispatcher, _ context.Context, _ *ActionConfig, _ *ActionPayload) error {
		return fmt.Errorf("firehose action not supported in vorpalstacks")
	},
	"timestream": (*Dispatcher).dispatchTimestream,
}
