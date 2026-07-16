// Package validator provides API Gateway request validation functionality for vorpalstacks.
package validator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	apigatewaystore "vorpalstacks/internal/store/aws/apigateway"

	"github.com/xeipuuv/gojsonschema"
)

// Validator validates API Gateway requests.
type Validator struct{}

// NewValidator creates a new validator.
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateRequest validates an API Gateway request against method configuration.
func (v *Validator) ValidateRequest(method *apigatewaystore.Method, restAPI *apigatewaystore.RestApi, headers map[string]string, queryParams map[string]string, pathParams map[string]string, body []byte) *ValidationError {
	if method == nil {
		return nil
	}

	validatorID := method.RequestValidatorId
	if validatorID == "" {
		return nil
	}

	validator, ok := restAPI.RequestValidators[validatorID]
	if !ok {
		return nil
	}

	if validator.ValidateRequestParameters {
		if err := v.validateRequestParameters(method.RequestParameters, headers, queryParams, pathParams); err != nil {
			return err
		}
	}

	if validator.ValidateRequestBody {
		if err := v.validateRequestBody(method.RequestModels, restAPI.Models, headers, body); err != nil {
			return err
		}
	}

	return nil
}

func (v *Validator) validateRequestParameters(requestParams map[string]bool, headers map[string]string, queryParams map[string]string, pathParams map[string]string) *ValidationError {
	missing := []string{}

	for param, required := range requestParams {
		if !required {
			continue
		}

		parts := strings.Split(param, ".")
		if len(parts) < 4 {
			continue
		}

		location := parts[2]
		name := parts[3]

		var present bool
		switch location {
		case "path":
			_, present = pathParams[name]
		case "querystring":
			_, present = queryParams[name]
		case "header":
			_, present = headers[name]
			if !present {
				_, present = headers[http.CanonicalHeaderKey(name)]
			}
		}

		if !present {
			missing = append(missing, name)
		}
	}

	if len(missing) > 0 {
		return &ValidationError{
			Message:  fmt.Sprintf("Missing required request parameters: [%s]", strings.Join(missing, ", ")),
			Type:     "BadRequestException",
			HTTPCode: http.StatusBadRequest,
		}
	}

	return nil
}

func (v *Validator) validateRequestBody(requestModels map[string]string, models map[string]*apigatewaystore.Model, headers map[string]string, body []byte) *ValidationError {
	if len(body) == 0 {
		return nil
	}

	contentType := headers["Content-Type"]
	if contentType == "" {
		contentType = "application/json"
	}

	modelName, ok := requestModels[contentType]
	if !ok {
		modelName, ok = requestModels["application/json"]
	}
	if !ok {
		return nil
	}

	model, ok := models[modelName]
	if !ok {
		return nil
	}

	if model.Schema == "" {
		var bodyObj interface{}
		if err := json.Unmarshal(body, &bodyObj); err != nil {
			return &ValidationError{
				Message:  fmt.Sprintf("Invalid request body: %v", err),
				Type:     "BadRequestException",
				HTTPCode: http.StatusBadRequest,
			}
		}
		return nil
	}

	schemaLoader := gojsonschema.NewStringLoader(model.Schema)
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		return nil
	}

	bodyLoader := gojsonschema.NewBytesLoader(body)

	result, err := schema.Validate(bodyLoader)
	if err != nil {
		return &ValidationError{
			Message:  fmt.Sprintf("Validation error: %v", err),
			Type:     "BadRequestException",
			HTTPCode: http.StatusBadRequest,
		}
	}

	if !result.Valid() {
		errors := []string{}
		for _, desc := range result.Errors() {
			errors = append(errors, desc.String())
		}
		return &ValidationError{
			Message:  fmt.Sprintf("Request body does not match model schema: %s", strings.Join(errors, "; ")),
			Type:     "BadRequestException",
			HTTPCode: http.StatusBadRequest,
		}
	}

	return nil
}

// ValidationError represents an error during request validation.
type ValidationError struct {
	Message  string
	Type     string
	HTTPCode int
}

// Error returns the error message for the validation error.
func (e *ValidationError) Error() string {
	return e.Message
}
