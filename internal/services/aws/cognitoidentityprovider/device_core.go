package cognitoidentityprovider

import (
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// ConfirmDeviceInput carries the wire parameters of ConfirmDevice. Params
// holds the raw request parameter map; the nested DeviceSecretVerifierConfig
// structure and DeviceAttributes list are read from it inside the Core.
type ConfirmDeviceInput struct {
	AccessToken string
	DeviceKey   string
	DeviceName  string
	Params      map[string]interface{}
}

// GetDeviceInput carries the wire parameters of GetDevice.
type GetDeviceInput struct {
	AccessToken string
	DeviceKey   string
}

// ForgetDeviceInput carries the wire parameters of ForgetDevice.
type ForgetDeviceInput struct {
	AccessToken string
	DeviceKey   string
}

// ListDevicesInput carries the wire parameters of ListDevices. Params holds
// the raw request parameter map for the Limit member.
type ListDevicesInput struct {
	AccessToken     string
	PaginationToken string
	Params          map[string]interface{}
}

// UpdateDeviceStatusInput carries the wire parameters of UpdateDeviceStatus.
type UpdateDeviceStatusInput struct {
	AccessToken            string
	DeviceKey              string
	DeviceRememberedStatus string
}

// AdminGetDeviceInput carries the wire parameters of AdminGetDevice.
type AdminGetDeviceInput struct {
	UserPoolID string
	Username   string
	DeviceKey  string
}

// AdminForgetDeviceInput carries the wire parameters of AdminForgetDevice.
type AdminForgetDeviceInput struct {
	UserPoolID string
	Username   string
	DeviceKey  string
}

// AdminListDevicesInput carries the wire parameters of AdminListDevices.
// Params holds the raw request parameter map for the Limit member.
type AdminListDevicesInput struct {
	UserPoolID      string
	Username        string
	PaginationToken string
	Params          map[string]interface{}
}

// AdminUpdateDeviceStatusInput carries the wire parameters of
// AdminUpdateDeviceStatus.
type AdminUpdateDeviceStatusInput struct {
	UserPoolID             string
	Username               string
	DeviceKey              string
	DeviceRememberedStatus string
}

// confirmDeviceCore registers a new device for the authenticated user.
func (s *CognitoService) confirmDeviceCore(reqCtx *request.RequestContext, in ConfirmDeviceInput) (interface{}, error) {
	if in.AccessToken == "" || in.DeviceKey == "" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, in.AccessToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUserByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	now := time.Now().UTC()
	device := &cognitostore.Device{
		DeviceKey:              in.DeviceKey,
		UserPoolID:             user.UserPoolID,
		UserID:                 user.ID,
		DeviceName:             in.DeviceName,
		DeviceCreateDate:       now,
		DeviceLastModifiedDate: now,
		DeviceRememberedStatus: "remembered",
	}

	if svc, ok := in.Params["DeviceSecretVerifierConfig"]; ok {
		if m, ok := svc.(map[string]interface{}); ok {
			if sigVerifier, ok := m["PasswordVerifier"].(string); ok {
				device.DeviceSecretVerifierB = sigVerifier
			}
			if salt, ok := m["Salt"].(string); ok {
				device.DeviceSaltVerifier = salt
			}
		}
	}

	if attrs := parseDeviceAttributes(in.Params); len(attrs) > 0 {
		device.DeviceAttributes = attrs
	}

	if err := store.CreateDevice(device); err != nil {
		return nil, ErrInternalError
	}

	result := map[string]interface{}{
		"UserConfirmationNecessary": false,
	}

	pool, _ := store.GetUserPool(user.UserPoolID)
	if pool != nil && pool.DeviceConfiguration != nil && pool.DeviceConfiguration.ChallengeRequiredOnNewDevice {
		result["UserConfirmationNecessary"] = true
	}

	return result, nil
}

// getDeviceCore retrieves a device by its key.
func (s *CognitoService) getDeviceCore(reqCtx *request.RequestContext, in GetDeviceInput) (interface{}, error) {
	if in.AccessToken == "" || in.DeviceKey == "" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, in.AccessToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUserByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	device, err := store.GetDevice(user.UserPoolID, user.ID, in.DeviceKey)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	return map[string]interface{}{"Device": formatDevice(device)}, nil
}

// forgetDeviceCore removes a device for the authenticated user.
func (s *CognitoService) forgetDeviceCore(reqCtx *request.RequestContext, in ForgetDeviceInput) (interface{}, error) {
	if in.AccessToken == "" || in.DeviceKey == "" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, in.AccessToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUserByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if _, err := store.GetDevice(user.UserPoolID, user.ID, in.DeviceKey); err != nil {
		return nil, ErrResourceNotFound
	}

	if err := store.DeleteDevice(user.UserPoolID, user.ID, in.DeviceKey); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// listDevicesCore lists devices for the authenticated user.
func (s *CognitoService) listDevicesCore(reqCtx *request.RequestContext, in ListDevicesInput) (interface{}, error) {
	if in.AccessToken == "" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, in.AccessToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUserByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Smithy QueryLimitType: range {min: 0, max: 60}
	limit, err := parseListLimit(in.Params, "Limit", 60)
	if err != nil {
		return nil, err
	}

	result, err := store.ListDevicesPaginated(user.UserPoolID, user.ID, storecommon.ListOptions{
		MaxItems: limit,
		Marker:   in.PaginationToken,
	})
	if err != nil {
		return nil, ErrInternalError
	}

	formatted := make([]map[string]interface{}, 0, len(result.Items))
	for _, d := range result.Items {
		formatted = append(formatted, formatDevice(d))
	}

	resp := map[string]interface{}{"Devices": formatted}
	if result.IsTruncated && result.NextMarker != "" {
		resp["PaginationToken"] = result.NextMarker
	}
	return resp, nil
}

// updateDeviceStatusCore updates the remembered status of a device.
func (s *CognitoService) updateDeviceStatusCore(reqCtx *request.RequestContext, in UpdateDeviceStatusInput) (interface{}, error) {
	if in.AccessToken == "" || in.DeviceKey == "" {
		return nil, ErrInvalidParameter
	}
	if in.DeviceRememberedStatus != "remembered" && in.DeviceRememberedStatus != "not_remembered" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, in.AccessToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUserByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	device, err := store.GetDevice(user.UserPoolID, user.ID, in.DeviceKey)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	device.DeviceRememberedStatus = in.DeviceRememberedStatus
	device.DeviceLastModifiedDate = time.Now().UTC()
	if err := store.UpdateDevice(device); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// adminGetDeviceCore retrieves a device by its key (admin).
func (s *CognitoService) adminGetDeviceCore(reqCtx *request.RequestContext, in AdminGetDeviceInput) (interface{}, error) {
	if in.UserPoolID == "" || in.Username == "" || in.DeviceKey == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUser(in.UserPoolID, in.Username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	device, err := store.GetDevice(in.UserPoolID, user.ID, in.DeviceKey)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	return map[string]interface{}{"Device": formatDevice(device)}, nil
}

// adminForgetDeviceCore removes a device (admin).
func (s *CognitoService) adminForgetDeviceCore(reqCtx *request.RequestContext, in AdminForgetDeviceInput) (interface{}, error) {
	if in.UserPoolID == "" || in.Username == "" || in.DeviceKey == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUser(in.UserPoolID, in.Username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if _, err := store.GetDevice(in.UserPoolID, user.ID, in.DeviceKey); err != nil {
		return nil, ErrResourceNotFound
	}

	if err := store.DeleteDevice(in.UserPoolID, user.ID, in.DeviceKey); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// adminListDevicesCore lists devices for a user (admin).
func (s *CognitoService) adminListDevicesCore(reqCtx *request.RequestContext, in AdminListDevicesInput) (interface{}, error) {
	if in.UserPoolID == "" || in.Username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUser(in.UserPoolID, in.Username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Smithy QueryLimitType: range {min: 0, max: 60}
	limit, err := parseListLimit(in.Params, "Limit", 60)
	if err != nil {
		return nil, err
	}

	result, err := store.ListDevicesPaginated(in.UserPoolID, user.ID, storecommon.ListOptions{
		MaxItems: limit,
		Marker:   in.PaginationToken,
	})
	if err != nil {
		return nil, ErrInternalError
	}

	formatted := make([]map[string]interface{}, 0, len(result.Items))
	for _, d := range result.Items {
		formatted = append(formatted, formatDevice(d))
	}

	resp := map[string]interface{}{"Devices": formatted}
	if result.IsTruncated && result.NextMarker != "" {
		resp["PaginationToken"] = result.NextMarker
	}
	return resp, nil
}

// adminUpdateDeviceStatusCore updates the remembered status of a device
// (admin).
func (s *CognitoService) adminUpdateDeviceStatusCore(reqCtx *request.RequestContext, in AdminUpdateDeviceStatusInput) (interface{}, error) {
	if in.UserPoolID == "" || in.Username == "" || in.DeviceKey == "" {
		return nil, ErrInvalidParameter
	}
	if in.DeviceRememberedStatus != "remembered" && in.DeviceRememberedStatus != "not_remembered" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUser(in.UserPoolID, in.Username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	device, err := store.GetDevice(in.UserPoolID, user.ID, in.DeviceKey)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	device.DeviceRememberedStatus = in.DeviceRememberedStatus
	device.DeviceLastModifiedDate = time.Now().UTC()
	if err := store.UpdateDevice(device); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// parseDeviceAttributes extracts the DeviceAttributes list from the raw
// request parameter map.
func parseDeviceAttributes(params map[string]interface{}) map[string]string {
	result := make(map[string]string)
	if val, ok := params["DeviceAttributes"]; ok {
		if slice, ok := val.([]interface{}); ok {
			for _, v := range slice {
				if m, ok := v.(map[string]interface{}); ok {
					name, _ := m["Name"].(string)
					value, _ := m["Value"].(string)
					if name != "" {
						result[name] = value
					}
				}
			}
		}
	}
	return result
}
