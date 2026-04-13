package timestreamquery

import (
	"context"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/timestreamquery"
	timestreamqueryconnect "vorpalstacks/internal/pb/aws/timestreamquery/timestreamqueryconnect"
	svccommon "vorpalstacks/internal/common"
	timestreamstore "vorpalstacks/internal/store/aws/timestream"
)

// AdminHandler provides Timestream Query service administration functionality.
// It implements the TimestreamQueryServiceHandler interface for gRPC-Web communication.
type AdminHandler struct {
	timestreamqueryconnect.UnimplementedTimestreamQueryServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
	dataPath       string
}

var _ timestreamqueryconnect.TimestreamQueryServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Timestream Query AdminHandler.
// It initialises the handler with the provided storage manager, account ID, and data path.
func NewAdminHandler(storageManager *storage.RegionStorageManager, accountId, dataPath string) *AdminHandler {
	return &AdminHandler{
		storageManager: storageManager,
		accountId:      accountId,
		dataPath:       dataPath,
	}
}

// getScheduledQueryStore retrieves the Timestream Query scheduled query store for the request.
// It extracts the region from the request header and creates a new ScheduledQueryStore instance.
func (h *AdminHandler) getScheduledQueryStore(req *connect.Request[pb.ListScheduledQueriesRequest]) (*timestreamstore.ScheduledQueryStore, error) {
	region := svccommon.GetRegionFromHeader(req.Header())
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return timestreamstore.NewScheduledQueryStore(regionStorage, h.accountId, region), nil
}

// ListScheduledQueries lists scheduled queries in Timestream Query.
func (h *AdminHandler) ListScheduledQueries(ctx context.Context, req *connect.Request[pb.ListScheduledQueriesRequest]) (*connect.Response[pb.ListScheduledQueriesResponse], error) {
	store, err := h.getScheduledQueryStore(req)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	queries, err := store.ListScheduledQueries()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var summaries []*pb.ScheduledQuery
	for _, sq := range queries {
		summary := &pb.ScheduledQuery{
			Arn:          sq.ARN,
			Name:         sq.Name,
			Creationtime: sq.CreationTime.Format("2006-01-02T15:04:05.000Z"),
		}
		if !sq.PreviousRunTime.IsZero() {
			summary.Previousinvocationtime = sq.PreviousRunTime.Format("2006-01-02T15:04:05.000Z")
		}
		if !sq.NextRunTime.IsZero() {
			summary.Nextinvocationtime = sq.NextRunTime.Format("2006-01-02T15:04:05.000Z")
		}
		summaries = append(summaries, summary)
	}

	return connect.NewResponse(&pb.ListScheduledQueriesResponse{
		Scheduledqueries: summaries,
	}), nil
}
