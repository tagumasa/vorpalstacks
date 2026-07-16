package acm

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
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
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	acmstorelib "vorpalstacks/internal/store/aws/acm"
	"vorpalstacks/internal/utils/aws/types"

	vcrypto "vorpalstacks/internal/utils/crypto"

	"golang.org/x/crypto/pbkdf2"
)

func (s *ACMService) acmTagConfig(stores *acmStores, arn string) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.TagOperationConfig{
			ResourceParam:    "CertificateArn",
			TagsParam:        "Tags",
			TagKeysParam:     "Tags",
			TagKeyName:       "Key",
			TagValueName:     "Value",
			RequireTags:      true,
			RequireTagKeys:   true,
			RequireResource:  true,
			UseQueryFallback: true,
		},
		ResourceKey: func(rawKey string) string {
			return arn
		},
		ValidateResource: func(ctx context.Context, resourceKey string) error {
			_, err := stores.certificates.Get(arn)
			if err != nil {
				if acmstorelib.IsNotFound(err) {
					return awserrors.NewResourceNotFoundException("certificate", arn)
				}
				return err
			}
			return nil
		},
		TagFunc: func(ctx context.Context, resourceKey string, tagList []types.Tag) error {
			cert, err := stores.certificates.Get(resourceKey)
			if err != nil {
				return err
			}
			cert.Tags = tagutil.Apply(cert.Tags, tagList)
			return stores.certificates.Update(cert)
		},
		ParseTagKeys: func(params map[string]interface{}) []string {
			tagsToRemove := tagutil.ParseTagsWithQueryFallback(params, "Tags")
			keys := make([]string, 0, len(tagsToRemove))
			for _, t := range tagsToRemove {
				keys = append(keys, t.Key)
			}
			return keys
		},
		UntagFunc: func(ctx context.Context, resourceKey string, tagKeys []string) error {
			cert, err := stores.certificates.Get(resourceKey)
			if err != nil {
				return err
			}
			tagKeySet := make(map[string]bool, len(tagKeys))
			for _, k := range tagKeys {
				tagKeySet[k] = true
			}
			cert.Tags = tagutil.Remove(cert.Tags, tagKeySet)
			return stores.certificates.Update(cert)
		},
		ListFunc: func(ctx context.Context, resourceKey string) ([]types.Tag, error) {
			cert, err := stores.certificates.Get(resourceKey)
			if err != nil {
				return nil, err
			}
			return cert.Tags, nil
		},
		FormatResponse: func(tagList []types.Tag, _ string) (interface{}, error) {
			return map[string]interface{}{
				"Tags": tagutil.ToResponse(tagList),
			}, nil
		},
		EmptyResponse: func() (interface{}, error) {
			return response.EmptyResponse(), nil
		},
		MapError: func(err error) error {
			switch err.(type) {
			case *tagutil.MissingResourceError:
				return awserrors.NewResourceNotFoundException("certificate", "")
			case *tagutil.MissingTagsError:
				return awserrors.NewValidationException("Tags are required")
			case *tagutil.MissingTagKeysError:
				return awserrors.NewValidationException("Tag keys are required")
			}
			return err
		},
	}
}

// RequestCertificate requests a new certificate from ACM.
func (s *ACMService) RequestCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	domainName, err := parseDomainName(params)
	if err != nil {
		return nil, err
	}

	certId := acmstorelib.GenerateCertificateId()
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	certificateArn := stores.arnBuilder.BuildCertificateARN(certId)

	validationMethod := parseValidationMethod(params)
	keyAlgorithm := parseKeyAlgorithm(params)

	now := time.Now().UTC()
	key, err := vcrypto.GenerateRSAKey(2048)
	if err != nil {
		return nil, awserrors.NewAWSError("InternalErrorException", "Failed to generate key", 500)
	}

	serialBigInt, err := vcrypto.GenerateSerialNumber()
	if err != nil {
		return nil, awserrors.NewAWSError("InternalErrorException", "Failed to generate serial", 500)
	}
	notAfter := now.AddDate(1, 0, 0)
	template := &x509.Certificate{
		SerialNumber: serialBigInt,
		Subject:      pkix.Name{CommonName: domainName},
		NotBefore:    now,
		NotAfter:     notAfter,
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		DNSNames:     []string{domainName},
	}

	certDER, err := vcrypto.CreateCertificate(template, template, &key.PublicKey, key)
	if err != nil {
		return nil, awserrors.NewAWSError("InternalErrorException", "Failed to create certificate", 500)
	}

	certPEM := vcrypto.EncodeCertificatePEM(certDER)

	serialStr := acmstorelib.GenerateCertificateSerial()
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
		Subject:            domainName,
		Issuer:             domainName,
		AccountID:          reqCtx.GetAccountID(),
		Region:             reqCtx.GetRegion(),
		Certificate:        certPEM,
		NotBefore:          now,
		NotAfter:           notAfter,
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
			return nil, awserrors.NewAWSError("ResourceConflictException", "Certificate already exists", 409)
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
			return nil, awserrors.NewResourceNotFoundException("certificate", arn)
		}
		return nil, err
	}

	result := map[string]interface{}{
		"Certificate": cert.Certificate,
	}
	if cert.CertificateChain != "" {
		result["CertificateChain"] = cert.CertificateChain
	}
	return result, nil
}

// ListCertificates retrieves a list of certificates for the account.
func (s *ACMService) ListCertificates(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	marker := request.GetStringParam(params, "NextToken")
	maxItems := getMaxItems(params)

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

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

	var result *acmstorelib.CertificateListResult
	if len(statuses) > 0 {
		result, err = stores.certificates.ListByStatus(statuses, marker, maxItems)
	} else {
		result, err = stores.certificates.List(marker, maxItems)
	}
	if err != nil {
		return nil, err
	}

	return listResultToResponse(result), nil
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
			return nil, awserrors.NewResourceNotFoundException("certificate", arn)
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
			return nil, awserrors.NewResourceNotFoundException("certificate", arn)
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
			return nil, awserrors.NewResourceNotFoundException("certificate", arn)
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
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	arn, _ := parseCertificateArn(req.Parameters, "CertificateArn")
	return tagutil.HandleTag(ctx, req, s.acmTagConfig(stores, arn))
}

// RemoveTagsFromCertificate removes one or more tags from a certificate.
func (s *ACMService) RemoveTagsFromCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	arn, _ := parseCertificateArn(req.Parameters, "CertificateArn")
	return tagutil.HandleUntag(ctx, req, s.acmTagConfig(stores, arn))
}

// ListTagsForCertificate lists the tags associated with a certificate.
func (s *ACMService) ListTagsForCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	arn, _ := parseCertificateArn(req.Parameters, "CertificateArn")
	return tagutil.HandleList(ctx, req, s.acmTagConfig(stores, arn))
}

// ImportCertificate imports a certificate into ACM.
func (s *ACMService) ImportCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	certificate := request.GetStringParam(params, "Certificate")
	if certificate == "" {
		return nil, awserrors.NewValidationException("Certificate is required")
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

	certId := acmstorelib.GenerateCertificateId()
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	certificateArn := stores.arnBuilder.BuildCertificateARN(certId)

	now := time.Now().UTC()
	cert := &acmstorelib.Certificate{
		CertificateArn:     certificateArn,
		DomainName:         extractDomainFromCert(certificate),
		Serial:             acmstorelib.GenerateCertificateSerial(),
		Status:             "ISSUED",
		Type:               "IMPORTED",
		KeyAlgorithm:       determineKeyAlgorithm(certificate),
		SignatureAlgorithm: "SHA256WITHRSA",
		RenewalEligibility: "INELIGIBLE",
		CreatedAt:          now,
		ImportedAt:         now,
		NotBefore:          now,
		NotAfter:           now.AddDate(1, 0, 0),
		Certificate:        certificate,
		CertificateChain:   certificateChain,
		PrivateKey:         privateKey,
		Tags:               tags,
		AccountID:          reqCtx.GetAccountID(),
		Region:             reqCtx.GetRegion(),
	}

	if parsedCert, _ := vcrypto.ParseCertificatePEM([]byte(certificate)); parsedCert != nil {
		cert.NotBefore = parsedCert.NotBefore
		cert.NotAfter = parsedCert.NotAfter
		cert.KeyAlgorithm = determineKeyAlgorithmFromParsed(parsedCert)
		cert.SignatureAlgorithm = determineSignatureAlgorithmFromParsed(parsedCert)
		cert.Subject = parsedCert.Subject.String()
		cert.Issuer = parsedCert.Issuer.String()
	}

	if err := stores.certificates.Create(cert); err != nil {
		if acmstorelib.IsAlreadyExists(err) {
			return nil, awserrors.NewAWSError("ResourceConflictException", "Certificate already exists", 409)
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
			return nil, awserrors.NewResourceNotFoundException("certificate", arn)
		}
		return nil, err
	}

	optionsMap, ok := params["Options"].(map[string]interface{})
	if !ok {
		return nil, awserrors.NewValidationException("Options are required")
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
			return nil, awserrors.NewResourceNotFoundException("certificate", arn)
		}
		return nil, err
	}

	if cert.Type != "AMAZON_ISSUED" {
		return nil, awserrors.NewValidationException("Certificate type is not supported. Only Amazon-issued certificates can be renewed.")
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
			return nil, awserrors.NewResourceNotFoundException("certificate", arn)
		}
		return nil, err
	}

	if cert.PrivateKey == "" {
		return nil, awserrors.NewValidationException("Certificate does not have an exportable private key. Only imported certificates with private keys can be exported.")
	}

	passphrase := request.GetStringParam(params, "Passphrase")
	if passphrase == "" {
		return nil, awserrors.NewValidationException("Passphrase is required")
	}

	encryptedKey, err := encryptPrivateKey(cert.PrivateKey, passphrase)
	if err != nil {
		return nil, awserrors.NewValidationException("Failed to encrypt private key")
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
			return nil, awserrors.NewResourceNotFoundException("certificate", arn)
		}
		return nil, err
	}

	if cert.Status == "REVOKED" {
		return nil, awserrors.NewResourceNotFoundException("certificate", arn)
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
			Name:  acmstorelib.GenerateDomainValidationRecordName(domainName),
			Type:  "CNAME",
			Value: acmstorelib.GenerateDomainValidationRecordValue(),
		}
	}

	return []*acmstorelib.DomainValidation{dv}
}

func extractDomainFromCert(cert string) string {
	parsed, err := vcrypto.ParseCertificatePEM([]byte(cert))
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
	parsed, err := vcrypto.ParseCertificatePEM([]byte(cert))
	if err != nil {
		return "RSA_2048"
	}
	return determineKeyAlgorithmFromParsed(parsed)
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
