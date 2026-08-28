package appsync

import (
	"strings"
	"testing"
)

// The typed AWS SDK validates the required members of these operations
// client-side, so the server-side rejections never surface through the SDK
// test suite; they are pinned here at the Core boundary instead. Every
// rejection below fires before any store access, so a nil store is safe.

func TestListTypesCoreRequiresFormat(t *testing.T) {
	svc := &AppSyncService{}

	_, _, err := svc.listTypesCore(nil, "api-id", "", 25, "")
	if err == nil || !strings.Contains(err.Error(), "format is required") {
		t.Fatalf("expected format-required rejection, got %v", err)
	}

	_, _, err = svc.listTypesCore(nil, "api-id", "YAML", 25, "")
	if err == nil || !strings.Contains(err.Error(), "Invalid format") {
		t.Fatalf("expected invalid-format rejection, got %v", err)
	}
}

func TestUpdateApiCacheCoreRequiresMembers(t *testing.T) {
	svc := &AppSyncService{}

	full := updateApiCacheInput{
		ApiId:                 "api-id",
		Type:                  "SMALL",
		HasType:               true,
		Ttl:                   300,
		HasTtl:                true,
		ApiCachingBehavior:    "PER_RESOLVER_CACHING",
		HasApiCachingBehavior: true,
	}

	in := full
	in.HasType = false
	if _, err := svc.updateApiCacheCore(nil, in); err == nil || !strings.Contains(err.Error(), "type is required") {
		t.Fatalf("expected type-required rejection, got %v", err)
	}

	in = full
	in.HasTtl = false
	if _, err := svc.updateApiCacheCore(nil, in); err == nil || !strings.Contains(err.Error(), "ttl is required") {
		t.Fatalf("expected ttl-required rejection, got %v", err)
	}

	in = full
	in.HasApiCachingBehavior = false
	if _, err := svc.updateApiCacheCore(nil, in); err == nil || !strings.Contains(err.Error(), "apiCachingBehavior is required") {
		t.Fatalf("expected apiCachingBehavior-required rejection, got %v", err)
	}

	// An explicitly supplied zero ttl passes the presence gate and is
	// rejected by the documented 1-3600 window.
	in = full
	in.Ttl = 0
	if _, err := svc.updateApiCacheCore(nil, in); err == nil || !strings.Contains(err.Error(), "ttl must be between 1 and 3600 seconds") {
		t.Fatalf("expected ttl-range rejection, got %v", err)
	}
}

func TestPutGraphqlApiEnvironmentVariablesCoreRequiresMap(t *testing.T) {
	svc := &AppSyncService{}

	_, err := svc.putGraphqlApiEnvironmentVariablesCore(nil, "api-id", nil)
	if err == nil || !strings.Contains(err.Error(), "environmentVariables is required") {
		t.Fatalf("expected environmentVariables-required rejection, got %v", err)
	}
}

func TestGetIntrospectionSchemaCoreRequiresFormat(t *testing.T) {
	svc := &AppSyncService{}

	_, err := svc.getIntrospectionSchemaCore(nil, getIntrospectionSchemaInput{ApiId: "api-id"})
	if err == nil || !strings.Contains(err.Error(), "format is required") {
		t.Fatalf("expected format-required rejection, got %v", err)
	}

	_, err = svc.getIntrospectionSchemaCore(nil, getIntrospectionSchemaInput{ApiId: "api-id", Format: "YAML"})
	if err == nil || !strings.Contains(err.Error(), "Invalid format: YAML") {
		t.Fatalf("expected format-enum rejection, got %v", err)
	}
}

func TestDisassociateMergedGraphqlApiCoreRequiresSourceApi(t *testing.T) {
	svc := &AppSyncService{}

	_, err := svc.disassociateMergedGraphqlApiCore(nil, "", "assoc-id")
	if err == nil || !strings.Contains(err.Error(), "sourceApiIdentifier is required") {
		t.Fatalf("expected sourceApiIdentifier-required rejection, got %v", err)
	}
}
