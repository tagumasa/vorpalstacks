package iot

import (
	"errors"
	"testing"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/storage"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// newTagValidationTestStore opens a throwaway IoT store seeded with one
// resource of a representative set of taggable kinds: a typed-store thing
// and the generic-KV families (custom metric, dimension, CA certificate).
func newTagValidationTestStore(t *testing.T) iotstore.IotStoreInterface {
	t.Helper()
	raw, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { raw.Close() })
	store := iotstore.NewIotStore(raw, "000000000000", "us-east-1", nil)

	if _, err := store.CreateThing(&iotstore.Thing{ThingName: "pin-thing"}); err != nil {
		t.Fatalf("seed thing: %v", err)
	}
	for key, value := range map[string]interface{}{
		"customMetric/cm-pin": map[string]interface{}{"name": "cm-pin"},
		"dimension/dim-pin":   map[string]interface{}{"name": "dim-pin"},
		"caCert/cacert-pin":   map[string]interface{}{"name": "cacert-pin"},
	} {
		if err := store.PutGeneric(key, value); err != nil {
			t.Fatalf("seed %s: %v", key, err)
		}
	}
	return store
}

// TestValidateTagResourceCore pins the tag-target resolution behind the
// IoT tag trio: a seeded thing and the generic-KV families resolve, a
// missing resource of a hosted kind is the modelled
// ResourceNotFoundException, and an ARN of any other kind addresses no
// resource this service can own.
func TestValidateTagResourceCore(t *testing.T) {
	svc := &IoTService{}
	store := newTagValidationTestStore(t)
	arn := func(kind, name string) string {
		return "arn:aws:iot:us-east-1:000000000000:" + kind + "/" + name
	}

	t.Run("seeded resources resolve", func(t *testing.T) {
		for _, resource := range []string{
			arn("thing", "pin-thing"),
			arn("custommetric", "cm-pin"),
			arn("dimension", "dim-pin"),
			arn("cacert", "cacert-pin"),
		} {
			if err := svc.validateTagResourceCore(store, resource); err != nil {
				t.Errorf("validateTagResourceCore(%q) = %v, want nil", resource, err)
			}
		}
	})

	t.Run("missing resource is ResourceNotFoundException", func(t *testing.T) {
		for _, resource := range []string{
			arn("thing", "absent"),
			arn("policy", "absent"),
			arn("cert", "absent"),
			arn("securityprofile", "absent"),
			arn("custommetric", "absent"),
			arn("scheduledaudit", "absent"),
			arn("fleetmetric", "absent"),
		} {
			assertTagResourceNotFound(t, svc, store, resource)
		}
	})

	t.Run("unrecognised ARN kind is rejected", func(t *testing.T) {
		assertTagResourceNotFound(t, svc, store, arn("something-else", "whatever"))
		assertTagResourceNotFound(t, svc, store, "arn:aws:iot:us-east-1:000000000000:thing")
	})
}

func assertTagResourceNotFound(t *testing.T, svc *IoTService, store iotstore.IotStoreInterface, resourceArn string) {
	t.Helper()
	err := svc.validateTagResourceCore(store, resourceArn)
	if err == nil {
		t.Fatalf("validateTagResourceCore(%q) = nil, want not-found", resourceArn)
	}
	var awsErr *awserrors.AWSError
	if !errors.As(err, &awsErr) {
		t.Fatalf("validateTagResourceCore(%q) = %T, want *awserrors.AWSError", resourceArn, err)
	}
	if awsErr.Code != "ResourceNotFoundException" {
		t.Errorf("code = %q, want ResourceNotFoundException", awsErr.Code)
	}
}
