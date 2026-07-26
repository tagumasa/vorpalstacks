package cognitoidentityprovider

import (
	"context"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// ConfirmDevice registers a new device for the authenticated user.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ConfirmDevice.html
func (s *CognitoService) ConfirmDevice(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessToken := getAccessToken(req)
	devKey := req.GetParam("DeviceKey")
	if accessToken == "" || devKey == "" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, accessToken)
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
		DeviceKey:              devKey,
		UserPoolID:             user.UserPoolID,
		UserID:                 user.ID,
		DeviceName:             req.GetParam("DeviceName"),
		DeviceCreateDate:       now,
		DeviceLastModifiedDate: now,
		DeviceRememberedStatus: "remembered",
	}

	if svc, ok := req.Parameters["DeviceSecretVerifierConfig"]; ok {
		if m, ok := svc.(map[string]interface{}); ok {
			if sigVerifier, ok := m["PasswordVerifier"].(string); ok {
				device.DeviceSecretVerifierB = sigVerifier
			}
			if saltVerifier, ok := m["SaltVerifier"].(string); ok {
				device.DeviceSaltVerifier = saltVerifier
			}
		}
	}

	if attrs := parseDeviceAttributes(req); len(attrs) > 0 {
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

// GetDevice retrieves a device by its key.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_GetDevice.html
func (s *CognitoService) GetDevice(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessToken := getAccessToken(req)
	devKey := req.GetParam("DeviceKey")
	if accessToken == "" || devKey == "" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, accessToken)
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

	device, err := store.GetDevice(user.UserPoolID, user.ID, devKey)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	return map[string]interface{}{"Device": formatDevice(device)}, nil
}

// ForgetDevice removes a device for the authenticated user.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ForgetDevice.html
func (s *CognitoService) ForgetDevice(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessToken := getAccessToken(req)
	devKey := req.GetParam("DeviceKey")
	if accessToken == "" || devKey == "" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, accessToken)
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

	if _, err := store.GetDevice(user.UserPoolID, user.ID, devKey); err != nil {
		return nil, ErrResourceNotFound
	}

	if err := store.DeleteDevice(user.UserPoolID, user.ID, devKey); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// ListDevices lists devices for the authenticated user.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ListDevices.html
func (s *CognitoService) ListDevices(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessToken := getAccessToken(req)
	if accessToken == "" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, accessToken)
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
	limit := 60
	if l := request.GetIntParam(req.Parameters, "Limit"); l > 0 && l <= 60 {
		limit = l
	}

	result, err := store.ListDevicesPaginated(user.UserPoolID, user.ID, storecommon.ListOptions{
		MaxItems: limit,
		Marker:   req.GetParam("PaginationToken"),
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

// UpdateDeviceStatus updates the remembered status of a device.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_UpdateDeviceStatus.html
func (s *CognitoService) UpdateDeviceStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessToken := getAccessToken(req)
	devKey := req.GetParam("DeviceKey")
	devStatus := req.GetParam("DeviceRememberedStatus")
	if accessToken == "" || devKey == "" {
		return nil, ErrInvalidParameter
	}
	if devStatus != "remembered" && devStatus != "not_remembered" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, accessToken)
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

	device, err := store.GetDevice(user.UserPoolID, user.ID, devKey)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	device.DeviceRememberedStatus = devStatus
	device.DeviceLastModifiedDate = time.Now().UTC()
	if err := store.UpdateDevice(device); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// AdminGetDevice retrieves a device by its key (admin).
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminGetDevice.html
func (s *CognitoService) AdminGetDevice(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	username := getUsername(req)
	devKey := req.GetParam("DeviceKey")
	if userPoolID == "" || username == "" || devKey == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUser(userPoolID, username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	device, err := store.GetDevice(userPoolID, user.ID, devKey)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	return map[string]interface{}{"Device": formatDevice(device)}, nil
}

// AdminForgetDevice removes a device (admin).
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminForgetDevice.html
func (s *CognitoService) AdminForgetDevice(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	username := getUsername(req)
	devKey := req.GetParam("DeviceKey")
	if userPoolID == "" || username == "" || devKey == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUser(userPoolID, username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if _, err := store.GetDevice(userPoolID, user.ID, devKey); err != nil {
		return nil, ErrResourceNotFound
	}

	if err := store.DeleteDevice(userPoolID, user.ID, devKey); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// AdminListDevices lists devices for a user (admin).
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminListDevices.html
func (s *CognitoService) AdminListDevices(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	username := getUsername(req)
	if userPoolID == "" || username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUser(userPoolID, username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Smithy QueryLimitType: range {min: 0, max: 60}
	limit := 60
	if l := request.GetIntParam(req.Parameters, "Limit"); l > 0 && l <= 60 {
		limit = l
	}

	result, err := store.ListDevicesPaginated(userPoolID, user.ID, storecommon.ListOptions{
		MaxItems: limit,
		Marker:   req.GetParam("PaginationToken"),
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

// AdminUpdateDeviceStatus updates the remembered status of a device (admin).
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminUpdateDeviceStatus.html
func (s *CognitoService) AdminUpdateDeviceStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := req.GetParam("UserPoolId")
	username := getUsername(req)
	devKey := req.GetParam("DeviceKey")
	devStatus := req.GetParam("DeviceRememberedStatus")
	if userPoolID == "" || username == "" || devKey == "" {
		return nil, ErrInvalidParameter
	}
	if devStatus != "remembered" && devStatus != "not_remembered" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUser(userPoolID, username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	device, err := store.GetDevice(userPoolID, user.ID, devKey)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	device.DeviceRememberedStatus = devStatus
	device.DeviceLastModifiedDate = time.Now().UTC()
	if err := store.UpdateDevice(device); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

func formatDevice(d *cognitostore.Device) map[string]interface{} {
	result := map[string]interface{}{
		"DeviceKey":              d.DeviceKey,
		"DeviceCreateDate":       d.DeviceCreateDate.Unix(),
		"DeviceLastModifiedDate": d.DeviceLastModifiedDate.Unix(),
	}
	if d.DeviceName != "" {
		result["DeviceName"] = d.DeviceName
	}
	if !d.DeviceLastAuthenticatedDate.IsZero() {
		result["DeviceLastAuthenticatedDate"] = d.DeviceLastAuthenticatedDate.Unix()
	}
	if len(d.DeviceAttributes) > 0 {
		attrs := make([]map[string]string, 0, len(d.DeviceAttributes))
		for k, v := range d.DeviceAttributes {
			attrs = append(attrs, map[string]string{"Name": k, "Value": v})
		}
		result["DeviceAttributes"] = attrs
	}
	return result
}

func parseDeviceAttributes(req *request.ParsedRequest) map[string]string {
	result := make(map[string]string)
	if val, ok := req.Parameters["DeviceAttributes"]; ok {
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
