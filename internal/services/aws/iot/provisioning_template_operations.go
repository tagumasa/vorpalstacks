package iot

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

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

	tmpl := &iotstore.ProvisioningTemplate{
		TemplateName:        name,
		Description:         request.GetParamCaseInsensitive(req.Parameters, "description"),
		Enabled:             true,
		ProvisioningRoleARN: request.GetParamCaseInsensitive(req.Parameters, "provisioningRoleArn"),
		TemplateBody:        request.GetParamCaseInsensitive(req.Parameters, "templateBody"),
		Type:                request.GetParamCaseInsensitive(req.Parameters, "type"),
		CreationDate:        time.Now().UTC(),
		LastModifiedDate:    time.Now().UTC(),
	}

	// Validate template body is well-formed JSON before storing.
	if tmpl.TemplateBody != "" {
		var v interface{}
		if err := json.Unmarshal([]byte(tmpl.TemplateBody), &v); err != nil {
			return nil, fmt.Errorf("invalid templateBody: %w", err)
		}
	}

	if request.HasParam(req.Parameters, "enabled") {
		tmpl.Enabled = request.GetBoolParam(req.Parameters, "enabled")
	}

	created, err := store.CreateProvisioningTemplate(tmpl)
	if err != nil {
		return nil, err
	}

	return provisioningTemplateResponse(created), nil
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
		return nil, err
	}

	return provisioningTemplateResponse(tmpl), nil
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

	opts := iotstore.ProvisioningTemplateUpdateOpts{
		Description: request.GetParamCaseInsensitive(req.Parameters, "description"),
		RoleARN:     request.GetParamCaseInsensitive(req.Parameters, "provisioningRoleArn"),
	}
	if request.HasParam(req.Parameters, "enabled") {
		v := request.GetBoolParam(req.Parameters, "enabled")
		opts.Enabled = &v
	}
	if tb := request.GetParamCaseInsensitive(req.Parameters, "templateBody"); tb != "" {
		var v interface{}
		if err := json.Unmarshal([]byte(tb), &v); err != nil {
			return nil, fmt.Errorf("invalid templateBody: %w", err)
		}
		opts.TemplateBody = &tb
	}

	_, err = store.UpdateProvisioningTemplate(name, opts)
	if err != nil {
		return nil, err
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

	arn := iotstore.BuildProvisioningTemplateARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), name)
	_ = store.DeleteAllTags(arn)

	if err := store.DeleteProvisioningTemplate(name); err != nil {
		return nil, err
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
		result = append(result, provisioningTemplateResponse(t))
	}

	return listResponse("templates", result, templates.NextMarker), nil
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

	versions, err := store.ListProvisioningTemplateVersions(name, parseListOptions(req.Parameters))
	if err != nil {
		return nil, err
	}

	items := make([]map[string]interface{}, 0, len(versions))
	for _, v := range versions {
		var vidInt int32
		fmt.Sscanf(v.VersionID, "%d", &vidInt)
		items = append(items, map[string]interface{}{
			"versionId":    vidInt,
			"templateName": name,
			"createdAt":    v.CreationDate.UTC().Unix(),
		})
	}

	return paginatedMaps("versions", items, req.Parameters), nil
}
