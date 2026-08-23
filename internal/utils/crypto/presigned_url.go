package crypto

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/common/auth"
)

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
			result.WriteString(rfc3986Encode(key))
			result.WriteString("=")
			result.WriteString(rfc3986Encode(value))
		}
	}

	return result.String()
}

// rfc3986Encode percent-encodes a string per RFC 3986, matching AWS SigV4
// canonical query string requirements. Differs from url.QueryEscape which
// encodes spaces as "+" and "~" as "%7E".
func rfc3986Encode(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') ||
			c == '-' || c == '_' || c == '.' || c == '~' {
			b.WriteByte(c)
		} else {
			fmt.Fprintf(&b, "%%%02X", c)
		}
	}
	return b.String()
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
	credentialsProvider auth.CredentialsProvider
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
func NewPresignedURLVerifier(provider auth.CredentialsProvider) *PresignedURLVerifier {
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

	if subtle.ConstantTimeCompare([]byte(calculatedSignature), []byte(params.Signature)) != 1 {
		return fmt.Errorf("signature mismatch")
	}

	return nil
}

func (v *PresignedURLVerifier) buildCanonicalRequestForVerification(r *http.Request, bucket, region, signedHeaders string) string {
	method := r.Method

	path := r.URL.EscapedPath()
	if path == "" {
		path = "/"
	}

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

// CheckPresignedURLFreshness validates only the freshness of a presigned
// URL's date and expiry parameters. It is used when the signature itself
// cannot be verified (no credentials provider wired): an expired URL is
// still rejected, which keeps the documented validity window enforced
// even without full SigV4 verification.
func CheckPresignedURLFreshness(query url.Values) error {
	params, err := ParsePresignedURL(query)
	if err != nil {
		return err
	}
	requestTime, err := time.Parse("20060102T150405Z", params.Date)
	if err != nil {
		return fmt.Errorf("invalid X-Amz-Date: %w", err)
	}
	if time.Now().UTC().After(requestTime.Add(time.Duration(params.Expires) * time.Second)) {
		return fmt.Errorf("presigned URL has expired")
	}
	return nil
}

// PresignS3URL builds a SigV4 presigned URL for an S3 object request. It
// is the counterpart of PresignedURLVerifier: the canonical request is
// constructed exactly as the verifier reconstructs it (virtual-host style
// canonical host, "host" as the only signed header, UNSIGNED-PAYLOAD), so
// a URL produced here passes VerifyPresignedURL until it expires. The
// endpoint host in the returned URL is the caller's own — typically the
// platform's S3 plane rather than the AWS DNS name used for signing.
func PresignS3URL(method, scheme, endpointHost, bucket, key, region string, expires time.Duration, accessKeyID, secretAccessKey string) string {
	now := time.Now().UTC()
	amzDate := now.Format("20060102T150405Z")
	dateStamp := now.Format("20060102")
	credentialScope := dateStamp + "/" + region + "/s3/aws4_request"
	credential := accessKeyID + "/" + credentialScope

	// The canonical host mirrors the verifier: virtual-host style unless
	// the bucket name itself contains dots.
	canonicalHost := bucket + ".s3." + region + ".amazonaws.com"
	if strings.Contains(bucket, ".") {
		canonicalHost = "s3." + region + ".amazonaws.com"
	}
	// The canonical request mirrors the verifier exactly, including its
	// canonical headers block: the host line's own newline plus the
	// block-terminating newline.
	canonicalRequest := method + "\n" +
		"/" + bucket + "/" + key + "\n" +
		"" + "\n" +
		"host:" + canonicalHost + "\n" +
		"\n" +
		"host" + "\n" +
		"UNSIGNED-PAYLOAD"

	stringToSign := buildPresignedStringToSign(amzDate, credentialScope, canonicalRequest)
	signingKey := DeriveSigningKey(secretAccessKey, dateStamp, region, "s3")
	signature := HMACSHA256HexString(signingKey, stringToSign)

	query := url.Values{}
	query.Set("X-Amz-Algorithm", Algorithm)
	query.Set("X-Amz-Credential", credential)
	query.Set("X-Amz-Date", amzDate)
	query.Set("X-Amz-Expires", strconv.Itoa(int(expires.Seconds())))
	query.Set("X-Amz-SignedHeaders", "host")
	query.Set("X-Amz-Signature", signature)

	return scheme + "://" + endpointHost + "/" + bucket + "/" + key + "?" + buildPresignedQueryEncoding(query)
}

// buildPresignedQueryEncoding renders the presign parameters in the
// RFC 3986 form the canonical query string comparison expects.
func buildPresignedQueryEncoding(query url.Values) string {
	keys := make([]string, 0, len(query))
	for key := range query {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var b strings.Builder
	for i, key := range keys {
		if i > 0 {
			b.WriteString("&")
		}
		b.WriteString(rfc3986Encode(key))
		b.WriteString("=")
		b.WriteString(rfc3986Encode(query.Get(key)))
	}
	return b.String()
}
