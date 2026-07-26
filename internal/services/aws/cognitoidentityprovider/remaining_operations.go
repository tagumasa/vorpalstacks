package cognitoidentityprovider

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// ===================== Phase 14: WebAuthn =====================

// StartWebAuthnRegistration starts a WebAuthn credential registration flow.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_StartWebAuthnRegistration.html
func (s *CognitoService) StartWebAuthnRegistration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessToken := getAccessToken(req)
	if accessToken == "" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, accessToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	challenge := make([]byte, 32)
	rand.Read(challenge)
	challengeB64 := base64.RawURLEncoding.EncodeToString(challenge)

	// M11: Store the challenge in a session for CompleteWebAuthnRegistration binding
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	user, err := store.GetUserByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	sessionID := generateSessionID()
	challengeSession := &cognitostore.ChallengeSession{
		SessionID:     sessionID,
		UserPoolID:    user.UserPoolID,
		ClientID:      "",
		Username:      user.Username,
		ChallengeName: "WEB_AUTHN_REGISTRATION",
		CreatedAt:     time.Now().UTC(),
		ExpiresAt:     time.Now().UTC().Add(5 * time.Minute),
	}
	if err := store.SaveChallengeSession(challengeSession); err != nil {
		return nil, ErrInternalError
	}

	// M10: Use custom domain as RP ID if configured, otherwise default
	rpID := cognitoIdpHost(s.region)
	if domain, err := store.GetUserPoolDomainByPool(user.UserPoolID); err == nil && domain.Domain != "" {
		rpID = domain.Domain
	}

	options := map[string]interface{}{
		"challenge": challengeB64,
		"rp": map[string]interface{}{
			"name": "Cognito",
			"id":   rpID,
		},
		"user": map[string]interface{}{
			"id":          base64.RawURLEncoding.EncodeToString([]byte(userID)),
			"name":        userID,
			"displayName": userID,
		},
		"pubKeyCredParams": []map[string]interface{}{
			{"type": "public-key", "alg": -7},
			{"type": "public-key", "alg": -257},
		},
		"timeout":     60000,
		"attestation": "none",
		"authenticatorSelection": map[string]interface{}{
			"authenticatorAttachment": "platform",
			"userVerification":        "preferred",
		},
	}

	return map[string]interface{}{
		"CredentialCreationOptions": options,
		"Session":                   sessionID,
	}, nil
}

// CompleteWebAuthnRegistration completes a WebAuthn credential registration.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_CompleteWebAuthnRegistration.html
func (s *CognitoService) CompleteWebAuthnRegistration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessToken := getAccessToken(req)
	if accessToken == "" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, accessToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUserByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	credentialRaw, ok := req.Parameters["Credential"]
	if !ok {
		return nil, ErrInvalidParameter
	}

	credentialBytes, _ := json.Marshal(credentialRaw)
	var credential struct {
		ID        string `json:"id"`
		PublicKey string `json:"publicKey"`
		Type      string `json:"type"`
	}
	json.Unmarshal(credentialBytes, &credential)
	if credential.ID == "" {
		return nil, ErrInvalidParameter
	}

	friendlyName := credential.ID
	if len(friendlyName) > 8 {
		friendlyName = friendlyName[:8]
	}

	cred := &cognitostore.WebAuthnCredential{
		CredentialID: credential.ID,
		FriendlyName: friendlyName,
		UserPoolID:   user.UserPoolID,
		UserID:       user.ID,
		PublicKey:    credential.PublicKey,
		CreatedAt:    time.Now().UTC(),
	}

	if err := store.CreateWebAuthnCredential(cred); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// ListWebAuthnCredentials lists registered WebAuthn credentials.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ListWebAuthnCredentials.html
func (s *CognitoService) ListWebAuthnCredentials(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessToken := getAccessToken(req)
	if accessToken == "" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, accessToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUserByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	maxResults := 20
	if mr := request.GetIntParam(req.Parameters, "MaxResults"); mr > 0 {
		maxResults = mr
	}
	// Smithy WebAuthnCredentialsQueryLimitType: range {min: 0, max: 20}
	if maxResults > 20 {
		return nil, ErrInvalidParameter
	}

	result, err := store.ListWebAuthnCredentialsPaginated(user.UserPoolID, user.ID, storecommon.ListOptions{
		MaxItems: maxResults,
		Marker:   req.GetParam("NextToken"),
	})
	if err != nil {
		return nil, ErrInternalError
	}

	formatted := make([]map[string]interface{}, 0, len(result.Items))
	for _, c := range result.Items {
		formatted = append(formatted, map[string]interface{}{
			"CredentialId":           c.CredentialID,
			"FriendlyCredentialName": c.FriendlyName,
			"CreatedAt":              c.CreatedAt.Unix(),
		})
	}

	resp := map[string]interface{}{"Credentials": formatted}
	if result.IsTruncated && result.NextMarker != "" {
		resp["NextToken"] = result.NextMarker
	}
	return resp, nil
}

// DeleteWebAuthnCredential deletes a WebAuthn credential.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_DeleteWebAuthnCredential.html
func (s *CognitoService) DeleteWebAuthnCredential(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessToken := getAccessToken(req)
	credID := req.GetParam("CredentialId")
	if accessToken == "" || credID == "" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, accessToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUserByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if err := store.DeleteWebAuthnCredential(user.UserPoolID, user.ID, credID); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// ===================== Phase 15: Managed Login Branding =====================

// CreateManagedLoginBranding creates a managed login branding configuration.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_CreateManagedLoginBranding.html
func (s *CognitoService) CreateManagedLoginBranding(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetUserPool(userPoolID); err != nil {
		return nil, ErrResourceNotFound
	}

	brandingID := "branding-" + generateID()
	b := &cognitostore.ManagedLoginBranding{
		ManagedLoginBrandingId: brandingID,
		UserPoolID:             userPoolID,
		ClientID:               req.GetParam("ClientId"),
	}
	if v, ok := req.Parameters["UseCognitoProvidedValues"].(bool); ok {
		b.UseCognitoProvidedValues = v
	}
	if settings, ok := req.Parameters["Settings"]; ok {
		if m, ok := settings.(map[string]interface{}); ok {
			b.Settings = m
		}
	}
	if err := parseBrandingAssets(req, b); err != nil {
		return nil, err
	}

	if err := store.SaveManagedLoginBranding(b); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{"ManagedLoginBranding": formatManagedLoginBranding(b)}, nil
}

// DescribeManagedLoginBranding describes a managed login branding configuration.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_DescribeManagedLoginBranding.html
func (s *CognitoService) DescribeManagedLoginBranding(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	brandingID := req.GetParam("ManagedLoginBrandingId")
	if userPoolID == "" || brandingID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	b, err := store.GetManagedLoginBranding(userPoolID, brandingID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	return map[string]interface{}{"ManagedLoginBranding": formatManagedLoginBranding(b)}, nil
}

// DescribeManagedLoginBrandingByClient describes branding by client ID.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_DescribeManagedLoginBrandingByClient.html
func (s *CognitoService) DescribeManagedLoginBrandingByClient(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	clientID := req.GetParam("ClientId")
	if userPoolID == "" || clientID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	b, err := store.GetManagedLoginBrandingByClient(userPoolID, clientID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	return map[string]interface{}{"ManagedLoginBranding": formatManagedLoginBranding(b)}, nil
}

// UpdateManagedLoginBranding updates a managed login branding configuration.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_UpdateManagedLoginBranding.html
func (s *CognitoService) UpdateManagedLoginBranding(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	brandingID := req.GetParam("ManagedLoginBrandingId")
	if userPoolID == "" || brandingID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	b, err := store.GetManagedLoginBranding(userPoolID, brandingID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if v, ok := req.Parameters["UseCognitoProvidedValues"].(bool); ok {
		b.UseCognitoProvidedValues = v
	}
	if settings, ok := req.Parameters["Settings"]; ok {
		if m, ok := settings.(map[string]interface{}); ok {
			b.Settings = m
		}
	}
	if err := parseBrandingAssets(req, b); err != nil {
		return nil, err
	}

	if err := store.SaveManagedLoginBranding(b); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{"ManagedLoginBranding": formatManagedLoginBranding(b)}, nil
}

// DeleteManagedLoginBranding deletes a managed login branding configuration.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_DeleteManagedLoginBranding.html
func (s *CognitoService) DeleteManagedLoginBranding(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	brandingID := req.GetParam("ManagedLoginBrandingId")
	if userPoolID == "" || brandingID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteManagedLoginBranding(userPoolID, brandingID); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

var validAssetCategories = map[string]bool{
	"AUTH_APP_GRAPHIC": true, "EMAIL_GRAPHIC": true, "FAVICON_ICO": true,
	"FAVICON_SVG": true, "FORM_BACKGROUND": true, "FORM_LOGO": true,
	"IDP_BUTTON_ICON": true, "PAGE_BACKGROUND": true, "PAGE_FOOTER_BACKGROUND": true,
	"PAGE_FOOTER_LOGO": true, "PAGE_HEADER_BACKGROUND": true, "PAGE_HEADER_LOGO": true,
	"PASSKEY_GRAPHIC": true, "PASSWORD_GRAPHIC": true, "SMS_GRAPHIC": true,
}

var validAssetExtensions = map[string]bool{
	"ICO": true, "JPEG": true, "PNG": true, "SVG": true, "WEBP": true,
}

var validColorModes = map[string]bool{
	"LIGHT": true, "DARK": true, "DYNAMIC": true,
}

var assetMagicBytes = map[string][]byte{
	"PNG":  {0x89, 0x50, 0x4E, 0x47},
	"JPEG": {0xFF, 0xD8, 0xFF},
	"ICO":  {0x00, 0x00, 0x01, 0x00},
	"GIF":  {0x47, 0x49, 0x46, 0x38},
}

func parseBrandingAssets(req *request.ParsedRequest, b *cognitostore.ManagedLoginBranding) error {
	if rawAssets, ok := req.Parameters["Assets"].([]interface{}); ok {
		b.Assets = nil
		for _, a := range rawAssets {
			m, ok := a.(map[string]interface{})
			if !ok {
				continue
			}
			category := getStringParam(m, "Category")
			extension := getStringParam(m, "Extension")
			colorMode := getStringParam(m, "Color")
			// L3: Validate enum values
			if category != "" && !validAssetCategories[category] {
				return ErrInvalidParameter
			}
			if extension != "" && !validAssetExtensions[extension] {
				return ErrInvalidParameter
			}
			if colorMode != "" && !validColorModes[colorMode] {
				return ErrInvalidParameter
			}
			// Validate Bytes magic bytes when extension is known
			bytesVal := getStringParam(m, "Bytes")
			if bytesVal != "" && extension != "" {
				if magic, ok := assetMagicBytes[extension]; ok {
					decoded, err := base64.StdEncoding.DecodeString(bytesVal)
					if err == nil && len(decoded) >= len(magic) {
						for i, b := range magic {
							if decoded[i] != b {
								return ErrInvalidParameter
							}
						}
					}
				}
			}
			asset := cognitostore.BrandingAsset{
				Category:  category,
				Color:     colorMode,
				Extension: extension,
				Bytes:     bytesVal,
			}
			b.Assets = append(b.Assets, asset)
		}
	}
	return nil
}

func formatManagedLoginBranding(b *cognitostore.ManagedLoginBranding) map[string]interface{} {
	result := map[string]interface{}{
		"ManagedLoginBrandingId":   b.ManagedLoginBrandingId,
		"UserPoolId":               b.UserPoolID,
		"UseCognitoProvidedValues": b.UseCognitoProvidedValues,
	}
	if b.ClientID != "" {
		result["ClientId"] = b.ClientID
	}
	if b.Settings != nil {
		result["Settings"] = b.Settings
	}
	if len(b.Assets) > 0 {
		assets := make([]map[string]interface{}, 0, len(b.Assets))
		for _, a := range b.Assets {
			assets = append(assets, map[string]interface{}{
				"Category":  a.Category,
				"Color":     a.Color,
				"Extension": a.Extension,
				"Bytes":     a.Bytes,
			})
		}
		result["Assets"] = assets
	}
	if !b.CreationDate.IsZero() {
		result["CreationDate"] = b.CreationDate.Unix()
	}
	if !b.LastModifiedDate.IsZero() {
		result["LastModifiedDate"] = b.LastModifiedDate.Unix()
	}
	return result
}

// ===================== Phase 16: Terms =====================

// CreateTerms creates a terms document.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_CreateTerms.html
func (s *CognitoService) CreateTerms(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetUserPool(userPoolID); err != nil {
		return nil, ErrResourceNotFound
	}

	termsID := "terms-" + generateID()
	t := &cognitostore.Terms{
		TermsID:    termsID,
		UserPoolID: userPoolID,
		ClientID:   req.GetParam("ClientId"),
		TermsName:  req.GetParam("TermsName"),
	}
	if source, ok := req.Parameters["TermsSource"].(map[string]interface{}); ok {
		t.TermsSource = source
	}
	if enforcement, ok := req.Parameters["Enforcement"].(map[string]interface{}); ok {
		t.Enforcement = enforcement
	}
	if links, ok := req.Parameters["Links"].(map[string]interface{}); ok {
		t.Links = links
	}

	if err := store.SaveTerms(t); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{"Terms": formatTerms(t)}, nil
}

// DescribeTerms describes a terms document.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_DescribeTerms.html
func (s *CognitoService) DescribeTerms(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	termsID := req.GetParam("TermsId")
	if userPoolID == "" || termsID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	t, err := store.GetTerms(userPoolID, termsID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	return map[string]interface{}{"Terms": formatTerms(t)}, nil
}

// ListTerms lists terms documents for a user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ListTerms.html
func (s *CognitoService) ListTerms(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Smithy ListTermsRequestMaxResultsInteger: range {min: 1, max: 60}
	maxResults := 60
	if mr := request.GetIntParam(req.Parameters, "MaxResults"); mr > 0 && mr <= 60 {
		maxResults = mr
	}

	result, err := store.ListTermsPaginated(userPoolID, storecommon.ListOptions{
		MaxItems: maxResults,
		Marker:   req.GetParam("NextToken"),
	})
	if err != nil {
		return nil, ErrInternalError
	}

	formatted := make([]map[string]interface{}, 0, len(result.Items))
	for _, t := range result.Items {
		formatted = append(formatted, formatTerms(t))
	}

	resp := map[string]interface{}{"Terms": formatted}
	if result.IsTruncated && result.NextMarker != "" {
		resp["NextToken"] = result.NextMarker
	}
	return resp, nil
}

// UpdateTerms updates a terms document.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_UpdateTerms.html
func (s *CognitoService) UpdateTerms(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	termsID := req.GetParam("TermsId")
	if userPoolID == "" || termsID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	t, err := store.GetTerms(userPoolID, termsID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if name := req.GetParam("TermsName"); name != "" {
		t.TermsName = name
	}
	if source, ok := req.Parameters["TermsSource"].(map[string]interface{}); ok {
		t.TermsSource = source
	}
	if enforcement, ok := req.Parameters["Enforcement"].(map[string]interface{}); ok {
		t.Enforcement = enforcement
	}
	if links, ok := req.Parameters["Links"].(map[string]interface{}); ok {
		t.Links = links
	}

	if err := store.SaveTerms(t); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{"Terms": formatTerms(t)}, nil
}

// DeleteTerms deletes a terms document.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_DeleteTerms.html
func (s *CognitoService) DeleteTerms(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	termsID := req.GetParam("TermsId")
	if userPoolID == "" || termsID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteTerms(userPoolID, termsID); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

func formatTerms(t *cognitostore.Terms) map[string]interface{} {
	result := map[string]interface{}{
		"TermsId":    t.TermsID,
		"UserPoolId": t.UserPoolID,
		"TermsName":  t.TermsName,
	}
	if t.ClientID != "" {
		result["ClientId"] = t.ClientID
	}
	if t.TermsSource != nil {
		result["TermsSource"] = t.TermsSource
	}
	if t.Enforcement != nil {
		result["Enforcement"] = t.Enforcement
	}
	if t.Links != nil {
		result["Links"] = t.Links
	}
	if !t.CreationDate.IsZero() {
		result["CreationDate"] = t.CreationDate.Unix()
	}
	if !t.LastModifiedDate.IsZero() {
		result["LastModifiedDate"] = t.LastModifiedDate.Unix()
	}
	return result
}

// ===================== Phase 17: Replicas =====================

// CreateUserPoolReplica creates a cross-region replica of a user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_CreateUserPoolReplica.html
func (s *CognitoService) CreateUserPoolReplica(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	regionName := req.GetParam("RegionName")
	if userPoolID == "" || regionName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	pool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	replica := &cognitostore.UserPoolReplica{
		UserPoolID:   userPoolID,
		RegionName:   regionName,
		Status:       "Active",
		Role:         "Full",
		UserPoolArn:  "arn:aws:cognito-idp:" + s.region + ":" + s.accountID + ":userpool/" + pool.ID,
		CreationDate: time.Now().UTC(),
	}

	if err := store.SaveUserPoolReplica(replica); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{"UserPoolReplica": formatUserPoolReplica(replica)}, nil
}

// ListUserPoolReplicas lists replicas for a user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ListUserPoolReplicas.html
func (s *CognitoService) ListUserPoolReplicas(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := store.ListUserPoolReplicasPaginated(userPoolID, storecommon.ListOptions{
		MaxItems: 60,
		Marker:   req.GetParam("NextToken"),
	})
	if err != nil {
		return nil, ErrInternalError
	}

	formatted := make([]map[string]interface{}, 0, len(result.Items))
	for _, r := range result.Items {
		formatted = append(formatted, formatUserPoolReplica(r))
	}

	resp := map[string]interface{}{"UserPoolReplicas": formatted}
	if result.IsTruncated && result.NextMarker != "" {
		resp["NextToken"] = result.NextMarker
	}
	return resp, nil
}

// DeleteUserPoolReplica deletes a cross-region replica.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_DeleteUserPoolReplica.html
func (s *CognitoService) DeleteUserPoolReplica(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	regionName := req.GetParam("RegionName")
	if userPoolID == "" || regionName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	replica, err := store.GetUserPoolReplica(userPoolID, regionName)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if err := store.DeleteUserPoolReplica(userPoolID, regionName); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{"UserPoolReplica": formatUserPoolReplica(replica)}, nil
}

// UpdateUserPoolReplica updates a cross-region replica.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_UpdateUserPoolReplica.html
func (s *CognitoService) UpdateUserPoolReplica(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	regionName := req.GetParam("RegionName")
	if userPoolID == "" || regionName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	replica, err := store.GetUserPoolReplica(userPoolID, regionName)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if status := req.GetParam("Status"); status != "" {
		replica.Status = status
	}

	if err := store.SaveUserPoolReplica(replica); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{"UserPoolReplica": formatUserPoolReplica(replica)}, nil
}

func formatUserPoolReplica(r *cognitostore.UserPoolReplica) map[string]interface{} {
	result := map[string]interface{}{
		"RegionName": r.RegionName,
		"Status":     r.Status,
		"Role":       r.Role,
	}
	if r.UserPoolArn != "" {
		result["UserPoolArn"] = r.UserPoolArn
	}
	return result
}
