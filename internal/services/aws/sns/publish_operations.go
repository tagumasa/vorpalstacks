package sns

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// Publish publishes a message to an SNS topic.
// https://docs.aws.amazon.com/sns/latest/api/API_Publish.html
//
// All validation, deduplication, and delivery fan-out live inside the Core.
func (s *SNSService) Publish(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.publishCore(store, reqCtx.GetRegion(), PublishInput{
		TopicArn:               request.GetParamLowerFirst(req.Parameters, "TopicArn"),
		TargetArn:              request.GetParamLowerFirst(req.Parameters, "TargetArn"),
		PhoneNumber:            request.GetParamLowerFirst(req.Parameters, "PhoneNumber"),
		Message:                request.GetParamLowerFirst(req.Parameters, "Message"),
		Subject:                request.GetParamLowerFirst(req.Parameters, "Subject"),
		MessageStructure:       request.GetParamLowerFirst(req.Parameters, "MessageStructure"),
		MessageGroupId:         request.GetParamLowerFirst(req.Parameters, "MessageGroupId"),
		MessageDeduplicationId: request.GetParamLowerFirst(req.Parameters, "MessageDeduplicationId"),
		Parameters:             req.Parameters,
	})
}

// PublishBatch publishes multiple messages to an SNS topic in a single request.
// https://docs.aws.amazon.com/sns/latest/api/API_PublishBatch.html
//
// Two-pass design (inside the Core) — all entries are validated before any
// entry is delivered, and the batch-level size check rejects the entire batch
// before any delivery or dedup recording.
func (s *SNSService) PublishBatch(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Wire parsing only; entry validation and execution live in the Core.
	entries := request.GetListParam(req.Parameters, "PublishBatchRequestEntries")

	return s.publishBatchCore(store, reqCtx, PublishBatchInput{
		TopicArn: request.GetParamLowerFirst(req.Parameters, "TopicArn"),
		Entries:  entries,
	})
}
