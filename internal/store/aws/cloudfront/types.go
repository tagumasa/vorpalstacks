// Package cloudfront provides CloudFront storage functionality for vorpalstacks.
package cloudfront

import (
	"time"

	"vorpalstacks/internal/store/aws/common"
)

// Distribution represents an Amazon CloudFront distribution.
type Distribution struct {
	ID                 string              `json:"id"`
	ARN                string              `json:"arn"`
	ETag               string              `json:"etag"`
	Status             string              `json:"status"`
	DomainName         string              `json:"domainName"`
	DistributionConfig *DistributionConfig `json:"distributionConfig"`
	CallerReference    string              `json:"callerReference"`
	Staging            bool                `json:"staging"`
	CreatedAt          time.Time           `json:"createdAt"`
	LastModifiedAt     time.Time           `json:"lastModifiedAt"`
	Enabled            bool                `json:"enabled"`
}

// DistributionConfig represents the configuration for a CloudFront distribution.
type DistributionConfig struct {
	CallerReference      string                `json:"callerReference"`
	Comment              string                `json:"comment"`
	DefaultRootObject    string                `json:"defaultRootObject"`
	Origins              Origins               `json:"origins"`
	DefaultCacheBehavior *CacheBehavior        `json:"defaultCacheBehavior"`
	CacheBehaviors       *CacheBehaviors       `json:"cacheBehaviors"`
	CustomErrorResponses *CustomErrorResponses `json:"customErrorResponses"`
	Logging              *LoggingConfig        `json:"logging"`
	PriceClass           string                `json:"priceClass"`
	Enabled              bool                  `json:"enabled"`
	ViewerCertificate    *ViewerCertificate    `json:"viewerCertificate"`
	HttpVersion          string                `json:"httpVersion"`
	IsIPV6Enabled        bool                  `json:"isIPV6Enabled"`
	AliasICPRecordal     string                `json:"aliasICPRecordal"`
	Aliases              *Aliases              `json:"aliases"`
	Staging              bool                  `json:"staging"`
	WebACLId             string                `json:"webACLId"`
	Restrictions         *Restrictions         `json:"restrictions"`
}

// Aliases represents the alternate domain names (CNAMEs) for a distribution.
type Aliases struct {
	Quantity int      `json:"quantity"`
	Items    []string `json:"items"`
}

// Restrictions represents the geo restriction configuration.
type Restrictions struct {
	GeoRestriction GeoRestriction `json:"geoRestriction"`
}

// GeoRestriction represents the geo restriction configuration.
type GeoRestriction struct {
	RestrictionType string `json:"restrictionType"`
	Quantity        int    `json:"quantity"`
}

// Origins represents a collection of origins for a CloudFront distribution.
type Origins struct {
	Quantity int       `json:"quantity"`
	Items    []*Origin `json:"items"`
}

// Origin represents an origin for a CloudFront distribution.
type Origin struct {
	ID                    string              `json:"id"`
	DomainName            string              `json:"domainName"`
	OriginPath            string              `json:"originPath"`
	CustomHeaders         *CustomHeaders      `json:"customHeaders"`
	S3OriginConfig        *S3OriginConfig     `json:"s3OriginConfig"`
	CustomOriginConfig    *CustomOriginConfig `json:"customOriginConfig"`
	ConnectionAttempts    int                 `json:"connectionAttempts"`
	ConnectionTimeout     int                 `json:"connectionTimeout"`
	OriginShield          *OriginShield       `json:"originShield"`
	OriginAccessControlId string              `json:"originAccessControlId"`
}

// CustomHeaders represents custom headers to be sent to the origin.
type CustomHeaders struct {
	Quantity int `json:"quantity"`
}

// S3OriginConfig represents the configuration for an S3 origin.
type S3OriginConfig struct {
	OriginAccessIdentity string `json:"originAccessIdentity"`
}

// CustomOriginConfig represents the configuration for a custom origin.
type CustomOriginConfig struct {
	HTTPPort               int      `json:"httpPort"`
	HTTPSPort              int      `json:"httpsPort"`
	OriginProtocolPolicy   string   `json:"originProtocolPolicy"`
	OriginSslProtocols     []string `json:"originSslProtocols"`
	OriginReadTimeout      int      `json:"originReadTimeout"`
	OriginKeepaliveTimeout int      `json:"originKeepaliveTimeout"`
}

// OriginShield represents the Origin Shield configuration.
type OriginShield struct {
	Enabled            bool   `json:"enabled"`
	OriginShieldRegion string `json:"originShieldRegion"`
}

// CacheBehaviors represents a collection of cache behaviors for a CloudFront distribution.
type CacheBehaviors struct {
	Quantity int              `json:"quantity"`
	Items    []*CacheBehavior `json:"items"`
}

// CacheBehavior represents a cache behaviour configuration for a CloudFront distribution.
type CacheBehavior struct {
	PathPattern                string                      `json:"pathPattern"`
	ViewerProtocolPolicy       string                      `json:"viewerProtocolPolicy"`
	AllowedMethods             *AllowedMethods             `json:"allowedMethods"`
	Compress                   bool                        `json:"compress"`
	CachePolicyId              string                      `json:"cachePolicyId"`
	OriginRequestPolicyId      string                      `json:"originRequestPolicyId"`
	ResponseHeadersPolicyId    string                      `json:"responseHeadersPolicyId"`
	LambdaFunctionAssociations *LambdaFunctionAssociations `json:"lambdaFunctionAssociations"`
	FunctionAssociations       *FunctionAssociations       `json:"functionAssociations"`
	TargetOriginId             string                      `json:"targetOriginId"`
	ForwardedValues            *ForwardedValues            `json:"forwardedValues"`
	MinTTL                     int                         `json:"minTTL"`
	DefaultTTL                 int                         `json:"defaultTTL"`
	MaxTTL                     int                         `json:"maxTTL"`
	SmoothStreaming            bool                        `json:"smoothStreaming"`
	TrustedSigners             *TrustedSigners             `json:"trustedSigners"`
	TrustedKeyGroups           *TrustedKeyGroups           `json:"trustedKeyGroups"`
}

// AllowedMethods represents the HTTP methods that CloudFront caches and forwards to the origin.
type AllowedMethods struct {
	Quantity      int      `json:"quantity"`
	CachedMethods []string `json:"cachedMethods"`
	Items         []string `json:"items"`
}

// ForwardedValues represents the query string parameters that CloudFront forwards to the origin.
type ForwardedValues struct {
	QueryString          bool                  `json:"queryString"`
	Cookies              *CookiePreferences    `json:"cookies"`
	Headers              *Headers              `json:"headers"`
	QueryStringCacheKeys *QueryStringCacheKeys `json:"queryStringCacheKeys"`
}

// CookiePreferences specifies how CloudFront handles cookies.
type CookiePreferences struct {
	Forward string `json:"forward"`
}

// Headers specifies the HTTP headers that CloudFront forwards to the origin.
type Headers struct {
	Quantity int      `json:"quantity"`
	Items    []string `json:"items"`
}

// QueryStringCacheKeys specifies the query string parameters that determine which objects CloudFront caches.
type QueryStringCacheKeys struct {
	Quantity int `json:"quantity"`
}

// LambdaFunctionAssociations represents Lambda@Edge functions associated with a cache behaviour.
type LambdaFunctionAssociations struct {
	Quantity int `json:"quantity"`
}

// FunctionAssociations represents CloudFront Functions associated with a cache behaviour.
type FunctionAssociations struct {
	Quantity int `json:"quantity"`
}

// TrustedSigners specifies the trusted signers for a cache behaviour.
type TrustedSigners struct {
	Enabled  bool     `json:"enabled"`
	Quantity int      `json:"quantity"`
	Items    []string `json:"items"`
}

// TrustedKeyGroups specifies the trusted key groups for a cache behaviour.
type TrustedKeyGroups struct {
	Enabled  bool     `json:"enabled"`
	Quantity int      `json:"quantity"`
	Items    []string `json:"items"`
}

// CustomErrorResponses represents custom error responses for a CloudFront distribution.
type CustomErrorResponses struct {
	Quantity int `json:"quantity"`
}

// LoggingConfig specifies logging configuration for the distribution.
type LoggingConfig struct {
	Enabled        bool   `json:"enabled"`
	IncludeCookies bool   `json:"includeCookies"`
	Bucket         string `json:"bucket"`
	Prefix         string `json:"prefix"`
}

// ViewerCertificate specifies the SSL/TLS certificate for the distribution.
type ViewerCertificate struct {
	CloudFrontDefaultCertificate bool   `json:"cloudFrontDefaultCertificate"`
	ACMCertificateArn            string `json:"acmCertificateArn"`
	IAMCertificateId             string `json:"iamCertificateId"`
	SSLSupportMethod             string `json:"sslSupportMethod"`
	MinimumProtocolVersion       string `json:"minimumProtocolVersion"`
	Certificate                  string `json:"certificate"`
	CertificateSource            string `json:"certificateSource"`
}

// CachePolicy represents a CloudFront cache policy.
type CachePolicy struct {
	ID                string             `json:"id"`
	ARN               string             `json:"arn"`
	ETag              string             `json:"etag"`
	Name              string             `json:"name"`
	CachePolicyConfig *CachePolicyConfig `json:"cachePolicyConfig"`
	CreatedAt         time.Time          `json:"createdAt"`
	ModifiedAt        time.Time          `json:"modifiedAt"`
}

// CachePolicyConfig represents the configuration for a cache policy.
type CachePolicyConfig struct {
	Name                                     string                `json:"name"`
	Comment                                  string                `json:"comment"`
	DefaultTTL                               int64                 `json:"defaultTTL"`
	MaxTTL                                   int64                 `json:"maxTTL"`
	MinTTL                                   int64                 `json:"minTTL"`
	ParametersInCacheKeyParametersInCacheKey *ParametersInCacheKey `json:"parametersInCacheKey"`
}

// ParametersInCacheKey specifies the parameters in the cache key that CloudFront uses to cache objects.
type ParametersInCacheKey struct {
	EnableAcceptEncodingGzip   bool               `json:"enableAcceptEncodingGzip"`
	EnableAcceptEncodingBrotli bool               `json:"enableAcceptEncodingBrotli"`
	QueryStringConfig          *QueryStringConfig `json:"queryStringConfig"`
	CookieConfig               *CookieConfig      `json:"cookieConfig"`
	HeaderConfig               *HeaderConfig      `json:"headerConfig"`
}

// QueryStringConfig specifies the query string parameters included in the cache key.
type QueryStringConfig struct {
	QueryStringBehavior string        `json:"queryStringBehavior"`
	QueryStrings        *QueryStrings `json:"queryStrings"`
}

// QueryStrings represents a collection of query string parameters.
type QueryStrings struct {
	Quantity int      `json:"quantity"`
	Items    []string `json:"items"`
}

// CookieConfig specifies the cookies included in the cache key.
type CookieConfig struct {
	CookieBehavior string   `json:"cookieBehavior"`
	Cookies        *Cookies `json:"cookies"`
}

// Cookies represents a collection of cookie names.
type Cookies struct {
	Quantity int      `json:"quantity"`
	Items    []string `json:"items"`
}

// HeaderConfig specifies the headers included in the cache key.
type HeaderConfig struct {
	HeaderBehavior string   `json:"headerBehavior"`
	Headers        *Headers `json:"headers"`
}

// OriginRequestPolicy represents an origin request policy.
type OriginRequestPolicy struct {
	ID                        string                     `json:"id"`
	ARN                       string                     `json:"arn"`
	ETag                      string                     `json:"etag"`
	Name                      string                     `json:"name"`
	OriginRequestPolicyConfig *OriginRequestPolicyConfig `json:"originRequestPolicyConfig"`
	CreatedAt                 time.Time                  `json:"createdAt"`
	ModifiedAt                time.Time                  `json:"modifiedAt"`
}

// OriginRequestPolicyConfig represents the configuration for an origin request policy.
type OriginRequestPolicyConfig struct {
	Name               string              `json:"name"`
	Comment            string              `json:"comment"`
	CookiesConfig      *CookiesConfig      `json:"cookiesConfig"`
	HeadersConfig      *HeadersConfig      `json:"headersConfig"`
	QueryStringsConfig *QueryStringsConfig `json:"queryStringsConfig"`
}

// CookiesConfig specifies the cookies to include in origin requests.
type CookiesConfig struct {
	CookieBehavior string   `json:"cookieBehavior"`
	Cookies        *Cookies `json:"cookies"`
}

// HeadersConfig specifies the headers to include in origin requests.
type HeadersConfig struct {
	HeaderBehavior string   `json:"headerBehavior"`
	Headers        *Headers `json:"headers"`
}

// QueryStringsConfig specifies the query strings to include in origin requests.
type QueryStringsConfig struct {
	QueryStringBehavior string        `json:"queryStringBehavior"`
	QueryStrings        *QueryStrings `json:"queryStrings"`
}

// Tag represents a tag for a CloudFront resource.
type Tag = common.Tag

// DistributionListResult represents the result of listing CloudFront distributions.
type DistributionListResult struct {
	Distributions []*Distribution
	IsTruncated   bool
	NextMarker    string
}

// CachePolicyListResult represents the result of listing CloudFront cache policies.
type CachePolicyListResult struct {
	CachePolicies []*CachePolicy
	IsTruncated   bool
	NextMarker    string
}

// OriginRequestPolicyListResult represents the result of listing CloudFront origin request policies.
type OriginRequestPolicyListResult struct {
	OriginRequestPolicies []*OriginRequestPolicy
	IsTruncated           bool
	NextMarker            string
}

// OriginAccessControl represents a CloudFront origin access control.
type OriginAccessControl struct {
	ID                            string    `json:"id"`
	ARN                           string    `json:"arn"`
	ETag                          string    `json:"etag"`
	Name                          string    `json:"name"`
	Description                   string    `json:"description"`
	OriginAccessControlOriginType string    `json:"originType"`
	SigningBehavior               string    `json:"signingBehavior"`
	SigningProtocol               string    `json:"signingProtocol"`
	CreatedAt                     time.Time `json:"createdAt"`
	LastModifiedAt                time.Time `json:"lastModifiedAt"`
}

// OriginAccessControlConfig represents the configuration for an origin access control.
type OriginAccessControlConfig struct {
	Name                          string `json:"name"`
	Description                   string `json:"description"`
	OriginAccessControlOriginType string `json:"originType"`
	SigningBehavior               string `json:"signingBehavior"`
	SigningProtocol               string `json:"signingProtocol"`
}

// OriginAccessControlListResult represents the result of listing CloudFront origin access controls.
type OriginAccessControlListResult struct {
	OriginAccessControls []*OriginAccessControl
	IsTruncated          bool
	NextMarker           string
}

// ResponseHeadersPolicy represents a CloudFront response headers policy.
type ResponseHeadersPolicy struct {
	ID                        string                                          `json:"id"`
	ARN                       string                                          `json:"arn"`
	ETag                      string                                          `json:"etag"`
	Name                      string                                          `json:"name"`
	Comment                   string                                          `json:"comment,omitempty"`
	CorsConfig                *ResponseHeadersPolicyCorsConfig                `json:"corsConfig,omitempty"`
	CustomHeadersConfig       *ResponseHeadersPolicyCustomHeadersConfig       `json:"customHeadersConfig,omitempty"`
	RemoveHeadersConfig       *ResponseHeadersPolicyRemoveHeadersConfig       `json:"removeHeadersConfig,omitempty"`
	SecurityHeadersConfig     *ResponseHeadersPolicySecurityHeadersConfig     `json:"securityHeadersConfig,omitempty"`
	ServerTimingHeadersConfig *ResponseHeadersPolicyServerTimingHeadersConfig `json:"serverTimingHeadersConfig,omitempty"`
	CreatedAt                 time.Time                                       `json:"createdAt"`
	LastModifiedAt            time.Time                                       `json:"lastModifiedAt"`
}

// ResponseHeadersPolicyConfig represents the configuration for a response headers policy.
type ResponseHeadersPolicyConfig struct {
	Name                      string                                          `json:"name"`
	Comment                   string                                          `json:"comment,omitempty"`
	CorsConfig                *ResponseHeadersPolicyCorsConfig                `json:"corsConfig,omitempty"`
	CustomHeadersConfig       *ResponseHeadersPolicyCustomHeadersConfig       `json:"customHeadersConfig,omitempty"`
	RemoveHeadersConfig       *ResponseHeadersPolicyRemoveHeadersConfig       `json:"removeHeadersConfig,omitempty"`
	SecurityHeadersConfig     *ResponseHeadersPolicySecurityHeadersConfig     `json:"securityHeadersConfig,omitempty"`
	ServerTimingHeadersConfig *ResponseHeadersPolicyServerTimingHeadersConfig `json:"serverTimingHeadersConfig,omitempty"`
}

// ResponseHeadersPolicyCorsConfig represents CORS configuration for a response headers policy.
type ResponseHeadersPolicyCorsConfig struct {
	AccessControlAllowCredentials bool                                             `json:"accessControlAllowCredentials"`
	AccessControlAllowHeaders     *ResponseHeadersPolicyAccessControlAllowHeaders  `json:"accessControlAllowHeaders"`
	AccessControlAllowMethods     *ResponseHeadersPolicyAccessControlAllowMethods  `json:"accessControlAllowMethods"`
	AccessControlAllowOrigins     *ResponseHeadersPolicyAccessControlAllowOrigins  `json:"accessControlAllowOrigins"`
	OriginOverride                bool                                             `json:"originOverride"`
	AccessControlExposeHeaders    *ResponseHeadersPolicyAccessControlExposeHeaders `json:"accessControlExposeHeaders,omitempty"`
	AccessControlMaxAgeSec        int32                                            `json:"accessControlMaxAgeSec,omitempty"`
}

// ResponseHeadersPolicyAccessControlAllowHeaders specifies the allowed headers for CORS.
type ResponseHeadersPolicyAccessControlAllowHeaders struct {
	Quantity int32    `json:"quantity"`
	Items    []string `json:"items"`
}

// ResponseHeadersPolicyAccessControlAllowMethods specifies the allowed methods for CORS.
type ResponseHeadersPolicyAccessControlAllowMethods struct {
	Quantity int32    `json:"quantity"`
	Items    []string `json:"items"`
}

// ResponseHeadersPolicyAccessControlAllowOrigins specifies the allowed origins for CORS.
type ResponseHeadersPolicyAccessControlAllowOrigins struct {
	Quantity int32    `json:"quantity"`
	Items    []string `json:"items"`
}

// ResponseHeadersPolicyAccessControlExposeHeaders specifies the exposed headers for CORS.
type ResponseHeadersPolicyAccessControlExposeHeaders struct {
	Quantity int32    `json:"quantity"`
	Items    []string `json:"items,omitempty"`
}

// ResponseHeadersPolicySecurityHeadersConfig specifies security headers configuration.
type ResponseHeadersPolicySecurityHeadersConfig struct {
	ContentSecurityPolicy   *ResponseHeadersPolicyContentSecurityPolicy   `json:"contentSecurityPolicy,omitempty"`
	ContentTypeOptions      *ResponseHeadersPolicyContentTypeOptions      `json:"contentTypeOptions,omitempty"`
	FrameOptions            *ResponseHeadersPolicyFrameOptions            `json:"frameOptions,omitempty"`
	ReferrerPolicy          *ResponseHeadersPolicyReferrerPolicy          `json:"referrerPolicy,omitempty"`
	StrictTransportSecurity *ResponseHeadersPolicyStrictTransportSecurity `json:"strictTransportSecurity,omitempty"`
	XSSProtection           *ResponseHeadersPolicyXSSProtection           `json:"xssProtection,omitempty"`
}

// ResponseHeadersPolicyContentSecurityPolicy specifies Content-Security-Policy header configuration.
type ResponseHeadersPolicyContentSecurityPolicy struct {
	ContentSecurityPolicy string `json:"contentSecurityPolicy"`
	Override              bool   `json:"override"`
}

// ResponseHeadersPolicyContentTypeOptions specifies X-Content-Type-Options header configuration.
type ResponseHeadersPolicyContentTypeOptions struct {
	Override bool `json:"override"`
}

// ResponseHeadersPolicyFrameOptions specifies X-Frame-Options header configuration.
type ResponseHeadersPolicyFrameOptions struct {
	FrameOption string `json:"frameOption"`
	Override    bool   `json:"override"`
}

// ResponseHeadersPolicyReferrerPolicy specifies Referrer-Policy header configuration.
type ResponseHeadersPolicyReferrerPolicy struct {
	Override       bool   `json:"override"`
	ReferrerPolicy string `json:"referrerPolicy"`
}

// ResponseHeadersPolicyStrictTransportSecurity specifies Strict-Transport-Security header configuration.
type ResponseHeadersPolicyStrictTransportSecurity struct {
	AccessControlMaxAgeSec int32 `json:"accessControlMaxAgeSec"`
	Override               bool  `json:"override"`
	IncludeSubdomains      bool  `json:"includeSubdomains,omitempty"`
	Preload                bool  `json:"preload,omitempty"`
}

// ResponseHeadersPolicyXSSProtection specifies X-XSS-Protection header configuration.
type ResponseHeadersPolicyXSSProtection struct {
	Override   bool   `json:"override"`
	Protection bool   `json:"protection"`
	ModeBlock  bool   `json:"modeBlock,omitempty"`
	ReportUri  string `json:"reportUri,omitempty"`
}

// ResponseHeadersPolicyCustomHeadersConfig specifies custom headers configuration.
type ResponseHeadersPolicyCustomHeadersConfig struct {
	Quantity int32                               `json:"quantity"`
	Items    []ResponseHeadersPolicyCustomHeader `json:"items,omitempty"`
}

// ResponseHeadersPolicyCustomHeader specifies a custom header to add to responses.
type ResponseHeadersPolicyCustomHeader struct {
	Header   string `json:"header"`
	Override bool   `json:"override"`
	Value    string `json:"value"`
}

// ResponseHeadersPolicyRemoveHeadersConfig specifies headers to remove from responses.
type ResponseHeadersPolicyRemoveHeadersConfig struct {
	Quantity int32                               `json:"quantity"`
	Items    []ResponseHeadersPolicyRemoveHeader `json:"items,omitempty"`
}

// ResponseHeadersPolicyRemoveHeader specifies a header to remove from responses.
type ResponseHeadersPolicyRemoveHeader struct {
	Header string `json:"header"`
}

// ResponseHeadersPolicyServerTimingHeadersConfig specifies Server-Timing header configuration.
type ResponseHeadersPolicyServerTimingHeadersConfig struct {
	Enabled      bool    `json:"enabled"`
	SamplingRate float64 `json:"samplingRate,omitempty"`
}

// ResponseHeadersPolicyListResult represents the result of listing response headers policies.
type ResponseHeadersPolicyListResult struct {
	ResponseHeadersPolicies []*ResponseHeadersPolicy
	IsTruncated             bool
	NextMarker              string
}
