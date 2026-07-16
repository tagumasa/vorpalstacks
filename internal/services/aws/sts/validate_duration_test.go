package sts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateDurationSeconds(t *testing.T) {
	tests := []struct {
		name          string
		duration      int
		expectedValue int
		expectError   bool
	}{
		{
			name:          "default duration when zero",
			duration:      0,
			expectedValue: 3600,
			expectError:   false,
		},
		{
			name:          "minimum duration",
			duration:      900,
			expectedValue: 900,
			expectError:   false,
		},
		{
			name:          "maximum duration",
			duration:      43200,
			expectedValue: 43200,
			expectError:   false,
		},
		{
			name:          "valid duration in middle",
			duration:      7200,
			expectedValue: 7200,
			expectError:   false,
		},
		{
			name:          "duration below minimum",
			duration:      899,
			expectedValue: 0,
			expectError:   true,
		},
		{
			name:          "duration above maximum",
			duration:      43201,
			expectedValue: 0,
			expectError:   true,
		},
		{
			name:          "duration exactly at minimum minus one",
			duration:      899,
			expectedValue: 0,
			expectError:   true,
		},
		{
			name:          "duration exactly at maximum plus one",
			duration:      43201,
			expectedValue: 0,
			expectError:   true,
		},
		{
			name:          "large valid duration",
			duration:      36000,
			expectedValue: 36000,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validateDurationSeconds(tt.duration)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, ErrInvalidDuration, err)
				assert.Equal(t, 0, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedValue, result)
			}
		})
	}
}

func TestDurationSecondsConstants(t *testing.T) {
	assert.Equal(t, 900, MinDurationSeconds)
	assert.Equal(t, 43200, MaxDurationSeconds)
	assert.Equal(t, 3600, DefaultDurationSeconds)

	assert.True(t, MinDurationSeconds < DefaultDurationSeconds)
	assert.True(t, DefaultDurationSeconds < MaxDurationSeconds)
}
