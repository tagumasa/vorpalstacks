package cloudfront

import (
	"context"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/cloudfront"
	cloudfrontconnect "vorpalstacks/internal/pb/aws/cloudfront/cloudfrontconnect"
)

type AdminHandler struct {
	cloudfrontconnect.UnimplementedCloudFrontServiceHandler
}

var _ cloudfrontconnect.CloudFrontServiceHandler = (*AdminHandler)(nil)

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func (h *AdminHandler) ListDistributions(ctx context.Context, req *connect.Request[pb.ListDistributionsRequest]) (*connect.Response[pb.ListDistributionsResult], error) {
	return connect.NewResponse(&pb.ListDistributionsResult{
		Distributionlist: &pb.DistributionList{
			Quantity:    0,
			Items:       []*pb.DistributionSummary{},
			Istruncated: false,
		},
	}), nil
}
