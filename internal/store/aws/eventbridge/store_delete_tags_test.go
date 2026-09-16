package eventbridge

import (
	"testing"

	"vorpalstacks/internal/core/storage"
)

func newDeleteTagsTestStore(t *testing.T) *EventsStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	return NewEventsStore(st, "000000000000", "us-east-1")
}

// TestDeleteCleansTagRows pins the delete-path tag hygiene for the four
// families that previously leaked TagStore rows: deleting a rule, archive,
// connection or API destination also deletes its tag records, so a
// recreate with the same name cannot observe the previous life's tags.
func TestDeleteCleansTagRows(t *testing.T) {
	ctx := t.Context()

	cases := []struct {
		name   string
		create func(store *EventsStore) (arn string, err error)
		delete func(store *EventsStore) error
	}{
		{
			name: "rule",
			create: func(store *EventsStore) (string, error) {
				rule := &Rule{Name: "tagged", EventBusName: "default"}
				if err := store.CreateRule(ctx, rule); err != nil {
					return "", err
				}
				return rule.ARN, nil
			},
			delete: func(store *EventsStore) error { return store.DeleteRule(ctx, "default", "tagged") },
		},
		{
			name: "archive",
			create: func(store *EventsStore) (string, error) {
				archive := &Archive{Name: "tagged", EventBusName: "default"}
				if err := store.CreateArchive(ctx, archive); err != nil {
					return "", err
				}
				return archive.ARN, nil
			},
			delete: func(store *EventsStore) error { return store.DeleteArchive(ctx, "tagged") },
		},
		{
			name: "connection",
			create: func(store *EventsStore) (string, error) {
				connection := &Connection{Name: "tagged", AuthorizationType: "BASIC"}
				if err := store.CreateConnection(ctx, connection); err != nil {
					return "", err
				}
				return connection.ARN, nil
			},
			delete: func(store *EventsStore) error { return store.DeleteConnection(ctx, "tagged") },
		},
		{
			name: "api-destination",
			create: func(store *EventsStore) (string, error) {
				apiDest := &ApiDestination{Name: "tagged", ConnectionARN: "arn:aws:events:us-east-1:000000000000:connection/c"}
				if err := store.CreateApiDestination(ctx, apiDest); err != nil {
					return "", err
				}
				return apiDest.ARN, nil
			},
			delete: func(store *EventsStore) error { return store.DeleteApiDestination(ctx, "tagged") },
		},
	}

	for _, tc := range cases {
		store := newDeleteTagsTestStore(t)
		arn, err := tc.create(store)
		if err != nil {
			t.Fatalf("%s: create: %v", tc.name, err)
		}
		if err := store.TagStore.Tag(arn, map[string]string{"env": "test"}); err != nil {
			t.Fatalf("%s: tag: %v", tc.name, err)
		}
		if err := tc.delete(store); err != nil {
			t.Fatalf("%s: delete: %v", tc.name, err)
		}
		tags, err := store.TagStore.ListAsSlice(arn)
		if err != nil {
			t.Fatalf("%s: list tags: %v", tc.name, err)
		}
		if len(tags) != 0 {
			t.Errorf("%s: %d tag rows survived the delete (first key %q)", tc.name, len(tags), tags[0].Key)
		}
	}
}

// TestCreateStampsLastModified pins the creation stamping of the two
// families that previously left LastModifiedAt zero: a freshly created
// connection and API destination carry both timestamps, so Create*
// responses echo a stored value instead of fabricating one.
func TestCreateStampsLastModified(t *testing.T) {
	ctx := t.Context()
	store := newDeleteTagsTestStore(t)

	connection := &Connection{Name: "conn", AuthorizationType: "BASIC"}
	if err := store.CreateConnection(ctx, connection); err != nil {
		t.Fatalf("create connection: %v", err)
	}
	if connection.LastModifiedAt.IsZero() {
		t.Error("created connection carries a zero LastModifiedAt")
	}
	if !connection.LastModifiedAt.Equal(connection.CreatedAt) {
		t.Errorf("connection LastModifiedAt = %v, want the creation instant %v", connection.LastModifiedAt, connection.CreatedAt)
	}

	apiDest := &ApiDestination{Name: "dest", ConnectionARN: "arn:aws:events:us-east-1:000000000000:connection/c"}
	if err := store.CreateApiDestination(ctx, apiDest); err != nil {
		t.Fatalf("create api destination: %v", err)
	}
	if apiDest.LastModifiedAt.IsZero() {
		t.Error("created api destination carries a zero LastModifiedAt")
	}
	if !apiDest.LastModifiedAt.Equal(apiDest.CreatedAt) {
		t.Errorf("api destination LastModifiedAt = %v, want the creation instant %v", apiDest.LastModifiedAt, apiDest.CreatedAt)
	}
}
