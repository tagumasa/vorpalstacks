package iam

import (
	"context"
	"fmt"

	"vorpalstacks/internal/common/request"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/timeutils"
)

// AcquireRole creates an IAM role from a role template version. The role's
// configuration comes from the template; ReplacementValues supplies the
// values substituted into the template's parameters.
func (s *IAMService) AcquireRole(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	replacementValues := map[string][]string{}
	for i := 1; ; i++ {
		name := request.GetStringParam(req.Parameters, fmt.Sprintf("ReplacementValues.entry.%d.key", i))
		if name == "" {
			break
		}
		var values []string
		for j := 1; ; j++ {
			value := request.GetStringParam(req.Parameters, fmt.Sprintf("ReplacementValues.entry.%d.value.Values.member.%d", i, j))
			if value == "" {
				break
			}
			values = append(values, value)
		}
		replacementValues[name] = values
	}

	input := &AcquireRoleInput{
		TemplateArn:          request.GetStringParam(req.Parameters, "TemplateArn"),
		TemplateMinorVersion: request.GetIntParam(req.Parameters, "TemplateMinorVersion"),
		ReplacementValues:    replacementValues,
	}
	role, err := s.acquireRoleCore(store, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Role": roleToResponse(role),
	}, nil
}

// GetRoleTemplateVersion retrieves a version of a role template. Without
// MinorVersion the template's default minor version is returned.
func (s *IAMService) GetRoleTemplateVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	input := &GetRoleTemplateVersionInput{
		TemplateArn:  request.GetStringParam(req.Parameters, "TemplateArn"),
		MinorVersion: request.GetIntParam(req.Parameters, "MinorVersion"),
	}
	template, err := s.getRoleTemplateVersionCore(store, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"RoleTemplateVersion": roleTemplateVersionToResponse(template),
	}, nil
}

// roleTemplateVersionToResponse renders a role template version in the
// RoleTemplateVersion wire shape; absent members are omitted, mirroring the
// documented response in which only the template's set members appear.
func roleTemplateVersionToResponse(template *iamstore.RoleTemplateVersion) map[string]interface{} {
	resp := map[string]interface{}{
		"TemplateArn":         template.TemplateArn,
		"TemplateName":        template.TemplateName,
		"MajorVersion":        template.MajorVersion,
		"MinorVersion":        template.MinorVersion,
		"DefaultMinorVersion": template.DefaultMinorVersion,
		"Enabled":             template.Enabled,
		"RoleNamePattern":     template.RoleNamePattern,
		"CreateTimestamp":     template.CreateTimestamp.Format(timeutils.ISO8601SimpleFormat),
	}

	if template.TemplateVersionId != "" {
		resp["TemplateVersionId"] = template.TemplateVersionId
	}
	if template.Description != "" {
		resp["Description"] = template.Description
	}
	if template.ManagedByType != "" {
		resp["ManagedByType"] = template.ManagedByType
	}
	if template.ManagedByValue != "" {
		resp["ManagedByValue"] = template.ManagedByValue
	}
	if template.RolePathPattern != "" {
		resp["RolePathPattern"] = template.RolePathPattern
	}
	if template.RoleDescriptionPattern != "" {
		resp["RoleDescriptionPattern"] = template.RoleDescriptionPattern
	}
	if template.AssumeRolePolicyDocumentTemplate != "" {
		resp["AssumeRolePolicyDocumentTemplate"] = template.AssumeRolePolicyDocumentTemplate
	}
	if len(template.InlinePolicyTemplates) > 0 {
		inlinePolicies := make([]map[string]interface{}, 0, len(template.InlinePolicyTemplates))
		for _, inlinePolicy := range template.InlinePolicyTemplates {
			inlinePolicies = append(inlinePolicies, map[string]interface{}{
				"PolicyName":     inlinePolicy.PolicyName,
				"PolicyDocument": inlinePolicy.PolicyDocument,
			})
		}
		resp["InlinePolicyTemplates"] = inlinePolicies
	}
	if len(template.ManagedPolicyArns) > 0 {
		resp["ManagedPolicyArns"] = template.ManagedPolicyArns
	}
	if template.PermissionBoundaryArn != "" {
		resp["PermissionBoundaryArn"] = template.PermissionBoundaryArn
	}
	if len(template.ParametersDefinition) > 0 {
		parameters := make([]map[string]interface{}, 0, len(template.ParametersDefinition))
		for _, definition := range template.ParametersDefinition {
			entry := map[string]interface{}{
				"Name":       definition.Name,
				"Type":       definition.Type,
				"IsRequired": definition.IsRequired,
			}
			if definition.SubType != "" {
				entry["SubType"] = definition.SubType
			}
			if definition.Description != "" {
				entry["Description"] = definition.Description
			}
			if definition.DefaultValue != "" {
				entry["DefaultValue"] = definition.DefaultValue
			}
			entry["Immutable"] = definition.Immutable
			parameters = append(parameters, entry)
		}
		resp["ParametersDefinition"] = parameters
	}
	if len(template.RoleTagsTemplate) > 0 {
		tagTemplates := make([]map[string]interface{}, 0, len(template.RoleTagsTemplate))
		for _, tagTemplate := range template.RoleTagsTemplate {
			tagTemplates = append(tagTemplates, map[string]interface{}{
				"Key":   tagTemplate.Key,
				"Value": tagTemplate.Value,
			})
		}
		resp["RoleTagsTemplate"] = tagTemplates
	}
	if template.MaxSessionDuration > 0 {
		resp["MaxSessionDuration"] = template.MaxSessionDuration
	}
	if template.VersionEnabled {
		resp["VersionEnabled"] = template.VersionEnabled
	}
	if !template.UpdateTimestamp.IsZero() {
		resp["UpdateTimestamp"] = template.UpdateTimestamp.Format(timeutils.ISO8601SimpleFormat)
	}

	return resp
}
