package grpcweb

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/config"
	"vorpalstacks/pkg/vsjwt"

	adminauthconnect "vorpalstacks/internal/pb/aws/admin_auth/admin_authconnect"
	adminconfigconnect "vorpalstacks/internal/pb/aws/admin_config/admin_configconnect"
	svcadminauth "vorpalstacks/internal/services/aws/admin_auth"
	svcadminconfig "vorpalstacks/internal/services/aws/admin_config"
)

// HandlerRegistration pairs a Connect RPC path with its HTTP handler.
type HandlerRegistration struct {
	Path    string
	Handler http.Handler
}

// RegisterAdminHandlers registers Connect RPC handlers for the admin console.
// Service-specific handlers are provided by the caller; only admin config and
// admin auth are created internally (they are grpcweb-owned concerns).
func RegisterAdminHandlers(s *Server, st storage.BasicStorage, accountID, region, dataPath string, handlers []HandlerRegistration) {
	configStore := config.NewStore(st)
	adminConfigService := svcadminconfig.NewService(configStore)
	path, handler := adminconfigconnect.NewAdminConfigServiceHandler(adminConfigService)
	s.Handle(path, handler)

	for _, h := range handlers {
		s.Handle(h.Path, h.Handler)
	}

	adminAuthKey, err := loadOrGenerateAdminAuthKey(dataPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to initialise admin auth JWT key: %v\n", err)
	} else {
		jwtManager := vsjwt.NewManager(adminAuthKey, "admin-auth-key", "vorpalstacks/admin-auth")
		adminAuthService := svcadminauth.NewService(st, jwtManager, accountID)
		s.Handle(adminauthconnect.NewAdminAuthServiceHandler(adminAuthService))
	}
}

func loadOrGenerateAdminAuthKey(dataPath string) (*rsa.PrivateKey, error) {
	keyDir := filepath.Join(dataPath, "admin_auth")
	keyFile := filepath.Join(keyDir, "jwt_key.pem")

	if pemData, err := os.ReadFile(keyFile); err == nil {
		key, err := vsjwt.DecodePrivateKeyFromPEM(string(pemData))
		if err == nil {
			return key, nil
		}
	}

	key, err := vsjwt.GenerateRSAKeyPair()
	if err != nil {
		return nil, fmt.Errorf("generate RSA key pair: %w", err)
	}

	if err := os.MkdirAll(keyDir, 0700); err != nil {
		return nil, fmt.Errorf("create key directory: %w", err)
	}

	if err := os.WriteFile(keyFile, []byte(vsjwt.EncodePrivateKeyToPEM(key)), 0600); err != nil {
		return nil, fmt.Errorf("write key file: %w", err)
	}

	return key, nil
}
