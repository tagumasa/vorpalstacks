// Package serviceports defines all port constants for vorpalstacks.
// Every consumer (server, CLI, SDK tests) imports from here.
// Changing a constant migrates that port everywhere.
package serviceports

// Fixed listener ports — change these to migrate all ports.
const (
	HTTP       = 50080
	GRPCWeb    = 50090
	TLS        = 50443
	Route53DNS = 50088
	Route53HC  = 50089 // default health check target port (not a listener)
)

// Reserved service port slots (50101-50107) — these are NOT used for
// actual listener binding. They exist for two purposes:
//  1. Migration reference: CLI "port list" shows these as the old defaults
//     so operators can map old 80xx ports to new 50xxx ports.
//  2. Neptune only: when Neptune has exactly 1 cluster (common case), the
//     first cluster uses port 50107 as its individual port instead of
//     consuming a dynamic-range slot. This keeps Neptune predictable.
//
// All other services (S3, APIGW, Cognito, CloudFront, Lambda, AppSync)
// default to FQDN mode and never bind individual ports unless explicitly
// switched. When switched to Individual, portalloc allocates from the
// dynamic range (50200-50400), NOT from these reserved slots.
const (
	S3Website  = 50101 // reference only — S3 Website defaults to FQDN
	APIGateway = 50102 // reference only — API Gateway defaults to FQDN
	Cognito    = 50103 // reference only — Cognito defaults to FQDN
	CloudFront = 50104 // reference only — CloudFront defaults to FQDN
	LambdaURL  = 50105 // reference only — Lambda URL defaults to FQDN
	AppSync    = 50106 // reference only — AppSync Events defaults to FQDN
	Neptune    = 50107 // used as first-cluster port in Individual mode
)

// Dynamic port allocation range — starts AFTER reserved service ports.
const (
	DynamicRangeStart = 50200
	DynamicRangeEnd   = 50400
)
