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
