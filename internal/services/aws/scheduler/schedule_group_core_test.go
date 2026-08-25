package scheduler

import "testing"

// TestScheduleGroupArnToNameCharset pins that the group-name portion of
// a schedule-group ARN is validated against the TagResourceArn pattern
// ([0-9a-zA-Z-_.]+): malformed characters are a validation failure, not
// a missing resource.
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
