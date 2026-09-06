package apigateway

import (
	"testing"

	"vorpalstacks/internal/store/aws/apigateway"
)

// TestApplyRestApiPatchEndpointAndCompression pins the endpoint and
// compression rows of the official UpdateRestApi table: the
// minimumCompressionSize footnote's disable form, the ipAddressType row, the
// vpcEndpointIds exact row path, and the binaryMediaTypes element-token gate.
func TestApplyRestApiPatchEndpointAndCompression(t *testing.T) {
	api := &apigateway.RestApi{}

	// The footnote: "To disable compression, apply a replace operation with
	// the value property set to null or omit the value property." A null or
	// absent value arrives as the empty string.
	size := int32(2048)
	api.MinimumCompressionSize = &size
	for _, value := range []string{"", "null"} {
		if _, err := applyRestApiPatch(api, PatchOperation{Op: "replace", Path: "/minimumCompressionSize", Value: value}); err != nil {
			t.Fatalf("disable with value %q failed: %v", value, err)
		}
		if api.MinimumCompressionSize != nil {
			t.Fatalf("disable with value %q left the setting: %v", value, *api.MinimumCompressionSize)
		}
	}
	if _, err := applyRestApiPatch(api, PatchOperation{Op: "replace", Path: "/minimumCompressionSize", Value: "1024"}); err != nil {
		t.Fatalf("re-enable failed: %v", err)
	}
	if api.MinimumCompressionSize == nil || *api.MinimumCompressionSize != 1024 {
		t.Fatalf("re-enable not applied: %v", api.MinimumCompressionSize)
	}

	// ipAddressType: replace only, ipv4|dualstack per the model enum.
	if _, err := applyRestApiPatch(api, PatchOperation{Op: "add", Path: "/endpointConfiguration/ipAddressType", Value: "dualstack"}); err == nil {
		t.Fatal("ipAddressType add accepted")
	}
	if _, err := applyRestApiPatch(api, PatchOperation{Op: "replace", Path: "/endpointConfiguration/ipAddressType", Value: "ipv6"}); err == nil {
		t.Fatal("ipAddressType ipv6 accepted")
	}
	if _, err := applyRestApiPatch(api, PatchOperation{Op: "replace", Path: "/endpointConfiguration/ipAddressType", Value: "dualstack"}); err != nil {
		t.Fatalf("ipAddressType replace failed: %v", err)
	}
	if api.EndpointConfiguration.IpAddressType != "dualstack" {
		t.Fatalf("ipAddressType not stored: %+v", api.EndpointConfiguration)
	}

	// The vpcEndpointIds row is the exact member path; a sub-path addresses
	// no documented row (unhandled paths surface the unknown-patch-path
	// error at the caller).
	api.EndpointConfiguration = &apigateway.EndpointConfiguration{Types: []string{"PRIVATE"}}
	if handled, err := applyRestApiPatch(api, PatchOperation{Op: "add", Path: "/endpointConfiguration/vpcEndpointIds/0", Value: "vpce-1"}); err == nil && handled {
		t.Fatal("vpcEndpointIds sub-path accepted")
	}
	if _, err := applyRestApiPatch(api, PatchOperation{Op: "add", Path: "/endpointConfiguration/vpcEndpointIds", Value: "vpce-1"}); err != nil {
		t.Fatalf("vpcEndpointIds add failed: %v", err)
	}
	if len(api.EndpointConfiguration.VpcEndpointIds) != 1 {
		t.Fatalf("vpcEndpointIds not stored: %+v", api.EndpointConfiguration)
	}

	// binaryMediaTypes element tokens: numeric indexes and the append
	// marker reject; the value-token remove keeps working.
	api.BinaryMediaTypes = []string{"image/png"}
	for _, po := range []PatchOperation{
		{Op: "remove", Path: "/binaryMediaTypes/0"},
		{Op: "add", Path: "/binaryMediaTypes/2", Value: "image/jpeg"},
		{Op: "remove", Path: "/binaryMediaTypes/-"},
		{Op: "add", Path: "/binaryMediaTypes/", Value: "image/jpeg"},
	} {
		if _, err := applyRestApiPatch(api, po); err == nil {
			t.Fatalf("op %s on %s accepted", po.Op, po.Path)
		}
	}
	if _, err := applyRestApiPatch(api, PatchOperation{Op: "remove", Path: "/binaryMediaTypes/image~1png"}); err != nil {
		t.Fatalf("value-token remove failed: %v", err)
	}
	if len(api.BinaryMediaTypes) != 0 {
		t.Fatalf("value-token remove did not clear the entry: %v", api.BinaryMediaTypes)
	}
}
