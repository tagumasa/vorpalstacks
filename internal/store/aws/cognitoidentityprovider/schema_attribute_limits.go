package cognitoidentityprovider

// Custom-attribute and schema-attribute bounds. Source: the Smithy model's
// CustomAttributesListType length trait (min 1, max 25) and
// SchemaAttributesListType length trait (min 1, max 50). These are the
// single definitions; every other site must reference them.

const (
	// MinCustomAttributesPerAdd is the minimum number of custom attribute
	// definitions one AddCustomAttributes request must carry.
	MinCustomAttributesPerAdd = 1
	// MaxCustomAttributesPerAdd is the maximum number of custom attribute
	// definitions one AddCustomAttributes request may carry.
	MaxCustomAttributesPerAdd = 25
	// MinSchemaAttributesPerPool is the minimum number of attribute
	// definitions a present CreateUserPool Schema member must carry.
	MinSchemaAttributesPerPool = 1
	// MaxSchemaAttributesPerPool is the maximum number of attribute
	// definitions a user pool's schema may hold.
	MaxSchemaAttributesPerPool = 50
)
