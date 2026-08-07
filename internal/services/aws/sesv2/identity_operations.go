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

// updateEmailIdentity is a common helper for PutEmailIdentity* operations.
// It retrieves the identity, applies the modifier, and persists.
func (s *SESv2Service) updateEmailIdentity(reqCtx *request.RequestContext, req *request.ParsedRequest, modify func(*sesv2store.EmailIdentity, map[string]interface{}) error) (interface{}, error) {
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

	if err := modify(identity, req.Parameters); err != nil {
		return nil, err
	}

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
// includes Tags and Policies alongside the identity attributes. The
// previous implementation omitted both, forcing clients to make separate
// ListTagsForResource and GetEmailIdentityPolicies round-trips.
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

	verificationStatus := "PENDING"
	if identity.VerifiedForSending {
		verificationStatus = "SUCCESS"
	}

	resp := map[string]interface{}{
		"IdentityType":             identity.IdentityType,
		"VerifiedForSendingStatus": identity.VerifiedForSending,
		"FeedbackForwardingStatus": identity.FeedbackForwarding,
		"VerificationStatus":       verificationStatus,
	}

	// VerificationInfo carries the last-success timestamp.  In our
	// edge/on-prem implementation every identity is auto-verified at
	// creation time, so LastSuccessTimestamp mirrors CreatedTimestamp.
	if identity.VerifiedForSending {
		resp["VerificationInfo"] = map[string]interface{}{
			"LastSuccessTimestamp": float64(identity.CreatedTimestamp.Unix()),
		}
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

	arn := store.BuildIdentityArn(emailIdentity)
	if tags, err := store.ListAsSlice(arn); err == nil && len(tags) > 0 {
		resp["Tags"] = tags
	}

	policies, err := store.ListEmailIdentityPolicies(emailIdentity)
	if err == nil && len(policies) > 0 {
		policyMap := make(map[string]string, len(policies))
		for _, p := range policies {
			policyMap[p.PolicyName] = p.Policy
		}
		resp["Policies"] = policyMap
	}

	return resp, nil
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
	return s.updateEmailIdentity(reqCtx, req, func(id *sesv2store.EmailIdentity, params map[string]interface{}) error {
		ensureDkimAttributes(id)
		id.DkimAttributes.SigningEnabled = request.GetBoolParam(params, "SigningEnabled")
		return nil
	})
}

// PutEmailIdentityDkimSigningAttributes updates the DKIM signing attributes for an email identity.
// Per Smithy com.amazonaws.sesv2#PutEmailIdentityDkimSigningAttributesRequest
// the BYODKIM selector/private key live inside the nested SigningAttributes
// structure (not at the top level). The top-level SigningAttributesOrigin
// member selects the method (AWS_SES / EXTERNAL). The previous
// implementation passed req.Parameters straight to parseDkimSigningAttributes,
// so BYODKIM via this operation was completely non-functional.
//
// Per Smithy com.amazonaws.sesv2#PutEmailIdentityDkimSigningAttributesResponse
// the response carries DkimStatus, DkimTokens, and SigningHostedZone so
// the caller can publish the DNS records needed to complete verification.
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

	ensureDkimAttributes(identity)
	signingAttrs := request.GetMapParam(req.Parameters, "SigningAttributes")
	byo := parseDkimSigningAttributes(signingAttrs)

	// SigningAttributesOrigin is Smithy @required.
	origin := request.GetStringParam(req.Parameters, "SigningAttributesOrigin")
	if origin == "" && byo == nil {
		return nil, ErrMissingParameter
	}
	if origin != "" && !validateDkimSigningAttributesOrigin(origin) {
		return nil, ErrBadRequest
	}

	if byo == nil {
		// Caller is switching back to AWS-managed DKIM (or only set
		// the top-level SigningAttributesOrigin parameter on its own).
		if origin != "" {
			identity.DkimAttributes.SigningAttributesOrigin = origin
		}
	} else {
		applyDkimSigningAttributes(identity.DkimAttributes, byo)
		if origin != "" {
			identity.DkimAttributes.SigningAttributesOrigin = origin
		}
	}

	if err := store.UpdateEmailIdentity(identity); err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"DkimStatus": identity.DkimAttributes.Status,
		"DkimTokens": identity.DkimAttributes.Tokens,
	}
	if identity.DkimAttributes.SigningHostedZone != "" {
		resp["SigningHostedZone"] = identity.DkimAttributes.SigningHostedZone
	}
	return resp, nil
}

// PutEmailIdentityFeedbackAttributes updates the feedback attributes for an email identity.
func (s *SESv2Service) PutEmailIdentityFeedbackAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.updateEmailIdentity(reqCtx, req, func(id *sesv2store.EmailIdentity, params map[string]interface{}) error {
		id.FeedbackForwarding = request.GetBoolParam(params, "EmailForwardingEnabled")
		return nil
	})
}

// PutEmailIdentityMailFromAttributes updates the MAIL FROM attributes for an email identity.
// Per AWS, MailFromDomainStatus transitions PENDING -> SUCCESS only after
// the required MX record is published. The previous impl hard-coded
// SUCCESS immediately; we now record PENDING so the response shape is
// AWS-compatible (callers that poll for SUCCESS will continue to see
// PENDING until they re-Put or AWS verifies, which on this platform is
// immediate in practice but the wire value must be PENDING initially).
func (s *SESv2Service) PutEmailIdentityMailFromAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.updateEmailIdentity(reqCtx, req, func(id *sesv2store.EmailIdentity, params map[string]interface{}) error {
		if id.MailFromAttributes == nil {
			id.MailFromAttributes = &sesv2store.MailFromAttributes{}
		}
		mailFromDomain := request.GetStringParam(params, "MailFromDomain")
		if mailFromDomain == "" {
			id.MailFromAttributes = &sesv2store.MailFromAttributes{}
			return nil
		}
		id.MailFromAttributes.MailFromDomain = mailFromDomain
		if v := request.GetStringParam(params, "BehaviorOnMxFailure"); v != "" {
			if !validateBehaviorOnMxFailure(v) {
				return ErrBadRequest
			}
			id.MailFromAttributes.BehaviorOnMxFailure = v
		}
		id.MailFromAttributes.MailFromDomainStatus = "PENDING"
		return nil
	})
}

// PutEmailIdentityConfigurationSetAttributes associates a configuration set with an email identity.
func (s *SESv2Service) PutEmailIdentityConfigurationSetAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// Check EmailIdentity presence before fetching the store
	// so a malformed request does not waste a store read on
	// ConfigurationSetExists.
	emailIdentity := request.GetStringParam(req.Parameters, "EmailIdentity")
	if emailIdentity == "" {
		return nil, ErrMissingParameter
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// Verify configuration-set existence when a name is supplied.
	configSetName := request.GetStringParam(req.Parameters, "ConfigurationSetName")
	if configSetName != "" && !store.ConfigurationSetExists(configSetName) {
		return nil, ErrConfigurationSetNotFound
	}
	return s.updateEmailIdentity(reqCtx, req, func(id *sesv2store.EmailIdentity, params map[string]interface{}) error {
		id.ConfigurationSetName = request.GetStringParam(params, "ConfigurationSetName")
		return nil
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
	return s.putEmailIdentityPolicy(reqCtx, req, true)
}

// UpdateEmailIdentityPolicy updates a sending authorisation policy for an email identity.
func (s *SESv2Service) UpdateEmailIdentityPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.putEmailIdentityPolicy(reqCtx, req, false)
}

func (s *SESv2Service) putEmailIdentityPolicy(reqCtx *request.RequestContext, req *request.ParsedRequest, isCreate bool) (interface{}, error) {
	emailIdentity := request.GetStringParam(req.Parameters, "EmailIdentity")
	policyName := request.GetStringParam(req.Parameters, "PolicyName")
	policy := request.GetStringParam(req.Parameters, "Policy")

	if emailIdentity == "" || policyName == "" {
		return nil, ErrMissingParameter
	}
	// Policy is Smithy @required.
	if policy == "" {
		return nil, ErrMissingParameter
	}
	// Validate that the policy is well-formed JSON.
	if err := validatePolicyJSON(policy); err != nil {
		return nil, ErrBadRequest
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	_, err = store.GetEmailIdentity(emailIdentity)
	if err != nil {
		return nil, err
	}

	// Distinguish Create (AlreadyExists) from Update (NotFound).
	existingPolicies, err := store.ListEmailIdentityPolicies(emailIdentity)
	if err != nil {
		return nil, err
	}
	policyExists := false
	for _, p := range existingPolicies {
		if p.PolicyName == policyName {
			policyExists = true
			break
		}
	}
	if isCreate && policyExists {
		return nil, NewAlreadyExistsException("policy")
	}
	if !isCreate && !policyExists {
		return nil, ErrNotFound
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

	// Check that the policy exists before deleting — AWS returns NotFound.
	existingPolicies, err := store.ListEmailIdentityPolicies(emailIdentity)
	if err != nil {
		return nil, err
	}
	policyExists := false
	for _, p := range existingPolicies {
		if p.PolicyName == policyName {
			policyExists = true
			break
		}
	}
	if !policyExists {
		return nil, ErrNotFound
	}

	if err := store.DeleteEmailIdentityPolicy(emailIdentity, policyName); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
