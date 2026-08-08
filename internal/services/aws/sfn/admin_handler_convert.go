package sfn

import (
	"net/http"

	svccommon "vorpalstacks/internal/common"
	sfnstore "vorpalstacks/internal/store/aws/sfn"

	pb "vorpalstacks/internal/pb/aws/sfn"
)

// This file is the sole admin handler file permitted to import store packages
// (sole exception to the store-import prohibition). It provides store retrieval and proto conversion helpers
// that admin_handler.go uses without directly importing any store package.

// getStore resolves the per-region StepFunctionStore from request headers.
func (h *AdminHandler) getStore(headers http.Header) (*sfnstore.StepFunctionStore, error) {
	region := svccommon.GetRegionFromHeader(headers)
	return h.service.getStoreForRegion(region)
}

// stateMachineListItemToProto converts a store StateMachine to the proto
// StateMachineListItem used by the admin console.
func stateMachineListItemToProto(sm *sfnstore.StateMachine) *pb.StateMachineListItem {
	return &pb.StateMachineListItem{
		Statemachinearn: sm.StateMachineArn,
		Name:            sm.Name,
		Type:            toPbStateMachineType(sm.Type),
	}
}

// toPbStateMachineType converts a Smithy string type to the proto enum.
func toPbStateMachineType(typ string) pb.StateMachineType {
	switch typ {
	case "EXPRESS":
		return pb.StateMachineType_STATE_MACHINE_TYPE_EXPRESS
	default:
		return pb.StateMachineType_STATE_MACHINE_TYPE_STANDARD
	}
}
