package cognitoidentityprovider

import (
	"vorpalstacks/internal/store/aws/common"
)

// CreateDevice stores a new device record.
func (s *CognitoStore) CreateDevice(d *Device) error {
	return s.devicesStore.Put(deviceKey(d.UserPoolID, d.UserID, d.DeviceKey), d)
}

// GetDevice retrieves a device by its key.
func (s *CognitoStore) GetDevice(userPoolID, userID, devKey string) (*Device, error) {
	var d Device
	if err := s.devicesStore.Get(deviceKey(userPoolID, userID, devKey), &d); err != nil {
		return nil, err
	}
	return &d, nil
}

// UpdateDevice updates an existing device record.
func (s *CognitoStore) UpdateDevice(d *Device) error {
	return s.devicesStore.Put(deviceKey(d.UserPoolID, d.UserID, d.DeviceKey), d)
}

// DeleteDevice removes a device record.
func (s *CognitoStore) DeleteDevice(userPoolID, userID, devKey string) error {
	return s.devicesStore.Delete(deviceKey(userPoolID, userID, devKey))
}

// ListDevicesPaginated lists devices for a user with server-side pagination.
func (s *CognitoStore) ListDevicesPaginated(userPoolID, userID string, opts common.ListOptions) (*common.ListResult[Device], error) {
	opts.Prefix = devicePrefix(userPoolID, userID)
	return common.List[Device](s.devicesStore, opts, nil)
}

// CreateAuthEvent stores a new authentication event.
func (s *CognitoStore) CreateAuthEvent(e *AuthEvent) error {
	return s.authEventsStore.Put(authEventKey(e.UserPoolID, e.UserID, e.EventID), e)
}

// GetAuthEvent retrieves an auth event by its ID.
func (s *CognitoStore) GetAuthEvent(userPoolID, userID, eventID string) (*AuthEvent, error) {
	var e AuthEvent
	if err := s.authEventsStore.Get(authEventKey(userPoolID, userID, eventID), &e); err != nil {
		return nil, err
	}
	return &e, nil
}

// UpdateAuthEvent updates an existing auth event record.
func (s *CognitoStore) UpdateAuthEvent(e *AuthEvent) error {
	return s.authEventsStore.Put(authEventKey(e.UserPoolID, e.UserID, e.EventID), e)
}

// ListAuthEventsPaginated lists auth events for a user with server-side pagination.
func (s *CognitoStore) ListAuthEventsPaginated(userPoolID, userID string, opts common.ListOptions) (*common.ListResult[AuthEvent], error) {
	opts.Prefix = authEventPrefix(userPoolID, userID)
	return common.List[AuthEvent](s.authEventsStore, opts, nil)
}

// ===================== WebAuthn =====================

func (s *CognitoStore) CreateWebAuthnCredential(c *WebAuthnCredential) error {
	return s.webauthnStore.Put(webauthnKey(c.UserPoolID, c.UserID, c.CredentialID), c)
}

func (s *CognitoStore) GetWebAuthnCredential(userPoolID, userID, credID string) (*WebAuthnCredential, error) {
	var c WebAuthnCredential
	if err := s.webauthnStore.Get(webauthnKey(userPoolID, userID, credID), &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *CognitoStore) DeleteWebAuthnCredential(userPoolID, userID, credID string) error {
	return s.webauthnStore.Delete(webauthnKey(userPoolID, userID, credID))
}

func (s *CognitoStore) ListWebAuthnCredentialsPaginated(userPoolID, userID string, opts common.ListOptions) (*common.ListResult[WebAuthnCredential], error) {
	opts.Prefix = webauthnPrefix(userPoolID, userID)
	return common.List[WebAuthnCredential](s.webauthnStore, opts, nil)
}
