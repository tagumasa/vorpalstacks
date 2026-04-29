// Package http provides HTTP server functionality for vorpalstacks.
package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"vorpalstacks/internal/common/defaults"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	"vorpalstacks/internal/server/authorization"
	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/server/http/chain"
	"vorpalstacks/internal/server/http/classifier"
	"vorpalstacks/internal/server/http/router"
	"vorpalstacks/internal/store/api"
	"vorpalstacks/internal/store/aws/iam"
	s3store "vorpalstacks/internal/store/aws/s3"
	stsstore "vorpalstacks/internal/store/aws/sts"
)

// NewServer creates and configures a new HTTP server with all required components.
func NewServer(cfg *Config) (*Server, error) {
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.DataPath == "" {
		cfg.DataPath = "./data"
	}

	storageMgr, err := storage.NewRegionStorageManager(&storage.Config{Path: cfg.DataPath})
	if err != nil {
		return nil, fmt.Errorf("failed to initialise storage manager: %w", err)
	}

	globalStore, err := storageMgr.GetGlobalStorage()
	if err != nil {
		storageMgr.Close()
		return nil, fmt.Errorf("failed to get global storage: %w", err)
	}

	serviceStore := api.NewServiceStore(globalStore)
	operationStore := api.NewOperationStore(globalStore)
	shapeStore := api.NewShapeStore(globalStore)
	memberStore := api.NewMemberStore(globalStore)
	configStore := api.NewConfigStore(globalStore)

	if err := configStore.LoadAll(); err != nil {
		storageMgr.Close()
		return nil, fmt.Errorf("failed to load service configs: %w", err)
	}

	resilienceConfig := &resilience.ServiceResilienceConfig{
		ServiceName: "aws-http-server",
		DefaultCircuitBreaker: &resilience.CircuitBreakerConfig{
			Name:             "default",
			MaxFailures:      resilience.DefaultCircuitBreakerMaxFailures,
			ResetTimeout:     resilience.DefaultCircuitBreakerResetTimeout,
			HalfOpenMaxCalls: resilience.DefaultCircuitBreakerHalfOpenRequests,
			SuccessThreshold: resilience.DefaultCircuitBreakerSuccessThreshold,
		},
		DefaultBulkhead: &resilience.BulkheadConfig{
			Name:          "default",
			MaxConcurrent: resilience.DefaultBulkheadMaxConcurrent,
			MaxWait:       resilience.DefaultBulkheadTimeout,
		},
		DefaultAdaptiveTimeout: &resilience.AdaptiveTimeoutConfig{
			Name:             "default",
			InitialTimeout:   resilience.DefaultAdaptiveTimeoutDefault,
			MinTimeout:       resilience.DefaultAdaptiveTimeoutMin,
			MaxTimeout:       resilience.DefaultAdaptiveTimeoutMax,
			SuccessThreshold: resilience.DefaultAdaptiveTimeoutSuccessThreshold,
			FailureThreshold: resilience.DefaultAdaptiveTimeoutFailureThreshold,
			AdjustmentFactor: resilience.DefaultAdaptiveAdjustmentFactor,
		},
	}

	// Create IAM store (always needed for services and authorization)
	iamStore := iam.NewIAMStore(globalStore, cfg.AccountID)

	// Create blob store for S3 and other services
	blobStore, err := storage.NewHybridBlobStore(globalStore, cfg.DataPath)
	if err != nil {
		storageMgr.Close()
		return nil, fmt.Errorf("failed to initialise blob store: %w", err)
	}

	// Create S3 store
	s3Store := s3store.NewS3Store(storageMgr, blobStore, cfg.AccountID)

	var authorizer dispatcher.Authorizer
	authEnabled := os.Getenv("AUTHORIZATION_ENABLED") == "true"
	if authEnabled {
		authorizer = authorization.NewAuthorizer(iamStore)
	}

	disp := dispatcher.NewDispatcher(
		serviceStore,
		operationStore,
		shapeStore,
		memberStore,
		configStore,
		resilienceConfig,
		storageMgr,
		iamStore,
		authorizer,
		cfg.AccountID,
	)

	serviceIndex, err := router.NewServiceIndex(serviceStore, operationStore)
	if err != nil {
		storageMgr.Close()
		return nil, fmt.Errorf("failed to build service index: %w", err)
	}
	classifier := classifier.NewClassifier(serviceIndex)

	region := cfg.Region
	if region == "" {
		region = defaults.DefaultRegion
	}

	// Create STS session store for temporary credential verification
	stsSessionStore := stsstore.NewSessionStore(globalStore, region)

	regionStorage, err := storageMgr.GetStorage(region)
	if err != nil {
		storageMgr.Close()
		return nil, fmt.Errorf("eventbus: get storage for region %s: %w", region, err)
	}

	pebbleStorage, ok := regionStorage.(*storage.PebbleStorage)
	if !ok {
		storageMgr.Close()
		return nil, fmt.Errorf("eventbus: storage is not PebbleStorage")
	}

	pebbleDB := pebbleStorage.DB()

	outbox := eventbus.NewPebbleOutboxStore(pebbleDB)

	registry := eventbus.NewEventRegistry()
	registry.Register("service:invoke", func() eventbus.Event { return &eventbus.ServiceInvokeRequest{} })
	registry.Register("sns:deliver", func() eventbus.Event { return &eventbus.SNSDeliveryEvent{} })
	registry.Register("events:deliver", func() eventbus.Event { return &eventbus.EventBridgeDeliveryEvent{} })
	registry.Register("scheduler:fired", func() eventbus.Event { return &eventbus.ScheduleFiredEvent{} })
	registry.Register("logs:deliver", func() eventbus.Event { return &eventbus.CloudWatchLogDeliveryEvent{} })
	registry.Register("s3:ObjectEvent", func() eventbus.Event { return &eventbus.S3ObjectEvent{} })
	registry.Register("cloudwatch:AlarmStateChange", func() eventbus.Event { return &eventbus.CloudWatchAlarmStateEvent{} })
	registry.Register("cognito:Trigger", func() eventbus.Event { return &eventbus.CognitoTriggerEvent{} })
	registry.Register("secretsmanager:Rotation", func() eventbus.Event { return &eventbus.SecretRotationEvent{} })
	registry.Register("dynamodb:Record", func() eventbus.Event { return &eventbus.DynamoDBRecordEvent{} })
	registry.Register("lambda:LogWrite", func() eventbus.Event { return &eventbus.LambdaLogWriteEvent{} })
	registry.Register("apigateway:AccessLog", func() eventbus.Event { return &eventbus.APIGatewayAccessLogEvent{} })
	registry.Register("logs:PutLogEvents", func() eventbus.Event { return &eventbus.CloudWatchLogsPutEvent{} })
	registry.Register("states:startExecution", func() eventbus.Event { return &eventbus.StepFunctionsStartExecutionEvent{} })
	registry.Register("events:putEvents", func() eventbus.Event { return &eventbus.EventBridgePutEventsEvent{} })

	bus := eventbus.NewEventBus(
		eventbus.WithOutbox(outbox),
		eventbus.WithEventRegistry(registry),
		eventbus.WithPolicyEvaluator(eventbus.NewSimplePolicyEvaluator()),
		eventbus.WithRoleResolver(eventbus.NewIAMRoleResolver(nil)),
	)

	shutdownCtx, shutdownCancel := context.WithCancel(context.Background())

	srv := &Server{
		config:          cfg,
		storageManager:  storageMgr,
		eventBus:        bus,
		dispatcher:      disp,
		classifier:      classifier,
		serviceStore:    serviceStore,
		operationStore:  operationStore,
		shutdownCtx:     shutdownCtx,
		shutdownCancel:  shutdownCancel,
		iamStore:        iamStore,
		stsSessionStore: stsSessionStore,
		s3Store:         s3Store,
		blobStore:       blobStore,
	}

	if cfg.TLSConfig.Enabled {
		tlsMgr, err := NewTLSManager(globalStore, cfg.TLSConfig.CertPath, cfg.TLSConfig.KeyPath, cfg.TLSConfig.Hostname)
		if err != nil {
			storageMgr.Close()
			return nil, fmt.Errorf("failed to initialise TLS manager: %w", err)
		}
		srv.tlsManager = tlsMgr
	}

	if cfg.UseChainGateway {
		srv.chainGateway = chain.NewGateway(classifier, disp)
	}

	return srv, nil
}

// Start starts the HTTP server.
//
// Returns:
//   - error: Any error that occurred during startup
func (s *Server) Start() error {
	var startErr error
	s.startOnce.Do(func() {
		r := chi.NewRouter()

		r.Use(middleware.Recoverer)
		r.Use(middleware.RequestID)
		r.Use(middleware.RealIP)

		const maxBodySize = 64 << 20
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
				next.ServeHTTP(w, r)
			})
		})

		if s.config.SignatureConfig.Enabled {
			r.Use(SignatureMiddleware(s.config.SignatureConfig, s.classifier, s.stsSessionStore))
		}

		r.Use(LoggingMiddleware)

		s.registerRoutes(r)

		h2s := &http2.Server{}

		s.httpServerMu.Lock()
		s.httpServer = &http.Server{
			Addr:              ":" + s.config.Port,
			Handler:           h2c.NewHandler(r, h2s),
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       15 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       120 * time.Second,
		}
		s.httpServerMu.Unlock()

		go s.handleShutdown()

		s.httpServerMu.Lock()
		srv := s.httpServer
		s.httpServerMu.Unlock()

		go func() {
			err := srv.ListenAndServe()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				if s.shutdownCancel != nil {
					s.shutdownCancel()
				}
				startErr = err
			}
		}()

		if s.config.TLSConfig.Enabled && s.tlsManager != nil {
			tlsConfig, err := s.tlsManager.GetTLSConfig()
			if err != nil {
				logs.Error("Failed to get TLS config", logs.Err(err))
				return
			}

			s.tlsServerMu.Lock()
			s.tlsServer = &http.Server{
				Addr:              ":" + s.config.TLSConfig.Port,
				Handler:           r,
				TLSConfig:         tlsConfig,
				ReadHeaderTimeout: 5 * time.Second,
				ReadTimeout:       15 * time.Second,
				WriteTimeout:      30 * time.Second,
				IdleTimeout:       120 * time.Second,
			}
			s.tlsServerMu.Unlock()

			go func() {
				s.tlsServerMu.Lock()
				tlsSrv := s.tlsServer
				s.tlsServerMu.Unlock()

				err := tlsSrv.ListenAndServeTLS("", "")
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					logs.Error("TLS server error", logs.Err(err))
				}
			}()
		}

		if s.listenerManager != nil {
			s.listenerManager.Start()
		}

		if s.eventBus != nil {
			if err := s.eventBus.Start(context.Background()); err != nil {
				logs.Error("Event bus start error", logs.Err(err))
			}
		}

		<-s.shutdownCtx.Done()
	})
	return startErr
}

func (s *Server) handleShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	select {
	case <-sigChan:
	case <-s.shutdownCtx.Done():
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	s.httpServerMu.Lock()
	srv := s.httpServer
	s.httpServerMu.Unlock()

	if srv != nil {
		if err := srv.Shutdown(ctx); err != nil {
			logs.Error("Server shutdown error", logs.Err(err))
		}
	}

	s.tlsServerMu.Lock()
	tlsSrv := s.tlsServer
	s.tlsServerMu.Unlock()

	if tlsSrv != nil {
		if err := tlsSrv.Shutdown(ctx); err != nil {
			logs.Error("TLS server shutdown error", logs.Err(err))
		}
	}

	if s.listenerManager != nil {
		s.listenerManager.Shutdown(ctx)
	}

	s.shutdownHooksMu.Lock()
	hooks := make([]ShutdownHook, len(s.shutdownHooks))
	copy(hooks, s.shutdownHooks)
	s.shutdownHooksMu.Unlock()

	for i := len(hooks) - 1; i >= 0; i-- {
		hooks[i](ctx)
	}

	if s.eventBus != nil {
		if err := s.eventBus.Shutdown(ctx); err != nil {
			logs.Error("Event bus shutdown error", logs.Err(err))
		}
	}

	s.closeOnce.Do(func() {
		if err := s.storageManager.Close(); err != nil {
			logs.Error("Storage manager close error", logs.Err(err))
		}
	})
}

// Close shuts down the HTTP server.
//
// Returns:
//   - error: Any error that occurred during shutdown
func (s *Server) Close() error {
	if s.shutdownCancel != nil {
		s.shutdownCancel()
	}
	s.httpServerMu.Lock()
	srv := s.httpServer
	s.httpServerMu.Unlock()
	if srv != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logs.Error("Server shutdown error", logs.Err(err))
		}
	}

	s.tlsServerMu.Lock()
	tlsSrv := s.tlsServer
	s.tlsServerMu.Unlock()
	if tlsSrv != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tlsSrv.Shutdown(ctx); err != nil {
			logs.Error("TLS server shutdown error", logs.Err(err))
		}
	}

	if s.listenerManager != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.listenerManager.Shutdown(ctx)
	}

	var err error
	s.closeOnce.Do(func() {
		err = s.storageManager.Close()
	})
	return err
}
