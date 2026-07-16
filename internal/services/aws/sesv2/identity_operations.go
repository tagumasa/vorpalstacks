package sesv2

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/store/aws/common"
	sesv2store "vorpalstacks/internal/store/aws/sesv2"
)

// dkimAttributesToMap converts DkimAttributes to a response map.
func dkimAttributesToMap(dkim *sesv2store.DkimAttributes) map[string]interface{} {
	if dkim == nil {
		return nil
	}
	return map[string]interface{}{
		"SigningEnabled":          dkim.SigningEnabled,
		"Status":                  dkim.Status,
		"Tokens":                  dkim.Tokens,
		"CurrentSigningKeyLength": dkim.CurrentSigningKeyLength,
		"SigningAttributesOrigin": dkim.SigningAttributesOrigin,
	}
}

// updateEmailIdentity is a common helper for PutEmailIdentity* operations.
// It retrieves the identity, applies the modifier, and persists.
func (s *SESv2Service) updateEmailIdentity(reqCtx *request.RequestContext, req *request.ParsedRequest, modify func(*sesv2store.EmailIdentity, map[string]interface{})) (interface{}, error) {
	emailIdentity := request.GetStringParam(req.Parameters, "EmailIdentity")
	if emailIdentity == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	identity, err := store.GetEmailIdentity(emailIdentity)
	if err != nil {
		return nil, err
	}

	modify(identity, req.Parameters)

	if err := store.UpdateEmailIdentity(identity); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ensureDkimAttributes initialises DkimAttributes if nil.
func ensureDkimAttributes(identity *sesv2store.EmailIdentity) {
	if identity.DkimAttributes == nil {
		identity.DkimAttributes = sesv2store.GenerateDkimAttributes(identity.IdentityType)
	}
}

// CreateEmailIdentity creates a new email identity for sending email.
func (s *SESv2Service) CreateEmailIdentity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	emailIdentity := request.GetStringParam(req.Parameters, "EmailIdentity")
	if emailIdentity == "" {
		return nil, ErrMissingParameter
	}

	parsedTags := tags.ParseTags(req.Parameters, "Tags")
	configurationSetName := request.GetStringParam(req.Parameters, "ConfigurationSetName")

	identity := sesv2store.NewEmailIdentity(emailIdentity)
	if configurationSetName != "" {
		identity.ConfigurationSetName = configurationSetName
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	created, err := store.CreateEmailIdentity(identity)
	if err != nil {
		return nil, err
	}

	if len(parsedTags) > 0 {
		arn := store.BuildIdentityArn(emailIdentity)
		if err := store.TagFromSlice(arn, parsedTags); err != nil {
			return nil, err
		}
	}

	resp := map[string]interface{}{
		"IdentityType":             created.IdentityType,
		"VerifiedForSendingStatus": created.VerifiedForSending,
	}

	if dkim := dkimAttributesToMap(created.DkimAttributes); dkim != nil {
		resp["DkimAttributes"] = dkim
	}

	return resp, nil
}

// GetEmailIdentity retrieves the details of an email identity.
func (s *SESv2Service) GetEmailIdentity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	emailIdentity := request.GetStringParam(req.Parameters, "EmailIdentity")
	if emailIdentity == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	identity, err := store.GetEmailIdentity(emailIdentity)
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"IdentityType":             identity.IdentityType,
		"VerifiedForSendingStatus": identity.VerifiedForSending,
		"FeedbackForwardingStatus": identity.FeedbackForwarding,
	}

	if dkim := dkimAttributesToMap(identity.DkimAttributes); dkim != nil {
		resp["DkimAttributes"] = dkim
	}

	if identity.ConfigurationSetName != "" {
		resp["ConfigurationSetName"] = identity.ConfigurationSetName
	}

	if identity.MailFromAttributes != nil {
		resp["MailFromAttributes"] = map[string]interface{}{
			"MailFromDomain":       identity.MailFromAttributes.MailFromDomain,
			"BehaviorOnMxFailure":  identity.MailFromAttributes.BehaviorOnMxFailure,
			"MailFromDomainStatus": identity.MailFromAttributes.MailFromDomainStatus,
		}
	}

	return resp, nil
}

// DeleteEmailIdentity deletes an email identity.
func (s *SESv2Service) DeleteEmailIdentity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	emailIdentity := request.GetStringParam(req.Parameters, "EmailIdentity")
	if emailIdentity == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteEmailIdentity(emailIdentity); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListEmailIdentities returns a list of email identities for the account.
func (s *SESv2Service) ListEmailIdentities(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	pageSize := request.GetIntParam(req.Parameters, "PageSize")
	if pageSize == 0 {
		pageSize = 100
	}
	nextToken := request.GetStringParam(req.Parameters, "NextToken")

	opts := common.ListOptions{
		MaxItems: pageSize,
		Marker:   nextToken,
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := store.ListEmailIdentities(opts)
	if err != nil {
		return nil, err
	}

	identities := make([]map[string]interface{}, 0, len(result.Items))
	for _, identity := range result.Items {
		verificationStatus := "PENDING"
		if identity.DkimAttributes != nil {
			verificationStatus = identity.DkimAttributes.Status
			if verificationStatus == "" {
				verificationStatus = "SUCCESS"
			}
		}
		if identity.VerifiedForSending {
			verificationStatus = "SUCCESS"
		}
		item := map[string]interface{}{
			"IdentityType":       identity.IdentityType,
			"IdentityName":       identity.Identity,
			"SendingEnabled":     identity.VerifiedForSending,
			"VerificationStatus": verificationStatus,
		}
		identities = append(identities, item)
	}

	resp := map[string]interface{}{
		"EmailIdentities": identities,
	}

	if result.IsTruncated {
		resp["NextToken"] = result.NextMarker
	}

	return resp, nil
}

// PutEmailIdentityDkimAttributes updates the DKIM attributes for an email identity.
func (s *SESv2Service) PutEmailIdentityDkimAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.updateEmailIdentity(reqCtx, req, func(id *sesv2store.EmailIdentity, params map[string]interface{}) {
		ensureDkimAttributes(id)
		id.DkimAttributes.SigningEnabled = request.GetBoolParam(params, "SigningEnabled")
	})
}

// PutEmailIdentityDkimSigningAttributes updates the DKIM signing attributes for an email identity.
func (s *SESv2Service) PutEmailIdentityDkimSigningAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.updateEmailIdentity(reqCtx, req, func(id *sesv2store.EmailIdentity, params map[string]interface{}) {
		ensureDkimAttributes(id)
		if v := request.GetStringParam(params, "SigningAttributesOrigin"); v != "" {
			id.DkimAttributes.SigningAttributesOrigin = v
		}
	})
}

// PutEmailIdentityFeedbackAttributes updates the feedback attributes for an email identity.
func (s *SESv2Service) PutEmailIdentityFeedbackAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.updateEmailIdentity(reqCtx, req, func(id *sesv2store.EmailIdentity, params map[string]interface{}) {
		id.FeedbackForwarding = request.GetBoolParam(params, "EmailForwardingEnabled")
	})
}

// PutEmailIdentityMailFromAttributes updates the MAIL FROM attributes for an email identity.
func (s *SESv2Service) PutEmailIdentityMailFromAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.updateEmailIdentity(reqCtx, req, func(id *sesv2store.EmailIdentity, params map[string]interface{}) {
		if id.MailFromAttributes == nil {
			id.MailFromAttributes = &sesv2store.MailFromAttributes{}
		}
		if v := request.GetStringParam(params, "MailFromDomain"); v != "" {
			id.MailFromAttributes.MailFromDomain = v
		}
		if v := request.GetStringParam(params, "BehaviorOnMxFailure"); v != "" {
			id.MailFromAttributes.BehaviorOnMxFailure = v
		}
		id.MailFromAttributes.MailFromDomainStatus = "SUCCESS"
	})
}

// PutEmailIdentityConfigurationSetAttributes associates a configuration set with an email identity.
func (s *SESv2Service) PutEmailIdentityConfigurationSetAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.updateEmailIdentity(reqCtx, req, func(id *sesv2store.EmailIdentity, params map[string]interface{}) {
		id.ConfigurationSetName = request.GetStringParam(params, "ConfigurationSetName")
	})
}

// GetEmailIdentityPolicies retrieves the policies for an email identity.
func (s *SESv2Service) GetEmailIdentityPolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	emailIdentity := request.GetStringParam(req.Parameters, "EmailIdentity")
	if emailIdentity == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	_, err = store.GetEmailIdentity(emailIdentity)
	if err != nil {
		return nil, err
	}

	policies, err := store.ListEmailIdentityPolicies(emailIdentity)
	if err != nil {
		return nil, err
	}

	policyMap := make(map[string]string)
	for _, p := range policies {
		policyMap[p.PolicyName] = p.Policy
	}

	return map[string]interface{}{
		"Policies": policyMap,
	}, nil
}

// CreateEmailIdentityPolicy creates a sending authorisation policy for an email identity.
func (s *SESv2Service) CreateEmailIdentityPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.putEmailIdentityPolicy(reqCtx, req)
}

// UpdateEmailIdentityPolicy updates a sending authorisation policy for an email identity.
func (s *SESv2Service) UpdateEmailIdentityPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.putEmailIdentityPolicy(reqCtx, req)
}

func (s *SESv2Service) putEmailIdentityPolicy(reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	emailIdentity := request.GetStringParam(req.Parameters, "EmailIdentity")
	policyName := request.GetStringParam(req.Parameters, "PolicyName")
	policy := request.GetStringParam(req.Parameters, "Policy")

	if emailIdentity == "" || policyName == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	_, err = store.GetEmailIdentity(emailIdentity)
	if err != nil {
		return nil, err
	}

	if err := store.PutEmailIdentityPolicy(emailIdentity, policyName, policy); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DeleteEmailIdentityPolicy deletes a sending authorisation policy for an email identity.
func (s *SESv2Service) DeleteEmailIdentityPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	emailIdentity := request.GetStringParam(req.Parameters, "EmailIdentity")
	policyName := request.GetStringParam(req.Parameters, "PolicyName")

	if emailIdentity == "" || policyName == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	_, err = store.GetEmailIdentity(emailIdentity)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteEmailIdentityPolicy(emailIdentity, policyName); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
