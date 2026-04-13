package neptune

import (
	"context"
	"fmt"

	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	neptunestore "vorpalstacks/internal/store/aws/neptune"
)

// CreateDBClusterParameterGroup creates a new DB cluster parameter group.
func (s *NeptuneService) CreateDBClusterParameterGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	name := request.GetStringParam(params, "DBClusterParameterGroupName")
	if name == "" {
		return nil, fmt.Errorf("neptune: DBClusterParameterGroupName is required")
	}
	family := request.GetStringParam(params, "DBParameterGroupFamily")
	if family == "" {
		family = "neptune1"
	}
	desc := request.GetStringParam(params, "Description")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	pg := &neptunestore.DBClusterParameterGroup{
		DBClusterParameterGroupName: name,
		DBParameterGroupFamily:      family,
		Description:                 desc,
		ARN:                         neptunestore.ClusterParameterGroupARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), name),
	}

	if err := store.CreateClusterParameterGroup(pg); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"DBClusterParameterGroup": pg,
	}, nil
}

// DeleteDBClusterParameterGroup deletes the specified DB cluster parameter group.
func (s *NeptuneService) DeleteDBClusterParameterGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	name := request.GetStringParam(params, "DBClusterParameterGroupName")
	if name == "" {
		return nil, fmt.Errorf("neptune: DBClusterParameterGroupName is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteClusterParameterGroup(name); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{}, nil
}

// DescribeDBClusterParameterGroups returns information about the specified
// cluster parameter group or lists all groups when no name is provided.
func (s *NeptuneService) DescribeDBClusterParameterGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	name := request.GetStringParam(params, "DBClusterParameterGroupName")
	if name != "" {
		pg, err := store.GetClusterParameterGroup(name)
		if err != nil {
			return nil, translateStoreError(err)
		}
		return map[string]interface{}{
			"DBClusterParameterGroups": protocol.XMLElements{ElementName: "DBClusterParameterGroup", Items: []interface{}{pg}},
		}, nil
	}

	groups, err := store.ListClusterParameterGroups()
	if err != nil {
		return nil, translateStoreError(err)
	}

	items := make([]interface{}, 0, len(groups))
	for _, g := range groups {
		items = append(items, g)
	}

	marker := request.GetStringParam(params, "Marker")
	maxRecords := request.GetIntParam(params, "MaxRecords")
	resultItems, nextMarker, isTruncated := paginateItems(items, marker, maxRecords, func(item interface{}) string {
		return item.(*neptunestore.DBClusterParameterGroup).DBClusterParameterGroupName
	})

	result := map[string]interface{}{
		"DBClusterParameterGroups": protocol.XMLElements{ElementName: "DBClusterParameterGroup", Items: resultItems},
	}
	if isTruncated {
		result["Marker"] = nextMarker
	}
	return result, nil
}

// DescribeDBClusterParameters returns the parameters contained in the
// specified DB cluster parameter group.
func (s *NeptuneService) DescribeDBClusterParameters(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	name := request.GetStringParam(params, "DBClusterParameterGroupName")
	if name == "" {
		return nil, fmt.Errorf("neptune: DBClusterParameterGroupName is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	pg, err := store.GetClusterParameterGroup(name)
	if err != nil {
		return nil, translateStoreError(err)
	}

	defaultParams := []map[string]interface{}{
		{"ParameterName": "neptune_query_timeout", "ParameterValue": "120000", "Description": "Query execution timeout in milliseconds", "Source": "system", "ApplyType": "dynamic", "DataType": "integer", "IsModifiable": true},
		{"ParameterName": "neptune_enable_audit_log", "ParameterValue": "0", "Description": "Enable audit logging", "Source": "system", "ApplyType": "static", "DataType": "boolean", "IsModifiable": true},
	}
	userMods := make(map[string]neptunestore.Parameter, len(pg.Parameters))
	for _, p := range pg.Parameters {
		userMods[p.ParameterName] = p
	}

	items := make([]interface{}, 0, len(defaultParams))
	for _, p := range defaultParams {
		name := p["ParameterName"].(string)
		if mod, ok := userMods[name]; ok {
			items = append(items, map[string]interface{}{
				"ParameterName": mod.ParameterName, "ParameterValue": mod.ParameterValue,
				"Description": mod.Description, "Source": mod.Source, "ApplyType": mod.ApplyType,
				"DataType": mod.DataType, "IsModifiable": mod.IsModifiable,
			})
			delete(userMods, name)
		} else {
			items = append(items, p)
		}
	}
	for _, p := range userMods {
		items = append(items, map[string]interface{}{
			"ParameterName": p.ParameterName, "ParameterValue": p.ParameterValue,
			"Description": p.Description, "Source": p.Source, "ApplyType": p.ApplyType,
			"DataType": p.DataType, "IsModifiable": p.IsModifiable,
		})
	}

	marker := request.GetStringParam(params, "Marker")
	maxRecords := request.GetIntParam(params, "MaxRecords")
	resultItems, nextMarker, isTruncated := paginateItems(items, marker, maxRecords, func(item interface{}) string {
		return item.(map[string]interface{})["ParameterName"].(string)
	})

	result := map[string]interface{}{
		"Parameters": protocol.XMLElements{ElementName: "Parameter", Items: resultItems},
	}
	if isTruncated {
		result["Marker"] = nextMarker
	}
	return result, nil
}

// ModifyDBClusterParameterGroup modifies the parameters of the specified DB
// cluster parameter group.
func (s *NeptuneService) ModifyDBClusterParameterGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	name := request.GetStringParam(params, "DBClusterParameterGroupName")
	if name == "" {
		return nil, fmt.Errorf("neptune: DBClusterParameterGroupName is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	pg, err := store.GetClusterParameterGroup(name)
	if err != nil {
		return nil, translateStoreError(err)
	}

	if modParams := getNeptuneParameterList(params); len(modParams) > 0 {
		existing := make(map[string]neptunestore.Parameter, len(pg.Parameters))
		for _, p := range pg.Parameters {
			existing[p.ParameterName] = p
		}
		for _, mp := range modParams {
			existing[mp.ParameterName] = mp
		}
		pg.Parameters = make([]neptunestore.Parameter, 0, len(existing))
		for _, p := range existing {
			pg.Parameters = append(pg.Parameters, p)
		}
		if err := store.UpdateClusterParameterGroup(pg); err != nil {
			return nil, translateStoreError(err)
		}
	}

	return map[string]interface{}{
		"DBClusterParameterGroupName": name,
	}, nil
}

// ResetDBClusterParameterGroup resets the parameters of the specified DB
// cluster parameter group to their default values.
func (s *NeptuneService) ResetDBClusterParameterGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	name := request.GetStringParam(params, "DBClusterParameterGroupName")
	if name == "" {
		return nil, fmt.Errorf("neptune: DBClusterParameterGroupName is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	pg, err := store.GetClusterParameterGroup(name)
	if err != nil {
		return nil, translateStoreError(err)
	}

	pg.Parameters = nil
	if err := store.UpdateClusterParameterGroup(pg); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"DBClusterParameterGroupName": name,
	}, nil
}

// CopyDBClusterParameterGroup creates a copy of the specified DB cluster
// parameter group.
func (s *NeptuneService) CopyDBClusterParameterGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	sourceName := request.GetStringParam(params, "SourceDBClusterParameterGroupIdentifier")
	if sourceName == "" {
		return nil, fmt.Errorf("neptune: SourceDBClusterParameterGroupIdentifier is required")
	}
	targetName := request.GetStringParam(params, "TargetDBClusterParameterGroupIdentifier")
	if targetName == "" {
		return nil, fmt.Errorf("neptune: TargetDBClusterParameterGroupIdentifier is required")
	}
	desc := request.GetStringParam(params, "TargetDBParameterGroupDescription")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	source, err := store.GetClusterParameterGroup(sourceName)
	if err != nil {
		return nil, translateStoreError(err)
	}

	if desc == "" {
		desc = source.Description
	}

	pg := &neptunestore.DBClusterParameterGroup{
		DBClusterParameterGroupName: targetName,
		DBParameterGroupFamily:      source.DBParameterGroupFamily,
		Description:                 desc,
		ARN:                         neptunestore.ClusterParameterGroupARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), targetName),
		Parameters:                  append([]neptunestore.Parameter(nil), source.Parameters...),
	}

	if err := store.CreateClusterParameterGroup(pg); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"DBClusterParameterGroup": pg,
	}, nil
}

// CreateDBParameterGroup creates a new DB parameter group.
func (s *NeptuneService) CreateDBParameterGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	name := request.GetStringParam(params, "DBParameterGroupName")
	if name == "" {
		return nil, fmt.Errorf("neptune: DBParameterGroupName is required")
	}
	family := request.GetStringParam(params, "DBParameterGroupFamily")
	if family == "" {
		family = "neptune1"
	}
	desc := request.GetStringParam(params, "Description")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	pg := &neptunestore.DBParameterGroup{
		DBParameterGroupName:   name,
		DBParameterGroupFamily: family,
		Description:            desc,
		ARN:                    neptunestore.ParameterGroupARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), name),
	}

	if err := store.CreateParameterGroup(pg); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"DBParameterGroup": pg,
	}, nil
}

// DeleteDBParameterGroup deletes the specified DB parameter group.
func (s *NeptuneService) DeleteDBParameterGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	name := request.GetStringParam(params, "DBParameterGroupName")
	if name == "" {
		return nil, fmt.Errorf("neptune: DBParameterGroupName is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteParameterGroup(name); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{}, nil
}

// DescribeDBParameterGroups returns information about the specified DB
// parameter group or lists all groups when no name is provided.
func (s *NeptuneService) DescribeDBParameterGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	name := request.GetStringParam(params, "DBParameterGroupName")
	if name != "" {
		pg, err := store.GetParameterGroup(name)
		if err != nil {
			return nil, translateStoreError(err)
		}
		return map[string]interface{}{
			"DBParameterGroups": protocol.XMLElements{ElementName: "DBParameterGroup", Items: []interface{}{pg}},
		}, nil
	}

	groups, err := store.ListParameterGroups()
	if err != nil {
		return nil, translateStoreError(err)
	}

	items := make([]interface{}, 0, len(groups))
	for _, g := range groups {
		items = append(items, g)
	}

	marker := request.GetStringParam(params, "Marker")
	maxRecords := request.GetIntParam(params, "MaxRecords")
	resultItems, nextMarker, isTruncated := paginateItems(items, marker, maxRecords, func(item interface{}) string {
		return item.(*neptunestore.DBParameterGroup).DBParameterGroupName
	})

	result := map[string]interface{}{
		"DBParameterGroups": protocol.XMLElements{ElementName: "DBParameterGroup", Items: resultItems},
	}
	if isTruncated {
		result["Marker"] = nextMarker
	}
	return result, nil
}

// DescribeDBParameters returns the parameters contained in the specified DB
// parameter group.
func (s *NeptuneService) DescribeDBParameters(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	name := request.GetStringParam(params, "DBParameterGroupName")
	if name == "" {
		return nil, fmt.Errorf("neptune: DBParameterGroupName is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	pg, err := store.GetParameterGroup(name)
	if err != nil {
		return nil, translateStoreError(err)
	}

	defaultParams := []map[string]interface{}{
		{"ParameterName": "neptune_query_timeout", "ParameterValue": "120000", "Description": "Query execution timeout", "Source": "system", "ApplyType": "dynamic", "DataType": "integer", "IsModifiable": true},
	}
	userMods := make(map[string]neptunestore.Parameter, len(pg.Parameters))
	for _, p := range pg.Parameters {
		userMods[p.ParameterName] = p
	}

	items := make([]interface{}, 0, len(defaultParams))
	for _, p := range defaultParams {
		name := p["ParameterName"].(string)
		if mod, ok := userMods[name]; ok {
			items = append(items, map[string]interface{}{
				"ParameterName": mod.ParameterName, "ParameterValue": mod.ParameterValue,
				"Description": mod.Description, "Source": mod.Source, "ApplyType": mod.ApplyType,
				"DataType": mod.DataType, "IsModifiable": mod.IsModifiable,
			})
			delete(userMods, name)
		} else {
			items = append(items, p)
		}
	}
	for _, p := range userMods {
		items = append(items, map[string]interface{}{
			"ParameterName": p.ParameterName, "ParameterValue": p.ParameterValue,
			"Description": p.Description, "Source": p.Source, "ApplyType": p.ApplyType,
			"DataType": p.DataType, "IsModifiable": p.IsModifiable,
		})
	}

	marker := request.GetStringParam(params, "Marker")
	maxRecords := request.GetIntParam(params, "MaxRecords")
	resultItems, nextMarker, isTruncated := paginateItems(items, marker, maxRecords, func(item interface{}) string {
		return item.(map[string]interface{})["ParameterName"].(string)
	})

	result := map[string]interface{}{
		"Parameters": protocol.XMLElements{ElementName: "Parameter", Items: resultItems},
	}
	if isTruncated {
		result["Marker"] = nextMarker
	}
	return result, nil
}

// ModifyDBParameterGroup modifies the parameters of the specified DB parameter
// group.
func (s *NeptuneService) ModifyDBParameterGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	name := request.GetStringParam(params, "DBParameterGroupName")
	if name == "" {
		return nil, fmt.Errorf("neptune: DBParameterGroupName is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	pg, err := store.GetParameterGroup(name)
	if err != nil {
		return nil, translateStoreError(err)
	}

	if modParams := getNeptuneParameterList(params); len(modParams) > 0 {
		existing := make(map[string]neptunestore.Parameter, len(pg.Parameters))
		for _, p := range pg.Parameters {
			existing[p.ParameterName] = p
		}
		for _, mp := range modParams {
			existing[mp.ParameterName] = mp
		}
		pg.Parameters = make([]neptunestore.Parameter, 0, len(existing))
		for _, p := range existing {
			pg.Parameters = append(pg.Parameters, p)
		}
		if err := store.UpdateParameterGroup(pg); err != nil {
			return nil, translateStoreError(err)
		}
	}

	return map[string]interface{}{
		"DBParameterGroupName": name,
	}, nil
}

// ResetDBParameterGroup resets the parameters of the specified DB parameter
// group to their default values.
func (s *NeptuneService) ResetDBParameterGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	name := request.GetStringParam(params, "DBParameterGroupName")
	if name == "" {
		return nil, fmt.Errorf("neptune: DBParameterGroupName is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	pg, err := store.GetParameterGroup(name)
	if err != nil {
		return nil, translateStoreError(err)
	}

	pg.Parameters = nil
	if err := store.UpdateParameterGroup(pg); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"DBParameterGroupName": name,
	}, nil
}

// CopyDBParameterGroup creates a copy of the specified DB parameter group.
func (s *NeptuneService) CopyDBParameterGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	sourceName := request.GetStringParam(params, "SourceDBParameterGroupIdentifier")
	if sourceName == "" {
		return nil, fmt.Errorf("neptune: SourceDBParameterGroupIdentifier is required")
	}
	targetName := request.GetStringParam(params, "TargetDBParameterGroupIdentifier")
	if targetName == "" {
		return nil, fmt.Errorf("neptune: TargetDBParameterGroupIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	source, err := store.GetParameterGroup(sourceName)
	if err != nil {
		return nil, translateStoreError(err)
	}

	desc := request.GetStringParam(params, "TargetDBParameterGroupDescription")
	if desc == "" {
		desc = source.Description
	}

	pg := &neptunestore.DBParameterGroup{
		DBParameterGroupName:   targetName,
		DBParameterGroupFamily: source.DBParameterGroupFamily,
		Description:            desc,
		ARN:                    neptunestore.ParameterGroupARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), targetName),
		Parameters:             append([]neptunestore.Parameter(nil), source.Parameters...),
	}

	if err := store.CreateParameterGroup(pg); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"DBParameterGroup": pg,
	}, nil
}

// DescribeEngineDefaultClusterParameters returns the default engine parameters
// for the specified cluster parameter group family.
func (s *NeptuneService) DescribeEngineDefaultClusterParameters(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	family := request.GetStringParam(params, "DBParameterGroupFamily")
	if family == "" {
		family = "neptune1"
	}

	return map[string]interface{}{
		"EngineDefaults": map[string]interface{}{
			"DBParameterGroupFamily": family,
			"Parameters": protocol.XMLElements{ElementName: "Parameter", Items: []interface{}{
				map[string]interface{}{"ParameterName": "neptune_query_timeout", "ParameterValue": "120000", "Source": "engine-default", "ApplyType": "dynamic", "DataType": "integer", "IsModifiable": true},
				map[string]interface{}{"ParameterName": "neptune_enable_audit_log", "ParameterValue": "0", "Source": "engine-default", "ApplyType": "static", "DataType": "boolean", "IsModifiable": true},
			}},
		},
	}, nil
}

// DescribeEngineDefaultParameters returns the default engine parameters for the
// specified DB parameter group family.
func (s *NeptuneService) DescribeEngineDefaultParameters(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	family := request.GetStringParam(params, "DBParameterGroupFamily")
	if family == "" {
		family = "neptune1"
	}

	return map[string]interface{}{
		"EngineDefaults": map[string]interface{}{
			"DBParameterGroupFamily": family,
			"Parameters": protocol.XMLElements{ElementName: "Parameter", Items: []interface{}{
				map[string]interface{}{"ParameterName": "neptune_query_timeout", "ParameterValue": "120000", "Source": "engine-default", "ApplyType": "dynamic", "DataType": "integer", "IsModifiable": true},
			}},
		},
	}, nil
}
