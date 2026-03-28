// Package http provides HTTP server functionality for vorpalstacks.
package http

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"vorpalstacks/internal/core/resilience"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/server/http/chain"
	"vorpalstacks/internal/server/http/router"
	"vorpalstacks/internal/services/aws/common/interfaces"
	"vorpalstacks/internal/store/api"
	acmstore "vorpalstacks/internal/store/aws/acm"
	apigatewaystore "vorpalstacks/internal/store/aws/apigateway"
	athenastore "vorpalstacks/internal/store/aws/athena"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
	cloudwatchstore "vorpalstacks/internal/store/aws/cloudwatch"
	cloudwatchlogsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	cognitoidentitystore "vorpalstacks/internal/store/aws/cognitoidentity"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	dynamodbstore "vorpalstacks/internal/store/aws/dynamodb"
	eventbridgestore "vorpalstacks/internal/store/aws/eventbridge"
	"vorpalstacks/internal/store/aws/iam"
	kmsstore "vorpalstacks/internal/store/aws/kms"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
	route53store "vorpalstacks/internal/store/aws/route53"
	s3store "vorpalstacks/internal/store/aws/s3"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
	sesv2store "vorpalstacks/internal/store/aws/sesv2"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
	snsstore "vorpalstacks/internal/store/aws/sns"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
	ssmstore "vorpalstacks/internal/store/aws/ssm"
	stsstore "vorpalstacks/internal/store/aws/sts"
	timestreamstore "vorpalstacks/internal/store/aws/timestream"
	wafstore "vorpalstacks/internal/store/aws/waf"
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

	// Create DynamoDB store (lazy initialization per region in service layer)
	// Note: DynamoDB store needs TransactionalStorageWith2PC which is obtained per-region
	// We pass nil here and let the service create it per-request for now
	// This will be refactored to use a factory pattern similar to S3
	var acmStore acmstore.CertificateStoreInterface
	var ssmStore ssmstore.SSMStoreInterface
	var stsStore stsstore.SessionStoreInterface
	var dynamoDBStore dynamodbstore.DynamoDBStoreInterface
	var lambdaStore lambdastore.LambdaStoreInterface
	var eventBridgeStore eventbridgestore.EventsStoreInterface
	var cloudWatchAlarmStore cloudwatchstore.AlarmStoreInterface
	var cloudWatchMetricStore cloudwatchstore.MetricChunkStoreInterface
	var cloudWatchLogsStore cloudwatchlogsstore.CloudWatchLogsStoreInterface
	var apiGatewayStore apigatewaystore.APIGatewayStoresInterface
	var sfnStore sfnstore.StepFunctionStoreInterface
	var athenaStore athenastore.AthenaStoresInterface
	var cognitoStore cognitostore.CognitoStoreInterface
	var timestreamStores timestreamstore.TimestreamStoresInterface
	var wafStores wafstore.WAFStoresInterface
	var kmsStores kmsstore.KMSStoresInterface
	var sesv2Store sesv2store.SESv2StoreInterface
	var schedulerStore schedulerstore.SchedulerStoreInterface
	var route53Stores route53store.Route53StoresInterface
	var cognitoIdentityStore cognitoidentitystore.CognitoIdentityStoreInterface
	var cloudFrontStores cloudfrontstore.CloudFrontStoresInterface
	var cloudTrailStore cloudtrailstore.CloudTrailStoreInterface
	var secretsManagerStore secretsmanagerstore.SecretStoreInterface
	var snsStore snsstore.SNSStoreInterface
	var sqsStore sqsstore.SQSStoreInterface

	storeProvider := interfaces.NewStoreProvider(
		iamStore,
		acmStore,
		s3Store,
		ssmStore,
		stsStore,
		dynamoDBStore,
		lambdaStore,
		eventBridgeStore,
		cloudWatchAlarmStore,
		cloudWatchMetricStore,
		cloudWatchLogsStore,
		apiGatewayStore,
		sfnStore,
		athenaStore,
		cognitoStore,
		timestreamStores,
		wafStores,
		kmsStores,
		sesv2Store,
		schedulerStore,
		route53Stores,
		cognitoIdentityStore,
		cloudFrontStores,
		cloudTrailStore,
		secretsManagerStore,
		snsStore,
		sqsStore,
	)

	var authorizer *dispatcher.Authorizer
	authEnabled := os.Getenv("AUTHORIZATION_ENABLED") == "true"
	if authEnabled {
		authorizer = dispatcher.NewAuthorizer(iamStore)
	}

	disp := dispatcher.NewDispatcher(
		serviceStore,
		operationStore,
		shapeStore,
		memberStore,
		configStore,
		resilienceConfig,
		storageMgr,
		storeProvider,
		authorizer,
		cfg.AccountID,
	)

	shapeTraitStore := api.NewShapeTraitStore(globalStore)
	memberTraitStore := api.NewMemberTraitStore(globalStore)
	operationErrorStore := api.NewOperationErrorStore(globalStore)
	enumValueStore := api.NewEnumValueStore(globalStore)

	if cfg.MetadataPath != "" {
		if err := loadMetadata(cfg.MetadataPath, serviceStore, operationStore, shapeStore, memberStore, shapeTraitStore, memberTraitStore, operationErrorStore, enumValueStore, configStore); err != nil {
			storageMgr.Close()
			return nil, fmt.Errorf("failed to load metadata: %w", err)
		}
	}

	serviceIndex, err := router.NewServiceIndex(serviceStore, operationStore)
	if err != nil {
		storageMgr.Close()
		return nil, fmt.Errorf("failed to build service index: %w", err)
	}
	serviceRouter := router.NewServiceRouter(serviceIndex)

	shutdownCtx, shutdownCancel := context.WithCancel(context.Background())

	srv := &Server{
		config:         cfg,
		storageManager: storageMgr,
		dispatcher:     disp,
		serviceRouter:  serviceRouter,
		serviceStore:   serviceStore,
		operationStore: operationStore,
		shutdownCtx:    shutdownCtx,
		shutdownCancel: shutdownCancel,
		iamStore:       iamStore,
		s3Store:        s3Store,
		blobStore:      blobStore,
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
		srv.chainGateway = chain.NewGateway(serviceRouter, disp)
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
			r.Use(SignatureMiddleware(s.config.SignatureConfig))
		}

		r.Use(LoggingMiddleware)

		s.registerRoutes(r)

		s.httpServerMu.Lock()
		s.httpServer = &http.Server{
			Addr:              ":" + s.config.Port,
			Handler:           r,
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
				log.Printf("Failed to get TLS config: %v", err)
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
					log.Printf("TLS server error: %v", err)
				}
			}()
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
			log.Printf("handleShutdown: server shutdown error: %v", err)
		}
	}

	s.tlsServerMu.Lock()
	tlsSrv := s.tlsServer
	s.tlsServerMu.Unlock()

	if tlsSrv != nil {
		if err := tlsSrv.Shutdown(ctx); err != nil {
			log.Printf("handleShutdown: TLS server shutdown error: %v", err)
		}
	}

	s.shutdownHooksMu.Lock()
	hooks := make([]ShutdownHook, len(s.shutdownHooks))
	copy(hooks, s.shutdownHooks)
	s.shutdownHooksMu.Unlock()

	for i := len(hooks) - 1; i >= 0; i-- {
		hooks[i](ctx)
	}

	s.closeOnce.Do(func() {
		if err := s.storageManager.Close(); err != nil {
			log.Printf("handleShutdown: storage manager close error: %v", err)
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
			log.Printf("Close: server shutdown error: %v", err)
		}
	}

	s.tlsServerMu.Lock()
	tlsSrv := s.tlsServer
	s.tlsServerMu.Unlock()
	if tlsSrv != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tlsSrv.Shutdown(ctx); err != nil {
			log.Printf("Close: TLS server shutdown error: %v", err)
		}
	}

	var err error
	s.closeOnce.Do(func() {
		err = s.storageManager.Close()
	})
	return err
}
