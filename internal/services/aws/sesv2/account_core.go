package sesv2

import (
	"context"

	"vorpalstacks/internal/common/request"
	sesv2store "vorpalstacks/internal/store/aws/sesv2"
)

// ---------------------------------------------------------------------------
// Input DTOs — account family
// ---------------------------------------------------------------------------

// PutAccountSendingAttributesInput carries the account sending toggle.
type PutAccountSendingAttributesInput struct {
	SendingEnabled bool
}

// PutAccountSuppressionAttributesInput carries the account suppression
// reasons and the raw ValidationAttributes wire map (the nested
// ConditionThreshold/OverallConfidenceThreshold structure is interpreted
// by the Core).
type PutAccountSuppressionAttributesInput struct {
	SuppressedReasons    []string
	ValidationAttributes map[string]interface{}
}

// PutAccountDetailsInput carries the account-details members. MailType and
// WebsiteURL are Smithy @required; ProductionAccessEnabled is an
// EnabledWrapper whose presence is carried by ProductionAccessProvided.
type PutAccountDetailsInput struct {
	MailType                        string
	WebsiteURL                      string
	UseCaseDescription              string
	ContactLanguage                 string
	AdditionalContactEmailAddresses []string
	ProductionAccessEnabled         bool
	ProductionAccessProvided        bool
}

// PutAccountVdmAttributesInput carries the raw VdmAttributes wire map.
type PutAccountVdmAttributesInput struct {
	VdmAttributes map[string]interface{}
}

// PutAccountDedicatedIpWarmupAttributesInput carries the dedicated IP auto
// warmup toggle.
type PutAccountDedicatedIpWarmupAttributesInput struct {
	AutoWarmupEnabled bool
}

// ---------------------------------------------------------------------------
// Core functions — account family
// ---------------------------------------------------------------------------

// getAccountCore is the single entry point for reading the SES v2 account
// state shared by the HTTP API handler.
func (s *SESv2Service) getAccountCore(ctx context.Context, store sesv2store.SESv2StoreInterface) (map[string]interface{}, error) {
	account, err := store.GetAccount()
	if err != nil {
		return nil, err
	}

	// Defensive nil guards for sub-structs — a stored account from an
	// older format may have nil SendingAttributes or SuppressionAttributes.
	if account.SendingAttributes == nil {
		account.SendingAttributes = &sesv2store.SendingAttributes{}
	}
	if account.SuppressionAttributes == nil {
		account.SuppressionAttributes = &sesv2store.SuppressionAttributes{}
	}

	result := map[string]interface{}{
		"DedicatedIpAutoWarmupEnabled": account.SendingAttributes.DedicatedIpAutoWarmup,
		"EnforcementStatus":            account.EnforcementStatus,
		"ProductionAccessEnabled":      account.ProductionAccessEnabled,
		"SendingEnabled":               account.SendingEnabled,
		"SuppressionAttributes": map[string]interface{}{
			"SuppressedReasons": account.SuppressionAttributes.SuppressedReasons,
		},
	}

	if account.Details != nil {
		details := map[string]interface{}{
			"MailType":                        account.Details.MailType,
			"AdditionalContactEmailAddresses": account.Details.AdditionalContactEmailAddresses,
			"UseCaseDescription":              account.Details.UseCaseDescription,
			"WebsiteURL":                      account.Details.WebsiteURL,
			"ContactLanguage":                 account.Details.ContactLanguage,
		}
		// ReviewDetails is a read-only status managed by the service.
		// In our edge/on-prem build production access is auto-granted.
		if account.ProductionAccessEnabled {
			details["ReviewDetails"] = map[string]interface{}{
				"Status": "GRANTED",
			}
		}
		result["Details"] = details
	}

	if account.VdmAttributes != nil {
		vdm := map[string]interface{}{}
		if account.VdmAttributes.VdmEnabled {
			vdm["VdmEnabled"] = "ENABLED"
		} else {
			vdm["VdmEnabled"] = "DISABLED"
		}
		if account.VdmAttributes.DashboardAttributes != "" {
			vdm["DashboardAttributes"] = map[string]interface{}{
				"EngagementMetrics": account.VdmAttributes.DashboardAttributes,
			}
		}
		if account.VdmAttributes.GuardianAttributes != "" {
			vdm["GuardianAttributes"] = map[string]interface{}{
				"OptimizedSharedDelivery": account.VdmAttributes.GuardianAttributes,
			}
		}
		result["VdmAttributes"] = vdm
	}

	if account.SendingAttributes != nil {
		result["SendQuota"] = map[string]interface{}{
			"Max24HourSend":   account.SendingAttributes.Max24HourSend,
			"MaxSendRate":     account.SendingAttributes.MaxSendRate,
			"SentLast24Hours": account.SendingAttributes.SentLast24Hours,
		}
	}

	return result, nil
}

// putAccountSendingAttributesCore is the single entry point for the account
// sending toggle.
func (s *SESv2Service) putAccountSendingAttributesCore(store sesv2store.SESv2StoreInterface, in PutAccountSendingAttributesInput) error {
	return store.PutSendingAttributes(in.SendingEnabled)
}

// putAccountSuppressionAttributesCore is the single entry point for the
// account suppression attributes, including the ValidationAttributes
// threshold structure.
func (s *SESv2Service) putAccountSuppressionAttributesCore(store sesv2store.SESv2StoreInterface, in PutAccountSuppressionAttributesInput) error {
	for _, r := range in.SuppressedReasons {
		if !validateSuppressionListReason(r) {
			return ErrBadRequest
		}
	}

	attrs := &sesv2store.SuppressionAttributes{
		SuppressedReasons: in.SuppressedReasons,
	}

	if va := in.ValidationAttributes; va != nil {
		if ct := request.GetMapParam(va, "ConditionThreshold"); ct != nil {
			thresholdEnabled := request.GetStringParam(ct, "ConditionThresholdEnabled")
			if thresholdEnabled != "" && !validateFeatureStatus(thresholdEnabled) {
				return ErrBadRequest
			}
			threshold := &sesv2store.SuppressionConditionThreshold{
				ConditionThresholdEnabled: thresholdEnabled,
			}
			if oct := request.GetMapParam(ct, "OverallConfidenceThreshold"); oct != nil {
				cvt := request.GetStringParam(oct, "ConfidenceVerdictThreshold")
				if cvt != "" && !validateSuppressionConfidenceVerdictThreshold(cvt) {
					return ErrBadRequest
				}
				threshold.OverallConfidenceThreshold = &sesv2store.SuppressionConfidenceThreshold{
					ConfidenceVerdictThreshold: cvt,
				}
			}
			attrs.ValidationAttributes = &sesv2store.SuppressionValidationAttributes{
				ConditionThreshold: threshold,
			}
		}
	}

	return store.PutSuppressionAttributes(attrs)
}

// putAccountDetailsCore is the single entry point for the account details.
// Per Smithy com.amazonaws.sesv2#PutAccountDetailsRequest, MailType and
// WebsiteURL are @required. ProductionAccessEnabled has Smithy type
// EnabledWrapper (a Boolean wrapper that distinguishes unset from explicit
// false); ProductionAccessProvided carries that distinction.
func (s *SESv2Service) putAccountDetailsCore(store sesv2store.SESv2StoreInterface, in PutAccountDetailsInput) error {
	// MailType and WebsiteURL are Smithy @required.
	if in.MailType == "" {
		return ErrMissingParameter
	}
	if in.WebsiteURL == "" {
		return ErrMissingParameter
	}
	// Validate enum values.
	if !validateMailType(in.MailType) {
		return ErrBadRequest
	}
	if in.ContactLanguage != "" {
		if !validateContactLanguage(in.ContactLanguage) {
			return ErrBadRequest
		}
	}
	// AdditionalContactEmailAddresses: max 5, validate email format.
	if len(in.AdditionalContactEmailAddresses) > 5 {
		return ErrBadRequest
	}
	for _, e := range in.AdditionalContactEmailAddresses {
		if !validateEmailAddressRFC5321(e) {
			return ErrBadRequest
		}
	}

	details := &sesv2store.AccountDetails{
		MailType:                        in.MailType,
		UseCaseDescription:              in.UseCaseDescription,
		WebsiteURL:                      in.WebsiteURL,
		ContactLanguage:                 in.ContactLanguage,
		AdditionalContactEmailAddresses: in.AdditionalContactEmailAddresses,
	}

	if in.ProductionAccessProvided {
		details.ProductionAccessEnabled = in.ProductionAccessEnabled
		details.ProductionAccessProvided = true
	}

	return store.PutAccountDetails(details)
}

// putAccountVdmAttributesCore is the single entry point for the account VDM
// attributes. Per Smithy `com.amazonaws.sesv2#VdmAttributes`,
// DashboardAttributes and GuardianAttributes are nested structs (not flat
// strings).
func (s *SESv2Service) putAccountVdmAttributesCore(store sesv2store.SESv2StoreInterface, in PutAccountVdmAttributesInput) error {
	vdmAttrs := in.VdmAttributes
	if vdmAttrs == nil {
		return ErrMissingParameter
	}

	vdmEnabledStr := request.GetStringParam(vdmAttrs, "VdmEnabled")
	if vdmEnabledStr != "ENABLED" && vdmEnabledStr != "DISABLED" {
		return ErrBadRequest
	}
	vdm := &sesv2store.VdmAttributes{
		VdmEnabled:                      vdmEnabledStr == "ENABLED",
		AdditionalContactEmailAddresses: request.GetStringList(vdmAttrs, "AdditionalContactEmailAddresses"),
	}

	if dashboardAttrs := request.GetMapParam(vdmAttrs, "DashboardAttributes"); dashboardAttrs != nil {
		em := request.GetStringParam(dashboardAttrs, "EngagementMetrics")
		if em != "" && !validateFeatureStatus(em) {
			return ErrBadRequest
		}
		vdm.DashboardAttributes = em
	}

	if guardianAttrs := request.GetMapParam(vdmAttrs, "GuardianAttributes"); guardianAttrs != nil {
		osd := request.GetStringParam(guardianAttrs, "OptimizedSharedDelivery")
		if osd != "" && !validateFeatureStatus(osd) {
			return ErrBadRequest
		}
		vdm.GuardianAttributes = osd
	}

	return store.PutVdmAttributes(vdm)
}

// putAccountDedicatedIpWarmupAttributesCore is the single entry point for
// the dedicated IP auto warmup setting.
func (s *SESv2Service) putAccountDedicatedIpWarmupAttributesCore(store sesv2store.SESv2StoreInterface, in PutAccountDedicatedIpWarmupAttributesInput) error {
	return store.PutDedicatedIpAutoWarmupEnabled(in.AutoWarmupEnabled)
}

// ---------------------------------------------------------------------------
// HTTP handlers — parse → DTO → Core → serialise
// ---------------------------------------------------------------------------
