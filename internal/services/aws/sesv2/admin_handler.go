package sesv2

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/email"
	emailconnect "vorpalstacks/internal/pb/aws/email/emailconnect"
	"vorpalstacks/internal/store/aws/common"
	sesv2store "vorpalstacks/internal/store/aws/sesv2"
)

// AdminHandler handles administrative operations for Amazon SES v2.
// It provides gRPC-Web handlers for managing email identities, configuration sets,
// contact lists, and other SES resources.
type AdminHandler struct {
	emailconnect.UnimplementedSESv2ServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
}

var _ emailconnect.SESv2ServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new AdminHandler with the given storage manager and account ID.
// The storage manager provides access to region-specific storage, while the account ID
// identifies the AWS account for which operations are performed.
func NewAdminHandler(storageManager *storage.RegionStorageManager, accountId string) *AdminHandler {
	return &AdminHandler{
		storageManager: storageManager,
		accountId:      accountId,
	}
}

// getSESv2Store retrieves the SES v2 store for the requested region.
// It extracts the region from the request header or defaults to us-east-1.
// Returns an error if the region storage cannot be obtained.
func (h *AdminHandler) getSESv2Store(req *connect.Request[pb.ListEmailIdentitiesRequest]) (*sesv2store.SESv2Store, error) {
	region := req.Header().Get("X-Aws-Region")
	if region == "" {
		region = "us-east-1"
	}
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return sesv2store.NewSESv2Store(regionStorage, h.accountId, region), nil
}

// ListEmailIdentities returns a list of all verified email identities for the account.
// It supports pagination via the NextToken marker and limits results to the specified page size.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListEmailIdentities(ctx context.Context, req *connect.Request[pb.ListEmailIdentitiesRequest]) (*connect.Response[pb.ListEmailIdentitiesResponse], error) {
	store, err := h.getSESv2Store(req)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	limit := int(req.Msg.Pagesize)
	if limit <= 0 {
		limit = 100
	}

	opts := common.ListOptions{
		MaxItems: limit,
		Marker:   req.Msg.Nexttoken,
	}

	result, err := store.ListEmailIdentities(opts)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var identities []*pb.IdentityInfo
	for _, identity := range result.Items {
		info := &pb.IdentityInfo{
			Identityname:       identity.Identity,
			Identitytype:       pb.IdentityType_IDENTITY_TYPE_EMAIL_ADDRESS,
			Sendingenabled:     identity.VerifiedForSending,
			Verificationstatus: pb.VerificationStatus_VERIFICATION_STATUS_SUCCESS,
		}
		if identity.IdentityType == "DOMAIN" {
			info.Identitytype = pb.IdentityType_IDENTITY_TYPE_DOMAIN
		}
		identities = append(identities, info)
	}

	return connect.NewResponse(&pb.ListEmailIdentitiesResponse{
		Emailidentities: identities,
		Nexttoken:       result.NextMarker,
	}), nil
}

// BatchGetMetricData retrieves metric data for multiple metrics in a single request.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) BatchGetMetricData(ctx context.Context, req *connect.Request[pb.BatchGetMetricDataRequest]) (*connect.Response[pb.BatchGetMetricDataResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CancelExportJob cancels an export job that is currently in progress.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CancelExportJob(ctx context.Context, req *connect.Request[pb.CancelExportJobRequest]) (*connect.Response[pb.CancelExportJobResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateConfigurationSet creates a new configuration set for governing email sending.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateConfigurationSet(ctx context.Context, req *connect.Request[pb.CreateConfigurationSetRequest]) (*connect.Response[pb.CreateConfigurationSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateConfigurationSetEventDestination configures an event destination for a configuration set.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateConfigurationSetEventDestination(ctx context.Context, req *connect.Request[pb.CreateConfigurationSetEventDestinationRequest]) (*connect.Response[pb.CreateConfigurationSetEventDestinationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateContact creates a contact within a contact list.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateContact(ctx context.Context, req *connect.Request[pb.CreateContactRequest]) (*connect.Response[pb.CreateContactResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateContactList creates a new contact list for email recipients.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateContactList(ctx context.Context, req *connect.Request[pb.CreateContactListRequest]) (*connect.Response[pb.CreateContactListResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateCustomVerificationEmailTemplate creates a custom verification email template.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateCustomVerificationEmailTemplate(ctx context.Context, req *connect.Request[pb.CreateCustomVerificationEmailTemplateRequest]) (*connect.Response[pb.CreateCustomVerificationEmailTemplateResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateDedicatedIpPool creates a new pool of dedicated IP addresses.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateDedicatedIpPool(ctx context.Context, req *connect.Request[pb.CreateDedicatedIpPoolRequest]) (*connect.Response[pb.CreateDedicatedIpPoolResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateDeliverabilityTestReport creates a new predictive inbox placement test report.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateDeliverabilityTestReport(ctx context.Context, req *connect.Request[pb.CreateDeliverabilityTestReportRequest]) (*connect.Response[pb.CreateDeliverabilityTestReportResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateEmailIdentity creates a new verified email identity or domain.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateEmailIdentity(ctx context.Context, req *connect.Request[pb.CreateEmailIdentityRequest]) (*connect.Response[pb.CreateEmailIdentityResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateEmailIdentityPolicy creates a sending authorisation policy for an identity.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateEmailIdentityPolicy(ctx context.Context, req *connect.Request[pb.CreateEmailIdentityPolicyRequest]) (*connect.Response[pb.CreateEmailIdentityPolicyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateEmailTemplate creates a new email template for use with SendEmail.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateEmailTemplate(ctx context.Context, req *connect.Request[pb.CreateEmailTemplateRequest]) (*connect.Response[pb.CreateEmailTemplateResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateExportJob initiates an export job for email data.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateExportJob(ctx context.Context, req *connect.Request[pb.CreateExportJobRequest]) (*connect.Response[pb.CreateExportJobResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateImportJob initiates an import job for email data.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateImportJob(ctx context.Context, req *connect.Request[pb.CreateImportJobRequest]) (*connect.Response[pb.CreateImportJobResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}
