package sqs

import (
	"net/http"

	svccommon "vorpalstacks/internal/common"
	pb "vorpalstacks/internal/pb/aws/sqs"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
)

// ---------------------------------------------------------------------------
// admin_handler_convert.go — the sole file in the SQS service package that
// imports store packages and performs proto↔DTO conversion. This enforces
// AGENTS.md rule #29: admin handlers must not import store packages directly.
// ---------------------------------------------------------------------------

// getQueueStore returns the SQS store for the region extracted from the
// request header.
func (h *AdminHandler) getQueueStore(headers http.Header) (sqsstore.SQSStoreInterface, error) {
	region := svccommon.GetRegionFromHeader(headers)
	return h.service.GetStoreForRegion(region)
}

// toPbCreateQueueResult converts a service-layer CreateQueueResult to the
// proto response type.
func toPbCreateQueueResult(r *CreateQueueResult) *pb.CreateQueueResult {
	return &pb.CreateQueueResult{
		Queueurl: r.QueueURL,
	}
}

// toPbListQueuesResult converts a service-layer ListQueuesResult to the
// proto response type.
func toPbListQueuesResult(r *ListQueuesResult) *pb.ListQueuesResult {
	return &pb.ListQueuesResult{
		Queueurls: r.QueueURLs,
		Nexttoken: r.NextToken,
	}
}

// toPbGetQueueUrlResult converts a service-layer GetQueueUrlResult to the
// proto response type.
func toPbGetQueueUrlResult(r *GetQueueUrlResult) *pb.GetQueueUrlResult {
	return &pb.GetQueueUrlResult{
		Queueurl: r.QueueURL,
	}
}
