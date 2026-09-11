package iam

import (
	"time"
)

// RoleTemplateParameterDefinition describes one parameter a role template
// version accepts (RoleTemplateVersion.ParametersDefinition member).
type RoleTemplateParameterDefinition struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	SubType      string `json:"sub_type,omitempty"`
	Description  string `json:"description,omitempty"`
	IsRequired   bool   `json:"is_required,omitempty"`
	DefaultValue string `json:"default_value,omitempty"`
	Immutable    bool   `json:"immutable,omitempty"`
}

// RoleTemplateInlinePolicy is one inline policy template embedded in roles
// created from the template (InlinePolicyTemplates member).
type RoleTemplateInlinePolicy struct {
	PolicyName     string `json:"policy_name"`
	PolicyDocument string `json:"policy_document"`
}

// RoleTemplateTag is one tag template applied to roles created from the
// template; Key and Value may carry @{parameter} placeholders.
type RoleTemplateTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// RoleTemplateVersion is one version of a role template: the reusable role
// configuration that AcquireRole materialises roles from.
type RoleTemplateVersion struct {
	TemplateArn                      string                            `json:"template_arn"`
	TemplateName                     string                            `json:"template_name"`
	TemplateVersionId                string                            `json:"template_version_id,omitempty"`
	Description                      string                            `json:"description,omitempty"`
	MajorVersion                     int                               `json:"major_version"`
	DefaultMinorVersion              int                               `json:"default_minor_version"`
	ManagedByType                    string                            `json:"managed_by_type,omitempty"`
	ManagedByValue                   string                            `json:"managed_by_value,omitempty"`
	Enabled                          bool                              `json:"enabled"`
	MinorVersion                     int                               `json:"minor_version"`
	RoleNamePattern                  string                            `json:"role_name_pattern"`
	RolePathPattern                  string                            `json:"role_path_pattern,omitempty"`
	RoleDescriptionPattern           string                            `json:"role_description_pattern,omitempty"`
	AssumeRolePolicyDocumentTemplate string                            `json:"assume_role_policy_document_template,omitempty"`
	InlinePolicyTemplates            []RoleTemplateInlinePolicy        `json:"inline_policy_templates,omitempty"`
	ManagedPolicyArns                []string                          `json:"managed_policy_arns,omitempty"`
	PermissionBoundaryArn            string                            `json:"permission_boundary_arn,omitempty"`
	ParametersDefinition             []RoleTemplateParameterDefinition `json:"parameters_definition,omitempty"`
	RoleTagsTemplate                 []RoleTemplateTag                 `json:"role_tags_template,omitempty"`
	MaxSessionDuration               int                               `json:"max_session_duration,omitempty"`
	VersionEnabled                   bool                              `json:"version_enabled"`
	CreateTimestamp                  time.Time                         `json:"create_timestamp"`
	UpdateTimestamp                  time.Time                         `json:"update_timestamp,omitempty"`
}

// roleTemplateCatalogue holds the platform-hosted role templates. Entries
// carry only content AWS documents: the Example template below reproduces
// the API Reference examples of GetRoleTemplateVersion (all shown members,
// including the exact trust policy document and name pattern) and AcquireRole
// (whose response pins the created role's path). Template versions are
// immutable, so the catalogue is versioned data, not mutable state — no
// storage bucket backs it and any additionally documented template is a data
// addition here.
var roleTemplateCatalogue = []RoleTemplateVersion{
	{
		TemplateArn:  "arn:aws:iam::aws:role-template/awsserviceprincipal/Example:1",
		TemplateName: "Example",
		MajorVersion: 1,
		MinorVersion: 1,
		// Minor 1 is both the first and the default minor version.
		DefaultMinorVersion:              1,
		Enabled:                          true,
		VersionEnabled:                   true,
		RoleNamePattern:                  "Example-@{Department}",
		RolePathPattern:                  "/awsserviceprincipal/",
		AssumeRolePolicyDocumentTemplate: `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":["ec2.amazonaws.com"]},"Action":["sts:AssumeRole"]}]}`,
		ParametersDefinition: []RoleTemplateParameterDefinition{
			{
				Name:       "Department",
				Type:       "String",
				IsRequired: true,
			},
		},
		CreateTimestamp: time.Date(2026, 6, 26, 20, 43, 32, 0, time.UTC),
	},
}

// RoleTemplateStore serves the role template catalogue. Template lookups key
// on the exact template ARN (including its major-version suffix); the minor
// version selects among the template's versions, defaulting to the
// template's default minor version.
type RoleTemplateStore struct {
	catalogue []RoleTemplateVersion
}

// NewRoleTemplateStore creates a new RoleTemplateStore over the seeded
// catalogue.
func NewRoleTemplateStore() *RoleTemplateStore {
	return &RoleTemplateStore{catalogue: roleTemplateCatalogue}
}

// Get resolves the template version addressed by templateArn and
// minorVersion among every catalogue entry of the ARN: a minorVersion of 0
// selects the template's default minor version. A template ARN absent from
// the catalogue yields ErrRoleTemplateNotFound; a known template without
// the requested minor version yields ErrRoleTemplateVersionNotFound.
func (s *RoleTemplateStore) Get(templateArn string, minorVersion int) (*RoleTemplateVersion, error) {
	var matches []*RoleTemplateVersion
	for i := range s.catalogue {
		if s.catalogue[i].TemplateArn == templateArn {
			matches = append(matches, &s.catalogue[i])
		}
	}
	if len(matches) == 0 {
		return nil, NewStoreError("get_role_template", ErrRoleTemplateNotFound)
	}

	if minorVersion == 0 {
		minorVersion = matches[0].DefaultMinorVersion
	}
	for _, template := range matches {
		if template.MinorVersion == minorVersion {
			// The catalogue slice is immutable after seeding; every lookup
			// returns the seeded entry itself, which keeps the served
			// template identical to its definition site.
			return template, nil
		}
	}
	return nil, NewStoreError("get_role_template", ErrRoleTemplateVersionNotFound)
}
