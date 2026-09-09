package s3

import (
	"strings"
	"testing"
	"time"

	s3store "vorpalstacks/internal/store/aws/s3"
)

// lifecycleDaysInstant follows the documented AWS timing calculation — days
// added to the base time and rounded up to the next midnight UTC — and the
// age bound keeps the real day unit in production while TEST_MODE compresses
// one day to one second, so a Days=1 rule acts within a test run instead of
// after 24 hours.
func TestLifecycleDaysInstant(t *testing.T) {
	origUnit := lifecycleDayUnit
	defer func() { lifecycleDayUnit = origUnit }()

	lifecycleDayUnit = 24 * time.Hour

	// The documented example: created 1/15/2014 10:30 UTC with 3 days →
	// 1/19/2014 00:00 UTC.
	created := time.Date(2014, 1, 15, 10, 30, 0, 0, time.UTC)
	if got := lifecycleDaysInstant(created, 3); !got.Equal(time.Date(2014, 1, 19, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("documented example = %v, want 2014-01-19T00:00:00Z", got)
	}

	// A result that already lands exactly on midnight stays there.
	if got := lifecycleDaysInstant(time.Date(2026, 9, 9, 0, 0, 0, 0, time.UTC), 1); !got.Equal(time.Date(2026, 9, 10, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("exact-midnight result = %v, want it unchanged", got)
	}
	// A base one nanosecond into the day rounds a full day ahead.
	if got := lifecycleDaysInstant(time.Date(2026, 9, 9, 0, 0, 0, 1, time.UTC), 1); !got.Equal(time.Date(2026, 9, 11, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("just-past-midnight base = %v, want 2026-09-11T00:00:00Z", got)
	}
	// The raw age bound is the candidate pre-filter edge.
	if got := lifecycleAgeBound(3, created); !got.Equal(created.AddDate(0, 0, -3)) {
		t.Fatalf("production age bound = %v, want now-3d", got)
	}

	// The compressed unit keeps the calculation shape: its "midnights" are
	// the second boundaries.
	lifecycleDayUnit = time.Second
	base := time.Date(2026, 9, 8, 12, 0, 0, 400, time.UTC)
	if got := lifecycleDaysInstant(base, 1); !got.Equal(time.Date(2026, 9, 8, 12, 0, 2, 0, time.UTC)) {
		t.Fatalf("compressed instant = %v, want 12:00:02", got)
	}
	if got := lifecycleAgeBound(3, base); !got.Equal(base.Add(-3 * time.Second)) {
		t.Fatalf("compressed age bound = %v, want now-3s", got)
	}
}

// filterPrefix prefers the And operator's prefix over the plain prefix, and
// an absent filter yields no prefix constraint.
func TestFilterPrefix(t *testing.T) {
	if got := filterPrefix(nil); got != "" {
		t.Fatalf("nil filter prefix = %q, want empty", got)
	}
	if got := filterPrefix(&s3store.LifecycleRuleFilter{Prefix: "plain/"}); got != "plain/" {
		t.Fatalf("plain prefix = %q, want plain/", got)
	}
	and := &s3store.LifecycleRuleFilter{Prefix: "ignored/", And: &s3store.LifecycleRuleAndOperator{Prefix: "and/"}}
	if got := filterPrefix(and); got != "and/" {
		t.Fatalf("And prefix = %q, want and/ (And takes precedence)", got)
	}
}

// The waterfall must accept exactly the documented source-to-destination
// moves and forbid every regression or sideways move.
func TestLifecycleTransitionAllowed(t *testing.T) {
	cases := []struct {
		from, to s3store.ObjectStorageClass
		want     bool
	}{
		{s3store.StorageClassStandard, s3store.StorageClassStandardIA, true},
		{s3store.StorageClassStandard, s3store.StorageClassIntelligentTiering, true},
		{s3store.StorageClassStandard, s3store.StorageClassOneZoneIA, true},
		{s3store.StorageClassStandard, s3store.StorageClassGlacierIR, true},
		{s3store.StorageClassStandard, s3store.StorageClassGlacier, true},
		{s3store.StorageClassStandard, s3store.StorageClassDeepArchive, true},
		{s3store.StorageClassStandard, s3store.StorageClassStandard, false},
		{s3store.StorageClassStandard, s3store.StorageClassReducedRedundancy, false},
		{s3store.StorageClassReducedRedundancy, s3store.StorageClassGlacier, true},
		{s3store.StorageClassStandardIA, s3store.StorageClassGlacierIR, true},
		{s3store.StorageClassIntelligentTiering, s3store.StorageClassOneZoneIA, true},
		{s3store.StorageClassOneZoneIA, s3store.StorageClassGlacierIR, false},
		{s3store.StorageClassOneZoneIA, s3store.StorageClassGlacier, true},
		{s3store.StorageClassGlacierIR, s3store.StorageClassGlacier, true},
		{s3store.StorageClassGlacier, s3store.StorageClassDeepArchive, true},
		{s3store.StorageClassGlacier, s3store.StorageClassStandardIA, false},
		{s3store.StorageClassDeepArchive, s3store.StorageClassGlacier, false},
	}
	for _, c := range cases {
		if got := lifecycleTransitionAllowed(c.from, c.to); got != c.want {
			t.Fatalf("transition %s -> %s = %v, want %v", c.from, c.to, got, c.want)
		}
	}
}

// Entries replay in rule order from the object's present class: a due entry
// whose destination the waterfall forbids from the class reached so far is
// skipped, so the object never regresses to an earlier class. Day windows
// gate on the rounded midnight instant, never the raw age.
func TestResolveTransitionTarget(t *testing.T) {
	origUnit := lifecycleDayUnit
	defer func() { lifecycleDayUnit = origUnit }()
	lifecycleDayUnit = 24 * time.Hour

	now := time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC)
	old := now.Add(-10 * 24 * time.Hour)

	regression := []transitionSpec{
		{days: ptrInt32(1), class: s3store.StorageClassGlacier},
		{days: ptrInt32(2), class: s3store.StorageClassStandardIA},
	}
	if got := resolveTransitionTarget(s3store.StorageClassStandard, regression, old, now); got != s3store.StorageClassGlacier {
		t.Fatalf("both entries due must not regress: got %s, want GLACIER", got)
	}

	chained := []transitionSpec{
		{days: ptrInt32(1), class: s3store.StorageClassStandardIA},
		{days: ptrInt32(2), class: s3store.StorageClassGlacier},
	}
	if got := resolveTransitionTarget(s3store.StorageClassStandard, chained, old, now); got != s3store.StorageClassGlacier {
		t.Fatalf("chained entries must land on the deepest due class: got %s, want GLACIER", got)
	}
	// 36 hours old: the 1-day window has rounded past (yesterday 00:00),
	// the 2-day window has not (tomorrow 00:00).
	young := now.Add(-36 * time.Hour)
	if got := resolveTransitionTarget(s3store.StorageClassStandard, chained, young, now); got != s3store.StorageClassStandardIA {
		t.Fatalf("only the first entry due: got %s, want STANDARD_IA", got)
	}
	// 35 hours old is inside the raw 1-day age but before the rounded
	// midnight instant: no entry is due yet.
	rawOnly := now.Add(-35 * time.Hour)
	if got := resolveTransitionTarget(s3store.StorageClassStandard, chained, rawOnly, now); got != s3store.StorageClassStandard {
		t.Fatalf("raw age without the rounded instant must not be due: got %s, want STANDARD", got)
	}

	past := now.Add(-time.Hour)
	future := now.Add(time.Hour)
	dateSpecs := []transitionSpec{{date: &past, class: s3store.StorageClassGlacierIR}}
	if got := resolveTransitionTarget(s3store.StorageClassStandard, dateSpecs, now, now); got != s3store.StorageClassGlacierIR {
		t.Fatalf("a passed date entry is due immediately: got %s, want GLACIER_IR", got)
	}
	dateSpecs = []transitionSpec{{date: &future, class: s3store.StorageClassGlacierIR}}
	if got := resolveTransitionTarget(s3store.StorageClassStandard, dateSpecs, now, now); got != s3store.StorageClassStandard {
		t.Fatalf("a future date entry is not due: got %s, want STANDARD", got)
	}

	alreadyThere := []transitionSpec{{days: ptrInt32(1), class: s3store.StorageClassGlacier}}
	if got := resolveTransitionTarget(s3store.StorageClassGlacier, alreadyThere, old, now); got != s3store.StorageClassGlacier {
		t.Fatalf("an object already in the target class stays put: got %s", got)
	}
}

// The default transition minimum holds only while the rule carries no
// explicit size bound; a bound, direct or inside And, takes precedence.
func TestTransitionSizeEligible(t *testing.T) {
	small := &s3store.Object{Size: 100}
	large := &s3store.Object{Size: s3store.MinTransitionObjectSize}

	if transitionSizeEligible(small, nil) {
		t.Fatal("small object without a size bound must not transition")
	}
	if !transitionSizeEligible(large, nil) {
		t.Fatal("an object at the default minimum must transition")
	}
	greater := &s3store.LifecycleRuleFilter{ObjectSizeGreaterThan: ptrInt64(50)}
	if !transitionSizeEligible(small, greater) {
		t.Fatal("an explicit size bound overrides the default minimum")
	}
	less := &s3store.LifecycleRuleFilter{ObjectSizeLessThan: ptrInt64(1024)}
	if !transitionSizeEligible(small, less) {
		t.Fatal("an explicit maximum size overrides the default minimum")
	}
	and := &s3store.LifecycleRuleFilter{And: &s3store.LifecycleRuleAndOperator{ObjectSizeLessThan: ptrInt64(1024)}}
	if !transitionSizeEligible(small, and) {
		t.Fatal("an And-carried size bound overrides the default minimum")
	}
}

func ptrInt32(v int32) *int32 { return &v }

func ptrInt64(v int64) *int64 { return &v }

// The projected expiry follows the AWS timing calculation: days added to
// LastModified and rounded up to the next midnight UTC, a Date rule at the
// configured date, the earliest applicable rule wins, and markers,
// noncurrent versions, filter mismatches, and unrelicated objects carry
// none.
func TestObjectExpiration(t *testing.T) {
	stamp := time.Date(2026, 9, 9, 12, 0, 0, 500000000, time.UTC)
	obj := &s3store.Object{
		Key:          "docs/a.txt",
		BucketName:   "b",
		LastModified: stamp,
		IsLatest:     true,
		Tags:         nil,
	}

	days := func(n int32) *s3store.LifecycleConfiguration {
		return &s3store.LifecycleConfiguration{Rules: []s3store.LifecycleRule{{
			ID: "r1", Status: "Enabled",
			Filter:     &s3store.LifecycleRuleFilter{Prefix: "docs/"},
			Expiration: &s3store.LifecycleExpiration{Days: ptrInt32(n)},
		}}}
	}

	// 12:00:00.5 + 1 day rounds up to the following midnight.
	expiry, ruleID, ok := objectExpiration(days(1), obj)
	if !ok || ruleID != "r1" || !expiry.Equal(time.Date(2026, 9, 11, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("days=1 projection = %v %q %v, want 2026-09-11T00:00:00Z r1", expiry, ruleID, ok)
	}

	// An exact-midnight LastModified lands on that day's midnight.
	midnightObj := *obj
	midnightObj.LastModified = time.Date(2026, 9, 9, 0, 0, 0, 0, time.UTC)
	expiry, _, ok = objectExpiration(days(1), &midnightObj)
	if !ok || !expiry.Equal(time.Date(2026, 9, 10, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("midnight projection = %v %v, want 2026-09-10T00:00:00Z", expiry, ok)
	}

	// A Date rule expires at the configured date; with a Days rule also
	// present, the earlier instant wins.
	date := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	both := &s3store.LifecycleConfiguration{Rules: []s3store.LifecycleRule{
		{ID: "days-rule", Status: "Enabled", Expiration: &s3store.LifecycleExpiration{Days: ptrInt32(30)}},
		{ID: "date-rule", Status: "Enabled", Expiration: &s3store.LifecycleExpiration{Date: &date}},
	}}
	expiry, ruleID, ok = objectExpiration(both, obj)
	if !ok || ruleID != "date-rule" || !expiry.Equal(date) {
		t.Fatalf("earliest-rule projection = %v %q %v, want the date rule", expiry, ruleID, ok)
	}

	// Marker-only rules and Disabled rules carry no data-object expiry.
	markerOnly := &s3store.LifecycleConfiguration{Rules: []s3store.LifecycleRule{{
		ID: "markers", Status: "Enabled",
		Expiration: &s3store.LifecycleExpiration{ExpiredObjectDeleteMarker: ptrBool(true)},
	}}}
	if _, _, ok := objectExpiration(markerOnly, obj); ok {
		t.Fatal("an ExpiredObjectDeleteMarker-only rule must not project a data-object expiry")
	}

	// Filter mismatch, noncurrent versions, and unrelicated objects project
	// nothing.
	otherPrefix := days(1)
	otherPrefix.Rules[0].Filter.Prefix = "other/"
	if _, _, ok := objectExpiration(otherPrefix, obj); ok {
		t.Fatal("an object outside the rule filter must not project an expiry")
	}
	nonCurrent := *obj
	nonCurrent.IsLatest = false
	if _, _, ok := objectExpiration(days(1), &nonCurrent); ok {
		t.Fatal("a noncurrent version must not project an expiry")
	}
	failed := *obj
	failed.ReplicationStatus = s3store.ReplicationStatusFailed
	if _, _, ok := objectExpiration(days(1), &failed); ok {
		t.Fatal("an object whose replication has not succeeded must not project an expiry")
	}

	// The header value renders the documented pair: HTTP-date expiry with
	// the URL-encoded rule id.
	header := objectExpirationHeaderValue(days(1), obj)
	want := `expiry-date="Fri, 11 Sep 2026 00:00:00 GMT", rule-id="r1"`
	if header != want {
		t.Fatalf("header = %q, want %q", header, want)
	}
	spaced := days(1)
	spaced.Rules[0].ID = "my rule"
	header = objectExpirationHeaderValue(spaced, obj)
	if !strings.Contains(header, `rule-id="my+rule"`) {
		t.Fatalf("rule id not URL-encoded: %q", header)
	}
}

func ptrBool(v bool) *bool { return &v }
