package neptunegraph

import (
	"time"
)

// Graph represents a NeptuneGraph graph resource with its configuration and metadata.
type Graph struct {
	Id                        string              `json:"id"`
	Name                      string              `json:"name"`
	Arn                       string              `json:"arn,omitempty"`
	Status                    string              `json:"status"`
	Endpoint                  string              `json:"endpoint,omitempty"`
	ProvisionedMemory         *int32              `json:"provisionedMemory"`
	ReplicaCount              *int32              `json:"replicaCount"`
	DeletionProtection        bool                `json:"deletionProtection"`
	PublicConnectivity        bool                `json:"publicConnectivity"`
	KmsKeyIdentifier          string              `json:"kmsKeyIdentifier,omitempty"`
	VectorSearchConfiguration *VectorSearchConfig `json:"vectorSearchConfiguration,omitempty"`
	BuildNumber               string              `json:"buildNumber,omitempty"`
	SourceSnapshotId          string              `json:"sourceSnapshotId,omitempty"`
	StatusReason              string              `json:"statusReason,omitempty"`
	CreateTime                *time.Time          `json:"createTime,omitempty"`
	AccountID                 string              `json:"accountId,omitempty"`
	Region                    string              `json:"region,omitempty"`
}

// VectorSearchConfig holds the dimension for vector search configuration on a graph.
type VectorSearchConfig struct {
	Dimension int32 `json:"dimension"`
}

// GraphSnapshot represents a point-in-time snapshot of a NeptuneGraph graph.
type GraphSnapshot struct {
	Id                 string     `json:"id"`
	Name               string     `json:"name"`
	Arn                string     `json:"arn,omitempty"`
	Status             string     `json:"status"`
	KmsKeyIdentifier   string     `json:"kmsKeyIdentifier,omitempty"`
	SourceGraphId      string     `json:"sourceGraphId"`
	SnapshotCreateTime *time.Time `json:"snapshotCreateTime,omitempty"`
	AccountID          string     `json:"accountId,omitempty"`
	Region             string     `json:"region,omitempty"`
}

// PrivateGraphEndpoint represents a VPC-private endpoint for accessing a NeptuneGraph graph.
type PrivateGraphEndpoint struct {
	GraphId       string   `json:"graphId"`
	VpcId         string   `json:"vpcId"`
	VpcEndpointId string   `json:"vpcEndpointId,omitempty"`
	SubnetIds     []string `json:"subnetIds"`
	Status        string   `json:"status"`
	AccountID     string   `json:"accountId,omitempty"`
	Region        string   `json:"region,omitempty"`
}

// ExportFilter defines edge and vertex filters applied during a graph export task.
type ExportFilter struct {
	EdgeFilter   map[string]ExportFilterElement `json:"edgeFilter,omitempty"`
	VertexFilter map[string]ExportFilterElement `json:"vertexFilter,omitempty"`
}

// ExportFilterElement defines property-level filter attributes for a single edge or vertex label.
type ExportFilterElement struct {
	Properties map[string]ExportFilterPropertyAttributes `json:"properties"`
}

// ExportFilterPropertyAttributes controls how individual properties are handled during export.
type ExportFilterPropertyAttributes struct {
	MultiValueHandling string  `json:"multiValueHandling,omitempty"`
	OutputType         *string `json:"outputType,omitempty"`
	SourcePropertyName *string `json:"sourcePropertyName,omitempty"`
}

// ImportOptions holds format-specific options for an import task.
type ImportOptions struct {
	Neptune *NeptuneImportOptions `json:"neptune,omitempty"`
}

// NeptuneImportOptions contains Neptune-specific import configuration such as S3 paths and label preservation.
type NeptuneImportOptions struct {
	S3ExportPath                string `json:"s3ExportPath"`
	S3ExportKmsKeyId            string `json:"s3ExportKmsKeyId"`
	PreserveDefaultVertexLabels *bool  `json:"preserveDefaultVertexLabels,omitempty"`
	PreserveEdgeIds             *bool  `json:"preserveEdgeIds,omitempty"`
}

// GraphDataSummary contains structural statistics about a graph's nodes, edges, and properties.
type GraphDataSummary struct {
	NumNodes                *int64             `json:"numNodes,omitempty"`
	NumEdges                *int64             `json:"numEdges,omitempty"`
	NumNodeLabels           *int64             `json:"numNodeLabels,omitempty"`
	NumEdgeLabels           *int64             `json:"numEdgeLabels,omitempty"`
	NumNodeProperties       *int64             `json:"numNodeProperties,omitempty"`
	NumEdgeProperties       *int64             `json:"numEdgeProperties,omitempty"`
	TotalNodePropertyValues *int64             `json:"totalNodePropertyValues,omitempty"`
	TotalEdgePropertyValues *int64             `json:"totalEdgePropertyValues,omitempty"`
	NodeLabels              []string           `json:"nodeLabels,omitempty"`
	EdgeLabels              []string           `json:"edgeLabels,omitempty"`
	NodeProperties          []map[string]int64 `json:"nodeProperties,omitempty"`
	EdgeProperties          []map[string]int64 `json:"edgeProperties,omitempty"`
	NodeStructures          []NodeStructure    `json:"nodeStructures,omitempty"`
	EdgeStructures          []EdgeStructure    `json:"edgeStructures,omitempty"`
}

// NodeStructure summarises the outgoing edge labels and properties for a node label group.
type NodeStructure struct {
	Count                      *int64   `json:"count,omitempty"`
	DistinctOutgoingEdgeLabels []string `json:"distinctOutgoingEdgeLabels,omitempty"`
	NodeProperties             []string `json:"nodeProperties,omitempty"`
}

// EdgeStructure summarises the properties associated with an edge label group.
type EdgeStructure struct {
	Count          *int64   `json:"count,omitempty"`
	EdgeProperties []string `json:"edgeProperties,omitempty"`
}

// QueryRecord stores metadata for an executed or in-flight graph query.
type QueryRecord struct {
	Id          string    `json:"id"`
	QueryString string    `json:"queryString"`
	Language    string    `json:"language,omitempty"`
	State       string    `json:"state"`
	GraphId     string    `json:"graphId,omitempty"`
	StartedAt   time.Time `json:"startedAt,omitempty"`
	Elapsed     int32     `json:"elapsed,omitempty"`
	Waited      int32     `json:"waited,omitempty"`
}

// ImportTask represents a NeptuneGraph bulk import task with its configuration and status.
type ImportTask struct {
	TaskId                    string              `json:"taskId"`
	GraphId                   string              `json:"graphId,omitempty"`
	Status                    string              `json:"status"`
	Format                    string              `json:"format,omitempty"`
	ParquetType               string              `json:"parquetType,omitempty"`
	Source                    string              `json:"source"`
	RoleArn                   string              `json:"roleArn"`
	AttemptNumber             int32               `json:"attemptNumber,omitempty"`
	ImportOptions             *ImportOptions      `json:"importOptions,omitempty"`
	ImportTaskDetails         *ImportTaskDetails  `json:"importTaskDetails,omitempty"`
	StatusReason              string              `json:"statusReason,omitempty"`
	GraphName                 string              `json:"graphName,omitempty"`
	DeletionProtection        bool                `json:"deletionProtection,omitempty"`
	KmsKeyIdentifier          string              `json:"kmsKeyIdentifier,omitempty"`
	PublicConnectivity        bool                `json:"publicConnectivity,omitempty"`
	ReplicaCount              *int32              `json:"replicaCount,omitempty"`
	VectorSearchConfiguration *VectorSearchConfig `json:"vectorSearchConfiguration,omitempty"`
	MinProvisionedMemory      *int32              `json:"minProvisionedMemory,omitempty"`
	MaxProvisionedMemory      *int32              `json:"maxProvisionedMemory,omitempty"`
	BlankNodeHandling         string              `json:"blankNodeHandling,omitempty"`
	FailOnError               bool                `json:"failOnError,omitempty"`
	StartTime                 *time.Time          `json:"startTime,omitempty"`
}

// ImportTaskDetails holds progress metrics and error information for an import task.
type ImportTaskDetails struct {
	ProgressPercentage   *int32     `json:"progressPercentage"`
	StartTime            *time.Time `json:"startTime"`
	TimeElapsedSeconds   *int64     `json:"timeElapsedSeconds"`
	StatementCount       *int64     `json:"statementCount"`
	DictionaryEntryCount *int64     `json:"dictionaryEntryCount"`
	ErrorCount           *int32     `json:"errorCount"`
	ErrorDetails         *string    `json:"errorDetails,omitempty"`
	Status               *string    `json:"status"`
}

// ExportTask represents a NeptuneGraph bulk export task with its configuration and status.
type ExportTask struct {
	TaskId            string             `json:"taskId"`
	GraphId           string             `json:"graphId,omitempty"`
	Status            string             `json:"status"`
	Format            string             `json:"format,omitempty"`
	ParquetType       string             `json:"parquetType,omitempty"`
	Destination       string             `json:"destination"`
	RoleArn           string             `json:"roleArn"`
	KmsKeyIdentifier  string             `json:"kmsKeyIdentifier"`
	StatusReason      string             `json:"statusReason,omitempty"`
	ExportFilter      *ExportFilter      `json:"exportFilter,omitempty"`
	ExportTaskDetails *ExportTaskDetails `json:"exportTaskDetails,omitempty"`
	StartTime         *time.Time         `json:"startTime,omitempty"`
}

// ExportTaskDetails holds progress metrics for an export task.
type ExportTaskDetails struct {
	ProgressPercentage *int32     `json:"progressPercentage"`
	StartTime          *time.Time `json:"startTime"`
	TimeElapsedSeconds *int64     `json:"timeElapsedSeconds"`
	NumEdgesWritten    *int64     `json:"numEdgesWritten,omitempty"`
	NumVerticesWritten *int64     `json:"numVerticesWritten,omitempty"`
}
