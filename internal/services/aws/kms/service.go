// Package kms provides KMS (Key Management Service) operations for vorpalstacks.
package kms

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"
	"vorpalstacks/internal/common/kmsutil"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/iam/policy"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	"vorpalstacks/internal/services/aws/kms/hsm"
	storecommon "vorpalstacks/internal/store/aws/common"
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

// CascadeDeleteKey removes the key and every related artefact (HSM key
// material, grants, aliases, key policies, tags) from all stores. It is
// the single, authoritative entry point for hard-deleting a KMS key and
// must be used by every rollback path as well as any future hard-delete
// worker that purges keys whose PendingWindow has elapsed. Failing to go
// through this helper leaves orphaned grants (security: authz remnants),
// policies, aliases or tags that reference a non-existent key.
//
// Errors from individual stages are accumulated and returned as a joined
// error so that a partial failure does not silently mask other failures.
// Callers that absolutely must continue on partial failure (for example a
// best-effort rollback) may log the returned error and proceed.
func (ks *kmsStores) CascadeDeleteKey(hsmBackend hsm.Backend, keyID string) error {
	var errs []error

	if hsmBackend != nil {
		if err := hsmBackend.DeleteKey(keyID); err != nil {
			errs = append(errs, fmt.Errorf("hsm delete: %w", err))
		}
	}

	if err := ks.grants.DeleteByKeyID(keyID); err != nil {
		errs = append(errs, fmt.Errorf("grants delete: %w", err))
	}

	aliases, err := ks.aliases.ListForKeyID(keyID)
	if err != nil {
		errs = append(errs, fmt.Errorf("list aliases: %w", err))
	} else {
		for _, a := range aliases {
			if err := ks.aliases.Delete(a.AliasName); err != nil {
				errs = append(errs, fmt.Errorf("alias %s delete: %w", a.AliasName, err))
			}
		}
	}

	if err := ks.keyPolicies.DeleteAllForKey(keyID); err != nil {
		errs = append(errs, fmt.Errorf("policies delete: %w", err))
	}

	if err := ks.keys.TagStore.Delete(keyID); err != nil {
		errs = append(errs, fmt.Errorf("tags delete: %w", err))
	}

	if err := ks.keys.Delete(keyID); err != nil {
		errs = append(errs, fmt.Errorf("key delete: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("cascade delete failed for key %s: %v", keyID, errs)
	}
	return nil
}

// KMSService provides AWS KMS key and encryption operations.
type KMSService struct {
	hsmBackend        hsm.Backend
	policyEvaluator   *policy.PolicyEvaluator
	principalResolver eventbus.IAMPrincipalResolver
	accountID         string
	region            string
	stores            sync.Map // region → *kmsStores
	storageManager    *storage.RegionStorageManager
}

// SetPrincipalResolver registers the IAM principal resolver for grant validation.
func (s *KMSService) SetPrincipalResolver(resolver eventbus.IAMPrincipalResolver) {
	s.principalResolver = resolver
}

// SetStorageManager injects the region storage manager for lazy store creation.
func (s *KMSService) SetStorageManager(sm *storage.RegionStorageManager) {
	s.storageManager = sm
}

// NewKMSService creates a new KMS service instance.
func NewKMSService(accountID, region string, hsmBackend hsm.Backend) *KMSService {
	return &KMSService{
		hsmBackend:      hsmBackend,
		policyEvaluator: policy.NewPolicyEvaluator(),
		accountID:       accountID,
		region:          region,
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

// GetStoreForRegion returns the cached KMS stores for the given region,
// creating new store instances if not already cached.
func (s *KMSService) GetStoreForRegion(region string) (*kmsStores, error) {
	if v, ok := s.stores.Load(region); ok {
		return v.(*kmsStores), nil
	}
	if s.storageManager == nil {
		return nil, fmt.Errorf("kms storage manager not initialised")
	}
	st, err := s.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	stores := &kmsStores{
		keys:        kmsstore.NewKeyStore(st, s.accountID, region),
		aliases:     kmsstore.NewAliasStore(st, s.accountID, region),
		grants:      kmsstore.NewGrantStore(st, s.accountID, region),
		keyPolicies: kmsstore.NewKeyPolicyStore(st, region),
	}
	actual, _ := s.stores.LoadOrStore(region, stores)
	return actual.(*kmsStores), nil
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
	d.RegisterHandlerForService("kms", "RotateKeyOnDemand", s.RotateKeyOnDemand)
	d.RegisterHandlerForService("kms", "GetKeyLastUsage", s.GetKeyLastUsage)
	// Custom Key Store operations require CloudHSM/XKS infrastructure
	// which is not implemented. Register stubs so that SDK clients receive
	// UnsupportedOperationException instead of a generic "operation not
	// found" error.
	d.RegisterHandlerForService("kms", "ConnectCustomKeyStore", s.unsupportedOperation)
	d.RegisterHandlerForService("kms", "CreateCustomKeyStore", s.unsupportedOperation)
	d.RegisterHandlerForService("kms", "DeleteCustomKeyStore", s.unsupportedOperation)
	d.RegisterHandlerForService("kms", "DescribeCustomKeyStores", s.unsupportedOperation)
	d.RegisterHandlerForService("kms", "DisconnectCustomKeyStore", s.unsupportedOperation)
	d.RegisterHandlerForService("kms", "UpdateCustomKeyStore", s.unsupportedOperation)
	// DeriveSharedSecret requires ECDH key agreement support which is not
	// implemented.
	d.RegisterHandlerForService("kms", "DeriveSharedSecret", s.unsupportedOperation)
}

func (s *KMSService) getKeyID(params map[string]interface{}) string {
	keyID := request.GetStringParam(params, "KeyId")
	if keyID == "" {
		keyID = request.GetStringParam(params, "KeyID")
	}
	return keyID
}

// unresolvedCallerPrincipal is returned by resolveCallerPrincipal when the
// caller's access key cannot be mapped to an IAM user. Returning the account
// root ARN here (the previous behaviour) would silently grant root privileges
// on every key whose default policy permits the root principal, which is the
// common case. AWS denies the request in this scenario, so we emit a
// sentinel ARN that will not match any statement in any reasonable policy.
const unresolvedCallerPrincipal = "arn:vorpalstacks:iam::000000000000:unresolved-principal"

func (s *KMSService) resolveCallerPrincipal(reqCtx *request.RequestContext, req *request.ParsedRequest) string {
	accessKeyId := req.AccessKeyID
	if accessKeyId == "" {
		accessKeyId = request.ExtractAccessKeyIDFromAuth(req.Headers.Get("Authorization"))
	}
	if accessKeyId == "" {
		accessKeyId = req.Headers.Get("X-Amz-Access-Key")
	}
	if accessKeyId == "" {
		// No access key at all — this is an unauthenticated request; deny.
		return unresolvedCallerPrincipal
	}
	if s.principalResolver == nil {
		// Without a resolver we cannot map the access key to a user; deny
		// rather than fall back to root.
		return unresolvedCallerPrincipal
	}
	username, err := s.principalResolver.ResolvePrincipal(reqCtx, accessKeyId)
	if err != nil || username == "" {
		// Resolution failed (transient store error, deleted user, etc.).
		// Do NOT grant root as a side effect; deny.
		return unresolvedCallerPrincipal
	}
	return arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").IAM().User(username)
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
	// ErrSerializationException is returned when a request member carries a
	// value the aws-json-1.1 deserialiser cannot accept — an enum member
	// violation such as an algorithm string outside the modelled spec. It is
	// a protocol-level error no service model enumerates on its operations.
	ErrSerializationException = awserrors.NewAWSError("SerializationException", "The request payload couldn't be deserialized.", http.StatusBadRequest)
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
	// NewValidationError returns a ValidationException with a specific
	// detail message, mirroring AWS's \"1 validation error detected: ...\"
	// output. Static ErrValidation is fine for opaque failure but loses
	// actionable context for SDK clients; prefer this helper at call
	// sites that can identify the failing parameter.
	NewValidationError = func(detail string) *awserrors.AWSError {
		return awserrors.NewAWSError("ValidationException", detail, http.StatusBadRequest)
	}
	// ErrUnsupportedOperation is returned when the operation is not supported for the key type.
	ErrUnsupportedOperation = awserrors.NewAWSError("UnsupportedOperationException", "Operation is not supported for this key type", http.StatusBadRequest)
	// ErrTagException is returned when a tag limit or format is violated.
	ErrTagException = awserrors.NewAWSError("TagException", "Tag validation failed", http.StatusBadRequest)
	// ErrDryRunOperation is returned when the DryRun parameter is set to
	// true. AWS KMS uses HTTP 412 (Precondition Failed) for this error,
	// and the AWS SDK treats it as a successful dry-run verification.
	ErrDryRunOperation = awserrors.NewAWSError("DryRunOperation", "Request would have succeeded, but the DryRun flag is set.", http.StatusPreconditionFailed)
	// ErrExpiredImportToken is returned when the ImportToken supplied to
	// ImportKeyMaterial has passed its 24-hour validity window. The token
	// and wrapping key pair are generated by GetParametersForImport and
	// expire after ParametersValidTo.
	ErrExpiredImportToken = awserrors.NewAWSError("ExpiredImportTokenException", "The import token is expired.", http.StatusBadRequest)
	// ErrKMSInvalidSignature is returned by Verify when the supplied
	// signature does not match the message. AWS KMS throws
	// KMSInvalidSignatureException rather than returning SignatureValid=false.
	ErrKMSInvalidSignature = awserrors.NewAWSError("KMSInvalidSignatureException", "The signature did not match the message.", http.StatusBadRequest)
	// ErrKMSInvalidMac is returned by VerifyMac when the supplied MAC
	// does not match the message. AWS KMS throws KMSInvalidMacException
	// rather than returning MacValid=false.
	ErrKMSInvalidMac = awserrors.NewAWSError("KMSInvalidMacException", "The MAC did not match the message.", http.StatusBadRequest)
)

// checkKMSDryRun returns ErrDryRunOperation when the DryRun parameter is
// set to true. Per the Smithy model, DryRun is supported on Encrypt,
// Decrypt, ReEncrypt, GenerateDataKey, GenerateDataKeyWithoutPlaintext,
// GenerateDataKeyPair, GenerateDataKeyPairWithoutPlaintext, GenerateMac,
// VerifyMac, Sign, Verify, CreateGrant, RevokeGrant, and RetireGrant.
// The call should be placed after all validation and authorisation but
// before the actual cryptographic or state-mutating operation.
func checkKMSDryRun(params map[string]interface{}) error {
	if v, ok := params["DryRun"]; ok {
		switch val := v.(type) {
		case string:
			if val == "true" {
				return ErrDryRunOperation
			}
		case bool:
			if val {
				return ErrDryRunOperation
			}
		}
	}
	return nil
}

// unsupportedOperation is a stub handler for operations that require
// infrastructure not implemented on this platform (CloudHSM custom key
// stores, XKS, ECDH key agreement). Returning UnsupportedOperationException
// gives SDK clients a clear, AWS-compatible error instead of a generic
// "operation not found" dispatch error.
func (s *KMSService) unsupportedOperation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return nil, ErrUnsupportedOperation
}

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

	result, err := s.hsmBackend.Encrypt(resolvedKeyID, []byte(plaintext), hsm.EncryptionAlgorithmSymmetricDefault, nil)
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

	result, err := s.hsmBackend.Decrypt(resolvedKeyID, ciphertextBytes, hsm.EncryptionAlgorithmSymmetricDefault, nil)
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

// GenerateDataKey generates a data key encrypted under the specified KMS key.
func (a *kmsBusAdapter) GenerateDataKey(ctx context.Context, keyID string, keySpec string, encryptionContext map[string]string, sourceArn string) (*eventbus.KMSDataKeyResult, error) {
	if a.hsmBackend == nil {
		return nil, fmt.Errorf("KMS HSM backend not configured")
	}
	keyID = a.resolveCanonicalKeyID(keyID)
	// Evaluate grant constraints with the caller's sourceArn so that
	// SourceArn-constrained grants are honoured for internal calls.
	if sourceArn != "" {
		if !a.grantAllowsSourceArn(keyID, sourceArn) {
			return nil, ErrAccessDenied
		}
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

// Decrypt decrypts ciphertext that was encrypted under a KMS key.
func (a *kmsBusAdapter) Decrypt(ctx context.Context, keyID string, ciphertext []byte, encryptionContext map[string]string, sourceArn string) ([]byte, error) {
	if a.hsmBackend == nil {
		return nil, fmt.Errorf("KMS HSM backend not configured")
	}
	keyID = a.resolveCanonicalKeyID(keyID)
	if sourceArn != "" {
		if !a.grantAllowsSourceArn(keyID, sourceArn) {
			return nil, ErrAccessDenied
		}
	}
	result, err := a.hsmBackend.Decrypt(keyID, ciphertext, hsm.EncryptionAlgorithmSymmetricDefault, encryptionContext)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}
	return result.Plaintext, nil
}

// KeyExists reports whether the specified KMS key exists.
func (a *kmsBusAdapter) KeyExists(ctx context.Context, keyID string) bool {
	if a.hsmBackend == nil {
		return false
	}
	return a.hsmBackend.KeyExists(a.resolveCanonicalKeyID(keyID))
}

// SymmetricEncryptionKeyExists reports whether the key exists and is a
// symmetric encryption key (KeySpec SYMMETRIC_DEFAULT, KeyUsage
// ENCRYPT_DECRYPT). Consumers such as EventBridge Scheduler validate
// this key class at configuration time, the way AWS validates it via
// kms:DescribeKey.
func (a *kmsBusAdapter) SymmetricEncryptionKeyExists(ctx context.Context, keyID string) bool {
	stores, err := a.GetStoreForRegion(a.region)
	if err != nil {
		return false
	}
	key, err := stores.keys.Get(a.resolveCanonicalKeyID(keyID))
	if err != nil {
		return false
	}
	return key.KeySpec == kmsstore.KeySpecSymmetricDefault && key.KeyUsage == kmsstore.KeyUsageEncryptDecrypt
}

// KMSBusInvoker returns an eventbus.KMSInvoker backed by this service.
func (s *KMSService) KMSBusInvoker() eventbus.KMSInvoker {
	return &kmsBusAdapter{s}
}

// ---------------------------------------------------------------------------
// KMSKeyChecker adapter — lets other services (SQS, SNS, etc.) validate a KMS
// key ID for their KmsMasterKeyId attribute without importing the KMS store
// directly.
// ---------------------------------------------------------------------------

// kmsKeyCheckerAdapter adapts KMSService to satisfy kmsutil.Checker.
type kmsKeyCheckerAdapter struct {
	s *KMSService
}

// NewKeyChecker returns a kmsutil.Checker backed by this service.
func (s *KMSService) NewKeyChecker() kmsutil.Checker {
	return &kmsKeyCheckerAdapter{s}
}

// CheckKey resolves the key by ID/alias/ARN and verifies that it exists, is
// enabled, and has the ENCRYPT_DECRYPT key usage. Returns sentinel errors
// from the common package so callers can map them to service-specific error
// codes.
func (a *kmsKeyCheckerAdapter) CheckKey(ctx context.Context, region, keyID string) error {
	stores, err := a.s.GetStoreForRegion(region)
	if err != nil {
		return kmsutil.ErrKeyNotFound
	}

	key, err := a.s.resolveKey(stores, map[string]interface{}{"KeyId": keyID})
	if err != nil {
		return kmsutil.ErrKeyNotFound
	}

	switch key.KeyState {
	case kmsstore.KeyStateDisabled:
		return kmsutil.ErrKeyDisabled
	case kmsstore.KeyStatePendingDeletion, kmsstore.KeyStatePendingImport:
		return kmsutil.ErrKeyInvalidState
	}

	if key.KeyUsage != kmsstore.KeyUsageEncryptDecrypt {
		return kmsutil.ErrKeyInvalidUsage
	}

	return nil
}
