package cognitoidentityprovider

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// stubImportS3 backs the import worker's CSV read/delete path without a
// real S3 plane.
type stubImportS3 struct {
	objects map[string][]byte
}

func (s *stubImportS3) GetObject(_ context.Context, _, bucket, key string, _ int64) ([]byte, error) {
	data, ok := s.objects[bucket+"/"+key]
	if !ok {
		return nil, fmt.Errorf("no such object %s/%s", bucket, key)
	}
	return data, nil
}

func (s *stubImportS3) GetObjectVersion(ctx context.Context, region, bucket, key, versionID string, maxBytes int64) ([]byte, error) {
	return s.GetObject(ctx, region, bucket, key, maxBytes)
}

func (s *stubImportS3) PutObject(_ context.Context, _, bucket, key string, data []byte, _ string) error {
	s.objects[bucket+"/"+key] = data
	return nil
}

func (s *stubImportS3) ListObjects(_ context.Context, _, bucket, prefix string, _ int) ([]string, error) {
	var keys []string
	for k := range s.objects {
		if strings.HasPrefix(k, bucket+"/"+prefix) {
			keys = append(keys, k)
		}
	}
	return keys, nil
}

func (s *stubImportS3) BucketExists(_ context.Context, _, bucket string) (bool, error) {
	return bucket == importBucketName, nil
}

func (s *stubImportS3) EnsureBucket(_ context.Context, _, _ string) error { return nil }

func (s *stubImportS3) DeleteObject(_ context.Context, _, bucket, key string) error {
	delete(s.objects, bucket+"/"+key)
	return nil
}

func (s *stubImportS3) ListObjectEntries(ctx context.Context, region, bucket, prefix string, maxKeys int) ([]invokers.S3ObjectEntry, error) {
	keys, err := s.ListObjects(ctx, region, bucket, prefix, maxKeys)
	if err != nil {
		return nil, err
	}
	entries := make([]invokers.S3ObjectEntry, 0, len(keys))
	for _, key := range keys {
		entries = append(entries, invokers.S3ObjectEntry{
			Key:          key,
			ETag:         "stub-etag",
			LastModified: 0,
			Size:         int64(len(s.objects[bucket+"/"+key])),
			StorageClass: "STANDARD",
		})
	}
	return entries, nil
}

func newImportTestService(t *testing.T) (*CognitoService, cognitostore.CognitoStoreInterface, *stubImportS3) {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	store := cognitostore.NewCognitoStore(st, "123456789012", "us-east-1")

	bus := eventbus.NewEventBus()
	if err := bus.Start(context.Background()); err != nil {
		t.Fatalf("start event bus: %v", err)
	}
	t.Cleanup(func() { _ = bus.Shutdown(context.Background()) })
	s3 := &stubImportS3{objects: map[string][]byte{}}
	bus.SetS3Invoker(s3)

	svc := NewCognitoService("123456789012", "us-east-1")
	svc.bus = bus
	svc.stores.Store("us-east-1", store)
	return svc, store, s3
}

// newImportTestPool creates a pool with one auto-verified attribute and a
// Created import job whose CSV is staged in the stub bucket.
func newImportTestPool(t *testing.T, store cognitostore.CognitoStoreInterface, s3 *stubImportS3, csv string) (poolID, jobID string) {
	t.Helper()
	pool, err := store.CreateUserPool(cognitostore.NewUserPool("import-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	if err := store.UpdateUserPoolFunc(pool.ID, func(p *cognitostore.UserPool) error {
		p.AutoVerifiedAttributes = []string{"email"}
		return nil
	}); err != nil {
		t.Fatalf("update pool: %v", err)
	}
	job := &cognitostore.UserImportJob{
		JobID:        "import-testjob",
		JobName:      "test-job",
		UserPoolID:   pool.ID,
		Status:       importJobStatusPending,
		CreationDate: time.Now().UTC(),
	}
	if err := store.CreateUserImportJob(job); err != nil {
		t.Fatalf("create job: %v", err)
	}
	s3.objects[importBucketName+"/"+importObjectKey(pool.ID, job.JobID)] = []byte(csv)
	return pool.ID, job.JobID
}

const importHappyCSV = "cognito:username,email,email_verified\n" +
	"johndoe,johndoe@example.com,TRUE\n" +
	"janedoe,janedoe@example.com,TRUE\n"

func runImportToCompletion(t *testing.T, svc *CognitoService, store cognitostore.CognitoStoreInterface, poolID, jobID string) *cognitostore.UserImportJob {
	t.Helper()
	svc.runUserImportJob("us-east-1", poolID, jobID)
	job, err := store.GetUserImportJob(poolID, jobID)
	if err != nil {
		t.Fatalf("load job after run: %v", err)
	}
	return job
}

// A straight run imports every user, reports the counters on a Succeeded
// job, and deletes the uploaded CSV.
func TestRunUserImportJobHappyPath(t *testing.T) {
	svc, store, s3 := newImportTestService(t)
	poolID, jobID := newImportTestPool(t, store, s3, importHappyCSV)

	job := runImportToCompletion(t, svc, store, poolID, jobID)

	if job.Status != importJobStatusSucceeded {
		t.Fatalf("status = %q, want Succeeded (message %q)", job.Status, job.CompletionMessage)
	}
	if job.ImportedUsers != 2 || job.SkippedUsers != 0 || job.FailedUsers != 0 {
		t.Fatalf("counters = %d/%d/%d, want 2/0/0", job.ImportedUsers, job.SkippedUsers, job.FailedUsers)
	}
	if !strings.Contains(job.CompletionMessage, "2 users were imported") {
		t.Fatalf("completion message = %q", job.CompletionMessage)
	}
	if _, ok := s3.objects[importBucketName+"/"+importObjectKey(poolID, jobID)]; ok {
		t.Fatal("CSV object still present after completion")
	}
	if _, err := store.GetUser(poolID, "johndoe"); err != nil {
		t.Fatalf("imported user missing: %v", err)
	}
}

// The worker batches its progress writes on the status-check cadence; a
// run spanning several cadence batches must still land the exact counters
// through the in-loop flushes and the final flush after the loop.
func TestRunUserImportJobCountersExactAcrossCadenceBatches(t *testing.T) {
	svc, store, s3 := newImportTestService(t)
	const rows = 3*importStatusCheckRows + 50
	var csv strings.Builder
	csv.WriteString("cognito:username,email,email_verified\n")
	for i := 0; i < rows; i++ {
		fmt.Fprintf(&csv, "user%04d,user%04d@example.com,TRUE\n", i, i)
	}
	poolID, jobID := newImportTestPool(t, store, s3, csv.String())

	job := runImportToCompletion(t, svc, store, poolID, jobID)

	if job.Status != importJobStatusSucceeded {
		t.Fatalf("status = %q, want Succeeded (message %q)", job.Status, job.CompletionMessage)
	}
	if job.ImportedUsers != rows || job.SkippedUsers != 0 || job.FailedUsers != 0 {
		t.Fatalf("counters = %d/%d/%d, want %d/0/0", job.ImportedUsers, job.SkippedUsers, job.FailedUsers, rows)
	}
}

// Usernames differing only in case are duplicates when the pool has not
// opted into case-sensitive usernames: the second row is skipped.
func TestRunUserImportJobCaseInsensitiveDuplicate(t *testing.T) {
	svc, store, s3 := newImportTestService(t)
	csv := "cognito:username,email,email_verified\n" +
		"John,john@example.com,TRUE\n" +
		"john,john2@example.com,TRUE\n"
	poolID, jobID := newImportTestPool(t, store, s3, csv)

	job := runImportToCompletion(t, svc, store, poolID, jobID)

	if job.ImportedUsers != 1 || job.SkippedUsers != 1 {
		t.Fatalf("counters = %d imported / %d skipped, want 1/1", job.ImportedUsers, job.SkippedUsers)
	}
}

// A row whose field count differs from the header fails for that user
// only; the job still succeeds.
func TestRunUserImportJobFieldCountMismatchFailsRow(t *testing.T) {
	svc, store, s3 := newImportTestService(t)
	csv := "cognito:username,email,email_verified\n" +
		"johndoe,johndoe@example.com\n" +
		"janedoe,janedoe@example.com,TRUE\n"
	poolID, jobID := newImportTestPool(t, store, s3, csv)

	job := runImportToCompletion(t, svc, store, poolID, jobID)

	if job.Status != importJobStatusSucceeded {
		t.Fatalf("status = %q, want Succeeded", job.Status)
	}
	if job.ImportedUsers != 1 || job.FailedUsers != 1 {
		t.Fatalf("counters = %d imported / %d failed, want 1/1", job.ImportedUsers, job.FailedUsers)
	}
}

// A header column that is not a pool attribute fails the whole job.
func TestRunUserImportJobUnknownHeaderColumnFailsJob(t *testing.T) {
	svc, store, s3 := newImportTestService(t)
	csv := "cognito:username,email,email_verified,not_an_attribute\n" +
		"johndoe,johndoe@example.com,TRUE,x\n"
	poolID, jobID := newImportTestPool(t, store, s3, csv)

	job := runImportToCompletion(t, svc, store, poolID, jobID)

	if job.Status != importJobStatusFailed {
		t.Fatalf("status = %q, want Failed", job.Status)
	}
	if !strings.Contains(job.CompletionMessage, "not_an_attribute") {
		t.Fatalf("completion message = %q", job.CompletionMessage)
	}
}

// Deleting the user pool removes the job record with it, so a worker that
// runs afterwards finds no job and leaves nothing behind to finalise.
func TestRunUserImportJobDeletedPoolRemovesJob(t *testing.T) {
	svc, store, s3 := newImportTestService(t)
	poolID, jobID := newImportTestPool(t, store, s3, importHappyCSV)
	if err := store.DeleteUserPool(poolID); err != nil {
		t.Fatalf("delete pool: %v", err)
	}

	svc.runUserImportJob("us-east-1", poolID, jobID)

	if _, err := store.GetUserImportJob(poolID, jobID); err == nil {
		t.Fatal("job record must be gone with the deleted pool")
	}
	if _, ok := s3.objects[importBucketName+"/"+importObjectKey(poolID, jobID)]; ok {
		t.Fatal("uploaded CSV must be removed when the pool is deleted before the worker runs")
	}
}

// A job record still in a running state after the worker exits (for
// example after a recovered panic) is finalised as Failed.
func TestEnsureTerminalImportStateFailsStrandedRunningJob(t *testing.T) {
	svc, store, _ := newImportTestService(t)
	pool, err := store.CreateUserPool(cognitostore.NewUserPool("stranded-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	job := &cognitostore.UserImportJob{
		JobID:        "import-stranded",
		JobName:      "stranded",
		UserPoolID:   pool.ID,
		Status:       importJobStatusInProgress,
		CreationDate: time.Now().UTC(),
	}
	if err := store.CreateUserImportJob(job); err != nil {
		t.Fatalf("create job: %v", err)
	}

	svc.ensureTerminalImportState("us-east-1", pool.ID, job.JobID)

	final, err := store.GetUserImportJob(pool.ID, job.JobID)
	if err != nil {
		t.Fatalf("load job: %v", err)
	}
	if final.Status != importJobStatusFailed {
		t.Fatalf("status = %q, want Failed", final.Status)
	}
	if final.CompletionDate.IsZero() {
		t.Fatal("completion date not set")
	}
}

// A job already finalised (Stopped) is left untouched by the terminal
// guard.
func TestEnsureTerminalImportStateKeepsTerminalJob(t *testing.T) {
	svc, store, _ := newImportTestService(t)
	pool, err := store.CreateUserPool(cognitostore.NewUserPool("stopped-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	job := &cognitostore.UserImportJob{
		JobID:        "import-stopped",
		JobName:      "stopped",
		UserPoolID:   pool.ID,
		Status:       importJobStatusStopped,
		CreationDate: time.Now().UTC(),
	}
	if err := store.CreateUserImportJob(job); err != nil {
		t.Fatalf("create job: %v", err)
	}

	svc.ensureTerminalImportState("us-east-1", pool.ID, job.JobID)

	final, err := store.GetUserImportJob(pool.ID, job.JobID)
	if err != nil {
		t.Fatalf("load job: %v", err)
	}
	if final.Status != importJobStatusStopped {
		t.Fatalf("status = %q, want Stopped", final.Status)
	}
}

// AWS deletes the uploaded CSV once a job reaches a terminal state, so the
// guard's abnormal-termination path must remove it together with failing
// the job.
func TestEnsureTerminalImportStateFailsJobAndDeletesCSV(t *testing.T) {
	svc, store, s3 := newImportTestService(t)
	pool, jobID := newImportTestPool(t, store, s3, importHappyCSV)
	job, err := store.GetUserImportJob(pool, jobID)
	if err != nil {
		t.Fatalf("load job: %v", err)
	}
	job.Status = importJobStatusInProgress
	if err := store.UpdateUserImportJob(job); err != nil {
		t.Fatalf("update job: %v", err)
	}

	svc.ensureTerminalImportState("us-east-1", pool, jobID)

	final, err := store.GetUserImportJob(pool, jobID)
	if err != nil {
		t.Fatalf("reload job: %v", err)
	}
	if final.Status != importJobStatusFailed {
		t.Fatalf("status = %q, want Failed", final.Status)
	}
	if _, ok := s3.objects[importBucketName+"/"+importObjectKey(pool, jobID)]; ok {
		t.Fatal("uploaded CSV must be deleted when the guard fails a job")
	}
}

// A job interrupted while Stopping is finalised as Stopped by the guard;
// the uploaded CSV goes with it.
func TestEnsureTerminalImportStateStoppedJobDeletesCSV(t *testing.T) {
	svc, store, s3 := newImportTestService(t)
	pool, jobID := newImportTestPool(t, store, s3, importHappyCSV)
	job, err := store.GetUserImportJob(pool, jobID)
	if err != nil {
		t.Fatalf("load job: %v", err)
	}
	job.Status = importJobStatusStopping
	if err := store.UpdateUserImportJob(job); err != nil {
		t.Fatalf("update job: %v", err)
	}

	svc.ensureTerminalImportState("us-east-1", pool, jobID)

	final, err := store.GetUserImportJob(pool, jobID)
	if err != nil {
		t.Fatalf("reload job: %v", err)
	}
	if final.Status != importJobStatusStopped {
		t.Fatalf("status = %q, want Stopped", final.Status)
	}
	if _, ok := s3.objects[importBucketName+"/"+importObjectKey(pool, jobID)]; ok {
		t.Fatal("uploaded CSV must be deleted when the guard stops a job")
	}
}

// A stop request that lands between the terminal guard's read of the job
// and its finalisation write must win: finaliseImportJob goes through the
// status-guarded transition, so a stale InProgress snapshot yields a
// tolerated conflict instead of overwriting the already-Stopped state
// with Failed. The CSV still goes, because the job is terminal either way.
func TestFinaliseImportJobDoesNotOverwriteConcurrentStop(t *testing.T) {
	svc, store, s3 := newImportTestService(t)
	pool, jobID := newImportTestPool(t, store, s3, importHappyCSV)
	seed, err := store.GetUserImportJob(pool, jobID)
	if err != nil {
		t.Fatalf("load job: %v", err)
	}
	seed.Status = importJobStatusInProgress
	if err := store.UpdateUserImportJob(seed); err != nil {
		t.Fatalf("seed in progress: %v", err)
	}

	// The snapshot the guard read before the stop request landed.
	snapshot, err := store.GetUserImportJob(pool, jobID)
	if err != nil {
		t.Fatalf("snapshot job: %v", err)
	}

	// The concurrent stop: the developer stops the job after the guard
	// read its snapshot but before the guard writes.
	if _, err := store.TransitionUserImportJobStatus(pool, jobID, importJobStatusInProgress, importJobStatusStopping, nil); err != nil {
		t.Fatalf("stop to Stopping: %v", err)
	}
	if _, err := store.TransitionUserImportJobStatus(pool, jobID, importJobStatusStopping, importJobStatusStopped, func(j *cognitostore.UserImportJob) {
		j.CompletionMessage = stoppedImportCompletionMessage
	}); err != nil {
		t.Fatalf("finalise Stopped: %v", err)
	}

	svc.finaliseImportJob(store, "us-east-1", snapshot)

	final, err := store.GetUserImportJob(pool, jobID)
	if err != nil {
		t.Fatalf("reload job: %v", err)
	}
	if final.Status != importJobStatusStopped {
		t.Fatalf("status = %q, want Stopped (the concurrent stop must not be overwritten)", final.Status)
	}
	if final.CompletionMessage != stoppedImportCompletionMessage {
		t.Fatalf("completion message = %q, want the developer-stopped message", final.CompletionMessage)
	}
	if _, ok := s3.objects[importBucketName+"/"+importObjectKey(pool, jobID)]; ok {
		t.Fatal("uploaded CSV must be deleted once the job is terminal")
	}
}

// The startup finaliser sweeps every job a previous process left in a
// non-terminal state to its terminal state, deletes the orphaned CSVs, and
// leaves terminal and never-started jobs untouched — after the sweep, the
// one-active-job start guard accepts a new import again.
func TestFinaliseStaleImportJobs(t *testing.T) {
	svc, store, s3 := newImportTestService(t)
	pool, err := store.CreateUserPool(cognitostore.NewUserPool("stale-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}

	seed := func(jobID, status string) {
		t.Helper()
		job := &cognitostore.UserImportJob{
			JobID:        jobID,
			JobName:      "stale",
			UserPoolID:   pool.ID,
			Status:       status,
			CreationDate: time.Now().UTC(),
		}
		if err := store.CreateUserImportJob(job); err != nil {
			t.Fatalf("create job %s: %v", jobID, err)
		}
		s3.objects[importBucketName+"/"+importObjectKey(pool.ID, jobID)] = []byte("stale csv")
	}
	seed("stale-inprogress", importJobStatusInProgress)
	seed("stale-pending", importJobStatusPending)
	seed("stale-stopping", importJobStatusStopping)
	seed("stale-created", importJobStatusCreated)
	terminal := &cognitostore.UserImportJob{
		JobID:      "stale-succeeded",
		JobName:    "done",
		UserPoolID: pool.ID,
		Status:     importJobStatusSucceeded,
	}
	if err := store.CreateUserImportJob(terminal); err != nil {
		t.Fatalf("create terminal job: %v", err)
	}

	svc.finaliseStaleImportJobs("us-east-1", store)

	for _, jobID := range []string{"stale-inprogress", "stale-pending"} {
		final, err := store.GetUserImportJob(pool.ID, jobID)
		if err != nil {
			t.Fatalf("load %s: %v", jobID, err)
		}
		if final.Status != importJobStatusFailed {
			t.Fatalf("%s status = %q, want Failed", jobID, final.Status)
		}
		if final.CompletionMessage != "The import job terminated before completion." {
			t.Fatalf("%s completion message = %q", jobID, final.CompletionMessage)
		}
		if final.CompletionDate.IsZero() {
			t.Fatalf("%s completion date not set", jobID)
		}
	}
	stopped, err := store.GetUserImportJob(pool.ID, "stale-stopping")
	if err != nil {
		t.Fatalf("load stale-stopping: %v", err)
	}
	if stopped.Status != importJobStatusStopped || stopped.CompletionMessage != stoppedImportCompletionMessage {
		t.Fatalf("stale-stopping = %q / %q", stopped.Status, stopped.CompletionMessage)
	}
	for _, jobID := range []string{"stale-inprogress", "stale-pending", "stale-stopping"} {
		if _, ok := s3.objects[importBucketName+"/"+importObjectKey(pool.ID, jobID)]; ok {
			t.Fatalf("CSV of %s must be deleted with the sweep", jobID)
		}
	}

	// The never-started and already-terminal records keep their state.
	created, err := store.GetUserImportJob(pool.ID, "stale-created")
	if err != nil {
		t.Fatalf("load stale-created: %v", err)
	}
	if created.Status != importJobStatusCreated {
		t.Fatalf("stale-created status = %q, want Created", created.Status)
	}
	succeeded, err := store.GetUserImportJob(pool.ID, "stale-succeeded")
	if err != nil {
		t.Fatalf("load stale-succeeded: %v", err)
	}
	if succeeded.Status != importJobStatusSucceeded {
		t.Fatalf("stale-succeeded status = %q, want Succeeded", succeeded.Status)
	}

	// With the orphans finalised, a new job passes the start guard.
	fresh := &cognitostore.UserImportJob{
		JobID:        "stale-fresh",
		JobName:      "fresh",
		UserPoolID:   pool.ID,
		Status:       importJobStatusCreated,
		CreationDate: time.Now().UTC(),
	}
	if err := store.CreateUserImportJob(fresh); err != nil {
		t.Fatalf("create fresh job: %v", err)
	}
	started, err := store.StartUserImportJobIfEligible(pool.ID, "stale-fresh")
	if err != nil {
		t.Fatalf("start after sweep: %v", err)
	}
	if started.Status != importJobStatusPending {
		t.Fatalf("fresh job status = %q, want Pending", started.Status)
	}
}

// AWS example responses always include the counters, including zeros.
func TestFormatUserImportJobAlwaysReportsCounters(t *testing.T) {
	job := &cognitostore.UserImportJob{
		JobID:        "import-zero",
		JobName:      "zero",
		UserPoolID:   "pool",
		Status:       importJobStatusCreated,
		CreationDate: time.Now().UTC(),
	}
	formatted := formatUserImportJob(job)
	for _, key := range []string{"ImportedUsers", "SkippedUsers", "FailedUsers"} {
		v, ok := formatted[key]
		if !ok {
			t.Fatalf("missing %s in %v", key, formatted)
		}
		if v.(int64) != 0 {
			t.Fatalf("%s = %v, want 0", key, v)
		}
	}
}

// An import that produced no users at all ends Failed with the AWS
// "Too many users..." wording; an import with at least one success ends
// Succeeded even with failures.
func TestRunUserImportJobAllRowsFailEndsFailed(t *testing.T) {
	svc, store, s3 := newImportTestService(t)
	csv := "cognito:username,email,email_verified\n" +
		"baduser1,bad1@example.com,false\n" +
		"baduser2,bad2@example.com,false\n"
	poolID, jobID := newImportTestPool(t, store, s3, csv)

	job := runImportToCompletion(t, svc, store, poolID, jobID)

	if job.Status != importJobStatusFailed {
		t.Fatalf("status = %q, want Failed", job.Status)
	}
	if job.CompletionMessage != "Too many users have failed or been skipped during the import." {
		t.Fatalf("completion message = %q", job.CompletionMessage)
	}
}

// A job that specifies a password hashing algorithm fails outright when
// the CSV header carries no password_hash column.
func TestRunUserImportJobMissingPasswordHashColumnFailsJob(t *testing.T) {
	svc, store, s3 := newImportTestService(t)
	poolID, jobID := newImportTestPool(t, store, s3, importHappyCSV)

	job, err := store.GetUserImportJob(poolID, jobID)
	if err != nil {
		t.Fatalf("load job: %v", err)
	}
	job.PasswordHashingAlgorithm = "BCRYPT"
	if err := store.UpdateUserImportJob(job); err != nil {
		t.Fatalf("update job: %v", err)
	}

	final := runImportToCompletion(t, svc, store, poolID, jobID)
	if final.Status != importJobStatusFailed {
		t.Fatalf("status = %q, want Failed", final.Status)
	}
	if _, ok := s3.objects[importBucketName+"/"+importObjectKey(poolID, jobID)]; ok {
		t.Fatal("uploaded CSV must be deleted when the job fails")
	}
}

// A failed job deletes the uploaded CSV alongside its terminal transition.
func TestRunUserImportJobFailureDeletesCSV(t *testing.T) {
	svc, store, s3 := newImportTestService(t)
	// An unknown header column fails the job before any row is applied.
	csv := "cognito:username,email,email_verified,not_an_attribute\n" +
		"johndoe,johndoe@example.com,TRUE,x\n"
	poolID, jobID := newImportTestPool(t, store, s3, csv)

	final := runImportToCompletion(t, svc, store, poolID, jobID)
	if final.Status != importJobStatusFailed {
		t.Fatalf("status = %q, want Failed", final.Status)
	}
	if _, ok := s3.objects[importBucketName+"/"+importObjectKey(poolID, jobID)]; ok {
		t.Fatal("uploaded CSV must be deleted when the job fails")
	}
}

// Expiring a stale Created job deletes its uploaded data, matching the
// documented Expired behaviour.
func TestExpireStaleImportJobDeletesCSV(t *testing.T) {
	svc, store, s3 := newImportTestService(t)
	poolID, jobID := newImportTestPool(t, store, s3, importHappyCSV)

	job, err := store.GetUserImportJob(poolID, jobID)
	if err != nil {
		t.Fatalf("load job: %v", err)
	}
	job.Status = importJobStatusCreated
	job.CreationDate = time.Now().UTC().Add(-25 * time.Hour)
	if err := store.UpdateUserImportJob(job); err != nil {
		t.Fatalf("update job: %v", err)
	}
	if _, ok := s3.objects[importBucketName+"/"+importObjectKey(poolID, jobID)]; !ok {
		t.Fatal("test setup: CSV should be staged")
	}

	svc.expireStaleImportJob(store, job, "us-east-1")

	expired, err := store.GetUserImportJob(poolID, jobID)
	if err != nil {
		t.Fatalf("reload job: %v", err)
	}
	if expired.Status != importJobStatusExpired {
		t.Fatalf("status = %q, want Expired", expired.Status)
	}
	if _, ok := s3.objects[importBucketName+"/"+importObjectKey(poolID, jobID)]; ok {
		t.Fatal("uploaded CSV must be deleted when the job expires")
	}
}

// Any flow that writes a native password clears the imported-hash format
// flag; a lingering flag would lock the user out after a reset.
func TestSetNativePasswordCredentialsClearsImportedFlag(t *testing.T) {
	user := cognitostore.NewUser("us-east-1_pool", "johndoe")
	user.PasswordHash = "$2b$10$CtA.Rcu/szzn9U00wpUjOuN3vrgJRZycv4aOzcP3GzqzO8UDPEFq6"
	user.PasswordHashAlgo = "BCRYPT"

	if err := setNativePasswordCredentials(user, nil, "NewPassword1!"); err != nil {
		t.Fatalf("set native credentials: %v", err)
	}
	if user.PasswordHashAlgo != "" {
		t.Fatalf("PasswordHashAlgo = %q, want cleared", user.PasswordHashAlgo)
	}
	if user.PasswordHash == "" || user.SrpSalt == "" || user.SrpVerifier == "" {
		t.Fatal("native hash and SRP material must be populated")
	}
}

// PasswordHistorySize governs what the native write retains and forbids: a
// depth of zero retains nothing and forbids nothing, a depth of one forbids
// only the password being replaced, and an imported hash never enters the
// remembered set because it was not generated by this platform.
func TestSetNativePasswordCredentialsHistoryBoundaries(t *testing.T) {
	unrestricted := cognitostore.NewUser("us-east-1_pool", "alice")
	noHistory := &cognitostore.PasswordPolicy{PasswordHistorySize: 0}
	if err := setNativePasswordCredentials(unrestricted, noHistory, "FirstPass1!"); err != nil {
		t.Fatalf("first write with no history configured: %v", err)
	}
	if err := setNativePasswordCredentials(unrestricted, noHistory, "FirstPass1!"); err != nil {
		t.Fatalf("rewrite with no history configured: %v", err)
	}
	if len(unrestricted.PasswordHistory) != 0 {
		t.Fatalf("history retained with PasswordHistorySize 0: %d entries", len(unrestricted.PasswordHistory))
	}

	// A depth of one forbids only the password being replaced: the
	// immediate rewrite is rejected, while an older password becomes
	// reusable once a newer one has superseded it.
	depthOne := cognitostore.NewUser("us-east-1_pool", "carol")
	oneDeep := &cognitostore.PasswordPolicy{PasswordHistorySize: 1}
	if err := setNativePasswordCredentials(depthOne, oneDeep, "FirstPass1!"); err != nil {
		t.Fatalf("first write at depth one: %v", err)
	}
	if err := setNativePasswordCredentials(depthOne, oneDeep, "FirstPass1!"); !errors.Is(err, ErrPasswordHistoryViolation) {
		t.Fatalf("rewrite of the current password at depth one: err = %v, want ErrPasswordHistoryViolation", err)
	}
	if err := setNativePasswordCredentials(depthOne, oneDeep, "SecondPass2!"); err != nil {
		t.Fatalf("rotation at depth one: %v", err)
	}
	if err := setNativePasswordCredentials(depthOne, oneDeep, "FirstPass1!"); err != nil {
		t.Fatalf("superseded password must be reusable at depth one: %v", err)
	}

	// The imported hash is excluded from the remembered set, so the first
	// native write after an import is never compared against it; the
	// resulting native hash itself is remembered from then on.
	imported := cognitostore.NewUser("us-east-1_pool", "bob")
	imported.PasswordHash = "$2b$10$CtA.Rcu/szzn9U00wpUjOuN3vrgJRZycv4aOzcP3GzqzO8UDPEFq6"
	imported.PasswordHashAlgo = "BCRYPT"
	twoDeep := &cognitostore.PasswordPolicy{PasswordHistorySize: 2}
	if err := setNativePasswordCredentials(imported, twoDeep, "NewPassword1!"); err != nil {
		t.Fatalf("first native write after import: %v", err)
	}
	if len(imported.PasswordHistory) != 0 {
		t.Fatalf("imported hash must not enter the remembered set: %d entries", len(imported.PasswordHistory))
	}
	if err := setNativePasswordCredentials(imported, twoDeep, "NewPassword1!"); !errors.Is(err, ErrPasswordHistoryViolation) {
		t.Fatalf("rewrite of the current native password after import: err = %v, want ErrPasswordHistoryViolation", err)
	}
}
