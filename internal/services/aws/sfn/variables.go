package sfn

import (
	"encoding/json"
	"fmt"
	"sync"
	"unicode"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// variableBudget is the execution-wide aggregate every scope of one
// execution shares: "The total size of all stored variables cannot exceed
// 10MiB per execution" — Parallel branches and Map iterations cannot see
// each other's values, but the total is per execution, not per branch, so
// the ceiling must be enforced against one shared counter rather than a
// snapshot of the ancestor chain taken at branch start. "Stored" counts
// the variables live scopes hold: "When a Parallel branch or Map
// iteration completes, its variables go out of scope and are no longer
// accessible", so Release returns a completed scope's bytes instead of
// billing them against the ceiling forever.
type variableBudget struct {
	mu    sync.Mutex
	total int64
}

// VariableScope manages a hierarchical collection of Step Functions state variables.
// Variables defined in inner scopes shadow those in outer scopes;
// lookups fall through to the parent when not found locally.
type VariableScope struct {
	mu        sync.RWMutex
	parent    *VariableScope
	variables map[string]interface{}
	budget    *variableBudget
}

// NewVariableScope creates a new variable scope, optionally nested under a parent scope.
// The execution-wide variable budget is shared with the parent — every
// scope of one execution draws from the same aggregate ceiling.
func NewVariableScope(parent *VariableScope) *VariableScope {
	budget := &variableBudget{}
	if parent != nil {
		// The budget pointer is immutable after creation, so the read
		// needs no lock.
		budget = parent.budget
	}
	return &VariableScope{
		parent:    parent,
		variables: make(map[string]interface{}),
		budget:    budget,
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
		if size > sfnstore.MaxPerVariableBytes {
			return fmt.Errorf("variable %q exceeds maximum size of %d bytes", name, sfnstore.MaxPerVariableBytes)
		}
		assignTotal += size
	}
	if assignTotal > sfnstore.MaxAssignTotalBytes {
		return fmt.Errorf("combined Assign exceeds maximum size of %d bytes", sfnstore.MaxAssignTotalBytes)
	}

	if s.parent != nil {
		for name := range assignments {
			if s.parent.isDefinedInScope(name) {
				return fmt.Errorf("cannot shadow variable %q defined in outer scope", name)
			}
		}
	}

	// The aggregate check draws on the shared execution budget: sibling
	// branches and concurrent iterations each hold their own values, but
	// the per-execution total bounds them all together. A variable
	// reassigned in this scope releases its prior value's bytes.
	s.mu.Lock()
	defer s.mu.Unlock()
	s.budget.mu.Lock()
	defer s.budget.mu.Unlock()

	delta := int64(0)
	for name, value := range assignments {
		if existing, ok := s.variables[name]; ok {
			delta -= valueSize(existing)
		}
		delta += valueSize(value)
	}
	if s.budget.total+delta > sfnstore.MaxTotalVariableBytes {
		return fmt.Errorf("total variable size would exceed maximum of %d bytes", sfnstore.MaxTotalVariableBytes)
	}
	s.budget.total += delta

	for name, value := range assignments {
		s.variables[name] = value
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

// Release retires this scope: its local variables go out of scope ("no
// longer accessible") and their bytes return to the execution-wide budget,
// so a Map iteration or Parallel branch that completed stops counting
// against the 10MiB total. The map is cleared, which makes the release
// idempotent and leaves post-release lookups falling through to the parent
// exactly as the out-of-scope wording describes. Root scopes never
// release — their budget dies with the execution itself.
func (s *VariableScope) Release() {
	s.mu.Lock()
	s.budget.mu.Lock()
	for _, v := range s.variables {
		s.budget.total -= valueSize(v)
	}
	s.variables = make(map[string]interface{})
	s.budget.mu.Unlock()
	s.mu.Unlock()
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
	// "the maximum length of a variable name is 80" characters — the count
	// is Unicode characters, so multi-byte CJK names are measured in runes
	// and the reported position is a rune index.
	runes := []rune(name)
	if len(runes) > sfnstore.MaxVariableNameLength {
		return fmt.Errorf("variable name %q exceeds maximum length of %d characters", name, sfnstore.MaxVariableNameLength)
	}
	for i, r := range runes {
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

// isIDStart reports whether a rune may begin a variable name: the UAX #31
// ID_Start derivation — letters, letter numbers and Other_ID_Start minus
// the Pattern_Syntax and Pattern_White_Space exclusions. The '$' reference
// sigil fails every arm (currency symbol, Pattern_Syntax), which is why a
// name never contains it. "The first character of a variable name MUST be
// a Unicode ID_Start character."
func isIDStart(r rune) bool {
	return (unicode.IsLetter(r) || unicode.Is(unicode.Nl, r) || unicode.Is(unicode.Other_ID_Start, r)) &&
		!unicode.Is(unicode.Pattern_Syntax, r) &&
		!unicode.Is(unicode.Pattern_White_Space, r)
}

// isIDContinue reports whether a rune may continue a variable name: ID_Start
// plus the ID_Continue derivation's additions — combining marks, decimal
// digits, connector punctuation and Other_ID_Continue. The underscore lives
// here (connector punctuation), not in ID_Start, so it may continue a name
// but never begin one.
func isIDContinue(r rune) bool {
	return isIDStart(r) || unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Mc, r) ||
		unicode.Is(unicode.Nd, r) || unicode.Is(unicode.Pc, r) || unicode.Is(unicode.Other_ID_Continue, r)
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
