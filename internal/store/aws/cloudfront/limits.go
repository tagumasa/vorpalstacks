package cloudfront

// AWS specification limit and default values for CloudFront. Each constant
// below is the single definition of its value; every other package must
// reference these instead of inlining the numbers.

// MaxInvalidationPathsPerRequest is the AWS quota for the number of paths
// in a single invalidation request.
const MaxInvalidationPathsPerRequest = 3000

// DefaultListMaxItems is the page size CloudFront list operations apply
// when the request omits MaxItems or passes a non-positive value. AWS
// publishes no upper bound beyond this default, except for
// ListDistributionsByWebACLId (see MaxListDistributionsByWebACLIdItems).
const DefaultListMaxItems = 100

// MaxListDistributionsByWebACLIdItems is the documented maximum value of
// the MaxItems request parameter for ListDistributionsByWebACLId, whose
// API reference states that the maximum and default values are both 100.
const MaxListDistributionsByWebACLIdItems = 100

// MaxPublicKeysPerKeyGroup is the AWS hard limit for public keys in a
// single key group.
const MaxPublicKeysPerKeyGroup = 5

// MaxAliasItemLength is the maximum length of a single CNAME alias per the
// Smithy aliasString shape (@length(min:0, max:253)).
const MaxAliasItemLength = 253

// LegacyDefaultTTLSeconds is the default TTL a cache behaviour applies when
// it carries no cache policy and the request omits DefaultTTL: the
// Developer Guide (Manage how long content stays in the cache) states the
// default TTL is 24 hours when a cache policy is not used.
const LegacyDefaultTTLSeconds = 86400

// DefaultErrorCachingTTLSeconds is the duration CloudFront caches a 4xx or
// 5xx origin response for when it carries no Cache-Control max-age or
// s-maxage directive and no ErrorCachingMinTTL is configured (How
// CloudFront processes HTTP 4xx and 5xx status codes from your origin).
const DefaultErrorCachingTTLSeconds = 10

// MaxContinuousDeploymentPolicies is the AWS quota for continuous
// deployment policies per account (Quotas and other considerations for
// continuous deployment).
const MaxContinuousDeploymentPolicies = 20

// MaxContinuousDeploymentWeight is the maximum fraction of traffic a
// weight-based continuous deployment policy may route to the staging
// distribution.
const MaxContinuousDeploymentWeight = 0.15

// MinSessionStickinessTTLSeconds and MaxSessionStickinessTTLSeconds bound
// both the idle duration and the maximum session duration of a sticky
// weight-based continuous deployment policy.
const (
	MinSessionStickinessTTLSeconds = 300
	MaxSessionStickinessTTLSeconds = 3600
)
