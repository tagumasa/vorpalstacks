package kinesis

import (
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
	storecommon "vorpalstacks/internal/store/aws/common"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// The HTTP statuses below follow the AWS Kinesis API reference's per-
// operation error tables: the awsJson1_1 protocol transmits the error
// identity in __type and carries every modelled client error at HTTP 400 —
// including ResourceNotFoundException, the KMS family, and
// ProvisionedThroughputExceededException, which REST-style services
// conventionally map to 404/403/429.
var (
	// ErrResourceNotFound is returned when the specified resource does not exist.
	ErrResourceNotFound = awserrors.NewAWSError("ResourceNotFoundException", "Requested resource not found.", http.StatusBadRequest)
	// ErrResourceInUse is returned when the resource is in use.
	ErrResourceInUse = awserrors.NewAWSError("ResourceInUseException", "The resource is in use.", http.StatusBadRequest)
	// ErrInvalidArgument is returned when an argument is invalid.
	ErrInvalidArgument = awserrors.NewAWSError("InvalidArgumentException", "Invalid argument.", http.StatusBadRequest)
	// ErrLimitExceeded is returned when a resource quota is exceeded —
	// the per-stream registered-consumer limit. The API reference maps
	// LimitExceededException to HTTP 400.
	ErrLimitExceeded = awserrors.NewAWSError("LimitExceededException", "The requested resource exceeds the maximum number allowed.", http.StatusBadRequest)
	// ErrProvisionedThroughputExceeded is returned when the provisioned throughput is exceeded.
	ErrProvisionedThroughputExceeded = awserrors.NewAWSError("ProvisionedThroughputExceededException", "Rate exceeded for this shard.", http.StatusBadRequest)
	// ErrExpiredIterator is returned when the iterator has expired.
	ErrExpiredIterator = awserrors.NewAWSError("ExpiredIteratorException", "Iterator expired.", http.StatusBadRequest)
	// ErrInvalidIterator is returned when the iterator is invalid.
	ErrInvalidIterator = awserrors.NewAWSError("InvalidArgumentException", "Invalid iterator.", http.StatusBadRequest)
	// ErrShardClosed is returned when the shard is closed.
	ErrShardClosed = awserrors.NewAWSError("InvalidArgumentException", "Shard is closed.", http.StatusBadRequest)
	// ErrExpiredNextToken is returned when a pagination token cannot be
	// decoded. Declared by ListStreams, ListShards and ListStreamConsumers.
	ErrExpiredNextToken = awserrors.NewAWSError("ExpiredNextTokenException", "The pagination token passed to the operation is expired.", http.StatusBadRequest)
	// ErrDryRunOperation is returned when a request carrying the DryRun
	// member passes its full validation — the model's own wording for the
	// shape: "The request was rejected because the DryRun parameter was
	// specified." Declared by PutRecord, PutRecords, GetRecords,
	// GetShardIterator and SubscribeToShard.
	ErrDryRunOperation = awserrors.NewAWSError("DryRunOperationException", "The request was rejected because the DryRun parameter was specified.", http.StatusBadRequest)
	// The KMS error family StartStreamEncryption declares: key-check
	// failures surface under their modelled identities (awsJson1_1 carries
	// every modelled client error at HTTP 400). KMSOptInRequired and
	// KMSThrottlingException are AWS-account subscription and throttle
	// states the platform's checker cannot produce; they stay declared for
	// the shapes that model them.
	ErrKMSNotFound      = awserrors.NewAWSError("KMSNotFoundException", "The request was rejected because the specified KMS key was not found.", http.StatusBadRequest)
	ErrKMSDisabled      = awserrors.NewAWSError("KMSDisabledException", "The request was rejected because the specified KMS key is disabled.", http.StatusBadRequest)
	ErrKMSInvalidState  = awserrors.NewAWSError("KMSInvalidStateException", "The request was rejected because the specified KMS key is in an invalid state.", http.StatusBadRequest)
	ErrKMSAccessDenied  = awserrors.NewAWSError("KMSAccessDeniedException", "The request was rejected because access to the specified KMS key was denied.", http.StatusBadRequest)
	ErrKMSThrottling    = awserrors.NewAWSError("KMSThrottlingException", "The request was rejected because the KMS throttling limit was exceeded.", http.StatusBadRequest)
	ErrKMSOptInRequired = awserrors.NewAWSError("KMSOptInRequired", "The request was rejected because the account is not opted in to the KMS feature.", http.StatusBadRequest)
	// ErrInternalFailure is the model-declared internal-error shape and the
	// fallback for every store error the table below does not map.
	ErrInternalFailure = awserrors.NewAWSError("InternalFailureException", "The processing of the request failed because of an unknown error, exception, or failure.", http.StatusInternalServerError)
)

// storeErrorMappings maps store-level sentinel errors to Kinesis API errors.
// A duplicate consumer name answers ResourceInUseException — the shape
// RegisterStreamConsumer declares; ResourceAlreadyExistsException is
// declared nowhere in the model.
var storeErrorMappings = []awserrors.StoreErrorMapping{
	{Store: kinesisstore.ErrStreamNotFound, AWS: ErrResourceNotFound},
	{Store: kinesisstore.ErrStreamAlreadyExists, AWS: ErrResourceInUse},
	{Store: kinesisstore.ErrShardNotFound, AWS: ErrResourceNotFound},
	{Store: kinesisstore.ErrConsumerNotFound, AWS: ErrResourceNotFound},
	{Store: kinesisstore.ErrConsumerAlreadyExists, AWS: ErrResourceInUse},
	{Store: kinesisstore.ErrConsumerQuotaExceeded, AWS: ErrLimitExceeded},
	// A tag write whose merged set would cross the documented per-resource
	// total answers the invalid-input identity.
	{Store: storecommon.ErrTagQuotaExceeded, AWS: ErrInvalidArgument},
	{Store: kinesisstore.ErrShardClosed, AWS: ErrShardClosed},
	{Store: kinesisstore.ErrNoActiveShards, AWS: ErrProvisionedThroughputExceeded},
	{Store: kinesisstore.ErrInvalidIterator, AWS: ErrInvalidIterator},
	{Store: kinesisstore.ErrExpiredIterator, AWS: ErrExpiredIterator},
	{Store: kinesisstore.ErrShardsNotAdjacent, AWS: ErrInvalidArgument},
	{Store: kinesisstore.ErrOperationNotSupportedOnOnDemandStream, AWS: ErrInvalidArgument},
	{Store: kinesisstore.ErrInvalidShardCount, AWS: ErrInvalidArgument},
	{Store: kinesisstore.ErrInvalidParameter, AWS: ErrInvalidArgument},
	{Store: kinesisstore.ErrResourceNotFound, AWS: ErrResourceNotFound},
}
