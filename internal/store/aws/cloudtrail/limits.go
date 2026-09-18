package cloudtrail

import "time"

// Service limit and default values, each defined once here and referenced
// by name at every enforcement site; raw numbers never appear elsewhere.
// Both the store and the service packages reference these constants.

const (
	// DefaultEventDataStoreRetentionDays is the retention applied to an
	// event data store created without an explicit RetentionPeriod.
	// MinEventDataStoreRetentionDays and MaxEventDataStoreRetentionDays
	// bound the accepted range.
	DefaultEventDataStoreRetentionDays = 366
	MinEventDataStoreRetentionDays     = 7
	MaxEventDataStoreRetentionDays     = 3653

	// EventDataStoreRestoreWindow is how long a deleted event data store
	// stays in PENDING_DELETION before it is removed for good: "the event
	// data store enters a PENDING_DELETION state, and is automatically
	// deleted after a wait period of seven days" (DeleteEventDataStore).
	EventDataStoreRestoreWindow = 7 * 24 * time.Hour

	// PublicKeyValidity is how long a generated log-validation public
	// key remains valid.
	PublicKeyValidity = 365 * 24 * time.Hour

	// DefaultLookupEventsResults is the page size used when a lookup
	// omits MaxResults; MaxLookupEventsResults is the hard cap AWS
	// documents for the operation. MaxLookupAttributeValueLength is the
	// model's LookupAttributeValue length cap.
	DefaultLookupEventsResults    = 50
	MaxLookupEventsResults        = 50
	MaxLookupAttributeValueLength = 2000

	// DefaultListTrailsMaxItems is the ListTrails page size when the
	// caller omits MaxItems; MaxListTrailsMaxItems is the hard cap and
	// the page size used when listing every trail internally.
	DefaultListTrailsMaxItems = 1000
	MaxListTrailsMaxItems     = 1000

	// DefaultListEventDataStoresResults and DefaultListChannelsResults
	// are the page sizes for the two list operations when the caller
	// omits MaxResults. MaxListChannelsResults,
	// MaxListEventDataStoresResults and MaxListImportsResults are the
	// model's MaxResults upper bounds (range 1-1000).
	DefaultListEventDataStoresResults = 100
	DefaultListChannelsResults        = 100
	MaxListChannelsResults            = 1000
	MaxListEventDataStoresResults     = 1000

	// DefaultListImportsResults is the page size for ListImports and
	// ListImportFailures when the caller omits MaxResults, including the
	// store-level fallback. MaxListImportsResults is the model bound.
	DefaultListImportsResults = 50
	MaxListImportsResults     = 1000

	// DefaultQueryResultsPageSize is the page size for GetQueryResults
	// and ListQueries when the caller omits the bound.
	// MaxGetQueryResultsResults and MaxListQueriesResults are the model's
	// MaxResults upper bounds for the two operations (range 1-1000).
	DefaultQueryResultsPageSize = 50
	MaxGetQueryResultsResults   = 1000
	MaxListQueriesResults       = 1000

	// DefaultSampleQueriesResults is the SearchSampleQueries page size
	// when the caller omits MaxResults; MaxSampleQueriesResults caps it.
	DefaultSampleQueriesResults = 10
	MaxSampleQueriesResults     = 50

	// LakeQueryScanBound is the per-page scan size the Lake query engine
	// requests from the store; the engine follows the store's NextToken to
	// exhaustion, so the bound limits page size, never the scanned set.
	LakeQueryScanBound = 1000

	// MaxQueryStatementChars is the QueryStatement length cap: the
	// QueryStatement shape is @length 1-10000 in the model, a count
	// of Unicode scalar values (the Smithy length trait measures
	// strings in scalar values, not bytes; the typed SDK rejects a
	// longer statement client-side, so the guard answers the
	// raw-HTTP path).
	MaxQueryStatementChars = 10000

	// ListQueriesWindow bounds how long a query record stays listable —
	// and, through the hourly retention sweep, stored: "Returns a list of
	// queries and query statuses for the past seven days" (ListQueries).
	ListQueriesWindow = 7 * 24 * time.Hour

	// MaxConcurrentQueries is the number of non-terminal Lake queries the
	// store admits before StartQuery is rejected: "The maximum number of
	// concurrent queries is 10" (MaxConcurrentQueriesException
	// documentation).
	MaxConcurrentQueries = 10

	// LakeQueryDeadline is how long a Lake query may run before its record
	// transitions to TIMED_OUT: "Queries that run for longer than one hour
	// might time out. You can still get partial results that were processed
	// before the query timed out" (Run a query, CloudTrail User Guide).
	LakeQueryDeadline = time.Hour

	// MaxQueryResultFileBytes is the documented ceiling on a delivered
	// query result file's size: "The maximum file size for a query
	// result file is 1 TB" (Download saved query results). Delivery
	// splits earlier, at the platform operational bound
	// MaxQueryResultFileChunkBytes.
	MaxQueryResultFileBytes = 1024 * 1024 * 1024 * 1024

	// MaxQueryResultFileChunkBytes is the platform operational bound on
	// the rendered CSV one delivered result file holds in memory before
	// it is compressed, uploaded and released: the delivery path carries
	// one file at a time, and numbered files smaller than the documented
	// ceiling (MaxQueryResultFileBytes) stay inside the documented form.
	MaxQueryResultFileChunkBytes = 64 * 1024 * 1024

	// MaxEventDataStoreRetentionDaysFixedPricing bounds RetentionPeriod
	// when the event data store's effective billing mode is
	// FIXED_RETENTION_PRICING: "you can set a retention period of up to
	// 2557 days, the equivalent of seven years" (UpdateEventDataStore).
	MaxEventDataStoreRetentionDaysFixedPricing = 2557

	// MaxTrailsPerRegion is the trail quota: "Trails per Region — 5 ...
	// The maximum number of trails per AWS Region. This quota cannot be
	// increased" (Quotas in AWS CloudTrail).
	MaxTrailsPerRegion = 5

	// MaxEventDataStoresPerRegion bounds the event data stores in a
	// region: "Event data stores — 10 ... This includes event data stores
	// in any lifecycle stage" (Quotas in AWS CloudTrail) — a
	// PENDING_DELETION store still occupies the quota.
	MaxEventDataStoresPerRegion = 10

	// MaxChannelsPerRegion bounds the channels used for CloudTrail Lake
	// integrations with event sources outside AWS: "Channels — 25 ...
	// This quota cannot be increased" (Quotas in AWS CloudTrail).
	MaxChannelsPerRegion = 25

	// MaxAuditEventsPerRequest bounds one cloudtrail-data PutAuditEvents
	// request: "You can add up to 100 of these events (or up to 1 MB) per
	// PutAuditEvents request" (Quotas in AWS CloudTrail); the
	// cloudtrail-data model marks the AuditEvents list @length max 100.
	MaxAuditEventsPerRequest = 100

	// MinExternalIDLength and MaxExternalIDLength bound the cloudtrail-data
	// externalId member, and ExternalIDPattern is its model pattern: the
	// ExternalId shape is @length 2-1224 with pattern
	// ^[\w+=,.@:\/-]*$.
	MinExternalIDLength = 2
	MaxExternalIDLength = 1224
	ExternalIDPattern   = `^[\w+=,.@:\/-]*$`

	// Prompt and SearchPhrase bounds from the cloudtrail-data-era AI
	// operations' model shapes: Prompt is @length 3-500,
	// SearchSampleQueriesSearchPhrase @length 2-1000, both with the
	// printable pattern.
	MinPromptLength       = 3
	MaxPromptLength       = 500
	MinSearchPhraseLength = 2
	MaxSearchPhraseLength = 1000

	// MaxS3KeyPrefixLength bounds a trail's S3KeyPrefix: "The maximum
	// length is 200 characters" (CreateTrail S3KeyPrefix).
	MaxS3KeyPrefixLength = 200

	// Selector bounds from the PutEventSelectors reference: a trail can
	// carry up to five event selectors and up to five advanced event
	// selectors (the CreateEventDataStore reference states the same
	// five-selector bound for event data stores); a trail's advanced
	// selectors hold at most 500 condition values in total; basic event
	// selectors are limited to 250 data resources across all selectors.
	MaxEventSelectorsPerTrail = 5
	MaxAdvancedEventSelectors = 5
	MaxAdvancedSelectorValues = 500
	MaxDataResourcesPerTrail  = 250
)
