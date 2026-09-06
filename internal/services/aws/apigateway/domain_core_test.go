package apigateway

import (
	"testing"

	"vorpalstacks/internal/store/aws/apigateway"
)

// TestApplyDomainNamePatch pins the UpdateDomainName table transcription:
// the certificate removes clear the addressed field, the add/remove
// same-request exclusion the cells state verbatim, the endpoint-types row
// (add appends / remove drops the addressed value, EDGE|REGIONAL), and the
// ipAddressType row.
func TestApplyDomainNamePatch(t *testing.T) {
	domain := &apigateway.DomainName{CertificateArn: "arn:aws:acm:us-east-1:1:certificate/edge"}

	// remove clears the addressed edge certificate; the regional one stays.
	domain.RegionalCertificateArn = "arn:aws:acm:us-east-1:1:certificate/regional"
	handled, err := applyDomainNamePatch(domain, PatchOperation{Op: "remove", Path: "/certificateArn"}, map[string]string{})
	if !handled || err != nil {
		t.Fatalf("edge remove failed: handled=%v err=%v", handled, err)
	}
	if domain.CertificateArn != "" || domain.RegionalCertificateArn == "" {
		t.Fatalf("edge remove did not clear the addressed field: %+v", domain)
	}

	// add then remove on the same path within one request rejects, per the
	// cells' "cannot be included with the ... operation in the same request".
	seen := map[string]string{}
	if _, err := applyDomainNamePatch(domain, PatchOperation{Op: "add", Path: "/regionalCertificateName", Value: "n"}, seen); err != nil {
		t.Fatalf("regional name add failed: %v", err)
	}
	if _, err := applyDomainNamePatch(domain, PatchOperation{Op: "remove", Path: "/regionalCertificateName"}, seen); err == nil {
		t.Fatal("same-request add+remove accepted")
	}

	// replace on the endpoint types row is Not supported; the add/remove
	// cells serve edge-optimized/regional updates only.
	if _, err := applyDomainNamePatch(domain, PatchOperation{Op: "replace", Path: "/endpointConfiguration/types", Value: "REGIONAL"}, map[string]string{}); err == nil {
		t.Fatal("types replace accepted")
	}
	if _, err := applyDomainNamePatch(domain, PatchOperation{Op: "add", Path: "/endpointConfiguration/types", Value: "PRIVATE"}, map[string]string{}); err == nil {
		t.Fatal("types add with PRIVATE accepted")
	}
	handled, err = applyDomainNamePatch(domain, PatchOperation{Op: "add", Path: "/endpointConfiguration/types", Value: "REGIONAL"}, map[string]string{})
	if !handled || err != nil {
		t.Fatalf("types add failed: %v", err)
	}
	if len(domain.EndpointConfiguration.Types) != 1 || domain.EndpointConfiguration.Types[0] != "REGIONAL" {
		t.Fatalf("types add not applied: %+v", domain.EndpointConfiguration)
	}

	// The developer guide's migration flow: the new type joins the existing
	// list (its output example shows "types": ["EDGE", "REGIONAL"] with both
	// coexisting until the DNS cutover) and removing the obsolete type
	// completes the transition.
	handled, err = applyDomainNamePatch(domain, PatchOperation{Op: "add", Path: "/endpointConfiguration/types", Value: "EDGE"}, map[string]string{})
	if !handled || err != nil {
		t.Fatalf("types add on a populated list failed: %v", err)
	}
	if len(domain.EndpointConfiguration.Types) != 2 ||
		domain.EndpointConfiguration.Types[0] != "REGIONAL" || domain.EndpointConfiguration.Types[1] != "EDGE" {
		t.Fatalf("types add did not append: %+v", domain.EndpointConfiguration)
	}
	handled, err = applyDomainNamePatch(domain, PatchOperation{Op: "add", Path: "/endpointConfiguration/types", Value: "EDGE"}, map[string]string{})
	if !handled || err != nil {
		t.Fatalf("repeated types add failed: %v", err)
	}
	if len(domain.EndpointConfiguration.Types) != 2 {
		t.Fatalf("repeated types add duplicated the value: %+v", domain.EndpointConfiguration)
	}
	handled, err = applyDomainNamePatch(domain, PatchOperation{Op: "remove", Path: "/endpointConfiguration/types", Value: "REGIONAL"}, map[string]string{})
	if !handled || err != nil {
		t.Fatalf("types remove failed: %v", err)
	}
	if len(domain.EndpointConfiguration.Types) != 1 || domain.EndpointConfiguration.Types[0] != "EDGE" {
		t.Fatalf("types remove did not drop the addressed value: %+v", domain.EndpointConfiguration)
	}

	// ipAddressType: replace only, ipv4|dualstack per the model enum.
	if _, err := applyDomainNamePatch(domain, PatchOperation{Op: "replace", Path: "/endpointConfiguration/ipAddressType", Value: "ipv6"}, map[string]string{}); err == nil {
		t.Fatal("ipAddressType ipv6 accepted")
	}
	if _, err := applyDomainNamePatch(domain, PatchOperation{Op: "add", Path: "/endpointConfiguration/ipAddressType", Value: "dualstack"}, map[string]string{}); err == nil {
		t.Fatal("ipAddressType add accepted")
	}
	if _, err := applyDomainNamePatch(domain, PatchOperation{Op: "replace", Path: "/endpointConfiguration/ipAddressType", Value: "dualstack"}, map[string]string{}); err != nil {
		t.Fatalf("ipAddressType replace failed: %v", err)
	}
	if domain.EndpointConfiguration.IpAddressType != "dualstack" {
		t.Fatalf("ipAddressType not stored: %+v", domain.EndpointConfiguration)
	}

	// The endpoint-types sub-path form is not the row the table documents.
	if handled, _ := applyDomainNamePatch(domain, PatchOperation{Op: "replace", Path: "/endpointConfiguration/types/REGIONAL"}, map[string]string{}); handled {
		t.Fatal("types sub-path handled")
	}

	// The invariant: a patch request may not leave the domain certificate-less.
	domain.CertificateArn = ""
	domain.RegionalCertificateArn = ""
	if err := requireDomainCertificate(domain); err == nil {
		t.Fatal("certificate-less domain accepted")
	}
}
