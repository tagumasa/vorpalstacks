package apigateway

import (
	"vorpalstacks/internal/store/aws/apigateway"
)

// --- DTOs ---

// RequestValidatorInput carries the parsed wire members of a
// CreateRequestValidator request.
type RequestValidatorInput struct {
	Name                      string
	ValidateRequestBody       bool
	ValidateRequestParameters bool
}

// ModelInput carries the parsed wire members of a CreateModel request.
type ModelInput struct {
	Name        string
	Description string
	Schema      string
	ContentType string
}

// createRequestValidatorCore validates and persists a request validator.
func (s *APIGatewayService) createRequestValidatorCore(
	stores *apiGatewayStores,
	apiId string,
	in *RequestValidatorInput,
) (*apigateway.RequestValidator, error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	if in.Name == "" {
		return nil, NewBadRequestException("name is required")
	}
	if len(in.Name) > 128 {
		return nil, NewBadRequestException("name must be 1 to 128 characters")
	}

	validator := &apigateway.RequestValidator{
		Name:                      in.Name,
		ValidateRequestBody:       in.ValidateRequestBody,
		ValidateRequestParameters: in.ValidateRequestParameters,
	}

	created, err := stores.restApis.CreateRequestValidator(apiId, validator)
	if err != nil {
		return nil, err
	}

	return created, nil
}

// getRequestValidatorCore retrieves a request validator.
func (s *APIGatewayService) getRequestValidatorCore(stores *apiGatewayStores, apiId, validatorId string) (*apigateway.RequestValidator, error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}
	if validatorId == "" {
		return nil, NewBadRequestException("requestValidatorId is required")
	}

	validator, err := stores.restApis.GetRequestValidator(apiId, validatorId)
	if err != nil {
		return nil, ErrNotFoundException
	}

	return validator, nil
}

// deleteRequestValidatorCore removes a request validator.
func (s *APIGatewayService) deleteRequestValidatorCore(stores *apiGatewayStores, apiId, validatorId string) error {
	if apiId == "" {
		return NewBadRequestException("restApiId is required")
	}
	if validatorId == "" {
		return NewBadRequestException("requestValidatorId is required")
	}

	if err := stores.restApis.DeleteRequestValidator(apiId, validatorId); err != nil {
		return ErrNotFoundException
	}

	return nil
}

// updateRequestValidatorCore applies JSON Patch operations to a request
// validator under the api-and-validator key lock.
func (s *APIGatewayService) updateRequestValidatorCore(
	stores *apiGatewayStores,
	apiId, validatorId string,
	ops []PatchOperation,
) (*apigateway.RequestValidator, error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}
	if validatorId == "" {
		return nil, NewBadRequestException("requestValidatorId is required")
	}

	stores.keyLocker.Lock(apiId + ":" + validatorId)
	defer stores.keyLocker.Unlock(apiId + ":" + validatorId)

	validator, err := stores.restApis.GetRequestValidator(apiId, validatorId)
	if err != nil {
		return nil, ErrNotFoundException
	}

	for _, po := range ops {
		switch po.Path {
		case "/name":
			validator.Name = po.Value
		case "/validateRequestBody":
			validator.ValidateRequestBody = po.Value == "true"
		case "/validateRequestParameters":
			validator.ValidateRequestParameters = po.Value == "true"
		}
	}

	if err := stores.restApis.UpdateRequestValidator(apiId, validator); err != nil {
		return nil, err
	}

	return validator, nil
}

// listRequestValidatorsCore returns every request validator of an API;
// pagination is applied by the caller.
func (s *APIGatewayService) listRequestValidatorsCore(stores *apiGatewayStores, apiId string) ([]*apigateway.RequestValidator, error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	return stores.restApis.ListRequestValidators(apiId)
}

// createModelCore validates and persists a model.
func (s *APIGatewayService) createModelCore(
	stores *apiGatewayStores,
	apiId string,
	in *ModelInput,
) (*apigateway.Model, error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	if in.Name == "" {
		return nil, NewBadRequestException("name is required")
	}
	if !validateModelName(in.Name) {
		return nil, NewBadRequestException("name must be alphanumeric")
	}

	model := &apigateway.Model{
		Name:        in.Name,
		Description: in.Description,
		Schema:      in.Schema,
		ContentType: in.ContentType,
	}
	if !validateModelSchemaSize(model.Schema) {
		return nil, NewBadRequestException("schema must not exceed 400 KB")
	}
	if model.ContentType == "" {
		return nil, NewBadRequestException("contentType is required")
	}

	created, err := stores.restApis.CreateModel(apiId, model)
	if err != nil {
		return nil, err
	}

	return created, nil
}

// getModelCore retrieves a model by name.
func (s *APIGatewayService) getModelCore(stores *apiGatewayStores, apiId, modelName string) (*apigateway.Model, error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}
	if modelName == "" {
		return nil, NewBadRequestException("modelName is required")
	}

	model, err := stores.restApis.GetModel(apiId, modelName)
	if err != nil {
		return nil, ErrNotFoundException
	}

	return model, nil
}

// deleteModelCore removes a model by name.
func (s *APIGatewayService) deleteModelCore(stores *apiGatewayStores, apiId, modelName string) error {
	if apiId == "" {
		return NewBadRequestException("restApiId is required")
	}
	if modelName == "" {
		return NewBadRequestException("modelName is required")
	}

	if err := stores.restApis.DeleteModel(apiId, modelName); err != nil {
		return ErrNotFoundException
	}

	return nil
}

// updateModelCore applies JSON Patch operations to a model under the api key
// lock.
func (s *APIGatewayService) updateModelCore(
	stores *apiGatewayStores,
	apiId, modelName string,
	ops []PatchOperation,
) (*apigateway.Model, error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}
	if modelName == "" {
		return nil, NewBadRequestException("modelName is required")
	}

	stores.keyLocker.Lock(apiId)
	defer stores.keyLocker.Unlock(apiId)

	model, err := stores.restApis.GetModel(apiId, modelName)
	if err != nil {
		return nil, ErrNotFoundException
	}

	for _, po := range ops {
		switch po.Path {
		case "/description":
			model.Description = po.Value
		case "/schema":
			if !validateModelSchemaSize(po.Value) {
				return nil, NewBadRequestException("schema must not exceed 400 KB")
			}
			model.Schema = po.Value
		}
	}

	if err := stores.restApis.UpdateModel(apiId, model); err != nil {
		return nil, toApiGatewayError(err)
	}

	return model, nil
}

// listModelsCore returns every model of an API; pagination is applied by the
// caller.
func (s *APIGatewayService) listModelsCore(stores *apiGatewayStores, apiId string) ([]*apigateway.Model, error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	return stores.restApis.ListModels(apiId)
}
