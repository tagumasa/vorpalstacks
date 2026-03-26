// Package validators provides AWS-specific utility functions for vorpalstacks.
package validators

var validConditionOperators = map[string]bool{
	"StringEquals":              true,
	"StringNotEquals":           true,
	"StringLike":                true,
	"StringNotLike":             true,
	"StringEqualsIgnoreCase":    true,
	"StringNotEqualsIgnoreCase": true,
	"NumericEquals":             true,
	"NumericNotEquals":          true,
	"NumericLessThan":           true,
	"NumericLessThanEquals":     true,
	"NumericGreaterThan":        true,
	"NumericGreaterThanEquals":  true,
	"DateEquals":                true,
	"DateNotEquals":             true,
	"DateLessThan":              true,
	"DateLessThanEquals":        true,
	"DateGreaterThan":           true,
	"DateGreaterThanEquals":     true,
	"Bool":                      true,
	"IpAddress":                 true,
	"NotIpAddress":              true,
	"ArnEquals":                 true,
	"ArnNotEquals":              true,
	"ArnLike":                   true,
	"ArnNotLike":                true,
	"Null":                      true,
	"BinaryEquals":              true,
	"ForAllValues:StringEquals": true,
	"ForAnyValue:StringEquals":  true,
}

// IsValidConditionOperator returns true if the given operator is a valid IAM condition operator.
//
// Valid operators include:
//   - StringEquals, StringNotEquals, StringLike, StringNotLike
//   - StringEqualsIgnoreCase, StringNotEqualsIgnoreCase
//   - NumericEquals, NumericNotEquals, NumericLessThan, NumericLessThanEquals
//   - NumericGreaterThan, NumericGreaterThanEquals
//   - DateEquals, DateNotEquals, DateLessThan, DateLessThanEquals
//   - DateGreaterThan, DateGreaterThanEquals, Bool
//   - IpAddress, NotIpAddress, ArnEquals, ArnNotEquals
//   - ArnLike, ArnNotLike, Null, BinaryEquals
//   - ForAllValues:StringEquals, ForAnyValue:StringEquals
//
// Parameters:
//   - operator: The condition operator to validate
//
// Returns:
//   - bool: True if the operator is valid
func IsValidConditionOperator(operator string) bool {
	return validConditionOperators[operator]
}
