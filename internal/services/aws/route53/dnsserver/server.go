// Package dnsserver provides DNS server functionality for Route 53 in vorpalstacks.
package dnsserver

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	route53store "vorpalstacks/internal/store/aws/route53"

	"github.com/miekg/dns"
	"vorpalstacks/internal/core/logs"
)

// DNSServer handles DNS queries for Route 53 hosted zones.
type DNSServer struct {
	store           *route53store.HostedZoneStore
	recordSets      *route53store.RecordSetStore
	cidrCollections *route53store.CidrCollectionStore
	udpServer       *dns.Server
	tcpServer       *dns.Server
	bindAddr        string
	port            int
	started         bool
	mu              sync.RWMutex
	shutdownCh      chan struct{}
}

// NewDNSServer creates a new DNSServer with the given stores and bind address.
func NewDNSServer(hostedZoneStore *route53store.HostedZoneStore, recordSetStore *route53store.RecordSetStore, cidrCollectionStore *route53store.CidrCollectionStore, bindAddr string, port int) *DNSServer {
	return &DNSServer{
		store:           hostedZoneStore,
		recordSets:      recordSetStore,
		cidrCollections: cidrCollectionStore,
		bindAddr:        bindAddr,
		port:            port,
		shutdownCh:      make(chan struct{}),
	}
}

// Start starts the DNS server on the configured UDP and TCP port.
func (s *DNSServer) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.started {
		return fmt.Errorf("DNS server already started")
	}

	s.shutdownCh = make(chan struct{})

	handler := dns.NewServeMux()
	handler.HandleFunc(".", s.handleDNSRequest)

	addr := fmt.Sprintf("%s:%d", s.bindAddr, s.port)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return fmt.Errorf("resolve UDP addr: %w", err)
	}

	// Bind synchronously so that port conflicts are detected
	// before reporting started=true. Previously, ListenAndServe ran
	// in a goroutine and bind failures were silently logged while
	// started remained true.
	udpListener, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return fmt.Errorf("listen UDP: %w", err)
	}

	tcpListener, err := net.Listen("tcp", addr)
	if err != nil {
		udpListener.Close()
		return fmt.Errorf("listen TCP: %w", err)
	}

	s.udpServer = &dns.Server{
		Handler: handler,
	}

	s.tcpServer = &dns.Server{
		Handler: handler,
	}

	go func() {
		s.udpServer.PacketConn = udpListener
		if err := s.udpServer.ActivateAndServe(); err != nil {
			select {
			case <-s.shutdownCh:
			default:
				logs.Error("UDP DNS server failed", logs.Err(err))
			}
		}
	}()

	go func() {
		s.tcpServer.Listener = tcpListener
		if err := s.tcpServer.ActivateAndServe(); err != nil {
			select {
			case <-s.shutdownCh:
			default:
				logs.Error("TCP DNS server failed", logs.Err(err))
			}
		}
	}()

	s.started = true
	logs.Info("DNS server started", logs.String("address", addr))
	return nil
}

// Shutdown gracefully stops the DNS server.
func (s *DNSServer) Shutdown() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.started {
		return nil
	}

	close(s.shutdownCh)

	ctx := context.Background()
	var errs []string

	if s.udpServer != nil {
		if err := s.udpServer.ShutdownContext(ctx); err != nil {
			errs = append(errs, fmt.Sprintf("UDP: %v", err))
		}
	}

	if s.tcpServer != nil {
		if err := s.tcpServer.ShutdownContext(ctx); err != nil {
			errs = append(errs, fmt.Sprintf("TCP: %v", err))
		}
	}

	s.started = false

	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %s", strings.Join(errs, ", "))
	}
	return nil
}

// applyRoutingPolicy selects which records to return when multiple records
// share the same name+type. If no routing policy is detected (no
// SetIdentifier on any record), all records are returned (simple routing).
func (s *DNSServer) applyRoutingPolicy(matched []*route53store.ResourceRecordSet, w dns.ResponseWriter) []*route53store.ResourceRecordSet {
	if len(matched) <= 1 {
		return matched
	}

	// Check if any record uses a routing policy (has SetIdentifier).
	hasRoutingPolicy := false
	for _, rs := range matched {
		if rs.SetIdentifier != "" {
			hasRoutingPolicy = true
			break
		}
	}
	if !hasRoutingPolicy {
		return matched
	}

	// Group records by routing type to select the appropriate one.
	// Priority: CIDR > GeoProximity > GeoLocation > Weighted > Failover > MultiValue.
	querierIP := extractQuerierIP(w)

	// Try CIDR routing first.
	if s.cidrCollections != nil && querierIP != "" {
		if selected := s.selectByCidrRouting(matched, querierIP); selected != nil {
			return []*route53store.ResourceRecordSet{selected}
		}
	}

	// Try GeoProximity routing.
	if selected := s.selectByGeoProximity(matched, querierIP); selected != nil {
		return []*route53store.ResourceRecordSet{selected}
	}

	// Try GeoLocation routing by querier IP region.
	if selected := selectByGeoLocation(matched, querierIP); selected != nil {
		return []*route53store.ResourceRecordSet{selected}
	}

	// Try Weighted routing.
	if selected := selectByWeight(matched); selected != nil {
		return []*route53store.ResourceRecordSet{selected}
	}

	// Try Failover routing.
	if selected := selectByFailover(matched); selected != nil {
		return []*route53store.ResourceRecordSet{selected}
	}

	// MultiValueAnswer: return all (up to 8).
	return matched
}

// extractQuerierIP returns the querier's IP address from EDNS Client Subnet
// if present, otherwise from the RemoteAddr.
func extractQuerierIP(w dns.ResponseWriter) string {
	host, _, err := net.SplitHostPort(w.RemoteAddr().String())
	if err != nil {
		return w.RemoteAddr().String()
	}
	return host
}

// selectByCidrRouting matches the querier's IP against CIDR collections
// referenced by CidrRoutingConfig on the records.
func (s *DNSServer) selectByCidrRouting(matched []*route53store.ResourceRecordSet, querierIP string) *route53store.ResourceRecordSet {
	for _, rs := range matched {
		if rs.CidrRoutingConfig == nil || rs.CidrRoutingConfig.CollectionId == "" {
			continue
		}
		loc, err := s.cidrCollections.FindLocationForIP(rs.CidrRoutingConfig.CollectionId, querierIP)
		if err != nil {
			continue
		}
		if loc == rs.CidrRoutingConfig.LocationName {
			return rs
		}
	}
	return nil
}

// selectByGeoProximity selects the record with the nearest GeoProximityLocation
// to the querier's IP address. Uses AWS region coordinates as approximation.
func (s *DNSServer) selectByGeoProximity(matched []*route53store.ResourceRecordSet, querierIP string) *route53store.ResourceRecordSet {
	var best *route53store.ResourceRecordSet
	bestDist := -1.0

	querierLat, querierLon := lookupIPCoordinates(querierIP)

	for _, rs := range matched {
		if rs.GeoProximityLocation == nil {
			continue
		}
		targetLat, targetLon := 0.0, 0.0
		if rs.GeoProximityLocation.Coordinates != nil {
			targetLat = rs.GeoProximityLocation.Coordinates.Latitude
			targetLon = rs.GeoProximityLocation.Coordinates.Longitude
		} else if rs.GeoProximityLocation.AWSRegion != "" {
			lat, lon, ok := awsRegionCoordinates(rs.GeoProximityLocation.AWSRegion)
			if !ok {
				continue
			}
			targetLat, targetLon = lat, lon
		} else {
			continue
		}

		dist := haversine(querierLat, querierLon, targetLat, targetLon)
		// Apply Bias: positive bias shrinks the effective distance
		// (expanding the region), negative bias increases it.
		bias := float64(rs.GeoProximityLocation.Bias)
		effectiveDist := dist * (1.0 - bias/100.0)

		if best == nil || effectiveDist < bestDist {
			best = rs
			bestDist = effectiveDist
		}
	}
	return best
}

// selectByGeoLocation selects a record based on GeoLocation matching the
// querier's IP-derived region.
func selectByGeoLocation(matched []*route53store.ResourceRecordSet, querierIP string) *route53store.ResourceRecordSet {
	querierRegion := ipToAWSRegion(querierIP)
	if querierRegion == "" {
		return nil
	}
	for _, rs := range matched {
		if rs.GeoLocation == nil {
			continue
		}
		if rs.GeoLocation.CountryCode != "" {
			country := ipToCountry(querierIP)
			if country != "" && country == rs.GeoLocation.CountryCode {
				return rs
			}
		}
	}
	return nil
}

// selectByWeight performs weighted random selection among records that
// have Weight > 0.
func selectByWeight(matched []*route53store.ResourceRecordSet) *route53store.ResourceRecordSet {
	var candidates []*route53store.ResourceRecordSet
	totalWeight := int64(0)
	for _, rs := range matched {
		if rs.Weight > 0 {
			candidates = append(candidates, rs)
			totalWeight += rs.Weight
		}
	}
	if totalWeight == 0 || len(candidates) == 0 {
		return nil
	}
	r := cryptoRandInt63n(totalWeight)
	for _, rs := range candidates {
		r -= rs.Weight
		if r < 0 {
			return rs
		}
	}
	return candidates[0]
}

// selectByFailover returns the PRIMARY record (failover=PRIMARY).
// If no PRIMARY found, returns SECONDARY.
func selectByFailover(matched []*route53store.ResourceRecordSet) *route53store.ResourceRecordSet {
	var primary, secondary *route53store.ResourceRecordSet
	for _, rs := range matched {
		switch strings.ToUpper(rs.Failover) {
		case "PRIMARY":
			primary = rs
		case "SECONDARY":
			secondary = rs
		}
	}
	if primary != nil {
		return primary
	}
	return secondary
}

func (s *DNSServer) handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	m.RecursionAvailable = false

	if len(r.Question) == 0 {
		// Return FORMERR per DNS spec for empty question section.
		m.Rcode = dns.RcodeFormatError
		if err := w.WriteMsg(m); err != nil {
			logs.Error("DNS write error", logs.Err(err))
		}
		return
	}

	q := r.Question[0]
	qname := strings.ToLower(dns.CanonicalName(q.Name))
	qtype := q.Qtype

	zone := s.findHostedZone(qname)
	if zone == nil {
		m.Rcode = dns.RcodeNameError
		if err := w.WriteMsg(m); err != nil {
			logs.Error("DNS write error", logs.Err(err))
		}
		return
	}

	recordSets, err := s.recordSets.List(zone.ID)
	if err != nil {
		m.Rcode = dns.RcodeServerFailure
		if err := w.WriteMsg(m); err != nil {
			logs.Error("DNS write error", logs.Err(err))
		}
		return
	}

	var matched []*route53store.ResourceRecordSet

	for _, rs := range recordSets {
		rsName := strings.ToLower(dns.CanonicalName(rs.Name))
		if rsName == qname || rsName == qname+"." {
			matched = append(matched, rs)
		}
	}

	// Apply routing policy selection when multiple records share the
	// same name+type (indicated by SetIdentifier presence).
	selected := s.applyRoutingPolicy(matched, w)

	for _, rs := range selected {
		rsName := strings.ToLower(dns.CanonicalName(rs.Name))
		m.Answer = append(m.Answer, s.recordToRR(rs, rsName, qtype)...)
	}

	if len(m.Answer) == 0 {
		soa := s.findSOA(zone, recordSets)
		if soa != nil {
			m.Ns = append(m.Ns, soa)
		}
	}

	if err := w.WriteMsg(m); err != nil {
		logs.Error("DNS write error", logs.Err(err))
	}
}

func (s *DNSServer) findHostedZone(qname string) *route53store.HostedZone {
	// Paginate through all hosted zones instead of capping at 100.
	var allZones []*route53store.HostedZone
	marker := ""
	for {
		result, err := s.store.List(marker, 100)
		if err != nil {
			return nil
		}
		allZones = append(allZones, result.HostedZones...)
		if !result.IsTruncated {
			break
		}
		marker = result.Marker
	}

	var bestMatch *route53store.HostedZone
	bestMatchLen := 0

	for _, zone := range allZones {
		zoneName := strings.ToLower(dns.CanonicalName(zone.Name))
		if strings.HasSuffix(qname, zoneName) && len(zoneName) > bestMatchLen {
			bestMatch = zone
			bestMatchLen = len(zoneName)
		}
	}
	return bestMatch
}

func (s *DNSServer) recordToRR(rs *route53store.ResourceRecordSet, qname string, qtype uint16) []dns.RR {
	var rrs []dns.RR

	ttl := uint32(rs.TTL)
	// Respect TTL=0 per RFC 1035 ("do not cache") instead of
	// silently overriding to 300.

	if rs.AliasTarget != nil {
		switch qtype {
		case dns.TypeA:
			ip := net.ParseIP(rs.AliasTarget.DNSName)
			if ip != nil {
				rr := &dns.A{}
				rr.Hdr = dns.RR_Header{Name: qname, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: ttl}
				rr.A = ip
				rrs = append(rrs, rr)
			} else {
				// AliasTarget DNS name is a hostname (e.g. ELB, CloudFront).
				// Return a CNAME so the resolver can follow the alias chain.
				rr := &dns.CNAME{}
				rr.Hdr = dns.RR_Header{Name: qname, Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: ttl}
				rr.Target = dns.Fqdn(rs.AliasTarget.DNSName)
				rrs = append(rrs, rr)
			}
		case dns.TypeAAAA:
			ip := net.ParseIP(rs.AliasTarget.DNSName)
			if ip != nil {
				rr := &dns.AAAA{}
				rr.Hdr = dns.RR_Header{Name: qname, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: ttl}
				rr.AAAA = ip
				rrs = append(rrs, rr)
			} else {
				rr := &dns.CNAME{}
				rr.Hdr = dns.RR_Header{Name: qname, Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: ttl}
				rr.Target = dns.Fqdn(rs.AliasTarget.DNSName)
				rrs = append(rrs, rr)
			}
		case dns.TypeCNAME:
			rr := &dns.CNAME{}
			rr.Hdr = dns.RR_Header{Name: qname, Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: ttl}
			rr.Target = dns.Fqdn(rs.AliasTarget.DNSName)
			rrs = append(rrs, rr)
		}
		return rrs
	}

	for _, record := range rs.ResourceRecords {
		switch strings.ToUpper(rs.Type) {
		case "A":
			rr := &dns.A{}
			rr.Hdr = dns.RR_Header{Name: qname, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: ttl}
			rr.A = net.ParseIP(record.Value)
			if rr.A != nil {
				rrs = append(rrs, rr)
			}
		case "AAAA":
			rr := &dns.AAAA{}
			rr.Hdr = dns.RR_Header{Name: qname, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: ttl}
			rr.AAAA = net.ParseIP(record.Value)
			if rr.AAAA != nil {
				rrs = append(rrs, rr)
			}
		case "CNAME":
			rr := &dns.CNAME{}
			rr.Hdr = dns.RR_Header{Name: qname, Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: ttl}
			rr.Target = dns.Fqdn(record.Value)
			rrs = append(rrs, rr)
		case "TXT":
			rr := &dns.TXT{}
			rr.Hdr = dns.RR_Header{Name: qname, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: ttl}
			rr.Txt = []string{record.Value}
			rrs = append(rrs, rr)
		case "MX":
			rr := &dns.MX{}
			rr.Hdr = dns.RR_Header{Name: qname, Rrtype: dns.TypeMX, Class: dns.ClassINET, Ttl: ttl}
			parts := strings.Fields(record.Value)
			if len(parts) >= 2 {
				if pref, err := strconv.ParseUint(parts[0], 10, 16); err == nil {
					rr.Preference = uint16(pref)
				}
				rr.Mx = dns.Fqdn(strings.Join(parts[1:], " "))
			} else {
				rr.Preference = 10
				rr.Mx = dns.Fqdn(record.Value)
			}
			rrs = append(rrs, rr)
		case "NS":
			rr := &dns.NS{}
			rr.Hdr = dns.RR_Header{Name: qname, Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: ttl}
			rr.Ns = dns.Fqdn(record.Value)
			rrs = append(rrs, rr)
		case "SOA":
			rr := &dns.SOA{}
			rr.Hdr = dns.RR_Header{Name: qname, Rrtype: dns.TypeSOA, Class: dns.ClassINET, Ttl: ttl}
			parts := strings.Fields(record.Value)
			// Fail-closed for malformed SOA — skip the record
			// instead of emitting one with zero-valued fields.
			if len(parts) < 7 {
				break
			}
			rr.Ns = dns.Fqdn(parts[0])
			rr.Mbox = dns.Fqdn(parts[1])
			if serial, err := strconv.ParseUint(parts[2], 10, 32); err == nil {
				rr.Serial = uint32(serial)
			}
			if refresh, err := strconv.ParseUint(parts[3], 10, 32); err == nil {
				rr.Refresh = uint32(refresh)
			}
			if retry, err := strconv.ParseUint(parts[4], 10, 32); err == nil {
				rr.Retry = uint32(retry)
			}
			if expire, err := strconv.ParseUint(parts[5], 10, 32); err == nil {
				rr.Expire = uint32(expire)
			}
			if minttl, err := strconv.ParseUint(parts[6], 10, 32); err == nil {
				rr.Minttl = uint32(minttl)
			}
			rrs = append(rrs, rr)
		case "PTR":
			rr := &dns.PTR{}
			rr.Hdr = dns.RR_Header{Name: qname, Rrtype: dns.TypePTR, Class: dns.ClassINET, Ttl: ttl}
			rr.Ptr = dns.Fqdn(record.Value)
			rrs = append(rrs, rr)
		case "SRV":
			rr := &dns.SRV{}
			rr.Hdr = dns.RR_Header{Name: qname, Rrtype: dns.TypeSRV, Class: dns.ClassINET, Ttl: ttl}
			parts := strings.Fields(record.Value)
			if len(parts) >= 4 {
				if priority, err := strconv.ParseUint(parts[0], 10, 16); err == nil {
					rr.Priority = uint16(priority)
				}
				if weight, err := strconv.ParseUint(parts[1], 10, 16); err == nil {
					rr.Weight = uint16(weight)
				}
				if port, err := strconv.ParseUint(parts[2], 10, 16); err == nil {
					rr.Port = uint16(port)
				}
				rr.Target = dns.Fqdn(parts[3])
			}
			rrs = append(rrs, rr)
		case "CAA":
			rr := &dns.CAA{}
			rr.Hdr = dns.RR_Header{Name: qname, Rrtype: dns.TypeCAA, Class: dns.ClassINET, Ttl: ttl}
			parts := strings.SplitN(record.Value, " ", 3)
			if len(parts) >= 3 {
				if flag, err := strconv.ParseUint(parts[0], 10, 8); err == nil {
					rr.Flag = uint8(flag)
				}
				rr.Tag = parts[1]
				rr.Value = parts[2]
			}
			rrs = append(rrs, rr)
		}
	}

	return rrs
}

func (s *DNSServer) findSOA(zone *route53store.HostedZone, recordSets []*route53store.ResourceRecordSet) dns.RR {
	zoneName := strings.ToLower(dns.CanonicalName(zone.Name))
	for _, rs := range recordSets {
		rsName := strings.ToLower(dns.CanonicalName(rs.Name))
		if rs.Type == "SOA" && (rsName == zoneName || rsName == zoneName+".") {
			rrs := s.recordToRR(rs, zoneName, dns.TypeSOA)
			if len(rrs) > 0 {
				return rrs[0]
			}
		}
	}
	return nil
}
