// This file holds the single deep-copy implementation for the Lambda
// store's Function and Version graphs. Every nested pointer, slice and
// map is reproduced so a copy is fully detached from the stored record —
// the async invocation paths rely on that detachment. The reflect
// field-walk tests in deepcopy_test.go fail whenever a new field is
// added without being copied here.
package lambda

import (
	types "vorpalstacks/internal/common/tags"
)

// DeepCopy returns a fully detached copy of the function, including every
// published version, alias, policy and concurrency record. The unexported
// lookup indexes are rebuilt for the copy.
func (f *Function) DeepCopy() *Function {
	if f == nil {
		return nil
	}

	result := &Function{
		FunctionName:               f.FunctionName,
		FunctionArn:                f.FunctionArn,
		Runtime:                    f.Runtime,
		Role:                       f.Role,
		Handler:                    f.Handler,
		CodeSize:                   f.CodeSize,
		CodeSha256:                 f.CodeSha256,
		CodeLocation:               f.CodeLocation,
		ImageUri:                   f.ImageUri,
		SourceCodeHash:             f.SourceCodeHash,
		Description:                f.Description,
		Timeout:                    f.Timeout,
		MemorySize:                 f.MemorySize,
		Publish:                    f.Publish,
		KMSKeyArn:                  f.KMSKeyArn,
		RevisionId:                 f.RevisionId,
		State:                      f.State,
		StateReason:                f.StateReason,
		StateReasonCode:            f.StateReasonCode,
		LastUpdateStatus:           f.LastUpdateStatus,
		LastUpdateReason:           f.LastUpdateReason,
		LastUpdateStatusReason:     f.LastUpdateStatusReason,
		LastUpdateStatusReasonCode: f.LastUpdateStatusReasonCode,
		LastModified:               f.LastModified,
		LastModifiedUser:           f.LastModifiedUser,
		PackageType:                f.PackageType,
		SigningProfileVersionArn:   f.SigningProfileVersionArn,
		SigningJobArn:              f.SigningJobArn,
		CodeSigningConfigArn:       f.CodeSigningConfigArn,
		ContainerID:                f.ContainerID,
		ContainerImageID:           f.ContainerImageID,
		PolicyRevisionId:           f.PolicyRevisionId,
	}

	result.EphemeralStorage = deepCopyEphemeralStorage(f.EphemeralStorage)
	result.Architectures = deepCopyStrings(f.Architectures)
	result.VpcConfig = deepCopyVpcConfig(f.VpcConfig)
	result.Environment = deepCopyEnvironment(f.Environment)
	result.DeadLetterConfig = deepCopyDeadLetterConfig(f.DeadLetterConfig)
	result.TracingConfig = deepCopyTracingConfig(f.TracingConfig)
	result.Layers = deepCopyLayers(f.Layers)
	result.Tags = deepCopyTags(f.Tags)
	result.SnapStart = deepCopySnapStart(f.SnapStart)
	result.UrlConfig = deepCopyFunctionUrlConfig(f.UrlConfig)
	result.LoggingConfig = deepCopyLoggingConfig(f.LoggingConfig)
	result.ImageConfig = deepCopyImageConfig(f.ImageConfig)
	result.FileSystemConfigs = deepCopyFileSystemConfigs(f.FileSystemConfigs)
	result.TenancyConfig = deepCopyTenancyConfig(f.TenancyConfig)
	result.CapacityProviderConfig = deepCopyStringAnyMap(f.CapacityProviderConfig)
	result.RuntimeVersionConfig = deepCopyRuntimeVersionConfig(f.RuntimeVersionConfig)
	result.DurableConfig = deepCopyStringAnyMap(f.DurableConfig)
	result.ReservedConcurrency = deepCopyInt64Ptr(f.ReservedConcurrency)

	if len(f.Versions) > 0 {
		result.Versions = make([]Version, len(f.Versions))
		for i := range f.Versions {
			result.Versions[i] = *f.Versions[i].DeepCopy()
		}
	}

	if len(f.Aliases) > 0 {
		result.Aliases = make([]Alias, len(f.Aliases))
		for i := range f.Aliases {
			result.Aliases[i] = *f.Aliases[i].DeepCopy()
		}
	}

	if len(f.Policies) > 0 {
		result.Policies = make([]FunctionPolicy, len(f.Policies))
		for i := range f.Policies {
			result.Policies[i] = *f.Policies[i].DeepCopy()
		}
	}

	if len(f.ProvisionedConcurrency) > 0 {
		result.ProvisionedConcurrency = make([]ProvisionedConcurrencyConfig, len(f.ProvisionedConcurrency))
		copy(result.ProvisionedConcurrency, f.ProvisionedConcurrency)
	}

	if len(f.EventInvokeConfigs) > 0 {
		result.EventInvokeConfigs = make([]EventInvokeConfig, len(f.EventInvokeConfigs))
		for i := range f.EventInvokeConfigs {
			result.EventInvokeConfigs[i] = *f.EventInvokeConfigs[i].DeepCopy()
		}
	}

	result.rebuildIndexes()
	return result
}

// DeepCopy returns a fully detached copy of the published version,
// including the members the historic services-side copier silently
// dropped (DeadLetterConfig, TracingConfig, SnapStart, ContainerID).
func (v *Version) DeepCopy() *Version {
	if v == nil {
		return nil
	}

	result := &Version{
		Version:                  v.Version,
		FunctionArn:              v.FunctionArn,
		Runtime:                  v.Runtime,
		Role:                     v.Role,
		Handler:                  v.Handler,
		CodeSize:                 v.CodeSize,
		CodeSha256:               v.CodeSha256,
		CodeLocation:             v.CodeLocation,
		ImageUri:                 v.ImageUri,
		Description:              v.Description,
		Timeout:                  v.Timeout,
		MemorySize:               v.MemorySize,
		KMSKeyArn:                v.KMSKeyArn,
		RevisionId:               v.RevisionId,
		State:                    v.State,
		StateReason:              v.StateReason,
		StateReasonCode:          v.StateReasonCode,
		LastUpdateStatus:         v.LastUpdateStatus,
		LastModified:             v.LastModified,
		PackageType:              v.PackageType,
		SigningProfileVersionArn: v.SigningProfileVersionArn,
		SigningJobArn:            v.SigningJobArn,
		ContainerID:              v.ContainerID,
		ContainerImageID:         v.ContainerImageID,
	}

	result.EphemeralStorage = deepCopyEphemeralStorage(v.EphemeralStorage)
	result.Architectures = deepCopyStrings(v.Architectures)
	result.VpcConfig = deepCopyVpcConfig(v.VpcConfig)
	result.Environment = deepCopyEnvironment(v.Environment)
	result.DeadLetterConfig = deepCopyDeadLetterConfig(v.DeadLetterConfig)
	result.TracingConfig = deepCopyTracingConfig(v.TracingConfig)
	result.Layers = deepCopyLayers(v.Layers)
	result.SnapStart = deepCopySnapStart(v.SnapStart)
	result.LoggingConfig = deepCopyLoggingConfig(v.LoggingConfig)
	result.ImageConfig = deepCopyImageConfig(v.ImageConfig)
	result.FileSystemConfigs = deepCopyFileSystemConfigs(v.FileSystemConfigs)

	return result
}

// DeepCopy returns a fully detached copy of the alias, including a copied
// routing-config weight map.
func (a *Alias) DeepCopy() *Alias {
	if a == nil {
		return nil
	}
	result := *a
	if a.RoutingConfig != nil {
		result.RoutingConfig = &RoutingConfig{
			AdditionalVersionWeights: make(map[string]float64, len(a.RoutingConfig.AdditionalVersionWeights)),
		}
		for k, w := range a.RoutingConfig.AdditionalVersionWeights {
			result.RoutingConfig.AdditionalVersionWeights[k] = w
		}
	}
	return &result
}

// DeepCopy returns a fully detached copy of the resource-based policy,
// including a copied condition map.
func (p *FunctionPolicy) DeepCopy() *FunctionPolicy {
	if p == nil {
		return nil
	}
	result := *p
	result.Condition = deepCopyStringAnyMap(p.Condition)
	return &result
}

// DeepCopy returns a fully detached copy of the event invoke config,
// reproducing the destination config at one level (its members are plain
// strings).
func (c *EventInvokeConfig) DeepCopy() *EventInvokeConfig {
	if c == nil {
		return nil
	}
	result := *c
	if c.DestinationConfig != nil {
		result.DestinationConfig = &DestinationConfig{}
		if c.DestinationConfig.OnSuccess != nil {
			result.DestinationConfig.OnSuccess = &OnSuccess{Destination: c.DestinationConfig.OnSuccess.Destination}
		}
		if c.DestinationConfig.OnFailure != nil {
			result.DestinationConfig.OnFailure = &OnFailure{Destination: c.DestinationConfig.OnFailure.Destination}
		}
	}
	return &result
}

func deepCopyFunctionUrlConfig(src *FunctionUrlConfig) *FunctionUrlConfig {
	if src == nil {
		return nil
	}
	result := &FunctionUrlConfig{
		FunctionUrl:      src.FunctionUrl,
		FunctionArn:      src.FunctionArn,
		AuthType:         src.AuthType,
		InvokeMode:       src.InvokeMode,
		CreationTime:     src.CreationTime,
		LastModifiedTime: src.LastModifiedTime,
		Qualifier:        src.Qualifier,
	}
	if src.Cors != nil {
		result.Cors = &CorsConfig{
			AllowCredentials: src.Cors.AllowCredentials,
			MaxAge:           src.Cors.MaxAge,
		}
		result.Cors.AllowHeaders = deepCopyStrings(src.Cors.AllowHeaders)
		result.Cors.AllowMethods = deepCopyStrings(src.Cors.AllowMethods)
		result.Cors.AllowOrigins = deepCopyStrings(src.Cors.AllowOrigins)
		result.Cors.ExposeHeaders = deepCopyStrings(src.Cors.ExposeHeaders)
	}
	return result
}

func deepCopyTags(src []types.Tag) []types.Tag {
	if src == nil {
		return nil
	}
	dst := make([]types.Tag, len(src))
	copy(dst, src)
	return dst
}

func deepCopyTenancyConfig(src *TenancyConfig) *TenancyConfig {
	if src == nil {
		return nil
	}
	dst := *src
	return &dst
}

func deepCopyRuntimeVersionConfig(src *RuntimeVersionConfig) *RuntimeVersionConfig {
	if src == nil {
		return nil
	}
	dst := *src
	dst.Error = deepCopyStringAnyMap(src.Error)
	return &dst
}

// deepCopyStringAnyMap copies the map itself; the interface{} values are
// shared and treated as immutable configuration documents.
func deepCopyStringAnyMap(src map[string]interface{}) map[string]interface{} {
	if src == nil {
		return nil
	}
	dst := make(map[string]interface{}, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func deepCopyInt64Ptr(src *int64) *int64 {
	if src == nil {
		return nil
	}
	dst := *src
	return &dst
}

func deepCopyVpcConfig(src *VpcConfig) *VpcConfig {
	if src == nil {
		return nil
	}
	dst := *src
	if src.SecurityGroupIds != nil {
		dst.SecurityGroupIds = make([]string, len(src.SecurityGroupIds))
		copy(dst.SecurityGroupIds, src.SecurityGroupIds)
	}
	if src.SubnetIds != nil {
		dst.SubnetIds = make([]string, len(src.SubnetIds))
		copy(dst.SubnetIds, src.SubnetIds)
	}
	return &dst
}

func deepCopyEnvironment(src *Environment) *Environment {
	if src == nil {
		return nil
	}
	dst := *src
	if src.Variables != nil {
		dst.Variables = make(map[string]string, len(src.Variables))
		for k, v := range src.Variables {
			dst.Variables[k] = v
		}
	}
	if src.Error != nil {
		dst.Error = &EnvironmentError{ErrorCode: src.Error.ErrorCode, Message: src.Error.Message}
	}
	return &dst
}

func deepCopyDeadLetterConfig(src *DeadLetterConfig) *DeadLetterConfig {
	if src == nil {
		return nil
	}
	dst := *src
	return &dst
}

func deepCopyTracingConfig(src *TracingConfig) *TracingConfig {
	if src == nil {
		return nil
	}
	dst := *src
	return &dst
}

func deepCopySnapStart(src *SnapStart) *SnapStart {
	if src == nil {
		return nil
	}
	dst := *src
	return &dst
}

func deepCopyEphemeralStorage(src *EphemeralStorage) *EphemeralStorage {
	if src == nil {
		return nil
	}
	dst := *src
	return &dst
}

func deepCopyLayers(src []LayerReference) []LayerReference {
	if src == nil {
		return nil
	}
	dst := make([]LayerReference, len(src))
	copy(dst, src)
	return dst
}

func deepCopyStrings(src []string) []string {
	if src == nil {
		return nil
	}
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}

func deepCopyLoggingConfig(src *LoggingConfig) *LoggingConfig {
	if src == nil {
		return nil
	}
	dst := *src
	return &dst
}

func deepCopyImageConfig(src *ImageConfig) *ImageConfig {
	if src == nil {
		return nil
	}
	dst := *src
	if src.EntryPoint != nil {
		dst.EntryPoint = make([]string, len(src.EntryPoint))
		copy(dst.EntryPoint, src.EntryPoint)
	}
	if src.Command != nil {
		dst.Command = make([]string, len(src.Command))
		copy(dst.Command, src.Command)
	}
	return &dst
}

func deepCopyFileSystemConfigs(src []FileSystemConfig) []FileSystemConfig {
	if src == nil {
		return nil
	}
	dst := make([]FileSystemConfig, len(src))
	copy(dst, src)
	for i := range dst {
		if dst[i].S3FilesConfig != nil {
			cfg := *dst[i].S3FilesConfig
			dst[i].S3FilesConfig = &cfg
		}
	}
	return dst
}
