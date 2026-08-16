package cloudfront

import (
	"context"
	"fmt"
	"time"

	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
)

// tooManyPublicKeysMsg is the error message for exceeding the AWS hard
// limit on public keys in a single key group.
func tooManyPublicKeysMsg() string {
	return fmt.Sprintf("Number of public keys in key group exceeds the maximum of %d", cloudfrontstore.MaxPublicKeysPerKeyGroup)
}

// CreateKeyGroup creates a new CloudFront key group.
func (s *CloudFrontService) CreateKeyGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	configMap := request.GetMapParam(req.Parameters, "KeyGroupConfig")
	if configMap == nil {
		configMap = req.Parameters
	}

	var keyItems []string
	if itemsMap := request.GetMapParam(configMap, "Items"); itemsMap != nil {
		keyItems = parseStringItemList(configMap, "Items", "PublicKey")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	kg, err := s.createKeyGroupCore(store, CreateKeyGroupInput{
		Config: &cloudfrontstore.KeyGroupConfig{
			Name:    request.GetStringParam(configMap, "Name"),
			Items:   keyItems,
			Comment: request.GetStringParam(configMap, "Comment"),
		},
	})
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
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	kg, err := s.getKeyGroupCore(store, request.GetStringParam(req.Parameters, "Id"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"KeyGroup": formatKeyGroupResponse(kg),
		"ETag":     kg.ETag,
	}, nil
}

// GetKeyGroupConfig retrieves the configuration of a CloudFront key group.
func (s *CloudFrontService) GetKeyGroupConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	kg, err := s.getKeyGroupCore(store, request.GetStringParam(req.Parameters, "Id"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"KeyGroupConfig": formatKeyGroupConfigResponse(kg),
		"ETag":           kg.ETag,
	}, nil
}

// UpdateKeyGroup updates an existing CloudFront key group.
func (s *CloudFrontService) UpdateKeyGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	configMap := request.GetMapParam(req.Parameters, "KeyGroupConfig")
	if configMap == nil {
		configMap = req.Parameters
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	kg, err := s.updateKeyGroupCore(store, UpdateKeyGroupInput{
		Id:      request.GetStringParam(req.Parameters, "Id"),
		IfMatch: getIfMatch(req),
		Config: &cloudfrontstore.KeyGroupConfig{
			Name:    request.GetStringParam(configMap, "Name"),
			Items:   parseStringItemList(configMap, "Items", "PublicKey"),
			Comment: request.GetStringParam(configMap, "Comment"),
		},
	})
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
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteKeyGroupCore(store,
		request.GetStringParam(req.Parameters, "Id"), getIfMatch(req)); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// ListKeyGroups lists CloudFront key groups.
func (s *CloudFrontService) ListKeyGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.listKeyGroupsCore(store,
		request.GetStringParam(req.Parameters, "Marker"), request.GetIntParam(req.Parameters, "MaxItems"))
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.Groups))
	for _, kg := range result.Groups {
		items = append(items, map[string]interface{}{
			"KeyGroup": map[string]interface{}{
				"Id":               kg.ID,
				"LastModifiedTime": kg.LastModifiedAt.Format(time.RFC3339),
				"KeyGroupConfig":   formatKeyGroupConfigResponse(kg),
			},
		})
	}

	kgList := map[string]interface{}{
		"MaxItems": result.EffectiveMaxItems,
		"Quantity": len(items),
		"Items":    protocol.XMLElements{ElementName: "KeyGroupSummary", Items: items},
	}
	if result.NextMarker != "" {
		kgList["NextMarker"] = result.NextMarker
	}
	return map[string]interface{}{"KeyGroupList": kgList}, nil
}

func formatKeyGroupResponse(kg *cloudfrontstore.KeyGroup) map[string]interface{} {
	return map[string]interface{}{
		"Id":               kg.ID,
		"LastModifiedTime": kg.LastModifiedAt.Format(time.RFC3339),
		"KeyGroupConfig":   formatKeyGroupConfigResponse(kg),
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
	}
	return m
}

// isKeyGroupReferenced checks whether any distribution references the key group
// via TrustedKeyGroups in its cache behaviours. CloudFront returns ResourceInUse
// (409) when attempting to delete a key group that is still referenced.
// A store failure is reported as an error so deletion can fail closed.
func isKeyGroupReferenced(store *cloudfrontStores, keyGroupID string) (bool, error) {
	return scanDistributions(store, func(dist *cloudfrontstore.Distribution) bool {
		if dist.DistributionConfig == nil {
			return false
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
		return false
	})
}
