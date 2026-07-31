package acm

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"
	svccommon "vorpalstacks/internal/common"
	svcerrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/utils/aws/types"
	"vorpalstacks/internal/utils/timeutils"

	pb "vorpalstacks/internal/pb/aws/acm"
	acmconnect "vorpalstacks/internal/pb/aws/acm/acmconnect"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	acmstore "vorpalstacks/internal/store/aws/acm"
)

// AdminHandler implements the ACM admin console gRPC-Web handler.
// It delegates to the shared ACMService store cache so that the same
// per-region store instances are used by both the HTTP API handlers and the
// admin console gRPC-Web handlers.
type AdminHandler struct {
	acmconnect.UnimplementedACMServiceHandler
	service *ACMService
}

var _ acmconnect.ACMServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new ACM admin handler backed by the given service instance.
func NewAdminHandler(svc *ACMService) *AdminHandler {
	return &AdminHandler{service: svc}
}

func (h *AdminHandler) getStoreFromHeaders(headers http.Header) (*acmStores, error) {
	region := svccommon.GetRegionFromHeader(headers)
	return h.service.GetStoreForRegion(region)
}

// ListCertificates returns all ACM certificates visible to the admin console.
func (h *AdminHandler) ListCertificates(ctx context.Context, req *connect.Request[pb.ListCertificatesRequest]) (*connect.Response[pb.ListCertificatesResponse], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	marker := req.Msg.GetNexttoken()
	maxItems := int(req.Msg.GetMaxitems())
	if maxItems <= 0 {
		maxItems = 100
	}
	// Smithy MaxItems: @range(1, 1000). Reject out-of-range instead of clamping.
	if maxItems > 1000 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("MaxItems must not exceed 1000"))
	}

	result, err := stores.certificates.List(marker, maxItems)
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
			Inuse:                                proto.Bool(s.InUse),
			Exported:                             proto.Bool(s.Exported),
			Hasadditionalsubjectalternativenames: proto.Bool(s.HasAdditionalSubjectAlternativeNames),
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
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	domain, err := validateDomainName(req.Msg.Domainname)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// Validate SubjectAlternativeNames (Smithy DomainList: min 1, max 100,
	// each entry must be a valid domain name).
	sans := req.Msg.Subjectalternativenames
	if len(sans) > 100 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("SubjectAlternativeNames must not exceed 100 entries"))
	}
	for _, san := range sans {
		if _, err := validateDomainName(san); err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid SAN %q: %w", san, err))
		}
	}

	// Smithy IdempotencyToken: @length(1-32) + @pattern(^\w+$).
	if err := validateIdempotencyToken(req.Msg.Idempotencytoken); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	keyAlgorithm := "RSA_2048"
	if req.Msg.Keyalgorithm != 0 {
		keyAlgorithm = keyAlgorithmFromProto(req.Msg.Keyalgorithm)
	}

	validationMethod := "DNS"
	if req.Msg.Validationmethod == pb.ValidationMethod_VALIDATION_METHOD_EMAIL {
		validationMethod = "EMAIL"
	}

	certId := acmstore.GenerateCertificateId()
	certArn := stores.arnBuilder.BuildCertificateARN(certId)

	cert, err := generateAmazonIssuedCert(certArn, domain, sans, keyAlgorithm, validationMethod)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	cert.AccountID = h.service.accountID

	// ManagedBy (Smithy CertificateManagedBy enum: CLOUDFRONT only).
	if req.Msg.Managedby == pb.CertificateManagedBy_CERTIFICATE_MANAGED_BY_CLOUDFRONT {
		cert.ManagedBy = "CLOUDFRONT"
	}

	// CertificateAuthorityArn (Smithy PcaArn: @length 20-2048 + pattern).
	pcaArn, err := validateCertificateAuthorityArn(req.Msg.Certificateauthorityarn)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	cert.CertificateAuthorityArn = pcaArn

	// Tags: convert proto Tags to types.Tag, validate, set on cert.
	if len(req.Msg.Tags) > 0 {
		tags := make([]types.Tag, 0, len(req.Msg.Tags))
		for _, t := range req.Msg.Tags {
			tags = append(tags, types.Tag{Key: t.Key, Value: t.Value})
		}
		if len(tags) > 50 {
			return nil, connect.NewError(connect.CodeInvalidArgument, NewTooManyTagsException("Tags must not exceed 50 entries"))
		}
		if err := validateACMTags(tags); err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		cert.Tags = tags
	}

	// DomainValidationOptions: merge user-provided ValidationDomain overrides
	// into the auto-generated DVOs from generateAmazonIssuedCert.
	for _, dvo := range req.Msg.Domainvalidationoptions {
		if err := validateDomainValidationFields(dvo.Domainname, dvo.Validationdomain); err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		for _, dv := range cert.DomainValidationOptions {
			if strings.EqualFold(dv.DomainName, dvo.Domainname) {
				dv.ValidationDomain = strings.ToLower(dvo.Validationdomain)
			}
		}
	}

	if req.Msg.Options != nil {
		ctlp := "DISABLED"
		if req.Msg.Options.Certificatetransparencyloggingpreference == pb.CertificateTransparencyLoggingPreference_CERTIFICATE_TRANSPARENCY_LOGGING_PREFERENCE_ENABLED {
			ctlp = "ENABLED"
		}
		exportOpt := "DISABLED"
		if req.Msg.Options.Export == pb.CertificateExport_CERTIFICATE_EXPORT_ENABLED {
			exportOpt = "ENABLED"
		}
		cert.Options = &acmstore.CertificateOptions{
			CertificateTransparencyLoggingPreference: ctlp,
			Export:                                   exportOpt,
		}
	} else {
		// Match the HTTP API default (certificate_operations.go:226-231):
		// CTLPreference=ENABLED, Export=DISABLED.
		cert.Options = &acmstore.CertificateOptions{
			CertificateTransparencyLoggingPreference: "ENABLED",
			Export:                                   "DISABLED",
		}
	}

	if err := stores.certificates.Create(cert); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.RequestCertificateResponse{
		Certificatearn: certArn,
	}), nil
}

// DeleteCertificate deletes an ACM certificate by its ARN.
func (h *AdminHandler) DeleteCertificate(ctx context.Context, req *connect.Request[pb.DeleteCertificateRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if req.Msg.Certificatearn == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("CertificateArn is required"))
	}
	if err := validateCertificateArn(req.Msg.Certificatearn); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	cert, err := stores.certificates.Get(req.Msg.Certificatearn)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}
	if len(cert.InUseBy) > 0 {
		return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf("certificate %s is in use", req.Msg.Certificatearn))
	}

	if err := stores.certificates.Delete(req.Msg.Certificatearn); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Acm admin console.
func NewConnectHandler(svc *ACMService) (string, http.Handler) {
	return acmconnect.NewACMServiceHandler(NewAdminHandler(svc))
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
