// Transport-agnostic Core function for ListPoliciesGrantingServiceAccess:
// the namespace-grant scan shared by the AWS-compatible HTTP API handlers
// and any admin plane paths (the xxxCore pattern).
package iam

import (
	"fmt"
	"regexp"
	"strings"

	"vorpalstacks/internal/common/iam/policy"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// serviceNamespacePattern matches the serviceNamespaceType constraint.
var serviceNamespacePattern = regexp.MustCompile(`^[\w-]*$`)

// ListPoliciesGrantingServiceAccessInput holds the parameters for listing
// the policies that grant a principal access to service namespaces.
type ListPoliciesGrantingServiceAccessInput struct {
	Arn               string
	ServiceNamespaces []string
}

// GrantingAccessPolicy is one policy reported as granting access to a
// service namespace, before wire serialisation.
type GrantingAccessPolicy struct {
	PolicyName string
	Managed    bool
	PolicyArn  string // managed policies only
	EntityType string // internal principal type, inline policies only
	EntityName string // inline policies only
}

// GrantingServiceAccessEntry is one service namespace and the policies that
// grant access to it.
type GrantingServiceAccessEntry struct {
	ServiceNamespace string
	Policies         []GrantingAccessPolicy
}

// ListPoliciesGrantingServiceAccessResult holds the entries produced by
// listPoliciesGrantingServiceAccessCore.
type ListPoliciesGrantingServiceAccessResult struct {
	Entries []GrantingServiceAccessEntry
}

// listPoliciesGrantingServiceAccessCore lists the permissions policies that
// let the specified identity (user, group, or role) access each requested
// service namespace. For a user, the managed and inline policies of the
// user's groups are reported with the group as the attached entity. Managed
// policies carry their ARN; inline policies carry the entity they are
// embedded in. Policies attached only as permissions boundaries are never
// returned. The request carries no page-size parameter, so a single result
// always holds every namespace entry.
func (s *IAMService) listPoliciesGrantingServiceAccessCore(store *iamstore.IAMStore, input *ListPoliciesGrantingServiceAccessInput) (*ListPoliciesGrantingServiceAccessResult, error) {
	if input.Arn == "" {
		return nil, NewValidationError("Arn")
	}
	if len(input.ServiceNamespaces) == 0 {
		return nil, NewValidationError("ServiceNamespaces")
	}
	if len(input.ServiceNamespaces) > MaxServiceNamespaces {
		return nil, NewInvalidInputError("ServiceNamespaces", fmt.Sprintf("must contain 1 to %d namespaces", MaxServiceNamespaces))
	}
	for _, ns := range input.ServiceNamespaces {
		if len(ns) == 0 || len(ns) > MaxServiceNamespaceLength || !serviceNamespacePattern.MatchString(ns) {
			return nil, NewInvalidInputError("ServiceNamespaces", fmt.Sprintf("each namespace must be 1 to %d characters matching [\\w-]", MaxServiceNamespaceLength))
		}
	}

	records, err := s.gatherPrincipalPolicyRecordsCore(store, input.Arn, "Arn")
	if err != nil {
		return nil, err
	}

	result := &ListPoliciesGrantingServiceAccessResult{
		Entries: make([]GrantingServiceAccessEntry, 0, len(input.ServiceNamespaces)),
	}
	for _, ns := range input.ServiceNamespaces {
		seen := map[string]bool{}
		granted := []GrantingAccessPolicy{}
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
				granted = append(granted, GrantingAccessPolicy{
					PolicyName: rec.PolicyName,
					Managed:    true,
					PolicyArn:  rec.PolicyArn,
				})
			} else {
				granted = append(granted, GrantingAccessPolicy{
					PolicyName: rec.PolicyName,
					EntityType: rec.EntityType,
					EntityName: rec.EntityName,
				})
			}
		}
		result.Entries = append(result.Entries, GrantingServiceAccessEntry{
			ServiceNamespace: ns,
			Policies:         granted,
		})
	}

	return result, nil
}

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
