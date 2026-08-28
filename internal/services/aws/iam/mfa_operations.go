// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"context"
	"encoding/base64"
	"strconv"
	"strings"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/timeutils"
)

// CreateVirtualMFADevice creates a new virtual MFA device.
// Tags are optional.
// Returns the virtual MFA device with base32 seed and secret.
func (s *IAMService) CreateVirtualMFADevice(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &CreateVirtualMFADeviceInput{
		VirtualMFADeviceName: request.GetStringParam(req.Parameters, "VirtualMFADeviceName"),
		Tags:                 tags.ParseTagsWithQueryFallback(req.Parameters, "Tags"),
	}
	device, err := s.createVirtualMFADeviceCore(store, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"VirtualMFADevice": mfaDeviceToResponse(device, nil, true),
	}, nil
}

// DeleteVirtualMFADevice deletes a virtual MFA device by its serial number.
// SerialNumber is required.
// Returns an error if the device is still assigned to a user.
func (s *IAMService) DeleteVirtualMFADevice(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteVirtualMFADeviceCore(store, request.GetStringParam(req.Parameters, "SerialNumber")); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// EnableMFADevice enables an MFA device for a user.
// UserName, SerialNumber, AuthenticationCode1, and AuthenticationCode2 are required.
func (s *IAMService) EnableMFADevice(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &EnableMFADeviceInput{
		UserName:            request.GetStringParam(req.Parameters, "UserName"),
		SerialNumber:        request.GetStringParam(req.Parameters, "SerialNumber"),
		AuthenticationCode1: request.GetStringParam(req.Parameters, "AuthenticationCode1"),
		AuthenticationCode2: request.GetStringParam(req.Parameters, "AuthenticationCode2"),
	}
	if err := s.enableMFADeviceCore(store, input); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// DeactivateMFADevice deactivates an MFA device for a user.
// UserName and SerialNumber are required.
func (s *IAMService) DeactivateMFADevice(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	serialNumber := request.GetStringParam(req.Parameters, "SerialNumber")

	userName, err := resolveUserName(reqCtx, userName)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deactivateMFADeviceCore(store, userName, serialNumber); err != nil {
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
	maxItems := pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.listMFADevicesCore(store, userName, marker, maxItems)
	if err != nil {
		return nil, err
	}

	devices := make([]interface{}, len(result.Devices))
	for i, device := range result.Devices {
		devices[i] = mfaDeviceToListResponse(device)
	}

	resp := map[string]interface{}{
		"MFADevices":  devices,
		"IsTruncated": result.IsTruncated,
	}

	if result.Marker != "" {
		resp["Marker"] = result.Marker
	}

	return resp, nil
}

// ListVirtualMFADevices lists all virtual MFA devices.
// Supports pagination via Marker and MaxItems.
func (s *IAMService) ListVirtualMFADevices(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	assignmentStatus := request.GetStringParam(req.Parameters, "AssignmentStatus")
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.listVirtualMFADevicesCore(store, assignmentStatus, marker, maxItems)
	if err != nil {
		return nil, err
	}

	devices := make([]interface{}, len(result.Devices))
	for i, device := range result.Devices {
		devices[i] = mfaDeviceToResponse(device, result.UsersBySerial[device.SerialNumber], false)
	}

	resp := map[string]interface{}{
		"VirtualMFADevices": devices,
		"IsTruncated":       result.IsTruncated,
	}

	if result.Marker != "" {
		resp["Marker"] = result.Marker
	}

	return resp, nil
}

// GetMFADevice retrieves information about an MFA device.
// SerialNumber is required.
// Since only virtual MFA devices are supported, returns device info if the
// serial number matches a virtual MFA device.
func (s *IAMService) GetMFADevice(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	device, err := s.getMFADeviceCore(store, request.GetStringParam(req.Parameters, "SerialNumber"))
	if err != nil {
		return nil, err
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

	// The operation output is flat: SerialNumber, EnableDate and UserName
	// sit at the top level without an MFADevice wrapper.
	return resp, nil
}

// ResyncMFADevice resynchronises an MFA device for a user.
// UserName, SerialNumber, AuthenticationCode1, and AuthenticationCode2 are required.
func (s *IAMService) ResyncMFADevice(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &ResyncMFADeviceInput{
		UserName:            request.GetStringParam(req.Parameters, "UserName"),
		SerialNumber:        request.GetStringParam(req.Parameters, "SerialNumber"),
		AuthenticationCode1: request.GetStringParam(req.Parameters, "AuthenticationCode1"),
		AuthenticationCode2: request.GetStringParam(req.Parameters, "AuthenticationCode2"),
	}
	if err := s.resyncMFADeviceCore(store, input); err != nil {
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
	policy, err := s.getAccountPasswordPolicyCore(store)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"PasswordPolicy": passwordPolicyToResponse(policy),
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
	input := &UpdateAccountPasswordPolicyInput{
		MinimumPasswordLength: request.GetIntParam(req.Parameters, "MinimumPasswordLength"),
	}
	if v, ok := req.Parameters["RequireSymbols"]; ok {
		b := toBool(v)
		input.RequireSymbols = &b
	}
	if v, ok := req.Parameters["RequireNumbers"]; ok {
		b := toBool(v)
		input.RequireNumbers = &b
	}
	if v, ok := req.Parameters["RequireUppercaseCharacters"]; ok {
		b := toBool(v)
		input.RequireUppercaseCharacters = &b
	}
	if v, ok := req.Parameters["RequireLowercaseCharacters"]; ok {
		b := toBool(v)
		input.RequireLowercaseCharacters = &b
	}
	if v, ok := req.Parameters["AllowUsersToChangePassword"]; ok {
		b := toBool(v)
		input.AllowUsersToChangePassword = &b
	}
	if v, ok := req.Parameters["HardExpiry"]; ok {
		b := toBool(v)
		input.HardExpiry = &b
	}
	if v, ok := req.Parameters["MaxPasswordAge"]; ok {
		n := toInt(v)
		input.MaxPasswordAge = &n
	}
	if v, ok := req.Parameters["PasswordReusePrevention"]; ok {
		n := toInt(v)
		input.PasswordReusePrevention = &n
	}
	if err := s.updateAccountPasswordPolicyCore(store, input); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
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
		i, err := strconv.Atoi(val)
		if err != nil {
			return 0
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
	if err := s.deleteAccountPasswordPolicyCore(store); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

var mfaDeviceTagOps = tagOps[*iamstore.VirtualMFADevice]{
	paramName:  "SerialNumber",
	emptyErr:   ErrNoSuchMFADevice,
	notFoundFn: func(n string) error { return NewNoSuchMFADeviceError(n) },
	getFn:      func(s *iamstore.IAMStore, n string) (*iamstore.VirtualMFADevice, error) { return s.MFADevices().Get(n) },
	putFn:      func(s *iamstore.IAMStore, r *iamstore.VirtualMFADevice) error { return s.MFADevices().Put(r) },
	tagsFn:     func(r *iamstore.VirtualMFADevice) *[]tags.Tag { return &r.Tags },
}

// TagMFADevice adds tags to a virtual MFA device.
// SerialNumber is required.
// Tags are provided as a list of key-value pairs.
func (s *IAMService) TagMFADevice(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return tagResource(ctx, s, reqCtx, req, mfaDeviceTagOps)
}

// UntagMFADevice removes tags from a virtual MFA device.
// SerialNumber is required.
// TagKeys specifies which tags to remove.
func (s *IAMService) UntagMFADevice(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return untagResource(ctx, s, reqCtx, req, mfaDeviceTagOps)
}

// ListMFADeviceTags lists the tags attached to a virtual MFA device.
// SerialNumber is required.
func (s *IAMService) ListMFADeviceTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return listResourceTags(ctx, s, reqCtx, req, mfaDeviceTagOps)
}

func mfaDeviceToResponse(device *iamstore.VirtualMFADevice, user *iamstore.User, includeSecret bool) map[string]interface{} {
	resp := map[string]interface{}{
		"SerialNumber": device.SerialNumber,
	}

	if device.FriendlyName != "" {
		resp["FriendlyName"] = device.FriendlyName
	}

	// Base32StringSeed is a Smithy blob (BootstrapDatum).  Per AWS docs
	// ("Type: Base64-encoded binary data object") the blob content is the
	// UTF-8 bytes of the base32 seed string.  The previous implementation
	// base32-decoded the seed to raw 20 bytes then base64-encoded those,
	// producing a value that TOTP clients could not consume.
	if includeSecret && device.Base32StringSeed != "" {
		resp["Base32StringSeed"] = base64.StdEncoding.EncodeToString([]byte(device.Base32StringSeed))
	}

	if device.EnableDate != nil {
		resp["EnableDate"] = device.EnableDate.Format(timeutils.ISO8601SimpleFormat)
	}

	if device.UserAssignment != nil && user != nil {
		resp["User"] = map[string]interface{}{
			"UserName":   user.UserName,
			"UserId":     user.ID,
			"Arn":        user.Arn,
			"Path":       user.Path,
			"CreateDate": user.CreateDate.Format(timeutils.ISO8601SimpleFormat),
		}
	}

	if tags := tags.ToResponse(device.Tags); tags != nil {
		resp["Tags"] = tags
	}

	return resp
}

func mfaDeviceToListResponse(device *iamstore.VirtualMFADevice) map[string]interface{} {
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

func passwordPolicyToResponse(policy *iamstore.AccountPasswordPolicy) map[string]interface{} {
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
	summaryMap, err := s.getAccountSummaryCore(store)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"SummaryMap": summaryMap,
	}, nil
}
