package iot

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

func (s *IoTService) CreatePolicyVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyName := request.GetParamCaseInsensitive(req.Parameters, "policyName")
	if policyName == "" {
		return nil, iotstore.ErrMissingParam
	}
	doc := request.GetParamCaseInsensitive(req.Parameters, "policyDocument")
	setAsDefault := request.GetBoolParam(req.Parameters, "setAsDefault")
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	versions := loadPolicyVersions(store, policyName)
	// AWS IoT limits policies to 5 versions.
	if len(versions) >= 5 {
		return nil, iotstore.ErrPolicyVersionLimitExceeded
	}
	maxVer := 0
	for _, v := range versions {
		id, _ := v["versionId"].(string)
		if len(id) > 1 && id[0] == 'v' {
			if n, err := strconv.Atoi(id[1:]); err == nil && n > maxVer {
				maxVer = n
			}
		}
	}
	newVersion := map[string]interface{}{
		"versionId":        "v" + strconv.Itoa(maxVer+1),
		"policyDocument":   doc,
		"isDefaultVersion": setAsDefault,
		"createDate":       time.Now().UTC().Unix(),
	}
	if setAsDefault {
		for i := range versions {
			versions[i]["isDefaultVersion"] = false
		}
	}
	versions = append(versions, newVersion)
	if err := store.PutGeneric("policyVersions/"+policyName, map[string]interface{}{"versions": versions}); err != nil {
		return nil, err
	}
	if setAsDefault {
		// Sync the Policy record's document so the MQTT broker uses the new default.
		p, err := store.GetPolicy(policyName)
		if err != nil {
			return nil, err
		}
		p.PolicyDocument = doc
		p.LastModifiedDate = time.Now().UTC()
		if err := store.UpdatePolicy(p); err != nil {
			return nil, err
		}
	}
	return map[string]interface{}{
		"policyArn":        iotstore.BuildPolicyARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), policyName),
		"policyDocument":   doc,
		"policyName":       policyName,
		"policyVersionId":  newVersion["versionId"],
		"isDefaultVersion": newVersion["isDefaultVersion"],
	}, nil
}

func (s *IoTService) DeletePolicyVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyName := request.GetParamCaseInsensitive(req.Parameters, "policyName")
	versionID := request.GetParamCaseInsensitive(req.Parameters, "policyVersionId")
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	versions := loadPolicyVersions(store, policyName)
	found := false
	for i, v := range versions {
		if id, _ := v["versionId"].(string); id == versionID {
			// AWS IoT rejects deletion of the default version.
			if isDefault, _ := v["isDefaultVersion"].(bool); isDefault {
				return nil, iotstore.ErrDeleteConflict
			}
			versions = append(versions[:i], versions[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return nil, iotstore.ErrPolicyVersionNotFound
	}
	if err := store.PutGeneric("policyVersions/"+policyName, map[string]interface{}{"versions": versions}); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) SetDefaultPolicyVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyName := request.GetParamCaseInsensitive(req.Parameters, "policyName")
	versionID := request.GetParamCaseInsensitive(req.Parameters, "policyVersionId")
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	versions := loadPolicyVersions(store, policyName)
	found := false
	var defaultDoc string
	for i := range versions {
		id, _ := versions[i]["versionId"].(string)
		isDefault := id == versionID
		versions[i]["isDefaultVersion"] = isDefault
		if isDefault {
			found = true
			defaultDoc, _ = versions[i]["policyDocument"].(string)
		}
	}
	if !found {
		return nil, iotstore.ErrPolicyVersionNotFound
	}
	if err := store.PutGeneric("policyVersions/"+policyName, map[string]interface{}{"versions": versions}); err != nil {
		return nil, err
	}
	// Sync the Policy record so the MQTT broker uses the new default document.
	if defaultDoc != "" {
		p, err := store.GetPolicy(policyName)
		if err != nil {
			return nil, err
		}
		p.PolicyDocument = defaultDoc
		p.LastModifiedDate = time.Now().UTC()
		if err := store.UpdatePolicy(p); err != nil {
			return nil, err
		}
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) ListTargetsForPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
	targets := make([]map[string]interface{}, 0, len(principals))
	for _, p := range principals {
		targets = append(targets, map[string]interface{}{
			"targetArn": p,
		})
	}
	return paginatedMaps("targets", targets, req.Parameters), nil
}

// AttachPrincipalPolicy is the legacy alias of AttachPolicy. AWS accepts a
// "principal" parameter instead of "target" (Smithy AttachPrincipalPolicyRequest
// vs AttachPolicyRequest). Map the parameter so the alias dispatches through
// AttachPolicy without losing the principal identifier.
func (s *IoTService) AttachPrincipalPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if _, ok := req.Parameters["target"]; !ok {
		if principal, ok := req.Parameters["principal"]; ok {
			req.Parameters["target"] = principal
		}
	}
	return s.AttachPolicy(ctx, reqCtx, req)
}

// DetachPrincipalPolicy is the legacy alias of DetachPolicy. Same parameter
// remap as AttachPrincipalPolicy.
func (s *IoTService) DetachPrincipalPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if _, ok := req.Parameters["target"]; !ok {
		if principal, ok := req.Parameters["principal"]; ok {
			req.Parameters["target"] = principal
		}
	}
	return s.DetachPolicy(ctx, reqCtx, req)
}

func (s *IoTService) SetDefaultAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "authorizerName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.PutGeneric("config/defaultAuthorizer", name); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"authorizerName": name,
	}, nil
}

func (s *IoTService) ClearDefaultAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.DeleteGeneric("config/defaultAuthorizer"); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) DescribeDefaultAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	name := ""
	exists, err := store.GetGenericExists("config/defaultAuthorizer", &name)
	if err != nil {
		return nil, err
	}
	if !exists || name == "" {
		return nil, iotstore.ErrDefaultAuthorizerNotFound
	}
	return map[string]interface{}{
		"authorizerName": name,
	}, nil
}

func (s *IoTService) TestInvokeAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "authorizerName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}
	return map[string]interface{}{
		"isAuthenticated":       true,
		"principalId":           "test-principal-" + uuid.New().String()[:8],
		"policyDocuments":       []map[string]interface{}{},
		"refreshAfterInSeconds": 300,
	}, nil
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
