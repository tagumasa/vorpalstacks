// Package apigateway provides API Gateway storage functionality for vorpalstacks.
package apigateway

import "errors"

var (
	// ErrRestApiNotFound is returned when the specified REST API does not exist.
	ErrRestApiNotFound = errors.New("rest api not found")

	// ErrRestApiAlreadyExists is returned when attempting to create a REST API
	// that already exists.
	ErrRestApiAlreadyExists = errors.New("rest api already exists")

	// ErrResourceNotFound is returned when the specified API resource
	// does not exist.
	ErrResourceNotFound = errors.New("resource not found")

	// ErrResourceAlreadyExists is returned when attempting to create a resource
	// that already exists.
	ErrResourceAlreadyExists = errors.New("resource already exists")

	// ErrMethodNotFound is returned when the specified API method does not exist.
	ErrMethodNotFound = errors.New("method not found")

	// ErrIntegrationNotFound is returned when the specified integration
	// does not exist.
	ErrIntegrationNotFound = errors.New("integration not found")

	// ErrDeploymentNotFound is returned when the specified deployment
	// does not exist.
	ErrDeploymentNotFound = errors.New("deployment not found")

	// ErrDeploymentAlreadyExists is returned when attempting to create a deployment
	// that already exists.
	ErrDeploymentAlreadyExists = errors.New("deployment already exists")

	// ErrDeploymentInUse is returned when the deployment is in use by a stage
	// and cannot be deleted.
	ErrDeploymentInUse = errors.New("deployment is in use by a stage")

	// ErrStageNotFound is returned when the specified stage does not exist.
	ErrStageNotFound = errors.New("stage not found")

	// ErrStageAlreadyExists is returned when attempting to create a stage
	// that already exists.
	ErrStageAlreadyExists = errors.New("stage already exists")

	// ErrRequestValidatorNotFound is returned when the specified request validator
	// does not exist.
	ErrRequestValidatorNotFound = errors.New("request validator not found")

	// ErrRequestValidatorAlreadyExists is returned when attempting to create a
	// request validator that already exists.
	ErrRequestValidatorAlreadyExists = errors.New("request validator already exists")

	// ErrModelNotFound is returned when the specified model does not exist.
	ErrModelNotFound = errors.New("model not found")

	// ErrModelAlreadyExists is returned when attempting to create a model
	// that already exists.
	ErrModelAlreadyExists = errors.New("model already exists")

	// ErrApiKeyNotFound is returned when the specified API key does not exist.
	ErrApiKeyNotFound = errors.New("api key not found")

	// ErrApiKeyAlreadyExists is returned when attempting to create an API key
	// that already exists.
	ErrApiKeyAlreadyExists = errors.New("api key already exists")

	// ErrUsagePlanNotFound is returned when the specified usage plan does not exist.
	ErrUsagePlanNotFound = errors.New("usage plan not found")

	// ErrUsagePlanAlreadyExists is returned when attempting to create a usage plan
	// that already exists.
	ErrUsagePlanAlreadyExists = errors.New("usage plan already exists")

	// ErrUsagePlanKeyNotFound is returned when the specified usage plan key
	// does not exist.
	ErrUsagePlanKeyNotFound = errors.New("usage plan key not found")

	// ErrUsagePlanKeyAlreadyExists is returned when attempting to create a usage
	// plan key that already exists.
	ErrUsagePlanKeyAlreadyExists = errors.New("usage plan key already exists")

	// ErrDomainNameNotFound is returned when the specified domain name
	// does not exist.
	ErrDomainNameNotFound = errors.New("domain name not found")

	// ErrDomainNameAlreadyExists is returned when attempting to create a domain name
	// that already exists.
	ErrDomainNameAlreadyExists = errors.New("domain name already exists")

	// ErrBasePathMappingNotFound is returned when the specified base path mapping
	// does not exist.
	ErrBasePathMappingNotFound = errors.New("base path mapping not found")

	// ErrBasePathMappingAlreadyExists is returned when attempting to create a base
	// path mapping that already exists.
	ErrBasePathMappingAlreadyExists = errors.New("base path mapping already exists")

	// ErrAuthorizerNotFound is returned when the specified authorizer does not exist.
	ErrAuthorizerNotFound = errors.New("authorizer not found")

	// ErrAuthorizerAlreadyExists is returned when attempting to create an authorizer
	// that already exists.
	ErrAuthorizerAlreadyExists = errors.New("authorizer already exists")

	// ErrMethodResponseNotFound is returned when the specified method response
	// does not exist.
	ErrMethodResponseNotFound = errors.New("method response not found")

	// ErrIntegrationResponseNotFound is returned when the specified integration
	// response does not exist.
	ErrIntegrationResponseNotFound = errors.New("integration response not found")

	// ErrInvalidParameter is returned when a parameter is not valid.
	ErrInvalidParameter = errors.New("invalid parameter")
)
