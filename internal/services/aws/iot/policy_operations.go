package iot

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

var policyNamePattern = regexp.MustCompile(`^[\w+=,.@-]+$`)

// validatePolicyDocument checks that the policy document is valid JSON with
// at least Version and Statement fields, matching AWS validation.
func validatePolicyDocument(doc string) error {
	var pd struct {
		Version   string          `json:"Version"`
		Statement json.RawMessage `json:"Statement"`
	}
	if err := json.Unmarshal([]byte(doc), &pd); err != nil {
		return err
	}
	if pd.Version == "" {
		return fmt.Errorf("policy document missing Version")
	}
	if len(pd.Statement) == 0 {
		return fmt.Errorf("policy document missing Statement")
	}
	return nil
}

func (s *IoTService) CreatePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.createPolicyCore(store, CreatePolicyInput{
		PolicyName:     request.GetParamCaseInsensitive(req.Parameters, "policyName"),
		PolicyDocument: request.GetParamCaseInsensitive(req.Parameters, "policyDocument"),
	})
	if err != nil {
		return nil, err
	}

	return policyResponse(result.Policy), nil
}

func (s *IoTService) GetPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.getPolicyCore(store, request.GetParamCaseInsensitive(req.Parameters, "policyName"))
	if err != nil {
		return nil, err
	}

	return policyResponse(result.Policy), nil
}

func (s *IoTService) DeletePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deletePolicyCore(store, request.GetParamCaseInsensitive(req.Parameters, "policyName"), reqCtx.GetAccountID(), reqCtx.GetRegion()); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) ListPolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	opts := parseListOptions(req.Parameters)
	result, err := s.listPoliciesCore(store, opts.Marker, opts.MaxItems)
	if err != nil {
		return nil, err
	}

	items := make([]map[string]interface{}, 0, len(result.Policies))
	for _, p := range result.Policies {
		resp := policyResponse(p)
		delete(resp, "policyDocument")
		items = append(items, resp)
	}

	return listResponse("policies", items, result.NextToken), nil
}

func (s *IoTService) AttachPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.attachPolicyCore(store, AttachPolicyInput{
		PolicyName: request.GetParamCaseInsensitive(req.Parameters, "policyName"),
		Target:     request.GetParamCaseInsensitive(req.Parameters, "target"),
	}); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) DetachPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.detachPolicyCore(store, AttachPolicyInput{
		PolicyName: request.GetParamCaseInsensitive(req.Parameters, "policyName"),
		Target:     request.GetParamCaseInsensitive(req.Parameters, "target"),
	}); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) AttachThingPrincipal(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.attachThingPrincipalCore(store, AttachThingPrincipalInput{
		ThingName: request.GetParamCaseInsensitive(req.Parameters, "thingName"),
		Principal: request.GetParamCaseInsensitive(req.Parameters, "principal"),
	}); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) DetachThingPrincipal(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.detachThingPrincipalCore(store, AttachThingPrincipalInput{
		ThingName: request.GetParamCaseInsensitive(req.Parameters, "thingName"),
		Principal: request.GetParamCaseInsensitive(req.Parameters, "principal"),
	}); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) ListPolicyPrincipals(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	principals, err := s.listPolicyPrincipalsCore(store, request.GetParamCaseInsensitive(req.Parameters, "policyName"))
	if err != nil {
		return nil, err
	}

	return paginatedStrings("principals", principals, req.Parameters)
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

	policyNames, err := s.listPrincipalPoliciesCore(store, principal)
	if err != nil {
		return nil, err
	}

	items := make([]map[string]interface{}, 0, len(policyNames))
	for _, name := range policyNames {
		items = append(items, map[string]interface{}{
			"policyName": name,
		})
	}

	return paginatedMaps("policies", items, req.Parameters)
}

func (s *IoTService) ListThingPrincipals(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	principals, err := s.listThingPrincipalsCore(store, request.GetParamCaseInsensitive(req.Parameters, "thingName"))
	if err != nil {
		return nil, err
	}

	return paginatedStrings("principals", principals, req.Parameters)
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

	things, err := s.listPrincipalThingsCore(store, principal)
	if err != nil {
		return nil, err
	}

	return paginatedStrings("things", things, req.Parameters)
}

func (s *IoTService) GetEffectivePolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	principal := request.GetParamCaseInsensitive(req.Parameters, "principal")
	if principal == "" {
		principal = request.GetParamCaseInsensitive(req.Parameters, "certificateArn")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	policies, err := s.getEffectivePoliciesCore(store, principal)
	if err != nil {
		return nil, err
	}

	effective := make([]map[string]interface{}, 0, len(policies))
	for _, p := range policies {
		effective = append(effective, map[string]interface{}{
			"policyName":     p.PolicyName,
			"policyArn":      p.PolicyARN,
			"policyDocument": p.PolicyDocument,
		})
	}

	return map[string]interface{}{
		"effectivePolicies": effective,
	}, nil
}

func (s *IoTService) ListPolicyVersions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.listPolicyVersionsCore(store, request.GetParamCaseInsensitive(req.Parameters, "policyName"))
	if err != nil {
		return nil, err
	}

	return paginatedMaps("policyVersions", result.PolicyVersions, req.Parameters)
}

func (s *IoTService) GetPolicyVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.getPolicyVersionCore(store,
		request.GetParamCaseInsensitive(req.Parameters, "policyName"),
		request.GetParamCaseInsensitive(req.Parameters, "policyVersionId"),
	)
}
