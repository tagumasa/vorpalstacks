package s3

import (
	"testing"

	"vorpalstacks/internal/core/storage"
)

// The EventBridge member must survive the persistence round-trip: a fresh
// store instance over the same Pebble reads back the configuration with the
// field intact, alongside the classic topic destinations.
func TestNotificationConfigurationEventBridgePersists(t *testing.T) {
	st, err := storage.NewPebbleStorage(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("pebble storage: %v", err)
	}
	defer st.Close()

	bucketStore := NewBucketStore(st, "123456789012", "test-region")
	if _, err := bucketStore.Create("persist-bucket", "test-region"); err != nil {
		t.Fatalf("Create: %v", err)
	}

	set := &NotificationConfiguration{
		TopicConfigurations: []TopicNotificationConfiguration{{
			Id:       "topic-1",
			TopicArn: "arn:aws:sns:us-east-1:123456789012:alerts",
			Events:   []string{"s3:ObjectCreated:*"},
			Filter: &NotificationConfigurationFilter{Key: &S3KeyFilter{
				FilterRules: []FilterRule{{Name: "prefix", Value: "logs/"}},
			}},
		}},
		EventBridgeConfiguration: &EventBridgeNotificationConfiguration{},
	}
	if err := bucketStore.SetNotificationConfiguration("persist-bucket", set); err != nil {
		t.Fatalf("SetNotificationConfiguration: %v", err)
	}

	// A second store instance models a process restart: everything it
	// returns comes from the persisted record.
	reopened := NewBucketStore(st, "123456789012", "test-region")
	got, err := reopened.GetNotificationConfiguration("persist-bucket")
	if err != nil {
		t.Fatalf("GetNotificationConfiguration: %v", err)
	}
	if got == nil {
		t.Fatal("notification configuration missing after reopen")
	}
	if got.EventBridgeConfiguration == nil {
		t.Fatal("EventBridgeConfiguration dropped on the persistence round-trip")
	}
	if len(got.TopicConfigurations) != 1 {
		t.Fatalf("topic configurations = %d, want 1", len(got.TopicConfigurations))
	}
	topic := got.TopicConfigurations[0]
	if topic.Id != "topic-1" || topic.TopicArn != "arn:aws:sns:us-east-1:123456789012:alerts" {
		t.Fatalf("topic round-trip = %+v", topic)
	}
	if topic.Filter == nil || topic.Filter.Key == nil || len(topic.Filter.Key.FilterRules) != 1 ||
		topic.Filter.Key.FilterRules[0].Name != "prefix" || topic.Filter.Key.FilterRules[0].Value != "logs/" {
		t.Fatalf("filter round-trip = %+v", topic.Filter)
	}
}
