package sfn

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/storage"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

func newResolverTestStore(t *testing.T) *sfnstore.StepFunctionStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	store := sfnstore.NewStepFunctionStore(st, "000000000000", "us-east-1")
	sm := &sfnstore.StateMachine{
		Name:       "resolver-sm",
		Definition: `{"StartAt":"A","States":{"A":{"Type":"Pass","End":true}}}`,
		RoleArn:    "arn:aws:iam::000000000000:role/sm",
	}
	if err := store.CreateStateMachine(context.Background(), sm); err != nil {
		t.Fatal(err)
	}
	return store
}

func TestResolveStateMachineReferenceBase(t *testing.T) {
	store := newResolverTestStore(t)
	arn := "arn:aws:states:us-east-1:000000000000:stateMachine:resolver-sm"

	ref, err := resolveStateMachineReference(context.Background(), store, arn)
	if err != nil {
		t.Fatalf("base resolution failed: %v", err)
	}
	if ref.Version != nil || ref.Alias != nil {
		t.Fatalf("unqualified ARN must not resolve to a version or alias")
	}
	if ref.definition() != ref.StateMachine.Definition {
		t.Fatalf("base reference must use the live definition")
	}
}

func TestResolveStateMachineReferenceVersion(t *testing.T) {
	store := newResolverTestStore(t)
	smArn := "arn:aws:states:us-east-1:000000000000:stateMachine:resolver-sm"

	version, err := store.PublishStateMachineVersion(context.Background(), smArn, "")
	if err != nil {
		t.Fatal(err)
	}

	ref, err := resolveStateMachineReference(context.Background(), store, version.StateMachineVersionArn)
	if err != nil {
		t.Fatalf("version resolution failed: %v", err)
	}
	if ref.Version == nil || ref.Version.StateMachineVersionArn != version.StateMachineVersionArn {
		t.Fatalf("version qualifier must resolve to the version record")
	}
	if ref.Alias != nil {
		t.Fatalf("version ARN must not resolve to an alias")
	}
}

func TestResolveStateMachineReferenceAlias(t *testing.T) {
	store := newResolverTestStore(t)
	smArn := "arn:aws:states:us-east-1:000000000000:stateMachine:resolver-sm"

	version, err := store.PublishStateMachineVersion(context.Background(), smArn, "")
	if err != nil {
		t.Fatal(err)
	}
	alias := &sfnstore.StateMachineAlias{
		StateMachineArn: smArn,
		Name:            "PROD",
		RoutingConfiguration: []sfnstore.RoutingConfiguration{
			{StateMachineVersionArn: version.StateMachineVersionArn, Weight: 100},
		},
	}
	if err := store.CreateStateMachineAlias(context.Background(), alias); err != nil {
		t.Fatal(err)
	}

	want := "arn:aws:states:us-east-1:000000000000:stateMachine:resolver-sm:PROD"
	if alias.StateMachineAliasArn != want {
		t.Fatalf("alias ARN = %q, want %q", alias.StateMachineAliasArn, want)
	}

	ref, err := resolveStateMachineReference(context.Background(), store, alias.StateMachineAliasArn)
	if err != nil {
		t.Fatalf("alias resolution failed: %v", err)
	}
	if ref.Alias == nil || ref.Alias.Name != "PROD" {
		t.Fatalf("alias qualifier must resolve to the alias record")
	}
	if ref.Version == nil || ref.Version.StateMachineVersionArn != version.StateMachineVersionArn {
		t.Fatalf("alias must route to the routed version")
	}
	if ref.definition() != version.Definition {
		t.Fatalf("alias reference must run the version snapshot")
	}
}

func TestResolveStateMachineReferenceUnknownQualifier(t *testing.T) {
	store := newResolverTestStore(t)

	unknownVersion := "arn:aws:states:us-east-1:000000000000:stateMachine:resolver-sm:99"
	if _, err := resolveStateMachineReference(context.Background(), store, unknownVersion); err == nil {
		t.Fatal("unknown version qualifier must fail")
	}

	unknownAlias := "arn:aws:states:us-east-1:000000000000:stateMachine:resolver-sm:NOSUCH"
	if _, err := resolveStateMachineReference(context.Background(), store, unknownAlias); err == nil {
		t.Fatal("unknown alias qualifier must fail")
	}
}

func TestResolveStateMachineReferenceMalformedAndLabel(t *testing.T) {
	store := newResolverTestStore(t)

	if _, err := resolveStateMachineReference(context.Background(), store, "arn:aws:iam::000000000000:role/x"); err == nil {
		t.Fatal("non-States ARN must be rejected")
	}
	label := "arn:aws:states:us-east-1:000000000000:stateMachine:resolver-sm/mapLabel"
	if _, err := resolveStateMachineReference(context.Background(), store, label); err == nil {
		t.Fatal("Distributed Map label ARN must be rejected")
	}
}

func TestSelectVersionByWeightSingleAndSplit(t *testing.T) {
	single := []sfnstore.RoutingConfiguration{{StateMachineVersionArn: "v1", Weight: 100}}
	for i := 0; i < 8; i++ {
		got, err := selectVersionByWeight(single)
		if err != nil || got != "v1" {
			t.Fatalf("single entry must be deterministic, got %q err %v", got, err)
		}
	}

	split := []sfnstore.RoutingConfiguration{
		{StateMachineVersionArn: "v1", Weight: 50},
		{StateMachineVersionArn: "v2", Weight: 50},
	}
	seen := map[string]bool{}
	for i := 0; i < 200 && len(seen) < 2; i++ {
		got, err := selectVersionByWeight(split)
		if err != nil {
			t.Fatalf("split selection failed: %v", err)
		}
		if got != "v1" && got != "v2" {
			t.Fatalf("unexpected version %q", got)
		}
		seen[got] = true
	}
	if len(seen) != 2 {
		t.Fatal("a 50/50 split must reach both versions over repeated picks")
	}
}

// TestDescribeExecutionRedriveStatusMatrix pins the derived redriveStatus
// across the whole contract matrix: the verdict comes from the same
// eligibility function RedriveExecution enforces, so an ABORTED execution
// (redrivable) and a Map Run child (redriven only by its parent Map Run)
// must not be reported through the wrong member.
func TestDescribeExecutionRedriveStatusMatrix(t *testing.T) {
	cases := []struct {
		name      string
		smType    string
		status    string
		stopAge   time.Duration
		mapRunArn string
		want      string
	}{
		{"failed standard", "STANDARD", "FAILED", 0, "", "REDRIVABLE"},
		{"timed out standard", "STANDARD", "TIMED_OUT", 0, "", "REDRIVABLE"},
		{"aborted standard", "STANDARD", "ABORTED", 0, "", "REDRIVABLE"},
		{"failed map child", "STANDARD", "FAILED", 0, "arn:aws:states:us-east-1:000000000000:mapRun:sm/M/1", "REDRIVABLE_BY_MAP_RUN"},
		{"timed out map child", "STANDARD", "TIMED_OUT", 0, "arn:aws:states:us-east-1:000000000000:mapRun:sm/M/1", "REDRIVABLE_BY_MAP_RUN"},
		{"pending map child", "STANDARD", "PENDING_REDRIVE", 0, "arn:aws:states:us-east-1:000000000000:mapRun:sm/M/1", "REDRIVABLE_BY_MAP_RUN"},
		{"succeeded map child", "STANDARD", "SUCCEEDED", 0, "arn:aws:states:us-east-1:000000000000:mapRun:sm/M/1", "NOT_REDRIVABLE"},
		{"window expired map child", "STANDARD", "FAILED", 15 * 24 * time.Hour, "arn:aws:states:us-east-1:000000000000:mapRun:sm/M/1", "NOT_REDRIVABLE"},
		{"running", "STANDARD", "RUNNING", 0, "", "NOT_REDRIVABLE"},
		{"succeeded", "STANDARD", "SUCCEEDED", 0, "", "NOT_REDRIVABLE"},
		{"express failed", "EXPRESS", "FAILED", 0, "", "NOT_REDRIVABLE"},
		{"window expired", "STANDARD", "FAILED", 15 * 24 * time.Hour, "", "NOT_REDRIVABLE"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc, store := newRecoveryService(t)
			ctx := t.Context()

			sm := &sfnstore.StateMachine{Name: "redrive-" + tc.name, Definition: recoveryDefinition, Type: tc.smType}
			if err := store.CreateStateMachine(ctx, sm); err != nil {
				t.Fatalf("create state machine: %v", err)
			}
			exec := sfnstore.NewExecution(sm.StateMachineArn, "exec-1", "{}", "")
			exec.ExecutionArn = sm.StateMachineArn + ":exec-1"
			exec.MapRunArn = tc.mapRunArn
			if err := store.CreateExecution(ctx, exec); err != nil {
				t.Fatalf("create execution: %v", err)
			}
			exec.Status = tc.status
			if tc.stopAge > 0 {
				exec.StopDate = time.Now().UTC().Add(-tc.stopAge)
			} else if tc.status != "RUNNING" && tc.status != "PENDING_REDRIVE" {
				exec.StopDate = time.Now().UTC()
			}
			if err := store.UpdateExecution(ctx, exec); err != nil {
				t.Fatalf("update execution: %v", err)
			}

			resp, err := svc.describeExecutionCore(ctx, store, DescribeExecutionInput{ExecutionArn: exec.ExecutionArn})
			if err != nil {
				t.Fatalf("describe execution: %v", err)
			}
			if resp["redriveStatus"] != tc.want {
				t.Errorf("redriveStatus = %v, want %s", resp["redriveStatus"], tc.want)
			}
			if tc.want == "NOT_REDRIVABLE" {
				if reason, ok := resp["redriveStatusReason"].(string); !ok || reason == "" {
					t.Errorf("NOT_REDRIVABLE must carry a redriveStatusReason, got %v", resp["redriveStatusReason"])
				}
			} else if _, present := resp["redriveStatusReason"]; present {
				t.Errorf("%s must not carry a redriveStatusReason, got %v", tc.want, resp["redriveStatusReason"])
			}
		})
	}
}

// TestRedriveExecutionSerialisesConcurrentTransitions pins the redrive
// transition's atomicity: concurrent RedriveExecution calls against the
// same terminal execution must produce exactly one transition — every
// other caller receives ExecutionNotRedrivable, the record is redriven
// once and a single ExecutionRedriven event lands in the history.
func TestRedriveExecutionSerialisesConcurrentTransitions(t *testing.T) {
	svc, store, arn := launchHistoryExecution(t, "hist-redrive-race",
		`{"StartAt":"W","States":{"W":{"Type":"Wait","Seconds":60,"End":true}}}`)
	t.Cleanup(func() {
		svc.stopExecutionCore(t.Context(), store, StopExecutionInput{ExecutionArn: arn})
		waitForExecution(t, store, arn)
	})

	// Stop mid-wait so the execution lands in the redrivable ABORTED
	// terminal status, then wait out the 60s wait on the redriven attempt.
	deadline := time.Now().Add(5 * time.Second)
	entered := false
	for !entered && time.Now().Before(deadline) {
		history, _, err := store.GetExecutionHistory(t.Context(), arn, 100, "", false)
		if err != nil {
			t.Fatalf("get history: %v", err)
		}
		for _, ev := range history {
			if ev.Type == "WaitStateEntered" {
				entered = true
				break
			}
		}
		if !entered {
			time.Sleep(10 * time.Millisecond)
		}
	}
	if !entered {
		t.Fatal("wait state was never entered before the stop")
	}
	if _, err := svc.stopExecutionCore(t.Context(), store, StopExecutionInput{ExecutionArn: arn}); err != nil {
		t.Fatalf("stop execution: %v", err)
	}
	if final := waitForExecution(t, store, arn); final.Status != "ABORTED" {
		t.Fatalf("status = %s, want ABORTED", final.Status)
	}

	const callers = 4
	errs := make([]error, callers)
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := range callers {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			_, errs[i] = svc.redriveExecutionCore(t.Context(), store, RedriveExecutionInput{ExecutionArn: arn})
		}(i)
	}
	close(start)
	wg.Wait()

	succeeded := 0
	for _, err := range errs {
		if err == nil {
			succeeded++
			continue
		}
		var awsErr *awserrors.AWSError
		if !errors.As(err, &awsErr) || awsErr.Code != "ExecutionNotRedrivable" {
			t.Fatalf("concurrent redrive error = %v, want ExecutionNotRedrivable", err)
		}
	}
	if succeeded != 1 {
		t.Fatalf("exactly one concurrent redrive must transition the execution, %d succeeded: %v", succeeded, errs)
	}

	if _, err := svc.stopExecutionCore(t.Context(), store, StopExecutionInput{ExecutionArn: arn}); err != nil {
		t.Fatalf("stop the redriven execution: %v", err)
	}
	final := waitForExecution(t, store, arn)
	if final.RedriveCount != 1 {
		t.Errorf("RedriveCount = %d, want exactly 1 across %d concurrent callers", final.RedriveCount, callers)
	}

	history, _, err := store.GetExecutionHistory(t.Context(), arn, 1000, "", false)
	if err != nil {
		t.Fatalf("get history: %v", err)
	}
	redriven := 0
	for _, ev := range history {
		if ev.Type == "ExecutionRedriven" {
			redriven++
		}
	}
	if redriven != 1 {
		t.Errorf("ExecutionRedriven events = %d, want 1", redriven)
	}
}

// TestRedrivenExecutionRemainsStoppable pins the registration generation
// guard: the draining goroutine of the stopped run must not delete the
// redrived run's registration, so a second StopExecution still cancels
// the resumed execution instead of leaving it running to its timeout.
func TestRedrivenExecutionRemainsStoppable(t *testing.T) {
	svc, store, arn := launchHistoryExecution(t, "hist-redrive-restop",
		`{"StartAt":"W","States":{"W":{"Type":"Wait","Seconds":60,"End":true}}}`)

	awaitWaitEntered := func(stage string) {
		t.Helper()
		deadline := time.Now().Add(5 * time.Second)
		for !t.Failed() && time.Now().Before(deadline) {
			history, _, err := store.GetExecutionHistory(t.Context(), arn, 100, "", false)
			if err != nil {
				t.Fatalf("%s: get history: %v", stage, err)
			}
			entered := 0
			for _, ev := range history {
				if ev.Type == "WaitStateEntered" {
					entered++
				}
			}
			if entered >= 1 {
				return
			}
			time.Sleep(10 * time.Millisecond)
		}
		if !t.Failed() {
			t.Fatalf("%s: the wait state was never entered", stage)
		}
	}

	awaitWaitEntered("initial run")
	if _, err := svc.stopExecutionCore(t.Context(), store, StopExecutionInput{ExecutionArn: arn}); err != nil {
		t.Fatalf("first stop: %v", err)
	}
	if final := waitForExecution(t, store, arn); final.Status != "ABORTED" {
		t.Fatalf("status after first stop = %s, want ABORTED", final.Status)
	}

	if _, err := svc.redriveExecutionCore(t.Context(), store, RedriveExecutionInput{ExecutionArn: arn}); err != nil {
		t.Fatalf("redrive: %v", err)
	}

	awaitWaitEntered("redriven run")
	if _, err := svc.stopExecutionCore(t.Context(), store, StopExecutionInput{ExecutionArn: arn}); err != nil {
		t.Fatalf("second stop: %v", err)
	}
	final := waitForExecution(t, store, arn)
	if final.Status != "ABORTED" {
		t.Fatalf("status after second stop = %s, want ABORTED (the redrived run must stay stoppable)", final.Status)
	}
	if final.RedriveCount != 1 {
		t.Fatalf("RedriveCount = %d, want 1 preserved across the second stop", final.RedriveCount)
	}
}

// TestRedriveExecutionRejectsMapRunChild pins the child side of the
// redrive contract: a Distributed Map child carries a mapRunArn and is
// redriven only through its parent Map Run, so a direct RedriveExecution
// against it is rejected with ExecutionNotRedrivable even in a redrivable
// terminal status.
func TestRedriveExecutionRejectsMapRunChild(t *testing.T) {
	svc, store := newRecoveryService(t)
	ctx := t.Context()

	sm := &sfnstore.StateMachine{Name: "redrive-child-sm", Definition: recoveryDefinition}
	if err := store.CreateStateMachine(ctx, sm); err != nil {
		t.Fatalf("create state machine: %v", err)
	}
	child := sfnstore.NewExecution(sm.StateMachineArn, "parent:M-0", `{"item":1}`, "")
	child.ExecutionArn = sm.StateMachineArn + ":parent:M-0"
	child.MapRunArn = sm.StateMachineArn + ":mapRun:M/1"
	if err := store.CreateExecution(ctx, child); err != nil {
		t.Fatalf("create child execution: %v", err)
	}
	child.Status = "FAILED"
	child.StopDate = time.Now().UTC()
	if err := store.UpdateExecution(ctx, child); err != nil {
		t.Fatalf("fail child execution: %v", err)
	}

	_, err := svc.redriveExecutionCore(ctx, store, RedriveExecutionInput{ExecutionArn: child.ExecutionArn})
	if err == nil {
		t.Fatal("direct redrive of a Map Run child must be rejected")
	}
	var awsErr *awserrors.AWSError
	if !errors.As(err, &awsErr) || awsErr.Code != "ExecutionNotRedrivable" {
		t.Fatalf("redrive child error = %v, want ExecutionNotRedrivable", err)
	}
	if !strings.Contains(awsErr.Message, "Map Run") {
		t.Errorf("rejection must state the Map Run ownership, got %q", awsErr.Message)
	}

	persisted, err := store.GetExecution(ctx, child.ExecutionArn)
	if err != nil {
		t.Fatalf("get child: %v", err)
	}
	if persisted.Status != "FAILED" || persisted.RedriveCount != 0 {
		t.Errorf("rejected redrive must leave the child untouched, got status %s redriveCount %d", persisted.Status, persisted.RedriveCount)
	}
}

// TestExpressNameReuseStartsFreshExecutions pins the EXPRESS name-reuse
// contract ("For EXPRESS workflows, execution names can be reused") on
// both start paths: a second same-name start persists a distinct execution
// under a suffixed ARN while the record keeps the caller's name, and a
// STANDARD same-name start with different input still answers
// ExecutionAlreadyExists.
func TestExpressNameReuseStartsFreshExecutions(t *testing.T) {
	svc, store := newRecoveryService(t)
	pass := `{"StartAt":"A","States":{"A":{"Type":"Pass","Result":"ok","End":true}}}`
	express := &sfnstore.StateMachine{Name: "reuse-express", Type: "EXPRESS", Definition: pass}
	if err := store.CreateStateMachine(t.Context(), express); err != nil {
		t.Fatalf("create express machine: %v", err)
	}

	first, err := svc.startExecutionCore(t.Context(), store, StartExecutionInput{
		StateMachineArn: express.StateMachineArn, Name: "nightly", Input: `{}`})
	if err != nil {
		t.Fatalf("first express start: %v", err)
	}
	second, err := svc.startExecutionCore(t.Context(), store, StartExecutionInput{
		StateMachineArn: express.StateMachineArn, Name: "nightly", Input: `{}`})
	if err != nil {
		t.Fatalf("express name reuse: %v", err)
	}
	if second.ExecutionArn == first.ExecutionArn {
		t.Error("the reused express name must start a distinct execution, not the same ARN")
	}
	for i, arn := range []string{first.ExecutionArn, second.ExecutionArn} {
		rec, gerr := store.GetExecution(t.Context(), arn)
		if gerr != nil || rec == nil {
			t.Fatalf("get execution %d: %v", i, gerr)
		}
		if rec.Name != "nightly" {
			t.Errorf("execution %d name = %s, want the caller's nightly", i, rec.Name)
		}
	}

	syncFirst, err := svc.startSyncExecutionCore(t.Context(), store, StartSyncExecutionInput{
		StateMachineArn: express.StateMachineArn, Name: "syncname", Input: `{}`})
	if err != nil {
		t.Fatalf("first sync start: %v", err)
	}
	syncSecond, err := svc.startSyncExecutionCore(t.Context(), store, StartSyncExecutionInput{
		StateMachineArn: express.StateMachineArn, Name: "syncname", Input: `{}`})
	if err != nil {
		t.Fatalf("sync name reuse: %v", err)
	}
	if syncSecond["executionArn"] == syncFirst["executionArn"] {
		t.Error("the reused sync name must start a distinct execution, not the same ARN")
	}
	if syncSecond["name"] != "syncname" {
		t.Errorf("sync name = %v, want the caller's syncname", syncSecond["name"])
	}

	standard := &sfnstore.StateMachine{Name: "reuse-std", Definition: pass}
	if err := store.CreateStateMachine(t.Context(), standard); err != nil {
		t.Fatalf("create standard machine: %v", err)
	}
	daily, err := svc.startExecutionCore(t.Context(), store, StartExecutionInput{
		StateMachineArn: standard.StateMachineArn, Name: "daily", Input: `{"v":1}`})
	if err != nil {
		t.Fatalf("first standard start: %v", err)
	}
	_, err = svc.startExecutionCore(t.Context(), store, StartExecutionInput{
		StateMachineArn: standard.StateMachineArn, Name: "daily", Input: `{"v":2}`})
	var awsErr *awserrors.AWSError
	if !errors.As(err, &awsErr) || awsErr.Code != "ExecutionAlreadyExists" {
		t.Fatalf("standard same-name different-input start = %v, want ExecutionAlreadyExists", err)
	}

	// The asynchronously launched executions are Pass machines; let them
	// finish before the harness closes the store.
	for _, arn := range []string{first.ExecutionArn, second.ExecutionArn, daily.ExecutionArn} {
		waitForExecution(t, store, arn)
	}
}

// TestDescribeExecutionOutputDetailsPresence pins the details presence
// contract at the Core level: in the default mode outputDetails appears
// only when an output exists (a RUNNING execution has none), and
// METADATA_ONLY reports both details members as included:false with the
// payloads withheld. The sync start follows the same presence rule.
func TestDescribeExecutionOutputDetailsPresence(t *testing.T) {
	svc, store := newRecoveryService(t)
	sm := &sfnstore.StateMachine{Name: "details-presence",
		Definition: `{"StartAt":"A","States":{"A":{"Type":"Pass","Result":"ok","End":true}}}`}
	if err := store.CreateStateMachine(t.Context(), sm); err != nil {
		t.Fatalf("create machine: %v", err)
	}
	exec := sfnstore.NewExecution(sm.StateMachineArn, "running-1", `{"in":1}`, "")
	exec.ExecutionArn = sm.StateMachineArn + ":running-1"
	if err := store.CreateExecution(t.Context(), exec); err != nil {
		t.Fatalf("create execution: %v", err)
	}

	resp, err := svc.describeExecutionCore(t.Context(), store, DescribeExecutionInput{ExecutionArn: exec.ExecutionArn})
	if err != nil {
		t.Fatalf("describe: %v", err)
	}
	if _, present := resp["outputDetails"]; present {
		t.Error("outputDetails must be absent while the execution has no output")
	}
	in, ok := resp["inputDetails"].(map[string]interface{})
	if !ok || in["included"] != true {
		t.Errorf("inputDetails = %v, want included:true (input always exists)", resp["inputDetails"])
	}

	meta, err := svc.describeExecutionCore(t.Context(), store, DescribeExecutionInput{
		ExecutionArn: exec.ExecutionArn, IncludedData: "METADATA_ONLY"})
	if err != nil {
		t.Fatalf("describe metadata-only: %v", err)
	}
	if _, present := meta["input"]; present {
		t.Error("METADATA_ONLY must withhold the input payload")
	}
	for _, member := range []string{"inputDetails", "outputDetails"} {
		details, ok := meta[member].(map[string]interface{})
		if !ok || details["included"] != false {
			t.Errorf("%s = %v, want included:false under METADATA_ONLY", member, meta[member])
		}
	}

	express := &sfnstore.StateMachine{Name: "details-sync", Type: "EXPRESS",
		Definition: `{"StartAt":"F","States":{"F":{"Type":"Fail","Error":"E","Cause":"c"}}}`}
	if err := store.CreateStateMachine(t.Context(), express); err != nil {
		t.Fatalf("create express machine: %v", err)
	}
	sync, err := svc.startSyncExecutionCore(t.Context(), store, StartSyncExecutionInput{
		StateMachineArn: express.StateMachineArn, Input: `{}`})
	if err != nil {
		t.Fatalf("sync start: %v", err)
	}
	if sync["status"] != "FAILED" {
		t.Fatalf("sync status = %v, want FAILED", sync["status"])
	}
	if _, present := sync["outputDetails"]; present {
		t.Error("a failed sync execution has no output — outputDetails must be absent")
	}
	if _, present := sync["inputDetails"]; !present {
		t.Error("inputDetails must always be reported")
	}
}

// failingGetBucket fails every Get with the injected error, standing in
// for a transient storage fault that is not the not-found condition.
type failingGetBucket struct {
	storage.Bucket
	fail error
}

func (b *failingGetBucket) Get(key []byte) ([]byte, error) { return nil, b.fail }

// transientFaultStorage routes one bucket's reads through the fault.
type transientFaultStorage struct {
	storage.BasicStorage
	failOn string
	fail   error
}

func (s *transientFaultStorage) Bucket(name string) storage.Bucket {
	b := s.BasicStorage.Bucket(name)
	if name == s.failOn {
		return &failingGetBucket{Bucket: b, fail: s.fail}
	}
	return b
}

// TestExecutionHistoryTransientFaultIsNotDoesNotExist pins the error class
// on both read Cores: a storage fault that is not the not-found sentinel
// surfaces as itself, never as the permanent ExecutionDoesNotExist verdict
// a client might answer by deleting resources.
func TestExecutionHistoryTransientFaultIsNotDoesNotExist(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	wrapped := &transientFaultStorage{
		BasicStorage: st,
		failOn:       "stepfunction-executions-us-east-1",
		fail:         errors.New("pebble: transient io fault"),
	}
	store := sfnstore.NewStepFunctionStore(wrapped, "000000000000", "us-east-1")
	svc := &StepFunctionService{}

	arn := "arn:aws:states:us-east-1:000000000000:execution:sm:none"
	_, histErr := svc.getExecutionHistoryCore(context.Background(), store, GetExecutionHistoryInput{ExecutionArn: arn})
	if histErr == nil {
		t.Fatal("a transient storage fault must fail the history read")
	}
	var awsErr *awserrors.AWSError
	if errors.As(histErr, &awsErr) && awsErr.Code == "ExecutionDoesNotExist" {
		t.Fatalf("history read misreported the transient fault as ExecutionDoesNotExist: %v", histErr)
	}

	_, runsErr := svc.listMapRunsCore(context.Background(), store, arn, 0, "")
	if runsErr == nil {
		t.Fatal("a transient storage fault must fail the map-run listing")
	}
	if errors.As(runsErr, &awsErr) && awsErr.Code == "ExecutionDoesNotExist" {
		t.Fatalf("map-run listing misreported the transient fault as ExecutionDoesNotExist: %v", runsErr)
	}
}

// TestStoreFaultsAnswerAsThemselves pins the fault class across the getter
// families, the task mutation paths and the redrive transition: a storage
// fault that is not the not-found sentinel surfaces as itself, never as a
// permanent-absence verdict a client might answer by deleting resources.
func TestStoreFaultsAnswerAsThemselves(t *testing.T) {
	fault := errors.New("pebble: transient io fault")
	newFaultedStore := func(t *testing.T, bucket string) *sfnstore.StepFunctionStore {
		st, err := storage.Open(t.TempDir())
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { st.Close() })
		return sfnstore.NewStepFunctionStore(&transientFaultStorage{
			BasicStorage: st,
			failOn:       bucket,
			fail:         fault,
		}, "000000000000", "us-east-1")
	}
	ctx := context.Background()
	svc := &StepFunctionService{}
	var awsErr *awserrors.AWSError

	// The store getters and the task mutations pass the fault through
	// instead of minting their not-found sentinels.
	storeLevel := []struct {
		name     string
		bucket   string
		call     func(*sfnstore.StepFunctionStore) error
		sentinel error
	}{
		{"GetStateMachine", "stepfunction-statemachines-us-east-1",
			func(s *sfnstore.StepFunctionStore) error {
				_, err := s.GetStateMachine(ctx, "arn:aws:states:us-east-1:000000000000:stateMachine:sm")
				return err
			},
			sfnstore.ErrStateMachineNotFound},
		{"GetActivity", "stepfunction-activities-us-east-1",
			func(s *sfnstore.StepFunctionStore) error {
				_, err := s.GetActivity(ctx, "arn:aws:states:us-east-1:000000000000:activity:a")
				return err
			},
			sfnstore.ErrActivityNotFound},
		{"GetActivityTaskByToken", "stepfunction-tasks-us-east-1",
			func(s *sfnstore.StepFunctionStore) error { _, err := s.GetActivityTaskByToken("token"); return err },
			sfnstore.ErrTaskNotFound},
		{"HeartbeatActivityTask", "stepfunction-tasks-us-east-1",
			func(s *sfnstore.StepFunctionStore) error { return s.HeartbeatActivityTask("token") },
			sfnstore.ErrTaskNotFound},
		{"CompleteActivityTask", "stepfunction-tasks-us-east-1",
			func(s *sfnstore.StepFunctionStore) error { return s.CompleteActivityTask("token", "{}") },
			sfnstore.ErrTaskNotFound},
		{"FailActivityTask", "stepfunction-tasks-us-east-1",
			func(s *sfnstore.StepFunctionStore) error { return s.FailActivityTask("token", "E", "c") },
			sfnstore.ErrTaskNotFound},
		{"WaitForTaskResultHeartbeatRecheck", "stepfunction-tasks-us-east-1",
			// The heartbeat re-read must fail as the fault itself: with the
			// tasks bucket faulted the re-read cannot distinguish absence,
			// so reporting TaskNotFound would blame the worker's token for
			// a server-side storage fault. The hbTimer fires well before
			// the overall timeout, so the call returns promptly.
			func(s *sfnstore.StepFunctionStore) error {
				_, err := s.WaitForTaskResult(ctx, "token", 10*time.Second, 5*time.Millisecond)
				return err
			},
			sfnstore.ErrTaskNotFound},
		{"GetMapRun", "stepfunction-mapruns-us-east-1",
			func(s *sfnstore.StepFunctionStore) error {
				_, err := s.GetMapRun(ctx, "arn:aws:states:us-east-1:000000000000:mapRun:sm/m")
				return err
			},
			sfnstore.ErrMapRunNotFound},
		{"GetStateMachineVersion", "stepfunction-versions-us-east-1",
			func(s *sfnstore.StepFunctionStore) error {
				_, err := s.GetStateMachineVersion(ctx, "arn:aws:states:us-east-1:000000000000:stateMachine:sm:1")
				return err
			},
			sfnstore.ErrStateMachineVersionNotFound},
		{"GetStateMachineAlias", "stepfunction-aliases-us-east-1",
			func(s *sfnstore.StepFunctionStore) error {
				_, err := s.GetStateMachineAlias(ctx, "arn:aws:states:us-east-1:000000000000:stateMachine:sm:prod")
				return err
			},
			sfnstore.ErrStateMachineAliasNotFound},
		{"TransitionExecutionForRedrive", "stepfunction-executions-us-east-1",
			func(s *sfnstore.StepFunctionStore) error {
				_, err := s.TransitionExecutionForRedrive(ctx, "arn:aws:states:us-east-1:000000000000:execution:sm:e",
					func(fresh *sfnstore.Execution) bool { return true },
					func(fresh *sfnstore.Execution) {})
				return err
			},
			sfnstore.ErrExecutionNotFound},
	}
	for _, tt := range storeLevel {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call(newFaultedStore(t, tt.bucket))
			if err == nil {
				t.Fatal("a transient storage fault must fail the read")
			}
			if errors.Is(err, tt.sentinel) {
				t.Fatalf("transient fault misreported as the absence sentinel: %v", err)
			}
			if !strings.Contains(err.Error(), fault.Error()) {
				t.Fatalf("the fault must surface as itself, got: %v", err)
			}
		})
	}

	// The Core surfaces answer the fault with the raw error, not with
	// their DoesNotExist-family verdicts.
	_, startErr := svc.startExecutionCore(ctx, newFaultedStore(t, "stepfunction-statemachines-us-east-1"),
		StartExecutionInput{StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm", Input: "{}"})
	if startErr == nil || (errors.As(startErr, &awsErr) && awsErr.Code == "StateMachineDoesNotExist") {
		t.Fatalf("StartExecution misreported the transient fault as StateMachineDoesNotExist: %v", startErr)
	}

	_, actErr := svc.describeActivityCore(ctx, newFaultedStore(t, "stepfunction-activities-us-east-1"),
		"arn:aws:states:us-east-1:000000000000:activity:a")
	if actErr == nil || (errors.As(actErr, &awsErr) && awsErr.Code == "ActivityDoesNotExist") {
		t.Fatalf("DescribeActivity misreported the transient fault as ActivityDoesNotExist: %v", actErr)
	}

	_, runErr := svc.describeMapRunCore(ctx, newFaultedStore(t, "stepfunction-mapruns-us-east-1"),
		"arn:aws:states:us-east-1:000000000000:mapRun:sm/m")
	if runErr == nil || (errors.As(runErr, &awsErr) && awsErr.Code == "ResourceNotFound") {
		t.Fatalf("DescribeMapRun misreported the transient fault as a missing Map Run: %v", runErr)
	}

	hbErr := svc.sendTaskHeartbeatCore(ctx, newFaultedStore(t, "stepfunction-tasks-us-east-1"), "token")
	if hbErr == nil || errors.Is(hbErr, ErrTaskDoesNotExist) {
		t.Fatalf("SendTaskHeartbeat misreported the transient fault as TaskDoesNotExist: %v", hbErr)
	}

	_, aliasErr := svc.describeStateMachineAliasCore(ctx, newFaultedStore(t, "stepfunction-aliases-us-east-1"),
		"arn:aws:states:us-east-1:000000000000:stateMachine:sm:prod")
	if aliasErr == nil || (errors.As(aliasErr, &awsErr) && awsErr.Code == "ResourceNotFound") {
		t.Fatalf("DescribeStateMachineAlias misreported the transient fault as a missing alias: %v", aliasErr)
	}
}
