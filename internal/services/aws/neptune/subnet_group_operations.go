package neptune

import (
	"context"
	"fmt"

	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	neptunestore "vorpalstacks/internal/store/aws/neptune"
)

func getNeptuneStringList(params map[string]interface{}, key string, memberNames ...string) []string {
	var result []string
	if len(memberNames) == 0 {
		memberNames = []string{"SubnetIdentifier", "member"}
	}
	for _, memberName := range memberNames {
		for i := 1; ; i++ {
			memberKey := fmt.Sprintf("%s.%s.%d", key, memberName, i)
			val := request.GetStringParam(params, memberKey)
			if val == "" {
				break
			}
			result = append(result, val)
		}
		if len(result) > 0 {
			return result
		}
	}
	return request.GetStringList(params, key)
}

func getNeptuneTagList(params map[string]interface{}) []map[string]interface{} {
	var tags []map[string]interface{}
	for i := 1; ; i++ {
		prefix := fmt.Sprintf("Tags.Tag.%d", i)
		key := request.GetStringParam(params, prefix+".Key")
		if key == "" {
			break
		}
		value := request.GetStringParam(params, prefix+".Value")
		tags = append(tags, map[string]interface{}{"Key": key, "Value": value})
	}
	return tags
}

func getNeptuneParameterList(params map[string]interface{}) []neptunestore.Parameter {
	var parameters []neptunestore.Parameter
	for i := 1; ; i++ {
		prefix := fmt.Sprintf("Parameters.Parameter.%d", i)
		paramName := request.GetStringParam(params, prefix+".ParameterName")
		if paramName == "" {
			break
		}
		paramValue := request.GetStringParam(params, prefix+".ParameterValue")
		description := request.GetStringParam(params, prefix+".Description")
		source := request.GetStringParam(params, prefix+".Source")
		applyType := request.GetStringParam(params, prefix+".ApplyType")
		dataType := request.GetStringParam(params, prefix+".DataType")
		parameters = append(parameters, neptunestore.Parameter{
			ParameterName:  paramName,
			ParameterValue: paramValue,
			Description:    description,
			Source:         source,
			ApplyType:      applyType,
			DataType:       dataType,
			IsModifiable:   true,
		})
	}
	return parameters
}

// CreateDBSubnetGroup creates a new DB subnet group with the specified subnets.
func (s *NeptuneService) CreateDBSubnetGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	name := request.GetStringParam(params, "DBSubnetGroupName")
	if name == "" {
		return nil, fmt.Errorf("neptune: DBSubnetGroupName is required")
	}
	desc := request.GetStringParam(params, "DBSubnetGroupDescription")
	subnetIds := getNeptuneStringList(params, "SubnetIds")
	if len(subnetIds) == 0 {
		return nil, fmt.Errorf("neptune: SubnetIds is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if len(subnetIds) > 26 {
		return nil, fmt.Errorf("neptune: cannot assign more than 26 subnets to a DB subnet group")
	}
	subnets := make([]neptunestore.Subnet, 0, len(subnetIds))
	region := reqCtx.GetRegion()
	for i, id := range subnetIds {
		az := fmt.Sprintf("%s%c", region, 'a'+byte(i))
		subnets = append(subnets, neptunestore.Subnet{
			SubnetIdentifier:       id,
			SubnetAvailabilityZone: az,
			SubnetStatus:           "Active",
		})
	}

	sg := &neptunestore.DBSubnetGroup{
		DBSubnetGroupName:        name,
		DBSubnetGroupDescription: desc,
		VpcId:                    fmt.Sprintf("vpc-%s", region),
		SubnetGroupStatus:        "Complete",
		Subnets:                  subnets,
		ARN:                      neptunestore.SubnetGroupARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), name),
	}

	if err := store.CreateSubnetGroup(sg); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"DBSubnetGroup": sg,
	}, nil
}

// DeleteDBSubnetGroup deletes the specified DB subnet group.
func (s *NeptuneService) DeleteDBSubnetGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	name := request.GetStringParam(params, "DBSubnetGroupName")
	if name == "" {
		return nil, fmt.Errorf("neptune: DBSubnetGroupName is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteSubnetGroup(name); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{}, nil
}

// DescribeDBSubnetGroups describes one or all DB subnet groups.
func (s *NeptuneService) DescribeDBSubnetGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	name := request.GetStringParam(params, "DBSubnetGroupName")
	if name != "" {
		sg, err := store.GetSubnetGroup(name)
		if err != nil {
			return nil, translateStoreError(err)
		}
		return map[string]interface{}{
			"DBSubnetGroups": protocol.XMLElements{ElementName: "DBSubnetGroup", Items: []interface{}{sg}},
		}, nil
	}

	groups, err := store.ListSubnetGroups()
	if err != nil {
		return nil, translateStoreError(err)
	}

	result := make([]interface{}, 0, len(groups))
	for _, g := range groups {
		result = append(result, g)
	}

	return map[string]interface{}{
		"DBSubnetGroups": protocol.XMLElements{ElementName: "DBSubnetGroup", Items: result},
	}, nil
}

// ModifyDBSubnetGroup modifies the description or subnet list of an existing DB subnet group.
func (s *NeptuneService) ModifyDBSubnetGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	name := request.GetStringParam(params, "DBSubnetGroupName")
	if name == "" {
		return nil, fmt.Errorf("neptune: DBSubnetGroupName is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	sg, err := store.GetSubnetGroup(name)
	if err != nil {
		return nil, translateStoreError(err)
	}

	if desc := request.GetStringParam(params, "DBSubnetGroupDescription"); desc != "" {
		sg.DBSubnetGroupDescription = desc
	}
	if subnetIds := getNeptuneStringList(params, "SubnetIds"); len(subnetIds) > 0 {
		if len(subnetIds) > 26 {
			return nil, fmt.Errorf("neptune: cannot assign more than 26 subnets to a DB subnet group")
		}
		subnets := make([]neptunestore.Subnet, 0, len(subnetIds))
		region := reqCtx.GetRegion()
		for i, id := range subnetIds {
			az := fmt.Sprintf("%s%c", region, 'a'+byte(i))
			subnets = append(subnets, neptunestore.Subnet{
				SubnetIdentifier:       id,
				SubnetAvailabilityZone: az,
				SubnetStatus:           "Active",
			})
		}
		sg.Subnets = subnets
	}

	if err := store.UpdateSubnetGroup(sg); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"DBSubnetGroup": sg,
	}, nil
}
