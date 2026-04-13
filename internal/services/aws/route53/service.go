// Package route53 provides AWS Route 53 DNS service operations for vorpalstacks.
package route53

import (
	"fmt"
	"os"
	"sync"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/route53/dnsserver"
	route53store "vorpalstacks/internal/store/aws/route53"
)

var route53FallbackKey struct{}

// Route53Service provides AWS Route 53 DNS operations.
type Route53Service struct {
	accountID string
	dnsStores *route53store.Route53Stores
	dnsServer *dnsserver.DNSServer
	stores    sync.Map // global — single cached instance for store() fallback
}

// NewRoute53Service creates a new Route 53 service instance.
func NewRoute53Service(store storage.BasicStorage, accountID string) *Route53Service {
	return &Route53Service{
		dnsStores: route53store.NewRoute53Stores(
			route53store.NewHostedZoneStore(store, accountID),
			route53store.NewHealthCheckStore(store, accountID),
			route53store.NewRecordSetStore(store),
			route53store.NewChangeStore(store),
			route53store.NewTagStore(store),
			route53store.NewARNBuilder(accountID),
		),
		accountID: accountID,
	}
}

// NewRoute53ServiceWithDNS creates a new Route 53 service with DNS server enabled.
func NewRoute53ServiceWithDNS(store storage.BasicStorage, accountID, bindAddr string, enableDNS bool) (*Route53Service, error) {
	svc := NewRoute53Service(store, accountID)

	if enableDNS {
		dnsServer := dnsserver.NewDNSServer(svc.dnsStores.HostedZones().Raw(), svc.dnsStores.RecordSets().Raw(), bindAddr)
		if err := dnsServer.Start(); err != nil {
			return nil, fmt.Errorf("failed to start DNS server: %w", err)
		}
		svc.dnsServer = dnsServer
		logs.Info("Route53 DNS server enabled", logs.String("address", bindAddr+":53"))
	}

	return svc, nil
}

func (s *Route53Service) store(reqCtx *request.RequestContext) (*route53store.Route53Stores, error) {
	if stores := reqCtx.GetRoute53Stores(); stores != nil {
		return stores.Raw(), nil
	}
	if cached, ok := s.stores.Load(route53FallbackKey); ok {
		if typed, ok := cached.(*route53store.Route53Stores); ok {
			return typed, nil
		}
	}
	storage, err := reqCtx.GetGlobalStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to get global storage: %w", err)
	}
	stores := route53store.NewRoute53Stores(
		route53store.NewHostedZoneStore(storage, s.accountID),
		route53store.NewHealthCheckStore(storage, s.accountID),
		route53store.NewRecordSetStore(storage),
		route53store.NewChangeStore(storage),
		route53store.NewTagStore(storage),
		route53store.NewARNBuilder(s.accountID),
	)
	if actual, loaded := s.stores.LoadOrStore(route53FallbackKey, stores); loaded {
		if typed, ok := actual.(*route53store.Route53Stores); ok {
			return typed, nil
		}
	}
	return stores, nil
}

// HostedZoneStore returns the hosted zone store.
func (s *Route53Service) HostedZoneStore() *route53store.HostedZoneStore {
	return s.dnsStores.HostedZones().Raw()
}

// HealthCheckStore returns the health check store.
func (s *Route53Service) HealthCheckStore() *route53store.HealthCheckStore {
	return s.dnsStores.HealthChecks().Raw()
}

// RecordSetStore returns the record set store.
func (s *Route53Service) RecordSetStore() *route53store.RecordSetStore {
	return s.dnsStores.RecordSets().Raw()
}

// ChangeStore returns the change store.
func (s *Route53Service) ChangeStore() *route53store.ChangeStore {
	return s.dnsStores.Changes().Raw()
}

// ARNBuilder returns the ARN builder.
func (s *Route53Service) ARNBuilder() *route53store.ARNBuilder {
	return s.dnsStores.ARNBuilder()
}

// AccountId returns the account ID.
func (s *Route53Service) AccountId() string {
	return s.accountID
}

// DNSServer returns the DNS server if enabled.
func (s *Route53Service) DNSServer() *dnsserver.DNSServer {
	return s.dnsServer
}

// Shutdown stops the Route 53 DNS server.
func (s *Route53Service) Shutdown() error {
	if s.dnsServer != nil {
		return s.dnsServer.Shutdown()
	}
	return nil
}

// RegisterHandlers registers Route 53 operation handlers with the dispatcher.
func (s *Route53Service) RegisterHandlers(d *dispatcher.Dispatcher) {
	d.RegisterHandlerForService("route53", "CreateHostedZone", s.CreateHostedZone)
	d.RegisterHandlerForService("route53", "GetHostedZone", s.GetHostedZone)
	d.RegisterHandlerForService("route53", "ListHostedZones", s.ListHostedZones)
	d.RegisterHandlerForService("route53", "ListHostedZonesByName", s.ListHostedZonesByName)
	d.RegisterHandlerForService("route53", "DeleteHostedZone", s.DeleteHostedZone)
	d.RegisterHandlerForService("route53", "UpdateHostedZoneComment", s.UpdateHostedZoneComment)
	d.RegisterHandlerForService("route53", "ChangeResourceRecordSets", s.ChangeResourceRecordSets)
	d.RegisterHandlerForService("route53", "ListResourceRecordSets", s.ListResourceRecordSets)
	d.RegisterHandlerForService("route53", "GetChange", s.GetChange)
	d.RegisterHandlerForService("route53", "CreateHealthCheck", s.CreateHealthCheck)
	d.RegisterHandlerForService("route53", "GetHealthCheck", s.GetHealthCheck)
	d.RegisterHandlerForService("route53", "ListHealthChecks", s.ListHealthChecks)
	d.RegisterHandlerForService("route53", "DeleteHealthCheck", s.DeleteHealthCheck)
	d.RegisterHandlerForService("route53", "UpdateHealthCheck", s.UpdateHealthCheck)
	d.RegisterHandlerForService("route53", "AssociateVPCWithHostedZone", s.AssociateVPCWithHostedZone)
	d.RegisterHandlerForService("route53", "DisassociateVPCFromHostedZone", s.DisassociateVPCFromHostedZone)
	d.RegisterHandlerForService("route53", "ChangeTagsForResource", s.ChangeTagsForResource)
	d.RegisterHandlerForService("route53", "ListTagsForResource", s.ListTagsForResource)
	d.RegisterHandlerForService("route53", "ListReusableDelegationSets", s.ListReusableDelegationSets)
	d.RegisterHandlerForService("route53", "GetDNSSEC", s.GetDNSSEC)
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
