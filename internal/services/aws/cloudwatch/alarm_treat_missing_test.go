package cloudwatch

import (
	"testing"

	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
)

func treatMissingTestAlarm(treatMissingData string) *cwstore.Alarm {
	a := cwstore.NewAlarm("missing-data", "ns", "metric")
	a.EvaluationPeriods = 3
	a.DatapointsToAlarm = 3
	a.TreatMissingData = treatMissingData
	a.State = "OK"
	return a
}

// An alarm without an explicit TreatMissingData follows the AWS default
// "missing" behaviour: absent periods are not breaches, and with no real
// data at all the alarm moves to INSUFFICIENT_DATA rather than ALARM.
func TestTreatMissingDataUnsetDefaultIsMissing(t *testing.T) {
	result := determineStateTransition(treatMissingTestAlarm(""), 0, 0)
	if result == nil {
		t.Fatal("no transition computed for fully missing data")
	}
	if result.newState != "INSUFFICIENT_DATA" {
		t.Fatalf("unset TreatMissingData with all periods missing -> %s, want INSUFFICIENT_DATA", result.newState)
	}
}

// "breaching" counts missing periods as breaches and alarms with no real
// data at all.
func TestTreatMissingDataBreachingCountsMissing(t *testing.T) {
	result := determineStateTransition(treatMissingTestAlarm("breaching"), 0, 0)
	if result == nil || result.newState != "ALARM" {
		t.Fatalf("breaching with all periods missing -> %+v, want ALARM", result)
	}
}

// Real breaching data still alarms regardless of the TreatMissingData
// setting; missing periods only fill the gaps.
func TestTreatMissingDataRealBreachesStillAlarm(t *testing.T) {
	for _, setting := range []string{"", "missing", "notBreaching", "breaching", "ignore"} {
		result := determineStateTransition(treatMissingTestAlarm(setting), 3, 3)
		if result == nil || result.newState != "ALARM" {
			t.Fatalf("TreatMissingData=%q with 3/3 real breaches -> %+v, want ALARM", setting, result)
		}
	}
}

// "notBreaching" fills missing periods as non-breaching, so partial real
// breaches below the threshold keep the alarm in its current (OK) state.
func TestTreatMissingDataNotBreachingFillsGaps(t *testing.T) {
	result := determineStateTransition(treatMissingTestAlarm("notBreaching"), 2, 2)
	if result != nil {
		t.Fatalf("notBreaching with 2/3 real breaches -> %+v, want no transition from OK", result)
	}
}

// A cycle that already exists between unrelated composites must not block
// creation of a new acyclic composite; only a rule that closes a path back
// to its own alarm is rejected.
func TestValidateCompositeAcyclicIgnoresUnrelatedCycle(t *testing.T) {
	store := newAlarmTestStore(t)

	composite := func(name, rule string) *cwstore.Alarm {
		a := cwstore.NewAlarm(name, "", "")
		a.AlarmType = cwstore.AlarmTypeCompositeAlarm
		a.AlarmRule = rule
		return a
	}
	// Legacy cyclic pair written directly to the store, predating creation
	// validation.
	if _, err := store.CreateAlarm(composite("x", `ALARM("y")`)); err != nil {
		t.Fatal(err)
	}
	if _, err := store.CreateAlarm(composite("y", `ALARM("x")`)); err != nil {
		t.Fatal(err)
	}

	freshRule, err := parseAlarmRule(`ALARM("unrelated-metric")`)
	if err != nil {
		t.Fatal(err)
	}
	if err := validateCompositeAcyclic(store, "fresh", freshRule); err != nil {
		t.Fatalf("unrelated existing cycle blocked a new acyclic alarm: %v", err)
	}

	// Referencing a member of the legacy cycle does not close a cycle back
	// to the new alarm, so creation proceeds; the cyclic pair itself stays
	// unevaluable on every tick.
	refRule, err := parseAlarmRule(`ALARM("x")`)
	if err != nil {
		t.Fatal(err)
	}
	if err := validateCompositeAcyclic(store, "fresh", refRule); err != nil {
		t.Fatalf("reference to a cyclic alarm rejected although it closes no cycle: %v", err)
	}
}
