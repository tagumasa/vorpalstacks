package sesv2

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	sesv2store "vorpalstacks/internal/store/aws/sesv2"
)

// dkimAttributesToMap converts DkimAttributes to a response map.
func dkimAttributesToMap(dkim *sesv2store.DkimAttributes) map[string]interface{} {
	if dkim == nil {
		return nil
	}
	m := map[string]interface{}{
		"SigningEnabled":          dkim.SigningEnabled,
		"Status":                  dkim.Status,
		"Tokens":                  dkim.Tokens,
		"CurrentSigningKeyLength": dkim.CurrentSigningKeyLength,
		"SigningAttributesOrigin": dkim.SigningAttributesOrigin,
	}
	if dkim.NextSigningKeyLength != "" {
		m["NextSigningKeyLength"] = dkim.NextSigningKeyLength
	}
	if dkim.LastKeyGenerationTimestamp != "" {
		m["LastKeyGenerationTimestamp"] = dkim.LastKeyGenerationTimestamp
	}
	if dkim.SigningHostedZone != "" {
		m["SigningHostedZone"] = dkim.SigningHostedZone
	}
	return m
}

// parseDkimSigningAttributes extracts BYODKIM signing attributes from a
// request map. Per Smithy com.amazonaws.sesv2#DkimSigningAttributes the
// shape carries DomainSigningSelector + DomainSigningPrivateKey (and
// optional NextSigningKeyLength + DomainSigningAttributesOrigin). A nil
// map or a map without any of those keys returns nil to indicate that
// AWS-managed DKIM should be used instead.
func parseDkimSigningAttributes(params map[string]interface{}) *sesv2store.DkimSigningAttributes {
	if params == nil {
		return nil
	}
	selector := request.GetStringParam(params, "DomainSigningSelector")
	privateKey := request.GetStringParam(params, "DomainSigningPrivateKey")
	nextLen := request.GetStringParam(params, "NextSigningKeyLength")
	origin := request.GetStringParam(params, "DomainSigningAttributesOrigin")
	if selector == "" && privateKey == "" && nextLen == "" && origin == "" {
		return nil
	}
	return &sesv2store.DkimSigningAttributes{
		DomainSigningSelector:         selector,
		DomainSigningPrivateKey:       privateKey,
		NextSigningKeyLength:          nextLen,
		DomainSigningAttributesOrigin: origin,
	}
}

// applyDkimSigningAttributes writes the BYODKIM caller-supplied fields
// onto an existing DkimAttributes. The caller is responsible for
// initialising DkimAttributes (typically via ensureDkimAttributes) before
// invoking this helper.
func applyDkimSigningAttributes(d *sesv2store.DkimAttributes, byo *sesv2store.DkimSigningAttributes) {
	if byo == nil {
		return
	}
	if byo.DomainSigningSelector != "" {
		d.DomainSigningSelector = byo.DomainSigningSelector
	}
	if byo.DomainSigningPrivateKey != "" {
		d.DomainSigningPrivateKey = byo.DomainSigningPrivateKey
	}
	if byo.NextSigningKeyLength != "" {
		d.NextSigningKeyLength = byo.NextSigningKeyLength
	}
	if byo.DomainSigningAttributesOrigin != "" {
		d.SigningAttributesOrigin = byo.DomainSigningAttributesOrigin
	} else {
		// When the caller supplies BYODKIM without an explicit origin the
		// service records EXTERNAL to indicate the keys are not AWS-managed.
		d.SigningAttributesOrigin = "EXTERNAL"
	}
}

// CreateEmailIdentity creates a new email identity for sending email.
// Per Smithy com.amazonaws.sesv2#CreateEmailIdentityRequest the caller may
// supply DkimSigningAttributes for BYODKIM (Bring Your Own DKIM). When
// present, the identity's DkimAttributes record the caller-supplied
// selector and private key and SigningAttributesOrigin becomes EXTERNAL;
// otherwise AWS-managed DKIM tokens are minted.
// Validation, DKIM generation, configuration-set existence checking, and
// atomic tag persistence are performed by createEmailIdentityCore.
func (s *SESv2Service) CreateEmailIdentity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	dkimMap := request.GetMapParam(req.Parameters, "DkimSigningAttributes")
	result, err := s.createEmailIdentityCore(ctx, store, CreateEmailIdentityInput{
		EmailIdentity:        request.GetStringParam(req.Parameters, "EmailIdentity"),
		Tags:                 tags.ParseTags(req.Parameters, "Tags"),
		ConfigurationSetName: request.GetStringParam(req.Parameters, "ConfigurationSetName"),
		DkimSigningAttrs:     dkimMap,
		DkimSigningProvided:  dkimMap != nil,
	})
	if err != nil {
		return nil, err
	}

	return identityResultToResponse(result), nil
}

// GetEmailIdentity retrieves the details of an email identity.
// Per Smithy `com.amazonaws.sesv2#GetEmailIdentityResponse`, the response
// includes Tags and Policies alongside the identity attributes.
func (s *SESv2Service) GetEmailIdentity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.getEmailIdentityCore(store, request.GetStringParam(req.Parameters, "EmailIdentity"))
}

// DeleteEmailIdentity deletes an email identity.
func (s *SESv2Service) DeleteEmailIdentity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteEmailIdentityCore(store, request.GetStringParam(req.Parameters, "EmailIdentity")); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListEmailIdentities returns a list of email identities for the account.
func (s *SESv2Service) ListEmailIdentities(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.listEmailIdentitiesCore(store, ListEmailIdentitiesInput{
		NextToken: request.GetStringParam(req.Parameters, "NextToken"),
		MaxItems:  request.GetIntParam(req.Parameters, "PageSize"),
	})
	if err != nil {
		return nil, err
	}

	return listEmailIdentitiesResultToResponse(result), nil
}

// PutEmailIdentityDkimAttributes updates the DKIM attributes for an email identity.
func (s *SESv2Service) PutEmailIdentityDkimAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.putEmailIdentityDkimAttributesCore(store, PutEmailIdentityDkimAttributesInput{
		EmailIdentity:  request.GetStringParam(req.Parameters, "EmailIdentity"),
		SigningEnabled: request.GetBoolParam(req.Parameters, "SigningEnabled"),
	}); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// PutEmailIdentityDkimSigningAttributes updates the DKIM signing attributes for an email identity.
// Per Smithy com.amazonaws.sesv2#PutEmailIdentityDkimSigningAttributesRequest
// the BYODKIM selector/private key live inside the nested SigningAttributes
// structure (not at the top level). The top-level SigningAttributesOrigin
// member selects the method (AWS_SES / EXTERNAL).
//
// Per Smithy com.amazonaws.sesv2#PutEmailIdentityDkimSigningAttributesResponse
// the response carries DkimStatus, DkimTokens, and SigningHostedZone so
// the caller can publish the DNS records needed to complete verification.
func (s *SESv2Service) PutEmailIdentityDkimSigningAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.putEmailIdentityDkimSigningAttributesCore(store, PutEmailIdentityDkimSigningAttributesInput{
		EmailIdentity:           request.GetStringParam(req.Parameters, "EmailIdentity"),
		SigningAttributes:       request.GetMapParam(req.Parameters, "SigningAttributes"),
		SigningAttributesOrigin: request.GetStringParam(req.Parameters, "SigningAttributesOrigin"),
	})
}

// PutEmailIdentityFeedbackAttributes updates the feedback attributes for an email identity.
func (s *SESv2Service) PutEmailIdentityFeedbackAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.putEmailIdentityFeedbackAttributesCore(store, PutEmailIdentityFeedbackAttributesInput{
		EmailIdentity:          request.GetStringParam(req.Parameters, "EmailIdentity"),
		EmailForwardingEnabled: request.GetBoolParam(req.Parameters, "EmailForwardingEnabled"),
	}); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// PutEmailIdentityMailFromAttributes updates the MAIL FROM attributes for an email identity.
// Per AWS, MailFromDomainStatus transitions PENDING -> SUCCESS only after
// the required MX record is published; the initial wire value must be
// PENDING.
func (s *SESv2Service) PutEmailIdentityMailFromAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.putEmailIdentityMailFromAttributesCore(store, PutEmailIdentityMailFromAttributesInput{
		EmailIdentity:       request.GetStringParam(req.Parameters, "EmailIdentity"),
		MailFromDomain:      request.GetStringParam(req.Parameters, "MailFromDomain"),
		BehaviorOnMxFailure: request.GetStringParam(req.Parameters, "BehaviorOnMxFailure"),
	}); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// PutEmailIdentityConfigurationSetAttributes associates a configuration set with an email identity.
func (s *SESv2Service) PutEmailIdentityConfigurationSetAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.putEmailIdentityConfigurationSetAttributesCore(store, PutEmailIdentityConfigurationSetAttributesInput{
		EmailIdentity:        request.GetStringParam(req.Parameters, "EmailIdentity"),
		ConfigurationSetName: request.GetStringParam(req.Parameters, "ConfigurationSetName"),
	}); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// GetEmailIdentityPolicies retrieves the policies for an email identity.
func (s *SESv2Service) GetEmailIdentityPolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.getEmailIdentityPoliciesCore(store, request.GetStringParam(req.Parameters, "EmailIdentity"))
}

// CreateEmailIdentityPolicy creates a sending authorisation policy for an email identity.
func (s *SESv2Service) CreateEmailIdentityPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.putEmailIdentityPolicyCore(store, PutEmailIdentityPolicyInput{
		EmailIdentity: request.GetStringParam(req.Parameters, "EmailIdentity"),
		PolicyName:    request.GetStringParam(req.Parameters, "PolicyName"),
		Policy:        request.GetStringParam(req.Parameters, "Policy"),
		IsCreate:      true,
	}); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// UpdateEmailIdentityPolicy updates a sending authorisation policy for an email identity.
func (s *SESv2Service) UpdateEmailIdentityPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.putEmailIdentityPolicyCore(store, PutEmailIdentityPolicyInput{
		EmailIdentity: request.GetStringParam(req.Parameters, "EmailIdentity"),
		PolicyName:    request.GetStringParam(req.Parameters, "PolicyName"),
		Policy:        request.GetStringParam(req.Parameters, "Policy"),
		IsCreate:      false,
	}); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// DeleteEmailIdentityPolicy deletes a sending authorisation policy for an email identity.
func (s *SESv2Service) DeleteEmailIdentityPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteEmailIdentityPolicyCore(store,
		request.GetStringParam(req.Parameters, "EmailIdentity"),
		request.GetStringParam(req.Parameters, "PolicyName")); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}
