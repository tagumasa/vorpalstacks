package cloudtrail

import (
	"fmt"
	"os"
	"time"

	"vorpalstacks/internal/core/logs"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

// eventHistoryRetentionEnv overrides the event-history retention period.
// It is platform operational configuration for dev/test boxes whose
// traffic is regression runs; unset or zero keeps the 90-day AWS default.
const eventHistoryRetentionEnv = "CLOUDTRAIL_EVENT_HISTORY_RETENTION"

// minEventHistoryRetention is the shortest retention the override accepts;
// anything below it is rejected at startup so a typo cannot flush recent
// history.
const minEventHistoryRetention = time.Hour

// purgeTickInterval is how often the retention worker sweeps the region
// storages: hourly in production, twice a second under TEST_MODE so the
// shortened restore window (edsRestoreWindow) keeps up with regression
// runs.
func purgeTickInterval() time.Duration {
	if os.Getenv("TEST_MODE") == "true" {
		return 500 * time.Millisecond
	}
	return time.Hour
}

// edsRestoreWindow is the effective event-data-store restore window: the
// documented seven days, shortened under TEST_MODE because the per-region
// event-data-store quota counts every lifecycle stage — a week-long window
// would hold the quota against regression suites after a single deleted
// store. The half-second window clears a deleted store between one test
// and the next; production keeps the documented window.
func edsRestoreWindow() time.Duration {
	if os.Getenv("TEST_MODE") == "true" {
		return 500 * time.Millisecond
	}
	return cloudtrailstore.EventDataStoreRestoreWindow
}

// resolveEventHistoryRetention returns the effective event-history
// retention: the ENV override when set (and at least one hour), otherwise
// the AWS 90-day default. An unset or zero value selects the default.
func resolveEventHistoryRetention() (time.Duration, error) {
	raw := os.Getenv(eventHistoryRetentionEnv)
	if raw == "" {
		return cloudtrailstore.EventHistoryRetention, nil
	}
	d, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("%s: invalid duration %q", eventHistoryRetentionEnv, raw)
	}
	if d == 0 {
		return cloudtrailstore.EventHistoryRetention, nil
	}
	if d < minEventHistoryRetention {
		return 0, fmt.Errorf("%s: %s is below the minimum of %s", eventHistoryRetentionEnv, d, minEventHistoryRetention)
	}
	return d, nil
}

// StartEventHistoryPurger validates the retention configuration and starts
// the hourly event-history retention worker. It returns an error (and
// starts nothing) when the ENV override is invalid, so a bad value is
// rejected at startup.
func (s *CloudTrailService) StartEventHistoryPurger() error {
	retention, err := resolveEventHistoryRetention()
	if err != nil {
		return err
	}
	s.retention = retention
	s.startEventHistoryPurger()
	return nil
}

func (s *CloudTrailService) startEventHistoryPurger() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer func() {
			if r := recover(); r != nil {
				logs.Error("PANIC in cloudtrail event-history purger, restarting",
					logs.Any("panic", r))
				// The restart runs synchronously so its wg.Add(1) lands
				// before this goroutine's deferred wg.Done(): an async
				// restart could race a concurrent Stop()'s wg.Wait past
				// zero and leak the restarted purger past shutdown.
				s.startEventHistoryPurger()
			}
		}()
		ticker := time.NewTicker(purgeTickInterval())
		defer ticker.Stop()

		for {
			select {
			case <-s.ctx.Done():
				return
			case <-ticker.C:
				s.purgeEventHistoryAllRegions()
				s.enforceEDSRetention()
				s.purgeAgedQueries()
			}
		}
	}()
}

// purgeEventHistoryAllRegions purges expired event history in every region
// known to the region storage manager. Enumeration deliberately comes from
// the manager and not from this service's store cache: a region whose
// storage is open must be purged even when no CloudTrail request has ever
// touched it, and a store is created here for exactly such regions.
func (s *CloudTrailService) purgeEventHistoryAllRegions() {
	if s.storageManager == nil {
		return
	}
	cutoff := time.Now().UTC().Add(-s.retention)
	s.sweepEventHistory("event-history purge", cutoff)
}

// sweepEventHistory purges entries older than cutoff across every open
// region storage, logging per-region counts under the given label. The
// per-region loop resolves stores through the region storage manager so
// configured-but-never-called regions are covered.
func (s *CloudTrailService) sweepEventHistory(label string, cutoff time.Time) int {
	if s.storageManager == nil {
		return 0
	}
	total := 0
	for _, region := range s.storageManager.GetActiveRegions() {
		store, err := s.GetStoreForRegion(region)
		if err != nil {
			logs.Error("cloudtrail "+label+": failed to resolve store",
				logs.String("region", region), logs.Err(err))
			continue
		}
		result, err := store.PurgeEventHistoryBefore(cutoff, 0)
		if err != nil {
			logs.Error("cloudtrail "+label+" failed",
				logs.String("region", region), logs.Err(err))
			continue
		}
		if result.Events > 0 {
			logs.Info("cloudtrail "+label,
				logs.String("region", region),
				logs.Int("events", result.Events),
				logs.Int("eventIDIndex", result.EventIDIndex),
				logs.Int("timeIndex", result.TimeIndex),
				logs.Int("eventNameIndex", result.EventNameIndex),
				logs.Int("usernameIndex", result.UsernameIndex),
				logs.Int("eventSourceIndex", result.EventSourceIndex),
				logs.Int("batches", result.Batches))
		}
		total += result.Events
	}
	return total
}

// testModeBootPurgeCutoff is the TEST_MODE boot bound: entries older than
// 24 hours are removed before serving. Regression suites record and look
// up immediately, so a one-day window never hides a suite's own events
// while the accumulated mass of prior runs is dropped at every boot.
const testModeBootPurgeCutoff = 24 * time.Hour

// BootPurgeEventHistory removes event-history entries older than the
// TEST_MODE boot bound from every open region storage, synchronously. It
// runs before the server starts serving so regression traffic always
// begins from a bounded history; non-TEST_MODE boots never call it.
func (s *CloudTrailService) BootPurgeEventHistory() {
	start := time.Now()
	cutoff := time.Now().UTC().Add(-testModeBootPurgeCutoff)
	total := s.sweepEventHistory("TEST_MODE boot purge", cutoff)
	total += s.bootPurgeEDSEvents(cutoff)
	logs.Info("cloudtrail TEST_MODE boot purge complete",
		logs.Int("events", total),
		logs.String("duration", time.Since(start).String()))
}

// enforceEDSRetention applies the per-store retention periods and the
// hard-delete restore window across every open region storage: each event
// data store's events older than its RetentionPeriod are purged, and a
// PENDING_DELETION store whose seven-day wait has closed is removed with
// its events.
func (s *CloudTrailService) enforceEDSRetention() {
	if s.storageManager == nil {
		return
	}
	for _, region := range s.storageManager.GetActiveRegions() {
		store, err := s.GetStoreForRegion(region)
		if err != nil {
			logs.Error("cloudtrail EDS retention: failed to resolve store",
				logs.String("region", region), logs.Err(err))
			continue
		}
		s.enforceEDSRetentionForStore(store)
	}
}

// purgeAgedQueries drops query records older than the ListQueries window
// across every open region storage: the listing is bounded to the past
// seven days, and the sweep keeps the stored records from outliving that
// bound.
func (s *CloudTrailService) purgeAgedQueries() {
	if s.storageManager == nil {
		return
	}
	cutoff := time.Now().UTC().Add(-cloudtrailstore.ListQueriesWindow)
	for _, region := range s.storageManager.GetActiveRegions() {
		store, err := s.GetStoreForRegion(region)
		if err != nil {
			logs.Error("cloudtrail query purge: failed to resolve store",
				logs.String("region", region), logs.Err(err))
			continue
		}
		purged, err := store.PurgeQueriesBefore(cutoff)
		if err != nil {
			logs.Error("cloudtrail query purge failed",
				logs.String("region", region), logs.Err(err))
			continue
		}
		if purged > 0 {
			logs.Info("cloudtrail query purge",
				logs.String("region", region), logs.Int("queries", purged))
		}
	}
}

// enforceEDSRetentionForStore runs one store's retention sweep: live stores
// keep only their retention window of events; deleted stores past the
// restore window are removed entirely, and those still inside it keep their
// events for a possible restore.
func (s *CloudTrailService) enforceEDSRetentionForStore(store cloudtrailstore.CloudTrailStoreInterface) {
	now := time.Now().UTC()
	edsList, err := store.ListEventDataStoresAll()
	if err != nil {
		logs.Error("cloudtrail EDS retention: failed to list event data stores", logs.Err(err))
		return
	}
	for _, eds := range edsList {
		if eds.Status == "PENDING_DELETION" {
			s.hardDeleteIfRestoreWindowClosed(store, eds, now)
			continue
		}
		cutoff := now.Add(-time.Duration(eds.RetentionPeriod) * 24 * time.Hour)
		removed, err := store.PurgeEDSEventsBefore(eds.EventDataStoreID, cutoff, 0)
		if err != nil {
			logs.Error("cloudtrail EDS retention purge failed",
				logs.String("eventDataStoreId", eds.EventDataStoreID), logs.Err(err))
			continue
		}
		if removed > 0 {
			logs.Info("cloudtrail EDS retention purge",
				logs.String("eventDataStoreId", eds.EventDataStoreID),
				logs.Int("events", removed))
		}
	}
}

// hardDeleteIfRestoreWindowClosed removes a PENDING_DELETION event data
// store whose restore-window wait has closed (the documented seven days;
// TEST_MODE shortens it). The record goes first, under the store mutex
// with the window re-checked in the guard, so a restore that lands first
// flips the record to ENABLED and the delete refuses; only a successful
// record delete drops the event bucket afterwards. A crash between the
// two leaves an orphan bucket no record points at — inert residue, never
// a resurrected store.
func (s *CloudTrailService) hardDeleteIfRestoreWindowClosed(store cloudtrailstore.CloudTrailStoreInterface, eds *cloudtrailstore.EventDataStore, now time.Time) {
	if eds.DeletedTimestamp == nil ||
		now.Sub(*eds.DeletedTimestamp) < edsRestoreWindow() {
		return
	}
	deleted, err := store.DeleteEventDataStoreIf(eds.EventDataStoreID, func(current *cloudtrailstore.EventDataStore) error {
		if current.Status != "PENDING_DELETION" || current.DeletedTimestamp == nil ||
			now.Sub(*current.DeletedTimestamp) < edsRestoreWindow() {
			return cloudtrailstore.ErrUnchanged
		}
		return nil
	})
	if err != nil {
		logs.Error("cloudtrail EDS hard delete failed",
			logs.String("eventDataStoreId", eds.EventDataStoreID), logs.Err(err))
		return
	}
	if !deleted {
		return
	}
	if err := store.DropEDSEvents(eds.EventDataStoreID); err != nil {
		logs.Error("cloudtrail EDS hard delete: failed to drop events",
			logs.String("eventDataStoreId", eds.EventDataStoreID), logs.Err(err))
		return
	}
	logs.Info("cloudtrail EDS hard delete complete",
		logs.String("eventDataStoreId", eds.EventDataStoreID))
}

// bootPurgeEDSEvents bounds every event data store's events at the TEST_MODE
// boot bound — pending-deletion stores included, whose events otherwise
// outlive the seven-day restore window a regression cadence never reaches.
// It returns the number of events removed.
func (s *CloudTrailService) bootPurgeEDSEvents(cutoff time.Time) int {
	if s.storageManager == nil {
		return 0
	}
	total := 0
	for _, region := range s.storageManager.GetActiveRegions() {
		store, err := s.GetStoreForRegion(region)
		if err != nil {
			logs.Error("cloudtrail TEST_MODE boot purge: failed to resolve store",
				logs.String("region", region), logs.Err(err))
			continue
		}
		edsList, err := store.ListEventDataStoresAll()
		if err != nil {
			logs.Error("cloudtrail TEST_MODE boot purge: failed to list event data stores",
				logs.String("region", region), logs.Err(err))
			continue
		}
		for _, eds := range edsList {
			removed, err := store.PurgeEDSEventsBefore(eds.EventDataStoreID, cutoff, 0)
			if err != nil {
				logs.Error("cloudtrail TEST_MODE boot purge failed",
					logs.String("region", region),
					logs.String("eventDataStoreId", eds.EventDataStoreID), logs.Err(err))
				continue
			}
			total += removed
		}
	}
	return total
}

// Stop shuts down the event-history retention worker.
func (s *CloudTrailService) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
	s.wg.Wait()
}
