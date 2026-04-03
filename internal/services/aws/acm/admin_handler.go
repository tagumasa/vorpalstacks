package acm

import (
	"context"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/acm"
	acmconnect "vorpalstacks/internal/pb/aws/acm/acmconnect"
)

type AdminHandler struct {
	acmconnect.UnimplementedACMServiceHandler
}

var _ acmconnect.ACMServiceHandler = (*AdminHandler)(nil)

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func (h *AdminHandler) ListCertificates(ctx context.Context, req *connect.Request[pb.ListCertificatesRequest]) (*connect.Response[pb.ListCertificatesResponse], error) {
	return connect.NewResponse(&pb.ListCertificatesResponse{
		Certificatesummarylist: []*pb.CertificateSummary{},
	}), nil
}
