package sfn

import (
	"context"
	"sync"
	"testing"
)

func newVersionTestMachine(t *testing.T, store *StepFunctionStore, name string) string {
	t.Helper()
	sm := &StateMachine{Name: name, Definition: `{"StartAt":"A","States":{"A":{"Type":"Pass","End":true}}}`}
	if err := store.CreateStateMachine(context.Background(), sm); err != nil {
		t.Fatal(err)
	}
	return sm.StateMachineArn
}

// TestPublishStateMachineVersionConcurrentIdempotent pins the publish
// serialisation: concurrent publishes of one revision must all return the
// SAME version record — publishing is idempotent per revision, and the
// duplicate record the unserialised find→number→Put window produced would
// inflate the version quota count.
func TestPublishStateMachineVersionConcurrentIdempotent(t *testing.T) {
	store := newHistoryTestStore(t)
	arn := newVersionTestMachine(t, store, "publish-race")

	const publishers = 8
	versions := make([]*StateMachineVersion, publishers)
	var wg sync.WaitGroup
	for i := range publishers {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			v, err := store.PublishStateMachineVersion(context.Background(), arn, "")
			if err != nil {
				t.Errorf("publish %d: %v", i, err)
				return
			}
			versions[i] = v
		}(i)
	}
	wg.Wait()

	for i, v := range versions {
		if v == nil {
			continue
		}
		if v.Version != versions[0].Version || v.StateMachineVersionArn != versions[0].StateMachineVersionArn {
			t.Errorf("publish %d returned version %d (%s), want the shared version %d — concurrent publishes must be idempotent per revision", i, v.Version, v.StateMachineVersionArn, versions[0].Version)
		}
	}
	count, err := store.CountStateMachineVersions(arn)
	if err != nil {
		t.Fatalf("count versions: %v", err)
	}
	if count != 1 {
		t.Errorf("CountStateMachineVersions = %d, want 1 — one revision must hold one version record", count)
	}
}
