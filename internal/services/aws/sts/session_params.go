package sts

import (
	"fmt"
	"regexp"

	"vorpalstacks/internal/common/request"
)

// externalIdPattern mirrors the Smithy externalIdType trait: [\w+=,.@:\/-]*.
// Combined with the length check (2-1224) it enforces the AWS AssumeRole
// ExternalId parameter constraints.
var externalIdPattern = regexp.MustCompile(`^[\w+=,.@:\/-]+$`)

// sessionTagKeyPattern mirrors the Smithy tagKeyType pattern.
var sessionTagKeyPattern = regexp.MustCompile(`^[\p{L}\p{Z}\p{N}_.:/=+\-@]+$`)

// sessionTagValuePattern mirrors the Smithy tagValueType pattern.
var sessionTagValuePattern = regexp.MustCompile(`^[\p{L}\p{Z}\p{N}_.:/=+\-@]*$`)

// extractSessionTags parses Tags.member.N.Key / Tags.member.N.Value pairs from
// the flat query-parameter map and validates them against the Smithy
// tagKeyType (1-128 chars) and tagValueType (0-256 chars) traits. At most 50
// tags are accepted per the tagListType length constraint.
func extractSessionTags(params map[string]interface{}) (map[string]string, error) {
	tags := make(map[string]string)
	for i := 1; ; i++ {
		key := request.GetStringParam(params, fmt.Sprintf("Tags.member.%d.Key", i))
		if key == "" {
			break
		}
		if i > 50 {
			return nil, ErrTooManySessionTags
		}
		if len(key) < 1 || len(key) > 128 || !sessionTagKeyPattern.MatchString(key) {
			return nil, ErrInvalidSessionTag
		}
		value := request.GetStringParam(params, fmt.Sprintf("Tags.member.%d.Value", i))
		if len(value) > 256 || !sessionTagValuePattern.MatchString(value) {
			return nil, ErrInvalidSessionTag
		}
		if _, exists := tags[key]; exists {
			return nil, ErrDuplicateSessionTagKey
		}
		tags[key] = value
	}
	if len(tags) == 0 {
		return nil, nil
	}
	return tags, nil
}

// extractTransitiveTagKeys parses TransitiveTagKeys.member.N from the flat
// query-parameter map. Each key is validated against the Smithy tagKeyType
// trait (1-128 chars). At most 50 keys are accepted. Duplicate keys are
// rejected.
func extractTransitiveTagKeys(params map[string]interface{}) ([]string, error) {
	var keys []string
	seen := make(map[string]bool)
	for i := 1; ; i++ {
		key := request.GetStringParam(params, fmt.Sprintf("TransitiveTagKeys.member.%d", i))
		if key == "" {
			break
		}
		if i > 50 {
			return nil, ErrTooManySessionTags
		}
		if len(key) < 1 || len(key) > 128 || !sessionTagKeyPattern.MatchString(key) {
			return nil, ErrInvalidSessionTag
		}
		if seen[key] {
			return nil, ErrDuplicateSessionTagKey
		}
		seen[key] = true
		keys = append(keys, key)
	}
	return keys, nil
}

// extractPolicyArns parses PolicyArns.member.N.arn from the flat query-
// parameter map. Each ARN is validated against the Smithy arnType
// constraint. At most 10 ARNs are accepted per the AWS documentation
// ("You can provide up to 10 managed policy ARNs"). Duplicate ARNs are
// rejected to mirror the ErrDuplicateSessionTagKey guard for session
// tags. The slice preserves caller order.
func extractPolicyArns(params map[string]interface{}) ([]string, error) {
	var arns []string
	seen := make(map[string]bool)
	for i := 1; ; i++ {
		key := fmt.Sprintf("PolicyArns.member.%d.arn", i)
		arn := request.GetStringParam(params, key)
		if arn == "" {
			break
		}
		if len(arns) >= maxPolicyArns {
			return nil, ErrTooManyPolicyArns
		}
		if err := validateARN(arn); err != nil {
			return nil, err
		}
		if seen[arn] {
			return nil, ErrDuplicatePolicyArn
		}
		seen[arn] = true
		arns = append(arns, arn)
	}
	return arns, nil
}

// extractProvidedContexts parses ProvidedContexts.member.N.ProviderArn and
// ProvidedContexts.member.N.ContextAssertion from the flat query-parameter
// map. Smithy ProvidedContextsListType limits the list to 1-5 entries; the
// trust policy evaluator inspects sts:ProvidedContextProviderArn and
// sts:ProvidedContextAssertion condition keys. Provided contexts are
// signed-and-encrypted by STS in real AWS; VorpalStacks does not verify the
// signature but exposes the values to the evaluation context for
// compatibility with policy templates that reference them.
func extractProvidedContexts(params map[string]interface{}) ([]ProvidedContextEntry, error) {
	var entries []ProvidedContextEntry
	for i := 1; ; i++ {
		providerKey := fmt.Sprintf("ProvidedContexts.member.%d.ProviderArn", i)
		assertionKey := fmt.Sprintf("ProvidedContexts.member.%d.ContextAssertion", i)
		providerArn := request.GetStringParam(params, providerKey)
		contextAssertion := request.GetStringParam(params, assertionKey)
		if providerArn == "" && contextAssertion == "" {
			break
		}
		// Per-entry validation: contextAssertionType length 4-2048
		// when non-empty; ProviderArn arnType format when non-empty.
		if contextAssertion != "" {
			if len(contextAssertion) < minContextAssertionLen || len(contextAssertion) > maxContextAssertionLen {
				return nil, ErrInvalidContextAssertion
			}
		}
		if providerArn != "" {
			if err := validateARN(providerArn); err != nil {
				return nil, ErrInvalidProviderArn
			}
		}
		entries = append(entries, ProvidedContextEntry{
			ProviderArn:      providerArn,
			ContextAssertion: contextAssertion,
		})
	}
	// Smithy ProvidedContextsListType: length 1-5 when provided.
	// An empty slice means the caller did not supply the parameter at
	// all, which is valid (ProvidedContexts is optional on AssumeRole).
	if len(entries) > maxProvidedContexts {
		return nil, ErrTooManyProvidedContexts
	}
	return entries, nil
}

func computePackedPolicySize(policy string, policyArns []string, transitiveKeys []string, tags map[string]string) int32 {
	const maxPolicySize = 2048

	// AWS serialises the session policy, managed-policy ARNs, session
	// tags, transitive tag keys and source identity into a packed JSON
	// document whose total byte length must not exceed maxPolicySize.
	// The raw character counts alone underestimate the packed size
	// because they ignore JSON structural overhead (quotes, brackets,
	// key names, separators). We add per-element overhead so the
	// reported percentage more closely tracks real AWS.
	totalSize := len(policy)

	// Managed-policy ARNs — each serialised as a JSON string element
	// inside an array: ["arn:...","arn:..."]. Per-element overhead is
	// 3 bytes (two quotes + one comma/bracket).
	if len(policyArns) > 0 {
		totalSize += 2 // array brackets
	}
	for _, arn := range policyArns {
		totalSize += len(arn) + 3
	}

	// Session tags — each serialised as {"Key":"...","Value":"..."}.
	// Structural overhead per tag: {"Key":"","Value":""} = 20 bytes.
	if len(tags) > 0 {
		totalSize += 2 // array brackets
	}
	for key, value := range tags {
		totalSize += len(key) + len(value) + 20
	}

	// Transitive tag keys — each serialised as a JSON string.
	if len(transitiveKeys) > 0 {
		totalSize += 2 // array brackets
	}
	for _, val := range transitiveKeys {
		totalSize += len(val) + 3
	}

	if totalSize <= 0 {
		return 0
	}
	// Use ceiling to prevent integer-truncation bypass: a 2049-byte
	// packed policy would floor-divide to exactly 100, passing the
	// >100 check. Ceiling ensures any value exceeding the limit reports
	// at least 101.
	pct := totalSize * 100
	result := pct / maxPolicySize
	if pct%maxPolicySize != 0 {
		result++
	}
	return int32(result)
}

// withSourceIdentity returns resp with the SourceIdentity field set only when
// si is non-empty. The Smithy AssumeRole*Response shapes declare SourceIdentity
// as an optional member; AWS Query protocol responses omit the field entirely
// when the caller did not supply it, rather than serialising an empty string.
func withSourceIdentity(resp map[string]interface{}, si string) map[string]interface{} {
	if si != "" {
		resp["SourceIdentity"] = si
	}
	return resp
}
