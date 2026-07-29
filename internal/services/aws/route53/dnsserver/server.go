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
	store      *route53store.HostedZoneStore
	recordSets *route53store.RecordSetStore
	udpServer  *dns.Server
	tcpServer  *dns.Server
	bindAddr   string
	port       int
	started    bool
	mu         sync.RWMutex
	shutdownCh chan struct{}
}

// NewDNSServer creates a new DNSServer with the given stores and bind address.
func NewDNSServer(hostedZoneStore *route53store.HostedZoneStore, recordSetStore *route53store.RecordSetStore, bindAddr string, port int) *DNSServer {
	return &DNSServer{
		store:      hostedZoneStore,
		recordSets: recordSetStore,
		bindAddr:   bindAddr,
		port:       port,
		shutdownCh: make(chan struct{}),
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

	// L9: Bind synchronously so that port conflicts are detected
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

func (s *DNSServer) handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	m.RecursionAvailable = false

	if len(r.Question) == 0 {
		// L6: Return FORMERR per DNS spec for empty question section.
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

	for _, rs := range recordSets {
		rsName := strings.ToLower(dns.CanonicalName(rs.Name))
		if rsName == qname || rsName == qname+"." {
			m.Answer = append(m.Answer, s.recordToRR(rs, qname, qtype)...)
		}
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
	// L7: Paginate through all hosted zones instead of capping at 100.
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
	// L5: Respect TTL=0 per RFC 1035 ("do not cache") instead of
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
			// L4: Fail-closed for malformed SOA — skip the record
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
