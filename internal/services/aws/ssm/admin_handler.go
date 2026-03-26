package ssm

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/ssm"
	ssmconnect "vorpalstacks/internal/pb/aws/ssm/ssmconnect"
	ssmstore "vorpalstacks/internal/store/aws/ssm"
)

// AdminHandler provides Systems Manager (SSM) service administration functionality.
// It implements the SSMServiceHandler interface for gRPC-Web communication.
type AdminHandler struct {
	ssmconnect.UnimplementedSSMServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
}

var _ ssmconnect.SSMServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Systems Manager AdminHandler.
// It initialises the handler with the provided storage manager and account ID.
func NewAdminHandler(storageManager *storage.RegionStorageManager, accountId string) *AdminHandler {
	return &AdminHandler{
		storageManager: storageManager,
		accountId:      accountId,
	}
}

// getStore retrieves the SSM store for the request.
// It extracts the region from the request header and creates a new Store instance.
func (h *AdminHandler) getStore(req *connect.Request[pb.DescribeParametersRequest]) (*ssmstore.Store, error) {
	region := req.Header().Get("X-Aws-Region")
	if region == "" {
		region = "us-east-1"
	}
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return ssmstore.NewStore(regionStorage, h.accountId, region), nil
}

// DescribeParameters lists parameters in Systems Manager Parameter Store.
// It supports pagination via the NextToken parameter and filtering by parameter attributes.
func (h *AdminHandler) DescribeParameters(ctx context.Context, req *connect.Request[pb.DescribeParametersRequest]) (*connect.Response[pb.DescribeParametersResult], error) {
	store, err := h.getStore(req)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	maxResults := req.Msg.Maxresults
	if maxResults <= 0 {
		maxResults = 50
	}

	filters := make(map[string]string)
	for _, f := range req.Msg.Filters {
		if len(f.Values) > 0 {
			filters[f.Key.String()] = f.Values[0]
		}
	}

	params, nextToken, err := store.DescribeParameters(filters, maxResults, req.Msg.Nexttoken)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var metadataList []*pb.ParameterMetadata
	for _, p := range params {
		meta := &pb.ParameterMetadata{
			Name:             p.Name,
			Type:             pb.ParameterType_PARAMETER_TYPE_STRING,
			Tier:             pb.ParameterTier_PARAMETER_TIER_STANDARD,
			Version:          p.Version,
			Lastmodifieddate: p.LastModifiedDate.Format("2006-01-02T15:04:05.000Z"),
			Datatype:         p.DataType,
			Arn:              p.ARN,
		}
		if p.Description != "" {
			meta.Description = p.Description
		}
		if p.KeyID != "" {
			meta.Keyid = p.KeyID
		}
		if p.AllowedPattern != "" {
			meta.Allowedpattern = p.AllowedPattern
		}
		switch p.Type {
		case ssmstore.ParameterTypeString:
			meta.Type = pb.ParameterType_PARAMETER_TYPE_STRING
		case ssmstore.ParameterTypeStringList:
			meta.Type = pb.ParameterType_PARAMETER_TYPE_STRING_LIST
		case ssmstore.ParameterTypeSecureString:
			meta.Type = pb.ParameterType_PARAMETER_TYPE_SECURE_STRING
		}
		switch p.Tier {
		case ssmstore.ParameterTierStandard:
			meta.Tier = pb.ParameterTier_PARAMETER_TIER_STANDARD
		case ssmstore.ParameterTierAdvanced:
			meta.Tier = pb.ParameterTier_PARAMETER_TIER_ADVANCED
		case ssmstore.ParameterTierIntelligentTiering:
			meta.Tier = pb.ParameterTier_PARAMETER_TIER_INTELLIGENT_TIERING
		}
		metadataList = append(metadataList, meta)
	}

	return connect.NewResponse(&pb.DescribeParametersResult{
		Parameters: metadataList,
		Nexttoken:  nextToken,
	}), nil
}

// AddTagsToResource adds tags to an SSM resource.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) AddTagsToResource(ctx context.Context, req *connect.Request[pb.AddTagsToResourceRequest]) (*connect.Response[pb.AddTagsToResourceResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// AssociateOpsItemRelatedItem associates an OpsItem with a related item.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) AssociateOpsItemRelatedItem(ctx context.Context, req *connect.Request[pb.AssociateOpsItemRelatedItemRequest]) (*connect.Response[pb.AssociateOpsItemRelatedItemResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CancelCommand cancels a command that has been sent to an instance.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CancelCommand(ctx context.Context, req *connect.Request[pb.CancelCommandRequest]) (*connect.Response[pb.CancelCommandResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CancelMaintenanceWindowExecution cancels a maintenance window execution.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CancelMaintenanceWindowExecution(ctx context.Context, req *connect.Request[pb.CancelMaintenanceWindowExecutionRequest]) (*connect.Response[pb.CancelMaintenanceWindowExecutionResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateActivation creates an activation for a managed instance.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateActivation(ctx context.Context, req *connect.Request[pb.CreateActivationRequest]) (*connect.Response[pb.CreateActivationResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateAssociation creates an association between a document and an instance.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateAssociation(ctx context.Context, req *connect.Request[pb.CreateAssociationRequest]) (*connect.Response[pb.CreateAssociationResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateAssociationBatch creates associations for multiple instances in a single request.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateAssociationBatch(ctx context.Context, req *connect.Request[pb.CreateAssociationBatchRequest]) (*connect.Response[pb.CreateAssociationBatchResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateDocument creates a new SSM document.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateDocument(ctx context.Context, req *connect.Request[pb.CreateDocumentRequest]) (*connect.Response[pb.CreateDocumentResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateMaintenanceWindow creates a maintenance window for managing instances.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateMaintenanceWindow(ctx context.Context, req *connect.Request[pb.CreateMaintenanceWindowRequest]) (*connect.Response[pb.CreateMaintenanceWindowResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateOpsItem creates an OpsItem in AWS Systems Manager OpsCenter.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateOpsItem(ctx context.Context, req *connect.Request[pb.CreateOpsItemRequest]) (*connect.Response[pb.CreateOpsItemResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateOpsMetadata creates operational metadata for an OpsItem.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateOpsMetadata(ctx context.Context, req *connect.Request[pb.CreateOpsMetadataRequest]) (*connect.Response[pb.CreateOpsMetadataResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreatePatchBaseline creates a patch baseline for managing patches on instances.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreatePatchBaseline(ctx context.Context, req *connect.Request[pb.CreatePatchBaselineRequest]) (*connect.Response[pb.CreatePatchBaselineResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateResourceDataSync creates a resource data sync for Systems Manager Inventory.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateResourceDataSync(ctx context.Context, req *connect.Request[pb.CreateResourceDataSyncRequest]) (*connect.Response[pb.CreateResourceDataSyncResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteActivation deletes an activation for a managed instance.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteActivation(ctx context.Context, req *connect.Request[pb.DeleteActivationRequest]) (*connect.Response[pb.DeleteActivationResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteAssociation deletes an association between a document and an instance.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteAssociation(ctx context.Context, req *connect.Request[pb.DeleteAssociationRequest]) (*connect.Response[pb.DeleteAssociationResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteDocument deletes an SSM document.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteDocument(ctx context.Context, req *connect.Request[pb.DeleteDocumentRequest]) (*connect.Response[pb.DeleteDocumentResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteInventory deletes custom inventory data from managed instances.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteInventory(ctx context.Context, req *connect.Request[pb.DeleteInventoryRequest]) (*connect.Response[pb.DeleteInventoryResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}
