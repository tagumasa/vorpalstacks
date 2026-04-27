// Package kms provides KMS (Key Management Service) operations for vorpalstacks.
package kms

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/iam/policy"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/eventbus"
	"vorpalstacks/internal/services/aws/kms/hsm"
	storecommon "vorpalstacks/internal/store/aws/common"
	iamstore "vorpalstacks/internal/store/aws/iam"
	kmsstore "vorpalstacks/internal/store/aws/kms"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// kmsStores holds the various KMS stores.
type kmsStores struct {
	keys        *kmsstore.KeyStore
	aliases     *kmsstore.AliasStore
	grants      *kmsstore.GrantStore
	keyPolicies *kmsstore.KeyPolicyStore
}

// KMSService provides AWS KMS key and encryption operations.
type KMSService struct {
	hsmBackend      hsm.Backend
	policyEvaluator *policy.PolicyEvaluator
	accountID       string
	stores          sync.Map // region → *kmsStores
}

// NewKMSService creates a new KMS service instance.
func NewKMSService(accountID, region string, hsmBackend hsm.Backend) *KMSService {
	return &KMSService{
		hsmBackend:      hsmBackend,
		policyEvaluator: policy.NewPolicyEvaluator(),
		accountID:       accountID,
	}
}

func (s *KMSService) store(reqCtx *request.RequestContext) (*kmsStores, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, reqCtx.GetRegion(), func() (*kmsStores, error) {
		storage, err := reqCtx.GetStorage()
		if err != nil {
			return nil, err
		}
		return &kmsStores{
			keys:        kmsstore.NewKeyStore(storage, s.accountID, reqCtx.GetRegion()),
			aliases:     kmsstore.NewAliasStore(storage, s.accountID, reqCtx.GetRegion()),
			grants:      kmsstore.NewGrantStore(storage, s.accountID, reqCtx.GetRegion()),
			keyPolicies: kmsstore.NewKeyPolicyStore(storage, reqCtx.GetRegion()),
		}, nil
	})
}

// RegisterHandlers registers the KMS service handlers with the dispatcher.
func (s *KMSService) RegisterHandlers(d handler.Registrar) {
	d.RegisterHandlerForService("kms", "CreateKey", s.CreateKey)
	d.RegisterHandlerForService("kms", "DescribeKey", s.DescribeKey)
	d.RegisterHandlerForService("kms", "ListKeys", s.ListKeys)
	d.RegisterHandlerForService("kms", "EnableKey", s.EnableKey)
	d.RegisterHandlerForService("kms", "DisableKey", s.DisableKey)
	d.RegisterHandlerForService("kms", "ScheduleKeyDeletion", s.ScheduleKeyDeletion)
	d.RegisterHandlerForService("kms", "CancelKeyDeletion", s.CancelKeyDeletion)
	d.RegisterHandlerForService("kms", "UpdateKeyDescription", s.UpdateKeyDescription)
	d.RegisterHandlerForService("kms", "Encrypt", s.Encrypt)
	d.RegisterHandlerForService("kms", "Decrypt", s.Decrypt)
	d.RegisterHandlerForService("kms", "ReEncrypt", s.ReEncrypt)
	d.RegisterHandlerForService("kms", "GenerateDataKey", s.GenerateDataKey)
	d.RegisterHandlerForService("kms", "GenerateDataKeyWithoutPlaintext", s.GenerateDataKeyWithoutPlaintext)
	d.RegisterHandlerForService("kms", "GenerateRandom", s.GenerateRandom)
	d.RegisterHandlerForService("kms", "GetKeyPolicy", s.GetKeyPolicy)
	d.RegisterHandlerForService("kms", "PutKeyPolicy", s.PutKeyPolicy)
	d.RegisterHandlerForService("kms", "ListKeyPolicies", s.ListKeyPolicies)
	d.RegisterHandlerForService("kms", "CreateGrant", s.CreateGrant)
	d.RegisterHandlerForService("kms", "ListGrants", s.ListGrants)
	d.RegisterHandlerForService("kms", "ListRetirableGrants", s.ListRetirableGrants)
	d.RegisterHandlerForService("kms", "RevokeGrant", s.RevokeGrant)
	d.RegisterHandlerForService("kms", "RetireGrant", s.RetireGrant)
	d.RegisterHandlerForService("kms", "CreateAlias", s.CreateAlias)
	d.RegisterHandlerForService("kms", "DeleteAlias", s.DeleteAlias)
	d.RegisterHandlerForService("kms", "ListAliases", s.ListAliases)
	d.RegisterHandlerForService("kms", "UpdateAlias", s.UpdateAlias)
	d.RegisterHandlerForService("kms", "TagResource", s.TagResource)
	d.RegisterHandlerForService("kms", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("kms", "ListResourceTags", s.ListResourceTags)
	d.RegisterHandlerForService("kms", "EnableKeyRotation", s.EnableKeyRotation)
	d.RegisterHandlerForService("kms", "DisableKeyRotation", s.DisableKeyRotation)
	d.RegisterHandlerForService("kms", "GetKeyRotationStatus", s.GetKeyRotationStatus)
	d.RegisterHandlerForService("kms", "Sign", s.Sign)
	d.RegisterHandlerForService("kms", "Verify", s.Verify)
	d.RegisterHandlerForService("kms", "GetPublicKey", s.GetPublicKey)
	d.RegisterHandlerForService("kms", "GenerateMac", s.GenerateMac)
	d.RegisterHandlerForService("kms", "VerifyMac", s.VerifyMac)
	d.RegisterHandlerForService("kms", "GetParametersForImport", s.GetParametersForImport)
	d.RegisterHandlerForService("kms", "ImportKeyMaterial", s.ImportKeyMaterial)
	d.RegisterHandlerForService("kms", "DeleteImportedKeyMaterial", s.DeleteImportedKeyMaterial)
	d.RegisterHandlerForService("kms", "ReplicateKey", s.ReplicateKey)
	d.RegisterHandlerForService("kms", "UpdatePrimaryRegion", s.UpdatePrimaryRegion)
	d.RegisterHandlerForService("kms", "GenerateDataKeyPair", s.GenerateDataKeyPair)
	d.RegisterHandlerForService("kms", "GenerateDataKeyPairWithoutPlaintext", s.GenerateDataKeyPairWithoutPlaintext)
	d.RegisterHandlerForService("kms", "ListKeyRotations", s.ListKeyRotations)
}

func (s *KMSService) authorizeOperation(stores *kmsStores, principal, action, keyID string, encryptionContext map[string]string) error {
	if keyID == "" {
		return ErrKeyNotFound
	}

	key, err := stores.keys.Get(keyID)
	if err != nil {
		return ErrKeyNotFound
	}

	policyDoc, err := stores.keyPolicies.GetDefault(key.KeyID)
	if err != nil {
		return ErrAccessDenied
	}

	parsedPolicy, err := policy.ParseDocument(policyDoc.PolicyDocument)
	if err != nil {
		return ErrMalformedPolicy
	}

	ctx := &policy.EvaluationContext{
		Principal:         principal,
		Action:            "kms:" + action,
		Resource:          key.Arn,
		EncryptionContext: encryptionContext,
		Variables: map[string]string{
			"kms:EncryptionContextKeys": getEncryptionContextKeys(encryptionContext),
		},
	}

	decision := s.policyEvaluator.Evaluate(ctx, []*policy.Document{parsedPolicy})
	if decision.Effect == policy.DecisionEffectAllow {
		return nil
	}

	return ErrAccessDenied
}

func (s *KMSService) validateKeyState(key *kmsstore.Key) error {
	if key.KeyState == kmsstore.KeyStatePendingDeletion {
		return ErrKeyPendingDeletion
	}
	if key.KeyState == kmsstore.KeyStatePendingImport {
		return ErrKeyPendingImport
	}
	if !key.Enabled || key.KeyState == kmsstore.KeyStateDisabled {
		return ErrKeyDisabled
	}
	return nil
}

func (s *KMSService) getKeyID(params map[string]interface{}) string {
	keyID := request.GetStringParam(params, "KeyId")
	if keyID == "" {
		keyID = request.GetStringParam(params, "KeyID")
	}
	return keyID
}

func (s *KMSService) resolveKey(stores *kmsStores, params map[string]interface{}) (*kmsstore.Key, error) {
	keyID := s.getKeyID(params)
	if keyID == "" {
		return nil, ErrKeyNotFound
	}

	if stores.keys.ARNBuilder().IsAlias(keyID) {
		alias, err := stores.aliases.Get(keyID)
		if err != nil {
			return nil, NewAliasNotFoundError(keyID)
		}
		keyID = alias.TargetKeyID
	}

	key, err := stores.keys.Get(keyID)
	if err != nil {
		return nil, NewKeyNotFoundError(keyID)
	}
	return key, nil
}

func (s *KMSService) resolveCallerPrincipal(reqCtx *request.RequestContext, req *request.ParsedRequest) string {
	accessKeyId := req.AccessKeyID
	if accessKeyId == "" {
		accessKeyId = request.ExtractAccessKeyIDFromAuth(req.Headers.Get("Authorization"))
	}
	if accessKeyId == "" {
		accessKeyId = req.Headers.Get("X-Amz-Access-Key")
	}
	if accessKeyId == "" {
		return arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").IAM().Root()
	}
	iamStoreAny := reqCtx.GetIAMStore()
	if iamStoreAny == nil {
		return arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").IAM().Root()
	}
	iamStore, ok := iamStoreAny.(iamstore.IAMStoreInterface)
	if !ok {
		return arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").IAM().Root()
	}
	accessKey, err := iamStore.AccessKeys().Get(accessKeyId)
	if err != nil || accessKey == nil {
		return arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").IAM().Root()
	}
	return arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").IAM().User(accessKey.UserName)
}

func (s *KMSService) resolveAndAuthorizeKey(reqCtx *request.RequestContext, req *request.ParsedRequest, stores *kmsStores, action string, encryptionContext map[string]string) (*kmsstore.Key, error) {
	key, err := s.resolveKey(stores, req.Parameters)
	if err != nil {
		return nil, err
	}
	if err := s.authorizeOperation(stores, s.resolveCallerPrincipal(reqCtx, req), action, key.KeyID, encryptionContext); err != nil {
		return nil, err
	}
	return key, nil
}

func getEncryptionContextKeys(ctx map[string]string) string {
	if ctx == nil {
		return ""
	}
	keys := ""
	for k := range ctx {
		if keys != "" {
			keys += ","
		}
		keys += k
	}
	return keys
}

func parseEncryptionContext(params map[string]interface{}) map[string]string {
	if ec, ok := params["EncryptionContext"]; ok {
		if ecMap, ok := ec.(map[string]interface{}); ok {
			return request.CopyStringMap(ecMap)
		}
	}
	return nil
}

// ErrKeyNotFound is returned when a requested key does not exist.
var (
	ErrKeyNotFound = awserrors.NewAWSError("NotFoundException", "Key not found", http.StatusNotFound)
	// ErrKeyAlreadyExists is returned when attempting to create a key that already exists.
	ErrKeyAlreadyExists = awserrors.NewAWSError("AlreadyExistsException", "Key already exists", http.StatusConflict)
	// ErrAliasNotFound is returned when a requested alias does not exist.
	ErrAliasNotFound = awserrors.NewAWSError("NotFoundException", "Alias not found", http.StatusNotFound)
	// ErrAliasAlreadyExists is returned when attempting to create an alias that already exists.
	ErrAliasAlreadyExists = awserrors.NewAWSError("AlreadyExistsException", "Alias already exists", http.StatusConflict)
	// ErrGrantNotFound is returned when a requested grant does not exist.
	ErrGrantNotFound = awserrors.NewAWSError("NotFoundException", "Grant not found", http.StatusNotFound)
	// ErrKeyDisabled is returned when attempting to use a disabled key.
	ErrKeyDisabled = awserrors.NewAWSError("DisabledException", "Key is disabled", http.StatusBadRequest)
	// ErrKeyPendingDeletion is returned when the key is pending deletion.
	ErrKeyPendingDeletion = awserrors.NewAWSError("KMSInvalidStateException", "Key is pending deletion", http.StatusBadRequest)
	// ErrKeyPendingImport is returned when the key is pending import of key material.
	ErrKeyPendingImport = awserrors.NewAWSError("KMSInvalidStateException", "Key is pending import", http.StatusBadRequest)
	// ErrInvalidKeyUsage is returned when the key usage is invalid for the operation.
	ErrInvalidKeyUsage = awserrors.NewAWSError("InvalidKeyUsageException", "Invalid key usage", http.StatusBadRequest)
	// ErrInvalidKeySpec is returned when the key spec is invalid.
	ErrInvalidKeySpec = awserrors.NewAWSError("InvalidKeySpecException", "Invalid key spec", http.StatusBadRequest)
	// ErrInvalidAlgorithm is returned when the algorithm is invalid.
	ErrInvalidAlgorithm = awserrors.NewAWSError("InvalidAlgorithmException", "Invalid algorithm", http.StatusBadRequest)
	// ErrAccessDenied is returned when access is denied.
	ErrAccessDenied = awserrors.NewAWSError("AccessDeniedException", "Access denied", http.StatusForbidden)
	// ErrInvalidCiphertext is returned when the ciphertext is invalid.
	ErrInvalidCiphertext = awserrors.NewAWSError("InvalidCiphertextException", "Invalid ciphertext", http.StatusBadRequest)
	// ErrInvalidGrantToken is returned when the grant token is invalid.
	ErrInvalidGrantToken = awserrors.NewAWSError("InvalidGrantTokenException", "Invalid grant token", http.StatusBadRequest)
	// ErrInvalidAliasName is returned when the alias name is invalid.
	ErrInvalidAliasName = awserrors.NewAWSError("InvalidAliasNameException", "Invalid alias name", http.StatusBadRequest)
	// ErrDependencyTimeout is returned when a dependency operation times out.
	ErrDependencyTimeout = awserrors.NewAWSError("DependencyTimeoutException", "Dependency timeout", http.StatusServiceUnavailable)
	// ErrKMSInternal is returned for internal KMS failures.
	ErrKMSInternal = awserrors.NewAWSError("KMSInternalException", "An internal error occurred", http.StatusInternalServerError)
	// ErrMalformedPolicy is returned when the policy document is malformed.
	ErrMalformedPolicy = awserrors.NewAWSError("MalformedPolicyDocumentException", "Malformed policy document", http.StatusBadRequest)
	// ErrValidation is returned when a parameter validation fails.
	ErrValidation = awserrors.NewAWSError("ValidationException", "Invalid parameter", http.StatusBadRequest)
)

// NewKeyNotFoundError creates a new key not found error.
func NewKeyNotFoundError(keyID string) *awserrors.AWSError {
	return awserrors.NewAWSError("NotFoundException", "Key '"+keyID+"' does not exist", http.StatusNotFound)
}

// NewAliasNotFoundError creates a new alias not found error.
func NewAliasNotFoundError(aliasName string) *awserrors.AWSError {
	return awserrors.NewAWSError("NotFoundException", "Alias '"+aliasName+"' does not exist", http.StatusNotFound)
}

// EncryptString encrypts a plaintext string using the specified key.
// If no key ID is provided, it uses the default AWS SSM key.
func (s *KMSService) EncryptString(ctx context.Context, keyID string, plaintext string) (string, error) {
	if s.hsmBackend == nil {
		return plaintext, nil
	}

	resolvedKeyID := keyID
	if keyID == "" {
		resolvedKeyID = "alias/aws/ssm"
	}

	result, err := s.hsmBackend.Encrypt(resolvedKeyID, []byte(plaintext), nil)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(result.Ciphertext), nil
}

// DecryptString decrypts a ciphertext string using the specified key.
// If no key ID is provided, it uses the default AWS SSM key.
func (s *KMSService) DecryptString(ctx context.Context, keyID string, ciphertext string) (string, error) {
	if s.hsmBackend == nil {
		return ciphertext, nil
	}

	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	resolvedKeyID := keyID
	if keyID == "" {
		resolvedKeyID = "alias/aws/ssm"
	}

	result, err := s.hsmBackend.Decrypt(resolvedKeyID, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(result.Plaintext), nil
}

// EnsureDefaultSSMKey ensures that the default AWS SSM key exists.
// This is used internally for encrypting secrets stored by other AWS services.
func (s *KMSService) EnsureDefaultSSMKey() error {
	if s.hsmBackend == nil {
		return nil
	}

	if s.hsmBackend.KeyExists("alias/aws/ssm") {
		return nil
	}

	if err := s.hsmBackend.GenerateKey("alias/aws/ssm", hsm.KeySpecSymmetricDefault); err != nil {
		return err
	}

	return nil
}

// kmsBusAdapter adapts KMSService to satisfy eventbus.KMSInvoker without
// conflicting with the existing GenerateDataKey/Decrypt handler methods.
type kmsBusAdapter struct {
	*KMSService
}

func (a *kmsBusAdapter) GenerateDataKey(ctx context.Context, keyID string, keySpec string, encryptionContext map[string]string) (*eventbus.KMSDataKeyResult, error) {
	if a.hsmBackend == nil {
		return nil, fmt.Errorf("KMS HSM backend not configured")
	}
	result, err := a.hsmBackend.GenerateDataKey(keyID, keySpec, 0, encryptionContext)
	if err != nil {
		return nil, fmt.Errorf("failed to generate data key: %w", err)
	}
	return &eventbus.KMSDataKeyResult{
		Plaintext:      result.Plaintext,
		CiphertextBlob: result.Ciphertext,
	}, nil
}

func (a *kmsBusAdapter) Decrypt(ctx context.Context, keyID string, ciphertext []byte, encryptionContext map[string]string) ([]byte, error) {
	if a.hsmBackend == nil {
		return nil, fmt.Errorf("KMS HSM backend not configured")
	}
	result, err := a.hsmBackend.Decrypt(keyID, ciphertext, encryptionContext)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}
	return result.Plaintext, nil
}

func (a *kmsBusAdapter) KeyExists(ctx context.Context, keyID string) bool {
	if a.hsmBackend == nil {
		return false
	}
	return a.hsmBackend.KeyExists(keyID)
}

// KMSBusInvoker returns an eventbus.KMSInvoker backed by this service.
func (s *KMSService) KMSBusInvoker() eventbus.KMSInvoker {
	return &kmsBusAdapter{s}
}
