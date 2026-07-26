// Package neptune implements the AWS Neptune Management API service handlers.
package neptune

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	rdssvc "vorpalstacks/internal/services/aws/rds"
	storecommon "vorpalstacks/internal/store/aws/common"
	neptunestore "vorpalstacks/internal/store/aws/rds/neptune"
)

type NeptuneService struct {
	accountID      string
	region         string
	serverHost     string
	storageManager *storage.RegionStorageManager
	stores         sync.Map
	eventBus       *eventbus.EventBus
	cancelCleanup  context.CancelFunc
	engine         rdssvc.Engine
	mysqlEngine    rdssvc.Engine
	porter         rdssvc.GetPorter
	snapOp         rdssvc.SnapshotOperator
}

// NewNeptuneService creates a new NeptuneService for the specified account and
// region. A background goroutine is started to periodically purge old events.
func NewNeptuneService(accountID, region string) *NeptuneService {
	ctx, cancel := context.WithCancel(context.Background())
	s := &NeptuneService{accountID: accountID, region: region, cancelCleanup: cancel}
	go s.cleanupOldEvents(ctx)
	return s
}

// Close stops the background event cleanup goroutine.
func (s *NeptuneService) Close() {
	if s.cancelCleanup != nil {
		s.cancelCleanup()
	}
}

// SetServerHost sets the server hostname used to construct cluster Endpoint addresses.
func (s *NeptuneService) SetServerHost(host string) {
	s.serverHost = host
}

// endpointAddress returns the hostname for a cluster's Endpoint field.
func (s *NeptuneService) endpointAddress(clusterID string) string {
	return s.endpointAddressFor(clusterID, "neptune")
}

// endpointAddressFor returns the hostname for a cluster/instance Endpoint
// field, choosing the appropriate AWS-style DNS suffix based on engine type.
func (s *NeptuneService) endpointAddressFor(resourceID, engineType string) string {
	if s.serverHost != "" {
		if host, _, err := net.SplitHostPort(s.serverHost); err == nil {
			return host
		}
		return s.serverHost
	}
	if engineType == "mysql" {
		return fmt.Sprintf("%s.%s.%s.rds.amazonaws.com", resourceID, s.accountID, s.region)
	}
	return fmt.Sprintf("%s.cluster-%s.%s.neptune.amazonaws.com", resourceID, s.accountID, s.region)
}

// cleanupOldEvents periodically purges events older than the retention period.
func (s *NeptuneService) cleanupOldEvents(ctx context.Context) {
	defer func() { resilience.RecoverPanic("Neptune cleanupOldEvents") }()
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.stores.Range(func(_, value any) bool {
				store, ok := value.(neptunestore.NeptuneStoreInterface)
				if !ok {
					return true
				}
				if purgeable, ok := store.(interface{ PurgeOldEvents() error }); ok {
					if err := purgeable.PurgeOldEvents(); err != nil {
						logs.Warn("failed to purge old events", logs.Err(err))
					}
				}
				return true
			})
		}
	}
}

// SetStorageManager injects the region storage manager for per-region store
// caching and admin console access.
func (s *NeptuneService) SetStorageManager(sm *storage.RegionStorageManager) {
	s.storageManager = sm
}

// SetEventBus injects the event bus for cross-service invocations such as
// EC2 subnet lookups.
func (s *NeptuneService) SetEventBus(bus *eventbus.EventBus) {
	s.eventBus = bus
}

func (s *NeptuneService) SetEngine(e rdssvc.Engine) {
	s.engine = e
}

// SetMysqlEngine wires the MySQL (vmysql) engine so the handler can
// serve RDS SDK clients that create MySQL instances/clusters. When the
// engine type of a resource is "mysql", operations dispatch to this
// engine instead of the Neptune engine.
func (s *NeptuneService) SetMysqlEngine(e rdssvc.Engine) {
	s.mysqlEngine = e
}

// SetSnapshotOperator wires the row-level snapshot operator (vmysql).
// When set, CreateDBSnapshot captures actual row data, and
// RestoreDBInstanceFromDBSnapshot / DeleteDBSnapshot recover / clean up
// that data. When nil, snapshots are metadata-only.
func (s *NeptuneService) SetSnapshotOperator(op rdssvc.SnapshotOperator) {
	s.snapOp = op
}

// engineFor returns the appropriate engine for the given engine type.
// Recognised values: "neptune" and "" (the default) map to the Neptune
// engine; "mysql" maps to the MySQL engine. All other values return nil
// so that engine open is a no-op rather than silently opening the wrong
// engine.
func (s *NeptuneService) engineFor(engineType string) rdssvc.Engine {
	switch engineType {
	case "", "neptune":
		return s.engine
	case "mysql":
		return s.mysqlEngine
	default:
		return nil
	}
}

func (s *NeptuneService) SetClusterPorter(p rdssvc.GetPorter) {
	s.porter = p
}

// GetStoreForRegion returns the cached Neptune store for the given region,
// creating one if not already cached. Used by both HTTP handlers and the
// admin console to ensure a single store instance per region.
func (s *NeptuneService) GetStoreForRegion(region string) (neptunestore.NeptuneStoreInterface, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, region, func() (neptunestore.NeptuneStoreInterface, error) {
		if s.storageManager == nil {
			return nil, fmt.Errorf("storage manager not set")
		}
		rs, err := s.storageManager.GetStorage(region)
		if err != nil {
			return nil, err
		}
		return neptunestore.NewNeptuneStore(rs), nil
	})
}

// store resolves the NeptuneStore for the current request context, using the
// per-region cache.
func (s *NeptuneService) store(reqCtx *request.RequestContext) (neptunestore.NeptuneStoreInterface, error) {
	region := reqCtx.GetRegion()
	return s.GetStoreForRegion(region)
}

// IsSubnetInUse implements eventbus.SubnetUsageChecker. It returns true if
// any DB subnet group in the given region references the specified subnet.
// EC2 calls this before deleting a subnet to prevent orphaned references.
func (s *NeptuneService) IsSubnetInUse(ctx context.Context, region, subnetId string) bool {
	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return false
	}
	groups, err := store.ListSubnetGroups()
	if err != nil {
		return false
	}
	for _, g := range groups {
		for _, sub := range g.Subnets {
			if sub.SubnetIdentifier == subnetId {
				return true
			}
		}
	}
	return false
}

func recordEvent(store neptunestore.NeptuneStoreInterface, sourceType, sourceID, sourceArn, message string, categories []string) {
	evt := &neptunestore.Event{
		Date:             time.Now().UTC(),
		EventCategories:  categories,
		Message:          message,
		SourceArn:        sourceArn,
		SourceIdentifier: sourceID,
		SourceType:       sourceType,
	}
	if err := store.RecordEvent(evt); err != nil {
		logs.Warn("failed to record event", logs.Err(err))
	}
}

// RegisterHandlers registers all Neptune Management API action handlers with
// the given dispatcher.
func (s *NeptuneService) RegisterHandlers(d handler.Registrar) {
	d.RegisterHandlerForService("neptune", "CreateDBCluster", s.CreateDBCluster)
	d.RegisterHandlerForService("neptune", "DeleteDBCluster", s.DeleteDBCluster)
	d.RegisterHandlerForService("neptune", "ModifyDBCluster", s.ModifyDBCluster)
	d.RegisterHandlerForService("neptune", "DescribeDBClusters", s.DescribeDBClusters)
	d.RegisterHandlerForService("neptune", "StartDBCluster", s.StartDBCluster)
	d.RegisterHandlerForService("neptune", "StopDBCluster", s.StopDBCluster)
	d.RegisterHandlerForService("neptune", "FailoverDBCluster", s.FailoverDBCluster)
	d.RegisterHandlerForService("neptune", "CreateDBClusterEndpoint", s.CreateDBClusterEndpoint)
	d.RegisterHandlerForService("neptune", "DescribeDBClusterEndpoints", s.DescribeDBClusterEndpoints)
	d.RegisterHandlerForService("neptune", "ModifyDBClusterEndpoint", s.ModifyDBClusterEndpoint)
	d.RegisterHandlerForService("neptune", "DeleteDBClusterEndpoint", s.DeleteDBClusterEndpoint)
	d.RegisterHandlerForService("neptune", "CreateDBInstance", s.CreateDBInstance)
	d.RegisterHandlerForService("neptune", "DeleteDBInstance", s.DeleteDBInstance)
	d.RegisterHandlerForService("neptune", "ModifyDBInstance", s.ModifyDBInstance)
	d.RegisterHandlerForService("neptune", "DescribeDBInstances", s.DescribeDBInstances)
	d.RegisterHandlerForService("neptune", "RebootDBInstance", s.RebootDBInstance)
	d.RegisterHandlerForService("neptune", "RestoreDBInstanceFromDBSnapshot", s.RestoreDBInstanceFromDBSnapshot)
	d.RegisterHandlerForService("neptune", "RestoreDBInstanceToPointInTime", s.RestoreDBInstanceToPointInTime)
	d.RegisterHandlerForService("neptune", "CreateDBInstanceReadReplica", s.CreateDBInstanceReadReplica)
	d.RegisterHandlerForService("neptune", "PromoteReadReplica", s.PromoteReadReplica)
	d.RegisterHandlerForService("neptune", "ModifyDBSnapshot", s.ModifyDBSnapshot)
	d.RegisterHandlerForService("neptune", "DownloadDBLogFilePortion", s.DownloadDBLogFilePortion)
	d.RegisterHandlerForService("neptune", "CreateDBSnapshot", s.CreateDBSnapshot)
	d.RegisterHandlerForService("neptune", "DescribeDBSnapshots", s.DescribeDBSnapshots)
	d.RegisterHandlerForService("neptune", "DeleteDBSnapshot", s.DeleteDBSnapshot)
	d.RegisterHandlerForService("neptune", "CopyDBSnapshot", s.CopyDBSnapshot)
	d.RegisterHandlerForService("neptune", "DescribeDBSnapshotAttributes", s.DescribeDBSnapshotAttributes)
	d.RegisterHandlerForService("neptune", "ModifyDBSnapshotAttribute", s.ModifyDBSnapshotAttribute)
	d.RegisterHandlerForService("neptune", "CreateDBClusterSnapshot", s.CreateDBClusterSnapshot)
	d.RegisterHandlerForService("neptune", "DeleteDBClusterSnapshot", s.DeleteDBClusterSnapshot)
	d.RegisterHandlerForService("neptune", "DescribeDBClusterSnapshots", s.DescribeDBClusterSnapshots)
	d.RegisterHandlerForService("neptune", "CopyDBClusterSnapshot", s.CopyDBClusterSnapshot)
	d.RegisterHandlerForService("neptune", "DescribeDBClusterSnapshotAttributes", s.DescribeDBClusterSnapshotAttributes)
	d.RegisterHandlerForService("neptune", "ModifyDBClusterSnapshotAttribute", s.ModifyDBClusterSnapshotAttribute)
	d.RegisterHandlerForService("neptune", "CreateDBClusterParameterGroup", s.CreateDBClusterParameterGroup)
	d.RegisterHandlerForService("neptune", "DeleteDBClusterParameterGroup", s.DeleteDBClusterParameterGroup)
	d.RegisterHandlerForService("neptune", "DescribeDBClusterParameterGroups", s.DescribeDBClusterParameterGroups)
	d.RegisterHandlerForService("neptune", "DescribeDBClusterParameters", s.DescribeDBClusterParameters)
	d.RegisterHandlerForService("neptune", "ModifyDBClusterParameterGroup", s.ModifyDBClusterParameterGroup)
	d.RegisterHandlerForService("neptune", "CreateDBParameterGroup", s.CreateDBParameterGroup)
	d.RegisterHandlerForService("neptune", "DeleteDBParameterGroup", s.DeleteDBParameterGroup)
	d.RegisterHandlerForService("neptune", "DescribeDBParameterGroups", s.DescribeDBParameterGroups)
	d.RegisterHandlerForService("neptune", "DescribeDBParameters", s.DescribeDBParameters)
	d.RegisterHandlerForService("neptune", "ModifyDBParameterGroup", s.ModifyDBParameterGroup)
	d.RegisterHandlerForService("neptune", "ResetDBClusterParameterGroup", s.ResetDBClusterParameterGroup)
	d.RegisterHandlerForService("neptune", "ResetDBParameterGroup", s.ResetDBParameterGroup)
	d.RegisterHandlerForService("neptune", "CopyDBClusterParameterGroup", s.CopyDBClusterParameterGroup)
	d.RegisterHandlerForService("neptune", "CopyDBParameterGroup", s.CopyDBParameterGroup)
	d.RegisterHandlerForService("neptune", "DescribeEngineDefaultClusterParameters", s.DescribeEngineDefaultClusterParameters)
	d.RegisterHandlerForService("neptune", "DescribeEngineDefaultParameters", s.DescribeEngineDefaultParameters)
	d.RegisterHandlerForService("neptune", "CreateGlobalCluster", s.CreateGlobalCluster)
	d.RegisterHandlerForService("neptune", "DeleteGlobalCluster", s.DeleteGlobalCluster)
	d.RegisterHandlerForService("neptune", "DescribeGlobalClusters", s.DescribeGlobalClusters)
	d.RegisterHandlerForService("neptune", "ModifyGlobalCluster", s.ModifyGlobalCluster)
	d.RegisterHandlerForService("neptune", "FailoverGlobalCluster", s.FailoverGlobalCluster)
	d.RegisterHandlerForService("neptune", "SwitchoverGlobalCluster", s.SwitchoverGlobalCluster)
	d.RegisterHandlerForService("neptune", "RemoveFromGlobalCluster", s.RemoveFromGlobalCluster)
	d.RegisterHandlerForService("neptune", "CreateDBSubnetGroup", s.CreateDBSubnetGroup)
	d.RegisterHandlerForService("neptune", "DeleteDBSubnetGroup", s.DeleteDBSubnetGroup)
	d.RegisterHandlerForService("neptune", "DescribeDBSubnetGroups", s.DescribeDBSubnetGroups)
	d.RegisterHandlerForService("neptune", "ModifyDBSubnetGroup", s.ModifyDBSubnetGroup)
	d.RegisterHandlerForService("neptune", "CreateEventSubscription", s.CreateEventSubscription)
	d.RegisterHandlerForService("neptune", "DeleteEventSubscription", s.DeleteEventSubscription)
	d.RegisterHandlerForService("neptune", "DescribeEventSubscriptions", s.DescribeEventSubscriptions)
	d.RegisterHandlerForService("neptune", "ModifyEventSubscription", s.ModifyEventSubscription)
	d.RegisterHandlerForService("neptune", "AddSourceIdentifierToSubscription", s.AddSourceIdentifierToSubscription)
	d.RegisterHandlerForService("neptune", "RemoveSourceIdentifierFromSubscription", s.RemoveSourceIdentifierFromSubscription)
	d.RegisterHandlerForService("neptune", "AddTagsToResource", s.AddTagsToResource)
	d.RegisterHandlerForService("neptune", "ListTagsForResource", s.ListTagsForResource)
	d.RegisterHandlerForService("neptune", "RemoveTagsFromResource", s.RemoveTagsFromResource)
	d.RegisterHandlerForService("neptune", "DescribeEvents", s.DescribeEvents)
	d.RegisterHandlerForService("neptune", "DescribeEventCategories", s.DescribeEventCategories)
	d.RegisterHandlerForService("neptune", "DescribePendingMaintenanceActions", s.DescribePendingMaintenanceActions)
	d.RegisterHandlerForService("neptune", "ApplyPendingMaintenanceAction", s.ApplyPendingMaintenanceAction)
	d.RegisterHandlerForService("neptune", "DescribeDBEngineVersions", s.DescribeDBEngineVersions)
	d.RegisterHandlerForService("neptune", "DescribeOrderableDBInstanceOptions", s.DescribeOrderableDBInstanceOptions)
	d.RegisterHandlerForService("neptune", "DescribeValidDBInstanceModifications", s.DescribeValidDBInstanceModifications)
	d.RegisterHandlerForService("neptune", "DescribeDBLogFiles", s.DescribeDBLogFiles)
	d.RegisterHandlerForService("neptune", "DescribeOptionGroups", s.DescribeOptionGroups)
	d.RegisterHandlerForService("neptune", "CreateOptionGroup", s.CreateOptionGroup)
	d.RegisterHandlerForService("neptune", "DeleteOptionGroup", s.DeleteOptionGroup)
	d.RegisterHandlerForService("neptune", "ModifyOptionGroup", s.ModifyOptionGroup)
	d.RegisterHandlerForService("neptune", "RestoreDBClusterFromSnapshot", s.RestoreDBClusterFromSnapshot)
	d.RegisterHandlerForService("neptune", "RestoreDBClusterToPointInTime", s.RestoreDBClusterToPointInTime)
	d.RegisterHandlerForService("neptune", "PromoteReadReplicaDBCluster", s.PromoteReadReplicaDBCluster)
	d.RegisterHandlerForService("neptune", "AddRoleToDBCluster", s.AddRoleToDBCluster)
	d.RegisterHandlerForService("neptune", "RemoveRoleFromDBCluster", s.RemoveRoleFromDBCluster)
}
