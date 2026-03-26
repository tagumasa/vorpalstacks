// Package smithy_reader provides functionality for reading Smithy JSON models directly.
package smithy_reader

// SmithyModel represents a Smithy service model.
type SmithyModel struct {
	Smithy string           `json:"smithy"`
	Shapes map[string]Shape `json:"shapes"`
}

// Shape represents a Smithy shape definition.
type Shape struct {
	Type       string                 `json:"type"`
	Version    string                 `json:"version,omitempty"`
	Operations []OperationRef         `json:"operations,omitempty"`
	Members    map[string]Member      `json:"members,omitempty"`
	Member     *MemberRef             `json:"member,omitempty"`
	Key        *MemberRef             `json:"key,omitempty"`
	Value      *MemberRef             `json:"value,omitempty"`
	Input      *MemberRef             `json:"input,omitempty"`
	Output     *MemberRef             `json:"output,omitempty"`
	Errors     []OperationRef         `json:"errors,omitempty"`
	Enum       []string               `json:"enum,omitempty"`
	Traits     map[string]interface{} `json:"traits,omitempty"`
	Resources  []OperationRef         `json:"resources,omitempty"`
	List       *MemberRef             `json:"list,omitempty"`
	Put        *MemberRef             `json:"put,omitempty"`
	Create     *MemberRef             `json:"create,omitempty"`
	Read       *MemberRef             `json:"read,omitempty"`
	Update     *MemberRef             `json:"update,omitempty"`
	Delete     *MemberRef             `json:"delete,omitempty"`
}

// Member represents a Smithy structure member.
type Member struct {
	Target string                 `json:"target"`
	Traits map[string]interface{} `json:"traits,omitempty"`
}

// MemberRef represents a reference to a Smithy shape member.
type MemberRef struct {
	Target string `json:"target"`
}

// OperationRef represents a reference to a Smithy operation.
type OperationRef struct {
	Target string `json:"target"`
}
