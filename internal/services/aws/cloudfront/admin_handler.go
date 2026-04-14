package cloudfront

import (
	"context"
	"net/http"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/cloudfront"
	cloudfrontconnect "vorpalstacks/internal/pb/aws/cloudfront/cloudfrontconnect"
)

// AdminHandler implements the CloudFront admin console gRPC-Web handler.
type AdminHandler struct {
	cloudfrontconnect.UnimplementedCloudFrontServiceHandler
}

var _ cloudfrontconnect.CloudFrontServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new CloudFront admin handler.
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// ListDistributions returns all CloudFront distributions visible to the admin console.
func (h *AdminHandler) ListDistributions(ctx context.Context, req *connect.Request[pb.ListDistributionsRequest]) (*connect.Response[pb.ListDistributionsResult], error) {
	return connect.NewResponse(&pb.ListDistributionsResult{
		Distributionlist: &pb.DistributionList{
			Quantity:    0,
			Items:       []*pb.DistributionSummary{},
			Istruncated: false,
		},
	}), nil
}

func NewConnectHandler() (string, http.Handler) {
	return cloudfrontconnect.NewCloudFrontServiceHandler(NewAdminHandler())
}
