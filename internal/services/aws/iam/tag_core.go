// Transport-agnostic Core functions and family adapters for the tag
// operations (Tag/Untag/ListTags) shared by the eight taggable IAM
// resource families. The handlers parse the wire request into the input
// DTOs below, call the Core, and serialise the result; validation and the
// store read-modify-write live here alone.
package iam

import (
	"sort"
	"strconv"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// tagOps adapts the generic tag cores to one resource family: the
// required-member validation error, the family's not-found error and its
// store sentinel for a missing entity, and the keyed read, persist and
// tags-field accessors.
type tagOps[T any] struct {
	emptyErr   error
	notFoundFn func(string) *awserrors.AWSError
	notFound   error
	getFn      func(*iamstore.IAMStore, string) (T, error)
	putFn      func(*iamstore.IAMStore, T) error
	tagsFn     func(T) *[]tags.Tag
}

// TagResourceInput carries the parsed Tag{Resource} request.
type TagResourceInput struct {
	ResourceName string
	Tags         []tags.Tag
}

// UntagResourceInput carries the parsed Untag{Resource} request.
type UntagResourceInput struct {
	ResourceName string
	TagKeys      []string
}

// ListResourceTagsInput carries the parsed List{Resource}Tags request.
type ListResourceTagsInput struct {
	ResourceName string
	Marker       string
	MaxItems     int
}

// ListResourceTagsResult is the Core's paginated tag listing, ready for
// wire serialisation by the handlers.
type ListResourceTagsResult struct {
	Tags        []tags.Tag
	IsTruncated bool
	Marker      string
}

// tagResourceCore validates the new tag entries, merges them with the
// resource's existing tags under the per-resource tag limit, and persists
// the resource.
func tagResourceCore[T any](store *iamstore.IAMStore, ops tagOps[T], input *TagResourceInput) error {
	if input.ResourceName == "" {
		return ops.emptyErr
	}
	res, err := ops.getFn(store, input.ResourceName)
	if err != nil {
		return storeReadError(err, ops.notFound, ops.notFoundFn(input.ResourceName))
	}
	currentTags := ops.tagsFn(res)
	if err := validateTagEntries(input.Tags); err != nil {
		return err
	}
	merged := tags.Apply(*currentTags, input.Tags)
	if len(merged) > tags.MaxTagsPerResource {
		return NewInvalidInputError("Tags", "exceeds maximum of "+strconv.Itoa(tags.MaxTagsPerResource)+" tags per resource")
	}
	*currentTags = merged
	return ops.putFn(store, res)
}

// untagResourceCore removes the addressed tag keys from the resource and
// persists it.
func untagResourceCore[T any](store *iamstore.IAMStore, ops tagOps[T], input *UntagResourceInput) error {
	if input.ResourceName == "" {
		return ops.emptyErr
	}
	res, err := ops.getFn(store, input.ResourceName)
	if err != nil {
		return storeReadError(err, ops.notFound, ops.notFoundFn(input.ResourceName))
	}
	currentTags := ops.tagsFn(res)
	*currentTags = tags.RemoveByTagKeys(*currentTags, input.TagKeys)
	return ops.putFn(store, res)
}

// listResourceTagsCore returns the resource's tags, paginated by tag key.
func listResourceTagsCore[T any](store *iamstore.IAMStore, ops tagOps[T], input *ListResourceTagsInput) (*ListResourceTagsResult, error) {
	if input.ResourceName == "" {
		return nil, ops.emptyErr
	}
	res, err := ops.getFn(store, input.ResourceName)
	if err != nil {
		return nil, storeReadError(err, ops.notFound, ops.notFoundFn(input.ResourceName))
	}
	// Marker pagination walks the slice in ascending key order — the
	// stored order is not a contract, so the listing imposes it.
	storedTags := *ops.tagsFn(res)
	sorted := make([]tags.Tag, len(storedTags))
	copy(sorted, storedTags)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Key < sorted[j].Key })
	paged := pagination.PaginateSlice(sorted, input.Marker, input.MaxItems, func(tag tags.Tag) string {
		return tag.Key
	})
	return &ListResourceTagsResult{
		Tags:        paged.Items,
		IsTruncated: paged.IsTruncated,
		Marker:      paged.NextMarker,
	}, nil
}

var userTagOps = tagOps[*iamstore.User]{
	emptyErr:   NewValidationError("UserName"),
	notFoundFn: NewNoSuchUserError,
	notFound:   iamstore.ErrUserNotFound,
	getFn:      func(s *iamstore.IAMStore, n string) (*iamstore.User, error) { return s.Users().Get(n) },
	putFn:      func(s *iamstore.IAMStore, r *iamstore.User) error { return s.Users().Put(r) },
	tagsFn:     func(r *iamstore.User) *[]tags.Tag { return &r.Tags },
}

var roleTagOps = tagOps[*iamstore.Role]{
	emptyErr:   NewValidationError("RoleName"),
	notFoundFn: NewNoSuchRoleError,
	notFound:   iamstore.ErrRoleNotFound,
	getFn:      func(s *iamstore.IAMStore, n string) (*iamstore.Role, error) { return s.Roles().Get(n) },
	putFn:      func(s *iamstore.IAMStore, r *iamstore.Role) error { return s.Roles().Put(r) },
	tagsFn:     func(r *iamstore.Role) *[]tags.Tag { return &r.Tags },
}

var policyTagOps = tagOps[*iamstore.Policy]{
	emptyErr:   NewValidationError("PolicyArn"),
	notFoundFn: NewNoSuchPolicyError,
	notFound:   iamstore.ErrPolicyNotFound,
	getFn:      func(s *iamstore.IAMStore, n string) (*iamstore.Policy, error) { return s.Policies().Get(n) },
	putFn:      func(s *iamstore.IAMStore, r *iamstore.Policy) error { return s.Policies().Put(r) },
	tagsFn:     func(r *iamstore.Policy) *[]tags.Tag { return &r.Tags },
}

var instanceProfileTagOps = tagOps[*iamstore.InstanceProfile]{
	emptyErr:   NewValidationError("InstanceProfileName"),
	notFoundFn: NewNoSuchInstanceProfileError,
	notFound:   iamstore.ErrInstanceProfileNotFound,
	getFn: func(s *iamstore.IAMStore, n string) (*iamstore.InstanceProfile, error) {
		return s.InstanceProfiles().Get(n)
	},
	putFn:  func(s *iamstore.IAMStore, r *iamstore.InstanceProfile) error { return s.InstanceProfiles().Put(r) },
	tagsFn: func(r *iamstore.InstanceProfile) *[]tags.Tag { return &r.Tags },
}

var mfaDeviceTagOps = tagOps[*iamstore.VirtualMFADevice]{
	emptyErr:   NewValidationError("SerialNumber"),
	notFoundFn: NewNoSuchMFADeviceError,
	notFound:   iamstore.ErrMFADeviceNotFound,
	getFn:      func(s *iamstore.IAMStore, n string) (*iamstore.VirtualMFADevice, error) { return s.MFADevices().Get(n) },
	putFn:      func(s *iamstore.IAMStore, r *iamstore.VirtualMFADevice) error { return s.MFADevices().Put(r) },
	tagsFn:     func(r *iamstore.VirtualMFADevice) *[]tags.Tag { return &r.Tags },
}

var samlProviderTagOps = tagOps[*iamstore.SAMLProvider]{
	emptyErr:   NewValidationError("SAMLProviderArn"),
	notFoundFn: func(n string) *awserrors.AWSError { return NewNoSuchEntityError("SAML provider", n) },
	notFound:   iamstore.ErrSAMLProviderNotFound,
	getFn:      func(s *iamstore.IAMStore, n string) (*iamstore.SAMLProvider, error) { return s.SAMLProviders().Get(n) },
	putFn:      func(s *iamstore.IAMStore, r *iamstore.SAMLProvider) error { return s.SAMLProviders().Put(r) },
	tagsFn:     func(r *iamstore.SAMLProvider) *[]tags.Tag { return &r.Tags },
}

var oidcProviderTagOps = tagOps[*iamstore.OpenIDConnectProvider]{
	emptyErr:   NewValidationError("OpenIDConnectProviderArn"),
	notFoundFn: func(n string) *awserrors.AWSError { return NewNoSuchEntityError("OpenID Connect provider", n) },
	notFound:   iamstore.ErrOpenIDConnectProviderNotFound,
	getFn: func(s *iamstore.IAMStore, n string) (*iamstore.OpenIDConnectProvider, error) {
		return s.OpenIDConnectProviders().Get(n)
	},
	putFn: func(s *iamstore.IAMStore, r *iamstore.OpenIDConnectProvider) error {
		return s.OpenIDConnectProviders().Put(r)
	},
	tagsFn: func(r *iamstore.OpenIDConnectProvider) *[]tags.Tag { return &r.Tags },
}

var serverCertificateTagOps = tagOps[*iamstore.ServerCertificate]{
	emptyErr:   NewValidationError("ServerCertificateName"),
	notFoundFn: func(n string) *awserrors.AWSError { return NewNoSuchEntityError("server certificate", n) },
	notFound:   iamstore.ErrServerCertificateNotFound,
	getFn: func(s *iamstore.IAMStore, n string) (*iamstore.ServerCertificate, error) {
		return s.ServerCertificates().Get(n)
	},
	putFn:  func(s *iamstore.IAMStore, r *iamstore.ServerCertificate) error { return s.ServerCertificates().Put(r) },
	tagsFn: func(r *iamstore.ServerCertificate) *[]tags.Tag { return &r.Tags },
}
