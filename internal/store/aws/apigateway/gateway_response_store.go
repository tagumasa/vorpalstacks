package apigateway

import "sort"

// PutGatewayResponse creates or replaces the gateway response of the given
// response type.
func (s *RestApiStore) PutGatewayResponse(apiId string, resp *GatewayResponse) (*GatewayResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return nil, err
	}
	if api.GatewayResponses == nil {
		api.GatewayResponses = make(map[string]*GatewayResponse)
	}
	api.GatewayResponses[resp.ResponseType] = resp
	if err := s.updateLocked(api); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetGatewayResponse retrieves the gateway response of a response type.
func (s *RestApiStore) GetGatewayResponse(apiId, responseType string) (*GatewayResponse, error) {
	api, err := s.Get(apiId)
	if err != nil {
		return nil, err
	}
	resp, ok := api.GatewayResponses[responseType]
	if !ok {
		return nil, ErrGatewayResponseNotFound
	}
	return resp, nil
}

// DeleteGatewayResponse removes the gateway response of a response type.
func (s *RestApiStore) DeleteGatewayResponse(apiId, responseType string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return err
	}
	if _, ok := api.GatewayResponses[responseType]; !ok {
		return ErrGatewayResponseNotFound
	}
	delete(api.GatewayResponses, responseType)
	return s.updateLocked(api)
}

// ListGatewayResponses returns the API's gateway responses sorted by
// response type.
func (s *RestApiStore) ListGatewayResponses(apiId string) ([]*GatewayResponse, error) {
	api, err := s.Get(apiId)
	if err != nil {
		return nil, err
	}
	responses := make([]*GatewayResponse, 0, len(api.GatewayResponses))
	for _, resp := range api.GatewayResponses {
		responses = append(responses, resp)
	}
	sort.Slice(responses, func(i, j int) bool {
		return responses[i].ResponseType < responses[j].ResponseType
	})
	return responses, nil
}
