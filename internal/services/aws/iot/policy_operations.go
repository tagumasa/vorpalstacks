package iot

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

func (s *IoTService) CreatePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyName := request.GetParamCaseInsensitive(req.Parameters, "policyName")
	policyDoc := request.GetParamCaseInsensitive(req.Parameters, "policyDocument")
	if policyName == "" || policyDoc == "" {
		return nil, iotstore.ErrMissingParam
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
		return nil, iotstore.ErrPolicyNotFound
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

	return map[string]interface{}{
		"principals": principals,
	}, nil
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

	return map[string]interface{}{
		"policies": items,
	}, nil
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

	return map[string]interface{}{
		"principals": principals,
	}, nil
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

	return map[string]interface{}{
		"things": things,
	}, nil
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
		return nil, iotstore.ErrPolicyNotFound
	}

	return map[string]interface{}{
		"policyVersions": []map[string]interface{}{
			{
				"versionId":        fmt.Sprintf("%d", p.Version),
				"isDefaultVersion": true,
				"createDate":       p.CreationDate.Unix(),
			},
		},
	}, nil
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
		return nil, iotstore.ErrPolicyNotFound
	}

	currentVersion := fmt.Sprintf("%d", p.Version)
	if versionId != currentVersion {
		return nil, iotstore.ErrPolicyVersionNotFound
	}

	return map[string]interface{}{
		"policyName":       p.PolicyName,
		"policyArn":        p.PolicyARN,
		"policyDocument":   p.PolicyDocument,
		"versionId":        fmt.Sprintf("%d", p.Version),
		"isDefaultVersion": true,
		"createDate":       p.CreationDate.Unix(),
	}, nil
}
