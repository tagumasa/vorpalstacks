package sesv2

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// GetAccount retrieves the details of the current SES v2 account.
func (s *SESv2Service) GetAccount(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.getAccountCore(ctx, store)
}

// PutAccountSendingAttributes updates the sending attributes for the SES v2 account.
func (s *SESv2Service) PutAccountSendingAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.putAccountSendingAttributesCore(store, PutAccountSendingAttributesInput{
		SendingEnabled: request.GetBoolParam(req.Parameters, "SendingEnabled"),
	}); err != nil {
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
	if err := s.putAccountSuppressionAttributesCore(store, PutAccountSuppressionAttributesInput{
		SuppressedReasons:    request.GetStringList(req.Parameters, "SuppressedReasons"),
		ValidationAttributes: request.GetMapParam(req.Parameters, "ValidationAttributes"),
	}); err != nil {
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
	in := PutAccountDetailsInput{
		MailType:                        request.GetStringParam(req.Parameters, "MailType"),
		WebsiteURL:                      request.GetStringParam(req.Parameters, "WebsiteURL"),
		UseCaseDescription:              request.GetStringParam(req.Parameters, "UseCaseDescription"),
		ContactLanguage:                 request.GetStringParam(req.Parameters, "ContactLanguage"),
		AdditionalContactEmailAddresses: request.GetStringList(req.Parameters, "AdditionalContactEmailAddresses"),
	}
	if _, ok := req.Parameters["ProductionAccessEnabled"]; ok {
		in.ProductionAccessEnabled = request.GetBoolParam(req.Parameters, "ProductionAccessEnabled")
		in.ProductionAccessProvided = true
	}
	if err := s.putAccountDetailsCore(store, in); err != nil {
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
	if err := s.putAccountVdmAttributesCore(store, PutAccountVdmAttributesInput{
		VdmAttributes: request.GetMapParam(req.Parameters, "VdmAttributes"),
	}); err != nil {
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
	if err := s.putAccountDedicatedIpWarmupAttributesCore(store, PutAccountDedicatedIpWarmupAttributesInput{
		AutoWarmupEnabled: request.GetBoolParam(req.Parameters, "AutoWarmupEnabled"),
	}); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}
