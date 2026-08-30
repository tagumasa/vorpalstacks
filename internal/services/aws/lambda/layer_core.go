package lambda

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	storecommon "vorpalstacks/internal/store/aws/common"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// LayerVersionCreateInput carries the wire members of a
// PublishLayerVersion request. Content keeps its raw wire form so the Core
// applies the ZipFile/S3 decode order; the runtimes and architectures
// arrive pre-collected as string slices and are validated at their
// original position.
type LayerVersionCreateInput struct {
	LayerName               string
	Description             string
	LicenseInfo             string
	CompatibleRuntimes      []string
	CompatibleArchitectures []string
	Content                 map[string]interface{}
	Region                  string
}

// LayerPermissionInput carries the wire members of an
// AddLayerVersionPermission request.
type LayerPermissionInput struct {
	LayerName     string
	VersionNumber int64
	StatementId   string
	Action        string
	Principal     string
}

// publishLayerVersionCore publishes a new version of a Lambda layer,
// creating the layer when it does not exist yet, and persists the version's
// code archive. It returns the updated layer and the created version.
func (s *LambdaService) publishLayerVersionCore(ctx context.Context, reqCtx *request.RequestContext, in *LayerVersionCreateInput) (*lambdastore.Layer, *lambdastore.LayerVersion, error) {
	if in.LayerName == "" {
		return nil, nil, NewInvalidParameter("LayerName", "Layer name is required")
	}

	if in.Content == nil {
		return nil, nil, NewInvalidParameter("Content", "Content is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, nil, err
	}
	layers := stores.Layers
	layer, err := layers.Get(in.LayerName)
	if err != nil {
		layer = &lambdastore.Layer{
			LayerName:   in.LayerName,
			CreatedDate: time.Now().UTC(),
		}
		layer, err = layers.Create(layer)
		if err != nil {
			if err == lambdastore.ErrResourceConflict {
				layer, err = layers.Get(in.LayerName)
				if err != nil {
					return nil, nil, err
				}
			} else {
				return nil, nil, err
			}
		}
	}

	version := &lambdastore.LayerVersion{
		Description: in.Description,
		LicenseInfo: in.LicenseInfo,
	}

	for _, c := range in.CompatibleRuntimes {
		if !ValidateRuntime(c) {
			return nil, nil, NewInvalidParameter("CompatibleRuntimes",
				fmt.Sprintf("Runtime '%s' is not supported", c))
		}
		version.CompatibleRuntimes = append(version.CompatibleRuntimes, lambdastore.Runtime(c))
	}

	version.CompatibleArchitectures = append(version.CompatibleArchitectures, in.CompatibleArchitectures...)

	// Decode ZipFile once and reuse for hash/size computation and persistence.
	var decodedZipFile []byte
	if zipFileStr, ok := in.Content["ZipFile"].(string); ok && zipFileStr != "" {
		decodedZipFile, err = base64.StdEncoding.DecodeString(zipFileStr)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid ZipFile encoding: %w", err)
		}
		version.CodeSize = int64(len(decodedZipFile))
		version.CodeSha256 = lambdastore.GenerateCodeHash(decodedZipFile)
	} else if s3Bucket, ok := in.Content["S3Bucket"].(string); ok && s3Bucket != "" {
		s3Key, _ := in.Content["S3Key"].(string)
		if s3Key == "" {
			return nil, nil, NewInvalidParameter("Content.S3Key", "S3Key is required when S3Bucket is specified")
		}
		s3Version, _ := in.Content["S3ObjectVersion"].(string)
		decodedZipFile, err = s.fetchCodeFromS3(ctx, s3Bucket, s3Key, s3Version, in.Region)
		if err != nil {
			return nil, nil, NewInvalidParameter("Content", err.Error())
		}
		version.CodeSize = int64(len(decodedZipFile))
		version.CodeSha256 = lambdastore.GenerateCodeHash(decodedZipFile)
	}

	created, err := layers.PublishVersion(layer, version)
	if err != nil {
		return nil, nil, err
	}

	if decodedZipFile != nil {
		codePath, storeErr := s.storeLayerCode(in.LayerName, created.Version, decodedZipFile, in.Region)
		if storeErr != nil {
			return nil, nil, fmt.Errorf("failed to persist layer code: %w", storeErr)
		}
		created.CodeLocation = codePath
		// Persist the updated CodeLocation so it survives server restarts.
		// created points into layer.Versions, so this call writes the
		// CodeLocation to PebbleDB alongside the rest of the layer.
		if err := layers.Update(layer); err != nil {
			return nil, nil, fmt.Errorf("failed to persist layer code location: %w", err)
		}
	}

	return layer, created, nil
}

// deleteLayerVersionCore deletes a specific version of a layer.
func (s *LambdaService) deleteLayerVersionCore(reqCtx *request.RequestContext, layerName string, versionNumber int64) error {
	if layerName == "" {
		return NewInvalidParameter("LayerName", "Layer name is required")
	}
	if versionNumber <= 0 {
		return NewInvalidParameter("VersionNumber", "Version number is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return err
	}
	if err := stores.Layers.DeleteVersion(layerName, versionNumber); err != nil {
		if errors.Is(err, lambdastore.ErrLayerNotFound) || errors.Is(err, lambdastore.ErrLayerVersionNotFound) {
			return NewResourceNotFound("LayerVersion", layerName)
		}
		return err
	}
	return nil
}

// getLayerVersionCore retrieves a layer and one of its versions.
func (s *LambdaService) getLayerVersionCore(reqCtx *request.RequestContext, layerName string, versionNumber int64) (*lambdastore.Layer, *lambdastore.LayerVersion, error) {
	if layerName == "" {
		return nil, nil, NewInvalidParameter("LayerName", "Layer name is required")
	}
	if versionNumber <= 0 {
		return nil, nil, NewInvalidParameter("VersionNumber", "Version number is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, nil, err
	}
	layer, err := stores.Layers.Get(layerName)
	if err != nil {
		return nil, nil, NewResourceNotFound("LayerVersion", layerName)
	}
	layerVersion, err := stores.Layers.GetVersion(layerName, versionNumber)
	if err != nil {
		return nil, nil, NewResourceNotFound("LayerVersion", layerName)
	}

	return layer, layerVersion, nil
}

// listLayersCore lists the layers in the account, optionally filtered by
// compatible runtime.
func (s *LambdaService) listLayersCore(stores *lambdaStore, compatibleRuntime, marker string, maxItems int) ([]*lambdastore.Layer, string, bool, error) {
	if compatibleRuntime != "" {
		result, err := stores.Layers.ListWithRuntimeFilter(lambdastore.Runtime(compatibleRuntime), storecommon.ListOptions{Marker: marker, MaxItems: maxItems})
		if err != nil {
			return nil, "", false, err
		}
		return result.Items, result.NextMarker, result.IsTruncated, nil
	}
	result, err := stores.Layers.List(storecommon.ListOptions{Marker: marker, MaxItems: maxItems})
	if err != nil {
		return nil, "", false, err
	}
	return result.Items, result.NextMarker, result.IsTruncated, nil
}

// listLayerVersionsCore lists every version of a layer.
func (s *LambdaService) listLayerVersionsCore(reqCtx *request.RequestContext, layerName, marker string, maxItems int) ([]*lambdastore.LayerVersion, string, bool, error) {
	if layerName == "" {
		return nil, "", false, NewInvalidParameter("LayerName", "Layer name is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, "", false, err
	}
	result, err := stores.Layers.ListVersions(layerName, storecommon.ListOptions{Marker: marker, MaxItems: maxItems})
	if err != nil {
		if errors.Is(err, lambdastore.ErrLayerNotFound) {
			return nil, "", false, NewResourceNotFound("Layer", layerName)
		}
		return nil, "", false, err
	}
	return result.Items, result.NextMarker, result.IsTruncated, nil
}

// getLayerVersionByArnCore resolves a layer version by its full ARN and
// also derives the unversioned layer ARN for the response.
func (s *LambdaService) getLayerVersionByArnCore(stores *lambdaStore, layerVersionArn string) (*lambdastore.LayerVersion, string, error) {
	layerVersion, err := stores.Layers.GetVersionByArn(layerVersionArn)
	if err != nil {
		return nil, "", NewResourceNotFound("LayerVersion", layerVersionArn)
	}

	layerArn := ""
	parts := strings.SplitN(layerVersionArn, ":", 7)
	if len(parts) >= 6 {
		layerName := parts[5]
		if layer, err := stores.Layers.Get(layerName); err == nil {
			layerArn = layer.LayerArn
		}
	}

	return layerVersion, layerArn, nil
}

// addLayerVersionPermissionCore adds a permission statement to a layer
// version's resource-based policy and returns the target version for the
// response.
func (s *LambdaService) addLayerVersionPermissionCore(reqCtx *request.RequestContext, in *LayerPermissionInput) (*lambdastore.LayerVersion, error) {
	if in.LayerName == "" {
		return nil, NewInvalidParameter("LayerName", "Layer name is required")
	}

	if in.VersionNumber <= 0 {
		return nil, NewInvalidParameter("VersionNumber", "Version number is required")
	}

	if in.StatementId == "" {
		return nil, NewInvalidParameter("StatementId", "Statement ID is required")
	}
	if err := validateStatementId(in.StatementId); err != nil {
		return nil, err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	layer, err := stores.Layers.Get(in.LayerName)
	if err != nil {
		return nil, NewResourceNotFound("LayerVersion", in.LayerName)
	}

	policy := &lambdastore.LayerPolicy{
		Id:        in.StatementId,
		Action:    in.Action,
		Principal: in.Principal,
	}

	if err := validateLayerPermission(policy); err != nil {
		return nil, err
	}

	if err := stores.Layers.AddPolicy(layer, in.VersionNumber, policy); err != nil {
		if errors.Is(err, lambdastore.ErrLayerNotFound) || errors.Is(err, lambdastore.ErrLayerVersionNotFound) {
			return nil, NewResourceNotFound("LayerVersion", in.LayerName)
		}
		if errors.Is(err, lambdastore.ErrPolicyAlreadyExists) {
			return nil, NewResourceConflict(fmt.Sprintf("StatementId %s already exists", in.StatementId))
		}
		return nil, mapStoreError(err)
	}

	targetVersion, err := stores.Layers.GetVersion(in.LayerName, in.VersionNumber)
	if err != nil {
		return nil, NewResourceNotFound("LayerVersion", in.LayerName)
	}
	if targetVersion == nil {
		return nil, NewResourceNotFound("LayerVersion", in.LayerName)
	}

	return targetVersion, nil
}

// removeLayerVersionPermissionCore removes a permission statement from a
// layer version's resource-based policy.
func (s *LambdaService) removeLayerVersionPermissionCore(reqCtx *request.RequestContext, layerName string, versionNumber int64, statementId string) error {
	if layerName == "" {
		return NewInvalidParameter("LayerName", "Layer name is required")
	}
	if versionNumber <= 0 {
		return NewInvalidParameter("VersionNumber", "Version number is required")
	}
	if statementId == "" {
		return NewInvalidParameter("StatementId", "Statement ID is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return err
	}
	if err := stores.Layers.RemovePolicy(layerName, versionNumber, statementId); err != nil {
		if errors.Is(err, lambdastore.ErrLayerNotFound) || errors.Is(err, lambdastore.ErrLayerVersionNotFound) {
			return NewResourceNotFound("LayerVersion", layerName)
		}
		if errors.Is(err, lambdastore.ErrPolicyNotFound) {
			return NewResourceNotFound("Statement", statementId)
		}
		return err
	}
	return nil
}

// getLayerVersionPolicyCore retrieves a layer version carrying its
// resource-based policy statements. It fails with ResourceNotFound when the
// version carries no policy, matching the AWS contract.
func (s *LambdaService) getLayerVersionPolicyCore(reqCtx *request.RequestContext, layerName string, versionNumber int64) (*lambdastore.LayerVersion, error) {
	if layerName == "" {
		return nil, NewInvalidParameter("LayerName", "Layer name is required")
	}
	if versionNumber <= 0 {
		return nil, NewInvalidParameter("VersionNumber", "Version number is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	layerVersion, err := stores.Layers.GetVersion(layerName, versionNumber)
	if err != nil {
		return nil, NewResourceNotFound("LayerVersion", layerName)
	}

	if len(layerVersion.Policies) == 0 {
		return nil, ErrResourceNotFound
	}

	return layerVersion, nil
}
