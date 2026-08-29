package sesv2

import (
	"context"
	"strings"

	types "vorpalstacks/internal/common/tags"
	storecommon "vorpalstacks/internal/store/aws/common"
	sesv2store "vorpalstacks/internal/store/aws/sesv2"
)

// ensureDkimAttributes initialises DkimAttributes if nil.
func ensureDkimAttributes(identity *sesv2store.EmailIdentity) {
	if identity.DkimAttributes == nil {
		identity.DkimAttributes = sesv2store.GenerateDkimAttributes(identity.IdentityType)
	}
}

// ---------------------------------------------------------------------------
// Input DTOs — identity attribute / policy family
// ---------------------------------------------------------------------------

// PutEmailIdentityDkimAttributesInput carries the DKIM signing toggle.
type PutEmailIdentityDkimAttributesInput struct {
	EmailIdentity  string
	SigningEnabled bool
}

// PutEmailIdentityDkimSigningAttributesInput carries the BYODKIM update
// members; SigningAttributes travels as the raw nested wire map.
type PutEmailIdentityDkimSigningAttributesInput struct {
	EmailIdentity           string
	SigningAttributes       map[string]interface{}
	SigningAttributesOrigin string
}

// PutEmailIdentityFeedbackAttributesInput carries the feedback forwarding
// toggle.
type PutEmailIdentityFeedbackAttributesInput struct {
	EmailIdentity          string
	EmailForwardingEnabled bool
}

// PutEmailIdentityMailFromAttributesInput carries the MAIL FROM members.
type PutEmailIdentityMailFromAttributesInput struct {
	EmailIdentity       string
	MailFromDomain      string
	BehaviorOnMxFailure string
}

// PutEmailIdentityConfigurationSetAttributesInput carries the
// configuration-set association member.
type PutEmailIdentityConfigurationSetAttributesInput struct {
	EmailIdentity        string
	ConfigurationSetName string
}

// PutEmailIdentityPolicyInput carries the policy document members; IsCreate
// distinguishes the create (AlreadyExists) from the update (NotFound)
// semantics.
type PutEmailIdentityPolicyInput struct {
	EmailIdentity string
	PolicyName    string
	Policy        string
	IsCreate      bool
}

// ---------------------------------------------------------------------------
// Core functions — identity attribute / policy family
// ---------------------------------------------------------------------------

// updateEmailIdentityCore is the shared skeleton for the
// PutEmailIdentity* family: it resolves the identity, applies the
// modifier, and persists.
func (s *SESv2Service) updateEmailIdentityCore(store sesv2store.SESv2StoreInterface, emailIdentity string, modify func(*sesv2store.EmailIdentity) error) error {
	if emailIdentity == "" {
		return ErrMissingParameter
	}

	identity, err := store.GetEmailIdentity(emailIdentity)
	if err != nil {
		return err
	}

	if err := modify(identity); err != nil {
		return err
	}

	return store.UpdateEmailIdentity(identity)
}

// getEmailIdentityCore is the single entry point for reading an email
// identity, including its tags and policies.
func (s *SESv2Service) getEmailIdentityCore(store sesv2store.SESv2StoreInterface, emailIdentity string) (map[string]interface{}, error) {
	if emailIdentity == "" {
		return nil, ErrMissingParameter
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

// putEmailIdentityDkimAttributesCore is the single entry point for the
// DKIM signing toggle.
func (s *SESv2Service) putEmailIdentityDkimAttributesCore(store sesv2store.SESv2StoreInterface, in PutEmailIdentityDkimAttributesInput) error {
	return s.updateEmailIdentityCore(store, in.EmailIdentity, func(id *sesv2store.EmailIdentity) error {
		ensureDkimAttributes(id)
		id.DkimAttributes.SigningEnabled = in.SigningEnabled
		return nil
	})
}

// putEmailIdentityDkimSigningAttributesCore is the single entry point for
// the BYODKIM signing-attributes update. Per Smithy
// com.amazonaws.sesv2#PutEmailIdentityDkimSigningAttributesRequest the
// BYODKIM selector/private key live inside the nested SigningAttributes
// structure (not at the top level). The top-level SigningAttributesOrigin
// member selects the method (AWS_SES / EXTERNAL) and is Smithy @required.
//
// Per Smithy com.amazonaws.sesv2#PutEmailIdentityDkimSigningAttributesResponse
// the response carries DkimStatus, DkimTokens, and SigningHostedZone so
// the caller can publish the DNS records needed to complete verification.
func (s *SESv2Service) putEmailIdentityDkimSigningAttributesCore(store sesv2store.SESv2StoreInterface, in PutEmailIdentityDkimSigningAttributesInput) (map[string]interface{}, error) {
	emailIdentity := in.EmailIdentity
	if emailIdentity == "" {
		return nil, ErrMissingParameter
	}

	identity, err := store.GetEmailIdentity(emailIdentity)
	if err != nil {
		return nil, err
	}

	ensureDkimAttributes(identity)
	byo := parseDkimSigningAttributes(in.SigningAttributes)

	// SigningAttributesOrigin is Smithy @required.
	origin := in.SigningAttributesOrigin
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

// putEmailIdentityFeedbackAttributesCore is the single entry point for the
// feedback forwarding toggle.
func (s *SESv2Service) putEmailIdentityFeedbackAttributesCore(store sesv2store.SESv2StoreInterface, in PutEmailIdentityFeedbackAttributesInput) error {
	return s.updateEmailIdentityCore(store, in.EmailIdentity, func(id *sesv2store.EmailIdentity) error {
		id.FeedbackForwarding = in.EmailForwardingEnabled
		return nil
	})
}

// putEmailIdentityMailFromAttributesCore is the single entry point for the
// MAIL FROM attributes. Per AWS, MailFromDomainStatus transitions
// PENDING -> SUCCESS only after the required MX record is published, so
// the wire value must be PENDING initially.
func (s *SESv2Service) putEmailIdentityMailFromAttributesCore(store sesv2store.SESv2StoreInterface, in PutEmailIdentityMailFromAttributesInput) error {
	return s.updateEmailIdentityCore(store, in.EmailIdentity, func(id *sesv2store.EmailIdentity) error {
		if id.MailFromAttributes == nil {
			id.MailFromAttributes = &sesv2store.MailFromAttributes{}
		}
		if in.MailFromDomain == "" {
			id.MailFromAttributes = &sesv2store.MailFromAttributes{}
			return nil
		}
		id.MailFromAttributes.MailFromDomain = in.MailFromDomain
		if in.BehaviorOnMxFailure != "" {
			if !validateBehaviorOnMxFailure(in.BehaviorOnMxFailure) {
				return ErrBadRequest
			}
			id.MailFromAttributes.BehaviorOnMxFailure = in.BehaviorOnMxFailure
		}
		id.MailFromAttributes.MailFromDomainStatus = "PENDING"
		return nil
	})
}

// putEmailIdentityConfigurationSetAttributesCore is the single entry point
// for associating a configuration set with an email identity.
func (s *SESv2Service) putEmailIdentityConfigurationSetAttributesCore(store sesv2store.SESv2StoreInterface, in PutEmailIdentityConfigurationSetAttributesInput) error {
	// Check EmailIdentity presence before the store reads so a malformed
	// request does not waste a store read on ConfigurationSetExists.
	if in.EmailIdentity == "" {
		return ErrMissingParameter
	}
	// Verify configuration-set existence when a name is supplied.
	if in.ConfigurationSetName != "" && !store.ConfigurationSetExists(in.ConfigurationSetName) {
		return ErrConfigurationSetNotFound
	}
	return s.updateEmailIdentityCore(store, in.EmailIdentity, func(id *sesv2store.EmailIdentity) error {
		id.ConfigurationSetName = in.ConfigurationSetName
		return nil
	})
}

// getEmailIdentityPoliciesCore is the single entry point for reading an
// identity's sending-authorization policies.
func (s *SESv2Service) getEmailIdentityPoliciesCore(store sesv2store.SESv2StoreInterface, emailIdentity string) (map[string]interface{}, error) {
	if emailIdentity == "" {
		return nil, ErrMissingParameter
	}

	if _, err := store.GetEmailIdentity(emailIdentity); err != nil {
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

// putEmailIdentityPolicyCore is the single entry point shared by
// CreateEmailIdentityPolicy and UpdateEmailIdentityPolicy, distinguishing
// the Create (AlreadyExists) from the Update (NotFound) semantics.
func (s *SESv2Service) putEmailIdentityPolicyCore(store sesv2store.SESv2StoreInterface, in PutEmailIdentityPolicyInput) error {
	emailIdentity := in.EmailIdentity
	policyName := in.PolicyName
	policy := in.Policy

	if emailIdentity == "" || policyName == "" {
		return ErrMissingParameter
	}
	// Policy is Smithy @required.
	if policy == "" {
		return ErrMissingParameter
	}
	// Validate that the policy is well-formed JSON.
	if err := validatePolicyJSON(policy); err != nil {
		return ErrBadRequest
	}

	if _, err := store.GetEmailIdentity(emailIdentity); err != nil {
		return err
	}

	// Distinguish Create (AlreadyExists) from Update (NotFound).
	existingPolicies, err := store.ListEmailIdentityPolicies(emailIdentity)
	if err != nil {
		return err
	}
	policyExists := false
	for _, p := range existingPolicies {
		if p.PolicyName == policyName {
			policyExists = true
			break
		}
	}
	if in.IsCreate && policyExists {
		return NewAlreadyExistsException("policy")
	}
	if !in.IsCreate && !policyExists {
		return ErrNotFound
	}

	return store.PutEmailIdentityPolicy(emailIdentity, policyName, policy)
}

// deleteEmailIdentityPolicyCore is the single entry point for deleting a
// sending-authorization policy. AWS returns NotFound for an unknown policy.
func (s *SESv2Service) deleteEmailIdentityPolicyCore(store sesv2store.SESv2StoreInterface, emailIdentity, policyName string) error {
	if emailIdentity == "" || policyName == "" {
		return ErrMissingParameter
	}

	if _, err := store.GetEmailIdentity(emailIdentity); err != nil {
		return err
	}

	// Check that the policy exists before deleting — AWS returns NotFound.
	existingPolicies, err := store.ListEmailIdentityPolicies(emailIdentity)
	if err != nil {
		return err
	}
	policyExists := false
	for _, p := range existingPolicies {
		if p.PolicyName == policyName {
			policyExists = true
			break
		}
	}
	if !policyExists {
		return ErrNotFound
	}

	return store.DeleteEmailIdentityPolicy(emailIdentity, policyName)
}

// ---------------------------------------------------------------------------
// Transport-agnostic Input / Result structs
// ---------------------------------------------------------------------------

// CreateEmailIdentityInput carries every field that CreateEmailIdentity
// needs, in a format independent of the wire protocol (HTTP Query/JSON vs
// gRPC-Web). Both the HTTP API handler (identity_operations.go) and the
// admin gRPC handler (admin_handler.go) build this struct from their
// respective request formats and delegate to createEmailIdentityCore,
// ensuring that identity validation, DKIM generation, configuration-set
// existence checks, and tag persistence follow a single code path.
type CreateEmailIdentityInput struct {
	EmailIdentity        string
	Tags                 []types.Tag
	ConfigurationSetName string
	DkimSigningAttrs     map[string]interface{}
	// DkimSigningProvided is true when the DkimSigningAttributes key was
	// present in the request (empty-map detection).
	DkimSigningProvided bool
}

// IdentityResult is the transport-agnostic result of creating or looking
// up an email identity.
type IdentityResult struct {
	IdentityType       string
	VerifiedForSending bool
	DkimAttributes     map[string]interface{}
}

// ListEmailIdentitiesInput carries the pagination parameters.
type ListEmailIdentitiesInput struct {
	NextToken string
	MaxItems  int
}

// IdentitySummary is the transport-agnostic summary of a single identity
// in list results.
type IdentitySummary struct {
	IdentityType       string
	IdentityName       string
	SendingEnabled     bool
	VerificationStatus string
}

// ListEmailIdentitiesResult is the transport-agnostic result of listing
// email identities.
type ListEmailIdentitiesResult struct {
	Identities []IdentitySummary
	NextToken  string
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// createEmailIdentityCore is the single entry point for email-identity
// creation logic shared by the HTTP API and the admin gRPC handler. It
// performs identity-format validation, configuration-set existence
// checking, BYODKIM attribute application, and atomic tag
// persistence (rolls back the identity on tag failure).
func (s *SESv2Service) createEmailIdentityCore(ctx context.Context, store sesv2store.SESv2StoreInterface, in CreateEmailIdentityInput) (*IdentityResult, error) {
	if in.EmailIdentity == "" {
		return nil, ErrMissingParameter
	}
	if !validateIdentityFormat(in.EmailIdentity) {
		return nil, ErrBadRequest
	}

	// AWS SES treats identity names as case-insensitive.
	// Normalise to lowercase so that lookups in identityExistsForEmail
	// match regardless of the case used at creation time.
	identityName := strings.ToLower(in.EmailIdentity)

	// Verify configuration-set existence when a name is supplied.
	if in.ConfigurationSetName != "" {
		if !store.ConfigurationSetExists(in.ConfigurationSetName) {
			return nil, ErrConfigurationSetNotFound
		}
	}

	identity := sesv2store.NewEmailIdentity(identityName)
	if in.ConfigurationSetName != "" {
		identity.ConfigurationSetName = in.ConfigurationSetName
	}

	if in.DkimSigningProvided {
		byoDkim := parseDkimSigningAttributes(in.DkimSigningAttrs)
		if byoDkim != nil {
			ensureDkimAttributes(identity)
			applyDkimSigningAttributes(identity.DkimAttributes, byoDkim)
		}
	}

	created, err := store.CreateEmailIdentity(identity)
	if err != nil {
		return nil, err
	}

	// Apply tags atomically — roll back the identity if tagging fails.
	if len(in.Tags) > 0 {
		arn := store.BuildIdentityArn(identityName)
		if err := store.TagFromSlice(arn, in.Tags); err != nil {
			_ = store.DeleteEmailIdentity(identityName)
			return nil, err
		}
	}

	return &IdentityResult{
		IdentityType:       created.IdentityType,
		VerifiedForSending: created.VerifiedForSending,
		DkimAttributes:     dkimAttributesToMap(created.DkimAttributes),
	}, nil
}

// deleteEmailIdentityCore is the single entry point for email-identity
// deletion shared by the HTTP API and the admin gRPC handler.
func (s *SESv2Service) deleteEmailIdentityCore(store sesv2store.SESv2StoreInterface, emailIdentity string) error {
	if emailIdentity == "" {
		return ErrMissingParameter
	}
	return store.DeleteEmailIdentity(emailIdentity)
}

// listEmailIdentitiesCore is the single entry point for email-identity
// listing shared by the HTTP API and the admin gRPC handler.
func (s *SESv2Service) listEmailIdentitiesCore(store sesv2store.SESv2StoreInterface, in ListEmailIdentitiesInput) (*ListEmailIdentitiesResult, error) {
	maxItems := in.MaxItems
	if maxItems <= 0 {
		maxItems = 100
	}
	if maxItems > 100 {
		maxItems = 100
	}

	result, err := store.ListEmailIdentities(storecommon.ListOptions{
		MaxItems: maxItems,
		Marker:   in.NextToken,
	})
	if err != nil {
		return nil, err
	}

	summaries := make([]IdentitySummary, 0, len(result.Items))
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
		summaries = append(summaries, IdentitySummary{
			IdentityType:       identity.IdentityType,
			IdentityName:       identity.Identity,
			SendingEnabled:     identity.VerifiedForSending,
			VerificationStatus: verificationStatus,
		})
	}

	nextToken := ""
	if result.IsTruncated {
		nextToken = result.NextMarker
	}

	return &ListEmailIdentitiesResult{
		Identities: summaries,
		NextToken:  nextToken,
	}, nil
}

// identityResultToResponse converts an IdentityResult to the HTTP API
// response map.
func identityResultToResponse(r *IdentityResult) map[string]interface{} {
	resp := map[string]interface{}{
		"IdentityType":             r.IdentityType,
		"VerifiedForSendingStatus": r.VerifiedForSending,
	}
	if r.DkimAttributes != nil {
		resp["DkimAttributes"] = r.DkimAttributes
	}
	return resp
}

// listEmailIdentitiesResultToResponse converts a ListEmailIdentitiesResult
// to the HTTP API response map.
func listEmailIdentitiesResultToResponse(r *ListEmailIdentitiesResult) map[string]interface{} {
	items := make([]map[string]interface{}, len(r.Identities))
	for i, s := range r.Identities {
		items[i] = map[string]interface{}{
			"IdentityType":       s.IdentityType,
			"IdentityName":       s.IdentityName,
			"SendingEnabled":     s.SendingEnabled,
			"VerificationStatus": s.VerificationStatus,
		}
	}
	resp := map[string]interface{}{
		"EmailIdentities": items,
	}
	if r.NextToken != "" {
		resp["NextToken"] = r.NextToken
	}
	return resp
}
