package sns

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/sns"
	snsconnect "vorpalstacks/internal/pb/aws/sns/snsconnect"
	svccommon "vorpalstacks/internal/services/aws/common"
	storecommon "vorpalstacks/internal/store/aws/common"
	snsstore "vorpalstacks/internal/store/aws/sns"
)

// AdminHandler provides administrative operations for SNS topics and subscriptions.
// It implements the SNSServiceHandler interface for gRPC-Web communication.
type AdminHandler struct {
	snsconnect.UnimplementedSNSServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
}

var _ snsconnect.SNSServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new AdminHandler with the given storage manager and account ID.
// The storage manager provides access to regional storage for SNS data.
func NewAdminHandler(storageManager *storage.RegionStorageManager, accountId string) *AdminHandler {
	return &AdminHandler{
		storageManager: storageManager,
		accountId:      accountId,
	}
}

// getSNSStoreByRegion returns an SNS store for the specified region.
// The store provides access to SNS topics and subscriptions in that region.
func (h *AdminHandler) getSNSStoreByRegion(region string) (*snsstore.SNSStore, error) {
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return snsstore.NewSNSStore(regionStorage, h.accountId, region), nil
}

// ListTopics retrieves a list of SNS topics in the specified region.
// It supports pagination via the NextToken parameter.
// Returns a list of Topic objects containing the TopicArn.
func (h *AdminHandler) ListTopics(ctx context.Context, req *connect.Request[pb.ListTopicsInput]) (*connect.Response[pb.ListTopicsResponse], error) {
	region := svccommon.GetRegionFromHeader(req.Header())
	store, err := h.getSNSStoreByRegion(region)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	opts := storecommon.ListOptions{
		Marker:   req.Msg.Nexttoken,
		MaxItems: 100,
	}

	result, err := store.ListTopics(opts)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	topics := make([]*pb.Topic, len(result.Items))
	for i, t := range result.Items {
		topics[i] = &pb.Topic{
			Topicarn: t.Arn,
		}
	}

	return connect.NewResponse(&pb.ListTopicsResponse{
		Topics:    topics,
		Nexttoken: result.NextMarker,
	}), nil
}

// AddPermission adds permissions to a topic for cross-account access.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) AddPermission(ctx context.Context, req *connect.Request[pb.AddPermissionInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CheckIfPhoneNumberIsOptedOut checks whether a phone number has opted out of receiving SMS messages.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CheckIfPhoneNumberIsOptedOut(ctx context.Context, req *connect.Request[pb.CheckIfPhoneNumberIsOptedOutInput]) (*connect.Response[pb.CheckIfPhoneNumberIsOptedOutResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ConfirmSubscription confirms a subscription by returning the subscription ARN.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ConfirmSubscription(ctx context.Context, req *connect.Request[pb.ConfirmSubscriptionInput]) (*connect.Response[pb.ConfirmSubscriptionResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreatePlatformApplication creates a platform application object for push notification services.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreatePlatformApplication(ctx context.Context, req *connect.Request[pb.CreatePlatformApplicationInput]) (*connect.Response[pb.CreatePlatformApplicationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreatePlatformEndpoint creates an endpoint for a device on a push notification platform.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreatePlatformEndpoint(ctx context.Context, req *connect.Request[pb.CreatePlatformEndpointInput]) (*connect.Response[pb.CreateEndpointResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateSMSSandboxPhoneNumber creates a new SMS sandbox phone number for testing.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateSMSSandboxPhoneNumber(ctx context.Context, req *connect.Request[pb.CreateSMSSandboxPhoneNumberInput]) (*connect.Response[pb.CreateSMSSandboxPhoneNumberResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateTopic creates a new SNS topic for publishing messages.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateTopic(ctx context.Context, req *connect.Request[pb.CreateTopicInput]) (*connect.Response[pb.CreateTopicResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteEndpoint deletes an endpoint for a push notification platform.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteEndpoint(ctx context.Context, req *connect.Request[pb.DeleteEndpointInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeletePlatformApplication deletes a platform application object.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeletePlatformApplication(ctx context.Context, req *connect.Request[pb.DeletePlatformApplicationInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteSMSSandboxPhoneNumber deletes an SMS sandbox phone number.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteSMSSandboxPhoneNumber(ctx context.Context, req *connect.Request[pb.DeleteSMSSandboxPhoneNumberInput]) (*connect.Response[pb.DeleteSMSSandboxPhoneNumberResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteTopic deletes an SNS topic and all its subscriptions.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteTopic(ctx context.Context, req *connect.Request[pb.DeleteTopicInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetDataProtectionPolicy retrieves the data protection policy for an SNS topic.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetDataProtectionPolicy(ctx context.Context, req *connect.Request[pb.GetDataProtectionPolicyInput]) (*connect.Response[pb.GetDataProtectionPolicyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetEndpointAttributes retrieves the attributes of an endpoint for push notifications.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetEndpointAttributes(ctx context.Context, req *connect.Request[pb.GetEndpointAttributesInput]) (*connect.Response[pb.GetEndpointAttributesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetPlatformApplicationAttributes retrieves the attributes of a platform application.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetPlatformApplicationAttributes(ctx context.Context, req *connect.Request[pb.GetPlatformApplicationAttributesInput]) (*connect.Response[pb.GetPlatformApplicationAttributesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetSMSAttributes retrieves the default settings for sending SMS messages.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetSMSAttributes(ctx context.Context, req *connect.Request[pb.GetSMSAttributesInput]) (*connect.Response[pb.GetSMSAttributesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetSMSSandboxAccountStatus retrieves the status of the SMS sandbox for the account.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetSMSSandboxAccountStatus(ctx context.Context, req *connect.Request[pb.GetSMSSandboxAccountStatusInput]) (*connect.Response[pb.GetSMSSandboxAccountStatusResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetSubscriptionAttributes retrieves the attributes of a subscription.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetSubscriptionAttributes(ctx context.Context, req *connect.Request[pb.GetSubscriptionAttributesInput]) (*connect.Response[pb.GetSubscriptionAttributesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetTopicAttributes retrieves the attributes of an SNS topic.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetTopicAttributes(ctx context.Context, req *connect.Request[pb.GetTopicAttributesInput]) (*connect.Response[pb.GetTopicAttributesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListEndpointsByPlatformApplication lists the endpoints associated with a platform application.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListEndpointsByPlatformApplication(ctx context.Context, req *connect.Request[pb.ListEndpointsByPlatformApplicationInput]) (*connect.Response[pb.ListEndpointsByPlatformApplicationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListOriginationNumbers lists the phone numbers that are opted in to the SMS sandbox.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListOriginationNumbers(ctx context.Context, req *connect.Request[pb.ListOriginationNumbersRequest]) (*connect.Response[pb.ListOriginationNumbersResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListPhoneNumbersOptedOut lists phone numbers that have opted out of receiving SMS messages.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListPhoneNumbersOptedOut(ctx context.Context, req *connect.Request[pb.ListPhoneNumbersOptedOutInput]) (*connect.Response[pb.ListPhoneNumbersOptedOutResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListPlatformApplications lists the platform applications registered for push notifications.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListPlatformApplications(ctx context.Context, req *connect.Request[pb.ListPlatformApplicationsInput]) (*connect.Response[pb.ListPlatformApplicationsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListSMSSandboxPhoneNumbers lists the SMS sandbox phone numbers in the account.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListSMSSandboxPhoneNumbers(ctx context.Context, req *connect.Request[pb.ListSMSSandboxPhoneNumbersInput]) (*connect.Response[pb.ListSMSSandboxPhoneNumbersResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListSubscriptions returns a list of all subscriptions in the account.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListSubscriptions(ctx context.Context, req *connect.Request[pb.ListSubscriptionsInput]) (*connect.Response[pb.ListSubscriptionsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListSubscriptionsByTopic returns a list of subscriptions to a specific SNS topic.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListSubscriptionsByTopic(ctx context.Context, req *connect.Request[pb.ListSubscriptionsByTopicInput]) (*connect.Response[pb.ListSubscriptionsByTopicResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListTagsForResource returns the tags attached to an SNS topic or platform application.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListTagsForResource(ctx context.Context, req *connect.Request[pb.ListTagsForResourceRequest]) (*connect.Response[pb.ListTagsForResourceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// OptInPhoneNumber opts in a phone number to receive SMS messages.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) OptInPhoneNumber(ctx context.Context, req *connect.Request[pb.OptInPhoneNumberInput]) (*connect.Response[pb.OptInPhoneNumberResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// Publish sends a message to an SNS topic or directly to an endpoint.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) Publish(ctx context.Context, req *connect.Request[pb.PublishInput]) (*connect.Response[pb.PublishResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// PublishBatch sends multiple messages to an SNS topic in a single request.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PublishBatch(ctx context.Context, req *connect.Request[pb.PublishBatchInput]) (*connect.Response[pb.PublishBatchResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// PutDataProtectionPolicy sets a data protection policy on an SNS topic.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutDataProtectionPolicy(ctx context.Context, req *connect.Request[pb.PutDataProtectionPolicyInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// RemovePermission removes permissions from a topic for cross-account access.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) RemovePermission(ctx context.Context, req *connect.Request[pb.RemovePermissionInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// SetEndpointAttributes sets the attributes of an endpoint for push notifications.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) SetEndpointAttributes(ctx context.Context, req *connect.Request[pb.SetEndpointAttributesInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// SetPlatformApplicationAttributes sets the attributes of a platform application.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) SetPlatformApplicationAttributes(ctx context.Context, req *connect.Request[pb.SetPlatformApplicationAttributesInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// SetSMSAttributes sets the default settings for sending SMS messages.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) SetSMSAttributes(ctx context.Context, req *connect.Request[pb.SetSMSAttributesInput]) (*connect.Response[pb.SetSMSAttributesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// SetSubscriptionAttributes sets the attributes of a subscription.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) SetSubscriptionAttributes(ctx context.Context, req *connect.Request[pb.SetSubscriptionAttributesInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// SetTopicAttributes sets the attributes of an SNS topic.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) SetTopicAttributes(ctx context.Context, req *connect.Request[pb.SetTopicAttributesInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// Subscribe creates a subscription to an SNS topic.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) Subscribe(ctx context.Context, req *connect.Request[pb.SubscribeInput]) (*connect.Response[pb.SubscribeResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// TagResource adds or overwrites tags on an SNS topic or platform application.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) TagResource(ctx context.Context, req *connect.Request[pb.TagResourceRequest]) (*connect.Response[pb.TagResourceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// Unsubscribe deletes a subscription to an SNS topic.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) Unsubscribe(ctx context.Context, req *connect.Request[pb.UnsubscribeInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UntagResource removes tags from an SNS topic or platform application.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UntagResource(ctx context.Context, req *connect.Request[pb.UntagResourceRequest]) (*connect.Response[pb.UntagResourceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// VerifySMSSandboxPhoneNumber verifies a phone number in the SMS sandbox.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) VerifySMSSandboxPhoneNumber(ctx context.Context, req *connect.Request[pb.VerifySMSSandboxPhoneNumberInput]) (*connect.Response[pb.VerifySMSSandboxPhoneNumberResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}
