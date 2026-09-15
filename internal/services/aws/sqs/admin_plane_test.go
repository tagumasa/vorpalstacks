package sqs

import (
	"context"
	"strings"
	"testing"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"

	"vorpalstacks/internal/common/request"
	pb "vorpalstacks/internal/pb/aws/sqs"
)

// TestAdminListQueuesMaxResultsSet pins the admin-plane MaxResults parity:
// the proto's optional Maxresults field distinguishes an explicitly supplied
// 0 (a set value the Core rejects below 1, exactly like the HTTP plane) from
// a nil field (omitted, default page). Dropping the nil distinction made the
// admin surface silently looser than the same Core's HTTP contract.
func TestAdminListQueuesMaxResultsSet(t *testing.T) {
	svc, reqCtx, _ := newQueryWireTestService(t)

	if _, err := svc.CreateQueue(context.Background(), reqCtx, &request.ParsedRequest{
		Operation:  "CreateQueue",
		Parameters: map[string]interface{}{"QueueName": "admin-lq"},
	}); err != nil {
		t.Fatalf("CreateQueue: %v", err)
	}

	handler := NewAdminHandler(svc)

	// Explicit 0 — rejected like the HTTP plane.
	zeroReq := connect.NewRequest(&pb.ListQueuesRequest{Maxresults: proto.Int32(0)})
	zeroReq.Header().Set("X-Aws-Region", "us-east-1")
	_, err := handler.ListQueues(context.Background(), zeroReq)
	if err == nil || connect.CodeOf(err) != connect.CodeInvalidArgument || !strings.Contains(err.Error(), "InvalidParameterValue") {
		t.Fatalf("explicit Maxresults 0: got %v, want InvalidArgument carrying InvalidParameterValue", err)
	}

	// Nil field — omitted, default page, the queue is listed.
	nilReq := connect.NewRequest(&pb.ListQueuesRequest{})
	nilReq.Header().Set("X-Aws-Region", "us-east-1")
	resp, err := handler.ListQueues(context.Background(), nilReq)
	if err != nil {
		t.Fatalf("ListQueues without Maxresults: %v", err)
	}
	if len(resp.Msg.Queueurls) != 1 {
		t.Fatalf("absent Maxresults: Queueurls = %v, want the created queue", resp.Msg.Queueurls)
	}
}

// TestGetQueueUrlOwnerAccount pins the QueueOwnerAWSAccountId contract on
// both planes: the platform runs a single account, so a queue under any other
// owner does not exist in this deployment — the Core reads the queue's own
// account from its ARN and answers QueueDoesNotExist for a mismatch instead
// of silently returning the local queue's URL. An empty member or the local
// account resolves normally.
func TestGetQueueUrlOwnerAccount(t *testing.T) {
	svc, reqCtx, _ := newQueryWireTestService(t)

	resp, err := svc.CreateQueue(context.Background(), reqCtx, &request.ParsedRequest{
		Operation:  "CreateQueue",
		Parameters: map[string]interface{}{"QueueName": "owner-check"},
	})
	if err != nil {
		t.Fatalf("CreateQueue: %v", err)
	}
	respMap, _ := resp.(map[string]interface{})
	queueURL, _ := respMap["QueueUrl"].(string)

	// Mismatched owner — the queue does not exist under that account.
	_, err = svc.GetQueueUrl(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "GetQueueUrl",
		Parameters: map[string]interface{}{
			"QueueName":              "owner-check",
			"QueueOwnerAWSAccountId": "999999999999",
		},
	})
	if err == nil || !strings.Contains(err.Error(), "QueueDoesNotExist") {
		t.Fatalf("foreign owner account: got %v, want QueueDoesNotExist", err)
	}

	// The local account and an empty member both resolve the queue.
	for _, owner := range []string{"123456789012", ""} {
		params := map[string]interface{}{"QueueName": "owner-check"}
		if owner != "" {
			params["QueueOwnerAWSAccountId"] = owner
		}
		urlResp, err := svc.GetQueueUrl(context.Background(), reqCtx, &request.ParsedRequest{
			Operation:  "GetQueueUrl",
			Parameters: params,
		})
		if err != nil {
			t.Fatalf("GetQueueUrl owner %q: %v", owner, err)
		}
		urlMap, _ := urlResp.(map[string]interface{})
		if got, _ := urlMap["QueueUrl"].(string); got != queueURL {
			t.Fatalf("GetQueueUrl owner %q: QueueUrl = %q, want %q", owner, got, queueURL)
		}
	}

	// Admin plane carries the same contract.
	handler := NewAdminHandler(svc)
	foreignReq := connect.NewRequest(&pb.GetQueueUrlRequest{
		Queuename:              "owner-check",
		Queueownerawsaccountid: proto.String("999999999999"),
	})
	foreignReq.Header().Set("X-Aws-Region", "us-east-1")
	_, err = handler.GetQueueUrl(context.Background(), foreignReq)
	if err == nil || connect.CodeOf(err) != connect.CodeInvalidArgument || !strings.Contains(err.Error(), "QueueDoesNotExist") {
		t.Fatalf("admin foreign owner account: got %v, want InvalidArgument carrying QueueDoesNotExist", err)
	}

	ownReq := connect.NewRequest(&pb.GetQueueUrlRequest{
		Queuename:              "owner-check",
		Queueownerawsaccountid: proto.String("123456789012"),
	})
	ownReq.Header().Set("X-Aws-Region", "us-east-1")
	ownResp, err := handler.GetQueueUrl(context.Background(), ownReq)
	if err != nil {
		t.Fatalf("admin GetQueueUrl local owner: %v", err)
	}
	if ownResp.Msg.GetQueueurl() != queueURL {
		t.Fatalf("admin GetQueueUrl local owner: Queueurl = %q, want %q", ownResp.Msg.GetQueueurl(), queueURL)
	}
}
