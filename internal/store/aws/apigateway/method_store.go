package apigateway

import "strings"

// resourceForUpdate loads the API and addressed resource for a write
// operation; the caller must hold s.mu and persist via updateLocked.
func (s *RestApiStore) resourceForUpdate(apiId, resourceId string) (*RestApi, *Resource, error) {
	api, err := s.Get(apiId)
	if err != nil {
		return nil, nil, err
	}
	resource, ok := api.Resources[resourceId]
	if !ok {
		return nil, nil, ErrResourceNotFound
	}
	return api, resource, nil
}

// methodForUpdate loads the API, resource, and addressed method for a
// write operation; the caller must hold s.mu.
func (s *RestApiStore) methodForUpdate(apiId, resourceId, httpMethod string) (*RestApi, *Resource, *Method, error) {
	api, resource, err := s.resourceForUpdate(apiId, resourceId)
	if err != nil {
		return nil, nil, nil, err
	}
	method, ok := resource.ResourceMethods[strings.ToUpper(httpMethod)]
	if !ok {
		return nil, nil, nil, ErrMethodNotFound
	}
	return api, resource, method, nil
}

// integrationForUpdate additionally resolves the method's integration for
// a write operation; the caller must hold s.mu.
func (s *RestApiStore) integrationForUpdate(apiId, resourceId, httpMethod string) (*RestApi, *Method, *Integration, error) {
	api, _, method, err := s.methodForUpdate(apiId, resourceId, httpMethod)
	if err != nil {
		return nil, nil, nil, err
	}
	if method.MethodIntegration == nil {
		return nil, nil, nil, ErrIntegrationNotFound
	}
	return api, method, method.MethodIntegration, nil
}

// PutMethod creates or updates a method for an API resource.
func (s *RestApiStore) PutMethod(apiId, resourceId string, method *Method) (*Method, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, resource, err := s.resourceForUpdate(apiId, resourceId)
	if err != nil {
		return nil, err
	}

	method.RestApiId = apiId
	method.ResourceId = resourceId
	method.HttpMethod = strings.ToUpper(method.HttpMethod)
	if resource.ResourceMethods == nil {
		resource.ResourceMethods = make(map[string]*Method)
	}
	if method.MethodResponses == nil {
		method.MethodResponses = make(map[string]*MethodResponse)
	}

	resource.ResourceMethods[method.HttpMethod] = method

	if err := s.updateLocked(api); err != nil {
		return nil, err
	}

	return method, nil
}

// GetMethod retrieves a method for an API resource.
func (s *RestApiStore) GetMethod(apiId, resourceId, httpMethod string) (*Method, error) {
	resource, err := s.GetResource(apiId, resourceId)
	if err != nil {
		return nil, err
	}

	method, ok := resource.ResourceMethods[strings.ToUpper(httpMethod)]
	if !ok {
		return nil, ErrMethodNotFound
	}
	return method, nil
}

// DeleteMethod deletes a method from an API resource.
func (s *RestApiStore) DeleteMethod(apiId, resourceId, httpMethod string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, resource, _, err := s.methodForUpdate(apiId, resourceId, httpMethod)
	if err != nil {
		return err
	}

	delete(resource.ResourceMethods, strings.ToUpper(httpMethod))
	return s.updateLocked(api)
}

// PutIntegration creates or updates an integration for a method.
func (s *RestApiStore) PutIntegration(apiId, resourceId, httpMethod string, integration *Integration) (*Integration, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, _, method, err := s.methodForUpdate(apiId, resourceId, httpMethod)
	if err != nil {
		return nil, err
	}

	integration.RestApiId = apiId
	integration.ResourceId = resourceId
	integration.HttpMethod = strings.ToUpper(httpMethod)
	if integration.IntegrationResponses == nil {
		integration.IntegrationResponses = make(map[string]*IntegrationResponse)
	}

	method.MethodIntegration = integration

	return integration, s.updateLocked(api)
}

// GetIntegration retrieves an integration for a method.
func (s *RestApiStore) GetIntegration(apiId, resourceId, httpMethod string) (*Integration, error) {
	method, err := s.GetMethod(apiId, resourceId, httpMethod)
	if err != nil {
		return nil, err
	}

	if method.MethodIntegration == nil {
		return nil, ErrIntegrationNotFound
	}
	return method.MethodIntegration, nil
}

// DeleteIntegration deletes an integration from a method.
func (s *RestApiStore) DeleteIntegration(apiId, resourceId, httpMethod string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, _, method, err := s.methodForUpdate(apiId, resourceId, httpMethod)
	if err != nil {
		return err
	}

	method.MethodIntegration = nil
	return s.updateLocked(api)
}

// UpdateIntegration updates an integration for a method.
func (s *RestApiStore) UpdateIntegration(apiId, resourceId, httpMethod string, integration *Integration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, method, _, err := s.integrationForUpdate(apiId, resourceId, httpMethod)
	if err != nil {
		return err
	}

	integration.RestApiId = apiId
	integration.ResourceId = resourceId
	integration.HttpMethod = strings.ToUpper(httpMethod)
	method.MethodIntegration = integration
	return s.updateLocked(api)
}

// PutIntegrationResponse creates or updates an integration response for a method.
func (s *RestApiStore) PutIntegrationResponse(apiId, resourceId, httpMethod, statusCode string, response *IntegrationResponse) (*IntegrationResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, _, integration, err := s.integrationForUpdate(apiId, resourceId, httpMethod)
	if err != nil {
		return nil, err
	}

	if integration.IntegrationResponses == nil {
		integration.IntegrationResponses = make(map[string]*IntegrationResponse)
	}
	integration.IntegrationResponses[statusCode] = response

	return response, s.updateLocked(api)
}

// GetIntegrationResponse retrieves an integration response for a method.
func (s *RestApiStore) GetIntegrationResponse(apiId, resourceId, httpMethod, statusCode string) (*IntegrationResponse, error) {
	integration, err := s.GetIntegration(apiId, resourceId, httpMethod)
	if err != nil {
		return nil, err
	}

	response, ok := integration.IntegrationResponses[statusCode]
	if !ok {
		return nil, ErrIntegrationNotFound
	}
	return response, nil
}

// DeleteIntegrationResponse deletes an integration response from a method.
func (s *RestApiStore) DeleteIntegrationResponse(apiId, resourceId, httpMethod, statusCode string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, _, integration, err := s.integrationForUpdate(apiId, resourceId, httpMethod)
	if err != nil {
		return err
	}

	delete(integration.IntegrationResponses, statusCode)
	return s.updateLocked(api)
}

// UpdateIntegrationResponse updates an integration response for a method.
func (s *RestApiStore) UpdateIntegrationResponse(apiId, resourceId, httpMethod, statusCode string, response *IntegrationResponse) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, _, integration, err := s.integrationForUpdate(apiId, resourceId, httpMethod)
	if err != nil {
		return err
	}

	if _, ok := integration.IntegrationResponses[statusCode]; !ok {
		return ErrIntegrationResponseNotFound
	}

	integration.IntegrationResponses[statusCode] = response
	return s.updateLocked(api)
}

// PutMethodResponse creates or updates a method response for a method.
func (s *RestApiStore) PutMethodResponse(apiId, resourceId, httpMethod, statusCode string, response *MethodResponse) (*MethodResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, _, method, err := s.methodForUpdate(apiId, resourceId, httpMethod)
	if err != nil {
		return nil, err
	}

	if method.MethodResponses == nil {
		method.MethodResponses = make(map[string]*MethodResponse)
	}
	// Normalise the nested maps so a response created without them never
	// persists nil maps — later patch appliers write into them.
	if response.ResponseParameters == nil {
		response.ResponseParameters = make(map[string]bool)
	}
	if response.ResponseModels == nil {
		response.ResponseModels = make(map[string]string)
	}
	method.MethodResponses[statusCode] = response

	if err := s.updateLocked(api); err != nil {
		return nil, err
	}
	return response, nil
}

// GetMethodResponse retrieves a method response for a method.
func (s *RestApiStore) GetMethodResponse(apiId, resourceId, httpMethod, statusCode string) (*MethodResponse, error) {
	method, err := s.GetMethod(apiId, resourceId, httpMethod)
	if err != nil {
		return nil, err
	}

	response, ok := method.MethodResponses[statusCode]
	if !ok {
		return nil, ErrMethodResponseNotFound
	}
	return response, nil
}

// DeleteMethodResponse deletes a method response from a method.
func (s *RestApiStore) DeleteMethodResponse(apiId, resourceId, httpMethod, statusCode string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, _, method, err := s.methodForUpdate(apiId, resourceId, httpMethod)
	if err != nil {
		return err
	}

	if _, ok := method.MethodResponses[statusCode]; !ok {
		return ErrMethodResponseNotFound
	}

	delete(method.MethodResponses, statusCode)
	return s.updateLocked(api)
}
