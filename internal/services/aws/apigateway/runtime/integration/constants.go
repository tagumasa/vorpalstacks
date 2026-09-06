package integration

// DefaultIntegrationTimeoutMillis is AWS's documented default for the
// integration timeoutInMillis member: 29 seconds, just under API Gateway's
// 30-second maximum. The services-layer integration Core and the HTTP
// executor both reference this single definition.
const DefaultIntegrationTimeoutMillis = 29000
