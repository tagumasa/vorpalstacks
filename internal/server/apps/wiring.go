package apps

import (
	"context"
	"fmt"

	"vorpalstacks/internal/core/logs"
	svccognitoidentity "vorpalstacks/internal/services/aws/cognitoidentity"
	"vorpalstacks/internal/services/aws/dynamodb"
	stsstore "vorpalstacks/internal/store/aws/sts"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

// cognitoCredentialAdapter bridges the STS SessionStore to the
// CredentialIssuer interface required by CognitoIdentityService.
type cognitoCredentialAdapter struct {
	store stsstore.SessionStoreInterface
}

func (a *cognitoCredentialAdapter) IssueSession(roleArn, roleSessionName string, durationSeconds int) (*svccognitoidentity.CredentialResult, error) {
	session, err := a.store.Create(stsstore.CreateSessionParams{
		PrincipalType:   "WebIdentity",
		PrincipalName:   roleSessionName,
		RoleArn:         roleArn,
		RoleSessionName: roleSessionName,
		DurationSeconds: durationSeconds,
	})
	if err != nil {
		return nil, err
	}
	return &svccognitoidentity.CredentialResult{
		AccessKeyID:     session.AccessKeyId,
		SecretAccessKey: session.SecretAccessKey,
		SessionToken:    session.SessionToken,
		Expiration:      session.Expiration,
	}, nil
}

func (a *App) wireCrossServiceDeps() {
	st := a.state
	eb := a.server.EventBus()

	if eb == nil {
		return
	}

	if st.lambdaService != nil {
		eb.SetLambdaInvoker(st.lambdaService)
	}
	if st.sqsStoreInstance != nil {
		eb.SetSQSInvoker(&sqsInvokerAdapter{store: st.sqsStoreInstance})
	}
	if st.snsStoreInstance != nil {
		var pub snsPublisher
		if st.snsService != nil {
			pub = st.snsService
		}
		eb.SetSNSInvoker(&snsInvokerAdapter{store: st.snsStoreInstance, publisher: pub})
	}
	if st.kinesisStoreInstance != nil {
		eb.SetKinesisInvoker(&kinesisInvokerAdapter{store: st.kinesisStoreInstance})
	}
	if st.eventsStoreInstance != nil {
		eb.SetEventsInvoker(&eventsInvokerAdapter{putFn: st.eventsStoreInstance.Put})
	}
	if st.dynamoDBService != nil {
		eb.SetDynamoDBInvoker(&dynamoDBInvokerAdapter{provider: st.dynamoDBService})
		eb.SetDynamoDBStreamsInvoker(dynamodb.NewDynamoDBStreamsInvoker(st.dynamoDBService))
	}

	sm := a.server.StorageManager()
	globalStorage, err := sm.GetGlobalStorage()
	if err == nil {
		eb.SetWAFInvoker(&wafInvokerAdapter{
			store: wafstore.NewWebACLAssociationStore(globalStorage),
		})
	} else {
		logs.Warn("WAFInvoker not initialised: failed to get global storage", logs.Err(err))
	}
	eb.SetCloudWatchMetricInvoker(&cloudWatchMetricInvokerAdapter{
		storageMgr: sm,
		dataPath:   a.cfg.DataPath,
	})
	eb.SetCloudWatchAlarmInvoker(&cloudWatchAlarmInvokerAdapter{
		storageMgr: sm,
	})
	eb.SetTimestreamInvoker(&timestreamInvokerAdapter{
		storageMgr: sm,
		dataPath:   a.cfg.DataPath,
	})
	eb.SetCloudTrailInvoker(&cloudTrailInvokerAdapter{
		storageMgr: sm,
		accountID:  st.accountID,
	})
	eb.SetLogsInvoker(&logsInvokerAdapter{
		storageMgr: sm,
		accountID:  st.accountID,
		dataPath:   a.cfg.DataPath,
	})

	if st.iamService != nil {
		st.iamService.SetCloudTrailInvoker(eb.CloudTrailInvoker())
	}

	if st.logsService != nil {
		st.logsService.SetCloudWatchMetricInvoker(eb.CloudWatchMetricInvoker())
		st.logsService.SetEventBus(eb)
	}

	if st.cloudWatchService != nil {
		st.cloudWatchService.SetEventBus(eb)
		st.cloudWatchService.StartEvaluator(context.Background())
		a.addShutdown("cloudwatch", func(ctx context.Context) error {
			st.cloudWatchService.StopEvaluator()
			return nil
		})
	}

	if st.cognitoIdentityService != nil {
		if stsStore := a.server.STSSessionStore(); stsStore != nil {
			st.cognitoIdentityService.SetCredentialIssuer(&cognitoCredentialAdapter{store: stsStore})
		}
	}

	if st.cognitoService != nil {
		st.cognitoService.SetEventBus(eb)
	}

	if st.eventBridgeService != nil {
		st.eventBridgeService.SetEventBus(eb)
	}

	if st.lambdaService != nil {
		st.lambdaService.SetS3Invoker(eb.S3Invoker())
		st.lambdaService.SetLogsInvoker(eb.LogsInvoker())
		st.lambdaService.SetEventBus(eb)
		st.lambdaService.StartESMPoller(context.Background())
		eb.RegisterSubnetUsageChecker(st.lambdaService)
		eb.RegisterSecurityGroupUsageChecker(st.lambdaService)
		a.addShutdown("lambda-esm", func(ctx context.Context) error {
			st.lambdaService.StopESMPoller()
			return nil
		})
	}

	if st.s3Service != nil {
		st.s3Service.SetEventBus(eb)
	}

	if st.schedulerService != nil {
		st.schedulerService.BuildEngine()
		st.schedulerService.SetEventBus(eb)
		if err := st.schedulerService.StartEngine(); err != nil {
			logs.Warn("failed to start scheduler engine", logs.Err(err))
		}
		a.addShutdown("scheduler", func(ctx context.Context) error {
			if err := st.schedulerService.StopEngine(); err != nil {
				return fmt.Errorf("scheduler shutdown: %w", err)
			}
			return nil
		})
	}

	if st.secretsManagerService != nil {
		st.secretsManagerService.SetEventBus(eb)
	}

	if st.stepFunctionService != nil {
		st.stepFunctionService.SetEventBus(eb)
	}

	if st.snsService != nil {
		st.snsService.SetEventBus(eb)
	}
}
