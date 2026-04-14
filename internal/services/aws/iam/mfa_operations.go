// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"context"
	"encoding/base64"
	"strconv"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/crypto"
	"vorpalstacks/internal/utils/timeutils"
)

// CreateVirtualMFADevice creates a new virtual MFA device.
// Tags are optional.
// Returns the virtual MFA device with base32 seed and secret.
func (s *IAMService) CreateVirtualMFADevice(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	newTags := tags.ParseTagsWithQueryFallback(req.Parameters, "Tags")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	device, err := store.MFADevices().Create(s.accountID, newTags)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"VirtualMFADevice": s.mfaDeviceToResponse(reqCtx, device, true),
	}, nil
}

// DeleteVirtualMFADevice deletes a virtual MFA device by its serial number.
// SerialNumber is required.
// Returns an error if the device is still assigned to a user.
func (s *IAMService) DeleteVirtualMFADevice(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	serialNumber := request.GetStringParam(req.Parameters, "SerialNumber")
	if serialNumber == "" {
		return nil, ErrNoSuchMFADevice
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.MFADevices().Exists(serialNumber) {
		return nil, NewNoSuchMFADeviceError(serialNumber)
	}

	device, err := store.MFADevices().Get(serialNumber)
	if err != nil {
		return nil, err
	}

	if device.UserAssignment != nil {
		return nil, NewMFADeviceStillAssignedError(serialNumber)
	}

	if err := store.MFADevices().Delete(serialNumber); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// EnableMFADevice enables an MFA device for a user.
// UserName, SerialNumber, AuthenticationCode1, and AuthenticationCode2 are required.
func (s *IAMService) EnableMFADevice(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	serialNumber := request.GetStringParam(req.Parameters, "SerialNumber")
	authCode1 := request.GetStringParam(req.Parameters, "AuthenticationCode1")
	authCode2 := request.GetStringParam(req.Parameters, "AuthenticationCode2")

	if userName == "" {
		return nil, ErrNoSuchUser
	}
	if serialNumber == "" {
		return nil, ErrNoSuchMFADevice
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	device, err := store.MFADevices().Get(serialNumber)
	if err != nil {
		return nil, NewNoSuchMFADeviceError(serialNumber)
	}

	if device.UserAssignment != nil {
		return nil, NewMFADeviceAlreadyAssignedError(serialNumber)
	}

	if device.Base32StringSeed == "" {
		return nil, ErrInvalidAuthenticationCode
	}

	if authCode1 == "" || authCode2 == "" {
		return nil, ErrInvalidAuthenticationCode
	}

	if err := crypto.ValidateConsecutiveTOTPCodes(device.Base32StringSeed, authCode1, authCode2); err != nil {
		return nil, ErrInvalidAuthenticationCode
	}

	if err := store.MFADevices().EnableForUser(serialNumber, userName); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DeactivateMFADevice deactivates an MFA device for a user.
// UserName and SerialNumber are required.
func (s *IAMService) DeactivateMFADevice(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	serialNumber := request.GetStringParam(req.Parameters, "SerialNumber")

	if userName == "" {
		return nil, ErrNoSuchUser
	}
	if serialNumber == "" {
		return nil, ErrNoSuchMFADevice
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	device, err := store.MFADevices().Get(serialNumber)
	if err != nil {
		return nil, NewNoSuchMFADeviceError(serialNumber)
	}

	if device.UserAssignment == nil || device.UserAssignment.UserName != userName {
		return nil, NewNoSuchMFADeviceError(serialNumber)
	}

	if err := store.MFADevices().Deactivate(serialNumber); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListMFADevices lists MFA devices.
// UserName is optional; if specified, only lists devices for that user.
// Supports pagination via Marker and MaxItems.
func (s *IAMService) ListMFADevices(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := getMaxItems(req.Parameters)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if userName != "" && !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	result, err := store.MFADevices().ListForUser(userName, marker, maxItems)
	if err != nil {
		return nil, err
	}

	devices := make([]interface{}, len(result.MFADevices))
	for i, device := range result.MFADevices {
		devices[i] = s.mfaDeviceToListResponse(device)
	}

	response := map[string]interface{}{
		"MFADevices":  devices,
		"IsTruncated": result.IsTruncated,
	}

	if result.Marker != "" {
		response["Marker"] = result.Marker
	}

	return response, nil
}

// ListVirtualMFADevices lists all virtual MFA devices.
// Supports pagination via Marker and MaxItems.
func (s *IAMService) ListVirtualMFADevices(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := getMaxItems(req.Parameters)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := store.MFADevices().ListVirtual(marker, maxItems)
	if err != nil {
		return nil, err
	}

	devices := make([]interface{}, len(result.MFADevices))
	for i, device := range result.MFADevices {
		devices[i] = s.mfaDeviceToResponse(reqCtx, device, false)
	}

	response := map[string]interface{}{
		"VirtualMFADevices": devices,
		"IsTruncated":       result.IsTruncated,
	}

	if result.Marker != "" {
		response["Marker"] = result.Marker
	}

	return response, nil
}

// GetMFADevice retrieves information about an MFA device.
// SerialNumber is required.
// Since only virtual MFA devices are supported, returns device info if the
// serial number matches a virtual MFA device.
func (s *IAMService) GetMFADevice(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	serialNumber := request.GetStringParam(req.Parameters, "SerialNumber")
	if serialNumber == "" {
		return nil, ErrNoSuchMFADevice
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	device, err := store.MFADevices().Get(serialNumber)
	if err != nil {
		return nil, NewNoSuchMFADeviceError(serialNumber)
	}

	resp := map[string]interface{}{
		"SerialNumber": device.SerialNumber,
	}

	if device.EnableDate != nil {
		resp["EnableDate"] = device.EnableDate.Format(timeutils.ISO8601SimpleFormat)
	}

	if device.UserAssignment != nil {
		resp["UserName"] = device.UserAssignment.UserName
	}

	return map[string]interface{}{
		"MFADevice": resp,
	}, nil
}

// ResyncMFADevice resynchronises an MFA device for a user.
// UserName, SerialNumber, AuthenticationCode1, and AuthenticationCode2 are required.
func (s *IAMService) ResyncMFADevice(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	serialNumber := request.GetStringParam(req.Parameters, "SerialNumber")
	authCode1 := request.GetStringParam(req.Parameters, "AuthenticationCode1")
	authCode2 := request.GetStringParam(req.Parameters, "AuthenticationCode2")

	if userName == "" {
		return nil, ErrNoSuchUser
	}
	if serialNumber == "" {
		return nil, ErrNoSuchMFADevice
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	device, err := store.MFADevices().Get(serialNumber)
	if err != nil {
		return nil, NewNoSuchMFADeviceError(serialNumber)
	}

	if device.UserAssignment == nil || device.UserAssignment.UserName != userName {
		return nil, NewNoSuchMFADeviceError(serialNumber)
	}

	if err := crypto.ValidateConsecutiveTOTPCodes(device.Base32StringSeed, authCode1, authCode2); err != nil {
		return nil, ErrInvalidAuthenticationCode
	}

	if err := store.MFADevices().Resync(serialNumber); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetAccountPasswordPolicy retrieves the account password policy.
// Returns an error if no password policy has been set.
func (s *IAMService) GetAccountPasswordPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.PasswordPolicy().Exists() {
		return nil, NewNoSuchPasswordPolicyError()
	}

	policy, err := store.PasswordPolicy().Get()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"PasswordPolicy": s.passwordPolicyToResponse(reqCtx, policy),
	}, nil
}

// UpdateAccountPasswordPolicy updates the account password policy.
// Supports various password policy parameters such as MinimumPasswordLength,
// RequireSymbols, RequireNumbers, RequireUppercaseCharacters, RequireLowercaseCharacters,
// AllowUsersToChangePassword, HardExpiry, MaxPasswordAge, and PasswordReusePrevention.
func (s *IAMService) UpdateAccountPasswordPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	policy := store.PasswordPolicy().GetOrDefault()

	minLength := request.GetIntParam(req.Parameters, "MinimumPasswordLength")
	if minLength > 0 {
		policy.MinimumPasswordLength = minLength
	}

	if requireSymbols, ok := req.Parameters["RequireSymbols"]; ok {
		policy.RequireSymbols = toBool(requireSymbols)
	}
	if requireNumbers, ok := req.Parameters["RequireNumbers"]; ok {
		policy.RequireNumbers = toBool(requireNumbers)
	}
	if requireUppercase, ok := req.Parameters["RequireUppercaseCharacters"]; ok {
		policy.RequireUppercaseCharacters = toBool(requireUppercase)
	}
	if requireLowercase, ok := req.Parameters["RequireLowercaseCharacters"]; ok {
		policy.RequireLowercaseCharacters = toBool(requireLowercase)
	}
	if allowChange, ok := req.Parameters["AllowUsersToChangePassword"]; ok {
		policy.AllowUsersToChangePassword = toBool(allowChange)
	}
	if hardExpiry, ok := req.Parameters["HardExpiry"]; ok {
		policy.HardExpiry = toBool(hardExpiry)
	}
	if maxAge, ok := req.Parameters["MaxPasswordAge"]; ok {
		policy.MaxPasswordAge = toInt(maxAge)
	}
	if reusePrevention, ok := req.Parameters["PasswordReusePrevention"]; ok {
		policy.PasswordReusePrevention = toInt(reusePrevention)
	}

	if err := store.PasswordPolicy().Put(policy); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"PasswordPolicy": s.passwordPolicyToResponse(reqCtx, policy),
	}, nil
}

func toBool(v interface{}) bool {
	switch val := v.(type) {
	case bool:
		return val
	case string:
		return strings.ToLower(val) == "true"
	}
	return false
}

func toInt(v interface{}) int {
	switch val := v.(type) {
	case int:
		return val
	case int64:
		return int(val)
	case float64:
		return int(val)
	case string:
		var i int
		for _, c := range val {
			if c >= '0' && c <= '9' {
				i = i*10 + int(c-'0')
			}
		}
		return i
	}
	return 0
}

// DeleteAccountPasswordPolicy deletes the account password policy.
// Returns an error if no password policy has been set.
func (s *IAMService) DeleteAccountPasswordPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.PasswordPolicy().Exists() {
		return nil, NewNoSuchPasswordPolicyError()
	}

	if err := store.PasswordPolicy().Delete(); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// TagMFADevice adds tags to a virtual MFA device.
// SerialNumber is required.
// Tags are provided as a list of key-value pairs.
func (s *IAMService) TagMFADevice(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	serialNumber := request.GetStringParam(req.Parameters, "SerialNumber")
	if serialNumber == "" {
		return nil, ErrNoSuchMFADevice
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	device, err := store.MFADevices().Get(serialNumber)
	if err != nil {
		return nil, NewNoSuchMFADeviceError(serialNumber)
	}

	device.Tags = tags.Apply(device.Tags, tags.ParseTagsWithQueryFallback(req.Parameters, "Tags"))

	if err := store.MFADevices().Put(device); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// UntagMFADevice removes tags from a virtual MFA device.
// SerialNumber is required.
// TagKeys specifies which tags to remove.
func (s *IAMService) UntagMFADevice(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	serialNumber := request.GetStringParam(req.Parameters, "SerialNumber")
	if serialNumber == "" {
		return nil, ErrNoSuchMFADevice
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	device, err := store.MFADevices().Get(serialNumber)
	if err != nil {
		return nil, NewNoSuchMFADeviceError(serialNumber)
	}

	device.Tags = removeTags(device.Tags, parseTagKeys(req.Parameters))

	if err := store.MFADevices().Put(device); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListMFADeviceTags lists the tags attached to a virtual MFA device.
// SerialNumber is required.
func (s *IAMService) ListMFADeviceTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	serialNumber := request.GetStringParam(req.Parameters, "SerialNumber")
	if serialNumber == "" {
		return nil, ErrNoSuchMFADevice
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	device, err := store.MFADevices().Get(serialNumber)
	if err != nil {
		return nil, NewNoSuchMFADeviceError(serialNumber)
	}

	return map[string]interface{}{
		"Tags":        tagsToResponse(device.Tags),
		"IsTruncated": false,
	}, nil
}

func (s *IAMService) mfaDeviceToResponse(reqCtx *request.RequestContext, device *iamstore.VirtualMFADevice, includeSecret bool) map[string]interface{} {
	resp := map[string]interface{}{
		"SerialNumber": device.SerialNumber,
	}

	if includeSecret && device.Base32StringSeed != "" {
		resp["Base32StringSeed"] = base64.StdEncoding.EncodeToString([]byte(device.Base32StringSeed))
	}

	if device.EnableDate != nil {
		resp["EnableDate"] = device.EnableDate.Format(timeutils.ISO8601SimpleFormat)
	}

	if device.UserAssignment != nil {
		store, err := s.store(reqCtx)
		if err == nil {
			if user, err := store.Users().Get(device.UserAssignment.UserName); err == nil {
				resp["User"] = map[string]interface{}{
					"UserName":   user.UserName,
					"UserId":     user.ID,
					"Arn":        user.Arn,
					"Path":       user.Path,
					"CreateDate": user.CreateDate.Format(timeutils.ISO8601SimpleFormat),
				}
			}
		}
	}

	if tags := tagsToResponse(device.Tags); tags != nil {
		resp["Tags"] = tags
	}

	return resp
}

func (s *IAMService) mfaDeviceToListResponse(device *iamstore.VirtualMFADevice) map[string]interface{} {
	resp := map[string]interface{}{
		"SerialNumber": device.SerialNumber,
	}

	if device.UserAssignment != nil {
		resp["UserName"] = device.UserAssignment.UserName
	}

	if device.EnableDate != nil {
		resp["EnableDate"] = device.EnableDate.Format(timeutils.ISO8601SimpleFormat)
	}

	return resp
}

func (s *IAMService) passwordPolicyToResponse(reqCtx *request.RequestContext, policy *iamstore.AccountPasswordPolicy) map[string]interface{} {
	return map[string]interface{}{
		"MinimumPasswordLength":      policy.MinimumPasswordLength,
		"RequireSymbols":             policy.RequireSymbols,
		"RequireNumbers":             policy.RequireNumbers,
		"RequireUppercaseCharacters": policy.RequireUppercaseCharacters,
		"RequireLowercaseCharacters": policy.RequireLowercaseCharacters,
		"AllowUsersToChangePassword": policy.AllowUsersToChangePassword,
		"MaxPasswordAge":             policy.MaxPasswordAge,
		"PasswordReusePrevention":    policy.PasswordReusePrevention,
		"HardExpiry":                 policy.HardExpiry,
		"ExpirePasswords":            policy.ExpirePasswords,
	}
}

// GetAccountSummary retrieves account-level summary information.
// Returns counts of users, groups, roles, policies, instance profiles, and MFA devices.
// Also returns quota limits for each resource type.
func (s *IAMService) GetAccountSummary(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	users := store.Users().Count()
	groups := store.Groups().Count()
	roles := store.Roles().Count()
	policies := store.Policies().Count()
	instanceProfiles := store.InstanceProfiles().Count()
	mfaDevices := store.MFADevices().Count()

	summaryMap := map[string]string{
		"Users":                     strconv.Itoa(users),
		"Groups":                    strconv.Itoa(groups),
		"Roles":                     strconv.Itoa(roles),
		"LocalManagedPolicies":      strconv.Itoa(policies),
		"InstanceProfiles":          strconv.Itoa(instanceProfiles),
		"MFADevices":                strconv.Itoa(mfaDevices),
		"ServerCertificates":        "0",
		"UsersQuota":                "5000",
		"GroupsQuota":               "300",
		"RolesQuota":                "1000",
		"InstanceProfilesQuota":     "500",
		"LocalManagedPoliciesQuota": "1500",
		"MFADevicesQuota":           "500",
		"ServerCertificatesQuota":   "20",
	}

	return map[string]interface{}{
		"SummaryMap": summaryMap,
	}, nil
}
