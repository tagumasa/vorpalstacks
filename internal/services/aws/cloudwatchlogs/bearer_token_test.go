package cloudwatchlogs

import (
	"strings"
	"testing"

	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// newBearerTestService builds a service over temporary storage carrying one
// log group, returning the service and the group name.
func newBearerTestService(t *testing.T) (*LogsService, string) {
	t.Helper()
	svc, store := newTestService(t)
	const group = "bearer-group"
	createTestLogGroup(t, store, group)
	return svc, group
}

func TestPutBearerTokenAuthentication(t *testing.T) {
	svc, group := newBearerTestService(t)
	store, err := svc.getLogsStoreByRegion("us-east-1")
	if err != nil {
		t.Fatalf("logs store: %v", err)
	}

	enabled := func(set bool, value bool) PutBearerTokenAuthenticationInput {
		return PutBearerTokenAuthenticationInput{
			LogGroupIdentifier:               group,
			BearerTokenAuthenticationEnabled: value,
			BearerTokenAuthenticationSet:     set,
			Region:                           "us-east-1",
		}
	}

	// Both members are required: an absent boolean rejects rather than
	// silently writing false, and an absent identifier rejects too.
	if err := svc.putBearerTokenAuthenticationCore(enabled(false, false)); err == nil {
		t.Fatal("absent bearerTokenAuthenticationEnabled accepted")
	} else if logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("absent enabled member: %v, want InvalidParameterException", err)
	}
	if err := svc.putBearerTokenAuthenticationCore(PutBearerTokenAuthenticationInput{
		BearerTokenAuthenticationEnabled: true,
		BearerTokenAuthenticationSet:     true,
		Region:                           "us-east-1",
	}); err == nil {
		t.Fatal("absent logGroupIdentifier accepted")
	}

	// An unknown group answers the operation's ResourceNotFoundException.
	err = svc.putBearerTokenAuthenticationCore(PutBearerTokenAuthenticationInput{
		LogGroupIdentifier:               "no-such-group",
		BearerTokenAuthenticationEnabled: true,
		BearerTokenAuthenticationSet:     true,
		Region:                           "us-east-1",
	})
	if err == nil {
		t.Fatal("unknown log group accepted")
	}
	if logsErrorCode(err) != "ResourceNotFoundException" {
		t.Fatalf("unknown log group: %v, want ResourceNotFoundException", err)
	}

	// Round-trip: a group that never enabled reports disabled; enabling
	// reports enabled; an explicit false (the documented disable
	// operation) reports disabled again. The switch rides the group
	// record, so the reads go through the record itself.
	switchState := func() bool {
		lg, err := store.GetLogGroup(group)
		if err != nil {
			t.Fatalf("read group record: %v", err)
		}
		return lg.BearerTokenAuthenticationEnabled
	}
	if switchState() {
		t.Fatal("never-enabled group reports bearer token authentication enabled")
	}
	if err := svc.putBearerTokenAuthenticationCore(enabled(true, true)); err != nil {
		t.Fatalf("enable: %v", err)
	}
	if !switchState() {
		t.Fatal("after enable: the group record still reports disabled")
	}
	if err := svc.putBearerTokenAuthenticationCore(enabled(true, false)); err != nil {
		t.Fatalf("disable: %v", err)
	}
	if switchState() {
		t.Fatal("after disable: the group record still reports enabled")
	}
}

func TestPutBearerTokenAuthenticationARNIdentifier(t *testing.T) {
	svc, group := newBearerTestService(t)
	// The identifier member is name-or-ARN: the group's ARN resolves to
	// the same per-group record.
	arn := "arn:aws:logs:us-east-1:000000000000:log-group:" + group
	if err := svc.putBearerTokenAuthenticationCore(PutBearerTokenAuthenticationInput{
		LogGroupIdentifier:               arn,
		BearerTokenAuthenticationEnabled: true,
		BearerTokenAuthenticationSet:     true,
		Region:                           "us-east-1",
	}); err != nil {
		t.Fatalf("enable by ARN: %v", err)
	}
	store, err := svc.getLogsStoreByRegion("us-east-1")
	if err != nil {
		t.Fatalf("logs store: %v", err)
	}
	lg, err := store.GetLogGroup(group)
	if err != nil {
		t.Fatalf("read group record: %v", err)
	}
	if !lg.BearerTokenAuthenticationEnabled {
		t.Fatal("after enable by ARN: the group record still reports disabled")
	}
}

func TestDeleteLogGroupTearsDownBearerTokenRecord(t *testing.T) {
	svc, group := newBearerTestService(t)
	store, err := svc.getLogsStoreByRegion("us-east-1")
	if err != nil {
		t.Fatalf("logs store: %v", err)
	}
	if err := svc.putBearerTokenAuthenticationCore(PutBearerTokenAuthenticationInput{
		LogGroupIdentifier:               group,
		BearerTokenAuthenticationEnabled: true,
		BearerTokenAuthenticationSet:     true,
		Region:                           "us-east-1",
	}); err != nil {
		t.Fatalf("enable: %v", err)
	}
	if err := store.DeleteLogGroup(group); err != nil {
		t.Fatalf("delete log group: %v", err)
	}
	// A same-named recreation starts with bearer token authentication
	// disabled: the switch died with the group record.
	if err := store.CreateLogGroup(logsstore.NewLogGroup(group, "us-east-1", "000000000000")); err != nil {
		t.Fatalf("recreate log group: %v", err)
	}
	recreated, err := store.GetLogGroup(group)
	if err != nil {
		t.Fatalf("read recreated group record: %v", err)
	}
	if recreated.BearerTokenAuthenticationEnabled {
		t.Fatal("after recreation: the fresh record reports bearer token authentication enabled")
	}
}

// The logGroupIdentifier member rides the shared LogGroupIdentifier
// shape (1..2048): a long-but-valid ARN must be accepted, and only the
// shape's own 2048 ceiling rejects.
func TestPutBearerTokenIdentifierLength(t *testing.T) {
	svc, _ := newBearerTestService(t)
	store, err := svc.getLogsStoreByRegion("us-east-1")
	if err != nil {
		t.Fatalf("logs store: %v", err)
	}
	longName := strings.Repeat("g", 500)
	if err := store.CreateLogGroup(logsstore.NewLogGroup(longName, "us-east-1", "000000000000")); err != nil {
		t.Fatalf("create long-named group: %v", err)
	}
	arn := store.ARNBuilder().CloudWatch().LogGroup(longName)
	if len(arn) <= 512 {
		t.Fatalf("test setup: ARN length %d does not exercise the old 512 ceiling", len(arn))
	}
	if err := svc.putBearerTokenAuthenticationCore(PutBearerTokenAuthenticationInput{
		LogGroupIdentifier:               arn,
		BearerTokenAuthenticationEnabled: true,
		BearerTokenAuthenticationSet:     true,
		Region:                           "us-east-1",
	}); err != nil {
		t.Fatalf("valid long ARN rejected: %v", err)
	}
	over := strings.Repeat("g", logsstore.MaxLogGroupIdentifierLength+1)
	err = svc.putBearerTokenAuthenticationCore(PutBearerTokenAuthenticationInput{
		LogGroupIdentifier:               over,
		BearerTokenAuthenticationEnabled: true,
		BearerTokenAuthenticationSet:     true,
		Region:                           "us-east-1",
	})
	if err == nil || !strings.Contains(err.Error(), "2048") {
		t.Fatalf("a 2049-character identifier must reject on the shape ceiling, got %v", err)
	}
}
