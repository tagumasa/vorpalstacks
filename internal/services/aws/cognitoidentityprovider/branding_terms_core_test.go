package cognitoidentityprovider

import (
	"errors"
	"testing"

	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// newBrandingTermsClient registers one app client on a fresh pool for the
// branding/terms cores.
func newBrandingTermsClient(t *testing.T) (*CognitoService, *request.RequestContext, string, string) {
	t.Helper()
	svc, reqCtx, store, pool := newTokensTestEnv(t)
	if err := store.CreateUserPoolClient(&cognitostore.UserPoolClient{
		UserPoolID: pool.ID,
		ClientID:   "branding-client",
		ClientName: "branding-client",
	}); err != nil {
		t.Fatal(err)
	}
	return svc, reqCtx, pool.ID, "branding-client"
}

// An unknown ClientId is refused with ResourceNotFoundException instead of
// proceeding to create a branding style for a client that does not exist.
func TestCreateManagedLoginBrandingRejectsUnknownClient(t *testing.T) {
	svc, reqCtx, poolID, _ := newBrandingTermsClient(t)

	if _, err := svc.createManagedLoginBrandingCore(reqCtx, CreateManagedLoginBrandingInput{
		UserPoolID: poolID,
		ClientID:   "absent-client",
		Params:     map[string]interface{}{},
	}); !errors.Is(err, ErrResourceNotFound) {
		t.Fatalf("branding for an unknown client returned %v, want ResourceNotFound", err)
	}
}

// A second branding style for one app client is refused, and the store
// guard holds even when the pre-check is bypassed (two direct saves).
func TestCreateManagedLoginBrandingDuplicateRefused(t *testing.T) {
	svc, reqCtx, poolID, clientID := newBrandingTermsClient(t)

	if _, err := svc.createManagedLoginBrandingCore(reqCtx, CreateManagedLoginBrandingInput{
		UserPoolID: poolID,
		ClientID:   clientID,
		Params:     map[string]interface{}{},
	}); err != nil {
		t.Fatalf("first branding create failed: %v", err)
	}
	if _, err := svc.createManagedLoginBrandingCore(reqCtx, CreateManagedLoginBrandingInput{
		UserPoolID: poolID,
		ClientID:   clientID,
		Params:     map[string]interface{}{},
	}); !errors.Is(err, ErrManagedLoginBrandingExists) {
		t.Fatalf("second branding create returned %v, want ManagedLoginBrandingExists", err)
	}
}

// Asset entries that are not maps, or whose Bytes payload does not decode
// to the extension's magic bytes, are rejected instead of silently skipped.
func TestCreateManagedLoginBrandingRejectsMalformedAssets(t *testing.T) {
	svc, reqCtx, poolID, clientID := newBrandingTermsClient(t)

	if _, err := svc.createManagedLoginBrandingCore(reqCtx, CreateManagedLoginBrandingInput{
		UserPoolID: poolID,
		ClientID:   clientID,
		Params: map[string]interface{}{
			"Assets": []interface{}{"not-a-map"},
		},
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("non-map asset returned %v, want InvalidParameter", err)
	}

	if _, err := svc.createManagedLoginBrandingCore(reqCtx, CreateManagedLoginBrandingInput{
		UserPoolID: poolID,
		ClientID:   clientID,
		Params: map[string]interface{}{
			"Assets": []interface{}{map[string]interface{}{
				"Category":  "FORM_LOGO",
				"Extension": "PNG",
				"Bytes":     "not-base64!!!",
			}},
		},
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("undecodable asset bytes returned %v, want InvalidParameter", err)
	}
}

// An unknown ClientId is refused for terms documents as well.
func TestCreateTermsRejectsUnknownClient(t *testing.T) {
	svc, reqCtx, poolID, _ := newBrandingTermsClient(t)

	if _, err := svc.createTermsCore(reqCtx, CreateTermsInput{
		UserPoolID:  poolID,
		ClientID:    "absent-client",
		TermsName:   "terms-of-use",
		TermsSource: "LINK",
		Enforcement: "NONE",
	}); !errors.Is(err, ErrResourceNotFound) {
		t.Fatalf("terms for an unknown client returned %v, want ResourceNotFound", err)
	}
}

// A second terms document of the same name for one app client is refused.
func TestCreateTermsDuplicateNameRefused(t *testing.T) {
	svc, reqCtx, poolID, clientID := newBrandingTermsClient(t)

	in := CreateTermsInput{
		UserPoolID:  poolID,
		ClientID:    clientID,
		TermsName:   "terms-of-use",
		TermsSource: "LINK",
		Enforcement: "NONE",
	}
	if _, err := svc.createTermsCore(reqCtx, in); err != nil {
		t.Fatalf("first terms create failed: %v", err)
	}
	if _, err := svc.createTermsCore(reqCtx, in); !errors.Is(err, ErrTermsExists) {
		t.Fatalf("duplicate terms name returned %v, want TermsExists", err)
	}
}
