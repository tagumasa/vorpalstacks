package s3

import (
	"context"
	"errors"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// LifecycleWorker periodically scans buckets for objects whose lifecycle
// rules have taken effect — expiring objects and noncurrent versions,
// transitioning them between storage classes, removing expired delete
// markers, and aborting incomplete multipart uploads.
type LifecycleWorker struct {
	svc       *S3Service
	interval  time.Duration
	testMode  bool
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	startOnce sync.Once
}

// NewLifecycleWorker creates a new LifecycleWorker with a 5-minute default
// interval, tightened to 1 second in TEST_MODE so expiry sweeps are
// observable within a test run.
func NewLifecycleWorker(svc *S3Service) *LifecycleWorker {
	testMode := os.Getenv("TEST_MODE") == "true"
	interval := 5 * time.Minute
	if testMode {
		interval = 1 * time.Second
	}
	return &LifecycleWorker{
		svc:      svc,
		interval: interval,
		testMode: testMode,
	}
}

// Start launches the lifecycle enforcement goroutine.
func (w *LifecycleWorker) Start() {
	w.startOnce.Do(func() {
		w.ctx, w.cancel = context.WithCancel(context.Background())
		w.wg.Add(1)
		go w.run()
		logs.Info("s3: lifecycle worker started", logs.Any("interval", w.interval))
	})
}

// Close gracefully stops the lifecycle worker.
func (w *LifecycleWorker) Close() {
	if w.cancel != nil {
		w.cancel()
	}
	w.wg.Wait()
	logs.Info("s3: lifecycle worker stopped")
}

func (w *LifecycleWorker) run() {
	defer w.wg.Done()

	// The initial delay lets the server settle before the first scan; the
	// TEST_MODE branch shortens it alongside the tightened interval.
	initial := 30 * time.Second
	if w.testMode {
		initial = 5 * time.Second
	}
	timer := time.NewTimer(initial)
	defer timer.Stop()

	for {
		select {
		case <-w.ctx.Done():
			return
		case <-timer.C:
			w.enforceLifecycle()
			timer.Reset(w.interval)
		}
	}
}

// lifecycleAgeBound returns the raw age boundary for a rule period
// expressed in days: the newest creation time whose rounded action instant —
// which lifecycleDaysInstant computes — can already have arrived. It is a
// cheap candidate pre-filter for scans; the exact gates compare the rounded
// instants. The day unit and its TEST_MODE compression are the single
// definition in lifecycleDayUnit. Restore expiry is not compressed: the
// restore request itself fixes the copy's absolute expiry instant, which the
// sweep reads as stored.
func lifecycleAgeBound(days int, now time.Time) time.Time {
	return now.Add(-time.Duration(days) * lifecycleDayUnit)
}

// enforceLifecycle iterates over all buckets with lifecycle configuration
// and deletes expired objects. The scan goes through the same per-region
// stores the request plane uses; an ad-hoc store keyed by the account id
// would read a different, empty keyspace.
func (w *LifecycleWorker) enforceLifecycle() {
	w.svc.s3RegionStores(func(region string, bucketStore *s3store.BucketStore, objectStore *s3store.ObjectStore) {
		buckets, err := bucketStore.List()
		if err != nil {
			logs.Warn("s3: lifecycle worker failed to list buckets", logs.String("region", region), logs.Err(err))
			return
		}

		// The restored-copy expiry sweep runs regardless of lifecycle
		// configuration, driven by the restore index.
		w.expireRestoredCopies(region, objectStore)

		for _, bucket := range buckets {
			if bucket.LifecycleConfiguration == nil || len(bucket.LifecycleConfiguration.Rules) == 0 {
				continue
			}

			w.processBucketLifecycle(bucket, objectStore)
		}
	})
}

// expireRestoredCopies clears the restore state of objects whose temporary
// restored copy has expired and publishes the ObjectRestore:Delete
// notification. It is driven by the store's restore index, so buckets
// without active restores cost a single index scan instead of a full
// version listing. Index entries whose object no longer exists (a deleted
// version or bucket) are dropped silently.
func (w *LifecycleWorker) expireRestoredCopies(region string, objectStore *s3store.ObjectStore) {
	now := time.Now()
	entries, err := objectStore.ActiveRestores()
	if err != nil {
		logs.Warn("s3: restore index scan failed", logs.Err(err))
		return
	}

	for _, entry := range entries {
		if now.Before(entry.Expiry) {
			continue
		}
		obj, err := objectStore.HeadWithVersion(context.Background(), entry.Bucket, entry.Key, entry.VersionID)
		if err != nil {
			if errors.Is(err, s3store.ErrObjectNotFound) {
				// The object version is gone; forget the orphaned entry.
				if err := objectStore.SetRestoreState(entry.Bucket, entry.Key, entry.VersionID, nil); err != nil {
					logs.Warn("s3: restore index cleanup failed", logs.String("bucket", entry.Bucket), logs.String("key", entry.Key), logs.Err(err))
				}
				continue
			}
			logs.Warn("s3: restore expiry head failed", logs.String("bucket", entry.Bucket), logs.String("key", entry.Key), logs.Err(err))
			continue
		}
		if err := objectStore.SetRestoreState(entry.Bucket, entry.Key, entry.VersionID, nil); err != nil {
			logs.Warn("s3: restore expiry revert failed", logs.String("bucket", entry.Bucket), logs.String("key", entry.Key), logs.Err(err))
			continue
		}
		reqCtx := request.NewRequestContext(context.Background(), w.svc.storageManager, w.svc.accountID, region)
		w.svc.publishObjectNotification(reqCtx, reqCtx, entry.Bucket, entry.Key, obj.Size, obj.VersionID, obj.ETag, eventbus.S3ObjectRestoreDelete)
	}
}

// processBucketLifecycle evaluates lifecycle rules for a single bucket
// and deletes expired objects.
func (w *LifecycleWorker) processBucketLifecycle(bucket *s3store.Bucket, objectStore *s3store.ObjectStore) {
	now := time.Now()

	for _, rule := range bucket.LifecycleConfiguration.Rules {
		if rule.Status != "Enabled" {
			continue
		}

		var filter *s3store.LifecycleRuleFilter
		if rule.Filter != nil {
			filter = rule.Filter
		}

		if rule.Expiration != nil {
			days := 0
			if rule.Expiration.Days != nil && *rule.Expiration.Days > 0 {
				days = int(*rule.Expiration.Days)
			}
			if days > 0 {
				w.expireObjectsByAge(bucket.Name, objectStore, filter, days, now)
			}
			if rule.Expiration.Date != nil && now.After(*rule.Expiration.Date) {
				w.expireObjectsAll(bucket.Name, objectStore, filter)
			}
			// The expired-delete-marker action removes latest markers that
			// are the object's only version; the expiry sweeps above never
			// touch delete markers.
			if rule.Expiration.ExpiredObjectDeleteMarker != nil && *rule.Expiration.ExpiredObjectDeleteMarker {
				w.removeExpiredDeleteMarkers(bucket.Name, objectStore, filter)
			}
		}

		if len(rule.Transitions) > 0 {
			w.transitionObjects(bucket.Name, objectStore, filter, rule.Transitions, now)
		}

		if rule.AbortIncompleteMultipartUpload != nil && rule.AbortIncompleteMultipartUpload.DaysAfterInitiation != nil && *rule.AbortIncompleteMultipartUpload.DaysAfterInitiation > 0 {
			w.abortIncompleteUploads(bucket.Name, objectStore, int(*rule.AbortIncompleteMultipartUpload.DaysAfterInitiation), now)
		}

		if rule.NoncurrentVersionExpiration != nil && rule.NoncurrentVersionExpiration.NoncurrentDays != nil && *rule.NoncurrentVersionExpiration.NoncurrentDays > 0 {
			w.expireNoncurrentVersions(bucket.Name, objectStore, filter, rule.NoncurrentVersionExpiration, now)
		}

		if len(rule.NoncurrentVersionTransitions) > 0 {
			w.transitionNoncurrentVersions(bucket.Name, objectStore, filter, rule.NoncurrentVersionTransitions, now)
		}
	}
}

// isProtectedByObjectLock checks whether an object is protected from
// deletion by a legal hold or an active retention period. Lifecycle
// expiration must respect Object Lock just like any other delete operation.
func isProtectedByObjectLock(obj *s3store.Object, now time.Time) bool {
	if obj.ObjectLockLegalHold != nil && obj.ObjectLockLegalHold.Status == s3store.ObjectLockLegalHoldOn {
		return true
	}
	if obj.ObjectLockRetention != nil {
		if obj.ObjectLockRetention.RetainUntilDate.After(now) {
			return true
		}
	}
	return false
}

// expireObjectsByAge deletes objects whose Days-based expiry has arrived.
func (w *LifecycleWorker) expireObjectsByAge(bucketName string, objectStore *s3store.ObjectStore, filter *s3store.LifecycleRuleFilter, days int, now time.Time) {
	w.expireObjects(bucketName, objectStore, filter, &days, now)
}

// expireObjectsAll deletes all objects matching the filter (for Date-based expiration).
func (w *LifecycleWorker) expireObjectsAll(bucketName string, objectStore *s3store.ObjectStore, filter *s3store.LifecycleRuleFilter) {
	w.expireObjects(bucketName, objectStore, filter, nil, time.Now())
}

// expireObjects deletes every object matching the filter: with a non-nil
// days only objects whose rounded Days instant has arrived, with a nil days
// every match (Date-based expiration). The raw age bound is only the
// candidate pre-filter — the gate is the rounded midnight instant, so
// nothing expires before the x-amz-expiration promise. Delete markers are
// skipped (nothing to expire), object-lock-protected objects are never
// touched, and the listing paginates until exhausted.
func (w *LifecycleWorker) expireObjects(bucketName string, objectStore *s3store.ObjectStore, filter *s3store.LifecycleRuleFilter, days *int, now time.Time) {
	prefix := filterPrefix(filter)
	cutoff := time.Time{}
	if days != nil {
		cutoff = lifecycleAgeBound(*days, now)
	}
	marker := ""
	for {
		result, err := objectStore.List(bucketName, prefix, "", marker, s3MaxKeys)
		if err != nil {
			logs.Warn("s3: lifecycle list failed", logs.String("bucket", bucketName), logs.Err(err))
			return
		}

		for _, obj := range result.Objects {
			if obj.IsDeleteMarker {
				continue
			}
			if !matchesLifecycleFilter(obj, filter) {
				continue
			}
			if isProtectedByObjectLock(obj, now) {
				continue
			}
			if !lifecycleReplicationEligible(obj) {
				continue
			}
			if days != nil {
				if !obj.LastModified.Before(cutoff) {
					continue
				}
				if lifecycleDaysInstant(obj.LastModified, int32(*days)).After(now) {
					continue
				}
			}
			if err := objectStore.Delete(context.Background(), bucketName, obj.Key); err != nil {
				logs.Warn("s3: lifecycle delete failed", logs.String("bucket", bucketName), logs.String("key", obj.Key), logs.Err(err))
			}
		}

		if !result.IsTruncated {
			break
		}
		marker = result.NextMarker
	}
}

// abortIncompleteUploads aborts multipart uploads older than the specified number of days.
func (w *LifecycleWorker) abortIncompleteUploads(bucketName string, objectStore *s3store.ObjectStore, daysAfterInit int, now time.Time) {
	cutoff := lifecycleAgeBound(daysAfterInit, now)
	keyMarker := ""
	uploadIdMarker := ""
	for {
		result, err := objectStore.ListMultipartUploads(bucketName, "", keyMarker, uploadIdMarker, s3MaxUploads)
		if err != nil {
			logs.Warn("s3: lifecycle multipart list failed", logs.String("bucket", bucketName), logs.Err(err))
			return
		}

		for _, upload := range result.Uploads {
			if upload.Initiated.Before(cutoff) {
				if err := objectStore.AbortMultipartUpload(context.Background(), bucketName, upload.Key, upload.UploadID); err != nil {
					logs.Warn("s3: lifecycle abort multipart failed", logs.String("bucket", bucketName), logs.String("key", upload.Key), logs.Err(err))
				}
			}
		}

		if !result.IsTruncated {
			break
		}
		keyMarker = result.NextKeyMarker
		uploadIdMarker = result.NextUploadIDMarker
	}
}

// expireNoncurrentVersions permanently deletes noncurrent versions past the
// rule window. A version expires only when BOTH gates are exceeded: the
// NoncurrentDays window, counted — per the AWS noncurrent-day calculation —
// from when the successor version was created, and, when configured, the
// NewerNoncurrentVersions retention count. Delete markers are noncurrent
// versions like any other and expire under the same rule.
func (w *LifecycleWorker) expireNoncurrentVersions(bucketName string, objectStore *s3store.ObjectStore, filter *s3store.LifecycleRuleFilter, nve *s3store.NoncurrentVersionExpiration, now time.Time) {
	retain := int32(0)
	if nve.NewerNoncurrentVersions != nil {
		retain = *nve.NewerNoncurrentVersions
	}
	w.forEachVersionGroup(bucketName, objectStore, filterPrefix(filter), func(key string, versions []*s3store.Object) {
		for i, obj := range versions {
			if i == 0 {
				continue // the current version is never a noncurrent expiry candidate
			}
			// versions[0] is current, so the versions newer than obj are
			// versions[1:i] plus versions[0]: obj has i-1 newer noncurrent
			// versions. AWS requires the retention count to be exceeded,
			// not merely met.
			if int32(i-1) <= retain {
				continue
			}
			if !matchesLifecycleFilter(obj, filter) {
				continue
			}
			if isProtectedByObjectLock(obj, now) {
				continue
			}
			if !lifecycleReplicationEligible(obj) {
				continue
			}
			// The version became noncurrent when its successor was created:
			// the immediately newer version's LastModified, and the
			// noncurrent age rounds up to the next midnight UTC exactly as
			// the expiration projection does.
			if lifecycleDaysInstant(versions[i-1].LastModified, *nve.NoncurrentDays).After(now) {
				continue
			}
			if _, err := objectStore.DeleteWithVersion(context.Background(), bucketName, obj.Key, obj.VersionID); err != nil {
				logs.Warn("s3: lifecycle version delete failed", logs.String("bucket", bucketName), logs.String("key", obj.Key), logs.String("version", obj.VersionID), logs.Err(err))
			}
		}
	})
}

// transitionObjects moves current object versions past the rule window to
// the configured storage class. The entries replay in rule order from the
// object's present class, so a move the storage-class waterfall forbids is
// skipped rather than regressing the object to an earlier class.
func (w *LifecycleWorker) transitionObjects(bucketName string, objectStore *s3store.ObjectStore, filter *s3store.LifecycleRuleFilter, transitions []s3store.LifecycleTransition, now time.Time) {
	specs := make([]transitionSpec, 0, len(transitions))
	for _, t := range transitions {
		specs = append(specs, transitionSpec{days: t.Days, date: t.Date, class: t.StorageClass})
	}
	prefix := filterPrefix(filter)
	marker := ""
	for {
		result, err := objectStore.List(bucketName, prefix, "", marker, s3MaxKeys)
		if err != nil {
			logs.Warn("s3: lifecycle list failed", logs.String("bucket", bucketName), logs.Err(err))
			return
		}

		for _, obj := range result.Objects {
			if obj.IsDeleteMarker {
				continue
			}
			if !matchesLifecycleFilter(obj, filter) {
				continue
			}
			if !transitionSizeEligible(obj, filter) {
				continue
			}
			if !lifecycleReplicationEligible(obj) {
				continue
			}
			target := resolveTransitionTarget(obj.StorageClass, specs, obj.LastModified, now)
			if target == obj.StorageClass {
				continue
			}
			if err := objectStore.SetStorageClass(bucketName, obj.Key, obj.VersionID, target); err != nil {
				logs.Warn("s3: lifecycle transition failed", logs.String("bucket", bucketName), logs.String("key", obj.Key), logs.String("version", obj.VersionID), logs.Err(err))
			}
		}

		if !result.IsTruncated {
			break
		}
		marker = result.NextMarker
	}
}

// transitionNoncurrentVersions moves noncurrent versions past the rule
// window to the configured storage class. The gates mirror
// NoncurrentVersionExpiration: for each entry both the NoncurrentDays
// window, counted from when the successor version was created, and, when
// configured, the NewerNoncurrentVersions retention count must be exceeded.
// Entries replay in rule order from the version's present class, so a move
// the storage-class waterfall forbids is skipped rather than regressing the
// version to an earlier class. Noncurrent delete markers carry no storage
// class and are never transition targets.
func (w *LifecycleWorker) transitionNoncurrentVersions(bucketName string, objectStore *s3store.ObjectStore, filter *s3store.LifecycleRuleFilter, transitions []s3store.NoncurrentVersionTransition, now time.Time) {
	w.forEachVersionGroup(bucketName, objectStore, filterPrefix(filter), func(key string, versions []*s3store.Object) {
		for i, obj := range versions {
			if i == 0 || obj.IsDeleteMarker {
				continue
			}
			if !matchesLifecycleFilter(obj, filter) {
				continue
			}
			if !transitionSizeEligible(obj, filter) {
				continue
			}
			if !lifecycleReplicationEligible(obj) {
				continue
			}
			// The version became noncurrent when its successor was created:
			// the immediately newer version's LastModified.
			becameNoncurrent := versions[i-1].LastModified
			// versions[0] is current, so the versions newer than obj are
			// versions[1:i] plus versions[0]: obj has i-1 newer noncurrent
			// versions.
			newerNoncurrent := int32(i - 1)
			class := obj.StorageClass
			for _, t := range transitions {
				if t.NoncurrentDays == nil || *t.NoncurrentDays <= 0 {
					continue
				}
				if t.NewerNoncurrentVersions != nil && newerNoncurrent <= *t.NewerNoncurrentVersions {
					continue
				}
				if lifecycleDaysInstant(becameNoncurrent, *t.NoncurrentDays).After(now) {
					continue
				}
				if lifecycleTransitionAllowed(class, t.StorageClass) {
					class = t.StorageClass
				}
			}
			if class == obj.StorageClass {
				continue
			}
			if err := objectStore.SetStorageClass(bucketName, obj.Key, obj.VersionID, class); err != nil {
				logs.Warn("s3: lifecycle version transition failed", logs.String("bucket", bucketName), logs.String("key", obj.Key), logs.String("version", obj.VersionID), logs.Err(err))
			}
		}
	})
}

// removeExpiredDeleteMarkers removes expired object delete markers: a
// latest delete marker with no noncurrent versions beneath it. Such a marker
// is the object's only version, so removing it removes the object entirely.
func (w *LifecycleWorker) removeExpiredDeleteMarkers(bucketName string, objectStore *s3store.ObjectStore, filter *s3store.LifecycleRuleFilter) {
	w.forEachVersionGroup(bucketName, objectStore, filterPrefix(filter), func(key string, versions []*s3store.Object) {
		if len(versions) != 1 || !versions[0].IsDeleteMarker {
			return
		}
		if !matchesLifecycleFilter(versions[0], filter) {
			return
		}
		if _, err := objectStore.DeleteWithVersion(context.Background(), bucketName, key, versions[0].VersionID); err != nil {
			logs.Warn("s3: lifecycle delete-marker removal failed", logs.String("bucket", bucketName), logs.String("key", key), logs.String("version", versions[0].VersionID), logs.Err(err))
		}
	})
}

// forEachVersionGroup enumerates every object version under a prefix,
// grouped by key and ordered newest-first within each group, and invokes fn
// once per key. A key's group may span listing pages; the continuation
// markers resume exactly after the last delivered record, so buffering the
// current group across pages neither drops nor repeats records.
func (w *LifecycleWorker) forEachVersionGroup(bucketName string, objectStore *s3store.ObjectStore, prefix string, fn func(key string, versions []*s3store.Object)) {
	keyMarker := ""
	versionIdMarker := ""
	curKey := ""
	var group []*s3store.Object
	flush := func() {
		if curKey == "" || len(group) == 0 {
			return
		}
		sort.Slice(group, func(i, j int) bool {
			if !group[i].LastModified.Equal(group[j].LastModified) {
				return group[i].LastModified.After(group[j].LastModified)
			}
			// Version IDs are UUIDs with no intrinsic order; the ID breaks
			// ties so identical timestamps from one write burst still yield
			// the same total order re-latesting uses.
			return group[i].VersionID > group[j].VersionID
		})
		fn(curKey, group)
		curKey = ""
		group = nil
	}
	for {
		result, err := objectStore.ListObjectVersions(bucketName, prefix, "", keyMarker, versionIdMarker, s3MaxKeys)
		if err != nil {
			logs.Warn("s3: lifecycle version list failed", logs.String("bucket", bucketName), logs.Err(err))
			return
		}
		for _, obj := range result.Objects {
			if obj.Key != curKey {
				flush()
				curKey = obj.Key
			}
			group = append(group, obj)
		}
		if !result.IsTruncated {
			break
		}
		keyMarker = result.NextVersionKeyMarker
		versionIdMarker = result.NextVersionIDMarker
	}
	flush()
}

// transitionSpec normalises a Transition entry for the replay that decides
// a current version's final storage class.
type transitionSpec struct {
	days  *int32
	date  *time.Time
	class s3store.ObjectStorageClass
}

// resolveTransitionTarget replays a rule's transition entries in order from
// the current version's present class and returns the class it ends in. An
// entry is due when its Days window — counted from since, the version's
// LastModified, and rounded up to the next midnight UTC like every other
// Days window — or its Date has passed; a due entry moves the version only
// when the storage-class waterfall allows the destination from the class
// reached so far.
func resolveTransitionTarget(current s3store.ObjectStorageClass, specs []transitionSpec, since time.Time, now time.Time) s3store.ObjectStorageClass {
	class := current
	for _, s := range specs {
		due := (s.days != nil && *s.days > 0 && !lifecycleDaysInstant(since, *s.days).After(now)) ||
			(s.date != nil && now.After(*s.date))
		if due && lifecycleTransitionAllowed(class, s.class) {
			class = s.class
		}
	}
	return class
}

// lifecycleTransitionTargets lists the storage classes a lifecycle
// transition may move an object to, keyed by the object's current class —
// the supported-transitions waterfall of the S3 User Guide. Reduced
// redundancy has no row of its own: it is a Standard-tier redundancy
// variant, so the Standard row applies. Intelligent-Tiering destinations
// depend on the access tier in AWS; this platform does not model tiers, so
// its Frequent/Infrequent Access row applies.
var lifecycleTransitionTargets = map[s3store.ObjectStorageClass][]s3store.ObjectStorageClass{
	s3store.StorageClassStandard: {
		s3store.StorageClassStandardIA,
		s3store.StorageClassIntelligentTiering,
		s3store.StorageClassOneZoneIA,
		s3store.StorageClassGlacierIR,
		s3store.StorageClassGlacier,
		s3store.StorageClassDeepArchive,
	},
	s3store.StorageClassReducedRedundancy: {
		s3store.StorageClassStandardIA,
		s3store.StorageClassIntelligentTiering,
		s3store.StorageClassOneZoneIA,
		s3store.StorageClassGlacierIR,
		s3store.StorageClassGlacier,
		s3store.StorageClassDeepArchive,
	},
	s3store.StorageClassStandardIA: {
		s3store.StorageClassIntelligentTiering,
		s3store.StorageClassOneZoneIA,
		s3store.StorageClassGlacierIR,
		s3store.StorageClassGlacier,
		s3store.StorageClassDeepArchive,
	},
	s3store.StorageClassIntelligentTiering: {
		s3store.StorageClassOneZoneIA,
		s3store.StorageClassGlacierIR,
		s3store.StorageClassGlacier,
		s3store.StorageClassDeepArchive,
	},
	s3store.StorageClassOneZoneIA: {
		s3store.StorageClassGlacier,
		s3store.StorageClassDeepArchive,
	},
	s3store.StorageClassGlacierIR: {
		s3store.StorageClassGlacier,
		s3store.StorageClassDeepArchive,
	},
	s3store.StorageClassGlacier: {
		s3store.StorageClassDeepArchive,
	},
	s3store.StorageClassDeepArchive: {},
}

// lifecycleTransitionAllowed reports whether the storage-class waterfall
// permits moving from one class to another.
func lifecycleTransitionAllowed(from, to s3store.ObjectStorageClass) bool {
	for _, target := range lifecycleTransitionTargets[from] {
		if target == to {
			return true
		}
	}
	return false
}

// transitionSizeEligible applies the AWS default transition minimum: objects
// smaller than 128 KiB do not transition to any storage class unless the
// rule's filter carries an explicit size bound, which then governs
// eligibility on its own.
func transitionSizeEligible(obj *s3store.Object, filter *s3store.LifecycleRuleFilter) bool {
	if filterHasSizeBound(filter) {
		return true
	}
	return obj.Size >= s3store.MinTransitionObjectSize
}

// filterHasSizeBound reports whether a lifecycle rule filter selects objects
// by size, directly or inside the And operator.
func filterHasSizeBound(filter *s3store.LifecycleRuleFilter) bool {
	if filter == nil {
		return false
	}
	if filter.ObjectSizeGreaterThan != nil || filter.ObjectSizeLessThan != nil {
		return true
	}
	return filter.And != nil && (filter.And.ObjectSizeGreaterThan != nil || filter.And.ObjectSizeLessThan != nil)
}

// lifecycleReplicationEligible enforces the AWS rule that S3 Lifecycle
// prevents expiration and transition actions on an object until its
// replication has succeeded: a Pending or Failed status blocks every
// lifecycle action.
func lifecycleReplicationEligible(obj *s3store.Object) bool {
	return obj.ReplicationStatus != s3store.ReplicationStatusPending &&
		obj.ReplicationStatus != s3store.ReplicationStatusFailed
}

// filterPrefix extracts the effective prefix from a lifecycle rule filter
// for efficient List queries. When the And operator is present with a
// non-empty prefix, that prefix takes precedence.
func filterPrefix(filter *s3store.LifecycleRuleFilter) string {
	if filter == nil {
		return ""
	}
	if filter.And != nil && filter.And.Prefix != "" {
		return filter.And.Prefix
	}
	return filter.Prefix
}

// matchesLifecycleFilter evaluates all filter criteria (Prefix, Tag,
// ObjectSizeGreaterThan, ObjectSizeLessThan, and And.*) against an object.
// Returns true when the object matches all applicable criteria.
func matchesLifecycleFilter(obj *s3store.Object, filter *s3store.LifecycleRuleFilter) bool {
	if filter == nil {
		return true
	}

	if filter.And != nil {
		and := filter.And
		if and.Prefix != "" && !strings.HasPrefix(obj.Key, and.Prefix) {
			return false
		}
		for _, tag := range and.Tags {
			if !objectHasTag(obj.Tags, tag.Key, tag.Value) {
				return false
			}
		}
		if and.ObjectSizeGreaterThan != nil && obj.Size <= *and.ObjectSizeGreaterThan {
			return false
		}
		if and.ObjectSizeLessThan != nil && obj.Size >= *and.ObjectSizeLessThan {
			return false
		}
		return true
	}

	if filter.Prefix != "" && !strings.HasPrefix(obj.Key, filter.Prefix) {
		return false
	}
	if filter.Tag != nil {
		if !objectHasTag(obj.Tags, filter.Tag.Key, filter.Tag.Value) {
			return false
		}
	}
	if filter.ObjectSizeGreaterThan != nil && obj.Size <= *filter.ObjectSizeGreaterThan {
		return false
	}
	if filter.ObjectSizeLessThan != nil && obj.Size >= *filter.ObjectSizeLessThan {
		return false
	}
	return true
}
