package cognitoidentityprovider

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// The user pool schema projection implements the Amazon Cognito attribute
// contract: the schema contains standard attributes, custom attributes with
// a custom: prefix, and developer-only attributes with a dev: prefix.
// Attribute definitions are supplied with bare names (a Name value of
// MyAttribute creates custom:MyAttribute, or dev:MyAttribute when
// DeveloperOnlyAttribute is true) and describe-style responses always
// return the full standard attribute set followed by the pool's custom
// attributes with their prefix applied.
const (
	customAttributePrefix    = "custom:"
	developerAttributePrefix = "dev:"
)

// standardSchemaAttributeOrder is the order in which Amazon Cognito
// returns the standard attribute set in DescribeUserPool responses.
var standardSchemaAttributeOrder = []string{
	"sub", "name", "given_name", "family_name", "middle_name",
	"nickname", "preferred_username", "profile", "picture", "website",
	"email", "email_verified", "gender", "birthdate", "zoneinfo",
	"locale", "phone_number", "phone_number_verified", "address",
	"updated_at", "identities",
}

// standardSchemaAttributeDefault returns the default properties Amazon
// Cognito reports for a standard attribute that the user pool did not
// customise: sub is required and immutable, birthdate values are exactly
// 10 characters, updated_at is a number with a zero minimum, identities
// carries no string constraints, and every other standard string
// attribute accepts up to 2048 characters.
func standardSchemaAttributeDefault(name string) (cognitostore.SchemaAttributeType, bool) {
	switch name {
	case "sub":
		return cognitostore.SchemaAttributeType{
			AttributeDataType:          "String",
			Mutable:                    false,
			Required:                   true,
			StringAttributeConstraints: &cognitostore.StringAttributeConstraints{MinLength: "1", MaxLength: "2048"},
		}, true
	case "email_verified", "phone_number_verified":
		return cognitostore.SchemaAttributeType{
			AttributeDataType: "Boolean",
			Mutable:           true,
		}, true
	case "birthdate":
		return cognitostore.SchemaAttributeType{
			AttributeDataType:          "String",
			Mutable:                    true,
			StringAttributeConstraints: &cognitostore.StringAttributeConstraints{MinLength: "10", MaxLength: "10"},
		}, true
	case "updated_at":
		return cognitostore.SchemaAttributeType{
			AttributeDataType:          "Number",
			Mutable:                    true,
			NumberAttributeConstraints: &cognitostore.NumberAttributeConstraints{MinValue: "0"},
		}, true
	case "identities":
		return cognitostore.SchemaAttributeType{
			AttributeDataType:          "String",
			Mutable:                    true,
			StringAttributeConstraints: &cognitostore.StringAttributeConstraints{},
		}, true
	case "name", "given_name", "family_name", "middle_name", "nickname",
		"preferred_username", "profile", "picture", "website", "email",
		"gender", "zoneinfo", "locale", "phone_number", "address":
		return cognitostore.SchemaAttributeType{
			AttributeDataType:          "String",
			Mutable:                    true,
			StringAttributeConstraints: &cognitostore.StringAttributeConstraints{MinLength: "0", MaxLength: "2048"},
		}, true
	}
	return cognitostore.SchemaAttributeType{}, false
}

// schemaAttributeWireName returns the attribute name as it appears on the
// wire: standard attribute names stay bare, everything else carries the
// custom: or dev: prefix. Amazon Cognito does not strip a prefix from the
// supplied name, so an attribute defined as "custom:222" with
// DeveloperOnlyAttribute becomes "dev:custom:222".
func schemaAttributeWireName(sa cognitostore.SchemaAttributeType) string {
	if standardSchemaAttributeNames[sa.Name] {
		return sa.Name
	}
	if sa.DeveloperOnlyAttribute {
		return developerAttributePrefix + sa.Name
	}
	return customAttributePrefix + sa.Name
}

// schemaAttributesForDescribe builds the SchemaAttributes list of a
// describe-style response: the standard attribute set in the documented
// order (a definition supplied at pool creation overrides the defaults for
// its standard attribute) followed by the pool's custom attributes with
// their wire names.
func schemaAttributesForDescribe(pool *cognitostore.UserPool) []cognitostore.SchemaAttributeType {
	supplied := make(map[string]cognitostore.SchemaAttributeType, len(pool.SchemaAttributes))
	for _, sa := range pool.SchemaAttributes {
		if sa.Name != "" {
			supplied[sa.Name] = sa
		}
	}

	result := make([]cognitostore.SchemaAttributeType, 0, len(standardSchemaAttributeOrder)+len(pool.SchemaAttributes))
	for _, name := range standardSchemaAttributeOrder {
		if sa, ok := supplied[name]; ok {
			entry := sa
			entry.Name = name
			result = append(result, entry)
			continue
		}
		def, _ := standardSchemaAttributeDefault(name)
		def.Name = name
		result = append(result, def)
	}
	for _, sa := range pool.SchemaAttributes {
		if sa.Name == "" || standardSchemaAttributeNames[sa.Name] {
			continue
		}
		entry := sa
		entry.Name = schemaAttributeWireName(entry)
		result = append(result, entry)
	}
	return result
}

// validateUserAttributesAgainstSchema is the single Core validation every
// user-attribute write passes: each supplied name must exist in the pool's
// schema (a custom: or dev: attribute must be defined by the pool), each
// value must satisfy its attribute's data type and constraints, update
// paths (isCreate=false) reject immutable attributes, and creation paths
// require a non-empty value for every attribute the pool marks required.
// sub is assigned by the service, never by a client: creation paths must
// strip it before calling, and on update paths it fails as immutable.
func validateUserAttributesAgainstSchema(pool *cognitostore.UserPool, attrs map[string]string, isCreate bool) error {
	ordered := schemaAttributesForDescribe(pool)
	schema := make(map[string]cognitostore.SchemaAttributeType, len(ordered))
	for _, sa := range ordered {
		schema[sa.Name] = sa
	}
	for name, value := range attrs {
		sa, ok := schema[name]
		if !ok {
			return fmt.Errorf("attribute %s does not exist in the pool schema", name)
		}
		if !isCreate && !sa.Mutable {
			return fmt.Errorf("attribute %s is immutable", name)
		}
		if err := validateAttributeValueAgainstType(sa, value); err != nil {
			return err
		}
	}
	if !isCreate {
		return nil
	}
	for _, sa := range ordered {
		if !sa.Required || sa.Name == "sub" {
			continue
		}
		if attrs[sa.Name] == "" {
			return fmt.Errorf("required attribute %s must have a value", sa.Name)
		}
	}
	return nil
}

// validateAttributeValueAgainstType checks one attribute value against its
// schema definition's data type and constraints.
func validateAttributeValueAgainstType(sa cognitostore.SchemaAttributeType, value string) error {
	switch sa.AttributeDataType {
	case "Boolean":
		if lower := strings.ToLower(value); lower != "true" && lower != "false" {
			return fmt.Errorf("attribute %s must be true or false", sa.Name)
		}
	case "Number":
		num, numErr := strconv.ParseFloat(value, 64)
		if numErr != nil {
			return fmt.Errorf("attribute %s must be a number", sa.Name)
		}
		if c := sa.NumberAttributeConstraints; c != nil {
			if c.MinValue != "" {
				if min, minErr := strconv.ParseFloat(c.MinValue, 64); minErr == nil && num < min {
					return fmt.Errorf("attribute %s must be at least %s", sa.Name, c.MinValue)
				}
			}
			if c.MaxValue != "" {
				if max, maxErr := strconv.ParseFloat(c.MaxValue, 64); maxErr == nil && num > max {
					return fmt.Errorf("attribute %s must be at most %s", sa.Name, c.MaxValue)
				}
			}
		}
	case "String":
		length := utf8.RuneCountInString(value)
		if c := sa.StringAttributeConstraints; c != nil {
			if c.MinLength != "" {
				if min, minErr := strconv.Atoi(c.MinLength); minErr == nil && length < min {
					return fmt.Errorf("attribute %s must be at least %s characters long", sa.Name, c.MinLength)
				}
			}
			if c.MaxLength != "" {
				if max, maxErr := strconv.Atoi(c.MaxLength); maxErr == nil && length > max {
					return fmt.Errorf("attribute %s must be at most %s characters long", sa.Name, c.MaxLength)
				}
			}
		}
	}
	return nil
}

// csvHeaderBase is the header row Amazon Cognito returns for a user pool
// without custom attributes: the standard attributes that an import file
// can populate, the MFA flag, and the username column.
var csvHeaderBase = []string{
	"name", "given_name", "family_name", "middle_name",
	"nickname", "preferred_username", "profile", "picture", "website",
	"email", "email_verified", "gender", "birthdate", "zoneinfo",
	"locale", "phone_number", "phone_number_verified", "address",
	"updated_at", "cognito:mfa_enabled", "cognito:username",
}

// csvHeaderForPool extends the standard CSV header with the pool's custom
// attribute columns. Standard attribute definitions never add a column of
// their own because every standard attribute is already present in the
// base header.
func csvHeaderForPool(pool *cognitostore.UserPool) []string {
	header := make([]string, len(csvHeaderBase))
	copy(header, csvHeaderBase)
	for _, sa := range pool.SchemaAttributes {
		if sa.Name == "" || standardSchemaAttributeNames[sa.Name] {
			continue
		}
		header = append(header, schemaAttributeWireName(sa))
	}
	return header
}
