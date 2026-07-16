package request

// AppSync v1 GraphQL API parser — handles /v1/apis/... paths.
// This is the largest path group, covering the core GraphQL API lifecycle
// (CRUD), nested resources (datasources, types, resolvers, functions, apikeys),
// schema management, caching, environment variables, and the GraphQL execution
// data-plane endpoint (sentinel operation).

import (
	"net/http"
	"strings"
)

// extractAppSyncV1ApisOperation handles v1 GraphQL API paths (/v1/apis/...).
// Includes nested resources: datasources, types, resolvers, functions, apikeys,
// ApiCaches, FlushCache, schemacreation, schema, environmentVariables,
// sourceApiAssociations, and the GraphQL execution sentinel.
func extractAppSyncV1ApisOperation(path, method string) string {
	rel := strings.TrimPrefix(path, "/v1/apis")
	if rel == "" || rel == "/" {
		switch method {
		case http.MethodGet:
			return "ListGraphqlApis"
		case http.MethodPost:
			return "CreateGraphqlApi"
		}
		return ""
	}

	rel = strings.TrimPrefix(rel, "/")
	parts := strings.Split(rel, "/")

	if len(parts) == 1 {
		switch method {
		case http.MethodGet:
			return "GetGraphqlApi"
		case http.MethodPost:
			return "UpdateGraphqlApi"
		case http.MethodDelete:
			return "DeleteGraphqlApi"
		}
		return ""
	}

	if len(parts) >= 2 {
		switch parts[1] {
		case "graphql":
			// Data-plane GraphQL execution — not an SDK operation, but a sentinel
			// that bypasses the REST-JSON pipeline and routes to the GraphQL engine.
			if method == http.MethodPost {
				return "GraphQLExecution"
			}

		case "datasources":
			if len(parts) == 2 {
				switch method {
				case http.MethodGet:
					return "ListDataSources"
				case http.MethodPost:
					return "CreateDataSource"
				}
			} else {
				switch method {
				case http.MethodGet:
					return "GetDataSource"
				case http.MethodPost:
					return "UpdateDataSource"
				case http.MethodDelete:
					return "DeleteDataSource"
				}
			}

		case "types":
			if len(parts) == 2 {
				switch method {
				case http.MethodGet:
					return "ListTypes"
				case http.MethodPost:
					return "CreateType"
				}
			} else if len(parts) >= 4 && parts[3] == "resolvers" {
				// Resolvers are nested under types: /v1/apis/{apiId}/types/{typeName}/resolvers[/{fieldName}]
				if len(parts) == 4 {
					switch method {
					case http.MethodGet:
						return "ListResolvers"
					case http.MethodPost:
						return "CreateResolver"
					}
				} else {
					switch method {
					case http.MethodGet:
						return "GetResolver"
					case http.MethodPost:
						return "UpdateResolver"
					case http.MethodDelete:
						return "DeleteResolver"
					}
				}
			} else {
				switch method {
				case http.MethodGet:
					return "GetType"
				case http.MethodPost:
					return "UpdateType"
				case http.MethodDelete:
					return "DeleteType"
				}
			}

		case "functions":
			// Functions share the /resolvers/ sub-path for ListResolversByFunction.
			if len(parts) == 2 {
				switch method {
				case http.MethodGet:
					return "ListFunctions"
				case http.MethodPost:
					return "CreateFunction"
				}
			} else if len(parts) >= 4 && parts[3] == "resolvers" {
				if method == http.MethodGet {
					return "ListResolversByFunction"
				}
			} else {
				switch method {
				case http.MethodGet:
					return "GetFunction"
				case http.MethodPost:
					return "UpdateFunction"
				case http.MethodDelete:
					return "DeleteFunction"
				}
			}

		case "apikeys":
			if len(parts) == 2 {
				switch method {
				case http.MethodGet:
					return "ListApiKeys"
				case http.MethodPost:
					return "CreateApiKey"
				}
			} else {
				switch method {
				case http.MethodPost:
					return "UpdateApiKey"
				case http.MethodDelete:
					return "DeleteApiKey"
				}
			}

		case "ApiCaches":
			// Case-sensitive: "ApiCaches" with uppercase A and C.
			// UpdateApiCache uses a literal "/update" suffix, not a {cacheId} parameter.
			if len(parts) == 2 {
				switch method {
				case http.MethodGet:
					return "GetApiCache"
				case http.MethodPost:
					return "CreateApiCache"
				case http.MethodDelete:
					return "DeleteApiCache"
				}
			} else if len(parts) == 3 && parts[2] == "update" {
				if method == http.MethodPost {
					return "UpdateApiCache"
				}
			}

		case "FlushCache":
			// Case-sensitive: "FlushCache" with uppercase F and C.
			// Path is /v1/apis/{apiId}/FlushCache, NOT /v1/apis/{apiId}/ApiCaches/flush.
			if method == http.MethodDelete {
				return "FlushApiCache"
			}

		case "schemacreation":
			// Case-sensitive: lowercase "schemacreation".
			switch method {
			case http.MethodGet:
				return "GetSchemaCreationStatus"
			case http.MethodPost:
				return "StartSchemaCreation"
			}

		case "schema":
			// GetIntrospectionSchema returns raw SDL bytes, not JSON-wrapped.
			if method == http.MethodGet {
				return "GetIntrospectionSchema"
			}

		case "environmentVariables":
			// PutGraphqlApiEnvironmentVariables uses PUT (the only PUT in the entire API).
			switch method {
			case http.MethodGet:
				return "GetGraphqlApiEnvironmentVariables"
			case http.MethodPut:
				return "PutGraphqlApiEnvironmentVariables"
			}

		case "sourceApiAssociations":
			if method == http.MethodGet {
				return "ListSourceApiAssociations"
			}
		}
	}

	return ""
}

// extractAppSyncV1ApisPathParams extracts path parameters for v1 GraphQL API requests.
// Handles nested resources: datasources, types, resolvers, functions, apikeys.
func extractAppSyncV1ApisPathParams(path string, params map[string]interface{}) {
	rel := strings.TrimPrefix(path, "/v1/apis")
	if rel == "" || rel == "/" {
		return
	}

	rel = strings.TrimPrefix(rel, "/")
	parts := strings.Split(rel, "/")

	if len(parts) >= 1 && parts[0] != "" {
		params["apiId"] = parts[0]
	}

	if len(parts) >= 2 {
		switch parts[1] {
		case "datasources":
			if len(parts) >= 3 && parts[2] != "" {
				params["name"] = parts[2]
			}
		case "types":
			if len(parts) >= 3 && parts[2] != "" {
				params["typeName"] = parts[2]
				if len(parts) >= 5 && parts[3] == "resolvers" && parts[4] != "" {
					params["fieldName"] = parts[4]
				}
			}
		case "functions":
			if len(parts) >= 3 && parts[2] != "" {
				params["functionId"] = parts[2]
			}
		case "apikeys":
			if len(parts) >= 3 && parts[2] != "" {
				params["id"] = parts[2]
			}
		}
	}
}
