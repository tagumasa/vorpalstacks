package sns

import (
	"fmt"
	"strings"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/store/aws/common"
	snsstore "vorpalstacks/internal/store/aws/sns"
)

// validPlatforms lists the push notification platforms accepted by SNS.
var validPlatforms = map[string]bool{
	"APNS":         true,
	"APNS_SANDBOX": true,
	"GCM":          true,
	"ADM":          true,
	"BAIDU":        true,
	"WNS":          true,
	"MPNS":         true,
}

// CreatePlatformApplicationInput carries the parsed parameters for creating a
// platform application.
type CreatePlatformApplicationInput struct {
	Name       string
	Platform   string
	Attributes map[string]string
}

// DeletePlatformApplicationInput carries the ARN of the platform application
// to delete.
type DeletePlatformApplicationInput struct {
	PlatformApplicationArn string
}

// GetPlatformApplicationAttributesInput carries the ARN whose attributes are
// read.
type GetPlatformApplicationAttributesInput struct {
	PlatformApplicationArn string
}

// SetPlatformApplicationAttributesInput carries the ARN and the attribute
// updates for a platform application.
type SetPlatformApplicationAttributesInput struct {
	PlatformApplicationArn string
	Attributes             map[string]string
}

// ListPlatformApplicationsInput carries the pagination token for listing
// platform applications.
type ListPlatformApplicationsInput struct {
	NextToken string
}

// CreatePlatformEndpointInput carries the parsed parameters for creating a
// platform endpoint.
type CreatePlatformEndpointInput struct {
	PlatformApplicationArn string
	Token                  string
	CustomUserData         string
	Attributes             map[string]string
}

// DeleteEndpointInput carries the ARN of the platform endpoint to delete.
type DeleteEndpointInput struct {
	EndpointArn string
}

// GetEndpointAttributesInput carries the endpoint ARN whose attributes are
// read.
type GetEndpointAttributesInput struct {
	EndpointArn string
}

// SetEndpointAttributesInput carries the endpoint ARN and the attribute
// updates.
type SetEndpointAttributesInput struct {
	EndpointArn string
	Attributes  map[string]string
}

// ListEndpointsByPlatformApplicationInput carries the platform application
// ARN and the pagination token for listing its endpoints.
type ListEndpointsByPlatformApplicationInput struct {
	PlatformApplicationArn string
	NextToken              string
}

// createPlatformApplicationCore is the single validation and persistence path
// for CreatePlatformApplication.
func (s *SNSService) createPlatformApplicationCore(store snsstore.SNSStoreInterface, in CreatePlatformApplicationInput) (interface{}, error) {
	if err := validatePlatformApplicationName(in.Name); err != nil {
		return nil, err
	}
	if in.Platform == "" {
		return nil, NewInvalidParameter("Platform is required")
	}
	normalisedPlatform := strings.ToUpper(in.Platform)
	if !validPlatforms[normalisedPlatform] {
		return nil, NewInvalidParameter(fmt.Sprintf("Invalid Platform: %s. Valid values: APNS, APNS_SANDBOX, GCM, ADM, BAIDU, WNS, MPNS", in.Platform))
	}

	for attrName, value := range in.Attributes {
		if err := validatePlatformAttributeValue(attrName, value); err != nil {
			return nil, err
		}
	}

	app := &snsstore.PlatformApplication{
		Name:       in.Name,
		Platform:   normalisedPlatform,
		Attributes: in.Attributes,
	}

	created, err := store.CreatePlatformApplication(app)
	if err != nil {
		if err == snsstore.ErrPlatformApplicationAlreadyExists {
			return nil, NewInvalidParameter("Platform application already exists with the same name")
		}
		return nil, err
	}

	return map[string]interface{}{
		"PlatformApplicationArn": created.PlatformApplicationArn,
	}, nil
}

// deletePlatformApplicationCore is the single validation and persistence path
// for DeletePlatformApplication.
func (s *SNSService) deletePlatformApplicationCore(store snsstore.SNSStoreInterface, in DeletePlatformApplicationInput) (interface{}, error) {
	if in.PlatformApplicationArn == "" {
		return nil, NewInvalidParameter("PlatformApplicationArn is required")
	}

	if err := store.DeletePlatformApplication(in.PlatformApplicationArn); err != nil {
		if err == snsstore.ErrPlatformApplicationNotFound {
			return nil, ErrPlatformAppNotFound
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// getPlatformApplicationAttributesCore is the single validation and
// persistence path for GetPlatformApplicationAttributes.
func (s *SNSService) getPlatformApplicationAttributesCore(store snsstore.SNSStoreInterface, in GetPlatformApplicationAttributesInput) (interface{}, error) {
	if in.PlatformApplicationArn == "" {
		return nil, NewInvalidParameter("PlatformApplicationArn is required")
	}

	attrs, err := store.GetPlatformApplicationAttributes(in.PlatformApplicationArn)
	if err != nil {
		if err == snsstore.ErrPlatformApplicationNotFound {
			return nil, ErrPlatformAppNotFound
		}
		return nil, err
	}

	return map[string]interface{}{
		"Attributes": attrs,
	}, nil
}

// setPlatformApplicationAttributesCore is the single validation and
// persistence path for SetPlatformApplicationAttributes.
func (s *SNSService) setPlatformApplicationAttributesCore(store snsstore.SNSStoreInterface, in SetPlatformApplicationAttributesInput) (interface{}, error) {
	if in.PlatformApplicationArn == "" {
		return nil, NewInvalidParameter("PlatformApplicationArn is required")
	}

	if len(in.Attributes) == 0 {
		return nil, NewInvalidParameter("Attributes is required")
	}

	// Enforce value length cap for DoS protection.
	for name, value := range in.Attributes {
		if err := validatePlatformAttributeValue(name, value); err != nil {
			return nil, err
		}
	}

	if err := store.SetPlatformApplicationAttributes(in.PlatformApplicationArn, in.Attributes); err != nil {
		if err == snsstore.ErrPlatformApplicationNotFound {
			return nil, ErrPlatformAppNotFound
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// listPlatformApplicationsCore is the single validation and persistence path
// for ListPlatformApplications.
func (s *SNSService) listPlatformApplicationsCore(store snsstore.SNSStoreInterface, in ListPlatformApplicationsInput) (interface{}, error) {
	result, err := store.ListPlatformApplications(common.ListOptions{Marker: in.NextToken})
	if err != nil {
		return nil, err
	}

	apps := make([]map[string]interface{}, 0, len(result.Items))
	for _, app := range result.Items {
		apps = append(apps, map[string]interface{}{
			"PlatformApplicationArn": app.PlatformApplicationArn,
			"Attributes":             app.Attributes,
		})
	}

	token := ""
	if result.IsTruncated && result.NextMarker != "" {
		token = result.NextMarker
	}
	return pagination.BuildListResponse("PlatformApplications", apps, token), nil
}

// createPlatformEndpointCore is the single validation and persistence path for
// CreatePlatformEndpoint.
func (s *SNSService) createPlatformEndpointCore(store snsstore.SNSStoreInterface, in CreatePlatformEndpointInput) (interface{}, error) {
	// Validate PlatformApplicationArn format before proceeding.
	if err := validatePlatformApplicationArn(in.PlatformApplicationArn); err != nil {
		return nil, err
	}
	if in.Token == "" {
		return nil, NewInvalidParameter("Token is required")
	}
	if len(in.Token) > 2048 {
		return nil, NewInvalidParameter(fmt.Sprintf("Token too long: %d bytes (maximum 2048)", len(in.Token)))
	}
	if len(in.CustomUserData) > 2048 {
		return nil, NewInvalidParameter(fmt.Sprintf("CustomUserData too long: %d bytes (maximum 2048)", len(in.CustomUserData)))
	}

	for attrName, value := range in.Attributes {
		if err := validatePlatformAttributeValue(attrName, value); err != nil {
			return nil, err
		}
	}

	endpoint := &snsstore.PlatformEndpoint{
		PlatformApplicationArn: in.PlatformApplicationArn,
		Token:                  in.Token,
		CustomUserData:         in.CustomUserData,
		Attributes:             in.Attributes,
	}

	created, err := store.CreatePlatformEndpoint(endpoint)
	if err != nil {
		if err == snsstore.ErrPlatformApplicationNotFound {
			return nil, ErrPlatformAppNotFound
		}
		return nil, err
	}

	return map[string]interface{}{
		"EndpointArn": created.EndpointArn,
	}, nil
}

// deleteEndpointCore is the single validation and persistence path for
// DeleteEndpoint.
func (s *SNSService) deleteEndpointCore(store snsstore.SNSStoreInterface, in DeleteEndpointInput) (interface{}, error) {
	if in.EndpointArn == "" {
		return nil, NewInvalidParameter("EndpointArn is required")
	}

	if err := store.DeleteEndpoint(in.EndpointArn); err != nil {
		if err == snsstore.ErrEndpointNotFound {
			return nil, ErrEndpointNotFound
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// getEndpointAttributesCore is the single validation and persistence path for
// GetEndpointAttributes.
func (s *SNSService) getEndpointAttributesCore(store snsstore.SNSStoreInterface, in GetEndpointAttributesInput) (interface{}, error) {
	if in.EndpointArn == "" {
		return nil, NewInvalidParameter("EndpointArn is required")
	}

	attrs, err := store.GetEndpointAttributes(in.EndpointArn)
	if err != nil {
		if err == snsstore.ErrEndpointNotFound {
			return nil, ErrEndpointNotFound
		}
		return nil, err
	}

	return map[string]interface{}{
		"Attributes": attrs,
	}, nil
}

// setEndpointAttributesCore is the single validation and persistence path for
// SetEndpointAttributes.
func (s *SNSService) setEndpointAttributesCore(store snsstore.SNSStoreInterface, in SetEndpointAttributesInput) (interface{}, error) {
	if in.EndpointArn == "" {
		return nil, NewInvalidParameter("EndpointArn is required")
	}

	if len(in.Attributes) == 0 {
		return nil, NewInvalidParameter("Attributes is required")
	}

	// Enforce value length cap for DoS protection.
	for name, value := range in.Attributes {
		if err := validatePlatformAttributeValue(name, value); err != nil {
			return nil, err
		}
	}

	if err := store.SetEndpointAttributes(in.EndpointArn, in.Attributes); err != nil {
		if err == snsstore.ErrEndpointNotFound {
			return nil, ErrEndpointNotFound
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// listEndpointsByPlatformApplicationCore is the single validation and
// persistence path for ListEndpointsByPlatformApplication.
func (s *SNSService) listEndpointsByPlatformApplicationCore(store snsstore.SNSStoreInterface, in ListEndpointsByPlatformApplicationInput) (interface{}, error) {
	if in.PlatformApplicationArn == "" {
		return nil, NewInvalidParameter("PlatformApplicationArn is required")
	}

	// Verify platform application existence before listing endpoints.
	// Without this check, a non-existent ARN returns an empty result instead
	// of NotFoundException.
	if _, err := store.GetPlatformApplication(in.PlatformApplicationArn); err != nil {
		if err == snsstore.ErrPlatformApplicationNotFound {
			return nil, ErrPlatformAppNotFound
		}
		return nil, err
	}

	result, err := store.ListEndpointsByPlatformApplication(in.PlatformApplicationArn, common.ListOptions{Marker: in.NextToken})
	if err != nil {
		return nil, err
	}

	endpoints := make([]map[string]interface{}, 0, len(result.Items))
	for _, ep := range result.Items {
		endpoints = append(endpoints, map[string]interface{}{
			"EndpointArn": ep.EndpointArn,
			"Attributes":  ep.Attributes,
		})
	}

	token := ""
	if result.IsTruncated && result.NextMarker != "" {
		token = result.NextMarker
	}
	return pagination.BuildListResponse("Endpoints", endpoints, token), nil
}
