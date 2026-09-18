package cloudtrail

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"vorpalstacks/internal/core/storage"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

// TestPurgeEventHistoryAllRegions_PurgesConfiguredButNeverCalledRegion is
// the region-enumeration pin: a region whose storage is open and seeded
// with expired events is purged even though no CloudTrail request has ever
// created a service store for it. Enumeration that ranged only over the
// service's own store cache would leave such a region unpurged.
func TestPurgeEventHistoryAllRegions_PurgesConfiguredButNeverCalledRegion(t *testing.T) {
	tmpDir := "./tmp/cloudtrail-purger-test"
	os.RemoveAll(tmpDir)
	t.Cleanup(func() { os.RemoveAll(tmpDir) })

	sm, err := storage.NewRegionStorageManager(&storage.Config{Path: tmpDir})
	require.NoError(t, err)
	t.Cleanup(func() { sm.Close() })

	// Configure a secondary region: its storage is opened and seeded with
	// backdated events directly at the store layer; no request ever
	// routes through the service for it.
	const quietRegion = "eu-west-1"
	quietStorage, err := sm.GetStorage(quietRegion)
	require.NoError(t, err)
	quietStore := cloudtrailstore.NewCloudTrailStore(quietStorage, "acc123", quietRegion)

	old := time.Now().UTC().Add(-48 * time.Hour).Truncate(time.Hour)
	for i := 0; i < 3; i++ {
		require.NoError(t, quietStore.PutEvent(&cloudtrailstore.Event{
			EventID:     fmt.Sprintf("quiet-old-%d", i),
			EventName:   "CreateTrail",
			EventSource: "cloudtrail.amazonaws.com",
			EventTime:   old.Add(time.Duration(i) * time.Minute),
			UserIdentity: &cloudtrailstore.UserIdentity{
				Type:     "IAMUser",
				UserName: "alice",
			},
		}))
	}
	require.NoError(t, quietStore.PutEvent(&cloudtrailstore.Event{
		EventID:     "quiet-fresh",
		EventName:   "CreateTrail",
		EventSource: "cloudtrail.amazonaws.com",
		EventTime:   time.Now().UTC().Add(-time.Minute),
		UserIdentity: &cloudtrailstore.UserIdentity{
			Type:     "IAMUser",
			UserName: "alice",
		},
	}))

	svc := NewCloudTrailService("acc123", "us-east-1")
	svc.SetStorageManager(sm)
	svc.retention = 24 * time.Hour

	// The sweep body the hourly tick runs; drives it directly so the
	// test does not wait for the ticker.
	svc.purgeEventHistoryAllRegions()

	_, err = quietStore.GetEventByID("quiet-old-0")
	assert.ErrorIs(t, err, cloudtrailstore.ErrEventNotFound, "expired event purged in the never-called region")
	_, err = quietStore.GetEventByID("quiet-old-1")
	assert.ErrorIs(t, err, cloudtrailstore.ErrEventNotFound)
	_, err = quietStore.GetEventByID("quiet-old-2")
	assert.ErrorIs(t, err, cloudtrailstore.ErrEventNotFound)

	fresh, err := quietStore.GetEventByID("quiet-fresh")
	require.NoError(t, err, "recent event retained")
	assert.Equal(t, "CreateTrail", fresh.EventName)

	events, nextToken, err := quietStore.LookupEvents(cloudtrailstore.EventQuery{MaxResults: 50})
	require.NoError(t, err)
	assert.Empty(t, nextToken)
	assert.Len(t, events, 1, "only the retained event is visible")

	svc.Stop()
}

// TestBootPurgeEventHistory pins the TEST_MODE boot bound: entries older
// than 24 hours are removed across every open region storage, newer
// entries survive, and the sweep is safe to repeat at every boot.
func TestBootPurgeEventHistory(t *testing.T) {
	tmpDir := "./tmp/cloudtrail-boot-purge-test"
	os.RemoveAll(tmpDir)
	t.Cleanup(func() { os.RemoveAll(tmpDir) })

	sm, err := storage.NewRegionStorageManager(&storage.Config{Path: tmpDir})
	require.NoError(t, err)
	t.Cleanup(func() { sm.Close() })

	regionStorage, err := sm.GetStorage("us-east-1")
	require.NoError(t, err)
	st := cloudtrailstore.NewCloudTrailStore(regionStorage, "acc123", "us-east-1")

	seed := func(id string, age time.Duration) {
		require.NoError(t, st.PutEvent(&cloudtrailstore.Event{
			EventID:     id,
			EventName:   "CreateTrail",
			EventSource: "cloudtrail.amazonaws.com",
			EventTime:   time.Now().UTC().Add(-age),
			UserIdentity: &cloudtrailstore.UserIdentity{
				Type:     "IAMUser",
				UserName: "alice",
			},
		}))
	}
	seed("purged-25h", 25*time.Hour)
	seed("purged-48h", 48*time.Hour)
	seed("retained-23h", 23*time.Hour)
	seed("retained-now", time.Minute)

	svc := NewCloudTrailService("acc123", "us-east-1")
	svc.SetStorageManager(sm)
	svc.BootPurgeEventHistory()
	svc.Stop()

	_, err = st.GetEventByID("purged-25h")
	assert.ErrorIs(t, err, cloudtrailstore.ErrEventNotFound)
	_, err = st.GetEventByID("purged-48h")
	assert.ErrorIs(t, err, cloudtrailstore.ErrEventNotFound)
	for _, id := range []string{"retained-23h", "retained-now"} {
		_, err := st.GetEventByID(id)
		assert.NoError(t, err, "%s survives the boot purge", id)
	}

	// A second boot over the same data removes nothing further.
	svc2 := NewCloudTrailService("acc123", "us-east-1")
	svc2.SetStorageManager(sm)
	svc2.BootPurgeEventHistory()
	svc2.Stop()

	events, _, err := st.LookupEvents(cloudtrailstore.EventQuery{MaxResults: 50})
	require.NoError(t, err)
	assert.Len(t, events, 2)
}

func TestResolveEventHistoryRetention(t *testing.T) {
	t.Run("default when unset", func(t *testing.T) {
		t.Setenv(eventHistoryRetentionEnv, "")
		d, err := resolveEventHistoryRetention()
		require.NoError(t, err)
		assert.Equal(t, cloudtrailstore.EventHistoryRetention, d)
	})

	t.Run("zero selects the default", func(t *testing.T) {
		t.Setenv(eventHistoryRetentionEnv, "0")
		d, err := resolveEventHistoryRetention()
		require.NoError(t, err)
		assert.Equal(t, cloudtrailstore.EventHistoryRetention, d)
	})

	t.Run("override honoured", func(t *testing.T) {
		t.Setenv(eventHistoryRetentionEnv, "48h")
		d, err := resolveEventHistoryRetention()
		require.NoError(t, err)
		assert.Equal(t, 48*time.Hour, d)
	})

	t.Run("minimum accepted", func(t *testing.T) {
		t.Setenv(eventHistoryRetentionEnv, "1h")
		d, err := resolveEventHistoryRetention()
		require.NoError(t, err)
		assert.Equal(t, time.Hour, d)
	})

	t.Run("invalid duration rejected", func(t *testing.T) {
		t.Setenv(eventHistoryRetentionEnv, "not-a-duration")
		_, err := resolveEventHistoryRetention()
		assert.Error(t, err)
	})

	t.Run("below minimum rejected", func(t *testing.T) {
		t.Setenv(eventHistoryRetentionEnv, "30m")
		_, err := resolveEventHistoryRetention()
		assert.Error(t, err)
	})
}

func TestStartEventHistoryPurger(t *testing.T) {
	t.Run("invalid override rejected at startup", func(t *testing.T) {
		t.Setenv(eventHistoryRetentionEnv, "30m")
		svc := NewCloudTrailService("acc123", "us-east-1")
		err := svc.StartEventHistoryPurger()
		assert.Error(t, err)
		svc.Stop()
	})

	t.Run("valid configuration starts and stops the worker", func(t *testing.T) {
		t.Setenv(eventHistoryRetentionEnv, "48h")
		svc := NewCloudTrailService("acc123", "us-east-1")
		require.NoError(t, svc.StartEventHistoryPurger())
		assert.Equal(t, 48*time.Hour, svc.retention)
		// Stop must return once the worker goroutine observes the
		// cancellation.
		done := make(chan struct{})
		go func() { svc.Stop(); close(done) }()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("Stop did not return after cancellation")
		}
	})
}
