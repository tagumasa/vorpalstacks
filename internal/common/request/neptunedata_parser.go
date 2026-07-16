package request

import (
	"net/http"
	"strings"
)

// IsNeptunedataPath returns true if the HTTP request path belongs to the Neptune
// Data API namespace. Used by the classifier and dispatcher to route requests
// to the neptunedata service. Uses strict matching for single-segment paths
// (e.g. "/status") and prefix matching for namespace paths (e.g. "/ml/").
func IsNeptunedataPath(path string) bool {
	if path == "/status" {
		return true
	}
	if path == "/system" {
		return true
	}
	if path == "/gremlin" || strings.HasPrefix(path, "/gremlin/") {
		return true
	}
	if path == "/opencypher" || strings.HasPrefix(path, "/opencypher/") {
		return true
	}
	if path == "/loader" || strings.HasPrefix(path, "/loader/") {
		return true
	}
	if path == "/propertygraph" || strings.HasPrefix(path, "/propertygraph/") {
		return true
	}
	if strings.HasPrefix(path, "/sparql/") {
		return true
	}
	if strings.HasPrefix(path, "/rdf/") {
		return true
	}
	if strings.HasPrefix(path, "/ml/") {
		return true
	}
	return false
}

// extractNeptunedataOperation maps an HTTP method and path to the corresponding
// Neptune Data API operation name. Returns an empty string if no match is found.
func extractNeptunedataOperation(r *http.Request) string {
	path := r.URL.Path
	method := r.Method

	switch {
	case path == "/gremlin" && method == http.MethodPost:
		return "ExecuteGremlinQuery"
	case path == "/gremlin/explain" && method == http.MethodPost:
		return "ExecuteGremlinExplainQuery"
	case path == "/gremlin/profile" && method == http.MethodPost:
		return "ExecuteGremlinProfileQuery"
	case path == "/gremlin/status" && method == http.MethodGet:
		return "ListGremlinQueries"
	case strings.HasPrefix(path, "/gremlin/status/") && method == http.MethodGet:
		return "GetGremlinQueryStatus"
	case strings.HasPrefix(path, "/gremlin/status/") && method == http.MethodDelete:
		return "CancelGremlinQuery"

	case path == "/opencypher" && method == http.MethodPost:
		return "ExecuteOpenCypherQuery"
	case path == "/opencypher/explain" && method == http.MethodPost:
		return "ExecuteOpenCypherExplainQuery"
	case path == "/opencypher/status" && method == http.MethodGet:
		return "ListOpenCypherQueries"
	case strings.HasPrefix(path, "/opencypher/status/") && method == http.MethodGet:
		return "GetOpenCypherQueryStatus"
	case strings.HasPrefix(path, "/opencypher/status/") && method == http.MethodDelete:
		return "CancelOpenCypherQuery"

	case path == "/status" && method == http.MethodGet:
		return "GetEngineStatus"
	case path == "/system" && method == http.MethodPost:
		return "ExecuteFastReset"

	case path == "/loader" && method == http.MethodPost:
		return "StartLoaderJob"
	case path == "/loader" && method == http.MethodGet:
		return "ListLoaderJobs"
	case strings.HasPrefix(path, "/loader/") && method == http.MethodGet:
		return "GetLoaderJobStatus"
	case strings.HasPrefix(path, "/loader/") && method == http.MethodDelete:
		return "CancelLoaderJob"

	case path == "/propertygraph/statistics" && method == http.MethodGet:
		return "GetPropertygraphStatistics"
	case path == "/propertygraph/statistics" && method == http.MethodPost:
		return "ManagePropertygraphStatistics"
	case path == "/propertygraph/statistics" && method == http.MethodDelete:
		return "DeletePropertygraphStatistics"
	case path == "/propertygraph/statistics/summary" && method == http.MethodGet:
		return "GetPropertygraphSummary"
	case path == "/propertygraph/stream" && method == http.MethodGet:
		return "GetPropertygraphStream"

	case path == "/sparql/statistics" && method == http.MethodGet:
		return "GetSparqlStatistics"
	case path == "/sparql/statistics" && method == http.MethodPost:
		return "ManageSparqlStatistics"
	case path == "/sparql/statistics" && method == http.MethodDelete:
		return "DeleteSparqlStatistics"
	case path == "/rdf/statistics/summary" && method == http.MethodGet:
		return "GetRDFGraphSummary"
	case path == "/sparql/stream" && method == http.MethodGet:
		return "GetSparqlStream"

	case path == "/ml/endpoints" && method == http.MethodPost:
		return "CreateMLEndpoint"
	case path == "/ml/endpoints" && method == http.MethodGet:
		return "ListMLEndpoints"
	case strings.HasPrefix(path, "/ml/endpoints/") && method == http.MethodGet:
		return "GetMLEndpoint"
	case strings.HasPrefix(path, "/ml/endpoints/") && method == http.MethodDelete:
		return "DeleteMLEndpoint"

	case path == "/ml/dataprocessing" && method == http.MethodPost:
		return "StartMLDataProcessingJob"
	case path == "/ml/dataprocessing" && method == http.MethodGet:
		return "ListMLDataProcessingJobs"
	case strings.HasPrefix(path, "/ml/dataprocessing/") && method == http.MethodGet:
		return "GetMLDataProcessingJob"
	case strings.HasPrefix(path, "/ml/dataprocessing/") && method == http.MethodDelete:
		return "CancelMLDataProcessingJob"

	case path == "/ml/modeltraining" && method == http.MethodPost:
		return "StartMLModelTrainingJob"
	case path == "/ml/modeltraining" && method == http.MethodGet:
		return "ListMLModelTrainingJobs"
	case strings.HasPrefix(path, "/ml/modeltraining/") && method == http.MethodGet:
		return "GetMLModelTrainingJob"
	case strings.HasPrefix(path, "/ml/modeltraining/") && method == http.MethodDelete:
		return "CancelMLModelTrainingJob"

	case path == "/ml/modeltransform" && method == http.MethodPost:
		return "StartMLModelTransformJob"
	case path == "/ml/modeltransform" && method == http.MethodGet:
		return "ListMLModelTransformJobs"
	case strings.HasPrefix(path, "/ml/modeltransform/") && method == http.MethodGet:
		return "GetMLModelTransformJob"
	case strings.HasPrefix(path, "/ml/modeltransform/") && method == http.MethodDelete:
		return "CancelMLModelTransformJob"
	}

	return ""
}

// extractNeptunedataPathParams extracts path parameters (queryId, loadId, id)
// from the Neptune Data API URL and writes them into the params map.
func extractNeptunedataPathParams(path string, params map[string]interface{}) {
	switch {
	case strings.HasPrefix(path, "/gremlin/status/") && len(path) > len("/gremlin/status/"):
		params["queryId"] = strings.TrimPrefix(path, "/gremlin/status/")

	case strings.HasPrefix(path, "/opencypher/status/") && len(path) > len("/opencypher/status/"):
		params["queryId"] = strings.TrimPrefix(path, "/opencypher/status/")

	case strings.HasPrefix(path, "/loader/") && len(path) > len("/loader/"):
		params["loadId"] = strings.TrimPrefix(path, "/loader/")

	case strings.HasPrefix(path, "/ml/endpoints/") && len(path) > len("/ml/endpoints/"):
		params["id"] = strings.TrimPrefix(path, "/ml/endpoints/")

	case strings.HasPrefix(path, "/ml/dataprocessing/") && len(path) > len("/ml/dataprocessing/"):
		params["id"] = strings.TrimPrefix(path, "/ml/dataprocessing/")

	case strings.HasPrefix(path, "/ml/modeltraining/") && len(path) > len("/ml/modeltraining/"):
		params["id"] = strings.TrimPrefix(path, "/ml/modeltraining/")

	case strings.HasPrefix(path, "/ml/modeltransform/") && len(path) > len("/ml/modeltransform/"):
		params["id"] = strings.TrimPrefix(path, "/ml/modeltransform/")
	}
}

// neptunedataRESTParser implements RESTServiceParser for Amazon Neptune Data.
type neptunedataRESTParser struct{}

// MatchPath returns true if the path belongs to Neptune Data (ML or query APIs).
func (p *neptunedataRESTParser) MatchPath(path string) bool {
	return IsNeptunedataPath(path)
}

// ExtractOperation returns the Neptune Data operation name, or empty if the path does not match.
func (p *neptunedataRESTParser) ExtractOperation(r *http.Request) string {
	return extractNeptunedataOperation(r)
}

// ExtractPathParams extracts URI-bound parameters from the Neptune Data request path.
func (p *neptunedataRESTParser) ExtractPathParams(r *http.Request, params map[string]interface{}) {
	extractNeptunedataPathParams(r.URL.Path, params)
}
