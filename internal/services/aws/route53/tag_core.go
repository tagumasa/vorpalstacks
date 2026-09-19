package route53

import (
	"errors"
	"fmt"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/tags"
	route53store "vorpalstacks/internal/store/aws/route53"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input / Result structs
// ---------------------------------------------------------------------------

// ChangeTagsForResourceInput carries the raw resource identifiers and the
// parsed tag operations for ChangeTagsForResource. ResourceType and
// ResourceId keep their wire form — Core normalises and validates them.
type ChangeTagsForResourceInput struct {
	ResourceType  string
	ResourceId    string
	AddTags       []tags.Tag
	RemoveTagKeys []string
}

// ListTagsForResourceInput carries the raw resource identifiers for
// ListTagsForResource.
type ListTagsForResourceInput struct {
	ResourceType string
	ResourceId   string
}

// ListTagsForResourceResult holds the tags of one resource together with the
// normalised resource identity echoed in the response.
type ListTagsForResourceResult struct {
	ResourceType string
	ResourceId   string
	Tags         []tags.Tag
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// maxTagsPerRoute53Resource is the ChangeTagsForResource API tag quota: the
// API and CLI accept at most 10 tags per hosted zone or health check (the
// newer console allows 50, but the API this server implements does not).
// The store package owns the definition the tag store's budget shares.
const maxTagsPerRoute53Resource = route53store.MaxTagsPerResource

// route53TagLimits applies the API tag count quota on top of the standard
// key and value length bounds.
var route53TagLimits = tags.TagLimits{
	MaxCount:       maxTagsPerRoute53Resource,
	MaxKeyLength:   tags.MaxTagKeyLength,
	MaxValueLength: tags.MaxTagValueLength,
}

// parseResourceParamsCore validates and normalises ResourceType/ResourceId,
// returning (normalisedType, bareId, resourceKey, error).
func parseResourceParamsCore(resourceType, resourceId string) (string, string, string, error) {
	if resourceType == "" || resourceId == "" {
		return "", "", "", awserrors.NewAWSError("InvalidParameter", "ResourceType and ResourceId are required", 400)
	}

	normalizedType := strings.ToLower(resourceType)
	if normalizedType != "hostedzone" && normalizedType != "healthcheck" {
		return "", "", "", awserrors.NewAWSError("InvalidParameter", "ResourceType must be 'hostedzone' or 'healthcheck'", 400)
	}

	resourceId = strings.TrimPrefix(resourceId, "/hostedzone/")
	resourceId = strings.TrimPrefix(resourceId, "/healthcheck/")

	return normalizedType, resourceId, normalizedType + "/" + resourceId, nil
}

// verifyResourceExists checks whether the specified resource exists
// in the store, returning NoSuchHostedZone or NoSuchHealthCheck if not.
func verifyResourceExists(st *route53store.Route53Stores, resourceType, resourceId string) error {
	switch resourceType {
	case "hostedzone":
		if !st.HostedZones().Exists(resourceId) {
			return NewNoSuchHostedZoneError(resourceId)
		}
	case "healthcheck":
		if !st.HealthChecks().Exists(resourceId) {
			return NewNoSuchHealthCheckError(resourceId)
		}
	}
	return nil
}

// changeTagsForResourceCore is the single entry point applying tag changes
// to a hosted zone or health check. It validates the resource identity,
// verifies the resource exists, enforces the per-request tag count quota and
// the resulting ten-tag ceiling, then applies the additions and removals.
func changeTagsForResourceCore(st *route53store.Route53Stores, input ChangeTagsForResourceInput) error {
	normalizedType, resourceId, resourceKey, err := parseResourceParamsCore(input.ResourceType, input.ResourceId)
	if err != nil {
		return err
	}

	// Verify resource exists before applying tag operations.
	if err := verifyResourceExists(st, normalizedType, resourceId); err != nil {
		return err
	}

	if v, _ := tags.CheckTags(input.AddTags, route53TagLimits); v != tags.OK {
		switch v {
		case tags.TooManyTags:
			return awserrors.NewAWSError("InvalidInput",
				fmt.Sprintf("Number of tags must not exceed %d", maxTagsPerRoute53Resource), 400)
		case tags.TagKeyTooLong:
			return awserrors.NewAWSError("InvalidInput", "Tag key must not exceed 128 characters", 400)
		case tags.TagValueTooLong:
			return awserrors.NewAWSError("InvalidInput", "Tag value must not exceed 256 characters", 400)
		}
	}

	// Enforce the resulting tag total — compute the resulting key set
	// after applying both AddTags and RemoveTagKeys and bound it by the
	// same quota the request check applies. The store's tag budget
	// re-enforces the bound under its lock; this pre-check answers with
	// the operation's declared InvalidInput instead of a server error.
	existingTags, _ := st.Tags().ListTagsForResource(resourceKey)
	keySet := make(map[string]bool)
	for _, t := range existingTags {
		keySet[t.Key] = true
	}
	for _, k := range input.RemoveTagKeys {
		delete(keySet, k)
	}
	for _, t := range input.AddTags {
		keySet[t.Key] = true
	}
	if len(keySet) > maxTagsPerRoute53Resource {
		return awserrors.NewAWSError("InvalidInput",
			fmt.Sprintf("Number of tags must not exceed %d", maxTagsPerRoute53Resource), 400)
	}

	if len(input.AddTags) > 0 {
		if err := st.Tags().Tag(resourceKey, input.AddTags); err != nil {
			return mapTagStoreError(err, "apply tags")
		}
	}

	if len(input.RemoveTagKeys) > 0 {
		if err := st.Tags().Raw().Untag(resourceKey, input.RemoveTagKeys); err != nil {
			return mapTagStoreError(err, "remove tags")
		}
	}

	return nil
}

// listTagsForResourceCore is the single entry point listing the tags of a
// hosted zone or health check.
func listTagsForResourceCore(st *route53store.Route53Stores, input ListTagsForResourceInput) (*ListTagsForResourceResult, error) {
	normalizedType, resourceId, resourceKey, err := parseResourceParamsCore(input.ResourceType, input.ResourceId)
	if err != nil {
		return nil, err
	}

	// Verify resource exists before listing tags.
	if err := verifyResourceExists(st, normalizedType, resourceId); err != nil {
		return nil, err
	}

	tagList, err := st.Tags().ListTagsForResource(resourceKey)
	if err != nil {
		return nil, mapTagStoreError(err, "list tags")
	}

	return &ListTagsForResourceResult{
		ResourceType: normalizedType,
		ResourceId:   resourceId,
		Tags:         tagList,
	}, nil
}

// mapTagStoreError carries a tag-store failure to the wire with the
// operation's declared identity: a store error that is already wire-shaped
// (the tag budget answers quota overflow as InvalidInput) passes through
// unchanged, and a raw storage failure answers the storage-failure wrap
// the route53 cores use — InvalidInput-coded with a server status.
func mapTagStoreError(err error, action string) error {
	var awsErr *awserrors.AWSError
	if errors.As(err, &awsErr) {
		return awsErr
	}
	return awserrors.NewAWSError("InvalidInput", fmt.Sprintf("Failed to %s: %v", action, err), 500)
}
