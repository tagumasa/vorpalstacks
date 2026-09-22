package cloudwatchlogs

import (
	"net"
	"strings"
)

// The IP address function family of the documented Logs Insights QL
// function library.

// evalIPFunctions holds the IP family's dispatch: false means the name
// is not this family's.
func evalIPFunctions(name string, f funcArgs) (interface{}, bool) {
	switch name {
	case "isvalidip":
		return net.ParseIP(f.str(0)) != nil, true
	case "isvalidipv4":
		ip := net.ParseIP(f.str(0))
		// An IPv6-mapped dotted form ("::ffff:1.2.3.4") parses with a
		// non-nil To4 yet is not an IPv4 address.
		return ip != nil && ip.To4() != nil && !strings.Contains(f.str(0), ":") && strings.Count(f.str(0), ".") == 3, true
	case "isvalidipv6":
		ip := net.ParseIP(f.str(0))
		return ip != nil && strings.Contains(f.str(0), ":"), true
	case "isipinsubnet":
		return ipInSubnet(f.str(0), f.str(1), false), true
	case "isipv4insubnet":
		return ipInSubnet(f.str(0), f.str(1), true), true
	case "isipv6insubnet":
		ip := net.ParseIP(f.str(0))
		_, cidr, err := net.ParseCIDR(f.str(1))
		return err == nil && ip != nil && strings.Contains(f.str(1), ":") && cidr.Contains(ip), true
	case "ipv4tonumber":
		ip := net.ParseIP(f.str(0))
		if ip == nil || ip.To4() == nil {
			return nil, true
		}
		v := uint32(0)
		for _, b := range ip.To4() {
			v = v<<8 | uint32(b)
		}
		return float64(v), true
	case "isprivateip":
		ip := net.ParseIP(f.str(0))
		if ip == nil || ip.To4() == nil {
			return false, true
		}
		for _, cidr := range []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"} {
			_, nw, _ := net.ParseCIDR(cidr)
			if nw.Contains(ip) {
				return true, true
			}
		}
		return false, true
	case "ispublicip":
		ip := net.ParseIP(f.str(0))
		if ip == nil {
			return false, true
		}
		private := []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "127.0.0.0/8", "169.254.0.0/16"}
		reserved := []string{"0.0.0.0/8", "192.0.2.0/24", "198.51.100.0/24", "203.0.113.0/24", "224.0.0.0/4", "240.0.0.0/4"}
		for _, cidr := range append(private, reserved...) {
			_, nw, _ := net.ParseCIDR(cidr)
			if nw.Contains(ip) {
				return false, true
			}
		}
		return true, true
	case "isreservedip":
		ip := net.ParseIP(f.str(0))
		if ip == nil {
			return false, true
		}
		reserved := []string{"0.0.0.0/8", "10.0.0.0/8", "127.0.0.0/8", "169.254.0.0/16", "172.16.0.0/12",
			"192.0.2.0/24", "192.168.0.0/16", "198.51.100.0/24", "203.0.113.0/24", "224.0.0.0/4", "240.0.0.0/4"}
		for _, cidr := range reserved {
			_, nw, _ := net.ParseCIDR(cidr)
			if nw.Contains(ip) {
				return true, true
			}
		}
		return false, true
	}
	return nil, false
}

// ipInSubnet reports whether the address is within the subnet; when v4Only
// is set the address must be IPv4 and the subnet IPv4 as well.
func ipInSubnet(addr, subnet string, v4Only bool) bool {
	ip := net.ParseIP(addr)
	_, cidr, err := net.ParseCIDR(subnet)
	if err != nil || ip == nil {
		return false
	}
	if v4Only && (ip.To4() == nil || strings.Contains(subnet, ":")) {
		return false
	}
	return cidr.Contains(ip)
}
