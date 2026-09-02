package inspection

import (
	"net/netip"
	"strings"

	wafstore "vorpalstacks/internal/store/aws/waf"
)

// parseIP parses an IPv4 or IPv6 address, returning the zero
// (invalid) Addr when the value is not a valid address.
func parseIP(s string) netip.Addr {
	addr, err := netip.ParseAddr(s)
	if err != nil {
		return netip.Addr{}
	}
	return addr
}

// parseIPPrefix parses a bare address as a host prefix or a CIDR,
// returning the zero prefix when the value is neither.
func parseIPPrefix(s string) (netip.Prefix, bool) {
	if prefix, err := netip.ParsePrefix(s); err == nil {
		return prefix, true
	}
	if addr, err := netip.ParseAddr(s); err == nil {
		return netip.PrefixFrom(addr, addr.BitLen()), true
	}
	return netip.Prefix{}, false
}

// ipSetMatches reports whether the request's client address is
// contained in the referenced IP set. With IPSetForwardedIPConfig the
// addresses are instead taken from the configured header: FIRST uses
// the first entry, LAST the last, and ANY matches when any entry is in
// the set (when the header carries more than ten addresses, only the
// last ten are inspected). A request whose configured header is absent
// leaves the rule unapplied — a state distinct from a non-match, so an
// enclosing NOT must not invert it; the fallback behaviour governs
// only a present header without a valid address in the configured
// position.
func (e *Evaluator) ipSetMatches(req *Request, stmt *wafstore.IPSetReferenceStatement) statementOutcome {
	ipSet, err := e.resolvers.IPSet(stmt.ARN)
	if err != nil || ipSet == nil {
		return outcomeNoMatch
	}
	members := compiledIPSet(ipSet)

	client := parseIP(req.SourceIP)
	if stmt.IPSetForwardedIPConfig == nil {
		if !client.IsValid() {
			return outcomeNoMatch
		}
		return boolOutcome(prefixContains(members, client))
	}
	cfg := stmt.IPSetForwardedIPConfig
	if !req.hasHeader(cfg.HeaderName) {
		// The API documents an absent configured header as leaving the
		// rule unapplied entirely.
		return outcomeNotApplied
	}
	entries := forwardedIPList(req, cfg.HeaderName)
	if len(entries) == 0 {
		return boolOutcome(cfg.FallbackBehavior == "MATCH")
	}
	switch cfg.Position {
	case "FIRST":
		return boolOutcome(prefixContains(members, entries[0]))
	case "LAST":
		return boolOutcome(prefixContains(members, entries[len(entries)-1]))
	default: // ANY
		for _, entry := range lastEntries(entries, maxForwardedIPEntries) {
			if prefixContains(members, entry) {
				return outcomeMatch
			}
		}
		return outcomeNoMatch
	}
}

// maxForwardedIPEntries is the number of trailing header entries the
// ANY position inspects when the header carries more addresses.
const maxForwardedIPEntries = 10

func lastEntries(entries []netip.Addr, n int) []netip.Addr {
	if len(entries) <= n {
		return entries
	}
	return entries[len(entries)-n:]
}

func prefixContains(members []netip.Prefix, addr netip.Addr) bool {
	for _, p := range members {
		if p.Contains(addr) {
			return true
		}
	}
	return false
}

// compiledIPSet converts a stored IP set's address strings into parsed
// prefixes, skipping unparseable entries.
func compiledIPSet(set *wafstore.IPSet) []netip.Prefix {
	out := make([]netip.Prefix, 0, len(set.Addresses))
	for _, a := range set.Addresses {
		if prefix, ok := parseIPPrefix(a); ok {
			out = append(out, prefix)
		}
	}
	return out
}

// forwardedIPList returns every parseable address in the named
// forwarded-IP header, in wire order.
func forwardedIPList(req *Request, headerName string) []netip.Addr {
	var out []netip.Addr
	for _, value := range req.headerValues(headerName) {
		for _, candidate := range splitComma(value) {
			if addr := parseForwardedAddress(candidate); addr.IsValid() {
				out = append(out, addr)
			}
		}
	}
	return out
}

// parseForwardedAddress parses one comma-separated entry of a
// forwarded-IP header. Bare IPv4 and IPv6 addresses parse as-is;
// bracketed addresses (with or without a port suffix) are unbracketed;
// an IPv4 address with a port suffix drops the port. Port stripping is
// restricted to the single-colon host:port form so a bare IPv6
// address, whose colons are part of the address, always survives
// intact. Entries that remain unparseable return the zero address.
func parseForwardedAddress(candidate string) netip.Addr {
	candidate = trimASCIISpace(candidate)
	if addr := parseIP(candidate); addr.IsValid() {
		return addr
	}
	if strings.HasPrefix(candidate, "[") {
		end := strings.IndexByte(candidate, ']')
		if end > 1 {
			return parseIP(candidate[1:end])
		}
		return netip.Addr{}
	}
	if idx := lastIndexByte(candidate, ':'); idx > 0 && strings.Count(candidate, ":") == 1 {
		return parseIP(candidate[:idx])
	}
	return netip.Addr{}
}

func splitComma(s string) []string {
	var out []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			out = append(out, trimASCIISpace(s[start:i]))
			start = i + 1
		}
	}
	return out
}

func trimASCIISpace(s string) string {
	for len(s) > 0 && isASCIISpaceByte(s[0]) {
		s = s[1:]
	}
	for len(s) > 0 && isASCIISpaceByte(s[len(s)-1]) {
		s = s[:len(s)-1]
	}
	return s
}

func lastIndexByte(s string, b byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == b {
			return i
		}
	}
	return -1
}
