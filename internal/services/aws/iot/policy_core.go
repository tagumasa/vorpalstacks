package iot

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/services/aws/iot/policy"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---------------------------------------------------------------------------
// Policy attachment, listing, version, and authorisation-evaluation Core
// ---------------------------------------------------------------------------

// AttachPolicyInput carries the fields for AttachPolicy and DetachPolicy.
type AttachPolicyInput struct {
	PolicyName string
	Target     string
}

// AttachThingPrincipalInput carries the fields for AttachThingPrincipal and
// DetachThingPrincipal.
type AttachThingPrincipalInput struct {
	ThingName string
	Principal string
}

// EffectivePolicy carries one resolved policy entry for
// GetEffectivePolicies.
type EffectivePolicy struct {
	PolicyName     string
	PolicyARN      string
	PolicyDocument string
}

// CreatePolicyVersionInput carries the fields for CreatePolicyVersion.
type CreatePolicyVersionInput struct {
	PolicyName     string
	PolicyDocument string
	SetAsDefault   bool
}

// CreatePolicyVersionResult is the transport-agnostic result of
// CreatePolicyVersion (the response carries the model's four members).
type CreatePolicyVersionResult struct {
	PolicyARN       string
	PolicyDocument  string
	PolicyVersionID string
	IsDefault       bool
}

// PolicyVersionInput carries the fields for DeletePolicyVersion and
// SetDefaultPolicyVersion.
type PolicyVersionInput struct {
	PolicyName      string
	PolicyVersionID string
}

// ListPolicyVersionsResult is the transport-agnostic result of
// ListPolicyVersions.
type ListPolicyVersionsResult struct {
	PolicyVersions []map[string]interface{}
}

// TestAuthorizationAuthInfo is one parsed authInfos entry of
// TestAuthorization.
type TestAuthorizationAuthInfo struct {
	ActionType string
	Resources  []string
}

// TestAuthorizationInput carries the fields for TestAuthorization.
type TestAuthorizationInput struct {
	Principal             string
	ClientID              string
	CognitoIdentityPoolID string
	PolicyNamesToAdd      []string
	PolicyNamesToSkip     []string
	AuthInfos             []TestAuthorizationAuthInfo
}

// attachPolicyCore attaches a policy to a principal (the AttachPolicy
// operation).
func (s *IoTService) attachPolicyCore(store iotstore.IotStoreInterface, in AttachPolicyInput) error {
	if in.PolicyName == "" || in.Target == "" {
		return iotstore.ErrMissingParam
	}
	return store.AttachPolicyToPrincipal(in.PolicyName, in.Target)
}

// detachPolicyCore detaches a policy from a principal (the DetachPolicy
// operation).
func (s *IoTService) detachPolicyCore(store iotstore.IotStoreInterface, in AttachPolicyInput) error {
	if in.PolicyName == "" || in.Target == "" {
		return iotstore.ErrMissingParam
	}
	return store.DetachPolicyFromPrincipal(in.PolicyName, in.Target)
}

// attachThingPrincipalCore binds a principal (certificate ARN) to a thing.
func (s *IoTService) attachThingPrincipalCore(store iotstore.IotStoreInterface, in AttachThingPrincipalInput) error {
	if in.ThingName == "" || in.Principal == "" {
		return iotstore.ErrMissingParam
	}
	return store.AttachThingPrincipal(in.ThingName, in.Principal)
}

// detachThingPrincipalCore unbinds a principal from a thing.
func (s *IoTService) detachThingPrincipalCore(store iotstore.IotStoreInterface, in AttachThingPrincipalInput) error {
	if in.ThingName == "" || in.Principal == "" {
		return iotstore.ErrMissingParam
	}
	return store.DetachThingPrincipal(in.ThingName, in.Principal)
}

// listPolicyPrincipalsCore lists the principals a policy is attached to (the
// ListPolicyPrincipals and ListTargetsForPolicy operations).
func (s *IoTService) listPolicyPrincipalsCore(store iotstore.IotStoreInterface, policyName string) ([]string, error) {
	if policyName == "" {
		return nil, iotstore.ErrMissingParam
	}
	return store.ListPrincipalsForPolicy(policyName)
}

// listPrincipalPoliciesCore lists the policies attached to a principal.
func (s *IoTService) listPrincipalPoliciesCore(store iotstore.IotStoreInterface, principal string) ([]string, error) {
	return store.ListPoliciesForPrincipal(principal)
}

// listThingPrincipalsCore lists the principals attached to a thing.
func (s *IoTService) listThingPrincipalsCore(store iotstore.IotStoreInterface, thingName string) ([]string, error) {
	if thingName == "" {
		return nil, iotstore.ErrMissingParam
	}
	return store.ListPrincipalsForThing(thingName)
}

// listPrincipalThingsCore lists the things bound to a principal. An empty
// principal yields an empty list without touching the store.
func (s *IoTService) listPrincipalThingsCore(store iotstore.IotStoreInterface, principal string) ([]string, error) {
	if principal == "" {
		return []string{}, nil
	}
	return store.ListThingsForPrincipal(principal)
}

// getEffectivePoliciesCore resolves every policy that applies to a principal:
// the policies attached directly, the policies attached to the principal's
// things, and the policies attached to those things' groups.
func (s *IoTService) getEffectivePoliciesCore(store iotstore.IotStoreInterface, principal string) ([]*EffectivePolicy, error) {
	if principal == "" {
		return []*EffectivePolicy{}, nil
	}
	accountID := store.GetAccountID()
	region := store.GetRegion()

	seen := make(map[string]bool)
	var policies []*EffectivePolicy

	addPolicies := func(policyNames []string) {
		for _, name := range policyNames {
			if seen[name] {
				continue
			}
			seen[name] = true
			pol, err := store.GetPolicy(name)
			if err != nil {
				slog.Warn("GetEffectivePolicies: failed to load policy", "policy", name, "error", err)
				continue
			}
			policies = append(policies, &EffectivePolicy{
				PolicyName:     pol.PolicyName,
				PolicyARN:      pol.PolicyARN,
				PolicyDocument: pol.PolicyDocument,
			})
		}
	}

	principalPolicies, err := store.ListPoliciesForPrincipal(principal)
	if err != nil {
		return nil, err
	}
	addPolicies(principalPolicies)

	things, err := store.ListThingsForPrincipal(principal)
	if err != nil {
		slog.Warn("GetEffectivePolicies: failed to list things for principal", "principal", principal, "error", err)
		return nil, fmt.Errorf("failed to list things for principal: %w", err)
	}
	for _, thingName := range things {
		thingARN := iotstore.BuildThingARN(accountID, region, thingName)
		thingPolicies, err := store.ListPoliciesForPrincipal(thingARN)
		if err != nil {
			slog.Warn("GetEffectivePolicies: failed to list policies for thing", "thing", thingARN, "error", err)
			continue
		}
		addPolicies(thingPolicies)
		groups, err := store.ListGroupsForThing(thingName)
		if err != nil {
			slog.Warn("GetEffectivePolicies: failed to list groups for thing", "thing", thingName, "error", err)
			continue
		}
		for _, groupName := range groups {
			groupARN := iotstore.BuildThingGroupARN(accountID, region, groupName)
			groupPolicies, err := store.ListPoliciesForPrincipal(groupARN)
			if err != nil {
				slog.Warn("GetEffectivePolicies: failed to list policies for group", "group", groupARN, "error", err)
				continue
			}
			addPolicies(groupPolicies)
		}
	}
	return policies, nil
}

// createPolicyVersionCore appends a new version to an existing policy. AWS
// IoT limits policies to 5 versions and numbers versions from the policy's
// own default version onward with purely numeric ids; a setAsDefault
// version also syncs the policy record's document so the MQTT broker uses
// the new default.
func (s *IoTService) createPolicyVersionCore(store iotstore.IotStoreInterface, in CreatePolicyVersionInput) (*CreatePolicyVersionResult, error) {
	if in.PolicyName == "" || in.PolicyDocument == "" {
		return nil, iotstore.ErrMissingParam
	}
	p, err := store.GetPolicy(in.PolicyName)
	if err != nil {
		return nil, err
	}
	if err := validatePolicyDocument(in.PolicyDocument); err != nil {
		return nil, iotstore.ErrMalformedPolicy
	}
	versions := loadPolicyVersions(store, in.PolicyName)
	if len(versions) >= 5 {
		return nil, iotstore.ErrPolicyVersionLimitExceeded
	}
	maxVer := int(p.Version)
	for _, v := range versions {
		if n := policyVersionNumber(v["versionId"]); n > maxVer {
			maxVer = n
		}
	}
	newVersion := map[string]interface{}{
		"versionId":        strconv.Itoa(maxVer + 1),
		"policyDocument":   in.PolicyDocument,
		"isDefaultVersion": in.SetAsDefault,
		"createDate":       time.Now().UTC().Unix(),
	}
	if in.SetAsDefault {
		for i := range versions {
			versions[i]["isDefaultVersion"] = false
		}
	}
	versions = append(versions, newVersion)
	if err := store.PutGeneric("policyVersions/"+in.PolicyName, map[string]interface{}{"versions": versions}); err != nil {
		return nil, err
	}
	if in.SetAsDefault {
		// Sync the Policy record's document and default-version pointer so
		// the MQTT broker uses the new default and GetPolicy reports it.
		p.PolicyDocument = in.PolicyDocument
		p.Version = int64(policyVersionNumber(newVersion["versionId"]))
		p.LastModifiedDate = time.Now().UTC()
		if err := store.UpdatePolicy(p); err != nil {
			return nil, err
		}
	}
	return &CreatePolicyVersionResult{
		PolicyARN:       iotstore.BuildPolicyARN(store.GetAccountID(), store.GetRegion(), in.PolicyName),
		PolicyDocument:  in.PolicyDocument,
		PolicyVersionID: newVersion["versionId"].(string),
		IsDefault:       in.SetAsDefault,
	}, nil
}

// policyVersionNumber extracts the numeric part of a stored version id.
// Records written by the earlier scheme carry a "v" prefix; both shapes are
// tolerated, the same way listPolicyVersionsCore tolerates older createDate
// shapes.
func policyVersionNumber(id interface{}) int {
	s, _ := id.(string)
	s = strings.TrimPrefix(s, "v")
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}

// deletePolicyVersionCore removes a non-default policy version. AWS IoT
// rejects deletion of the default version.
func (s *IoTService) deletePolicyVersionCore(store iotstore.IotStoreInterface, in PolicyVersionInput) error {
	if in.PolicyName == "" || in.PolicyVersionID == "" {
		return iotstore.ErrMissingParam
	}
	versions := loadPolicyVersions(store, in.PolicyName)
	found := false
	for i, v := range versions {
		if id, _ := v["versionId"].(string); id == in.PolicyVersionID {
			if isDefault, _ := v["isDefaultVersion"].(bool); isDefault {
				return iotstore.ErrDeleteConflict
			}
			versions = append(versions[:i], versions[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return iotstore.ErrPolicyVersionNotFound
	}
	return store.PutGeneric("policyVersions/"+in.PolicyName, map[string]interface{}{"versions": versions})
}

// setDefaultPolicyVersionCore marks one stored version as the default and
// syncs the Policy record so the MQTT broker uses the new default document.
func (s *IoTService) setDefaultPolicyVersionCore(store iotstore.IotStoreInterface, in PolicyVersionInput) error {
	if in.PolicyName == "" || in.PolicyVersionID == "" {
		return iotstore.ErrMissingParam
	}
	versions := loadPolicyVersions(store, in.PolicyName)
	found := false
	var defaultDoc string
	for i := range versions {
		id, _ := versions[i]["versionId"].(string)
		isDefault := id == in.PolicyVersionID
		versions[i]["isDefaultVersion"] = isDefault
		if isDefault {
			found = true
			defaultDoc, _ = versions[i]["policyDocument"].(string)
		}
	}
	if !found {
		return iotstore.ErrPolicyVersionNotFound
	}
	if err := store.PutGeneric("policyVersions/"+in.PolicyName, map[string]interface{}{"versions": versions}); err != nil {
		return err
	}
	if defaultDoc != "" {
		p, err := store.GetPolicy(in.PolicyName)
		if err != nil {
			return err
		}
		p.PolicyDocument = defaultDoc
		p.Version = int64(policyVersionNumber(in.PolicyVersionID))
		p.LastModifiedDate = time.Now().UTC()
		if err := store.UpdatePolicy(p); err != nil {
			return err
		}
	}
	return nil
}

// listPolicyVersionsCore lists every version of a policy from the persisted
// version set — the single source of version identity, so exactly one entry
// can carry the default flag at a time.
func (s *IoTService) listPolicyVersionsCore(store iotstore.IotStoreInterface, policyName string) (*ListPolicyVersionsResult, error) {
	if policyName == "" {
		return nil, iotstore.ErrMissingParam
	}
	if _, err := store.GetPolicy(policyName); err != nil {
		return nil, err
	}
	versions := make([]map[string]interface{}, 0)
	for _, v := range loadPolicyVersions(store, policyName) {
		isDefault, _ := v["isDefaultVersion"].(bool)
		entry := map[string]interface{}{
			"versionId":        v["versionId"],
			"isDefaultVersion": isDefault,
		}
		// createDate must be a JSON number (unix seconds) for the SDK; older
		// records may have stored an RFC3339 string, so coerce both shapes.
		switch cd := v["createDate"].(type) {
		case float64:
			entry["createDate"] = int64(cd)
		case int64:
			entry["createDate"] = cd
		case int:
			entry["createDate"] = int64(cd)
		case string:
			if t, err := time.Parse(time.RFC3339Nano, cd); err == nil {
				entry["createDate"] = t.Unix()
			}
		}
		versions = append(versions, entry)
	}
	return &ListPolicyVersionsResult{PolicyVersions: versions}, nil
}

// getPolicyVersionCore resolves any version of a policy from the persisted
// version set, which holds every version including the initial default.
func (s *IoTService) getPolicyVersionCore(store iotstore.IotStoreInterface, policyName, versionID string) (map[string]interface{}, error) {
	if policyName == "" || versionID == "" {
		return nil, iotstore.ErrMissingParam
	}
	p, err := store.GetPolicy(policyName)
	if err != nil {
		return nil, err
	}
	for _, v := range loadPolicyVersions(store, policyName) {
		if id, _ := v["versionId"].(string); id == versionID {
			return map[string]interface{}{
				"policyName":       p.PolicyName,
				"policyArn":        p.PolicyARN,
				"policyDocument":   v["policyDocument"],
				"policyVersionId":  v["versionId"],
				"isDefaultVersion": v["isDefaultVersion"],
				"createDate":       v["createDate"],
			}, nil
		}
	}
	return nil, iotstore.ErrPolicyVersionNotFound
}

// loadPolicyVersions reads the persisted version slice for a policy. Returns
// an empty slice when the policy has no recorded versions yet.
func loadPolicyVersions(store iotstore.IotStoreInterface, policyName string) []map[string]interface{} {
	wrap := map[string]interface{}{}
	if err := store.GetGeneric("policyVersions/"+policyName, &wrap); err != nil {
		return []map[string]interface{}{}
	}
	raw, ok := wrap["versions"].([]interface{})
	if !ok {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, 0, len(raw))
	for _, item := range raw {
		if m, ok := item.(map[string]interface{}); ok {
			out = append(out, m)
		}
	}
	return out
}

// testAuthorizationCore resolves the policies attached to a principal and
// evaluates each requested action/resource combination against them,
// mirroring the MQTT broker auth hook (auth_hooks.go).
func (s *IoTService) testAuthorizationCore(store iotstore.IotStoreInterface, in TestAuthorizationInput) ([]map[string]interface{}, error) {
	if in.Principal == "" {
		return []map[string]interface{}{}, nil
	}
	policyNames, err := store.ListPoliciesForPrincipal(in.Principal)
	if err != nil {
		return nil, err
	}
	// Merge additional policies supplied by the caller (policyNamesToAdd).
	skipSet := make(map[string]bool, len(in.PolicyNamesToSkip))
	for _, n := range in.PolicyNamesToSkip {
		skipSet[n] = true
	}
	for _, n := range in.PolicyNamesToAdd {
		if !skipSet[n] {
			policyNames = append(policyNames, n)
		}
	}
	// When cognitoIdentityPoolId is supplied, also include policies attached
	// to the cognito identity pool principal.
	if in.CognitoIdentityPoolID != "" {
		cognitoPolicies, err := store.ListPoliciesForPrincipal(in.CognitoIdentityPoolID)
		if err == nil {
			for _, n := range cognitoPolicies {
				if !skipSet[n] {
					policyNames = append(policyNames, n)
				}
			}
		}
	}
	var versions []*policy.PolicyVersion
	matchedNames := make([]string, 0, len(policyNames))
	for _, name := range policyNames {
		if skipSet[name] {
			continue
		}
		p, gErr := store.GetPolicy(name)
		if gErr != nil {
			continue
		}
		pv, pErr := policy.ParsePolicyVersion([]byte(p.PolicyDocument))
		if pErr != nil {
			continue
		}
		versions = append(versions, pv)
		matchedNames = append(matchedNames, name)
	}
	policyEntries := make([]map[string]interface{}, 0, len(matchedNames))
	for _, n := range matchedNames {
		policyEntries = append(policyEntries, map[string]interface{}{
			"policyName": n,
			"policyArn":  iotstore.BuildPolicyARN(store.GetAccountID(), store.GetRegion(), n),
		})
	}
	results := make([]map[string]interface{}, 0, len(in.AuthInfos))
	for _, info := range in.AuthInfos {
		action := info.ActionType
		resources := info.Resources
		evalResource := "*"
		if len(resources) > 0 {
			evalResource = resources[0]
		}
		// Normalise the action to "iot:TitleCase" for policy evaluation.
		// The CLI sends lowercase (e.g. "connect") but policies use "iot:Connect".
		iotAction := action
		if action != "" && !strings.HasPrefix(action, "iot:") {
			iotAction = "iot:" + strings.ToUpper(action[:1]) + strings.ToLower(action[1:])
		}
		allowed, _ := policy.Evaluate(&policy.EvaluateRequest{
			Policies: versions,
			Action:   iotAction,
			Resource: evalResource,
			ClientID: in.ClientID,
		})
		entry := map[string]interface{}{
			"authInfo": map[string]interface{}{
				"actionType": action,
				"resources":  resources,
			},
			"matchedPolicies": matchedNames,
		}
		if allowed {
			entry["allowed"] = map[string]interface{}{"policies": policyEntries}
			entry["decision"] = "ALLOWED"
		} else {
			entry["denied"] = map[string]interface{}{
				"implicitDeny": map[string]interface{}{"policies": policyEntries},
			}
			entry["decision"] = "IMPLICIT_DENY"
		}
		results = append(results, entry)
	}
	return results, nil
}
