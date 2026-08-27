package iot

import "testing"

// TestValidateDomainConfigEnum pins the enum sets enforced on the
// domain-configuration members.
func TestValidateDomainConfigEnum(t *testing.T) {
	if err := validateDomainConfigEnum("DATA", domainConfigServiceTypes); err != nil {
		t.Fatal("service type DATA rejected")
	}
	if err := validateDomainConfigEnum("BIT_BUCKET_OVER_SMS", domainConfigServiceTypes); err == nil {
		t.Fatal("undocumented service type accepted")
	}
	if err := validateDomainConfigEnum("DISABLED", domainConfigStatuses); err != nil {
		t.Fatal("status DISABLED rejected")
	}
	if err := validateDomainConfigEnum("PAUSED", domainConfigStatuses); err == nil {
		t.Fatal("undocumented status accepted")
	}
	if err := validateDomainConfigEnum("AWS_SIGV4", domainConfigAuthTypes); err != nil {
		t.Fatal("authentication type AWS_SIGV4 rejected")
	}
	if err := validateDomainConfigEnum("SIGV4", domainConfigAuthTypes); err == nil {
		t.Fatal("undocumented authentication type accepted")
	}
	if err := validateDomainConfigEnum("SECURE_MQTT", domainConfigAppProtocols); err != nil {
		t.Fatal("application protocol SECURE_MQTT rejected")
	}
	if err := validateDomainConfigEnum("MQTT", domainConfigAppProtocols); err == nil {
		t.Fatal("undocumented application protocol accepted")
	}
	// Empty values are "not supplied" and must pass every set.
	for _, set := range []map[string]struct{}{domainConfigServiceTypes, domainConfigStatuses, domainConfigAuthTypes, domainConfigAppProtocols} {
		if err := validateDomainConfigEnum("", set); err != nil {
			t.Fatal("empty value rejected")
		}
	}
}
