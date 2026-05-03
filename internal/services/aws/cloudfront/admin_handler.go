package cloudfront

import (
	"context"
	"net/http"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/cloudfront"
	cloudfrontconnect "vorpalstacks/internal/pb/aws/cloudfront/cloudfrontconnect"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
)

// AdminHandler implements the CloudFront admin console gRPC-Web handler.
type AdminHandler struct {
	cloudfrontconnect.UnimplementedCloudFrontServiceHandler
	store     storage.BasicStorage
	accountId string
}

var _ cloudfrontconnect.CloudFrontServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new CloudFront admin handler with the given storage and account ID.
func NewAdminHandler(store storage.BasicStorage, accountId string) *AdminHandler {
	return &AdminHandler{
		store:     store,
		accountId: accountId,
	}
}

// ListDistributions returns all CloudFront distributions visible to the admin console.
func (h *AdminHandler) ListDistributions(ctx context.Context, req *connect.Request[pb.ListDistributionsRequest]) (*connect.Response[pb.ListDistributionsResult], error) {
	distStore := cloudfrontstore.NewDistributionStore(h.store, h.accountId)

	maxItems := int(req.Msg.Maxitems)
	if maxItems <= 0 {
		maxItems = 100
	}

	result, err := distStore.List(req.Msg.Marker, maxItems)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var items []*pb.DistributionSummary
	for _, d := range result.Distributions {
		summary := &pb.DistributionSummary{
			Id:       d.ID,
			Arn:      d.ARN,
			Status:   d.Status,
			Enabled:  d.Enabled,
			Staging:  d.Staging,
			Etag:     d.ETag,
			Comment:  d.DistributionConfig.Comment,
			Domainname: d.DomainName,
		}
		if !d.LastModifiedAt.IsZero() {
			summary.Lastmodifiedtime = d.LastModifiedAt.Format("2006-01-02T15:04:05.000Z")
		}
		items = append(items, summary)
	}

	return connect.NewResponse(&pb.ListDistributionsResult{
		Distributionlist: &pb.DistributionList{
			Quantity:    int32(len(items)),
			Items:       items,
			Istruncated: result.IsTruncated,
			Nextmarker:  result.NextMarker,
			Marker:      req.Msg.Marker,
			Maxitems:    int32(maxItems),
		},
	}), nil
}

// DeleteDistribution deletes a CloudFront distribution via the admin console.
func (h *AdminHandler) DeleteDistribution(ctx context.Context, req *connect.Request[pb.DeleteDistributionRequest]) (*connect.Response[pbcommon.Empty], error) {
	distStore := cloudfrontstore.NewDistributionStore(h.store, h.accountId)

	if err := distStore.Delete(req.Msg.Id); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Cloudfront admin console.
func NewConnectHandler(store storage.BasicStorage, accountID string) (string, http.Handler) {
	return cloudfrontconnect.NewCloudFrontServiceHandler(NewAdminHandler(store, accountID))
}
