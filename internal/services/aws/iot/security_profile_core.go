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

// ListSecurityProfilesInput carries the parsed ListSecurityProfiles request.
// The dimensionName and metricName filters are mutually exclusive per the
// model documentation.
type ListSecurityProfilesInput struct {
	DimensionName string
	MetricName    string
}

// listSecurityProfilesCore lists security profiles with their identifiers,
// optionally restricted to the profiles that reference the given dimension
// or custom metric.
func (s *IoTService) listSecurityProfilesCore(store iotstore.IotStoreInterface, opts storecommon.ListOptions, in ListSecurityProfilesInput) ([]SecurityProfileListItem, string, error) {
	if in.DimensionName != "" && in.MetricName != "" {
		return nil, "", iotstore.ErrInvalidRequest
	}
	profiles, err := store.ListSecurityProfiles(opts)
	if err != nil {
		return nil, "", err
	}
	items := make([]SecurityProfileListItem, 0, len(profiles.Items))
	for _, sp := range profiles.Items {
		if !securityProfileMatchesFilter(sp, in) {
			continue
		}
		items = append(items, SecurityProfileListItem{
			Name: sp.SecurityProfileName,
			Arn:  sp.SecurityProfileARN,
		})
	}
	return items, profiles.NextMarker, nil
}

// securityProfileMatchesFilter reports whether a profile references the
// filtered dimension (through a behaviour metric dimension or a retained
// metric's dimension) or the filtered metric (through a behaviour metric or
// a retained metric).
func securityProfileMatchesFilter(sp *iotstore.SecurityProfile, in ListSecurityProfilesInput) bool {
	if in.DimensionName == "" && in.MetricName == "" {
		return true
	}
	for _, b := range sp.Behaviors {
		if b == nil {
			continue
		}
		if in.DimensionName != "" && b.MetricDimension == in.DimensionName {
			return true
		}
		if in.MetricName != "" && b.Metric == in.MetricName {
			return true
		}
	}
	for _, m := range sp.AdditionalMetricsToRetainV2 {
		if m == nil {
			continue
		}
		if in.DimensionName != "" && m.MetricDimension == in.DimensionName {
			return true
		}
		if in.MetricName != "" && m.Metric == in.MetricName {
			return true
		}
	}
	return false
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

// ListActiveViolationsInput carries the parsed ListActiveViolations request
// filters.
type ListActiveViolationsInput struct {
	ThingName            string
	SecurityProfileName  string
	BehaviorCriteriaType string
	VerificationState    string
	ListSuppressedAlerts *bool
}

// behaviorCriteriaTypes is the BehaviorCriteriaType enum member set.
var behaviorCriteriaTypes = map[string]bool{
	"STATIC": true, "STATISTICAL": true, "MACHINE_LEARNING": true,
}

// listActiveViolationsCore lists active violations filtered by thing,
// profile, criteria type, verification state and the suppressed-alert
// selector. Suppressed alerts (behaviours flagged suppressAlerts) are only
// returned when listSuppressedAlerts is explicitly set.
func (s *IoTService) listActiveViolationsCore(store iotstore.IotStoreInterface, in ListActiveViolationsInput) ([]*iotstore.ViolationEvent, error) {
	if in.BehaviorCriteriaType != "" && !behaviorCriteriaTypes[in.BehaviorCriteriaType] {
		return nil, iotstore.ErrInvalidRequest
	}
	violations, err := store.ListActiveViolations(in.ThingName)
	if err != nil {
		return nil, err
	}
	filtered := make([]*iotstore.ViolationEvent, 0, len(violations))
	for _, v := range violations {
		if violationMatchesFilters(v, in.SecurityProfileName, in.BehaviorCriteriaType, in.VerificationState, in.ListSuppressedAlerts) {
			filtered = append(filtered, v)
		}
	}
	return filtered, nil
}

// ListViolationEventsInput carries the parsed ListViolationEvents request.
// The model marks both time-range members required.
type ListViolationEventsInput struct {
	StartTime            int64
	EndTime              int64
	StartTimeProvided    bool
	EndTimeProvided      bool
	SecurityProfileName  string
	ThingName            string
	BehaviorCriteriaType string
	VerificationState    string
	ListSuppressedAlerts *bool
}

// listViolationEventsCore lists violation events inside the required
// [startTime, endTime] window with the same filter set as the active
// violations list.
func (s *IoTService) listViolationEventsCore(store iotstore.IotStoreInterface, opts storecommon.ListOptions, in ListViolationEventsInput) ([]*iotstore.ViolationEvent, error) {
	if !in.StartTimeProvided || !in.EndTimeProvided {
		return nil, iotstore.ErrInvalidRequest
	}
	if in.BehaviorCriteriaType != "" && !behaviorCriteriaTypes[in.BehaviorCriteriaType] {
		return nil, iotstore.ErrInvalidRequest
	}
	events, err := store.ListViolationEvents(opts, in.SecurityProfileName, in.ThingName)
	if err != nil {
		return nil, err
	}
	filtered := make([]*iotstore.ViolationEvent, 0, len(events))
	for _, e := range events {
		eventTime := e.ViolationEventTime.Unix()
		if eventTime < in.StartTime || eventTime > in.EndTime {
			continue
		}
		if violationMatchesFilters(e, "", in.BehaviorCriteriaType, in.VerificationState, in.ListSuppressedAlerts) {
			filtered = append(filtered, e)
		}
	}
	return filtered, nil
}

// violationMatchesFilters applies the shared security-profile, criteria-type,
// verification-state and suppressed-alert filters to one violation.
func violationMatchesFilters(v *iotstore.ViolationEvent, securityProfileName, behaviorCriteriaType, verificationState string, listSuppressed *bool) bool {
	if securityProfileName != "" && v.SecurityProfileName != securityProfileName {
		return false
	}
	if behaviorCriteriaType != "" && behaviorCriteriaTypeOf(v.Behavior) != behaviorCriteriaType {
		return false
	}
	if verificationState != "" && v.VerificationState != verificationState {
		return false
	}
	if listSuppressed == nil || !*listSuppressed {
		if v.Behavior != nil && v.Behavior.SuppressAlerts {
			return false
		}
	}
	return true
}

// behaviorCriteriaTypeOf classifies a behaviour's criteria onto the
// BehaviorCriteriaType enum: an ML detection config is MACHINE_LEARNING, a
// statistical threshold is STATISTICAL, any other criteria is STATIC, and a
// behaviour without criteria matches none of the typed filters.
func behaviorCriteriaTypeOf(b *iotstore.Behavior) string {
	if b == nil || b.Criteria == nil {
		return ""
	}
	if b.Criteria.MLDetectionConfig != nil {
		return "MACHINE_LEARNING"
	}
	if b.Criteria.StatisticalThreshold != nil {
		return "STATISTICAL"
	}
	return "STATIC"
}
