package acm

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"connectrpc.com/connect"
	svcerrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/utils/timeutils"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/acm"
	acmconnect "vorpalstacks/internal/pb/aws/acm/acmconnect"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	acmstore "vorpalstacks/internal/store/aws/acm"
)

// AdminHandler implements the ACM admin console gRPC-Web handler.
type AdminHandler struct {
	acmconnect.UnimplementedACMServiceHandler
	store      acmstore.CertificateStoreInterface
	arnBuilder *acmstore.ARNBuilder
	accountID  string
}

var _ acmconnect.ACMServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new ACM admin handler backed by the given certificate store.
func NewAdminHandler(store acmstore.CertificateStoreInterface, arnBuilder *acmstore.ARNBuilder, accountID string) *AdminHandler {
	return &AdminHandler{store: store, arnBuilder: arnBuilder, accountID: accountID}
}

// ListCertificates returns all ACM certificates visible to the admin console.
func (h *AdminHandler) ListCertificates(ctx context.Context, req *connect.Request[pb.ListCertificatesRequest]) (*connect.Response[pb.ListCertificatesResponse], error) {
	marker := req.Msg.GetNexttoken()
	maxItems := int(req.Msg.GetMaxitems())
	if maxItems <= 0 {
		maxItems = 100
	}

	result, err := h.store.List(marker, maxItems)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	summaries := make([]*pb.CertificateSummary, 0, len(result.Certificates))
	for _, s := range result.Certificates {
		summary := &pb.CertificateSummary{
			Certificatearn:                       s.CertificateArn,
			Domainname:                           s.DomainName,
			Status:                               certificateStatusToProto(s.Status),
			Type:                                 certificateTypeToProto(s.Type),
			Renewaleligibility:                   renewalEligibilityToProto(s.RenewalEligibility),
			Keyalgorithm:                         keyAlgorithmToProto(s.KeyAlgorithm),
			Inuse:                                s.InUse,
			Exported:                             s.Exported,
			Hasadditionalsubjectalternativenames: s.HasAdditionalSubjectAlternativeNames,
		}
		if len(s.SubjectAlternativeNameSummaries) > 0 {
			summary.Subjectalternativenamesummaries = s.SubjectAlternativeNameSummaries
		}
		if s.NotBefore != 0 {
			summary.Notbefore = time.Unix(int64(s.NotBefore), 0).UTC().Format(timeutils.ISO8601UTCFormat)
		}
		if s.NotAfter != 0 {
			summary.Notafter = time.Unix(int64(s.NotAfter), 0).UTC().Format(timeutils.ISO8601UTCFormat)
		}
		if s.CreatedAt != 0 {
			summary.Createdat = time.Unix(int64(s.CreatedAt), 0).UTC().Format(timeutils.ISO8601UTCFormat)
		}
		if s.IssuedAt != 0 {
			summary.Issuedat = time.Unix(int64(s.IssuedAt), 0).UTC().Format(timeutils.ISO8601UTCFormat)
		}
		if s.ImportedAt != 0 {
			summary.Importedat = time.Unix(int64(s.ImportedAt), 0).UTC().Format(timeutils.ISO8601UTCFormat)
		}
		summaries = append(summaries, summary)
	}

	return connect.NewResponse(&pb.ListCertificatesResponse{
		Certificatesummarylist: summaries,
		Nexttoken:              result.NextToken,
	}), nil
}

// RequestCertificate creates a new ACM certificate request with the given domain name.
func (h *AdminHandler) RequestCertificate(ctx context.Context, req *connect.Request[pb.RequestCertificateRequest]) (*connect.Response[pb.RequestCertificateResponse], error) {
	if req.Msg.Domainname == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("DomainName is required"))
	}

	certId := acmstore.GenerateCertificateId()
	certArn := h.arnBuilder.BuildCertificateARN(certId)

	cert := &acmstore.Certificate{
		CertificateArn:          certArn,
		DomainName:              req.Msg.Domainname,
		SubjectAlternativeNames: req.Msg.Subjectalternativenames,
		Status:                  "PENDING_VALIDATION",
		Type:                    "AMAZON_ISSUED",
		KeyAlgorithm:            "RSA_2048",
		RenewalEligibility:      "INELIGIBLE",
		AccountID:               h.accountID,
	}
	if req.Msg.Keyalgorithm != 0 {
		cert.KeyAlgorithm = keyAlgorithmFromProto(req.Msg.Keyalgorithm)
	}
	if req.Msg.Options != nil {
		ctlp := "DISABLED"
		if req.Msg.Options.Certificatetransparencyloggingpreference == pb.CertificateTransparencyLoggingPreference_CERTIFICATE_TRANSPARENCY_LOGGING_PREFERENCE_ENABLED {
			ctlp = "ENABLED"
		}
		cert.Options = &acmstore.CertificateOptions{
			CertificateTransparencyLoggingPreference: ctlp,
		}
	}

	if err := h.store.Create(cert); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.RequestCertificateResponse{
		Certificatearn: certArn,
	}), nil
}

// DeleteCertificate deletes an ACM certificate by its ARN.
func (h *AdminHandler) DeleteCertificate(ctx context.Context, req *connect.Request[pb.DeleteCertificateRequest]) (*connect.Response[pbcommon.Empty], error) {
	if req.Msg.Certificatearn == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("CertificateArn is required"))
	}

	if err := h.store.Delete(req.Msg.Certificatearn); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Acm admin console.
func NewConnectHandler(st storage.BasicStorage, accountID, region string) (string, http.Handler) {
	store := acmstore.NewCertificateStore(st, accountID, region)
	return acmconnect.NewACMServiceHandler(NewAdminHandler(store, acmstore.NewARNBuilder(accountID, region), accountID))
}

func certificateStatusToProto(status string) pb.CertificateStatus {
	switch status {
	case "ISSUED":
		return pb.CertificateStatus_CERTIFICATE_STATUS_ISSUED
	case "PENDING_VALIDATION":
		return pb.CertificateStatus_CERTIFICATE_STATUS_PENDING_VALIDATION
	case "VALIDATION_TIMED_OUT":
		return pb.CertificateStatus_CERTIFICATE_STATUS_VALIDATION_TIMED_OUT
	case "EXPIRED":
		return pb.CertificateStatus_CERTIFICATE_STATUS_EXPIRED
	case "INACTIVE":
		return pb.CertificateStatus_CERTIFICATE_STATUS_INACTIVE
	case "FAILED":
		return pb.CertificateStatus_CERTIFICATE_STATUS_FAILED
	case "REVOKED":
		return pb.CertificateStatus_CERTIFICATE_STATUS_REVOKED
	default:
		return pb.CertificateStatus_CERTIFICATE_STATUS_ISSUED
	}
}

func certificateTypeToProto(t string) pb.CertificateType {
	switch t {
	case "IMPORTED":
		return pb.CertificateType_CERTIFICATE_TYPE_IMPORTED
	case "PRIVATE":
		return pb.CertificateType_CERTIFICATE_TYPE_PRIVATE
	case "AMAZON_ISSUED":
		return pb.CertificateType_CERTIFICATE_TYPE_AMAZON_ISSUED
	default:
		return pb.CertificateType_CERTIFICATE_TYPE_AMAZON_ISSUED
	}
}

func renewalEligibilityToProto(r string) pb.RenewalEligibility {
	switch r {
	case "ELIGIBLE":
		return pb.RenewalEligibility_RENEWAL_ELIGIBILITY_ELIGIBLE
	case "INELIGIBLE":
		return pb.RenewalEligibility_RENEWAL_ELIGIBILITY_INELIGIBLE
	default:
		return pb.RenewalEligibility_RENEWAL_ELIGIBILITY_ELIGIBLE
	}
}

func keyAlgorithmToProto(a string) pb.KeyAlgorithm {
	switch a {
	case "RSA_2048":
		return pb.KeyAlgorithm_KEY_ALGORITHM_RSA_2048
	case "RSA_3072":
		return pb.KeyAlgorithm_KEY_ALGORITHM_RSA_3072
	case "RSA_4096":
		return pb.KeyAlgorithm_KEY_ALGORITHM_RSA_4096
	case "RSA_1024":
		return pb.KeyAlgorithm_KEY_ALGORITHM_RSA_1024
	case "EC_prime256v1":
		return pb.KeyAlgorithm_KEY_ALGORITHM_EC_PRIME256V1
	case "EC_secp384r1":
		return pb.KeyAlgorithm_KEY_ALGORITHM_EC_SECP384R1
	case "EC_secp521r1":
		return pb.KeyAlgorithm_KEY_ALGORITHM_EC_SECP521R1
	default:
		return pb.KeyAlgorithm_KEY_ALGORITHM_RSA_2048
	}
}

func keyAlgorithmFromProto(a pb.KeyAlgorithm) string {
	switch a {
	case pb.KeyAlgorithm_KEY_ALGORITHM_RSA_2048:
		return "RSA_2048"
	case pb.KeyAlgorithm_KEY_ALGORITHM_RSA_3072:
		return "RSA_3072"
	case pb.KeyAlgorithm_KEY_ALGORITHM_RSA_4096:
		return "RSA_4096"
	case pb.KeyAlgorithm_KEY_ALGORITHM_RSA_1024:
		return "RSA_1024"
	case pb.KeyAlgorithm_KEY_ALGORITHM_EC_PRIME256V1:
		return "EC_prime256v1"
	case pb.KeyAlgorithm_KEY_ALGORITHM_EC_SECP384R1:
		return "EC_secp384r1"
	case pb.KeyAlgorithm_KEY_ALGORITHM_EC_SECP521R1:
		return "EC_secp521r1"
	default:
		return "RSA_2048"
	}
}
