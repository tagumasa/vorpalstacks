package acm

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	acmstorelib "vorpalstacks/internal/store/aws/acm"

	vcrypto "vorpalstacks/internal/utils/crypto"

	"golang.org/x/crypto/pbkdf2"
)

// renewCertificateMaterial regenerates the X.509 certificate material (new
// key pair, new serial, new PEM blob) for an existing AMAZON_ISSUED
// certificate during renewal. The ARN is preserved per AWS specification.
// The caller is responsible for setting RenewalSummary and persisting via
// stores.certificates.Update afterwards.
func renewCertificateMaterial(cert *acmstorelib.Certificate) error {
	key, err := generateKeyForKeyAlgorithm(cert.KeyAlgorithm)
	if err != nil {
		return err
	}

	serialBigInt, err := vcrypto.GenerateSerialNumber()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	notAfter := now.AddDate(1, 0, 0)

	dnsNames := make([]string, 0, 1+len(cert.SubjectAlternativeNames))
	dnsNames = append(dnsNames, cert.DomainName)
	dnsNames = append(dnsNames, cert.SubjectAlternativeNames...)

	template := &x509.Certificate{
		SerialNumber: serialBigInt,
		Subject:      pkix.Name{CommonName: cert.DomainName},
		NotBefore:    now,
		NotAfter:     notAfter,
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		DNSNames:     dnsNames,
	}

	certDER, err := vcrypto.CreateCertificate(template, template, key.Public(), key)
	if err != nil {
		return err
	}

	cert.Serial = serialBigInt.String()
	cert.Certificate = vcrypto.EncodeCertificatePEM(certDER)
	cert.Status = "ISSUED"
	cert.NotBefore = now
	cert.NotAfter = notAfter
	cert.IssuedAt = now
	cert.SignatureAlgorithm = signatureAlgorithmForKeyAlgorithm(cert.KeyAlgorithm)
	return nil
}

// generateAmazonIssuedCert generates a self-signed AMAZON_ISSUED certificate
// with the given parameters. Shared by the HTTP API (RequestCertificate) and
// the admin handler. The caller is responsible for setting AccountID, Region,
// Tags, ManagedBy, CertificateAuthorityArn, and Options afterwards.
func generateAmazonIssuedCert(certArn, domainName string, sans []string, keyAlgorithm, validationMethod string) (*acmstorelib.Certificate, error) {
	// Build the complete DNSNames list for the certificate template.
	dnsNames := make([]string, 0, 1+len(sans))
	dnsNames = append(dnsNames, domainName)
	dnsNames = append(dnsNames, sans...)

	now := time.Now().UTC()
	key, err := generateKeyForKeyAlgorithm(keyAlgorithm)
	if err != nil {
		return nil, NewInvalidParameterError(fmt.Sprintf("Unsupported KeyAlgorithm: %s", keyAlgorithm))
	}

	serialBigInt, err := vcrypto.GenerateSerialNumber()
	if err != nil {
		return nil, NewInternalServerException("Failed to generate serial")
	}
	notAfter := now.AddDate(1, 0, 0)
	template := &x509.Certificate{
		SerialNumber: serialBigInt,
		Subject:      pkix.Name{CommonName: domainName},
		NotBefore:    now,
		NotAfter:     notAfter,
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		DNSNames:     dnsNames,
	}

	certDER, err := vcrypto.CreateCertificate(template, template, key.Public(), key)
	if err != nil {
		return nil, NewInternalServerException("Failed to create certificate")
	}

	certPEM := vcrypto.EncodeCertificatePEM(certDER)

	// Build domain validation options. ACM creates certs in PENDING_VALIDATION
	// state; the platform validates synchronously (self-signed) so domain
	// validation options transition to SUCCESS immediately.
	domainValidationOptions := buildDomainValidationOptions(domainName, validationMethod, sans)
	for i := range domainValidationOptions {
		domainValidationOptions[i].ValidationStatus = "SUCCESS"
	}

	return &acmstorelib.Certificate{
		CertificateArn:           certArn,
		DomainName:               domainName,
		Serial:                   serialBigInt.String(),
		Status:                   "ISSUED",
		Type:                     "AMAZON_ISSUED",
		KeyAlgorithm:             keyAlgorithm,
		SignatureAlgorithm:       signatureAlgorithmForKeyAlgorithm(keyAlgorithm),
		RenewalEligibility:       "ELIGIBLE",
		CreatedAt:                now,
		Subject:                  domainName,
		Issuer:                   domainName,
		Certificate:              certPEM,
		NotBefore:                now,
		NotAfter:                 notAfter,
		IssuedAt:                 now,
		SubjectAlternativeNames:  sans,
		DomainValidationOptions:  domainValidationOptions,
		CertificateKeyPairOrigin: "AWS_MANAGED",
	}, nil
}

func buildDomainValidationOptions(domainName, validationMethod string, sans []string) []*acmstorelib.DomainValidation {
	// Build validation entries for the primary domain and all SANs.
	allDomains := make([]string, 0, 1+len(sans))
	allDomains = append(allDomains, domainName)
	allDomains = append(allDomains, sans...)

	options := make([]*acmstorelib.DomainValidation, 0, len(allDomains))
	for _, d := range allDomains {
		dv := &acmstorelib.DomainValidation{
			DomainName:       d,
			ValidationDomain: d,
			ValidationMethod: validationMethod,
			ValidationStatus: "PENDING_VALIDATION",
		}

		if validationMethod == "DNS" {
			dv.ResourceRecord = &acmstorelib.ResourceRecord{
				Name:  acmstorelib.GenerateDomainValidationRecordName(d),
				Type:  "CNAME",
				Value: acmstorelib.GenerateDomainValidationRecordValue(),
			}
		}

		options = append(options, dv)
	}

	return options
}

// generateKeyForKeyAlgorithm generates a private key matching the ACM
// KeyAlgorithm enum value. Returns a crypto.Signer that works with both
// RSA and ECDSA key types for x509 certificate creation.
func generateKeyForKeyAlgorithm(keyAlgorithm string) (crypto.Signer, error) {
	switch keyAlgorithm {
	case "", "RSA_2048":
		return vcrypto.GenerateRSAKey(2048)
	case "RSA_3072":
		return vcrypto.GenerateRSAKey(3072)
	case "RSA_4096":
		return vcrypto.GenerateRSAKey(4096)
	case "EC_prime256v1":
		return vcrypto.GenerateECDSAKey(elliptic.P256())
	case "EC_secp384r1":
		return vcrypto.GenerateECDSAKey(elliptic.P384())
	case "EC_secp521r1":
		return vcrypto.GenerateECDSAKey(elliptic.P521())
	default:
		return nil, fmt.Errorf("unsupported key algorithm: %s", keyAlgorithm)
	}
}

// signatureAlgorithmForKeyAlgorithm returns the ACM SignatureAlgorithm
// string that corresponds to the given KeyAlgorithm.
func signatureAlgorithmForKeyAlgorithm(keyAlgorithm string) string {
	switch keyAlgorithm {
	case "RSA_3072":
		return "SHA384WITHRSA"
	case "RSA_4096":
		return "SHA512WITHRSA"
	case "EC_prime256v1":
		return "SHA256WITHECDSA"
	case "EC_secp384r1":
		return "SHA384WITHECDSA"
	case "EC_secp521r1":
		return "SHA512WITHECDSA"
	default:
		return "SHA256WITHRSA"
	}
}

func extractDomainFromCert(cert string) (string, error) {
	parsed, err := vcrypto.ParseCertificatePEM([]byte(cert))
	if err != nil {
		return "", fmt.Errorf("failed to parse certificate: %w", err)
	}

	if len(parsed.DNSNames) > 0 {
		return parsed.DNSNames[0], nil
	}

	if parsed.Subject.CommonName != "" {
		return parsed.Subject.CommonName, nil
	}

	return "", fmt.Errorf("no domain name found in certificate")
}

// extractSerialFromPEM parses the certificate PEM and returns the X.509
// SerialNumber as a decimal string. Returns an empty string if parsing fails.
func extractSerialFromPEM(certPEM string) string {
	parsed, err := vcrypto.ParseCertificatePEM([]byte(certPEM))
	if err != nil || parsed == nil {
		return ""
	}
	return parsed.SerialNumber.String()
}

func determineKeyAlgorithm(cert string) string {
	parsed, err := vcrypto.ParseCertificatePEM([]byte(cert))
	if err != nil || parsed == nil {
		return "RSA_2048"
	}
	return determineKeyAlgorithmFromParsed(parsed)
}

func determineKeyAlgorithmFromParsed(cert *x509.Certificate) string {
	switch cert.PublicKeyAlgorithm {
	case x509.ECDSA:
		if key, ok := cert.PublicKey.(*ecdsa.PublicKey); ok {
			switch key.Curve {
			case elliptic.P384():
				return "EC_secp384r1"
			case elliptic.P521():
				return "EC_secp521r1"
			}
		}
		return "EC_prime256v1"
	case x509.RSA:
		bits := 0
		if key, ok := cert.PublicKey.(*rsa.PublicKey); ok {
			bits = key.N.BitLen()
		}
		switch bits {
		case 1024:
			return "RSA_1024"
		case 2048:
			return "RSA_2048"
		case 3072:
			return "RSA_3072"
		case 4096:
			return "RSA_4096"
		default:
			return "RSA_2048"
		}
	default:
		return "RSA_2048"
	}
}

func determineSignatureAlgorithmFromParsed(cert *x509.Certificate) string {
	switch cert.SignatureAlgorithm {
	case x509.SHA256WithRSA:
		return "SHA256WITHRSA"
	case x509.SHA384WithRSA:
		return "SHA384WITHRSA"
	case x509.SHA512WithRSA:
		return "SHA512WITHRSA"
	case x509.ECDSAWithSHA256:
		return "SHA256WITHECDSA"
	case x509.ECDSAWithSHA384:
		return "SHA384WITHECDSA"
	case x509.ECDSAWithSHA512:
		return "SHA512WITHECDSA"
	default:
		return "SHA256WITHRSA"
	}
}

func encryptPrivateKey(privateKeyPEM, passphrase string) (string, error) {
	keyBytes := []byte(privateKeyPEM)
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		if decoded, err := base64.StdEncoding.DecodeString(privateKeyPEM); err == nil {
			block, _ = pem.Decode(decoded)
		}
	}
	if block == nil {
		if restored := pemFixNewlines(privateKeyPEM); restored != privateKeyPEM {
			block, _ = pem.Decode([]byte(restored))
		}
	}
	if block == nil {
		return "", fmt.Errorf("failed to decode private key PEM")
	}

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	key := pbkdf2.Key([]byte(passphrase), salt, 100000, 32, sha256.New)
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return "", err
	}

	encrypted := gcm.Seal(nil, nonce, block.Bytes, nil)
	result := append(salt, nonce...)
	result = append(result, encrypted...)
	return base64.StdEncoding.EncodeToString(result), nil
}

func decodeBase64PEM(s string) (string, error) {
	if strings.Contains(s, "-----BEGIN") {
		return s, nil
	}
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", awserrors.NewValidationException(fmt.Sprintf("Failed to decode base64-encoded data: %v", err))
	}
	return strings.TrimSpace(string(decoded)), nil
}

func pemFixNewlines(pemStr string) string {
	pemStr = strings.TrimSpace(pemStr)
	if strings.Contains(pemStr, "\n") {
		return pemStr
	}
	begin := strings.Index(pemStr, "-----BEGIN")
	if begin < 0 {
		return pemStr
	}
	end := strings.Index(pemStr, "-----END")
	if end < 0 {
		return pemStr
	}
	end += len("-----END")
	typeEnd := strings.Index(pemStr[begin:], "-----")
	if typeEnd < 0 {
		return pemStr
	}
	typeEnd += begin
	var b strings.Builder
	b.WriteString(pemStr[begin : typeEnd+5])
	b.WriteByte('\n')
	b64 := pemStr[typeEnd+5 : end]
	for i := 0; i < len(b64); i += 64 {
		if i+64 > len(b64) {
			b.WriteString(b64[i:])
		} else {
			b.WriteString(b64[i : i+64])
			b.WriteByte('\n')
		}
	}
	b.WriteString(pemStr[end:])
	return b.String()
}
