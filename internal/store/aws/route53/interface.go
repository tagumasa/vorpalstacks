package route53

import (
	"vorpalstacks/internal/store/aws/common"
)

// HostedZoneStoreInterface defines operations for managing Route 53 hosted zones.
type HostedZoneStoreInterface interface {
	Get(id string) (*HostedZone, error)
	GetByName(name string) (*HostedZone, error)
	List(marker string, maxItems int) (*HostedZoneListResult, error)
	ListByName() ([]*HostedZone, error)
	Create(zone *HostedZone) error
	Update(zone *HostedZone) error
	Delete(id string) error
	Exists(id string) bool
	Raw() *HostedZoneStore
}

// HealthCheckStoreInterface defines operations for managing Route 53 health checks.
type HealthCheckStoreInterface interface {
	Get(id string) (*HealthCheck, error)
	List(marker string, maxItems int) (*HealthCheckListResult, error)
	Create(healthCheck *HealthCheck) error
	Update(healthCheck *HealthCheck) error
	Delete(id string) error
	Exists(id string) bool
	Raw() *HealthCheckStore
}

// RecordSetStoreInterface defines operations for managing Route 53 record sets.
type RecordSetStoreInterface interface {
	Get(hostedZoneId, name, recordType, setIdentifier string) (*ResourceRecordSet, error)
	List(hostedZoneId string) ([]*ResourceRecordSet, error)
	Create(hostedZoneId string, recordSet *ResourceRecordSet) error
	Upsert(hostedZoneId string, recordSet *ResourceRecordSet) error
	Delete(hostedZoneId, name, recordType, setIdentifier string) error
	Exists(key string) bool
	Raw() *RecordSetStore
}

// ChangeStoreInterface defines operations for managing Route 53 change operations.
type ChangeStoreInterface interface {
	Get(id string) (*ChangeInfo, error)
	Create(change *ChangeInfo) error
	UpdateStatus(id, status string) error
	UpdateStatusAndGet(id, status string) (*ChangeInfo, error)
	Exists(id string) bool
	Raw() *ChangeStore
}

// TagStoreInterface defines operations for managing Route 53 tags.
type TagStoreInterface interface {
	ListTagsForResource(resourceKey string) ([]common.Tag, error)
	TagResource(resourceKey string, tags []common.Tag) error
	// Raw returns the underlying tag store.
	Raw() *TagStore
}

// Route53StoresInterface defines access to all Route 53 stores.
type Route53StoresInterface interface {
	HostedZones() HostedZoneStoreInterface
	HealthChecks() HealthCheckStoreInterface
	RecordSets() RecordSetStoreInterface
	Changes() ChangeStoreInterface
	Tags() TagStoreInterface
	ARNBuilder() *ARNBuilder
	Raw() *Route53Stores
}

// Route53Stores provides access to all Route 53 stores.
type Route53Stores struct {
	hostedZones  *HostedZoneStore
	healthChecks *HealthCheckStore
	recordSets   *RecordSetStore
	changes      *ChangeStore
	tags         *TagStore
	arnBuilder   *ARNBuilder
}

// NewRoute53Stores creates a new Route53Stores with the given stores.
func NewRoute53Stores(
	hostedZones *HostedZoneStore,
	healthChecks *HealthCheckStore,
	recordSets *RecordSetStore,
	changes *ChangeStore,
	tags *TagStore,
	arnBuilder *ARNBuilder,
) *Route53Stores {
	return &Route53Stores{
		hostedZones:  hostedZones,
		healthChecks: healthChecks,
		recordSets:   recordSets,
		changes:      changes,
		tags:         tags,
		arnBuilder:   arnBuilder,
	}
}

// HostedZones returns the hosted zone store.
func (s *Route53Stores) HostedZones() HostedZoneStoreInterface {
	return s.hostedZones
}

// HealthChecks returns the health check store.
func (s *Route53Stores) HealthChecks() HealthCheckStoreInterface {
	return s.healthChecks
}

// RecordSets returns the record set store.
func (s *Route53Stores) RecordSets() RecordSetStoreInterface {
	return s.recordSets
}

// Changes returns the change store.
func (s *Route53Stores) Changes() ChangeStoreInterface {
	return s.changes
}

// Tags returns the tag store.
func (s *Route53Stores) Tags() TagStoreInterface {
	return s.tags
}

// ARNBuilder returns the ARN builder.
func (s *Route53Stores) ARNBuilder() *ARNBuilder {
	return s.arnBuilder
}

// Raw returns the Route 53 stores.
func (s *Route53Stores) Raw() *Route53Stores {
	return s
}

// Raw returns the hosted zone store.
func (s *HostedZoneStore) Raw() *HostedZoneStore {
	return s
}

// Raw returns the health check store.
func (s *HealthCheckStore) Raw() *HealthCheckStore {
	return s
}

// Raw returns the record set store.
func (s *RecordSetStore) Raw() *RecordSetStore {
	return s
}

// Raw returns the change store.
func (s *ChangeStore) Raw() *ChangeStore {
	return s
}

// Raw returns the tag store.
func (s *TagStore) Raw() *TagStore {
	return s
}
