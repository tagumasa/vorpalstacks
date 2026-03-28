package cloudwatch

import (
	"context"

	"vorpalstacks/internal/services/aws/common/request"
)

// PutDashboard creates or updates a CloudWatch dashboard.
func (s *CloudWatchService) PutDashboard(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "DashboardName")
	body := request.GetStringParam(req.Parameters, "DashboardBody")
	if name == "" || body == "" {
		return nil, ErrInvalidParameter
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	dashboard, err := stores.dashboards.PutDashboard(name, body)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"DashboardArn": dashboard.ARN,
	}, nil
}

// GetDashboard retrieves a CloudWatch dashboard.
func (s *CloudWatchService) GetDashboard(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "DashboardName")
	if name == "" {
		return nil, ErrInvalidParameter
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	dashboard, err := stores.dashboards.GetDashboard(name)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	return map[string]interface{}{
		"DashboardArn":  dashboard.ARN,
		"DashboardBody": dashboard.Body,
		"LastModified":  dashboard.UpdatedAt,
	}, nil
}

// ListDashboards lists CloudWatch dashboards, optionally filtered by name prefix.
func (s *CloudWatchService) ListDashboards(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	prefix := request.GetStringParam(req.Parameters, "DashboardNamePrefix")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	dashboards, err := stores.dashboards.ListDashboards(prefix)
	if err != nil {
		return nil, err
	}

	entries := make([]map[string]interface{}, 0, len(dashboards))
	for _, d := range dashboards {
		entries = append(entries, map[string]interface{}{
			"DashboardName": d.Name,
			"DashboardArn":  d.ARN,
			"LastModified":  d.UpdatedAt,
			"Size":          len(d.Body),
		})
	}
	if entries == nil {
		entries = []map[string]interface{}{}
	}

	return map[string]interface{}{
		"DashboardEntries": entries,
	}, nil
}

// DeleteDashboards deletes one or more CloudWatch dashboards.
func (s *CloudWatchService) DeleteDashboards(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var names []string
	if namesRaw, ok := req.Parameters["DashboardNames"]; ok {
		if namesList, ok := namesRaw.([]interface{}); ok {
			for _, n := range namesList {
				if ns, ok := n.(string); ok {
					names = append(names, ns)
				}
			}
		}
	}
	if len(names) == 0 {
		return nil, ErrInvalidParameter
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	_, notFound, err := stores.dashboards.DeleteDashboards(names)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{}
	if len(notFound) > 0 {
		result["Errors"] = notFound
	}
	return result, nil
}
