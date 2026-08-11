package neptune

import (
	"context"
	"fmt"
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	rdssvc "vorpalstacks/internal/services/aws/rds"
	neptunestore "vorpalstacks/internal/store/aws/rds/neptune"
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

// resolveSubnets queries the EC2 service via the event bus to resolve each
// subnet ID to its VPC ID and availability zone. Returns the resolved subnet
// slice and the common VPC ID (all subnets must belong to the same VPC).
func (s *NeptuneService) resolveSubnets(ctx context.Context, region string, subnetIds []string) ([]neptunestore.Subnet, string, error) {
	if s.eventBus == nil {
		return nil, "", awserrors.NewAWSError("InvalidSubnet", "EC2 event bus not available for subnet lookup", http.StatusBadRequest)
	}
	ec2 := s.eventBus.EC2Invoker()
	if ec2 == nil {
		return nil, "", awserrors.NewAWSError("InvalidSubnet", "EC2 service not available for subnet lookup", http.StatusBadRequest)
	}

	subnets := make([]neptunestore.Subnet, 0, len(subnetIds))
	var vpcId string

	for _, subnetId := range subnetIds {
		subnetVpcId, az, err := ec2.LookupSubnet(ctx, region, subnetId)
		if err != nil {
			return nil, "", awserrors.NewAWSError("InvalidSubnet", fmt.Sprintf("Failed to resolve subnet %s: %v", subnetId, err), http.StatusBadRequest)
		}
		if vpcId == "" {
			vpcId = subnetVpcId
		} else if subnetVpcId != vpcId {
			return nil, "", awserrors.NewValidationException(fmt.Sprintf("All subnets must belong to the same VPC (expected %s, got %s)", vpcId, subnetVpcId))
		}
		subnets = append(subnets, neptunestore.Subnet{
			SubnetIdentifier:       subnetId,
			SubnetAvailabilityZone: az,
			SubnetStatus:           "Active",
		})
	}
	return subnets, vpcId, nil
}

// resolveSecurityGroups validates that each security group ID exists via the
// EC2 invoker and returns the common VPC ID. All security groups must belong
// to the same VPC.
func (s *NeptuneService) resolveSecurityGroups(ctx context.Context, region string, sgIds []string) (string, error) {
	if s.eventBus == nil {
		return "", awserrors.NewAWSError("InvalidParameterValue", "EC2 event bus not available for security group lookup", http.StatusBadRequest)
	}
	ec2 := s.eventBus.EC2Invoker()
	if ec2 == nil {
		return "", awserrors.NewAWSError("InvalidParameterValue", "EC2 service not available for security group lookup", http.StatusBadRequest)
	}

	var vpcId string
	for _, sgId := range sgIds {
		sgVpcId, err := ec2.LookupSecurityGroup(ctx, region, sgId)
		if err != nil {
			return "", awserrors.NewAWSError("InvalidParameterValue", fmt.Sprintf("Failed to resolve security group %s: %v", sgId, err), http.StatusBadRequest)
		}
		if vpcId == "" {
			vpcId = sgVpcId
		} else if sgVpcId != vpcId {
			return "", awserrors.NewValidationException(fmt.Sprintf("All security groups must belong to the same VPC (expected %s, got %s)", vpcId, sgVpcId))
		}
	}
	return vpcId, nil
}

// CreateDBSubnetGroup creates a new DB subnet group with the specified subnets.
func (s *NeptuneService) CreateDBSubnetGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	name := request.GetStringParam(params, "DBSubnetGroupName")
	if name == "" {
		return nil, awserrors.NewMissingParameter("DBSubnetGroupName is required")
	}
	desc := request.GetStringParam(params, "DBSubnetGroupDescription")
	subnetIds := getNeptuneStringList(params, "SubnetIds")
	if len(subnetIds) == 0 {
		return nil, awserrors.NewMissingParameter("SubnetIds is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := validateSubnetCount(len(subnetIds)); err != nil {
		return nil, awserrors.NewAWSError("DBSubnetGroupQuotaExceededFault", err.Error(), http.StatusBadRequest)
	}

	region := reqCtx.GetRegion()
	subnets, vpcId, err := s.resolveSubnets(ctx, region, subnetIds)
	if err != nil {
		return nil, err
	}

	sg := &neptunestore.DBSubnetGroup{
		DBSubnetGroupName:        name,
		DBSubnetGroupDescription: desc,
		VpcId:                    vpcId,
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
		return nil, awserrors.NewMissingParameter("DBSubnetGroupName is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	sg, err := store.GetSubnetGroup(name)
	if err != nil {
		return nil, translateStoreError(err)
	}
	if err := store.DeleteSubnetGroup(name); err != nil {
		return nil, translateStoreError(err)
	}
	removeTagsForResource(store, sg.ARN)
	return map[string]interface{}{}, nil
}

// DescribeDBSubnetGroups describes one or all DB subnet groups.
func (s *NeptuneService) DescribeDBSubnetGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	groups, nextMarker, err := rdssvc.QuerySubnetGroups(store, rdssvc.DescribeDBSubnetGroupsInput{
		DBSubnetGroupName: request.GetStringParam(params, "DBSubnetGroupName"),
		Filters:           nil,
		Marker:            request.GetStringParam(params, "Marker"),
		MaxRecords:        int32(request.GetIntParam(params, "MaxRecords")),
	})
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(groups))
	for _, g := range groups {
		items = append(items, g)
	}

	result := map[string]interface{}{
		"DBSubnetGroups": protocol.XMLElements{ElementName: "DBSubnetGroup", Items: items},
	}
	if nextMarker != "" {
		result["Marker"] = nextMarker
	}
	return result, nil
}

// ModifyDBSubnetGroup modifies the description or subnet list of an existing DB subnet group.
func (s *NeptuneService) ModifyDBSubnetGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	name := request.GetStringParam(params, "DBSubnetGroupName")
	if name == "" {
		return nil, awserrors.NewMissingParameter("DBSubnetGroupName is required")
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
			return nil, awserrors.NewAWSError("DBSubnetGroupQuotaExceededFault", "Cannot assign more than 26 subnets to a DB subnet group", http.StatusBadRequest)
		}
		region := reqCtx.GetRegion()
		subnets, vpcId, err := s.resolveSubnets(ctx, region, subnetIds)
		if err != nil {
			return nil, err
		}
		sg.Subnets = subnets
		sg.VpcId = vpcId
	}

	if err := store.UpdateSubnetGroup(sg); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"DBSubnetGroup": sg,
	}, nil
}
