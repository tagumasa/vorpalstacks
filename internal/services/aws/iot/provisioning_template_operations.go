package iot

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sort"
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
	if dvid := request.GetIntParam(req.Parameters, "defaultVersionId"); dvid > 0 {
		opts.DefaultVersionID = int64(dvid)
	}
	if hook := request.GetParamCaseInsensitive(req.Parameters, "preProvisioningHook"); hook != "" {
		opts.PreProvisioningHook = hook
	}
	if request.GetBoolParam(req.Parameters, "removePreProvisioningHook") {
		opts.RemovePreProvisioningHook = true
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
		if _, scanErr := fmt.Sscanf(v.VersionID, "%d", &vidInt); scanErr != nil {
			slog.Warn("provisioning template version ID parse failed", "versionID", v.VersionID, "error", scanErr)
			continue
		}
		items = append(items, map[string]interface{}{
			"versionId":    vidInt,
			"templateName": name,
			"createdAt":    v.CreationDate.UTC().Unix(),
		})
	}

	// The store iterates keys lexicographically ("1", "10", "11", "2");
	// version IDs are numeric and must be listed in numeric order.
	sort.Slice(items, func(i, j int) bool {
		return items[i]["versionId"].(int32) < items[j]["versionId"].(int32)
	})

	return paginatedMaps("versions", items, req.Parameters), nil
}

// ---------------------------------------------------------------------------
// ProvisioningTemplateVersion operations.
// ---------------------------------------------------------------------------

func (s *IoTService) CreateProvisioningTemplateVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	templateName := request.GetParamCaseInsensitive(req.Parameters, "templateName")
	if templateName == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if tmpl, err := store.GetProvisioningTemplate(templateName); err != nil {
		return nil, err
	} else if tmpl == nil || tmpl.TemplateName == "" {
		return nil, iotstore.ErrTemplateNotFound
	}
	// Determine the next version ID by scanning existing versions.
	existing, err := store.ListProvisioningTemplateVersions(templateName, parseListOptions(req.Parameters))
	if err != nil {
		return nil, err
	}
	maxVersion := 0
	for _, v := range existing {
		var n int
		fmt.Sscanf(v.VersionID, "%d", &n)
		if n > maxVersion {
			maxVersion = n
		}
	}
	versionID := fmt.Sprintf("%d", maxVersion+1)
	v := &iotstore.ProvisioningTemplateVersion{
		VersionID:        versionID,
		TemplateBody:     request.GetParamCaseInsensitive(req.Parameters, "templateBody"),
		IsDefaultVersion: request.GetBoolParam(req.Parameters, "setAsDefault"),
		CreationDate:     time.Now().UTC(),
	}
	if _, err := store.CreateProvisioningTemplateVersion(templateName, v); err != nil {
		return nil, err
	}
	if v.IsDefaultVersion {
		// setAsDefault must repoint the template's default version and
		// clear the sibling versions' flags in one step; otherwise the new
		// version claims to be the default while the template still
		// resolves to the previous one, and two versions would report
		// themselves as the default at once.
		if err := store.SetDefaultProvisioningTemplateVersion(templateName, int64(maxVersion+1)); err != nil {
			return nil, err
		}
	}
	versionIDInt := int32(maxVersion + 1)
	return map[string]interface{}{
		"templateArn":      iotstore.BuildProvisioningTemplateARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), templateName),
		"templateName":     templateName,
		"versionId":        versionIDInt,
		"isDefaultVersion": v.IsDefaultVersion,
	}, nil
}

func (s *IoTService) DeleteProvisioningTemplateVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	templateName := request.GetParamCaseInsensitive(req.Parameters, "templateName")
	versionID := request.GetParamCaseInsensitive(req.Parameters, "versionId")
	if templateName == "" || versionID == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetProvisioningTemplateVersion(templateName, versionID); err != nil {
		return nil, err
	}
	if err := store.DeleteProvisioningTemplateVersion(templateName, versionID); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) DescribeProvisioningTemplateVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	templateName := request.GetParamCaseInsensitive(req.Parameters, "templateName")
	versionID := request.GetParamCaseInsensitive(req.Parameters, "versionId")
	if templateName == "" || versionID == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	v, err := store.GetProvisioningTemplateVersion(templateName, versionID)
	if err != nil {
		return nil, err
	}
	var versionIDInt int32
	fmt.Sscanf(v.VersionID, "%d", &versionIDInt)
	return map[string]interface{}{
		"templateName":     templateName,
		"versionId":        versionIDInt,
		"templateBody":     v.TemplateBody,
		"isDefaultVersion": v.IsDefaultVersion,
		"creationDate":     v.CreationDate.UTC().Unix(),
	}, nil
}
