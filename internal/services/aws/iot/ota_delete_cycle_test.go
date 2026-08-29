package iot

import (
	"testing"

	"vorpalstacks/internal/core/storage"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// TestOTAUpdateDeleteCycle pins the delete cycle for an OTA update: the
// create materialises the named IoT job, so the delete forces the
// still-in-progress job away with the update.
func TestOTAUpdateDeleteCycle(t *testing.T) {
	rawStore, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { rawStore.Close() })
	store := iotstore.NewIotStore(rawStore, "000000000000", "us-east-1", nil)
	svc := &IoTService{}

	created, err := svc.createOTAUpdateCore(store, CreateOTAUpdateInput{
		OtaUpdateID: "ota-cycle",
		Targets:     []string{"arn:aws:iot:us-east-1:000000000000:thing/p"},
		Files:       []interface{}{map[string]interface{}{"fileName": "f.bin"}},
		RoleArn:     "arn:aws:iam::000000000000:role/ota",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.AwsIotJobID == "" {
		t.Fatal("expected a minted awsIotJobId")
	}

	if _, err := svc.getOTAUpdateCore(store, "ota-cycle"); err != nil {
		t.Fatalf("get: %v", err)
	}
	if err := svc.deleteOTAUpdateCore(store, DeleteOTAUpdateInput{OtaUpdateID: "ota-cycle", ForceDeleteAWSJob: true}); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := svc.getOTAUpdateCore(store, "ota-cycle"); err == nil {
		t.Fatal("expected the deleted update to be gone")
	}
}

// TestOTAUpdateDeleteStreamIgnoredForUserStreams pins the documented
// deleteStream semantics: the member only removes streams the OTAUpdate
// process itself created, and this platform's CreateOTAUpdate never
// generates streams, so a user-supplied stream referenced through the
// files' stream location survives the OTA update deletion.
func TestOTAUpdateDeleteStreamIgnoredForUserStreams(t *testing.T) {
	rawStore, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { rawStore.Close() })
	store := iotstore.NewIotStore(rawStore, "000000000000", "us-east-1", nil)
	svc := &IoTService{}

	streamID := "ota-user-stream"
	if err := store.PutGeneric("stream/"+streamID, map[string]interface{}{
		"streamId": streamID,
		"roleArn":  "arn:aws:iam::000000000000:role/ota",
	}); err != nil {
		t.Fatalf("seed stream: %v", err)
	}
	files := []interface{}{map[string]interface{}{
		"fileName":     "firmware.bin",
		"fileLocation": map[string]interface{}{"stream": map[string]interface{}{"streamId": streamID}},
	}}
	if _, err := svc.createOTAUpdateCore(store, CreateOTAUpdateInput{
		OtaUpdateID: "ota-user-stream-update",
		Targets:     []string{"arn:aws:iot:us-east-1:000000000000:thing/dev"},
		Files:       files,
		RoleArn:     "arn:aws:iam::000000000000:role/ota",
	}); err != nil {
		t.Fatalf("create OTA update: %v", err)
	}
	if err := svc.deleteOTAUpdateCore(store, DeleteOTAUpdateInput{
		OtaUpdateID:       "ota-user-stream-update",
		DeleteStream:      true,
		ForceDeleteAWSJob: true,
	}); err != nil {
		t.Fatalf("delete OTA update: %v", err)
	}
	streamExists, err := store.GetGenericExists("stream/"+streamID, &map[string]interface{}{})
	if err != nil {
		t.Fatalf("load stream: %v", err)
	}
	if !streamExists {
		t.Fatal("user-supplied stream must survive deleteStream")
	}
	otaExists, err := store.GetGenericExists("otaUpdate/ota-user-stream-update", &map[string]interface{}{})
	if err != nil {
		t.Fatalf("load OTA update: %v", err)
	}
	if otaExists {
		t.Fatal("expected the OTA update record to be deleted")
	}
}
