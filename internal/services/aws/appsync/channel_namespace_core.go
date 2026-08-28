package appsync

import (
	appsyncstore "vorpalstacks/internal/store/aws/appsync"
)

// channelNamespaceWithTags pairs a channel namespace with its tag-store view.
type channelNamespaceWithTags struct {
	Namespace *appsyncstore.ChannelNamespace
	Tags      map[string]string
}

// createChannelNamespaceInput carries the parsed CreateChannelNamespace
// request payload.
type createChannelNamespaceInput struct {
	ApiId              string
	Name               string
	CodeHandlers       string
	PublishAuthModes   []appsyncstore.AuthMode
	SubscribeAuthModes []appsyncstore.AuthMode
	Tags               map[string]string
	HandlerConfigs     *appsyncstore.HandlerConfigs
}

// updateChannelNamespaceInput carries the parsed UpdateChannelNamespace
// request payload. HasCodeHandlers distinguishes an explicitly supplied
// codeHandlers member from an omitted one.
type updateChannelNamespaceInput struct {
	ApiId              string
	Name               string
	CodeHandlers       string
	HasCodeHandlers    bool
	PublishAuthModes   []appsyncstore.AuthMode
	SubscribeAuthModes []appsyncstore.AuthMode
	HandlerConfigs     *appsyncstore.HandlerConfigs
}

// createChannelNamespaceCore validates the request and persists a new channel
// namespace within an Event API (v2), applying create-time tags.
func (s *AppSyncService) createChannelNamespaceCore(store *appsyncstore.AppSyncStore, in createChannelNamespaceInput) (*appsyncstore.ChannelNamespace, map[string]string, error) {
	if in.ApiId == "" {
		return nil, nil, NewBadRequestException("apiId is required")
	}
	if err := validateEventApiExists(store, in.ApiId); err != nil {
		return nil, nil, err
	}

	if in.Name == "" {
		return nil, nil, NewBadRequestException("name is required")
	}
	if err := validateNamespace(in.Name); err != nil {
		return nil, nil, err
	}

	if in.CodeHandlers != "" {
		if err := validateCode(in.CodeHandlers); err != nil {
			return nil, nil, err
		}
	}

	ns := &appsyncstore.ChannelNamespace{
		ApiId:              in.ApiId,
		Name:               in.Name,
		CodeHandlers:       in.CodeHandlers,
		HandlerConfigs:     in.HandlerConfigs,
		PublishAuthModes:   in.PublishAuthModes,
		SubscribeAuthModes: in.SubscribeAuthModes,
		Tags:               in.Tags,
	}

	created, err := store.CreateChannelNamespace(ns)
	if err != nil {
		return nil, nil, mapStoreErrorE(err)
	}

	if len(created.Tags) > 0 {
		tagMap := make(map[string]string, len(created.Tags))
		for k, v := range created.Tags {
			tagMap[k] = v
		}
		if err := store.TagStore.Tag(created.ChannelNamespaceArn, tagMap); err != nil {
			return nil, nil, err
		}
	}

	return created, listTagsIfAny(store, created.ChannelNamespaceArn), nil
}

// getChannelNamespaceCore fetches a channel namespace by API ID and name.
func (s *AppSyncService) getChannelNamespaceCore(store *appsyncstore.AppSyncStore, apiId, name string) (*appsyncstore.ChannelNamespace, map[string]string, error) {
	if apiId == "" || name == "" {
		return nil, nil, NewBadRequestException("apiId and name are required")
	}

	ns, err := store.GetChannelNamespace(apiId, name)
	if err != nil {
		return nil, nil, mapStoreErrorE(err)
	}

	return ns, listTagsIfAny(store, ns.ChannelNamespaceArn), nil
}

// updateChannelNamespaceCore applies an update to an existing channel
// namespace.
func (s *AppSyncService) updateChannelNamespaceCore(store *appsyncstore.AppSyncStore, in updateChannelNamespaceInput) (*appsyncstore.ChannelNamespace, map[string]string, error) {
	if in.ApiId == "" || in.Name == "" {
		return nil, nil, NewBadRequestException("apiId and name are required")
	}
	if err := validateNamespace(in.Name); err != nil {
		return nil, nil, err
	}

	ns := &appsyncstore.ChannelNamespace{
		ApiId:              in.ApiId,
		Name:               in.Name,
		PublishAuthModes:   in.PublishAuthModes,
		SubscribeAuthModes: in.SubscribeAuthModes,
	}

	if in.HasCodeHandlers {
		if in.CodeHandlers != "" {
			if err := validateCode(in.CodeHandlers); err != nil {
				return nil, nil, err
			}
		}
		ns.CodeHandlers = in.CodeHandlers
		ns.CodeHandlersSet = true
	}

	ns.HandlerConfigs = in.HandlerConfigs

	updated, err := store.UpdateChannelNamespace(ns)
	if err != nil {
		return nil, nil, mapStoreErrorE(err)
	}

	return updated, listTagsIfAny(store, updated.ChannelNamespaceArn), nil
}

// deleteChannelNamespaceCore removes a channel namespace and cleans up its
// active subscriptions to prevent stale data delivery.
func (s *AppSyncService) deleteChannelNamespaceCore(store *appsyncstore.AppSyncStore, apiId, name string) error {
	if apiId == "" || name == "" {
		return NewBadRequestException("apiId and name are required")
	}

	if _, err := store.GetChannelNamespace(apiId, name); err != nil {
		return mapStoreErrorE(err)
	}

	if err := store.DeleteChannelNamespace(apiId, name); err != nil {
		return mapStoreErrorE(err)
	}

	// Clean up active subscriptions on the deleted namespace to prevent
	// stale data delivery.
	s.eventServer.RemoveSubscriptionsByNamespace(name)

	return nil
}

// listChannelNamespacesCore lists the channel namespaces of an Event API,
// enriching each entry with its tag-store view.
func (s *AppSyncService) listChannelNamespacesCore(store *appsyncstore.AppSyncStore, apiId string, maxResults int, nextToken string) ([]channelNamespaceWithTags, string, error) {
	if apiId == "" {
		return nil, "", NewBadRequestException("apiId is required")
	}

	opts, err := listOptionsFromParams(maxResults, nextToken)
	if err != nil {
		return nil, "", err
	}

	namespaces, nextToken, err := store.ListChannelNamespaces(apiId, opts)
	if err != nil {
		return nil, "", mapStoreErrorE(err)
	}

	entries := make([]channelNamespaceWithTags, 0, len(namespaces))
	for _, ns := range namespaces {
		entries = append(entries, channelNamespaceWithTags{Namespace: ns, Tags: listTagsIfAny(store, ns.ChannelNamespaceArn)})
	}
	return entries, nextToken, nil
}
