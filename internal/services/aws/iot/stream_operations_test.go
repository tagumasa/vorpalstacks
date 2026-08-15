package iot

import (
	"context"
	"testing"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// streamId keys the stream catalog, so CreateStream with an existing id
// must fail with ResourceAlreadyExistsException instead of overwriting the
// record (which would reset streamVersion to one); the original record
// stays intact.
func TestCreateStreamRejectsDuplicateStreamID(t *testing.T) {
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	reqCtx := request.NewRequestContext(context.Background(), mgr, "000000000000", "us-east-1")
	svc := NewIoTService("000000000000")

	create := func() (interface{}, error) {
		return svc.CreateStream(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"streamId":    "firmware-stream",
			"description": "original description",
			"files":       []interface{}{map[string]interface{}{"fileId": 1}},
			"roleArn":     "arn:aws:iam::000000000000:role/deliver",
		}})
	}

	if _, err := create(); err != nil {
		t.Fatal(err)
	}
	if _, err := create(); err != iotstore.ErrStreamAlreadyExists {
		t.Fatalf("duplicate streamId error = %v, want ErrStreamAlreadyExists", err)
	}

	desc, err := svc.DescribeStream(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"streamId": "firmware-stream",
	}})
	if err != nil {
		t.Fatal(err)
	}
	info, _ := desc.(map[string]interface{})["streamInfo"].(map[string]interface{})
	if info == nil {
		t.Fatalf("streamInfo missing after duplicate rejection: %#v", desc)
	}
	if info["description"] != "original description" {
		t.Fatalf("description overwritten: %#v", info["description"])
	}
	// The generic KV store decodes numbers as float64, mirroring the dual
	// handling UpdateStream applies to the stored version.
	switch v := info["streamVersion"].(type) {
	case int64:
		if v != 1 {
			t.Fatalf("streamVersion must stay 1 after duplicate rejection, got %v", v)
		}
	case float64:
		if int64(v) != 1 {
			t.Fatalf("streamVersion must stay 1 after duplicate rejection, got %v", v)
		}
	default:
		t.Fatalf("unexpected streamVersion type %T", info["streamVersion"])
	}
}
