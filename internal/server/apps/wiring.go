package apps

import (
	"context"
	"fmt"
)

func (a *App) wireCrossServiceDeps() {
	st := a.state

	// --- CloudWatch ---
	if st.cloudWatchService != nil {
		st.cloudWatchService.SetEventBus(a.server.EventBus())
		if st.lambdaService != nil {
			st.cloudWatchService.SetLambdaInvoker(st.lambdaService)
		}
		st.cloudWatchService.StartEvaluator(context.Background())
		a.addShutdown("cloudwatch", func(ctx context.Context) error {
			st.cloudWatchService.StopEvaluator()
			return nil
		})
	}

	// --- CloudWatchLogs ---
	if st.logsService != nil {
		if st.lambdaService != nil {
			st.logsService.SetLambdaInvoker(st.region, st.lambdaService)
		}
		if st.kinesisStoreInstance != nil {
			st.logsService.SetKinesisStore(st.region, st.kinesisStoreInstance)
		}
	}

	// --- Cognito ---
	if st.cognitoService != nil {
		if st.lambdaService != nil {
			st.cognitoService.SetLambdaInvoker(st.lambdaService)
		}
	}

	// --- EventBridge ---
	if st.eventBridgeService != nil {
		if st.sqsStoreInstance != nil {
			st.eventBridgeService.SetSQSStore(st.region, st.sqsStoreInstance)
		}
		if st.snsStoreInstance != nil {
			st.eventBridgeService.SetSNSStore(st.region, st.snsStoreInstance)
		}
		if st.snsService != nil {
			st.eventBridgeService.SetSNSPublisher(st.snsService)
		}
		if st.lambdaService != nil {
			st.eventBridgeService.SetLambdaInvoker(st.lambdaService)
		}
	}

	// --- Lambda ---
	if st.lambdaService != nil {
		if st.s3ObjectStore != nil {
			st.lambdaService.SetS3ObjectStore(st.region, st.s3ObjectStore)
		}
		if st.sqsStoreInstance != nil {
			st.lambdaService.SetSQSStore(st.sqsStoreInstance)
		}
		if st.kinesisStoreInstance != nil {
			st.lambdaService.SetKinesisStore(st.kinesisStoreInstance)
		}
		st.lambdaService.SetEventBus(a.server.EventBus())
	}

	// --- S3 ---
	if st.s3Service != nil {
		if st.sqsStoreInstance != nil {
			st.s3Service.SetSQSStore(st.sqsStoreInstance)
		}
		if st.lambdaService != nil {
			st.s3Service.SetLambdaInvoker(st.lambdaService)
		}
	}

	// --- Scheduler ---
	if st.schedulerService != nil {
		if st.sqsStoreInstance != nil {
			st.schedulerService.SetSQSStore(st.sqsStoreInstance)
		}
		if st.snsStoreInstance != nil {
			st.schedulerService.SetSNSStore(st.snsStoreInstance)
		}
		if st.lambdaService != nil {
			st.schedulerService.SetLambdaInvoker(st.lambdaService)
		}
		st.schedulerService.BuildEngine()
		st.schedulerService.SetEventBus(a.server.EventBus())
		st.schedulerService.StartEngine()
		a.addShutdown("scheduler", func(ctx context.Context) error {
			if err := st.schedulerService.StopEngine(); err != nil {
				return fmt.Errorf("scheduler shutdown: %w", err)
			}
			return nil
		})
	}

	// --- SecretsManager ---
	if st.secretsManagerService != nil {
		if st.lambdaService != nil {
			st.secretsManagerService.SetLambdaInvoker(st.lambdaService)
		}
	}

	// --- SFN ---
	if st.stepFunctionService != nil {
		if st.lambdaService != nil {
			st.stepFunctionService.SetLambdaInvoker(st.lambdaService)
		}
		if st.sqsStoreInstance != nil {
			st.stepFunctionService.SetSQSStore(st.sqsStoreInstance)
		}
		if st.snsStoreInstance != nil {
			st.stepFunctionService.SetSNSStore(st.snsStoreInstance)
		}
		if st.eventsStoreInstance != nil {
			st.stepFunctionService.SetEventsStore(st.eventsStoreInstance)
		}
	}

	// --- SNS ---
	if st.snsService != nil {
		if st.sqsStoreInstance != nil {
			st.snsService.SetSQSStore(st.sqsStoreInstance)
		}
		if st.lambdaService != nil {
			st.snsService.SetLambdaInvoker(st.lambdaService)
		}
	}

	// --- AppSync ---
	if st.appSyncService != nil {
		if st.lambdaService != nil {
			st.appSyncService.SetLambdaInvoker(st.lambdaService)
		}
	}
}
