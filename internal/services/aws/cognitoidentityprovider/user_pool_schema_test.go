package cognitoidentityprovider

import (
	"slices"
	"testing"

	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// TestFormatUserPoolEmitsSchemaAttributesWireName pins the wire contract of
// the user pool schema: the response member is SchemaAttributes (the model
// member name), the full standard attribute set is always present, and
// custom attribute names carry the custom: prefix on output.
func TestFormatUserPoolEmitsSchemaAttributesWireName(t *testing.T) {
	pool := cognitostore.NewUserPool("wire-name-pool", "us-east-1")
	pool.SchemaAttributes = []cognitostore.SchemaAttributeType{
		{Name: "rank", AttributeDataType: "String", Mutable: true},
	}

	formatted := formatUserPool(pool)
	if _, ok := formatted["Schema"]; ok {
		t.Error("emitted the non-model Schema member; the wire member is SchemaAttributes")
	}
	attrs, ok := formatted["SchemaAttributes"].([]map[string]interface{})
	if !ok {
		t.Fatalf("SchemaAttributes member missing from user pool response")
	}

	var names []string
	for _, a := range attrs {
		name, _ := a["Name"].(string)
		names = append(names, name)
	}
	for _, standard := range []string{"sub", "email", "phone_number_verified", "updated_at", "identities"} {
		if !slices.Contains(names, standard) {
			t.Errorf("standard attribute %q missing from SchemaAttributes", standard)
		}
	}
	if !slices.Contains(names, "custom:rank") {
		t.Error("custom attribute name not returned with the custom: prefix")
	}
	if slices.Contains(names, "rank") {
		t.Error("custom attribute returned without the custom: prefix")
	}
}

// TestSchemaAttributesForDescribeStandardSet pins the describe projection:
// the documented standard attribute order with per-attribute defaults, a
// supplied standard attribute definition merged into its slot, and the
// pool's custom attributes appended with their wire names.
func TestSchemaAttributesForDescribeStandardSet(t *testing.T) {
	pool := cognitostore.NewUserPool("describe-projection", "us-east-1")
	pool.SchemaAttributes = []cognitostore.SchemaAttributeType{
		{
			Name:                       "email",
			AttributeDataType:          "String",
			Required:                   true,
			Mutable:                    true,
			StringAttributeConstraints: &cognitostore.StringAttributeConstraints{MinLength: "0", MaxLength: "2048"},
		},
		{Name: "rank", AttributeDataType: "String", Mutable: true},
		{Name: "secret_level", AttributeDataType: "Number", DeveloperOnlyAttribute: true},
	}

	attrs := schemaAttributesForDescribe(pool)
	if len(attrs) != len(standardSchemaAttributeOrder)+2 {
		t.Fatalf("expected %d standard + 2 custom entries, got %d", len(standardSchemaAttributeOrder), len(attrs))
	}
	for i, name := range standardSchemaAttributeOrder {
		if attrs[i].Name != name {
			t.Fatalf("standard entry %d is %q, want %q", i, attrs[i].Name, name)
		}
	}
	if attrs[0].Name != "sub" || !attrs[0].Required || attrs[0].Mutable {
		t.Errorf("sub must be required and immutable: %+v", attrs[0])
	}
	birthdate := attrs[slices.Index(standardSchemaAttributeOrder, "birthdate")]
	if birthdate.StringAttributeConstraints == nil ||
		birthdate.StringAttributeConstraints.MinLength != "10" || birthdate.StringAttributeConstraints.MaxLength != "10" {
		t.Errorf("birthdate constraints must be exactly 10-10: %+v", birthdate.StringAttributeConstraints)
	}
	updatedAt := attrs[slices.Index(standardSchemaAttributeOrder, "updated_at")]
	if updatedAt.AttributeDataType != "Number" ||
		updatedAt.NumberAttributeConstraints == nil || updatedAt.NumberAttributeConstraints.MinValue != "0" {
		t.Errorf("updated_at must be a number with a zero minimum: %+v", updatedAt)
	}
	identities := attrs[slices.Index(standardSchemaAttributeOrder, "identities")]
	if identities.StringAttributeConstraints == nil ||
		identities.StringAttributeConstraints.MinLength != "" || identities.StringAttributeConstraints.MaxLength != "" {
		t.Errorf("identities must carry empty string constraints: %+v", identities.StringAttributeConstraints)
	}
	email := attrs[slices.Index(standardSchemaAttributeOrder, "email")]
	if !email.Required {
		t.Errorf("supplied email definition must keep Required=true: %+v", email)
	}
	if last := attrs[len(attrs)-2]; last.Name != "custom:rank" {
		t.Errorf("custom attribute not projected with the custom: prefix: %+v", last)
	}
	if last := attrs[len(attrs)-1]; last.Name != "dev:secret_level" {
		t.Errorf("developer-only attribute not projected with the dev: prefix: %+v", last)
	}
}

// TestCSVHeaderForPoolStandardBaseAndPrefixedCustoms pins the CSV header:
// the documented base columns (including cognito:mfa_enabled) followed by
// the pool's custom attribute columns; standard attribute definitions must
// not duplicate a base column.
func TestCSVHeaderForPoolStandardBaseAndPrefixedCustoms(t *testing.T) {
	pool := cognitostore.NewUserPool("csv-header", "us-east-1")
	pool.SchemaAttributes = []cognitostore.SchemaAttributeType{
		{Name: "email", AttributeDataType: "String", Required: true, Mutable: true},
		{Name: "rank", AttributeDataType: "String", Mutable: true},
	}

	header := csvHeaderForPool(pool)
	if len(header) != len(csvHeaderBase)+1 {
		t.Fatalf("expected %d base columns + 1 custom column, got %d (%v)", len(csvHeaderBase), len(header), header)
	}
	if header[0] != "name" || !slices.Contains(header, "cognito:mfa_enabled") || !slices.Contains(header, "cognito:username") {
		t.Errorf("base header must start with the standard names and include cognito:mfa_enabled and cognito:username: %v", header)
	}
	if last := header[len(header)-1]; last != "custom:rank" {
		t.Errorf("custom attribute column missing or not appended after the base header: %v", header)
	}
	emailCount := 0
	for _, column := range header {
		if column == "email" {
			emailCount++
		}
	}
	if emailCount != 1 {
		t.Errorf("standard attribute column duplicated: email appears %d times", emailCount)
	}
	if !slices.Contains(header, "custom:rank") {
		t.Errorf("custom attribute column missing or unprefixed: %v", header)
	}
}

// TestParseSchemaAttributesFillsStandardDefaults pins that a standard
// attribute definition supplied through the Schema parameter inherits the
// documented default properties (mutability, string constraints) for the
// members the request left unset, while custom attributes inherit nothing.
func TestParseSchemaAttributesFillsStandardDefaults(t *testing.T) {
	req := &request.ParsedRequest{
		Parameters: map[string]interface{}{
			"Schema": []interface{}{
				map[string]interface{}{"Name": "email", "AttributeDataType": "String", "Required": true},
				map[string]interface{}{"Name": "rank", "AttributeDataType": "String"},
			},
		},
	}

	attrs := parseSchemaAttributes(req)
	if len(attrs) != 2 {
		t.Fatalf("expected 2 parsed schema attributes, got %d", len(attrs))
	}
	email := attrs[0]
	if !email.Required || !email.Mutable {
		t.Errorf("standard attribute did not inherit the documented defaults: %+v", email)
	}
	if email.StringAttributeConstraints == nil || email.StringAttributeConstraints.MaxLength != "2048" {
		t.Errorf("standard attribute missing default string constraints: %+v", email.StringAttributeConstraints)
	}
	if attrs[1].Mutable || attrs[1].StringAttributeConstraints != nil {
		t.Errorf("custom attribute must not inherit standard defaults: %+v", attrs[1])
	}
}

// TestApplyUserPoolUpdatesIgnoresSchema pins that UpdateUserPool — whose
// model request has no Schema member — never mutates the pool schema:
// custom attributes added at creation or through AddCustomAttributes
// survive an update request that carries a stray Schema parameter, while
// the model members the operation does carry still apply.
func TestApplyUserPoolUpdatesIgnoresSchema(t *testing.T) {
	pool := cognitostore.NewUserPool("schema-update", "us-east-1")
	pool.SchemaAttributes = []cognitostore.SchemaAttributeType{
		{Name: "rank", AttributeDataType: "String", Mutable: true},
	}
	req := &request.ParsedRequest{
		Operation: "UpdateUserPool",
		Parameters: map[string]interface{}{
			"UserPoolId": "us-east-1_schema-update",
			"PoolName":   "schema-update-renamed",
			"Schema": []interface{}{
				map[string]interface{}{"Name": "intruder", "AttributeDataType": "String"},
			},
		},
	}

	if err := applyUserPoolUpdates(pool, req); err != nil {
		t.Fatalf("applyUserPoolUpdates failed: %v", err)
	}
	if len(pool.SchemaAttributes) != 1 || pool.SchemaAttributes[0].Name != "rank" {
		t.Errorf("update request replaced the pool schema: %+v", pool.SchemaAttributes)
	}
	if pool.Name != "schema-update-renamed" {
		t.Error("the update did not apply the model members it does carry")
	}
}
