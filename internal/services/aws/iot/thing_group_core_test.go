package iot

import (
	"errors"
	"testing"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

// The update Cores must reject requests that omit the properties member
// before any store access; the store argument is therefore never touched
// on these paths and is passed as nil.
func TestUpdateGroupCoresRejectMissingProperties(t *testing.T) {
	svc := &IoTService{}

	tests := []struct {
		name string
		call func() error
	}{
		{"thing group without name", func() error {
			_, err := svc.updateThingGroupCore(nil, UpdateThingGroupInput{PropertiesProvided: true})
			return err
		}},
		{"thing group without properties", func() error {
			_, err := svc.updateThingGroupCore(nil, UpdateThingGroupInput{GroupName: "group"})
			return err
		}},
		{"billing group without name", func() error {
			_, err := svc.updateBillingGroupCore(nil, UpdateBillingGroupInput{PropertiesProvided: true})
			return err
		}},
		{"billing group without properties", func() error {
			_, err := svc.updateBillingGroupCore(nil, UpdateBillingGroupInput{GroupName: "group"})
			return err
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			if !errors.Is(err, iotstore.ErrMissingParam) {
				t.Fatalf("expected ErrMissingParam, got %v", err)
			}
		})
	}
}
