package sns

import "testing"

func TestValidateDataProtectionPolicy(t *testing.T) {
	valid := []string{
		`{"Version":"2012-10-17","Statement":[]}`,
		`{}`,
	}
	for _, policy := range valid {
		if err := validateDataProtectionPolicy(policy); err != nil {
			t.Errorf("validateDataProtectionPolicy(%q) = %v, want nil", policy, err)
		}
	}

	invalid := []string{
		`not json`,
		`{"Statement":`,
		// Exceeds the 30,720-byte cap (31 one-character keys are enough
		// once wrapped in JSON syntax overhead).
		`{"pad":"` + string(make([]byte, maxTopicAttributeValueLength)) + `"}`,
	}
	for _, policy := range invalid {
		if err := validateDataProtectionPolicy(policy); err == nil {
			t.Errorf("validateDataProtectionPolicy(%q) = nil, want error", policy)
		}
	}
}
