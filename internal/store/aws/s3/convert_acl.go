package s3

import (
	pb "vorpalstacks/internal/pb/storage/storage_s3"
)

func aclOwnerToProto(o *ACLOwner) *pb.ACLOwner {
	if o == nil {
		return nil
	}
	return &pb.ACLOwner{Id: o.ID, DisplayName: o.DisplayName}
}

func protoToACLOwner(p *pb.ACLOwner) *ACLOwner {
	if p == nil {
		return nil
	}
	return &ACLOwner{ID: p.Id, DisplayName: p.DisplayName}
}

func granteeToProto(g *Grantee) *pb.Grantee {
	if g == nil {
		return nil
	}
	return &pb.Grantee{
		Type:        granteeTypeToProto(g.Type),
		Id:          g.ID,
		DisplayName: g.DisplayName,
		Email:       g.Email,
		Uri:         g.URI,
	}
}

func protoToGrantee(p *pb.Grantee) *Grantee {
	if p == nil {
		return nil
	}
	return &Grantee{
		Type:        protoToGranteeType(p.Type),
		ID:          p.Id,
		DisplayName: p.DisplayName,
		Email:       p.Email,
		URI:         p.Uri,
	}
}

func granteeTypeToProto(t GranteeType) pb.GranteeType {
	switch t {
	case GranteeTypeCanonicalUser:
		return pb.GranteeType_GRANTEE_TYPE_CANONICAL_USER
	case GranteeTypeGroup:
		return pb.GranteeType_GRANTEE_TYPE_GROUP
	case GranteeTypeAmazonCustomerByEmail:
		return pb.GranteeType_GRANTEE_TYPE_AMAZON_CUSTOMER_BY_EMAIL
	default:
		return pb.GranteeType_GRANTEE_TYPE_UNSPECIFIED
	}
}

func protoToGranteeType(p pb.GranteeType) GranteeType {
	switch p {
	case pb.GranteeType_GRANTEE_TYPE_CANONICAL_USER:
		return GranteeTypeCanonicalUser
	case pb.GranteeType_GRANTEE_TYPE_GROUP:
		return GranteeTypeGroup
	case pb.GranteeType_GRANTEE_TYPE_AMAZON_CUSTOMER_BY_EMAIL:
		return GranteeTypeAmazonCustomerByEmail
	default:
		return ""
	}
}

func permissionToProto(p Permission) pb.Permission {
	switch p {
	case PermissionFullControl:
		return pb.Permission_PERMISSION_FULL_CONTROL
	case PermissionRead:
		return pb.Permission_PERMISSION_READ
	case PermissionWrite:
		return pb.Permission_PERMISSION_WRITE
	case PermissionReadACP:
		return pb.Permission_PERMISSION_READ_ACP
	case PermissionWriteACP:
		return pb.Permission_PERMISSION_WRITE_ACP
	default:
		return pb.Permission_PERMISSION_UNSPECIFIED
	}
}

func protoToPermission(p pb.Permission) Permission {
	switch p {
	case pb.Permission_PERMISSION_FULL_CONTROL:
		return PermissionFullControl
	case pb.Permission_PERMISSION_READ:
		return PermissionRead
	case pb.Permission_PERMISSION_WRITE:
		return PermissionWrite
	case pb.Permission_PERMISSION_READ_ACP:
		return PermissionReadACP
	case pb.Permission_PERMISSION_WRITE_ACP:
		return PermissionWriteACP
	default:
		return ""
	}
}

func grantToProto(g *Grant) *pb.Grant {
	if g == nil {
		return nil
	}
	return &pb.Grant{
		Grantee:    granteeToProto(g.Grantee),
		Permission: permissionToProto(g.Permission),
	}
}

func protoToGrant(p *pb.Grant) *Grant {
	if p == nil {
		return nil
	}
	return &Grant{
		Grantee:    protoToGrantee(p.Grantee),
		Permission: protoToPermission(p.Permission),
	}
}

func grantsToProto(g []*Grant) []*pb.Grant {
	if g == nil {
		return nil
	}
	result := make([]*pb.Grant, len(g))
	for i, grant := range g {
		result[i] = grantToProto(grant)
	}
	return result
}

func protoToGrants(p []*pb.Grant) []*Grant {
	if p == nil {
		return nil
	}
	result := make([]*Grant, len(p))
	for i, grant := range p {
		result[i] = protoToGrant(grant)
	}
	return result
}

func accessControlPolicyToProto(a *AccessControlPolicy) *pb.AccessControlPolicy {
	if a == nil {
		return nil
	}
	return &pb.AccessControlPolicy{
		Owner:  aclOwnerToProto(a.Owner),
		Grants: grantsToProto(a.Grants),
	}
}

func protoToAccessControlPolicy(p *pb.AccessControlPolicy) *AccessControlPolicy {
	if p == nil {
		return nil
	}
	return &AccessControlPolicy{
		Owner:  protoToACLOwner(p.Owner),
		Grants: protoToGrants(p.Grants),
	}
}
