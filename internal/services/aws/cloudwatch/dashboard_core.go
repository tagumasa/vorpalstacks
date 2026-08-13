package cloudwatch

import (
	"encoding/json"

	awserrors "vorpalstacks/internal/common/errors"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
	"vorpalstacks/internal/store/aws/common"
)

// PutDashboardInput holds parameters for PutDashboard in a
// transport-agnostic format.
type PutDashboardInput struct {
	DashboardName string
	DashboardBody string
	Tags          map[string]string
}

// GetDashboardInput holds parameters for GetDashboard.
type GetDashboardInput struct {
	DashboardName string
}

// ListDashboardsInput holds parameters for ListDashboards.
type ListDashboardsInput struct {
	DashboardNamePrefix string
	NextToken           string
}

// DeleteDashboardsInput holds parameters for DeleteDashboards.
type DeleteDashboardsInput struct {
	DashboardNames []string
}

// putDashboardCore validates input, stores the dashboard, and returns
// the dashboard ARN.
func (s *CloudWatchService) putDashboardCore(stores *cloudwatchStores, input *PutDashboardInput) (string, error) {
	if err := validateDashboardName(input.DashboardName); err != nil {
		return "", err
	}
	if input.DashboardBody == "" {
		return "", ErrDashboardInvalidInput
	}
	if !json.Valid([]byte(input.DashboardBody)) {
		return "", awserrors.NewInvalidParameterValueException(
			"DashboardBody must be valid JSON")
	}

	dashboard, err := stores.dashboards.PutDashboard(input.DashboardName, input.DashboardBody, input.Tags)
	if err != nil {
		return "", err
	}
	return dashboard.ARN, nil
}

// getDashboardCore retrieves a dashboard by name.
func (s *CloudWatchService) getDashboardCore(stores *cloudwatchStores, input *GetDashboardInput) (*cwstore.Dashboard, error) {
	if input.DashboardName == "" {
		return nil, ErrInvalidParameter
	}
	dashboard, err := stores.dashboards.GetDashboard(input.DashboardName)
	if err != nil {
		return nil, ErrResourceNotFound
	}
	return dashboard, nil
}

// listDashboardsCore lists dashboards with optional prefix filter.
func (s *CloudWatchService) listDashboardsCore(stores *cloudwatchStores, input *ListDashboardsInput) ([]*cwstore.Dashboard, string, error) {
	opts := common.ListOptions{Marker: input.NextToken, MaxItems: 1000}
	result, err := stores.dashboards.ListDashboardsPaginated(input.DashboardNamePrefix, opts)
	if err != nil {
		return nil, "", err
	}
	return result.Items, result.NextMarker, nil
}

// deleteDashboardsCore deletes dashboards and returns a list of any
// names that were not found.
func (s *CloudWatchService) deleteDashboardsCore(stores *cloudwatchStores, input *DeleteDashboardsInput) ([]string, error) {
	if len(input.DashboardNames) == 0 {
		return nil, ErrInvalidParameter
	}
	if len(input.DashboardNames) > maxDashboardNames {
		return nil, awserrors.NewInvalidParameterValueException(
			"Number of DashboardNames must not exceed 100")
	}

	_, notFound, err := stores.dashboards.DeleteDashboards(input.DashboardNames)
	if err != nil {
		return nil, err
	}
	return notFound, nil
}
