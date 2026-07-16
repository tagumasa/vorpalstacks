package iotevents

import (
	"testing"
	"time"

	"vorpalstacks/internal/common/iotutil"
	storecommon "vorpalstacks/internal/store/aws/common"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

func TestDetectorModelConfig(t *testing.T) {
	now := time.Now()
	dm := &iotstore.DetectorModel{
		DetectorModelName:        "dm-1",
		DetectorModelARN:         "arn:aws:iotevents:us-east-1:123:detectorModel/dm-1",
		DetectorModelDescription: "test detector",
		RoleARN:                  "arn:aws:iam::123:role/test",
		Key:                      "device_id",
		EvaluationMethod:         "BATCH",
		DetectorModelVersion:     "1",
		Status:                   "ACTIVE",
		CreationDate:             now,
		LastModifiedDate:         now,
	}

	config := detectorModelConfig(dm)

	if config["detectorModelName"] != "dm-1" {
		t.Errorf("name = %v", config["detectorModelName"])
	}
	if config["detectorModelArn"] != dm.DetectorModelARN {
		t.Errorf("arn = %v", config["detectorModelArn"])
	}
	if config["key"] != "device_id" {
		t.Errorf("key = %v", config["key"])
	}
	if config["evaluationMethod"] != "BATCH" {
		t.Errorf("evalMethod = %v", config["evaluationMethod"])
	}
	if config["detectorModelVersion"] != "1" {
		t.Errorf("version = %v", config["detectorModelVersion"])
	}
	if config["creationTime"] != now.Unix() {
		t.Errorf("creationTime = %v, want %d", config["creationTime"], now.Unix())
	}
}

func TestInputConfig(t *testing.T) {
	now := time.Now()
	input := &iotstore.Input{
		InputName:        "temp-input",
		InputARN:         "arn:aws:iotevents:us-east-1:123:input/temp-input",
		InputDescription: "temperature sensor",
		Status:           "CREATING",
		CreationDate:     now,
		LastModifiedDate: now,
	}

	config := inputConfig(input)

	if config["inputName"] != "temp-input" {
		t.Errorf("name = %v", config["inputName"])
	}
	if config["inputArn"] != input.InputARN {
		t.Errorf("arn = %v", config["inputArn"])
	}
	if config["status"] != "CREATING" {
		t.Errorf("status = %v", config["status"])
	}
	if config["creationTime"] != now.Unix() {
		t.Errorf("creationTime = %v", config["creationTime"])
	}
}

func TestAlarmModelDescribeResponse(t *testing.T) {
	now := time.Now()
	am := &iotstore.AlarmModel{
		AlarmModelName:        "alarm-1",
		AlarmModelARN:         "arn:aws:iotevents:us-east-1:123:alarmModel/alarm-1",
		AlarmModelDescription: "temp alarm",
		RoleARN:               "arn:aws:iam::123:role/alarm",
		Status:                "ACTIVE",
		Severity:              "HIGH",
		AlarmModelVersion:     "1",
		CreationDate:          now,
		LastModifiedDate:      now,
	}

	resp := alarmModelDescribeResponse(am)

	if resp["alarmModelName"] != "alarm-1" {
		t.Errorf("name = %v", resp["alarmModelName"])
	}
	if resp["severity"] != "HIGH" {
		t.Errorf("severity = %v", resp["severity"])
	}
	if resp["roleArn"] != am.RoleARN {
		t.Errorf("roleArn = %v", resp["roleArn"])
	}
	if _, ok := resp["alarmModelDefinition"]; ok {
		t.Error("expected no alarmModelDefinition when nil")
	}

	am.AlarmModelDefinition = map[string]interface{}{"key": "value"}
	resp = alarmModelDescribeResponse(am)
	if resp["alarmModelDefinition"] == nil {
		t.Error("expected alarmModelDefinition when set")
	}
}

func TestAlarmModelSummary(t *testing.T) {
	now := time.Now()
	am := &iotstore.AlarmModel{
		AlarmModelName:    "alarm-1",
		AlarmModelARN:     "arn:aws:iotevents:us-east-1:123:alarmModel/alarm-1",
		AlarmModelVersion: "2",
		CreationDate:      now,
		LastModifiedDate:  now,
	}

	summary := alarmModelSummary(am)

	if summary["alarmModelName"] != "alarm-1" {
		t.Errorf("name = %v", summary["alarmModelName"])
	}
	if summary["alarmModelVersion"] != "2" {
		t.Errorf("version = %v", summary["alarmModelVersion"])
	}
	if _, ok := summary["severity"]; ok {
		t.Error("summary should not include severity")
	}
}

func TestParseAlarmTagsParam(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]interface{}
		want   map[string]string
	}{
		{
			name:   "no tags key",
			params: map[string]interface{}{},
			want:   nil,
		},
		{
			name: "valid tags lowercase keys",
			params: map[string]interface{}{
				"tags": []interface{}{
					map[string]interface{}{"key": "env", "value": "prod"},
					map[string]interface{}{"key": "team", "value": "iot"},
				},
			},
			want: map[string]string{"env": "prod", "team": "iot"},
		},
		{
			name: "tags with capitalised keys",
			params: map[string]interface{}{
				"tags": []interface{}{
					map[string]interface{}{"Key": "env", "Value": "prod"},
				},
			},
			want: map[string]string{"env": "prod"},
		},
		{
			name: "tags not an array",
			params: map[string]interface{}{
				"tags": "not-an-array",
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseAlarmTagsParam(tt.params)
			if tt.want == nil {
				if got != nil {
					t.Errorf("expected nil, got %v", got)
				}
				return
			}
			for k, v := range tt.want {
				if got[k] != v {
					t.Errorf("tag[%s] = %s, want %s", k, got[k], v)
				}
			}
		})
	}
}

func TestListResponse(t *testing.T) {
	items := []map[string]interface{}{
		{"name": "a"},
		{"name": "b"},
	}

	resp := listResponse("detectorModels", items, "")
	if len(resp["detectorModels"].([]map[string]interface{})) != 2 {
		t.Error("expected 2 items")
	}
	if _, ok := resp["nextToken"]; ok {
		t.Error("expected no nextToken")
	}

	resp = listResponse("detectorModels", items, "marker-abc")
	if resp["nextToken"] != "marker-abc" {
		t.Errorf("nextToken = %v", resp["nextToken"])
	}
}

func TestParseBatchMessagesInterfaceSlice(t *testing.T) {
	raw := []interface{}{
		map[string]interface{}{
			"inputName":      "sensor",
			"messagePayload": `{"temp":25}`,
		},
	}
	msgs, err := parseBatchMessages(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	if msgs[0].InputName != "sensor" {
		t.Errorf("inputName = %s", msgs[0].InputName)
	}
	if msgs[0].Payload["temp"] != float64(25) {
		t.Errorf("payload = %v", msgs[0].Payload)
	}
}

func TestParseBatchMessagesMapSlice(t *testing.T) {
	raw := []map[string]interface{}{
		{
			"inputName":      "sensor",
			"messagePayload": map[string]interface{}{"value": 42},
		},
	}
	msgs, err := parseBatchMessages(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	if msgs[0].Payload["value"] != 42 {
		t.Errorf("payload = %v", msgs[0].Payload)
	}
}

func TestParseBatchMessagesEmpty(t *testing.T) {
	msgs, err := parseBatchMessages([]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(msgs) != 0 {
		t.Errorf("expected 0 messages, got %d", len(msgs))
	}
}

func TestParseBatchMessagesInvalidType(t *testing.T) {
	_, err := parseBatchMessages("not-an-array")
	if err == nil {
		t.Fatal("expected error for non-array")
	}
}

func TestParseBatchMessagesNil(t *testing.T) {
	_, err := parseBatchMessages(nil)
	if err == nil {
		t.Fatal("expected error for nil")
	}
}

func TestStrFromMap(t *testing.T) {
	m := map[string]interface{}{
		"roleArn":     "arn:aws:iam::123:role/test",
		"queueUrl":    "https://sqs.us-east-1.amazonaws.com/123/queue",
		"targetArn":   "arn:aws:sns:us-east-1:123:topic",
		"functionArn": "arn:aws:lambda:us-east-1:123:function:fn",
	}

	tests := []struct {
		keys []string
		want string
	}{
		{[]string{"roleArn"}, "arn:aws:iam::123:role/test"},
		{[]string{"queueUrl"}, "https://sqs.us-east-1.amazonaws.com/123/queue"},
		{[]string{"nonexistent", "roleArn"}, "arn:aws:iam::123:role/test"},
		{[]string{"nonexistent"}, ""},
	}

	for _, tt := range tests {
		got := iotutil.StrFromMap(m, tt.keys...)
		if got != tt.want {
			t.Errorf("StrFromMap(%v) = %s, want %s", tt.keys, got, tt.want)
		}
	}
}

func TestNewDetectorActionAdapter(t *testing.T) {
	adapter := NewDetectorActionAdapter(nil, nil)
	if adapter == nil {
		t.Fatal("expected non-nil adapter")
	}
}

func TestParseAlarmListOptionsEmpty(t *testing.T) {
	opts := parseAlarmListOptions(map[string]interface{}{})
	if opts.Marker != "" {
		t.Errorf("marker = %s, want empty", opts.Marker)
	}
	if opts.MaxItems != 0 {
		t.Errorf("maxItems = %d, want 0", opts.MaxItems)
	}
}

func TestParseAlarmListOptionsWithToken(t *testing.T) {
	opts := parseAlarmListOptions(map[string]interface{}{
		"nextToken": "page-2",
	})
	if opts.Marker != "page-2" {
		t.Errorf("marker = %s, want page-2", opts.Marker)
	}
}

func TestParseAlarmListOptionsWithMaxResults(t *testing.T) {
	opts := parseAlarmListOptions(map[string]interface{}{
		"maxResults": "50",
	})
	if opts.MaxItems != 50 {
		t.Errorf("maxItems = %d, want 50", opts.MaxItems)
	}
}

func TestParseAlarmListOptionsInvalidMaxResults(t *testing.T) {
	opts := parseAlarmListOptions(map[string]interface{}{
		"maxResults": "not-a-number",
	})
	if opts.MaxItems != 0 {
		t.Errorf("maxItems = %d, want 0 for invalid", opts.MaxItems)
	}
}

func TestParseAlarmListOptionsZeroMaxResults(t *testing.T) {
	opts := parseAlarmListOptions(map[string]interface{}{
		"maxResults": "0",
	})
	if opts.MaxItems != 0 {
		t.Errorf("maxItems = %d, want 0 for zero", opts.MaxItems)
	}
}

var _ storecommon.ListOptions
