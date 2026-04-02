package sesv2

import (
	"context"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	"vorpalstacks/internal/services/aws/common/tags"
	"vorpalstacks/internal/store/aws/common"
	sesv2store "vorpalstacks/internal/store/aws/sesv2"
)

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
		if err := store.TagResourceFromSlice(arn, parsedTags); err != nil {
			return nil, err
		}
	}

	response := map[string]interface{}{
		"IdentityType":             created.IdentityType,
		"VerifiedForSendingStatus": created.VerifiedForSending,
	}

	if created.DkimAttributes != nil {
		response["DkimAttributes"] = map[string]interface{}{
			"SigningEnabled":          created.DkimAttributes.SigningEnabled,
			"Status":                  created.DkimAttributes.Status,
			"Tokens":                  created.DkimAttributes.Tokens,
			"CurrentSigningKeyLength": created.DkimAttributes.CurrentSigningKeyLength,
			"SigningAttributesOrigin": created.DkimAttributes.SigningAttributesOrigin,
		}
	}

	return response, nil
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

	response := map[string]interface{}{
		"IdentityType":             identity.IdentityType,
		"VerifiedForSendingStatus": identity.VerifiedForSending,
		"FeedbackForwardingStatus": identity.FeedbackForwarding,
	}

	if identity.DkimAttributes != nil {
		response["DkimAttributes"] = map[string]interface{}{
			"SigningEnabled":          identity.DkimAttributes.SigningEnabled,
			"Status":                  identity.DkimAttributes.Status,
			"Tokens":                  identity.DkimAttributes.Tokens,
			"CurrentSigningKeyLength": identity.DkimAttributes.CurrentSigningKeyLength,
			"SigningAttributesOrigin": identity.DkimAttributes.SigningAttributesOrigin,
		}
	}

	if identity.ConfigurationSetName != "" {
		response["ConfigurationSetName"] = identity.ConfigurationSetName
	}

	if identity.MailFromAttributes != nil {
		response["MailFromAttributes"] = map[string]interface{}{
			"MailFromDomain":       identity.MailFromAttributes.MailFromDomain,
			"BehaviorOnMxFailure":  identity.MailFromAttributes.BehaviorOnMxFailure,
			"MailFromDomainStatus": identity.MailFromAttributes.MailFromDomainStatus,
		}
	}

	return response, nil
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
		item := map[string]interface{}{
			"IdentityType":       identity.IdentityType,
			"IdentityName":       identity.Identity,
			"SendingEnabled":     identity.VerifiedForSending,
			"VerificationStatus": "SUCCESS",
		}
		identities = append(identities, item)
	}

	response := map[string]interface{}{
		"EmailIdentities": identities,
	}

	if result.IsTruncated {
		response["NextToken"] = result.NextMarker
	}

	return response, nil
}

// PutEmailIdentityDkimAttributes updates the DKIM attributes for an email identity.
func (s *SESv2Service) PutEmailIdentityDkimAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	emailIdentity := request.GetStringParam(req.Parameters, "EmailIdentity")
	if emailIdentity == "" {
		return nil, ErrMissingParameter
	}

	signingEnabled := request.GetBoolParam(req.Parameters, "SigningEnabled")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	identity, err := store.GetEmailIdentity(emailIdentity)
	if err != nil {
		return nil, err
	}

	if identity.DkimAttributes == nil {
		identity.DkimAttributes = sesv2store.GenerateDkimAttributes(identity.IdentityType)
	}
	identity.DkimAttributes.SigningEnabled = signingEnabled

	if err := store.UpdateEmailIdentity(identity); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// PutEmailIdentityDkimSigningAttributes updates the DKIM signing attributes for an email identity.
func (s *SESv2Service) PutEmailIdentityDkimSigningAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	if identity.DkimAttributes == nil {
		identity.DkimAttributes = sesv2store.GenerateDkimAttributes(identity.IdentityType)
	}

	if err := store.UpdateEmailIdentity(identity); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// PutEmailIdentityFeedbackAttributes updates the feedback attributes for an email identity.
func (s *SESv2Service) PutEmailIdentityFeedbackAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	emailIdentity := request.GetStringParam(req.Parameters, "EmailIdentity")
	if emailIdentity == "" {
		return nil, ErrMissingParameter
	}

	feedbackForwardingEnabled := request.GetBoolParam(req.Parameters, "EmailFeedbackForwardingEnabled")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	identity, err := store.GetEmailIdentity(emailIdentity)
	if err != nil {
		return nil, err
	}

	identity.FeedbackForwarding = feedbackForwardingEnabled

	if err := store.UpdateEmailIdentity(identity); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// PutEmailIdentityMailFromAttributes updates the MAIL FROM attributes for an email identity.
func (s *SESv2Service) PutEmailIdentityMailFromAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	emailIdentity := request.GetStringParam(req.Parameters, "EmailIdentity")
	if emailIdentity == "" {
		return nil, ErrMissingParameter
	}

	mailFromDomain := request.GetStringParam(req.Parameters, "MailFromDomain")
	behaviorOnMxFailure := request.GetStringParam(req.Parameters, "BehaviorOnMxFailure")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	identity, err := store.GetEmailIdentity(emailIdentity)
	if err != nil {
		return nil, err
	}

	if identity.MailFromAttributes == nil {
		identity.MailFromAttributes = &sesv2store.MailFromAttributes{}
	}

	if mailFromDomain != "" {
		identity.MailFromAttributes.MailFromDomain = mailFromDomain
	}
	if behaviorOnMxFailure != "" {
		identity.MailFromAttributes.BehaviorOnMxFailure = behaviorOnMxFailure
	}
	identity.MailFromAttributes.MailFromDomainStatus = "SUCCESS"

	if err := store.UpdateEmailIdentity(identity); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// PutEmailIdentityConfigurationSetAttributes associates a configuration set with an email identity.
func (s *SESv2Service) PutEmailIdentityConfigurationSetAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	emailIdentity := request.GetStringParam(req.Parameters, "EmailIdentity")
	if emailIdentity == "" {
		return nil, ErrMissingParameter
	}

	configurationSetName := request.GetStringParam(req.Parameters, "ConfigurationSetName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	identity, err := store.GetEmailIdentity(emailIdentity)
	if err != nil {
		return nil, err
	}

	identity.ConfigurationSetName = configurationSetName

	if err := store.UpdateEmailIdentity(identity); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
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

// UpdateEmailIdentityPolicy updates a sending authorisation policy for an email identity.
func (s *SESv2Service) UpdateEmailIdentityPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	_ = store.DeleteEmailIdentityPolicy(emailIdentity, policyName)

	return response.EmptyResponse(), nil
}
