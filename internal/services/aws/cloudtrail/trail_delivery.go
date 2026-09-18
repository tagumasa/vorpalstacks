package cloudtrail

import (
	"bytes"
	"compress/gzip"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/common/invokers"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
	"vorpalstacks/internal/utils/aws/arn"
)

// Trail log-file delivery. A logging trail periodically aggregates the
// events recorded since its last delivery into gzip JSON log files under
// [prefix/]AWSLogs/<account>/CloudTrail/<region>/<Y>/<M>/<D>/ named
// <account>_CloudTrail_<region>_<YYYYMMDDTHHmmZ>_<UniqueString>.json.gz;
// validation-enabled trails additionally deliver an hourly digest file
// under CloudTrail-Digest/ chaining every delivered log file and signed
// with the trail's validation key (the signature travels as S3 object
// metadata, x-amz-meta-signature). Delivery outcomes are recorded on the
// trail and reported by GetTrailStatus — never synthesised.

const (
	// logFlushInterval is how often a logging trail's events are
	// aggregated into log files. CloudTrail "typically delivers logs
	// within an average of about 5 minutes of an API call".
	logFlushInterval = 5 * time.Minute
	// digestInterval is the digest cadence: "Each digest file contains
	// the names of the log files that were delivered to your Amazon S3
	// bucket during the last hour", and a digest is delivered even when
	// no API activity occurred.
	digestInterval = time.Hour
	// testModeLogFlushInterval and testModeDigestInterval shorten the
	// cadences in TEST_MODE so regression suites observe deliveries within
	// a bounded wait instead of minutes.
	testModeLogFlushInterval = time.Second
	testModeDigestInterval   = 3 * time.Second
	// logFileUniqueStringLength is the documented "16-character
	// UniqueString component of the log file name ... there to prevent
	// overwriting of files. It has no meaning".
	logFileUniqueStringLength = 16
	// trailDeliveryPageSize is the store page size while walking a
	// delivery window to exhaustion.
	trailDeliveryPageSize = 500
	// MaxTrailLogFileEvents bounds the events carried by one delivered
	// log file. This is a platform operational bound — it caps a delivery
	// window's in-memory footprint at one file's events plus a page — not
	// a documented AWS quota; AWS itself delivers more than one log file
	// per interval when the volume is high, so a window split into
	// several files stays compatible.
	MaxTrailLogFileEvents = 10000
)

// trailDeliveryIntervals resolves the effective cadences: the TEST_MODE
// values when the server runs in test mode, otherwise the documented ones.
func trailDeliveryIntervals() (logFlush, digest time.Duration) {
	if os.Getenv("TEST_MODE") == "true" {
		return testModeLogFlushInterval, testModeDigestInterval
	}
	return logFlushInterval, digestInterval
}

// StartTrailDeliveryWorker starts the periodic trail delivery worker. It is
// started alongside the event-history purger and stops with the service.
func (s *CloudTrailService) StartTrailDeliveryWorker() {
	logFlush, _ := trailDeliveryIntervals()
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer func() {
			if r := recover(); r != nil {
				slog.Error("PANIC in cloudtrail trail delivery worker, restarting", "panic", r)
				s.StartTrailDeliveryWorker()
			}
		}()
		ticker := time.NewTicker(logFlush)
		defer ticker.Stop()

		for {
			select {
			case <-s.ctx.Done():
				return
			case <-ticker.C:
				s.deliverTrailsAllRegions()
			}
		}
	}()
}

// deliverTrailsAllRegions flushes every logging trail's pending windows in
// every open region storage. A trail record lives in its home region's
// store and all bookkeeping is written there; the events it delivers for a
// region come from that region's store — a multi-region trail ("trails
// that log events in all Regions", IsMultiRegionTrail) flushes every
// active region's window, a single-region trail only its home region.
// One goroutine per trail walks that trail's regions sequentially, so a
// single trail still never overlaps itself (its watermark and digest
// chain writes stay ordered) while different trails flush concurrently:
// one trail's slow leg (a flood-sized window's gzip and object writes, a
// stalled destination) must not hold back every other trail's tick —
// there is no cross-trail ordering contract. Each region's flush reloads
// the trail record from the home store: the walk's earlier regions have
// already moved the digest chain and the pending set onto that record,
// and a digest built from the listing's stale snapshot would claim the
// same predecessor as an already-delivered digest — a forked chain.
func (s *CloudTrailService) deliverTrailsAllRegions() {
	if s.storageManager == nil {
		return
	}
	activeRegions := s.storageManager.GetActiveRegions()
	stores := make([]cloudtrailstore.CloudTrailStoreInterface, 0, len(activeRegions))
	for _, region := range activeRegions {
		store, err := s.GetStoreForRegion(region)
		if err != nil {
			slog.Error("cloudtrail trail delivery: failed to resolve store", "region", region, "error", err)
			continue
		}
		stores = append(stores, store)
	}
	var wg sync.WaitGroup
	for _, homeStore := range stores {
		trails, err := listAllTrails(homeStore)
		if err != nil {
			slog.Error("cloudtrail trail delivery: failed to list trails", "region", homeStore.GetRegion(), "error", err)
			continue
		}
		for _, trail := range trails {
			if !trail.IsLogging {
				continue
			}
			wg.Add(1)
			go func(homeStore cloudtrailstore.CloudTrailStoreInterface, name string) {
				defer wg.Done()
				for _, eventStore := range stores {
					region := eventStore.GetRegion()
					trail, err := homeStore.GetTrail(name)
					if err != nil {
						// A trail deleted (or unreadable) mid-walk stops the
						// walk: no later region may deliver against state it
						// could not read.
						if !errors.Is(err, cloudtrailstore.ErrTrailNotFound) {
							slog.Error("cloudtrail trail delivery: failed to reload trail",
								"trail", name, "region", region, "error", err)
						}
						return
					}
					if region != trail.HomeRegion && !trail.IsMultiRegionTrail {
						continue
					}
					s.flushTrailGuarded(homeStore, eventStore, region, trail)
				}
			}(homeStore, trail.Name)
		}
	}
	wg.Wait()
}

// flushTrailGuarded runs one trail's flush for one region with panic
// containment: a panic in a single trail's delivery (a store or
// serialisation fault) must not kill the process every other trail shares,
// nor silently swallow the window — the outcome records a delivery error,
// holding the window for the tick's retry and surfacing through
// GetTrailStatus like any delivery failure. The outcome object belongs to
// this guard, so the panic path records the state the aborted flush had
// already accumulated: log files written to S3 before the panic join the
// pending digest set (their objects exist; dropping them would punch a
// chain hole) and a succeeded CloudWatch Logs leg keeps its watermark
// advance. Recording the outcome runs under its own guard: a fault class
// that panicked the flush can panic the recording path too.
func (s *CloudTrailService) flushTrailGuarded(homeStore, eventStore cloudtrailstore.CloudTrailStoreInterface, region string, trail *cloudtrailstore.Trail) {
	outcome := &trailDeliveryOutcome{}
	defer func() {
		if r := recover(); r != nil {
			slog.Error("cloudtrail trail delivery: flush panicked",
				"trail", trail.Name, "region", region, "panic", r)
			now := time.Now().UTC()
			if outcome.windowEnd.IsZero() {
				outcome.windowEnd = now
			}
			if outcome.now.IsZero() {
				outcome.now = now
			}
			if outcome.deliveryErr == "" {
				outcome.deliveryErr = fmt.Sprintf("delivery panic: %v", r)
			}
			if !outcome.recorded {
				func() {
					defer func() {
						if r2 := recover(); r2 != nil {
							slog.Error("cloudtrail trail delivery: recording the panic outcome panicked",
								"trail", trail.Name, "region", region, "panic", r2)
						}
					}()
					recordTrailDeliveryOutcome(homeStore, trail.Name, region, trail.LogFileValidationEnabled, *outcome)
				}()
			}
		}
	}()
	s.flushTrail(homeStore, eventStore, region, trail, outcome)
}

// trailDeliveryOutcome carries one flush's result into the trail record.
type trailDeliveryOutcome struct {
	windowEnd       time.Time
	logDelivered    bool
	delivered       []cloudtrailstore.DigestLogFile
	deliveredKeys   []string
	deliveryErr     string
	digestDelivered bool
	digestEntry     *cloudtrailstore.DigestLogFile
	digestErr       string
	cwlogsDelivered bool
	cwlogsErr       string
	// cwlogsWindowStart is the span the stream leg attempted, recorded so
	// a failed stream window pins the stream watermark to its own start —
	// everything before it stands delivered and the next flush re-reads
	// from there rather than falling back to the advanced S3 window.
	cwlogsWindowStart time.Time
	notified          bool
	notificationErr   string
	now               time.Time
	// recorded marks that the flush completed its own bookkeeping write;
	// the panic guard consults it so an outcome is never recorded twice.
	recorded bool
}

// flushTrail delivers one trail's pending window for one region: the log
// files when events match the trail's selectors, and the digest when the
// digest cadence elapsed on a validation-enabled trail. The trail record
// and every bookkeeping write live in the trail's home store (homeStore);
// the events delivered are the target region's, read from eventStore. All
// bookkeeping is written in one MutateTrail step so concurrent API updates
// cannot lose it. The outcome the flush accumulates is owned by the caller
// (flushTrailGuarded) so a mid-flush panic can still record it.
func (s *CloudTrailService) flushTrail(homeStore, eventStore cloudtrailstore.CloudTrailStoreInterface, region string, trail *cloudtrailstore.Trail, outcome *trailDeliveryOutcome) {
	if !trail.IsLogging {
		return
	}

	watermark := trailWindowStart(trail, region)
	windowEnd := time.Now().UTC()
	if !windowEnd.After(watermark) {
		return
	}

	// The window is consumed page-wise: the store never holds the whole
	// window in memory, and the matched events accumulate into one log
	// file's buffer, finalised and written whenever it reaches
	// MaxTrailLogFileEvents events — a flood-sized window lands as several
	// log files (AWS itself delivers more than one file per interval on
	// high volume) instead of one unbounded buffer and object.
	outcome.windowEnd = windowEnd
	outcome.now = time.Now().UTC()
	fileBuf := make([]*cloudtrailstore.Event, 0, trailDeliveryPageSize)
	finaliseFileBuf := func() {
		if len(fileBuf) == 0 {
			return
		}
		file, err := s.deliverLogFile(homeStore, trail, region, fileBuf, windowEnd)
		if err != nil {
			if outcome.deliveryErr == "" {
				outcome.deliveryErr = err.Error()
				// Warn, not Error: the outcome is recorded on the trail
				// (GetTrailStatus reports it) and the tick repeats, so
				// error-level spam would add nothing.
				slog.Warn("cloudtrail trail log-file delivery failed",
					"trail", trail.Name, "region", region, "error", err)
			}
		} else {
			outcome.logDelivered = true
			outcome.deliveredKeys = append(outcome.deliveredKeys, file.ObjectKey)
			if trail.LogFileValidationEnabled {
				outcome.delivered = append(outcome.delivered, file)
			}
		}
		fileBuf = make([]*cloudtrailstore.Event, 0, trailDeliveryPageSize)
	}

	// The read starts at the watermark's whole second (the time index
	// buckets by second); the inclusive millisecond bound is the in-memory
	// filter. A failed page read records the failure and holds the window:
	// a partially read window must never be treated as complete, or the
	// watermark would advance past events that were never delivered. Files
	// already finalised before the failure stay delivered — the retry's
	// copies are the duplicates AWS itself tolerates, and the pending-file
	// stash keeps them in the digest chain.
	readFrom := watermark.Truncate(time.Second)
	token := ""
	for {
		query := cloudtrailstore.EventQuery{
			StartTime:  &readFrom,
			EndTime:    &windowEnd,
			MaxResults: trailDeliveryPageSize,
			NextToken:  token,
		}
		events, next, err := eventStore.LookupEvents(query)
		if err != nil {
			outcome.deliveryErr = fmt.Errorf("event read failed: %w", err).Error()
			slog.Warn("cloudtrail trail delivery: window read failed",
				"trail", trail.Name, "region", region, "error", err)
			recordTrailDeliveryOutcome(homeStore, trail.Name, region, trail.LogFileValidationEnabled, *outcome)
			outcome.recorded = true
			return
		}
		for _, e := range events {
			// The bound is INCLUSIVE at the watermark millisecond: stored
			// event times are millisecond-truncated, so an event in the
			// watermark's own millisecond equals the watermark and a strict
			// after-bound would lose it forever. The mirror-image risk — an
			// event duplicated across two windows' shared boundary
			// millisecond — is what AWS itself tolerates ("Although
			// uncommon, you may receive log files that contain one or more
			// duplicate events").
			if !e.EventTime.Before(watermark) && trailSelectorMatches(trail, e) {
				fileBuf = append(fileBuf, e)
				if len(fileBuf) >= MaxTrailLogFileEvents {
					finaliseFileBuf()
				}
			}
		}
		if next == "" {
			break
		}
		token = next
	}
	finaliseFileBuf()

	// CloudWatch Logs delivery is an independent destination: it reads its
	// own window — from its own per-region watermark to this flush's end,
	// not the S3 window — so the two destinations never lose each other's
	// failed windows. A stream failure holds the stream watermark only and
	// reports through LatestCloudWatchLogsDeliveryError, never the S3
	// watermark; the failed span is re-read whole on a later flush.
	if trail.CloudWatchLogsLogGroupARN != "" {
		s.deliverCWLogsWindow(homeStore, eventStore, trail, region, windowEnd, outcome)
	}

	// The notification announces the log files written to the bucket this
	// flush, so it fires only after a successful write.
	if outcome.logDelivered && trail.SnsTopicARN != "" {
		if err := s.deliverLogNotification(homeStore, trail, region, outcome.deliveredKeys); err != nil {
			outcome.notificationErr = err.Error()
			slog.Warn("cloudtrail sns notification failed",
				"trail", trail.Name, "region", region, "error", err)
		} else {
			outcome.notified = true
		}
	}

	_, digestEvery := trailDeliveryIntervals()
	if trail.LogFileValidationEnabled && dueForDigest(trail, windowEnd, digestEvery) {
		// The digest covers every file delivered since the previous
		// digest — including the ones this same flush delivered.
		pending := append(append([]cloudtrailstore.DigestLogFile{}, trail.PendingDigestFiles...), outcome.delivered...)
		entry, err := s.deliverDigestFile(homeStore, trail, region, windowEnd, pending)
		if err != nil {
			outcome.digestErr = err.Error()
			slog.Warn("cloudtrail digest delivery failed",
				"trail", trail.Name, "region", region, "error", err)
		} else {
			outcome.digestDelivered = true
			outcome.digestEntry = &entry
		}
	}

	recordTrailDeliveryOutcome(homeStore, trail.Name, region, trail.LogFileValidationEnabled, *outcome)
	outcome.recorded = true
}

// deliverCWLogsWindow delivers the CloudWatch Logs leg of one flush: the
// selector-matched events of [the stream's own watermark, windowEnd],
// paged from the event store independently of the S3 window. The
// independence is the recovery contract — an S3 failure advances nothing
// on the stream and a stream failure advances nothing on S3, so each
// destination re-reads exactly its own undelivered span on a later flush.
func (s *CloudTrailService) deliverCWLogsWindow(homeStore, eventStore cloudtrailstore.CloudTrailStoreInterface, trail *cloudtrailstore.Trail, region string, windowEnd time.Time, outcome *trailDeliveryOutcome) {
	cwStart := cwlogsWindowStart(trail, region)
	outcome.cwlogsWindowStart = cwStart
	if !windowEnd.After(cwStart) {
		return
	}
	readFrom := cwStart.Truncate(time.Second)
	token := ""
	cwWindow := make([]*cloudtrailstore.Event, 0, trailDeliveryPageSize)
	for {
		query := cloudtrailstore.EventQuery{
			StartTime:  &readFrom,
			EndTime:    &windowEnd,
			MaxResults: trailDeliveryPageSize,
			NextToken:  token,
		}
		events, next, err := eventStore.LookupEvents(query)
		if err != nil {
			outcome.cwlogsErr = fmt.Errorf("event read failed: %w", err).Error()
			slog.Warn("cloudtrail cloudwatch-logs delivery: window read failed",
				"trail", trail.Name, "region", region, "error", err)
			return
		}
		for _, e := range events {
			// The bound is INCLUSIVE at the watermark millisecond,
			// mirroring the S3 window: stored event times are
			// millisecond-truncated, and the boundary millisecond
			// tolerates the duplicate AWS itself documents for log files.
			if !e.EventTime.Before(cwStart) && trailSelectorMatches(trail, e) {
				cwWindow = append(cwWindow, e)
			}
		}
		if next == "" {
			break
		}
		token = next
	}
	if len(cwWindow) == 0 {
		return
	}
	if err := s.deliverToCloudWatchLogs(homeStore, trail, region, cwWindow); err != nil {
		outcome.cwlogsErr = err.Error()
		slog.Warn("cloudtrail cloudwatch-logs delivery failed",
			"trail", trail.Name, "region", region, "error", err)
		return
	}
	outcome.cwlogsDelivered = true
}

// gzipAndHash compresses payload and returns the compressed bytes with
// their SHA-256 digest — the framing every S3 delivery writes with (log
// files, digest files, query result files).
func gzipAndHash(payload []byte) ([]byte, [sha256.Size]byte, error) {
	var out bytes.Buffer
	zw := gzip.NewWriter(&out)
	if _, err := zw.Write(payload); err != nil {
		return nil, [sha256.Size]byte{}, err
	}
	if err := zw.Close(); err != nil {
		return nil, [sha256.Size]byte{}, err
	}
	return out.Bytes(), sha256.Sum256(out.Bytes()), nil
}

// trailWindowStart is the exclusive lower bound of the next delivery
// window: the per-region watermark once set, else the trail's logging
// start (events before StartLogging were never recorded for this trail).
func trailWindowStart(trail *cloudtrailstore.Trail, region string) time.Time {
	if ms, ok := trail.DeliveryWatermarks[region]; ok {
		return time.UnixMilli(ms)
	}
	if trail.StartedLoggingAt != nil {
		return trail.StartedLoggingAt.UTC()
	}
	return trail.CreatedAt.UTC()
}

// cwlogsWindowStart is the CloudWatch Logs leg's own window start: the
// per-region CW Logs watermark once set, else the S3 window's start — a
// destination attached mid-life begins with the events the current S3
// window still holds, never a replay of already-delivered history.
func cwlogsWindowStart(trail *cloudtrailstore.Trail, region string) time.Time {
	if ms, ok := trail.CwlogsWatermarks[region]; ok {
		return time.UnixMilli(ms)
	}
	return trailWindowStart(trail, region)
}

// dueForDigest reports whether a digest covering up to windowEnd is due:
// the first digest fires as soon as the trail has a delivered window, the
// rest at the digest cadence from the previous digest's end.
func dueForDigest(trail *cloudtrailstore.Trail, windowEnd time.Time, every time.Duration) bool {
	if trail.LastDigestEnd == nil {
		return true
	}
	return windowEnd.Sub(*trail.LastDigestEnd) >= every
}

// deliverLogFile renders the window's records as one gzip JSON log file
// and stores it in the trail's bucket, returning the digest-bookkeeping
// entry. The file name and path follow the documented conventions.
func (s *CloudTrailService) deliverLogFile(store cloudtrailstore.CloudTrailStoreInterface, trail *cloudtrailstore.Trail, region string, events []*cloudtrailstore.Event, deliveredAt time.Time) (cloudtrailstore.DigestLogFile, error) {
	entry := cloudtrailstore.DigestLogFile{}
	invoker := s.s3Invoker()
	if invoker == nil {
		return entry, fmt.Errorf("s3 delivery unavailable")
	}

	// The AWS references do not state a record order within a log file;
	// files are written in ascending event time, the natural reading
	// order of a log (the store's lookup order is most-recent-first).
	sort.SliceStable(events, func(i, j int) bool { return events[i].EventTime.Before(events[j].EventTime) })

	var buf bytes.Buffer
	buf.WriteString(`{"Records":[`)
	newest, oldest := time.Time{}, time.Time{}
	for i, e := range events {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(e.CloudTrailEvent)
		if newest.IsZero() || e.EventTime.After(newest) {
			newest = e.EventTime
		}
		if oldest.IsZero() || e.EventTime.Before(oldest) {
			oldest = e.EventTime
		}
	}
	buf.WriteString("]}")

	key := trailObjectKey(trail.S3KeyPrefix, store.GetAccountID(), "CloudTrail", region, deliveredAt)
	name := logFileName(store.GetAccountID(), "CloudTrail", region, deliveredAt, "")
	key = key + name

	compressed, sum, err := gzipAndHash(buf.Bytes())
	if err != nil {
		return entry, err
	}
	if err := invoker.PutObject(s.ctx, region, trail.S3BucketName, key, compressed, "application/gzip"); err != nil {
		return entry, err
	}
	return cloudtrailstore.DigestLogFile{
		Bucket:        trail.S3BucketName,
		ObjectKey:     key,
		HashHex:       hex.EncodeToString(sum[:]),
		NewestEventAt: newest.UTC(),
		OldestEventAt: oldest.UTC(),
	}, nil
}

// digestFile is the digest file JSON: the AWS field spellings of the
// digest file reference's sample, including digestPublicKeyFingerprint
// and the previous-digest chain members (null in a starting digest).
type digestFile struct {
	AWSAccountId                string             `json:"awsAccountId"`
	DigestStartTime             string             `json:"digestStartTime"`
	DigestEndTime               string             `json:"digestEndTime"`
	DigestS3Bucket              string             `json:"digestS3Bucket"`
	DigestS3Object              string             `json:"digestS3Object"`
	DigestPublicKeyFingerprint  string             `json:"digestPublicKeyFingerprint"`
	DigestSignatureAlgorithm    string             `json:"digestSignatureAlgorithm"`
	NewestEventTime             *string            `json:"newestEventTime"`
	OldestEventTime             *string            `json:"oldestEventTime"`
	PreviousDigestS3Bucket      *string            `json:"previousDigestS3Bucket"`
	PreviousDigestS3Object      *string            `json:"previousDigestS3Object"`
	PreviousDigestHashValue     *string            `json:"previousDigestHashValue"`
	PreviousDigestHashAlgorithm *string            `json:"previousDigestHashAlgorithm"`
	PreviousDigestSignature     *string            `json:"previousDigestSignature"`
	LogFiles                    []digestLogFileRef `json:"logFiles"`
}

// digestLogFileRef is one delivered log file's digest entry.
type digestLogFileRef struct {
	S3Bucket        string  `json:"s3Bucket"`
	S3Object        string  `json:"s3Object"`
	HashValue       string  `json:"hashValue"`
	HashAlgorithm   string  `json:"hashAlgorithm"`
	NewestEventTime *string `json:"newestEventTime"`
	OldestEventTime *string `json:"oldestEventTime"`
}

// deliverDigestFile renders and delivers the digest covering the given
// log files, signs it with the trail's validation key (the signature
// travels as S3 object metadata), and returns the delivered digest's own
// chain entry. The chain state on the trail record is advanced by the
// caller after the outcome is recorded.
func (s *CloudTrailService) deliverDigestFile(store cloudtrailstore.CloudTrailStoreInterface, trail *cloudtrailstore.Trail, region string, windowEnd time.Time, pending []cloudtrailstore.DigestLogFile) (cloudtrailstore.DigestLogFile, error) {
	none := cloudtrailstore.DigestLogFile{}
	invoker := s.s3Invoker()
	if invoker == nil {
		return none, fmt.Errorf("s3 delivery unavailable")
	}
	pub, priv, err := store.LoadTrailSigningKey(trail.Name)
	if err != nil {
		return none, err
	}

	digestStart := windowEnd.Add(-digestInterval)
	if trail.LastDigestEnd != nil {
		digestStart = *trail.LastDigestEnd
	} else if trail.StartedLoggingAt != nil && trail.StartedLoggingAt.After(digestStart) {
		digestStart = trail.StartedLoggingAt.UTC()
	}

	key := trailObjectKey(trail.S3KeyPrefix, store.GetAccountID(), "CloudTrail-Digest", region, windowEnd)
	name := digestFileName(store.GetAccountID(), region, trail.Name, trail.HomeRegion, windowEnd)
	objectKey := key + name

	refs := make([]digestLogFileRef, 0, len(pending))
	var newest, oldest *time.Time
	for _, f := range pending {
		n, o := f.NewestEventAt.UTC(), f.OldestEventAt.UTC()
		nn, oo := n, o
		if newest == nil || nn.After(*newest) {
			newest = &nn
		}
		if oldest == nil || oo.Before(*oldest) {
			oldest = &oo
		}
		newestStr, oldestStr := nn.Format(time.RFC3339), oo.Format(time.RFC3339)
		refs = append(refs, digestLogFileRef{
			S3Bucket:        f.Bucket,
			S3Object:        f.ObjectKey,
			HashValue:       f.HashHex,
			HashAlgorithm:   "SHA-256",
			NewestEventTime: &newestStr,
			OldestEventTime: &oldestStr,
		})
	}
	var newestStr, oldestStr *string
	if newest != nil {
		s := newest.Format(time.RFC3339)
		newestStr = &s
	}
	if oldest != nil {
		s := oldest.Format(time.RFC3339)
		oldestStr = &s
	}

	digest := digestFile{
		AWSAccountId:               store.GetAccountID(),
		DigestStartTime:            digestStart.UTC().Format(time.RFC3339),
		DigestEndTime:              windowEnd.UTC().Format(time.RFC3339),
		DigestS3Bucket:             trail.S3BucketName,
		DigestS3Object:             objectKey,
		DigestPublicKeyFingerprint: pub.Fingerprint(),
		DigestSignatureAlgorithm:   "SHA256withRSA",
		NewestEventTime:            newestStr,
		OldestEventTime:            oldestStr,
		LogFiles:                   refs,
	}
	if trail.PreviousDigestObject != "" {
		prevBucket, prevObject, prevHash, prevSig, prevAlgo :=
			trail.PreviousDigestBucket, trail.PreviousDigestObject,
			trail.PreviousDigestHash, trail.PreviousDigestSignature, "SHA-256"
		digest.PreviousDigestS3Bucket = &prevBucket
		digest.PreviousDigestS3Object = &prevObject
		digest.PreviousDigestHashValue = &prevHash
		digest.PreviousDigestHashAlgorithm = &prevAlgo
		digest.PreviousDigestSignature = &prevSig
	}

	body, err := json.Marshal(digest)
	if err != nil {
		return none, err
	}
	digestSum := sha256.Sum256(body)
	signature, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, digestSum[:])
	if err != nil {
		return none, err
	}

	compressed, _, err := gzipAndHash(body)
	if err != nil {
		return none, err
	}
	// The S3 layer serves stored metadata keys under the x-amz-meta-
	// prefix, so the keys here are the bare suffixes.
	metadata := map[string]string{
		"signature":           hex.EncodeToString(signature),
		"signature-algorithm": "SHA256withRSA",
	}
	if err := invoker.PutObjectWithMetadata(s.ctx, region, trail.S3BucketName, objectKey, compressed, "application/gzip", metadata); err != nil {
		return none, err
	}
	return cloudtrailstore.DigestLogFile{
		Bucket:        trail.S3BucketName,
		ObjectKey:     objectKey,
		HashHex:       hex.EncodeToString(digestSum[:]),
		NewestEventAt: newestOrTime(newest, windowEnd),
		OldestEventAt: oldestOrTime(oldest, digestStart),
		SignatureHex:  hex.EncodeToString(signature),
	}, nil
}

// newestOrTime dereferences an optional aggregate time, falling back to
// the given bound when the digest carried no log files.
func newestOrTime(t *time.Time, fallback time.Time) time.Time {
	if t != nil {
		return *t
	}
	return fallback
}

// oldestOrTime is newestOrTime's lower-bound counterpart.
func oldestOrTime(t *time.Time, fallback time.Time) time.Time {
	if t != nil {
		return *t
	}
	return fallback
}

// cloudTrailLogStreamName renders the CloudWatch Logs stream CloudTrail
// delivers a trail's events to: {{account_ID}}_CloudTrail_{{trail_region}}
// ("CloudWatch log group and log stream naming for CloudTrail") — the
// trail's home region for every delivering region, so a multi-region
// trail's events from all enabled regions flow into the one stream of the
// one log group ("A multi-Region trail sends log files from all enabled
// Regions in your AWS account to the CloudWatch Logs log group that you
// specify").
const cloudTrailLogStreamFormat = "%s_CloudTrail_%s"

// cloudWatchLogsDestination resolves the trail's CloudWatch Logs
// destination from its ARN (arn:aws:logs:<region>:<account>:log-group:
// <name>): the group lives in the ARN's region, and every invoker call
// addresses that region whichever region's window is being delivered. An
// unparseable ARN keeps the historical textual extraction and addresses
// the trail's home region — the group's own region in the documented
// console flow ("choose an existing one in the same Region as the
// trail").
func cloudWatchLogsDestination(logGroupARN, homeRegion string) (group, region string) {
	if parsed, err := arn.ParseARN(logGroupARN); err == nil && strings.HasPrefix(parsed.Resource, "log-group:") {
		return strings.TrimPrefix(parsed.Resource, "log-group:"), parsed.Region
	}
	if _, rest, ok := strings.Cut(logGroupARN, "log-group:"); ok {
		return rest, homeRegion
	}
	return logGroupARN, homeRegion
}

// deliverToCloudWatchLogs writes the window's events to the trail's
// CloudWatch Logs log group. Each event is one log entry carrying the event
// record JSON at the event's own time; the stream is named with the trail's
// account and home region — one stream receives every delivering region's
// events — and is created on first delivery (the documented role grants
// CloudTrail logs:CreateLogStream and logs:PutLogEvents on exactly it).
// Every invoker call addresses the group's own region, resolved from the
// group ARN. The log group itself is not created: it must exist in the
// account, and a missing group reports through
// LatestCloudWatchLogsDeliveryError.
// The entries go out in chronological order in batches bounded by the
// documented PutLogEvents limits; an entry whose message alone exceeds the
// per-event bound can never fit any batch and is skipped with a warning —
// its size is a recorder-side defect (AWS caps audit event size), not a
// destination failure, so it does not poison LatestCloudWatchLogsDeliveryError.
func (s *CloudTrailService) deliverToCloudWatchLogs(store cloudtrailstore.CloudTrailStoreInterface, trail *cloudtrailstore.Trail, region string, events []*cloudtrailstore.Event) error {
	invoker := s.logsInvoker()
	if invoker == nil {
		return fmt.Errorf("cloudwatch logs delivery unavailable")
	}
	group, groupRegion := cloudWatchLogsDestination(trail.CloudWatchLogsLogGroupARN, trail.HomeRegion)
	stream := fmt.Sprintf(cloudTrailLogStreamFormat, store.GetAccountID(), trail.HomeRegion)
	if err := invoker.EnsureLogStream(s.ctx, groupRegion, group, stream); err != nil {
		return err
	}
	entries := make([]invokers.LogsLogEntry, len(events))
	for i, e := range events {
		entries[i] = invokers.LogsLogEntry{Timestamp: e.EventTime.UnixMilli(), Message: e.CloudTrailEvent}
	}
	// The window walk serves events most-recent-first; a batch must be
	// chronological by timestamp (PutLogEvents reference).
	sort.Slice(entries, func(i, j int) bool { return entries[i].Timestamp < entries[j].Timestamp })
	for start := 0; start < len(entries); {
		batch, consumed := takePutLogEventsBatch(entries[start:])
		if skipped := consumed - len(batch); skipped > 0 {
			slog.Warn("cloudtrail cloudwatch-logs delivery: event message exceeds the per-event bound, skipping",
				"trail", trail.Name, "region", region,
				"messageBytes", len(entries[start+consumed-1].Message))
		}
		if len(batch) > 0 {
			if err := invoker.PutLogEvents(s.ctx, groupRegion, group, stream, batch); err != nil {
				return err
			}
		}
		start += consumed
	}
	return nil
}

// takePutLogEventsBatch gathers the next PutLogEvents batch from entries:
// up to PutLogEventsBatchMaxEvents entries whose messages plus the
// 26-byte per-event allowance stay within PutLogEventsBatchMaxBytes. An
// entry that alone exceeds the per-event bound is skipped: it can never
// fit any batch, and holding the delivery on it would block every later
// event. The return values are the batch to send and the number of input
// entries consumed (skipped entries included).
func takePutLogEventsBatch(entries []invokers.LogsLogEntry) ([]invokers.LogsLogEntry, int) {
	batch := make([]invokers.LogsLogEntry, 0, len(entries))
	batchBytes := 0
	for i, e := range entries {
		if i == invokers.PutLogEventsBatchMaxEvents {
			break
		}
		n := len(e.Message) + invokers.PutLogEventsEventOverhead
		if n > invokers.PutLogEventsBatchMaxBytes {
			return batch, i + 1
		}
		if batchBytes+n > invokers.PutLogEventsBatchMaxBytes && len(batch) > 0 {
			break
		}
		batch = append(batch, e)
		batchBytes += n
	}
	return batch, len(batch)
}

// snsNotificationMessage is the JSON object CloudTrail publishes to the
// trail's SNS topic when log files are delivered: the bucket and the
// delivered object keys.
type snsNotificationMessage struct {
	S3Bucket    string   `json:"s3Bucket"`
	S3ObjectKey []string `json:"s3ObjectKey"`
}

// deliverLogNotification publishes the delivered-files notification to the
// trail's resolved SNS topic.
func (s *CloudTrailService) deliverLogNotification(store cloudtrailstore.CloudTrailStoreInterface, trail *cloudtrailstore.Trail, region string, keys []string) error {
	invoker := s.snsInvoker()
	if invoker == nil {
		return fmt.Errorf("sns delivery unavailable")
	}
	message, err := json.Marshal(snsNotificationMessage{S3Bucket: trail.S3BucketName, S3ObjectKey: keys})
	if err != nil {
		return err
	}
	_, err = invoker.PublishToTopic(s.ctx, trail.SnsTopicARN, string(message), "", nil)
	return err
}

// recordTrailDeliveryOutcome writes one flush's bookkeeping to the trail
// under the store lock: the watermark advances only when the window was
// delivered (or held nothing to deliver), the Latest* members report the
// outcomes GetTrailStatus serves, and the digest chain moves to the
// digest just delivered.
func recordTrailDeliveryOutcome(store cloudtrailstore.CloudTrailStoreInterface, trailName, region string, validationEnabled bool, outcome trailDeliveryOutcome) {
	_, err := store.MutateTrail(trailName, func(trail *cloudtrailstore.Trail) error {
		if !trail.IsLogging {
			// Logging stopped while the window was being delivered; the
			// window is redelivered when logging restarts.
			return cloudtrailstore.ErrUnchanged
		}
		deliveryHeld := outcome.deliveryErr == ""
		if deliveryHeld {
			if trail.DeliveryWatermarks == nil {
				trail.DeliveryWatermarks = make(map[string]int64)
			}
			trail.DeliveryWatermarks[region] = outcome.windowEnd.UnixMilli()
			if outcome.logDelivered {
				now := outcome.now
				trail.LatestDeliveryTime = &now
				trail.LatestDeliveryError = ""
				trail.LatestDeliveryAttemptTime = &now
				trail.LatestDeliveryAttemptSuccess = true
			}
			if validationEnabled && !outcome.digestDelivered && len(outcome.delivered) > 0 {
				trail.PendingDigestFiles = append(trail.PendingDigestFiles, outcome.delivered...)
			}
		} else {
			now := outcome.now
			trail.LatestDeliveryAttemptTime = &now
			trail.LatestDeliveryAttemptSuccess = false
			trail.LatestDeliveryError = outcome.deliveryErr
			// Files finalised before the failure stay in the digest
			// chain: their objects exist in the bucket, so dropping them
			// from the pending set would punch a hole the retry's
			// redelivered copies cannot close.
			if validationEnabled && len(outcome.delivered) > 0 {
				trail.PendingDigestFiles = append(trail.PendingDigestFiles, outcome.delivered...)
			}
		}
		if outcome.digestDelivered {
			now := outcome.now
			trail.LatestDigestTime = &now
			trail.LatestDigestError = ""
			if outcome.digestEntry != nil {
				trail.PreviousDigestBucket = outcome.digestEntry.Bucket
				trail.PreviousDigestObject = outcome.digestEntry.ObjectKey
				trail.PreviousDigestHash = outcome.digestEntry.HashHex
				trail.PreviousDigestSignature = outcome.digestEntry.SignatureHex
			}
			end := outcome.windowEnd
			trail.LastDigestEnd = &end
			trail.PendingDigestFiles = nil
		} else if outcome.digestErr != "" {
			trail.LatestDigestError = outcome.digestErr
		}
		// The CloudWatch Logs and SNS destinations record their own
		// outcomes independently of the S3 watermark: their failures
		// report through their Latest* members and never hold the S3
		// window. CloudWatch Logs additionally advances its own
		// per-region watermark only on its own success, so an S3 retry
		// never re-sends the stream's events and a failed stream window
		// stays below its watermark, to be re-read whole on a later
		// flush.
		if outcome.cwlogsDelivered {
			now := outcome.now
			trail.LatestCWLogsDeliveryTime = &now
			trail.LatestCWLogsDeliveryError = ""
			if trail.CwlogsWatermarks == nil {
				trail.CwlogsWatermarks = make(map[string]int64)
			}
			trail.CwlogsWatermarks[region] = outcome.windowEnd.UnixMilli()
		} else if outcome.cwlogsErr != "" {
			trail.LatestCWLogsDeliveryError = outcome.cwlogsErr
			// The failed stream window pins the stream watermark to the
			// span's start: the events before it stand delivered, and the
			// next flush re-reads the span from there instead of falling
			// back to the already-advanced S3 window start.
			if !outcome.cwlogsWindowStart.IsZero() {
				if trail.CwlogsWatermarks == nil {
					trail.CwlogsWatermarks = make(map[string]int64)
				}
				trail.CwlogsWatermarks[region] = outcome.cwlogsWindowStart.UnixMilli()
			}
		}
		if outcome.notified {
			now := outcome.now
			trail.LatestNotificationTime = &now
			trail.LatestNotificationError = ""
			trail.LatestNotificationAttemptTime = &now
			trail.LatestNotificationAttemptSuccess = true
		} else if outcome.notificationErr != "" {
			now := outcome.now
			trail.LatestNotificationAttemptTime = &now
			trail.LatestNotificationAttemptSuccess = false
			trail.LatestNotificationError = outcome.notificationErr
		}
		return nil
	})
	// A trail deleted mid-flight has nothing to record against; any other
	// failure to persist the outcome is a store fault worth the log.
	if err != nil && !errors.Is(err, cloudtrailstore.ErrTrailNotFound) {
		slog.Error("cloudtrail trail delivery: failed to record outcome",
			"trail", trailName, "region", region, "error", err)
	}
}

// trailObjectKey builds the delivered-object folder path:
// [prefix/]AWSLogs/<account>/<family>/<region>/<Y>/<M>/<D>/. The trailing
// file name is appended by the caller (log files and digests name
// themselves differently).
func trailObjectKey(prefix, accountID, family, region string, at time.Time) string {
	var b strings.Builder
	if prefix != "" {
		b.WriteString(strings.Trim(prefix, "/"))
		b.WriteByte('/')
	}
	fmt.Fprintf(&b, "AWSLogs/%s/%s/%s/%04d/%02d/%02d/",
		accountID, family, region, at.Year(), at.Month(), at.Day())
	return b.String()
}

// logFileName renders <account>_CloudTrail_<region>_<YYYYMMDDTHHmmZ>_
// <UniqueString>.json.gz — the documented log-file name format, with the
// 16-character unique suffix appended when absent.
func logFileName(accountID, family, region string, deliveredAt time.Time, uniqueIn string) string {
	unique := uniqueIn
	if unique == "" {
		unique = randomLogFileNameSuffix()
	}
	return fmt.Sprintf("%s_%s_%s_%s_%s.json.gz",
		accountID, family, region, deliveredAt.UTC().Format("20060102T1504Z"), unique)
}

// digestFileName renders <account>_CloudTrail-Digest_<region>_<trail>_
// <homeRegion>_<YYYYMMDDTHHMMSS>Z.json.gz — the digest file location
// syntax, whose end timestamp carries seconds.
func digestFileName(accountID, region, trailName, homeRegion string, endAt time.Time) string {
	return fmt.Sprintf("%s_CloudTrail-Digest_%s_%s_%s_%s.json.gz",
		accountID, region, trailName, homeRegion, endAt.UTC().Format("20060102T150405Z"))
}

// logFileNameSuffixAlphabet is the alphanumeric set of the documented
// 16-character unique file-name component.
const logFileNameSuffixAlphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// randomLogFileNameSuffix generates the 16-character unique component.
func randomLogFileNameSuffix() string {
	raw := make([]byte, logFileUniqueStringLength)
	if _, err := rand.Read(raw); err != nil {
		// Random generation never fails on the supported platforms; a
		// time-derived suffix still disambiguates the file.
		now := time.Now().UTC().UnixNano()
		for i := range raw {
			raw[i] = logFileNameSuffixAlphabet[now%int64(len(logFileNameSuffixAlphabet))]
			now /= int64(len(logFileNameSuffixAlphabet))
		}
		return string(raw)
	}
	for i, b := range raw {
		raw[i] = logFileNameSuffixAlphabet[int(b)%len(logFileNameSuffixAlphabet)]
	}
	return string(raw)
}

// trailSelectorMatches reports whether an event falls inside the trail's
// selector configuration. With no custom selectors the default basic
// selector applies (all management events); basic selectors match on the
// read/write type and management inclusion, data resources match nothing
// (the store records management events); advanced selectors match when any
// selector's field selectors all hold — the shared store-package
// evaluator, the same one event data store ingestion applies.
func trailSelectorMatches(trail *cloudtrailstore.Trail, e *cloudtrailstore.Event) bool {
	if len(trail.AdvancedEventSelectors) > 0 {
		for _, sel := range trail.AdvancedEventSelectors {
			if cloudtrailstore.AdvancedSelectorMatches(sel, e) {
				return true
			}
		}
		return false
	}
	selectors := trail.EventSelectors
	if len(selectors) == 0 {
		return true
	}
	for _, sel := range selectors {
		if !sel.IncludeManagementEvents {
			continue
		}
		if eventSourceExcluded(sel.ExcludeManagementEventSources, e.EventSource) {
			continue
		}
		switch sel.ReadWriteType {
		case "ReadOnly":
			if e.ReadOnly != "true" {
				continue
			}
		case "WriteOnly":
			if e.ReadOnly == "true" {
				continue
			}
		}
		return true
	}
	return false
}

func eventSourceExcluded(excluded []string, source string) bool {
	for _, s := range excluded {
		if s == source {
			return true
		}
	}
	return false
}
