package actions

import (
	"context"
	"encoding/json"
	"testing"

	"vorpalstacks/internal/common/iotutil"
)

func TestNewActionConfigFromMapLambda(t *testing.T) {
	m := map[string]interface{}{"functionArn": "arn:aws:lambda:us-east-1:123:function:fn"}
	ac := NewActionConfigFromMap("lambda", m)
	if ac.Type != "lambda" {
		t.Errorf("type = %s", ac.Type)
	}
	if ac.FunctionName != "arn:aws:lambda:us-east-1:123:function:fn" {
		t.Errorf("functionName = %s", ac.FunctionName)
	}
}

func TestNewActionConfigFromMapLambdaByName(t *testing.T) {
	m := map[string]interface{}{"functionName": "my-fn"}
	ac := NewActionConfigFromMap("lambda", m)
	if ac.FunctionName != "my-fn" {
		t.Errorf("functionName = %s", ac.FunctionName)
	}
}

func TestNewActionConfigFromMapSQS(t *testing.T) {
	m := map[string]interface{}{
		"queueUrl": "https://sqs.us-east-1.amazonaws.com/123/queue",
		"queueArn": "arn:aws:sqs:us-east-1:123:queue",
		"roleArn":  "arn:aws:iam::123:role/test",
	}
	ac := NewActionConfigFromMap("sqs", m)
	if ac.QueueURL != "https://sqs.us-east-1.amazonaws.com/123/queue" {
		t.Errorf("queueUrl = %s", ac.QueueURL)
	}
	if ac.TargetARN != "arn:aws:sqs:us-east-1:123:queue" {
		t.Errorf("targetARN = %s", ac.TargetARN)
	}
	if ac.RoleARN != "arn:aws:iam::123:role/test" {
		t.Errorf("roleARN = %s", ac.RoleARN)
	}
}

func TestNewActionConfigFromMapSNS(t *testing.T) {
	m := map[string]interface{}{"targetArn": "arn:aws:sns:us-east-1:123:topic"}
	ac := NewActionConfigFromMap("sns", m)
	if ac.TopicARN != "arn:aws:sns:us-east-1:123:topic" {
		t.Errorf("topicARN = %s", ac.TopicARN)
	}
}

func TestNewActionConfigFromMapDynamoDB(t *testing.T) {
	m := map[string]interface{}{
		"tableName": "my-table",
		"tableArn":  "arn:aws:dynamodb:us-east-1:123:table/my-table",
	}
	ac := NewActionConfigFromMap("dynamoDB", m)
	if ac.TableName != "my-table" {
		t.Errorf("tableName = %s", ac.TableName)
	}
	if ac.TargetARN != "arn:aws:dynamodb:us-east-1:123:table/my-table" {
		t.Errorf("targetARN = %s", ac.TargetARN)
	}
}

func TestNewActionConfigFromMapS3(t *testing.T) {
	m := map[string]interface{}{
		"bucketName": "my-bucket",
		"key":        "data/${timestamp()}.json",
	}
	ac := NewActionConfigFromMap("s3", m)
	if ac.BucketName != "my-bucket" {
		t.Errorf("bucketName = %s", ac.BucketName)
	}
	if ac.ObjectKey != "data/${timestamp()}.json" {
		t.Errorf("objectKey = %s", ac.ObjectKey)
	}
}

func TestNewActionConfigFromMapKinesis(t *testing.T) {
	m := map[string]interface{}{
		"streamName": "my-stream",
		"streamArn":  "arn:aws:kinesis:us-east-1:123:stream/my-stream",
	}
	ac := NewActionConfigFromMap("kinesis", m)
	if ac.StreamName != "my-stream" {
		t.Errorf("streamName = %s", ac.StreamName)
	}
}

func TestNewActionConfigFromMapRepublish(t *testing.T) {
	m := map[string]interface{}{"topic": "device/output"}
	ac := NewActionConfigFromMap("republish", m)
	if ac.RepublishTopic != "device/output" {
		t.Errorf("republishTopic = %s", ac.RepublishTopic)
	}
}

func TestNewActionConfigFromMapStepFunctions(t *testing.T) {
	m := map[string]interface{}{"stateMachineName": "my-stm"}
	ac := NewActionConfigFromMap("stepFunctions", m)
	if ac.TargetARN != "my-stm" {
		t.Errorf("targetARN = %s", ac.TargetARN)
	}
}

func TestNewActionConfigFromMapStepFunctionsByArn(t *testing.T) {
	m := map[string]interface{}{"stateMachineArn": "arn:aws:states:us-east-1:123:stm:my-stm"}
	ac := NewActionConfigFromMap("stepFunctions", m)
	if ac.TargetARN != "arn:aws:states:us-east-1:123:stm:my-stm" {
		t.Errorf("targetARN = %s", ac.TargetARN)
	}
}

func TestNewActionConfigFromMapUnknownType(t *testing.T) {
	m := map[string]interface{}{"foo": "bar"}
	ac := NewActionConfigFromMap("unknown", m)
	if ac.Type != "unknown" {
		t.Errorf("type = %s", ac.Type)
	}
}

func TestNewActionConfigFromMapRoleArnFallback(t *testing.T) {
	m := map[string]interface{}{"roleArn": "arn:aws:iam::123:role/test"}
	ac := NewActionConfigFromMap("lambda", m)
	if ac.RoleARN != "arn:aws:iam::123:role/test" {
		t.Errorf("roleARN = %s", ac.RoleARN)
	}
}

func TestBuildDynamoDBItemDefaultHash(t *testing.T) {
	payload := map[string]interface{}{"temp": 25}
	extra := map[string]interface{}{}
	key, attrs := buildDynamoDBItem(payload, extra)

	if _, ok := key["id"]; !ok {
		t.Error("expected default hash field 'id'")
	}
	if attrs["temp"] != 25 {
		t.Errorf("attrs[temp] = %v", attrs["temp"])
	}
}

func TestBuildDynamoDBItemCustomHash(t *testing.T) {
	payload := map[string]interface{}{"device_id": "dev-001", "temp": 30}
	extra := map[string]interface{}{
		"hashKeyField": "device_id",
	}
	key, attrs := buildDynamoDBItem(payload, extra)

	if key["device_id"] != "dev-001" {
		t.Errorf("hash value = %v, want dev-001", key["device_id"])
	}
	if attrs["device_id"] != "dev-001" {
		t.Errorf("attrs[device_id] = %v", attrs["device_id"])
	}
	if attrs["temp"] != 30 {
		t.Errorf("attrs[temp] = %v", attrs["temp"])
	}
}

func TestBuildDynamoDBItemCustomHashValue(t *testing.T) {
	payload := map[string]interface{}{"data": "x"}
	extra := map[string]interface{}{
		"hashKeyField":  "device_id",
		"hashKeyValue":  "fixed-id",
		"rangeKeyField": "timestamp",
		"rangeKeyValue": "12345",
	}
	key, attrs := buildDynamoDBItem(payload, extra)

	if key["device_id"] != "fixed-id" {
		t.Errorf("hash value = %v", key["device_id"])
	}
	if key["timestamp"] != "12345" {
		t.Errorf("range value = %v", key["timestamp"])
	}
	if attrs["data"] != "x" {
		t.Errorf("attrs[data] = %v", attrs["data"])
	}
}

func TestBuildDynamoDBItemMissingPayloadHash(t *testing.T) {
	payload := map[string]interface{}{"other": "val"}
	extra := map[string]interface{}{"hashKeyField": "device_id"}
	key, _ := buildDynamoDBItem(payload, extra)

	val, ok := key["device_id"]
	if !ok {
		t.Fatal("expected hash key to exist")
	}
	if val == "" {
		t.Error("expected non-empty auto-generated hash value")
	}
}

func TestMustJSON(t *testing.T) {
	data := map[string]interface{}{"key": "value", "num": float64(42)}
	result := mustJSON(data)

	var parsed map[string]interface{}
	if err := json.Unmarshal(result, &parsed); err != nil {
		t.Fatalf("mustJSON produced invalid JSON: %v", err)
	}
	if parsed["key"] != "value" {
		t.Errorf("key = %v", parsed["key"])
	}
}

func TestMustJSONString(t *testing.T) {
	result := mustJSON("hello")
	if string(result) != `"hello"` {
		t.Errorf("result = %s", string(result))
	}
}

func TestMustJSONNil(t *testing.T) {
	result := mustJSON(nil)
	if string(result) != "null" {
		t.Errorf("result = %s", string(result))
	}
}

func TestStrFromMapDispatcher(t *testing.T) {
	m := map[string]interface{}{
		"key1": "val1",
		"key2": "val2",
	}
	if v := iotutil.StrFromMap(m, "key1"); v != "val1" {
		t.Errorf("StrFromMap(key1) = %s", v)
	}
	if v := iotutil.StrFromMap(m, "missing", "key2"); v != "val2" {
		t.Errorf("StrFromMap(missing, key2) = %s", v)
	}
	if v := iotutil.StrFromMap(m, "missing"); v != "" {
		t.Errorf("StrFromMap(missing) = %s, want empty", v)
	}
}

func TestNewDispatcher(t *testing.T) {
	d := NewDispatcher(nil, nil)
	if d == nil {
		t.Fatal("expected non-nil dispatcher")
	}
}

func TestDispatcherSetFunctions(t *testing.T) {
	d := NewDispatcher(nil, nil)
	d.SetRepublishFn(func(_ context.Context, _ string, _ map[string]interface{}) error { return nil })
	d.SetHTTPPostFn(func(_ context.Context, _ string, _ []byte) error { return nil })
}

func TestDispatchNilConfig(t *testing.T) {
	d := NewDispatcher(nil, nil)
	err := d.Dispatch(context.Background(), nil, "topic", map[string]interface{}{"data": "x"})
	if err == nil {
		t.Fatal("expected error for nil bus and config")
	}
}
