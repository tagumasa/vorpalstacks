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

func (a *App) wireCrossServiceDeps() error {
	st := a.state
	eb := a.server.EventBus()

	if eb == nil {
		return nil
	}

	if st.lambdaService != nil {
		eb.SetLambdaInvoker(st.lambdaService)
	}
	if st.snsStoreInstance != nil {
		var pub snsPublisher
		if st.snsService != nil {
			pub = st.snsService
		}
		eb.SetSNSInvoker(&snsInvokerAdapter{store: st.snsStoreInstance, kvStore: st.snsStoreInstance.BaseStore, publisher: pub})
	}
	if st.kinesisService != nil {
		eb.SetKinesisInvoker(&kinesisInvokerAdapter{
			provider:      st.kinesisService,
			defaultRegion: st.region,
		})
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

	if st.sqsService != nil {
		eb.SetSQSInvoker(&sqsInvokerAdapter{provider: st.sqsService})
	}
	if st.cloudTrailService != nil {
		eb.SetCloudTrailInvoker(&cloudTrailInvokerAdapter{provider: st.cloudTrailService})
	}
	if st.iamService != nil {
		st.iamService.SetCloudTrailInvoker(eb.CloudTrailInvoker())
	}

	// Register the CloudWatch invokers before the logs service reads them
	// back: its metric-filter evaluation publishes through the metric
	// invoker at PutLogEvents time.
	if st.cloudWatchService != nil {
		eb.SetCloudWatchMetricInvoker(&cloudWatchMetricInvokerAdapter{provider: st.cloudWatchService})
		eb.SetCloudWatchAlarmInvoker(&cloudWatchAlarmInvokerAdapter{provider: st.cloudWatchService})
	}

	if st.logsService != nil {
		st.logsService.SetCloudWatchMetricInvoker(eb.CloudWatchMetricInvoker())
		if err := st.logsService.SetEventBus(eb); err != nil {
			return err
		}
		eb.SetLogsInvoker(&logsInvokerAdapter{provider: st.logsService})
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
		if err := st.cognitoService.SetEventBus(eb); err != nil {
			return err
		}
		eb.SetCognitoTokenValidator(st.cognitoService)
	}

	if st.eventBridgeService != nil {
		if err := st.eventBridgeService.SetEventBus(eb); err != nil {
			return err
		}
	}

	// DynamoDB needs the bus for its S3 invoker: table exports upload to
	// S3 and table imports read from it through the event bus invoker.
	if st.dynamoDBService != nil {
		st.dynamoDBService.SetEventBus(eb)
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
		if err := st.s3Service.SetEventBus(eb); err != nil {
			return err
		}
	}

	if st.schedulerService != nil {
		st.schedulerService.BuildEngine()
		if err := st.schedulerService.SetEventBus(eb); err != nil {
			return err
		}
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
		if err := st.stepFunctionService.SetEventBus(eb); err != nil {
			return err
		}
		st.stepFunctionService.RecoverRunningExecutions()
	}

	if st.snsService != nil {
		if err := st.snsService.SetEventBus(eb); err != nil {
			return err
		}
	}

	return nil
}
