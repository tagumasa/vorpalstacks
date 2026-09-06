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
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	acmstorelib "vorpalstacks/internal/store/aws/acm"

	vcrypto "vorpalstacks/internal/utils/crypto"

	"golang.org/x/crypto/pbkdf2"
)

// errUnsupportedKeyAlgorithm marks a key-algorithm string the issuance
// pipeline cannot generate a key pair for; issuance maps it to
// InvalidParameterException while renewal maps every pipeline failure to an
// internal error.
var errUnsupportedKeyAlgorithm = errors.New("unsupported key algorithm")

// issuedCertMaterial is the freshly generated X.509 material produced by
// issueCertificateMaterial.
type issuedCertMaterial struct {
	certPEM   string
	keyPEM    string
	serial    string
	notBefore time.Time
	notAfter  time.Time
	// keyUsages carries the KeyUsageName strings granted by the issuance
	// template so callers persist them alongside the PEM.
	keyUsages []string
}

// issueCertificateMaterial runs the issuance pipeline shared by certificate
// issuance and renewal: a new key pair, a new serial, and a one-year
// self-signed certificate for domainName plus its SANs, returned as PEM
// blobs. Callers stamp their own certificate fields onto the result.
func issueCertificateMaterial(domainName string, sans []string, keyAlgorithm string) (*issuedCertMaterial, error) {
	key, err := generateKeyForKeyAlgorithm(keyAlgorithm)
	if err != nil {
		return nil, fmt.Errorf("%w %s: %v", errUnsupportedKeyAlgorithm, keyAlgorithm, err)
	}

	serialBigInt, err := vcrypto.GenerateSerialNumber()
	if err != nil {
		return nil, fmt.Errorf("generate serial: %w", err)
	}

	now := time.Now().UTC()
	notAfter := now.AddDate(1, 0, 0)

	dnsNames := make([]string, 0, 1+len(sans))
	dnsNames = append(dnsNames, domainName)
	dnsNames = append(dnsNames, sans...)

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
		return nil, fmt.Errorf("create certificate: %w", err)
	}

	keyPEM, err := vcrypto.EncodePrivateKeyPEM(key)
	if err != nil {
		return nil, fmt.Errorf("encode private key: %w", err)
	}

	return &issuedCertMaterial{
		certPEM:   vcrypto.EncodeCertificatePEM(certDER),
		keyPEM:    keyPEM,
		serial:    formatSerialNumberHex(serialBigInt),
		notBefore: now,
		notAfter:  notAfter,
		keyUsages: keyUsageFlagsToStrings(template.KeyUsage),
	}, nil
}

// formatSerialNumberHex renders an X.509 serial number in the ACM wire
// form: colon-separated lowercase hex byte pairs, the single representation
// the Smithy SerialNumber shape defines (pattern
// ^[0-9a-f]{2}(:[0-9a-f]{2}){1,19}$ — at least two byte pairs). Storage,
// the DescribeCertificate Serial field, the X509Attributes.SerialNumber
// output, and the search filter all carry this form, so a filter value
// compares equal to what the response emitted. Serials shorter than two
// bytes are zero-padded so every rendered value satisfies the pattern.
func formatSerialNumberHex(n *big.Int) string {
	serialBytes := n.Bytes()
	for len(serialBytes) < 2 {
		serialBytes = append([]byte{0}, serialBytes...)
	}
	pairs := make([]string, len(serialBytes))
	for i, b := range serialBytes {
		pairs[i] = fmt.Sprintf("%02x", b)
	}
	return strings.Join(pairs, ":")
}

// renewCertificateMaterial regenerates the X.509 certificate material (new
// key pair, new serial, new PEM blob) for an existing AMAZON_ISSUED
// certificate during renewal. The ARN is preserved per AWS specification.
// The caller is responsible for setting RenewalSummary and persisting via
// stores.certificates.Update afterwards.
func renewCertificateMaterial(cert *acmstorelib.Certificate) error {
	material, err := issueCertificateMaterial(cert.DomainName, cert.SubjectAlternativeNames, cert.KeyAlgorithm)
	if err != nil {
		return err
	}

	// The renewed key pair replaces the stored one so TLS termination on
	// cross-service listeners keeps serving the renewed certificate. The
	// renewal template grants no extended usages, so the reissued material
	// determines both usage fields.
	cert.Serial = material.serial
	cert.Certificate = material.certPEM
	cert.PrivateKey = material.keyPEM
	cert.Status = "ISSUED"
	cert.NotBefore = material.notBefore
	cert.NotAfter = material.notAfter
	cert.IssuedAt = material.notBefore
	cert.SignatureAlgorithm = signatureAlgorithmForKeyAlgorithm(cert.KeyAlgorithm)
	cert.KeyUsages = keyUsageList(material.keyUsages)
	cert.ExtendedKeyUsages = nil
	return nil
}

// generateAmazonIssuedCert generates a self-signed AMAZON_ISSUED certificate
// with the given parameters. Shared by the HTTP API (RequestCertificate) and
// the admin handler. The caller is responsible for setting AccountID, Region,
// Tags, ManagedBy, CertificateAuthorityArn, and Options afterwards.
func generateAmazonIssuedCert(certArn, domainName string, sans []string, keyAlgorithm, validationMethod string) (*acmstorelib.Certificate, error) {
	material, err := issueCertificateMaterial(domainName, sans, keyAlgorithm)
	if err != nil {
		if errors.Is(err, errUnsupportedKeyAlgorithm) {
			return nil, NewInvalidParameterError(fmt.Sprintf("Unsupported KeyAlgorithm: %s", keyAlgorithm))
		}
		return nil, NewInternalServerException(fmt.Sprintf("Failed to generate certificate material: %v", err))
	}

	// Build domain validation options. ACM creates certs in PENDING_VALIDATION
	// state; the platform validates synchronously (self-signed) so domain
	// validation options transition to SUCCESS immediately.
	domainValidationOptions := buildDomainValidationOptions(domainName, validationMethod, sans)
	for i := range domainValidationOptions {
		domainValidationOptions[i].ValidationStatus = "SUCCESS"
	}

	// The issuing key pair is persisted with the certificate: TLS
	// termination on cross-service listeners (CloudFront, API Gateway)
	// needs it to serve the certificate, mirroring how AWS deploys
	// ACM-held key pairs to its edges. ExportCertificate keeps the
	// material unavailable to callers through the certificate-type guard.
	return &acmstorelib.Certificate{
		CertificateArn:           certArn,
		DomainName:               domainName,
		Serial:                   material.serial,
		Status:                   "ISSUED",
		Type:                     "AMAZON_ISSUED",
		KeyAlgorithm:             keyAlgorithm,
		SignatureAlgorithm:       signatureAlgorithmForKeyAlgorithm(keyAlgorithm),
		RenewalEligibility:       "ELIGIBLE",
		CreatedAt:                material.notBefore,
		Subject:                  domainName,
		Issuer:                   domainName,
		Certificate:              material.certPEM,
		PrivateKey:               material.keyPEM,
		NotBefore:                material.notBefore,
		NotAfter:                 material.notAfter,
		IssuedAt:                 material.notBefore,
		KeyUsages:                keyUsageList(material.keyUsages),
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

// domainFromParsedCert selects the certificate's primary domain from an
// already-parsed certificate: the first DNS SAN, falling back to the
// subject CommonName.
func domainFromParsedCert(parsed *x509.Certificate) (string, error) {
	if len(parsed.DNSNames) > 0 {
		return parsed.DNSNames[0], nil
	}

	if parsed.Subject.CommonName != "" {
		return parsed.Subject.CommonName, nil
	}

	return "", fmt.Errorf("no domain name found in certificate")
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
