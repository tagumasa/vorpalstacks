// Distribution TLS serving: the SNI certificate resolver for the
// distribution plane's TLS listener and the ViewerProtocolPolicy
// enforcement shared by both planes.
package cloudfront

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/common/serviceports"
	"vorpalstacks/internal/config"
	"vorpalstacks/internal/utils/aws/arn"
	vcrypto "vorpalstacks/internal/utils/crypto"
)

// tlsCertificateCacheTTL bounds how long a parsed certificate or a
// synthesised default certificate stays cached. Re-imported material and
// renewed keys therefore propagate within this window without a restart.
const tlsCertificateCacheTTL = 5 * time.Minute

// tlsCertificateEntry is one cached, ready-to-serve certificate.
type tlsCertificateEntry struct {
	cert      *tls.Certificate
	expiresAt time.Time
}

// tlsCertificateCache memoises parsed certificates by source key so the
// handshake path avoids re-parsing PEM material on every connection.
type tlsCertificateCache struct {
	mu      sync.Mutex
	entries map[string]*tlsCertificateEntry
}

func newTLSCertificateCache() *tlsCertificateCache {
	return &tlsCertificateCache{entries: make(map[string]*tlsCertificateEntry)}
}

func (c *tlsCertificateCache) get(key string) *tls.Certificate {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.entries[key]
	if !ok || time.Now().After(entry.expiresAt) {
		return nil
	}
	return entry.cert
}

func (c *tlsCertificateCache) put(key string, cert *tls.Certificate) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = &tlsCertificateEntry{cert: cert, expiresAt: time.Now().Add(tlsCertificateCacheTTL)}
}

func (c *tlsCertificateCache) purge() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]*tlsCertificateEntry)
}

// SetTLSCertificateProviders injects the certificate material providers the
// TLS plane resolves viewer certificates through. Either may be nil, in
// which case the corresponding certificate source falls back to the
// synthesised default certificate.
func (s *DistributionServer) SetTLSCertificateProviders(acm invokers.ACMCertificateProvider, iam invokers.IAMServerCertificateProvider) {
	s.providerMu.Lock()
	s.acmCertificates = acm
	s.iamCertificates = iam
	s.providerMu.Unlock()
}

func (s *DistributionServer) certificateProviders() (invokers.ACMCertificateProvider, invokers.IAMServerCertificateProvider) {
	s.providerMu.RLock()
	defer s.providerMu.RUnlock()
	return s.acmCertificates, s.iamCertificates
}

// PurgeCertificates drops every cached TLS certificate. Distribution
// configuration changes call it so a switched certificate serves
// immediately.
func (s *DistributionServer) PurgeCertificates() {
	s.certCache.purge()
}

// TLSCertificate resolves the certificate for a distribution-plane TLS
// handshake. The SNI name maps to a distribution exactly like the Host
// header does on the plain plane (cloudfront.net first label or CNAME
// alias); the distribution's ViewerCertificate then selects the served
// material: the referenced ACM certificate, the referenced IAM server
// certificate, or the CloudFront default certificate, which this platform
// synthesises as a self-signed certificate for the requested name.
func (s *DistributionServer) TLSCertificate(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
	name := strings.ToLower(strings.TrimSpace(info.ServerName))
	name = strings.TrimSuffix(name, ".")

	if name != "" {
		store := s.getDistributionStore()
		if store != nil {
			if id := s.resolveDistributionID(name); id != "" {
				if dist, err := store.Get(id); err == nil && dist != nil && dist.DistributionConfig != nil && dist.DistributionConfig.ViewerCertificate != nil {
					vc := dist.DistributionConfig.ViewerCertificate
					acmProvider, iamProvider := s.certificateProviders()
					if vc.ACMCertificateArn != "" && acmProvider != nil {
						return s.cachedACMCertificate(acmProvider, vc.ACMCertificateArn)
					}
					if vc.IAMCertificateId != "" && iamProvider != nil {
						return s.cachedIAMCertificate(iamProvider, vc.IAMCertificateId)
					}
				}
			}
		}
	}

	return s.defaultViewerCertificate(name)
}

// cachedACMCertificate resolves an ACM certificate by ARN through the
// provider, deriving the region from the ARN itself, and caches the parsed
// result.
func (s *DistributionServer) cachedACMCertificate(provider invokers.ACMCertificateProvider, certArn string) (*tls.Certificate, error) {
	key := "acm:" + certArn
	if cert := s.certCache.get(key); cert != nil {
		return cert, nil
	}
	region := ""
	if parsed, err := arn.ParseARN(certArn); err == nil {
		region = parsed.Region
	}
	material, err := provider.CertificateMaterial(context.Background(), region, certArn)
	if err != nil {
		return nil, fmt.Errorf("cloudfront: ACM certificate material for %s: %w", certArn, err)
	}
	cert, err := parseTLSCertificateMaterial(material)
	if err != nil {
		return nil, fmt.Errorf("cloudfront: ACM certificate %s: %w", certArn, err)
	}
	s.certCache.put(key, cert)
	return cert, nil
}

// cachedIAMCertificate resolves an IAM server certificate by its unique
// certificate ID through the provider and caches the parsed result.
func (s *DistributionServer) cachedIAMCertificate(provider invokers.IAMServerCertificateProvider, certID string) (*tls.Certificate, error) {
	key := "iam:" + certID
	if cert := s.certCache.get(key); cert != nil {
		return cert, nil
	}
	material, err := provider.ServerCertificateMaterial(context.Background(), certID)
	if err != nil {
		return nil, fmt.Errorf("cloudfront: IAM server certificate material for %s: %w", certID, err)
	}
	cert, err := parseTLSCertificateMaterial(material)
	if err != nil {
		return nil, fmt.Errorf("cloudfront: IAM server certificate %s: %w", certID, err)
	}
	s.certCache.put(key, cert)
	return cert, nil
}

// defaultViewerCertificate synthesises the self-signed certificate the
// platform serves for distribution names without attached material — the
// equivalent of the CloudFront default certificate that covers the assigned
// cloudfront.net domain. An empty SNI name still completes the handshake;
// Host-based routing happens after the handshake.
func (s *DistributionServer) defaultViewerCertificate(name string) (*tls.Certificate, error) {
	if name == "" {
		name = "localhost"
	}
	key := "default:" + name
	if cert := s.certCache.get(key); cert != nil {
		return cert, nil
	}

	rsaKey, err := vcrypto.GenerateRSAKey(2048)
	if err != nil {
		return nil, fmt.Errorf("cloudfront: default viewer certificate key: %w", err)
	}
	cert, err := synthesiseViewerCertificate(name, rsaKey)
	if err != nil {
		return nil, err
	}
	s.certCache.put(key, cert)
	return cert, nil
}

// synthesiseViewerCertificate builds the self-signed certificate covering
// the requested viewer name.
func synthesiseViewerCertificate(name string, key *rsa.PrivateKey) (*tls.Certificate, error) {
	serial, err := vcrypto.GenerateSerialNumber()
	if err != nil {
		return nil, fmt.Errorf("cloudfront: default viewer certificate serial: %w", err)
	}
	now := time.Now().UTC()
	template := &x509.Certificate{
		SerialNumber: serial,
		Subject:      pkix.Name{CommonName: name},
		NotBefore:    now,
		NotAfter:     now.AddDate(1, 0, 0),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		DNSNames:     []string{name},
	}
	certDER, err := vcrypto.CreateCertificate(template, template, key.Public(), key)
	if err != nil {
		return nil, fmt.Errorf("cloudfront: default viewer certificate: %w", err)
	}
	return &tls.Certificate{
		Certificate: [][]byte{certDER},
		PrivateKey:  key,
	}, nil
}

// parseTLSCertificateMaterial turns PEM certificate material into a
// serveable tls.Certificate: the certificate blocks of the leaf material
// lead (the first is the leaf), followed by any chain blocks.
func parseTLSCertificateMaterial(material invokers.TLSCertificateMaterial) (*tls.Certificate, error) {
	if material.PrivateKey == "" {
		return nil, fmt.Errorf("private key is missing")
	}
	keyBlock, _ := pem.Decode([]byte(material.PrivateKey))
	if keyBlock == nil {
		return nil, fmt.Errorf("private key is not PEM-encoded")
	}
	var privateKey any
	if key, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes); err == nil {
		privateKey = key
	} else if key, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes); err == nil {
		privateKey = key
	} else if key, err := x509.ParseECPrivateKey(keyBlock.Bytes); err == nil {
		privateKey = key
	} else {
		return nil, fmt.Errorf("private key is not parseable")
	}

	var derChain [][]byte
	appendCerts := func(pemText string) {
		rest := []byte(pemText)
		for {
			var block *pem.Block
			block, rest = pem.Decode(rest)
			if block == nil {
				return
			}
			if block.Type == "CERTIFICATE" {
				derChain = append(derChain, block.Bytes)
			}
		}
	}
	appendCerts(material.Certificate)
	if len(derChain) == 0 {
		return nil, fmt.Errorf("certificate is not PEM-encoded")
	}
	appendCerts(material.CertificateChain)
	return &tls.Certificate{
		Certificate: derChain,
		PrivateKey:  privateKey,
	}, nil
}

// enforceViewerProtocolPolicy applies a cache behaviour's viewer protocol
// policy to plain-HTTP requests. redirect-to-https answers GET and HEAD
// with 301 (Moved Permanently) per the AWS contract, the remaining methods
// with 307 (Temporary Redirect) when the request's protocol version can
// repeat the method and payload, and 403 below HTTP/1.1; https-only answers
// 403 (Forbidden) and never returns the object. It returns false when it
// wrote the response itself.
func (s *DistributionServer) enforceViewerProtocolPolicy(w http.ResponseWriter, r *http.Request, policy string) bool {
	switch policy {
	case "https-only":
		http.Error(w, `{"message":"Viewer protocol policy requires HTTPS"}`, http.StatusForbidden)
		return false
	case "redirect-to-https":
		status := http.StatusMovedPermanently
		switch r.Method {
		case http.MethodGet, http.MethodHead:
		default:
			if r.ProtoMajor == 1 && r.ProtoMinor == 0 {
				http.Error(w, `{"message":"Viewer protocol policy requires a protocol version that can repeat the request"}`, http.StatusForbidden)
				return false
			}
			status = http.StatusTemporaryRedirect
		}
		w.Header().Set("Location", s.httpsRedirectTarget(r))
		w.WriteHeader(status)
		return false
	}
	return true
}

// httpsRedirectTarget builds the HTTPS URL a redirect-to-https response
// points at: the viewer's host and the request URI, redirected to this
// platform's distribution TLS port. The explicit port disappears only for
// the standard HTTPS port 443, matching URL semantics.
func (s *DistributionServer) httpsRedirectTarget(r *http.Request) string {
	host := r.Host
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}
	port := distributionTLSPort()
	if port != 443 {
		host = net.JoinHostPort(host, strconv.Itoa(port))
	}
	return "https://" + host + r.URL.RequestURI()
}

// distributionTLSPort resolves the configured distribution TLS port with
// the compiled default as the fallback for stores seeded before the key
// existed.
func distributionTLSPort() int {
	if port := config.GetInt("ports.cloudfront_tls"); port > 0 {
		return port
	}
	return serviceports.CloudFrontTLS
}
