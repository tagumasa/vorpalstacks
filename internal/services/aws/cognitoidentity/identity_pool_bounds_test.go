package cognitoidentity

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
)

// The IdentityProviders map behind SupportedLoginProviders carries the
// model's @length max of 10 on both the create and the update face.
func TestIdentityPoolLoginProviderBound(t *testing.T) {
	svc, real := newMergeTestService(t)

	providers := map[string]string{}
	for i := 0; i < 11; i++ {
		providers[fmt.Sprintf("provider%d.example.com", i)] = fmt.Sprintf("client-%d", i)
	}

	_, err := svc.createIdentityPoolCore(real, CreateIdentityPoolInput{
		IdentityPoolName:        "providers-bound",
		SupportedLoginProviders: providers,
	})
	if !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("create with 11 login providers returned %v, want InvalidParameterException", err)
	}

	pool, err := svc.createIdentityPoolCore(real, CreateIdentityPoolInput{
		IdentityPoolName: "providers-ok",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = svc.updateIdentityPoolCore(&request.RequestContext{Region: "us-east-1"}, UpdateIdentityPoolInput{
		IdentityPoolID:          pool.ID,
		SupportedLoginProviders: providers,
	})
	if !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("update with 11 login providers returned %v, want InvalidParameterException", err)
	}
}

// The platform-wide 50-tag quota applies to create, update and the
// TagResource merge.
func TestIdentityPoolTagQuota(t *testing.T) {
	svc, real := newMergeTestService(t)

	tags := func(n int) map[string]string {
		m := map[string]string{}
		for i := 0; i < n; i++ {
			m[fmt.Sprintf("key%d", i)] = "value"
		}
		return m
	}

	if _, err := svc.createIdentityPoolCore(real, CreateIdentityPoolInput{
		IdentityPoolName: "tag-quota",
		TagsProvided:     true,
		Tags:             tags(51),
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("create with 51 tags returned %v, want InvalidParameterException", err)
	}

	pool, err := svc.createIdentityPoolCore(real, CreateIdentityPoolInput{
		IdentityPoolName: "tag-quota-ok",
		TagsProvided:     true,
		Tags:             tags(50),
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := svc.updateIdentityPoolCore(&request.RequestContext{Region: "us-east-1"}, UpdateIdentityPoolInput{
		IdentityPoolID: pool.ID,
		TagsProvided:   true,
		Tags:           tags(51),
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("update with 51 tags returned %v, want InvalidParameterException", err)
	}

	// The TagResource merge: a 51st distinct key is refused on a resource
	// already holding 50.
	cfg := cognitoIdentityTagConfig(real)
	if err := cfg.TagFunc(nil, pool.Arn, []tagutil.Tag{{Key: "overflow", Value: "v"}}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("51st tag key returned %v, want InvalidParameterException", err)
	}
}

// GetOpenIdTokenForDeveloperIdentity applies the LoginsMap bounds: at most
// 10 entries with provider-name keys.
func TestDeveloperIdentityLoginsBounds(t *testing.T) {
	svc, _ := newMergeTestService(t)

	logins := map[string]string{}
	for i := 0; i < 11; i++ {
		logins[fmt.Sprintf("provider%d.example.com", i)] = "token"
	}
	_, err := svc.getOpenIdTokenForDeveloperIdentityCore(&request.RequestContext{Region: "us-east-1"}, GetOpenIdTokenForDeveloperIdentityInput{
		IdentityPoolID: "us-east-1:0123456789abcdef0123456789abcdef01234567",
		Logins:         logins,
	})
	if !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("11 logins returned %v, want InvalidParameterException", err)
	}

	longKey := strings.Repeat("k", 129)
	_, err = svc.getOpenIdTokenForDeveloperIdentityCore(&request.RequestContext{Region: "us-east-1"}, GetOpenIdTokenForDeveloperIdentityInput{
		IdentityPoolID: "us-east-1:0123456789abcdef0123456789abcdef01234567",
		Logins:         map[string]string{longKey: "token"},
	})
	if !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("129-character provider key returned %v, want InvalidParameterException", err)
	}
}
