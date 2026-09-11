package iam

import (
	types "vorpalstacks/internal/common/tags"
)

// RoleStoreInterface defines operations for managing IAM roles.
type RoleStoreInterface interface {
	Get(roleName string) (*Role, error)
	GetByID(roleID string) (*Role, error)
	List(pathPrefix string, marker string, maxItems int) (*RoleListResult, error)
	Put(role *Role) error
	Delete(roleName string) error
	Exists(roleName string) bool
	Create(roleName, path, accountId, assumeRolePolicyDocument, description string, maxSessionDuration int, tags []types.Tag, sourceRoleTemplate *SourceRoleTemplate) (*Role, error)
	Update(role *Role) error
	UpdateRoleLastUsed(roleName, region string) error
	Count() int
	GetAssumeRolePolicyDocument(roleName string) (string, error)
}

// AccessKeyStoreInterface defines operations for managing IAM access keys.
type AccessKeyStoreInterface interface {
	Get(accessKeyId string) (*AccessKey, error)
	GetBySecretKey(secretAccessKey string) (*AccessKey, error)
	ListByUserName(userName string) ([]*AccessKey, error)
	ListByUserNameWithSecret(userName string) ([]*AccessKey, error)
	Put(key *AccessKey) error
	Delete(accessKeyId string) error
	UpdateStatus(accessKeyId string, status AccessKeyStatus) error
	UpdateLastUsed(accessKeyId, region, service string) error
	Create(userName string) (*AccessKey, error)
	CreateWithLimit(userName string, maxKeys int) (*AccessKey, error)
	Exists(accessKeyId string) bool
	DeleteByUserName(userName string) error
	Count() int
	CountByUserName(userName string) (int, error)
}

// LoginProfileStoreInterface defines operations for managing IAM user login profiles.
type LoginProfileStoreInterface interface {
	Get(userName string) (*LoginProfile, error)
	Put(profile *LoginProfile) error
	Delete(userName string) error
	Exists(userName string) bool
	Create(userName, password string, passwordResetRequired bool) (*LoginProfile, error)
	UpdatePassword(userName, password string) error
	UpdatePasswordResetRequired(userName string, required bool) error
	VerifyPassword(userName, password string) (bool, error)
	CheckPasswordReuse(userName, password string, maxHistory int) (bool, error)
	Count() int
}

// UserStoreInterface defines operations for managing IAM users.
type UserStoreInterface interface {
	Get(userName string) (*User, error)
	GetByID(userID string) (*User, error)
	List(pathPrefix string, marker string, maxItems int) (*UserListResult, error)
	Put(user *User) error
	Delete(userName string) error
	Exists(userName string) bool
	Create(userName, path, accountId string, tags []types.Tag) (*User, error)
	UpdatePasswordLastUsed(userName string) error
	Count() int
}

// UserGroupStoreInterface defines operations for managing IAM user-group memberships.
type UserGroupStoreInterface interface {
	AddUserToGroup(userName, groupName string) error
	AddUserToGroupWithLimit(userName, groupName string, maxGroups int) error
	RemoveUserFromGroup(userName, groupName string) error
	MigrateUserGroup(userName, fromGroup, toGroup string) error
	IsUserInGroup(userName, groupName string) bool
	ListGroupsForUser(userName string) ([]string, error)
	ListUsersInGroup(groupName string) ([]string, error)
	RemoveAllGroupsForUser(userName string) error
	RemoveAllUsersFromGroup(groupName string) error
	CountUsersInGroup(groupName string) (int, error)
	MigrateUser(oldUserName, newUserName string) error
}

// AttachedPolicyStoreInterface defines operations for managing IAM attached policies.
type AttachedPolicyStoreInterface interface {
	Attach(principalType, principalName, policyArn string) error
	Detach(principalType, principalName, policyArn string) error
	IsAttached(principalType, principalName, policyArn string) bool
	ListAttachedPolicies(principalType, principalName string) ([]string, error)
	ListPrincipalsForPolicy(policyArn string) ([]AttachedPolicyRef, error)
	DetachAllForPrincipal(principalType, principalName string) error
	DetachAllForPolicy(policyArn string) error
	MigratePrincipal(oldName, newName, principalType string) error
}

// GroupStoreInterface defines operations for managing IAM groups.
type GroupStoreInterface interface {
	Get(groupName string) (*Group, error)
	List(pathPrefix string, marker string, maxItems int) (*GroupListResult, error)
	Put(group *Group) error
	Delete(groupName string) error
	Exists(groupName string) bool
	Create(groupName, path, accountId string) (*Group, error)
	Update(group *Group) error
	Count() int
}

// InstanceProfileStoreInterface defines operations for managing IAM instance profiles.
type InstanceProfileStoreInterface interface {
	Get(instanceProfileName string) (*InstanceProfile, error)
	GetByID(instanceProfileID string) (*InstanceProfile, error)
	List(pathPrefix string, marker string, maxItems int) (*InstanceProfileListResult, error)
	Put(profile *InstanceProfile) error
	Delete(instanceProfileName string) error
	Exists(instanceProfileName string) bool
	Create(instanceProfileName, path, accountId string, tags []types.Tag) (*InstanceProfile, error)
	AddRole(instanceProfileName, roleName string) error
	RemoveRole(instanceProfileName, roleName string) error
	HasRole(instanceProfileName, roleName string) (bool, error)
	Count() int
	ListForRole(roleName string, marker string, maxItems int) (*InstanceProfileListResult, error)
}

// PolicyStoreInterface defines operations for managing IAM policies.
type PolicyStoreInterface interface {
	Get(policyArn string) (*Policy, error)
	List(scope, pathPrefix string, onlyAttached bool, marker string, maxItems int) (*PolicyListResult, error)
	Put(policy *Policy) error
	Delete(policyArn string) error
	Exists(policyArn string) bool
	Create(policyName, path, accountId, document, description string, tags []types.Tag) (*Policy, error)
	Count() int
	PutVersion(version *PolicyVersion) error
	GetVersion(policyArn, versionId string) (*PolicyVersion, error)
	DeleteVersion(policyArn, versionId string) error
	ListVersions(policyArn string, marker string, maxItems int) (*PolicyVersionListResult, error)
	GetDefaultVersion(policyArn string) (*PolicyVersion, error)
	SetDefaultVersion(policyArn, versionId string) error
	CreateVersion(policyArn, document string, setAsDefault bool, maxVersions int) (*PolicyVersion, error)
	GetMaxVersion(policyArn string) (int, error)
	CountVersions(policyArn string) (int, error)
	DeleteAllVersions(policyArn string) error
	IncrementAttachmentCount(policyArn string) error
	DecrementAttachmentCount(policyArn string) error
	IncrementPermissionsBoundaryUsageCount(policyArn string) error
	DecrementPermissionsBoundaryUsageCount(policyArn string) error
}

// InlinePolicyStoreInterface defines operations for managing IAM inline policies.
type InlinePolicyStoreInterface interface {
	Get(principalType, principalName, policyName string) (*InlinePolicy, error)
	List(principalType, principalName string) ([]string, error)
	Put(principalType, principalName, policyName, document string) error
	Delete(principalType, principalName, policyName string) error
	Exists(principalType, principalName, policyName string) bool
	DeleteAllForPrincipal(principalType, principalName string) error
	MigratePrincipal(oldName, newName, principalType string) error
}

// MFADeviceStoreInterface defines operations for managing IAM MFA devices.
type MFADeviceStoreInterface interface {
	Get(serialNumber string) (*VirtualMFADevice, error)
	Put(device *VirtualMFADevice) error
	Delete(serialNumber string) error
	Exists(serialNumber string) bool
	Create(accountId string, deviceName string, tags []types.Tag) (*VirtualMFADevice, error)
	EnableForUser(serialNumber, userName string) error
	EnableForUserWithLimit(serialNumber, userName string, maxDevices int) error
	Deactivate(serialNumber string) error
	ListForUser(userName string, marker string, maxItems int) (*MFADeviceListResult, error)
	ListVirtual(assignmentStatus, marker string, maxItems int) (*MFADeviceListResult, error)
	Resync(serialNumber string) error
	Count() int
	MigrateUser(oldUserName, newUserName string) error
	CountForUser(userName string) (int, error)
}

// PasswordPolicyStoreInterface defines operations for managing IAM account password policies.
type PasswordPolicyStoreInterface interface {
	Get() (*AccountPasswordPolicy, error)
	Put(policy *AccountPasswordPolicy) error
	Exists() bool
	Delete() error
	GetOrDefault() *AccountPasswordPolicy
	DefaultPolicy() *AccountPasswordPolicy
	ParameterDefaults() *AccountPasswordPolicy
}

// IAMStoreInterface is the cross-package IAM store contract: the STS cores
// and the server wiring (dispatcher, authoriser, HTTP assembly, optional
// apps) hold this interface instead of the concrete *IAMStore. It covers
// the store domains those consumers reach plus AccountID; the later
// domains (account aliases, certificates, providers, settings and the
// service-linked-role tasks) have no cross-package consumer and are
// reached on the concrete type — widening the facade ahead of a consumer
// only forks it.
type IAMStoreInterface interface {
	Users() UserStoreInterface
	AccessKeys() AccessKeyStoreInterface
	LoginProfiles() LoginProfileStoreInterface
	Groups() GroupStoreInterface
	UserGroups() UserGroupStoreInterface
	Roles() RoleStoreInterface
	InstanceProfiles() InstanceProfileStoreInterface
	Policies() PolicyStoreInterface
	InlinePolicies() InlinePolicyStoreInterface
	AttachedPolicies() AttachedPolicyStoreInterface
	MFADevices() MFADeviceStoreInterface
	PasswordPolicy() PasswordPolicyStoreInterface
	AccountID() string
}

var (
	_ RoleStoreInterface            = (*RoleStore)(nil)
	_ AccessKeyStoreInterface       = (*AccessKeyStore)(nil)
	_ LoginProfileStoreInterface    = (*LoginProfileStore)(nil)
	_ UserStoreInterface            = (*UserStore)(nil)
	_ UserGroupStoreInterface       = (*UserGroupStore)(nil)
	_ AttachedPolicyStoreInterface  = (*AttachedPolicyStore)(nil)
	_ GroupStoreInterface           = (*GroupStore)(nil)
	_ InstanceProfileStoreInterface = (*InstanceProfileStore)(nil)
	_ PolicyStoreInterface          = (*PolicyStore)(nil)
	_ InlinePolicyStoreInterface    = (*InlinePolicyStore)(nil)
	_ MFADeviceStoreInterface       = (*MFADeviceStore)(nil)
	_ PasswordPolicyStoreInterface  = (*PasswordPolicyStore)(nil)
	_ IAMStoreInterface             = (*IAMStore)(nil)
)
