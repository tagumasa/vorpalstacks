package neptune

import (
	"context"
	"sort"

	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	rdssvc "vorpalstacks/internal/services/aws/rds"
	neptunestore "vorpalstacks/internal/store/aws/rds/neptune"
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
	sort.Slice(result, func(i, j int) bool { return result[i].ParameterName < result[j].ParameterName })
	return result
}

// resetNamedParameters removes the user-modified entries named by the reset
// list so the engine defaults show through again; every other modification
// is preserved.
func resetNamedParameters(existing []neptunestore.Parameter, reset []neptunestore.Parameter) []neptunestore.Parameter {
	names := make(map[string]bool, len(reset))
	for _, p := range reset {
		names[p.ParameterName] = true
	}
	result := make([]neptunestore.Parameter, 0, len(existing))
	for _, p := range existing {
		if names[p.ParameterName] {
			continue
		}
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

// CreateDBClusterParameterGroup creates a new DB cluster parameter group.
func (s *NeptuneService) CreateDBClusterParameterGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &CreateDBClusterParameterGroupInput{
		DBClusterParameterGroupName: request.GetStringParam(params, "DBClusterParameterGroupName"),
		DBParameterGroupFamily:      request.GetStringParam(params, "DBParameterGroupFamily"),
		Description:                 request.GetStringParam(params, "Description"),
		AccountID:                   reqCtx.GetAccountID(),
		Region:                      reqCtx.GetRegion(),
	}
	return s.createDBClusterParameterGroupCore(ctx, store, in)
}

// DeleteDBClusterParameterGroup deletes the specified DB cluster parameter group.
func (s *NeptuneService) DeleteDBClusterParameterGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &DeleteDBClusterParameterGroupInput{
		DBClusterParameterGroupName: request.GetStringParam(req.Parameters, "DBClusterParameterGroupName"),
	}
	return s.deleteDBClusterParameterGroupCore(ctx, store, in)
}

// DescribeDBClusterParameterGroups returns information about the specified
// cluster parameter group, or lists all groups when no name is provided.
func (s *NeptuneService) DescribeDBClusterParameterGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	groups, nextMarker, err := rdssvc.QueryClusterParameterGroups(store, rdssvc.DescribeDBClusterParameterGroupsInput{
		DBClusterParameterGroupName: request.GetStringParam(req.Parameters, "DBClusterParameterGroupName"),
		Filters:                     nil,
		Marker:                      request.GetStringParam(req.Parameters, "Marker"),
		MaxRecords:                  int32(request.GetIntParam(req.Parameters, "MaxRecords")),
	})
	if err != nil {
		return nil, err
	}
	items := make([]interface{}, 0, len(groups))
	for _, g := range groups {
		items = append(items, g)
	}
	result := map[string]interface{}{
		"DBClusterParameterGroups": protocol.XMLElements{ElementName: "DBClusterParameterGroup", Items: items},
	}
	if nextMarker != "" {
		result["Marker"] = nextMarker
	}
	return result, nil
}

// DescribeDBClusterParameters returns the parameters contained in the
// specified DB cluster parameter group.
func (s *NeptuneService) DescribeDBClusterParameters(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &DescribeDBClusterParametersInput{
		DBClusterParameterGroupName: request.GetStringParam(req.Parameters, "DBClusterParameterGroupName"),
	}
	pg, err := s.getClusterParameterGroupCore(ctx, store, in)
	if err != nil {
		return nil, translateStoreError(err)
	}
	items := buildParameterItems(clusterDefaultParams, pg.Parameters)
	result, _ := paginateParameters(items, req.Parameters)
	return result, nil
}

// ModifyDBClusterParameterGroup modifies the parameters of the specified DB
// cluster parameter group.
func (s *NeptuneService) ModifyDBClusterParameterGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &ModifyDBClusterParameterGroupInput{
		DBClusterParameterGroupName: request.GetStringParam(req.Parameters, "DBClusterParameterGroupName"),
		Parameters:                  getNeptuneParameterList(req.Parameters),
	}
	return s.modifyDBClusterParameterGroupCore(ctx, store, in)
}

// ResetDBClusterParameterGroup resets the parameters of the specified DB
// cluster parameter group to their default values.
func (s *NeptuneService) ResetDBClusterParameterGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &ResetDBClusterParameterGroupInput{
		DBClusterParameterGroupName: request.GetStringParam(req.Parameters, "DBClusterParameterGroupName"),
		Parameters:                  getNeptuneParameterList(req.Parameters),
	}
	return s.resetDBClusterParameterGroupCore(ctx, store, in)
}

// CopyDBClusterParameterGroup creates a copy of the specified DB cluster
// parameter group.
func (s *NeptuneService) CopyDBClusterParameterGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &CopyDBClusterParameterGroupInput{
		SourceDBClusterParameterGroupIdentifier:  request.GetStringParam(params, "SourceDBClusterParameterGroupIdentifier"),
		TargetDBClusterParameterGroupIdentifier:  request.GetStringParam(params, "TargetDBClusterParameterGroupIdentifier"),
		TargetDBClusterParameterGroupDescription: request.GetStringParam(params, "TargetDBClusterParameterGroupDescription"),
		AccountID:                                reqCtx.GetAccountID(),
		Region:                                   reqCtx.GetRegion(),
	}
	return s.copyDBClusterParameterGroupCore(ctx, store, in)
}

// --- Instance Parameter Group handlers ---

// CreateDBParameterGroup creates a new DB parameter group.
func (s *NeptuneService) CreateDBParameterGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &CreateDBParameterGroupInput{
		DBParameterGroupName:   request.GetStringParam(params, "DBParameterGroupName"),
		DBParameterGroupFamily: request.GetStringParam(params, "DBParameterGroupFamily"),
		Description:            request.GetStringParam(params, "Description"),
		AccountID:              reqCtx.GetAccountID(),
		Region:                 reqCtx.GetRegion(),
	}
	return s.createDBParameterGroupCore(ctx, store, in)
}

// DeleteDBParameterGroup deletes the specified DB parameter group.
func (s *NeptuneService) DeleteDBParameterGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &DeleteDBParameterGroupInput{
		DBParameterGroupName: request.GetStringParam(req.Parameters, "DBParameterGroupName"),
	}
	return s.deleteDBParameterGroupCore(ctx, store, in)
}

// DescribeDBParameterGroups returns information about the specified DB
// parameter group, or lists all groups when no name is provided.
func (s *NeptuneService) DescribeDBParameterGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	groups, nextMarker, err := rdssvc.QueryParameterGroups(store, rdssvc.DescribeDBParameterGroupsInput{
		DBParameterGroupName: request.GetStringParam(req.Parameters, "DBParameterGroupName"),
		Filters:              nil,
		Marker:               request.GetStringParam(req.Parameters, "Marker"),
		MaxRecords:           int32(request.GetIntParam(req.Parameters, "MaxRecords")),
	})
	if err != nil {
		return nil, err
	}
	items := make([]interface{}, 0, len(groups))
	for _, g := range groups {
		items = append(items, g)
	}
	result := map[string]interface{}{
		"DBParameterGroups": protocol.XMLElements{ElementName: "DBParameterGroup", Items: items},
	}
	if nextMarker != "" {
		result["Marker"] = nextMarker
	}
	return result, nil
}

// DescribeDBParameters returns the parameters contained in the specified DB
// parameter group.
func (s *NeptuneService) DescribeDBParameters(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &DescribeDBParametersInput{
		DBParameterGroupName: request.GetStringParam(req.Parameters, "DBParameterGroupName"),
	}
	pg, err := s.getParameterGroupCore(ctx, store, in)
	if err != nil {
		return nil, translateStoreError(err)
	}
	items := buildParameterItems(instanceDefaultParams, pg.Parameters)
	result, _ := paginateParameters(items, req.Parameters)
	return result, nil
}

// ModifyDBParameterGroup modifies the parameters of the specified DB parameter
// group.
func (s *NeptuneService) ModifyDBParameterGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &ModifyDBParameterGroupInput{
		DBParameterGroupName: request.GetStringParam(req.Parameters, "DBParameterGroupName"),
		Parameters:           getNeptuneParameterList(req.Parameters),
	}
	return s.modifyDBParameterGroupCore(ctx, store, in)
}

// ResetDBParameterGroup resets the parameters of the specified DB parameter
// group to their default values.
func (s *NeptuneService) ResetDBParameterGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &ResetDBParameterGroupInput{
		DBParameterGroupName: request.GetStringParam(req.Parameters, "DBParameterGroupName"),
		Parameters:           getNeptuneParameterList(req.Parameters),
	}
	return s.resetDBParameterGroupCore(ctx, store, in)
}

// CopyDBParameterGroup creates a copy of the specified DB parameter group.
func (s *NeptuneService) CopyDBParameterGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &CopyDBParameterGroupInput{
		SourceDBParameterGroupIdentifier:  request.GetStringParam(params, "SourceDBParameterGroupIdentifier"),
		TargetDBParameterGroupIdentifier:  request.GetStringParam(params, "TargetDBParameterGroupIdentifier"),
		TargetDBParameterGroupDescription: request.GetStringParam(params, "TargetDBParameterGroupDescription"),
		AccountID:                         reqCtx.GetAccountID(),
		Region:                            reqCtx.GetRegion(),
	}
	return s.copyDBParameterGroupCore(ctx, store, in)
}

// DescribeEngineDefaultClusterParameters returns the default engine parameters
// for the specified cluster parameter group family.
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

// DescribeEngineDefaultParameters returns the default engine parameters for the
// specified DB parameter group family.
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
