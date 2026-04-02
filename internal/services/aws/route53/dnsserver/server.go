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
	started    bool
	mu         sync.RWMutex
	shutdownCh chan struct{}
}

// NewDNSServer creates a new DNSServer with the given stores and bind address.
func NewDNSServer(hostedZoneStore *route53store.HostedZoneStore, recordSetStore *route53store.RecordSetStore, bindAddr string) *DNSServer {
	return &DNSServer{
		store:      hostedZoneStore,
		recordSets: recordSetStore,
		bindAddr:   bindAddr,
		shutdownCh: make(chan struct{}),
	}
}

// Start starts the DNS server on UDP and TCP port 53.
func (s *DNSServer) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.started {
		return fmt.Errorf("DNS server already started")
	}

	s.shutdownCh = make(chan struct{})

	handler := dns.NewServeMux()
	handler.HandleFunc(".", s.handleDNSRequest)

	s.udpServer = &dns.Server{
		Addr:    s.bindAddr + ":53",
		Net:     "udp",
		Handler: handler,
	}

	s.tcpServer = &dns.Server{
		Addr:    s.bindAddr + ":53",
		Net:     "tcp",
		Handler: handler,
	}

	go func() {
		if err := s.udpServer.ListenAndServe(); err != nil {
			select {
			case <-s.shutdownCh:
			default:
				logs.Error("UDP DNS server failed", logs.Err(err))
			}
		}
	}()

	go func() {
		if err := s.tcpServer.ListenAndServe(); err != nil {
			select {
			case <-s.shutdownCh:
			default:
				logs.Error("TCP DNS server failed", logs.Err(err))
			}
		}
	}()

	s.started = true
	logs.Info("DNS server started", logs.String("address", s.bindAddr+":53"))
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
	zones, err := s.store.List("", 100)
	if err != nil {
		return nil
	}

	var bestMatch *route53store.HostedZone
	bestMatchLen := 0

	for _, zone := range zones.HostedZones {
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
	if ttl == 0 {
		ttl = 300
	}

	if rs.AliasTarget != nil {
		switch qtype {
		case dns.TypeA:
			rr := &dns.A{}
			rr.Hdr = dns.RR_Header{Name: qname, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: ttl}
			rr.A = net.ParseIP(rs.AliasTarget.DNSName)
			if rr.A != nil {
				rrs = append(rrs, rr)
			}
		case dns.TypeAAAA:
			rr := &dns.AAAA{}
			rr.Hdr = dns.RR_Header{Name: qname, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: ttl}
			rr.AAAA = net.ParseIP(rs.AliasTarget.DNSName)
			if rr.AAAA != nil {
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
			if len(parts) >= 7 {
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
