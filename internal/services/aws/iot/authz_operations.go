package iot

import (
	"context"
	"time"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// CreateAuthorizer creates a custom authorizer for MQTT connections.
func (s *IoTService) CreateAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "authorizerName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetAuthorizer(name); err == nil {
		return nil, iotstore.ErrResourceAlreadyExists
	}

	auth := &iotstore.Authorizer{
		AuthorizerName:        name,
		AuthorizerFunctionARN: request.GetParamCaseInsensitive(req.Parameters, "authorizerFunctionArn"),
		TokenName:             request.GetParamCaseInsensitive(req.Parameters, "tokenKeyName"),
		TokenSignature:        request.GetParamCaseInsensitive(req.Parameters, "tokenSignature"),
		Status:                true,
		EnableCachingForHTTP:  true,
		CreationDate:          time.Now().UTC(),
		LastModifiedDate:      time.Now().UTC(),
	}

	if statusStr := request.GetParamCaseInsensitive(req.Parameters, "status"); statusStr == "INACTIVE" {
		auth.Status = false
	}

	auth.EnableCachingForHTTP = request.GetBoolParam(req.Parameters, "enableCachingForHttp")

	created, err := store.CreateAuthorizer(auth)
	if err != nil {
		return nil, iotstore.ErrInvalidRequest
	}

	status := "INACTIVE"
	if created.Status {
		status = "ACTIVE"
	}
	return map[string]interface{}{
		"authorizerName":        created.AuthorizerName,
		"authorizerArn":         created.AuthorizerARN,
		"authorizerFunctionArn": created.AuthorizerFunctionARN,
		"tokenKeyName":          created.TokenName,
		"status":                status,
		"creationDate":          created.CreationDate.Unix(),
		"lastModifiedDate":      created.LastModifiedDate.Unix(),
	}, nil
}

// DescribeAuthorizer retrieves details of a custom authorizer.
func (s *IoTService) DescribeAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "authorizerName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	auth, err := store.GetAuthorizer(name)
	if err != nil {
		return nil, iotstore.ErrAuthorizerNotFound
	}

	status := "INACTIVE"
	if auth.Status {
		status = "ACTIVE"
	}

	return map[string]interface{}{
		"authorizerName":        auth.AuthorizerName,
		"authorizerArn":         auth.AuthorizerARN,
		"authorizerFunctionArn": auth.AuthorizerFunctionARN,
		"tokenKeyName":          auth.TokenName,
		"tokenSignature":        auth.TokenSignature,
		"status":                status,
		"enableCachingForHttp":  auth.EnableCachingForHTTP,
		"creationDate":          auth.CreationDate.Unix(),
		"lastModifiedDate":      auth.LastModifiedDate.Unix(),
	}, nil
}

// UpdateAuthorizer modifies a custom authorizer configuration.
func (s *IoTService) UpdateAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "authorizerName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	existing, err := store.GetAuthorizer(name)
	if err != nil {
		return nil, iotstore.ErrAuthorizerNotFound
	}

	if fn := request.GetParamCaseInsensitive(req.Parameters, "authorizerFunctionArn"); fn != "" {
		existing.AuthorizerFunctionARN = fn
	}
	if tkn := request.GetParamCaseInsensitive(req.Parameters, "tokenKeyName"); tkn != "" {
		existing.TokenName = tkn
	}
	if sig := request.GetParamCaseInsensitive(req.Parameters, "tokenSignature"); sig != "" {
		existing.TokenSignature = sig
	}
	existing.EnableCachingForHTTP = request.GetBoolParam(req.Parameters, "enableCachingForHttp")
	if statusStr := request.GetParamCaseInsensitive(req.Parameters, "status"); statusStr != "" {
		existing.Status = statusStr == "ACTIVE"
	}

	existing.LastModifiedDate = time.Now().UTC()

	if err := store.UpdateAuthorizer(existing); err != nil {
		return nil, iotstore.ErrInvalidRequest
	}

	return map[string]interface{}{
		"authorizerName": existing.AuthorizerName,
	}, nil
}

// DeleteAuthorizer removes a custom authorizer.
func (s *IoTService) DeleteAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "authorizerName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteAuthorizer(name); err != nil {
		return nil, iotstore.ErrAuthorizerNotFound
	}

	return map[string]interface{}{}, nil
}

// ListAuthorizers returns all custom authorizers.
func (s *IoTService) ListAuthorizers(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	auths, err := store.ListAuthorizers(parseListOptions(req.Parameters))
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0, len(auths.Items))
	for _, a := range auths.Items {
		status := "INACTIVE"
		if a.Status {
			status = "ACTIVE"
		}
		result = append(result, map[string]interface{}{
			"authorizerName":        a.AuthorizerName,
			"authorizerArn":         a.AuthorizerARN,
			"authorizerFunctionArn": a.AuthorizerFunctionARN,
			"tokenKeyName":          a.TokenName,
			"status":                status,
			"creationDate":          a.CreationDate.Unix(),
			"lastModifiedDate":      a.LastModifiedDate.Unix(),
		})
	}

	resp := map[string]interface{}{
		"authorizers": result,
	}
	if auths.NextMarker != "" {
		resp["nextToken"] = auths.NextMarker
	}
	return resp, nil
}

// CreateProvisioningTemplate creates a fleet provisioning template.
func (s *IoTService) CreateProvisioningTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "templateName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetProvisioningTemplate(name); err == nil {
		return nil, iotstore.ErrResourceAlreadyExists
	}

	tmpl := &iotstore.ProvisioningTemplate{
		TemplateName:        name,
		Description:         request.GetParamCaseInsensitive(req.Parameters, "description"),
		Enabled:             true,
		ProvisioningRoleARN: request.GetParamCaseInsensitive(req.Parameters, "provisioningRoleArn"),
		CreationDate:        time.Now().UTC(),
		LastModifiedDate:    time.Now().UTC(),
	}

	tmpl.Enabled = request.GetBoolParam(req.Parameters, "enabled")

	created, err := store.CreateProvisioningTemplate(tmpl)
	if err != nil {
		return nil, iotstore.ErrInvalidRequest
	}

	return map[string]interface{}{
		"templateName":        created.TemplateName,
		"templateArn":         created.TemplateARN,
		"description":         created.Description,
		"enabled":             created.Enabled,
		"provisioningRoleArn": created.ProvisioningRoleARN,
		"creationDate":        created.CreationDate.Unix(),
		"lastModifiedDate":    created.LastModifiedDate.Unix(),
	}, nil
}

// DescribeProvisioningTemplate retrieves a provisioning template.
func (s *IoTService) DescribeProvisioningTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "templateName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	tmpl, err := store.GetProvisioningTemplate(name)
	if err != nil {
		return nil, iotstore.ErrTemplateNotFound
	}

	return map[string]interface{}{
		"templateName":        tmpl.TemplateName,
		"templateArn":         tmpl.TemplateARN,
		"description":         tmpl.Description,
		"enabled":             tmpl.Enabled,
		"provisioningRoleArn": tmpl.ProvisioningRoleARN,
		"creationDate":        tmpl.CreationDate.Unix(),
		"lastModifiedDate":    tmpl.LastModifiedDate.Unix(),
	}, nil
}

// UpdateProvisioningTemplate modifies a provisioning template.
func (s *IoTService) UpdateProvisioningTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "templateName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	existing, err := store.GetProvisioningTemplate(name)
	if err != nil {
		return nil, iotstore.ErrTemplateNotFound
	}

	if desc := request.GetParamCaseInsensitive(req.Parameters, "description"); desc != "" {
		existing.Description = desc
	}
	if role := request.GetParamCaseInsensitive(req.Parameters, "provisioningRoleArn"); role != "" {
		existing.ProvisioningRoleARN = role
	}
	existing.Enabled = request.GetBoolParam(req.Parameters, "enabled")

	existing.LastModifiedDate = time.Now().UTC()

	if err := store.UpdateProvisioningTemplate(existing); err != nil {
		return nil, iotstore.ErrInvalidRequest
	}

	return map[string]interface{}{}, nil
}

// DeleteProvisioningTemplate removes a provisioning template.
func (s *IoTService) DeleteProvisioningTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "templateName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteProvisioningTemplate(name); err != nil {
		return nil, iotstore.ErrTemplateNotFound
	}

	return map[string]interface{}{}, nil
}

// ListProvisioningTemplates returns all provisioning templates.
func (s *IoTService) ListProvisioningTemplates(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	templates, err := store.ListProvisioningTemplates(parseListOptions(req.Parameters))
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0, len(templates.Items))
	for _, t := range templates.Items {
		result = append(result, map[string]interface{}{
			"templateName":        t.TemplateName,
			"templateArn":         t.TemplateARN,
			"description":         t.Description,
			"enabled":             t.Enabled,
			"provisioningRoleArn": t.ProvisioningRoleARN,
			"creationDate":        t.CreationDate.Unix(),
			"lastModifiedDate":    t.LastModifiedDate.Unix(),
		})
	}

	resp := map[string]interface{}{
		"templates": result,
	}
	if templates.NextMarker != "" {
		resp["nextToken"] = templates.NextMarker
	}
	return resp, nil
}

// ListProvisioningTemplateVersions returns all versions of a template.
func (s *IoTService) ListProvisioningTemplateVersions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "templateName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	tmpl, err := store.GetProvisioningTemplate(name)
	if err != nil {
		return nil, iotstore.ErrTemplateNotFound
	}

	return map[string]interface{}{
		"versions": []map[string]interface{}{
			{
				"versionId":    int64(1),
				"templateName": name,
				"createdAt":    tmpl.CreationDate.UTC().Unix(),
			},
		},
	}, nil
}
