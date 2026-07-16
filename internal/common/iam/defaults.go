package iam

import "time"

// RootUserName is the special user name used for root user access keys.
// Access keys with this user name are treated as root credentials, bypassing
// all IAM policy evaluation in the authoriser. This constant is shared across
// services and stores to ensure consistent root user identification.
const RootUserName = "<root>"

// Default values for IAM caching and validation.
const (
	// DefaultCacheNumCounters is the default number of cache counters.
	DefaultCacheNumCounters = 1e4
	// DefaultCacheMaxCost is the default maximum cache cost.
	DefaultCacheMaxCost = 1 << 24
	// DefaultCacheBufferItems is the default number of buffer items per cache entry.
	DefaultCacheBufferItems = 64
	// DefaultCacheTTL is the default time-to-live for cache entries.
	DefaultCacheTTL = 10 * time.Minute

	// DefaultMaxRoleNameLength is the default maximum length for role names.
	DefaultMaxRoleNameLength = 64
	// DefaultMaxPathLength is the default maximum length for IAM paths.
	DefaultMaxPathLength = 512
)
