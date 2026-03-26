package naming

// Package naming provides naming utility functions for vorpalstacks.

import "strings"

// ParseShapeID parses a shape ID into namespace and name parts
// Example: "com.amazonaws.lambda#DeleteFunction" -> ["com.amazonaws.lambda", "DeleteFunction"]
func ParseShapeID(shapeID string) (namespace, name string) {
	parts := strings.Split(shapeID, "#")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "", shapeID
}

// ExtractShapeIDName extracts the shape name from a shape ID
// Example: "com.amazonaws.lambda#DeleteFunction" -> "DeleteFunction"
func ExtractShapeIDName(shapeID string) string {
	_, name := ParseShapeID(shapeID)
	return name
}
