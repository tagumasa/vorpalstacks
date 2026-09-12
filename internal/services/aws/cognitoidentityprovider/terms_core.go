package cognitoidentityprovider

import (
	"errors"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	storecommon "vorpalstacks/internal/store/aws/common"

	"github.com/google/uuid"
)

// CreateTermsInput carries the wire parameters of CreateTerms. Params holds
// the raw request parameter map; the nested Links map is read from it inside
// the Core.
type CreateTermsInput struct {
	UserPoolID  string
	ClientID    string
	TermsName   string
	TermsSource string
	Enforcement string
	Params      map[string]interface{}
}

// DescribeTermsInput carries the wire parameters of DescribeTerms.
type DescribeTermsInput struct {
	UserPoolID string
	TermsID    string
}

// DescribeTermsByClientInput carries the wire parameters of
// DescribeTermsByClient.
type DescribeTermsByClientInput struct {
	UserPoolID string
	ClientID   string
	TermsName  string
}

// ListTermsInput carries the wire parameters of ListTerms. Params holds the
// raw request parameter map for the MaxResults member.
type ListTermsInput struct {
	UserPoolID string
	NextToken  string
	Params     map[string]interface{}
}

// UpdateTermsInput carries the wire parameters of UpdateTerms. Params holds
// the raw request parameter map for the Links member.
type UpdateTermsInput struct {
	UserPoolID  string
	TermsID     string
	TermsName   string
	TermsSource string
	Enforcement string
	Params      map[string]interface{}
}

// DeleteTermsInput carries the wire parameters of DeleteTerms.
type DeleteTermsInput struct {
	UserPoolID string
	TermsID    string
}

// validTermsSources is the set of Smithy TermsSourceType enum values.
var validTermsSources = map[string]bool{
	"LINK": true,
}

// validTermsEnforcements is the set of Smithy TermsEnforcementType enum
// values.
var validTermsEnforcements = map[string]bool{
	"NONE": true,
}

// findTermsByName returns the terms document an app client of the user pool
// holds under the given name, walking all pages of the pool's documents.
// Terms document names are unique to the app client, so at most one document
// can match.
func findTermsByName(store cognitostore.CognitoStoreInterface, userPoolID, clientID, termsName string) (*cognitostore.Terms, error) {
	marker := ""
	for {
		result, err := store.ListTermsPaginated(userPoolID, storecommon.ListOptions{MaxItems: listLimitMax, Marker: marker})
		if err != nil {
			return nil, err
		}
		for _, t := range result.Items {
			if t.ClientID == clientID && t.TermsName == termsName {
				return t, nil
			}
		}
		if !result.IsTruncated || result.NextMarker == "" {
			return nil, nil
		}
		marker = result.NextMarker
	}
}

// termsNameExists reports whether an app client of the user pool already
// holds a terms document with the given name, excluding the document with
// the given terms ID (empty excludes nothing).
func termsNameExists(store cognitostore.CognitoStoreInterface, userPoolID, clientID, termsName, excludeTermsID string) (bool, error) {
	t, err := findTermsByName(store, userPoolID, clientID, termsName)
	if err != nil {
		return false, err
	}
	return t != nil && t.TermsID != excludeTermsID, nil
}

// createTermsCore creates a terms document.
func (s *CognitoService) createTermsCore(reqCtx *request.RequestContext, in CreateTermsInput) (interface{}, error) {
	// CreateTermsRequest marks UserPoolId, ClientId, TermsName, TermsSource
	// and Enforcement required.
	if in.UserPoolID == "" || in.ClientID == "" || in.TermsName == "" ||
		in.TermsSource == "" || in.Enforcement == "" {
		return nil, ErrInvalidParameter
	}
	if !validateTermsName(in.TermsName) {
		return nil, ErrInvalidParameter
	}
	if _, ok := validTermsSources[in.TermsSource]; !ok {
		return nil, ErrInvalidParameter
	}
	if _, ok := validTermsEnforcements[in.Enforcement]; !ok {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetUserPool(in.UserPoolID); err != nil {
		return nil, ErrResourceNotFound
	}
	// The terms document's client must reference an existing app client:
	// the model reports an unknown ClientId as ResourceNotFoundException.
	if _, err := store.GetUserPoolClient(in.UserPoolID, in.ClientID); err != nil {
		if errors.Is(err, cognitostore.ErrClientNotFound) {
			return nil, ErrResourceNotFound
		}
		return nil, ErrInternalError
	}

	// Terms document names are unique to the app client; a storage failure
	// in the duplicate check fails closed instead of reading as absence.
	dup, err := termsNameExists(store, in.UserPoolID, in.ClientID, in.TermsName, "")
	if err != nil {
		return nil, ErrInternalError
	}
	if dup {
		return nil, ErrTermsExists
	}

	// TermsIdType is a v4 UUID pattern; the minted identifier must be a
	// value AWS itself could have issued.
	termsID := uuid.NewString()
	t := &cognitostore.Terms{
		TermsID:     termsID,
		UserPoolID:  in.UserPoolID,
		ClientID:    in.ClientID,
		TermsName:   in.TermsName,
		TermsSource: in.TermsSource,
		Enforcement: in.Enforcement,
	}
	if links, ok := in.Params["Links"].(map[string]interface{}); ok {
		t.Links = links
	}

	if err := store.SaveTerms(t); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{"Terms": formatTerms(t)}, nil
}

// describeTermsCore describes a terms document.
func (s *CognitoService) describeTermsCore(reqCtx *request.RequestContext, in DescribeTermsInput) (interface{}, error) {
	if in.UserPoolID == "" || in.TermsID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	t, err := store.GetTerms(in.UserPoolID, in.TermsID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	return map[string]interface{}{"Terms": formatTerms(t)}, nil
}

// describeTermsByClientCore describes the terms document an app client holds
// under a terms name — the client-keyed counterpart of describeTermsCore's
// terms-ID lookup.
func (s *CognitoService) describeTermsByClientCore(reqCtx *request.RequestContext, in DescribeTermsByClientInput) (interface{}, error) {
	// DescribeTermsByClientRequest marks ClientId, UserPoolId and TermsName
	// required.
	if in.UserPoolID == "" || in.ClientID == "" || in.TermsName == "" {
		return nil, ErrInvalidParameter
	}
	if !validateTermsName(in.TermsName) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	t, err := findTermsByName(store, in.UserPoolID, in.ClientID, in.TermsName)
	if err != nil {
		return nil, ErrInternalError
	}
	if t == nil {
		return nil, ErrResourceNotFound
	}

	return map[string]interface{}{"Terms": formatTerms(t)}, nil
}

// listTermsCore lists terms documents for a user pool.
func (s *CognitoService) listTermsCore(reqCtx *request.RequestContext, in ListTermsInput) (interface{}, error) {
	if in.UserPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Smithy ListTermsRequestMaxResultsInteger: range {min: 1, max: 60}
	maxResults, err := parseStrictListLimit(in.Params, "MaxResults", listLimitMax)
	if err != nil {
		return nil, err
	}

	result, err := store.ListTermsPaginated(in.UserPoolID, storecommon.ListOptions{
		MaxItems: maxResults,
		Marker:   in.NextToken,
	})
	if err != nil {
		return nil, ErrInternalError
	}

	formatted := make([]map[string]interface{}, 0, len(result.Items))
	for _, t := range result.Items {
		formatted = append(formatted, formatTermsDescription(t))
	}

	resp := map[string]interface{}{"Terms": formatted}
	if result.IsTruncated && result.NextMarker != "" {
		resp["NextToken"] = result.NextMarker
	}
	return resp, nil
}

// updateTermsCore updates a terms document. The update members are optional:
// an omitted member leaves the stored value untouched.
func (s *CognitoService) updateTermsCore(reqCtx *request.RequestContext, in UpdateTermsInput) (interface{}, error) {
	if in.UserPoolID == "" || in.TermsID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	t, err := store.GetTerms(in.UserPoolID, in.TermsID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if in.TermsName != "" {
		if !validateTermsName(in.TermsName) {
			return nil, ErrInvalidParameter
		}
		// Renaming onto a name another terms document of the same app
		// client already holds is a duplicate; a storage failure in the
		// check fails closed.
		dup, err := termsNameExists(store, in.UserPoolID, t.ClientID, in.TermsName, in.TermsID)
		if err != nil {
			return nil, ErrInternalError
		}
		if dup {
			return nil, ErrTermsExists
		}
		t.TermsName = in.TermsName
	}
	if in.TermsSource != "" {
		if _, ok := validTermsSources[in.TermsSource]; !ok {
			return nil, ErrInvalidParameter
		}
		t.TermsSource = in.TermsSource
	}
	if in.Enforcement != "" {
		if _, ok := validTermsEnforcements[in.Enforcement]; !ok {
			return nil, ErrInvalidParameter
		}
		t.Enforcement = in.Enforcement
	}
	if links, ok := in.Params["Links"].(map[string]interface{}); ok {
		t.Links = links
	}

	if err := store.SaveTerms(t); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{"Terms": formatTerms(t)}, nil
}

// deleteTermsCore deletes a terms document. The operation's error surface
// reports a missing document as ResourceNotFoundException, so the delete is
// guarded by an existence check instead of silently succeeding on an absent
// record.
func (s *CognitoService) deleteTermsCore(reqCtx *request.RequestContext, in DeleteTermsInput) (interface{}, error) {
	if in.UserPoolID == "" || in.TermsID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetTerms(in.UserPoolID, in.TermsID); err != nil {
		return nil, ErrResourceNotFound
	}

	if err := store.DeleteTerms(in.UserPoolID, in.TermsID); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}
