package iot

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"strings"
	"time"

	"github.com/google/uuid"

	storecommon "vorpalstacks/internal/store/aws/common"
	iotstore "vorpalstacks/internal/store/aws/iot"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// ---------------------------------------------------------------------------
// Custom authorizer Core — CRUD plus the Lambda test-invocation path
// ---------------------------------------------------------------------------

// CreateAuthorizerInput carries the fields for CreateAuthorizer. Status is
// applied only when StatusProvided is true; EnableCachingForHTTP only when
// EnableCachingProvided is true (an explicit false must survive the trip).
type CreateAuthorizerInput struct {
	AuthorizerName         string
	AuthorizerFunctionARN  string
	TokenKeyName           string
	TokenSigningPublicKeys map[string]string
	SigningDisabled        bool
	Status                 string
	StatusProvided         bool
	EnableCachingForHTTP   bool
	EnableCachingProvided  bool
}

// UpdateAuthorizerInput carries the fields for UpdateAuthorizer. Empty string
// fields are left untouched; EnableCaching is applied only when non-nil.
type UpdateAuthorizerInput struct {
	AuthorizerName         string
	AuthorizerFunctionARN  string
	TokenKeyName           string
	TokenSigningPublicKeys map[string]string
	Status                 string
	EnableCaching          *bool
}

// ListAuthorizersResult is the transport-agnostic result of ListAuthorizers.
type ListAuthorizersResult struct {
	Authorizers []*iotstore.Authorizer
	NextMarker  string
}

// TestInvokeAuthorizerInput carries the fields for TestInvokeAuthorizer.
type TestInvokeAuthorizerInput struct {
	AuthorizerName string
	Token          string
	TokenSignature string
	MqttContext    map[string]interface{}
	HttpContext    map[string]interface{}
	TlsContext     map[string]interface{}
}

// createAuthorizerCore validates and persists a custom authorizer for MQTT
// connections.
func (s *IoTService) createAuthorizerCore(store iotstore.IotStoreInterface, in CreateAuthorizerInput) (*iotstore.Authorizer, error) {
	if in.AuthorizerName == "" {
		return nil, iotstore.ErrMissingParam
	}

	auth := &iotstore.Authorizer{
		AuthorizerName:         in.AuthorizerName,
		AuthorizerFunctionARN:  in.AuthorizerFunctionARN,
		TokenName:              in.TokenKeyName,
		TokenSigningPublicKeys: in.TokenSigningPublicKeys,
		SigningDisabled:        in.SigningDisabled,
		Status:                 true,
		EnableCachingForHTTP:   true,
		CreationDate:           time.Now().UTC(),
		LastModifiedDate:       time.Now().UTC(),
	}

	if in.StatusProvided {
		if err := ValidateAuthorizerStatus(in.Status); err != nil {
			return nil, err
		}
		auth.Status = in.Status == "ACTIVE"
	}

	if in.EnableCachingProvided {
		auth.EnableCachingForHTTP = in.EnableCachingForHTTP
	}

	return store.CreateAuthorizer(auth)
}

// describeAuthorizerCore retrieves a custom authorizer by name.
func (s *IoTService) describeAuthorizerCore(store iotstore.IotStoreInterface, name string) (*iotstore.Authorizer, error) {
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}
	return store.GetAuthorizer(name)
}

// updateAuthorizerCore applies the supplied fields to an existing custom
// authorizer.
func (s *IoTService) updateAuthorizerCore(store iotstore.IotStoreInterface, in UpdateAuthorizerInput) (*iotstore.Authorizer, error) {
	if in.AuthorizerName == "" {
		return nil, iotstore.ErrMissingParam
	}

	// ParseAttributes always returns a non-nil (but possibly empty) map.
	// Guard with len() so absent keys are not overwritten with an empty map.
	var signingKeys map[string]string
	if len(in.TokenSigningPublicKeys) > 0 {
		signingKeys = in.TokenSigningPublicKeys
	}
	opts := iotstore.AuthorizerUpdateOpts{
		FunctionARN:            in.AuthorizerFunctionARN,
		TokenName:              in.TokenKeyName,
		TokenSigningPublicKeys: signingKeys,
		Status:                 in.Status,
	}
	if opts.Status != "" {
		if err := ValidateAuthorizerStatus(opts.Status); err != nil {
			return nil, err
		}
	}
	if in.EnableCaching != nil {
		opts.EnableCaching = in.EnableCaching
	}

	return store.UpdateAuthorizer(in.AuthorizerName, opts)
}

// deleteAuthorizerCore removes a custom authorizer and its tags.
func (s *IoTService) deleteAuthorizerCore(store iotstore.IotStoreInterface, name string) error {
	if name == "" {
		return iotstore.ErrMissingParam
	}

	arn := iotstore.BuildAuthorizerARN(store.GetAccountID(), store.GetRegion(), name)
	_ = store.DeleteAllTags(arn)

	return store.DeleteAuthorizer(name)
}

// listAuthorizersCore lists custom authorizers with pagination.
func (s *IoTService) listAuthorizersCore(store iotstore.IotStoreInterface, marker string, maxItems int) (*ListAuthorizersResult, error) {
	maxItems = listMaxItems(maxItems)
	auths, err := store.ListAuthorizers(iotstoreListOpts(maxItems, marker))
	if err != nil {
		return nil, err
	}
	return &ListAuthorizersResult{
		Authorizers: auths.Items,
		NextMarker:  auths.NextMarker,
	}, nil
}

// testInvokeAuthorizerCore verifies the supplied token against the
// authorizer's registered public keys, then dispatches the authorizer
// invocation event to its Lambda function.
func (s *IoTService) testInvokeAuthorizerCore(ctx context.Context, store iotstore.IotStoreInterface, in TestInvokeAuthorizerInput) (map[string]interface{}, error) {
	if in.AuthorizerName == "" {
		return nil, iotstore.ErrMissingParam
	}

	auth, err := store.GetAuthorizer(in.AuthorizerName)
	if err != nil {
		return nil, err
	}
	if auth == nil || auth.AuthorizerName == "" {
		return nil, iotstore.ErrAuthorizerNotFound
	}

	// When signing is enabled, the token must be accompanied by a valid
	// signature verified against one of the registered public keys.
	signatureVerified := false
	if !auth.SigningDisabled {
		if in.Token == "" {
			return nil, iotstore.ErrInvalidRequest
		}
		if len(auth.TokenSigningPublicKeys) > 0 && in.TokenSignature != "" {
			sigBytes, sigErr := hex.DecodeString(in.TokenSignature)
			if sigErr == nil {
				for _, pubKeyPEM := range auth.TokenSigningPublicKeys {
					block, _ := pem.Decode([]byte(pubKeyPEM))
					if block == nil {
						continue
					}
					pub, pubErr := x509.ParsePKIXPublicKey(block.Bytes)
					if pubErr != nil {
						continue
					}
					switch pk := pub.(type) {
					case *ecdsa.PublicKey:
						if ecdsa.VerifyASN1(pk, cryptoSHA256([]byte(in.Token)), sigBytes) {
							signatureVerified = true
						}
					case *rsa.PublicKey:
						if rsa.VerifyPKCS1v15(pk, crypto.SHA256, cryptoSHA256([]byte(in.Token)), sigBytes) == nil {
							signatureVerified = true
						}
					}
					if signatureVerified {
						break
					}
				}
			}
		}
	} else {
		signatureVerified = true
	}

	// Build the authorizer invocation event and dispatch to Lambda.
	protocolData := map[string]interface{}{}
	if len(in.MqttContext) > 0 {
		protocolData["mqtt"] = in.MqttContext
	}
	if len(in.HttpContext) > 0 {
		protocolData["http"] = in.HttpContext
	}
	if len(in.TlsContext) > 0 {
		protocolData["tls"] = in.TlsContext
	}

	event := map[string]interface{}{
		"token":             in.Token,
		"signatureVerified": signatureVerified,
		"protocols":         []string{"tls", "http", "mqtt"},
		"protocolData":      protocolData,
		"connectionMetadata": map[string]interface{}{
			"sessionId": uuid.New().String(),
		},
	}
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return nil, iotstore.ErrInternalFailure
	}

	// Resolve the Lambda function name from the ARN.
	fnName := lambdaFunctionNameFromARN(auth.AuthorizerFunctionARN)
	if fnName == "" {
		return nil, iotstore.ErrInvalidRequest
	}

	invoker := s.deps.EventBus.LambdaInvoker()
	if invoker == nil {
		return nil, iotstore.ErrInternalFailure
	}

	_, respBytes, invokeErr := invoker.InvokeForGateway(ctx, fnName, eventJSON)
	if invokeErr != nil {
		return nil, iotstore.ErrInternalFailure
	}

	// Parse the Lambda response into the TestInvokeAuthorizerResponse shape.
	var lambdaResp struct {
		IsAuthenticated       bool     `json:"isAuthenticated"`
		PrincipalID           string   `json:"principalId"`
		PolicyDocuments       []string `json:"policyDocuments"`
		RefreshAfterInSeconds int64    `json:"refreshAfterInSeconds"`
		DisconnectAfterInSecs int64    `json:"disconnectAfterInSeconds"`
	}
	if err := json.Unmarshal(respBytes, &lambdaResp); err != nil {
		return nil, iotstore.ErrInternalFailure
	}

	result := map[string]interface{}{
		"isAuthenticated":          lambdaResp.IsAuthenticated,
		"principalId":              lambdaResp.PrincipalID,
		"refreshAfterInSeconds":    lambdaResp.RefreshAfterInSeconds,
		"disconnectAfterInSeconds": lambdaResp.DisconnectAfterInSecs,
	}
	if len(lambdaResp.PolicyDocuments) > 0 {
		policyDocs := make([]interface{}, 0, len(lambdaResp.PolicyDocuments))
		for _, pd := range lambdaResp.PolicyDocuments {
			policyDocs = append(policyDocs, pd)
		}
		result["policyDocuments"] = policyDocs
	} else {
		result["policyDocuments"] = []interface{}{}
	}
	return result, nil
}

// cryptoSHA256 returns the SHA-256 digest of the given data.
func cryptoSHA256(data []byte) []byte {
	h := sha256.Sum256(data)
	return h[:]
}

// listMaxItems applies the default page size for list cores.
func listMaxItems(maxItems int) int {
	if maxItems <= 0 {
		return storecommon.DefaultMaxItems
	}
	return maxItems
}

// lambdaFunctionNameFromARN extracts the function name (or ARN suffix)
// from a Lambda ARN for internal invocation.
func lambdaFunctionNameFromARN(arn string) string {
	if arn == "" {
		return ""
	}
	_, _, _, _, resource := svcarn.SplitARN(arn)
	segs := strings.SplitN(resource, ":", 3)
	if len(segs) < 2 {
		// Function name without ARN prefix.
		return arn
	}
	return segs[1]
}

// ---------------------------------------------------------------------------
// Default authorizer Core — the account's default-authorizer pointer
// ---------------------------------------------------------------------------

// DefaultAuthorizerResult is the transport-agnostic result of
// SetDefaultAuthorizer.
type DefaultAuthorizerResult struct {
	AuthorizerName string
	AuthorizerARN  string
}

// setDefaultAuthorizerCore validates that the named authorizer exists, then
// records it as the account's default authorizer.
func (s *IoTService) setDefaultAuthorizerCore(store iotstore.IotStoreInterface, authorizerName string) (*DefaultAuthorizerResult, error) {
	if authorizerName == "" {
		return nil, iotstore.ErrMissingParam
	}
	auth, err := store.GetAuthorizer(authorizerName)
	if err != nil {
		return nil, err
	}
	if auth == nil || auth.AuthorizerName == "" {
		return nil, iotstore.ErrAuthorizerNotFound
	}
	if err := store.PutGeneric("config/defaultAuthorizer", authorizerName); err != nil {
		return nil, err
	}
	return &DefaultAuthorizerResult{
		AuthorizerName: authorizerName,
		AuthorizerARN:  auth.AuthorizerARN,
	}, nil
}

// clearDefaultAuthorizerCore removes the default authorizer pointer.
func (s *IoTService) clearDefaultAuthorizerCore(store iotstore.IotStoreInterface) error {
	return store.DeleteGeneric("config/defaultAuthorizer")
}

// describeDefaultAuthorizerCore resolves the configured default authorizer
// record. The response shape is the authorizerDescription structure, the
// same as DescribeAuthorizer's output.
func (s *IoTService) describeDefaultAuthorizerCore(store iotstore.IotStoreInterface) (*iotstore.Authorizer, error) {
	name := ""
	exists, err := store.GetGenericExists("config/defaultAuthorizer", &name)
	if err != nil {
		return nil, err
	}
	if !exists || name == "" {
		return nil, iotstore.ErrDefaultAuthorizerNotFound
	}
	return store.GetAuthorizer(name)
}
