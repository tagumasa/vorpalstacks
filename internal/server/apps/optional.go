package apps

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	appconfig "vorpalstacks/internal/config"
	"vorpalstacks/internal/eventbus"
	"vorpalstacks/internal/server/grpcweb"
	"vorpalstacks/internal/server/listener"
	svcappsync "vorpalstacks/internal/services/aws/appsync"
	svcathena "vorpalstacks/internal/services/aws/athena"
	svccloudfront "vorpalstacks/internal/services/aws/cloudfront"
	svclambda "vorpalstacks/internal/services/aws/lambda"
	svcneptune "vorpalstacks/internal/services/aws/neptune"
	svcneptunedata "vorpalstacks/internal/services/aws/neptunedata"
	svcneptuneGraph "vorpalstacks/internal/services/aws/neptunegraph"
	svcroute53 "vorpalstacks/internal/services/aws/route53"
	svcs3 "vorpalstacks/internal/services/aws/s3"
	svctimestreamquery "vorpalstacks/internal/services/aws/timestreamquery"
	svctimestreamwrite "vorpalstacks/internal/services/aws/timestreamwrite"
	svcwafv2 "vorpalstacks/internal/services/aws/wafv2"
	svcarn "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/pkg/graphengine"
)

func (a *App) initOptionalServices() {
	st := a.state

	if a.cfg.CloudFront {
		a.initCloudFront(st)
	}
	if a.cfg.WAFv2 {
		a.initWAFv2(st)
	}
	if a.cfg.Route53 {
		a.initRoute53(st)
	}
	if a.cfg.Neptune {
		a.initNeptune(st)
	}
	if a.cfg.NeptuneData {
		a.initNeptuneData(st)
	}
	if a.cfg.NeptuneGraph {
		a.initNeptuneGraph(st)
	}
	if a.cfg.AppSync {
		a.initAppSync(st)
	}
	if a.cfg.TimestreamWrite {
		a.initTimestreamWrite(st)
	}
	if a.cfg.TimestreamQuery {
		a.initTimestreamQuery(st)
	}
	if a.cfg.Athena {
		a.initAthena(st)
	}
}

// --- Athena (optional) ---

func (a *App) initAthena(st *serviceState) {
	athenaService := svcathena.NewServiceWithS3(st.accountID, st.s3ObjectStore)
	athenaService.RegisterHandlers(a.server.Dispatcher())
	a.addShutdown("athena", func(ctx context.Context) error {
		athenaService.Shutdown()
		return nil
	})
}

// --- AppSync (optional) ---

func (a *App) initAppSync(st *serviceState) {
	st.appSyncService = svcappsync.NewAppSyncService(st.accountID)
	st.appSyncService.SetEventBus(a.server.EventBus())
	st.appSyncService.RegisterHandlers(a.server.Dispatcher())
	a.addShutdown("appsync", func(ctx context.Context) error {
		st.appSyncService.ShutdownEventServer()
		return nil
	})
}

// --- CloudFront (optional) ---

func (a *App) initCloudFront(st *serviceState) {
	st.cloudFrontService = svccloudfront.NewCloudFrontService(st.accountID)
	st.cloudFrontService.SetRegionAndStorage(st.region, a.server.StorageManager())
	st.cloudFrontService.RegisterHandlers(a.server.Dispatcher())
	st.cloudFrontService.InitDistributionServer()
}

// --- Neptune (optional) ---

func (a *App) initNeptune(st *serviceState) {
	st.neptuneService = svcneptune.NewNeptuneService(st.accountID, st.region)
	st.neptuneService.SetStorageManager(a.server.StorageManager())
	st.neptuneService.RegisterHandlers(a.server.Dispatcher())
	a.addShutdown("neptune", func(ctx context.Context) error {
		st.neptuneService.Close()
		return nil
	})
}

// --- NeptuneData (optional) ---

func (a *App) initNeptuneData(st *serviceState) {
	st.neptuneDataService = svcneptunedata.NewNeptuneDataService()
	st.neptuneDataService.SetStorageManager(a.server.StorageManager())
	st.neptuneDataService.RegisterHandlers(a.server.Dispatcher())
	a.addShutdown("neptunedata", func(ctx context.Context) error {
		st.neptuneDataService.Close()
		return nil
	})
}

// --- NeptuneGraph (optional) ---

func (a *App) initNeptuneGraph(st *serviceState) {
	st.neptuneGraphService = svcneptuneGraph.NewNeptuneGraphService(st.accountID, st.region, a.cfg.DataPath)
	st.neptuneGraphService.SetStorageManager(a.server.StorageManager())
	st.neptuneGraphService.RegisterHandlers(a.server.Dispatcher())
	st.neptuneGraphService.RestoreEngines()
	a.addShutdown("neptunegraph", func(ctx context.Context) error {
		st.neptuneGraphService.Close()
		return nil
	})
}

// --- Route53 (optional) ---

func (a *App) initRoute53(st *serviceState) {
	route53Service, err := svcroute53.NewRoute53ServiceWithDNS(a.server.Storage(), st.accountID, a.cfg.Route53DNSBindAddr, a.cfg.Route53DNSEnabled)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create Route53 service: %v\n", err)
		os.Exit(1)
	}
	st.route53Service = route53Service
	st.route53Service.RegisterHandlers(a.server.Dispatcher())
	a.addShutdown("route53", func(ctx context.Context) error {
		return st.route53Service.Shutdown()
	})
}

// --- TimestreamQuery (optional) ---

func (a *App) initTimestreamQuery(st *serviceState) {
	timestreamQueryService := svctimestreamquery.NewService(st.accountID, a.cfg.ServerHost(), a.cfg.DataPath)
	timestreamQueryService.RegisterHandlers(a.server.Dispatcher())
}

// --- TimestreamWrite (optional) ---

func (a *App) initTimestreamWrite(st *serviceState) {
	timestreamWriteService := svctimestreamwrite.NewService(st.accountID, a.cfg.ServerHost(), a.cfg.DataPath)
	timestreamWriteService.RegisterHandlers(a.server.Dispatcher())
}

// --- WAFv2 (optional) ---

func (a *App) initWAFv2(st *serviceState) {
	wafv2Service := svcwafv2.NewWAFv2Service(st.accountID, st.region)
	wafv2Service.RegisterHandlers(a.server.Dispatcher())
}

// --- GraphDB ---

func (a *App) initGraphDB() {
	graphDir := a.cfg.DataPath + "/graph"
	graphDB, err := graphengine.Open(graphDir, graphengine.DefaultOptions())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to open graph DB at %s: %v\n", graphDir, err)
		return
	}
	a.server.Dispatcher().SetGraphDB(graphDB)
	a.graphDB = graphDB
	a.addShutdown("graphdb", func(ctx context.Context) error {
		graphDB.Close()
		return nil
	})
}

// --- gRPC-Web Admin ---

func (a *App) initGRPCWebAdmin() {
	st := a.state
	grpcWebServer := grpcweb.NewServer(&grpcweb.Config{Port: a.cfg.GRPCWebPort, BindAddr: a.cfg.GRPCWebBindAddr})
	grpcweb.RegisterAllAdminHandlers(grpcWebServer, grpcweb.AdminDeps{
		Server:            a.server,
		AccountID:         st.accountID,
		Region:            st.region,
		DataPath:          a.cfg.DataPath,
		BaseURL:           appconfig.BaseURL(),
		NeptuneDataAdmin:  st.neptuneDataService,
		NeptuneAdmin:      st.neptuneService,
		NeptuneGraphAdmin: st.neptuneGraphService,
		AppSyncAdmin:      st.appSyncService,
		LogsAdmin:         st.logsService,
		EventsAdmin:       st.eventBridgeService,
		SFNAdmin:          st.stepFunctionService,
		SNSAdmin:          st.snsService,
		SQSAdmin:          st.sqsService,
	})
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
			a.server.StorageManager(), a.server.BlobStore(), st.accountID, st.region,
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
		lambdaUrlServer := svclambda.NewFunctionURLServer(a.server.StorageManager(), st.accountID, st.region, st.lambdaService)
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

func (a *App) registerListener(cfg listener.ListenerConfig) {
	a.lm.Register(cfg)
}
