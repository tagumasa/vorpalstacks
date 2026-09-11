// Transport-agnostic Core functions for IAM managed policies and policy
// versions: validation and store operations shared by the AWS-compatible
// HTTP API handlers and the admin gRPC-Web handler (the xxxCore pattern).
package iam

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"vorpalstacks/internal/common/iam/policy"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	iamstore "vorpalstacks/internal/store/aws/iam"
	awsarn "vorpalstacks/internal/utils/aws/arn"
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

// getPolicyCore returns the IAM managed policy with the given ARN; an
// empty ARN is rejected as a validation error.
func (s *IAMService) getPolicyCore(store *iamstore.IAMStore, policyArn string) (*iamstore.Policy, error) {
	if policyArn == "" {
		return nil, NewValidationError("PolicyArn")
	}
	policy, err := store.Policies().Get(policyArn)
	if err != nil {
		return nil, storeReadError(err, iamstore.ErrPolicyNotFound, NewNoSuchPolicyError(policyArn))
	}
	return policy, nil
}

// listPoliciesCore returns a paginated list of IAM managed policies. An
// omitted Scope defaults to All, per the ListPolicies contract: "If it is
// not included, or if it is set to All, all policies are returned."
func (s *IAMService) listPoliciesCore(store *iamstore.IAMStore, scope, pathPrefix, marker string, onlyAttached bool, maxItems int) (*iamstore.PolicyListResult, error) {
	if scope == "" {
		scope = "All"
	}
	if !validatePolicyScope(scope) {
		return nil, NewInvalidInputError("Scope", "must be one of: All, AWS, Local")
	}
	return store.Policies().List(scope, pathPrefix, onlyAttached, marker, maxItems)
}

// deletePolicyCore validates input and deletes an IAM managed policy.
// Rejects if the policy has active attachments (AttachmentCount > 0).
// All policy versions are cleaned up before the policy record is removed.
func (s *IAMService) deletePolicyCore(store *iamstore.IAMStore, input *DeletePolicyInput) error {
	if input.PolicyArn == "" {
		return NewValidationError("PolicyArn")
	}

	policy, err := store.Policies().Get(input.PolicyArn)
	if err != nil {
		return storeReadError(err, iamstore.ErrPolicyNotFound, NewNoSuchPolicyError(input.PolicyArn))
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

	// CreateVersion does not signal a missing policy (its version scan of a
	// phantom policy is empty), so the policy record is resolved here; the
	// resolved read keeps an outage from masquerading as a missing policy.
	if _, err := store.Policies().Get(input.PolicyArn); err != nil {
		return nil, storeReadError(err, iamstore.ErrPolicyNotFound, NewNoSuchPolicyError(input.PolicyArn))
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

// getPolicyVersionCore retrieves a specific policy version; an empty ARN
// is rejected as a validation error before the store is consulted.
func (s *IAMService) getPolicyVersionCore(store *iamstore.IAMStore, policyArn, versionId string) (*iamstore.PolicyVersion, error) {
	if policyArn == "" {
		return nil, NewValidationError("PolicyArn")
	}
	version, err := store.Policies().GetVersion(policyArn, versionId)
	if err != nil {
		return nil, storeReadError(err, iamstore.ErrPolicyNotFound, NewNoSuchPolicyVersionError(versionId))
	}
	return version, nil
}

// deletePolicyVersionCore deletes a non-default policy version; an empty
// ARN is rejected as a validation error before the store is consulted.
func (s *IAMService) deletePolicyVersionCore(store *iamstore.IAMStore, policyArn, versionId string) error {
	if policyArn == "" {
		return NewValidationError("PolicyArn")
	}
	policy, err := store.Policies().Get(policyArn)
	if err != nil {
		return storeReadError(err, iamstore.ErrPolicyNotFound, NewNoSuchPolicyError(policyArn))
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
		return storeReadError(err, iamstore.ErrPolicyNotFound, NewNoSuchPolicyVersionError(versionId))
	}
	return nil
}

// listPolicyVersionsCore returns a paginated list of policy versions; an
// empty ARN is rejected as a validation error before the store is
// consulted.
func (s *IAMService) listPolicyVersionsCore(store *iamstore.IAMStore, policyArn, marker string, maxItems int) (*iamstore.PolicyVersionListResult, error) {
	if policyArn == "" {
		return nil, NewValidationError("PolicyArn")
	}
	// ListVersions cannot signal a missing policy (an empty scan is a valid
	// result), so the policy record is resolved first; the resolved read
	// keeps an outage from masquerading as a missing policy.
	if _, err := store.Policies().Get(policyArn); err != nil {
		return nil, storeReadError(err, iamstore.ErrPolicyNotFound, NewNoSuchPolicyError(policyArn))
	}
	return store.Policies().ListVersions(policyArn, marker, maxItems)
}

// setDefaultPolicyVersionCore sets the default version for a policy; an
// empty ARN is rejected as a validation error before the store is
// consulted.
func (s *IAMService) setDefaultPolicyVersionCore(store *iamstore.IAMStore, policyArn, versionId string) error {
	if policyArn == "" {
		return NewValidationError("PolicyArn")
	}
	// The policy record is resolved rather than probed: the resolved read
	// keeps an outage from masquerading as a missing policy.
	if _, err := store.Policies().Get(policyArn); err != nil {
		return storeReadError(err, iamstore.ErrPolicyNotFound, NewNoSuchPolicyError(policyArn))
	}
	// A missing policy is reported before modifiability so that an
	// AWS-managed ARN that is not present yields NoSuchEntity.
	if iamstore.IsAWSManagedPolicyARN(policyArn) {
		return NewInvalidInputError("PolicyArn", "AWS managed policies cannot be modified")
	}
	if err := store.Policies().SetDefaultVersion(policyArn, versionId); err != nil {
		return storeReadError(err, iamstore.ErrPolicyNotFound, NewNoSuchPolicyVersionError(versionId))
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
// and cross-type pagination. An empty ARN is rejected as a validation
// error before the store is consulted.
func (s *IAMService) listEntitiesForPolicyCore(store *iamstore.IAMStore, policyArn, entityFilter, marker string, maxItems int) (*ListEntitiesForPolicyResult, error) {
	if policyArn == "" {
		return nil, NewValidationError("PolicyArn")
	}
	// ListPrincipalsForPolicy cannot signal a missing policy (an empty scan
	// is a valid result), so the policy record is resolved first; the
	// resolved read keeps an outage from masquerading as a missing policy.
	if _, err := store.Policies().Get(policyArn); err != nil {
		return nil, storeReadError(err, iamstore.ErrPolicyNotFound, NewNoSuchPolicyError(policyArn))
	}

	refs, err := store.AttachedPolicies().ListPrincipalsForPolicy(policyArn)
	if err != nil {
		return nil, storeListError(err)
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
			// A listed attachment whose entity cannot be read — an
			// infrastructure fault, or an entity that vanished between the
			// attachment listing and this read — fails the response instead
			// of silently omitting the entity.
			user, err := store.Users().Get(ref.PrincipalName)
			if err != nil {
				return nil, storeListError(err)
			}
			entry := map[string]interface{}{
				"UserName": user.UserName,
				"UserId":   user.ID,
				"Arn":      user.Arn,
			}
			combined = append(combined, entityEntry{"User", user.UserName, entry})
		case PrincipalTypeGroup:
			if entityFilter != "" && entityFilter != "Group" {
				continue
			}
			group, err := store.Groups().Get(ref.PrincipalName)
			if err != nil {
				return nil, storeListError(err)
			}
			entry := map[string]interface{}{
				"GroupName": group.GroupName,
				"GroupId":   group.ID,
				"Arn":       group.Arn,
			}
			combined = append(combined, entityEntry{"Group", group.GroupName, entry})
		case PrincipalTypeRole:
			if entityFilter != "" && entityFilter != "Role" {
				continue
			}
			role, err := store.Roles().Get(ref.PrincipalName)
			if err != nil {
				return nil, storeListError(err)
			}
			entry := map[string]interface{}{
				"RoleName": role.RoleName,
				"RoleId":   role.ID,
				"Arn":      role.Arn,
			}
			combined = append(combined, entityEntry{"Role", role.RoleName, entry})
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
		// The entity is resolved rather than probed: the resolved read
		// keeps an outage from masquerading as a missing entity.
		if _, err := store.Users().Get(entityName); err != nil {
			return nil, storeReadError(err, iamstore.ErrUserNotFound, NewNoSuchUserError(entityName))
		}
		var err error
		records, err = collectInlinePolicies(records, store, PrincipalTypeUser, entityName)
		if err != nil {
			return nil, err
		}
		records, err = collectAttachedPolicies(records, store, PrincipalTypeUser, entityName)
		if err != nil {
			return nil, err
		}
		groupNames, err := store.UserGroups().ListGroupsForUser(entityName)
		if err != nil {
			return nil, storeListError(err)
		}
		for _, groupName := range groupNames {
			records, err = collectInlinePolicies(records, store, PrincipalTypeGroup, groupName)
			if err != nil {
				return nil, err
			}
			records, err = collectAttachedPolicies(records, store, PrincipalTypeGroup, groupName)
			if err != nil {
				return nil, err
			}
		}

	case "Role":
		if _, err := store.Roles().Get(entityName); err != nil {
			return nil, storeReadError(err, iamstore.ErrRoleNotFound, NewNoSuchRoleError(entityName))
		}
		var err error
		records, err = collectInlinePolicies(records, store, PrincipalTypeRole, entityName)
		if err != nil {
			return nil, err
		}
		records, err = collectAttachedPolicies(records, store, PrincipalTypeRole, entityName)
		if err != nil {
			return nil, err
		}

	case "Group":
		if _, err := store.Groups().Get(entityName); err != nil {
			return nil, storeReadError(err, iamstore.ErrGroupNotFound, NewNoSuchGroupError(entityName))
		}
		var err error
		records, err = collectInlinePolicies(records, store, PrincipalTypeGroup, entityName)
		if err != nil {
			return nil, err
		}
		records, err = collectAttachedPolicies(records, store, PrincipalTypeGroup, entityName)
		if err != nil {
			return nil, err
		}

	default:
		// Unknown or non-principal ARN (e.g. policy, server-certificate).
		// Fail-closed instead of silently returning an empty policy list.
		return nil, NewInvalidInputError(paramName, "must be a user, role, or group ARN")
	}

	return records, nil
}

// collectInlinePolicies appends the principal's inline policies to records.
// Store read failures abort the collection — the evaluated policy set must
// never be silently narrowed. A stored document that fails to parse is
// skipped (fail-closed for that policy: it can withhold permissions, never
// grant them) with a warning naming the principal and the policy.
func collectInlinePolicies(into []principalPolicy, store *iamstore.IAMStore, principalType, principalName string) ([]principalPolicy, error) {
	policyNames, err := store.InlinePolicies().List(principalType, principalName)
	if err != nil {
		return nil, storeListError(err)
	}
	for _, pn := range policyNames {
		ip, err := store.InlinePolicies().Get(principalType, principalName, pn)
		if err != nil {
			return nil, storeListError(err)
		}
		doc, err := policy.ParseDocument(ip.PolicyDocument)
		if err != nil {
			logs.Warn("skipping unparseable inline policy document in policy collection",
				logs.String("principalType", principalType),
				logs.String("principal", principalName),
				logs.String("policy", pn),
				logs.Err(err))
			continue
		}
		into = append(into, principalPolicy{
			Document:   doc,
			PolicyName: pn,
			EntityType: principalType,
			EntityName: principalName,
		})
	}
	return into, nil
}

// collectAttachedPolicies appends the managed policies attached to the
// principal to records, with the same read-failure semantics as
// collectInlinePolicies: any store failure aborts; an unparseable document
// is skipped fail-closed with a warning.
func collectAttachedPolicies(into []principalPolicy, store *iamstore.IAMStore, principalType, principalName string) ([]principalPolicy, error) {
	arns, err := store.AttachedPolicies().ListAttachedPolicies(principalType, principalName)
	if err != nil {
		return nil, storeListError(err)
	}
	for _, arn := range arns {
		version, err := store.Policies().GetDefaultVersion(arn)
		if err != nil {
			return nil, storeListError(err)
		}
		if version == nil {
			return nil, storeListError(errors.New("attached policy " + arn + " has no default version"))
		}
		doc, err := policy.ParseDocument(version.Document)
		if err != nil {
			logs.Warn("skipping unparseable managed policy document in policy collection",
				logs.String("principalType", principalType),
				logs.String("principal", principalName),
				logs.String("policy", arn),
				logs.Err(err))
			continue
		}
		p, err := store.Policies().Get(arn)
		if err != nil {
			return nil, storeListError(err)
		}
		into = append(into, principalPolicy{
			Document:   doc,
			PolicyName: p.PolicyName,
			PolicyArn:  arn,
		})
	}
	return into, nil
}

// SimulationContextEntry is one condition context key supplied to a policy
// simulation, with every value of the (possibly multi-valued) key.
type SimulationContextEntry struct {
	ContextKeyName   string
	ContextKeyValues []string
}

// SimulatePrincipalPolicyInput holds the parameters of a principal-policy
// simulation shared by both protocol planes.
type SimulatePrincipalPolicyInput struct {
	PolicySourceArn                    string
	ActionNames                        []string
	ResourceArns                       []string
	PolicyInputList                    []string
	PermissionsBoundaryPolicyInputList []string
	ContextEntries                     []SimulationContextEntry
	PolicyExclusionList                []SimulationPolicyIdentifier
	ResourcePolicy                     string
	ResourceOwner                      string
	CallerArn                          string
	ResourceHandlingOption             string
}

// SimulateCustomPolicyInput holds the parameters of a custom-policy
// simulation shared by both protocol planes.
type SimulateCustomPolicyInput struct {
	PolicyInputList                    []string
	PermissionsBoundaryPolicyInputList []string
	OrderedOrganizationPolicyInputList []string
	ActionNames                        []string
	ResourceArns                       []string
	ResourcePolicy                     string
	ResourceOwner                      string
	CallerArn                          string
	ContextEntries                     []SimulationContextEntry
	ResourceHandlingOption             string
}

// SimulationPolicyIdentifier is one PolicyExclusionList entry in the
// PolicyIdentifier union: exactly one of PolicyType, PolicyArn or the
// inline-policy discriminator (InlinePolicyName with optional attachment
// type and name) is set.
type SimulationPolicyIdentifier struct {
	PolicyType           string
	PolicyArn            string
	InlinePolicyName     string
	InlineAttachmentType string
	InlineAttachmentName string
}

// simulationSourcePolicy identifies the policy one simulation document was
// taken from: the wire Statement members SourcePolicyId and SourcePolicyType.
// PolicyType is the PolicySourceType wire value (user, group, role,
// aws-managed, user-managed, resource); an empty PolicyType omits the
// member, which is how policies passed in PolicyInputList are reported.
type simulationSourcePolicy struct {
	PolicyId   string
	PolicyType string
}

// simulationPolicySet pairs every document included in a simulation with
// the policy it came from at the same index, so a matched statement
// resolves to its source policy.
type simulationPolicySet struct {
	Documents []*policy.Document
	Sources   []simulationSourcePolicy
}

func (set *simulationPolicySet) add(doc *policy.Document, source simulationSourcePolicy) {
	set.Documents = append(set.Documents, doc)
	set.Sources = append(set.Sources, source)
}

// simulationSourcePolicyFor renders the wire source of one gathered policy:
// an inline policy is identified by its attachment (entity and policy name)
// and typed by the entity kind it is attached to; an attached managed
// policy is identified by its ARN and typed by its managed-ness.
func simulationSourcePolicyFor(rec principalPolicy) simulationSourcePolicy {
	if rec.PolicyArn == "" {
		return simulationSourcePolicy{
			PolicyId:   rec.EntityName + "_" + rec.PolicyName,
			PolicyType: rec.EntityType,
		}
	}
	policyType := "user-managed"
	if iamstore.IsAWSManagedPolicyARN(rec.PolicyArn) {
		policyType = "aws-managed"
	}
	return simulationSourcePolicy{PolicyId: rec.PolicyArn, PolicyType: policyType}
}

// resolveMatchedStatements maps the evaluator's matched statements to their
// source policies. Statements of the same policy stay separate entries —
// each is one matched statement of the document.
func resolveMatchedStatements(matched []policy.MatchedStatement, sources []simulationSourcePolicy) []simulationSourcePolicy {
	resolved := make([]simulationSourcePolicy, 0, len(matched))
	for _, m := range matched {
		if m.PolicyIndex < 0 || m.PolicyIndex >= len(sources) {
			continue
		}
		resolved = append(resolved, sources[m.PolicyIndex])
	}
	return resolved
}

// simulationResourceResult is the per-resource evaluation of one action.
type simulationResourceResult struct {
	EvalResourceName     string
	EvalResourceDecision string // allowed, explicitDeny, or implicitDeny
	MatchedStatements    []simulationSourcePolicy
	MissingContextValues []string
	AllowedByBoundary    bool
	// EvalDecisionDetails carries the per-policy-type decisions of a
	// cross-account simulation; nil unless the organisation policies
	// allowed the request.
	EvalDecisionDetails []simulationDecisionDetail
}

// simulationDecisionDetail is one EvalDecisionDetails entry: the decision
// one policy type reached (IAM Policy, Resource Policy, Permissions
// Boundary Policy).
type simulationDecisionDetail struct {
	PolicyType string
	Decision   string
}

// SimulationEvaluation is one action's evaluation: the most restrictive
// decision across every simulated resource, the union of the statements
// that determined that decision, and each resource's own decision under
// ResourceSpecificResults.
type SimulationEvaluation struct {
	EvalActionName       string
	EvalResourceName     string // the ARN template; "*" when none is defined
	EvalDecision         string // allowed, explicitDeny, or implicitDeny
	MatchedStatements    []simulationSourcePolicy
	MissingContextValues []string
	AllowedByBoundary    bool
	HasBoundary          bool
	// HasOrgPolicies reports that an ordered organisation policy list
	// participated; AllowedByOrganizations then carries whether those
	// policies allowed the operation, most restrictive across resources.
	HasOrgPolicies          bool
	AllowedByOrganizations  bool
	ResourceSpecificResults []simulationResourceResult
	// EvalDecisionDetails carries the per-policy-type decisions of a
	// cross-account simulation, aggregated to the most restrictive
	// decision per policy type across resources. It is present but empty
	// when a same-account simulation names a resource ARN, and nil
	// otherwise: same-account all-resources simulations, and cross-account
	// simulations an SCP denied.
	EvalDecisionDetails []simulationDecisionDetail
}

// SimulatePrincipalPolicyResult holds the evaluations produced by
// simulatePrincipalPolicyCore.
type SimulatePrincipalPolicyResult struct {
	Evaluations []SimulationEvaluation
}

// simulationDecisionRank orders the wire decisions by restrictiveness for
// the most-restrictive aggregate: an explicit deny overrides everything,
// and an implicit deny overrides an allow.
func simulationDecisionRank(decision string) int {
	switch decision {
	case "explicitDeny":
		return 0
	case "implicitDeny":
		return 1
	default:
		return 2
	}
}

// simulatePrincipalPolicyCore gathers every identity-based policy that
// applies to the principal (plus the caller-supplied PolicyInputList
// documents) and runs the shared simulation engine over them.
func (s *IAMService) simulatePrincipalPolicyCore(store *iamstore.IAMStore, input *SimulatePrincipalPolicyInput) (*SimulatePrincipalPolicyResult, error) {
	if input.PolicySourceArn == "" {
		return nil, NewValidationError("PolicySourceArn")
	}
	if len(input.ActionNames) == 0 {
		return nil, NewValidationError("ActionNames")
	}

	records, err := s.gatherPrincipalPolicyRecordsCore(store, input.PolicySourceArn, "PolicySourceArn")
	if err != nil {
		return nil, err
	}

	return s.runSimulationCore(&simulationRequest{
		SourceArn:         input.PolicySourceArn,
		SourceEntityType:  resolveEntityType(input.PolicySourceArn),
		CallerArn:         input.CallerArn,
		ActionNames:       input.ActionNames,
		Resources:         input.ResourceArns,
		Records:           records,
		PolicyInputList:   input.PolicyInputList,
		BoundaryInputList: input.PermissionsBoundaryPolicyInputList,
		ContextEntries:    input.ContextEntries,
		Exclusions:        input.PolicyExclusionList,
		ResourcePolicy:    input.ResourcePolicy,
		ResourceOwner:     input.ResourceOwner,
		HandlingOption:    input.ResourceHandlingOption,
	})
}

// simulateCustomPolicyCore runs the shared simulation engine over the
// caller-supplied policy documents only — no principal's policies are
// gathered — with the ordered organisation policy list as an additional
// upper bound.
func (s *IAMService) simulateCustomPolicyCore(_ *iamstore.IAMStore, input *SimulateCustomPolicyInput) (*SimulatePrincipalPolicyResult, error) {
	if len(input.PolicyInputList) == 0 {
		return nil, NewValidationError("PolicyInputList")
	}
	if len(input.ActionNames) == 0 {
		return nil, NewValidationError("ActionNames")
	}

	return s.runSimulationCore(&simulationRequest{
		CallerArn:          input.CallerArn,
		ActionNames:        input.ActionNames,
		Resources:          input.ResourceArns,
		PolicyInputList:    input.PolicyInputList,
		BoundaryInputList:  input.PermissionsBoundaryPolicyInputList,
		OrgPolicyInputList: input.OrderedOrganizationPolicyInputList,
		ContextEntries:     input.ContextEntries,
		ResourcePolicy:     input.ResourcePolicy,
		ResourceOwner:      input.ResourceOwner,
		HandlingOption:     input.ResourceHandlingOption,
	})
}

// simulationRequest is the evaluated simulation request shared by both
// simulate operations after operation-specific validation: the identity
// policies gathered from the principal (empty for the custom-policy
// operation), the caller-supplied documents, and the evaluation modifiers.
type simulationRequest struct {
	SourceArn          string // PolicySourceArn; empty on the custom path
	SourceEntityType   string // User, Group or Role; empty on the custom path
	CallerArn          string
	ActionNames        []string
	Resources          []string
	Records            []principalPolicy
	PolicyInputList    []string
	BoundaryInputList  []string
	OrgPolicyInputList []string
	ContextEntries     []SimulationContextEntry
	Exclusions         []SimulationPolicyIdentifier
	ResourcePolicy     string
	ResourceOwner      string
	HandlingOption     string
}

// runSimulationCore evaluates each action once across the whole resource
// set. The result follows the EvaluationResult contract: one evaluation
// per action, the top-level decision the most restrictive across
// resources, and per-resource decisions under ResourceSpecificResults.
// Identity-policy allows are capped by the permissions boundary and, on
// the custom path, by the ordered organisation policy list; a resource
// policy participates when supplied — same-account allows union, cross-
// account allows intersect — with the per-policy-type decisions reported
// under EvalDecisionDetails for cross-account simulations, aggregated to
// the most restrictive decision per policy type across resources; an SCP
// that denies ends the evaluation and suppresses the details.
func (s *IAMService) runSimulationCore(req *simulationRequest) (*SimulatePrincipalPolicyResult, error) {
	resources := req.Resources
	if len(resources) == 0 {
		resources = []string{"*"}
	}
	if err := validateSimulationHandlingOption(req.HandlingOption, resources); err != nil {
		return nil, err
	}
	if err := validateSimulationPolicyIdentifiers(req.Exclusions); err != nil {
		return nil, err
	}

	// The simulated caller defaults to the policy source; the custom
	// operation has no source and no default.
	caller := req.CallerArn
	if caller == "" {
		caller = req.SourceArn
	}
	if caller != "" {
		switch resolveEntityType(caller) {
		case "User", "Group", "Role":
		default:
			return nil, NewInvalidInputError("CallerArn", "must be a user, group, or role ARN")
		}
	}
	// Resource policies are not simulated for role sources.
	if req.ResourcePolicy != "" && req.SourceEntityType == "Role" {
		return nil, NewInvalidInputError("ResourcePolicy", "simulation of resource-based policies isn't supported for IAM roles")
	}

	policySet := simulationPolicySet{}
	for _, rec := range req.Records {
		if simulationExcluded(req.Exclusions, rec) {
			continue
		}
		policySet.add(rec.Document, simulationSourcePolicyFor(rec))
	}
	for i, pDoc := range req.PolicyInputList {
		doc, pErr := policy.ParseDocument(pDoc)
		if pErr != nil {
			return nil, NewInvalidInputError("PolicyInputList", "contains a malformed policy document")
		}
		// An input-list policy has no source type; its identifier is the
		// positional one the response reports it by.
		policySet.add(doc, simulationSourcePolicy{PolicyId: fmt.Sprintf("PolicyInputList.%d", i+1)})
	}

	var boundaryDocs []*policy.Document
	for _, bDoc := range req.BoundaryInputList {
		doc, bErr := policy.ParseDocument(bDoc)
		if bErr != nil {
			return nil, NewInvalidInputError("PermissionsBoundaryPolicyInputList", "contains a malformed policy document")
		}
		boundaryDocs = append(boundaryDocs, doc)
	}
	// A permission-boundary exclusion drops the caller-supplied boundary
	// from the simulation.
	if simulationExcludesBoundary(req.Exclusions) {
		boundaryDocs = nil
	}

	var orgDocs []*policy.Document
	for _, oDoc := range req.OrgPolicyInputList {
		doc, oErr := policy.ParseDocument(oDoc)
		if oErr != nil {
			return nil, NewInvalidInputError("OrderedOrganizationPolicyInputList", "contains a malformed policy document")
		}
		orgDocs = append(orgDocs, doc)
	}

	var resourcePolicyDocs []*policy.Document
	if req.ResourcePolicy != "" {
		doc, rErr := policy.ParseDocument(req.ResourcePolicy)
		if rErr != nil {
			return nil, NewInvalidInputError("ResourcePolicy", "contains a malformed policy document")
		}
		resourcePolicyDocs = []*policy.Document{doc}
	}

	// Condition key lookups go through EvaluationContext.GetContextValue
	// and ContextValues, which lowercase the key before consulting
	// SessionContext, so the map must be keyed lowercase for the supplied
	// values to be seen.
	sessionCtx := make(map[string][]string, len(req.ContextEntries))
	for _, entry := range req.ContextEntries {
		if entry.ContextKeyName == "" || len(entry.ContextKeyValues) == 0 {
			continue
		}
		// Each supplied value stays one element of the key's value list:
		// set operators consume the list structurally, so a value that
		// contains a comma is never mistaken for two values.
		sessionCtx[strings.ToLower(entry.ContextKeyName)] = entry.ContextKeyValues
	}

	// Missing context values cover identity-based and resource-based
	// policies only: context keys referenced solely by a permissions
	// boundary or an SCP are not reported, the same taxonomy that keeps
	// boundary and SCP statements out of MatchedStatements.
	missingContext := missingContextValues(append(policySet.Documents, resourcePolicyDocs...), sessionCtx)
	perResourceMissing := !allResourcesWildcard(resources)

	// The resource policy's account decides same-account union against
	// cross-account intersection semantics; an unspecified owner defaults
	// to the caller's own account.
	callerAccount := extractAccountFromARN(caller)
	resourceOwnerAccount := callerAccount
	if req.ResourceOwner != "" {
		resourceOwnerAccount = simulationOwnerAccount(req.ResourceOwner)
		if resourceOwnerAccount == "" {
			return nil, NewInvalidInputError("ResourceOwner", "must be an account ID or account root ARN")
		}
	}
	crossAccount := len(resourcePolicyDocs) > 0 && resourceOwnerAccount != callerAccount

	principalName := extractPrincipalNameFromARN(caller)

	evaluator := policy.NewPolicyEvaluator()

	result := &SimulatePrincipalPolicyResult{
		Evaluations: make([]SimulationEvaluation, 0, len(req.ActionNames)),
	}
	for _, action := range req.ActionNames {
		// The first resource sets the aggregate; a later resource can
		// only tighten it, and equally-deciding resources union their
		// determining statements into the top-level list. The initial
		// rank sits beyond every real decision rank.
		aggRank := 3
		evaluation := SimulationEvaluation{
			EvalActionName:    action,
			EvalResourceName:  "*", // no per-action ARN-template registry exists
			AllowedByBoundary: true,
			HasBoundary:       len(boundaryDocs) > 0,
			HasOrgPolicies:    len(orgDocs) > 0,
			// An organisation policy list allows an operation only when it
			// allows it for every simulated resource.
			AllowedByOrganizations: true,
		}
		if !perResourceMissing {
			evaluation.MissingContextValues = missingContext
		}
		if !crossAccount && !allResourcesWildcard(resources) {
			// A same-account simulation that names a resource ARN reports
			// the member present but empty; per-resource entries omit it
			// entirely. A same-account simulation over only the
			// all-resources wildcard does not return the member at all.
			evaluation.EvalDecisionDetails = []simulationDecisionDetail{}
		}
		// The top-level details aggregate, per policy type, the most
		// restrictive decision across resources — the deciding resource of
		// one type can differ from the resource that set the overall
		// decision.
		detailAggregate := make(map[string]string, len(simulationDetailPolicyTypes))

		for _, resource := range resources {
			evalCtx := &policy.EvaluationContext{
				Principal:        caller,
				PrincipalAccount: callerAccount,
				Action:           action,
				Resource:         resource,
				RequestTime:      time.Now(),
				UserName:         principalName,
				SessionContext:   sessionCtx,
			}

			// Evaluate identity-based policies.
			identityDecision := evaluator.Evaluate(evalCtx, policySet.Documents)
			identityMatched := resolveMatchedStatements(identityDecision.Matched, policySet.Sources)
			identityWire := simulationWireDecision(identityDecision.Effect)

			// The boundary is evaluated for its own decision detail even
			// when the identity policies deny, but it only caps the
			// decision when they allow.
			allowedByBoundary := true
			boundaryWire := "implicitDeny"
			if len(boundaryDocs) > 0 {
				boundaryDecision := evaluator.Evaluate(evalCtx, boundaryDocs)
				allowedByBoundary = boundaryDecision.Effect == policy.DecisionEffectAllow
				boundaryWire = simulationWireDecision(boundaryDecision.Effect)
			}
			identityCapped := identityWire == "allowed" && len(boundaryDocs) > 0 && !allowedByBoundary

			// The organisation policy list is evaluated in order and acts
			// as an upper bound like the boundary: every level must allow.
			orgAllowed := true
			for _, orgDoc := range orgDocs {
				if evaluator.Evaluate(evalCtx, []*policy.Document{orgDoc}).Effect != policy.DecisionEffectAllow {
					orgAllowed = false
					break
				}
			}
			evaluation.AllowedByOrganizations = evaluation.AllowedByOrganizations && orgAllowed

			// The resource policy participates per resource, evaluated
			// for the simulated caller.
			resourcePolicyWire := ""
			var resourcePolicyMatched []simulationSourcePolicy
			if len(resourcePolicyDocs) > 0 {
				rpDecision := evaluator.Evaluate(evalCtx, resourcePolicyDocs)
				resourcePolicyMatched = resolveMatchedStatements(rpDecision.Matched, []simulationSourcePolicy{
					{PolicyId: "ResourcePolicy", PolicyType: "resource"},
				})
				resourcePolicyWire = simulationWireDecision(rpDecision.Effect)
			}

			// Combine: an explicit deny in the identity or resource policy
			// overrides everything; otherwise same-account allows union
			// and cross-account allows intersect. Identity-sourced allows
			// are capped by the boundary and the organisation policies;
			// resource-policy grants are not bounded (a permissions
			// boundary limits identity-based permissions only).
			identityAllows := identityWire == "allowed" && !identityCapped && orgAllowed
			var decisionWire string
			var contributing, matched []simulationSourcePolicy
			switch {
			case identityWire == "explicitDeny" || resourcePolicyWire == "explicitDeny":
				decisionWire = "explicitDeny"
				// Only the denying statements determine an explicit deny.
				if identityWire == "explicitDeny" {
					contributing = append(contributing, identityMatched...)
				}
				if resourcePolicyWire == "explicitDeny" {
					contributing = append(contributing, resourcePolicyMatched...)
				}
			case len(resourcePolicyDocs) == 0:
				if identityAllows {
					decisionWire = "allowed"
					contributing = identityMatched
				} else {
					decisionWire = "implicitDeny"
				}
			case crossAccount:
				if identityAllows && resourcePolicyWire == "allowed" {
					decisionWire = "allowed"
					contributing = append(append(contributing, identityMatched...), resourcePolicyMatched...)
				} else {
					decisionWire = "implicitDeny"
				}
			default: // same-account resource policy
				if identityAllows || resourcePolicyWire == "allowed" {
					decisionWire = "allowed"
					if identityAllows {
						contributing = append(contributing, identityMatched...)
					}
					if resourcePolicyWire == "allowed" {
						contributing = append(contributing, resourcePolicyMatched...)
					}
				} else {
					decisionWire = "implicitDeny"
				}
			}
			// Per-resource MatchedStatements carries every matched
			// statement of the identity and resource-policy evaluations,
			// including ones that did not decide the outcome.
			matched = append(append([]simulationSourcePolicy{}, identityMatched...), resourcePolicyMatched...)

			resourceResult := simulationResourceResult{
				EvalResourceName:     resource,
				EvalResourceDecision: decisionWire,
				MatchedStatements:    matched,
				AllowedByBoundary:    allowedByBoundary,
			}
			if perResourceMissing {
				resourceResult.MissingContextValues = missingContext
			}
			// An SCP that denies ends the evaluation for this resource:
			// its decision stands, but no policy-type details are reported
			// for it and it contributes nothing to the aggregate.
			if crossAccount && orgAllowed {
				resourceResult.EvalDecisionDetails = simulationDecisionDetails(identityWire, resourcePolicyWire, boundaryWire, len(boundaryDocs) > 0)
				mergeMostRestrictiveDetail(detailAggregate, resourceResult.EvalDecisionDetails)
			}
			evaluation.ResourceSpecificResults = append(evaluation.ResourceSpecificResults, resourceResult)

			rank := simulationDecisionRank(decisionWire)
			if rank < aggRank {
				aggRank = rank
				evaluation.EvalDecision = decisionWire
				evaluation.MatchedStatements = contributing
				evaluation.AllowedByBoundary = allowedByBoundary
			} else if rank == aggRank {
				evaluation.MatchedStatements = append(evaluation.MatchedStatements, contributing...)
				evaluation.AllowedByBoundary = evaluation.AllowedByBoundary && allowedByBoundary
			}
		}

		if crossAccount && len(detailAggregate) > 0 {
			// An empty aggregate means every resource ended on an SCP
			// deny, and the member is not returned.
			evaluation.EvalDecisionDetails = simulationDetailsFromAggregate(detailAggregate)
		}

		result.Evaluations = append(result.Evaluations, evaluation)
	}

	return result, nil
}

// simulationDetailPolicyTypes are the policy types of the cross-account
// decision-detail map, in emission order.
var simulationDetailPolicyTypes = [...]string{"IAM Policy", "Resource Policy", "Permissions Boundary Policy"}

// simulationDecisionDetails renders the per-policy-type decisions of a
// cross-account simulation.
func simulationDecisionDetails(identity, resourcePolicy, boundary string, hasBoundary bool) []simulationDecisionDetail {
	details := []simulationDecisionDetail{
		{PolicyType: simulationDetailPolicyTypes[0], Decision: identity},
		{PolicyType: simulationDetailPolicyTypes[1], Decision: resourcePolicy},
	}
	if hasBoundary {
		details = append(details, simulationDecisionDetail{PolicyType: simulationDetailPolicyTypes[2], Decision: boundary})
	}
	return details
}

// mergeMostRestrictiveDetail folds one resource's details into the
// aggregate, keeping per policy type the most restrictive decision seen.
func mergeMostRestrictiveDetail(aggregate map[string]string, details []simulationDecisionDetail) {
	for _, detail := range details {
		if current, ok := aggregate[detail.PolicyType]; !ok || simulationDecisionRank(detail.Decision) < simulationDecisionRank(current) {
			aggregate[detail.PolicyType] = detail.Decision
		}
	}
}

// simulationDetailsFromAggregate renders the aggregated detail map in the
// emission order of simulationDecisionDetails.
func simulationDetailsFromAggregate(aggregate map[string]string) []simulationDecisionDetail {
	details := make([]simulationDecisionDetail, 0, len(aggregate))
	for _, policyType := range simulationDetailPolicyTypes {
		if decision, ok := aggregate[policyType]; ok {
			details = append(details, simulationDecisionDetail{PolicyType: policyType, Decision: decision})
		}
	}
	return details
}

// simulationOwnerAccount extracts the account identifier a ResourceOwner
// carries — a bare account ID or an account root ARN — or "" when the
// value is neither.
func simulationOwnerAccount(owner string) string {
	if strings.Contains(owner, ":") {
		if account := extractAccountFromARN(owner); account != "" {
			return account
		}
		return ""
	}
	return owner
}

// simulationHandlingOptions are the documented ResourceHandlingOption
// scenario values with the EC2 resource types each requires present in
// ResourceArns.
var simulationHandlingOptions = map[string][]string{
	"EC2-VPC-InstanceStore":        {"instance", "image", "security-group", "network-interface"},
	"EC2-VPC-InstanceStore-Subnet": {"instance", "image", "security-group", "network-interface", "subnet"},
	"EC2-VPC-EBS":                  {"instance", "image", "security-group", "network-interface", "volume"},
	"EC2-VPC-EBS-Subnet":           {"instance", "image", "security-group", "network-interface", "subnet", "volume"},
}

// validateSimulationHandlingOption rejects an unknown scenario value and
// enforces the resource types the documented scenarios require.
func validateSimulationHandlingOption(option string, resources []string) error {
	if option == "" {
		return nil
	}
	required, ok := simulationHandlingOptions[option]
	if !ok {
		return NewInvalidInputError("ResourceHandlingOption", "must be one of the documented simulation scenarios")
	}
	present := make(map[string]bool, len(resources))
	for _, resource := range resources {
		present[simulationEC2ResourceType(resource)] = true
	}
	for _, resourceType := range required {
		if !present[resourceType] {
			return NewInvalidInputError("ResourceArns", "the "+option+" scenario requires a "+resourceType+" resource")
		}
	}
	return nil
}

// simulationEC2ResourceType returns the resource-type segment of an EC2
// ARN (arn:...:ec2:region:account:instance/i-...), or "" for any other
// shape.
func simulationEC2ResourceType(resource string) string {
	_, service, _, _, res := awsarn.SplitARN(resource)
	if service != "ec2" {
		return ""
	}
	seg := strings.SplitN(res, "/", 2)
	return seg[0]
}

// simulationPolicyIdentifierTypes are the documented PolicyIdentifier
// PolicyType values.
var simulationPolicyIdentifierTypes = map[string]bool{
	"inline":              true,
	"aws-managed":         true,
	"user-managed":        true,
	"permission-boundary": true,
	"scp":                 true,
	"rcp":                 true,
}

// validateSimulationPolicyIdentifiers enforces the PolicyIdentifier union:
// exactly one discriminator per entry, a documented policy type, and a
// syntactically usable ARN or inline-policy shape. Syntactically valid
// identifiers that match nothing are ignored at evaluation.
func validateSimulationPolicyIdentifiers(list []SimulationPolicyIdentifier) error {
	for _, id := range list {
		discriminators := 0
		if id.PolicyType != "" {
			discriminators++
			if !simulationPolicyIdentifierTypes[id.PolicyType] {
				return NewInvalidInputError("PolicyExclusionList", "unknown policy type "+id.PolicyType)
			}
		}
		if id.PolicyArn != "" {
			discriminators++
		}
		if id.InlinePolicyName != "" {
			discriminators++
		}
		switch {
		case discriminators == 0:
			return NewInvalidInputError("PolicyExclusionList", "each identifier must set PolicyType, PolicyArn, or InlinePolicyIdentifier")
		case discriminators > 1:
			return NewInvalidInputError("PolicyExclusionList", "each identifier must set exactly one of PolicyType, PolicyArn, or InlinePolicyIdentifier")
		}
		if id.PolicyArn != "" {
			if _, err := awsarn.ParseARN(id.PolicyArn); err != nil {
				return NewInvalidInputError("PolicyExclusionList", "contains a malformed policy ARN")
			}
		}
	}
	return nil
}

// simulationExcluded reports whether one gathered policy matches any
// exclusion identifier. Types with no gathered counterpart (scp, rcp) and
// identifiers that match nothing leave the policy in the simulation.
func simulationExcluded(list []SimulationPolicyIdentifier, rec principalPolicy) bool {
	for _, id := range list {
		switch {
		case id.PolicyType == "inline":
			if rec.PolicyArn == "" {
				return true
			}
		case id.PolicyType == "aws-managed":
			if rec.PolicyArn != "" && iamstore.IsAWSManagedPolicyARN(rec.PolicyArn) {
				return true
			}
		case id.PolicyType == "user-managed":
			if rec.PolicyArn != "" && !iamstore.IsAWSManagedPolicyARN(rec.PolicyArn) {
				return true
			}
		case id.PolicyArn != "":
			if rec.PolicyArn != "" && simulationWildcardMatch(id.PolicyArn, rec.PolicyArn) {
				return true
			}
		case id.InlinePolicyName != "":
			if rec.PolicyArn != "" {
				continue
			}
			if id.InlinePolicyName != rec.PolicyName {
				continue
			}
			if id.InlineAttachmentType != "" && id.InlineAttachmentType != rec.EntityType {
				continue
			}
			if id.InlineAttachmentName != "" && !simulationWildcardMatch(id.InlineAttachmentName, rec.EntityName) {
				continue
			}
			return true
		}
	}
	return false
}

// simulationExcludesBoundary reports whether the exclusion list carries
// the permission-boundary type identifier, which drops the caller-supplied
// boundary documents from the simulation.
func simulationExcludesBoundary(list []SimulationPolicyIdentifier) bool {
	for _, id := range list {
		if id.PolicyType == "permission-boundary" {
			return true
		}
	}
	return false
}

// simulationWildcardMatch reports whether name matches pattern, where the
// pattern may carry * (any sequence) and ? (single character) — the
// wildcard grammar PolicyIdentifier documents for ARNs and entity names.
func simulationWildcardMatch(pattern, name string) bool {
	// Segment the pattern on * and anchor the literal runs.
	i, j := 0, 0
	star, mark := -1, 0
	for i < len(name) {
		if j < len(pattern) && (pattern[j] == '?' || pattern[j] == name[i]) {
			i++
			j++
		} else if j < len(pattern) && pattern[j] == '*' {
			star = j
			j++
			mark = i
		} else if star >= 0 {
			j = star + 1
			mark++
			i = mark
		} else {
			return false
		}
	}
	for j < len(pattern) && pattern[j] == '*' {
		j++
	}
	return j == len(pattern)
}

// simulationWireDecision maps an evaluation effect to its wire decision.
func simulationWireDecision(effect policy.DecisionEffect) string {
	switch effect {
	case policy.DecisionEffectAllow:
		return "allowed"
	case policy.DecisionEffectDeny:
		return "explicitDeny"
	default:
		return "implicitDeny"
	}
}

// allResourcesWildcard reports whether every simulated resource is the
// wildcard, which is where MissingContextValues belongs at the top level.
func allResourcesWildcard(resources []string) bool {
	for _, resource := range resources {
		if resource != "*" {
			return false
		}
	}
	return true
}

// simulationProvidedContextKeys are the condition keys the simulator itself
// materialises for every evaluation, so a policy requiring one of them is
// never reported as missing context.
var simulationProvidedContextKeys = map[string]bool{
	"aws:principalarn":     true,
	"aws:principalaccount": true,
	"aws:action":           true,
	"aws:resource":         true,
	"aws:currenttime":      true,
	"aws:username":         true,
}

// missingContextValues returns the condition keys the included policy
// documents require but neither the ContextEntries parameter nor the
// simulator's own materialised keys provide, in the order the documents
// first reference them.
func missingContextValues(documents []*policy.Document, provided map[string][]string) []string {
	var missing []string
	seen := make(map[string]bool)
	isProvided := func(key string) bool {
		key = strings.ToLower(key)
		if simulationProvidedContextKeys[key] {
			return true
		}
		_, ok := provided[key]
		return ok
	}
	for _, doc := range documents {
		if doc == nil {
			continue
		}
		for _, stmt := range doc.Statement {
			for _, keyValues := range stmt.Condition {
				for key := range keyValues {
					lower := strings.ToLower(key)
					if seen[lower] || isProvided(key) {
						continue
					}
					seen[lower] = true
					missing = append(missing, key)
				}
			}
		}
	}
	return missing
}

// extractPrincipalNameFromARN returns the entity name segment of an IAM
// principal ARN.
func extractPrincipalNameFromARN(arn string) string {
	_, name := parseIAMARNResource(arn)
	return name
}

// extractAccountFromARN returns the account segment of an IAM principal ARN.
func extractAccountFromARN(arn string) string {
	_, _, _, accountID, _ := awsarn.SplitARN(arn)
	return accountID
}
