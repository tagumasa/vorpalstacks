package eventbridge

import (
	"net/http"
	"vorpalstacks/internal/common/defaults"

	pb "vorpalstacks/internal/pb/aws/cloudwatchevents"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// getStore resolves the per-region EventsStore for the admin handler.
func (h *AdminHandler) getStore(header http.Header) (*eventsstore.EventsStore, error) {
	region := defaults.GetRegionFromHeader(header)
	return h.service.GetStoreForRegion(region)
}

// toPbEventBus converts a store EventBus to the proto representation.
func toPbEventBus(eb *eventsstore.EventBus) *pb.EventBus {
	return &pb.EventBus{
		Name:   eb.Name,
		Arn:    eb.ARN,
		Policy: eb.Policy,
	}
}

// toPbRule converts a store Rule to the proto representation.
func toPbRule(r *eventsstore.Rule) *pb.Rule {
	return &pb.Rule{
		Name:               r.Name,
		Arn:                r.ARN,
		Eventbusname:       r.EventBusName,
		Description:        r.Description,
		Eventpattern:       r.EventPattern,
		Scheduleexpression: r.ScheduleExpression,
		State:              toPbRuleState(r.State),
		Managedby:          r.ManagedBy,
		Rolearn:            r.RoleARN,
	}
}

// toPbRuleState maps a store RuleState to the proto enum.
// ENABLED_WITH_ALL_CLOUDTRAIL_MANAGEMENT_EVENTS maps to RULE_STATE_ENABLED
// because the proto schema does not define a separate enum value for it;
// both are functionally "active" states.
func toPbRuleState(state eventsstore.RuleState) pb.RuleState {
	switch state {
	case eventsstore.RuleStateEnabled, eventsstore.RuleStateEnabledWithAllCloudtrailManagementEvents:
		return pb.RuleState_RULE_STATE_ENABLED
	default:
		return pb.RuleState_RULE_STATE_DISABLED
	}
}
