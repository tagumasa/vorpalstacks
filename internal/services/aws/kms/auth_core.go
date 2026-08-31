package kms

// auth_core.go carries the key-resolution and authorisation layer of the
// KMS Core: resolving a wire key identifier (ID, ARN or alias) to the store
// key, evaluating the key policy / IAM / grant authorisation ladder, and
// recording key usage. Everything in this file sits behind the Core
// boundary — handlers and admin callers reach it as service-layer logic
// and never touch the store packages directly.

import (
	"sort"
	"strings"

	"vorpalstacks/internal/common/iam/policy"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	kmsstore "vorpalstacks/internal/store/aws/kms"
)

func (s *KMSService) authorizeOperation(stores *kmsStores, principal, action, keyID string, encryptionContext map[string]string) error {
	if keyID == "" {
		return ErrKeyNotFound
	}

	key, err := stores.keys.Get(keyID)
	if err != nil {
		return ErrKeyNotFound
	}

	// KMS authorisation order per AWS docs: explicit Deny wins; otherwise
	// key policy → IAM policy → grants. The policy evaluator honours both
	// Allow and Deny effects. An explicit Deny in the key policy is
	// definitive — it cannot be overridden by a grant. Only when the
	// policy returns DefaultDeny (no matching statement) do we fall
	// through to grant evaluation.

	policyDoc, err := stores.keyPolicies.GetDefault(key.KeyID)
	if err == nil {
		if parsedPolicy, perr := policy.ParseDocument(policyDoc.PolicyDocument); perr == nil {
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
			if decision.Effect == policy.DecisionEffectDeny {
				// Explicit Deny is definitive and cannot be overridden by
				// grants. This mirrors the AWS evaluation order where Deny
				// always wins.
				return ErrAccessDenied
			}
			if decision.Effect == policy.DecisionEffectAllow {
				return nil
			}
		}
	}

	if s.authorizeViaGrant(stores, principal, action, keyID, encryptionContext) {
		return nil
	}

	return ErrAccessDenied
}

// sourceArnMatches reports whether the requested resource ARN satisfies
// the grant's SourceArn constraint. AWS allows wildcard globbing in the
// final segment of the ARN (the resource portion). The constraint must
// match either exactly or via prefix wildcard. The constraint value
// comes from the aws:SourceArn global condition semantics.
func sourceArnMatches(constraintArn, requestArn string) bool {
	if constraintArn == "" {
		return true
	}
	if constraintArn == requestArn {
		return true
	}
	// AWS allows the constraint ARN to end with "*" to match any
	// resource under the prefix (e.g.
	// "arn:aws:sqs:us-east-1:111122223333:*").
	if strings.HasSuffix(constraintArn, "*") {
		prefix := constraintArn[:len(constraintArn)-1]
		return strings.HasPrefix(requestArn, prefix)
	}
	return false
}

// authorizeViaGrant implements the KMS grant evaluation path. A grant
// authorises the operation when all of the following hold:
//   - the grant is attached to the key being operated on
//   - the GranteePrincipal matches the caller's principal ARN exactly
//     (KMS does not currently support principal wildcards in grants)
//   - the grant Operations list contains the requested action
//   - every EncryptionContextEquals constraint is satisfied by an exact
//     match in the request encryption context, and every
//     EncryptionContextSubset constraint is satisfied by a subset match
//     (request context contains at least the constraint pairs)
func (s *KMSService) authorizeViaGrant(stores *kmsStores, principal, action, keyID string, encryptionContext map[string]string) bool {
	result, err := stores.grants.List(keyID, principal, "", 1000)
	if err != nil || result == nil {
		return false
	}
	// sourceArn comes from the calling service (e.g. S3 bucket ARN). When
	// non-empty, SourceArn-constrained grants are evaluated against it.
	// When empty, SourceArn-constrained grants are denied by default
	// (matching AWS deny-by-default).
	const requestArn = ""
	for _, g := range result.Grants {
		if g.GranteePrincipal != principal {
			continue
		}
		if !grantAllowsAction(g.Operations, action) {
			continue
		}
		if !grantConstraintsSatisfied(g.Constraints, encryptionContext, requestArn) {
			continue
		}
		return true
	}
	return false
}

// grantAllowsAction reports whether the grant's Operations list explicitly
// authorises the requested action. KMS uses the bare operation names
// (e.g. "Encrypt", "Decrypt", "CreateGrant") without the "kms:" prefix.
func grantAllowsAction(operations []string, action string) bool {
	for _, op := range operations {
		if op == action {
			return true
		}
	}
	return false
}

// grantConstraintsSatisfied evaluates EncryptionContextEquals (exact match
// of the entire constraint map against the request), EncryptionContextSubset
// (request must contain at least every constraint pair), and SourceArn
// (request resource ARN must match the constraint ARN, with optional
// trailing-segment wildcard). Returns true when the constraint block is
// nil or when all applicable constraints pass.
//
// The SourceArn constraint requires a request resource ARN. The KMS
// request pipeline does not currently expose aws:SourceArn (it is set by
// upstream AWS services like S3 when calling KMS on behalf of a resource).
// When SourceArn is configured on the grant and no requestArn is supplied,
// the constraint is treated as not satisfied, mirroring AWS deny-by-default
// behaviour.
func grantConstraintsSatisfied(constraints *kmsstore.GrantConstraints, encryptionContext map[string]string, requestArn string) bool {
	if constraints == nil {
		return true
	}
	if len(constraints.EncryptionContextEquals) > 0 {
		if len(encryptionContext) != len(constraints.EncryptionContextEquals) {
			return false
		}
		for k, v := range constraints.EncryptionContextEquals {
			if encryptionContext[k] != v {
				return false
			}
		}
	}
	if len(constraints.EncryptionContextSubset) > 0 {
		for k, v := range constraints.EncryptionContextSubset {
			if encryptionContext[k] != v {
				return false
			}
		}
	}
	if constraints.SourceArn != "" {
		// Deny-by-default: without a requestArn to match against, a
		// SourceArn-constrained grant cannot be satisfied.
		if !sourceArnMatches(constraints.SourceArn, requestArn) {
			return false
		}
	}
	return true
}

// validateKeyState reports whether the key is in a state that permits
// cryptographic operations. Read-only metadata operations (DescribeKey,
// ListResourceTags) bypass this check. The Enabled flag and KeyState field
// are kept in sync by the store layer, so checking KeyState alone is
// sufficient — the previous `!key.Enabled || KeyState == Disabled` test
// was redundant.
func (s *KMSService) validateKeyState(key *kmsstore.Key) error {
	switch key.KeyState {
	case kmsstore.KeyStatePendingDeletion:
		return ErrKeyPendingDeletion
	case kmsstore.KeyStatePendingImport:
		return ErrKeyPendingImport
	case kmsstore.KeyStateDisabled:
		return ErrKeyDisabled
	}
	return nil
}

func (s *KMSService) resolveKey(stores *kmsStores, params map[string]interface{}) (*kmsstore.Key, error) {
	keyID := s.getKeyID(params)
	if keyID == "" {
		return nil, ErrKeyNotFound
	}
	if err := validateKeyIdLength(keyID); err != nil {
		return nil, err
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
	// AWS rejects operations on a PendingDeletion key with
	// KMSInvalidStateException at resolution time; the granular
	// per-operation state check in validateKeyState still runs in the
	// handler, but resolving to a PendingDeletion key for read-only
	// metadata calls (DescribeKey, ListResourceTags) is permitted.
	// PendingImport keys likewise resolve successfully; the caller is
	// responsible for invoking validateKeyState when about to use the
	// key for crypto operations.
	return key, nil
}

// resolveKeyByKeyID resolves an already-extracted key identifier to the
// store key, sharing the resolveKey path (alias resolution included).
func (s *KMSService) resolveKeyByKeyID(stores *kmsStores, keyID string) (*kmsstore.Key, error) {
	return s.resolveKey(stores, map[string]interface{}{"KeyId": keyID})
}

func (s *KMSService) resolveAndAuthorizeKey(reqCtx *request.RequestContext, req *request.ParsedRequest, stores *kmsStores, action string, encryptionContext map[string]string) (*kmsstore.Key, error) {
	key, err := s.resolveKey(stores, req.Parameters)
	if err != nil {
		return nil, err
	}
	if err := s.authorizeOperation(stores, s.resolveCallerPrincipal(reqCtx, req), action, key.KeyID, encryptionContext); err != nil {
		return nil, err
	}

	// LastUsedAt tracking has been moved to markKeyLastUsed, which is
	// called by each crypto operation handler only after the HSM
	// operation succeeds. Updating LastUsedAt here would record failed
	// attempts as "usage", and the silent error swallow (_ =) masked
	// store update failures.

	return key, nil
}

// markKeyLastUsed records the timestamp and operation name of a
// successful cryptographic operation on the key. Called by Encrypt,
// Decrypt, ReEncrypt, GenerateDataKey, GenerateDataKeyPair, Sign,
// Verify, GenerateMac, and VerifyMac after the HSM operation succeeds.
// Errors are logged but not returned: LastUsedAt is a best-effort
// telemetry field and must not cause a successful crypto operation to
// fail.
func (s *KMSService) markKeyLastUsed(stores *kmsStores, keyID string, action string) {
	if err := stores.keys.UpdateLastUsed(keyID, action); err != nil {
		logs.Error("markKeyLastUsed: failed to update key", logs.String("keyId", keyID), logs.Err(err))
	}
}

func getEncryptionContextKeys(ctx map[string]string) string {
	if ctx == nil {
		return ""
	}
	// AWS requires the kms:EncryptionContextKeys condition key to be in
	// canonical (lexically sorted) form so that policy authors can match
	// the key set deterministically regardless of map iteration order.
	keys := make([]string, 0, len(ctx))
	for k := range ctx {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return strings.Join(keys, ",")
}

// resolveCanonicalKeyID normalises a key identifier (key ID, key ARN, or
// alias) to the canonical key ID the HSM backend indexes keys by. The store
// layer performs the same normalisation for direct KMS API calls; the bus
// adapter must match it because internal callers such as S3 pass whatever
// identifier the user supplied (SSEKMSKeyId, UpdateObjectEncryption KMSKeyArn).
// When the key cannot be resolved the original value is returned so the HSM
// lookup reports the not-found error.
func (a *kmsBusAdapter) resolveCanonicalKeyID(keyID string) string {
	stores, err := a.GetStoreForRegion(a.region)
	if err != nil {
		return keyID
	}
	key, err := stores.keys.Get(keyID)
	if err != nil {
		return keyID
	}
	return key.KeyID
}

// grantAllowsSourceArn checks whether any SourceArn-constrained grant on
// the key allows the given sourceArn. This enables the S3→KMS pattern
// where S3 creates a grant with SourceArn = bucket ARN.
//
// Returns true when no SourceArn-constrained grants exist (the common
// case), because without such constraints there is nothing to enforce.
func (s *KMSService) grantAllowsSourceArn(keyID, sourceArn string) bool {
	stores, err := s.GetStoreForRegion(s.region)
	if err != nil {
		// Cannot verify — allow to avoid blocking legitimate traffic.
		return true
	}
	result, err := stores.grants.List(keyID, "", "", 1000)
	if err != nil || result == nil || len(result.Grants) == 0 {
		// No grants at all — nothing to enforce.
		return true
	}
	hasSourceArnConstraint := false
	for _, g := range result.Grants {
		if g.Constraints != nil && g.Constraints.SourceArn != "" {
			hasSourceArnConstraint = true
			if sourceArnMatches(g.Constraints.SourceArn, sourceArn) {
				return true
			}
		} else {
			// Grant without SourceArn constraint always allows.
			return true
		}
	}
	// If there are SourceArn-constrained grants but none matched, deny.
	// If there are grants but none have SourceArn constraints, allow.
	return !hasSourceArnConstraint
}
