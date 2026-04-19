package lambda

import (
	"vorpalstacks/internal/store/aws/common"
)

// FunctionStoreInterface defines operations for managing Lambda functions.
type FunctionStoreInterface interface {
	Create(function *Function) (*Function, error)
	Get(functionName string) (*Function, error)
	GetByArn(arn string) (*Function, error)
	Update(function *Function) error
	UpdateAtomically(functionName string, modifier func(*Function) error) (*Function, error)
	Delete(functionName string) error
	List(opts common.ListOptions) (*common.ListResult[Function], error)
	ListByPrefix(prefix string) ([]*Function, error)

	PublishVersion(function *Function, description string) (*Version, error)
	GetVersion(functionName, version string) (*Version, error)
	DeleteVersion(functionName, version string) error

	CreateAlias(function *Function, alias *Alias) (*Alias, error)
	GetAlias(functionName, aliasName string) (*Alias, error)
	UpdateAlias(function *Function, alias *Alias) error
	CreateAliasAtomically(functionName string, creator func(*Function) (*Alias, error)) (*Alias, error)
	UpdateAliasAtomically(functionName, aliasName string, modifier func(*Function, *Alias) error) (*Alias, error)
	DeleteAlias(functionName, aliasName string) error

	AddPolicy(function *Function, policy *FunctionPolicy) error
	RemovePolicy(functionName, statementId string) error
	GetPolicy(functionName string) ([]FunctionPolicy, error)

	SetReservedConcurrency(functionName string, concurrency *int64) error
	GetReservedConcurrency(functionName string) (*int64, error)

	SetContainerInfo(functionName, containerID, containerImageID string) error

	SetFunctionUrlConfig(functionName string, config *FunctionUrlConfig) error
	GetFunctionUrlConfig(functionName string) (*FunctionUrlConfig, error)
	DeleteFunctionUrlConfig(functionName string) error

	ResolveQualifier(functionName, qualifier string) (*Function, *Version, *Alias, error)

	SetProvisionedConcurrency(functionName, qualifier string, concurrentExecutions int32) error
	GetProvisionedConcurrency(functionName, qualifier string) (*ProvisionedConcurrencyConfig, error)
	DeleteProvisionedConcurrency(functionName, qualifier string) error
	ListProvisionedConcurrency(functionName string) ([]ProvisionedConcurrencyConfig, error)

	SetEventInvokeConfig(functionName, qualifier string, config *EventInvokeConfig) error
	GetEventInvokeConfig(functionName, qualifier string) (*EventInvokeConfig, error)
	DeleteEventInvokeConfig(functionName, qualifier string) error
	ListEventInvokeConfigs(functionName string) ([]EventInvokeConfig, error)
}

// EventSourceStoreInterface defines operations for managing Lambda event source mappings.
type EventSourceStoreInterface interface {
	Create(mapping *EventSourceMapping) (*EventSourceMapping, error)
	Get(uuid string) (*EventSourceMapping, error)
	Update(mapping *EventSourceMapping) error
	Delete(uuid string) error
	List(opts common.ListOptions) (*common.ListResult[EventSourceMapping], error)
	ListByFunction(functionArn string) ([]*EventSourceMapping, error)
	FindByEventSourceAndFunction(eventSourceArn, functionArn string) (*EventSourceMapping, error)
	SetState(uuid, state, reason string) error
	Activate(uuid string) error
	Deactivate(uuid string) error
}

// LayerStoreInterface defines operations for managing Lambda layers.
type LayerStoreInterface interface {
	Create(layer *Layer) (*Layer, error)
	Get(layerName string) (*Layer, error)
	Delete(layerName string) error
	List(opts common.ListOptions) (*common.ListResult[Layer], error)
	Update(layer *Layer) error

	PublishVersion(layer *Layer, version *LayerVersion) (*LayerVersion, error)
	GetVersion(layerName string, versionNumber int64) (*LayerVersion, error)
	GetVersionByArn(arn string) (*LayerVersion, error)
	DeleteVersion(layerName string, versionNumber int64) error
	ListVersions(layerName string, opts common.ListOptions) (*common.ListResult[LayerVersion], error)

	AddPolicy(layer *Layer, versionNumber int64, policy *LayerPolicy) error
	RemovePolicy(layerName string, versionNumber int64, statementId string) error

	ListByCompatibleRuntime(runtime Runtime) ([]*Layer, error)
}

// LambdaStoreInterface defines access to all Lambda stores.
type LambdaStoreInterface interface {
	Functions() FunctionStoreInterface
	Layers() LayerStoreInterface
	EventSources() EventSourceStoreInterface
}
