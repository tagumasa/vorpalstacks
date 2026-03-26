package apigateway

import "strings"

// PutMethod creates or updates a method for an API resource.
func (s *RestApiStore) PutMethod(apiId, resourceId string, method *Method) (*Method, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return nil, err
	}

	resource, ok := api.Resources[resourceId]
	if !ok {
		return nil, ErrResourceNotFound
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

	if err := s.Update(api); err != nil {
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

	api, err := s.Get(apiId)
	if err != nil {
		return err
	}

	resource, ok := api.Resources[resourceId]
	if !ok {
		return ErrResourceNotFound
	}

	httpMethod = strings.ToUpper(httpMethod)
	if _, ok := resource.ResourceMethods[httpMethod]; !ok {
		return ErrMethodNotFound
	}

	delete(resource.ResourceMethods, httpMethod)
	return s.Update(api)
}

// PutIntegration creates or updates an integration for a method.
func (s *RestApiStore) PutIntegration(apiId, resourceId, httpMethod string, integration *Integration) (*Integration, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return nil, err
	}

	resource, ok := api.Resources[resourceId]
	if !ok {
		return nil, ErrResourceNotFound
	}

	httpMethod = strings.ToUpper(httpMethod)
	method, ok := resource.ResourceMethods[httpMethod]
	if !ok {
		return nil, ErrMethodNotFound
	}

	integration.RestApiId = apiId
	integration.ResourceId = resourceId
	integration.HttpMethod = httpMethod
	if integration.IntegrationResponses == nil {
		integration.IntegrationResponses = make(map[string]*IntegrationResponse)
	}

	method.MethodIntegration = integration

	return integration, s.Update(api)
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

	api, err := s.Get(apiId)
	if err != nil {
		return err
	}

	resource, ok := api.Resources[resourceId]
	if !ok {
		return ErrResourceNotFound
	}

	httpMethod = strings.ToUpper(httpMethod)
	method, ok := resource.ResourceMethods[httpMethod]
	if !ok {
		return ErrMethodNotFound
	}

	method.MethodIntegration = nil
	return s.Update(api)
}

// UpdateIntegration updates an integration for a method.
func (s *RestApiStore) UpdateIntegration(apiId, resourceId, httpMethod string, integration *Integration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return err
	}

	resource, ok := api.Resources[resourceId]
	if !ok {
		return ErrResourceNotFound
	}

	httpMethod = strings.ToUpper(httpMethod)
	method, ok := resource.ResourceMethods[httpMethod]
	if !ok {
		return ErrMethodNotFound
	}

	if method.MethodIntegration == nil {
		return ErrIntegrationNotFound
	}

	integration.RestApiId = apiId
	integration.ResourceId = resourceId
	integration.HttpMethod = httpMethod
	method.MethodIntegration = integration
	return s.Update(api)
}

// PutIntegrationResponse creates or updates an integration response for a method.
func (s *RestApiStore) PutIntegrationResponse(apiId, resourceId, httpMethod, statusCode string, response *IntegrationResponse) (*IntegrationResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return nil, err
	}

	resource, ok := api.Resources[resourceId]
	if !ok {
		return nil, ErrResourceNotFound
	}

	httpMethod = strings.ToUpper(httpMethod)
	method, ok := resource.ResourceMethods[httpMethod]
	if !ok {
		return nil, ErrMethodNotFound
	}

	if method.MethodIntegration == nil {
		return nil, ErrIntegrationNotFound
	}

	if method.MethodIntegration.IntegrationResponses == nil {
		method.MethodIntegration.IntegrationResponses = make(map[string]*IntegrationResponse)
	}
	method.MethodIntegration.IntegrationResponses[statusCode] = response

	return response, s.Update(api)
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

	api, err := s.Get(apiId)
	if err != nil {
		return err
	}

	resource, ok := api.Resources[resourceId]
	if !ok {
		return ErrResourceNotFound
	}

	httpMethod = strings.ToUpper(httpMethod)
	method, ok := resource.ResourceMethods[httpMethod]
	if !ok {
		return ErrMethodNotFound
	}

	if method.MethodIntegration == nil {
		return ErrIntegrationNotFound
	}

	delete(method.MethodIntegration.IntegrationResponses, statusCode)
	return s.Update(api)
}

// UpdateIntegrationResponse updates an integration response for a method.
func (s *RestApiStore) UpdateIntegrationResponse(apiId, resourceId, httpMethod, statusCode string, response *IntegrationResponse) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return err
	}

	resource, ok := api.Resources[resourceId]
	if !ok {
		return ErrResourceNotFound
	}

	httpMethod = strings.ToUpper(httpMethod)
	method, ok := resource.ResourceMethods[httpMethod]
	if !ok {
		return ErrMethodNotFound
	}

	if method.MethodIntegration == nil {
		return ErrIntegrationNotFound
	}

	if _, ok := method.MethodIntegration.IntegrationResponses[statusCode]; !ok {
		return ErrIntegrationResponseNotFound
	}

	method.MethodIntegration.IntegrationResponses[statusCode] = response
	return s.Update(api)
}

// PutMethodResponse creates or updates a method response for a method.
func (s *RestApiStore) PutMethodResponse(apiId, resourceId, httpMethod, statusCode string, response *MethodResponse) (*MethodResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return nil, err
	}

	resource, ok := api.Resources[resourceId]
	if !ok {
		return nil, ErrResourceNotFound
	}

	httpMethod = strings.ToUpper(httpMethod)
	method, ok := resource.ResourceMethods[httpMethod]
	if !ok {
		return nil, ErrMethodNotFound
	}

	if method.MethodResponses == nil {
		method.MethodResponses = make(map[string]*MethodResponse)
	}
	method.MethodResponses[statusCode] = response

	if err := s.Update(api); err != nil {
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

	api, err := s.Get(apiId)
	if err != nil {
		return err
	}

	resource, ok := api.Resources[resourceId]
	if !ok {
		return ErrResourceNotFound
	}

	httpMethod = strings.ToUpper(httpMethod)
	method, ok := resource.ResourceMethods[httpMethod]
	if !ok {
		return ErrMethodNotFound
	}

	if _, ok := method.MethodResponses[statusCode]; !ok {
		return ErrMethodResponseNotFound
	}

	delete(method.MethodResponses, statusCode)
	return s.Update(api)
}
