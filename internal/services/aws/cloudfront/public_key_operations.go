package cloudfront

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	awserrors "vorpalstacks/internal/common/errors"
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

	name := request.GetStringParam(configMap, "Name")
	if name == "" {
		return nil, awserrors.NewAWSError("InvalidArgument", "Name is required", 400)
	}

	encodedKey := request.GetStringParam(configMap, "EncodedKey")
	if encodedKey == "" {
		return nil, awserrors.NewAWSError("InvalidArgument", "EncodedKey is required", 400)
	}

	callerRef := request.GetStringParam(configMap, "CallerReference")
	if callerRef == "" {
		// The caller reference is CloudFront's idempotency token; a
		// default must be unique per creation, so it is minted from
		// crypto/rand rather than the clock.
		callerRef = fmt.Sprintf("ref-%s", uuid.New().String())
	}

	config := &cloudfrontstore.PublicKeyConfig{
		CallerReference: callerRef,
		Name:            name,
		EncodedKey:      encodedKey,
		Comment:         request.GetStringParam(configMap, "Comment"),
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	existing, _ := store.publicKeys.GetByName(name)
	if existing != nil {
		return nil, awserrors.NewAWSError("PublicKeyAlreadyExists", "Public key with this name already exists", 409)
	}

	pk, err := store.publicKeys.Create(config)
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
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, awserrors.NewAWSError("InvalidArgument", "Id is required", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	pk, err := store.publicKeys.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchPublicKey", "Public key not found", 404)
		}
		return nil, err
	}

	return map[string]interface{}{
		"PublicKey": formatPublicKeyResponse(pk),
		"ETag":      pk.ETag,
	}, nil
}

// GetPublicKeyConfig retrieves the configuration of a CloudFront public key.
func (s *CloudFrontService) GetPublicKeyConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, awserrors.NewAWSError("InvalidArgument", "Id is required", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	pk, err := store.publicKeys.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchPublicKey", "Public key not found", 404)
		}
		return nil, err
	}

	return map[string]interface{}{
		"PublicKeyConfig": pk.PublicKeyConfig,
		"ETag":            pk.ETag,
	}, nil
}

// UpdatePublicKey updates an existing CloudFront public key.
func (s *CloudFrontService) UpdatePublicKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	existing, err := store.publicKeys.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchPublicKey", "Public key not found: "+id, 404)
		}
		return nil, err
	}

	if ifMatch != "*" && existing.ETag != ifMatch {
		return nil, awserrors.NewAWSError("PreconditionFailed", preconditionFailedETagMsg, 412)
	}

	configMap := request.GetMapParam(req.Parameters, "PublicKeyConfig")
	if configMap == nil {
		configMap = req.Parameters
	}

	name := request.GetStringParam(configMap, "Name")
	if name == "" {
		return nil, awserrors.NewAWSError("InvalidArgument", "Name is required", 400)
	}

	encodedKey := request.GetStringParam(configMap, "EncodedKey")
	if encodedKey == "" {
		return nil, awserrors.NewAWSError("InvalidArgument", "EncodedKey is required", 400)
	}
	if encodedKey != existing.PublicKeyConfig.EncodedKey {
		return nil, awserrors.NewAWSError("CannotChangeImmutablePublicKeyFields",
			"You can't change the value of a public key.", 400)
	}

	callerRef := request.GetStringParam(configMap, "CallerReference")
	if callerRef == "" {
		callerRef = existing.PublicKeyConfig.CallerReference
	}

	config := &cloudfrontstore.PublicKeyConfig{
		CallerReference: callerRef,
		Name:            name,
		EncodedKey:      encodedKey,
		Comment:         request.GetStringParam(configMap, "Comment"),
	}

	pk, err := store.publicKeys.Update(id, config)
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

	existing, err := store.publicKeys.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchPublicKey", "Public key not found: "+id, 404)
		}
		return nil, err
	}

	if ifMatch != "*" && existing.ETag != ifMatch {
		return nil, awserrors.NewAWSError("PreconditionFailed", preconditionFailedETagMsg, 412)
	}

	if isPublicKeyReferenced(store, id) {
		return nil, awserrors.NewAWSError("PublicKeyInUse",
			"The specified public key is in use by one or more key groups", 409)
	}

	if err := store.publicKeys.Delete(id); err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchPublicKey", "Public key not found: "+id, 404)
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListPublicKeys lists CloudFront public keys.
func (s *CloudFrontService) ListPublicKeys(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := request.GetIntParam(req.Parameters, "MaxItems")
	if maxItems == 0 {
		maxItems = 100
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := store.publicKeys.List(marker, maxItems)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.PublicKeys))
	for _, pk := range result.PublicKeys {
		items = append(items, map[string]interface{}{
			"Id":          pk.ID,
			"Name":        pk.PublicKeyConfig.Name,
			"CreatedTime": pk.CreatedTime.Format(time.RFC3339),
			"EncodedKey":  pk.PublicKeyConfig.EncodedKey,
			"Comment":     pk.PublicKeyConfig.Comment,
		})
	}

	return map[string]interface{}{
		"PublicKeyList": map[string]interface{}{
			"MaxItems":    maxItems,
			"Quantity":    len(items),
			"IsTruncated": result.IsTruncated,
			"NextMarker":  result.NextMarker,
			"Items":       protocol.XMLElements{ElementName: "PublicKeySummary", Items: items},
		},
	}, nil
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
func isPublicKeyReferenced(store *cloudfrontStores, publicKeyID string) bool {
	result, err := store.keyGroups.List("", 10000)
	if err != nil {
		return false
	}
	for _, kg := range result.KeyGroups {
		for _, pkID := range kg.KeyGroupConfig.Items {
			if pkID == publicKeyID {
				return true
			}
		}
	}
	return false
}
