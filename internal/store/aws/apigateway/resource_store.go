package apigateway

import (
	"errors"

	"vorpalstacks/internal/store/aws/common"
)

// CreateResource creates a new resource for a REST API.
func (s *RestApiStore) CreateResource(apiId string, resource *Resource) (*Resource, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return nil, err
	}

	for _, existing := range api.Resources {
		if existing.Path == resource.Path {
			return nil, ErrResourceAlreadyExists
		}
	}

	if resource.ParentId != "" {
		if _, ok := api.Resources[resource.ParentId]; !ok {
			return nil, common.NewStoreError("apigateway", "create_resource", errors.New("parent resource not found"))
		}
	}

	if resource.Id == "" {
		resource.Id = s.arnBuilder.GenerateResourceId()
	}
	resource.RestApiId = apiId
	if resource.ResourceMethods == nil {
		resource.ResourceMethods = make(map[string]*Method)
	}

	api.Resources[resource.Id] = resource

	if err := s.Update(api); err != nil {
		return nil, err
	}

	return resource, nil
}

// GetResource retrieves a resource by its ID.
func (s *RestApiStore) GetResource(apiId, resourceId string) (*Resource, error) {
	api, err := s.Get(apiId)
	if err != nil {
		return nil, err
	}

	resource, ok := api.Resources[resourceId]
	if !ok {
		return nil, ErrResourceNotFound
	}
	return resource, nil
}

// DeleteResource deletes a resource from a REST API.
func (s *RestApiStore) DeleteResource(apiId, resourceId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return err
	}

	if _, ok := api.Resources[resourceId]; !ok {
		return ErrResourceNotFound
	}

	childResourceIds := []string{}
	for _, resource := range api.Resources {
		if resource.ParentId == resourceId {
			childResourceIds = append(childResourceIds, resource.Id)
		}
	}

	if len(childResourceIds) > 0 {
		return common.NewStoreError("apigateway", "delete_resource", errors.New("cannot delete resource with child resources"))
	}

	delete(api.Resources, resourceId)
	return s.Update(api)
}

// ListResources returns all resources for a REST API.
func (s *RestApiStore) ListResources(apiId string) ([]*Resource, error) {
	api, err := s.Get(apiId)
	if err != nil {
		return nil, err
	}

	resources := make([]*Resource, 0, len(api.Resources))
	for _, r := range api.Resources {
		resources = append(resources, r)
	}
	return resources, nil
}

// UpdateResource updates an existing resource.
func (s *RestApiStore) UpdateResource(apiId string, resource *Resource) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return err
	}

	if _, ok := api.Resources[resource.Id]; !ok {
		return ErrResourceNotFound
	}

	resource.RestApiId = apiId
	api.Resources[resource.Id] = resource
	return s.Update(api)
}
