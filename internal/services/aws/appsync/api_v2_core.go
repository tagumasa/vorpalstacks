package appsync

import (
	appsyncstore "vorpalstacks/internal/store/aws/appsync"

	storecommon "vorpalstacks/internal/store/aws/common"
)

// createApiInput carries the parsed CreateApi (Event API, v2) request payload.
type createApiInput struct {
	Name         string
	OwnerContact string
	WafWebAclArn string
	XrayEnabled  bool
}

// updateApiInput carries the parsed UpdateApi (Event API, v2) request payload.
// The Has* flags distinguish explicitly supplied members from omitted ones.
type updateApiInput struct {
	ApiId           string
	Name            string
	OwnerContact    string
	WafWebAclArn    string
	HasWafWebAclArn bool
	XrayEnabled     bool
	HasXrayEnabled  bool
}

// createApiCore validates the request and persists a new Event API (v2),
// applying create-time tags to the ARN-keyed tag store. The returned tags are
// the tag-store view read back after creation (nil when empty).
func (s *AppSyncService) createApiCore(store *appsyncstore.AppSyncStore, in createApiInput, eventConfig *appsyncstore.EventConfig, tagMap map[string]string) (*appsyncstore.Api, map[string]string, error) {
	if in.Name == "" {
		return nil, nil, NewBadRequestException("name is required")
	}
	if err := validateApiName(in.Name); err != nil {
		return nil, nil, err
	}

	api := &appsyncstore.Api{
		Name:         in.Name,
		EventConfig:  eventConfig,
		OwnerContact: in.OwnerContact,
		Tags:         tagMap,
		WafWebAclArn: in.WafWebAclArn,
		XrayEnabled:  in.XrayEnabled,
	}

	created, err := store.CreateApi(api)
	if err != nil {
		return nil, nil, mapStoreErrorE(err)
	}

	if len(created.Tags) > 0 {
		tagMap := make(map[string]string, len(created.Tags))
		for k, v := range created.Tags {
			tagMap[k] = v
		}
		if err := store.TagStore.Tag(created.Arn, tagMap); err != nil {
			return nil, nil, err
		}
	}

	return created, listTagsIfAny(store, created.Arn), nil
}

// getApiCore fetches an Event API (v2) by ID together with its tags.
func (s *AppSyncService) getApiCore(store *appsyncstore.AppSyncStore, apiId string) (*appsyncstore.Api, map[string]string, error) {
	if apiId == "" {
		return nil, nil, NewBadRequestException("apiId is required")
	}

	api, err := store.GetApiById(apiId)
	if err != nil {
		return nil, nil, mapStoreErrorE(err)
	}

	return api, listTagsIfAny(store, api.Arn), nil
}

// updateApiCore applies an update to an existing Event API (v2), preserving
// members that were not provided in the request.
func (s *AppSyncService) updateApiCore(store *appsyncstore.AppSyncStore, in updateApiInput, eventConfig *appsyncstore.EventConfig) (*appsyncstore.Api, map[string]string, error) {
	if in.ApiId == "" {
		return nil, nil, NewBadRequestException("apiId is required")
	}

	// Per Smithy model, name is required for UpdateApiRequest.
	if in.Name == "" {
		return nil, nil, NewBadRequestException("name is required")
	}
	if err := validateApiName(in.Name); err != nil {
		return nil, nil, err
	}

	// Fetch existing to preserve fields that were not provided in the request.
	// Without this, WafWebAclArn and XrayEnabled would be overwritten with
	// Go zero values on every update call that omits them.
	existing, err := store.GetApiById(in.ApiId)
	if err != nil {
		return nil, nil, mapStoreErrorE(err)
	}

	wafWebAclArn := existing.WafWebAclArn
	if in.HasWafWebAclArn {
		wafWebAclArn = in.WafWebAclArn
	}

	xrayEnabled := existing.XrayEnabled
	if in.HasXrayEnabled {
		xrayEnabled = in.XrayEnabled
	}

	api := &appsyncstore.Api{
		Name:         in.Name,
		OwnerContact: in.OwnerContact,
		WafWebAclArn: wafWebAclArn,
		XrayEnabled:  xrayEnabled,
	}

	// Per Smithy UpdateApiRequest, eventConfig is @required.
	if eventConfig == nil {
		return nil, nil, NewBadRequestException("eventConfig is required")
	}
	api.EventConfig = eventConfig

	updated, err := store.UpdateApiById(in.ApiId, api)
	if err != nil {
		return nil, nil, mapStoreErrorE(err)
	}

	return updated, listTagsIfAny(store, updated.Arn), nil
}

// deleteApiCore removes an Event API (v2) and disconnects its realtime
// event-server connections.
func (s *AppSyncService) deleteApiCore(store *appsyncstore.AppSyncStore, apiId string) error {
	if apiId == "" {
		return NewBadRequestException("apiId is required")
	}

	if _, err := store.GetApiById(apiId); err != nil {
		return mapStoreErrorE(err)
	}

	if err := store.DeleteApiById(apiId); err != nil {
		return mapStoreErrorE(err)
	}

	s.eventServer.DisconnectByApiId(apiId)

	return nil
}

// listTagsIfAny reads the tag-store view of a resource, returning nil when the
// read fails or the tag set is empty so callers can omit the member.
func listTagsIfAny(store *appsyncstore.AppSyncStore, arn string) map[string]string {
	if tags, err := store.TagStore.List(arn); err == nil && len(tags) > 0 {
		return tags
	}
	return nil
}

// apiWithTags pairs a listed Event API with its tag-store view.
type apiWithTags struct {
	Api  *appsyncstore.Api
	Tags map[string]string
}

// listApisCore lists Event APIs (v2) with pagination, enriching each entry
// with its tag-store view.
func (s *AppSyncService) listApisCore(store *appsyncstore.AppSyncStore, maxResults int, nextToken string) ([]apiWithTags, string, error) {
	if maxResults < 0 {
		maxResults = 0
	}
	if maxResults == 0 {
		maxResults = 25
	}
	if maxResults > 25 {
		return nil, "", NewBadRequestException("maxResults must be between 1 and 25")
	}
	apis, nextToken, err := store.ListApis(storecommon.ListOptions{
		MaxItems: maxResults,
		Marker:   nextToken,
	})
	if err != nil {
		return nil, "", mapStoreErrorE(err)
	}
	entries := make([]apiWithTags, 0, len(apis))
	for _, api := range apis {
		entries = append(entries, apiWithTags{Api: api, Tags: listTagsIfAny(store, api.Arn)})
	}
	return entries, nextToken, nil
}
