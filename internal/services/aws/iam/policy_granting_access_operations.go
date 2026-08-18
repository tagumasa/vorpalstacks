package iam

import (
	"context"
	"regexp"
	"strings"

	"vorpalstacks/internal/common/iam/policy"
	"vorpalstacks/internal/common/request"
)

// serviceNamespacePattern matches the serviceNamespaceType constraint.
var serviceNamespacePattern = regexp.MustCompile(`^[\w-]*$`)

// policyGrantsServiceNamespace reports whether a permissions policy can
// grant access to the given service namespace at namespace granularity:
// an Allow statement whose Action list covers the namespace (the action's
// service prefix, or a global wildcard). An Allow statement phrased with
// NotAction allows everything except the listed actions, so it grants
// every namespace. Deny statements do not remove a policy from the
// results: the operation reports which policies can grant access, not the
// effective decision after explicit denies.
func policyGrantsServiceNamespace(doc *policy.Document, namespace string) bool {
	for i := range doc.Statement {
		stmt := &doc.Statement[i]
		if stmt.Effect != policy.EffectAllow {
			continue
		}
		if len(stmt.NotAction) > 0 {
			return true
		}
		for _, action := range stmt.Action {
			if action == "*" || action == "*:*" {
				return true
			}
			if idx := strings.Index(action, ":"); idx > 0 && action[:idx] == namespace {
				return true
			}
		}
	}
	return false
}

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
// user's groups are reported with the group as the attached entity; managed
// policies carry their ARN, inline policies carry the entity they are
// embedded in. Policies attached only as permissions boundaries are not
// returned. The request carries no page-size parameter, so a single
// response always holds every namespace entry.
func (s *IAMService) ListPoliciesGrantingServiceAccess(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetStringParam(req.Parameters, "Arn")
	if arn == "" {
		return nil, NewValidationError("Arn")
	}
	namespaces := request.GetStringList(req.Parameters, "ServiceNamespaces")
	if len(namespaces) == 0 {
		return nil, NewValidationError("ServiceNamespaces")
	}
	if len(namespaces) > 200 {
		return nil, NewInvalidInputError("ServiceNamespaces", "must contain 1 to 200 namespaces")
	}
	for _, ns := range namespaces {
		if len(ns) == 0 || len(ns) > 64 || !serviceNamespacePattern.MatchString(ns) {
			return nil, NewInvalidInputError("ServiceNamespaces", "each namespace must be 1 to 64 characters matching [\\w-]")
		}
	}

	records, err := s.gatherPrincipalPolicyRecords(reqCtx, arn, "Arn")
	if err != nil {
		return nil, err
	}

	entries := make([]map[string]interface{}, 0, len(namespaces))
	for _, ns := range namespaces {
		seen := map[string]bool{}
		granted := []map[string]interface{}{}
		for _, rec := range records {
			if !policyGrantsServiceNamespace(rec.Document, ns) {
				continue
			}
			if rec.PolicyArn != "" {
				// The same managed policy reached through several
				// attachment points is reported once per namespace.
				if seen[rec.PolicyArn] {
					continue
				}
				seen[rec.PolicyArn] = true
				granted = append(granted, map[string]interface{}{
					"PolicyName": rec.PolicyName,
					"PolicyType": "MANAGED",
					"PolicyArn":  rec.PolicyArn,
				})
			} else {
				granted = append(granted, map[string]interface{}{
					"PolicyName": rec.PolicyName,
					"PolicyType": "INLINE",
					"EntityType": principalWireEntityType(rec.EntityType),
					"EntityName": rec.EntityName,
				})
			}
		}
		entries = append(entries, map[string]interface{}{
			"ServiceNamespace": ns,
			"Policies":         granted,
		})
	}

	return map[string]interface{}{
		"PoliciesGrantingServiceAccess": entries,
		"IsTruncated":                   false,
	}, nil
}
