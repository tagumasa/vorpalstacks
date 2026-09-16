package scheduler

import (
	"strings"
	"testing"
)

// TestScheduleGroupArnToNameCharset pins that the group-name portion of
// a schedule-group ARN is validated against the single namePattern
// definition (the ScheduleGroupName charset): malformed characters are a
// validation failure, not a missing resource.
func TestScheduleGroupArnToNameCharset(t *testing.T) {
	valid := "arn:aws:scheduler:us-east-1:000000000000:schedule-group/team-a_prod.1"
	if name, err := scheduleGroupArnToName(valid); err != nil || name != "team-a_prod.1" {
		t.Errorf("scheduleGroupArnToName(%q) = %q, %v; want team-a_prod.1, nil", valid, name, err)
	}

	for _, arn := range []string{
		"arn:aws:scheduler:us-east-1:000000000000:schedule-group/bad@name",
		"arn:aws:scheduler:us-east-1:000000000000:schedule-group/spa ce",
		"arn:aws:scheduler:us-east-1:000000000000:schedule-group/sl/ash",
		"arn:aws:scheduler:us-east-1:000000000000:schedule-group/",
		"arn:aws:events:us-east-1:000000000000:schedule-group/valid-name",
	} {
		if _, err := scheduleGroupArnToName(arn); err == nil {
			t.Errorf("scheduleGroupArnToName(%q) accepted a malformed ARN", arn)
		}
	}
}

// TestScheduleGroupArnToNameLength pins both length bounds on the tag
// target ARN: TagResourceArn @length(1, 1011) rejects an over-length
// ARN, and the group-name portion carries the name shapes' 1-64 bound —
// a conforming name can never push its ARN near the 1011 bound, so the
// longest accepted name is 64 characters.
func TestScheduleGroupArnToNameLength(t *testing.T) {
	prefix := "arn:aws:scheduler:us-east-1:000000000000:schedule-group/"
	over := prefix + strings.Repeat("a", maxTagResourceArnLength-len(prefix)+1)
	if _, err := scheduleGroupArnToName(over); err == nil {
		t.Errorf("scheduleGroupArnToName accepted an ARN over the %d maximum", maxTagResourceArnLength)
	}
	atNameBound := prefix + strings.Repeat("a", 64)
	if _, err := scheduleGroupArnToName(atNameBound); err != nil {
		t.Errorf("scheduleGroupArnToName rejected a 64-character group name: %v", err)
	}
	overName := prefix + strings.Repeat("a", 65)
	if _, err := scheduleGroupArnToName(overName); err == nil {
		t.Errorf("scheduleGroupArnToName accepted a 65-character group name")
	}
}
