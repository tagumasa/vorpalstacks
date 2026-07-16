package validators

import "testing"

func TestCheckCondition(t *testing.T) {
	tests := []struct {
		name          string
		condition     bool
		errorMessage  string
		shouldHaveErr bool
	}{
		{"True condition", true, "error", false},
		{"False condition", false, "test error", true},
		{"True with empty message", true, "", false},
		{"False with empty message", false, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckCondition(tt.condition, tt.errorMessage)
			if tt.shouldHaveErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.shouldHaveErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}
