package iot

import (
	"strings"
	"time"

	storecommon "vorpalstacks/internal/store/aws/common"
	iotstore "vorpalstacks/internal/store/aws/iot"
	"vorpalstacks/internal/utils/aws/arn"
)

// ---------------------------------------------------------------------------
// Security Profile Core (Device Defender security profiles, their target
// attachments and violation records). Profile definitions use the typed
// SecurityProfile store; attachments persist as bidirectional generic-KV
// mappings so ListSecurityProfilesForTarget (target->profile) and
// ListTargetsForSecurityProfile (profile->target) can scan a single
// prefix; violations use the typed ViolationEvent store.
// ---------------------------------------------------------------------------

// CreateSecurityProfileInput carries the parsed CreateSecurityProfile
// request. The *Malformed flags preserve the parse-error error precedence
// (behaviors, then alertTargets) that the wire parsing detects.
type CreateSecurityProfileInput struct {
	Name                        string
	Description                 string
	Behaviors                   []*iotstore.Behavior
	BehaviorsMalformed          bool
	AlertTargets                map[string]*iotstore.AlertTarget
	AlertTargetsMalformed       bool
	AdditionalMetricsToRetainV2 []*iotstore.MetricToRetain
	AdditionalMetricsToRetain   []string
	MetricsExportConfig         string
	Tags                        map[string]string
}

// dimensionValueOperators is the DimensionValueOperator enum member set.
var dimensionValueOperators = map[string]bool{"IN": true, "NOT_IN": true}

// validateMetricsToRetain enforces the model's required metric member and
// the DimensionValueOperator enum on every retained-metric entry (an
// omitted operator is valid — the model documents an implicit IN).
func validateMetricsToRetain(entries []*iotstore.MetricToRetain) error {
	for _, m := range entries {
		if m.Metric == "" {
			return iotstore.ErrInvalidRequest
		}
		if m.Operator != "" && !dimensionValueOperators[m.Operator] {
			return iotstore.ErrInvalidRequest
		}
	}
	return nil
}

// createSecurityProfileCore validates and persists a security profile and
// returns the created record.
func (s *IoTService) createSecurityProfileCore(store iotstore.IotStoreInterface, in CreateSecurityProfileInput) (*iotstore.SecurityProfile, error) {
	if in.Name == "" {
		return nil, iotstore.ErrMissingParam
	}
	if in.BehaviorsMalformed {
		return nil, iotstore.ErrInvalidRequest
	}
	if in.AlertTargetsMalformed {
		return nil, iotstore.ErrInvalidRequest
	}
	if err := validateMetricsToRetain(in.AdditionalMetricsToRetainV2); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	sp := &iotstore.SecurityProfile{
		SecurityProfileName:         in.Name,
		SecurityProfileDescription:  in.Description,
		Behaviors:                   in.Behaviors,
		AlertTargets:                in.AlertTargets,
		AdditionalMetricsToRetainV2: in.AdditionalMetricsToRetainV2,
		AdditionalMetricsToRetain:   in.AdditionalMetricsToRetain,
		MetricsExportConfig:         in.MetricsExportConfig,
		Tags:                        in.Tags,
		Version:                     1,
		CreationDate:                now,
		LastModifiedDate:            now,
	}
	created, err := store.CreateSecurityProfile(sp)
	if err != nil {
		return nil, err
	}
	// Create-time tags live in the ARN-keyed tag store so
	// ListTagsForResource observes them; the delete path already clears
	// them with DeleteAllTags.
	if len(in.Tags) > 0 {
		arn := iotstore.BuildSecurityProfileARN(store.GetAccountID(), store.GetRegion(), in.Name)
		if err := store.TagResource(arn, in.Tags); err != nil {
			return nil, err
		}
	}
	return created, nil
}

// describeSecurityProfileCore loads a security profile.
func (s *IoTService) describeSecurityProfileCore(store iotstore.IotStoreInterface, name string) (*iotstore.SecurityProfile, error) {
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}
	return store.GetSecurityProfile(name)
}

// UpdateSecurityProfileInput carries the parsed UpdateSecurityProfile
// request. The *Malformed flags preserve the parse-error error precedence;
// nil/empty optional members mean "leave unchanged" and the Delete* flags
// follow the model contract (a delete flag plus a replacement value in the
// same invocation is rejected).
type UpdateSecurityProfileInput struct {
	Name                            string
	ExpectedVersion                 int64
	Description                     string
	Behaviors                       []*iotstore.Behavior
	BehaviorsMalformed              bool
	DeleteBehaviors                 bool
	AlertTargets                    map[string]*iotstore.AlertTarget
	AlertTargetsMalformed           bool
	DeleteAlertTargets              bool
	AdditionalMetricsToRetainV2     []*iotstore.MetricToRetain
	AdditionalMetricsToRetain       []string
	DeleteAdditionalMetricsToRetain bool
	MetricsExportConfig             string
	DeleteMetricsExportConfig       bool
}

// updateSecurityProfileCore applies a partial update to an existing
// security profile, bumps its version and returns the updated record.
func (s *IoTService) updateSecurityProfileCore(store iotstore.IotStoreInterface, in UpdateSecurityProfileInput) (*iotstore.SecurityProfile, error) {
	if in.Name == "" {
		return nil, iotstore.ErrMissingParam
	}
	if in.BehaviorsMalformed {
		return nil, iotstore.ErrInvalidRequest
	}
	if in.AlertTargetsMalformed {
		return nil, iotstore.ErrInvalidRequest
	}
	if err := validateMetricsToRetain(in.AdditionalMetricsToRetainV2); err != nil {
		return nil, err
	}
	if in.DeleteBehaviors && in.Behaviors != nil {
		return nil, iotstore.ErrInvalidRequest
	}
	if in.DeleteAlertTargets && in.AlertTargets != nil {
		return nil, iotstore.ErrInvalidRequest
	}
	if in.DeleteAdditionalMetricsToRetain && (in.AdditionalMetricsToRetain != nil || in.AdditionalMetricsToRetainV2 != nil) {
		return nil, iotstore.ErrInvalidRequest
	}
	if in.DeleteMetricsExportConfig && in.MetricsExportConfig != "" {
		return nil, iotstore.ErrInvalidRequest
	}
	existing, err := store.GetSecurityProfile(in.Name)
	if err != nil {
		return nil, err
	}
	if in.ExpectedVersion > 0 && in.ExpectedVersion != existing.Version {
		return nil, iotstore.ErrVersionConflict
	}
	if in.Description != "" {
		existing.SecurityProfileDescription = in.Description
	}
	if in.DeleteBehaviors {
		existing.Behaviors = nil
	} else if in.Behaviors != nil {
		existing.Behaviors = in.Behaviors
	}
	if in.DeleteAlertTargets {
		existing.AlertTargets = nil
	} else if in.AlertTargets != nil {
		existing.AlertTargets = in.AlertTargets
	}
	if in.AdditionalMetricsToRetainV2 != nil {
		existing.AdditionalMetricsToRetainV2 = in.AdditionalMetricsToRetainV2
	}
	if in.DeleteAdditionalMetricsToRetain {
		existing.AdditionalMetricsToRetain = nil
	} else if in.AdditionalMetricsToRetain != nil {
		existing.AdditionalMetricsToRetain = in.AdditionalMetricsToRetain
	}
	if in.DeleteMetricsExportConfig {
		existing.MetricsExportConfig = ""
	} else if in.MetricsExportConfig != "" {
		existing.MetricsExportConfig = in.MetricsExportConfig
	}
	existing.Version++
	existing.LastModifiedDate = time.Now().UTC()
	if err := store.UpdateSecurityProfile(in.Name, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// deleteSecurityProfileCore removes a security profile and its tags.
func (s *IoTService) deleteSecurityProfileCore(store iotstore.IotStoreInterface, name string) error {
	if name == "" {
		return iotstore.ErrMissingParam
	}
	arn := iotstore.BuildSecurityProfileARN(store.GetAccountID(), store.GetRegion(), name)
	_ = store.DeleteAllTags(arn)
	return store.DeleteSecurityProfile(name)
}

// SecurityProfileListItem is one ListSecurityProfiles entry.
type SecurityProfileListItem struct {
	Name string
	Arn  string
}

// listSecurityProfilesCore lists security profiles with their identifiers.
func (s *IoTService) listSecurityProfilesCore(store iotstore.IotStoreInterface, opts storecommon.ListOptions) ([]SecurityProfileListItem, string, error) {
	profiles, err := store.ListSecurityProfiles(opts)
	if err != nil {
		return nil, "", err
	}
	items := make([]SecurityProfileListItem, 0, len(profiles.Items))
	for _, sp := range profiles.Items {
		items = append(items, SecurityProfileListItem{
			Name: sp.SecurityProfileName,
			Arn:  sp.SecurityProfileARN,
		})
	}
	return items, profiles.NextMarker, nil
}

// attachSecurityProfileCore records a profile<->target association. Both
// endpoints must exist: an unknown profile name or thing target raises
// ResourceNotFoundException. Forward mapping: profile -> target. Reverse
// mapping: target -> profile. Both are stored so that
// ListSecurityProfilesForTarget (target->profile) and
// ListTargetsForSecurityProfile (profile->target) can scan a single prefix.
func (s *IoTService) attachSecurityProfileCore(store iotstore.IotStoreInterface, profileName, targetArn string) error {
	if profileName == "" || targetArn == "" {
		return iotstore.ErrMissingParam
	}
	if _, err := store.GetSecurityProfile(profileName); err != nil {
		return err
	}
	if err := validateSecurityProfileTarget(store, targetArn); err != nil {
		return err
	}
	forwardKey := "secProfileTarget/" + profileName + "/" + targetArn
	reverseKey := "secTargetProfile/" + targetArn + "/" + profileName
	assocValue := map[string]interface{}{
		"securityProfileName":      profileName,
		"securityProfileTargetArn": targetArn,
	}
	if err := store.PutGeneric(forwardKey, assocValue); err != nil {
		return err
	}
	if err := store.PutGeneric(reverseKey, assocValue); err != nil {
		// Rollback forward write to maintain bidirectional consistency.
		_ = store.DeleteGeneric(forwardKey)
		return err
	}
	return nil
}

// validateSecurityProfileTarget checks the attach target ARN. Thing
// targets must reference an existing thing; account-scoped targets address
// the account itself and need no separate record. Other IoT resources
// (thing groups, certificate targets) follow once their attachment
// semantics exist; their ARN form is accepted.
func validateSecurityProfileTarget(store iotstore.IotStoreInterface, targetArn string) error {
	parsed, err := arn.ParseARN(targetArn)
	if err != nil {
		return iotstore.ErrInvalidRequest
	}
	if parsed.Service != "iot" {
		return iotstore.ErrInvalidRequest
	}
	// The resource is "thing/<name>", "account" or another IoT resource
	// form; only the thing form maps to a checkable record.
	if name, ok := strings.CutPrefix(parsed.Resource, "thing/"); ok && name != "" {
		if _, err := store.GetThing(name); err != nil {
			return err
		}
	}
	return nil
}

// detachSecurityProfileCore removes a profile<->target association. An
// unknown association yields ErrSecurityProfileAttachmentNotFound.
func (s *IoTService) detachSecurityProfileCore(store iotstore.IotStoreInterface, profileName, targetArn string) error {
	if profileName == "" || targetArn == "" {
		return iotstore.ErrMissingParam
	}
	forwardKey := "secProfileTarget/" + profileName + "/" + targetArn
	reverseKey := "secTargetProfile/" + targetArn + "/" + profileName
	exists, err := store.GetGenericExists(forwardKey, &map[string]interface{}{})
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrSecurityProfileAttachmentNotFound
	}
	// Attempt both deletes so a partial failure does not leave stale mappings
	// that block subsequent retries (the existence check above would reject
	// a retry after a partial delete).
	errForward := store.DeleteGeneric(forwardKey)
	errReverse := store.DeleteGeneric(reverseKey)
	if errForward != nil {
		return errForward
	}
	return errReverse
}

// listSecurityProfilesForTargetCore lists the profile names attached to a
// target.
func (s *IoTService) listSecurityProfilesForTargetCore(store iotstore.IotStoreInterface, targetArn string) ([]string, error) {
	if targetArn == "" {
		return nil, iotstore.ErrMissingParam
	}
	items, err := store.ListGeneric("secTargetProfile/" + targetArn + "/")
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(items))
	for _, rec := range items {
		profileName, _ := rec["securityProfileName"].(string)
		names = append(names, profileName)
	}
	return names, nil
}

// listTargetsForSecurityProfileCore lists the target ARNs a profile is
// attached to.
func (s *IoTService) listTargetsForSecurityProfileCore(store iotstore.IotStoreInterface, profileName string) ([]string, error) {
	if profileName == "" {
		return nil, iotstore.ErrMissingParam
	}
	items, err := store.ListGeneric("secProfileTarget/" + profileName + "/")
	if err != nil {
		return nil, err
	}
	targets := make([]string, 0, len(items))
	for _, rec := range items {
		targetArn, _ := rec["securityProfileTargetArn"].(string)
		targets = append(targets, targetArn)
	}
	return targets, nil
}

// verificationStates is the VerificationState enum member set.
var verificationStates = map[string]bool{
	"FALSE_POSITIVE": true, "BENIGN_POSITIVE": true, "TRUE_POSITIVE": true, "UNKNOWN": true,
}

// putVerificationStateOnViolationCore records the verification state on a
// violation. No Device Defender engine generates violations, so the record
// usually does not exist. AWS lists only InvalidRequestException in the
// Smithy errors trait (not ResourceNotFoundException), so an unknown
// violation id yields InvalidRequest rather than 404.
func (s *IoTService) putVerificationStateOnViolationCore(store iotstore.IotStoreInterface, violationId, verificationState, description string) error {
	if violationId == "" || verificationState == "" {
		return iotstore.ErrMissingParam
	}
	if !verificationStates[verificationState] {
		return iotstore.ErrInvalidRequest
	}
	key := "violation/" + violationId
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrInvalidRequest
	}
	rec["verificationState"] = verificationState
	if description != "" {
		rec["verificationStateDescription"] = description
	}
	return store.PutGeneric(key, rec)
}

// listActiveViolationsCore lists active violations, optionally restricted
// to one thing and one security profile.
func (s *IoTService) listActiveViolationsCore(store iotstore.IotStoreInterface, thingName, securityProfileName string) ([]*iotstore.ViolationEvent, error) {
	violations, err := store.ListActiveViolations(thingName)
	if err != nil {
		return nil, err
	}
	filtered := make([]*iotstore.ViolationEvent, 0, len(violations))
	for _, v := range violations {
		if securityProfileName != "" && v.SecurityProfileName != securityProfileName {
			continue
		}
		filtered = append(filtered, v)
	}
	return filtered, nil
}

// listViolationEventsCore lists violation events with the profile/thing
// filters.
func (s *IoTService) listViolationEventsCore(store iotstore.IotStoreInterface, opts storecommon.ListOptions, securityProfileName, thingName string) ([]*iotstore.ViolationEvent, error) {
	return store.ListViolationEvents(opts, securityProfileName, thingName)
}
