package cloudfront

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"vorpalstacks/internal/common/invokers"
	cfstore "vorpalstacks/internal/store/aws/cloudfront"
	vcrypto "vorpalstacks/internal/utils/crypto"
)

// mustParseURL parses a test path; a malformed literal is a test bug.
func mustParseURL(path string) *url.URL {
	u, err := url.Parse(path)
	if err != nil {
		panic(err)
	}
	return u
}

// acmMaterialProvider is a test ACM certificate provider serving fixed
// material.
type acmMaterialProvider struct {
	material map[string]invokers.TLSCertificateMaterial
	requests []string
}

func (p *acmMaterialProvider) CertificateMaterial(ctx context.Context, region, certArn string) (invokers.TLSCertificateMaterial, error) {
	p.requests = append(p.requests, region+"/"+certArn)
	material, ok := p.material[certArn]
	if !ok {
		return invokers.TLSCertificateMaterial{}, fmt.Errorf("certificate not found: %s", certArn)
	}
	return material, nil
}

// iamMaterialProvider is a test IAM server certificate provider serving
// fixed material.
type iamMaterialProvider struct {
	material map[string]invokers.TLSCertificateMaterial
}

func (p *iamMaterialProvider) ServerCertificateMaterial(ctx context.Context, serverCertificateId string) (invokers.TLSCertificateMaterial, error) {
	material, ok := p.material[serverCertificateId]
	if !ok {
		return invokers.TLSCertificateMaterial{}, fmt.Errorf("server certificate not found: %s", serverCertificateId)
	}
	return material, nil
}

// testCertificatePEM generates a self-signed certificate and returns its
// PEM material.
func testCertificatePEM(t *testing.T, name string) (invokers.TLSCertificateMaterial, *rsa.PrivateKey) {
	t.Helper()
	key, err := vcrypto.GenerateRSAKey(2048)
	if err != nil {
		t.Fatal(err)
	}
	serial, err := vcrypto.GenerateSerialNumber()
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	template := &x509.Certificate{
		SerialNumber: serial,
		Subject:      pkix.Name{CommonName: name},
		NotBefore:    now,
		NotAfter:     now.AddDate(1, 0, 0),
		DNSNames:     []string{name},
	}
	certDER, err := vcrypto.CreateCertificate(template, template, key.Public(), key)
	if err != nil {
		t.Fatal(err)
	}
	keyPEM, err := vcrypto.EncodePrivateKeyPEM(key)
	if err != nil {
		t.Fatal(err)
	}
	return invokers.TLSCertificateMaterial{
		Certificate: string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})),
		PrivateKey:  keyPEM,
	}, key
}

func TestParseTLSCertificateMaterial(t *testing.T) {
	material, key := testCertificatePEM(t, "parse.example.com")

	cert, err := parseTLSCertificateMaterial(material)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	leaf, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		t.Fatalf("leaf not parseable: %v", err)
	}
	if leaf.Subject.CommonName != "parse.example.com" {
		t.Fatalf("leaf CN = %q", leaf.Subject.CommonName)
	}
	if cert.PrivateKey == nil {
		t.Fatal("private key missing")
	}
	if rsaKey, ok := cert.PrivateKey.(*rsa.PrivateKey); !ok || rsaKey.N.Cmp(key.N) != 0 {
		t.Fatal("private key does not match the material")
	}

	// A chain certificate is appended after the leaf.
	chainMaterial, _ := testCertificatePEM(t, "chain.example.com")
	withChain := invokers.TLSCertificateMaterial{
		Certificate:      material.Certificate,
		PrivateKey:       material.PrivateKey,
		CertificateChain: chainMaterial.Certificate,
	}
	cert, err = parseTLSCertificateMaterial(withChain)
	if err != nil {
		t.Fatalf("parse with chain failed: %v", err)
	}
	if len(cert.Certificate) != 2 {
		t.Fatalf("chain not appended: %d blocks", len(cert.Certificate))
	}

	if _, err := parseTLSCertificateMaterial(invokers.TLSCertificateMaterial{Certificate: material.Certificate}); err == nil {
		t.Fatal("missing private key must fail")
	}
	if _, err := parseTLSCertificateMaterial(invokers.TLSCertificateMaterial{Certificate: material.Certificate, PrivateKey: "not pem"}); err == nil {
		t.Fatal("non-PEM private key must fail")
	}
	if _, err := parseTLSCertificateMaterial(invokers.TLSCertificateMaterial{Certificate: "not pem", PrivateKey: material.PrivateKey}); err == nil {
		t.Fatal("non-PEM certificate must fail")
	}
}

// newTLSTestServer builds a DistributionServer with the certificate
// providers attached and one enabled distribution whose configuration the
// mutation callback adjusts.
func newTLSTestServer(t *testing.T, acm *acmMaterialProvider, iam *iamMaterialProvider, mutateConfig func(*cfstore.DistributionConfig)) (*DistributionServer, *cfstore.Distribution) {
	t.Helper()
	server, dist, _ := newCacheTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}, mutateConfig)
	server.certCache = newTLSCertificateCache()
	server.SetTLSCertificateProviders(acm, iam)
	return server, dist
}

func TestTLSCertificateResolvesViewerCertificates(t *testing.T) {
	acmMaterial, _ := testCertificatePEM(t, "acm.example.com")
	iamMaterial, _ := testCertificatePEM(t, "iam.example.com")
	acm := &acmMaterialProvider{material: map[string]invokers.TLSCertificateMaterial{
		"arn:aws:acm:us-east-1:123456789012:certificate/aaa": acmMaterial,
	}}
	iam := &iamMaterialProvider{material: map[string]invokers.TLSCertificateMaterial{
		"ASCAEXAMPLE0001": iamMaterial,
	}}

	cases := []struct {
		name    string
		mutate  func(*cfstore.DistributionConfig)
		sniName string
		wantCN  string
	}{
		{
			name: "ACM certificate by alias SNI",
			mutate: func(cfg *cfstore.DistributionConfig) {
				cfg.Aliases = &cfstore.Aliases{Quantity: 1, Items: []string{"acm.example.com"}}
				cfg.ViewerCertificate = &cfstore.ViewerCertificate{ACMCertificateArn: "arn:aws:acm:us-east-1:123456789012:certificate/aaa"}
			},
			sniName: "acm.example.com",
			wantCN:  "acm.example.com",
		},
		{
			name: "ACM certificate by cloudfront.net SNI",
			mutate: func(cfg *cfstore.DistributionConfig) {
				cfg.ViewerCertificate = &cfstore.ViewerCertificate{ACMCertificateArn: "arn:aws:acm:us-east-1:123456789012:certificate/aaa"}
			},
			sniName: "by-id.cloudfront.net", // replaced with the real domain below
			wantCN:  "acm.example.com",
		},
		{
			name: "IAM server certificate",
			mutate: func(cfg *cfstore.DistributionConfig) {
				cfg.ViewerCertificate = &cfstore.ViewerCertificate{IAMCertificateId: "ASCAEXAMPLE0001"}
			},
			sniName: "by-id.cloudfront.net",
			wantCN:  "iam.example.com",
		},
		{
			name: "CloudFront default certificate synthesises for the name",
			mutate: func(cfg *cfstore.DistributionConfig) {
				cfg.ViewerCertificate = &cfstore.ViewerCertificate{CloudFrontDefaultCertificate: true}
			},
			sniName: "by-id.cloudfront.net",
			wantCN:  "", // replaced with the real domain below
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server, dist := newTLSTestServer(t, acm, iam, tc.mutate)
			sniName := strings.Replace(tc.sniName, "by-id", dist.ID, 1)
			wantCN := tc.wantCN
			if wantCN == "" {
				wantCN = sniName
			}
			cert, err := server.TLSCertificate(&tls.ClientHelloInfo{ServerName: sniName})
			if err != nil {
				t.Fatalf("resolution failed: %v", err)
			}
			leaf, err := x509.ParseCertificate(cert.Certificate[0])
			if err != nil {
				t.Fatalf("leaf not parseable: %v", err)
			}
			if leaf.Subject.CommonName != wantCN {
				t.Fatalf("served CN = %q, want %q", leaf.Subject.CommonName, wantCN)
			}
		})
	}

	if len(acm.requests) == 0 {
		t.Fatal("ACM provider was never consulted")
	}
	if !strings.Contains(strings.Join(acm.requests, ","), "us-east-1/") {
		t.Fatalf("ACM provider region not derived from the ARN: %v", acm.requests)
	}
}

func TestTLSCertificateCachesAndPurges(t *testing.T) {
	acmMaterial, _ := testCertificatePEM(t, "cached.example.com")
	acm := &acmMaterialProvider{material: map[string]invokers.TLSCertificateMaterial{
		"arn:aws:acm:us-east-1:123456789012:certificate/bbb": acmMaterial,
	}}
	server, dist := newTLSTestServer(t, acm, nil, func(cfg *cfstore.DistributionConfig) {
		cfg.ViewerCertificate = &cfstore.ViewerCertificate{ACMCertificateArn: "arn:aws:acm:us-east-1:123456789012:certificate/bbb"}
	})
	name := dist.ID + ".cloudfront.net"

	first, err := server.TLSCertificate(&tls.ClientHelloInfo{ServerName: name})
	if err != nil {
		t.Fatal(err)
	}
	calls := len(acm.requests)
	second, err := server.TLSCertificate(&tls.ClientHelloInfo{ServerName: name})
	if err != nil {
		t.Fatal(err)
	}
	if len(acm.requests) != calls {
		t.Fatalf("cached resolution hit the provider again: %d -> %d", calls, len(acm.requests))
	}
	if first != second {
		t.Fatal("cached certificate identity changed")
	}

	server.PurgeCertificates()
	if _, err := server.TLSCertificate(&tls.ClientHelloInfo{ServerName: name}); err != nil {
		t.Fatal(err)
	}
	if len(acm.requests) == calls {
		t.Fatal("purge did not drop the cached certificate")
	}
}

func TestTLSCertificateUnknownNameServesDefault(t *testing.T) {
	server, _ := newTLSTestServer(t, &acmMaterialProvider{}, &iamMaterialProvider{}, nil)
	cert, err := server.TLSCertificate(&tls.ClientHelloInfo{ServerName: "unknown.example.com"})
	if err != nil {
		t.Fatalf("unknown SNI name must still serve the default certificate: %v", err)
	}
	leaf, _ := x509.ParseCertificate(cert.Certificate[0])
	if leaf.Subject.CommonName != "unknown.example.com" {
		t.Fatalf("default certificate CN = %q", leaf.Subject.CommonName)
	}

	// An empty SNI name keeps the handshake possible.
	if _, err := server.TLSCertificate(&tls.ClientHelloInfo{}); err != nil {
		t.Fatalf("empty SNI name must serve a fallback certificate: %v", err)
	}
}

func TestEnforceViewerProtocolPolicyRedirectsAndRefuses(t *testing.T) {
	server := &DistributionServer{}

	recorder := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://d123.cloudfront.net/obj?a=1", nil)
	if server.enforceViewerProtocolPolicy(recorder, r, "redirect-to-https") {
		t.Fatal("redirect-to-https must answer the request")
	}
	if recorder.Code != http.StatusMovedPermanently {
		t.Fatalf("GET redirect status = %d, want 301", recorder.Code)
	}
	want := "https://d123.cloudfront.net:50108/obj?a=1"
	if got := recorder.Header().Get("Location"); got != want {
		t.Fatalf("Location = %q, want %q", got, want)
	}

	recorder = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodPost, "http://d123.cloudfront.net/obj", strings.NewReader("data"))
	if server.enforceViewerProtocolPolicy(recorder, r, "redirect-to-https") {
		t.Fatal("redirect-to-https must answer the request")
	}
	if recorder.Code != http.StatusTemporaryRedirect {
		t.Fatalf("POST redirect status = %d, want 307", recorder.Code)
	}

	// Below HTTP/1.1 the unsafe method cannot be repeated, so the policy
	// refuses the request instead of redirecting.
	recorder = httptest.NewRecorder()
	r = &http.Request{
		Method:     http.MethodPost,
		URL:        mustParseURL("/obj"),
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		ProtoMinor: 0,
		Header:     http.Header{},
		Host:       "d123.cloudfront.net",
		Body:       http.NoBody,
	}
	if server.enforceViewerProtocolPolicy(recorder, r, "redirect-to-https") {
		t.Fatal("redirect-to-https must answer the request")
	}
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("HTTP/1.0 POST status = %d, want 403", recorder.Code)
	}

	recorder = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "http://d123.cloudfront.net/obj", nil)
	if server.enforceViewerProtocolPolicy(recorder, r, "https-only") {
		t.Fatal("https-only must answer the request")
	}
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("https-only status = %d, want 403", recorder.Code)
	}

	recorder = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "http://d123.cloudfront.net/obj", nil)
	if !server.enforceViewerProtocolPolicy(recorder, r, "allow-all") {
		t.Fatal("allow-all must not answer the request")
	}
	if recorder.Code != http.StatusOK {
		t.Fatalf("allow-all must not write a response, got %d", recorder.Code)
	}
}

func TestHandleRequestAppliesViewerProtocolPolicy(t *testing.T) {
	server, dist, _ := newCacheTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}, func(cfg *cfstore.DistributionConfig) {
		cfg.DefaultCacheBehavior.ViewerProtocolPolicy = "redirect-to-https"
	})

	recorder := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://edge/obj", nil)
	r.Host = dist.ID + ".cloudfront.net:50104"
	server.HandleRequest(recorder, r)
	if recorder.Code != http.StatusMovedPermanently {
		t.Fatalf("plain-HTTP redirect status = %d, want 301", recorder.Code)
	}
	want := fmt.Sprintf("https://%s.cloudfront.net:50108/obj", dist.ID)
	if got := recorder.Header().Get("Location"); got != want {
		t.Fatalf("Location = %q, want %q", got, want)
	}

	// A TLS request serves normally under the same policy.
	r = httptest.NewRequest(http.MethodGet, "https://edge/obj", nil)
	r.Host = dist.ID + ".cloudfront.net"
	r.TLS = &tls.ConnectionState{}
	recorder = httptest.NewRecorder()
	server.HandleRequest(recorder, r)
	if recorder.Code != http.StatusOK {
		t.Fatalf("TLS request status = %d, want 200", recorder.Code)
	}

	httpsOnly, dist2, _ := newCacheTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}, func(cfg *cfstore.DistributionConfig) {
		cfg.DefaultCacheBehavior.ViewerProtocolPolicy = "https-only"
	})
	recorder = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "http://edge/obj", nil)
	r.Host = dist2.ID + ".cloudfront.net"
	httpsOnly.HandleRequest(recorder, r)
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("https-only plain-HTTP status = %d, want 403", recorder.Code)
	}
}
