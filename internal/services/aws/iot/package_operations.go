package iot

import (
	"context"
	"time"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---------------------------------------------------------------------------
// Package and PackageVersion operations (IoT Device Management — software
// package distribution). Packages contain versions; versions hold an artifact
// reference (S3 location), recipe, attributes, and an optional SBOM blob.
// ---------------------------------------------------------------------------

func (s *IoTService) CreatePackage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "packageName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// Reject duplicate package names.
	exists, err := store.GetGenericExists("iot-package/"+name, &map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, iotstore.ErrPackageAlreadyExists
	}
	now := time.Now().UTC().Unix()
	rec := map[string]interface{}{
		"packageName":      name,
		"packageArn":       iotstore.BuildPackageARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), name),
		"description":      request.GetParamCaseInsensitive(req.Parameters, "description"),
		"creationDate":     now,
		"lastModifiedDate": now,
	}
	if err := store.PutGeneric("iot-package/"+name, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"packageName": name,
		"packageArn":  rec["packageArn"],
		"description": rec["description"],
	}, nil
}

func (s *IoTService) GetPackage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "packageName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("iot-package/"+name, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrPackageNotFound
	}
	return rec, nil
}

func (s *IoTService) UpdatePackage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "packageName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("iot-package/"+name, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrPackageNotFound
	}
	if desc := request.GetParamCaseInsensitive(req.Parameters, "description"); desc != "" {
		rec["description"] = desc
	}
	if dv := request.GetParamCaseInsensitive(req.Parameters, "defaultVersionName"); dv != "" {
		rec["defaultVersionName"] = dv
	}
	rec["lastModifiedDate"] = time.Now().UTC().Unix()
	if err := store.PutGeneric("iot-package/"+name, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) DeletePackage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "packageName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	exists, err := store.GetGenericExists("iot-package/"+name, &map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrPackageNotFound
	}
	// Cascade delete all versions under this package.
	prefix := "iot-package-version/" + name + "/"
	versions, err := store.ListGeneric(prefix)
	if err != nil {
		return nil, err
	}
	for _, v := range versions {
		if vn, ok := v["versionName"].(string); ok && vn != "" {
			if err := store.DeleteGeneric(prefix + vn); err != nil {
				return nil, err
			}
		}
	}
	arn := iotstore.BuildPackageARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), name)
	_ = store.DeleteAllTags(arn)

	if err := store.DeleteGeneric("iot-package/" + name); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) ListPackages(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := store.ListGeneric("iot-package/")
	if err != nil {
		return nil, err
	}
	summaries := make([]map[string]interface{}, 0, len(items))
	for _, rec := range items {
		summaries = append(summaries, map[string]interface{}{
			"packageName":      rec["packageName"],
			"packageArn":       rec["packageArn"],
			"description":      rec["description"],
			"creationDate":     rec["creationDate"],
			"lastModifiedDate": rec["lastModifiedDate"],
		})
	}
	return paginatedMaps("packageSummaries", summaries, req.Parameters), nil
}

// --- PackageVersion ---

func (s *IoTService) CreatePackageVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	pkgName := request.GetParamCaseInsensitive(req.Parameters, "packageName")
	versionName := request.GetParamCaseInsensitive(req.Parameters, "versionName")
	if pkgName == "" || versionName == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// Validate parent package exists.
	pkgExists, err := store.GetGenericExists("iot-package/"+pkgName, &map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	if !pkgExists {
		return nil, iotstore.ErrPackageNotFound
	}
	// Reject duplicate version names.
	verExists, err := store.GetGenericExists("iot-package-version/"+pkgName+"/"+versionName, &map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	if verExists {
		return nil, iotstore.ErrPackageVersionAlreadyExists
	}
	now := time.Now().UTC().Unix()
	rec := map[string]interface{}{
		"packageName":       pkgName,
		"versionName":       versionName,
		"packageVersionArn": iotstore.BuildPackageARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), pkgName) + ":version/" + versionName,
		"description":       request.GetParamCaseInsensitive(req.Parameters, "description"),
		"attributes":        req.Parameters["attributes"],
		"artifact":          req.Parameters["artifact"],
		"recipe":            request.GetParamCaseInsensitive(req.Parameters, "recipe"),
		"status":            "DRAFT",
		"errorReason":       "",
		"creationDate":      now,
		"lastModifiedDate":  now,
	}
	if err := store.PutGeneric("iot-package-version/"+pkgName+"/"+versionName, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"packageVersionArn": rec["packageVersionArn"],
		"packageName":       pkgName,
		"versionName":       versionName,
		"description":       rec["description"],
		"attributes":        rec["attributes"],
		"status":            rec["status"],
		"errorReason":       rec["errorReason"],
	}, nil
}

func (s *IoTService) GetPackageVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	pkgName := request.GetParamCaseInsensitive(req.Parameters, "packageName")
	versionName := request.GetParamCaseInsensitive(req.Parameters, "versionName")
	if pkgName == "" || versionName == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("iot-package-version/"+pkgName+"/"+versionName, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrPackageVersionNotFound
	}
	return rec, nil
}

func (s *IoTService) UpdatePackageVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	pkgName := request.GetParamCaseInsensitive(req.Parameters, "packageName")
	versionName := request.GetParamCaseInsensitive(req.Parameters, "versionName")
	if pkgName == "" || versionName == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key := "iot-package-version/" + pkgName + "/" + versionName
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrPackageVersionNotFound
	}
	if desc := request.GetParamCaseInsensitive(req.Parameters, "description"); desc != "" {
		rec["description"] = desc
	}
	if attr, ok := req.Parameters["attributes"]; ok {
		rec["attributes"] = attr
	}
	if artifact, ok := req.Parameters["artifact"]; ok {
		rec["artifact"] = artifact
	}
	if recipe := request.GetParamCaseInsensitive(req.Parameters, "recipe"); recipe != "" {
		rec["recipe"] = recipe
	}
	// Handle lifecycle action: PUBLISH or DEPRECATE.
	switch request.GetParamCaseInsensitive(req.Parameters, "action") {
	case "PUBLISH":
		rec["status"] = "PUBLISHED"
	case "DEPRECATE":
		rec["status"] = "DEPRECATED"
	}
	rec["lastModifiedDate"] = time.Now().UTC().Unix()
	if err := store.PutGeneric(key, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) DeletePackageVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	pkgName := request.GetParamCaseInsensitive(req.Parameters, "packageName")
	versionName := request.GetParamCaseInsensitive(req.Parameters, "versionName")
	if pkgName == "" || versionName == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key := "iot-package-version/" + pkgName + "/" + versionName
	exists, err := store.GetGenericExists(key, &map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrPackageVersionNotFound
	}
	if err := store.DeleteGeneric(key); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) ListPackageVersions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	pkgName := request.GetParamCaseInsensitive(req.Parameters, "packageName")
	if pkgName == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := store.ListGeneric("iot-package-version/" + pkgName + "/")
	if err != nil {
		return nil, err
	}
	statusFilter := request.GetParamCaseInsensitive(req.Parameters, "status")
	summaries := make([]map[string]interface{}, 0, len(items))
	for _, rec := range items {
		if statusFilter != "" {
			if st, _ := rec["status"].(string); st != statusFilter {
				continue
			}
		}
		summaries = append(summaries, map[string]interface{}{
			"packageName":  rec["packageName"],
			"versionName":  rec["versionName"],
			"status":       rec["status"],
			"creationDate": rec["creationDate"],
		})
	}
	return paginatedMaps("packageVersionSummaries", summaries, req.Parameters), nil
}

// --- PackageConfiguration ---

func (s *IoTService) GetPackageConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	_, err = store.GetGenericExists("iot-package-config", &rec)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func (s *IoTService) UpdatePackageConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	_, err = store.GetGenericExists("iot-package-config", &rec)
	if err != nil {
		return nil, err
	}
	if vj, ok := req.Parameters["versionUpdateByJobsConfig"]; ok {
		rec["versionUpdateByJobsConfig"] = vj
	}
	if err := store.PutGeneric("iot-package-config", rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

// --- SBOM ---

func (s *IoTService) AssociateSbomWithPackageVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	pkgName := request.GetParamCaseInsensitive(req.Parameters, "packageName")
	versionName := request.GetParamCaseInsensitive(req.Parameters, "versionName")
	if pkgName == "" || versionName == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key := "iot-package-version/" + pkgName + "/" + versionName
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrPackageVersionNotFound
	}
	rec["sbom"] = req.Parameters["sbom"]
	rec["sbomValidationStatus"] = "SUCCEEDED"
	rec["lastModifiedDate"] = time.Now().UTC().Unix()
	if err := store.PutGeneric(key, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"packageName":          pkgName,
		"versionName":          versionName,
		"sbom":                 rec["sbom"],
		"sbomValidationStatus": "SUCCEEDED",
	}, nil
}

func (s *IoTService) DisassociateSbomFromPackageVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	pkgName := request.GetParamCaseInsensitive(req.Parameters, "packageName")
	versionName := request.GetParamCaseInsensitive(req.Parameters, "versionName")
	if pkgName == "" || versionName == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key := "iot-package-version/" + pkgName + "/" + versionName
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrPackageVersionNotFound
	}
	delete(rec, "sbom")
	delete(rec, "sbomValidationStatus")
	rec["lastModifiedDate"] = time.Now().UTC().Unix()
	if err := store.PutGeneric(key, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) ListSbomValidationResults(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// SBOM validation requires an external validator service which is not
	// part of the on-prem platform. Return an empty list per AWS behaviour
	// when no validation results exist.
	return paginatedMaps("validationResultSummaries", []map[string]interface{}{}, req.Parameters), nil
}
