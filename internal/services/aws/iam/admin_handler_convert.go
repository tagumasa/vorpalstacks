package iam

import (
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/timeutils"

	"google.golang.org/protobuf/proto"

	pb "vorpalstacks/internal/pb/aws/iam"
)

// getStore returns the global IAM store for admin console operations.
// IAM is a global service; the region parameter is ignored.
func (h *AdminHandler) getStore() (*iamstore.IAMStore, error) {
	return h.service.GetStoreForRegion("")
}

// toPbUser converts a store-layer User to the proto representation.
func toPbUser(user *iamstore.User) *pb.User {
	pbUser := &pb.User{
		Username:   user.UserName,
		Userid:     user.ID,
		Arn:        user.Arn,
		Path:       user.Path,
		Createdate: user.CreateDate.Format(timeutils.ISO8601UTCFormat),
	}

	if user.PasswordLastUsed != nil {
		pbUser.Passwordlastused = user.PasswordLastUsed.Format(timeutils.ISO8601UTCFormat)
	}

	if user.PermissionsBoundary != nil {
		pbUser.Permissionsboundary = &pb.AttachedPermissionsBoundary{
			Permissionsboundaryarn:  user.PermissionsBoundary.PermissionsBoundaryArn,
			Permissionsboundarytype: pb.PermissionsBoundaryAttachmentType_PERMISSIONS_BOUNDARY_ATTACHMENT_TYPE_POLICY,
		}
	}

	if len(user.Tags) > 0 {
		pbUser.Tags = make([]*pb.Tag, len(user.Tags))
		for i, tag := range user.Tags {
			pbUser.Tags[i] = &pb.Tag{Key: tag.Key, Value: tag.Value}
		}
	}

	return pbUser
}

// toPbRole converts a store-layer Role to the proto representation.
func toPbRole(role *iamstore.Role) *pb.Role {
	pbRole := &pb.Role{
		Rolename:                 role.RoleName,
		Roleid:                   role.ID,
		Arn:                      role.Arn,
		Path:                     role.Path,
		Createdate:               role.CreateDate.Format(timeutils.ISO8601UTCFormat),
		Assumerolepolicydocument: role.AssumeRolePolicyDocument,
		Description:              role.Description,
		Maxsessionduration:       proto.Int32(int32(role.MaxSessionDuration)),
	}

	if role.PermissionsBoundary != nil {
		pbRole.Permissionsboundary = &pb.AttachedPermissionsBoundary{
			Permissionsboundaryarn:  role.PermissionsBoundary.PermissionsBoundaryArn,
			Permissionsboundarytype: pb.PermissionsBoundaryAttachmentType_PERMISSIONS_BOUNDARY_ATTACHMENT_TYPE_POLICY,
		}
	}

	if role.RoleLastUsed != nil {
		pbRole.Rolelastused = &pb.RoleLastUsed{
			Region: role.RoleLastUsed.Region,
		}
		if role.RoleLastUsed.LastUsedDate != nil {
			pbRole.Rolelastused.Lastuseddate = role.RoleLastUsed.LastUsedDate.Format(timeutils.ISO8601UTCFormat)
		}
	}

	if len(role.Tags) > 0 {
		pbRole.Tags = make([]*pb.Tag, len(role.Tags))
		for i, tag := range role.Tags {
			pbRole.Tags[i] = &pb.Tag{Key: tag.Key, Value: tag.Value}
		}
	}

	return pbRole
}

// toPbPolicy converts a store-layer Policy to the proto representation.
func toPbPolicy(policy *iamstore.Policy) *pb.Policy {
	pbPolicy := &pb.Policy{
		Policyname:                    policy.PolicyName,
		Policyid:                      policy.ID,
		Arn:                           policy.Arn,
		Path:                          policy.Path,
		Createdate:                    policy.CreateDate.Format(timeutils.ISO8601UTCFormat),
		Updatedate:                    policy.UpdateDate.Format(timeutils.ISO8601UTCFormat),
		Defaultversionid:              policy.DefaultVersionId,
		Attachmentcount:               proto.Int32(int32(policy.AttachmentCount)),
		Permissionsboundaryusagecount: proto.Int32(int32(policy.PermissionsBoundaryUsageCount)),
		Isattachable:                  proto.Bool(policy.IsAttachable),
		Description:                   policy.Description,
	}

	if len(policy.Tags) > 0 {
		pbPolicy.Tags = make([]*pb.Tag, len(policy.Tags))
		for i, tag := range policy.Tags {
			pbPolicy.Tags[i] = &pb.Tag{Key: tag.Key, Value: tag.Value}
		}
	}

	return pbPolicy
}

// toPbGroup converts a store-layer Group to the proto representation.
func toPbGroup(group *iamstore.Group) *pb.Group {
	pbGroup := &pb.Group{
		Groupid:    group.ID,
		Groupname:  group.GroupName,
		Arn:        group.Arn,
		Path:       group.Path,
		Createdate: group.CreateDate.Format(timeutils.ISO8601UTCFormat),
	}

	return pbGroup
}
