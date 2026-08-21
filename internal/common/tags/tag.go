package tags

// Tag represents an AWS resource tag (the wire form shared across
// services: {"Key": ..., "Value": ...}).
type Tag struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}
