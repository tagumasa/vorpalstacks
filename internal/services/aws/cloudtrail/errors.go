package cloudtrail

import (
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

// This file is the single error vocabulary for the CloudTrail API surface.
// Every constructor spells one model error shape — wire code, the model's
// httpError status, and its awsQueryError code — exactly once; call sites
// build errors through the constructors and shared sentinels below, never
// through inline awserrors.NewAWSError literals. Status codes and query
// codes come from the vendored Smithy model, whose httpError trait is
// authoritative over the HTTP-status column of the API reference (the
// not-found family is 404 there, not the reference's generic 400).

// ---------------------------------------------------------------------------
// Constructors — one per produced model error shape
// ---------------------------------------------------------------------------

func newInvalidParameterException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidParameterException", message, http.StatusBadRequest).
		SetQueryErrorCode("InvalidParameter")
}

func newInvalidTrailNameException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidTrailNameException", message, http.StatusBadRequest).
		SetQueryErrorCode("InvalidTrailName")
}

func newInvalidS3BucketNameException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidS3BucketNameException", message, http.StatusBadRequest).
		SetQueryErrorCode("InvalidS3BucketName")
}

func newInvalidSnsTopicNameException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidSnsTopicNameException", message, http.StatusBadRequest).
		SetQueryErrorCode("InvalidSnsTopicName")
}

func newInvalidKmsKeyIdException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidKmsKeyIdException", message, http.StatusBadRequest).
		SetQueryErrorCode("InvalidKmsKeyId")
}

func newInvalidCloudWatchLogsLogGroupArnException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidCloudWatchLogsLogGroupArnException", message, http.StatusBadRequest).
		SetQueryErrorCode("InvalidCloudWatchLogsLogGroupArn")
}

func newInvalidCloudWatchLogsRoleArnException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidCloudWatchLogsRoleArnException", message, http.StatusBadRequest).
		SetQueryErrorCode("InvalidCloudWatchLogsRoleArn")
}

func newInvalidEventSelectorsException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidEventSelectorsException", message, http.StatusBadRequest).
		SetQueryErrorCode("InvalidEventSelectors")
}

func newInvalidInsightSelectorsException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidInsightSelectorsException", message, http.StatusBadRequest).
		SetQueryErrorCode("InvalidInsightSelectors")
}

func newInvalidLookupAttributesException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidLookupAttributesException", message, http.StatusBadRequest).
		SetQueryErrorCode("InvalidLookupAttributes")
}

func newInvalidTimeRangeException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidTimeRangeException", message, http.StatusBadRequest).
		SetQueryErrorCode("InvalidTimeRange")
}

func newInvalidNextTokenException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidNextTokenException", message, http.StatusBadRequest).
		SetQueryErrorCode("InvalidNextToken")
}

func newInvalidMaxResultsException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidMaxResultsException", message, http.StatusBadRequest).
		SetQueryErrorCode("InvalidMaxResults")
}

func newInvalidEventDataStoreStatusException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidEventDataStoreStatusException", message, http.StatusBadRequest).
		SetQueryErrorCode("InvalidEventDataStoreStatus")
}

func newInvalidEventDataStoreCategoryException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidEventDataStoreCategoryException", message, http.StatusBadRequest).
		SetQueryErrorCode("InvalidEventDataStoreCategory")
}

func newInvalidQueryStatementException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidQueryStatementException", message, http.StatusBadRequest).
		SetQueryErrorCode("InvalidQueryStatement")
}

func newInvalidQueryStatusException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidQueryStatusException", message, http.StatusBadRequest).
		SetQueryErrorCode("InvalidQueryStatus")
}

func newInvalidTagParameterException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidTagParameterException", message, http.StatusBadRequest).
		SetQueryErrorCode("InvalidTagParameter")
}

func newInvalidImportSourceException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidImportSourceException", message, http.StatusBadRequest).
		SetQueryErrorCode("InvalidImportSource")
}

func newInvalidEventCategoryException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidEventCategoryException", message, http.StatusBadRequest).
		SetQueryErrorCode("InvalidEventCategory")
}

func newInactiveEventDataStoreException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InactiveEventDataStoreException", message, http.StatusBadRequest).
		SetQueryErrorCode("InactiveEventDataStore")
}

func newOperationNotPermittedException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("OperationNotPermittedException", message, http.StatusBadRequest).
		SetQueryErrorCode("OperationNotPermitted")
}

func newTagsLimitExceededException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("TagsLimitExceededException", message, http.StatusBadRequest).
		SetQueryErrorCode("TagsLimitExceeded")
}

func newEventDataStoreTerminationProtectedException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("EventDataStoreTerminationProtectedException", message, http.StatusBadRequest).
		SetQueryErrorCode("EventDataStoreTerminationProtectedException")
}

func newEventDataStoreFederationEnabledException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("EventDataStoreFederationEnabledException", message, http.StatusBadRequest).
		SetQueryErrorCode("EventDataStoreFederationEnabled")
}

func newEventDataStoreHasOngoingImportException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("EventDataStoreHasOngoingImportException", message, http.StatusBadRequest).
		SetQueryErrorCode("EventDataStoreHasOngoingImport")
}

func newAccountHasOngoingImportException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("AccountHasOngoingImportException", message, http.StatusBadRequest).
		SetQueryErrorCode("AccountHasOngoingImport")
}

func newChannelExistsForEDSException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("ChannelExistsForEDSException", message, http.StatusBadRequest).
		SetQueryErrorCode("ChannelExistsForEDS")
}

func newResourceARNNotValidException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("ResourceARNNotValidException", message, http.StatusBadRequest).
		SetQueryErrorCode("ResourceARNNotValid")
}

func newTrailNotFoundException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("TrailNotFoundException", message, http.StatusNotFound).
		SetQueryErrorCode("TrailNotFound")
}

func newTrailAlreadyExistsException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("TrailAlreadyExistsException", message, http.StatusBadRequest).
		SetQueryErrorCode("TrailAlreadyExists")
}

func newEventDataStoreNotFoundException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("EventDataStoreNotFoundException", message, http.StatusNotFound).
		SetQueryErrorCode("EventDataStoreNotFound")
}

func newEventDataStoreAlreadyExistsException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("EventDataStoreAlreadyExistsException", message, http.StatusBadRequest).
		SetQueryErrorCode("EventDataStoreAlreadyExists")
}

func newChannelNotFoundException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("ChannelNotFoundException", message, http.StatusNotFound).
		SetQueryErrorCode("ChannelNotFound")
}

// newChannelARNInvalidException answers a Channel selector that carries
// the ARN prefix but is not a well-formed channel ARN — GetChannel words
// the trigger as "the specified value of ChannelARN is not valid".
func newChannelARNInvalidException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("ChannelARNInvalidException", message, http.StatusBadRequest).
		SetQueryErrorCode("ChannelARNInvalid")
}

// newCloudTrailARNInvalidException answers a trail selector that carries
// the ARN prefix but is not a well-formed trail ARN — "This exception is
// thrown when an operation is called with an ARN that is not valid"
// (DeleteTrail, which lists the trail ARN format).
func newCloudTrailARNInvalidException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("CloudTrailARNInvalidException", message, http.StatusBadRequest).
		SetQueryErrorCode("CloudTrailARNInvalid")
}

func newChannelAlreadyExistsException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("ChannelAlreadyExistsException", message, http.StatusBadRequest).
		SetQueryErrorCode("ChannelAlreadyExists")
}

func newQueryIdNotFoundException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("QueryIdNotFoundException", message, http.StatusNotFound).
		SetQueryErrorCode("QueryIdNotFound")
}

func newImportNotFoundException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("ImportNotFoundException", message, http.StatusNotFound).
		SetQueryErrorCode("ImportNotFound")
}

func newResourcePolicyNotFoundException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("ResourcePolicyNotFoundException", message, http.StatusNotFound).
		SetQueryErrorCode("ResourcePolicyNotFound")
}

func newUnsupportedOperationException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("UnsupportedOperationException", message, http.StatusBadRequest).
		SetQueryErrorCode("UnsupportedOperation")
}

func newS3BucketDoesNotExistException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("S3BucketDoesNotExistException", message, http.StatusNotFound).
		SetQueryErrorCode("S3BucketDoesNotExist")
}

func newInsufficientS3BucketPolicyException(message string) *awserrors.AWSError {
	// The model's httpError trait is 403; the API reference's error table
	// lists 400, and the trait governs.
	return awserrors.NewAWSError("InsufficientS3BucketPolicyException", message, http.StatusForbidden).
		SetQueryErrorCode("InsufficientS3BucketPolicy")
}

func newInsufficientSnsTopicPolicyException(message string) *awserrors.AWSError {
	// The model's httpError trait is 403; the API reference's error table
	// lists 400, and the trait governs.
	return awserrors.NewAWSError("InsufficientSnsTopicPolicyException", message, http.StatusForbidden).
		SetQueryErrorCode("InsufficientSnsTopicPolicy")
}

func newMaxConcurrentQueriesException(message string) *awserrors.AWSError {
	// The model's httpError trait is 429 (Too Many Requests); the API
	// reference's error table misprints it as 400, and the trait governs.
	return awserrors.NewAWSError("MaxConcurrentQueriesException", message, http.StatusTooManyRequests).
		SetQueryErrorCode("MaxConcurrentQueries")
}

func newInvalidSourceException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidSourceException", message, http.StatusBadRequest).
		SetQueryErrorCode("InvalidSource")
}

func newInvalidParameterCombinationException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidParameterCombinationException", message, http.StatusBadRequest).
		SetQueryErrorCode("InvalidParameterCombination")
}

func newInvalidDateRangeException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidDateRangeException", message, http.StatusBadRequest).
		SetQueryErrorCode("InvalidDateRange")
}

func newResourceTypeNotSupportedException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("ResourceTypeNotSupportedException", message, http.StatusBadRequest).
		SetQueryErrorCode("ResourceTypeNotSupported")
}

func newResourcePolicyNotValidException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("ResourcePolicyNotValidException", message, http.StatusBadRequest).
		SetQueryErrorCode("ResourcePolicyNotValid")
}

func newResourceNotFoundException(message string) *awserrors.AWSError {
	// The model's httpError trait for this shape is 400, not the 404 the
	// not-found family otherwise carries; the trait governs.
	return awserrors.NewAWSError("ResourceNotFoundException", message, http.StatusBadRequest).
		SetQueryErrorCode("ResourceNotFound")
}

func newInsightNotEnabledException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InsightNotEnabledException", message, http.StatusBadRequest).
		SetQueryErrorCode("InsightNotEnabled")
}

func newEventDataStoreARNInvalidException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("EventDataStoreARNInvalidException", message, http.StatusBadRequest).
		SetQueryErrorCode("EventDataStoreARNInvalid")
}

// The trails quota is the one resource-quota shape the model binds to
// Forbidden (httpError 403); the event-data-store and channel quotas
// carry 400 like the other client errors above.

func newMaximumNumberOfTrailsExceededException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("MaximumNumberOfTrailsExceededException", message, http.StatusForbidden).
		SetQueryErrorCode("MaximumNumberOfTrailsExceeded")
}

func newEventDataStoreMaxLimitExceededException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("EventDataStoreMaxLimitExceededException", message, http.StatusBadRequest).
		SetQueryErrorCode("EventDataStoreMaxLimitExceeded")
}

func newChannelMaxLimitExceededException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("ChannelMaxLimitExceededException", message, http.StatusBadRequest).
		SetQueryErrorCode("ChannelMaxLimitExceeded")
}

func newInvalidS3PrefixException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidS3PrefixException", message, http.StatusBadRequest).
		SetQueryErrorCode("InvalidS3Prefix")
}

// InactiveQueryException carries no httpError trait; the 400 default
// applies. "The specified query cannot be canceled because it is in the
// FINISHED, FAILED, TIMED_OUT, or CANCELLED state" (CancelQuery).
func newInactiveQueryException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InactiveQueryException", message, http.StatusBadRequest).
		SetQueryErrorCode("InactiveQuery")
}

// OrganizationsNotInUseException is the delegated-administration refusal
// for an account outside an organization: "This exception is thrown when
// the request is made from an Amazon Web Services account that is not a
// member of an organization" (RegisterOrganizationDelegatedAdmin). The
// shape carries httpError 404.
func newOrganizationsNotInUseException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("OrganizationsNotInUseException", message, http.StatusNotFound).
		SetQueryErrorCode("OrganizationsNotInUse")
}

// The cloudtrail-data service (PutAuditEvents) declares its error shapes
// without the management API's Exception suffix and without httpError
// traits; a REST-JSON client error answers 400 with the bare shape name as
// the wire code and no query-protocol code.

func newChannelNotFound(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("ChannelNotFound", message, http.StatusBadRequest)
}

// ChannelInsufficientPermission ("The caller's account ID must be the same
// as the channel owner's account ID", PutAuditEvents reference) has no
// constructor here: the platform holds a single account, every principal
// the authorizer resolves belongs to it, and a foreign access key fails
// authentication before the service is reached — the trigger has no
// substrate. A constructor is added when the shape gains one.

func newChannelUnsupportedSchema(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("ChannelUnsupportedSchema", message, http.StatusBadRequest)
}

func newDuplicatedAuditEventId(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("DuplicatedAuditEventId", message, http.StatusBadRequest)
}

func newInvalidChannelARN(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidChannelARN", message, http.StatusBadRequest)
}

// ---------------------------------------------------------------------------
// Shared sentinels — fixed-message errors and store-mapping targets
// ---------------------------------------------------------------------------

var (
	// ErrTrailNotFound is returned when the specified trail does not exist.
	ErrTrailNotFound = newTrailNotFoundException("Trail not found.")
	// ErrTrailAlreadyExists is returned when attempting to create a trail that already exists.
	ErrTrailAlreadyExists = newTrailAlreadyExistsException("Trail already exists.")
	// ErrInvalidParameter is returned when a parameter is invalid.
	ErrInvalidParameter = newInvalidParameterException("Invalid parameter.")
	// ErrInvalidTrailName is the mapping target for the store's
	// invalid-trail-name sentinel.
	ErrInvalidTrailName = newInvalidTrailNameException("Invalid trail name.")
	// ErrInternalError is returned when an internal server error occurs.
	// InternalFailure is the only 500 code the CloudTrail contract
	// documents (Common Error Types).
	ErrInternalError = awserrors.NewAWSError("InternalFailure", "The request can't be processed right now because of an internal server issue.", http.StatusInternalServerError)
	// ErrEventDataStoreNotFoundException is returned when the specified event data store does not exist.
	ErrEventDataStoreNotFoundException = newEventDataStoreNotFoundException("Event data store not found.")
	// ErrEventDataStoreAlreadyExists is returned when creating an event data store with a name that is already in use.
	ErrEventDataStoreAlreadyExists = newEventDataStoreAlreadyExistsException("Event data store already exists.")
	// ErrChannelNotFound is returned when the specified channel does not exist.
	ErrChannelNotFound = newChannelNotFoundException("Channel not found")
	// ErrChannelAlreadyExists is returned when creating or renaming a channel to a name that is already in use.
	ErrChannelAlreadyExists = newChannelAlreadyExistsException("Channel already exists")
	// ErrQueryIdNotFound is returned when the specified query does not exist.
	ErrQueryIdNotFound = newQueryIdNotFoundException("Query not found")
	// ErrImportNotFound is returned when the specified import does not exist.
	ErrImportNotFound = newImportNotFoundException("Import not found")
	// ErrResourcePolicyNotFound is returned when no resource policy is stored for the resource.
	ErrResourcePolicyNotFound = newResourcePolicyNotFoundException("Resource policy not found")
	// ErrOperationNotPermitted is returned when the operation is not permitted in the current state.
	ErrOperationNotPermitted = newOperationNotPermittedException("Operation not permitted.")
	// ErrInvalidEventCategory is returned when an aggregation event category is outside the model enum.
	ErrInvalidEventCategory = newInvalidEventCategoryException("EventCategory must be one of: insight, lap, management, data")
	// ErrInvalidNextToken is returned when a lookup continuation token was
	// not issued by this service or does not belong to the call's
	// parameterisation.
	ErrInvalidNextTokenException = newInvalidNextTokenException("The token is not valid.")
	// ErrMaxConcurrentQueries is returned when a query admission would
	// exceed the concurrent-query bound.
	ErrMaxConcurrentQueries = newMaxConcurrentQueriesException(
		"You are already running the maximum number of concurrent queries. Wait for some queries to finish, and then run the query again.")
	// ErrInvalidSource is the mapping target for the store's
	// channel-source-in-use sentinel and the Source-bound violations.
	ErrInvalidSource = newInvalidSourceException("Source is not valid.")
)

// storeErrorMappings maps store-level sentinel errors to CloudTrail API
// errors; mapStoreError consults it for every store failure.
var storeErrorMappings = []awserrors.StoreErrorMapping{
	{Store: cloudtrailstore.ErrTrailNotFound, AWS: ErrTrailNotFound},
	{Store: cloudtrailstore.ErrTrailAlreadyExists, AWS: ErrTrailAlreadyExists},
	{Store: cloudtrailstore.ErrInvalidTrailName, AWS: ErrInvalidTrailName},
	{Store: cloudtrailstore.ErrEventDataStoreNotFound, AWS: ErrEventDataStoreNotFoundException},
	{Store: cloudtrailstore.ErrEventDataStoreAlreadyExists, AWS: ErrEventDataStoreAlreadyExists},
	{Store: cloudtrailstore.ErrEventDataStoreNotPendingDeletion, AWS: ErrOperationNotPermitted},
	{Store: cloudtrailstore.ErrChannelNotFound, AWS: ErrChannelNotFound},
	{Store: cloudtrailstore.ErrChannelAlreadyExists, AWS: ErrChannelAlreadyExists},
	{Store: cloudtrailstore.ErrChannelSourceInUse, AWS: ErrInvalidSource},
	{Store: cloudtrailstore.ErrQueryNotFound, AWS: ErrQueryIdNotFound},
	{Store: cloudtrailstore.ErrImportNotFound, AWS: ErrImportNotFound},
	{Store: cloudtrailstore.ErrInvalidNextToken, AWS: ErrInvalidNextTokenException},
	{Store: cloudtrailstore.ErrMaxConcurrentQueries, AWS: ErrMaxConcurrentQueries},
}
