package acm

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	acmstorelib "vorpalstacks/internal/store/aws/acm"

	"golang.org/x/crypto/pbkdf2"
)

func generateCertificateId() string {
	return acmstorelib.GenerateCertificateId()
}

func generateCertificateSerial() string {
	return acmstorelib.GenerateCertificateSerial()
}

func generateDomainValidationRecordName(domain string) string {
	return acmstorelib.GenerateDomainValidationRecordName(domain)
}

func generateDomainValidationRecordValue() string {
	return acmstorelib.GenerateDomainValidationRecordValue()
}

// RequestCertificate requests a new certificate from ACM.
func (s *ACMService) RequestCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	domainName, err := parseDomainName(params)
	if err != nil {
		return nil, err
	}

	certId := generateCertificateId()
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	certificateArn := stores.arnBuilder.BuildCertificateARN(certId)

	validationMethod := parseValidationMethod(params)
	keyAlgorithm := parseKeyAlgorithm(params)

	now := time.Now().UTC()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, NewAPIError("InternalErrorException", "Failed to generate key", 500)
	}

	serialBigInt, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, NewAPIError("InternalErrorException", "Failed to generate serial", 500)
	}
	template := &x509.Certificate{
		SerialNumber: serialBigInt,
		NotBefore:    now,
		NotAfter:     now.AddDate(1, 0, 0),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		DNSNames:     []string{domainName},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return nil, NewAPIError("InternalErrorException", "Failed to create certificate", 500)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	serialStr := generateCertificateSerial()
	cert := &acmstorelib.Certificate{
		CertificateArn:     certificateArn,
		DomainName:         domainName,
		Serial:             serialStr,
		Status:             "ISSUED",
		Type:               "AMAZON_ISSUED",
		KeyAlgorithm:       keyAlgorithm,
		SignatureAlgorithm: "SHA256WITHRSA",
		RenewalEligibility: "ELIGIBLE",
		CreatedAt:          now,
		AccountID:          reqCtx.GetAccountID(),
		Region:             reqCtx.GetRegion(),
		Certificate:        string(certPEM),
		NotBefore:          now,
		NotAfter:           now.AddDate(1, 0, 0),
		IssuedAt:           now,
	}

	if sansRaw := params["SubjectAlternativeNames"]; sansRaw != nil {
		if sans, ok := sansRaw.([]interface{}); ok {
			for _, san := range sans {
				if s, ok := san.(string); ok {
					cert.SubjectAlternativeNames = append(cert.SubjectAlternativeNames, s)
				}
			}
		}
	}

	tags := tagutil.ParseTagsWithQueryFallback(params, "Tags")
	cert.Tags = tags

	domainValidationOptions := buildDomainValidationOptions(domainName, validationMethod)
	for i := range domainValidationOptions {
		domainValidationOptions[i].ValidationStatus = "SUCCESS"
	}
	cert.DomainValidationOptions = domainValidationOptions

	if optionsMap, ok := params["Options"].(map[string]interface{}); ok {
		cert.Options = &acmstorelib.CertificateOptions{
			CertificateTransparencyLoggingPreference: parseCertificateTransparencyLoggingPreference(optionsMap),
			Export:                                   "DISABLED",
		}
	} else {
		cert.Options = &acmstorelib.CertificateOptions{
			CertificateTransparencyLoggingPreference: "ENABLED",
			Export:                                   "DISABLED",
		}
	}

	if err := stores.certificates.Create(cert); err != nil {
		if acmstorelib.IsAlreadyExists(err) {
			return nil, NewAPIError("ResourceConflictException", "Certificate already exists", 409)
		}
		return nil, err
	}

	return map[string]interface{}{
		"CertificateArn": certificateArn,
	}, nil
}

// GetCertificate retrieves the certificate and certificate chain for the specified ARN.
func (s *ACMService) GetCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	arn, err := parseCertificateArn(params, "CertificateArn")
	if err != nil {
		return nil, err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cert, err := stores.certificates.Get(arn)
	if err != nil {
		if acmstorelib.IsNotFound(err) {
			return nil, NewResourceNotFoundError("certificate", arn)
		}
		return nil, err
	}

	return map[string]interface{}{
		"Certificate":      cert.Certificate,
		"CertificateChain": cert.CertificateChain,
	}, nil
}

// ListCertificates retrieves a list of certificates for the account.
func (s *ACMService) ListCertificates(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	marker := request.GetStringParam(params, "NextToken")
	maxItems := getMaxItems(params)

	var statuses []string
	if raw, ok := params["CertificateStatuses"]; ok {
		if arr, ok := raw.([]interface{}); ok {
			for _, v := range arr {
				if s, ok := v.(string); ok {
					statuses = append(statuses, s)
				}
			}
		}
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if len(statuses) > 0 {
		allCerts, err := stores.certificates.ListAll()
		if err != nil {
			return nil, err
		}
		statusSet := make(map[string]bool)
		for _, s := range statuses {
			statusSet[s] = true
		}
		var filtered []*acmstorelib.CertificateSummary
		for _, cert := range allCerts {
			if statusSet[cert.Status] {
				filtered = append(filtered, acmstorelib.CertificateToSummary(cert))
			}
		}
		return filteredListToResponse(filtered, marker, maxItems), nil
	}

	result, err := stores.certificates.List(marker, maxItems)
	if err != nil {
		return nil, err
	}

	return listCertificatesToResponse(result), nil
}

// DeleteCertificate deletes the specified certificate.
func (s *ACMService) DeleteCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	arn, err := parseCertificateArn(params, "CertificateArn")
	if err != nil {
		return nil, err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cert, err := stores.certificates.Get(arn)
	if err != nil {
		if acmstorelib.IsNotFound(err) {
			return nil, NewResourceNotFoundError("certificate", arn)
		}
		return nil, err
	}

	if len(cert.InUseBy) > 0 {
		return nil, NewResourceInUseError("certificate", arn)
	}

	if err := stores.certificates.Delete(arn); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DescribeCertificate retrieves detailed information about a certificate.
func (s *ACMService) DescribeCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	arn, err := parseCertificateArn(params, "CertificateArn")
	if err != nil {
		return nil, err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cert, err := stores.certificates.Get(arn)
	if err != nil {
		if acmstorelib.IsNotFound(err) {
			return nil, NewResourceNotFoundError("certificate", arn)
		}
		return nil, err
	}

	return map[string]interface{}{
		"Certificate": certificateToDetailResponse(cert),
	}, nil
}

// ResendValidationEmail resends the domain validation email for a certificate.
func (s *ACMService) ResendValidationEmail(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	arn, err := parseCertificateArn(params, "CertificateArn")
	if err != nil {
		return nil, err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cert, err := stores.certificates.Get(arn)
	if err != nil {
		if acmstorelib.IsNotFound(err) {
			return nil, NewResourceNotFoundError("certificate", arn)
		}
		return nil, err
	}

	if cert.Type == "IMPORTED" {
		return nil, NewInvalidStateException("Certificate is not in PENDING_VALIDATION state")
	}
	if cert.Status != "PENDING_VALIDATION" && cert.Status != "ISSUED" {
		return nil, NewInvalidStateException("Certificate is not in PENDING_VALIDATION state")
	}

	return response.EmptyResponse(), nil
}

// AddTagsToCertificate adds one or more tags to a certificate.
func (s *ACMService) AddTagsToCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	arn, err := parseCertificateArn(params, "CertificateArn")
	if err != nil {
		return nil, err
	}

	newTags := tagutil.ParseTagsWithQueryFallback(params, "Tags")
	if len(newTags) == 0 {
		return nil, NewValidationException("Tags are required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cert, err := stores.certificates.Get(arn)
	if err != nil {
		if acmstorelib.IsNotFound(err) {
			return nil, NewResourceNotFoundError("certificate", arn)
		}
		return nil, err
	}

	cert.Tags = tagutil.Apply(cert.Tags, newTags)
	if err := stores.certificates.Update(cert); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// RemoveTagsFromCertificate removes one or more tags from a certificate.
func (s *ACMService) RemoveTagsFromCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	arn, err := parseCertificateArn(params, "CertificateArn")
	if err != nil {
		return nil, err
	}

	tagsToRemove := tagutil.ParseTagsWithQueryFallback(params, "Tags")
	if len(tagsToRemove) == 0 {
		return nil, NewValidationException("Tags are required")
	}

	tagKeys := make(map[string]bool)
	for _, t := range tagsToRemove {
		tagKeys[t.Key] = true
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cert, err := stores.certificates.Get(arn)
	if err != nil {
		if acmstorelib.IsNotFound(err) {
			return nil, NewResourceNotFoundError("certificate", arn)
		}
		return nil, err
	}

	cert.Tags = tagutil.Remove(cert.Tags, tagKeys)

	if err := stores.certificates.Update(cert); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListTagsForCertificate lists the tags associated with a certificate.
func (s *ACMService) ListTagsForCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	arn, err := parseCertificateArn(params, "CertificateArn")
	if err != nil {
		return nil, err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cert, err := stores.certificates.Get(arn)
	if err != nil {
		if acmstorelib.IsNotFound(err) {
			return nil, NewResourceNotFoundError("certificate", arn)
		}
		return nil, err
	}

	return map[string]interface{}{
		"Tags": tagutil.ToResponse(cert.Tags),
	}, nil
}

// ImportCertificate imports a certificate into ACM.
func (s *ACMService) ImportCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	certificate := request.GetStringParam(params, "Certificate")
	if certificate == "" {
		return nil, NewValidationException("Certificate is required")
	}
	certificate = decodeBase64PEM(certificate)

	privateKey := request.GetStringParam(params, "PrivateKey")
	if privateKey != "" {
		privateKey = decodeBase64PEM(privateKey)
	}
	certificateChain := request.GetStringParam(params, "CertificateChain")
	if certificateChain != "" {
		certificateChain = decodeBase64PEM(certificateChain)
	}

	tags := tagutil.ParseTagsWithQueryFallback(params, "Tags")

	certId := generateCertificateId()
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	certificateArn := stores.arnBuilder.BuildCertificateARN(certId)

	cert := &acmstorelib.Certificate{
		CertificateArn:     certificateArn,
		DomainName:         extractDomainFromCert(certificate),
		Serial:             generateCertificateSerial(),
		Status:             "ISSUED",
		Type:               "IMPORTED",
		KeyAlgorithm:       determineKeyAlgorithm(certificate),
		SignatureAlgorithm: "SHA256WITHRSA",
		RenewalEligibility: "INELIGIBLE",
		CreatedAt:          time.Now(),
		ImportedAt:         time.Now(),
		NotBefore:          time.Now(),
		NotAfter:           time.Now().AddDate(1, 0, 0),
		Certificate:        certificate,
		CertificateChain:   certificateChain,
		PrivateKey:         privateKey,
		Tags:               tags,
		AccountID:          reqCtx.GetAccountID(),
		Region:             reqCtx.GetRegion(),
	}

	if parsedCert := parseCertificatePEM(certificate); parsedCert != nil {
		cert.NotBefore = parsedCert.NotBefore
		cert.NotAfter = parsedCert.NotAfter
		cert.KeyAlgorithm = determineKeyAlgorithmFromParsed(parsedCert)
		cert.SignatureAlgorithm = determineSignatureAlgorithmFromParsed(parsedCert)
	}

	if err := stores.certificates.Create(cert); err != nil {
		if acmstorelib.IsAlreadyExists(err) {
			return nil, NewAPIError("ResourceConflictException", "Certificate already exists", 409)
		}
		return nil, err
	}

	return map[string]interface{}{
		"CertificateArn": certificateArn,
	}, nil
}

// UpdateCertificateOptions updates the certificate options.
func (s *ACMService) UpdateCertificateOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	arn, err := parseCertificateArn(params, "CertificateArn")
	if err != nil {
		return nil, err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cert, err := stores.certificates.Get(arn)
	if err != nil {
		if acmstorelib.IsNotFound(err) {
			return nil, NewResourceNotFoundError("certificate", arn)
		}
		return nil, err
	}

	optionsMap, ok := params["Options"].(map[string]interface{})
	if !ok {
		return nil, NewValidationException("Options are required")
	}

	cert.Options = &acmstorelib.CertificateOptions{
		CertificateTransparencyLoggingPreference: parseCertificateTransparencyLoggingPreference(optionsMap),
		Export:                                   "DISABLED",
	}

	if err := stores.certificates.Update(cert); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// RenewCertificate renews an ACM certificate.
func (s *ACMService) RenewCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn, err := parseCertificateArn(req.Parameters, "CertificateArn")
	if err != nil {
		return nil, err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cert, err := stores.certificates.Get(arn)
	if err != nil {
		if acmstorelib.IsNotFound(err) {
			return nil, NewResourceNotFoundError("certificate", arn)
		}
		return nil, err
	}

	if cert.Type != "AMAZON_ISSUED" {
		return nil, NewValidationException("Certificate type is not supported. Only Amazon-issued certificates can be renewed.")
	}

	if cert.Status != "ISSUED" && cert.Status != "EXPIRED" {
		return nil, NewInvalidStateException("Certificate is not in a valid state for renewal.")
	}

	if cert.RenewalEligibility == "INELIGIBLE" {
		return nil, NewInvalidStateException("Certificate is not eligible for renewal.")
	}

	now := time.Now().UTC()
	cert.NotBefore = now
	cert.NotAfter = now.AddDate(1, 0, 0)
	cert.RenewalSummary = &acmstorelib.RenewalSummary{
		RenewalStatus: "PENDING",
		UpdatedAt:     now,
	}

	if err := stores.certificates.Update(cert); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ExportCertificate exports a private key and certificate.
func (s *ACMService) ExportCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	arn, err := parseCertificateArn(params, "CertificateArn")
	if err != nil {
		return nil, err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cert, err := stores.certificates.Get(arn)
	if err != nil {
		if acmstorelib.IsNotFound(err) {
			return nil, NewResourceNotFoundError("certificate", arn)
		}
		return nil, err
	}

	if cert.PrivateKey == "" {
		return nil, NewValidationException("Certificate does not have an exportable private key. Only imported certificates with private keys can be exported.")
	}

	passphrase := request.GetStringParam(params, "Passphrase")
	if passphrase == "" {
		return nil, NewValidationException("Passphrase is required")
	}

	encryptedKey, err := encryptPrivateKey(cert.PrivateKey, passphrase)
	if err != nil {
		return nil, NewValidationException("Failed to encrypt private key")
	}

	return map[string]interface{}{
		"Certificate":      cert.Certificate,
		"CertificateChain": cert.CertificateChain,
		"PrivateKey":       encryptedKey,
	}, nil
}

// RevokeCertificate revokes an ACM certificate.
func (s *ACMService) RevokeCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn, err := parseCertificateArn(req.Parameters, "CertificateArn")
	if err != nil {
		return nil, err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cert, err := stores.certificates.Get(arn)
	if err != nil {
		if acmstorelib.IsNotFound(err) {
			return nil, NewResourceNotFoundError("certificate", arn)
		}
		return nil, err
	}

	if cert.Status == "REVOKED" {
		return nil, NewResourceNotFoundError("certificate", arn)
	}

	if cert.Status != "ISSUED" {
		return nil, NewInvalidStateException("Certificate is not in a valid state for revocation.")
	}

	if reasonRaw, ok := req.Parameters["RevocationReason"]; ok {
		if reason, ok := reasonRaw.(string); ok {
			cert.RevocationReason = reason
		}
	}
	cert.Status = "REVOKED"
	cert.RevokedAt = time.Now().UTC()

	if err := stores.certificates.Update(cert); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

func buildDomainValidationOptions(domainName, validationMethod string) []*acmstorelib.DomainValidation {
	dv := &acmstorelib.DomainValidation{
		DomainName:       domainName,
		ValidationDomain: domainName,
		ValidationMethod: validationMethod,
		ValidationStatus: "PENDING",
	}

	if validationMethod == "DNS" {
		dv.ResourceRecord = &acmstorelib.ResourceRecord{
			Name:  generateDomainValidationRecordName(domainName),
			Type:  "CNAME",
			Value: generateDomainValidationRecordValue(),
		}
	}

	return []*acmstorelib.DomainValidation{dv}
}

func extractDomainFromCert(cert string) string {
	block, _ := pem.Decode([]byte(cert))
	if block == nil {
		return "unknown"
	}

	parsed, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "unknown"
	}

	if len(parsed.DNSNames) > 0 {
		return parsed.DNSNames[0]
	}

	if parsed.Subject.CommonName != "" {
		return parsed.Subject.CommonName
	}

	return "unknown"
}

func determineKeyAlgorithm(cert string) string {
	parsed := parseCertificatePEM(cert)
	if parsed == nil {
		return "RSA_2048"
	}
	return determineKeyAlgorithmFromParsed(parsed)
}

func parseCertificatePEM(cert string) *x509.Certificate {
	block, _ := pem.Decode([]byte(cert))
	if block == nil {
		return nil
	}
	parsed, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil
	}
	return parsed
}

func determineKeyAlgorithmFromParsed(cert *x509.Certificate) string {
	switch cert.PublicKeyAlgorithm {
	case x509.ECDSA:
		return "EC_prime256v1"
	case x509.RSA:
		bits := 0
		if key, ok := cert.PublicKey.(*rsa.PublicKey); ok {
			bits = key.N.BitLen()
		}
		switch bits {
		case 0:
			return "RSA_2048"
		case 2048:
			return "RSA_2048"
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
	case x509.SHA256WithRSA, x509.SHA384WithRSA, x509.SHA512WithRSA:
		return "SHA256WITHRSA"
	case x509.ECDSAWithSHA256:
		return "ECDSA_SHA_256"
	case x509.ECDSAWithSHA384:
		return "ECDSA_SHA_384"
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

func decodeBase64PEM(s string) string {
	if strings.Contains(s, "-----BEGIN") {
		return s
	}
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return s
	}
	return strings.TrimSpace(string(decoded))
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
