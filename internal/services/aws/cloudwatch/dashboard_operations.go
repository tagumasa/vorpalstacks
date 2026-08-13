package cloudwatch

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
)

// PutDashboard creates or updates a CloudWatch dashboard.
func (s *CloudWatchService) PutDashboard(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "DashboardName")
	body := request.GetStringParam(req.Parameters, "DashboardBody")
	tags := parseAlarmTags(req.Parameters)

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	arn, err := s.putDashboardCore(stores, &PutDashboardInput{
		DashboardName: name,
		DashboardBody: body,
		Tags:          tags,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"DashboardArn": arn,
	}, nil
}

// GetDashboard retrieves a CloudWatch dashboard.
func (s *CloudWatchService) GetDashboard(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "DashboardName")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	dashboard, err := s.getDashboardCore(stores, &GetDashboardInput{DashboardName: name})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"DashboardArn":  dashboard.ARN,
		"DashboardBody": dashboard.Body,
		"DashboardName": dashboard.Name,
	}, nil
}

// ListDashboards lists CloudWatch dashboards, optionally filtered by name prefix.
func (s *CloudWatchService) ListDashboards(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	prefix := request.GetStringParam(req.Parameters, "DashboardNamePrefix")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	marker := pagination.GetMarker(req.Parameters, "NextToken")

	items, nextMarker, err := s.listDashboardsCore(stores, &ListDashboardsInput{
		DashboardNamePrefix: prefix,
		NextToken:           marker,
	})
	if err != nil {
		return nil, err
	}

	entries := make([]map[string]interface{}, 0, len(items))
	for _, d := range items {
		entries = append(entries, map[string]interface{}{
			"DashboardName": d.Name,
			"DashboardArn":  d.ARN,
			"LastModified":  d.UpdatedAt,
			"Size":          len(d.Body),
		})
	}

	resp := map[string]interface{}{
		"DashboardEntries": entries,
	}
	if nextMarker != "" {
		resp["NextToken"] = nextMarker
	}
	return resp, nil
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

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	notFound, err := s.deleteDashboardsCore(stores, &DeleteDashboardsInput{DashboardNames: names})
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{}
	if len(notFound) > 0 {
		result["Errors"] = notFound
	}
	return result, nil
}
