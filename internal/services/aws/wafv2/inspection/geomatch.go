// GeoMatchStatement and AsnMatchStatement evaluation. Both resolve the
// request's address — the client IP by default, or with a
// ForwardedIPConfig the first address of the configured header, using
// the same fallback semantics as the IP-set path: a request whose
// configured header is absent leaves the rule unapplied entirely (a
// state that propagates through NOT and logical statements unchanged),
// and the fallback behaviour governs only a present header without a
// parseable address.
package inspection

import (
	"strconv"
	"strings"

	wafstore "vorpalstacks/internal/store/aws/waf"
)

// addressState classifies the outcome of the statement address
// resolution.
type addressState int

const (
	addressNotApplied      addressState = iota // the rule is not applied
	addressFallbackMatched                     // the fallback matched without an address
	addressResolved                            // an address was resolved
)

// statementAddress resolves the address a geo or ASN statement
// inspects.
func (ctx *evalCtx) statementAddress(cfg *wafstore.ForwardedIPConfig) (addrString string, state addressState) {
	if cfg == nil {
		if parseIP(ctx.req.SourceIP).IsValid() {
			return ctx.req.SourceIP, addressResolved
		}
		return "", addressNotApplied
	}
	ip, matched := ctx.req.forwardedIP(cfg)
	if !matched {
		return "", addressNotApplied
	}
	if ip == "" {
		// The forwarded-IP helper signals the MATCH fallback as a
		// match without an address.
		return "", addressFallbackMatched
	}
	return ip, addressResolved
}

// geoMatch reports whether the request's resolved country is listed in
// the statement's country codes (compared case-insensitively; the API
// documents upper-case alpha-2 codes). A forwarded configuration whose
// header is absent leaves the rule unapplied rather than non-matching.
func (ctx *evalCtx) geoMatch(stmt *wafstore.GeoMatchStatement) statementOutcome {
	addrString, state := ctx.statementAddress(stmt.ForwardedIPConfig)
	if state == addressNotApplied {
		return outcomeNotApplied
	}
	if state == addressFallbackMatched {
		return outcomeMatch
	}
	addr := parseIP(addrString)
	if !addr.IsValid() {
		return outcomeNoMatch
	}
	country, ok := CountryForIP(addr)
	if !ok {
		return outcomeNoMatch
	}
	for _, code := range stmt.CountryCodes {
		if strings.EqualFold(code, country) {
			return outcomeMatch
		}
	}
	return outcomeNoMatch
}

// asnMatch reports whether the request's origin ASN is listed in the
// statement's AsnList. A forwarded configuration whose header is
// absent leaves the rule unapplied rather than non-matching.
func (ctx *evalCtx) asnMatch(stmt *wafstore.AsnMatchStatement) statementOutcome {
	addrString, state := ctx.statementAddress(stmt.ForwardedIPConfig)
	if state == addressNotApplied {
		return outcomeNotApplied
	}
	if state == addressFallbackMatched {
		return outcomeMatch
	}
	addr := parseIP(addrString)
	if !addr.IsValid() {
		return outcomeNoMatch
	}
	asn, ok := ASNForIP(addr)
	if !ok {
		return outcomeNoMatch
	}
	for _, listed := range stmt.AsnList {
		value, err := strconv.ParseUint(listed, 10, 32)
		if err != nil {
			continue
		}
		if uint32(value) == asn {
			return outcomeMatch
		}
	}
	return outcomeNoMatch
}
