package iot

import (
	"errors"
	"strings"
	"testing"
	"time"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

// Shutdown's worker wait must return immediately when nothing is running,
// honour the bound while a worker is blocked, and succeed again once the
// worker finishes.
func TestWaitRegistrationWorkers(t *testing.T) {
	if !waitRegistrationWorkers(time.Second) {
		t.Fatal("no worker running: the wait must return immediately")
	}
	registrationTaskWg.Add(1)
	start := time.Now()
	if waitRegistrationWorkers(50 * time.Millisecond) {
		t.Fatal("a blocked worker must overrun a short bound")
	}
	if elapsed := time.Since(start); elapsed < 50*time.Millisecond {
		t.Fatalf("the bound was not honoured: returned after %v", elapsed)
	}
	registrationTaskWg.Done()
	if !waitRegistrationWorkers(time.Second) {
		t.Fatal("after Done the wait must succeed")
	}
}

func TestParseProvisioningTemplate(t *testing.T) {
	// The guide's own examples use arbitrary logical names — resources are
	// identified by Type, never by the literal key.
	lowercaseThing := `{"Parameters":{"ThingName":{"Type":"String"}},"Resources":{"thing":{"Type":"AWS::IoT::Thing","Properties":{"ThingName":{"Ref":"{{ThingName}}"}}}}}`
	camelCaseThing := `{"Resources":{"my-device":{"Type":"AWS::IoT::Thing","Properties":{"ThingName":"lit-name"}}}}`
	fullTemplate := `{"Parameters":{"ThingName":{"Type":"String"},"CSR":{"Type":"String"}},
		"Resources":{
			"thing":{"Type":"AWS::IoT::Thing","Properties":{"ThingName":{"Ref":"{{ThingName}}"}}},
			"certificate":{"Type":"AWS::IoT::Certificate","Properties":{"CertificateSigningRequest":{"Ref":"CSR"},"Status":"ACTIVE"},"OverrideSettings":{"Status":"REPLACE"}},
			"policy":{"Type":"AWS::IoT::Policy","Properties":{"PolicyDocument":"{\"Version\":\"2012-10-17\",\"Statement\":[]}"}}
		}}`
	tests := []struct {
		name          string
		body          string
		expectMissing bool
		expectInvalid bool
	}{
		{"lowercase logical name accepted", lowercaseThing, false, false},
		{"arbitrary logical name accepted", camelCaseThing, false, false},
		{"full template with cert and policy", fullTemplate, false, false},
		{"empty body", "", true, false},
		{"malformed json", "{not json", false, true},
		{"missing resources", `{"Parameters":{}}`, false, true},
		{"resources without a thing resource", `{"Resources":{"policy":{"Type":"AWS::IoT::Policy","Properties":{"PolicyName":"p"}}}}`, false, true},
		{"resources not an object", `{"Resources":[]}`, false, true},
		{"unknown resource type", `{"Resources":{"x":{"Type":"AWS::IoT::Other"}}}`, false, true},
		{"certificate with two declarations", `{"Resources":{"thing":{"Type":"AWS::IoT::Thing"},"cert":{"Type":"AWS::IoT::Certificate","Properties":{"CertificateId":"a","CertificatePem":"b"}}}}`, false, true},
		{"policy with both name and document", `{"Resources":{"thing":{"Type":"AWS::IoT::Thing"},"policy":{"Type":"AWS::IoT::Policy","Properties":{"PolicyName":"p","PolicyDocument":"d"}}}}`, false, true},
		{"policy with neither name nor document", `{"Resources":{"thing":{"Type":"AWS::IoT::Thing"},"policy":{"Type":"AWS::IoT::Policy","Properties":{}}}}`, false, true},
		{"merge override on non-mergeable property", `{"Resources":{"thing":{"Type":"AWS::IoT::Thing","OverrideSettings":{"ThingTypeName":"MERGE"}}}}`, false, true},
		{"unknown override action", `{"Resources":{"thing":{"Type":"AWS::IoT::Thing","OverrideSettings":{"ThingGroups":"SHATTER"}}}}`, false, true},
		{"override settings on a policy resource", `{"Resources":{"thing":{"Type":"AWS::IoT::Thing"},"policy":{"Type":"AWS::IoT::Policy","Properties":{"PolicyName":"p"},"OverrideSettings":{"ThingGroups":"REPLACE"}}}}`, false, true},
		{"over length bound", `{"Resources":{"thing":{"Type":"AWS::IoT::Thing"}},"padding":"` + strings.Repeat("x", MaxTemplateBodyLength) + `"}`, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseProvisioningTemplate(tt.body)
			if tt.expectMissing && !errors.Is(err, iotstore.ErrMissingParam) {
				t.Fatalf("expected ErrMissingParam, got %v", err)
			}
			if tt.expectInvalid && !errors.Is(err, iotstore.ErrInvalidRequest) {
				t.Fatalf("expected ErrInvalidRequest, got %v", err)
			}
			if !tt.expectMissing && !tt.expectInvalid && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// Ref resolution must accept both documented forms — {"Ref":"{{ThingName}}"}
// and the bare {"Ref":"SerialNumber"} — and fall back to a parameter's
// declared Default when the caller omits the value.
func TestProvisioningTemplateRefResolution(t *testing.T) {
	tpl, err := parseProvisioningTemplate(`{"Parameters":{"ThingName":{"Type":"String"},"Location":{"Type":"String","Default":"WA"}},"Resources":{"thing":{"Type":"AWS::IoT::Thing","Properties":{"ThingName":{"Ref":"{{ThingName}}"},"AttributePayload":{"loc":{"Ref":"Location"}}}}}}`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	var thingRes *templateResource
	for i := range tpl.resources {
		if tpl.resources[i].Type == thingResourceType {
			thingRes = &tpl.resources[i]
		}
	}
	if thingRes == nil {
		t.Fatal("template carried no thing resource")
	}
	params := map[string]string{"ThingName": "resolved-name"}
	if got := tpl.resolveValue(thingRes.Properties["ThingName"], params); got != "resolved-name" {
		t.Fatalf("expected Ref {{ThingName}} to resolve, got %q", got)
	}
	if got := tpl.resolveValue(thingRes.Properties["AttributePayload"].(map[string]interface{})["loc"], params); got != "WA" {
		t.Fatalf("expected parameter default WA, got %q", got)
	}
	// The override default without an OverrideSettings section is FAIL.
	if got := thingRes.overrideAction("AttributePayload"); got != overrideFail {
		t.Fatalf("expected FAIL default without OverrideSettings, got %q", got)
	}
}

// startThingRegistrationTaskCore must reject the missing of any of the
// four model-required members before store access (the Go SDK validates
// them client-side, so the path is only reachable by direct HTTP).
func TestStartThingRegistrationTaskRequiresAllMembers(t *testing.T) {
	svc := &IoTService{}
	full := StartThingRegistrationTaskInput{
		TemplateBody:    "{}",
		InputFileBucket: "bucket",
		InputFileKey:    "key",
		RoleArn:         "arn:aws:iam::000000000000:role/r",
	}
	missing := []struct {
		name string
		in   StartThingRegistrationTaskInput
	}{
		{"templateBody", StartThingRegistrationTaskInput{InputFileBucket: "b", InputFileKey: "k", RoleArn: "r"}},
		{"inputFileBucket", StartThingRegistrationTaskInput{TemplateBody: "{}", InputFileKey: "k", RoleArn: "r"}},
		{"inputFileKey", StartThingRegistrationTaskInput{TemplateBody: "{}", InputFileBucket: "b", RoleArn: "r"}},
		{"roleArn", StartThingRegistrationTaskInput{TemplateBody: "{}", InputFileBucket: "b", InputFileKey: "k"}},
	}
	for _, tt := range missing {
		t.Run("missing "+tt.name, func(t *testing.T) {
			if _, err := svc.startThingRegistrationTaskCore(nil, tt.in); !errors.Is(err, iotstore.ErrMissingParam) {
				t.Fatalf("expected ErrMissingParam, got %v", err)
			}
		})
	}
	_ = full
}
