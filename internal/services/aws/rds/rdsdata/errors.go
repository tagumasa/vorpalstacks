package rdsdata

import (
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
)

func badRequest(msg string) *awserrors.AWSError {
	return awserrors.NewAWSError("BadRequestException", msg, http.StatusBadRequest)
}

func invalidParam(msg string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidParameterException", msg, http.StatusBadRequest)
}

func accessDenied(msg string) *awserrors.AWSError {
	return awserrors.NewAWSError("AccessDeniedException", msg, http.StatusForbidden)
}

func internalError(msg string) *awserrors.AWSError {
	return awserrors.NewAWSError("InternalServerException", msg, http.StatusInternalServerError)
}

func transactionNotFound(msg string) *awserrors.AWSError {
	return awserrors.NewAWSError("TransactionNotFoundException", msg, http.StatusNotFound)
}
