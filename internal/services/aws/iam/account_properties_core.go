package iam

import (
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	iamstore "vorpalstacks/internal/store/aws/iam"
)

// accountPropertyKeyPattern is the Smithy accountPropertyKeyType pattern:
// keys start with a letter and continue with letters, digits, slashes,
// underscores and hyphens.
var accountPropertyKeyPattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9/_-]*$`)

// accountPropertyValueTypes registers every known account property with the
// value type its key's value is validated against — the service-side
// property registry the API documentation describes ("The service validates
// each value based on the property key's expected type"). RoleManager/
// Enabled is the one documented property; its value is a boolean.
var accountPropertyValueTypes = map[string]string{
	"RoleManager/Enabled": "boolean",
}

// getAccountPropertiesCore retrieves the account property map (Namespace/
// PropertyName keys), overlaid on the documented defaults.
func (s *IAMService) getAccountPropertiesCore(store *iamstore.IAMStore) (map[string]string, error) {
	return store.AccountProperties().Get()
}

// putAccountPropertiesCore validates and sets account properties. Every key
// must use the documented Namespace/PropertyName form (exactly one slash,
// neither prefix nor suffix slash, letter-first, within the length bound),
// every value must be within its length bound and parse as the key's
// registered type, and all keys in one request must share a namespace.
func (s *IAMService) putAccountPropertiesCore(store *iamstore.IAMStore, properties map[string]string) error {
	if len(properties) == 0 {
		return NewValidationError("Properties")
	}

	namespace := ""
	for key, value := range properties {
		if utf8.RuneCountInString(key) < 1 || utf8.RuneCountInString(key) > iamstore.MaxAccountPropertyKeyLength {
			return NewInvalidInputError("Properties", "property keys must be 1 to "+strconv.Itoa(iamstore.MaxAccountPropertyKeyLength)+" characters")
		}
		if !accountPropertyKeyPattern.MatchString(key) {
			return NewInvalidInputError("Properties", "property key "+key+" must match ^[A-Za-z][A-Za-z0-9/_-]*$")
		}
		if strings.Count(key, "/") != 1 || strings.HasPrefix(key, "/") || strings.HasSuffix(key, "/") {
			return NewInvalidInputError("Properties", "property key "+key+" must contain exactly one / separating the namespace from the property name, and cannot start or end with /")
		}
		if utf8.RuneCountInString(value) < 1 || utf8.RuneCountInString(value) > iamstore.MaxAccountPropertyValueLength {
			return NewInvalidInputError("Properties", "property values must be 1 to "+strconv.Itoa(iamstore.MaxAccountPropertyValueLength)+" characters")
		}

		keyNamespace := key[:strings.Index(key, "/")]
		if namespace == "" {
			namespace = keyNamespace
		} else if keyNamespace != namespace {
			return NewInvalidInputError("Properties", "all properties in a single request must belong to the same namespace")
		}

		valueType, known := accountPropertyValueTypes[key]
		if !known {
			return NewInvalidInputError("Properties", "unrecognized property key "+key)
		}
		if valueType == "boolean" && value != "true" && value != "false" {
			return NewInvalidInputError("Properties", "boolean property "+key+" expects true or false")
		}
	}

	return store.AccountProperties().Put(properties)
}
