package dynamodb

import (
	"regexp"
	"strings"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// ---------------------------------------------------------------------------
// Smithy / AWS Docs constraint constants
// ---------------------------------------------------------------------------

// TableName: 3-255 characters, [a-zA-Z0-9_.-]+
const (
	tableNameMinLength = 3
	tableNameMaxLength = 255
)

// ProjectionType enum values (Smithy ProjectionType).
const (
	ProjectionTypeAll      = "ALL"
	ProjectionTypeKeysOnly = "KEYS_ONLY"
	ProjectionTypeInclude  = "INCLUDE"
)

var tableNameValidationRegex = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)

// validProjectionTypes is the set of allowed ProjectionType values.
var validProjectionTypes = map[string]bool{
	ProjectionTypeAll:      true,
	ProjectionTypeKeysOnly: true,
	ProjectionTypeInclude:  true,
}

// ---------------------------------------------------------------------------
// Reusable validators (centralised — used by CreateTable, UpdateTable, and admin handler)
// ---------------------------------------------------------------------------

// validateTableName checks the DynamoDB table name constraints:
// 3-255 characters, characters [a-zA-Z0-9_.-].
func validateTableName(name string) error {
	if name == "" {
		return ErrInvalidParameter
	}
	if len(name) < tableNameMinLength || len(name) > tableNameMaxLength {
		return ErrInvalidParameter
	}
	if !tableNameValidationRegex.MatchString(name) {
		return ErrInvalidParameter
	}
	return nil
}

// validateProjectionType checks the Smithy ProjectionType enum.
func validateProjectionType(pt string) error {
	if pt == "" {
		return nil
	}
	if !validProjectionTypes[pt] {
		return ErrInvalidParameter
	}
	return nil
}

// validateProjectionRequired ensures a non-nil Projection is present
// when parsing GSI or LSI definitions (Smithy: Projection is required).
func validateProjectionRequired(projMap map[string]interface{}) error {
	if projMap == nil {
		return ErrInvalidParameter
	}
	if _, ok := projMap["ProjectionType"]; !ok {
		return ErrInvalidParameter
	}
	return nil
}

// validateKeyAttributeValue ensures key attribute values use only the
// types allowed for DynamoDB keys: S, N, or B (Smithy KeySchemaAttribute).
func validateKeyAttributeValue(key map[string]*dbstore.AttributeValue) error {
	for _, av := range key {
		if av == nil {
			return ErrInvalidParameter
		}
		// Key attributes must be one of S, N, or B.
		// BOOL, NULL, SS, NS, BS, M, L are not permitted.
		if av.S == nil && av.N == nil && av.B == nil {
			return ErrInvalidParameter
		}
		// Reject empty values (H6: extend to N type).
		if av.S != nil && *av.S == "" {
			return ErrInvalidParameter
		}
		if av.N != nil && *av.N == "" {
			return ErrInvalidParameter
		}
		if av.B != nil && len(av.B) == 0 {
			return ErrInvalidParameter
		}
	}
	return nil
}

// validateGSIDeleteExists checks that a GSI to be deleted actually exists
// before calling delete() on the map (Smithy: returns
// GlobalSecondaryIndexNotFoundException on unknown GSI).
func validateGSIDeleteExists(gsiMap map[string]*dbstore.GlobalSecondaryIndex, indexName string) error {
	if _, exists := gsiMap[indexName]; !exists {
		return ErrIndexNotFound
	}
	return nil
}

// validateGSICreateRequired checks that a GSI Create request includes
// the required IndexName and KeySchema (Smithy required traits).
func validateGSICreateRequired(create map[string]interface{}) error {
	indexName := request.GetStringParam(create, "IndexName")
	if indexName == "" {
		return ErrInvalidParameter
	}
	keySchema := parseKeySchema(create)
	if len(keySchema) == 0 {
		return ErrInvalidParameter
	}
	if err := validateKeySchema(keySchema); err != nil {
		return err
	}
	return nil
}

// validateBillingModeConsistency ensures that PROVISIONED billing mode
// includes ProvisionedThroughput, and PAY_PER_REQUEST does not.
func validateBillingModeConsistency(billingMode dbstore.BillingMode, provThroughput *dbstore.ProvisionedThroughput) error {
	if billingMode == dbstore.BillingModeProvisioned {
		if provThroughput == nil {
			return ErrInvalidParameter
		}
	}
	return nil
}

// validateBoolParam extracts a bool from a parameter map. If the key is
// absent, the default value is returned. If the key is present but not a
// bool, an error is returned (rejects malformed requests that would
// otherwise be silently coerced to the default).
func validateBoolParam(params map[string]interface{}, key string, defaultVal bool) (bool, error) {
	v, ok := params[key]
	if !ok {
		return defaultVal, nil
	}
	b, ok := v.(bool)
	if !ok {
		return false, ErrInvalidParameter
	}
	return b, nil
}

// validateBracketIndex parses a bracket-enclosed list index (e.g. "[3]")
// and returns the integer index. Returns an error for empty, non-numeric,
// or negative values (Smithy document path spec).
func validateBracketIndex(idxStr string) (int, error) {
	idxStr = strings.TrimSpace(idxStr)
	if idxStr == "" {
		return 0, ErrInvalidParameter
	}
	var idx int
	for _, ch := range idxStr {
		if ch < '0' || ch > '9' {
			return 0, ErrInvalidParameter
		}
		idx = idx*10 + int(ch-'0')
	}
	if idx < 0 {
		return 0, ErrInvalidParameter
	}
	return idx, nil
}
