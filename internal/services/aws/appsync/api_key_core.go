package appsync

import (
	"time"

	"github.com/google/uuid"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"
)

// DefaultApiKeyValidityDays is the documented CreateApiKey default expiry:
// 7 days from creation time when the expires member is omitted.
const DefaultApiKeyValidityDays = 7

// createApiKeyInput carries the parsed CreateApiKey request payload.
type createApiKeyInput struct {
	ApiId       string
	Description string
	Expires     int64
}

// updateApiKeyInput carries the parsed UpdateApiKey request payload.
type updateApiKeyInput struct {
	ApiId       string
	Id          string
	Description string
	Expires     int64
}

// createApiKeyCore validates the request and persists a new API key for a
// GraphQL or Event API.
func (s *AppSyncService) createApiKeyCore(store *appsyncstore.AppSyncStore, in createApiKeyInput) (*appsyncstore.ApiKey, error) {
	if in.ApiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	// API keys can be created for both GraphQL APIs and Event APIs.
	if err := validateGraphqlApiExists(store, in.ApiId); err != nil {
		if err := validateEventApiExists(store, in.ApiId); err != nil {
			return nil, mapStoreErrorE(err)
		}
	}

	expires := in.Expires
	if expires == 0 {
		expires = time.Now().Add(DefaultApiKeyValidityDays * 24 * time.Hour).Unix()
	} else {
		if err := validateApiKeyExpiry(expires); err != nil {
			return nil, err
		}
	}

	apiKey := &appsyncstore.ApiKey{
		Id:          uuid.New().String(),
		Description: in.Description,
		Expires:     expires,
		// AWS returns Deletes equal to Expires at creation time.
		Deletes: expires,
	}

	if err := store.CreateApiKey(in.ApiId, apiKey); err != nil {
		return nil, mapStoreErrorE(err)
	}

	return apiKey, nil
}

// listApiKeysCore lists the API keys attached to a GraphQL or Event API.
func (s *AppSyncService) listApiKeysCore(store *appsyncstore.AppSyncStore, apiId string, maxResults int, nextToken string) ([]*appsyncstore.ApiKey, string, error) {
	if apiId == "" {
		return nil, "", NewBadRequestException("apiId is required")
	}

	opts, err := listOptionsFromParams(maxResults, nextToken)
	if err != nil {
		return nil, "", err
	}

	keys, nextToken, err := store.ListApiKeys(apiId, opts)
	if err != nil {
		return nil, "", mapStoreErrorE(err)
	}

	return keys, nextToken, nil
}

// updateApiKeyCore applies description/expiry updates to an existing API key.
func (s *AppSyncService) updateApiKeyCore(store *appsyncstore.AppSyncStore, in updateApiKeyInput) (*appsyncstore.ApiKey, error) {
	if in.ApiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}
	if in.Id == "" {
		return nil, NewBadRequestException("id is required")
	}

	apiKey, err := store.GetApiKey(in.ApiId, in.Id)
	if err != nil {
		return nil, mapStoreErrorE(err)
	}

	if in.Description != "" {
		apiKey.Description = in.Description
	}
	if in.Expires != 0 {
		if err := validateApiKeyExpiry(in.Expires); err != nil {
			return nil, err
		}
		apiKey.Expires = in.Expires
		// Deleting an API key happens after the expiry TTL; sync Deletes with Expires.
		apiKey.Deletes = in.Expires
	}

	if err := store.UpdateApiKey(in.ApiId, apiKey); err != nil {
		return nil, mapStoreErrorE(err)
	}

	return apiKey, nil
}

// deleteApiKeyCore removes an API key from its API.
func (s *AppSyncService) deleteApiKeyCore(store *appsyncstore.AppSyncStore, apiId, id string) error {
	if apiId == "" {
		return NewBadRequestException("apiId is required")
	}
	if id == "" {
		return NewBadRequestException("id is required")
	}

	if err := store.DeleteApiKey(apiId, id); err != nil {
		return mapStoreErrorE(err)
	}

	return nil
}
