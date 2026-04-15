package apps

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"vorpalstacks/internal/common/audit"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	"vorpalstacks/internal/server/grpcweb"
	"vorpalstacks/internal/server/listener"
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
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
	svcarn "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/pkg/graphengine"
)

// initOptionalServices initialises optional services. Returns an error on the
// first failure; the caller should invoke Shutdown to clean up services that
// were already initialised.
func (a *App) initOptionalServices() error {
	st := a.state

	initers := []struct {
		enabled bool
		name    string
		fn      func(*serviceState) error
	}{
		{a.cfg.CloudFront, "CloudFront", a.initCloudFront},
		{a.cfg.WAFv2, "WAFv2", a.initWAFv2},
		{a.cfg.Route53, "Route53", a.initRoute53},
		{a.cfg.Neptune, "Neptune", a.initNeptune},
		{a.cfg.NeptuneData, "NeptuneData", a.initNeptuneData},
		{a.cfg.NeptuneGraph, "NeptuneGraph", a.initNeptuneGraph},
		{a.cfg.AppSync, "AppSync", a.initAppSync},
		{a.cfg.TimestreamWrite, "TimestreamWrite", a.initTimestreamWrite},
		{a.cfg.TimestreamQuery, "TimestreamQuery", a.initTimestreamQuery},
		{a.cfg.Athena, "Athena", a.initAthena},
	}

	for _, init := range initers {
		if init.enabled {
			if err := init.fn(st); err != nil {
				return fmt.Errorf("failed to initialise %s: %w", init.name, err)
			}
		}
	}

	a.initCloudTrailRecorderFactory(st)
	return nil
}

// --- Athena (optional) ---

// initAthena creates the Athena service. Skipped with a warning if S3 is not
// enabled, because Athena requires S3 for query result storage.
func (a *App) initAthena(st *serviceState) error {
	if st.s3ObjectStore == nil {
		logs.Warn("Athena skipped: S3 not enabled")
		return nil
	}
	athenaService := svcathena.NewServiceWithS3(st.accountID, st.s3ObjectStore)
	athenaService.RegisterHandlers(a.server.Dispatcher())
	a.addShutdown("athena", func(ctx context.Context) error {
		athenaService.Shutdown()
		return nil
	})
	return nil
}

// --- AppSync (optional) ---

func (a *App) initAppSync(st *serviceState) error {
	st.appSyncService = svcappsync.NewAppSyncService(st.accountID)
	st.appSyncService.SetEventBus(a.server.EventBus())
	st.appSyncService.RegisterHandlers(a.server.Dispatcher())
	a.addShutdown("appsync", func(ctx context.Context) error {
		st.appSyncService.ShutdownEventServer()
		return nil
	})
	return nil
}

// --- CloudFront (optional) ---

func (a *App) initCloudFront(st *serviceState) error {
	st.cloudFrontService = svccloudfront.NewCloudFrontService(st.accountID)
	st.cloudFrontService.SetRegionAndStorage(st.region, a.server.StorageManager())
	st.cloudFrontService.RegisterHandlers(a.server.Dispatcher())
	st.cloudFrontService.InitDistributionServer()
	return nil
}

// --- Neptune (optional) ---

func (a *App) initNeptune(st *serviceState) error {
	st.neptuneService = svcneptune.NewNeptuneService(st.accountID, st.region)
	st.neptuneService.SetStorageManager(a.server.StorageManager())
	st.neptuneService.RegisterHandlers(a.server.Dispatcher())
	a.addShutdown("neptune", func(ctx context.Context) error {
		st.neptuneService.Close()
		return nil
	})
	return nil
}

// --- NeptuneData (optional) ---

func (a *App) initNeptuneData(st *serviceState) error {
	st.neptuneDataService = svcneptunedata.NewNeptuneDataService()
	st.neptuneDataService.SetStorageManager(a.server.StorageManager())
	st.neptuneDataService.RegisterHandlers(a.server.Dispatcher())
	a.addShutdown("neptunedata", func(ctx context.Context) error {
		st.neptuneDataService.Close()
		return nil
	})
	return nil
}

// --- Route53 (optional) ---

func (a *App) initRoute53(st *serviceState) error {
	route53Service, err := svcroute53.NewRoute53ServiceWithDNS(a.server.Storage(), st.accountID, a.cfg.Route53DNSBindAddr, a.cfg.Route53DNSEnabled)
	if err != nil {
		return fmt.Errorf("failed to create Route53 service: %w", err)
	}
	st.route53Service = route53Service
	st.route53Service.RegisterHandlers(a.server.Dispatcher())
	a.addShutdown("route53", func(ctx context.Context) error {
		return st.route53Service.Shutdown()
	})
	return nil
}

// --- TimestreamQuery (optional) ---

func (a *App) initTimestreamQuery(st *serviceState) error {
	timestreamQueryService := svctimestreamquery.NewService(st.accountID, a.cfg.ServerHost(), a.cfg.DataPath)
	timestreamQueryService.RegisterHandlers(a.server.Dispatcher())
	return nil
}

// --- TimestreamWrite (optional) ---

func (a *App) initTimestreamWrite(st *serviceState) error {
	timestreamWriteService := svctimestreamwrite.NewService(st.accountID, a.cfg.ServerHost(), a.cfg.DataPath)
	timestreamWriteService.RegisterHandlers(a.server.Dispatcher())
	st.timestreamWriteService = timestreamWriteService
	a.addShutdown("timestreamwrite", func(ctx context.Context) error {
		st.timestreamWriteService.Close()
		return nil
	})
	return nil
}

// --- WAFv2 (optional) ---

func (a *App) initWAFv2(st *serviceState) error {
	wafv2Service := svcwafv2.NewWAFv2Service(st.accountID, st.region)
	wafv2Service.RegisterHandlers(a.server.Dispatcher())
	return nil
}

// --- GraphDB ---

func (a *App) initGraphDB() {
	graphDir := a.cfg.DataPath + "/graph"
	graphDB, err := graphengine.Open(graphDir, graphengine.DefaultOptions())
	if err != nil {
		logs.Warn("Failed to open graph DB",
			logs.String("path", graphDir),
			logs.String("error", err.Error()),
		)
		return
	}
	a.server.Dispatcher().SetGraphDB(graphDB)
	a.graphDB = graphDB
	a.addShutdown("graphdb", func(ctx context.Context) error {
		graphDB.Close()
		return nil
	})
}

// --- NeptuneGraph (optional) ---

func (a *App) initNeptuneGraph(st *serviceState) error {
	graphCache := graphengine.NewSharedCache(graphengine.DefaultCacheSize)
	st.neptuneGraphService = svcneptuneGraph.NewNeptuneGraphService(st.accountID, st.region, a.cfg.DataPath)
	st.neptuneGraphService.SetStorageManager(a.server.StorageManager())
	st.neptuneGraphService.SetGraphCache(graphCache)
	st.neptuneGraphService.RegisterHandlers(a.server.Dispatcher())
	st.neptuneGraphService.RestoreEngines()
	a.addShutdown("neptunegraph", func(ctx context.Context) error {
		st.neptuneGraphService.Close()
		graphCache.Release()
		return nil
	})
	return nil
}

// --- gRPC-Web Admin ---

func (a *App) initGRPCWebAdmin() {
	st := a.state
	grpcWebServer := grpcweb.NewServer(&grpcweb.Config{Port: a.cfg.GRPCWebPort, BindAddr: a.cfg.GRPCWebBindAddr})

	aid := st.accountID
	reg := st.region
	dp := a.cfg.DataPath
	sm := a.server.StorageManager()

	handlers := make([]grpcweb.HandlerRegistration, 0, 32)

	var p string
	var h http.Handler

	p, h = svcacm.NewConnectHandler()
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcapigateway.NewConnectHandler()
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svccloudfront.NewConnectHandler()
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcroute53.NewConnectHandler()
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcsecretsmanager.NewConnectHandler()
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcsts.NewConnectHandler(aid)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcwafv2.NewConnectHandler()
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcscheduler.NewConnectHandler()
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcdynamodb.NewConnectHandler(sm, aid)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svckinesis.NewConnectHandler(sm, aid)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svclambda.NewConnectHandler(sm, aid)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcs3.NewConnectHandler(a.server.S3Store(), aid)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcstepfunction.NewConnectHandler(st.stepFunctionService)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcsqs.NewConnectHandler(st.sqsService)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcevents.NewConnectHandler(st.eventBridgeService, sm)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svclogs.NewConnectHandler(st.logsService)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcsns.NewConnectHandler(st.snsService)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svciam.NewConnectHandler(a.server.Storage(), aid)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svckms.NewConnectHandler(a.server.Storage(), aid, reg)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svccloudwatch.NewConnectHandler(a.server.Storage(), aid, reg, dp)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svccognito.NewConnectHandler(a.server.Storage(), aid, reg)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svccognitoidentity.NewConnectHandler(a.server.Storage(), aid, reg)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcathena.NewConnectHandler(sm, aid)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svccloudtrail.NewConnectHandler(sm, aid)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcsesv2.NewConnectHandler(sm, aid)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcssm.NewConnectHandler(sm, aid)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svctimestreamquery.NewConnectHandler(sm, aid, dp)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svctimestreamwrite.NewConnectHandler(sm, aid, dp)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcneptune.NewConnectHandler(st.neptuneService, aid)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcneptuneGraph.NewConnectHandler(st.neptuneGraphService, aid)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	if st.neptuneDataService != nil {
		p, h = svcneptunedata.NewConnectHandler(st.neptuneDataService)
		handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	}
	p, h = svcappsync.NewConnectHandler(st.appSyncService, sm)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})

	grpcweb.RegisterAdminHandlers(grpcWebServer, a.server.Storage(), aid, reg, dp, handlers)
	a.grpcWeb = grpcWebServer
	a.addShutdown("grpcweb", func(ctx context.Context) error {
		return grpcWebServer.Shutdown(ctx)
	})
}

// --- EventBus Resource Policy Functions ---

func (a *App) initEventBusPolicies() {
	st := a.state
	eb := a.server.EventBus()
	if eb == nil {
		return
	}

	if st.lambdaService != nil {
		eb.SetResourcePolicyFunc("lambda", eventbus.LambdaResourcePolicyFn(
			func(ctx context.Context, functionARN string) ([]eventbus.LambdaPolicyEntry, error) {
				fnName := svcarn.ExtractFunctionNameFromARN(functionARN)
				policies, err := st.lambdaService.GetFunctionPolicy(fnName)
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

	if st.sqsStoreInstance != nil {
		eb.SetResourcePolicyFunc("sqs", eventbus.SQSResourcePolicyFn(
			func(ctx context.Context, queueARN string) (string, error) {
				queueName := svcarn.ExtractQueueNameFromARN(queueARN)
				queue, err := st.sqsStoreInstance.GetQueueByName(queueName)
				if err != nil {
					return "", err
				}
				return queue.Policy, nil
			},
		))
	}

	if st.snsStoreInstance != nil {
		eb.SetResourcePolicyFunc("sns", eventbus.SNSTopicResourcePolicyFn(
			func(ctx context.Context, topicARN string) (string, error) {
				topic, err := st.snsStoreInstance.GetTopic(topicARN)
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

	if st.eventsStoreInstance != nil {
		eb.SetResourcePolicyFunc("events", eventbus.EventBridgeBusResourcePolicyFn(
			func(ctx context.Context, busName string) (string, error) {
				bus, err := st.eventsStoreInstance.GetEventBus(context.Background(), busName)
				if err != nil {
					return "", err
				}
				return bus.Policy, nil
			},
		))
	}
}

// --- Listeners ---

func (a *App) registerListeners() {
	st := a.state

	a.registerListener(listener.ListenerConfig{
		Name:        "s3_website",
		PortKey:     "ports.s3_website",
		DefaultPort: 8081,
		Handler: http.HandlerFunc(svcs3.NewWebsiteServer(
			a.server.S3Store(), st.accountID, st.region,
		).HandleRequest),
	})

	if a.cfg.CloudFront && st.cloudFrontService != nil {
		if handler := st.cloudFrontService.DistributionHandler(); handler != nil {
			a.registerListener(listener.ListenerConfig{
				Name:        "cloudfront",
				PortKey:     "ports.cloudfront",
				DefaultPort: 8084,
				Handler:     handler,
			})
		}
	}

	if st.lambdaService != nil {
		lambdaUrlServer := svclambda.NewFunctionURLServer(st.lambdaService, st.accountID, st.region, st.lambdaService)
		a.registerListener(listener.ListenerConfig{
			Name:        "lambda_url",
			PortKey:     "ports.lambda_url",
			DefaultPort: 8085,
			Handler:     http.HandlerFunc(lambdaUrlServer.HandleRequest),
		})
	}

	if st.apiGatewayService != nil && st.lambdaService != nil {
		st.apiGatewayService.InitRuntimeServer(
			a.server.StorageManager(),
			st.lambdaService,
			st.sqsStoreInstance,
			st.snsStoreInstance,
			a.server.EventBus(),
		)
		if handler := st.apiGatewayService.RuntimeHandler(); handler != nil {
			a.server.RegisterAPIGatewayRuntimeHandler(handler)
			a.registerListener(listener.ListenerConfig{
				Name:        "apigateway",
				PortKey:     "ports.apigateway",
				DefaultPort: 8082,
				Handler:     handler,
			})
			a.addShutdown("apigateway", func(ctx context.Context) error {
				st.apiGatewayService.CloseRuntimeServer()
				return nil
			})
		}
	}

	if st.cognitoService != nil {
		a.registerListener(listener.ListenerConfig{
			Name:        "cognito_hosted",
			PortKey:     "ports.cognito_hosted",
			DefaultPort: 8083,
			Handler:     http.HandlerFunc(st.cognitoService.HostedUIHandler),
		})
	}

	if a.cfg.AppSync && st.appSyncService != nil {
		a.registerListener(listener.ListenerConfig{
			Name:        "appsync_events",
			PortKey:     "ports.appsync_events",
			DefaultPort: 8086,
			Handler:     st.appSyncService.EventServerHandler(),
			Timeouts: &listener.ListenerTimeouts{
				ReadHeaderTimeout: 5 * time.Second,
				ReadTimeout:       0,
				WriteTimeout:      0,
				IdleTimeout:       0,
			},
		})
	}
}

// --- CloudTrail Audit Recorder ---

func (a *App) initCloudTrailRecorderFactory(st *serviceState) {
	a.server.Dispatcher().SetCloudTrailRecorderFactory(func(region, accountID string) request.AuditRecorder {
		regionalStorage, err := a.server.StorageManager().GetStorage(region)
		if err != nil {
			return nil
		}
		ctStore := cloudtrailstore.NewCloudTrailStore(regionalStorage, accountID, region)
		return audit.NewCloudTrailRecorder(ctStore)
	})
}

func (a *App) registerListener(cfg listener.ListenerConfig) {
	a.lm.Register(cfg)
}
