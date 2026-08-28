package neptune

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/tags"
	rdssvc "vorpalstacks/internal/services/aws/rds"
	neptunestore "vorpalstacks/internal/store/aws/rds/neptune"
	"vorpalstacks/internal/utils/timeutils"
)

// AddTagsToResource adds metadata tags to the specified Neptune resource.
func (s *NeptuneService) AddTagsToResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tags.HandleTag(ctx, req, s.neptuneTagConfig(store))
}

// ListTagsForResource returns the tags attached to the specified Neptune resource.
func (s *NeptuneService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tags.HandleList(ctx, req, s.neptuneTagConfig(store))
}

// RemoveTagsFromResource removes metadata tags from the specified Neptune resource.
func (s *NeptuneService) RemoveTagsFromResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tags.HandleUntag(ctx, req, s.neptuneTagConfig(store))
}

// CreateDBClusterEndpoint creates a new custom endpoint for the specified DB
// cluster.
func (s *NeptuneService) CreateDBClusterEndpoint(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &CreateDBClusterEndpointInput{
		DBClusterEndpointIdentifier: request.GetStringParam(params, "DBClusterEndpointIdentifier"),
		DBClusterIdentifier:         request.GetStringParam(params, "DBClusterIdentifier"),
		EndpointType:                request.GetStringParam(params, "EndpointType"),
		ExcludedMembers:             request.GetStringList(params, "ExcludedMembers"),
		StaticMembers:               request.GetStringList(params, "StaticMembers"),
		AccountID:                   reqCtx.GetAccountID(),
		Region:                      reqCtx.GetRegion(),
	}
	ep, err := s.createDBClusterEndpointCore(ctx, store, in)
	if err != nil {
		return nil, err
	}
	return clusterEndpointToResponse(ep), nil
}

// DescribeDBClusterEndpoints returns information about the custom endpoints for
// the specified DB cluster.
func (s *NeptuneService) DescribeDBClusterEndpoints(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &DescribeDBClusterEndpointsInput{
		DBClusterIdentifier:         request.GetStringParam(params, "DBClusterIdentifier"),
		DBClusterEndpointIdentifier: request.GetStringParam(params, "DBClusterEndpointIdentifier"),
	}
	result, err := s.describeDBClusterEndpointsCore(ctx, store, in)
	if err != nil {
		return nil, err
	}

	if result.Endpoint != nil {
		ep := result.Endpoint
		return map[string]interface{}{
			"DBClusterEndpoints": protocol.XMLElements{ElementName: "DBClusterEndpointList", Items: []interface{}{
				clusterEndpointToResponse(ep),
			}},
		}, nil
	}

	items := make([]interface{}, 0, len(result.Endpoints))
	for _, ep := range result.Endpoints {
		items = append(items, clusterEndpointToResponse(ep))
	}

	return map[string]interface{}{
		"DBClusterEndpoints": protocol.XMLElements{ElementName: "DBClusterEndpointList", Items: items},
	}, nil
}

// clusterEndpointToResponse projects a cluster endpoint record onto the
// operation's response member shape, emitting every documented member the
// record carries.
func clusterEndpointToResponse(ep *neptunestore.DBClusterEndpoint) map[string]interface{} {
	m := map[string]interface{}{
		"DBClusterEndpointIdentifier": ep.DBClusterEndpointIdentifier,
		"DBClusterIdentifier":         ep.DBClusterIdentifier,
		"Endpoint":                    ep.Endpoint,
		"Status":                      ep.Status,
		"EndpointType":                ep.EndpointType,
		"DBClusterEndpointArn":        ep.DBClusterEndpointArn,
	}
	if len(ep.ExcludedMembers) > 0 {
		m["ExcludedMembers"] = protocol.XMLElements{ElementName: "member", Items: stringSliceToInterface(ep.ExcludedMembers)}
	}
	if len(ep.StaticMembers) > 0 {
		m["StaticMembers"] = protocol.XMLElements{ElementName: "member", Items: stringSliceToInterface(ep.StaticMembers)}
	}
	return m
}

// ModifyDBClusterEndpoint modifies the properties of the specified DB cluster
// endpoint.
func (s *NeptuneService) ModifyDBClusterEndpoint(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &ModifyDBClusterEndpointInput{
		DBClusterEndpointIdentifier: request.GetStringParam(params, "DBClusterEndpointIdentifier"),
		EndpointType:                request.GetStringParam(params, "EndpointType"),
		ExcludedMembers:             request.GetStringList(params, "ExcludedMembers"),
		StaticMembers:               request.GetStringList(params, "StaticMembers"),
	}
	ep, err := s.modifyDBClusterEndpointCore(ctx, store, in)
	if err != nil {
		return nil, err
	}
	return clusterEndpointToResponse(ep), nil
}

// DeleteDBClusterEndpoint deletes the specified custom DB cluster endpoint.
func (s *NeptuneService) DeleteDBClusterEndpoint(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &DeleteDBClusterEndpointInput{
		DBClusterEndpointIdentifier: request.GetStringParam(req.Parameters, "DBClusterEndpointIdentifier"),
	}
	ep, err := s.deleteDBClusterEndpointCore(ctx, store, in)
	if err != nil {
		return nil, err
	}
	resp := clusterEndpointToResponse(ep)
	resp["Status"] = "deleting"
	return resp, nil
}

// DescribeEvents returns Neptune events matching the specified filter criteria.
func (s *NeptuneService) DescribeEvents(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &DescribeEventsInput{
		SourceType:       request.GetStringParam(params, "SourceType"),
		SourceIdentifier: request.GetStringParam(params, "SourceIdentifier"),
		StartTime:        request.GetStringParam(params, "StartTime"),
		EndTime:          request.GetStringParam(params, "EndTime"),
		Duration:         request.GetIntParam(params, "Duration"),
		Marker:           pagination.GetMarker(params, "Marker"),
		MaxRecords:       request.GetIntParam(params, "MaxRecords"),
	}
	result, err := s.describeEventsCore(ctx, store, in)
	if err != nil {
		return nil, err
	}

	filtered := make([]interface{}, 0, len(result.Events))
	for _, evt := range result.Events {
		filtered = append(filtered, map[string]interface{}{
			"Date":             evt.Date.UTC().Format(timeutils.ISO8601UTCFormat),
			"EventCategories":  protocol.XMLElements{ElementName: "EventCategory", Items: stringSliceToInterface(evt.EventCategories)},
			"Message":          evt.Message,
			"SourceArn":        evt.SourceArn,
			"SourceIdentifier": evt.SourceIdentifier,
			"SourceType":       evt.SourceType,
		})
	}

	resp := map[string]interface{}{
		"Events": protocol.XMLElements{ElementName: "Event", Items: filtered},
	}
	pagination.SetNextToken(resp, "Marker", result.Marker)
	return resp, nil
}

func stringSliceToInterface(s []string) []interface{} {
	if s == nil {
		return nil
	}
	result := make([]interface{}, len(s))
	for i, v := range s {
		result[i] = v
	}
	return result
}

// DescribeEventCategories returns the available event categories for Neptune
// source types.
func (s *NeptuneService) DescribeEventCategories(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{
		"EventCategoriesMapList": protocol.XMLElements{ElementName: "EventCategoriesMap", Items: []interface{}{
			map[string]interface{}{
				"SourceType": "db-cluster",
				"EventCategories": protocol.XMLElements{ElementName: "EventCategory", Items: []interface{}{
					"creation", "deletion", "failover", "failure",
					"maintenance", "notification", "read replica",
					"recovery", "restoration", "backup",
				}},
			},
			map[string]interface{}{
				"SourceType": "db-instance",
				"EventCategories": protocol.XMLElements{ElementName: "EventCategory", Items: []interface{}{
					"creation", "deletion", "failure",
					"maintenance", "notification", "recovery",
				}},
			},
			map[string]interface{}{
				"SourceType": "db-snapshot",
				"EventCategories": protocol.XMLElements{ElementName: "EventCategory", Items: []interface{}{
					"creation", "deletion", "restoration",
				}},
			},
		}},
	}, nil
}

// DescribePendingMaintenanceActions returns the pending maintenance actions for
// the specified Neptune resources.
func (s *NeptuneService) DescribePendingMaintenanceActions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{
		"PendingMaintenanceActions": protocol.XMLElements{ElementName: "ResourcePendingMaintenanceAction", Items: []interface{}{}},
	}, nil
}

// ApplyPendingMaintenanceAction applies a pending maintenance action to the
// specified Neptune resource.
func (s *NeptuneService) ApplyPendingMaintenanceAction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &ApplyPendingMaintenanceActionInput{
		ResourceIdentifier: request.GetStringParam(params, "ResourceIdentifier"),
		ApplyAction:        request.GetStringParam(params, "ApplyAction"),
		OptInType:          request.GetStringParam(params, "OptInType"),
	}
	if err := s.applyPendingMaintenanceActionCore(ctx, store, in); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ResourcePendingMaintenanceActions": map[string]interface{}{
			"ResourceIdentifier": in.ResourceIdentifier,
		},
	}, nil
}

// DescribeDBEngineVersions returns the available engine versions for
// Neptune and MySQL (when the MySQL engine is wired). The version list
// is driven by rdssvc.SupportedNeptuneVersions / SupportedMysqlVersions
// so that the list a client sees here is exactly the list that
// rdssvc.ValidateEngineVersion accepts — no version can be advertised
// but rejected (or vice versa).
func (s *NeptuneService) DescribeDBEngineVersions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	engineFilter := request.GetStringParam(params, "Engine")
	versionFilter := request.GetStringParam(params, "EngineVersion")

	result := make([]interface{}, 0)

	// Neptune versions (always available — this is the Neptune handler).
	if engineFilter == "" || engineFilter == "neptune" {
		for _, v := range rdssvc.SupportedNeptuneVersions() {
			if versionFilter != "" && v.Version != versionFilter {
				continue
			}
			result = append(result, map[string]interface{}{
				"Engine":                             "neptune",
				"EngineVersion":                      v.Version,
				"DBParameterGroupFamily":             v.Family,
				"DBEngineDescription":                "Amazon Neptune",
				"DBEngineVersionDescription":         "Neptune " + v.Version,
				"ValidUpgradeTarget":                 protocol.XMLElements{ElementName: "UpgradeTarget", Items: []interface{}{}},
				"ExportableLogTypes":                 []interface{}{"audit", "slowquery"},
				"SupportsLogExportsToCloudwatchLogs": true,
				"SupportsGlobalDatabases":            false,
				"SupportsParallelQuery":              v.ParallelQuery,
				"SupportsReadReplica":                true,
				"Status":                             "available",
			})
		}
	}

	// MySQL versions (only when the MySQL engine is wired).
	if (engineFilter == "" || engineFilter == "mysql") && s.mysqlEngine != nil {
		for _, v := range rdssvc.SupportedMysqlVersions() {
			if versionFilter != "" && v.Version != versionFilter {
				continue
			}
			result = append(result, map[string]interface{}{
				"Engine":                             "mysql",
				"EngineVersion":                      v.Version,
				"DBParameterGroupFamily":             v.Family,
				"DBEngineDescription":                "MySQL",
				"DBEngineVersionDescription":         v.DescShort,
				"ValidUpgradeTarget":                 protocol.XMLElements{ElementName: "UpgradeTarget", Items: []interface{}{}},
				"ExportableLogTypes":                 []interface{}{"error", "general", "slowquery"},
				"SupportsLogExportsToCloudwatchLogs": true,
				"SupportsGlobalDatabases":            false,
				"SupportsParallelQuery":              false,
				"SupportsReadReplica":                true,
				"Status":                             "available",
			})
		}
	}

	return map[string]interface{}{
		"DBEngineVersions": protocol.XMLElements{ElementName: "DBEngineVersion", Items: result},
	}, nil
}

// DescribeOrderableDBInstanceOptions returns the available DB instance classes
// and configurations for Neptune.
func (s *NeptuneService) DescribeOrderableDBInstanceOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	engine := request.GetStringParam(params, "Engine")
	if engine == "" {
		engine = "neptune"
	}
	if engine != "neptune" {
		return map[string]interface{}{
			"OrderableDBInstanceOptions": protocol.XMLElements{ElementName: "OrderableDBInstanceOption", Items: []interface{}{}},
		}, nil
	}
	region := reqCtx.GetRegion()

	classes := []string{"db.r5.large", "db.r5.xlarge", "db.r5.2xlarge", "db.r5.4xlarge", "db.r6g.large", "db.r6g.xlarge", "db.r6g.2xlarge", "db.t3.medium"}

	options := make([]interface{}, 0, len(classes))
	for _, cls := range classes {
		options = append(options, map[string]interface{}{
			"Engine":          engine,
			"EngineVersion":   rdssvc.DefaultEngineVersion("neptune"),
			"DBInstanceClass": cls,
			"LicenseModel":    "bring-your-own-license",
			"AvailabilityZones": protocol.XMLElements{ElementName: "AvailabilityZone", Items: []interface{}{
				map[string]interface{}{"Name": region + "a"},
				map[string]interface{}{"Name": region + "b"},
				map[string]interface{}{"Name": region + "c"},
			}},
			"MultiAZCapable":            true,
			"ReadReplicaCapable":        true,
			"StorageType":               "standard",
			"SupportsStorageEncryption": true,
		})
	}

	return map[string]interface{}{
		"OrderableDBInstanceOptions": protocol.XMLElements{ElementName: "OrderableDBInstanceOption", Items: options},
	}, nil
}

// DescribeValidDBInstanceModifications returns the valid modifications that can
// be applied to the specified DB instance.
func (s *NeptuneService) DescribeValidDBInstanceModifications(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &DescribeValidDBInstanceModificationsInput{
		DBInstanceIdentifier: request.GetStringParam(req.Parameters, "DBInstanceIdentifier"),
	}
	return s.describeValidDBInstanceModificationsCore(ctx, store, in)
}
