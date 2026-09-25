// Package pbutil carries small helpers for the generated protobuf types
// the admin plane converts through.
package pbutil

// Enum returns a pointer to v, the construction form a generated enum
// member with explicit presence takes (a non-required enum member is
// emitted optional, so its field is a pointer type while the converters
// hold the value form). The int32-kind constraint matches every generated
// proto enum.
func Enum[T ~int32](v T) *T {
	return &v
}
