package eventbridge

import (
	"google.golang.org/protobuf/proto"
	"net/http"
	"vorpalstacks/internal/common/defaults"
	"vorpalstacks/internal/common/pbutil"

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
		Name:   proto.String(eb.Name),
		Arn:    proto.String(eb.ARN),
		Policy: proto.String(eb.Policy),
	}
}

// toPbRule converts a store Rule to the proto representation.
func toPbRule(r *eventsstore.Rule) *pb.Rule {
	return &pb.Rule{
		Name:               proto.String(r.Name),
		Arn:                proto.String(r.ARN),
		Eventbusname:       proto.String(r.EventBusName),
		Description:        proto.String(r.Description),
		Eventpattern:       proto.String(r.EventPattern),
		Scheduleexpression: proto.String(r.ScheduleExpression),
		State:              pbutil.Enum(toPbRuleState(r.State)),
		Managedby:          proto.String(r.ManagedBy),
		Rolearn:            proto.String(r.RoleARN),
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
