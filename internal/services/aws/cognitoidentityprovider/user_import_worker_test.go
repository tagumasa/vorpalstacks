package cognitoidentityprovider

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

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
	pool.AutoVerifiedAttributes = []string{"email"}
	if err := store.UpdateUserPool(pool); err != nil {
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

// Deleting the user pool before the worker runs fails the job with the
// documented pool-deleted message.
func TestRunUserImportJobDeletedPoolFailsJob(t *testing.T) {
	svc, store, s3 := newImportTestService(t)
	poolID, jobID := newImportTestPool(t, store, s3, importHappyCSV)
	if err := store.DeleteUserPool(poolID); err != nil {
		t.Fatalf("delete pool: %v", err)
	}

	job := runImportToCompletion(t, svc, store, poolID, jobID)

	if job.Status != importJobStatusFailed {
		t.Fatalf("status = %q, want Failed", job.Status)
	}
	if !strings.Contains(job.CompletionMessage, "user pool was deleted") {
		t.Fatalf("completion message = %q", job.CompletionMessage)
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

	if err := setNativePasswordCredentials(user, "us-east-1_pool", "johndoe", "NewPassword1!"); err != nil {
		t.Fatalf("set native credentials: %v", err)
	}
	if user.PasswordHashAlgo != "" {
		t.Fatalf("PasswordHashAlgo = %q, want cleared", user.PasswordHashAlgo)
	}
	if user.PasswordHash == "" || user.SrpSalt == "" || user.SrpVerifier == "" {
		t.Fatal("native hash and SRP material must be populated")
	}
}
