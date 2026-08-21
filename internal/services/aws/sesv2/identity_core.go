package sesv2

import (
	"context"
	"strings"

	types "vorpalstacks/internal/common/tags"
	storecommon "vorpalstacks/internal/store/aws/common"
	sesv2store "vorpalstacks/internal/store/aws/sesv2"
)

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
