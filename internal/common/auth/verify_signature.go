// Package auth provides AWS authentication utilities for vorpalstacks.
package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	algorithm     = "AWS4-HMAC-SHA256"
	servicePrefix = "aws4_request"
)

// SignatureV4Verifier verifies AWS Signature Version 4 requests.
type SignatureV4Verifier struct {
	credentialsProvider CredentialsProvider
}

// NewSignatureV4Verifier creates a new Signature V4 verifier with the specified credentials provider.
func NewSignatureV4Verifier(provider CredentialsProvider) *SignatureV4Verifier {
	return &SignatureV4Verifier{
		credentialsProvider: provider,
	}
}

// VerifyRequest verifies an AWS Signature Version 4 request.
func (v *SignatureV4Verifier) VerifyRequest(r *http.Request, service, region string) error {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ErrMissingAuthHeader
	}

	if !strings.HasPrefix(authHeader, algorithm+" ") {
		return ErrInvalidAlgorithm
	}

	authHeader = strings.TrimPrefix(authHeader, algorithm+" ")

	var credential, signedHeaders, signature string
	parts := strings.Split(authHeader, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "Credential=") {
			credential = strings.TrimPrefix(part, "Credential=")
		} else if strings.HasPrefix(part, "SignedHeaders=") {
			signedHeaders = strings.TrimPrefix(part, "SignedHeaders=")
		} else if strings.HasPrefix(part, "Signature=") {
			signature = strings.TrimPrefix(part, "Signature=")
		}
	}

	if credential == "" || signedHeaders == "" || signature == "" {
		return ErrMissingAuthFields
	}

	amzDate := r.Header.Get("X-Amz-Date")
	if amzDate == "" {
		return ErrMissingDateHeader
	}

	var bodyBytes []byte
	if r.Body != nil {
		var readErr error
		bodyBytes, readErr = io.ReadAll(r.Body)
		if readErr != nil {
			return ErrReadBodyFailed
		}
		r.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
	}

	canonicalRequest, err := v.buildCanonicalRequest(r, signedHeaders, bodyBytes)
	if err != nil {
		return err
	}

	credentialScope, stringToSign, err := v.buildStringToSign(amzDate, credential, canonicalRequest)
	if err != nil {
		return err
	}

	credentials, err := v.credentialsProvider.GetCredentials()
	if err != nil {
		return ErrGetCredentialsFailed
	}

	calculatedSignature, err := v.calculateSignature(amzDate, credentialScope.Region, credentialScope.Service, stringToSign, credentials.SecretAccessKey)
	if err != nil {
		return ErrCalculateSignatureFailed
	}

	if subtle.ConstantTimeCompare([]byte(calculatedSignature), []byte(signature)) != 1 {
		return ErrSignatureMismatch
	}

	return nil
}

// CredentialScope represents the scope of AWS credentials.
type CredentialScope struct {
	Date    string
	Region  string
	Service string
}

func (v *SignatureV4Verifier) buildCanonicalRequest(r *http.Request, signedHeaders string, body []byte) (string, error) {
	method := r.Method

	uri := r.URL.EscapedPath()
	if uri == "" {
		uri = "/"
	}

	query := v.buildCanonicalQueryString(r.URL.Query())

	host := r.Host
	if host == "" {
		host = r.Header.Get("Host")
	}

	headers := make(map[string]string)
	for _, header := range strings.Split(signedHeaders, ";") {
		header = strings.TrimSpace(header)
		if header != "" {
			lowerHeader := strings.ToLower(header)
			if lowerHeader == "host" {
				headers[lowerHeader] = host
			} else {
				value := r.Header.Get(header)
				headers[lowerHeader] = strings.TrimSpace(value)
			}
		}
	}

	canonicalHeaders := ""
	sortedHeaders := make([]string, 0, len(headers))
	for header := range headers {
		sortedHeaders = append(sortedHeaders, header)
	}
	sort.Strings(sortedHeaders)

	for _, header := range sortedHeaders {
		canonicalHeaders += header + ":" + headers[header] + "\n"
	}

	signedHeadersList := signedHeaders
	payloadHash := v.sha256Hash(string(body))

	canonicalRequest := method + "\n" + uri + "\n" + query + "\n" + canonicalHeaders + "\n" + signedHeadersList + "\n" + payloadHash

	return canonicalRequest, nil
}

func (v *SignatureV4Verifier) buildCanonicalQueryString(query url.Values) string {
	if len(query) == 0 {
		return ""
	}

	keys := make([]string, 0, len(query))
	for key := range query {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var result strings.Builder
	for i, key := range keys {
		if i > 0 {
			result.WriteString("&")
		}
		result.WriteString(url.QueryEscape(key))
		result.WriteString("=")
		values := make([]string, len(query[key]))
		copy(values, query[key])
		sort.Strings(values)
		for j, value := range values {
			if j > 0 {
				result.WriteString("&")
				result.WriteString(url.QueryEscape(key))
				result.WriteString("=")
			}
			result.WriteString(url.QueryEscape(value))
		}
	}

	return result.String()
}

func (v *SignatureV4Verifier) buildStringToSign(amzDate, credential, canonicalRequest string) (*CredentialScope, string, error) {
	parts := strings.Split(credential, "/")
	if len(parts) < 5 {
		return nil, "", ErrInvalidCredentialFormat
	}
	scope := &CredentialScope{
		Date:    parts[1],
		Region:  parts[2],
		Service: parts[3],
	}

	canonicalRequestHash := v.sha256Hash(canonicalRequest)
	stringToSign := algorithm + "\n" + amzDate + "\n" + scope.Date + "/" + scope.Region + "/" + scope.Service + "/" + servicePrefix + "\n" + canonicalRequestHash

	return scope, stringToSign, nil
}

func (v *SignatureV4Verifier) calculateSignature(amzDate, region, service, stringToSign, secretKey string) (string, error) {
	date, err := time.Parse("20060102T150405Z", amzDate)
	if err != nil {
		return "", ErrParseDateFailed
	}
	dateStr := date.Format("20060102")

	kDate := v.hmacSHA256([]byte("AWS4"+secretKey), dateStr)
	kRegion := v.hmacSHA256(kDate, region)
	kService := v.hmacSHA256(kRegion, service)
	kSigning := v.hmacSHA256(kService, servicePrefix)

	signature := v.hmacSHA256(kSigning, stringToSign)

	return hex.EncodeToString(signature), nil
}

func (v *SignatureV4Verifier) sha256Hash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func (v *SignatureV4Verifier) hmacSHA256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}
