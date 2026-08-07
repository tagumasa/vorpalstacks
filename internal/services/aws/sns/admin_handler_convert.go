package sns

import (
	"net/http"

	svccommon "vorpalstacks/internal/common"
	pb "vorpalstacks/internal/pb/aws/sns"
	snsstore "vorpalstacks/internal/store/aws/sns"
)

// ---------------------------------------------------------------------------
// admin_handler_convert.go — the sole file in the SNS service package that
// imports store packages and performs proto↔DTO conversion. This enforces
// AGENTS.md rule #29: admin handlers must not import store packages directly.
// ---------------------------------------------------------------------------

// getSNSStore returns the SNS store for the region extracted from the
// request header.
func (h *AdminHandler) getSNSStore(headers http.Header) (snsstore.SNSStoreInterface, error) {
	region := svccommon.GetRegionFromHeader(headers)
	return h.service.getSNSStoreByRegion(region)
}

// toPbCreateTopicResponse converts a service-layer TopicResult to the proto
// response type.
func toPbCreateTopicResponse(r *TopicResult) *pb.CreateTopicResponse {
	return &pb.CreateTopicResponse{
		Topicarn: r.Arn,
	}
}

// toPbListTopicsResponse converts a service-layer ListTopicsResult to the
// proto response type.
func toPbListTopicsResponse(r *ListTopicsResult) *pb.ListTopicsResponse {
	topics := make([]*pb.Topic, len(r.Topics))
	for i, t := range r.Topics {
		topics[i] = &pb.Topic{Topicarn: t.TopicArn}
	}
	return &pb.ListTopicsResponse{
		Topics:    topics,
		Nexttoken: r.NextToken,
	}
}
