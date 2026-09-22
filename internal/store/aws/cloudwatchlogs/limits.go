package cloudwatchlogs

import (
	"time"
)

// The store's bounds register: every documented limit, ceiling, quota and
// internal capacity cap the CloudWatch Logs store enforces, one value one
// definition — the raw number exists at its defining constant alone. The
// record-status vocabularies are not bounds and live with their record
// families in the operations files.
const (
	// MaxChunkSize is the maximum number of log entries per chunk.
	MaxChunkSize = 10000
	// MaxRetentionDays is the maximum retention period in days.
	MaxRetentionDays = 3653
	// MaxLookupTables is the documented per-account, per-Region quota of
	// lookup tables.
	MaxLookupTables = 100
	// MaxActiveImportTasks is the documented account-level concurrency
	// quota of import tasks: the CreateImportTask operation documentation
	// states "There can be no more than 3 active imports per account at a
	// given time". Active means IN_PROGRESS — the ImportStatus vocabulary
	// the model enumerates carries no PENDING.
	MaxActiveImportTasks = 3
	// DefaultDescribeLookupTablesResults is the documented default of the
	// maxResults parameter of DescribeLookupTables.
	DefaultDescribeLookupTablesResults = 50
	// MaxDescribeLookupTablesResults is the documented maximum of the
	// maxResults parameter of DescribeLookupTables.
	MaxDescribeLookupTablesResults = 100
	// MaxLookupTableBodyBytes is the documented size ceiling of a lookup
	// table's CSV content (10 MB).
	MaxLookupTableBodyBytes = 10485760
	// MaxLookupTableNameLength is the documented maximum length of a lookup
	// table name.
	MaxLookupTableNameLength = 256
	// MaxLookupTableDescriptionLength is the documented maximum length of a
	// lookup table description.
	MaxLookupTableDescriptionLength = 1024
	// MaxKmsKeyIdLength is the documented maximum length of the kmsKeyId
	// parameter of lookup table and scheduled-query destination operations
	// (the shared KmsKeyId shape).
	MaxKmsKeyIdLength = 256
	// MaxLogGroupIdentifierLength is the documented maximum length of the
	// shared LogGroupIdentifier shape (a log group name or ARN), carried
	// by the logGroupIdentifier/logGroupIdentifiers members.
	MaxLogGroupIdentifierLength = 2048
	// MaxLookupTableTags is the documented maximum number of tags attached
	// to one lookup table resource.
	MaxLookupTableTags = 50
	// MaxLogEventMessageBytes is the documented per-event ceiling of the
	// message member: "Each log event can be no larger than 1 MB" (model
	// member documentation and the PutLogEvents API reference agree).
	MaxLogEventMessageBytes = 1048576
	// MaxPutLogEventsBatchBytes is the documented batch ceiling: "The
	// maximum batch size is 1,048,576 bytes. This size is calculated as the
	// sum of all event messages in UTF-8, plus 26 bytes for each log event"
	// (PutLogEvents API reference, batch constraints).
	MaxPutLogEventsBatchBytes = 1048576
	// PutLogEventOverheadBytes is the per-event overhead the documented
	// batch-size accounting adds to every message.
	PutLogEventOverheadBytes = 26
	// MaxBatchLogEvents is the documented request-plane batch bound on
	// PutLogEvents ("The maximum batch size is 10,000 log events", API
	// reference, batch constraints): the operation's admission check and
	// every caller that sizes batches for the shared ingestion seam read
	// this bound. The value coincides with MaxChunkSize, the storage
	// plane's per-chunk capacity, but the two traits are distinct and
	// carry separate definitions.
	MaxBatchLogEvents = 10000
	// MaxReadResponseBytes is the documented default page budget of
	// GetLogEvents: "If you don't specify a limit, the default is as many
	// log events as can fit in a response size of 1 MB (up to 10,000 log
	// events)". AWS leaves the byte accounting undefined; the platform
	// counts each event's message bytes plus a fixed JSON-envelope
	// estimate (EventResponseEnvelopeBytes).
	MaxReadResponseBytes = 1048576
	// EventResponseEnvelopeBytes estimates the JSON wire cost of one output
	// log event outside its message (the timestamp and ingestionTime
	// numbers, keys, quotes and braces).
	EventResponseEnvelopeBytes = 64
	// DeliveryNameMax is the documented maximum length of a delivery
	// source or delivery destination name (the shared
	// DeliverySourceName/DeliveryDestinationName length trait, 1..60,
	// pattern [\w-]).
	DeliveryNameMax = 60
	// DeliveryIdMax is the documented maximum length of a delivery id
	// (DeliveryId length trait 1..64, pattern alphanumeric).
	DeliveryIdMax = 64
	// DeliveryLogTypeMax is the documented maximum length of a delivery
	// source's logType (LogType length trait, pattern [\w]).
	DeliveryLogTypeMax = 255
	// DeliverySourceConfigurationMemberMax is the documented maximum
	// length of a deliverySourceConfiguration key or value ("Both keys
	// and values must be between 1 and 255 characters in length.",
	// DeliverySourceConfigurationKey/Value length traits).
	DeliverySourceConfigurationMemberMax = 255
	// DeliveryRecordFieldsMax is the documented maximum number of
	// recordFields entries (RecordFields length trait 0..128; each element
	// is 1..64 characters).
	DeliveryRecordFieldsMax = 128
	// DeliveryRecordFieldLenMax is the documented maximum length of one
	// recordFields element.
	DeliveryRecordFieldLenMax = 64
	// DeliveryFieldDelimiterMax is the documented maximum length of the
	// fieldDelimiter member (FieldDelimiter length trait 0..5).
	DeliveryFieldDelimiterMax = 5
	// DeliverySuffixPathMax is the documented maximum length of the S3
	// delivery suffixPath member (DeliverySuffixPath length trait).
	DeliverySuffixPathMax = 256
	// IndexPolicyFieldsMax is the documented per-policy field quota:
	// "As many as 20 fields can be included in the policy" (Field index
	// syntax and quotas).
	IndexPolicyFieldsMax = 20
	// IndexPolicyFieldNameMax is the documented per-field name length
	// bound: "Each field name can include as many as 100 characters".
	IndexPolicyFieldNameMax = 100
	// DescribeFieldIndexesGroupsMax bounds the DescribeFieldIndexes
	// identifier array (DescribeFieldIndexesLogGroupIdentifiers length
	// trait 1..100).
	DescribeFieldIndexesGroupsMax = 100
	// FieldIndexPolicyPrefixQuota is the documented account-level quota
	// of prefix-selected FIELD_INDEX_POLICY policies ("of these policies
	// 20 can use log group name prefix selection criteria").
	FieldIndexPolicyPrefixQuota = 20
	// ListLogGroupsForQueryMinResults and MaxResults bound the listing's
	// maxResults member (ListLogGroupsForQueryMaxResults range trait
	// 50..500); the default page equals the floor.
	ListLogGroupsForQueryMinResults = 50
	ListLogGroupsForQueryMaxResults = 500
	// AggregateAccountIdentifiersMax bounds the accountIdentifiers array
	// of ListAggregateLogGroupSummaries (AccountIds length trait 0..20).
	AggregateAccountIdentifiersMax = 20
	// AggregateDataSourceFiltersMax bounds the dataSources array of
	// ListAggregateLogGroupSummaries (DataSourceFilters length trait 1..5).
	AggregateDataSourceFiltersMax = 5
	// LiveTailSessionLogGroupsMax bounds StartLiveTail's
	// logGroupIdentifiers array (StartLiveTailLogGroupIdentifiers length
	// trait 1..10).
	LiveTailSessionLogGroupsMax = 10
	// QueryLogGroupsMax bounds the groups one StartQuery may name across
	// its logGroupNames and logGroupIdentifiers members ("You can include
	// up to 50 log groups").
	QueryLogGroupsMax = 50
	// MaxConcurrentQueries is the account's concurrent Insights query
	// quota ("You can have up to 100 concurrent CloudWatch Logs insights
	// queries, including queries that have been added to dashboards").
	MaxConcurrentQueries = 100
	// MaxTopkK is the stats topk function's documented k ceiling: "Valid
	// values for k range from 1 to 10000".
	MaxTopkK = 10000
	// MaxJoinKeyValues is the join command's documented cap on the
	// secondary source's unique key values: "The number of unique key
	// values in the secondary data source is limited to 50,000".
	MaxJoinKeyValues = 50000
	// MaxScheduledQueryDescriptionLength bounds a scheduled query's
	// description (the ScheduledQueryDescription length trait 0-1024).
	MaxScheduledQueryDescriptionLength = 1024
	// MaxScheduledQueryTags bounds the tags attached to one scheduled
	// query (the shared Tags shape's length trait 1-50).
	MaxScheduledQueryTags = 50
	// MaxQueryDefinitionIdLength bounds the queryDefinitionId a saved
	// query is updated by (the PutQueryDefinitionInput member's length
	// trait 1-256).
	MaxQueryDefinitionIdLength = 256
	// LiveTailStreamFiltersMax bounds StartLiveTail's logStreamNames and
	// logStreamNamePrefixes arrays (InputLogStreamNames length trait
	// 1..100).
	LiveTailStreamFiltersMax = 100
	// LiveTailFilterPatternMax is the length ceiling of StartLiveTail's
	// logEventFilterPattern member (FilterPattern length trait 0..1024).
	LiveTailFilterPatternMax = 1024
	// LiveTailUpdateEventsMax is the per-update event ceiling: "The array
	// of log events contained in a LiveTailSessionUpdate can include as
	// many as 500 log events" (StartLiveTail operation documentation).
	LiveTailUpdateEventsMax = 500
	// LiveTailBufferUpdatesMax is the slow-client buffer bound: "CloudWatch
	// Logs buffers up to 10 LiveTailSessionUpdate events or 5000 log
	// events, after which it starts dropping the oldest events" — the
	// event cap is the update cap times the per-update ceiling.
	LiveTailBufferUpdatesMax = 10
	// LiveTailConcurrentSessionsMax is the documented concurrent Live Tail
	// session quota per account (quotas page, recorded at the C5
	// adjudication).
	LiveTailConcurrentSessionsMax = 15
	// LiveTailSessionTimeout is the documented session lifetime: "Live
	// Tail sessions time out after three hours" (SessionTimeoutException
	// documentation).
	LiveTailSessionTimeout = 3 * time.Hour
	// The HTTP ingestion endpoints' documented request limits (developer
	// guide, "Limitations" on every endpoint page): a request body of at
	// most 1 MB, at most 256 KB per event and at most 10,000 events.
	HTTPIngestionMaxRequestBytes = 1 << 20
	HTTPIngestionMaxEventBytes   = 256 << 10
	HTTPIngestionMaxEvents       = 10000
)

// Wire-side limit bounds re-derived from the vendored Smithy model
// (cloudwatch-logs 2014-03-28): every Max is the range trait of the named
// shape, every Default the value served when the member is absent.
const (
	// EventsLimit — the limit member of GetLogEvents and FilterLogEvents
	// (range 1..10000).
	DefaultEventsLimit = 10000
	MaxEventsLimit     = 10000
	// MaxQueryLanguageLimit is the documented ceiling of the query
	// language's limit command ("You can specify a limit value of up to
	// 100,000", Logs Insights query syntax — limit).
	MaxQueryLanguageLimit = 100000
	// EventsLimitStartQuery — the limit member of StartQuery (1..100000).
	DefaultStartQueryLimit = 10000
	MaxStartQueryLimit     = 100000
	// GetQueryResultsMaxItems (0..10000).
	DefaultQueryResultsItems = 10000
	MaxQueryResultsItems     = 10000
	// QueryResultRetentionPeriod is how long a finished query stays
	// retrievable through GetQueryResults. The GetQueryResults operation
	// documentation directs to the CloudWatch Logs quotas page for the
	// availability window; the current quotas table carries no explicit
	// row, and the 7-day value follows the corroborated AWS statement of
	// record (AWS re:Post staff answers, secondary references) —
	// re-verify against the quotas page if it regains the row.
	QueryResultRetentionPeriod = 7 * 24 * time.Hour
	// DescribeLimit — the limit member of the Describe* operations
	// (1..50).
	DefaultDescribeLimit = 50
	MaxDescribeLimit     = 50
	// ListLimit — the limit member of ListLogGroups (1..1000).
	DefaultListLogGroupsLimit = 50
	MaxListLogGroupsLimit     = 1000
	// MaxLogGroupsPerRegion — "You can create up to 1,000,000 log groups
	// per Region per account" (CreateLogGroup); the ceiling rejects with
	// LimitExceededException.
	MaxLogGroupsPerRegion = 1_000_000
	// DescribeQueriesMaxResults, ListScheduledQueriesMaxResults,
	// GetScheduledQueryHistoryMaxResults and QueryListMaxResults all carry
	// range 1..1000; DescribeQueryDefinitions serves the full ceiling as
	// its default.
	DefaultListMaxResults = 50
	MaxListMaxResults     = 1000
)

// Platform-internal capacity caps: no AWS wire limit governs these; each
// bounds one internal read or sampling step, never a total.
const (
	// ListingPageSize bounds one store listing page while a full-scan walk
	// (fetch-all streams or filters, S3 object listing) pages through to
	// exhaustion.
	ListingPageSize = 1000
	// GetLogGroupFieldsSampleSize and GetLogFieldsSampleSize bound the
	// streams visited and the events read per stream when the field-scan
	// operations sample a group's schema.
	GetLogGroupFieldsSampleSize = 50
	GetLogFieldsSampleSize      = 20
	// MaxImportObjectBytes caps one imported S3 object read (100 MB).
	MaxImportObjectBytes = 100 * 1024 * 1024
)

// MaxTestMetricFilterLogEventMessages is the model's length trait on
// TestEventMessages (1..50).
const MaxTestMetricFilterLogEventMessages = 50

// MaxLogStreamNameLength is the model's length trait on LogStreamName
// (1..512), expressed as its maximum.
const MaxLogStreamNameLength = 512

// MaxInputLogStreamNames is the model's @length max on the
// InputLogStreamNames array ("Array Members: Minimum number of 1 item.
// Maximum number of 100 items", FilterLogEvents and CreateExportTask);
// an absent array means every stream, so the maximum alone rejects.
const MaxInputLogStreamNames = 100

// MaxListLogGroupTagFilters is the ListLogGroups logGroupTags array's
// @length max ("Array Members: Minimum number of 1 item. Maximum number
// of 5 items").
const MaxListLogGroupTagFilters = 5

// MaxTagFilterValues is the TagFilter values array's @length max ("Array
// Members: Minimum number of 0 items. Maximum number of 5 items").
const MaxTagFilterValues = 5

// MaxListDataSourceFilters is the ListLogGroups dataSources array's
// @length max ("Array Members: Minimum number of 1 item. Maximum number
// of 5 items").
const MaxListDataSourceFilters = 5

// MaxListFieldIndexNames is the ListLogGroups fieldIndexNames array's
// @length max ("You can specify 1 to 20 field index names").
const MaxListFieldIndexNames = 20

// MaxFieldIndexNameFilterLength is the fieldIndexNames element's length
// ceiling ("each with 1 to 512 characters").
const MaxFieldIndexNameFilterLength = 512

// MaxLogGroupNameLength is the model's length trait on LogGroupName
// (1..512), expressed as its maximum.
const MaxLogGroupNameLength = 512

// MaxFilterNameLength is the model's length trait on FilterName
// (1..512), the metric and subscription filter name shape.
const MaxFilterNameLength = 512

// MaxDestinationNameLength is the model's length trait on
// DestinationName (1..512).
const MaxDestinationNameLength = 512

// MaxFilterPatternLength is the model's length trait on FilterPattern
// (0..1024), counted in Unicode characters.
const MaxFilterPatternLength = 1024

// MaxPolicyDocumentLength is the model's length trait on the
// policyDocument string (1..51200), the generic form the policy-family
// operations share; the destination policy and field-index documents
// carry their own equal-valued ceilings separately.
const MaxPolicyDocumentLength = 51200

// MaxQueryDefinitionNameLength is the model's length trait on
// QueryDefinitionName (1..255).
const MaxQueryDefinitionNameLength = 255

// MaxQueryStringLength is the model's length ceiling every query string
// member carries (queryString, queryDefinitionString and the scheduled
// query's queryString: 0/1..10000).
const MaxQueryStringLength = 10000

// MaxScheduledQueryNameLength is the model's length trait on the
// scheduled query's name member (1..300).
const MaxScheduledQueryNameLength = 300

// MaxExportDestinationLength is the model's length trait on
// CreateExportTask's destination bucket member (1..512); the member's
// own ceiling, not the S3 bucket-name limit the shared bucketname
// validator enforces elsewhere.
const MaxExportDestinationLength = 512

// MaxMetricNameLength is the model's length trait on MetricName and
// MetricNamespace (0..255), counted in Unicode characters.
const MaxMetricNameLength = 255

// MaxMetricValueLength is the model's length trait on MetricValue
// (0..100), counted in Unicode characters.
const MaxMetricValueLength = 100

// MaxMetricDimensionLength is the length ceiling a metric filter
// dimension's key and value each carry (0..255), counted in Unicode
// characters.
const MaxMetricDimensionLength = 255

// MaxTagFilterValueLength is the model's length trait on TagFilterValues
// elements (0..259) — deliberately not the TagValue ceiling, whose 256
// the filter form does not share.
const MaxTagFilterValueLength = 259

// MaxQueryIdLength is the model's length trait on QueryId (1..256), the
// shape StartQuery/StopQuery/DescribeQueries address.
const MaxQueryIdLength = 256

// MaxLookupTableDestinationIdentifierLength is the model's length trait
// on the lookup table destination's S3 destinationIdentifier (1..1024).
const MaxLookupTableDestinationIdentifierLength = 1024

// ConfigurationTemplatesLogTypesMax and
// ConfigurationTemplatesDestinationTypesMax are the @length maxima on
// DescribeConfigurationTemplates' logTypes/resourceTypes arrays (10) and
// deliveryDestinationTypes array (4).
const (
	ConfigurationTemplatesLogTypesMax         = 10
	ConfigurationTemplatesDestinationTypesMax = 4
)

// The query language's documented filter-list ceilings (the Query syntax
// page's own limits, distinct from the API members that share their
// values): as many as 10 data sources, 5 name prefixes and 20 account
// identifiers per query.
const (
	MaxQueryDataSources        = 10
	MaxQueryNamePrefixes       = 5
	MaxQueryAccountIdentifiers = 20
)

// MaxAccessPolicyBytes is the PutDestinationPolicy accessPolicy byte cap:
// "An IAM policy document that authorizes cross-account users to deliver
// their log events to the associated destination. This can be up to 5120
// bytes."
const MaxAccessPolicyBytes = 5120

// MaxAccountScopeResourcePolicies is the account-wide resource-policy
// ceiling ("Policy limits - An account can have a maximum of 10 policies
// without resourceARN and one per LogGroup resourceARN").
const MaxAccountScopeResourcePolicies = 10

// MaxDataProtectionPolicyDocumentBytes is the PutDataProtectionPolicy
// policyDocument character ceiling ("The JSON specified in policyDocument
// can be up to 30,720 characters").
const MaxDataProtectionPolicyDocumentBytes = 30720

// MaxListAccountIdentifiers is the accountIdentifiers array's @length max
// shared by ListLogGroups and DescribeLogGroups ("You can specify as many
// as 20 account IDs in the array").
const MaxListAccountIdentifiers = 20

// MaxSubscriptionFiltersPerLogGroup is the per-log-group quota from the
// CloudWatch Logs quotas page ("Subscription filters per log group — Each
// supported Region: 5"). The PutSubscriptionFilter documentation still
// carries the stale sentence "Each log group can have up to two
// subscription filters associated with it"; the quotas registry is the
// current constraint set and wins.
const MaxSubscriptionFiltersPerLogGroup = 5

// MaxMetricFiltersPerGroup is the per-log-group quota from the CloudWatch
// Logs quotas page ("Metrics filters per log group — Each supported
// Region: 100", not adjustable). PutMetricFilter declares
// LimitExceededException, the quota-shaped error of the operation.
const MaxMetricFiltersPerGroup = 100

// MaxRegexFilterPatternsPerGroup is the per-log-group ceiling on filter
// patterns containing regex: "There is a maximum of 5 filter patterns
// containing regex for each log group when creating metric filters or
// subscription filters" (Filter and pattern syntax). The census spans
// both filter families, so the service core counts it above either
// family's store.
const MaxRegexFilterPatternsPerGroup = 5

// LookupTableNamePattern is the documented character set of lookup table
// names: alphanumeric characters and underscores.
const LookupTableNamePattern = `^[a-zA-Z0-9_]+$`

// MaxMetricFilterDimensions is the dimension ceiling a metric filter's
// transformation carries ("One metric filter can include as many as three
// dimensions"; the emitSystemFieldDimensions members "count toward the
// total dimension limit for metric filters").
const MaxMetricFilterDimensions = 3

// Documented bounds of one saved query's parameter list and members
// (PutQueryDefinition parameter constraints; the list length is the
// QueryParameterList length trait).
const (
	MaxQueryParameters               = 20
	MaxQueryParameterNameLength      = 128
	MaxQueryParameterDefaultValueLen = 1024
	MaxQueryParameterDescriptionLen  = 512
)

// Export and import member limits: the ExportTaskName and ImportId length
// traits ("Length Constraints: Minimum length of 1. Maximum length of
// 512." / 256, CreateExportTask and CancelImportTask API references) and
// the documented export deadline ("Export tasks time out after 24 hours.",
// Exporting log data to Amazon S3, CloudWatch Logs User Guide).
const (
	MaxExportTaskNameLength = 512
	MaxImportIdLength       = 256
	MaxExportTaskRuntime    = 24 * time.Hour
)
