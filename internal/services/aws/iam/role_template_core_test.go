package iam

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/storage"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

const trustPolicyPattern = `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":["@{Service}"]},"Action":["sts:AssumeRole"]}]}`

// A replacement value substituted into a policy-document template is data
// inside the document: quoting, escaping and statement-shaped payloads must
// stay inside the placeholder's string position, and the resolved document
// must decode with the value intact and no statement added.
func TestSubstituteTemplateParametersDocumentValuesEscaped(t *testing.T) {
	payload := `ec2.amazonaws.com"],"Effect":"Allow","Action":["*"],"Resource":["*"]},{"Effect":"Deny"`
	defs := []iamstore.RoleTemplateParameterDefinition{{Name: "Service", Type: "String"}}
	values := map[string][]string{"Service": {payload}}

	resolved, err := substituteTemplateParameters(trustPolicyPattern, "AssumeRolePolicyDocumentTemplate", defs, values, true)
	require.NoError(t, err)
	assert.True(t, json.Valid([]byte(resolved)), "resolved document must be well-formed JSON: %s", resolved)

	var doc map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(resolved), &doc))
	statements, ok := doc["Statement"].([]interface{})
	require.True(t, ok, "Statement must remain an array")
	require.Len(t, statements, 1, "no statement may be injected")
	stmt, ok := statements[0].(map[string]interface{})
	require.True(t, ok)

	service, ok := stmt["Principal"].(map[string]interface{})["Service"].([]interface{})
	require.True(t, ok, "Principal.Service must remain an array")
	assert.Equal(t, []interface{}{payload}, service, "the value must land verbatim as string content")

	assert.Equal(t, []interface{}{"sts:AssumeRole"}, stmt["Action"], "the template's own Action is unchanged")
	assert.NotContains(t, stmt, "Resource", "the payload's Resource must not become a statement member")
}

// Quotes and backslashes in a value must survive a document-context
// substitution byte-for-byte once the document is decoded.
func TestSubstituteTemplateParametersDocumentRoundTrip(t *testing.T) {
	value := `a\b "quoted" \u0041 tab	_and_Ümläut`
	defs := []iamstore.RoleTemplateParameterDefinition{{Name: "Service", Type: "String"}}
	values := map[string][]string{"Service": {value}}

	resolved, err := substituteTemplateParameters(trustPolicyPattern, "AssumeRolePolicyDocumentTemplate", defs, values, true)
	require.NoError(t, err)

	var doc map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(resolved), &doc))
	service := doc["Statement"].([]interface{})[0].(map[string]interface{})["Principal"].(map[string]interface{})["Service"].([]interface{})
	assert.Equal(t, []interface{}{value}, service)
}

// Plain-text patterns (name, path, description, tags) substitute verbatim:
// their own field validators reject invalid resolutions, so escaping here
// would corrupt legal values.
func TestSubstituteTemplateParametersTextValuesVerbatim(t *testing.T) {
	value := `R&D "Ops" \ backslash`
	resolved, err := substituteTemplateParameters("Dept: @{Department}", "RoleDescriptionPattern", nil, map[string][]string{"Department": {value}}, false)
	require.NoError(t, err)
	assert.Equal(t, "Dept: "+value, resolved)
}

// The declared parameter type is the constraint a template places on its
// values: Number must be numeric, Arn must parse as an ARN, String (and an
// omitted type) accept anything, and a type outside the documented enum is
// a template defect. The default value is held to the same type as a
// supplied one.
func TestTemplateParameterTypeValidation(t *testing.T) {
	cases := []struct {
		name      string
		typeName  string
		value     string
		wantError bool
	}{
		{"number integer", "Number", "42", false},
		{"number negative decimal", "Number", "-0.5", false},
		{"number exponent", "Number", "1e5", false},
		{"number fraction only", "Number", ".5", false},
		{"number trailing dot", "Number", "5.", false},
		{"number word", "Number", "NaN", true},
		{"number empty", "Number", "", true},
		{"number list form", "NumberList", "12", false},
		{"arn role", "Arn", "arn:aws:iam::123456789012:role/example", false},
		{"arn service principal is not an ARN", "Arn", "iam.amazonaws.com", true},
		{"arn list form", "ArnList", "arn:aws:s3:::bucket", false},
		{"string quoting is legal", "String", `quo"ted`, false},
		{"string list", "StringList", "a,b", false},
		{"omitted type is unconstrained", "", `anything "goes"`, false},
		{"unsupported type", "Bogus", "x", true},
	}
	for _, tc := range cases {
		defs := []iamstore.RoleTemplateParameterDefinition{{Name: "P", Type: tc.typeName}}
		_, err := substituteTemplateParameters("@{P}", "RoleNamePattern", defs, map[string][]string{"P": {tc.value}}, false)
		if tc.wantError {
			assert.Error(t, err, "%s", tc.name)
		} else {
			assert.NoError(t, err, "%s", tc.name)
		}
	}
}

// A default value that violates its own declared type is rejected at
// resolution, exactly as a supplied value would be.
func TestTemplateParameterDefaultHeldToType(t *testing.T) {
	defs := []iamstore.RoleTemplateParameterDefinition{{Name: "P", Type: "Number", DefaultValue: "not-a-number", IsRequired: true}}
	_, err := substituteTemplateParameters("@{P}", "RoleNamePattern", defs, nil, false)
	require.Error(t, err)
	awsErr, ok := err.(*awserrors.AWSError)
	require.True(t, ok)
	assert.Equal(t, "InvalidInput", awsErr.Code)
	assert.Contains(t, awsErr.Message, "numeric")
}

// The unresolved-parameter and empty-list rejections predate the type
// validation and must hold unchanged.
func TestTemplateParameterResolutionErrorsUnchanged(t *testing.T) {
	_, err := substituteTemplateParameters("@{P}", "RoleNamePattern", nil, nil, false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no value supplied for template parameter P")

	_, err = substituteTemplateParameters("@{P}", "RoleNamePattern", nil, map[string][]string{"P": {}}, false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "carries an empty value list")
}

// End to end against the catalogue's Example template: a hostile
// Department value cannot pass the role-name pattern's field validation,
// and a benign value materialises the templated role unchanged.
func TestAcquireRoleCoreNameContextAndBenignPath(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()
	store := iamstore.NewIAMStore(st, "123456789012")
	s := NewIAMService("123456789012")

	const templateArn = "arn:aws:iam::aws:role-template/awsserviceprincipal/Example:1"

	_, err = s.acquireRoleCore(store, &AcquireRoleInput{
		TemplateArn: templateArn,
		ReplacementValues: map[string][]string{
			"Department": {`E"},{"Effect":"Allow","Action":"*","Resource":"*"`},
		},
	})
	require.Error(t, err, "a value that breaks the role-name rules must be rejected")
	awsErr, ok := err.(*awserrors.AWSError)
	require.True(t, ok)
	assert.Equal(t, "InvalidInput", awsErr.Code)

	role, err := s.acquireRoleCore(store, &AcquireRoleInput{
		TemplateArn:       templateArn,
		ReplacementValues: map[string][]string{"Department": {"Engineering"}},
	})
	require.NoError(t, err)
	assert.Equal(t, "Example-Engineering", role.RoleName)
	assert.Equal(t, "/awsserviceprincipal/", role.Path)
	assert.True(t, json.Valid([]byte(role.AssumeRolePolicyDocument)), "the template trust policy must stay well-formed")
}

// A template-layer failure must reverse the layers applied so far and
// delete the role: AcquireRole is all-or-nothing, so a retried request —
// the documented remedy for a mid-creation rejection — never collides
// with a half-provisioned role name.
func TestApplyRoleTemplateLayersRollback(t *testing.T) {
	secondPolicyKey := PrincipalTypeRole + ":RollbackRole:SecondLayer"
	s, store := faultTestStore(t, "iam_inline_policies", &readFault{
		failPutKey: secondPolicyKey,
		err:        errors.New("simulated inline policy write failure"),
	})

	trustPolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":["ec2.amazonaws.com"]},"Action":["sts:AssumeRole"]}]}`
	role, err := store.Roles().Create("RollbackRole", "/", store.AccountID(), trustPolicy, "", iamstore.DefaultRoleSessionDuration, nil, nil)
	require.NoError(t, err)

	inlineDoc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:ListBucket"],"Resource":["*"]}]}`
	template := &iamstore.RoleTemplateVersion{
		InlinePolicyTemplates: []iamstore.RoleTemplateInlinePolicy{
			{PolicyName: "FirstLayer", PolicyDocument: inlineDoc},
			{PolicyName: "SecondLayer", PolicyDocument: inlineDoc},
		},
	}

	err = s.applyRoleTemplateLayers(store, role, template, nil)
	require.Error(t, err, "the second inline policy write must fail")

	assert.False(t, store.Roles().Exists("RollbackRole"),
		"the role must be deleted after a layer failure")
	names, err := store.InlinePolicies().List(PrincipalTypeRole, "RollbackRole")
	require.NoError(t, err)
	assert.Empty(t, names, "the applied inline policy must be reversed")

	// With the fault gone, the same acquire succeeds — no name collision
	// with a half-provisioned role remains.
	cleanStore := quotaTestStore(t)
	cleanRole, err := cleanStore.Roles().Create("RollbackRole", "/", cleanStore.AccountID(), trustPolicy, "", iamstore.DefaultRoleSessionDuration, nil, nil)
	require.NoError(t, err)
	require.NoError(t, s.applyRoleTemplateLayers(cleanStore, cleanRole, template, nil))
	names, err = cleanStore.InlinePolicies().List(PrincipalTypeRole, "RollbackRole")
	require.NoError(t, err)
	assert.Len(t, names, 2)
}

// AcquireRole records the originating template on the created role: the
// template ARN and the minor version actually resolved, so a request that
// selected the default reports the default's version.
func TestAcquireRoleCoreRecordsSourceTemplate(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()
	store := iamstore.NewIAMStore(st, "123456789012")
	s := NewIAMService("123456789012")

	const templateArn = "arn:aws:iam::aws:role-template/awsserviceprincipal/Example:1"

	role, err := s.acquireRoleCore(store, &AcquireRoleInput{
		TemplateArn:       templateArn,
		ReplacementValues: map[string][]string{"Department": {"Engineering"}},
	})
	require.NoError(t, err)
	if role.SourceRoleTemplate == nil {
		t.Fatal("the created role must record its source template")
	}
	assert.Equal(t, templateArn, role.SourceRoleTemplate.TemplateArn)
	assert.Equal(t, 1, role.SourceRoleTemplate.TemplateMinorVersion)

	fetched, err := s.getRoleCore(store, role.RoleName)
	require.NoError(t, err)
	require.NotNil(t, fetched.SourceRoleTemplate, "the persisted role must keep the provenance")
	assert.Equal(t, templateArn, fetched.SourceRoleTemplate.TemplateArn)

	// A plain CreateRole never carries provenance.
	plain, err := s.createRoleCore(store, &CreateRoleInput{
		RoleName:                 "PlainRole",
		AssumeRolePolicyDocument: `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":["ec2.amazonaws.com"]},"Action":["sts:AssumeRole"]}]}`,
	})
	require.NoError(t, err)
	assert.Nil(t, plain.SourceRoleTemplate)
}

// A replacement value list outside the documented 1..20 bound is rejected
// before the template resolves.
func TestAcquireRoleCoreReplacementValueListBounds(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()
	store := iamstore.NewIAMStore(st, "123456789012")
	s := NewIAMService("123456789012")

	const templateArn = "arn:aws:iam::aws:role-template/awsserviceprincipal/Example:1"
	base := &AcquireRoleInput{TemplateArn: templateArn}

	over := *base
	over.ReplacementValues = map[string][]string{"Department": make([]string, iamstore.MaxRoleTemplateReplacementValues+1)}
	_, err = s.acquireRoleCore(store, &over)
	require.Error(t, err, "a 21-value list must be rejected")
	assert.Contains(t, err.Error(), "ReplacementValues")

	empty := *base
	empty.ReplacementValues = map[string][]string{"Department": {}}
	_, err = s.acquireRoleCore(store, &empty)
	require.Error(t, err, "a zero-value list must be rejected")

	atBound := *base
	atBound.ReplacementValues = map[string][]string{"Department": make([]string, iamstore.MaxRoleTemplateReplacementValues)}
	_, err = s.acquireRoleCore(store, &atBound)
	require.NoError(t, err, "a 20-value list is within the bound")
}

// An unterminated placeholder is reported against the pattern member under
// resolution, whatever that pattern is.
func TestUnterminatedPlaceholderNamesItsMember(t *testing.T) {
	_, err := substituteTemplateParameters("@{P", "RoleTagsTemplate", nil, map[string][]string{"P": {"v"}}, false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "RoleTagsTemplate")

	_, err = substituteTemplateParameters("@{P", "AssumeRolePolicyDocumentTemplate", nil, map[string][]string{"P": {"v"}}, true)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "AssumeRolePolicyDocumentTemplate")
}
