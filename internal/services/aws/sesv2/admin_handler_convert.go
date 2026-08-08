package sesv2

import (
	"net/http"

	svccommon "vorpalstacks/internal/common"
	pb "vorpalstacks/internal/pb/aws/sesv2"
	sesv2store "vorpalstacks/internal/store/aws/sesv2"
)

// ---------------------------------------------------------------------------
// admin_handler_convert.go — the sole file in the SESv2 service package
// that imports store packages and performs proto↔DTO conversion. This
// enforces the store-import prohibition: admin handlers must not import store
// packages directly.
// ---------------------------------------------------------------------------

// getSESv2Store returns the SESv2 store for the region extracted from the
// request header.
func (h *AdminHandler) getSESv2Store(headers http.Header) (sesv2store.SESv2StoreInterface, error) {
	region := svccommon.GetRegionFromHeader(headers)
	return h.service.GetStoreForRegion(region)
}

// toPbIdentityInfos converts service-layer IdentitySummary values to the
// proto IdentityInfo list.
func toPbIdentityInfos(items []IdentitySummary) []*pb.IdentityInfo {
	result := make([]*pb.IdentityInfo, 0, len(items))
	for _, item := range items {
		info := &pb.IdentityInfo{
			Identityname:       item.IdentityName,
			Identitytype:       pb.IdentityType_IDENTITY_TYPE_EMAIL_ADDRESS,
			Sendingenabled:     boolPtr(item.SendingEnabled),
			Verificationstatus: verificationStatusToProtoFromString(item.VerificationStatus),
		}
		if item.IdentityType == "DOMAIN" {
			info.Identitytype = pb.IdentityType_IDENTITY_TYPE_DOMAIN
		}
		result = append(result, info)
	}
	return result
}

// toPbCreateEmailIdentityResponse converts a service-layer IdentityResult
// to the proto response type.
func toPbCreateEmailIdentityResponse(r *IdentityResult) *pb.CreateEmailIdentityResponse {
	resp := &pb.CreateEmailIdentityResponse{
		Identitytype:             pb.IdentityType_IDENTITY_TYPE_EMAIL_ADDRESS,
		Verifiedforsendingstatus: boolPtr(r.VerifiedForSending),
	}
	if r.IdentityType == "DOMAIN" {
		resp.Identitytype = pb.IdentityType_IDENTITY_TYPE_DOMAIN
	}
	if r.DkimAttributes != nil {
		tokens, _ := r.DkimAttributes["Tokens"].([]string)
		statusStr, _ := r.DkimAttributes["Status"].(string)
		// Read the actual SigningEnabled value from the
		// identity map instead of hardcoding true. The HTTP API path
		// correctly surfaces dkim.SigningEnabled; the admin path must
		// do the same.
		signingEnabled := true
		if v, ok := r.DkimAttributes["SigningEnabled"].(bool); ok {
			signingEnabled = v
		}
		resp.Dkimattributes = &pb.DkimAttributes{
			Signingenabled: boolPtr(signingEnabled),
			Tokens:         tokens,
			Status:         dkimStatusToProtoFromString(statusStr),
		}
	}
	return resp
}

func boolPtr(v bool) *bool {
	return &v
}

// verificationStatusToProtoFromString converts a string status to the
// proto VerificationStatus enum.
func verificationStatusToProtoFromString(status string) pb.VerificationStatus {
	switch status {
	case "SUCCESS":
		return pb.VerificationStatus_VERIFICATION_STATUS_SUCCESS
	case "PENDING":
		return pb.VerificationStatus_VERIFICATION_STATUS_PENDING
	case "FAILED":
		return pb.VerificationStatus_VERIFICATION_STATUS_FAILED
	case "TEMPORARY_FAILURE":
		return pb.VerificationStatus_VERIFICATION_STATUS_TEMPORARY_FAILURE
	case "NOT_STARTED":
		return pb.VerificationStatus_VERIFICATION_STATUS_NOT_STARTED
	default:
		return pb.VerificationStatus_VERIFICATION_STATUS_PENDING
	}
}

// dkimStatusToProtoFromString converts a string status to the proto
// DkimStatus enum.
func dkimStatusToProtoFromString(status string) pb.DkimStatus {
	switch status {
	case "SUCCESS":
		return pb.DkimStatus_DKIM_STATUS_SUCCESS
	case "PENDING":
		return pb.DkimStatus_DKIM_STATUS_PENDING
	case "FAILED":
		return pb.DkimStatus_DKIM_STATUS_FAILED
	case "TEMPORARY_FAILURE":
		return pb.DkimStatus_DKIM_STATUS_TEMPORARY_FAILURE
	case "NOT_STARTED":
		return pb.DkimStatus_DKIM_STATUS_NOT_STARTED
	default:
		return pb.DkimStatus_DKIM_STATUS_PENDING
	}
}
