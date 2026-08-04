// Package rds provides the shared RDS service layer and admin console
// handler that serves management RPCs for both Neptune and MySQL database
// engines via the gRPC-Web admin interface.
package rds

import (
	"fmt"

	storerds "vorpalstacks/internal/store/aws/rds"
)

// rdsStores holds the common RDS store and the resolved region for a
// single admin-handler request. It is the analogue of ACM's acmStores:
// a local wrapper type that lets admin_handler.go avoid importing the
// store package directly (AGENTS.md #29).
type rdsStores struct {
	store  storerds.StoreInterface
	region string
}

// StoreProvider returns the common RDS store for a given region. This
// decouples the service layer from any specific engine sub-service
// (NeptuneService, vmysqlService, etc.), allowing the admin handler to
// serve data from all RDS engines through the common store interface.
type StoreProvider func(region string) (storerds.StoreInterface, error)

// EngineProvider supplies per-engine lifecycle managers keyed by engine
// type (e.g. "mysql", "neptune").
type EngineProvider func(engineType string) (Engine, error)

// RDSService is the shared service layer for the RDS admin console. It
// owns the store provider, engine provider, account ID, and optional
// snapshot operator. Admin handler RPC methods delegate to *Core methods
// on this type, keeping admin_handler.go free of store-package imports
// per AGENTS.md #29.
type RDSService struct {
	stores    StoreProvider
	engines   EngineProvider
	accountId string
	snapOp    SnapshotOperator
}

// NewRDSService creates a new RDS service backed by the given store and
// engine providers.
func NewRDSService(stores StoreProvider, engines EngineProvider, accountId string, snapOp SnapshotOperator) *RDSService {
	return &RDSService{
		stores:    stores,
		engines:   engines,
		accountId: accountId,
		snapOp:    snapOp,
	}
}

// GetStoreForRegion returns the common RDS store wrapper for the given
// region, creating the underlying store via the StoreProvider.
func (s *RDSService) GetStoreForRegion(region string) (*rdsStores, error) {
	store, err := s.stores(region)
	if err != nil {
		return nil, fmt.Errorf("rds store for region %q: %w", region, err)
	}
	return &rdsStores{store: store, region: region}, nil
}

// AccountID returns the AWS account ID associated with this service.
func (s *RDSService) AccountID() string { return s.accountId }

// Engines returns the engine provider function used by Create/Delete Core
// methods to start or stop database engines.
func (s *RDSService) Engines() EngineProvider { return s.engines }

// SnapOp returns the snapshot operator (may be nil).
func (s *RDSService) SnapOp() SnapshotOperator { return s.snapOp }
