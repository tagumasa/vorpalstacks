package cloudfront

import (
	"context"
	"time"

	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
)

// CreatePublicKey creates a new CloudFront public key.
func (s *CloudFrontService) CreatePublicKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	configMap := request.GetMapParam(req.Parameters, "PublicKeyConfig")
	if configMap == nil {
		configMap = req.Parameters
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	pk, err := s.createPublicKeyCore(store, CreatePublicKeyInput{
		Config: &cloudfrontstore.PublicKeyConfig{
			CallerReference: request.GetStringParam(configMap, "CallerReference"),
			Name:            request.GetStringParam(configMap, "Name"),
			EncodedKey:      request.GetStringParam(configMap, "EncodedKey"),
			Comment:         request.GetStringParam(configMap, "Comment"),
		},
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"PublicKey": formatPublicKeyResponse(pk),
		"Location":  pk.ARN,
		"ETag":      pk.ETag,
	}, nil
}

// GetPublicKey retrieves a CloudFront public key by ID.
func (s *CloudFrontService) GetPublicKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	pk, err := s.getPublicKeyCore(store, request.GetStringParam(req.Parameters, "Id"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"PublicKey": formatPublicKeyResponse(pk),
		"ETag":      pk.ETag,
	}, nil
}

// GetPublicKeyConfig retrieves the configuration of a CloudFront public key.
func (s *CloudFrontService) GetPublicKeyConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	pk, err := s.getPublicKeyCore(store, request.GetStringParam(req.Parameters, "Id"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"PublicKeyConfig": pk.PublicKeyConfig,
		"ETag":            pk.ETag,
	}, nil
}

// UpdatePublicKey updates an existing CloudFront public key.
func (s *CloudFrontService) UpdatePublicKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	configMap := request.GetMapParam(req.Parameters, "PublicKeyConfig")
	if configMap == nil {
		configMap = req.Parameters
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	pk, err := s.updatePublicKeyCore(store, UpdatePublicKeyInput{
		Id:      request.GetStringParam(req.Parameters, "Id"),
		IfMatch: getIfMatch(req),
		Config: &cloudfrontstore.PublicKeyConfig{
			CallerReference: request.GetStringParam(configMap, "CallerReference"),
			Name:            request.GetStringParam(configMap, "Name"),
			EncodedKey:      request.GetStringParam(configMap, "EncodedKey"),
			Comment:         request.GetStringParam(configMap, "Comment"),
		},
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"PublicKey": formatPublicKeyResponse(pk),
		"ETag":      pk.ETag,
	}, nil
}

// DeletePublicKey deletes a CloudFront public key.
func (s *CloudFrontService) DeletePublicKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deletePublicKeyCore(store,
		request.GetStringParam(req.Parameters, "Id"), getIfMatch(req)); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// ListPublicKeys lists CloudFront public keys.
func (s *CloudFrontService) ListPublicKeys(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.listPublicKeysCore(store,
		request.GetStringParam(req.Parameters, "Marker"), request.GetIntParam(req.Parameters, "MaxItems"))
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.Keys))
	for _, pk := range result.Keys {
		items = append(items, map[string]interface{}{
			"Id":          pk.ID,
			"Name":        pk.PublicKeyConfig.Name,
			"CreatedTime": pk.CreatedTime.Format(time.RFC3339),
			"EncodedKey":  pk.PublicKeyConfig.EncodedKey,
			"Comment":     pk.PublicKeyConfig.Comment,
		})
	}

	pkList := map[string]interface{}{
		"MaxItems": result.EffectiveMaxItems,
		"Quantity": len(items),
		"Items":    protocol.XMLElements{ElementName: "PublicKeySummary", Items: items},
	}
	if result.NextMarker != "" {
		pkList["NextMarker"] = result.NextMarker
	}
	return map[string]interface{}{"PublicKeyList": pkList}, nil
}

func formatPublicKeyResponse(pk *cloudfrontstore.PublicKey) map[string]interface{} {
	return map[string]interface{}{
		"Id":          pk.ID,
		"CreatedTime": pk.CreatedTime.Format(time.RFC3339),
		"PublicKeyConfig": map[string]interface{}{
			"CallerReference": pk.PublicKeyConfig.CallerReference,
			"Name":            pk.PublicKeyConfig.Name,
			"EncodedKey":      pk.PublicKeyConfig.EncodedKey,
			"Comment":         pk.PublicKeyConfig.Comment,
		},
	}
}

// isPublicKeyReferenced checks whether any key group references the public key
// ID in its Items list. CloudFront returns PublicKeyInUse (409) when
// attempting to delete a public key that is still referenced by a key group.
// A store failure is reported as an error so deletion can fail closed.
func isPublicKeyReferenced(store *cloudfrontStores, publicKeyID string) (bool, error) {
	return scanKeyGroups(store, func(kg *cloudfrontstore.KeyGroup) bool {
		for _, pkID := range kg.KeyGroupConfig.Items {
			if pkID == publicKeyID {
				return true
			}
		}
		return false
	})
}
