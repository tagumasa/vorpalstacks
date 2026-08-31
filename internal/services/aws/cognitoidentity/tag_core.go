package cognitoidentity

import (
	"context"
	"errors"
	"strings"

	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	cognitoidentitystore "vorpalstacks/internal/store/aws/cognitoidentity"
	arn "vorpalstacks/internal/utils/aws/arn"
)

// cognitoIdentityMapError maps tag errors to the modelled Cognito Identity
// error codes.
func cognitoIdentityMapError(err error) error {
	// errors.As rather than a type switch so wrapped tag errors map the same
	// way as bare ones.
	var missingResource *tagutil.MissingResourceError
	var missingTags *tagutil.MissingTagsError
	var missingTagKeys *tagutil.MissingTagKeysError
	if errors.As(err, &missingResource) || errors.As(err, &missingTags) || errors.As(err, &missingTagKeys) {
		return ErrInvalidParameter
	}
	if errors.Is(err, cognitoidentitystore.ErrIdentityPoolNotFound) {
		return ErrResourceNotFound
	}
	return err
}

// cognitoIdentityPoolIDFromARN extracts the identity pool ID from the pool
// ARN that forms a tag resource key
// ("arn:aws:cognito-identity:<region>:<account>:identitypool/<id>").
func cognitoIdentityPoolIDFromARN(resourceArn string) string {
	resource := arn.ExtractResourceFromARN(resourceArn)
	if id, ok := strings.CutPrefix(resource, "identitypool/"); ok {
		return id
	}
	return ""
}

// cognitoIdentityTagConfig builds the shared tag handler configuration for
// Cognito Identity pools: the tag operations on both the HTTP API and any
// internal callers route through these store-backed closures. Every tag
// operation first resolves the pool behind the ARN, so a tag against a
// nonexistent pool fails with ResourceNotFoundException (404) as the
// service model specifies, instead of silently persisting tags under an
// unowned key.
func cognitoIdentityTagConfig(store cognitoidentitystore.CognitoIdentityStoreInterface) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.StandardConfig,
		ValidateResource: func(ctx context.Context, rawResourceKey string) error {
			poolID := cognitoIdentityPoolIDFromARN(rawResourceKey)
			if poolID == "" {
				// A resource key that is not an identity pool ARN is a
				// malformed parameter, not a missing resource.
				return ErrInvalidParameter
			}
			// The mapped error passes through the shared MapError unchanged
			// (an AWSError matches none of its sentinel branches): a missing
			// pool becomes ResourceNotFoundException, anything else is an
			// internal failure.
			_, err := store.GetIdentityPool(poolID)
			if err != nil {
				return cognitoIdentityMapError(err)
			}
			return nil
		},
		TagFunc: func(ctx context.Context, resourceKey string, tags []tagutil.Tag) error {
			if err := store.Tag(resourceKey, tagutil.ToMap(tags)); err != nil {
				return ErrInternalError
			}
			return nil
		},
		UntagFunc: func(ctx context.Context, resourceKey string, tagKeys []string) error {
			if err := store.Untag(resourceKey, tagKeys); err != nil {
				return ErrInternalError
			}
			return nil
		},
		ListFunc: func(ctx context.Context, resourceKey string) ([]tagutil.Tag, error) {
			m, err := store.List(resourceKey)
			if err != nil {
				return nil, ErrInternalError
			}
			return tagutil.MapToTags(m), nil
		},
		FormatResponse: func(tags []tagutil.Tag, _ string) (interface{}, error) {
			return map[string]interface{}{"Tags": tagutil.ToMap(tags)}, nil
		},
		EmptyResponse: func() (interface{}, error) {
			return response.EmptyResponse(), nil
		},
		MapError: cognitoIdentityMapError,
	}
}
