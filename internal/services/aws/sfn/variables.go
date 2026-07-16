package sfn

import (
	"encoding/json"
	"fmt"
	"sync"
	"unicode"
)

const (
	maxVariableNameLength = 80
	maxPerVariableBytes   = 256 * 1024
	maxAssignTotalBytes   = 256 * 1024
	maxTotalVariableBytes = 10 * 1024 * 1024
)

// VariableScope manages a hierarchical collection of Step Functions state variables.
// Variables defined in inner scopes shadow those in outer scopes;
// lookups fall through to the parent when not found locally.
type VariableScope struct {
	mu        sync.RWMutex
	parent    *VariableScope
	variables map[string]interface{}
	totalSize int64
}

// NewVariableScope creates a new variable scope, optionally nested under a parent scope.
// The total size is initialised from the parent to enforce aggregate limits.
func NewVariableScope(parent *VariableScope) *VariableScope {
	var total int64
	if parent != nil {
		parent.mu.RLock()
		total = parent.totalSize
		parent.mu.RUnlock()
	}
	return &VariableScope{
		parent:    parent,
		variables: make(map[string]interface{}),
		totalSize: total,
	}
}

// Get retrieves the value of a variable by name, searching parent scopes if not found locally.
func (s *VariableScope) Get(name string) (interface{}, bool) {
	s.mu.RLock()
	if v, ok := s.variables[name]; ok {
		s.mu.RUnlock()
		return v, true
	}
	s.mu.RUnlock()

	if s.parent != nil {
		return s.parent.Get(name)
	}
	return nil, false
}

// GetAll returns a flattened map of all variables visible in this scope, including those from parent scopes.
func (s *VariableScope) GetAll() map[string]interface{} {
	s.mu.RLock()
	result := make(map[string]interface{}, len(s.variables))
	for k, v := range s.variables {
		result[k] = v
	}
	s.mu.RUnlock()

	if s.parent != nil {
		for k, v := range s.parent.GetAll() {
			if _, exists := result[k]; !exists {
				result[k] = v
			}
		}
	}
	return result
}

// SetAll atomically assigns multiple variables, enforcing size limits and preventing shadowing of outer scope variables.
func (s *VariableScope) SetAll(assignments map[string]interface{}) error {
	if len(assignments) == 0 {
		return nil
	}

	assignTotal := int64(0)
	for name, value := range assignments {
		if err := ValidateVariableName(name); err != nil {
			return err
		}
		size := valueSize(value)
		if size > maxPerVariableBytes {
			return fmt.Errorf("variable %q exceeds maximum size of %d bytes", name, maxPerVariableBytes)
		}
		assignTotal += size
	}
	if assignTotal > maxAssignTotalBytes {
		return fmt.Errorf("combined Assign exceeds maximum size of %d bytes", maxAssignTotalBytes)
	}

	if s.parent != nil {
		for name := range assignments {
			if s.parent.isDefinedInScope(name) {
				return fmt.Errorf("cannot shadow variable %q defined in outer scope", name)
			}
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	newTotal := s.totalSize
	for _, value := range assignments {
		newTotal += valueSize(value)
	}
	if newTotal > maxTotalVariableBytes {
		return fmt.Errorf("total variable size would exceed maximum of %d bytes", maxTotalVariableBytes)
	}

	for name, value := range assignments {
		if existing, ok := s.variables[name]; ok {
			s.totalSize -= valueSize(existing)
		}
		s.variables[name] = value
		s.totalSize += valueSize(value)
	}

	return nil
}

func (s *VariableScope) isDefinedInScope(name string) bool {
	s.mu.RLock()
	_, ok := s.variables[name]
	s.mu.RUnlock()
	if ok {
		return true
	}
	if s.parent != nil {
		return s.parent.isDefinedInScope(name)
	}
	return false
}

// NewChild creates a nested child scope that inherits lookups from this scope.
func (s *VariableScope) NewChild() *VariableScope {
	return NewVariableScope(s)
}

// Snapshot returns a deep copy of all variables defined directly in this scope (excluding parent).
func (s *VariableScope) Snapshot() map[string]interface{} {
	s.mu.RLock()
	result := make(map[string]interface{}, len(s.variables))
	for k, v := range s.variables {
		result[k] = deepCopyValue(v)
	}
	s.mu.RUnlock()
	return result
}

// ValidateVariableName checks that a variable name conforms to Step Functions
// naming rules (ID_Start followed by ID_Continue characters, max 80 chars).
func ValidateVariableName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("variable name must not be empty")
	}
	if len(name) > maxVariableNameLength {
		return fmt.Errorf("variable name %q exceeds maximum length of %d characters", name, maxVariableNameLength)
	}
	for i, r := range name {
		if i == 0 {
			if !isIDStart(r) {
				return fmt.Errorf("variable name %q must start with ID_Start character", name)
			}
		} else {
			if !isIDContinue(r) {
				return fmt.Errorf("variable name %q contains invalid character at position %d", name, i)
			}
		}
	}
	return nil
}

func isIDStart(r rune) bool {
	return unicode.IsLetter(r) || r == '_' || r == '$'
}

func isIDContinue(r rune) bool {
	return isIDStart(r) || unicode.IsDigit(r)
}

func valueSize(v interface{}) int64 {
	b, err := json.Marshal(v)
	if err != nil {
		return 0
	}
	return int64(len(b))
}

func deepCopyValue(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{}, len(val))
		for k, v := range val {
			result[k] = deepCopyValue(v)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(val))
		for i, v := range val {
			result[i] = deepCopyValue(v)
		}
		return result
	default:
		return v
	}
}
