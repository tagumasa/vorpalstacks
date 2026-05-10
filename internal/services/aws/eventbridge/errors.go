package eventbridge

import (
	arnutil "vorpalstacks/internal/utils/aws/arn"

	awserrors "vorpalstacks/internal/common/errors"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

var (
	ErrValidation            = awserrors.NewValidationException("Validation error")
	ErrResourceNotFound      = awserrors.NewResourceNotFoundException("Resource", "")
	ErrResourceAlreadyExists = awserrors.NewResourceAlreadyExistsException("Resource")
	ErrInvalidParameter      = awserrors.NewInvalidParameterException("Invalid parameter")
)

func NewResourceNotFoundException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("ResourceNotFoundException", message, 404)
}

func BuildEventBusARN(accountID, region, name string) string {
	return arnutil.NewARNBuilder(accountID, region).Events().EventBus(name)
}

// mapStoreError translates store-level errors into AWS API errors.
func mapStoreError(err error, resourceDesc string) error {
	switch err {
	case eventsstore.ErrEventBusNotFound:
		return NewResourceNotFoundException("Event bus '" + resourceDesc + "' does not exist")
	case eventsstore.ErrEventBusAlreadyExists:
		return awserrors.NewResourceAlreadyExistsException("Event bus '" + resourceDesc + "'")
	case eventsstore.ErrRuleNotFound:
		return NewResourceNotFoundException("Rule '" + resourceDesc + "' does not exist")
	case eventsstore.ErrRuleAlreadyExists:
		return awserrors.NewResourceAlreadyExistsException("Rule '" + resourceDesc + "'")
	case eventsstore.ErrArchiveNotFound:
		return NewResourceNotFoundException("Archive '" + resourceDesc + "' does not exist")
	case eventsstore.ErrArchiveAlreadyExists:
		return awserrors.NewResourceAlreadyExistsException("Archive '" + resourceDesc + "'")
	case eventsstore.ErrConnectionNotFound:
		return NewResourceNotFoundException("Connection '" + resourceDesc + "' does not exist")
	case eventsstore.ErrConnectionAlreadyExists:
		return awserrors.NewResourceAlreadyExistsException("Connection '" + resourceDesc + "'")
	case eventsstore.ErrApiDestinationNotFound:
		return NewResourceNotFoundException("Api destination '" + resourceDesc + "' does not exist")
	case eventsstore.ErrApiDestinationAlreadyExists:
		return awserrors.NewResourceAlreadyExistsException("Api destination '" + resourceDesc + "'")
	case eventsstore.ErrReplayNotFound:
		return NewResourceNotFoundException("Replay '" + resourceDesc + "' does not exist")
	case eventsstore.ErrReplayAlreadyExists:
		return awserrors.NewResourceAlreadyExistsException("Replay '" + resourceDesc + "' already exists")
	default:
		return err
	}
}
