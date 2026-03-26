package athena

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/athena"
	athenaconnect "vorpalstacks/internal/pb/aws/athena/athenaconnect"
	athenastore "vorpalstacks/internal/store/aws/athena"
)

// AdminHandler provides Athena service administration functionality.
// It implements the AthenaServiceHandler interface for gRPC-Web communication.
type AdminHandler struct {
	athenaconnect.UnimplementedAthenaServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
}

var _ athenaconnect.AthenaServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Athena AdminHandler.
// It initialises the handler with the provided storage manager and account ID.
func NewAdminHandler(storageManager *storage.RegionStorageManager, accountId string) *AdminHandler {
	return &AdminHandler{
		storageManager: storageManager,
		accountId:      accountId,
	}
}

// getWorkGroupStore retrieves the Athena work group store for the request.
// It extracts the region from the request header and creates a new WorkGroupStore instance.
func (h *AdminHandler) getWorkGroupStore(req *connect.Request[pb.ListWorkGroupsInput]) (*athenastore.WorkGroupStore, error) {
	region := req.Header().Get("X-Aws-Region")
	if region == "" {
		region = "us-east-1"
	}
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return athenastore.NewWorkGroupStore(regionStorage, h.accountId, region), nil
}

// ListWorkGroups lists work groups in Athena.
func (h *AdminHandler) ListWorkGroups(ctx context.Context, req *connect.Request[pb.ListWorkGroupsInput]) (*connect.Response[pb.ListWorkGroupsOutput], error) {
	store, err := h.getWorkGroupStore(req)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	workGroups, err := store.ListWorkGroups()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var summaries []*pb.WorkGroupSummary
	for _, wg := range workGroups {
		state := pb.WorkGroupState_WORK_GROUP_STATE_DISABLED
		if wg.State == "ENABLED" {
			state = pb.WorkGroupState_WORK_GROUP_STATE_ENABLED
		}
		summary := &pb.WorkGroupSummary{
			Name:  wg.Name,
			State: state,
		}
		if wg.Description != "" {
			summary.Description = wg.Description
		}
		if !wg.CreatedTime.IsZero() {
			summary.Creationtime = wg.CreatedTime.Format("2006-01-02T15:04:05.000Z")
		}
		summaries = append(summaries, summary)
	}

	return connect.NewResponse(&pb.ListWorkGroupsOutput{
		Workgroups: summaries,
	}), nil
}

// BatchGetNamedQuery retrieves multiple named queries in a single request.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) BatchGetNamedQuery(ctx context.Context, req *connect.Request[pb.BatchGetNamedQueryInput]) (*connect.Response[pb.BatchGetNamedQueryOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// BatchGetPreparedStatement retrieves multiple prepared statements in a single request.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) BatchGetPreparedStatement(ctx context.Context, req *connect.Request[pb.BatchGetPreparedStatementInput]) (*connect.Response[pb.BatchGetPreparedStatementOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// BatchGetQueryExecution retrieves details of multiple query executions in a single request.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) BatchGetQueryExecution(ctx context.Context, req *connect.Request[pb.BatchGetQueryExecutionInput]) (*connect.Response[pb.BatchGetQueryExecutionOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CancelCapacityReservation cancels a capacity reservation.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CancelCapacityReservation(ctx context.Context, req *connect.Request[pb.CancelCapacityReservationInput]) (*connect.Response[pb.CancelCapacityReservationOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateCapacityReservation creates a capacity reservation for query processing.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateCapacityReservation(ctx context.Context, req *connect.Request[pb.CreateCapacityReservationInput]) (*connect.Response[pb.CreateCapacityReservationOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateDataCatalog creates a data catalog for the Athena metadata repository.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateDataCatalog(ctx context.Context, req *connect.Request[pb.CreateDataCatalogInput]) (*connect.Response[pb.CreateDataCatalogOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateNamedQuery creates a named query in the specified work group.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateNamedQuery(ctx context.Context, req *connect.Request[pb.CreateNamedQueryInput]) (*connect.Response[pb.CreateNamedQueryOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateNotebook creates a notebook in the specified work group.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateNotebook(ctx context.Context, req *connect.Request[pb.CreateNotebookInput]) (*connect.Response[pb.CreateNotebookOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreatePreparedStatement creates a prepared statement in the specified work group.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreatePreparedStatement(ctx context.Context, req *connect.Request[pb.CreatePreparedStatementInput]) (*connect.Response[pb.CreatePreparedStatementOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreatePresignedNotebookUrl generates a presigned URL for accessing a notebook.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreatePresignedNotebookUrl(ctx context.Context, req *connect.Request[pb.CreatePresignedNotebookUrlRequest]) (*connect.Response[pb.CreatePresignedNotebookUrlResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateWorkGroup creates a work group for running queries.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateWorkGroup(ctx context.Context, req *connect.Request[pb.CreateWorkGroupInput]) (*connect.Response[pb.CreateWorkGroupOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteCapacityReservation deletes a capacity reservation.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteCapacityReservation(ctx context.Context, req *connect.Request[pb.DeleteCapacityReservationInput]) (*connect.Response[pb.DeleteCapacityReservationOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteDataCatalog deletes a data catalog.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteDataCatalog(ctx context.Context, req *connect.Request[pb.DeleteDataCatalogInput]) (*connect.Response[pb.DeleteDataCatalogOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteNamedQuery deletes a named query.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteNamedQuery(ctx context.Context, req *connect.Request[pb.DeleteNamedQueryInput]) (*connect.Response[pb.DeleteNamedQueryOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteNotebook deletes a notebook.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteNotebook(ctx context.Context, req *connect.Request[pb.DeleteNotebookInput]) (*connect.Response[pb.DeleteNotebookOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeletePreparedStatement deletes a prepared statement.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeletePreparedStatement(ctx context.Context, req *connect.Request[pb.DeletePreparedStatementInput]) (*connect.Response[pb.DeletePreparedStatementOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteWorkGroup deletes a work group.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteWorkGroup(ctx context.Context, req *connect.Request[pb.DeleteWorkGroupInput]) (*connect.Response[pb.DeleteWorkGroupOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}
