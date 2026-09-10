package apigateway

// Service limit values, each defined once here and referenced by name at
// every enforcement site; raw numbers never appear elsewhere.

const (
	// ApiKeyValueMinLength and ApiKeyValueMaxLength bound the explicit
	// value of an API key (the value the caller supplies, e.g. an
	// imported key's Key column).
	ApiKeyValueMinLength = 20
	ApiKeyValueMaxLength = 128

	// ApiKeyNameMaxLength bounds an API key's name.
	ApiKeyNameMaxLength = 1024

	// ImportDefinitionMaxBytes bounds an imported API definition body
	// (OpenAPI JSON or YAML) for ImportRestApi and PutRestApi.
	ImportDefinitionMaxBytes = 6 << 20

	// EdgeHostedZoneID is the Route 53 hosted zone that serves
	// CloudFront-backed (edge) custom domain names.
	EdgeHostedZoneID = "Z2FDTNDATAQYW2"

	// RegionalHostedZoneID is the Route 53 hosted zone that serves
	// regional execute-api custom domain names.
	RegionalHostedZoneID = "Z2OJLY3DKBEYEU"

	// AccountDefaultRateLimit and AccountDefaultBurstLimit are the
	// documented account-level throttle defaults: 10,000 requests per
	// second with a token-bucket capacity of 5,000, per account per
	// region — the limits GetAccount reports for an account that has
	// never stored a configuration.
	AccountDefaultRateLimit  = 10000.0
	AccountDefaultBurstLimit = 5000
)
