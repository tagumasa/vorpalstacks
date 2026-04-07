package request

import (
	"net/http"
	"strings"
)

func IsNeptuneGraphPath(path string) bool {
	return strings.HasPrefix(path, "/graphs") ||
		strings.HasPrefix(path, "/snapshots") ||
		strings.HasPrefix(path, "/tags/arn:aws:neptune-graph") ||
		strings.HasPrefix(path, "/importtasks") ||
		strings.HasPrefix(path, "/exporttasks") ||
		strings.HasPrefix(path, "/queries") ||
		strings.HasPrefix(path, "/summary")
}

func extractNeptuneGraphOperation(r *http.Request) string {
	path := r.URL.Path
	method := r.Method

	if strings.HasPrefix(path, "/graphs") {
		parts := strings.Split(strings.Trim(path, "/"), "/")
		if len(parts) >= 2 && parts[1] != "" {
			if len(parts) >= 3 {
				switch parts[2] {
				case "start":
					if method == "POST" {
						return "StartGraph"
					}
				case "stop":
					if method == "POST" {
						return "StopGraph"
					}
				case "endpoints":
					if len(parts) >= 4 {
						if method == "GET" {
							return "GetPrivateGraphEndpoint"
						}
						if method == "DELETE" {
							return "DeletePrivateGraphEndpoint"
						}
					}
					if method == "POST" {
						return "CreatePrivateGraphEndpoint"
					}
					if method == "GET" {
						return "ListPrivateGraphEndpoints"
					}
				case "importtasks":
					if method == "POST" {
						return "StartImportTask"
					}
				}
			}
			switch method {
			case "GET":
				return "GetGraph"
			case "DELETE":
				return "DeleteGraph"
			case "PATCH":
				return "UpdateGraph"
			case "PUT":
				return "ResetGraph"
			}
			return ""
		}
		switch method {
		case "POST":
			return "CreateGraph"
		case "GET":
			return "ListGraphs"
		}
		return ""
	}

	if strings.HasPrefix(path, "/snapshots") {
		parts := strings.Split(strings.Trim(path, "/"), "/")
		if len(parts) >= 2 && parts[1] != "" {
			if len(parts) >= 3 && parts[2] == "restore" {
				if method == "POST" {
					return "RestoreGraphFromSnapshot"
				}
				return ""
			}
			switch method {
			case "GET":
				return "GetGraphSnapshot"
			case "DELETE":
				return "DeleteGraphSnapshot"
			}
			return ""
		}
		switch method {
		case "POST":
			return "CreateGraphSnapshot"
		case "GET":
			return "ListGraphSnapshots"
		}
		return ""
	}

	if strings.HasPrefix(path, "/tags/arn:aws:neptune-graph") {
		switch method {
		case "GET":
			return "ListTagsForResource"
		case "POST":
			return "TagResource"
		case "DELETE":
			return "UntagResource"
		}
		return ""
	}

	if strings.HasPrefix(path, "/importtasks") {
		parts := strings.Split(strings.Trim(path, "/"), "/")
		if len(parts) >= 2 && parts[1] != "" {
			switch method {
			case "GET":
				return "GetImportTask"
			case "DELETE":
				return "CancelImportTask"
			}
			return ""
		}
		if method == "GET" {
			return "ListImportTasks"
		}
		if method == "POST" {
			return "CreateGraphUsingImportTask"
		}
		return ""
	}

	if strings.HasPrefix(path, "/exporttasks") {
		parts := strings.Split(strings.Trim(path, "/"), "/")
		if len(parts) >= 2 && parts[1] != "" {
			switch method {
			case "GET":
				return "GetExportTask"
			case "DELETE":
				return "CancelExportTask"
			}
			return ""
		}
		if method == "GET" {
			return "ListExportTasks"
		}
		if method == "POST" {
			return "StartExportTask"
		}
		return ""
	}

	if strings.HasPrefix(path, "/queries") {
		parts := strings.Split(strings.Trim(path, "/"), "/")
		if len(parts) >= 2 && parts[1] != "" {
			switch method {
			case "GET":
				return "GetQuery"
			case "DELETE":
				return "CancelQuery"
			}
			return ""
		}
		if method == "GET" {
			return "ListQueries"
		}
		if method == "POST" {
			return "ExecuteQuery"
		}
		return ""
	}

	if path == "/summary" && method == "GET" {
		return "GetGraphSummary"
	}

	return ""
}

func extractNeptuneGraphPathParams(path string, params map[string]interface{}) {
	parts := strings.Split(strings.Trim(path, "/"), "/")

	if len(parts) >= 1 && parts[0] == "graphs" {
		if len(parts) >= 2 && parts[1] != "" {
			params["graphIdentifier"] = parts[1]
		}
		if len(parts) >= 4 && parts[2] == "endpoints" && parts[3] != "" {
			params["vpcId"] = parts[3]
		}
		return
	}

	if len(parts) >= 1 && parts[0] == "snapshots" {
		if len(parts) >= 2 && parts[1] != "" {
			params["snapshotIdentifier"] = parts[1]
		}
		return
	}

	if len(parts) >= 1 && parts[0] == "tags" {
		if len(parts) >= 2 && parts[1] != "" {
			params["resourceArn"] = strings.Join(parts[1:], "/")
		}
		return
	}

	if len(parts) >= 1 && parts[0] == "importtasks" {
		if len(parts) >= 2 && parts[1] != "" {
			params["taskIdentifier"] = parts[1]
		}
		return
	}

	if len(parts) >= 1 && parts[0] == "exporttasks" {
		if len(parts) >= 2 && parts[1] != "" {
			params["taskIdentifier"] = parts[1]
		}
		return
	}

	if len(parts) >= 1 && parts[0] == "queries" {
		if len(parts) >= 2 && parts[1] != "" {
			params["queryId"] = parts[1]
		}
		return
	}
}
