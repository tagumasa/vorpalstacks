package cloudfront

import (
	"context"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
)

// maxPublicKeysPerKeyGroup is the AWS hard limit for public keys in a single key group.
const maxPublicKeysPerKeyGroup = 5

// CreateKeyGroup creates a new CloudFront key group.
func (s *CloudFrontService) CreateKeyGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	configMap := request.GetMapParam(req.Parameters, "KeyGroupConfig")
	if configMap == nil {
		configMap = req.Parameters
	}

	name := request.GetStringParam(configMap, "Name")
	if name == "" {
		return nil, awserrors.NewAWSError("InvalidArgument", "Name is required", 400)
	}

	var keyItems []string
	if itemsMap := request.GetMapParam(configMap, "Items"); itemsMap != nil {
		keyItems = parseStringItemList(configMap, "Items", "PublicKey")
	}

	if len(keyItems) == 0 {
		return nil, awserrors.NewAWSError("InvalidArgument", "At least one public key ID is required in Items", 400)
	}

	if len(keyItems) > maxPublicKeysPerKeyGroup {
		return nil, awserrors.NewAWSError("TooManyPublicKeysInKeyGroup",
			"Number of public keys in key group exceeds the maximum of 5", 400)
	}

	config := &cloudfrontstore.KeyGroupConfig{
		Name:    name,
		Items:   keyItems,
		Comment: request.GetStringParam(configMap, "Comment"),
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	existing, _ := store.keyGroups.GetByName(name)
	if existing != nil {
		return nil, awserrors.NewAWSError("KeyGroupAlreadyExists", "Key group with this name already exists", 409)
	}

	kg, err := store.keyGroups.Create(config)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"KeyGroup": formatKeyGroupResponse(kg),
		"Location": kg.ARN,
		"ETag":     kg.ETag,
	}, nil
}

// GetKeyGroup retrieves a CloudFront key group by ID.
func (s *CloudFrontService) GetKeyGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, awserrors.NewAWSError("InvalidArgument", "Id is required", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	kg, err := store.keyGroups.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchResource", "Key group not found: "+id, 404)
		}
		return nil, err
	}

	return map[string]interface{}{
		"KeyGroup": formatKeyGroupResponse(kg),
		"ETag":     kg.ETag,
	}, nil
}

// GetKeyGroupConfig retrieves the configuration of a CloudFront key group.
func (s *CloudFrontService) GetKeyGroupConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, awserrors.NewAWSError("InvalidArgument", "Id is required", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	kg, err := store.keyGroups.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchResource", "Key group not found: "+id, 404)
		}
		return nil, err
	}

	return map[string]interface{}{
		"KeyGroupConfig": formatKeyGroupConfigResponse(kg),
		"ETag":           kg.ETag,
	}, nil
}

// UpdateKeyGroup updates an existing CloudFront key group.
func (s *CloudFrontService) UpdateKeyGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, awserrors.NewAWSError("InvalidArgument", "Id is required", 400)
	}

	ifMatch := getIfMatch(req)
	if ifMatch == "" {
		return nil, awserrors.NewAWSError("InvalidIfMatchVersion",
			"The If-Match version is missing or not valid", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	existing, err := store.keyGroups.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchResource", "Key group not found: "+id, 404)
		}
		return nil, err
	}

	if ifMatch != "*" && existing.ETag != ifMatch {
		return nil, awserrors.NewAWSError("PreconditionFailed", preconditionFailedETagMsg, 412)
	}

	configMap := request.GetMapParam(req.Parameters, "KeyGroupConfig")
	if configMap == nil {
		configMap = req.Parameters
	}

	name := request.GetStringParam(configMap, "Name")
	if name == "" {
		return nil, awserrors.NewAWSError("InvalidArgument", "Name is required", 400)
	}

	var keyItems []string
	keyItems = parseStringItemList(configMap, "Items", "PublicKey")

	if len(keyItems) == 0 {
		return nil, awserrors.NewAWSError("InvalidArgument", "At least one public key ID is required in Items", 400)
	}

	if len(keyItems) > maxPublicKeysPerKeyGroup {
		return nil, awserrors.NewAWSError("TooManyPublicKeysInKeyGroup",
			"Number of public keys in key group exceeds the maximum of 5", 400)
	}

	if name != existing.KeyGroupConfig.Name {
		dup, _ := store.keyGroups.GetByName(name)
		if dup != nil {
			return nil, awserrors.NewAWSError("KeyGroupAlreadyExists",
				"Key group with this name already exists", 409)
		}
	}

	config := &cloudfrontstore.KeyGroupConfig{
		Name:    name,
		Items:   keyItems,
		Comment: request.GetStringParam(configMap, "Comment"),
	}

	kg, err := store.keyGroups.Update(id, config)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"KeyGroup": formatKeyGroupResponse(kg),
		"ETag":     kg.ETag,
	}, nil
}

// DeleteKeyGroup deletes a CloudFront key group.
func (s *CloudFrontService) DeleteKeyGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, awserrors.NewAWSError("InvalidArgument", "Id is required", 400)
	}

	ifMatch := getIfMatch(req)
	if ifMatch == "" {
		return nil, awserrors.NewAWSError("InvalidIfMatchVersion",
			"The If-Match version is missing or not valid", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	existing, err := store.keyGroups.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchResource", "Key group not found: "+id, 404)
		}
		return nil, err
	}

	if ifMatch != "*" && existing.ETag != ifMatch {
		return nil, awserrors.NewAWSError("PreconditionFailed", preconditionFailedETagMsg, 412)
	}

	if isKeyGroupReferenced(store, id) {
		return nil, awserrors.NewAWSError("ResourceInUse",
			"Cannot delete this key group because it is referenced by one or more distributions", 409)
	}

	if err := store.keyGroups.Delete(id); err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchResource", "Key group not found: "+id, 404)
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListKeyGroups lists CloudFront key groups.
func (s *CloudFrontService) ListKeyGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := request.GetIntParam(req.Parameters, "MaxItems")
	if maxItems == 0 {
		maxItems = 100
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := store.keyGroups.List(marker, maxItems)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.KeyGroups))
	for _, kg := range result.KeyGroups {
		items = append(items, map[string]interface{}{
			"KeyGroup": map[string]interface{}{
				"Id":               kg.ID,
				"LastModifiedTime": kg.LastModifiedAt.Format(time.RFC3339),
				"KeyGroupConfig":   formatKeyGroupConfigResponse(kg),
			},
		})
	}

	return map[string]interface{}{
		"KeyGroupList": map[string]interface{}{
			"MaxItems":    maxItems,
			"Quantity":    len(items),
			"IsTruncated": result.IsTruncated,
			"NextMarker":  result.NextMarker,
			"Items":       protocol.XMLElements{ElementName: "KeyGroupSummary", Items: items},
		},
	}, nil
}

func formatKeyGroupResponse(kg *cloudfrontstore.KeyGroup) map[string]interface{} {
	return map[string]interface{}{
		"Id":              kg.ID,
		"LastModifiedTime": kg.LastModifiedAt.Format(time.RFC3339),
		"KeyGroupConfig":  formatKeyGroupConfigResponse(kg),
	}
}

func formatKeyGroupConfigResponse(kg *cloudfrontstore.KeyGroup) map[string]interface{} {
	m := map[string]interface{}{
		"Name": kg.KeyGroupConfig.Name,
	}
	if kg.KeyGroupConfig.Comment != "" {
		m["Comment"] = kg.KeyGroupConfig.Comment
	}
	if len(kg.KeyGroupConfig.Items) > 0 {
		items := make([]interface{}, len(kg.KeyGroupConfig.Items))
		for i, item := range kg.KeyGroupConfig.Items {
			items[i] = item
		}
		m["Items"] = protocol.XMLElements{ElementName: "PublicKey", Items: items}
		m["Quantity"] = len(kg.KeyGroupConfig.Items)
	}
	return m
}

// isKeyGroupReferenced checks whether any distribution references the key group
// via TrustedKeyGroups in its cache behaviours. CloudFront returns ResourceInUse
// (409) when attempting to delete a key group that is still referenced.
func isKeyGroupReferenced(store *cloudfrontStores, keyGroupID string) bool {
	result, err := store.distributions.List("", 10000)
	if err != nil {
		return false
	}
	for _, dist := range result.Distributions {
		if dist.DistributionConfig == nil {
			continue
		}
		cfg := dist.DistributionConfig
		if cb := cfg.DefaultCacheBehavior; cb != nil && cb.TrustedKeyGroups != nil {
			for _, kgID := range cb.TrustedKeyGroups.Items {
				if kgID == keyGroupID {
					return true
				}
			}
		}
		if cbs := cfg.CacheBehaviors; cbs != nil {
			for _, cb := range cbs.Items {
				if cb != nil && cb.TrustedKeyGroups != nil {
					for _, kgID := range cb.TrustedKeyGroups.Items {
						if kgID == keyGroupID {
							return true
						}
					}
				}
			}
		}
	}
	return false
}
