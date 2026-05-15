// Package apps manages the application lifecycle: service creation, cross-service
// wiring, listener registration, and graceful shutdown.
package apps

import (
	"context"
	"embed"
	"fmt"
	"os"
	"sync"
	"time"

	"vorpalstacks/internal/common/serviceports"
	appconfig "vorpalstacks/internal/config"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	chihttp "vorpalstacks/internal/server/http"
	"vorpalstacks/internal/server/listener"
	"vorpalstacks/internal/server/portalloc"
)

// Config holds all configuration for service initialisation and wiring.
type Config struct {
	Port                  int
	GRPCWebPort           int
	GRPCWebBindAddr       string
	BindAddr              string
	DataPath              string
	AccountID             string
	Region                string
	AccessKeyID           string
	SecretAccessKey       string
	SignatureVerification bool
	UseChainGateway       bool
	TLSEnabled            bool
	TLSPort               int
	TLSCertPath           string
	TLSKeyPath            string
	TLSHostname           string
	Route53DNSEnabled     bool
	DockerHost            string
	ConsoleAssets         embed.FS

	ACM             bool
	APIGateway      bool
	Athena          bool
	AppSync         bool
	CloudFront      bool
	CloudTrail      bool
	CloudWatch      bool
	CloudWatchLogs  bool
	Cognito         bool
	CognitoIdentity bool
	DynamoDB        bool
	EC2             bool
	EventBridge     bool
	Kinesis         bool
	KMS             bool
	Lambda          bool
	Neptune         bool
	NeptuneData     bool
	NeptuneGraph    bool
	Route53         bool
	S3              bool
	Scheduler       bool
	SecretsManager  bool
	SESv2           bool
	SFN             bool
	SNS             bool
	SQS             bool
	SSM             bool
	STS             bool
	IAM             bool
	TimestreamQuery bool
	TimestreamWrite bool
	WAFv2           bool
	RDSMySQL        bool
}

// FromBootstrap converts a BootstrapConfig into an apps.Config.
func FromBootstrap(bc *appconfig.BootstrapConfig) *Config {
	bindAddr := portalloc.ResolveBindAddr(appconfig.GetStore())
	return &Config{
		Port:                  bc.Port,
		GRPCWebPort:           bc.GRPCWebPort,
		GRPCWebBindAddr:       bindAddr,
		BindAddr:              bindAddr,
		DataPath:              bc.DataPath,
		AccountID:             bc.AccountID,
		Region:                bc.Region,
		AccessKeyID:           bc.AccessKeyID,
		SecretAccessKey:       bc.SecretAccessKey,
		SignatureVerification: bc.SignatureVerification,
		UseChainGateway:       bc.UseChainGateway,
		TLSEnabled:            bc.TLSEnabled,
		TLSPort:               bc.TLSPort,
		TLSCertPath:           bc.TLSCertPath,
		TLSKeyPath:            bc.TLSKeyPath,
		TLSHostname:           bc.TLSHostname,
		Route53DNSEnabled:     bc.Route53DNSEnabled,
		DockerHost:            bc.DockerHost,
		ACM:                   bc.ACM,
		APIGateway:            bc.APIGateway,
		Athena:                bc.Athena,
		AppSync:               bc.AppSync,
		CloudFront:            bc.CloudFront,
		CloudTrail:            bc.CloudTrail,
		CloudWatch:            bc.CloudWatch,
		CloudWatchLogs:        bc.Logs,
		Cognito:               bc.Cognito,
		CognitoIdentity:       bc.CognitoIdentity,
		DynamoDB:              bc.DynamoDB,
		EC2:                   bc.EC2,
		RDSMySQL:              bc.RDSMySQL,
		EventBridge:           bc.Events,
		Kinesis:               bc.Kinesis,
		KMS:                   bc.KMS,
		Lambda:                bc.Lambda,
		Neptune:               bc.Neptune,
		NeptuneData:           bc.NeptuneData,
		NeptuneGraph:          bc.NeptuneGraph,
		Route53:               bc.Route53,
		S3:                    bc.S3,
		Scheduler:             bc.Scheduler,
		SecretsManager:        bc.SecretsManager,
		SESv2:                 bc.SESv2,
		SFN:                   bc.StepFunctions,
		SNS:                   bc.SNS,
		SQS:                   bc.SQS,
		SSM:                   bc.SSM,
		STS:                   bc.STS,
		IAM:                   bc.IAM,
		TimestreamQuery:       bc.TimestreamQuery,
		TimestreamWrite:       bc.TimestreamWrite,
		WAFv2:                 bc.WAFv2,
	}
}

// PrintStartupBanner writes the server startup information to stdout.
func (c *Config) PrintStartupBanner() {
	fmt.Printf("Starting AWS-compatible server on %s:%d\n", c.BindAddr, c.Port)
	fmt.Printf("Data path: %s\n", c.DataPath)
	fmt.Printf("Account ID: %s\n", c.AccountID)
	if c.UseChainGateway {
		fmt.Printf("Gateway mode: Chain (metadata-based routing)\n")
	} else {
		fmt.Printf("Gateway mode: Standard (hardcoded routing)\n")
	}
	services := []struct {
		name string
		on   bool
	}{
		{"ACM", c.ACM}, {"APIGateway", c.APIGateway}, {"Athena", c.Athena},
		{"AppSync", c.AppSync}, {"CloudFront", c.CloudFront},
		{"CloudTrail", c.CloudTrail}, {"CloudWatch", c.CloudWatch},
		{"CloudWatchLogs", c.CloudWatchLogs}, {"Cognito", c.Cognito},
		{"CognitoIdentity", c.CognitoIdentity}, {"DynamoDB", c.DynamoDB},
		{"EventBridge", c.EventBridge}, {"IAM", c.IAM}, {"Kinesis", c.Kinesis},
		{"KMS", c.KMS}, {"Lambda", c.Lambda}, {"Neptune", c.Neptune},
		{"NeptuneData", c.NeptuneData}, {"NeptuneGraph", c.NeptuneGraph},
		{"Route53", c.Route53}, {"S3", c.S3}, {"Scheduler", c.Scheduler},
		{"SecretsManager", c.SecretsManager}, {"SESv2", c.SESv2},
		{"SFN", c.SFN}, {"SNS", c.SNS}, {"SQS", c.SQS}, {"SSM", c.SSM},
		{"STS", c.STS}, {"TimestreamQuery", c.TimestreamQuery},
		{"TimestreamWrite", c.TimestreamWrite}, {"WAFv2", c.WAFv2},
		{"RDSMySQL", c.RDSMySQL},
	}
	fmt.Print("Services:")
	for _, s := range services {
		if s.on {
			fmt.Printf(" %s", s.name)
		}
	}
	fmt.Println()
	if c.Route53DNSEnabled {
		fmt.Printf("Route53 DNS server: enabled (port from config)\n")
	}
	if c.SignatureVerification {
		fmt.Printf("Signature verification: enabled\n")
	} else {
		fmt.Printf("Signature verification: disabled\n")
	}
	if c.AccessKeyID == "" {
		fmt.Println("Warning: AWS_ACCESS_KEY_ID not set")
	}
	if c.SecretAccessKey == "" {
		fmt.Println("Warning: AWS_SECRET_ACCESS_KEY not set")
	}
}

// ServerHost returns the server hostname suitable for self-referencing URLs.
func (c *Config) ServerHost() string {
	addr := c.BindAddr
	if addr == "" {
		addr = "127.0.0.1"
	}
	return fmt.Sprintf("%s:%d", addr, c.Port)
}

// shutdownEntry pairs a service name with its shutdown function.
type shutdownEntry struct {
	name string
	fn   func(context.Context) error
}

// App owns the HTTP server, gRPC-Web admin server, all services, and their
// lifecycle. It is created by New and started by Run.
type App struct {
	cfg           *Config
	server        *chihttp.Server
	grpcWeb       grpcWebStarter
	lm            *listener.Manager
	state         *serviceState
	consoleAssets embed.FS

	mu        sync.Mutex
	shutdowns []shutdownEntry
}

// grpcWebStarter abstracts the Start method of *grpcweb.Server.
type grpcWebStarter interface {
	Start() error
	Shutdown(ctx context.Context) error
}

// New creates the App, initialising the HTTP server and all services.
func New(cfg *Config) (*App, error) {
	serverCfg := &chihttp.Config{
		Port:      cfg.Port,
		BindAddr:  cfg.BindAddr,
		DataPath:  cfg.DataPath,
		AccountID: cfg.AccountID,
		SignatureConfig: chihttp.SignatureConfig{
			Enabled:         cfg.SignatureVerification,
			Region:          cfg.Region,
			AccessKeyID:     cfg.AccessKeyID,
			SecretAccessKey: cfg.SecretAccessKey,
		},
		UseChainGateway: cfg.UseChainGateway,
		TLSConfig: chihttp.TLSConfig{
			Enabled:  cfg.TLSEnabled,
			Port:     cfg.TLSPort,
			CertPath: cfg.TLSCertPath,
			KeyPath:  cfg.TLSKeyPath,
			Hostname: cfg.TLSHostname,
		},
	}

	srv, err := chihttp.NewServer(serverCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create server: %w", err)
	}

	appconfig.Initialise(srv.Storage())

	a := &App{
		cfg:           cfg,
		server:        srv,
		consoleAssets: cfg.ConsoleAssets,
	}

	mainPort := cfg.Port
	if mainPort == 0 {
		mainPort = serviceports.HTTP
	}
	a.lm = listener.NewManager(mainPort)
	srv.SetListenerManager(a.lm)

	if err := a.initAlwaysOnServices(); err != nil {
		return nil, err
	}
	a.wireCrossServiceDeps()
	if err := a.initOptionalServices(); err != nil {
		a.Shutdown(context.Background())
		return nil, err
	}
	a.initGRPCWebAdmin()
	a.initEventBusPolicies()

	a.registerListeners()

	return a, nil
}

// NewWithStorage creates the App using a pre-created storage manager.
// The caller must have already called config.Initialise() with the global
// storage. This is the production entry point used by main.go.
func NewWithStorage(cfg *Config, sm *storage.RegionStorageManager) (*App, error) {
	serverCfg := &chihttp.Config{
		Port:      cfg.Port,
		BindAddr:  cfg.BindAddr,
		DataPath:  cfg.DataPath,
		AccountID: cfg.AccountID,
		SignatureConfig: chihttp.SignatureConfig{
			Enabled:         cfg.SignatureVerification,
			Region:          cfg.Region,
			AccessKeyID:     cfg.AccessKeyID,
			SecretAccessKey: cfg.SecretAccessKey,
		},
		UseChainGateway: cfg.UseChainGateway,
		TLSConfig: chihttp.TLSConfig{
			Enabled:  cfg.TLSEnabled,
			Port:     cfg.TLSPort,
			CertPath: cfg.TLSCertPath,
			KeyPath:  cfg.TLSKeyPath,
			Hostname: cfg.TLSHostname,
		},
		StorageManager: sm,
	}

	srv, err := chihttp.NewServer(serverCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create server: %w", err)
	}

	a := &App{
		cfg:           cfg,
		server:        srv,
		consoleAssets: cfg.ConsoleAssets,
	}

	mainPort := cfg.Port
	if mainPort == 0 {
		mainPort = serviceports.HTTP
	}
	a.lm = listener.NewManager(mainPort)
	srv.SetListenerManager(a.lm)

	if err := a.initAlwaysOnServices(); err != nil {
		return nil, err
	}
	a.wireCrossServiceDeps()
	if err := a.initOptionalServices(); err != nil {
		a.Shutdown(context.Background())
		return nil, err
	}
	a.initGRPCWebAdmin()
	a.initEventBusPolicies()

	a.registerListeners()

	return a, nil
}

// Run starts the HTTP server, gRPC-Web server, and prints the startup banner.
// It blocks until the server shuts down.
func (a *App) Run() error {
	a.cfg.PrintStartupBanner()

	go func() {
		defer func() { recover() }()
		fmt.Printf("Starting gRPC-Web admin server on %s:%d\n", a.cfg.GRPCWebBindAddr, a.cfg.GRPCWebPort)
		if err := a.grpcWeb.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "gRPC-Web server error: %v\n", err)
		}
	}()

	err := a.server.Start()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	a.Shutdown(ctx)

	if a.grpcWeb != nil {
		a.grpcWeb.Shutdown(ctx)
	}

	return err
}

// Shutdown performs a graceful shutdown of all services, the gRPC-Web server,
// listeners, event bus, and storage in LIFO order.
func (a *App) Shutdown(ctx context.Context) error {
	a.mu.Lock()
	entries := make([]shutdownEntry, len(a.shutdowns))
	copy(entries, a.shutdowns)
	a.mu.Unlock()

	for i := len(entries) - 1; i >= 0; i-- {
		if err := entries[i].fn(ctx); err != nil {
			logs.Error("Shutdown error",
				logs.String("service", entries[i].name),
				logs.Err(err),
			)
		}
	}

	if a.lm != nil {
		a.lm.Shutdown(ctx)
	}

	return nil
}

// addShutdown registers a shutdown function for a named service.
func (a *App) addShutdown(name string, fn func(context.Context) error) {
	a.mu.Lock()
	a.shutdowns = append(a.shutdowns, shutdownEntry{name: name, fn: fn})
	a.mu.Unlock()
}
