package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// ConfirmDevice registers a new device for the authenticated user.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ConfirmDevice.html
func (s *CognitoService) ConfirmDevice(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.confirmDeviceCore(reqCtx, ConfirmDeviceInput{
		AccessToken: getAccessToken(req),
		DeviceKey:   req.GetParam("DeviceKey"),
		DeviceName:  req.GetParam("DeviceName"),
		Params:      req.Parameters,
	})
}

// GetDevice retrieves a device by its key.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_GetDevice.html
func (s *CognitoService) GetDevice(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.getDeviceCore(reqCtx, GetDeviceInput{
		AccessToken: getAccessToken(req),
		DeviceKey:   req.GetParam("DeviceKey"),
	})
}

// ForgetDevice removes a device for the authenticated user.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ForgetDevice.html
func (s *CognitoService) ForgetDevice(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.forgetDeviceCore(reqCtx, ForgetDeviceInput{
		AccessToken: getAccessToken(req),
		DeviceKey:   req.GetParam("DeviceKey"),
	})
}

// ListDevices lists devices for the authenticated user.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ListDevices.html
func (s *CognitoService) ListDevices(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.listDevicesCore(reqCtx, ListDevicesInput{
		AccessToken:     getAccessToken(req),
		PaginationToken: req.GetParam("PaginationToken"),
		Params:          req.Parameters,
	})
}

// UpdateDeviceStatus updates the remembered status of a device.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_UpdateDeviceStatus.html
func (s *CognitoService) UpdateDeviceStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.updateDeviceStatusCore(reqCtx, UpdateDeviceStatusInput{
		AccessToken:            getAccessToken(req),
		DeviceKey:              req.GetParam("DeviceKey"),
		DeviceRememberedStatus: req.GetParam("DeviceRememberedStatus"),
	})
}

// AdminGetDevice retrieves a device by its key (admin).
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminGetDevice.html
func (s *CognitoService) AdminGetDevice(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.adminGetDeviceCore(reqCtx, AdminGetDeviceInput{
		UserPoolID: getUserPoolID(req),
		Username:   getUsername(req),
		DeviceKey:  req.GetParam("DeviceKey"),
	})
}

// AdminForgetDevice removes a device (admin).
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminForgetDevice.html
func (s *CognitoService) AdminForgetDevice(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.adminForgetDeviceCore(reqCtx, AdminForgetDeviceInput{
		UserPoolID: getUserPoolID(req),
		Username:   getUsername(req),
		DeviceKey:  req.GetParam("DeviceKey"),
	})
}

// AdminListDevices lists devices for a user (admin).
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminListDevices.html
func (s *CognitoService) AdminListDevices(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.adminListDevicesCore(reqCtx, AdminListDevicesInput{
		UserPoolID:      getUserPoolID(req),
		Username:        getUsername(req),
		PaginationToken: req.GetParam("PaginationToken"),
		Params:          req.Parameters,
	})
}

// AdminUpdateDeviceStatus updates the remembered status of a device (admin).
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminUpdateDeviceStatus.html
func (s *CognitoService) AdminUpdateDeviceStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.adminUpdateDeviceStatusCore(reqCtx, AdminUpdateDeviceStatusInput{
		UserPoolID:             getUserPoolID(req),
		Username:               getUsername(req),
		DeviceKey:              req.GetParam("DeviceKey"),
		DeviceRememberedStatus: req.GetParam("DeviceRememberedStatus"),
	})
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
