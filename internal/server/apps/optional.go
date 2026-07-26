package apps

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"time"
	"vorpalstacks/internal/core/logs"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/sql"
	"vorpalstacks/internal/common/audit"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/serviceports"
	appconfig "vorpalstacks/internal/config"
	"vorpalstacks/internal/core/storage/graphengine"
	"vorpalstacks/internal/eventbus"
	"vorpalstacks/internal/server/grpcweb"
	"vorpalstacks/internal/server/listener"
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
	svcevents "vorpalstacks/internal/services/aws/eventbridge"
	svciam "vorpalstacks/internal/services/aws/iam"
	svciot "vorpalstacks/internal/services/aws/iot"
	"vorpalstacks/internal/services/aws/iot/broker"
	"vorpalstacks/internal/services/aws/iot/ca"
	"vorpalstacks/internal/services/aws/iot/rules/actions"
	svciotevents "vorpalstacks/internal/services/aws/iotevents"
	svckinesis "vorpalstacks/internal/services/aws/kinesis"
	svckms "vorpalstacks/internal/services/aws/kms"
	svclambda "vorpalstacks/internal/services/aws/lambda"
	svcrds "vorpalstacks/internal/services/aws/rds"
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
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
	storecommon "vorpalstacks/internal/store/aws/common"
	iamstore "vorpalstacks/internal/store/aws/iam"
	iotstore "vorpalstacks/internal/store/aws/iot"
	storerds "vorpalstacks/internal/store/aws/rds"
	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// initOptionalServices initialises optional services. Returns an error on the
// first failure; the caller should invoke Shutdown to clean up services that
// were already initialised.
func (a *App) initOptionalServices() error {
	st := a.state
	st.portAllocator = portalloc.New(appconfig.GetStore())

	// Initialise the standalone RDS store lookup before any service init.
	// This is used by the RDS admin handler and the RDS Data API and is
	// always available, even when Neptune is disabled.
	a.initRDSStoreLookup(st)

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
		{a.cfg.RDSMySQL, "RDSMySQL", a.initRDSMySQL},
		{a.cfg.RDSMySQL, "RDSData", a.initRDSData},
		{a.cfg.IoT, "IoT", a.initIoT},
		{a.cfg.IoTEvents, "IoTEvents", a.initIoTEvents},
	}

	for _, init := range initers {
		if init.enabled {
			if err := init.fn(st); err != nil {
				return fmt.Errorf("failed to initialise %s: %w", init.name, err)
			}
		}
	}

	if a.cfg.CloudTrail {
		a.initCloudTrailRecorderFactory(st)
		a.injectS3AuditRecorder(st)
	}
	a.initPrincipalResolver()

	if st.neptuneService != nil && st.neptuneDataService != nil {
		st.neptuneService.SetEngine(st.neptuneDataService)
		st.neptuneService.SetClusterPorter(st.neptuneDataService)
	}

	// Wire the MySQL (vmysql) engine and snapshot operator into the
	// Neptune SDK handler so it can serve RDS SDK clients that create
	// MySQL instances/clusters. The Neptune handler dispatches to the
	// MySQL engine when Engine is "mysql" (see NeptuneService.engineFor),
	// and captures / restores row-level snapshot data through the
	// snapshot operator.
	if st.neptuneService != nil && st.vmysqlService != nil {
		st.neptuneService.SetMysqlEngine(st.vmysqlService)
		st.neptuneService.SetSnapshotOperator(&vmysqlSnapshotOperatorAdapter{svc: st.vmysqlService})
	}

	if st.kmsService != nil {
		st.kmsService.SetPrincipalResolver(a.server.Dispatcher().PrincipalResolver())
	}

	return nil
}

// --- Athena (optional) ---

// initAthena creates the Athena service. S3 invoker is injected from the
// event bus immediately after creation.
func (a *App) initAthena(st *serviceState) error {
	athenaService := svcathena.NewAthenaService(st.accountID)
	athenaService.SetRegion(st.region)
	athenaService.SetStorageManager(a.server.StorageManager())
	st.athenaService = athenaService
	if eb := a.server.EventBus(); eb != nil {
		st.athenaService.SetS3Invoker(eb.S3Invoker())
	}
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
	st.appSyncService.SetStorageManager(a.server.StorageManager())
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
	if eb := a.server.EventBus(); eb != nil {
		st.cloudFrontService.SetWAFInvoker(eb.WAFInvoker())
	}
	st.cloudFrontService.RegisterHandlers(a.server.Dispatcher())
	st.cloudFrontService.InitDistributionServer()
	return nil
}

// --- Neptune (optional) ---

func (a *App) initNeptune(st *serviceState) error {
	st.neptuneService = svcneptune.NewNeptuneService(st.accountID, st.region)
	st.neptuneService.SetStorageManager(a.server.StorageManager())
	st.neptuneService.SetEventBus(a.server.EventBus())
	st.neptuneService.SetServerHost(a.cfg.ServerHost())
	st.neptuneService.RegisterHandlers(a.server.Dispatcher())
	if eb := a.server.EventBus(); eb != nil {
		eb.RegisterSubnetUsageChecker(st.neptuneService)
	}
	a.addShutdown("neptune", func(ctx context.Context) error {
		st.neptuneService.Close()
		return nil
	})
	return nil
}

// --- NeptuneData (optional) ---

func (a *App) initNeptuneData(st *serviceState) error {
	graphCache := graphengine.NewSharedCache(graphengine.DefaultCacheSize)
	st.neptuneDataService = svcneptunedata.NewNeptuneDataService(st.portAllocator)
	st.neptuneDataService.SetStorageManager(a.server.StorageManager())
	st.neptuneDataService.SetGraphCache(graphCache)
	st.neptuneDataService.SetListenerManager(a.lm)
	st.neptuneDataService.SetDispatcherHandler(a.server.MainHandler)
	if eb := a.server.EventBus(); eb != nil {
		st.neptuneDataService.SetS3Invoker(eb.S3Invoker())
	}
	st.neptuneDataService.RegisterHandlers(a.server.Dispatcher())
	st.neptuneDataService.RestoreEngines()
	st.neptuneDataService.RegisterClusterListeners()

	a.addShutdown("neptunedata", func(ctx context.Context) error {
		st.neptuneDataService.Shutdown()
		graphCache.Release()
		return nil
	})
	return nil
}

// --- Route53 (optional) ---

func (a *App) initRoute53(st *serviceState) error {
	dnsPort := a.resolvedPort("ports.route53_dns", serviceports.Route53DNS)
	route53Service, err := svcroute53.NewRoute53ServiceWithDNS(a.server.Storage(), st.accountID, "0.0.0.0", dnsPort, a.cfg.Route53DNSEnabled)
	if err != nil {
		return fmt.Errorf("failed to create Route53 service: %w", err)
	}
	if eb := a.server.EventBus(); eb != nil {
		route53Service.SetEventBus(eb)
	}
	hcPort := a.resolvedPort("ports.route53_healthcheck", serviceports.Route53HC)
	if hcPort != serviceports.Route53HC {
		route53Service.SetDefaultHCPort(hcPort)
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
	timestreamQueryService := svctimestreamquery.NewTimestreamQueryService(st.accountID, a.cfg.ServerHost(), a.cfg.DataPath)
	timestreamQueryService.SetStorageManager(a.server.StorageManager())
	st.timestreamQueryService = timestreamQueryService
	timestreamQueryService.RegisterHandlers(a.server.Dispatcher())
	return nil
}

// --- TimestreamWrite (optional) ---

func (a *App) initTimestreamWrite(st *serviceState) error {
	timestreamWriteService := svctimestreamwrite.NewTimestreamWriteService(st.accountID, a.cfg.ServerHost(), a.cfg.DataPath)
	timestreamWriteService.SetStorageManager(a.server.StorageManager())
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
	st.wafv2Service = svcwafv2.NewWAFv2Service(st.accountID, st.region)
	st.wafv2Service.SetStorageManager(a.server.StorageManager())
	st.wafv2Service.RegisterHandlers(a.server.Dispatcher())
	return nil
}

// --- NeptuneGraph (optional) ---

func (a *App) initNeptuneGraph(st *serviceState) error {
	graphCache := graphengine.NewSharedCache(graphengine.DefaultCacheSize)
	st.neptuneGraphService = svcneptuneGraph.NewNeptuneGraphService(st.accountID, st.region, a.cfg.DataPath)
	st.neptuneGraphService.SetStorageManager(a.server.StorageManager())
	st.neptuneGraphService.SetGraphCache(graphCache)
	st.neptuneGraphService.SetEventBus(a.server.EventBus())
	st.neptuneGraphService.RegisterHandlers(a.server.Dispatcher())
	st.neptuneGraphService.RestoreEngines()
	if eb := a.server.EventBus(); eb != nil {
		eb.SetNeptuneGraphInvoker(&neptuneGraphInvokerAdapter{service: st.neptuneGraphService})
	}
	a.addShutdown("neptunegraph", func(ctx context.Context) error {
		st.neptuneGraphService.Close()
		graphCache.Release()
		return nil
	})
	return nil
}

// --- RDS MySQL (optional) ---

// initRDSStoreLookup creates a per-region RDS store factory that is
// independent of the Neptune service. It uses the same Pebble bucket
// configuration as the Neptune store (neptune_clusters, neptune_instances,
// etc.), so data written via the RDS admin API is visible to all consumers
// regardless of which services are enabled. This eliminates the
// Neptune-dependency that previously caused cluster ARN resolution to fail
// with InternalServerError when Neptune was disabled (M-2) and the nil
// dereference in the RDS admin handler under the same condition (L-5).
func (a *App) initRDSStoreLookup(st *serviceState) {
	sm := a.server.StorageManager()
	var cache sync.Map
	st.rdsStoreLookup = func(region string) (storerds.StoreInterface, error) {
		return storecommon.GetOrCreateStoreE(&cache, region, func() (storerds.StoreInterface, error) {
			rs, err := sm.GetStorage(region)
			if err != nil {
				return nil, fmt.Errorf("failed to get regional storage for RDS store: %w", err)
			}
			return storerds.NewRDSStore(rs, storerds.BucketConfig{
				ClusterBucket:           "neptune_clusters",
				InstanceBucket:          "neptune_instances",
				SnapshotBucket:          "neptune_snapshots",
				InstSnapshotBucket:      "neptune_inst_snapshots",
				ClusterParamGroupBucket: "neptune_cluster_param_groups",
				ParamGroupBucket:        "neptune_param_groups",
				SubnetGroupBucket:       "neptune_subnet_groups",
				GlobalClusterBucket:     "neptune_global_clusters",
				EventSubBucket:          "neptune_event_subscriptions",
				EventsBucket:            "neptune_events",
				TagsBucket:              "neptune_tags",
				Namespace:               "neptune",
			}), nil
		})
	}
}

func (a *App) initRDSMySQL(st *serviceState) error {
	svc := svcvmysql.NewService(st.portAllocator)
	svc.SetStorageManager(a.server.StorageManager())
	st.vmysqlService = svc

	// The hardcoded "test-instance" engine exists to support the rdsdata
	// SDK test suite, which references arn:...:db:test-instance directly.
	// In any deployment that is not running the SDK test suite the
	// engine would leak a Pebble bucket and a TCP listener for no
	// benefit, so gate its creation on TEST_MODE.
	if os.Getenv("TEST_MODE") == "true" || os.Getenv("TEST_MODE") == "1" {
		if _, err := svc.Open(st.region, "test-instance"); err != nil {
			return fmt.Errorf("failed to open test MySQL instance: %w", err)
		}
	}

	a.addShutdown("rds-mysql", func(ctx context.Context) error {
		svc.Shutdown()
		return nil
	})
	return nil
}

// --- RDS Data API (optional, gated on RDSMySQL) ---

func (a *App) initRDSData(st *serviceState) error {
	if st.vmysqlService == nil {
		return nil // vmysql not initialised, skip
	}

	rdsDataSvc := svcrdsdata.NewRDSDataService()
	rdsDataSvc.SetVmysqlProvider(&vmysqlEngineProvider{svc: st.vmysqlService})

	// Cluster ARN resolution uses the standalone RDS store lookup that
	// is always available (initialised in initRDSStoreLookup), regardless
	// of whether Neptune is enabled. The store reads/writes the same
	// Pebble buckets as the Neptune store, ensuring data compatibility.
	rdsDataSvc.SetRDSStoreProvider(&rdsStoreProvider{
		region: st.region,
		lookup: st.rdsStoreLookup,
	})

	// Secrets Manager validation: when secretsmanager is enabled, wire it in
	// so secretArn is actually looked up rather than accepted verbatim.
	if st.secretsManagerService != nil {
		secretLookup := st.secretsManagerService.GetStoreForRegion
		rdsDataSvc.SetSecretResolver(&rdsDataSecretResolver{
			region: st.region,
			lookup: func(region string) (secretsmanagerstore.SecretStoreInterface, error) {
				return secretLookup(region)
			},
		})
	}

	if eb := a.server.EventBus(); eb != nil {
		rdsDataSvc.SetEventBus(eb)
		eb.SetRDSDataInvoker(&rdsDataInvokerAdapter{service: rdsDataSvc})
	}
	rdsDataSvc.RegisterHandlers(a.server.Dispatcher())
	rdsDataSvc.StartCleanup()
	st.rdsDataService = rdsDataSvc

	a.addShutdown("rdsdata", func(ctx context.Context) error {
		rdsDataSvc.Shutdown()
		return nil
	})
	return nil
}

// vmysqlEngineProvider adapts vmysql.Service to the rdsdata.VmysqlProvider interface.
type vmysqlEngineProvider struct {
	svc *svcvmysql.Service
}

func (p *vmysqlEngineProvider) GetEngine(instanceID string) *sqle.Engine {
	return p.svc.GetEngine(instanceID)
}

func (p *vmysqlEngineProvider) NewContext(instanceID string, database string) *sql.Context {
	return p.svc.NewContext(instanceID, database)
}

// vmysqlSnapshotOperatorAdapter adapts vmysql.Service to the
// SnapshotOperator interface so CreateDBSnapshot, DeleteDBInstance
// final snapshots, RestoreDBInstanceFromDBSnapshot, and DeleteDBSnapshot
// can capture / restore / clean up row data.
type vmysqlSnapshotOperatorAdapter struct {
	svc *svcvmysql.Service
}

func (a *vmysqlSnapshotOperatorAdapter) SnapshotData(instanceID, snapshotID string) error {
	if a == nil || a.svc == nil {
		return fmt.Errorf("vmysql snapshot operator: service is nil")
	}
	return a.svc.SnapshotData(instanceID, snapshotID)
}

func (a *vmysqlSnapshotOperatorAdapter) RestoreData(snapshotID, instanceID string) error {
	if a == nil || a.svc == nil {
		return fmt.Errorf("vmysql snapshot operator: service is nil")
	}
	return a.svc.RestoreData(snapshotID, instanceID)
}

func (a *vmysqlSnapshotOperatorAdapter) DeleteSnapshotData(snapshotID string) error {
	if a == nil || a.svc == nil {
		return fmt.Errorf("vmysql snapshot operator: service is nil")
	}
	return a.svc.DeleteSnapshotData(snapshotID)
}

// rdsStoreProvider resolves RDS resource ARNs to live DB instance IDs by
// consulting the unified RDS store. It implements the
// svcrdsdata.RDSStoreProvider interface and replaces the previous stub that
// returned the cluster ID verbatim.
//
// For cluster ARNs the provider lists all DB instances whose
// DBClusterIdentifier matches and returns the first writer instance. AWS
// Data API semantics require the writer (not a reader) for SQL execution.
type rdsStoreProvider struct {
	region string
	lookup func(region string) (storerds.StoreInterface, error)
}

func (p *rdsStoreProvider) GetInstance(identifier string) (string, string, string, error) {
	if p == nil || p.lookup == nil {
		return "", "", "", fmt.Errorf("rds store not configured")
	}
	store, err := p.lookup(p.region)
	if err != nil {
		return "", "", "", err
	}
	inst, err := store.GetInstance(identifier)
	if err != nil {
		return "", "", "", err
	}
	return inst.DBInstanceIdentifier, inst.Engine, inst.DBInstanceStatus, nil
}

func (p *rdsStoreProvider) GetCluster(identifier string) (string, []string, error) {
	if p == nil || p.lookup == nil {
		return "", nil, fmt.Errorf("rds store not configured")
	}
	store, err := p.lookup(p.region)
	if err != nil {
		return "", nil, err
	}
	if _, err := store.GetCluster(identifier); err != nil {
		return "", nil, err
	}
	// Find every DB instance whose DBClusterIdentifier matches the requested
	// cluster, then prefer the writer. Aurora exposes the writer via the
	// instance endpoint mode; our store does not yet model reader/writer
	// distinction explicitly, so the first member is returned. This matches
	// the documented Data API behaviour of routing to a writer instance and
	// is sufficient for the single-writer topology used by the local
	// vmysql engine.
	instances, err := store.ListInstances()
	if err != nil {
		return "", nil, err
	}
	var members []string
	for _, inst := range instances {
		if inst.DBClusterIdentifier == identifier {
			members = append(members, inst.DBInstanceIdentifier)
		}
	}
	if len(members) == 0 {
		return "", nil, fmt.Errorf("cluster %s has no instances", identifier)
	}
	return identifier, members, nil
}

// rdsDataSecretResolver implements svcrdsdata.SecretResolver by delegating
// to the Secrets Manager service. The secret's SecretString is parsed as a
// JSON credential document; failures map to the AWS exception types
// (SecretsErrorException for retrieval, InvalidSecretException for malformed
// credentials).
type rdsDataSecretResolver struct {
	region string
	lookup func(region string) (secretsmanagerstore.SecretStoreInterface, error)
}

func (r *rdsDataSecretResolver) ResolveSecret(ctx context.Context, secretArn string) (*svcrdsdata.SecretCredential, error) {
	if r == nil || r.lookup == nil {
		return nil, svcrdsdata.SecretsError("secrets manager not configured")
	}
	store, err := r.lookup(r.region)
	if err != nil {
		return nil, svcrdsdata.SecretsError(fmt.Sprintf("failed to load secrets store for region %s: %v", r.region, err))
	}
	secret, err := store.GetSecretByARN(secretArn)
	if err != nil {
		// Could be NotFound, deleted, or service error — collapse to
		// SecretsErrorException, matching AWS which uses the same exception
		// name for "secret wasn't found / couldn't be decrypted / timed out".
		return nil, svcrdsdata.SecretsError(fmt.Sprintf("secret %s could not be retrieved: %v", secretArn, err))
	}
	if secret == nil || secret.SecretString == "" {
		return nil, svcrdsdata.InvalidSecret(fmt.Sprintf("secret %s has no SecretString", secretArn))
	}
	var cred svcrdsdata.SecretCredential
	if err := json.Unmarshal([]byte(secret.SecretString), &cred); err != nil {
		return nil, svcrdsdata.InvalidSecret(fmt.Sprintf("secret %s is not valid JSON: %v", secretArn, err))
	}
	return &cred, nil
}

// --- IoT (optional) ---

func (a *App) initIoT(st *serviceState) error {
	eb := a.server.EventBus()

	regions := a.server.StorageManager().ListRegions()
	cas := make(map[string]*ca.CertificateAuthority, len(regions))
	brokers := make(map[string]*broker.Broker, len(regions))
	dispatchers := make(map[string]*actions.Dispatcher, len(regions))

	for _, region := range regions {
		regionalStorage, err := a.server.StorageManager().GetStorage(region)
		if err != nil {
			for _, b := range brokers {
				b.Stop()
			}
			return fmt.Errorf("failed to get regional storage for IoT CA in %s: %w", region, err)
		}
		certAuth, err := ca.NewCertificateAuthority(regionalStorage)
		if err != nil {
			for _, b := range brokers {
				b.Stop()
			}
			return fmt.Errorf("failed to create IoT CA in %s: %w", region, err)
		}
		cas[region] = certAuth

		mqttBroker := broker.NewBroker(certAuth, nil)
		port, err := st.portAllocator.Get("iot-mqtt", region)
		if err != nil {
			for _, b := range brokers {
				b.Stop()
			}
			return fmt.Errorf("failed to allocate MQTT port in %s: %w", region, err)
		}
		if err := mqttBroker.Start(port); err != nil {
			// In TEST_MODE the port may already be in use by a previous test
			// run; tolerate that per region so the rest of IoT can still
			// initialise. The failed broker is simply not registered.
			if os.Getenv("TEST_MODE") == "true" {
				logs.Warn("iot: MQTT broker failed to start in TEST_MODE (continuing)", logs.String("region", region), logs.Err(err))
				continue
			}
			for _, b := range brokers {
				b.Stop()
			}
			return fmt.Errorf("failed to start MQTT broker on port %d in %s: %w", port, region, err)
		}
		brokers[region] = mqttBroker

		if eb != nil {
			dispatchers[region] = actions.NewDispatcher(eb, nil)
		}
	}

	iotService := svciot.NewIoTService(st.accountID)
	iotService.Init(svciot.IoTServiceDeps{
		StorageManager: a.server.StorageManager(),
		CAs:            cas,
		Brokers:        brokers,
		EventBus:       eb,
		Dispatchers:    dispatchers,
	})

	// Each regional broker authenticates against its own region's store
	// and CA, so a certificate created in one region is never accepted by
	// another region's broker.
	for region, mqttBroker := range brokers {
		mqttBroker.SetAuthProvider(&svciot.IoTAuthBrokerProvider{Service: iotService, Region: region})
	}

	st.iotService = iotService
	st.iotService.RegisterHandlers(a.server.Dispatcher())
	a.addShutdown("iot", func(ctx context.Context) error {
		iotService.Shutdown()
		return nil
	})
	return nil
}

func (a *App) initIoTEvents(st *serviceState) error {
	iotEventsService := svciotevents.NewIoTEventsService(st.accountID, st.region)
	iotEventsService.SetStorageManager(a.server.StorageManager())

	eb := a.server.EventBus()
	if eb != nil {
		iotEventsService.Init(svciotevents.IoTEventsServiceDeps{
			EventBus:       eb,
			StorageManager: a.server.StorageManager(),
		})
	}

	st.iotEventsService = iotEventsService
	st.iotEventsService.RegisterHandlers(a.server.Dispatcher())

	// Wire the rule action dispatcher callbacks now that both IoT and IoT Events
	// services are initialised. The dispatcher was created in initIoT but the
	// callbacks require the broker (Republish) and the IoT Events state machine
	// (BatchPutMessage), which are only available at this point.
	if st.iotService != nil {
		// Wire the rule action dispatcher callbacks for every region so
		// Republish routes to the correct regional broker.
		for _, region := range a.server.StorageManager().ListRegions() {
			var batchPutFn func(context.Context, []map[string]interface{}) error
			if st.iotEventsService != nil {
				batchPutFn = makeIoTEventsBatchPutFn(st, region)
			}
			st.iotService.WireActionCallbacksForRegion(region, batchPutFn)
		}

		// Wire detector model iotTopicPublish action to the MQTT broker.
		if st.iotEventsService != nil {
			if brk := st.iotService.BrokerForRegion(st.region); brk != nil {
				st.iotEventsService.SetRepublishFn(func(topic string, payload []byte) error {
					return brk.Publish(topic, payload)
				})
			}
		}
	}

	return nil
}

// makeIoTEventsBatchPutFn returns a function that bridges the IoT rule
// dispatcher to the IoT Events detector state machine. It resolves the
// IotStore from the iotevents service and calls BatchEvaluate on the
// state machine for each incoming message batch.
func makeIoTEventsBatchPutFn(st *serviceState, region string) func(context.Context, []map[string]interface{}) error {
	return func(ctx context.Context, messages []map[string]interface{}) error {
		store, err := st.iotEventsService.GetStoreForRegion(region)
		if err != nil {
			return fmt.Errorf("iot events batch put: %w", err)
		}
		concrete, ok := store.(*iotstore.IotStore)
		if !ok || concrete.StateMachine() == nil {
			return nil
		}
		inputs := make([]iotstore.InputMessage, 0, len(messages))
		for _, m := range messages {
			name, _ := m["inputName"].(string)
			payload := iotstore.ExtractPayload(m["payload"])
			inputs = append(inputs, iotstore.InputMessage{InputName: name, Payload: payload})
		}
		errs := concrete.StateMachine().BatchEvaluate(ctx, inputs)
		if len(errs) > 0 {
			return fmt.Errorf("iot events batch evaluate: %d message(s) produced errors", len(errs))
		}
		return nil
	}
}

// --- gRPC-Web Admin ---

func (a *App) initGRPCWebAdmin() {
	st := a.state
	grpcWebServer := grpcweb.NewServer(&grpcweb.Config{Port: a.cfg.GRPCWebPort, BindAddr: a.cfg.GRPCWebBindAddr})

	aid := st.accountID
	reg := st.region
	dp := a.cfg.DataPath

	handlers := make([]grpcweb.HandlerRegistration, 0, 32)

	var p string
	var h http.Handler

	p, h = svcacm.NewConnectHandler(st.acmService)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcapigateway.NewConnectHandler(st.apiGatewayService)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svccloudfront.NewConnectHandler(st.cloudFrontService)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcroute53.NewConnectHandler(st.route53Service)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcsecretsmanager.NewConnectHandler(st.secretsManagerService)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcsts.NewConnectHandler(aid)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcwafv2.NewConnectHandler(st.wafv2Service)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcscheduler.NewConnectHandler(st.schedulerService)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcdynamodb.NewConnectHandler(st.dynamoDBService)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svckinesis.NewConnectHandler(st.kinesisService)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svclambda.NewConnectHandler(st.lambdaService)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcs3.NewConnectHandler(a.server.S3Store(), aid, st.s3Service.EncryptionManager())
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcstepfunction.NewConnectHandler(st.stepFunctionService)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcsqs.NewConnectHandler(st.sqsService)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcevents.NewConnectHandler(st.eventBridgeService)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svclogs.NewConnectHandler(st.logsService)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcsns.NewConnectHandler(st.snsService)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svciam.NewConnectHandler(st.iamService)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svckms.NewConnectHandler(st.kmsService)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svccloudwatch.NewConnectHandler(st.cloudWatchService)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svccognito.NewConnectHandler(st.cognitoService)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svccognitoidentity.NewConnectHandler(st.cognitoIdentityService)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcathena.NewConnectHandler(st.athenaService)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svccloudtrail.NewConnectHandler(st.cloudTrailService)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcsesv2.NewConnectHandler(st.sesv2Service)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcssm.NewConnectHandler(st.ssmService)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svctimestreamquery.NewConnectHandler(st.timestreamQueryService)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svctimestreamwrite.NewConnectHandler(st.timestreamWriteService)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcneptuneGraph.NewConnectHandler(st.neptuneGraphService, aid)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	if st.neptuneDataService != nil {
		p, h = svcneptunedata.NewConnectHandler(st.neptuneDataService)
		handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	}
	// RDS admin handler — serves both Neptune and MySQL data via the
	// standalone RDS store lookup (always available, independent of
	// Neptune service initialisation).
	p, h = svcrds.NewConnectHandler(
		svcrds.StoreProvider(func(region string) (storerds.StoreInterface, error) {
			return st.rdsStoreLookup(region)
		}),
		svcrds.EngineProvider(func(engineType string) (svcrds.Engine, error) {
			switch engineType {
			case "mysql":
				if st.vmysqlService != nil {
					return st.vmysqlService, nil
				}
				return nil, fmt.Errorf("mysql engine not available")
			default:
				return nil, fmt.Errorf("unsupported engine: %s", engineType)
			}
		}),
		aid,
		&vmysqlSnapshotOperatorAdapter{svc: st.vmysqlService})
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})
	p, h = svcappsync.NewConnectHandler(st.appSyncService)
	handlers = append(handlers, grpcweb.HandlerRegistration{Path: p, Handler: h})

	grpcweb.RegisterAdminHandlers(grpcWebServer, a.server.Storage(), aid, reg, dp, handlers, a.server.TriggerShutdown)
	grpcWebServer.ServeConsole(a.consoleAssets)
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

	if a.cfg.S3 {
		a.registerListener(listener.ListenerConfig{
			Name:        "s3_website",
			PortKey:     "ports.s3_website",
			DefaultPort: serviceports.S3Website,
			HostSuffix:  fmt.Sprintf(".s3-website.%s.amazonaws.com", st.region),
			Handler: http.HandlerFunc(svcs3.NewWebsiteServer(
				a.server.S3Store(), st.accountID, st.region,
			).HandleRequest),
		})
	}

	if a.cfg.CloudFront && st.cloudFrontService != nil {
		if handler := st.cloudFrontService.DistributionHandler(); handler != nil {
			a.registerListener(listener.ListenerConfig{
				Name:        "cloudfront",
				PortKey:     "ports.cloudfront",
				DefaultPort: serviceports.CloudFront,
				HostSuffix:  ".cloudfront.net",
				Handler:     handler,
			})
		}
	}

	if st.lambdaService != nil {
		lambdaUrlServer := svclambda.NewFunctionURLServer(st.lambdaService, st.accountID, st.region, st.lambdaService)
		a.registerListener(listener.ListenerConfig{
			Name:        "lambda_url",
			PortKey:     "ports.lambda_url",
			DefaultPort: serviceports.LambdaURL,
			HostSuffix:  fmt.Sprintf(".lambda-url.%s.on.aws", st.region),
			Handler:     http.HandlerFunc(lambdaUrlServer.HandleRequest),
		})
	}

	if st.apiGatewayService != nil && st.lambdaService != nil {
		st.apiGatewayService.InitRuntimeServer(a.server.EventBus())
		if handler := st.apiGatewayService.RuntimeHandler(); handler != nil {
			a.server.RegisterAPIGatewayRuntimeHandler(handler)
			a.registerListener(listener.ListenerConfig{
				Name:        "apigateway",
				PortKey:     "ports.apigateway",
				DefaultPort: serviceports.APIGateway,
				HostSuffix:  fmt.Sprintf(".execute-api.%s.amazonaws.com", st.region),
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
			DefaultPort: serviceports.Cognito,
			HostSuffix:  fmt.Sprintf(".auth.%s.amazoncognito.com", st.region),
			Handler:     http.HandlerFunc(st.cognitoService.HostedUIHandler),
		})
	}

	if a.cfg.AppSync && st.appSyncService != nil {
		a.registerListener(listener.ListenerConfig{
			Name:        "appsync_events",
			PortKey:     "ports.appsync_events",
			DefaultPort: serviceports.AppSync,
			HostSuffix:  fmt.Sprintf(".appsync-api.%s.amazonaws.com", st.region),
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
		ctStore := st.cloudTrailService.GetEventStore(regionalStorage, region)
		if ctStore == nil {
			return nil
		}
		return audit.NewCloudTrailRecorder(&cloudTrailStoreAdapter{store: ctStore})
	})
}

// cloudTrailStoreAdapter adapts a CloudTrailStoreInterface to the audit.EventStore interface,
// bridging the audit package's UserIdentity type to the store package's concrete UserIdentity type.
type cloudTrailStoreAdapter struct {
	store cloudtrailstore.CloudTrailStoreInterface
}

// RecordServiceEvent translates an audit UserIdentity to a cloudtrailstore UserIdentity and
// delegates to the underlying CloudTrail store.
func (a *cloudTrailStoreAdapter) RecordServiceEvent(eventName, eventSource string, userIdentity *audit.UserIdentity, sourceIP, accessKeyID string, requestParams, responseElements map[string]interface{}, resources []audit.ResourceEntry) error {
	var storeResources []cloudtrailstore.Resource
	for _, r := range resources {
		storeResources = append(storeResources, cloudtrailstore.Resource{ResourceType: r.ResourceType, ResourceName: r.ResourceName})
	}
	return a.store.RecordServiceEvent(eventName, eventSource, &cloudtrailstore.UserIdentity{
		Type:        userIdentity.Type,
		PrincipalID: userIdentity.PrincipalID,
		ARN:         userIdentity.ARN,
		AccountID:   userIdentity.AccountID,
		UserName:    userIdentity.UserName,
	}, sourceIP, accessKeyID, requestParams, responseElements, storeResources)
}

func (a *App) initPrincipalResolver() {
	iamStore := a.server.IAMStore()
	if iamStore == nil {
		return
	}
	a.server.Dispatcher().SetPrincipalResolver(&iamPrincipalResolverAdapter{store: iamStore})
}

type iamPrincipalResolverAdapter struct {
	store iamstore.IAMStoreInterface
}

// ResolvePrincipal returns the IAM user name associated with the given access key ID.
func (r *iamPrincipalResolverAdapter) ResolvePrincipal(ctx context.Context, accessKeyID string) (string, error) {
	accessKey, err := r.store.AccessKeys().Get(accessKeyID)
	if err != nil || accessKey == nil || accessKey.UserName == "" {
		return "", nil
	}
	return accessKey.UserName, nil
}

func (a *App) registerListener(cfg listener.ListenerConfig) {
	cfg.DefaultPort = a.resolvedPort(cfg.PortKey, cfg.DefaultPort)
	a.lm.Register(cfg)
}

func (a *App) resolvedPort(portKey string, defaultPort int) int {
	if portKey == "" {
		return defaultPort
	}
	cfgStore := appconfig.GetStore()
	if entry, err := cfgStore.Get(portKey); err == nil {
		if p, ok := entry.Value.(int); ok && p > 0 {
			return p
		}
		if p, ok := entry.Value.(float64); ok && p > 0 {
			return int(p)
		}
	}
	return defaultPort
}

// injectS3AuditRecorder creates a CloudTrail audit recorder for the S3 handler.
// This is needed because S3 bypasses the Dispatcher and handles requests via ServeHTTP.
func (a *App) injectS3AuditRecorder(st *serviceState) {
	if a.server.S3Handler() == nil || st.cloudTrailService == nil {
		return
	}
	regionalStorage, err := a.server.StorageManager().GetStorage(st.region)
	if err != nil {
		return
	}
	ctStore := st.cloudTrailService.GetEventStore(regionalStorage, st.region)
	if ctStore == nil {
		return
	}
	recorder := audit.NewCloudTrailRecorder(&cloudTrailStoreAdapter{store: ctStore})
	if handler, ok := a.server.S3Handler().(*svcs3.S3Handler); ok {
		handler.SetAuditRecorder(recorder)
	}
}
