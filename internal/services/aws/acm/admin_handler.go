package acm

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/acm"
	acmconnect "vorpalstacks/internal/pb/aws/acm/acmconnect"
	"vorpalstacks/internal/pb/aws/common"
)

// AdminHandler provides ACM service administration functionality.
type AdminHandler struct{}

var _ acmconnect.ACMServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new ACM AdminHandler.
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// ListCertificates lists certificates in ACM.
func (h *AdminHandler) ListCertificates(ctx context.Context, req *connect.Request[pb.ListCertificatesRequest]) (*connect.Response[pb.ListCertificatesResponse], error) {
	return connect.NewResponse(&pb.ListCertificatesResponse{
		Certificatesummarylist: []*pb.CertificateSummary{},
	}), nil
}

// AddTagsToCertificate adds one or more tags to an ACM certificate.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) AddTagsToCertificate(ctx context.Context, req *connect.Request[pb.AddTagsToCertificateRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("acm.ACMService.AddTagsToCertificate is not implemented"))
}

// DeleteCertificate deletes an ACM certificate.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteCertificate(ctx context.Context, req *connect.Request[pb.DeleteCertificateRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("acm.ACMService.DeleteCertificate is not implemented"))
}

// DescribeCertificate returns detailed information about an ACM certificate.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeCertificate(ctx context.Context, req *connect.Request[pb.DescribeCertificateRequest]) (*connect.Response[pb.DescribeCertificateResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("acm.ACMService.DescribeCertificate is not implemented"))
}

// ExportCertificate exports an ACM certificate and its private key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ExportCertificate(ctx context.Context, req *connect.Request[pb.ExportCertificateRequest]) (*connect.Response[pb.ExportCertificateResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("acm.ACMService.ExportCertificate is not implemented"))
}

// GetAccountConfiguration retrieves the ACM account configuration.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetAccountConfiguration(ctx context.Context, req *connect.Request[common.Empty]) (*connect.Response[pb.GetAccountConfigurationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("acm.ACMService.GetAccountConfiguration is not implemented"))
}

// GetCertificate retrieves an ACM certificate and its certificate chain.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetCertificate(ctx context.Context, req *connect.Request[pb.GetCertificateRequest]) (*connect.Response[pb.GetCertificateResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("acm.ACMService.GetCertificate is not implemented"))
}

// ImportCertificate imports an SSL/TLS certificate into ACM.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ImportCertificate(ctx context.Context, req *connect.Request[pb.ImportCertificateRequest]) (*connect.Response[pb.ImportCertificateResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("acm.ACMService.ImportCertificate is not implemented"))
}

// ListTagsForCertificate lists the tags associated with an ACM certificate.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListTagsForCertificate(ctx context.Context, req *connect.Request[pb.ListTagsForCertificateRequest]) (*connect.Response[pb.ListTagsForCertificateResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("acm.ACMService.ListTagsForCertificate is not implemented"))
}

// PutAccountConfiguration updates the ACM account configuration.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutAccountConfiguration(ctx context.Context, req *connect.Request[pb.PutAccountConfigurationRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("acm.ACMService.PutAccountConfiguration is not implemented"))
}

// RemoveTagsFromCertificate removes one or more tags from an ACM certificate.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) RemoveTagsFromCertificate(ctx context.Context, req *connect.Request[pb.RemoveTagsFromCertificateRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("acm.ACMService.RemoveTagsFromCertificate is not implemented"))
}

// RenewCertificate renews an ACM certificate.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) RenewCertificate(ctx context.Context, req *connect.Request[pb.RenewCertificateRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("acm.ACMService.RenewCertificate is not implemented"))
}

// RequestCertificate requests an ACM certificate.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) RequestCertificate(ctx context.Context, req *connect.Request[pb.RequestCertificateRequest]) (*connect.Response[pb.RequestCertificateResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("acm.ACMService.RequestCertificate is not implemented"))
}

// ResendValidationEmail resends the validation email for an ACM certificate.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ResendValidationEmail(ctx context.Context, req *connect.Request[pb.ResendValidationEmailRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("acm.ACMService.ResendValidationEmail is not implemented"))
}

// RevokeCertificate revokes an ACM certificate.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) RevokeCertificate(ctx context.Context, req *connect.Request[pb.RevokeCertificateRequest]) (*connect.Response[pb.RevokeCertificateResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("acm.ACMService.RevokeCertificate is not implemented"))
}

// UpdateCertificateOptions updates the options for an ACM certificate.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateCertificateOptions(ctx context.Context, req *connect.Request[pb.UpdateCertificateOptionsRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("acm.ACMService.UpdateCertificateOptions is not implemented"))
}
