// Transport-agnostic Core function for account-wide IAM entity details:
// validation and store operations shared by the AWS-compatible HTTP API
// handlers and the admin gRPC-Web handler (the xxxCore pattern).
package iam

import (
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/timeutils"
)

// AccountAuthorizationDetailsInput holds the parameters for
// GetAccountAuthorizationDetails.
type AccountAuthorizationDetailsInput struct {
	Filters  map[string]bool
	Marker   string
	MaxItems int
}

// getAccountAuthorizationDetailsCore aggregates all IAM entities (users,
// groups, roles, local managed policies) with their inline and attached
// policy relationships.  The logic is consolidated in a core function so
// that both the HTTP API and future admin-handler paths delegate here.
func (s *IAMService) getAccountAuthorizationDetailsCore(reqCtx *request.RequestContext, store *iamstore.IAMStore, input *AccountAuthorizationDetailsInput) (interface{}, error) {
	filters := input.Filters
	if len(filters) == 0 {
		filters = map[string]bool{
			"User":               true,
			"Group":              true,
			"Role":               true,
			"LocalManagedPolicy": true,
		}
	}

	marker := input.Marker
	maxItems := input.MaxItems

	type section struct {
		key    string
		filter string
		items  []interface{}
		marker func(i int) string
	}
	sections := make([]section, 0, 4)

	if filters["User"] {
		users, err := listAllUsers(store, "")
		if err != nil {
			return nil, err
		}
		items := make([]interface{}, 0, len(users))
		for _, user := range users {
			detail := map[string]interface{}{
				"UserId":     user.ID,
				"Path":       user.Path,
				"UserName":   user.UserName,
				"Arn":        user.Arn,
				"CreateDate": user.CreateDate.Format(timeutils.ISO8601SimpleFormat),
			}
			groupNames, _ := store.UserGroups().ListGroupsForUser(user.UserName)
			groupList := make([]interface{}, 0, len(groupNames))
			for _, gn := range groupNames {
				groupList = append(groupList, gn)
			}
			detail["GroupList"] = groupList
			detail["AttachedManagedPolicies"] = buildAttachedManagedPolicies(store, PrincipalTypeUser, user.UserName)
			detail["UserPolicyList"] = buildInlinePolicyList(store, PrincipalTypeUser, user.UserName)
			if user.PermissionsBoundary != nil {
				detail["PermissionsBoundary"] = map[string]interface{}{
					"PermissionsBoundaryType": user.PermissionsBoundary.PermissionsBoundaryType,
					"PermissionsBoundaryArn":  user.PermissionsBoundary.PermissionsBoundaryArn,
				}
			}
			if tagList := tags.ToResponse(user.Tags); tagList != nil {
				detail["Tags"] = tagList
			}
			items = append(items, detail)
		}
		sections = append(sections, section{
			key: "UserDetailList", filter: "User", items: items,
			marker: func(i int) string { return "user:" + users[i].UserName },
		})
	}

	if filters["Group"] {
		groups, err := listAllGroups(store, "")
		if err != nil {
			return nil, err
		}
		items := make([]interface{}, 0, len(groups))
		for _, group := range groups {
			detail := map[string]interface{}{
				"GroupId":    group.ID,
				"Path":       group.Path,
				"GroupName":  group.GroupName,
				"Arn":        group.Arn,
				"CreateDate": group.CreateDate.Format(timeutils.ISO8601SimpleFormat),
			}
			detail["GroupPolicyList"] = buildInlinePolicyList(store, PrincipalTypeGroup, group.GroupName)
			detail["AttachedManagedPolicies"] = buildAttachedManagedPolicies(store, PrincipalTypeGroup, group.GroupName)
			items = append(items, detail)
		}
		sections = append(sections, section{
			key: "GroupDetailList", filter: "Group", items: items,
			marker: func(i int) string { return "group:" + groups[i].GroupName },
		})
	}

	if filters["Role"] {
		roles, err := listAllRoles(store, "")
		if err != nil {
			return nil, err
		}
		items := make([]interface{}, 0, len(roles))
		for _, role := range roles {
			detail := map[string]interface{}{
				"RoleId":                   role.ID,
				"Path":                     role.Path,
				"RoleName":                 role.RoleName,
				"Arn":                      role.Arn,
				"CreateDate":               role.CreateDate.Format(timeutils.ISO8601SimpleFormat),
				"AssumeRolePolicyDocument": role.AssumeRolePolicyDocument,
			}
			detail["RolePolicyList"] = buildInlinePolicyList(store, PrincipalTypeRole, role.RoleName)
			detail["AttachedManagedPolicies"] = buildAttachedManagedPolicies(store, PrincipalTypeRole, role.RoleName)
			if role.PermissionsBoundary != nil {
				detail["PermissionsBoundary"] = map[string]interface{}{
					"PermissionsBoundaryType": role.PermissionsBoundary.PermissionsBoundaryType,
					"PermissionsBoundaryArn":  role.PermissionsBoundary.PermissionsBoundaryArn,
				}
			}
			if tagList := tags.ToResponse(role.Tags); tagList != nil {
				detail["Tags"] = tagList
			}
			if role.RoleLastUsed != nil {
				lastUsed := map[string]interface{}{}
				if role.RoleLastUsed.LastUsedDate != nil {
					lastUsed["LastUsedDate"] = role.RoleLastUsed.LastUsedDate.Format(timeutils.ISO8601SimpleFormat)
				}
				if role.RoleLastUsed.Region != "" {
					lastUsed["Region"] = role.RoleLastUsed.Region
				}
				detail["RoleLastUsed"] = lastUsed
			}
			profileList, err := store.InstanceProfiles().ListForRole(role.RoleName, "", 1000)
			if err != nil {
				return nil, err
			}
			instanceProfiles := make([]interface{}, 0, len(profileList.InstanceProfiles))
			for _, profile := range profileList.InstanceProfiles {
				instanceProfiles = append(instanceProfiles, s.instanceProfileToResponseWithRoles(profile, resolveInstanceProfileRoles(store, profile)))
			}
			detail["InstanceProfileList"] = instanceProfiles
			items = append(items, detail)
		}
		sections = append(sections, section{
			key: "RoleDetailList", filter: "Role", items: items,
			marker: func(i int) string { return "role:" + roles[i].RoleName },
		})
	}

	if filters["LocalManagedPolicy"] {
		policies, err := listAllPolicies(store, "Local", "", false)
		if err != nil {
			return nil, err
		}
		items := make([]interface{}, 0, len(policies))
		for _, policy := range policies {
			versions, err := store.Policies().ListVersions(policy.Arn, "", 1000)
			if err != nil {
				return nil, err
			}
			versionList := make([]interface{}, 0, len(versions.Versions))
			for _, version := range versions.Versions {
				versionList = append(versionList, s.policyVersionToResponse(version))
			}
			detail := map[string]interface{}{
				"PolicyName":                    policy.PolicyName,
				"PolicyId":                      policy.ID,
				"Arn":                           policy.Arn,
				"Path":                          policy.Path,
				"DefaultVersionId":              policy.DefaultVersionId,
				"AttachmentCount":               policy.AttachmentCount,
				"CreateDate":                    policy.CreateDate.Format(timeutils.ISO8601SimpleFormat),
				"UpdateDate":                    policy.UpdateDate.Format(timeutils.ISO8601SimpleFormat),
				"IsAttachable":                  policy.IsAttachable,
				"PermissionsBoundaryUsageCount": policy.PermissionsBoundaryUsageCount,
				"PolicyVersionList":             versionList,
			}
			if policy.Description != "" {
				detail["Description"] = policy.Description
			}
			items = append(items, detail)
		}
		sections = append(sections, section{
			key: "Policies", filter: "LocalManagedPolicy", items: items,
			marker: func(i int) string { return "policy:" + policies[i].Arn },
		})
	}

	// Apply pagination across all sections combined.
	// Marker format: "<sectionType>:<itemName>" (e.g. "user:admin", "role:MyRole").
	resp := map[string]interface{}{}
	skipUntilMarker := marker != ""
	count := 0
	isTruncated := false
	nextMarker := ""

	for _, sec := range sections {
		// Always emit the key so AWS SDK deserialisation gets an empty list
		// instead of null for skipped/partial sections.
		secItems := make([]interface{}, 0)

		for i, item := range sec.items {
			itemMarker := sec.marker(i)

			if skipUntilMarker {
				if itemMarker == marker {
					skipUntilMarker = false
				}
				continue
			}

			if count >= maxItems {
				isTruncated = true
				nextMarker = itemMarker
				break
			}
			secItems = append(secItems, item)
			count++
		}

		resp[sec.key] = secItems

		if isTruncated {
			break
		}
	}

	if skipUntilMarker {
		for _, sec := range sections {
			resp[sec.key] = []interface{}{}
		}
	}

	resp["IsTruncated"] = isTruncated
	if isTruncated && nextMarker != "" {
		resp["Marker"] = nextMarker
	}

	return resp, nil
}
