package cognitoidentityprovider

import (
	"errors"
	"strings"

	tagutil "vorpalstacks/internal/common/tags"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	arn "vorpalstacks/internal/utils/aws/arn"
)

// ---------------------------------------------------------------------------
// Core functions
// ---------------------------------------------------------------------------

// userPoolIDFromTagARN extracts the user pool ID from the user pool ARN that
// forms a tag resource key
// ("arn:<partition>:cognito-idp:<region>:<account>:userpool/<id>"). A resource
// key naming another service or resource type is a malformed parameter, not a
// missing resource.
func userPoolIDFromTagARN(resourceArn string) (string, error) {
	parsed, err := arn.ParseARN(resourceArn)
	if err != nil {
		return "", ErrInvalidParameter
	}
	if parsed.Service != "cognito-idp" {
		return "", ErrInvalidParameter
	}
	poolID, ok := strings.CutPrefix(parsed.Resource, "userpool/")
	if !ok || poolID == "" {
		return "", ErrInvalidParameter
	}
	return poolID, nil
}

// resolveTagTargetStore validates the ResourceArn against the Smithy ArnType
// constraints (length and pattern) and resolves the store of the user pool the
// ARN names: the tag model defines the resource as a user pool ARN, so every
// tag operation first confirms the pool exists — a well-formed ARN naming no
// pool fails with ResourceNotFoundException instead of persisting tags under
// an unowned key.
func (s *CognitoService) resolveTagTargetStore(region, resourceArn string) (cognitostore.CognitoStoreInterface, error) {
	if !validateArnType(resourceArn) {
		return nil, ErrInvalidParameter
	}
	poolID, err := userPoolIDFromTagARN(resourceArn)
	if err != nil {
		return nil, err
	}
	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetUserPool(poolID); err != nil {
		if errors.Is(err, cognitostore.ErrUserPoolNotFound) {
			return nil, ErrResourceNotFound
		}
		return nil, ErrInternalError
	}
	return store, nil
}

// tagResourceCore applies tags to a user pool. Validates the resource ARN
// against the Smithy ArnType constraints, confirms the named pool exists, and
// validates the tag set against Smithy constraints before delegating to the
// store.
func (s *CognitoService) tagResourceCore(region, resourceArn string, tags map[string]string) error {
	store, err := s.resolveTagTargetStore(region, resourceArn)
	if err != nil {
		return err
	}
	if len(tags) == 0 {
		return ErrInvalidParameter
	}
	if err := validateCognitoTags(tags); err != nil {
		return err
	}
	if err := store.Tag(resourceArn, tags); err != nil {
		return ErrInternalError
	}
	return nil
}

// untagResourceCore removes tags from a user pool. Validates the resource ARN
// against the Smithy ArnType constraints, confirms the named pool exists, and
// validates the tag-key set against Smithy constraints before delegating to
// the store.
func (s *CognitoService) untagResourceCore(region, resourceArn string, tagKeys []string) error {
	store, err := s.resolveTagTargetStore(region, resourceArn)
	if err != nil {
		return err
	}
	if len(tagKeys) == 0 {
		return ErrInvalidParameter
	}
	if err := validateCognitoTagKeys(tagKeys); err != nil {
		return err
	}
	if err := store.Untag(resourceArn, tagKeys); err != nil {
		return ErrInternalError
	}
	return nil
}

// listTagsForResourceCore retrieves all tags for a user pool. Validates the
// resource ARN against the Smithy ArnType constraints and confirms the named
// pool exists before delegating to the store.
func (s *CognitoService) listTagsForResourceCore(region, resourceArn string) (map[string]string, error) {
	store, err := s.resolveTagTargetStore(region, resourceArn)
	if err != nil {
		return nil, err
	}
	tags, err := store.ListAsSlice(resourceArn)
	if err != nil {
		return nil, err
	}
	return tagutil.ToMap(tags), nil
}
