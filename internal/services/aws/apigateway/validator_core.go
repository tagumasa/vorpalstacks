package apigateway

import (
	"encoding/json"

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
	if len(in.Name) > maxRequestValidatorNameLength {
		return nil, NewBadRequestException("name must be 1 to 128 characters")
	}

	validator := &apigateway.RequestValidator{
		Name:                      in.Name,
		ValidateRequestBody:       in.ValidateRequestBody,
		ValidateRequestParameters: in.ValidateRequestParameters,
	}

	created, err := stores.restApis.CreateRequestValidator(apiId, validator)
	if err != nil {
		return nil, toApiGatewayError(err)
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
		return nil, toApiGatewayError(err)
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
		return toApiGatewayError(err)
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
		return nil, toApiGatewayError(err)
	}

	for _, po := range ops {
		handled := false
		switch po.Path {
		case "/name":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			validator.Name = po.Value
		case "/validateRequestBody":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			validator.ValidateRequestBody = po.Value == "true"
		case "/validateRequestParameters":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			validator.ValidateRequestParameters = po.Value == "true"
		}
		if !handled {
			return nil, unknownPatchPathError(po)
		}
	}

	if err := stores.restApis.UpdateRequestValidator(apiId, validator); err != nil {
		return nil, toApiGatewayError(err)
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
		return nil, toApiGatewayError(err)
	}

	return created, nil
}

// getModelCore retrieves a model by name. With flatten set, every external
// model reference in the schema is resolved against the API's other
// models, yielding a self-contained schema.
func (s *APIGatewayService) getModelCore(stores *apiGatewayStores, apiId, modelName string, flatten bool) (*apigateway.Model, error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}
	if modelName == "" {
		return nil, NewBadRequestException("modelName is required")
	}

	model, err := stores.restApis.GetModel(apiId, modelName)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	if !flatten {
		return model, nil
	}

	models, err := stores.restApis.ListModels(apiId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	flattened := *model
	flattened.Schema = flattenModelSchema(model.Schema, models)
	return &flattened, nil
}

// flattenModelSchema resolves external model references in a JSON schema:
// every {"$ref": "<ModelName>"} naming another model of the same API is
// replaced by that model's schema, recursively, so the result is
// self-contained. Internal JSON Pointer references and unknown names are
// left untouched; a schema that is not valid JSON is returned unchanged.
func flattenModelSchema(schema string, models []*apigateway.Model) string {
	byName := make(map[string]*apigateway.Model, len(models))
	for _, m := range models {
		byName[m.Name] = m
	}
	var root interface{}
	if err := json.Unmarshal([]byte(schema), &root); err != nil {
		return schema
	}
	root = resolveSchemaRefs(root, byName, make(map[string]bool))
	out, err := json.Marshal(root)
	if err != nil {
		return schema
	}
	return string(out)
}

// resolveSchemaRefs walks a decoded schema and inlines the schemas of
// referenced models. The visiting set guards against reference cycles;
// a name is released on exit so sibling references to the same model
// still resolve.
func resolveSchemaRefs(node interface{}, models map[string]*apigateway.Model, visiting map[string]bool) interface{} {
	switch v := node.(type) {
	case map[string]interface{}:
		if ref, ok := v["$ref"].(string); ok {
			if target, exists := models[ref]; exists && target.Schema != "" && !visiting[ref] {
				var resolved interface{}
				if err := json.Unmarshal([]byte(target.Schema), &resolved); err == nil {
					visiting[ref] = true
					inlined := resolveSchemaRefs(resolved, models, visiting)
					delete(visiting, ref)
					return inlined
				}
			}
		}
		for k, child := range v {
			v[k] = resolveSchemaRefs(child, models, visiting)
		}
		return v
	case []interface{}:
		for i, child := range v {
			v[i] = resolveSchemaRefs(child, models, visiting)
		}
		return v
	default:
		return node
	}
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
		return toApiGatewayError(err)
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
		return nil, toApiGatewayError(err)
	}

	for _, po := range ops {
		handled := false
		switch po.Path {
		case "/description":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			model.Description = po.Value
		case "/schema":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			if !validateModelSchemaSize(po.Value) {
				return nil, NewBadRequestException("schema must not exceed 400 KB")
			}
			model.Schema = po.Value
		}
		if !handled {
			return nil, unknownPatchPathError(po)
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
