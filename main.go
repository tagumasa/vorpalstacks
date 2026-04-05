// Vorpalstacks is an AWS-compatible cloud platform for edge and on-premises environments.
//
// It provides 31 service APIs covering 27 AWS services with a single binary,
// using CockroachDB Pebble for persistent storage and supporting both
// JSON and Query AWS API protocols.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"vorpalstacks/internal/client/mobyclient"
	appconfig "vorpalstacks/internal/config"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/server/eventbus"
	"vorpalstacks/internal/server/grpcweb"
	chihttp "vorpalstacks/internal/server/http"
	"vorpalstacks/internal/server/listener"
	svcacm "vorpalstacks/internal/services/aws/acm"
	svcapigateway "vorpalstacks/internal/services/aws/apigateway"
	svcapigatewayruntime "vorpalstacks/internal/services/aws/apigateway/runtime"
	svcappsync "vorpalstacks/internal/services/aws/appsync"
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
	svcneptune "vorpalstacks/internal/services/aws/neptune"
	svcneptunedata "vorpalstacks/internal/services/aws/neptunedata"
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
	svcwafv2 "vorpalstacks/internal/services/aws/wafv2"
	storeapigateway "vorpalstacks/internal/store/aws/apigateway"
	storeevents "vorpalstacks/internal/store/aws/eventbridge"
	iamstore "vorpalstacks/internal/store/aws/iam"
	storekinesis "vorpalstacks/internal/store/aws/kinesis"
	neptunestore "vorpalstacks/internal/store/aws/neptune"
	storesns "vorpalstacks/internal/store/aws/sns"
	storesqs "vorpalstacks/internal/store/aws/sqs"
	svcarn "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/pkg/graphengine"
)

func main() {
	cfg := appconfig.LoadBootstrapConfig()

	serverCfg := &chihttp.Config{
		Port:      cfg.Port,
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

	server, err := chihttp.NewServer(serverCfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create server: %v\n", err)
		os.Exit(1)
	}

	appconfig.Initialise(server.Storage())

	var neptunedataService *svcneptunedata.NeptuneDataService
	var sharedNeptuneStore *neptunestore.NeptuneStore
	var neptuneService *svcneptune.NeptuneService
	var appsyncService *svcappsync.AppSyncService
	var logsService *svclogs.LogsService
	var eventsService *svcevents.EventsService
	var stepFunctionService *svcstepfunction.StepFunctionService
	var snsService *svcsns.SNSService
	var sqsService *svcsqs.SQSService

	if cfg.Neptune {
		neptuneService = svcneptune.NewNeptuneService(cfg.AccountID, cfg.Region)
		neptuneService.SetStorageManager(server.StorageManager())
		neptuneService.RegisterHandlers(server.Dispatcher())
	}

	if cfg.NeptuneData {
		sharedNeptuneStore = neptunestore.NewNeptuneStore(server.Storage())
		neptunedataService = svcneptunedata.NewNeptuneDataService(sharedNeptuneStore)
		neptunedataService.RegisterHandlers(server.Dispatcher())
	}

	grpcWebServer := grpcweb.NewServer(&grpcweb.Config{Port: cfg.GRPCWebPort, BindAddr: cfg.GRPCWebBindAddr})
	grpcweb.RegisterAllAdminHandlers(grpcWebServer, grpcweb.AdminDeps{
		Server:           server,
		AccountID:        cfg.AccountID,
		Region:           cfg.Region,
		DataPath:         cfg.DataPath,
		BaseURL:          appconfig.BaseURL(),
		NeptuneDataAdmin: neptunedataService,
		NeptuneAdmin:     neptuneService,
		AppSyncAdmin:     appsyncService,
		LogsAdmin:        logsService,
		EventsAdmin:      eventsService,
		SFNAdmin:         stepFunctionService,
		SNSAdmin:         snsService,
		SQSAdmin:         sqsService,
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

	if eb := server.EventBus(); eb != nil {
		if rr, ok := eb.RoleResolver().(*eventbus.IAMRoleResolver); ok {
			globalStore := server.Storage()
			iamStoreInstance := iamstore.NewIAMStore(globalStore, cfg.AccountID)
			rr.SetLookup(func(ctx context.Context, roleARN string) (string, error) {
				roleName := svcarn.ExtractRoleNameFromARN(roleARN)
				if roleName == "" {
					return "", fmt.Errorf("invalid role ARN format: %q", roleARN)
				}
				role, err := iamStoreInstance.Roles().Get(roleName)
				if err != nil {
					return "", err
				}
				return role.AssumeRolePolicyDocument, nil
			})
		}
	}

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
	s3Service.SetEventBus(server.EventBus())

	s3ObjectStore := s3Store.Objects(cfg.Region)

	mainPort, _ := strconv.Atoi(cfg.Port)
	if mainPort == 0 {
		mainPort = 8080
	}
	lm := listener.NewManager(mainPort)

	s3WebsiteServer := svcs3.NewWebsiteServer(server.StorageManager(), blobStore, cfg.AccountID, cfg.Region)
	s3WebsiteHandler := http.HandlerFunc(s3WebsiteServer.HandleRequest)
	s3WebsitePort := appconfig.GetInt("ports.s3_website")
	if s3WebsitePort == 0 {
		s3WebsitePort = 8081
	}
	lm.Register(listener.ListenerConfig{
		Name:        "s3_website",
		PortKey:     "ports.s3_website",
		DefaultPort: s3WebsitePort,
		Handler:     s3WebsiteHandler,
	})

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
	cloudWatchService.SetStorageManager(server.StorageManager())
	cloudWatchService.RegisterHandlers(server.Dispatcher())

	dynamodbService := svcdynamodb.NewDynamoDBService(cfg.AccountID)
	dynamodbService.RegisterHandlers(server.Dispatcher())

	acmService := svcacm.NewACMService(server.Storage(), cfg.AccountID, cfg.Region)
	acmService.RegisterHandlers(server.Dispatcher())

	cloudfrontService := svccloudfront.NewCloudFrontService(cfg.AccountID)
	cloudfrontService.RegisterHandlers(server.Dispatcher())

	cloudfrontDistServer := svccloudfront.NewDistributionServer(server.StorageManager(), cfg.AccountID, cfg.Region)
	cloudfrontHandler := http.HandlerFunc(cloudfrontDistServer.HandleRequest)
	cloudfrontPort := appconfig.GetInt("ports.cloudfront")
	if cloudfrontPort == 0 {
		cloudfrontPort = 8084
	}
	lm.Register(listener.ListenerConfig{
		Name:        "cloudfront",
		PortKey:     "ports.cloudfront",
		DefaultPort: cloudfrontPort,
		Handler:     cloudfrontHandler,
	})

	wafv2Service := svcwafv2.NewWAFv2Service(server.Storage(), cfg.AccountID, cfg.Region)
	wafv2Service.RegisterHandlers(server.Dispatcher())

	secretsManagerService := svcsecretsmanager.NewSecretsManagerService(cfg.AccountID)
	secretsManagerService.SetRegion(cfg.Region)
	secretsManagerService.SetStorageManager(server.StorageManager())
	secretsManagerService.RegisterHandlers(server.Dispatcher())

	var sqsStoreInstance *storesqs.SQSStore
	if cfg.SQS {
		regionalStorage, _ := server.StorageManager().GetStorage(cfg.Region)
		sqsStoreInstance = storesqs.NewSQSStore(regionalStorage, cfg.AccountID, cfg.Region, appconfig.BaseURL())
		sqsService = svcsqs.NewSQSServiceWithStore(sqsStoreInstance)
		sqsService.SetStorageManager(server.StorageManager())
		sqsService.RegisterHandlers(server.Dispatcher())
		s3Service.SetSQSStore(sqsStoreInstance)
	}

	var kinesisStoreInstance *storekinesis.KinesisStore

	if cfg.Kinesis {
		kinesisRegionalStorage, _ := server.StorageManager().GetStorage(cfg.Region)
		tstore, ok := kinesisRegionalStorage.(storage.TransactionalStorageWith2PC)
		if ok {
			kinesisStoreInstance = storekinesis.NewKinesisStore(tstore, cfg.AccountID, cfg.Region)
		}
		kinesisService := svckinesis.NewKinesisService(kinesisRegionalStorage, cfg.AccountID, cfg.Region)
		kinesisService.RegisterHandlers(server.Dispatcher())
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
			lambdaService.SetStorageManager(server.StorageManager())
			lambdaService.SetHostEndpoint(fmt.Sprintf("http://host.docker.internal:%s", cfg.Port))
			lambdaService.SetS3ObjectStore(cfg.Region, s3Store.Objects(cfg.Region))
			if sqsStoreInstance != nil {
				lambdaService.SetSQSStore(sqsStoreInstance)
			}
			if kinesisStoreInstance != nil {
				lambdaService.SetKinesisStore(kinesisStoreInstance)
			}
			lambdaService.RegisterHandlers(server.Dispatcher())
			s3Service.SetLambdaInvoker(lambdaService)

			lambdaUrlServer := svclambda.NewFunctionURLServer(server.StorageManager(), cfg.AccountID, cfg.Region, lambdaService)
			lambdaUrlHandler := http.HandlerFunc(lambdaUrlServer.HandleRequest)
			lambdaUrlPort := appconfig.GetInt("ports.lambda_url")
			if lambdaUrlPort == 0 {
				lambdaUrlPort = 8085
			}
			lm.Register(listener.ListenerConfig{
				Name:        "lambda_url",
				PortKey:     "ports.lambda_url",
				DefaultPort: lambdaUrlPort,
				Handler:     lambdaUrlHandler,
			})

			server.RegisterShutdownHook(func(ctx context.Context) {
				lambdaService.Shutdown()
			})
			lambdaService.StartESMPoller(context.Background())
			server.RegisterShutdownHook(func(ctx context.Context) {
				lambdaService.StopESMPoller()
			})
		}
	}

	cloudWatchService.SetEventBus(server.EventBus())
	if lambdaService != nil {
		cloudWatchService.SetLambdaInvoker(lambdaService)
		secretsManagerService.SetLambdaInvoker(lambdaService)
	}
	cloudWatchService.StartEvaluator(context.Background())
	secretsManagerService.StartRotationChecker(context.Background())
	server.RegisterShutdownHook(func(ctx context.Context) {
		cloudWatchService.StopEvaluator()
	})
	server.RegisterShutdownHook(func(ctx context.Context) {
		secretsManagerService.StopRotationChecker()
	})

	if cfg.Logs {
		logsService = svclogs.NewLogsService(server.StorageManager(), cfg.AccountID, cfg.DataPath)
		if lambdaService != nil {
			logsService.SetLambdaInvoker(cfg.Region, lambdaService)
		}
		if kinesisStoreInstance != nil {
			logsService.SetKinesisStore(cfg.Region, kinesisStoreInstance)
		}
		logsService.SetEventBus(server.EventBus())
		logsService.RegisterHandlers(server.Dispatcher())
		if lambdaService != nil {
			lambdaService.SetEventBus(server.EventBus())
		}
		server.RegisterShutdownHook(func(ctx context.Context) {
			logsService.Stop()
		})
	}

	var snsStoreInstance *storesns.SNSStore
	if cfg.SNS {
		snsRegionalStorage, _ := server.StorageManager().GetStorage(cfg.Region)
		snsStoreInstance = storesns.NewSNSStore(snsRegionalStorage, cfg.AccountID, cfg.Region)
		server.RegisterShutdownHook(func(ctx context.Context) {
			snsStoreInstance.Close()
		})
		snsService = svcsns.NewSNSService(server.StorageManager(), cfg.AccountID, cfg.Region)
		if sqsStoreInstance != nil {
			snsService.SetSQSStore(sqsStoreInstance)
		}
		if lambdaService != nil {
			snsService.SetLambdaInvoker(lambdaService)
		}
		snsService.SetSNSStore(cfg.Region, snsStoreInstance)
		snsService.SetEventBus(server.EventBus())
		snsService.RegisterHandlers(server.Dispatcher())
	}

	var eventsStoreInstance *storeevents.EventsStore
	if cfg.Events {
		eventsRegionalStorage, _ := server.StorageManager().GetStorage(cfg.Region)
		eventsStoreInstance = storeevents.NewEventsStore(eventsRegionalStorage, cfg.AccountID, cfg.Region)
		eventsService = svcevents.NewEventsService(server.StorageManager(), cfg.AccountID)
		eventsService.SetEventsStore(cfg.Region, eventsStoreInstance)
		if sqsStoreInstance != nil {
			eventsService.SetSQSStore(cfg.Region, sqsStoreInstance)
		}
		if snsStoreInstance != nil {
			eventsService.SetSNSStore(cfg.Region, snsStoreInstance)
		}
		if snsService != nil {
			eventsService.SetSNSPublisher(snsService)
		}
		if lambdaService != nil {
			eventsService.SetLambdaInvoker(lambdaService)
		}
		eventsService.SetEventBus(server.EventBus())
		eventsService.RegisterHandlers(server.Dispatcher())
	}

	if cfg.StepFunctions {
		stepFunctionService = svcstepfunction.NewStepFunctionService(server.StorageManager(), cfg.AccountID)
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
		stepFunctionService.SetEventBus(server.EventBus())
		stepFunctionService.RegisterHandlers(server.Dispatcher())
		server.RegisterShutdownHook(func(ctx context.Context) {
			stepFunctionService.Shutdown()
		})
	}

	if cfg.APIGateway {
		apiGatewayService := svcapigateway.NewAPIGatewayService(server.Storage(), cfg.AccountID, cfg.Region)
		apiGatewayService.RegisterHandlers(server.Dispatcher())

		if lambdaService != nil {
			apiGatewayStorage, err := server.StorageManager().GetStorage(cfg.Region)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to get API Gateway storage: %v\n", err)
				os.Exit(1)
			}
			restApiStore := storeapigateway.NewRestApiStore(apiGatewayStorage, cfg.AccountID, cfg.Region)
			usageStore := storeapigateway.NewUsageStore(apiGatewayStorage, cfg.AccountID, cfg.Region)
			runtimeServer := svcapigatewayruntime.NewRuntimeServer(restApiStore, usageStore, lambdaService)
			runtimeServer.SetAccountID(cfg.AccountID)
			if sqsStoreInstance != nil {
				runtimeServer.SetSQSStore(sqsStoreInstance, cfg.AccountID, cfg.Region)
			}
			if snsStoreInstance != nil {
				runtimeServer.SetSNSStore(snsStoreInstance, cfg.AccountID, cfg.Region)
			}
			runtimeServer.SetEventBus(server.EventBus())
			runtimeHandler := http.HandlerFunc(runtimeServer.HandleRequest)
			server.RegisterAPIGatewayRuntimeHandler(runtimeHandler)

			apiGatewayPort := appconfig.GetInt("ports.apigateway")
			if apiGatewayPort == 0 {
				apiGatewayPort = 8082
			}
			lm.Register(listener.ListenerConfig{
				Name:        "apigateway",
				PortKey:     "ports.apigateway",
				DefaultPort: apiGatewayPort,
				Handler:     runtimeHandler,
			})
		}
	}

	var cognitoService *svccognito.CognitoService
	if cfg.Cognito {
		cognitoService = svccognito.NewCognitoService(server.Storage(), cfg.AccountID, cfg.Region)
		cognitoService.SetStorageManager(server.StorageManager())
		cognitoService.SetEventBus(server.EventBus())
		if lambdaService != nil {
			cognitoService.SetLambdaInvoker(lambdaService)
		}
		cognitoService.RegisterHandlers(server.Dispatcher())
		server.RegisterJWKSHandler(http.HandlerFunc(cognitoService.JWKSHandler))

		cognitoUIHandler := http.HandlerFunc(cognitoService.HostedUIHandler)
		cognitoUIPort := appconfig.GetInt("ports.cognito_hosted")
		if cognitoUIPort == 0 {
			cognitoUIPort = 8083
		}
		lm.Register(listener.ListenerConfig{
			Name:        "cognito_hosted",
			PortKey:     "ports.cognito_hosted",
			DefaultPort: cognitoUIPort,
			Handler:     cognitoUIHandler,
		})
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
		schedulerService.SetEventBus(server.EventBus())
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

	if cfg.AppSync {
		appsyncService = svcappsync.NewAppSyncService(cfg.AccountID)
		if lambdaService != nil {
			appsyncService.SetLambdaInvoker(lambdaService)
		}
		appsyncService.SetEventBus(server.EventBus())
		appsyncService.RegisterHandlers(server.Dispatcher())

		appsyncEventsServer := svcappsync.NewEventServer()
		appsyncEventsServer.SetEventBus(server.EventBus())
		appsyncEventsPort := appconfig.GetInt("ports.appsync_events")
		if appsyncEventsPort == 0 {
			appsyncEventsPort = 8086
		}
		lm.Register(listener.ListenerConfig{
			Name:        "appsync_events",
			PortKey:     "ports.appsync_events",
			DefaultPort: appsyncEventsPort,
			Handler:     appsyncEventsServer,
			Timeouts: &listener.ListenerTimeouts{
				ReadHeaderTimeout: 5 * time.Second,
				ReadTimeout:       0,
				WriteTimeout:      0,
				IdleTimeout:       0,
			},
		})
	}

	graphDir := cfg.DataPath + "/graph"
	graphDB, err := graphengine.Open(graphDir, graphengine.DefaultOptions())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to open graph DB at %s: %v\n", graphDir, err)
	} else {
		server.Dispatcher().SetGraphDB(graphDB)
		server.RegisterShutdownHook(func(ctx context.Context) {
			graphDB.Close()
		})
	}

	if eb := server.EventBus(); eb != nil {
		if lambdaService != nil {
			eb.SetResourcePolicyFunc("lambda", eventbus.LambdaResourcePolicyFn(
				func(ctx context.Context, functionARN string) ([]eventbus.LambdaPolicyEntry, error) {
					fnName := svcarn.ExtractFunctionNameFromARN(functionARN)
					policies, err := lambdaService.GetFunctionPolicy(fnName)
					if err != nil {
						return nil, err
					}
					entries := make([]eventbus.LambdaPolicyEntry, len(policies))
					for i, p := range policies {
						entries[i] = eventbus.LambdaPolicyEntry{
							Statement: p.Statement,
							Principal: p.Principal,
							Action:    p.Action,
							Resource:  p.Resource,
						}
					}
					return entries, nil
				},
			))
		}

		if sqsStoreInstance != nil {
			eb.SetResourcePolicyFunc("sqs", eventbus.SQSResourcePolicyFn(
				func(ctx context.Context, queueARN string) (string, error) {
					queueName := svcarn.ExtractQueueNameFromARN(queueARN)
					queue, err := sqsStoreInstance.GetQueueByName(queueName)
					if err != nil {
						return "", err
					}
					return queue.Policy, nil
				},
			))
		}

		if snsStoreInstance != nil {
			eb.SetResourcePolicyFunc("sns", eventbus.SNSTopicResourcePolicyFn(
				func(ctx context.Context, topicARN string) (string, error) {
					topic, err := snsStoreInstance.GetTopic(topicARN)
					if err != nil {
						return "", err
					}
					policyJSON, ok := topic.Attributes["Policy"]
					if !ok {
						return "", nil
					}
					return policyJSON, nil
				},
			))
		}

		if eventsStoreInstance != nil {
			eb.SetResourcePolicyFunc("events", eventbus.EventBridgeBusResourcePolicyFn(
				func(ctx context.Context, busName string) (string, error) {
					bus, err := eventsStoreInstance.GetEventBus(context.Background(), busName)
					if err != nil {
						return "", err
					}
					return bus.Policy, nil
				},
			))
		}
	}

	cfg.PrintStartupBanner()

	server.SetListenerManager(lm)

	if err := server.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
