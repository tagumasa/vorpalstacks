package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// ---------------------------------------------------------------------------
// Package and PackageVersion operations (IoT Device Management — software
// package distribution). Packages contain versions; versions hold an artifact
// reference (S3 location), recipe, attributes, and an optional SBOM blob.
// ---------------------------------------------------------------------------

func (s *IoTService) CreatePackage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.createPackageCore(store, CreatePackageInput{
		PackageName: request.GetParamCaseInsensitive(req.Parameters, "packageName"),
		Description: request.GetParamCaseInsensitive(req.Parameters, "description"),
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"packageName": result.PackageName,
		"packageArn":  result.PackageArn,
		"description": result.Description,
	}, nil
}

func (s *IoTService) GetPackage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.getPackageCore(store, request.GetParamCaseInsensitive(req.Parameters, "packageName"))
}

func (s *IoTService) UpdatePackage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.updatePackageCore(store, UpdatePackageInput{
		PackageName:         request.GetParamCaseInsensitive(req.Parameters, "packageName"),
		Description:         request.GetParamCaseInsensitive(req.Parameters, "description"),
		DefaultVersionName:  request.GetParamCaseInsensitive(req.Parameters, "defaultVersionName"),
		UnsetDefaultVersion: request.GetBoolParam(req.Parameters, "unsetDefaultVersion"),
	}); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) DeletePackage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deletePackageCore(store, request.GetParamCaseInsensitive(req.Parameters, "packageName")); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) ListPackages(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	summaries, err := s.listPackagesCore(store)
	if err != nil {
		return nil, err
	}
	return paginatedMaps("packageSummaries", summaries, req.Parameters)
}

// --- PackageVersion ---

func (s *IoTService) CreatePackageVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.createPackageVersionCore(store, CreatePackageVersionInput{
		PackageName: request.GetParamCaseInsensitive(req.Parameters, "packageName"),
		VersionName: request.GetParamCaseInsensitive(req.Parameters, "versionName"),
		Description: request.GetParamCaseInsensitive(req.Parameters, "description"),
		Attributes:  req.Parameters["attributes"],
		Artifact:    req.Parameters["artifact"],
		Recipe:      request.GetParamCaseInsensitive(req.Parameters, "recipe"),
	})
}

func (s *IoTService) GetPackageVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.getPackageVersionCore(store,
		request.GetParamCaseInsensitive(req.Parameters, "packageName"),
		request.GetParamCaseInsensitive(req.Parameters, "versionName"))
}

func (s *IoTService) UpdatePackageVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	attributes, attributesProvided := req.Parameters["attributes"]
	artifact, artifactProvided := req.Parameters["artifact"]
	if err := s.updatePackageVersionCore(store, UpdatePackageVersionInput{
		PackageName:        request.GetParamCaseInsensitive(req.Parameters, "packageName"),
		VersionName:        request.GetParamCaseInsensitive(req.Parameters, "versionName"),
		Description:        request.GetParamCaseInsensitive(req.Parameters, "description"),
		Attributes:         attributes,
		AttributesProvided: attributesProvided,
		Artifact:           artifact,
		ArtifactProvided:   artifactProvided,
		Recipe:             request.GetParamCaseInsensitive(req.Parameters, "recipe"),
		Action:             request.GetParamCaseInsensitive(req.Parameters, "action"),
	}); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) DeletePackageVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deletePackageVersionCore(store,
		request.GetParamCaseInsensitive(req.Parameters, "packageName"),
		request.GetParamCaseInsensitive(req.Parameters, "versionName")); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) ListPackageVersions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	summaries, err := s.listPackageVersionsCore(store,
		request.GetParamCaseInsensitive(req.Parameters, "packageName"),
		request.GetParamCaseInsensitive(req.Parameters, "status"))
	if err != nil {
		return nil, err
	}
	return paginatedMaps("packageVersionSummaries", summaries, req.Parameters)
}

// --- PackageConfiguration ---

func (s *IoTService) GetPackageConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.getPackageConfigurationCore(store)
}

func (s *IoTService) UpdatePackageConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	config, configProvided := req.Parameters["versionUpdateByJobsConfig"]
	if err := s.updatePackageConfigurationCore(store, UpdatePackageConfigurationInput{
		VersionUpdateByJobsConfig: config,
		ConfigProvided:            configProvided,
	}); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

// --- SBOM ---

func (s *IoTService) AssociateSbomWithPackageVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.associateSbomCore(store, AssociateSbomInput{
		PackageName: request.GetParamCaseInsensitive(req.Parameters, "packageName"),
		VersionName: request.GetParamCaseInsensitive(req.Parameters, "versionName"),
		Sbom:        req.Parameters["sbom"],
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"packageName":          result.PackageName,
		"versionName":          result.VersionName,
		"sbom":                 result.Sbom,
		"sbomValidationStatus": result.SbomValidationStatus,
	}, nil
}

func (s *IoTService) DisassociateSbomFromPackageVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.disassociateSbomCore(store,
		request.GetParamCaseInsensitive(req.Parameters, "packageName"),
		request.GetParamCaseInsensitive(req.Parameters, "versionName")); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) ListSbomValidationResults(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	summaries, err := s.listSbomValidationResultsCore(
		request.GetParamCaseInsensitive(req.Parameters, "packageName"),
		request.GetParamCaseInsensitive(req.Parameters, "versionName"))
	if err != nil {
		return nil, err
	}
	return paginatedMaps("validationResultSummaries", summaries, req.Parameters)
}
