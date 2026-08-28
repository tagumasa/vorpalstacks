package neptune

import (
	"context"
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
	neptunestore "vorpalstacks/internal/store/aws/rds/neptune"
)

// CreateDBSubnetGroupInput carries the wire-parsed CreateDBSubnetGroup
// request.
type CreateDBSubnetGroupInput struct {
	DBSubnetGroupName        string
	DBSubnetGroupDescription string
	SubnetIds                []string
	AccountID                string
	Region                   string
}

// DeleteDBSubnetGroupInput carries the wire-parsed DeleteDBSubnetGroup
// request.
type DeleteDBSubnetGroupInput struct {
	DBSubnetGroupName string
}

// ModifyDBSubnetGroupInput carries the wire-parsed ModifyDBSubnetGroup
// request. An empty SubnetIds slice leaves the stored subnet list unchanged.
type ModifyDBSubnetGroupInput struct {
	DBSubnetGroupName        string
	DBSubnetGroupDescription string
	SubnetIds                []string
	Region                   string
}

// createDBSubnetGroupCore resolves the subnets through EC2, enforces the
// subnet-count quota and persists the new subnet group.
func (s *NeptuneService) createDBSubnetGroupCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *CreateDBSubnetGroupInput) (interface{}, error) {
	name := in.DBSubnetGroupName
	if name == "" {
		return nil, awserrors.NewMissingParameter("DBSubnetGroupName is required")
	}
	if len(in.SubnetIds) == 0 {
		return nil, awserrors.NewMissingParameter("SubnetIds is required")
	}

	if err := validateSubnetCount(len(in.SubnetIds)); err != nil {
		return nil, awserrors.NewAWSError("DBSubnetGroupQuotaExceededFault", err.Error(), http.StatusBadRequest)
	}

	subnets, vpcId, err := s.resolveSubnets(ctx, in.Region, in.SubnetIds)
	if err != nil {
		return nil, err
	}

	sg := &neptunestore.DBSubnetGroup{
		DBSubnetGroupName:        name,
		DBSubnetGroupDescription: in.DBSubnetGroupDescription,
		VpcId:                    vpcId,
		SubnetGroupStatus:        "Complete",
		Subnets:                  subnets,
		ARN:                      neptunestore.SubnetGroupARN(in.AccountID, in.Region, name),
	}

	if err := store.CreateSubnetGroup(sg); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"DBSubnetGroup": sg,
	}, nil
}

// deleteDBSubnetGroupCore deletes a subnet group and clears its tags.
func (s *NeptuneService) deleteDBSubnetGroupCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *DeleteDBSubnetGroupInput) (interface{}, error) {
	name := in.DBSubnetGroupName
	if name == "" {
		return nil, awserrors.NewMissingParameter("DBSubnetGroupName is required")
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

// modifyDBSubnetGroupCore updates the description and, when new subnets are
// supplied, re-resolves and replaces the subnet list.
func (s *NeptuneService) modifyDBSubnetGroupCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *ModifyDBSubnetGroupInput) (interface{}, error) {
	name := in.DBSubnetGroupName
	if name == "" {
		return nil, awserrors.NewMissingParameter("DBSubnetGroupName is required")
	}

	sg, err := store.GetSubnetGroup(name)
	if err != nil {
		return nil, translateStoreError(err)
	}

	if desc := in.DBSubnetGroupDescription; desc != "" {
		sg.DBSubnetGroupDescription = desc
	}
	if len(in.SubnetIds) > 0 {
		if err := validateSubnetCount(len(in.SubnetIds)); err != nil {
			return nil, awserrors.NewAWSError("DBSubnetGroupQuotaExceededFault", err.Error(), http.StatusBadRequest)
		}
		subnets, vpcId, err := s.resolveSubnets(ctx, in.Region, in.SubnetIds)
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
