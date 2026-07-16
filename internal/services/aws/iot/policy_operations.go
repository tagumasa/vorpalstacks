package iot

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"time"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

var policyNamePattern = regexp.MustCompile(`^[\w+=,.@-]+$`)

func (s *IoTService) CreatePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyName := request.GetParamCaseInsensitive(req.Parameters, "policyName")
	policyDoc := request.GetParamCaseInsensitive(req.Parameters, "policyDocument")
	if policyName == "" || policyDoc == "" {
		return nil, iotstore.ErrMissingParam
	}
	if !policyNamePattern.MatchString(policyName) {
		return nil, iotstore.ErrInvalidRequest
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	p := &iotstore.Policy{
		PolicyName:       policyName,
		PolicyDocument:   policyDoc,
		CreationDate:     time.Now().UTC(),
		LastModifiedDate: time.Now().UTC(),
		Version:          1,
	}

	created, err := store.CreatePolicy(p)
	if err != nil {
		return nil, err
	}

	return policyResponse(created), nil
}

func (s *IoTService) GetPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyName := request.GetParamCaseInsensitive(req.Parameters, "policyName")
	if policyName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	p, err := store.GetPolicy(policyName)
	if err != nil {
		return nil, err
	}

	return policyResponse(p), nil
}

func (s *IoTService) DeletePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyName := request.GetParamCaseInsensitive(req.Parameters, "policyName")
	if policyName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	arn := iotstore.BuildPolicyARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), policyName)
	_ = store.DeleteAllTags(arn)

	if err := store.DeletePolicy(policyName); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) ListPolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	policies, err := store.ListPolicies(parseListOptions(req.Parameters))
	if err != nil {
		return nil, err
	}

	items := make([]map[string]interface{}, 0, len(policies.Items))
	for _, p := range policies.Items {
		resp := policyResponse(p)
		delete(resp, "policyDocument")
		items = append(items, resp)
	}

	return listResponse("policies", items, policies.NextMarker), nil
}

func (s *IoTService) AttachPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyName := request.GetParamCaseInsensitive(req.Parameters, "policyName")
	target := request.GetParamCaseInsensitive(req.Parameters, "target")
	if policyName == "" || target == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.AttachPolicyToPrincipal(policyName, target); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) DetachPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyName := request.GetParamCaseInsensitive(req.Parameters, "policyName")
	target := request.GetParamCaseInsensitive(req.Parameters, "target")
	if policyName == "" || target == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DetachPolicyFromPrincipal(policyName, target); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) AttachThingPrincipal(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	principal := request.GetParamCaseInsensitive(req.Parameters, "principal")
	if thingName == "" || principal == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.AttachThingPrincipal(thingName, principal); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) DetachThingPrincipal(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	principal := request.GetParamCaseInsensitive(req.Parameters, "principal")
	if thingName == "" || principal == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DetachThingPrincipal(thingName, principal); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) ListPolicyPrincipals(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyName := request.GetParamCaseInsensitive(req.Parameters, "policyName")
	if policyName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	principals, err := store.ListPrincipalsForPolicy(policyName)
	if err != nil {
		return nil, err
	}

	return paginatedStrings("principals", principals, req.Parameters), nil
}

func (s *IoTService) ListPrincipalPolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	principal, ok := principalFromParams(req.Parameters)
	if !ok {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	policyNames, err := store.ListPoliciesForPrincipal(principal)
	if err != nil {
		return nil, err
	}

	items := make([]map[string]interface{}, 0, len(policyNames))
	for _, name := range policyNames {
		items = append(items, map[string]interface{}{
			"policyName": name,
		})
	}

	return paginatedMaps("policies", items, req.Parameters), nil
}

func (s *IoTService) ListThingPrincipals(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	if thingName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	principals, err := store.ListPrincipalsForThing(thingName)
	if err != nil {
		return nil, err
	}

	return paginatedStrings("principals", principals, req.Parameters), nil
}

// ListPrincipalThings returns things attached to the specified principal (certificate).
func (s *IoTService) ListPrincipalThings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	principal := request.GetParamCaseInsensitive(req.Parameters, "principal")
	if principal == "" {
		principal = request.GetParamCaseInsensitive(req.Parameters, "certificateArn")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if principal == "" {
		return map[string]interface{}{
			"things": []string{},
		}, nil
	}

	things, err := store.ListThingsForPrincipal(principal)
	if err != nil {
		return nil, err
	}

	return paginatedStrings("things", things, req.Parameters), nil
}

func (s *IoTService) GetEffectivePolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	principal := request.GetParamCaseInsensitive(req.Parameters, "principal")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if principal == "" {
		principal = request.GetParamCaseInsensitive(req.Parameters, "certificateArn")
	}

	if principal == "" {
		return map[string]interface{}{
			"effectivePolicies": []map[string]interface{}{},
		}, nil
	}

	seen := make(map[string]bool)
	var policies []map[string]interface{}

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
			policies = append(policies, map[string]interface{}{
				"policyName":     pol.PolicyName,
				"policyArn":      pol.PolicyARN,
				"policyDocument": pol.PolicyDocument,
			})
		}
	}

	principalPolicies, err := store.ListPoliciesForPrincipal(principal)
	if err != nil {
		return nil, err
	}
	addPolicies(principalPolicies)

	region := reqCtx.GetRegion()
	things, err := store.ListThingsForPrincipal(principal)
	if err != nil {
		slog.Warn("GetEffectivePolicies: failed to list things for principal", "principal", principal, "error", err)
		return nil, fmt.Errorf("failed to list things for principal: %w", err)
	}
	for _, thingName := range things {
		thingARN := iotstore.BuildThingARN(s.accountID, region, thingName)
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
			groupARN := iotstore.BuildThingGroupARN(s.accountID, region, groupName)
			groupPolicies, err := store.ListPoliciesForPrincipal(groupARN)
			if err != nil {
				slog.Warn("GetEffectivePolicies: failed to list policies for group", "group", groupARN, "error", err)
				continue
			}
			addPolicies(groupPolicies)
		}
	}

	if len(policies) == 0 {
		return map[string]interface{}{
			"effectivePolicies": []map[string]interface{}{},
		}, nil
	}

	return map[string]interface{}{
		"effectivePolicies": policies,
	}, nil
}

func (s *IoTService) ListPolicyVersions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyName := request.GetParamCaseInsensitive(req.Parameters, "policyName")
	if policyName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	p, err := store.GetPolicy(policyName)
	if err != nil {
		return nil, err
	}

	// The default version lives on the policy record; additional non-default
	// versions are persisted in the generic-KV store by CreatePolicyVersion.
	// Merge both so List reflects every version actually created.
	versions := []map[string]interface{}{
		{
			"versionId":        fmt.Sprintf("%d", p.Version),
			"isDefaultVersion": true,
			"createDate":       p.CreationDate.Unix(),
		},
	}
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

	return paginatedMaps("policyVersions", versions, req.Parameters), nil
}

func (s *IoTService) GetPolicyVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyName := request.GetParamCaseInsensitive(req.Parameters, "policyName")
	versionId := request.GetParamCaseInsensitive(req.Parameters, "policyVersionId")
	if policyName == "" || versionId == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	p, err := store.GetPolicy(policyName)
	if err != nil {
		return nil, err
	}

	// The default version lives on the policy record (numeric id); additional
	// non-default versions ("vN") are persisted in the generic-KV store by
	// CreatePolicyVersion. Check both so GetPolicyVersion resolves any version.
	defaultID := fmt.Sprintf("%d", p.Version)
	if versionId == defaultID {
		return map[string]interface{}{
			"policyName":       p.PolicyName,
			"policyArn":        p.PolicyARN,
			"policyDocument":   p.PolicyDocument,
			"policyVersionId":  defaultID,
			"isDefaultVersion": true,
			"createDate":       p.CreationDate.Unix(),
		}, nil
	}

	for _, v := range loadPolicyVersions(store, policyName) {
		if id, _ := v["versionId"].(string); id == versionId {
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
