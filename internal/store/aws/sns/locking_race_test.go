package sns

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

// TestListEndpointsVersusDeleteApplicationCompletes pins the platform lock
// order: ListEndpointsByPlatformApplication and DeletePlatformApplication
// acquire platformAppMu and platformEndpointMu in ONE order (app before
// endpoint). With opposing orders the pair deadlocks permanently — every
// concurrent list+delete interleaving must instead complete.
func TestListEndpointsVersusDeleteApplicationCompletes(t *testing.T) {
	store := newTestSNSStore(t)

	const iterations = 30
	for i := 0; i < iterations; i++ {
		app, err := store.CreatePlatformApplication(&PlatformApplication{
			Name:     fmt.Sprintf("deadlock-app-%d", i),
			Platform: "GCM",
		})
		if err != nil {
			t.Fatalf("create application %d: %v", i, err)
		}
		if _, err := store.CreatePlatformEndpoint(&PlatformEndpoint{
			PlatformApplicationArn: app.PlatformApplicationArn,
			Token:                  fmt.Sprintf("token-%d", i),
		}); err != nil {
			t.Fatalf("create endpoint %d: %v", i, err)
		}

		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			_, _ = store.ListEndpointsByPlatformApplication(app.PlatformApplicationArn, common.ListOptions{})
		}()
		go func() {
			defer wg.Done()
			_ = store.DeletePlatformApplication(app.PlatformApplicationArn)
		}()

		done := make(chan struct{})
		go func() { wg.Wait(); close(done) }()
		select {
		case <-done:
		case <-time.After(15 * time.Second):
			t.Fatalf("list deadlocked against delete (iteration %d) — the platform lock order inverted", i)
		}
	}
}

// TestCreateEndpointVersusDeleteApplicationLeavesNoOrphans pins the
// CreatePlatformEndpoint TOCTOU closure: the application check and every
// write run under one platformAppMu+platformEndpointMu section, so a
// concurrent DeletePlatformApplication can never leave a fully-written
// endpoint (or index entry) behind for a deleted application.
func TestCreateEndpointVersusDeleteApplicationLeavesNoOrphans(t *testing.T) {
	store := newTestSNSStore(t)

	app, err := store.CreatePlatformApplication(&PlatformApplication{
		Name:     "orphan-app",
		Platform: "GCM",
	})
	if err != nil {
		t.Fatalf("create application: %v", err)
	}
	arn := app.PlatformApplicationArn

	var wg sync.WaitGroup
	for g := 0; g < 4; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				_, _ = store.CreatePlatformEndpoint(&PlatformEndpoint{
					PlatformApplicationArn: arn,
					Token:                  fmt.Sprintf("tok-%d-%d", g, i),
				})
			}
		}(g)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = store.DeletePlatformApplication(arn)
	}()

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(30 * time.Second):
		t.Fatal("create/delete interleaving did not settle — lock order regression")
	}

	if _, err := store.GetPlatformApplication(arn); err != ErrPlatformApplicationNotFound {
		t.Fatalf("application still present after delete: %v", err)
	}
	if n := store.platformEndpointsStore.Count(); n != 0 {
		t.Fatalf("%d endpoint records survived their application's deletion — a create wrote outside the delete's lock window", n)
	}
	for name, bucket := range map[string]storage.Bucket{
		"app-index":   store.platformAppEndpointsIndex,
		"token-index": store.endpointTokenIndex,
	} {
		iter := bucket.ScanPrefix([]byte(arn + "\x00"))
		leaked := iter.Next()
		iter.Close()
		if leaked {
			t.Fatalf("%s retains an entry for the deleted application", name)
		}
	}
}

// TestCreateTopicIdempotentRetryAppliesTags pins the exists-path tag
// application: a retried CreateTopic whose first attempt failed between
// the record write and the tag write must still land its tags — the
// idempotent short-circuit merges the requested set instead of discarding
// it.
func TestCreateTopicIdempotentRetryAppliesTags(t *testing.T) {
	store := newTestSNSStore(t)

	first, err := store.CreateTopic(&Topic{Name: "retry-tags"}, map[string]string{"team": "core"})
	if err != nil {
		t.Fatalf("first create: %v", err)
	}
	second, err := store.CreateTopic(&Topic{Name: "retry-tags"}, map[string]string{"env": "prod"})
	if err != nil {
		t.Fatalf("idempotent create: %v", err)
	}
	if second.Arn != first.Arn {
		t.Fatalf("idempotent create returned %s, want the existing %s", second.Arn, first.Arn)
	}

	tags, err := store.ListTagsForResource(first.Arn)
	if err != nil {
		t.Fatalf("list tags: %v", err)
	}
	got := make(map[string]string, len(tags))
	for _, tag := range tags {
		got[tag.Key] = tag.Value
	}
	want := map[string]string{"team": "core", "env": "prod"}
	if len(got) != len(want) {
		t.Fatalf("topic carries %v, want the merged set %v", got, want)
	}
	for k, v := range want {
		if got[k] != v {
			t.Fatalf("topic carries %v, want the merged set %v", got, want)
		}
	}
}

// TestCreatePlatformEndpointMergeBackfillsAppIndex pins the merge-time
// self-heal — the single index-reconciliation mechanism: when a same-token
// create merges into an existing endpoint whose app-index entry is
// missing, the merge backfills it, returning the endpoint to
// ListEndpointsByPlatformApplication's view.
func TestCreatePlatformEndpointMergeBackfillsAppIndex(t *testing.T) {
	store := newTestSNSStore(t)

	app, err := store.CreatePlatformApplication(&PlatformApplication{
		Name:     "backfill-app",
		Platform: "GCM",
	})
	if err != nil {
		t.Fatalf("create application: %v", err)
	}
	ep, err := store.CreatePlatformEndpoint(&PlatformEndpoint{
		PlatformApplicationArn: app.PlatformApplicationArn,
		Token:                  "tok-merge",
	})
	if err != nil {
		t.Fatalf("create endpoint: %v", err)
	}

	// Simulate the lost app-index entry and prove the precondition: without
	// the index, the list cannot see the endpoint.
	idxKey := []byte(app.PlatformApplicationArn + "\x00" + ep.EndpointArn)
	if err := store.platformAppEndpointsIndex.Delete(idxKey); err != nil {
		t.Fatalf("drop index entry: %v", err)
	}
	hidden, err := store.ListEndpointsByPlatformApplication(app.PlatformApplicationArn, common.ListOptions{})
	if err != nil {
		t.Fatalf("list without index entry: %v", err)
	}
	if len(hidden.Items) != 0 {
		t.Fatalf("precondition: list saw %d endpoints without an index entry", len(hidden.Items))
	}

	merged, err := store.CreatePlatformEndpoint(&PlatformEndpoint{
		PlatformApplicationArn: app.PlatformApplicationArn,
		Token:                  "tok-merge",
		CustomUserData:         "updated",
	})
	if err != nil {
		t.Fatalf("merge create: %v", err)
	}
	if merged.EndpointArn != ep.EndpointArn {
		t.Fatalf("same-token create returned new ARN %s, want the existing %s", merged.EndpointArn, ep.EndpointArn)
	}

	listed, err := store.ListEndpointsByPlatformApplication(app.PlatformApplicationArn, common.ListOptions{})
	if err != nil {
		t.Fatalf("list after merge: %v", err)
	}
	if len(listed.Items) != 1 || listed.Items[0].EndpointArn != ep.EndpointArn {
		t.Fatalf("merge did not backfill the app-index entry: list returned %v", listed.Items)
	}
	if listed.Items[0].CustomUserData != "updated" {
		t.Fatalf("merge did not apply CustomUserData: %q", listed.Items[0].CustomUserData)
	}
}
