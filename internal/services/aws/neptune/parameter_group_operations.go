package neptune

import (
	"context"
	"fmt"

	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	neptunestore "vorpalstacks/internal/store/aws/neptune"
)

var clusterDefaultParams = []map[string]interface{}{
	{"ParameterName": "neptune_query_timeout", "ParameterValue": "120000", "Description": "Query execution timeout in milliseconds", "Source": "system", "ApplyType": "dynamic", "DataType": "integer", "IsModifiable": true},
	{"ParameterName": "neptune_enable_audit_log", "ParameterValue": "0", "Description": "Enable audit logging", "Source": "system", "ApplyType": "static", "DataType": "boolean", "IsModifiable": true},
}

var instanceDefaultParams = []map[string]interface{}{
	{"ParameterName": "neptune_query_timeout", "ParameterValue": "120000", "Description": "Query execution timeout", "Source": "system", "ApplyType": "dynamic", "DataType": "integer", "IsModifiable": true},
}

// buildParameterItems merges user-modified parameters with the default
// parameter templates, overriding defaults where user values exist.
func buildParameterItems(defaults []map[string]interface{}, userParams []neptunestore.Parameter) []interface{} {
	userMods := make(map[string]neptunestore.Parameter, len(userParams))
	for _, p := range userParams {
		userMods[p.ParameterName] = p
	}

	items := make([]interface{}, 0, len(defaults)+len(userMods))
	for _, p := range defaults {
		pName := p["ParameterName"].(string)
		if mod, ok := userMods[pName]; ok {
			items = append(items, mergeParameter(p, mod))
			delete(userMods, pName)
		} else {
			items = append(items, p)
		}
	}
	for _, p := range userMods {
		if p.Source == "" {
			p.Source = "user"
		}
		items = append(items, map[string]interface{}{
			"ParameterName": p.ParameterName, "ParameterValue": p.ParameterValue,
			"Description": p.Description, "Source": p.Source, "ApplyType": p.ApplyType,
			"DataType": p.DataType, "IsModifiable": p.IsModifiable,
		})
	}
	return items
}

// mergeParameter combines a default parameter template with user modifications,
// falling back to the template value when the user modification is empty.
func mergeParameter(def map[string]interface{}, mod neptunestore.Parameter) map[string]interface{} {
	desc := mod.Description
	if desc == "" {
		desc = def["Description"].(string)
	}
	source := mod.Source
	if source == "" {
		source = "user"
	}
	applyType := mod.ApplyType
	if applyType == "" {
		applyType = def["ApplyType"].(string)
	}
	dataType := mod.DataType
	if dataType == "" {
		dataType = def["DataType"].(string)
	}
	return map[string]interface{}{
		"ParameterName": mod.ParameterName, "ParameterValue": mod.ParameterValue,
		"Description": desc, "Source": source, "ApplyType": applyType,
		"DataType": dataType, "IsModifiable": mod.IsModifiable,
	}
}

// applyParameterModifications merges request parameters into an existing
// parameter group, returning the updated parameter slice.
func applyParameterModifications(existing []neptunestore.Parameter, mods []neptunestore.Parameter) []neptunestore.Parameter {
	m := make(map[string]neptunestore.Parameter, len(existing))
	for _, p := range existing {
		m[p.ParameterName] = p
	}
	for _, mp := range mods {
		if prev, ok := m[mp.ParameterName]; ok {
			prev.ParameterValue = mp.ParameterValue
			if mp.ApplyMethod != "" {
				prev.ApplyMethod = mp.ApplyMethod
			}
			m[mp.ParameterName] = prev
		} else {
			m[mp.ParameterName] = mp
		}
	}
	result := make([]neptunestore.Parameter, 0, len(m))
	for _, p := range m {
		result = append(result, p)
	}
	return result
}

// paginateParameters applies client-side pagination to a parameter item list.
func paginateParameters(items []interface{}, params map[string]interface{}) (map[string]interface{}, bool) {
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
	return result, isTruncated
}

// --- Cluster Parameter Group handlers ---

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

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	pg := &neptunestore.DBClusterParameterGroup{
		DBClusterParameterGroupName: name,
		DBParameterGroupFamily:      family,
		Description:                 request.GetStringParam(params, "Description"),
		ARN:                         neptunestore.ClusterParameterGroupARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), name),
	}
	if err := store.CreateClusterParameterGroup(pg); err != nil {
		return nil, translateStoreError(err)
	}
	return map[string]interface{}{"DBClusterParameterGroup": pg}, nil
}

func (s *NeptuneService) DeleteDBClusterParameterGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "DBClusterParameterGroupName")
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

func (s *NeptuneService) DescribeDBClusterParameterGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	name := request.GetStringParam(req.Parameters, "DBClusterParameterGroupName")
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
	items := make([]interface{}, len(groups))
	for i, g := range groups {
		items[i] = g
	}
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxRecords := request.GetIntParam(req.Parameters, "MaxRecords")
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

func (s *NeptuneService) DescribeDBClusterParameters(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "DBClusterParameterGroupName")
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
	items := buildParameterItems(clusterDefaultParams, pg.Parameters)
	result, _ := paginateParameters(items, req.Parameters)
	return result, nil
}

func (s *NeptuneService) ModifyDBClusterParameterGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "DBClusterParameterGroupName")
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
	if modParams := getNeptuneParameterList(req.Parameters); len(modParams) > 0 {
		pg.Parameters = applyParameterModifications(pg.Parameters, modParams)
		if err := store.UpdateClusterParameterGroup(pg); err != nil {
			return nil, translateStoreError(err)
		}
	}
	return map[string]interface{}{"DBClusterParameterGroupName": name}, nil
}

func (s *NeptuneService) ResetDBClusterParameterGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "DBClusterParameterGroupName")
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
	return map[string]interface{}{"DBClusterParameterGroupName": name}, nil
}

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
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	source, err := store.GetClusterParameterGroup(sourceName)
	if err != nil {
		return nil, translateStoreError(err)
	}
	desc := request.GetStringParam(params, "TargetDBParameterGroupDescription")
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
	return map[string]interface{}{"DBClusterParameterGroup": pg}, nil
}

// --- Instance Parameter Group handlers ---

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
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	pg := &neptunestore.DBParameterGroup{
		DBParameterGroupName:   name,
		DBParameterGroupFamily: family,
		Description:            request.GetStringParam(params, "Description"),
		ARN:                    neptunestore.ParameterGroupARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), name),
	}
	if err := store.CreateParameterGroup(pg); err != nil {
		return nil, translateStoreError(err)
	}
	return map[string]interface{}{"DBParameterGroup": pg}, nil
}

func (s *NeptuneService) DeleteDBParameterGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "DBParameterGroupName")
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

func (s *NeptuneService) DescribeDBParameterGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	name := request.GetStringParam(req.Parameters, "DBParameterGroupName")
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
	items := make([]interface{}, len(groups))
	for i, g := range groups {
		items[i] = g
	}
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxRecords := request.GetIntParam(req.Parameters, "MaxRecords")
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

func (s *NeptuneService) DescribeDBParameters(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "DBParameterGroupName")
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
	items := buildParameterItems(instanceDefaultParams, pg.Parameters)
	result, _ := paginateParameters(items, req.Parameters)
	return result, nil
}

func (s *NeptuneService) ModifyDBParameterGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "DBParameterGroupName")
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
	if modParams := getNeptuneParameterList(req.Parameters); len(modParams) > 0 {
		pg.Parameters = applyParameterModifications(pg.Parameters, modParams)
		if err := store.UpdateParameterGroup(pg); err != nil {
			return nil, translateStoreError(err)
		}
	}
	return map[string]interface{}{"DBParameterGroupName": name}, nil
}

func (s *NeptuneService) ResetDBParameterGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "DBParameterGroupName")
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
	return map[string]interface{}{"DBParameterGroupName": name}, nil
}

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
	return map[string]interface{}{"DBParameterGroup": pg}, nil
}

func (s *NeptuneService) DescribeEngineDefaultClusterParameters(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	family := request.GetStringParam(req.Parameters, "DBParameterGroupFamily")
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

func (s *NeptuneService) DescribeEngineDefaultParameters(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	family := request.GetStringParam(req.Parameters, "DBParameterGroupFamily")
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
