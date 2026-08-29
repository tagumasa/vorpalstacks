package iot

import (
	"time"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---------------------------------------------------------------------------
// Package and PackageVersion Core (IoT Device Management — software package
// distribution). Packages contain versions; versions hold an artifact
// reference (S3 location), recipe, attributes, and an optional SBOM blob.
// Records live under "iot-package/<name>" and
// "iot-package-version/<pkg>/<version>" keys; the account-level package
// configuration lives under the "iot-package-config" key.
// ---------------------------------------------------------------------------

// CreatePackageInput carries the fields for CreatePackage.
type CreatePackageInput struct {
	PackageName string
	Description string
}

// CreatePackageResult is the transport-agnostic result of CreatePackage.
type CreatePackageResult struct {
	PackageName string
	PackageArn  string
	Description string
}

// UpdatePackageInput carries the fields for UpdatePackage. Empty values
// leave the stored fields untouched.
type UpdatePackageInput struct {
	PackageName         string
	Description         string
	DefaultVersionName  string
	UnsetDefaultVersion bool
}

// CreatePackageVersionInput carries the fields for CreatePackageVersion.
// Attributes and Artifact keep the raw wire values (nested structures).
type CreatePackageVersionInput struct {
	PackageName string
	VersionName string
	Description string
	Attributes  interface{}
	Artifact    interface{}
	Recipe      string
}

// UpdatePackageVersionInput carries the fields for UpdatePackageVersion.
// AttributesProvided/ArtifactProvided distinguish an explicitly supplied
// (possibly null-clearing) member from an omitted one; the lifecycle action
// is PUBLISH or DEPRECATE.
type UpdatePackageVersionInput struct {
	PackageName        string
	VersionName        string
	Description        string
	Attributes         interface{}
	AttributesProvided bool
	Artifact           interface{}
	ArtifactProvided   bool
	Recipe             string
	Action             string
}

// AssociateSbomInput carries the fields for the SBOM association operations.
type AssociateSbomInput struct {
	PackageName string
	VersionName string
	Sbom        interface{}
}

// AssociateSbomResult is the transport-agnostic result of
// AssociateSbomWithPackageVersion.
type AssociateSbomResult struct {
	PackageName          string
	VersionName          string
	Sbom                 interface{}
	SbomValidationStatus string
}

// UpdatePackageConfigurationInput carries the fields for
// UpdatePackageConfiguration.
type UpdatePackageConfigurationInput struct {
	VersionUpdateByJobsConfig interface{}
	ConfigProvided            bool
}

// createPackageCore validates and persists a package record, rejecting
// duplicate names.
func (s *IoTService) createPackageCore(store iotstore.IotStoreInterface, in CreatePackageInput) (*CreatePackageResult, error) {
	if in.PackageName == "" {
		return nil, iotstore.ErrMissingParam
	}
	// Reject duplicate package names.
	exists, err := store.GetGenericExists("iot-package/"+in.PackageName, &map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, iotstore.ErrPackageAlreadyExists
	}
	now := time.Now().UTC().Unix()
	arn := iotstore.BuildPackageARN(store.GetAccountID(), store.GetRegion(), in.PackageName)
	rec := map[string]interface{}{
		"packageName":      in.PackageName,
		"packageArn":       arn,
		"description":      in.Description,
		"creationDate":     now,
		"lastModifiedDate": now,
	}
	if err := store.PutGeneric("iot-package/"+in.PackageName, rec); err != nil {
		return nil, err
	}
	return &CreatePackageResult{
		PackageName: in.PackageName,
		PackageArn:  arn,
		Description: in.Description,
	}, nil
}

// getPackageCore retrieves a package record.
func (s *IoTService) getPackageCore(store iotstore.IotStoreInterface, name string) (map[string]interface{}, error) {
	if name == "" {
		return nil, iotstore.ErrMissingParam
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

// updatePackageCore applies the supplied fields to an existing package.
// Setting and unsetting the default version in one request is rejected;
// unsetDefaultVersion alone clears the stored default.
func (s *IoTService) updatePackageCore(store iotstore.IotStoreInterface, in UpdatePackageInput) error {
	if in.PackageName == "" {
		return iotstore.ErrMissingParam
	}
	if in.UnsetDefaultVersion && in.DefaultVersionName != "" {
		return iotstore.ErrInvalidRequest
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("iot-package/"+in.PackageName, &rec)
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrPackageNotFound
	}
	if in.Description != "" {
		rec["description"] = in.Description
	}
	if in.DefaultVersionName != "" {
		rec["defaultVersionName"] = in.DefaultVersionName
	}
	if in.UnsetDefaultVersion {
		delete(rec, "defaultVersionName")
	}
	rec["lastModifiedDate"] = time.Now().UTC().Unix()
	return store.PutGeneric("iot-package/"+in.PackageName, rec)
}

// listSbomValidationResultsCore enforces the model-required packageName and
// versionName path labels. SBOM validation requires an external validator
// service which is not part of the on-prem platform, so no validation
// results ever exist and the summary list stays empty — per AWS behaviour
// when no validation results exist.
func (s *IoTService) listSbomValidationResultsCore(packageName, versionName string) ([]map[string]interface{}, error) {
	if packageName == "" || versionName == "" {
		return nil, iotstore.ErrValidation
	}
	return []map[string]interface{}{}, nil
}

// deletePackageCore removes a package, cascading to its versions and tags.
func (s *IoTService) deletePackageCore(store iotstore.IotStoreInterface, name string) error {
	if name == "" {
		return iotstore.ErrMissingParam
	}
	exists, err := store.GetGenericExists("iot-package/"+name, &map[string]interface{}{})
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrPackageNotFound
	}
	// Cascade delete all versions under this package.
	prefix := "iot-package-version/" + name + "/"
	versions, err := store.ListGeneric(prefix)
	if err != nil {
		return err
	}
	for _, v := range versions {
		if vn, ok := v["versionName"].(string); ok && vn != "" {
			if err := store.DeleteGeneric(prefix + vn); err != nil {
				return err
			}
		}
	}
	arn := iotstore.BuildPackageARN(store.GetAccountID(), store.GetRegion(), name)
	_ = store.DeleteAllTags(arn)

	return store.DeleteGeneric("iot-package/" + name)
}

// listPackagesCore lists package summaries.
func (s *IoTService) listPackagesCore(store iotstore.IotStoreInterface) ([]map[string]interface{}, error) {
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
	return summaries, nil
}

// createPackageVersionCore validates and persists a package-version record.
func (s *IoTService) createPackageVersionCore(store iotstore.IotStoreInterface, in CreatePackageVersionInput) (map[string]interface{}, error) {
	if in.PackageName == "" || in.VersionName == "" {
		return nil, iotstore.ErrMissingParam
	}
	// Validate parent package exists.
	pkgExists, err := store.GetGenericExists("iot-package/"+in.PackageName, &map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	if !pkgExists {
		return nil, iotstore.ErrPackageNotFound
	}
	// Reject duplicate version names.
	key := "iot-package-version/" + in.PackageName + "/" + in.VersionName
	verExists, err := store.GetGenericExists(key, &map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	if verExists {
		return nil, iotstore.ErrPackageVersionAlreadyExists
	}
	now := time.Now().UTC().Unix()
	rec := map[string]interface{}{
		"packageName":       in.PackageName,
		"versionName":       in.VersionName,
		"packageVersionArn": iotstore.BuildPackageARN(store.GetAccountID(), store.GetRegion(), in.PackageName) + ":version/" + in.VersionName,
		"description":       in.Description,
		"attributes":        in.Attributes,
		"artifact":          in.Artifact,
		"recipe":            in.Recipe,
		"status":            "DRAFT",
		"errorReason":       "",
		"creationDate":      now,
		"lastModifiedDate":  now,
	}
	if err := store.PutGeneric(key, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"packageVersionArn": rec["packageVersionArn"],
		"packageName":       in.PackageName,
		"versionName":       in.VersionName,
		"description":       rec["description"],
		"attributes":        rec["attributes"],
		"status":            rec["status"],
		"errorReason":       rec["errorReason"],
	}, nil
}

// getPackageVersionCore retrieves a package-version record.
func (s *IoTService) getPackageVersionCore(store iotstore.IotStoreInterface, pkgName, versionName string) (map[string]interface{}, error) {
	if pkgName == "" || versionName == "" {
		return nil, iotstore.ErrMissingParam
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

// updatePackageVersionCore applies the supplied fields and lifecycle action
// to an existing package version.
func (s *IoTService) updatePackageVersionCore(store iotstore.IotStoreInterface, in UpdatePackageVersionInput) error {
	if in.PackageName == "" || in.VersionName == "" {
		return iotstore.ErrMissingParam
	}
	key := "iot-package-version/" + in.PackageName + "/" + in.VersionName
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrPackageVersionNotFound
	}
	if in.Description != "" {
		rec["description"] = in.Description
	}
	if in.AttributesProvided {
		rec["attributes"] = in.Attributes
	}
	if in.ArtifactProvided {
		rec["artifact"] = in.Artifact
	}
	if in.Recipe != "" {
		rec["recipe"] = in.Recipe
	}
	// Handle lifecycle action: PUBLISH or DEPRECATE.
	switch in.Action {
	case "PUBLISH":
		rec["status"] = "PUBLISHED"
	case "DEPRECATE":
		rec["status"] = "DEPRECATED"
	}
	rec["lastModifiedDate"] = time.Now().UTC().Unix()
	return store.PutGeneric(key, rec)
}

// deletePackageVersionCore removes a package-version record.
func (s *IoTService) deletePackageVersionCore(store iotstore.IotStoreInterface, pkgName, versionName string) error {
	if pkgName == "" || versionName == "" {
		return iotstore.ErrMissingParam
	}
	key := "iot-package-version/" + pkgName + "/" + versionName
	exists, err := store.GetGenericExists(key, &map[string]interface{}{})
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrPackageVersionNotFound
	}
	return store.DeleteGeneric(key)
}

// listPackageVersionsCore lists the version summaries of a package with an
// optional status filter.
func (s *IoTService) listPackageVersionsCore(store iotstore.IotStoreInterface, pkgName, statusFilter string) ([]map[string]interface{}, error) {
	if pkgName == "" {
		return nil, iotstore.ErrMissingParam
	}
	items, err := store.ListGeneric("iot-package-version/" + pkgName + "/")
	if err != nil {
		return nil, err
	}
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
	return summaries, nil
}

// getPackageConfigurationCore retrieves the account-level package
// configuration record.
func (s *IoTService) getPackageConfigurationCore(store iotstore.IotStoreInterface) (map[string]interface{}, error) {
	rec := map[string]interface{}{}
	_, err := store.GetGenericExists("iot-package-config", &rec)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

// updatePackageConfigurationCore applies the supplied configuration members
// to the account-level package configuration record.
func (s *IoTService) updatePackageConfigurationCore(store iotstore.IotStoreInterface, in UpdatePackageConfigurationInput) error {
	rec := map[string]interface{}{}
	_, err := store.GetGenericExists("iot-package-config", &rec)
	if err != nil {
		return err
	}
	if in.ConfigProvided {
		rec["versionUpdateByJobsConfig"] = in.VersionUpdateByJobsConfig
	}
	return store.PutGeneric("iot-package-config", rec)
}

// associateSbomCore attaches an SBOM blob to a package version.
func (s *IoTService) associateSbomCore(store iotstore.IotStoreInterface, in AssociateSbomInput) (*AssociateSbomResult, error) {
	if in.PackageName == "" || in.VersionName == "" {
		return nil, iotstore.ErrMissingParam
	}
	key := "iot-package-version/" + in.PackageName + "/" + in.VersionName
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrPackageVersionNotFound
	}
	rec["sbom"] = in.Sbom
	rec["sbomValidationStatus"] = "SUCCEEDED"
	rec["lastModifiedDate"] = time.Now().UTC().Unix()
	if err := store.PutGeneric(key, rec); err != nil {
		return nil, err
	}
	return &AssociateSbomResult{
		PackageName:          in.PackageName,
		VersionName:          in.VersionName,
		Sbom:                 rec["sbom"],
		SbomValidationStatus: "SUCCEEDED",
	}, nil
}

// disassociateSbomCore removes the SBOM blob from a package version.
func (s *IoTService) disassociateSbomCore(store iotstore.IotStoreInterface, pkgName, versionName string) error {
	if pkgName == "" || versionName == "" {
		return iotstore.ErrMissingParam
	}
	key := "iot-package-version/" + pkgName + "/" + versionName
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrPackageVersionNotFound
	}
	delete(rec, "sbom")
	delete(rec, "sbomValidationStatus")
	rec["lastModifiedDate"] = time.Now().UTC().Unix()
	return store.PutGeneric(key, rec)
}
