package acm

import (
	"context"
	"net/http"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/acm"
	acmconnect "vorpalstacks/internal/pb/aws/acm/acmconnect"
)

// AdminHandler implements the ACM admin console gRPC-Web handler.
type AdminHandler struct {
	acmconnect.UnimplementedACMServiceHandler
}

var _ acmconnect.ACMServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new ACM admin handler.
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// ListCertificates returns all ACM certificates visible to the admin console.
func (h *AdminHandler) ListCertificates(ctx context.Context, req *connect.Request[pb.ListCertificatesRequest]) (*connect.Response[pb.ListCertificatesResponse], error) {
	return connect.NewResponse(&pb.ListCertificatesResponse{
		Certificatesummarylist: []*pb.CertificateSummary{},
	}), nil
}

func NewConnectHandler() (string, http.Handler) {
	return acmconnect.NewACMServiceHandler(NewAdminHandler())
}
