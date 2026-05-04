package grpcweb

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"

	"connectrpc.com/connect"
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
//
// All handlers except admin auth are protected by JWT authentication.
// The admin auth service is intentionally excluded because clients must be
// able to log in before they possess a valid token.
func RegisterAdminHandlers(s *Server, st storage.BasicStorage, accountID, region, dataPath string, handlers []HandlerRegistration, shutdownFunc func()) {
	configStore := config.NewStore(st)
	version := buildVersion()
	adminConfigService := svcadminconfig.NewAdminConfigService(configStore, shutdownFunc, dataPath, version)

	adminAuthKey, err := loadOrGenerateAdminAuthKey(dataPath)
	if err != nil {
		// Without a JWT key we cannot enforce authentication.  Register all
		// handlers without protection so the server remains functional.
		fmt.Fprintf(os.Stderr, "Warning: Failed to initialise admin auth JWT key: %v\n", err)
		path, handler := adminconfigconnect.NewAdminConfigServiceHandler(adminConfigService)
		s.Handle(path, handler)
		for _, h := range handlers {
			s.Handle(h.Path, h.Handler)
		}
		return
	}

	jwtManager := vsjwt.NewManager(adminAuthKey, "admin-auth-key", "vorpalstacks/admin-auth")
	authInterceptor := NewAuthInterceptor(jwtManager)

	// Admin config service — protected by the ConnectRPC auth interceptor.
	path, handler := adminconfigconnect.NewAdminConfigServiceHandler(
		adminConfigService,
		connect.WithInterceptors(authInterceptor),
	)
	s.Handle(path, handler)

	// External service handlers — protected by HTTP auth middleware.
	for _, h := range handlers {
		s.Handle(h.Path, newAuthHTTPMiddleware(jwtManager, h.Handler))
	}

	// Admin auth service — intentionally excluded from auth checks.
	adminAuthService := svcadminauth.NewAdminAuthService(st, jwtManager, accountID)
	s.Handle(adminauthconnect.NewAdminAuthServiceHandler(adminAuthService))
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

func buildVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, s := range info.Settings {
			if s.Key == "vcs.revision" {
				return info.GoVersion + " " + s.Value[:12]
			}
		}
		return info.GoVersion
	}
	return "unknown"
}
