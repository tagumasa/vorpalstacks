// Vorpalstacks is an AWS-compatible cloud platform for edge and on-premises environments.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"vorpalstacks/internal/client/mobyclient"
	appconfig "vorpalstacks/internal/config"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	acmconnect "vorpalstacks/internal/pb/aws/acm/acmconnect"
	adminconfigconnect "vorpalstacks/internal/pb/aws/admin_config/admin_configconnect"
	apigatewayconnect "vorpalstacks/internal/pb/aws/apigateway/apigatewayconnect"
	athenaconnect "vorpalstacks/internal/pb/aws/athena/athenaconnect"
	cloudfrontconnect "vorpalstacks/internal/pb/aws/cloudfront/cloudfrontconnect"
	cloudtrailconnect "vorpalstacks/internal/pb/aws/cloudtrail/cloudtrailconnect"
	cognitoidentityconnect "vorpalstacks/internal/pb/aws/cognito_identity/cognito_identityconnect"
	cognitoidpconnect "vorpalstacks/internal/pb/aws/cognito_idp/cognito_idpconnect"
	dynamodbconnect "vorpalstacks/internal/pb/aws/dynamodb/dynamodbconnect"
	emailconnect "vorpalstacks/internal/pb/aws/email/emailconnect"
	eventsconnect "vorpalstacks/internal/pb/aws/events/eventsconnect"
	iamconnect "vorpalstacks/internal/pb/aws/iam/iamconnect"
	ingesttimestreamconnect "vorpalstacks/internal/pb/aws/ingest.timestream/ingest_timestreamconnect"
	kinesisconnect "vorpalstacks/internal/pb/aws/kinesis/kinesisconnect"
	kmsconnect "vorpalstacks/internal/pb/aws/kms/kmsconnect"
	lambdaconnect "vorpalstacks/internal/pb/aws/lambda/lambdaconnect"
	logsconnect "vorpalstacks/internal/pb/aws/logs/logsconnect"
	monitoringconnect "vorpalstacks/internal/pb/aws/monitoring/monitoringconnect"
	querytimestreamconnect "vorpalstacks/internal/pb/aws/query.timestream/query_timestreamconnect"
	route53connect "vorpalstacks/internal/pb/aws/route53/route53connect"
	s3connect "vorpalstacks/internal/pb/aws/s3/s3connect"
	schedulerconnect "vorpalstacks/internal/pb/aws/scheduler/schedulerconnect"
	secretsmanagerconnect "vorpalstacks/internal/pb/aws/secretsmanager/secretsmanagerconnect"
	snsconnect "vorpalstacks/internal/pb/aws/sns/snsconnect"
	sqsconnect "vorpalstacks/internal/pb/aws/sqs/sqsconnect"
	ssmconnect "vorpalstacks/internal/pb/aws/ssm/ssmconnect"
	statesconnect "vorpalstacks/internal/pb/aws/states/statesconnect"
	stsconnect "vorpalstacks/internal/pb/aws/sts/stsconnect"
	wafconnect "vorpalstacks/internal/pb/aws/waf/wafconnect"
	wafv2connect "vorpalstacks/internal/pb/aws/wafv2/wafv2connect"
	"vorpalstacks/internal/server/grpcweb"
	chihttp "vorpalstacks/internal/server/http"
	svcacm "vorpalstacks/internal/services/aws/acm"
	svcadminconfig "vorpalstacks/internal/services/aws/admin_config"
	svcapigateway "vorpalstacks/internal/services/aws/apigateway"
	svcapigatewayruntime "vorpalstacks/internal/services/aws/apigateway/runtime"
	svcathena "vorpalstacks/internal/services/aws/athena"
	svccloudfront "vorpalstacks/internal/services/aws/cloudfront"
	svccloudtrail "vorpalstacks/internal/services/aws/cloudtrail"
	svccloudwatch "vorpalstacks/internal/services/aws/cloudwatch"
	svclogs "vorpalstacks/internal/services/aws/cloudwatchlogs"
	svccognitoidentity "vorpalstacks/internal/services/aws/cognitoidentity"
	svccognito "vorpalstacks/internal/services/aws/cognitoidentityprovider"
	"vorpalstacks/internal/services/aws/common/request"
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
	storecloudwatch "vorpalstacks/internal/store/aws/cloudwatch"
	storelogs "vorpalstacks/internal/store/aws/cloudwatchlogs"
	storecognitoidentity "vorpalstacks/internal/store/aws/cognitoidentity"
	storecognitoidp "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	storeevents "vorpalstacks/internal/store/aws/eventbridge"
	storekinesis "vorpalstacks/internal/store/aws/kinesis"
	storekms "vorpalstacks/internal/store/aws/kms"
	storesns "vorpalstacks/internal/store/aws/sns"
	storesqs "vorpalstacks/internal/store/aws/sqs"
	storeconfig "vorpalstacks/internal/store/config"
)

func main() {
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := getEnvOrDefault("AWS_REGION", "us-east-1")
	accountId := os.Getenv("AWS_ACCOUNT_ID")
	if accountId == "" {
		accountId = "000000000000"
	}
	port := getEnvOrDefault("PORT", "8080")
	dataPath := getEnvOrDefault("DATA_PATH", "./data")
	metadataPath := getEnvOrDefault("METADATA_PATH", "")
	signatureEnabled := getEnvOrDefault("SIGNATURE_VERIFICATION_ENABLED", "true") == "true"
	useChainGateway := getEnvOrDefault("USE_CHAIN_GATEWAY", "false") == "true"
	tlsEnabled := getEnvOrDefault("TLS_ENABLED", "false") == "true"
	tlsPort := getEnvOrDefault("TLS_PORT", "8443")
	tlsCertPath := getEnvOrDefault("TLS_CERT_PATH", "auto")
	tlsKeyPath := getEnvOrDefault("TLS_KEY_PATH", "auto")
	tlsHostname := getEnvOrDefault("TLS_HOSTNAME", "")

	cfg := &chihttp.Config{
		Port:         port,
		DataPath:     dataPath,
		MetadataPath: metadataPath,
		AccountID:    accountId,
		SignatureConfig: chihttp.SignatureConfig{
			Enabled:         signatureEnabled,
			Region:          region,
			AccessKeyID:     accessKeyID,
			SecretAccessKey: secretAccessKey,
		},
		UseChainGateway: useChainGateway,
		TLSConfig: chihttp.TLSConfig{
			Enabled:  tlsEnabled,
			Port:     tlsPort,
			CertPath: tlsCertPath,
			KeyPath:  tlsKeyPath,
			Hostname: tlsHostname,
		},
	}

	server, err := chihttp.NewServer(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create server: %v\n", err)
		os.Exit(1)
	}

	appconfig.Initialise(server.Storage())

	grpcWebPort := getEnvOrDefault("GRPC_WEB_PORT", "9090")
	grpcWebBindAddr := getEnvOrDefault("GRPC_WEB_BIND_ADDR", "127.0.0.1")
	grpcWebServer := grpcweb.NewServer(&grpcweb.Config{Port: grpcWebPort, BindAddr: grpcWebBindAddr})

	configStore := storeconfig.NewStore(server.Storage())
	adminConfigService := svcadminconfig.NewService(configStore)
	path, handler := adminconfigconnect.NewAdminConfigServiceHandler(adminConfigService)
	grpcWebServer.Handle(path, handler)

	grpcWebServer.Handle(acmconnect.NewACMServiceHandler(svcacm.NewAdminHandler()))
	grpcWebServer.Handle(apigatewayconnect.NewAPIGatewayServiceHandler(svcapigateway.NewAdminHandler()))
	grpcWebServer.Handle(cloudfrontconnect.NewCloudFrontServiceHandler(svccloudfront.NewAdminHandler()))
	grpcWebServer.Handle(route53connect.NewRoute53ServiceHandler(svcroute53.NewAdminHandler()))
	grpcWebServer.Handle(secretsmanagerconnect.NewSecretsManagerServiceHandler(svcsecretsmanager.NewAdminHandler()))
	grpcWebServer.Handle(stsconnect.NewSTSServiceHandler(svcsts.NewAdminHandler(accountId)))
	grpcWebServer.Handle(wafconnect.NewWAFServiceHandler(svcwaf.NewAdminHandler()))
	grpcWebServer.Handle(wafv2connect.NewWAFV2ServiceHandler(svcwafv2.NewAdminHandler()))
	grpcWebServer.Handle(schedulerconnect.NewSchedulerServiceHandler(svcscheduler.NewAdminHandler()))
	grpcWebServer.Handle(dynamodbconnect.NewDynamoDBServiceHandler(svcdynamodb.NewAdminHandler(server.StorageManager(), accountId)))
	grpcWebServer.Handle(kinesisconnect.NewKinesisServiceHandler(svckinesis.NewAdminHandler(server.StorageManager(), accountId)))
	grpcWebServer.Handle(lambdaconnect.NewLambdaServiceHandler(svclambda.NewAdminHandler(server.StorageManager(), accountId)))
	grpcWebServer.Handle(s3connect.NewS3ServiceHandler(svcs3.NewAdminHandler(server.StorageManager(), accountId)))
	grpcWebServer.Handle(statesconnect.NewSFNServiceHandler(svcstepfunction.NewAdminHandler(server.StorageManager(), accountId)))
	grpcWebServer.Handle(sqsconnect.NewSQSServiceHandler(svcsqs.NewAdminHandler(server.StorageManager(), accountId, appconfig.BaseURL())))
	grpcWebServer.Handle(eventsconnect.NewCloudWatchEventsServiceHandler(svcevents.NewAdminHandler(server.StorageManager(), accountId)))
	grpcWebServer.Handle(logsconnect.NewCloudWatchLogsServiceHandler(svclogs.NewAdminHandler(server.StorageManager(), accountId, dataPath)))
	grpcWebServer.Handle(snsconnect.NewSNSServiceHandler(svcsns.NewAdminHandler(server.StorageManager(), accountId)))
	grpcWebServer.Handle(iamconnect.NewIAMServiceHandler(svciam.NewAdminHandler(server.Storage(), accountId)))
	grpcWebServer.Handle(kmsconnect.NewKMSServiceHandler(svckms.NewAdminHandler(storekms.NewKeyStore(server.Storage(), accountId, region))))
	grpcWebServer.Handle(monitoringconnect.NewCloudWatchServiceHandler(svccloudwatch.NewAdminHandler(
		storecloudwatch.NewAlarmStore(server.Storage(), accountId, region),
		storecloudwatch.NewMetricChunkStore(server.Storage(), region, dataPath),
	)))
	grpcWebServer.Handle(cognitoidpconnect.NewCognitoIdentityProviderServiceHandler(svccognito.NewAdminHandler(
		storecognitoidp.NewCognitoStore(server.Storage(), accountId, region),
	)))
	grpcWebServer.Handle(cognitoidentityconnect.NewCognitoIdentityServiceHandler(svccognitoidentity.NewAdminHandler(
		storecognitoidentity.NewCognitoIdentityStore(server.Storage(), accountId, region),
	)))

	grpcWebServer.Handle(athenaconnect.NewAthenaServiceHandler(svcathena.NewAdminHandler(server.StorageManager(), accountId)))
	grpcWebServer.Handle(cloudtrailconnect.NewCloudTrailServiceHandler(svccloudtrail.NewAdminHandler(server.StorageManager(), accountId)))
	grpcWebServer.Handle(emailconnect.NewSESv2ServiceHandler(svcsesv2.NewAdminHandler(server.StorageManager(), accountId)))
	grpcWebServer.Handle(ssmconnect.NewSSMServiceHandler(svcssm.NewAdminHandler(server.StorageManager(), accountId)))
	grpcWebServer.Handle(querytimestreamconnect.NewTimestreamQueryServiceHandler(svctimestreamquery.NewAdminHandler(server.StorageManager(), accountId, dataPath)))
	grpcWebServer.Handle(ingesttimestreamconnect.NewTimestreamWriteServiceHandler(svctimestreamwrite.NewAdminHandler(server.StorageManager(), accountId, dataPath)))

	go func() {
		fmt.Printf("Starting gRPC-Web admin server on :%s\n", grpcWebPort)
		if err := grpcWebServer.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "gRPC-Web server error: %v\n", err)
		}
	}()

	iamService := svciam.NewIAMService(accountId)
	iamService.RegisterHandlers(server.Dispatcher())

	stsService := svcsts.NewSTSService()
	stsService.RegisterHandlers(server.Dispatcher())

	hsmBackend, err := svckmshsm.NewPersistentBackend(dataPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create persistent HSM backend: %v\n", err)
		os.Exit(1)
	}
	kmsService := svckms.NewKMSService(server.Storage(), accountId, region, hsmBackend)
	kmsService.RegisterHandlers(server.Dispatcher())

	if err := kmsService.EnsureDefaultSSMKey(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to create default SSM key: %v\n", err)
	}

	s3Store := server.S3Store()
	blobStore := server.BlobStore()
	s3Service := svcs3.NewS3Service(s3Store, blobStore, accountId)
	s3Handler := svcs3.NewS3Handler(s3Service, region, server.StorageManager())
	server.RegisterS3Handler(s3Handler)

	s3ObjectStore := s3Store.Objects(region)

	dnsEnabled := getEnvOrDefault("ROUTE53_DNS_ENABLED", "false") == "true"
	dnsBindAddr := getEnvOrDefault("ROUTE53_DNS_BIND_ADDR", "127.0.0.1")

	route53Service, err := svcroute53.NewRoute53ServiceWithDNS(server.Storage(), accountId, dnsBindAddr, dnsEnabled)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create Route53 service: %v\n", err)
		os.Exit(1)
	}
	route53Service.RegisterHandlers(server.Dispatcher())

	cloudWatchService := svccloudwatch.NewCloudWatchService(server.Storage(), accountId, region, dataPath)
	cloudWatchService.RegisterHandlers(server.Dispatcher())

	dynamodbService := svcdynamodb.NewDynamoDBService(accountId)
	dynamodbService.RegisterHandlers(server.Dispatcher())

	acmService := svcacm.NewACMService(server.Storage(), accountId, region)
	acmService.RegisterHandlers(server.Dispatcher())

	cloudfrontService := svccloudfront.NewCloudFrontService(accountId)
	cloudfrontService.RegisterHandlers(server.Dispatcher())

	wafService := svcwaf.NewWAFService(server.Storage(), accountId, region)
	wafService.RegisterHandlers(server.Dispatcher())

	wafv2Service := svcwafv2.NewWAFv2Service(server.Storage(), accountId, region)
	wafv2Service.RegisterHandlers(server.Dispatcher())

	secretsManagerService := svcsecretsmanager.NewSecretsManagerService(accountId)
	secretsManagerService.RegisterHandlers(server.Dispatcher())

	snsEnabled := getEnvOrDefault("SNS_ENABLED", "true") == "true"
	sqsEnabled := getEnvOrDefault("SQS_ENABLED", "true") == "true"
	lambdaEnabled := getEnvOrDefault("LAMBDA_ENABLED", "true") == "true"
	eventsEnabled := getEnvOrDefault("EVENTS_ENABLED", "true") == "true"
	logsEnabled := getEnvOrDefault("LOGS_ENABLED", "true") == "true"

	var sqsStoreInstance *storesqs.SQSStore
	if sqsEnabled {
		sqsStoreInstance = storesqs.NewSQSStore(server.Storage(), accountId, region, appconfig.BaseURL())
		sqsService := svcsqs.NewSQSServiceWithStore(sqsStoreInstance)
		sqsService.RegisterHandlers(server.Dispatcher())
	}

	var logsStoreInstance *storelogs.Store
	var kinesisStoreInstance *storekinesis.KinesisStore

	kinesisEnabled := getEnvOrDefault("KINESIS_ENABLED", "true") == "true"
	if kinesisEnabled {
		tstore, ok := server.Storage().(storage.TransactionalStorageWith2PC)
		if ok {
			kinesisStoreInstance = storekinesis.NewKinesisStore(tstore, accountId, region)
		}
		kinesisService := svckinesis.NewKinesisService(server.Storage(), accountId, region)
		kinesisService.RegisterHandlers(server.Dispatcher())
	}

	if logsEnabled {
		logsStoreInstance = storelogs.NewStore(server.Storage().Bucket("logs"), accountId, region, dataPath)
	}

	var lambdaService *svclambda.LambdaService
	if lambdaEnabled {
		dockerHost := getEnvOrDefault("DOCKER_HOST", "unix:///var/run/docker.sock")
		dockerCfg := mobyclient.Config{
			Host:    dockerHost,
			Version: "1.44",
		}
		dockerLogger := logs.NewLogger(&logs.Config{Level: logs.LevelInfo})
		dockerClient, err := mobyclient.NewClient(dockerCfg, dockerLogger)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to create Docker client for Lambda: %v\n", err)
		} else {
			if logsStoreInstance != nil {
				lambdaService = svclambda.NewLambdaServiceWithLogs(server.Storage(), dockerClient, logsStoreInstance, accountId, region, dataPath)
			} else {
				lambdaService = svclambda.NewLambdaService(server.Storage(), dockerClient, accountId, region, dataPath)
			}
			lambdaService.SetS3ObjectStore(region, s3Store.Objects(region))
			lambdaService.RegisterHandlers(server.Dispatcher())
		}
	}

	if logsEnabled {
		var logsService *svclogs.LogsService
		if lambdaService != nil && kinesisStoreInstance != nil {
			logsService = svclogs.NewLogsServiceWithClients(server.StorageManager(), accountId, region, dataPath, lambdaService, kinesisStoreInstance)
		} else if lambdaService != nil {
			logsService = svclogs.NewLogsServiceWithClients(server.StorageManager(), accountId, region, dataPath, lambdaService, nil)
		} else if kinesisStoreInstance != nil {
			logsService = svclogs.NewLogsServiceWithClients(server.StorageManager(), accountId, region, dataPath, nil, kinesisStoreInstance)
		} else {
			logsService = svclogs.NewLogsService(server.StorageManager(), accountId, dataPath)
		}
		logsService.RegisterHandlers(server.Dispatcher())
	}

	var snsStoreInstance *storesns.SNSStore
	if snsEnabled {
		snsStoreInstance = storesns.NewSNSStore(server.Storage(), accountId, region)
		var snsService *svcsns.SNSService
		if sqsStoreInstance != nil && lambdaService != nil {
			snsService = svcsns.NewSNSServiceWithClients(server.StorageManager(), accountId, region, sqsStoreInstance, lambdaService)
		} else if sqsStoreInstance != nil {
			snsService = svcsns.NewSNSServiceWithSQS(server.StorageManager(), accountId, region, sqsStoreInstance)
		} else if lambdaService != nil {
			snsService = svcsns.NewSNSServiceWithClients(server.StorageManager(), accountId, region, nil, lambdaService)
		} else {
			snsService = svcsns.NewSNSServiceWithStore(server.StorageManager(), accountId, region, snsStoreInstance)
		}
		snsService.RegisterHandlers(server.Dispatcher())
	}

	var eventsStoreInstance *storeevents.EventsStore
	if eventsEnabled {
		eventsStoreInstance = storeevents.NewEventsStore(server.Storage(), accountId, region)
		var eventsService *svcevents.EventsService
		if sqsStoreInstance != nil && snsStoreInstance != nil && lambdaService != nil {
			eventsService = svcevents.NewEventsServiceWithClients(server.StorageManager(), accountId, region, sqsStoreInstance, snsStoreInstance, lambdaService)
		} else if sqsStoreInstance != nil {
			eventsService = svcevents.NewEventsServiceWithClients(server.StorageManager(), accountId, region, sqsStoreInstance, nil, nil)
		} else {
			eventsService = svcevents.NewEventsServiceWithStore(server.StorageManager(), accountId, region, eventsStoreInstance)
		}
		eventsService.RegisterHandlers(server.Dispatcher())
	}

	stepFunctionsEnabled := getEnvOrDefault("STEPFUNCTIONS_ENABLED", "true") == "true"
	if stepFunctionsEnabled {
		var stepFunctionService *svcstepfunction.StepFunctionService
		if sqsStoreInstance != nil && snsStoreInstance != nil && eventsStoreInstance != nil {
			stepFunctionService = svcstepfunction.NewStepFunctionServiceWithStores(server.StorageManager(), accountId, lambdaService, sqsStoreInstance, snsStoreInstance, eventsStoreInstance)
		} else if lambdaService != nil {
			stepFunctionService = svcstepfunction.NewStepFunctionServiceWithClients(server.StorageManager(), accountId, lambdaService)
		} else {
			stepFunctionService = svcstepfunction.NewStepFunctionService(server.StorageManager(), accountId)
		}
		stepFunctionService.RegisterHandlers(server.Dispatcher())
	}

	apiGatewayEnabled := getEnvOrDefault("APIGATEWAY_ENABLED", "true") == "true"
	if apiGatewayEnabled {
		apiGatewayService := svcapigateway.NewAPIGatewayService(server.Storage(), accountId, region)
		apiGatewayService.RegisterHandlers(server.Dispatcher())

		if lambdaService != nil {
			restApiStore := storeapigateway.NewRestApiStore(server.Storage(), accountId, region)
			usageStore := storeapigateway.NewUsageStore(server.Storage(), accountId, region)
			var runtimeServer *svcapigatewayruntime.RuntimeServer
			if sqsStoreInstance != nil || snsStoreInstance != nil {
				runtimeServer = svcapigatewayruntime.NewRuntimeServerWithStores(restApiStore, usageStore, lambdaService, sqsStoreInstance, snsStoreInstance, accountId, region)
			} else {
				runtimeServer = svcapigatewayruntime.NewRuntimeServer(restApiStore, usageStore, lambdaService)
			}
			server.RegisterAPIGatewayRuntimeHandler(http.HandlerFunc(runtimeServer.HandleRequest))
		}
	}

	cognitoEnabled := getEnvOrDefault("COGNITO_ENABLED", "true") == "true"
	var cognitoService *svccognito.CognitoService
	if cognitoEnabled {
		cognitoService = svccognito.NewCognitoService(server.Storage(), accountId, region)
		cognitoService.RegisterHandlers(server.Dispatcher())
		server.RegisterJWKSHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqCtx := request.NewRequestContext(r.Context(), server.StorageManager(), accountId, region)
			userPoolID := r.URL.Query().Get("userPoolId")
			if userPoolID == "" {
				pools, _ := cognitoService.ListUserPoolsRaw(reqCtx)
				if len(pools) > 0 {
					userPoolID = pools[0].ID
				}
			}
			if userPoolID == "" {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"keys":[]}`))
				return
			}
			jwks, err := cognitoService.GetJWKS(reqCtx, userPoolID)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"keys":[]}`))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(jwks)
		}))
	}

	cognitoIdentityEnabled := getEnvOrDefault("COGNITO_IDENTITY_ENABLED", "true") == "true"
	if cognitoIdentityEnabled {
		cognitoIdentityService := svccognitoidentity.NewCognitoIdentityService(server.Storage(), accountId, region)
		cognitoIdentityService.RegisterHandlers(server.Dispatcher())
	}

	ssmEnabled := getEnvOrDefault("SSM_ENABLED", "true") == "true"
	if ssmEnabled {
		ssmService := svcssm.NewSSMServiceWithKMS(accountId, kmsService)
		ssmService.RegisterHandlers(server.Dispatcher())
	}

	sesv2Enabled := getEnvOrDefault("SESV2_ENABLED", "true") == "true"
	if sesv2Enabled {
		sesv2Service := svcsesv2.NewSESv2Service(accountId)
		sesv2Service.RegisterHandlers(server.Dispatcher())
	}

	schedulerEnabled := getEnvOrDefault("SCHEDULER_ENABLED", "true") == "true"
	if schedulerEnabled {
		var schedulerService *svcscheduler.SchedulerService
		if sqsStoreInstance != nil && snsStoreInstance != nil && lambdaService != nil {
			schedulerService = svcscheduler.NewSchedulerServiceWithClients(server.StorageManager(), accountId, sqsStoreInstance, snsStoreInstance, lambdaService)
		} else if sqsStoreInstance != nil {
			schedulerService = svcscheduler.NewSchedulerServiceWithClients(server.StorageManager(), accountId, sqsStoreInstance, snsStoreInstance, nil)
		} else {
			schedulerService = svcscheduler.NewSchedulerService(server.StorageManager(), accountId, region)
		}
		schedulerService.RegisterHandlers(server.Dispatcher())
		schedulerService.StartEngine()
	}

	cloudtrailEnabled := getEnvOrDefault("CLOUDTRAIL_ENABLED", "true") == "true"
	if cloudtrailEnabled {
		cloudtrailService := svccloudtrail.NewCloudTrailService(server.Storage(), accountId, region)
		cloudtrailService.RegisterHandlers(server.Dispatcher())
	}

	timestreamWriteEnabled := getEnvOrDefault("TIMESTREAM_WRITE_ENABLED", "true") == "true"
	if timestreamWriteEnabled {
		serverHost := fmt.Sprintf("localhost:%s", port)
		timestreamWriteService := svctimestreamwrite.NewService(accountId, serverHost, dataPath)
		timestreamWriteService.RegisterHandlers(server.Dispatcher())
	}

	timestreamQueryEnabled := getEnvOrDefault("TIMESTREAM_QUERY_ENABLED", "true") == "true"
	if timestreamQueryEnabled {
		serverHost := fmt.Sprintf("localhost:%s", port)
		timestreamQueryService := svctimestreamquery.NewService(accountId, serverHost, dataPath)
		timestreamQueryService.RegisterHandlers(server.Dispatcher())
	}

	athenaEnabled := getEnvOrDefault("ATHENA_ENABLED", "true") == "true"
	if athenaEnabled {
		serverHost := fmt.Sprintf("localhost:%s", port)
		athenaService := svcathena.NewServiceWithS3(accountId, serverHost, s3ObjectStore)
		athenaService.RegisterHandlers(server.Dispatcher())
	}

	fmt.Printf("Starting AWS-compatible server on :%s\n", port)
	fmt.Printf("Data path: %s\n", dataPath)
	fmt.Printf("Account ID: %s\n", accountId)
	if useChainGateway {
		fmt.Printf("Gateway mode: Chain (metadata-based routing)\n")
	} else {
		fmt.Printf("Gateway mode: Standard (hardcoded routing)\n")
	}
	fmt.Printf("Services: IAM, STS, KMS, S3, Route53, CloudWatch, DynamoDB, ACM, CloudFront, WAF, SecretsManager, SNS, SQS, Events, Lambda, APIGateway, StepFunctions, Cognito, SSM, Logs, SESv2, Scheduler, Kinesis, CloudTrail, TimestreamWrite, TimestreamQuery, Athena\n")
	if dnsEnabled {
		fmt.Printf("Route53 DNS server: enabled on %s:53\n", dnsBindAddr)
	}
	if signatureEnabled {
		fmt.Printf("Signature verification: enabled\n")
	} else {
		fmt.Printf("Signature verification: disabled\n")
	}
	if accessKeyID == "" {
		fmt.Println("Warning: AWS_ACCESS_KEY_ID not set")
	}
	if secretAccessKey == "" {
		fmt.Println("Warning: AWS_SECRET_ACCESS_KEY not set")
	}

	if err := server.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
