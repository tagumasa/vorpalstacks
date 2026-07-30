package sesv2

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	sesv2store "vorpalstacks/internal/store/aws/sesv2"
)

// GetAccount retrieves the details of the current SES v2 account.
func (s *SESv2Service) GetAccount(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

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

// PutAccountSendingAttributes updates the sending attributes for the SES v2 account.
func (s *SESv2Service) PutAccountSendingAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	sendingEnabled := request.GetBoolParam(req.Parameters, "SendingEnabled")
	if err := store.PutSendingAttributes(sendingEnabled); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// PutAccountSuppressionAttributes updates the suppression attributes for the SES v2 account.
// Per Smithy com.amazonaws.sesv2#PutAccountSuppressionAttributesRequest
// the input carries SuppressedReasons and ValidationAttributes (Auto
// Validation threshold settings).
func (s *SESv2Service) PutAccountSuppressionAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	attrs := &sesv2store.SuppressionAttributes{
		SuppressedReasons: request.GetStringList(req.Parameters, "SuppressedReasons"),
	}

	if va := request.GetMapParam(req.Parameters, "ValidationAttributes"); va != nil {
		if ct := request.GetMapParam(va, "ConditionThreshold"); ct != nil {
			threshold := &sesv2store.SuppressionConditionThreshold{
				ConditionThresholdEnabled: request.GetStringParam(ct, "ConditionThresholdEnabled"),
			}
			if oct := request.GetMapParam(ct, "OverallConfidenceThreshold"); oct != nil {
				threshold.OverallConfidenceThreshold = &sesv2store.SuppressionConfidenceThreshold{
					ConfidenceVerdictThreshold: request.GetStringParam(oct, "ConfidenceVerdictThreshold"),
				}
			}
			attrs.ValidationAttributes = &sesv2store.SuppressionValidationAttributes{
				ConditionThreshold: threshold,
			}
		}
	}

	if err := store.PutSuppressionAttributes(attrs); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// PutAccountDetails updates the account details for the SES v2 account.
// Per Smithy com.amazonaws.sesv2#PutAccountDetailsRequest,
// ProductionAccessEnabled has Smithy type EnabledWrapper (a Boolean
// wrapper that distinguishes unset from explicit false). We honour that
// by only overwriting the stored value when the caller actually supplied
// the parameter.
func (s *SESv2Service) PutAccountDetails(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	details := &sesv2store.AccountDetails{
		MailType:           request.GetStringParam(req.Parameters, "MailType"),
		UseCaseDescription: request.GetStringParam(req.Parameters, "UseCaseDescription"),
		WebsiteURL:         request.GetStringParam(req.Parameters, "WebsiteURL"),
		ContactLanguage:    request.GetStringParam(req.Parameters, "ContactLanguage"),
	}

	additionalEmails := request.GetStringList(req.Parameters, "AdditionalContactEmailAddresses")
	if len(additionalEmails) > 0 {
		details.AdditionalContactEmailAddresses = additionalEmails
	}

	if _, ok := req.Parameters["ProductionAccessEnabled"]; ok {
		details.ProductionAccessEnabled = request.GetBoolParam(req.Parameters, "ProductionAccessEnabled")
		details.ProductionAccessProvided = true
	}

	if err := store.PutAccountDetails(details); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// PutAccountVdmAttributes updates the VDM attributes for the SES v2 account.
// Per Smithy `com.amazonaws.sesv2#VdmAttributes`, DashboardAttributes and
// GuardianAttributes are nested structs (not flat strings). The previous
// implementation read them via GetStringParam, silently dropping the nested
// EngagementMetrics / OptimizedSharedDelivery values.
func (s *SESv2Service) PutAccountVdmAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	vdmAttrs := request.GetMapParam(req.Parameters, "VdmAttributes")
	if vdmAttrs == nil {
		return nil, ErrMissingParameter
	}

	vdmEnabledStr := request.GetStringParam(vdmAttrs, "VdmEnabled")
	if vdmEnabledStr != "ENABLED" && vdmEnabledStr != "DISABLED" {
		return nil, ErrBadRequest
	}
	vdm := &sesv2store.VdmAttributes{
		VdmEnabled:                      vdmEnabledStr == "ENABLED",
		AdditionalContactEmailAddresses: request.GetStringList(vdmAttrs, "AdditionalContactEmailAddresses"),
	}

	if dashboardAttrs := request.GetMapParam(vdmAttrs, "DashboardAttributes"); dashboardAttrs != nil {
		vdm.DashboardAttributes = request.GetStringParam(dashboardAttrs, "EngagementMetrics")
	}

	if guardianAttrs := request.GetMapParam(vdmAttrs, "GuardianAttributes"); guardianAttrs != nil {
		vdm.GuardianAttributes = request.GetStringParam(guardianAttrs, "OptimizedSharedDelivery")
	}

	if err := store.PutVdmAttributes(vdm); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// PutAccountDedicatedIpWarmupAttributes updates the dedicated IP auto warmup setting.
func (s *SESv2Service) PutAccountDedicatedIpWarmupAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	enabled := request.GetBoolParam(req.Parameters, "AutoWarmupEnabled")
	if err := store.PutDedicatedIpAutoWarmupEnabled(enabled); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
