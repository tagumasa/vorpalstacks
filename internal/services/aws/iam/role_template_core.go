// Transport-agnostic Core functions for the IAM role template operations:
// validation, template-parameter substitution and store operations shared by
// the AWS-compatible HTTP API handlers and any admin surface (the xxxCore
// pattern).
package iam

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/aws/arn"
)

// AcquireRoleInput carries the parsed AcquireRole request.
type AcquireRoleInput struct {
	TemplateArn          string
	TemplateMinorVersion int
	ReplacementValues    map[string][]string
}

// GetRoleTemplateVersionInput carries the parsed GetRoleTemplateVersion
// request. A MinorVersion of 0 means the template's default minor version.
type GetRoleTemplateVersionInput struct {
	TemplateArn  string
	MinorVersion int
}

// acquireRoleCore validates input and materialises a role from the addressed
// role template version: every templated field (name, path, description,
// trust policy, tags) is parameter-substituted with the request's
// replacement values, then the role is created through the standard role
// creation path with the template's inline policies, managed policies and
// permissions boundary applied on top.
func (s *IAMService) acquireRoleCore(store *iamstore.IAMStore, input *AcquireRoleInput) (*iamstore.Role, error) {
	if input.TemplateArn == "" {
		return nil, NewValidationError("TemplateArn")
	}
	if len(input.ReplacementValues) > iamstore.MaxRoleTemplateReplacementEntries {
		return nil, NewInvalidInputError("ReplacementValues",
			fmt.Sprintf("exceeds the maximum of %d replacement entries", iamstore.MaxRoleTemplateReplacementEntries))
	}
	for name, values := range input.ReplacementValues {
		if len(values) < iamstore.MinRoleTemplateReplacementValues || len(values) > iamstore.MaxRoleTemplateReplacementValues {
			return nil, NewInvalidInputError("ReplacementValues",
				fmt.Sprintf("parameter %s must carry between %d and %d values", name,
					iamstore.MinRoleTemplateReplacementValues, iamstore.MaxRoleTemplateReplacementValues))
		}
	}

	template, err := store.RoleTemplates().Get(input.TemplateArn, input.TemplateMinorVersion)
	if err != nil {
		if errors.Is(err, iamstore.ErrRoleTemplateNotFound) || errors.Is(err, iamstore.ErrRoleTemplateVersionNotFound) {
			return nil, NewNoSuchRoleTemplateError(input.TemplateArn)
		}
		return nil, err
	}
	if !template.Enabled || !template.VersionEnabled {
		return nil, NewRoleTemplateDisabledError(input.TemplateArn)
	}

	// Plain-text patterns (name, path, description, tags) substitute the
	// value verbatim: each resolved output passes its own field validator,
	// which rejects any resolution that breaks the field rules. member
	// names the pattern member under resolution, so a resolution fault is
	// reported against the pattern that carries it.
	resolveText := func(member, pattern string) (string, error) {
		return substituteTemplateParameters(pattern, member, template.ParametersDefinition, input.ReplacementValues, false)
	}
	// Policy-document templates substitute the value as JSON string
	// content: a replacement value is data inside the document, never
	// document structure, so an AcquireRole caller cannot exceed the
	// template's policy intent by crafting values.
	resolveDocument := func(member, pattern string) (string, error) {
		return substituteTemplateParameters(pattern, member, template.ParametersDefinition, input.ReplacementValues, true)
	}

	roleName, err := resolveText("RoleNamePattern", template.RoleNamePattern)
	if err != nil {
		return nil, err
	}

	path := template.RolePathPattern
	if path == "" {
		path = "/"
	}
	if path, err = resolveText("RolePathPattern", path); err != nil {
		return nil, err
	}
	// A pattern that resolves to the empty string is invalid on this path:
	// the shared creation sequence would otherwise default it to "/".
	if path == "" {
		return nil, templateRoleCreationMessages.pathInvalid
	}

	description := ""
	if template.RoleDescriptionPattern != "" {
		if description, err = resolveText("RoleDescriptionPattern", template.RoleDescriptionPattern); err != nil {
			return nil, err
		}
	}

	trustPolicy, err := resolveDocument("AssumeRolePolicyDocumentTemplate", template.AssumeRolePolicyDocumentTemplate)
	if err != nil {
		return nil, err
	}

	var roleTags []tagutil.Tag
	for _, tagTemplate := range template.RoleTagsTemplate {
		key, err := resolveText("RoleTagsTemplate", tagTemplate.Key)
		if err != nil {
			return nil, err
		}
		value, err := resolveText("RoleTagsTemplate", tagTemplate.Value)
		if err != nil {
			return nil, err
		}
		roleTags = append(roleTags, tagutil.Tag{Key: key, Value: value})
	}

	// The resolved values enter the standard creation sequence; name,
	// trust-policy, tag and duration predicates are checked there, with the
	// pattern-member error shapes this operation documents. The creation
	// records the template and the minor version actually resolved — a
	// request that selected the default reports the default's version.
	role, err := createValidatedRole(store, &CreateRoleInput{
		RoleName:                 roleName,
		Path:                     path,
		AssumeRolePolicyDocument: trustPolicy,
		Description:              description,
		MaxSessionDuration:       template.MaxSessionDuration,
		Tags:                     roleTags,
		SourceRoleTemplate: &iamstore.SourceRoleTemplate{
			TemplateArn:          input.TemplateArn,
			TemplateMinorVersion: template.MinorVersion,
		},
	}, templateRoleCreationMessages)
	if err != nil {
		return nil, err
	}

	if err := s.applyRoleTemplateLayers(store, role, template, input.ReplacementValues); err != nil {
		return nil, err
	}

	return role, nil
}

// applyRoleTemplateLayers applies the template's post-creation layers —
// inline policies, managed-policy attachments and the permissions
// boundary — to the freshly created role, all-or-nothing: a failure
// applying any layer reverses the layers applied so far and deletes the
// role, so retrying AcquireRole (the documented remedy for a mid-creation
// rejection) never collides with a half-provisioned role name. A reverse
// failure leaves partial state, so it is logged and joined into the
// returned error rather than discarded; the original fault stays first in
// the join so the wire mapping still resolves to it.
func (s *IAMService) applyRoleTemplateLayers(store *iamstore.IAMStore, role *iamstore.Role, template *iamstore.RoleTemplateVersion, replacementValues map[string][]string) error {
	roleName := role.RoleName
	resolveDocument := func(member, pattern string) (string, error) {
		return substituteTemplateParameters(pattern, member, template.ParametersDefinition, replacementValues, true)
	}

	var appliedInlinePolicies []string
	var appliedPolicyArns []string
	boundaryApplied := false
	rollback := func(failErr error) error {
		rollbackErrs := []error{failErr}
		if boundaryApplied {
			if err := s.deleteRolePermissionsBoundaryCore(store, roleName); err != nil {
				logs.Warn("iam: permissions boundary removal failed during AcquireRole rollback",
					logs.String("roleName", roleName), logs.Err(err))
				rollbackErrs = append(rollbackErrs, err)
			}
		}
		for i := len(appliedPolicyArns) - 1; i >= 0; i-- {
			if err := s.detachPolicyCore(store, &AttachPolicyInput{
				PrincipalType: PrincipalTypeRole,
				PrincipalName: roleName,
				PolicyArn:     appliedPolicyArns[i],
			}); err != nil {
				logs.Warn("iam: managed policy detach failed during AcquireRole rollback",
					logs.String("roleName", roleName), logs.String("policyArn", appliedPolicyArns[i]), logs.Err(err))
				rollbackErrs = append(rollbackErrs, err)
			}
		}
		for i := len(appliedInlinePolicies) - 1; i >= 0; i-- {
			if err := s.deleteInlinePolicyCore(store, PrincipalTypeRole, roleName, appliedInlinePolicies[i]); err != nil {
				logs.Warn("iam: inline policy removal failed during AcquireRole rollback",
					logs.String("roleName", roleName), logs.String("policyName", appliedInlinePolicies[i]), logs.Err(err))
				rollbackErrs = append(rollbackErrs, err)
			}
		}
		if err := store.Roles().Delete(roleName); err != nil {
			logs.Warn("iam: role deletion failed during AcquireRole rollback",
				logs.String("roleName", roleName), logs.Err(err))
			rollbackErrs = append(rollbackErrs, err)
		}
		return errors.Join(rollbackErrs...)
	}

	for _, inlineTemplate := range template.InlinePolicyTemplates {
		document, err := resolveDocument("InlinePolicyTemplates", inlineTemplate.PolicyDocument)
		if err != nil {
			return rollback(err)
		}
		if err := s.putInlinePolicyCore(store, &PutInlinePolicyInput{
			PrincipalType:  PrincipalTypeRole,
			PrincipalName:  roleName,
			PolicyName:     inlineTemplate.PolicyName,
			PolicyDocument: document,
		}); err != nil {
			return rollback(err)
		}
		appliedInlinePolicies = append(appliedInlinePolicies, inlineTemplate.PolicyName)
	}

	for _, policyArn := range template.ManagedPolicyArns {
		if err := s.attachPolicyCore(store, &AttachPolicyInput{
			PrincipalType: PrincipalTypeRole,
			PrincipalName: roleName,
			PolicyArn:     policyArn,
		}); err != nil {
			return rollback(err)
		}
		appliedPolicyArns = append(appliedPolicyArns, policyArn)
	}

	if template.PermissionBoundaryArn != "" {
		if err := attachRolePermissionsBoundaryCore(store, role, template.PermissionBoundaryArn); err != nil {
			return rollback(err)
		}
		boundaryApplied = true
	}

	return nil
}

// getRoleTemplateVersionCore validates input and resolves the addressed
// role template version.
func (s *IAMService) getRoleTemplateVersionCore(store *iamstore.IAMStore, input *GetRoleTemplateVersionInput) (*iamstore.RoleTemplateVersion, error) {
	if input.TemplateArn == "" {
		return nil, NewValidationError("TemplateArn")
	}

	template, err := store.RoleTemplates().Get(input.TemplateArn, input.MinorVersion)
	if err != nil {
		if errors.Is(err, iamstore.ErrRoleTemplateNotFound) || errors.Is(err, iamstore.ErrRoleTemplateVersionNotFound) {
			return nil, NewNoSuchRoleTemplateError(input.TemplateArn)
		}
		return nil, err
	}
	return template, nil
}

// substituteTemplateParameters replaces every @{name} placeholder in
// pattern with its template parameter value: the first supplied replacement
// value, else the parameter's documented default, else — for a parameter
// that is neither supplied nor defined — an InvalidInput error, because the
// pattern cannot resolve without it. Every resolved value is checked
// against the parameter's declared type. With inDocument set the value is
// additionally escaped as JSON string content, because the pattern is a
// policy document and the value must remain data inside it. member names
// the pattern member under resolution, so an unterminated placeholder is
// reported against the pattern that carries it.
func substituteTemplateParameters(pattern, member string, definitions []iamstore.RoleTemplateParameterDefinition, values map[string][]string, inDocument bool) (string, error) {
	if !strings.Contains(pattern, "@{") {
		return pattern, nil
	}

	var builder strings.Builder
	for {
		start := strings.Index(pattern, "@{")
		if start < 0 {
			break
		}
		end := strings.Index(pattern[start:], "}")
		if end < 0 {
			return "", NewInvalidInputError(member, "template carries an unterminated parameter placeholder")
		}
		name := pattern[start+2 : start+end]

		replacement, err := templateParameterValue(name, definitions, values)
		if err != nil {
			return "", err
		}
		if inDocument {
			replacement = escapeJSONStringContent(replacement)
		}
		builder.WriteString(pattern[:start])
		builder.WriteString(replacement)
		pattern = pattern[start+end+1:]
	}
	builder.WriteString(pattern)
	return builder.String(), nil
}

// templateParameterValue resolves one template parameter to its replacement
// string. A supplied value wins (the first list element — the AcquireRole
// wire format carries a value list per parameter), then the parameter's
// default value; an unsupplied parameter with no default and no value is an
// InvalidInput error. Both the supplied value and the default must satisfy
// the parameter's declared type.
func templateParameterValue(name string, definitions []iamstore.RoleTemplateParameterDefinition, values map[string][]string) (string, error) {
	var definition *iamstore.RoleTemplateParameterDefinition
	for i := range definitions {
		if definitions[i].Name == name {
			definition = &definitions[i]
			break
		}
	}
	if entry, ok := values[name]; ok {
		if len(entry) == 0 {
			return "", NewInvalidInputError("ReplacementValues", "parameter "+name+" carries an empty value list")
		}
		if err := validateTemplateParameterValue(name, definition, entry[0]); err != nil {
			return "", err
		}
		return entry[0], nil
	}
	if definition != nil {
		if definition.DefaultValue != "" {
			if err := validateTemplateParameterValue(name, definition, definition.DefaultValue); err != nil {
				return "", err
			}
			return definition.DefaultValue, nil
		}
		if !definition.IsRequired {
			return "", nil
		}
	}
	return "", NewInvalidInputError("ReplacementValues", "no value supplied for template parameter "+name)
}

// templateNumberPattern is the numeric grammar a Number (or NumberList)
// parameter value must satisfy: the JSON number grammar, since the model
// defines the type only by name and every value lands in a string context.
var templateNumberPattern = regexp.MustCompile(`^[+-]?(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+)(?:[eE][+-]?[0-9]+)?$`)

// validateTemplateParameterValue checks one resolved value against the
// parameter's declared type — the documented constraint a template places
// on what its parameters accept. Number values must be numeric and Arn
// values must parse as ARNs; String values (and parameters without a
// declared type) are unconstrained on the wire. A type outside the
// documented enum means the template itself is invalid.
func validateTemplateParameterValue(name string, definition *iamstore.RoleTemplateParameterDefinition, value string) error {
	if definition == nil {
		return nil
	}
	switch definition.Type {
	case "", "String", "StringList":
		return nil
	case "Number", "NumberList":
		if !templateNumberPattern.MatchString(value) {
			return NewInvalidInputError("ReplacementValues", "parameter "+name+" requires a numeric value")
		}
	case "Arn", "ArnList":
		if _, err := arn.ParseARN(value); err != nil {
			return NewInvalidInputError("ReplacementValues", "parameter "+name+" requires an ARN value")
		}
	default:
		return NewInvalidInputError("ReplacementValues", "parameter "+name+" declares the unsupported type "+definition.Type)
	}
	return nil
}

// escapeJSONStringContent escapes s so that splicing it inside a JSON
// string literal yields exactly s when the document is decoded: quote,
// backslash and control bytes are encoded, every other byte — including
// multi-byte UTF-8 — passes through verbatim. It never renders an invalid
// value well-formed: bytes that are not valid UTF-8 pass through and the
// resulting document fails parsing downstream.
func escapeJSONStringContent(s string) string {
	var builder strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '"':
			builder.WriteString(`\"`)
		case '\\':
			builder.WriteString(`\\`)
		case '\n':
			builder.WriteString(`\n`)
		case '\r':
			builder.WriteString(`\r`)
		case '\t':
			builder.WriteString(`\t`)
		case '\b':
			builder.WriteString(`\b`)
		case '\f':
			builder.WriteString(`\f`)
		default:
			if c < 0x20 {
				fmt.Fprintf(&builder, `\u%04x`, c)
				continue
			}
			builder.WriteByte(c)
		}
	}
	return builder.String()
}
