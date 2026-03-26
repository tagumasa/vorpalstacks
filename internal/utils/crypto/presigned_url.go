package crypto

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

func buildPresignedCanonicalRequest(method, path string, query url.Values, bucket, region, contentType, signedHeaders string) string {
	canonicalQueryString := buildCanonicalQueryStringForPresigned(query)

	host := bucket + ".s3." + region + ".amazonaws.com"
	if strings.Contains(bucket, ".") {
		host = "s3." + region + ".amazonaws.com"
	}

	canonicalHeaders := "host:" + host + "\n"
	if contentType != "" && strings.Contains(signedHeaders, "content-type") {
		canonicalHeaders += "content-type:" + contentType + "\n"
	}

	payloadHash := "UNSIGNED-PAYLOAD"

	return method + "\n" +
		path + "\n" +
		canonicalQueryString + "\n" +
		canonicalHeaders + "\n" +
		signedHeaders + "\n" +
		payloadHash
}

func buildCanonicalQueryStringForPresigned(query url.Values) string {
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
		values := query[key]
		for j, value := range values {
			if j > 0 {
				result.WriteString("&")
			}
			result.WriteString(url.QueryEscape(key))
			result.WriteString("=")
			result.WriteString(url.QueryEscape(value))
		}
	}

	return result.String()
}

func buildPresignedStringToSign(amzDate, credentialScope, canonicalRequest string) string {
	canonicalRequestHash := SHA256HashString(canonicalRequest)

	return Algorithm + "\n" +
		amzDate + "\n" +
		credentialScope + "\n" +
		canonicalRequestHash
}

// PresignedURLVerifier verifies S3 presigned URL requests.
type PresignedURLVerifier struct {
	credentialsProvider CredentialsProvider
}

// NewPresignedURLVerifier creates a new presigned URL verifier.
//
// Parameters:
//   - provider: The credentials provider
//
// Returns:
//   - *PresignedURLVerifier: A new verifier instance
//
// Example:
//
//	verifier := NewPresignedURLVerifier(credentialsProvider)
func NewPresignedURLVerifier(provider CredentialsProvider) *PresignedURLVerifier {
	return &PresignedURLVerifier{
		credentialsProvider: provider,
	}
}

// PresignedURLParams contains the parsed parameters from a presigned URL.
type PresignedURLParams struct {
	// Algorithm is the signing algorithm (AWS4-HMAC-SHA256)
	Algorithm string
	// Credential contains the access key and credential scope
	Credential string
	// Date is the ISO 8601 format date from the URL
	Date string
	// Expires is the number of seconds the URL is valid for
	Expires int
	// SignedHeaders is the list of headers that were signed
	SignedHeaders string
	// Signature is the computed request signature
	Signature string
}

// ParsePresignedURL parses the query parameters from a presigned URL.
//
// Parameters:
//   - query: The URL query values
//
// Returns:
//   - *PresignedURLParams: The parsed parameters
//   - error: An error if parsing fails
//
// Example:
//
//	params, err := ParsePresignedURL(r.URL.Query())
func ParsePresignedURL(query url.Values) (*PresignedURLParams, error) {
	params := &PresignedURLParams{
		Algorithm:     query.Get("X-Amz-Algorithm"),
		Credential:    query.Get("X-Amz-Credential"),
		Date:          query.Get("X-Amz-Date"),
		SignedHeaders: query.Get("X-Amz-SignedHeaders"),
		Signature:     query.Get("X-Amz-Signature"),
	}

	expiresStr := query.Get("X-Amz-Expires")
	if expiresStr != "" {
		var err error
		params.Expires, err = parseInt(expiresStr)
		if err != nil {
			return nil, fmt.Errorf("invalid X-Amz-Expires: %w", err)
		}
	}

	if params.Algorithm == "" || params.Credential == "" || params.Date == "" ||
		params.SignedHeaders == "" || params.Signature == "" {
		return nil, fmt.Errorf("missing required presigned URL parameters")
	}

	return params, nil
}

// VerifyPresignedURL verifies an incoming request with a presigned URL.
// It validates the signature, expiration, and credentials.
//
// Parameters:
//   - r: The HTTP request
//   - bucket: The expected S3 bucket name
//   - region: The expected AWS region
//
// Returns:
//   - error: An error if verification fails, nil if successful
//
// Example:
//
//	err := verifier.VerifyPresignedURL(r, "my-bucket", "us-east-1")
func (v *PresignedURLVerifier) VerifyPresignedURL(r *http.Request, bucket, region string) error {
	query := r.URL.Query()

	params, err := ParsePresignedURL(query)
	if err != nil {
		return err
	}

	if params.Algorithm != Algorithm {
		return fmt.Errorf("unsupported algorithm: %s", params.Algorithm)
	}

	requestTime, err := time.Parse("20060102T150405Z", params.Date)
	if err != nil {
		return fmt.Errorf("invalid X-Amz-Date: %w", err)
	}

	now := time.Now().UTC()
	expiresAt := requestTime.Add(time.Duration(params.Expires) * time.Second)
	if now.After(expiresAt) {
		return fmt.Errorf("presigned URL has expired")
	}

	credentialParts := strings.Split(params.Credential, "/")
	if len(credentialParts) < 5 {
		return fmt.Errorf("invalid credential format")
	}
	accessKey := credentialParts[0]

	credentials, err := v.credentialsProvider.GetCredentials()
	if err != nil {
		return fmt.Errorf("credentials not found: %w", err)
	}

	if credentials.AccessKeyID != accessKey {
		return fmt.Errorf("access key mismatch")
	}

	canonicalRequest := v.buildCanonicalRequestForVerification(r, bucket, region, params.SignedHeaders)

	dateStr := credentialParts[1]
	credentialScope := strings.Join(credentialParts[1:], "/")
	stringToSign := buildPresignedStringToSign(params.Date, credentialScope, canonicalRequest)

	signingKey := DeriveSigningKey(credentials.SecretAccessKey, dateStr, region, "s3")
	calculatedSignature := HMACSHA256HexString(signingKey, stringToSign)

	if calculatedSignature != params.Signature {
		return fmt.Errorf("signature mismatch")
	}

	return nil
}

func (v *PresignedURLVerifier) buildCanonicalRequestForVerification(r *http.Request, bucket, region, signedHeaders string) string {
	method := r.Method

	path := r.URL.Path

	queryForSignature := url.Values{}
	for key, values := range r.URL.Query() {
		if !isPresignedParam(key) {
			queryForSignature[key] = values
		}
	}
	canonicalQueryString := buildCanonicalQueryStringForPresigned(queryForSignature)

	host := bucket + ".s3." + region + ".amazonaws.com"
	if strings.Contains(bucket, ".") {
		host = "s3." + region + ".amazonaws.com"
	}

	canonicalHeaders := "host:" + host + "\n"

	signedHeadersList := signedHeaders

	payloadHash := "UNSIGNED-PAYLOAD"

	return method + "\n" +
		path + "\n" +
		canonicalQueryString + "\n" +
		canonicalHeaders + "\n" +
		signedHeadersList + "\n" +
		payloadHash
}

func isPresignedParam(key string) bool {
	switch strings.ToLower(key) {
	case "x-amz-algorithm", "x-amz-credential", "x-amz-date",
		"x-amz-expires", "x-amz-signedheaders", "x-amz-signature":
		return true
	default:
		return false
	}
}

func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}
