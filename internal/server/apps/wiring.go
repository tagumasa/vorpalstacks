package apps

import (
	"context"
	"fmt"

	"vorpalstacks/internal/core/logs"
)

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
	}
	if st.neptuneGraphService != nil {
		eb.SetNeptuneGraphInvoker(&neptuneGraphInvokerAdapter{service: st.neptuneGraphService})
	}

	if st.cloudWatchService != nil {
		st.cloudWatchService.SetEventBus(eb)
		st.cloudWatchService.StartEvaluator(context.Background())
		a.addShutdown("cloudwatch", func(ctx context.Context) error {
			st.cloudWatchService.StopEvaluator()
			return nil
		})
	}

	if st.logsService != nil {
		st.logsService.SetEventBus(eb)
	}

	if st.cognitoService != nil {
		st.cognitoService.SetEventBus(eb)
	}

	if st.eventBridgeService != nil {
		st.eventBridgeService.SetEventBus(eb)
	}

	if st.lambdaService != nil {
		if st.s3ObjectStore != nil {
			st.lambdaService.SetS3ObjectStore(st.region, st.s3ObjectStore)
		}
		st.lambdaService.SetEventBus(eb)
		st.lambdaService.StartESMPoller(context.Background())
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

	if st.appSyncService != nil {
		st.appSyncService.SetEventBus(eb)
	}
}
