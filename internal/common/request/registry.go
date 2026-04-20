package request

import "net/http"

// RESTServiceParser defines the interface for REST-protocol service parsers.
// Each REST service (Lambda, API Gateway, Scheduler, etc.) implements this
// interface to register its path matching, operation extraction, and
// path parameter extraction logic without modifying the central orchestrator.
type RESTServiceParser interface {
	// MatchPath returns true if the given URL path belongs to this service.
	MatchPath(path string) bool

	// ExtractOperation inspects the HTTP request and returns the operation name,
	// or empty string if the request does not match this service.
	ExtractOperation(r *http.Request) string

	// ExtractPathParams extracts URI-bound parameters from the request path
	// (and headers where applicable) into the params map.
	ExtractPathParams(r *http.Request, params map[string]interface{})
}

// restParsers holds the registered REST service parsers in priority order.
// The order matters: earlier parsers are checked first, and the first match wins.
var restParsers []RESTServiceParser

// RegisterRESTParser appends a REST service parser to the registry.
// Should be called from init() functions in each parser file.
func RegisterRESTParser(p RESTServiceParser) {
	restParsers = append(restParsers, p)
}

func init() {
	RegisterRESTParser(&lambdaRESTParser{})
	RegisterRESTParser(&apiGatewayRESTParser{})
	RegisterRESTParser(&schedulerRESTParser{})
	RegisterRESTParser(&sesv2RESTParser{})
	RegisterRESTParser(&route53RESTParser{})
	RegisterRESTParser(&cloudFrontRESTParser{})
	RegisterRESTParser(&neptunedataRESTParser{})
	RegisterRESTParser(&appSyncRESTParser{})
	RegisterRESTParser(&neptuneGraphRESTParser{})
}
