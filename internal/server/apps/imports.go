package apps

import (
	"vorpalstacks/internal/client/mobyclient"
	"vorpalstacks/internal/server/portalloc"
	svcacm "vorpalstacks/internal/services/aws/acm"
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
	svcec2 "vorpalstacks/internal/services/aws/ec2"
	svcevents "vorpalstacks/internal/services/aws/eventbridge"
	svciam "vorpalstacks/internal/services/aws/iam"
	svckinesis "vorpalstacks/internal/services/aws/kinesis"
	svckms "vorpalstacks/internal/services/aws/kms"
	svclambda "vorpalstacks/internal/services/aws/lambda"
	svcneptune "vorpalstacks/internal/services/aws/rds/neptune"
	svcneptunedata "vorpalstacks/internal/services/aws/rds/neptunedata"
	svcneptuneGraph "vorpalstacks/internal/services/aws/rds/neptunegraph"
	svcrdsdata "vorpalstacks/internal/services/aws/rds/rdsdata"
	svcvmysql "vorpalstacks/internal/services/aws/rds/vmysql"
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
	storeevents "vorpalstacks/internal/store/aws/eventbridge"
	storekinesis "vorpalstacks/internal/store/aws/kinesis"
	storesns "vorpalstacks/internal/store/aws/sns"
	storesqs "vorpalstacks/internal/store/aws/sqs"
)

type serviceState struct {
	accountID    string
	region       string
	dockerClient mobyclient.ContainerLifecycle

	acmService             *svcacm.ACMService
	apiGatewayService      *svcapigateway.APIGatewayService
	cloudFrontService      *svccloudfront.CloudFrontService
	cloudTrailService      *svccloudtrail.CloudTrailService
	cloudWatchService      *svccloudwatch.CloudWatchService
	cognitoService         *svccognito.CognitoService
	cognitoIdentityService *svccognitoidentity.CognitoIdentityService
	dynamoDBService        *svcdynamodb.DynamoDBService
	ec2Service             *svcec2.EC2Service
	eventBridgeService     *svcevents.EventsService
	iamService             *svciam.IAMService
	kinesisService         *svckinesis.KinesisService
	kmsService             *svckms.KMSService
	lambdaService          *svclambda.LambdaService
	neptuneService         *svcneptune.NeptuneService
	neptuneDataService     *svcneptunedata.NeptuneDataService
	neptuneGraphService    *svcneptuneGraph.NeptuneGraphService
	route53Service         *svcroute53.Route53Service
	s3Service              *svcs3.S3Service
	schedulerService       *svcscheduler.SchedulerService
	secretsManagerService  *svcsecretsmanager.SecretsManagerService
	sesv2Service           *svcsesv2.SESv2Service
	stepFunctionService    *svcstepfunction.StepFunctionService
	snsService             *svcsns.SNSService
	sqsService             *svcsqs.SQSService
	ssmService             *svcssm.SSMService
	stsService             *svcsts.STSService
	logsService            *svclogs.LogsService
	appSyncService         *svcappsync.AppSyncService
	timestreamWriteService *svctimestreamwrite.TimestreamWriteService
	timestreamQueryService *svctimestreamquery.TimestreamQueryService
	wafv2Service           *svcwafv2.WAFv2Service
	athenaService          *svcathena.AthenaService
	vmysqlService          *svcvmysql.Service
	rdsDataService         *svcrdsdata.RDSDataService
	portAllocator          *portalloc.Allocator

	sqsStoreInstance     *storesqs.SQSStore
	snsStoreInstance     *storesns.SNSStore
	kinesisStoreInstance *storekinesis.KinesisStore
	eventsStoreInstance  *storeevents.EventsStore
}
