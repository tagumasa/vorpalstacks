package neptune

import (
	"context"
	"fmt"

	"vorpalstacks/internal/services/aws/common/protocol"
	"vorpalstacks/internal/services/aws/common/request"
	neptunestore "vorpalstacks/internal/store/aws/neptune"
)

// getNeptuneStringList extracts a string list from Neptune/RDS Query protocol parameters.
// Neptune uses {ListKey}.{MemberName}.{N} (e.g. SubnetIds.SubnetIdentifier.1)
// rather than the EC2-style {ListKey}.member.{N} used by some other services.
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

// getNeptuneTagList extracts a list of tags from Neptune/RDS Query protocol parameters.
// Wire format: Tags.Tag.1.Key=k1, Tags.Tag.1.Value=v1, Tags.Tag.2.Key=k2, ...
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

	subnets := make([]neptunestore.Subnet, 0, len(subnetIds))
	for i, id := range subnetIds {
		az := fmt.Sprintf("us-east-1%c", 'a'+i)
		subnets = append(subnets, neptunestore.Subnet{
			SubnetIdentifier:       id,
			SubnetAvailabilityZone: az,
			SubnetStatus:           "Active",
		})
	}

	sg := &neptunestore.DBSubnetGroup{
		DBSubnetGroupName:        name,
		DBSubnetGroupDescription: desc,
		VpcId:                    "vpc-default",
		SubnetGroupStatus:        "Complete",
		Subnets:                  subnets,
		ARN:                      neptunestore.SubnetGroupARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), name),
	}

	if err := store.CreateSubnetGroup(sg); err != nil {
		return nil, err
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
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// DescribeDBSubnetGroups returns information about the specified DB subnet
// group or lists all groups when no name is provided.
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
			return nil, err
		}
		return map[string]interface{}{
			"DBSubnetGroups": protocol.XMLElements{ElementName: "DBSubnetGroup", Items: []interface{}{sg}},
		}, nil
	}

	groups, err := store.ListSubnetGroups()
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, 0, len(groups))
	for _, g := range groups {
		result = append(result, g)
	}

	return map[string]interface{}{
		"DBSubnetGroups": protocol.XMLElements{ElementName: "DBSubnetGroup", Items: result},
	}, nil
}

// ModifyDBSubnetGroup updates the configuration of the specified DB subnet group.
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
		return nil, err
	}

	if desc := request.GetStringParam(params, "DBSubnetGroupDescription"); desc != "" {
		sg.DBSubnetGroupDescription = desc
	}
	if subnetIds := getNeptuneStringList(params, "SubnetIds"); len(subnetIds) > 0 {
		subnets := make([]neptunestore.Subnet, 0, len(subnetIds))
		for i, id := range subnetIds {
			az := fmt.Sprintf("us-east-1%c", 'a'+i)
			subnets = append(subnets, neptunestore.Subnet{
				SubnetIdentifier:       id,
				SubnetAvailabilityZone: az,
				SubnetStatus:           "Active",
			})
		}
		sg.Subnets = subnets
	}

	_ = store.UpdateSubnetGroup(sg)

	return map[string]interface{}{
		"DBSubnetGroup": sg,
	}, nil
}
