package neptune

import (
	"context"
	"fmt"
	"net/http"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/tags"
	neptunestore "vorpalstacks/internal/store/aws/neptune"
	arnutil "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/aws/types"
)

func (s *NeptuneService) neptuneTagConfig(store neptunestore.NeptuneStoreInterface) tags.TagHandlerConfig {
	return tags.TagHandlerConfig{
		Param: tags.TagOperationConfig{
			ResourceParam:  "ResourceName",
			TagsParam:      "Tags",
			TagKeysParam:   "TagKeys",
			TagKeyName:     "Key",
			TagValueName:   "Value",
			RequireTags:    false,
			RequireTagKeys: false,
			RequireResource: true,
		},
		ParseTags: func(params map[string]interface{}) []types.Tag {
			rawTags := getNeptuneTagList(params)
			result := make([]types.Tag, 0, len(rawTags))
			for _, t := range rawTags {
				key, _ := t["Key"].(string)
				value, _ := t["Value"].(string)
				if key != "" {
					result = append(result, types.Tag{Key: key, Value: value})
				}
			}
			return result
		},
		ParseTagKeys: func(params map[string]interface{}) []string {
			return request.GetStringList(params, "TagKeys")
		},
		TagFunc: func(ctx context.Context, resourceKey string, tagList []types.Tag) error {
			return store.AddTags(resourceKey, tagList)
		},
		UntagFunc: func(ctx context.Context, resourceKey string, tagKeys []string) error {
			return store.RemoveTags(resourceKey, tagKeys)
		},
		ListFunc: func(ctx context.Context, resourceKey string) ([]types.Tag, error) {
			return store.GetTags(resourceKey)
		},
		FormatResponse: func(tagList []types.Tag, _ string) (interface{}, error) {
			items := make([]interface{}, 0, len(tagList))
			for _, t := range tagList {
				items = append(items, map[string]interface{}{"Key": t.Key, "Value": t.Value})
			}
			return map[string]interface{}{
				"TagList": protocol.XMLElements{ElementName: "Tag", Items: items},
			}, nil
		},
		EmptyResponse: func() (interface{}, error) {
			return map[string]interface{}{}, nil
		},
		MapError: func(err error) error {
			return err
		},
	}
}

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
	epID := request.GetStringParam(params, "DBClusterEndpointIdentifier")
	if epID == "" {
		return nil, fmt.Errorf("neptune: DBClusterEndpointIdentifier is required")
	}
	clusterID := request.GetStringParam(params, "DBClusterIdentifier")
	if clusterID == "" {
		return nil, fmt.Errorf("neptune: DBClusterIdentifier is required")
	}
	epType := request.GetStringParam(params, "EndpointType")
	if epType == "" {
		return nil, fmt.Errorf("neptune: EndpointType is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetCluster(clusterID); err != nil {
		return nil, fmt.Errorf("neptune: DBCluster %s not found", clusterID)
	}

	accountID := reqCtx.GetAccountID()
	region := reqCtx.GetRegion()

	excludedMembers := request.GetStringList(params, "ExcludedMembers")
	staticMembers := request.GetStringList(params, "StaticMembers")

	ep := &neptunestore.DBClusterEndpoint{
		DBClusterEndpointIdentifier: epID,
		DBClusterIdentifier:         clusterID,
		Endpoint:                    fmt.Sprintf("%s.cluster-%s.%s.amazonaws.com", epID, clusterID, region),
		Status:                      "available",
		EndpointType:                epType,
		ExcludedMembers:             excludedMembers,
		StaticMembers:               staticMembers,
		DBClusterEndpointArn:        arnutil.NewARNBuilder(accountID, region).RDS().ClusterEndpoint(epID),
	}

	if err := store.CreateClusterEndpoint(ep); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"DBClusterEndpointIdentifier": epID,
		"DBClusterIdentifier":         clusterID,
		"Endpoint":                    ep.Endpoint,
		"EndpointType":                epType,
		"Status":                      "available",
		"DBClusterEndpointArn":        ep.DBClusterEndpointArn,
	}, nil
}

// DescribeDBClusterEndpoints returns information about the custom endpoints for
// the specified DB cluster.
func (s *NeptuneService) DescribeDBClusterEndpoints(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	clusterID := request.GetStringParam(params, "DBClusterIdentifier")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	endpoints, err := store.ListClusterEndpoints(clusterID)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(endpoints))
	for _, ep := range endpoints {
		items = append(items, map[string]interface{}{
			"DBClusterEndpointIdentifier": ep.DBClusterEndpointIdentifier,
			"DBClusterIdentifier":         ep.DBClusterIdentifier,
			"Endpoint":                    ep.Endpoint,
			"Status":                      ep.Status,
			"EndpointType":                ep.EndpointType,
			"DBClusterEndpointArn":        ep.DBClusterEndpointArn,
		})
	}

	return map[string]interface{}{
		"DBClusterEndpoints": protocol.XMLElements{ElementName: "DBClusterEndpointList", Items: items},
	}, nil
}

// ModifyDBClusterEndpoint modifies the properties of the specified DB cluster
// endpoint.
func (s *NeptuneService) ModifyDBClusterEndpoint(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	epID := request.GetStringParam(params, "DBClusterEndpointIdentifier")
	if epID == "" {
		return nil, fmt.Errorf("neptune: DBClusterEndpointIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	ep, err := store.GetClusterEndpoint(epID)
	if err != nil {
		return nil, err
	}

	if request.HasParam(params, "ExcludedMembers") {
		ep.ExcludedMembers = request.GetStringList(params, "ExcludedMembers")
	}
	if request.HasParam(params, "StaticMembers") {
		ep.StaticMembers = request.GetStringList(params, "StaticMembers")
	}
	if newType := request.GetStringParam(params, "EndpointType"); newType != "" {
		ep.EndpointType = newType
	}

	ep.Status = "available"
	if err := store.UpdateClusterEndpoint(ep); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"DBClusterEndpointIdentifier": epID,
		"DBClusterIdentifier":         ep.DBClusterIdentifier,
		"Endpoint":                    ep.Endpoint,
		"EndpointType":                ep.EndpointType,
		"Status":                      "available",
		"DBClusterEndpointArn":        ep.DBClusterEndpointArn,
	}, nil
}

// DeleteDBClusterEndpoint deletes the specified custom DB cluster endpoint.
func (s *NeptuneService) DeleteDBClusterEndpoint(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	epID := request.GetStringParam(params, "DBClusterEndpointIdentifier")
	if epID == "" {
		return nil, fmt.Errorf("neptune: DBClusterEndpointIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	ep, err := store.GetClusterEndpoint(epID)
	if err != nil {
		return nil, err
	}
	epArn := ep.DBClusterEndpointArn

	if err := store.DeleteClusterEndpoint(epID); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"DBClusterEndpointIdentifier": epID,
		"Status":                      "deleting",
		"DBClusterEndpointArn":        epArn,
	}, nil
}

// DescribeEvents returns Neptune events matching the specified filter criteria.
func (s *NeptuneService) DescribeEvents(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var startTime time.Time
	if stStr := request.GetStringParam(params, "StartTime"); stStr != "" {
		var err error
		startTime, err = time.Parse(time.RFC3339, stStr)
		if err != nil {
			return nil, awserrors.NewAWSError("InvalidParameterValue", "Invalid StartTime format: use RFC3339", http.StatusBadRequest)
		}
	}

	duration := request.GetIntParam(params, "Duration")
	var endTime time.Time
	if duration > 0 && !startTime.IsZero() {
		endTime = startTime.Add(time.Duration(duration) * time.Minute)
	} else if etStr := request.GetStringParam(params, "EndTime"); etStr != "" {
		var err error
		endTime, err = time.Parse(time.RFC3339, etStr)
		if err != nil {
			return nil, awserrors.NewAWSError("InvalidParameterValue", "Invalid EndTime format: use RFC3339", http.StatusBadRequest)
		}
	}

	opts := neptunestore.EventListOptions{
		SourceType:       request.GetStringParam(params, "SourceType"),
		SourceIdentifier: request.GetStringParam(params, "SourceIdentifier"),
		StartTime:        startTime,
		EndTime:          endTime,
		Marker:           pagination.GetMarker(req.Parameters, "Marker"),
		MaxRecords:       pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems),
	}

	result, err := store.ListEvents(opts)
	if err != nil {
		return nil, err
	}

	filtered := make([]interface{}, 0, len(result.Events))
	for _, evt := range result.Events {
		filtered = append(filtered, map[string]interface{}{
			"Date":             evt.Date.UTC().Format("2006-01-02T15:04:05.000Z"),
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
	resourceID := request.GetStringParam(params, "ResourceIdentifier")
	if resourceID == "" {
		return nil, fmt.Errorf("neptune: ResourceIdentifier is required")
	}
	action := request.GetStringParam(params, "ApplyAction")
	if action == "" {
		return nil, fmt.Errorf("neptune: ApplyAction is required")
	}
	optIn := request.GetStringParam(params, "OptInType")
	if optIn == "" {
		return nil, fmt.Errorf("neptune: OptInType is required")
	}

	return map[string]interface{}{
		"ResourcePendingMaintenanceActions": map[string]interface{}{
			"ResourceIdentifier": resourceID,
		},
	}, nil
}

// DescribeDBEngineVersions returns the available Neptune engine versions.
func (s *NeptuneService) DescribeDBEngineVersions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	engine := request.GetStringParam(params, "Engine")
	if engine == "" {
		engine = "neptune"
	}
	if engine != "neptune" {
		return map[string]interface{}{
			"DBEngineVersions": protocol.XMLElements{ElementName: "DBEngineVersion", Items: []interface{}{}},
		}, nil
	}

	versions := []map[string]interface{}{
		{
			"Engine":                             engine,
			"EngineVersion":                      "1.3.2.0",
			"DBParameterGroupFamily":             "neptune1",
			"DBEngineDescription":                "Amazon Neptune",
			"DBEngineVersionDescription":         "Neptune 1.3.2.0",
			"ValidUpgradeTarget":                 protocol.XMLElements{ElementName: "UpgradeTarget", Items: []interface{}{}},
			"ExportableLogTypes":                 []interface{}{"audit", "slowquery"},
			"SupportsLogExportsToCloudwatchLogs": true,
			"SupportsGlobalDatabases":            false,
			"SupportsParallelQuery":              true,
			"SupportsReadReplica":                true,
			"Status":                             "available",
		},
		{
			"Engine":                             engine,
			"EngineVersion":                      "1.3.1.0",
			"DBParameterGroupFamily":             "neptune1",
			"DBEngineDescription":                "Amazon Neptune",
			"DBEngineVersionDescription":         "Neptune 1.3.1.0",
			"ValidUpgradeTarget":                 protocol.XMLElements{ElementName: "UpgradeTarget", Items: []interface{}{}},
			"ExportableLogTypes":                 []interface{}{"audit", "slowquery"},
			"SupportsLogExportsToCloudwatchLogs": true,
			"SupportsGlobalDatabases":            false,
			"SupportsParallelQuery":              true,
			"SupportsReadReplica":                true,
			"Status":                             "available",
		},
		{
			"Engine":                             engine,
			"EngineVersion":                      "1.2.1.0",
			"DBParameterGroupFamily":             "neptune1",
			"DBEngineDescription":                "Amazon Neptune",
			"DBEngineVersionDescription":         "Neptune 1.2.1.0",
			"ValidUpgradeTarget":                 protocol.XMLElements{ElementName: "UpgradeTarget", Items: []interface{}{}},
			"ExportableLogTypes":                 []interface{}{"audit", "slowquery"},
			"SupportsLogExportsToCloudwatchLogs": true,
			"SupportsGlobalDatabases":            false,
			"SupportsParallelQuery":              false,
			"SupportsReadReplica":                true,
			"Status":                             "available",
		},
	}

	result := make([]interface{}, 0, len(versions))
	for _, v := range versions {
		result = append(result, v)
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
			"EngineVersion":   "1.3.2.0",
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
	params := req.Parameters
	instanceID := request.GetStringParam(params, "DBInstanceIdentifier")
	if instanceID == "" {
		return nil, fmt.Errorf("neptune: DBInstanceIdentifier is required")
	}

	return map[string]interface{}{
		"ValidDBInstanceModificationsMessage": map[string]interface{}{
			"Storage": protocol.XMLElements{ElementName: "ValidStorageOptions", Items: []interface{}{}},
		},
	}, nil
}
