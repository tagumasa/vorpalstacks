package iam

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// principalWireEntityType maps an internal principal type to the
// policyOwnerEntityType enumeration used on the wire.
func principalWireEntityType(principalType string) string {
	switch principalType {
	case PrincipalTypeUser:
		return "USER"
	case PrincipalTypeGroup:
		return "GROUP"
	case PrincipalTypeRole:
		return "ROLE"
	}
	return ""
}

// ListPoliciesGrantingServiceAccess lists the permissions policies that let
// the specified identity (user, group, or role) access each requested
// service namespace. For a user, the managed and inline policies of the
// user's groups are reported with the group as the attached entity. The
// request carries no page-size parameter, so a single response always holds
// every namespace entry.
func (s *IAMService) ListPoliciesGrantingServiceAccess(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := &ListPoliciesGrantingServiceAccessInput{
		Arn:               request.GetStringParam(req.Parameters, "Arn"),
		ServiceNamespaces: request.GetStringList(req.Parameters, "ServiceNamespaces"),
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.listPoliciesGrantingServiceAccessCore(store, input)
	if err != nil {
		return nil, err
	}

	entries := make([]map[string]interface{}, 0, len(result.Entries))
	for _, entry := range result.Entries {
		policies := make([]map[string]interface{}, 0, len(entry.Policies))
		for _, p := range entry.Policies {
			wirePolicy := map[string]interface{}{
				"PolicyName": p.PolicyName,
			}
			if p.Managed {
				wirePolicy["PolicyType"] = "MANAGED"
				wirePolicy["PolicyArn"] = p.PolicyArn
			} else {
				wirePolicy["PolicyType"] = "INLINE"
				wirePolicy["EntityType"] = principalWireEntityType(p.EntityType)
				wirePolicy["EntityName"] = p.EntityName
			}
			policies = append(policies, wirePolicy)
		}
		entries = append(entries, map[string]interface{}{
			"ServiceNamespace": entry.ServiceNamespace,
			"Policies":         policies,
		})
	}

	return map[string]interface{}{
		"PoliciesGrantingServiceAccess": entries,
		"IsTruncated":                   false,
	}, nil
}
