package s3

import (
	"context"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// LifecycleWorker periodically scans buckets for objects that have expired
// according to their LifecycleConfiguration and deletes them.
type LifecycleWorker struct {
	svc            *S3Service
	interval       time.Duration
	storageManager *storage.RegionStorageManager
	accountID      string
	ctx            context.Context
	cancel         context.CancelFunc
	wg             sync.WaitGroup
	startOnce      sync.Once
}

// NewLifecycleWorker creates a new LifecycleWorker with a 5-minute default interval.
func NewLifecycleWorker(svc *S3Service, sm *storage.RegionStorageManager, accountID string) *LifecycleWorker {
	return &LifecycleWorker{
		svc:            svc,
		interval:       5 * time.Minute,
		storageManager: sm,
		accountID:      accountID,
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

	timer := time.NewTimer(30 * time.Second)
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

// enforceLifecycle iterates over all buckets with lifecycle configuration
// and deletes expired objects.
func (w *LifecycleWorker) enforceLifecycle() {
	regionStore, err := w.storageManager.GetStorage(w.accountID)
	if err != nil {
		logs.Warn("s3: lifecycle worker failed to get storage", logs.Err(err))
		return
	}

	bucketStore := s3store.NewBucketStore(regionStore, w.accountID, w.accountID)
	buckets, err := bucketStore.List()
	if err != nil {
		logs.Warn("s3: lifecycle worker failed to list buckets", logs.Err(err))
		return
	}

	for _, bucket := range buckets {
		if bucket.LifecycleConfiguration == nil || len(bucket.LifecycleConfiguration.Rules) == 0 {
			continue
		}

		w.processBucketLifecycle(bucket)
	}
}

// processBucketLifecycle evaluates lifecycle rules for a single bucket
// and deletes expired objects.
func (w *LifecycleWorker) processBucketLifecycle(bucket *s3store.Bucket) {
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
				w.expireObjectsByAge(bucket.Name, filter, days, now)
			}
			if rule.Expiration.Date != nil && now.After(*rule.Expiration.Date) {
				w.expireObjectsAll(bucket.Name, filter)
			}
		}

		if rule.AbortIncompleteMultipartUpload != nil && rule.AbortIncompleteMultipartUpload.DaysAfterInitiation != nil && *rule.AbortIncompleteMultipartUpload.DaysAfterInitiation > 0 {
			w.abortIncompleteUploads(bucket.Name, int(*rule.AbortIncompleteMultipartUpload.DaysAfterInitiation), now)
		}

		if rule.NoncurrentVersionExpiration != nil && rule.NoncurrentVersionExpiration.NoncurrentDays != nil && *rule.NoncurrentVersionExpiration.NoncurrentDays > 0 {
			w.expireNoncurrentVersions(bucket.Name, filter, int(*rule.NoncurrentVersionExpiration.NoncurrentDays), now)
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

// expireObjectsByAge deletes objects older than the specified number of days.
func (w *LifecycleWorker) expireObjectsByAge(bucketName string, filter *s3store.LifecycleRuleFilter, days int, now time.Time) {
	regionStore, err := w.storageManager.GetStorage(w.accountID)
	if err != nil {
		return
	}

	objectStore, err := s3store.NewObjectStore(regionStore, w.svc.blobStore, s3store.NewBucketStore(regionStore, w.accountID, w.accountID), w.accountID, w.accountID)
	if err != nil {
		return
	}

	cutoff := now.AddDate(0, 0, -days)
	prefix := filterPrefix(filter)
	marker := ""
	for {
		result, err := objectStore.List(bucketName, prefix, "", marker, 1000)
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
			if obj.LastModified.Before(cutoff) {
				if err := objectStore.Delete(context.Background(), bucketName, obj.Key); err != nil {
					logs.Warn("s3: lifecycle delete failed", logs.String("bucket", bucketName), logs.String("key", obj.Key), logs.Err(err))
				}
			}
		}

		if !result.IsTruncated {
			break
		}
		marker = result.NextMarker
	}
}

// expireObjectsAll deletes all objects matching the filter (for Date-based expiration).
func (w *LifecycleWorker) expireObjectsAll(bucketName string, filter *s3store.LifecycleRuleFilter) {
	regionStore, err := w.storageManager.GetStorage(w.accountID)
	if err != nil {
		return
	}

	objectStore, err := s3store.NewObjectStore(regionStore, w.svc.blobStore, s3store.NewBucketStore(regionStore, w.accountID, w.accountID), w.accountID, w.accountID)
	if err != nil {
		return
	}

	prefix := filterPrefix(filter)
	marker := ""
	for {
		result, err := objectStore.List(bucketName, prefix, "", marker, 1000)
		if err != nil {
			logs.Warn("s3: lifecycle list (date) failed", logs.String("bucket", bucketName), logs.Err(err))
			return
		}

		for _, obj := range result.Objects {
			if obj.IsDeleteMarker {
				continue
			}
			if !matchesLifecycleFilter(obj, filter) {
				continue
			}
			if isProtectedByObjectLock(obj, time.Now()) {
				continue
			}
			if err := objectStore.Delete(context.Background(), bucketName, obj.Key); err != nil {
				logs.Warn("s3: lifecycle delete (date) failed", logs.String("bucket", bucketName), logs.String("key", obj.Key), logs.Err(err))
			}
		}

		if !result.IsTruncated {
			break
		}
		marker = result.NextMarker
	}
}

// abortIncompleteUploads aborts multipart uploads older than the specified number of days.
func (w *LifecycleWorker) abortIncompleteUploads(bucketName string, daysAfterInit int, now time.Time) {
	regionStore, err := w.storageManager.GetStorage(w.accountID)
	if err != nil {
		return
	}

	objectStore, err := s3store.NewObjectStore(regionStore, w.svc.blobStore, s3store.NewBucketStore(regionStore, w.accountID, w.accountID), w.accountID, w.accountID)
	if err != nil {
		return
	}

	cutoff := now.AddDate(0, 0, -daysAfterInit)
	keyMarker := ""
	uploadIdMarker := ""
	for {
		result, err := objectStore.ListMultipartUploads(bucketName, "", keyMarker, uploadIdMarker, 1000)
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

// expireNoncurrentVersions deletes non-current (old) versions older than the
// specified number of days.
func (w *LifecycleWorker) expireNoncurrentVersions(bucketName string, filter *s3store.LifecycleRuleFilter, noncurrentDays int, now time.Time) {
	regionStore, err := w.storageManager.GetStorage(w.accountID)
	if err != nil {
		return
	}

	objectStore, err := s3store.NewObjectStore(regionStore, w.svc.blobStore, s3store.NewBucketStore(regionStore, w.accountID, w.accountID), w.accountID, w.accountID)
	if err != nil {
		return
	}

	cutoff := now.AddDate(0, 0, -noncurrentDays)
	prefix := filterPrefix(filter)
	keyMarker := ""
	versionIdMarker := ""
	for {
		result, err := objectStore.ListObjectVersions(bucketName, prefix, "", keyMarker, versionIdMarker, 1000)
		if err != nil {
			logs.Warn("s3: lifecycle version list failed", logs.String("bucket", bucketName), logs.Err(err))
			return
		}

		for _, obj := range result.Objects {
			if obj.IsLatest || obj.IsDeleteMarker {
				continue
			}
			if !matchesLifecycleFilter(obj, filter) {
				continue
			}
			if isProtectedByObjectLock(obj, now) {
				continue
			}
			if obj.LastModified.Before(cutoff) {
				if _, err := objectStore.DeleteWithVersion(context.Background(), bucketName, obj.Key, obj.VersionID); err != nil {
					logs.Warn("s3: lifecycle version delete failed", logs.String("bucket", bucketName), logs.String("key", obj.Key), logs.String("version", obj.VersionID), logs.Err(err))
				}
			}
		}

		if !result.IsTruncated {
			break
		}
		keyMarker = result.NextVersionKeyMarker
		versionIdMarker = result.NextVersionIDMarker
	}
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
