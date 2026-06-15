package iot

import (
	"context"
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
		CreationDate:        time.Now().UTC(),
		LastModifiedDate:    time.Now().UTC(),
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

	opts := iotstore.ProvisioningTemplateUpdateOpts{
		Description: request.GetParamCaseInsensitive(req.Parameters, "description"),
		RoleARN:     request.GetParamCaseInsensitive(req.Parameters, "provisioningRoleArn"),
	}
	if request.HasParam(req.Parameters, "enabled") {
		v := request.GetBoolParam(req.Parameters, "enabled")
		opts.Enabled = &v
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
		items = append(items, map[string]interface{}{
			"versionId":    v.VersionID,
			"templateName": name,
			"createdAt":    v.CreationDate.UTC().Unix(),
		})
	}

	return map[string]interface{}{
		"versions": items,
	}, nil
}
