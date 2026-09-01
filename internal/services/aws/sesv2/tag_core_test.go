package sesv2

import (
	"errors"
	"testing"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/storage"
	sesv2store "vorpalstacks/internal/store/aws/sesv2"
)

// newTagTargetTestStore opens a throwaway SESv2 store seeded with one
// resource of every taggable kind.
func newTagTargetTestStore(t *testing.T) sesv2store.SESv2StoreInterface {
	t.Helper()
	raw, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { raw.Close() })
	store := sesv2store.NewSESv2Store(raw, "000000000000", "us-east-1")

	seed := map[string]func() error{
		"identity": func() error {
			_, err := store.CreateEmailIdentity(sesv2store.NewEmailIdentity("example.com"))
			return err
		},
		"configuration-set": func() error {
			_, err := store.CreateConfigurationSet(sesv2store.NewConfigurationSet("daily"))
			return err
		},
		"template": func() error {
			_, err := store.CreateEmailTemplate(sesv2store.NewEmailTemplate("welcome"))
			return err
		},
		"contact-list": func() error {
			_, err := store.CreateContactList(sesv2store.NewContactList("newsletter"))
			return err
		},
		"dedicated-ip-pool": func() error {
			return store.CreateDedicatedIpPool(&sesv2store.DedicatedIpPool{PoolName: "marketing"})
		},
	}
	for kind, create := range seed {
		if err := create(); err != nil {
			t.Fatalf("seed %s: %v", kind, err)
		}
	}
	return store
}

// TestValidateSesv2TagTarget pins the tag-target resolution: each taggable
// ARN kind resolves its seeded resource, a missing resource of any kind is
// the modelled NotFoundException, and an ARN of any other kind addresses no
// resource this service can own.
func TestValidateSesv2TagTarget(t *testing.T) {
	store := newTagTargetTestStore(t)
	arn := func(kind, name string) string {
		return "arn:aws:ses:us-east-1:000000000000:" + kind + "/" + name
	}

	t.Run("seeded resources resolve", func(t *testing.T) {
		for _, resource := range []string{
			arn("identity", "example.com"),
			arn("configuration-set", "daily"),
			arn("template", "welcome"),
			arn("contact-list", "newsletter"),
			arn("dedicated-ip-pool", "marketing"),
		} {
			if err := validateSesv2TagTarget(store, resource); err != nil {
				t.Errorf("validateSesv2TagTarget(%q) = %v, want nil", resource, err)
			}
		}
	})

	t.Run("missing resource is NotFoundException", func(t *testing.T) {
		for _, resource := range []string{
			arn("identity", "absent.example"),
			arn("configuration-set", "absent"),
			arn("template", "absent"),
			arn("contact-list", "absent"),
			arn("dedicated-ip-pool", "absent"),
		} {
			assertTagTargetNotFound(t, store, resource)
		}
	})

	t.Run("unrecognised ARN kind is rejected", func(t *testing.T) {
		assertTagTargetNotFound(t, store, arn("something-else", "whatever"))
		assertTagTargetNotFound(t, store, "arn:aws:ses:us-east-1:000000000000:identity")
	})
}

func assertTagTargetNotFound(t *testing.T, store sesv2store.SESv2StoreInterface, resourceArn string) {
	t.Helper()
	err := validateSesv2TagTarget(store, resourceArn)
	if err == nil {
		t.Fatalf("validateSesv2TagTarget(%q) = nil, want not-found", resourceArn)
	}
	var awsErr *awserrors.AWSError
	if !errors.As(err, &awsErr) {
		t.Fatalf("validateSesv2TagTarget(%q) = %T, want *awserrors.AWSError", resourceArn, err)
	}
	if awsErr.Code != "NotFoundException" {
		t.Errorf("code = %q, want NotFoundException", awsErr.Code)
	}
}
