package sesv2

import (
	"context"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	sesv2store "vorpalstacks/internal/store/aws/sesv2"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// validateSesv2TagTarget resolves the SESv2 resource behind a tag resource
// ARN before any tag read or write. SESv2 tags five resource kinds, each in
// its own ARN namespace (identity/, configuration-set/, template/,
// contact-list/, dedicated-ip-pool/); every kind routes to its owning
// lookup so a tag against a nonexistent resource fails with the modelled
// NotFoundException. An ARN naming anything else addresses no resource this
// service can own.
func validateSesv2TagTarget(store sesv2store.SESv2StoreInterface, resourceArn string) error {
	_, _, _, _, resource := svcarn.SplitARN(resourceArn)
	kind, name, ok := strings.Cut(resource, "/")
	if !ok || name == "" {
		return ErrNotFound
	}
	switch kind {
	case "identity":
		_, err := store.GetEmailIdentity(name)
		return err
	case "configuration-set":
		_, err := store.GetConfigurationSet(name)
		return err
	case "template":
		_, err := store.GetEmailTemplate(name)
		return err
	case "contact-list":
		_, err := store.GetContactList(name)
		return err
	case "dedicated-ip-pool":
		_, err := store.GetDedicatedIpPool(name)
		return err
	default:
		return ErrNotFound
	}
}

// sesv2TagConfig wires the shared tag-trio handlers onto the SESv2 store.
// It lives in the Core layer so the TagFunc/UntagFunc/ListFunc store calls
// stay out of every handler closure.
func sesv2TagConfig(s *SESv2Service, reqCtx *request.RequestContext) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.StandardConfig,
		ValidateResource: func(ctx context.Context, resourceKey string) error {
			store, err := s.store(reqCtx)
			if err != nil {
				return err
			}
			return validateSesv2TagTarget(store, resourceKey)
		},
		TagFunc: func(ctx context.Context, resourceKey string, tags []tagutil.Tag) error {
			store, err := s.store(reqCtx)
			if err != nil {
				return err
			}
			return store.TagFromSlice(resourceKey, tags)
		},
		UntagFunc: func(ctx context.Context, resourceKey string, tagKeys []string) error {
			store, err := s.store(reqCtx)
			if err != nil {
				return err
			}
			return store.Untag(resourceKey, tagKeys)
		},
		ListFunc: func(ctx context.Context, resourceKey string) ([]tagutil.Tag, error) {
			store, err := s.store(reqCtx)
			if err != nil {
				return nil, err
			}
			return store.ListAsSlice(resourceKey)
		},
		EmptyResponse: func() (interface{}, error) { return response.EmptyResponse(), nil },
	}
}
