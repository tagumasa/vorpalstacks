package iam

import (
	"strconv"

	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/crypto"
)

// CreateVirtualMFADeviceInput carries the parsed CreateVirtualMFADevice
// request.
type CreateVirtualMFADeviceInput struct {
	VirtualMFADeviceName string
	Tags                 []tags.Tag
}

// EnableMFADeviceInput carries the parsed EnableMFADevice request.
type EnableMFADeviceInput struct {
	UserName            string
	SerialNumber        string
	AuthenticationCode1 string
	AuthenticationCode2 string
}

// ResyncMFADeviceInput carries the parsed ResyncMFADevice request.
type ResyncMFADeviceInput struct {
	UserName            string
	SerialNumber        string
	AuthenticationCode1 string
	AuthenticationCode2 string
}

// UpdateAccountPasswordPolicyInput carries the parsed
// UpdateAccountPasswordPolicy request.  Pointer members distinguish an
// omitted parameter (nil) from an explicit value so the replace semantics
// keep the documented default for omitted members.
type UpdateAccountPasswordPolicyInput struct {
	MinimumPasswordLength      int
	RequireSymbols             *bool
	RequireNumbers             *bool
	RequireUppercaseCharacters *bool
	RequireLowercaseCharacters *bool
	AllowUsersToChangePassword *bool
	HardExpiry                 *bool
	MaxPasswordAge             *int
	PasswordReusePrevention    *int
}

// MFADevicesResult is the Core result for ListMFADevices.
type MFADevicesResult struct {
	Devices     []*iamstore.VirtualMFADevice
	IsTruncated bool
	Marker      string
}

// VirtualMFADeviceListResult is the Core result for ListVirtualMFADevices.
// UsersBySerial carries the owning user for each assigned device so the
// serialiser can populate the User member without store access.
type VirtualMFADeviceListResult struct {
	Devices       []*iamstore.VirtualMFADevice
	UsersBySerial map[string]*iamstore.User
	IsTruncated   bool
	Marker        string
}

// createVirtualMFADeviceCore validates input and creates a virtual MFA
// device.
func (s *IAMService) createVirtualMFADeviceCore(store *iamstore.IAMStore, input *CreateVirtualMFADeviceInput) (*iamstore.VirtualMFADevice, error) {
	if input.VirtualMFADeviceName == "" {
		return nil, ErrInvalidInput
	}
	if err := validateNewTags(input.Tags); err != nil {
		return nil, err
	}
	return store.MFADevices().Create(s.accountID, input.VirtualMFADeviceName, input.Tags)
}

// deleteVirtualMFADeviceCore validates input and deletes an unassigned
// virtual MFA device.
func (s *IAMService) deleteVirtualMFADeviceCore(store *iamstore.IAMStore, serialNumber string) error {
	if serialNumber == "" {
		return NewValidationError("SerialNumber")
	}
	if !store.MFADevices().Exists(serialNumber) {
		return NewNoSuchMFADeviceError(serialNumber)
	}

	device, err := store.MFADevices().Get(serialNumber)
	if err != nil {
		return err
	}

	if device.UserAssignment != nil {
		return NewMFADeviceStillAssignedError(serialNumber)
	}

	return store.MFADevices().Delete(serialNumber)
}

// enableMFADeviceCore validates input and assigns the virtual MFA device to
// the user after verifying the two consecutive authentication codes.
func (s *IAMService) enableMFADeviceCore(store *iamstore.IAMStore, input *EnableMFADeviceInput) error {
	if input.UserName == "" {
		return NewValidationError("UserName")
	}
	if input.SerialNumber == "" {
		return NewValidationError("SerialNumber")
	}

	if !store.Users().Exists(input.UserName) {
		return NewNoSuchUserError(input.UserName)
	}

	device, err := store.MFADevices().Get(input.SerialNumber)
	if err != nil {
		return NewNoSuchMFADeviceError(input.SerialNumber)
	}

	if device.UserAssignment != nil {
		return NewMFADeviceAlreadyAssignedError(input.SerialNumber)
	}

	if device.Base32StringSeed == "" {
		return ErrInvalidAuthenticationCode
	}

	if input.AuthenticationCode1 == "" || input.AuthenticationCode2 == "" {
		return ErrInvalidAuthenticationCode
	}

	if err := crypto.ValidateConsecutiveTOTPCodes(device.Base32StringSeed, input.AuthenticationCode1, input.AuthenticationCode2); err != nil {
		return ErrInvalidAuthenticationCode
	}

	return store.MFADevices().EnableForUser(input.SerialNumber, input.UserName)
}

// deactivateMFADeviceCore validates input and deactivates the MFA device
// assigned to the named user.
func (s *IAMService) deactivateMFADeviceCore(store *iamstore.IAMStore, userName, serialNumber string) error {
	if serialNumber == "" {
		return NewValidationError("SerialNumber")
	}
	if !store.Users().Exists(userName) {
		return NewNoSuchUserError(userName)
	}

	device, err := store.MFADevices().Get(serialNumber)
	if err != nil {
		return NewNoSuchMFADeviceError(serialNumber)
	}

	if device.UserAssignment == nil || device.UserAssignment.UserName != userName {
		return NewNoSuchMFADeviceError(serialNumber)
	}

	return store.MFADevices().Deactivate(serialNumber)
}

// listMFADevicesCore lists the MFA devices, optionally scoped to a user.
func (s *IAMService) listMFADevicesCore(store *iamstore.IAMStore, userName, marker string, maxItems int) (*MFADevicesResult, error) {
	if userName != "" && !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	result, err := store.MFADevices().ListForUser(userName, marker, maxItems)
	if err != nil {
		return nil, err
	}
	return &MFADevicesResult{
		Devices:     result.MFADevices,
		IsTruncated: result.IsTruncated,
		Marker:      result.Marker,
	}, nil
}

// listVirtualMFADevicesCore lists the virtual MFA devices matching the
// assignment status, resolving the owning user for each assigned device.
func (s *IAMService) listVirtualMFADevicesCore(store *iamstore.IAMStore, assignmentStatus, marker string, maxItems int) (*VirtualMFADeviceListResult, error) {
	if assignmentStatus == "" {
		assignmentStatus = "Any"
	}
	if assignmentStatus != "Any" && assignmentStatus != "Assigned" && assignmentStatus != "Unassigned" {
		return nil, NewInvalidInputError("AssignmentStatus", "must be one of: Any, Assigned, Unassigned")
	}

	result, err := store.MFADevices().ListVirtual(assignmentStatus, marker, maxItems)
	if err != nil {
		return nil, err
	}

	usersBySerial := make(map[string]*iamstore.User)
	for _, device := range result.MFADevices {
		if device.UserAssignment != nil {
			if user, err := store.Users().Get(device.UserAssignment.UserName); err == nil {
				usersBySerial[device.SerialNumber] = user
			}
		}
	}
	return &VirtualMFADeviceListResult{
		Devices:       result.MFADevices,
		UsersBySerial: usersBySerial,
		IsTruncated:   result.IsTruncated,
		Marker:        result.Marker,
	}, nil
}

// getMFADeviceCore validates input and retrieves the MFA device.
func (s *IAMService) getMFADeviceCore(store *iamstore.IAMStore, serialNumber string) (*iamstore.VirtualMFADevice, error) {
	if serialNumber == "" {
		return nil, NewValidationError("SerialNumber")
	}
	device, err := store.MFADevices().Get(serialNumber)
	if err != nil {
		return nil, NewNoSuchMFADeviceError(serialNumber)
	}
	return device, nil
}

// resyncMFADeviceCore validates input and resynchronises the MFA device
// assigned to the named user.
func (s *IAMService) resyncMFADeviceCore(store *iamstore.IAMStore, input *ResyncMFADeviceInput) error {
	if input.UserName == "" {
		return NewValidationError("UserName")
	}
	if input.SerialNumber == "" {
		return NewValidationError("SerialNumber")
	}

	if !store.Users().Exists(input.UserName) {
		return NewNoSuchUserError(input.UserName)
	}

	device, err := store.MFADevices().Get(input.SerialNumber)
	if err != nil {
		return NewNoSuchMFADeviceError(input.SerialNumber)
	}

	if device.UserAssignment == nil || device.UserAssignment.UserName != input.UserName {
		return NewNoSuchMFADeviceError(input.SerialNumber)
	}

	if err := crypto.ValidateConsecutiveTOTPCodesForResync(device.Base32StringSeed, input.AuthenticationCode1, input.AuthenticationCode2); err != nil {
		return ErrInvalidAuthenticationCode
	}

	return store.MFADevices().Resync(input.SerialNumber)
}

// getAccountPasswordPolicyCore retrieves the account password policy.
func (s *IAMService) getAccountPasswordPolicyCore(store *iamstore.IAMStore) (*iamstore.AccountPasswordPolicy, error) {
	if !store.PasswordPolicy().Exists() {
		return nil, NewNoSuchPasswordPolicyError()
	}
	return store.PasswordPolicy().Get()
}

// updateAccountPasswordPolicyCore validates input and replaces the account
// password policy.  The operation replaces the whole policy: every
// parameter the request omits reverts to its documented default value
// instead of being merged with the previously stored policy.
func (s *IAMService) updateAccountPasswordPolicyCore(store *iamstore.IAMStore, input *UpdateAccountPasswordPolicyInput) error {
	policy := store.PasswordPolicy().ParameterDefaults()

	if input.MinimumPasswordLength != 0 {
		if input.MinimumPasswordLength < 6 || input.MinimumPasswordLength > 128 {
			return NewInvalidInputError("MinimumPasswordLength", "must be between 6 and 128")
		}
		policy.MinimumPasswordLength = input.MinimumPasswordLength
	}

	if input.RequireSymbols != nil {
		policy.RequireSymbols = *input.RequireSymbols
	}
	if input.RequireNumbers != nil {
		policy.RequireNumbers = *input.RequireNumbers
	}
	if input.RequireUppercaseCharacters != nil {
		policy.RequireUppercaseCharacters = *input.RequireUppercaseCharacters
	}
	if input.RequireLowercaseCharacters != nil {
		policy.RequireLowercaseCharacters = *input.RequireLowercaseCharacters
	}
	if input.AllowUsersToChangePassword != nil {
		policy.AllowUsersToChangePassword = *input.AllowUsersToChangePassword
	}
	if input.HardExpiry != nil {
		policy.HardExpiry = *input.HardExpiry
	}
	if input.MaxPasswordAge != nil {
		if *input.MaxPasswordAge < 1 || *input.MaxPasswordAge > 1095 {
			return NewInvalidInputError("MaxPasswordAge", "must be between 1 and 1095")
		}
		policy.MaxPasswordAge = *input.MaxPasswordAge
	}
	if input.PasswordReusePrevention != nil {
		if *input.PasswordReusePrevention < 1 || *input.PasswordReusePrevention > 24 {
			return NewInvalidInputError("PasswordReusePrevention", "must be between 1 and 24")
		}
		policy.PasswordReusePrevention = *input.PasswordReusePrevention
	}

	return store.PasswordPolicy().Put(policy)
}

// deleteAccountPasswordPolicyCore deletes the account password policy.
func (s *IAMService) deleteAccountPasswordPolicyCore(store *iamstore.IAMStore) error {
	if !store.PasswordPolicy().Exists() {
		return NewNoSuchPasswordPolicyError()
	}
	return store.PasswordPolicy().Delete()
}

// getAccountSummaryCore assembles the account-level summary counts and
// quota map.
func (s *IAMService) getAccountSummaryCore(store *iamstore.IAMStore) (map[string]string, error) {
	users := store.Users().Count()
	groups := store.Groups().Count()
	roles := store.Roles().Count()
	policies := store.Policies().Count()
	instanceProfiles := store.InstanceProfiles().Count()
	mfaDevices := store.MFADevices().Count()
	serverCertificates := store.ServerCertificates().Count()

	return map[string]string{
		"Users":                     strconv.Itoa(users),
		"Groups":                    strconv.Itoa(groups),
		"Roles":                     strconv.Itoa(roles),
		"LocalManagedPolicies":      strconv.Itoa(policies),
		"InstanceProfiles":          strconv.Itoa(instanceProfiles),
		"MFADevices":                strconv.Itoa(mfaDevices),
		"ServerCertificates":        strconv.Itoa(serverCertificates),
		"UsersQuota":                strconv.Itoa(iamstore.QuotaUsersPerAccount),
		"GroupsQuota":               strconv.Itoa(iamstore.QuotaGroupsPerAccount),
		"RolesQuota":                strconv.Itoa(iamstore.QuotaRolesPerAccount),
		"InstanceProfilesQuota":     strconv.Itoa(iamstore.QuotaInstanceProfilesPerAccount),
		"LocalManagedPoliciesQuota": strconv.Itoa(iamstore.QuotaLocalManagedPoliciesPerAccount),
		"MFADevicesQuota":           strconv.Itoa(iamstore.QuotaMFADevicesPerAccount),
		"ServerCertificatesQuota":   strconv.Itoa(iamstore.QuotaServerCertificatesPerAccount),
	}, nil
}
