package grpcweb

import (
	"crypto/rsa"
	"fmt"
	"os"
	"path/filepath"

	chihttp "vorpalstacks/internal/server/http"
	"vorpalstacks/internal/store/config"
	"vorpalstacks/pkg/vsjwt"

	acmconnect "vorpalstacks/internal/pb/aws/acm/acmconnect"
	adminauthconnect "vorpalstacks/internal/pb/aws/admin_auth/admin_authconnect"
	adminconfigconnect "vorpalstacks/internal/pb/aws/admin_config/admin_configconnect"
	apigatewayconnect "vorpalstacks/internal/pb/aws/apigateway/apigatewayconnect"
	appsyncconnect "vorpalstacks/internal/pb/aws/appsync/appsyncconnect"
	athenaconnect "vorpalstacks/internal/pb/aws/athena/athenaconnect"
	cloudfrontconnect "vorpalstacks/internal/pb/aws/cloudfront/cloudfrontconnect"
	cloudtrailconnect "vorpalstacks/internal/pb/aws/cloudtrail/cloudtrailconnect"
	cloudwatchconnect "vorpalstacks/internal/pb/aws/cloudwatch/cloudwatchconnect"
	cloudwatcheventsconnect "vorpalstacks/internal/pb/aws/cloudwatchevents/cloudwatcheventsconnect"
	cloudwatchlogsconnect "vorpalstacks/internal/pb/aws/cloudwatchlogs/cloudwatchlogsconnect"
	cognitoidentityconnect "vorpalstacks/internal/pb/aws/cognitoidentity/cognitoidentityconnect"
	cognitoidentityproviderconnect "vorpalstacks/internal/pb/aws/cognitoidentityprovider/cognitoidentityproviderconnect"
	dynamodbconnect "vorpalstacks/internal/pb/aws/dynamodb/dynamodbconnect"
	iamconnect "vorpalstacks/internal/pb/aws/iam/iamconnect"
	kinesisconnect "vorpalstacks/internal/pb/aws/kinesis/kinesisconnect"
	kmsconnect "vorpalstacks/internal/pb/aws/kms/kmsconnect"
	lambdaconnect "vorpalstacks/internal/pb/aws/lambda/lambdaconnect"
	neptuneconnect "vorpalstacks/internal/pb/aws/neptune/neptuneconnect"
	neptunedataconnect "vorpalstacks/internal/pb/aws/neptunedata/neptunedataconnect"
	neptunegraphconnect "vorpalstacks/internal/pb/aws/neptunegraph/neptunegraphconnect"
	route53connect "vorpalstacks/internal/pb/aws/route53/route53connect"
	s3connect "vorpalstacks/internal/pb/aws/s3/s3connect"
	schedulerconnect "vorpalstacks/internal/pb/aws/scheduler/schedulerconnect"
	secretsmanagerconnect "vorpalstacks/internal/pb/aws/secretsmanager/secretsmanagerconnect"
	sesv2connect "vorpalstacks/internal/pb/aws/sesv2/sesv2connect"
	sfnconnect "vorpalstacks/internal/pb/aws/sfn/sfnconnect"
	snsconnect "vorpalstacks/internal/pb/aws/sns/snsconnect"
	sqsconnect "vorpalstacks/internal/pb/aws/sqs/sqsconnect"
	ssmconnect "vorpalstacks/internal/pb/aws/ssm/ssmconnect"
	stsconnect "vorpalstacks/internal/pb/aws/sts/stsconnect"
	timestreamqueryconnect "vorpalstacks/internal/pb/aws/timestreamquery/timestreamqueryconnect"
	timestreamwriteconnect "vorpalstacks/internal/pb/aws/timestreamwrite/timestreamwriteconnect"
	wafv2connect "vorpalstacks/internal/pb/aws/wafv2/wafv2connect"
	svcacm "vorpalstacks/internal/services/aws/acm"
	svcadminauth "vorpalstacks/internal/services/aws/admin_auth"
	svcadminconfig "vorpalstacks/internal/services/aws/admin_config"
	svcapigateway "vorpalstacks/internal/services/aws/apigateway"
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
	svclambda "vorpalstacks/internal/services/aws/lambda"
	svcneptune "vorpalstacks/internal/services/aws/neptune"
	svcneptunedata "vorpalstacks/internal/services/aws/neptunedata"
	svcneptuneGraph "vorpalstacks/internal/services/aws/neptunegraph"
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
	storecloudwatch "vorpalstacks/internal/store/aws/cloudwatch"
	storecognitoidentity "vorpalstacks/internal/store/aws/cognitoidentity"
	storecognitoidp "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	storekms "vorpalstacks/internal/store/aws/kms"
)

// AdminDeps bundles the dependencies required to register all admin console handlers.
type AdminDeps struct {
	Server            *chihttp.Server
	AccountID         string
	Region            string
	DataPath          string
	BaseURL           string
	NeptuneDataAdmin  *svcneptunedata.NeptuneDataService
	AppSyncAdmin      *svcappsync.AppSyncService
	LogsAdmin         *svclogs.LogsService
	EventsAdmin       *svcevents.EventsService
	SFNAdmin          *svcstepfunction.StepFunctionService
	SNSAdmin          *svcsns.SNSService
	SQSAdmin          *svcsqs.SQSService
	NeptuneAdmin      *svcneptune.NeptuneService
	NeptuneGraphAdmin *svcneptuneGraph.NeptuneGraphService
}

// RegisterAllAdminHandlers registers Connect RPC handlers for all services
// on the gRPC-Web admin console server.
func RegisterAllAdminHandlers(s *Server, deps AdminDeps) {
	st := deps.Server.Storage()
	sm := deps.Server.StorageManager()
	aid := deps.AccountID
	reg := deps.Region
	dp := deps.DataPath

	configStore := config.NewStore(st)
	adminConfigService := svcadminconfig.NewService(configStore)
	path, handler := adminconfigconnect.NewAdminConfigServiceHandler(adminConfigService)
	s.Handle(path, handler)

	s.Handle(acmconnect.NewACMServiceHandler(svcacm.NewAdminHandler()))
	s.Handle(apigatewayconnect.NewAPIGatewayServiceHandler(svcapigateway.NewAdminHandler()))
	s.Handle(cloudfrontconnect.NewCloudFrontServiceHandler(svccloudfront.NewAdminHandler()))
	s.Handle(route53connect.NewRoute53ServiceHandler(svcroute53.NewAdminHandler()))
	s.Handle(secretsmanagerconnect.NewSecretsManagerServiceHandler(svcsecretsmanager.NewAdminHandler()))
	s.Handle(stsconnect.NewSTSServiceHandler(svcsts.NewAdminHandler(aid)))
	s.Handle(wafv2connect.NewWAFV2ServiceHandler(svcwafv2.NewAdminHandler()))
	s.Handle(schedulerconnect.NewSchedulerServiceHandler(svcscheduler.NewAdminHandler()))
	s.Handle(dynamodbconnect.NewDynamoDBServiceHandler(svcdynamodb.NewAdminHandler(sm, aid)))
	s.Handle(kinesisconnect.NewKinesisServiceHandler(svckinesis.NewAdminHandler(sm, aid)))
	s.Handle(lambdaconnect.NewLambdaServiceHandler(svclambda.NewAdminHandler(sm, aid)))
	s.Handle(s3connect.NewS3ServiceHandler(svcs3.NewAdminHandler(sm, aid)))
	s.Handle(sfnconnect.NewSFNServiceHandler(svcstepfunction.NewAdminHandler(deps.SFNAdmin)))
	s.Handle(sqsconnect.NewSQSServiceHandler(svcsqs.NewAdminHandler(deps.SQSAdmin)))
	s.Handle(cloudwatcheventsconnect.NewCloudWatchEventsServiceHandler(svcevents.NewAdminHandler(deps.EventsAdmin, sm)))
	s.Handle(cloudwatchlogsconnect.NewCloudWatchLogsServiceHandler(svclogs.NewAdminHandler(deps.LogsAdmin)))
	s.Handle(snsconnect.NewSNSServiceHandler(svcsns.NewAdminHandler(deps.SNSAdmin)))
	s.Handle(iamconnect.NewIAMServiceHandler(svciam.NewAdminHandler(st, aid)))
	s.Handle(kmsconnect.NewKMSServiceHandler(svckms.NewAdminHandler(storekms.NewKeyStore(st, aid, reg))))
	s.Handle(cloudwatchconnect.NewCloudWatchServiceHandler(svccloudwatch.NewAdminHandler(
		storecloudwatch.NewAlarmStore(st, aid, reg),
		storecloudwatch.NewMetricChunkStore(st, reg, dp),
	)))
	s.Handle(cognitoidentityproviderconnect.NewCognitoIdentityProviderServiceHandler(svccognito.NewAdminHandler(
		storecognitoidp.NewCognitoStore(st, aid, reg),
	)))
	s.Handle(cognitoidentityconnect.NewCognitoIdentityServiceHandler(svccognitoidentity.NewAdminHandler(
		storecognitoidentity.NewCognitoIdentityStore(st, aid, reg),
	)))

	s.Handle(athenaconnect.NewAthenaServiceHandler(svcathena.NewAdminHandler(sm, aid)))
	s.Handle(cloudtrailconnect.NewCloudTrailServiceHandler(svccloudtrail.NewAdminHandler(sm, aid)))
	s.Handle(sesv2connect.NewSESv2ServiceHandler(svcsesv2.NewAdminHandler(sm, aid)))
	s.Handle(ssmconnect.NewSSMServiceHandler(svcssm.NewAdminHandler(sm, aid)))
	s.Handle(timestreamqueryconnect.NewTimestreamQueryServiceHandler(svctimestreamquery.NewAdminHandler(sm, aid, dp)))
	s.Handle(timestreamwriteconnect.NewTimestreamWriteServiceHandler(svctimestreamwrite.NewAdminHandler(sm, aid, dp)))

	s.Handle(neptuneconnect.NewNeptuneServiceHandler(svcneptune.NewAdminHandler(deps.NeptuneAdmin, aid)))
	s.Handle(neptunegraphconnect.NewNeptuneGraphServiceHandler(svcneptuneGraph.NewAdminHandler(deps.NeptuneGraphAdmin, aid)))
	if deps.NeptuneDataAdmin != nil {
		s.Handle(neptunedataconnect.NewNeptunedataServiceHandler(svcneptunedata.NewAdminHandler(deps.NeptuneDataAdmin)))
	}
	s.Handle(appsyncconnect.NewAppSyncServiceHandler(svcappsync.NewAdminHandler(deps.AppSyncAdmin, sm)))

	adminAuthKey, err := loadOrGenerateAdminAuthKey(dp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to initialise admin auth JWT key: %v\n", err)
	} else {
		jwtManager := vsjwt.NewManager(adminAuthKey, "admin-auth-key", "vorpalstacks/admin-auth")
		adminAuthService := svcadminauth.NewService(st, jwtManager, aid)
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
