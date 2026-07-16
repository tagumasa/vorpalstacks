package types

// Tag represents an AWS resource tag.
type Tag struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}
