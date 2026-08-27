package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// CreateProvisioningTemplate creates a fleet provisioning template.
func (s *IoTService) CreateProvisioningTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := s.createProvisioningTemplateCore(store, CreateProvisioningTemplateInput{
		TemplateName:        request.GetParamCaseInsensitive(req.Parameters, "templateName"),
		Description:         request.GetParamCaseInsensitive(req.Parameters, "description"),
		TemplateBody:        request.GetParamCaseInsensitive(req.Parameters, "templateBody"),
		Enabled:             request.GetBoolParam(req.Parameters, "enabled"),
		EnabledProvided:     request.HasParam(req.Parameters, "enabled"),
		ProvisioningRoleARN: request.GetParamCaseInsensitive(req.Parameters, "provisioningRoleArn"),
		Type:                request.GetParamCaseInsensitive(req.Parameters, "type"),
		Tags:                tagListParam(req.Parameters),
	})
	if err != nil {
		return nil, err
	}
	// CreateProvisioningTemplateResponse carries only the ARN, name and
	// default version ID.
	return map[string]interface{}{
		"templateArn":      created.TemplateARN,
		"templateName":     created.TemplateName,
		"defaultVersionId": created.DefaultVersionID,
	}, nil
}

// DescribeProvisioningTemplate retrieves a provisioning template.
func (s *IoTService) DescribeProvisioningTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	tmpl, err := s.describeProvisioningTemplateCore(store, request.GetParamCaseInsensitive(req.Parameters, "templateName"))
	if err != nil {
		return nil, err
	}
	return provisioningTemplateDetailResponse(tmpl), nil
}

// UpdateProvisioningTemplate modifies a provisioning template.
func (s *IoTService) UpdateProvisioningTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := UpdateProvisioningTemplateInput{
		TemplateName:              request.GetParamCaseInsensitive(req.Parameters, "templateName"),
		Description:               request.GetParamCaseInsensitive(req.Parameters, "description"),
		RoleARN:                   request.GetParamCaseInsensitive(req.Parameters, "provisioningRoleArn"),
		TemplateBody:              request.GetParamCaseInsensitive(req.Parameters, "templateBody"),
		PreProvisioningHook:       request.GetParamCaseInsensitive(req.Parameters, "preProvisioningHook"),
		RemovePreProvisioningHook: request.GetBoolParam(req.Parameters, "removePreProvisioningHook"),
	}
	if request.HasParam(req.Parameters, "enabled") {
		v := request.GetBoolParam(req.Parameters, "enabled")
		in.Enabled = &v
	}
	if dvid := request.GetIntParam(req.Parameters, "defaultVersionId"); dvid > 0 {
		in.DefaultVersionID = int64(dvid)
	}
	if err := s.updateProvisioningTemplateCore(store, in); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

// DeleteProvisioningTemplate removes a provisioning template.
func (s *IoTService) DeleteProvisioningTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteProvisioningTemplateCore(store, request.GetParamCaseInsensitive(req.Parameters, "templateName")); err != nil {
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
	items, nextMarker, err := s.listProvisioningTemplatesCore(store, parseListOptions(req.Parameters))
	if err != nil {
		return nil, err
	}
	result := make([]map[string]interface{}, 0, len(items))
	for _, t := range items {
		result = append(result, provisioningTemplateSummaryResponse(t))
	}
	return listResponse("templates", result, nextMarker), nil
}

// ListProvisioningTemplateVersions returns all versions of a template.
func (s *IoTService) ListProvisioningTemplateVersions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	versions, err := s.listProvisioningTemplateVersionsCore(store,
		request.GetParamCaseInsensitive(req.Parameters, "templateName"),
		parseListOptions(req.Parameters))
	if err != nil {
		return nil, err
	}
	items := make([]map[string]interface{}, 0, len(versions))
	for _, v := range versions {
		items = append(items, map[string]interface{}{
			"versionId":        v.VersionID,
			"creationDate":     v.CreationDate,
			"isDefaultVersion": v.IsDefaultVersion,
		})
	}
	return paginatedMaps("versions", items, req.Parameters)
}

// ---------------------------------------------------------------------------
// ProvisioningTemplateVersion operations.
// ---------------------------------------------------------------------------

func (s *IoTService) CreateProvisioningTemplateVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.createProvisioningTemplateVersionCore(store, CreateProvisioningTemplateVersionInput{
		TemplateName: request.GetParamCaseInsensitive(req.Parameters, "templateName"),
		TemplateBody: request.GetParamCaseInsensitive(req.Parameters, "templateBody"),
		SetAsDefault: request.GetBoolParam(req.Parameters, "setAsDefault"),
		ListOpts:     parseListOptions(req.Parameters),
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"templateArn":      result.TemplateArn,
		"templateName":     result.TemplateName,
		"versionId":        result.VersionID,
		"isDefaultVersion": result.IsDefaultVersion,
	}, nil
}

func (s *IoTService) DeleteProvisioningTemplateVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteProvisioningTemplateVersionCore(store,
		request.GetParamCaseInsensitive(req.Parameters, "templateName"),
		request.GetParamCaseInsensitive(req.Parameters, "versionId")); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) DescribeProvisioningTemplateVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.describeProvisioningTemplateVersionCore(store,
		request.GetParamCaseInsensitive(req.Parameters, "templateName"),
		request.GetParamCaseInsensitive(req.Parameters, "versionId"))
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"templateName":     result.TemplateName,
		"versionId":        result.VersionID,
		"templateBody":     result.TemplateBody,
		"isDefaultVersion": result.IsDefaultVersion,
		"creationDate":     result.CreationDate,
	}, nil
}
