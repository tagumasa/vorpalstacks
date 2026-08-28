// Transport-agnostic Core functions for IAM managed policies and policy
// versions: validation and store operations shared by the AWS-compatible
// HTTP API handlers and the admin gRPC-Web handler (the xxxCore pattern).
package iam

import (
	"errors"
	"unicode/utf8"

	"vorpalstacks/internal/common/iam/policy"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// CreatePolicyInput holds the parameters for creating an IAM managed policy.
type CreatePolicyInput struct {
	PolicyName     string
	Path           string
	PolicyDocument string
	Description    string
	Tags           []tags.Tag
}

// DeletePolicyInput holds the parameters for deleting an IAM managed policy.
// Both AWS API and admin handler enforce the AttachmentCount > 0 check.
type DeletePolicyInput struct {
	PolicyArn string
}

// createPolicyCore validates input and creates an IAM managed policy in the
// store.  Returns the created policy or an IAM-formatted error.
func (s *IAMService) createPolicyCore(store *iamstore.IAMStore, input *CreatePolicyInput) (*iamstore.Policy, error) {
	if input.PolicyName == "" {
		return nil, NewInvalidInputError("PolicyName", "cannot be empty")
	}
	if err := validateEntityName128(input.PolicyName, "PolicyName"); err != nil {
		return nil, err
	}

	path := input.Path
	if path == "" {
		path = "/"
	}
	if !validatePath(path) {
		return nil, NewInvalidInputError("Path", "must be a valid path starting and ending with /")
	}

	if !validatePolicyDocument(input.PolicyDocument) {
		return nil, ErrMalformedPolicyDocument
	}

	if utf8.RuneCountInString(input.Description) > maxPolicyDescriptionLength {
		return nil, NewInvalidInputError("Description", "must be 0 to 1000 characters")
	}

	if err := validateNewTags(input.Tags); err != nil {
		return nil, err
	}

	policy, err := store.Policies().Create(input.PolicyName, path, store.AccountID(), input.PolicyDocument, input.Description, input.Tags)
	if err != nil {
		if errors.Is(err, iamstore.ErrPolicyAlreadyExists) {
			return nil, NewPolicyAlreadyExistsError(input.PolicyName)
		}
		return nil, err
	}
	return policy, nil
}

// getPolicyCore returns the IAM managed policy with the given ARN.
// Callers must validate that policyArn is non-empty before calling.
func (s *IAMService) getPolicyCore(store *iamstore.IAMStore, policyArn string) (*iamstore.Policy, error) {
	policy, err := store.Policies().Get(policyArn)
	if err != nil {
		return nil, NewNoSuchPolicyError(policyArn)
	}
	return policy, nil
}

// listPoliciesCore returns a paginated list of IAM managed policies.
func (s *IAMService) listPoliciesCore(store *iamstore.IAMStore, scope, pathPrefix, marker string, onlyAttached bool, maxItems int) (*iamstore.PolicyListResult, error) {
	return store.Policies().List(scope, pathPrefix, onlyAttached, marker, maxItems)
}

// deletePolicyCore validates input and deletes an IAM managed policy.
// Rejects if the policy has active attachments (AttachmentCount > 0).
// All policy versions are cleaned up before the policy record is removed.
func (s *IAMService) deletePolicyCore(store *iamstore.IAMStore, input *DeletePolicyInput) error {
	if input.PolicyArn == "" {
		return ErrNoSuchPolicy
	}

	policy, err := store.Policies().Get(input.PolicyArn)
	if err != nil {
		return NewNoSuchPolicyError(input.PolicyArn)
	}

	// The permissions defined in AWS managed policies cannot be changed,
	// so the policies themselves cannot be deleted. A missing policy is
	// reported before modifiability so that an AWS-managed ARN that is
	// not present yields NoSuchEntity.
	if iamstore.IsAWSManagedPolicyARN(input.PolicyArn) {
		return NewInvalidInputError("PolicyArn", "AWS managed policies cannot be modified")
	}

	if policy.AttachmentCount > 0 {
		return NewDeletePolicyConflictError(input.PolicyArn)
	}

	if err := store.Policies().DeleteAllVersions(input.PolicyArn); err != nil {
		return err
	}

	return store.Policies().Delete(input.PolicyArn)
}

// CreatePolicyVersionInput holds the parameters for creating a managed
// policy version.
type CreatePolicyVersionInput struct {
	PolicyArn      string
	PolicyDocument string
	SetAsDefault   bool
}

// createPolicyVersionCore validates input and creates a new version of a
// managed policy. A missing policy is reported before modifiability so that
// an AWS-managed ARN that is not present yields NoSuchEntity. CreateVersion
// atomically enforces the policy version quota and performs the
// default-version swap inside a single lock scope, eliminating the race
// condition where concurrent requests could both observe a version count
// below the limit.
func (s *IAMService) createPolicyVersionCore(store *iamstore.IAMStore, input *CreatePolicyVersionInput) (*iamstore.PolicyVersion, error) {
	if input.PolicyArn == "" {
		return nil, NewValidationError("PolicyArn")
	}

	if !validatePolicyDocument(input.PolicyDocument) {
		return nil, ErrMalformedPolicyDocument
	}

	if !store.Policies().Exists(input.PolicyArn) {
		return nil, NewNoSuchPolicyError(input.PolicyArn)
	}
	// The permissions defined in AWS managed policies cannot be changed.
	if iamstore.IsAWSManagedPolicyARN(input.PolicyArn) {
		return nil, NewInvalidInputError("PolicyArn", "AWS managed policies cannot be modified")
	}

	version, err := store.Policies().CreateVersion(input.PolicyArn, input.PolicyDocument, input.SetAsDefault, iamstore.MaxPolicyVersions)
	if err != nil {
		if errors.Is(err, iamstore.ErrPolicyVersionLimitExceeded) {
			return nil, ErrLimitExceededPolicyVersions
		}
		return nil, err
	}
	return version, nil
}

// getPolicyVersionCore retrieves a specific policy version.
func (s *IAMService) getPolicyVersionCore(store *iamstore.IAMStore, policyArn, versionId string) (*iamstore.PolicyVersion, error) {
	version, err := store.Policies().GetVersion(policyArn, versionId)
	if err != nil {
		return nil, NewNoSuchPolicyVersionError(versionId)
	}
	return version, nil
}

// deletePolicyVersionCore deletes a non-default policy version.
func (s *IAMService) deletePolicyVersionCore(store *iamstore.IAMStore, policyArn, versionId string) error {
	policy, err := store.Policies().Get(policyArn)
	if err != nil {
		return NewNoSuchPolicyError(policyArn)
	}
	// A missing policy is reported before modifiability so that an
	// AWS-managed ARN that is not present yields NoSuchEntity.
	if iamstore.IsAWSManagedPolicyARN(policyArn) {
		return NewInvalidInputError("PolicyArn", "AWS managed policies cannot be modified")
	}

	if policy.DefaultVersionId == versionId {
		return NewDeleteConflictError("Cannot delete the default policy version.")
	}

	if err := store.Policies().DeleteVersion(policyArn, versionId); err != nil {
		return NewNoSuchPolicyVersionError(versionId)
	}
	return nil
}

// listPolicyVersionsCore returns a paginated list of policy versions.
func (s *IAMService) listPolicyVersionsCore(store *iamstore.IAMStore, policyArn, marker string, maxItems int) (*iamstore.PolicyVersionListResult, error) {
	if !store.Policies().Exists(policyArn) {
		return nil, NewNoSuchPolicyError(policyArn)
	}
	return store.Policies().ListVersions(policyArn, marker, maxItems)
}

// setDefaultPolicyVersionCore sets the default version for a policy.
func (s *IAMService) setDefaultPolicyVersionCore(store *iamstore.IAMStore, policyArn, versionId string) error {
	if !store.Policies().Exists(policyArn) {
		return NewNoSuchPolicyError(policyArn)
	}
	// A missing policy is reported before modifiability so that an
	// AWS-managed ARN that is not present yields NoSuchEntity.
	if iamstore.IsAWSManagedPolicyARN(policyArn) {
		return NewInvalidInputError("PolicyArn", "AWS managed policies cannot be modified")
	}
	if err := store.Policies().SetDefaultVersion(policyArn, versionId); err != nil {
		return NewNoSuchPolicyVersionError(versionId)
	}
	return nil
}

// ListEntitiesForPolicyResult holds the aggregated entity lists returned
// by listEntitiesForPolicyCore.
type ListEntitiesForPolicyResult struct {
	PolicyUsers  []map[string]interface{}
	PolicyGroups []map[string]interface{}
	PolicyRoles  []map[string]interface{}
	IsTruncated  bool
	Marker       string
}

// listEntitiesForPolicyCore lists all principals that the specified
// managed policy is attached to, with optional entity-type filtering
// and cross-type pagination.
func (s *IAMService) listEntitiesForPolicyCore(store *iamstore.IAMStore, policyArn, entityFilter, marker string, maxItems int) (*ListEntitiesForPolicyResult, error) {
	if !store.Policies().Exists(policyArn) {
		return nil, NewNoSuchPolicyError(policyArn)
	}

	refs, err := store.AttachedPolicies().ListPrincipalsForPolicy(policyArn)
	if err != nil {
		return nil, err
	}

	type entityEntry struct {
		entityType string
		name       string
		data       map[string]interface{}
	}

	combined := make([]entityEntry, 0)

	for _, ref := range refs {
		switch ref.PrincipalType {
		case PrincipalTypeUser:
			if entityFilter != "" && entityFilter != "User" {
				continue
			}
			if user, err := store.Users().Get(ref.PrincipalName); err == nil {
				entry := map[string]interface{}{
					"UserName": user.UserName,
					"UserId":   user.ID,
					"Arn":      user.Arn,
				}
				combined = append(combined, entityEntry{"User", user.UserName, entry})
			}
		case PrincipalTypeGroup:
			if entityFilter != "" && entityFilter != "Group" {
				continue
			}
			if group, err := store.Groups().Get(ref.PrincipalName); err == nil {
				entry := map[string]interface{}{
					"GroupName": group.GroupName,
					"GroupId":   group.ID,
					"Arn":       group.Arn,
				}
				combined = append(combined, entityEntry{"Group", group.GroupName, entry})
			}
		case PrincipalTypeRole:
			if entityFilter != "" && entityFilter != "Role" {
				continue
			}
			if role, err := store.Roles().Get(ref.PrincipalName); err == nil {
				entry := map[string]interface{}{
					"RoleName": role.RoleName,
					"RoleId":   role.ID,
					"Arn":      role.Arn,
				}
				combined = append(combined, entityEntry{"Role", role.RoleName, entry})
			}
		}
	}

	paged := pagination.PaginateSlice(combined, marker, maxItems, func(item entityEntry) string {
		return item.entityType + ":" + item.name
	})

	result := &ListEntitiesForPolicyResult{
		PolicyUsers:  make([]map[string]interface{}, 0),
		PolicyGroups: make([]map[string]interface{}, 0),
		PolicyRoles:  make([]map[string]interface{}, 0),
		IsTruncated:  paged.IsTruncated,
		Marker:       paged.NextMarker,
	}

	for _, entry := range paged.Items {
		switch entry.entityType {
		case "User":
			result.PolicyUsers = append(result.PolicyUsers, entry.data)
		case "Group":
			result.PolicyGroups = append(result.PolicyGroups, entry.data)
		case "Role":
			result.PolicyRoles = append(result.PolicyRoles, entry.data)
		}
	}

	return result, nil
}

// principalPolicy pairs a parsed permissions policy with the identity of
// its attachment point: the managed policy ARN for attached policies, or
// the owning entity for inline policies.
type principalPolicy struct {
	Document   *policy.Document
	PolicyName string
	PolicyArn  string // empty for inline policies
	EntityType string // PrincipalTypeUser, PrincipalTypeGroup, or PrincipalTypeRole
	EntityName string
}

// gatherPrincipalPoliciesCore collects all identity-based policies
// applicable to the given principal.
func (s *IAMService) gatherPrincipalPoliciesCore(store *iamstore.IAMStore, principalArn string) ([]*policy.Document, error) {
	records, err := s.gatherPrincipalPolicyRecordsCore(store, principalArn, "PolicySourceArn")
	if err != nil {
		return nil, err
	}
	docs := make([]*policy.Document, 0, len(records))
	for _, rec := range records {
		docs = append(docs, rec.Document)
	}
	return docs, nil
}

// gatherPrincipalPolicyRecordsCore collects the permissions policies that
// apply to the principal identified by the ARN, preserving where each
// policy is attached. For a user this includes the managed and inline
// policies of every group the user belongs to. Permissions boundaries are
// stored separately from permission attachments, so they are never
// collected. paramName names the request parameter the ARN was taken from,
// for error messages.
func (s *IAMService) gatherPrincipalPolicyRecordsCore(store *iamstore.IAMStore, principalArn, paramName string) ([]principalPolicy, error) {
	entityType := resolveEntityType(principalArn)
	entityName := resolveEntityName(principalArn)

	var records []principalPolicy

	switch entityType {
	case "User":
		if !store.Users().Exists(entityName) {
			return nil, NewNoSuchUserError(entityName)
		}
		records = append(records, collectInlinePolicies(store, PrincipalTypeUser, entityName)...)
		records = append(records, collectAttachedPolicies(store, PrincipalTypeUser, entityName)...)
		groupNames, err := store.UserGroups().ListGroupsForUser(entityName)
		if err != nil {
			return nil, err
		}
		for _, groupName := range groupNames {
			records = append(records, collectInlinePolicies(store, PrincipalTypeGroup, groupName)...)
			records = append(records, collectAttachedPolicies(store, PrincipalTypeGroup, groupName)...)
		}

	case "Role":
		if !store.Roles().Exists(entityName) {
			return nil, NewNoSuchRoleError(entityName)
		}
		records = append(records, collectInlinePolicies(store, PrincipalTypeRole, entityName)...)
		records = append(records, collectAttachedPolicies(store, PrincipalTypeRole, entityName)...)

	case "Group":
		if !store.Groups().Exists(entityName) {
			return nil, NewNoSuchGroupError(entityName)
		}
		records = append(records, collectInlinePolicies(store, PrincipalTypeGroup, entityName)...)
		records = append(records, collectAttachedPolicies(store, PrincipalTypeGroup, entityName)...)

	default:
		// Unknown or non-principal ARN (e.g. policy, server-certificate).
		// Fail-closed instead of silently returning an empty policy list.
		return nil, NewInvalidInputError(paramName, "must be a user, role, or group ARN")
	}

	return records, nil
}

func collectInlinePolicies(store *iamstore.IAMStore, principalType, principalName string) []principalPolicy {
	policyNames, err := store.InlinePolicies().List(principalType, principalName)
	if err != nil {
		return nil
	}
	var records []principalPolicy
	for _, pn := range policyNames {
		ip, err := store.InlinePolicies().Get(principalType, principalName, pn)
		if err != nil || ip == nil {
			continue
		}
		doc, err := policy.ParseDocument(ip.PolicyDocument)
		if err != nil {
			continue
		}
		records = append(records, principalPolicy{
			Document:   doc,
			PolicyName: pn,
			EntityType: principalType,
			EntityName: principalName,
		})
	}
	return records
}

func collectAttachedPolicies(store *iamstore.IAMStore, principalType, principalName string) []principalPolicy {
	arns, err := store.AttachedPolicies().ListAttachedPolicies(principalType, principalName)
	if err != nil {
		return nil
	}
	var records []principalPolicy
	for _, arn := range arns {
		version, err := store.Policies().GetDefaultVersion(arn)
		if err != nil || version == nil {
			continue
		}
		doc, err := policy.ParseDocument(version.Document)
		if err != nil {
			continue
		}
		policyName := arn
		if p, err := store.Policies().Get(arn); err == nil && p != nil {
			policyName = p.PolicyName
		}
		records = append(records, principalPolicy{
			Document:   doc,
			PolicyName: policyName,
			PolicyArn:  arn,
		})
	}
	return records
}
