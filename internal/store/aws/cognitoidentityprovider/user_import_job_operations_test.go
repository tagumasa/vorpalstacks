package cognitoidentityprovider

import (
	"errors"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
)

func newImportJobTestStore(t *testing.T) *CognitoStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return NewCognitoStore(st, "000000000000", "us-east-1")
}

func seedImportJob(t *testing.T, s *CognitoStore, poolID, jobID, status string) *UserImportJob {
	t.Helper()
	job := &UserImportJob{
		JobID:        jobID,
		JobName:      "job",
		UserPoolID:   poolID,
		Status:       status,
		CreationDate: time.Now().UTC(),
	}
	if err := s.CreateUserImportJob(job); err != nil {
		t.Fatalf("create job: %v", err)
	}
	return job
}

// StartUserImportJobIfEligible moves a Created job to Pending and refuses
// a second start while any job in the account is active.
func TestStartUserImportJobIfEligible(t *testing.T) {
	s := newImportJobTestStore(t)
	seedImportJob(t, s, "pool-a", "job-1", "Created")

	started, err := s.StartUserImportJobIfEligible("pool-a", "job-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if started.Status != "Pending" || started.StartDate.IsZero() {
		t.Fatalf("started job = %+v", started)
	}

	// A second active job blocks starting anything else.
	seedImportJob(t, s, "pool-b", "job-2", "Created")
	if _, err := s.StartUserImportJobIfEligible("pool-b", "job-2"); !errors.Is(err, ErrImportJobActiveExists) {
		t.Fatalf("want ErrImportJobActiveExists, got %v", err)
	}

	// Starting the already-started job again conflicts.
	if _, err := s.StartUserImportJobIfEligible("pool-a", "job-1"); !errors.Is(err, ErrImportJobStatusConflict) {
		t.Fatalf("want ErrImportJobStatusConflict, got %v", err)
	}
}

// The region's active-import marker is released when the job it names
// reaches a terminal status, so the start guard does not refuse new jobs
// on account of a finished one.
func TestStartUserImportJobMarkerReleasedOnTerminalStatus(t *testing.T) {
	s := newImportJobTestStore(t)
	seedImportJob(t, s, "pool-a", "job-1", "Created")

	if _, err := s.StartUserImportJobIfEligible("pool-a", "job-1"); err != nil {
		t.Fatalf("start job-1: %v", err)
	}
	seedImportJob(t, s, "pool-b", "job-2", "Created")
	if _, err := s.StartUserImportJobIfEligible("pool-b", "job-2"); !errors.Is(err, ErrImportJobActiveExists) {
		t.Fatalf("want ErrImportJobActiveExists while job-1 is Pending, got %v", err)
	}

	if _, err := s.TransitionUserImportJobStatus("pool-a", "job-1", "Pending", "Failed", nil); err != nil {
		t.Fatalf("terminalise job-1: %v", err)
	}
	started, err := s.StartUserImportJobIfEligible("pool-b", "job-2")
	if err != nil {
		t.Fatalf("start after marker release: %v", err)
	}
	if started.Status != "Pending" {
		t.Fatalf("job-2 status = %q, want Pending", started.Status)
	}
}

// A marker naming a job that no longer exists (a pool cascade wiped it) or
// never became active is stale: the next start takes the slot over instead
// of being blocked by a dead pointer.
func TestStartUserImportJobStaleMarkerIsTakenOver(t *testing.T) {
	s := newImportJobTestStore(t)
	if err := s.userImportJobsStore.Put(activeImportJobKey, "pool-gone#job-gone"); err != nil {
		t.Fatalf("seed stale marker: %v", err)
	}

	seedImportJob(t, s, "pool-a", "job-1", "Created")
	started, err := s.StartUserImportJobIfEligible("pool-a", "job-1")
	if err != nil {
		t.Fatalf("start over stale marker: %v", err)
	}
	if started.Status != "Pending" {
		t.Fatalf("status = %q, want Pending", started.Status)
	}
	var activeKey string
	if err := s.userImportJobsStore.Get(activeImportJobKey, &activeKey); err != nil || activeKey != "pool-a#job-1" {
		t.Fatalf("marker = (%v, %q), want pool-a#job-1", err, activeKey)
	}
}

// TransitionUserImportJobStatus only succeeds from the expected state and
// applies the mutation atomically with the status change.
func TestTransitionUserImportJobStatus(t *testing.T) {
	s := newImportJobTestStore(t)
	seedImportJob(t, s, "pool-a", "job-1", "InProgress")

	if _, err := s.TransitionUserImportJobStatus("pool-a", "job-1", "Created", "Pending", nil); !errors.Is(err, ErrImportJobStatusConflict) {
		t.Fatalf("want conflict from wrong source status, got %v", err)
	}

	final, err := s.TransitionUserImportJobStatus("pool-a", "job-1", "InProgress", "Succeeded", func(j *UserImportJob) {
		j.CompletionDate = time.Now().UTC()
		j.CompletionMessage = "done"
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if final.Status != "Succeeded" || final.CompletionMessage != "done" || final.CompletionDate.IsZero() {
		t.Fatalf("final job = %+v", final)
	}

	// A terminal job no longer transitions.
	if _, err := s.TransitionUserImportJobStatus("pool-a", "job-1", "InProgress", "Failed", nil); !errors.Is(err, ErrImportJobStatusConflict) {
		t.Fatalf("want conflict from terminal state, got %v", err)
	}
}

// UpdateUserImportJobProgress applies counter mutations only while the
// job is still running.
func TestUpdateUserImportJobProgress(t *testing.T) {
	s := newImportJobTestStore(t)
	seedImportJob(t, s, "pool-a", "job-1", "InProgress")

	if err := s.UpdateUserImportJobProgress("pool-a", "job-1", func(j *UserImportJob) {
		j.ImportedUsers++
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	job, err := s.GetUserImportJob("pool-a", "job-1")
	if err != nil {
		t.Fatalf("load job: %v", err)
	}
	if job.ImportedUsers != 1 {
		t.Fatalf("imported = %d, want 1", job.ImportedUsers)
	}

	if _, err := s.TransitionUserImportJobStatus("pool-a", "job-1", "InProgress", "Stopped", nil); err != nil {
		t.Fatalf("stop: %v", err)
	}
	if err := s.UpdateUserImportJobProgress("pool-a", "job-1", func(j *UserImportJob) {
		j.ImportedUsers++
	}); !errors.Is(err, ErrImportJobStatusConflict) {
		t.Fatalf("want conflict on terminal job, got %v", err)
	}
}

// Usernames are unique and resolvable case-insensitively unless the pool
// opts into case-sensitive usernames.
func TestUsernameCaseInsensitivity(t *testing.T) {
	s := newImportJobTestStore(t)
	pool, err := s.CreateUserPool(NewUserPool("ci-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}

	// Default configuration: case-insensitive.
	if err := s.CreateUser(NewUser(pool.ID, "John")); err != nil {
		t.Fatalf("create John: %v", err)
	}
	if err := s.CreateUser(NewUser(pool.ID, "john")); err == nil {
		t.Fatal("creating john after John must fail case-insensitively")
	}
	if _, err := s.GetUser(pool.ID, "JOHN"); err != nil {
		t.Fatalf("lookup with different case: %v", err)
	}

	// Deleting through any case frees the name again.
	if err := s.DeleteUser(pool.ID, "jOhN"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if err := s.CreateUser(NewUser(pool.ID, "john")); err != nil {
		t.Fatalf("recreate after delete: %v", err)
	}

	// A case-sensitive pool allows both spellings.
	csPool, err := s.CreateUserPool(NewUserPool("cs-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create cs pool: %v", err)
	}
	if err := s.UpdateUserPoolFunc(csPool.ID, func(p *UserPool) error {
		p.UsernameConfiguration = &UsernameConfiguration{CaseSensitive: true}
		return nil
	}); err != nil {
		t.Fatalf("update cs pool: %v", err)
	}
	if err := s.CreateUser(NewUser(csPool.ID, "Jane")); err != nil {
		t.Fatalf("create Jane: %v", err)
	}
	if err := s.CreateUser(NewUser(csPool.ID, "jane")); err != nil {
		t.Fatalf("case-sensitive pool must allow both spellings: %v", err)
	}
	if _, err := s.GetUser(csPool.ID, "JANE"); err == nil {
		t.Fatal("case-sensitive pool must not resolve JANE to Jane")
	}
}

// The cached username-case sensitivity does not outlive the pool record it
// was derived from: deletion evicts the entry, so a record later written
// under the same pool ID with a different configuration resolves to the
// fresh value instead of the cached one.
func TestUsernameCaseCacheEvictedOnPoolDeletion(t *testing.T) {
	s := newImportJobTestStore(t)
	pool, err := s.CreateUserPool(NewUserPool("evict-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	// Warm the cache through the username-keyed path (case-insensitive).
	if err := s.CreateUser(NewUser(pool.ID, "Warm")); err != nil {
		t.Fatalf("create user: %v", err)
	}
	if _, err := s.GetUser(pool.ID, "warm"); err != nil {
		t.Fatalf("lookup through warm cache: %v", err)
	}

	if err := s.DeleteUserPool(pool.ID); err != nil {
		t.Fatalf("delete pool: %v", err)
	}
	resurrected := NewUserPool("evict-pool-resurrected", "us-east-1")
	resurrected.ID = pool.ID
	resurrected.UsernameConfiguration = &UsernameConfiguration{CaseSensitive: true}
	if err := s.BaseStore.Put(pool.ID, resurrected); err != nil {
		t.Fatalf("resurrect pool record: %v", err)
	}
	if !s.usernameCaseSensitive(pool.ID) {
		t.Fatal("stale cached flag served after pool deletion")
	}
}
