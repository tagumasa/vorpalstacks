package apps

import (
	"context"
	"fmt"
	"net/http"

	"vorpalstacks/internal/client/mobyclient"
	appconfig "vorpalstacks/internal/config"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	svcacm "vorpalstacks/internal/services/aws/acm"
	svcapigateway "vorpalstacks/internal/services/aws/apigateway"
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
	svcs3 "vorpalstacks/internal/services/aws/s3"
	svcscheduler "vorpalstacks/internal/services/aws/scheduler"
	svcsecretsmanager "vorpalstacks/internal/services/aws/secretsmanager"
	svcsesv2 "vorpalstacks/internal/services/aws/sesv2"
	svcstepfunction "vorpalstacks/internal/services/aws/sfn"
	svcsns "vorpalstacks/internal/services/aws/sns"
	svcsqs "vorpalstacks/internal/services/aws/sqs"
	svcssm "vorpalstacks/internal/services/aws/ssm"
	svcsts "vorpalstacks/internal/services/aws/sts"
	storeevents "vorpalstacks/internal/store/aws/eventbridge"
	storekinesis "vorpalstacks/internal/store/aws/kinesis"
	storesns "vorpalstacks/internal/store/aws/sns"
	storesqs "vorpalstacks/internal/store/aws/sqs"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// initAlwaysOnServices initialises all always-on services. Returns an error
// on the first initialisation failure; the caller is responsible for cleaning
// up any services that were already created.
func (a *App) initAlwaysOnServices() error {
	st := a.newServiceState()

	initers := []struct {
		enabled bool
		name    string
		fn      func(*serviceState) error
	}{
		{a.cfg.ACM, "ACM", a.initACM},
		{a.cfg.APIGateway, "APIGateway", a.initAPIGateway},
		{a.cfg.CloudTrail, "CloudTrail", a.initCloudTrail},
		{a.cfg.CloudWatch, "CloudWatch", a.initCloudWatch},
		{a.cfg.CloudWatchLogs, "CloudWatchLogs", a.initCloudWatchLogs},
		{a.cfg.Cognito, "Cognito", a.initCognito},
		{a.cfg.CognitoIdentity, "CognitoIdentity", a.initCognitoIdentity},
		{a.cfg.DynamoDB, "DynamoDB", a.initDynamoDB},
		{a.cfg.EventBridge, "EventBridge", a.initEventBridge},
		{a.cfg.IAM, "IAM", a.initIAM},
		{a.cfg.Kinesis, "Kinesis", a.initKinesis},
		{a.cfg.KMS, "KMS", a.initKMS},
		{a.cfg.Lambda, "Lambda", a.initLambda},
		{a.cfg.S3, "S3", a.initS3},
		{a.cfg.Scheduler, "Scheduler", a.initScheduler},
		{a.cfg.SecretsManager, "SecretsManager", a.initSecretsManager},
		{a.cfg.SESv2, "SESv2", a.initSESv2},
		{a.cfg.SFN, "SFN", a.initSFN},
		{a.cfg.SNS, "SNS", a.initSNS},
		{a.cfg.SQS, "SQS", a.initSQS},
		{a.cfg.SSM, "SSM", a.initSSM},
		{a.cfg.STS, "STS", a.initSTS},
	}

	for _, init := range initers {
		if init.enabled {
			if err := init.fn(st); err != nil {
				return fmt.Errorf("failed to initialise %s: %w", init.name, err)
			}
		}
	}

	a.state = st
	return nil
}

func (a *App) newServiceState() *serviceState {
	return &serviceState{
		accountID: a.cfg.AccountID,
		region:    a.cfg.Region,
	}
}

// --- ACM ---

func (a *App) initACM(st *serviceState) error {
	st.acmService = svcacm.NewACMService(st.accountID, st.region)
	st.acmService.RegisterHandlers(a.server.Dispatcher())
	return nil
}

// --- APIGateway ---

func (a *App) initAPIGateway(st *serviceState) error {
	st.apiGatewayService = svcapigateway.NewAPIGatewayService(st.accountID, st.region)
	st.apiGatewayService.RegisterHandlers(a.server.Dispatcher())
	return nil
}

// --- CloudTrail ---

func (a *App) initCloudTrail(st *serviceState) error {
	st.cloudTrailService = svccloudtrail.NewCloudTrailService(st.accountID, st.region)
	st.cloudTrailService.RegisterHandlers(a.server.Dispatcher())
	return nil
}

// --- CloudWatch ---

func (a *App) initCloudWatch(st *serviceState) error {
	st.cloudWatchService = svccloudwatch.NewCloudWatchService(st.accountID, st.region, a.cfg.DataPath)
	st.cloudWatchService.SetStorageManager(a.server.StorageManager())
	st.cloudWatchService.RegisterHandlers(a.server.Dispatcher())
	return nil
}

// --- CloudWatchLogs ---

func (a *App) initCloudWatchLogs(st *serviceState) error {
	st.logsService = svclogs.NewLogsService(a.server.StorageManager(), st.accountID, a.cfg.DataPath)
	st.logsService.SetEventBus(a.server.EventBus())
	st.logsService.RegisterHandlers(a.server.Dispatcher())
	a.addShutdown("cloudwatchlogs", func(ctx context.Context) error {
		st.logsService.Stop()
		return nil
	})
	return nil
}

// --- Cognito ---

func (a *App) initCognito(st *serviceState) error {
	st.cognitoService = svccognito.NewCognitoService(st.accountID, st.region)
	st.cognitoService.SetStorageManager(a.server.StorageManager())
	st.cognitoService.SetEventBus(a.server.EventBus())
	st.cognitoService.RegisterHandlers(a.server.Dispatcher())
	a.server.RegisterJWKSHandler(http.HandlerFunc(st.cognitoService.JWKSHandler))
	return nil
}

// --- CognitoIdentity ---

func (a *App) initCognitoIdentity(st *serviceState) error {
	st.cognitoIdentityService = svccognitoidentity.NewCognitoIdentityService(st.accountID, st.region)
	st.cognitoIdentityService.RegisterHandlers(a.server.Dispatcher())
	return nil
}

// --- DynamoDB ---

func (a *App) initDynamoDB(st *serviceState) error {
	st.dynamoDBService = svcdynamodb.NewDynamoDBService(st.accountID)
	st.dynamoDBService.RegisterHandlers(a.server.Dispatcher())
	return nil
}

// --- EventBridge ---

func (a *App) initEventBridge(st *serviceState) error {
	eventsRegionalStorage, err := a.server.StorageManager().GetStorage(st.region)
	if err != nil {
		return fmt.Errorf("failed to get EventBridge regional storage: %w", err)
	}
	st.eventsStoreInstance = storeevents.NewEventsStore(eventsRegionalStorage, st.accountID, st.region)
	st.eventBridgeService = svcevents.NewEventsService(a.server.StorageManager(), st.accountID)
	st.eventBridgeService.SetEventsStore(st.region, st.eventsStoreInstance)
	st.eventBridgeService.SetEventBus(a.server.EventBus())
	st.eventBridgeService.RegisterHandlers(a.server.Dispatcher())
	return nil
}

// --- IAM ---

func (a *App) initIAM(st *serviceState) error {
	st.iamService = svciam.NewIAMService(st.accountID)
	st.iamService.RegisterHandlers(a.server.Dispatcher())
	a.addShutdown("iam", func(ctx context.Context) error {
		st.iamService.WaitForReport()
		svciam.ShutdownSLRoleCleanup()
		return nil
	})

	if eb := a.server.EventBus(); eb != nil {
		if rr, ok := eb.RoleResolver().(*eventbus.IAMRoleResolver); ok {
			iamStore := a.server.IAMStore()
			rr.SetLookup(func(ctx context.Context, roleARN string) (string, error) {
				roleName := svcarn.ExtractRoleNameFromARN(roleARN)
				if roleName == "" {
					return "", fmt.Errorf("invalid role ARN format: %q", roleARN)
				}
				role, err := iamStore.Roles().Get(roleName)
				if err != nil {
					return "", err
				}
				return role.AssumeRolePolicyDocument, nil
			})
		}
	}
	return nil
}

// --- Kinesis ---

func (a *App) initKinesis(st *serviceState) error {
	kinesisRegionalStorage, err := a.server.StorageManager().GetStorage(st.region)
	if err != nil {
		return fmt.Errorf("failed to get Kinesis regional storage: %w", err)
	}
	tstore, ok := kinesisRegionalStorage.(storage.TransactionalStorageWith2PC)
	if ok {
		st.kinesisStoreInstance = storekinesis.NewKinesisStore(tstore, st.accountID, st.region)
	}
	st.kinesisService = svckinesis.NewKinesisService(st.accountID, st.region)
	if st.kinesisStoreInstance != nil {
		st.kinesisService.SetKinesisStore(st.region, st.kinesisStoreInstance)
	}
	st.kinesisService.RegisterHandlers(a.server.Dispatcher())
	return nil
}

// --- KMS ---

func (a *App) initKMS(st *serviceState) error {
	hsmBackend, err := svckmshsm.NewPersistentBackend(a.cfg.DataPath)
	if err != nil {
		return fmt.Errorf("failed to create persistent HSM backend: %w", err)
	}
	st.kmsService = svckms.NewKMSService(st.accountID, st.region, hsmBackend)
	st.kmsService.RegisterHandlers(a.server.Dispatcher())

	if err := st.kmsService.EnsureDefaultSSMKey(); err != nil {
		logs.Warn("Failed to create default SSM key",
			logs.String("error", err.Error()),
		)
	}
	return nil
}

// --- Lambda ---

func (a *App) initLambda(st *serviceState) error {
	dockerCfg := mobyclient.Config{
		Host:    a.cfg.DockerHost,
		Version: "1.44",
	}
	dockerLogger := logs.NewLogger(&logs.Config{Level: logs.LevelInfo})
	dockerClient, err := mobyclient.NewClient(dockerCfg, dockerLogger)
	if err != nil {
		logs.Warn("Failed to create Docker client for Lambda",
			logs.String("error", err.Error()),
		)
	}

	if dockerClient != nil {
		st.lambdaService = svclambda.NewLambdaService(dockerClient, st.accountID, st.region, a.cfg.DataPath)
		st.lambdaService.SetStorageManager(a.server.StorageManager())
		st.lambdaService.SetHostEndpoint(fmt.Sprintf("http://host.docker.internal:%s", a.cfg.Port))
		st.lambdaService.RegisterHandlers(a.server.Dispatcher())
		st.lambdaService.StartESMPoller(context.Background())

		a.addShutdown("lambda", func(ctx context.Context) error {
			st.lambdaService.Shutdown()
			return nil
		})
		a.addShutdown("lambda-esm", func(ctx context.Context) error {
			st.lambdaService.StopESMPoller()
			return nil
		})
	}
	return nil
}

// --- S3 ---

func (a *App) initS3(st *serviceState) error {
	s3Store := a.server.S3Store()
	blobStore := a.server.BlobStore()
	st.s3Service = svcs3.NewS3Service(s3Store, blobStore, st.accountID)
	st.s3Service.SetStorageManager(a.server.StorageManager())
	st.s3Service.RestoreSSE3Keys()
	s3Handler := svcs3.NewS3Handler(st.s3Service, st.region, a.server.StorageManager())
	a.server.RegisterS3Handler(s3Handler)
	st.s3Service.SetEventBus(a.server.EventBus())

	st.s3ObjectStore = s3Store.Objects(st.region)
	return nil
}

// --- Scheduler ---

func (a *App) initScheduler(st *serviceState) error {
	st.schedulerService = svcscheduler.NewSchedulerService(a.server.StorageManager(), st.accountID)
	st.schedulerService.RegisterHandlers(a.server.Dispatcher())
	return nil
}

// --- SecretsManager ---

func (a *App) initSecretsManager(st *serviceState) error {
	st.secretsManagerService = svcsecretsmanager.NewSecretsManagerService(st.accountID)
	st.secretsManagerService.SetRegion(st.region)
	st.secretsManagerService.SetStorageManager(a.server.StorageManager())
	st.secretsManagerService.RegisterHandlers(a.server.Dispatcher())
	st.secretsManagerService.StartRotationChecker(context.Background())
	a.addShutdown("secretsmanager", func(ctx context.Context) error {
		st.secretsManagerService.StopRotationChecker()
		return nil
	})
	return nil
}

// --- SESv2 ---

func (a *App) initSESv2(st *serviceState) error {
	st.sesv2Service = svcsesv2.NewSESv2Service(st.accountID)
	st.sesv2Service.RegisterHandlers(a.server.Dispatcher())
	return nil
}

// --- SFN ---

func (a *App) initSFN(st *serviceState) error {
	st.stepFunctionService = svcstepfunction.NewStepFunctionService(a.server.StorageManager(), st.accountID)
	st.stepFunctionService.SetEventBus(a.server.EventBus())
	st.stepFunctionService.RegisterHandlers(a.server.Dispatcher())
	a.addShutdown("sfn", func(ctx context.Context) error {
		st.stepFunctionService.Shutdown()
		return nil
	})
	return nil
}

// --- SNS ---

func (a *App) initSNS(st *serviceState) error {
	snsRegionalStorage, err := a.server.StorageManager().GetStorage(st.region)
	if err != nil {
		return fmt.Errorf("failed to get SNS regional storage: %w", err)
	}
	st.snsStoreInstance = storesns.NewSNSStore(snsRegionalStorage, st.accountID, st.region)
	st.snsService = svcsns.NewSNSService(a.server.StorageManager(), st.accountID, st.region)
	st.snsService.SetSNSStore(st.region, st.snsStoreInstance)
	st.snsService.SetEventBus(a.server.EventBus())
	st.snsService.RegisterHandlers(a.server.Dispatcher())
	a.addShutdown("sns", func(ctx context.Context) error {
		st.snsService.Close()
		st.snsStoreInstance.Close()
		return nil
	})
	return nil
}

// --- SQS ---

func (a *App) initSQS(st *serviceState) error {
	regionalStorage, err := a.server.StorageManager().GetStorage(st.region)
	if err != nil {
		return fmt.Errorf("failed to get SQS regional storage: %w", err)
	}
	st.sqsStoreInstance = storesqs.NewSQSStore(regionalStorage, st.accountID, st.region, appconfig.BaseURL())
	st.sqsService = svcsqs.NewSQSServiceWithStore(st.sqsStoreInstance)
	st.sqsService.SetStorageManager(a.server.StorageManager())
	st.sqsService.RegisterHandlers(a.server.Dispatcher())
	return nil
}

// --- SSM ---

func (a *App) initSSM(st *serviceState) error {
	st.ssmService = svcssm.NewSSMServiceWithKMS(st.accountID, st.kmsService)
	st.ssmService.RegisterHandlers(a.server.Dispatcher())
	return nil
}

// --- STS ---

func (a *App) initSTS(st *serviceState) error {
	st.stsService = svcsts.NewSTSService()
	st.stsService.RegisterHandlers(a.server.Dispatcher())
	return nil
}
