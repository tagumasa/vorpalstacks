package dnsserver

import (
	"crypto/rand"
	"math"
	"math/big"
	"net"
)

// haversine calculates the great-circle distance between two points
// (lat1,lon1) and (lat2,lon2) in kilometres.
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const earthKm = 6371.0
	rlat1 := lat1 * math.Pi / 180.0
	rlat2 := lat2 * math.Pi / 180.0
	dlat := (lat2 - lat1) * math.Pi / 180.0
	dlon := (lon2 - lon1) * math.Pi / 180.0
	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(rlat1)*math.Cos(rlat2)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthKm * c
}

// awsRegionCoords maps AWS region codes to approximate geographic
// coordinates (latitude, longitude). Used for GeoProximity routing
// when a record specifies AWSRegion instead of explicit Coordinates.
var awsRegionCoords = map[string][2]float64{
	"us-east-1":      {38.13, -78.45},
	"us-east-2":      {39.96, -83.00},
	"us-west-1":      {37.35, -121.96},
	"us-west-2":      {45.52, -122.68},
	"ap-northeast-1": {35.69, 139.69},
	"ap-northeast-2": {37.57, 126.98},
	"ap-northeast-3": {34.69, 135.50},
	"ap-southeast-1": {1.35, 103.82},
	"ap-southeast-2": {-33.86, 151.21},
	"ap-south-1":     {19.08, 72.88},
	"ca-central-1":   {45.42, -75.69},
	"eu-central-1":   {50.11, 8.68},
	"eu-west-1":      {53.35, -6.26},
	"eu-west-2":      {51.51, -0.13},
	"eu-west-3":      {48.86, 2.35},
	"eu-north-1":     {59.33, 18.07},
	"eu-south-1":     {41.90, 12.50},
	"sa-east-1":      {-23.53, -46.63},
	"me-south-1":     {26.23, 50.59},
	"af-south-1":     {-33.92, 18.42},
}

func awsRegionCoordinates(region string) (float64, float64, bool) {
	coords, ok := awsRegionCoords[region]
	if !ok {
		return 0, 0, false
	}
	return coords[0], coords[1], true
}

// lookupIPCoordinates maps a querier IP to approximate coordinates.
// For private/local IPs (edge/on-premises platform), it returns the
// default data centre coordinates. For public IPs, it attempts a
// best-effort region lookup; if no match, returns (0,0).
func lookupIPCoordinates(ip string) (float64, float64) {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return 0, 0
	}
	if parsed.IsPrivate() || parsed.IsLoopback() || parsed.IsLinkLocalUnicast() {
		// For local queries, use us-east-1 as default origin.
		return 38.13, -78.45
	}
	return 0, 0
}

// ipToAWSRegion maps a querier IP to an AWS region (best-effort).
// Returns empty string if no match.
func ipToAWSRegion(ip string) string {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return ""
	}
	if parsed.IsPrivate() || parsed.IsLoopback() {
		return "us-east-1"
	}
	return ""
}

// ipToCountry maps a querier IP to a country code (best-effort).
// Returns empty string if no match.
func ipToCountry(ip string) string {
	return ""
}

// cryptoRandInt63n returns a uniform random int64 in [0, n).
func cryptoRandInt63n(n int64) int64 {
	if n <= 0 {
		return 0
	}
	max := big.NewInt(n)
	r, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0
	}
	return r.Int64()
}
