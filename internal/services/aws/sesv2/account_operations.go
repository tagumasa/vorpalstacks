package sesv2

import (
	"context"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
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
		result["Details"] = map[string]interface{}{
			"MailType":                        account.Details.MailType,
			"AdditionalContactEmailAddresses": account.Details.AdditionalContactEmailAddresses,
			"UseCaseDescription":              account.Details.UseCaseDescription,
			"WebsiteURL":                      account.Details.WebsiteURL,
			"ContactLanguage":                 account.Details.ContactLanguage,
		}
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
func (s *SESv2Service) PutAccountSuppressionAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	reasons := request.GetStringList(req.Parameters, "SuppressedReasons")
	if err := store.PutSuppressionAttributes(reasons); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// PutAccountDetails updates the account details for the SES v2 account.
func (s *SESv2Service) PutAccountDetails(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	details := &sesv2store.AccountDetails{
		MailType:                request.GetStringParam(req.Parameters, "MailType"),
		UseCaseDescription:      request.GetStringParam(req.Parameters, "UseCaseDescription"),
		WebsiteURL:              request.GetStringParam(req.Parameters, "WebsiteUrl"),
		ContactLanguage:         request.GetStringParam(req.Parameters, "ContactLanguage"),
		ProductionAccessEnabled: request.GetBoolParam(req.Parameters, "ProductionAccessEnabled"),
	}

	additionalEmails := request.GetStringList(req.Parameters, "AdditionalContactEmailAddresses")
	if len(additionalEmails) > 0 {
		details.AdditionalContactEmailAddresses = additionalEmails
	}

	if err := store.PutAccountDetails(details); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// PutAccountVdmAttributes updates the VDM attributes for the SES v2 account.
func (s *SESv2Service) PutAccountVdmAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	vdm := &sesv2store.VdmAttributes{
		VdmEnabled:          request.GetBoolParam(req.Parameters, "VdmEnabled"),
		DashboardAttributes: request.GetStringParam(req.Parameters, "DashboardAttributes"),
		GuardianAttributes:  request.GetStringParam(req.Parameters, "GuardianAttributes"),
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
