package cognitoidentityprovider

import (
	tagutil "vorpalstacks/internal/common/tags"
)

// ---------------------------------------------------------------------------
// Core functions
// ---------------------------------------------------------------------------

// tagResourceCore applies tags to a Cognito resource. Validates the ARN and
// tag set against Smithy constraints before delegating to the store.
func (s *CognitoService) tagResourceCore(region, resourceArn string, tags map[string]string) error {
	if err := validateCognitoResourceArn(resourceArn); err != nil {
		return err
	}
	if len(tags) == 0 {
		return ErrInvalidParameter
	}
	if err := validateCognitoTags(tags); err != nil {
		return err
	}
	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return err
	}
	if err := store.Tag(resourceArn, tags); err != nil {
		return ErrInternalError
	}
	return nil
}

// untagResourceCore removes tags from a Cognito resource. Validates the ARN
// and tag-key set against Smithy constraints before delegating to the store.
func (s *CognitoService) untagResourceCore(region, resourceArn string, tagKeys []string) error {
	if err := validateCognitoResourceArn(resourceArn); err != nil {
		return err
	}
	if len(tagKeys) == 0 {
		return ErrInvalidParameter
	}
	if err := validateCognitoTagKeys(tagKeys); err != nil {
		return err
	}
	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return err
	}
	if err := store.Untag(resourceArn, tagKeys); err != nil {
		return ErrInternalError
	}
	return nil
}

// listTagsForResourceCore retrieves all tags for a Cognito resource. Validates
// the ARN against the Smithy ArnType length constraint before delegation.
func (s *CognitoService) listTagsForResourceCore(region, resourceArn string) (map[string]string, error) {
	if err := validateCognitoResourceArn(resourceArn); err != nil {
		return nil, err
	}
	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}
	tags, err := store.ListAsSlice(resourceArn)
	if err != nil {
		return nil, err
	}
	return tagutil.ToMap(tags), nil
}
