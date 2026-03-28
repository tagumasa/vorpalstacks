// Vorpalstacks is an AWS-compatible cloud platform for edge and on-premises environments.
//
// It provides 29 service APIs covering 25 AWS services with a single binary,
// using CockroachDB Pebble for persistent storage and supporting both
// JSON and Query AWS API protocols.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"vorpalstacks/internal/client/mobyclient"
	appconfig "vorpalstacks/internal/config"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/server/grpcweb"
	chihttp "vorpalstacks/internal/server/http"
	svcacm "vorpalstacks/internal/services/aws/acm"
	svcapigateway "vorpalstacks/internal/services/aws/apigateway"
	svcapigatewayruntime "vorpalstacks/internal/services/aws/apigateway/runtime"
	svcathena "vorpalstacks/internal/services/aws/athena"
	svccloudfront "vorpalstacks/internal/services/aws/cloudfront"
	svccloudtrail "vorpalstacks/internal/services/aws/cloudtrail"
	svccloudwatch "vorpalstacks/internal/services/aws/cloudwatch"
	svclogs "vorpalstacks/internal/services/aws/cloudwatchlogs"
	svccognitoidentity "vorpalstacks/internal/services/aws/cognitoidentity"
	svccognito "vorpalstacks/internal/services/aws/cognitoidentityprovider"
	svcdynamodb "vorpalstacks/internal/services/aws/dynamodb"
	svcevents "vorpalstacks/internal/services/aws/eventbridge"
	svciam "vorpalstacks/internal/services/aws/iam"
	svckinesis "vorpalstacks/internal/services/aws/kinesis"
	svckms "vorpalstacks/internal/services/aws/kms"
	svckmshsm "vorpalstacks/internal/services/aws/kms/hsm"
	svclambda "vorpalstacks/internal/services/aws/lambda"
	svcroute53 "vorpalstacks/internal/services/aws/route53"
	svcs3 "vorpalstacks/internal/services/aws/s3"
	svcscheduler "vorpalstacks/internal/services/aws/scheduler"
	svcsecretsmanager "vorpalstacks/internal/services/aws/secretsmanager"
	svcsesv2 "vorpalstacks/internal/services/aws/sesv2"
	svcstepfunction "vorpalstacks/internal/services/aws/sfn"
	svcsns "vorpalstacks/internal/services/aws/sns"
	svcsqs "vorpalstacks/internal/services/aws/sqs"
	svcssm "vorpalstacks/internal/services/aws/ssm"
	svcsts "vorpalstacks/internal/services/aws/sts"
	svctimestreamquery "vorpalstacks/internal/services/aws/timestreamquery"
	svctimestreamwrite "vorpalstacks/internal/services/aws/timestreamwrite"
	svcwaf "vorpalstacks/internal/services/aws/waf"
	svcwafv2 "vorpalstacks/internal/services/aws/wafv2"
	storeapigateway "vorpalstacks/internal/store/aws/apigateway"
	storelogs "vorpalstacks/internal/store/aws/cloudwatchlogs"
	storeevents "vorpalstacks/internal/store/aws/eventbridge"
	storekinesis "vorpalstacks/internal/store/aws/kinesis"
	storesns "vorpalstacks/internal/store/aws/sns"
	storesqs "vorpalstacks/internal/store/aws/sqs"
)

func main() {
	cfg := appconfig.LoadBootstrapConfig()

	serverCfg := &chihttp.Config{
		Port:         cfg.Port,
		DataPath:     cfg.DataPath,
		MetadataPath: cfg.MetadataPath,
		AccountID:    cfg.AccountID,
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

	server, err := chihttp.NewServer(serverCfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create server: %v\n", err)
		os.Exit(1)
	}

	appconfig.Initialise(server.Storage())

	grpcWebServer := grpcweb.NewServer(&grpcweb.Config{Port: cfg.GRPCWebPort, BindAddr: cfg.GRPCWebBindAddr})
	grpcweb.RegisterAllAdminHandlers(grpcWebServer, grpcweb.AdminDeps{
		Server:    server,
		AccountID: cfg.AccountID,
		Region:    cfg.Region,
		DataPath:  cfg.DataPath,
		BaseURL:   appconfig.BaseURL(),
	})

	server.RegisterShutdownHook(func(ctx context.Context) {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := grpcWebServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "gRPC-Web server shutdown error: %v\n", err)
		}
	})

	go func() {
		fmt.Printf("Starting gRPC-Web admin server on :%s\n", cfg.GRPCWebPort)
		if err := grpcWebServer.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "gRPC-Web server error: %v\n", err)
		}
	}()

	iamService := svciam.NewIAMService(cfg.AccountID)
	iamService.RegisterHandlers(server.Dispatcher())
	server.RegisterShutdownHook(func(ctx context.Context) {
		svciam.ShutdownSLRoleCleanup()
	})

	stsService := svcsts.NewSTSService()
	stsService.RegisterHandlers(server.Dispatcher())

	hsmBackend, err := svckmshsm.NewPersistentBackend(cfg.DataPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create persistent HSM backend: %v\n", err)
		os.Exit(1)
	}
	kmsService := svckms.NewKMSService(server.Storage(), cfg.AccountID, cfg.Region, hsmBackend)
	kmsService.RegisterHandlers(server.Dispatcher())

	if err := kmsService.EnsureDefaultSSMKey(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to create default SSM key: %v\n", err)
	}

	s3Store := server.S3Store()
	blobStore := server.BlobStore()
	s3Service := svcs3.NewS3Service(s3Store, blobStore, cfg.AccountID)
	s3Handler := svcs3.NewS3Handler(s3Service, cfg.Region, server.StorageManager())
	server.RegisterS3Handler(s3Handler)

	s3ObjectStore := s3Store.Objects(cfg.Region)

	route53Service, err := svcroute53.NewRoute53ServiceWithDNS(server.Storage(), cfg.AccountID, cfg.Route53DNSBindAddr, cfg.Route53DNSEnabled)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create Route53 service: %v\n", err)
		os.Exit(1)
	}
	route53Service.RegisterHandlers(server.Dispatcher())
	server.RegisterShutdownHook(func(ctx context.Context) {
		if err := route53Service.Shutdown(); err != nil {
			fmt.Fprintf(os.Stderr, "Route53 shutdown error: %v\n", err)
		}
	})

	cloudWatchService := svccloudwatch.NewCloudWatchService(server.Storage(), cfg.AccountID, cfg.Region, cfg.DataPath)
	cloudWatchService.RegisterHandlers(server.Dispatcher())

	dynamodbService := svcdynamodb.NewDynamoDBService(cfg.AccountID)
	dynamodbService.RegisterHandlers(server.Dispatcher())

	acmService := svcacm.NewACMService(server.Storage(), cfg.AccountID, cfg.Region)
	acmService.RegisterHandlers(server.Dispatcher())

	cloudfrontService := svccloudfront.NewCloudFrontService(cfg.AccountID)
	cloudfrontService.RegisterHandlers(server.Dispatcher())

	wafService := svcwaf.NewWAFService(server.Storage(), cfg.AccountID, cfg.Region)
	wafService.RegisterHandlers(server.Dispatcher())

	wafv2Service := svcwafv2.NewWAFv2Service(server.Storage(), cfg.AccountID, cfg.Region)
	wafv2Service.RegisterHandlers(server.Dispatcher())

	secretsManagerService := svcsecretsmanager.NewSecretsManagerService(cfg.AccountID)
	secretsManagerService.SetStorageManager(server.StorageManager())
	secretsManagerService.RegisterHandlers(server.Dispatcher())

	var sqsStoreInstance *storesqs.SQSStore
	if cfg.SQS {
		sqsStoreInstance = storesqs.NewSQSStore(server.Storage(), cfg.AccountID, cfg.Region, appconfig.BaseURL())
		sqsService := svcsqs.NewSQSServiceWithStore(sqsStoreInstance)
		sqsService.RegisterHandlers(server.Dispatcher())
	}

	var logsStoreInstance *storelogs.Store
	var kinesisStoreInstance *storekinesis.KinesisStore

	if cfg.Kinesis {
		tstore, ok := server.Storage().(storage.TransactionalStorageWith2PC)
		if ok {
			kinesisStoreInstance = storekinesis.NewKinesisStore(tstore, cfg.AccountID, cfg.Region)
		}
		kinesisService := svckinesis.NewKinesisService(server.Storage(), cfg.AccountID, cfg.Region)
		kinesisService.RegisterHandlers(server.Dispatcher())
	}

	if cfg.Logs {
		logsStoreInstance = storelogs.NewStore(server.Storage().Bucket("logs"), cfg.AccountID, cfg.Region, cfg.DataPath)
	}

	var lambdaService *svclambda.LambdaService
	if cfg.Lambda {
		dockerCfg := mobyclient.Config{
			Host:    cfg.DockerHost,
			Version: "1.44",
		}
		dockerLogger := logs.NewLogger(&logs.Config{Level: logs.LevelInfo})
		dockerClient, err := mobyclient.NewClient(dockerCfg, dockerLogger)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to create Docker client for Lambda: %v\n", err)
		} else {
			lambdaService = svclambda.NewLambdaService(server.Storage(), dockerClient, cfg.AccountID, cfg.Region, cfg.DataPath)
			if logsStoreInstance != nil {
				lambdaService.SetLogsStore(cfg.Region, logsStoreInstance)
			}
			lambdaService.SetS3ObjectStore(cfg.Region, s3Store.Objects(cfg.Region))
			lambdaService.RegisterHandlers(server.Dispatcher())
			server.RegisterShutdownHook(func(ctx context.Context) {
				lambdaService.Shutdown()
			})
		}
	}

	if cfg.Logs {
		logsService := svclogs.NewLogsService(server.StorageManager(), cfg.AccountID, cfg.DataPath)
		if lambdaService != nil {
			logsService.SetLambdaInvoker(cfg.Region, lambdaService)
		}
		if kinesisStoreInstance != nil {
			logsService.SetKinesisStore(cfg.Region, kinesisStoreInstance)
		}
		logsService.RegisterHandlers(server.Dispatcher())
		server.RegisterShutdownHook(func(ctx context.Context) {
			logsService.Stop()
		})
	}

	var snsStoreInstance *storesns.SNSStore
	if cfg.SNS {
		snsStoreInstance = storesns.NewSNSStore(server.Storage(), cfg.AccountID, cfg.Region)
		snsService := svcsns.NewSNSService(server.StorageManager(), cfg.AccountID)
		if sqsStoreInstance != nil {
			snsService.SetSQSStore(sqsStoreInstance)
		}
		if lambdaService != nil {
			snsService.SetLambdaInvoker(lambdaService)
		}
		snsService.SetSNSStore(cfg.Region, snsStoreInstance)
		snsService.RegisterHandlers(server.Dispatcher())
	}

	var eventsStoreInstance *storeevents.EventsStore
	if cfg.Events {
		eventsStoreInstance = storeevents.NewEventsStore(server.Storage(), cfg.AccountID, cfg.Region)
		eventsService := svcevents.NewEventsService(server.StorageManager(), cfg.AccountID, cfg.Region)
		eventsService.SetEventsStore(eventsStoreInstance)
		if sqsStoreInstance != nil {
			eventsService.SetSQSStore(cfg.Region, sqsStoreInstance)
		}
		if snsStoreInstance != nil {
			eventsService.SetSNSStore(cfg.Region, snsStoreInstance)
		}
		if lambdaService != nil {
			eventsService.SetLambdaInvoker(lambdaService)
		}
		eventsService.RegisterHandlers(server.Dispatcher())
	}

	if cfg.StepFunctions {
		stepFunctionService := svcstepfunction.NewStepFunctionService(server.StorageManager(), cfg.AccountID)
		if lambdaService != nil {
			stepFunctionService.SetLambdaInvoker(lambdaService)
		}
		if sqsStoreInstance != nil {
			stepFunctionService.SetSQSStore(sqsStoreInstance)
		}
		if snsStoreInstance != nil {
			stepFunctionService.SetSNSStore(snsStoreInstance)
		}
		if eventsStoreInstance != nil {
			stepFunctionService.SetEventsStore(eventsStoreInstance)
		}
		stepFunctionService.RegisterHandlers(server.Dispatcher())
		server.RegisterShutdownHook(func(ctx context.Context) {
			stepFunctionService.Shutdown()
		})
	}

	if cfg.APIGateway {
		apiGatewayService := svcapigateway.NewAPIGatewayService(server.Storage(), cfg.AccountID, cfg.Region)
		apiGatewayService.RegisterHandlers(server.Dispatcher())

		if lambdaService != nil {
			restApiStore := storeapigateway.NewRestApiStore(server.Storage(), cfg.AccountID, cfg.Region)
			usageStore := storeapigateway.NewUsageStore(server.Storage(), cfg.AccountID, cfg.Region)
			runtimeServer := svcapigatewayruntime.NewRuntimeServer(restApiStore, usageStore, lambdaService)
			if sqsStoreInstance != nil {
				runtimeServer.SetSQSStore(sqsStoreInstance, cfg.AccountID, cfg.Region)
			}
			if snsStoreInstance != nil {
				runtimeServer.SetSNSStore(snsStoreInstance, cfg.AccountID, cfg.Region)
			}
			server.RegisterAPIGatewayRuntimeHandler(http.HandlerFunc(runtimeServer.HandleRequest))
		}
	}

	var cognitoService *svccognito.CognitoService
	if cfg.Cognito {
		cognitoService = svccognito.NewCognitoService(server.Storage(), cfg.AccountID, cfg.Region)
		cognitoService.SetStorageManager(server.StorageManager())
		cognitoService.RegisterHandlers(server.Dispatcher())
		server.RegisterJWKSHandler(http.HandlerFunc(cognitoService.JWKSHandler))
	}

	if cfg.CognitoIdentity {
		cognitoIdentityService := svccognitoidentity.NewCognitoIdentityService(server.Storage(), cfg.AccountID, cfg.Region)
		cognitoIdentityService.RegisterHandlers(server.Dispatcher())
	}

	if cfg.SSM {
		ssmService := svcssm.NewSSMServiceWithKMS(cfg.AccountID, kmsService)
		ssmService.RegisterHandlers(server.Dispatcher())
	}

	if cfg.SESv2 {
		sesv2Service := svcsesv2.NewSESv2Service(cfg.AccountID)
		sesv2Service.RegisterHandlers(server.Dispatcher())
	}

	if cfg.Scheduler {
		schedulerService := svcscheduler.NewSchedulerService(server.StorageManager(), cfg.AccountID)
		if sqsStoreInstance != nil {
			schedulerService.SetSQSStore(sqsStoreInstance)
		}
		if snsStoreInstance != nil {
			schedulerService.SetSNSStore(snsStoreInstance)
		}
		if lambdaService != nil {
			schedulerService.SetLambdaInvoker(lambdaService)
		}
		schedulerService.BuildEngine()
		schedulerService.RegisterHandlers(server.Dispatcher())
		schedulerService.StartEngine()
		server.RegisterShutdownHook(func(ctx context.Context) {
			if err := schedulerService.StopEngine(); err != nil {
				fmt.Fprintf(os.Stderr, "Scheduler shutdown error: %v\n", err)
			}
		})
	}

	if cfg.CloudTrail {
		cloudtrailService := svccloudtrail.NewCloudTrailService(server.Storage(), cfg.AccountID, cfg.Region)
		cloudtrailService.RegisterHandlers(server.Dispatcher())
	}

	if cfg.TimestreamWrite {
		timestreamWriteService := svctimestreamwrite.NewService(cfg.AccountID, cfg.ServerHost(), cfg.DataPath)
		timestreamWriteService.RegisterHandlers(server.Dispatcher())
	}

	if cfg.TimestreamQuery {
		timestreamQueryService := svctimestreamquery.NewService(cfg.AccountID, cfg.ServerHost(), cfg.DataPath)
		timestreamQueryService.RegisterHandlers(server.Dispatcher())
	}

	if cfg.Athena {
		athenaService := svcathena.NewServiceWithS3(cfg.AccountID, cfg.ServerHost(), s3ObjectStore)
		athenaService.RegisterHandlers(server.Dispatcher())
		server.RegisterShutdownHook(func(ctx context.Context) {
			athenaService.Shutdown()
		})
	}

	cfg.PrintStartupBanner()

	if err := server.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
