// Package lambda provides AWS Lambda store functionality for vorpalstacks.
package lambda

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"

	"github.com/google/uuid"
)

// FunctionStore provides Lambda function storage functionality.
type FunctionStore struct {
	*common.BaseStore
	TagStore   *common.TagStore
	arnBuilder *ARNBuilder
	accountId  string
	region     string
	mu         sync.Mutex
}

// NewFunctionStore creates a new Lambda function store.
func NewFunctionStore(store storage.BasicStorage, accountId, region string) *FunctionStore {
	bucket := store.Bucket("lambda-functions-" + region)
	return &FunctionStore{
		BaseStore:  common.NewBaseStore(bucket, "lambda-functions"),
		TagStore:   common.NewTagStoreWithRegion(store, "lambda", region),
		arnBuilder: NewARNBuilder(accountId, region),
		accountId:  accountId,
		region:     region,
	}
}

// Create creates a new Lambda function in the store.
// Returns the created function or an error if it already exists.
func (s *FunctionStore) Create(function *Function) (*Function, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Exists(function.FunctionName) {
		return nil, ErrFunctionAlreadyExists
	}

	now := time.Now().UTC()
	function.FunctionArn = s.arnBuilder.FunctionArn(function.FunctionName)
	function.LastModified = now
	function.State = StateActive
	function.LastUpdateStatus = LastUpdateStatusSuccessful
	function.RevisionId = uuid.New().String()

	if function.Timeout == 0 {
		function.Timeout = 3
	}
	if function.MemorySize == 0 {
		function.MemorySize = 128
	}
	if function.EphemeralStorage == nil {
		function.EphemeralStorage = &EphemeralStorage{Size: 512}
	}
	if function.PackageType == "" {
		function.PackageType = "Zip"
	}
	if function.CurrentVersion == "" {
		function.CurrentVersion = "$LATEST"
	}

	if function.CodeSize > 0 && function.CodeSha256 == "" {
		function.CodeSha256 = GenerateCodeHash([]byte(fmt.Sprintf("%s-%d-%d", function.FunctionName, time.Now().UnixNano(), function.CodeSize)))
	}

	if err := s.Put(function.FunctionName, function); err != nil {
		return nil, err
	}

	return function, nil
}

// Get retrieves a Lambda function by its name from the store.
// Returns the function or an error if not found.
func (s *FunctionStore) Get(functionName string) (*Function, error) {
	var function Function
	if err := s.BaseStore.Get(functionName, &function); err != nil {
		return nil, ErrFunctionNotFound
	}
	function.rebuildIndexes()
	return &function, nil
}

// rebuildIndexes reconstructs the lookup indexes from the slices.
func (f *Function) rebuildIndexes() {
	if len(f.Versions) > 0 && f.versionsByNum == nil {
		f.versionsByNum = make(map[string]*Version, len(f.Versions))
		maxVersion := 0
		for i := range f.Versions {
			f.versionsByNum[f.Versions[i].Version] = &f.Versions[i]
			if vInt, err := strconv.Atoi(f.Versions[i].Version); err == nil && vInt > maxVersion {
				maxVersion = vInt
			}
		}
		f.latestVersionNum = maxVersion
	}
	if len(f.Aliases) > 0 && f.aliasesByName == nil {
		f.aliasesByName = make(map[string]*Alias, len(f.Aliases))
		for i := range f.Aliases {
			f.aliasesByName[f.Aliases[i].Name] = &f.Aliases[i]
		}
	}
}

// GetByArn retrieves a Lambda function by its ARN from the store.
// Returns the function or an error if not found.
func (s *FunctionStore) GetByArn(arn string) (*Function, error) {
	functionName := s.arnBuilder.ParseFunctionNameFromArn(arn)
	if functionName == "" {
		return nil, ErrFunctionNotFound
	}
	return s.Get(functionName)
}

// Update updates an existing Lambda function in the store.
// Returns an error if the function does not exist.
func (s *FunctionStore) Update(function *Function) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.updateInternal(function)
}

// UpdateAtomically performs an atomic read-modify-write operation on a function.
// The modifier function is called while holding the lock, preventing race conditions.
// Returns the updated function or an error.
func (s *FunctionStore) UpdateAtomically(functionName string, modifier func(*Function) error) (*Function, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	function, err := s.getInternal(functionName)
	if err != nil {
		return nil, err
	}

	if err := modifier(function); err != nil {
		return nil, err
	}

	if err := s.updateInternal(function); err != nil {
		return nil, err
	}

	return function, nil
}

// getInternal retrieves a function without mutex (caller must hold lock).
func (s *FunctionStore) getInternal(functionName string) (*Function, error) {
	var function Function
	if err := s.BaseStore.Get(functionName, &function); err != nil {
		return nil, ErrFunctionNotFound
	}
	function.rebuildIndexes()
	return &function, nil
}

// updateInternal updates without mutex (caller must hold lock).
func (s *FunctionStore) updateInternal(function *Function) error {
	if !s.Exists(function.FunctionName) {
		return ErrFunctionNotFound
	}

	function.LastModified = time.Now().UTC()
	function.LastUpdateStatus = LastUpdateStatusSuccessful
	function.RevisionId = uuid.New().String()

	return s.Put(function.FunctionName, function)
}

// Delete deletes a Lambda function by its name from the store.
// Returns an error if the function does not exist.
func (s *FunctionStore) Delete(functionName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.Exists(functionName) {
		return ErrFunctionNotFound
	}
	return s.BaseStore.Delete(functionName)
}

// List returns a list of Lambda functions from the store with the specified list options.
func (s *FunctionStore) List(opts common.ListOptions) (*common.ListResult[Function], error) {
	return common.List[Function](s.BaseStore, opts, nil)
}

// ListByPrefix returns a list of Lambda functions filtered by name prefix.
func (s *FunctionStore) ListByPrefix(prefix string) ([]*Function, error) {
	var functions []*Function
	err := s.ScanPrefix("", func(key string, value []byte) error {
		var fn Function
		if err := json.Unmarshal(value, &fn); err != nil {
			return err
		}
		if prefix == "" || (len(fn.FunctionName) >= len(prefix) && fn.FunctionName[:len(prefix)] == prefix) {
			functions = append(functions, &fn)
		}
		return nil
	})
	return functions, err
}

// PublishVersion publishes a new version of a Lambda function.
// Returns the published version or an error if the function does not exist.
func (s *FunctionStore) PublishVersion(function *Function, description string) (*Version, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	function.rebuildIndexes()

	versionNum := "1"
	if function.latestVersionNum > 0 || len(function.Versions) > 0 {
		versionNum = strconv.Itoa(function.latestVersionNum + 1)
	}

	version := &Version{
		Version:          versionNum,
		FunctionArn:      s.arnBuilder.FunctionVersionArn(function.FunctionName, versionNum),
		Role:             function.Role,
		Runtime:          function.Runtime,
		Handler:          function.Handler,
		CodeSize:         function.CodeSize,
		CodeSha256:       function.CodeSha256,
		CodeLocation:     function.CodeLocation,
		ImageUri:         function.ImageUri,
		Description:      description,
		Timeout:          function.Timeout,
		MemorySize:       function.MemorySize,
		EphemeralStorage: deepCopyEphemeralStorage(function.EphemeralStorage),
		Architectures:    deepCopyArchitectures(function.Architectures),
		KMSKeyArn:        function.KMSKeyArn,
		RevisionId:       uuid.New().String(),
		State:            StateActive,
		LastUpdateStatus: LastUpdateStatusSuccessful,
		LastModified:     time.Now().UTC(),
		VpcConfig:        deepCopyVpcConfig(function.VpcConfig),
		Environment:      deepCopyEnvironment(function.Environment),
		DeadLetterConfig: deepCopyDeadLetterConfig(function.DeadLetterConfig),
		TracingConfig:    deepCopyTracingConfig(function.TracingConfig),
		Layers:           deepCopyLayers(function.Layers),
		SnapStart:        deepCopySnapStart(function.SnapStart),
		PackageType:      function.PackageType,
		ContainerID:      function.ContainerID,
		ContainerImageID: function.ContainerImageID,
	}

	function.Versions = append(function.Versions, *version)
	function.CurrentVersion = versionNum

	if function.versionsByNum == nil {
		function.versionsByNum = make(map[string]*Version)
	}
	function.versionsByNum[versionNum] = &function.Versions[len(function.Versions)-1]
	if vInt, err := strconv.Atoi(versionNum); err == nil && vInt > function.latestVersionNum {
		function.latestVersionNum = vInt
	}

	if err := s.updateInternal(function); err != nil {
		return nil, err
	}

	return version, nil
}

// GetVersion retrieves a specific version of a Lambda function.
// Returns the version or an error if not found.
func (s *FunctionStore) GetVersion(functionName, version string) (*Version, error) {
	function, err := s.Get(functionName)
	if err != nil {
		return nil, err
	}

	if v, ok := function.versionsByNum[version]; ok {
		return v, nil
	}

	return nil, ErrVersionNotFound
}

// DeleteVersion deletes a specific version of a Lambda function.
// Returns an error if the version does not exist.
func (s *FunctionStore) DeleteVersion(functionName, version string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var function Function
	if err := s.BaseStore.Get(functionName, &function); err != nil {
		return ErrFunctionNotFound
	}

	function.rebuildIndexes()

	for _, alias := range function.Aliases {
		if alias.FunctionVersion == version {
			return ErrResourceConflict
		}
	}

	for i, v := range function.Versions {
		if v.Version == version {
			function.Versions = append(function.Versions[:i], function.Versions[i+1:]...)
			return s.updateInternal(&function)
		}
	}

	return ErrVersionNotFound
}

// CreateAlias creates a new alias for a Lambda function.
// Returns the created alias or an error if it already exists.
func (s *FunctionStore) CreateAlias(function *Function, alias *Alias) (*Alias, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	function.rebuildIndexes()

	if _, exists := function.aliasesByName[alias.Name]; exists {
		return nil, ErrAliasAlreadyExists
	}

	alias.AliasArn = s.arnBuilder.FunctionAliasArn(function.FunctionName, alias.Name)
	alias.FunctionName = function.FunctionName
	alias.RevisionId = uuid.New().String()

	if alias.FunctionVersion == "" {
		alias.FunctionVersion = "$LATEST"
	}

	function.Aliases = append(function.Aliases, *alias)

	if function.aliasesByName == nil {
		function.aliasesByName = make(map[string]*Alias)
	}
	function.aliasesByName[alias.Name] = &function.Aliases[len(function.Aliases)-1]

	if err := s.updateInternal(function); err != nil {
		return nil, err
	}

	return alias, nil
}

// GetAlias retrieves an alias for a Lambda function by name.
// Returns the alias or an error if not found.
func (s *FunctionStore) GetAlias(functionName, aliasName string) (*Alias, error) {
	function, err := s.Get(functionName)
	if err != nil {
		return nil, err
	}

	if a, ok := function.aliasesByName[aliasName]; ok {
		return a, nil
	}

	return nil, ErrAliasNotFound
}

// UpdateAlias updates an existing alias for a Lambda function.
// Returns an error if the alias does not exist.
// Deprecated: Use UpdateAliasAtomically for race-free updates.
func (s *FunctionStore) UpdateAlias(function *Function, alias *Alias) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	function.rebuildIndexes()

	existing, ok := function.aliasesByName[alias.Name]
	if !ok {
		return ErrAliasNotFound
	}

	alias.AliasArn = existing.AliasArn
	alias.FunctionName = existing.FunctionName
	alias.RevisionId = uuid.New().String()

	for i := range function.Aliases {
		if function.Aliases[i].Name == alias.Name {
			function.Aliases[i] = *alias
			function.aliasesByName[alias.Name] = &function.Aliases[i]
			return s.updateInternal(function)
		}
	}

	return ErrAliasNotFound
}

// CreateAliasAtomically creates an alias atomically (race-free).
// The creator function is called while holding the lock.
func (s *FunctionStore) CreateAliasAtomically(functionName string, creator func(*Function) (*Alias, error)) (*Alias, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	function, err := s.getInternal(functionName)
	if err != nil {
		return nil, err
	}

	alias, err := creator(function)
	if err != nil {
		return nil, err
	}

	function.rebuildIndexes()

	if _, exists := function.aliasesByName[alias.Name]; exists {
		return nil, ErrAliasAlreadyExists
	}

	alias.AliasArn = s.arnBuilder.FunctionAliasArn(function.FunctionName, alias.Name)
	alias.FunctionName = function.FunctionName
	alias.RevisionId = uuid.New().String()

	if alias.FunctionVersion == "" {
		alias.FunctionVersion = "$LATEST"
	}

	function.Aliases = append(function.Aliases, *alias)

	if function.aliasesByName == nil {
		function.aliasesByName = make(map[string]*Alias)
	}
	function.aliasesByName[alias.Name] = &function.Aliases[len(function.Aliases)-1]

	if err := s.updateInternal(function); err != nil {
		return nil, err
	}

	return alias, nil
}

// UpdateAliasAtomically updates an alias atomically (race-free).
// The modifier function is called while holding the lock.
func (s *FunctionStore) UpdateAliasAtomically(functionName, aliasName string, modifier func(*Function, *Alias) error) (*Alias, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	function, err := s.getInternal(functionName)
	if err != nil {
		return nil, err
	}

	function.rebuildIndexes()

	existing, ok := function.aliasesByName[aliasName]
	if !ok {
		return nil, ErrAliasNotFound
	}

	if err := modifier(function, existing); err != nil {
		return nil, err
	}

	existing.RevisionId = uuid.New().String()

	for i := range function.Aliases {
		if function.Aliases[i].Name == aliasName {
			function.Aliases[i] = *existing
			function.aliasesByName[aliasName] = &function.Aliases[i]
			if err := s.updateInternal(function); err != nil {
				return nil, err
			}
			return existing, nil
		}
	}

	return nil, ErrAliasNotFound
}

// DeleteAlias deletes an alias from a Lambda function.
// Returns an error if the alias does not exist.
func (s *FunctionStore) DeleteAlias(functionName, aliasName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var function Function
	if err := s.BaseStore.Get(functionName, &function); err != nil {
		return ErrFunctionNotFound
	}

	function.rebuildIndexes()

	if _, ok := function.aliasesByName[aliasName]; !ok {
		return ErrAliasNotFound
	}

	for i := range function.Aliases {
		if function.Aliases[i].Name == aliasName {
			function.Aliases = append(function.Aliases[:i], function.Aliases[i+1:]...)
			delete(function.aliasesByName, aliasName)
			return s.updateInternal(&function)
		}
	}

	return ErrAliasNotFound
}

// AddPolicy adds a resource-based policy to a Lambda function.
func (s *FunctionStore) AddPolicy(function *Function, policy *FunctionPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if policy.Id == "" {
		policy.Id = uuid.New().String()
	}
	function.Policies = append(function.Policies, *policy)
	return s.updateInternal(function)
}

// AddPolicyAtomically adds a resource-based policy to a Lambda function atomically.
func (s *FunctionStore) AddPolicyAtomically(functionName string, policy *FunctionPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	function, err := s.getInternal(functionName)
	if err != nil {
		return err
	}

	if policy.Id == "" {
		policy.Id = uuid.New().String()
	}
	function.Policies = append(function.Policies, *policy)
	return s.updateInternal(function)
}

// RemovePolicy removes a resource-based policy from a Lambda function.
func (s *FunctionStore) RemovePolicy(functionName, statementId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var function Function
	if err := s.BaseStore.Get(functionName, &function); err != nil {
		return ErrFunctionNotFound
	}

	for i, p := range function.Policies {
		if p.Id == statementId {
			function.Policies = append(function.Policies[:i], function.Policies[i+1:]...)
			return s.updateInternal(&function)
		}
	}

	return ErrPolicyNotFound
}

// GetPolicy retrieves the resource-based policy for a Lambda function.
func (s *FunctionStore) GetPolicy(functionName string) ([]FunctionPolicy, error) {
	function, err := s.Get(functionName)
	if err != nil {
		return nil, err
	}
	return function.Policies, nil
}

// SetReservedConcurrency sets the reserved concurrency for a Lambda function.
func (s *FunctionStore) SetReservedConcurrency(functionName string, concurrency *int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var function Function
	if err := s.BaseStore.Get(functionName, &function); err != nil {
		return ErrFunctionNotFound
	}
	function.ReservedConcurrency = concurrency
	return s.updateInternal(&function)
}

// GetReservedConcurrency retrieves the reserved concurrency for a Lambda function.
func (s *FunctionStore) GetReservedConcurrency(functionName string) (*int64, error) {
	function, err := s.Get(functionName)
	if err != nil {
		return nil, err
	}
	return function.ReservedConcurrency, nil
}

// SetContainerInfo sets the container information for a Lambda function.
func (s *FunctionStore) SetContainerInfo(functionName, containerID, containerImageID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var function Function
	if err := s.BaseStore.Get(functionName, &function); err != nil {
		return ErrFunctionNotFound
	}
	function.ContainerID = containerID
	function.ContainerImageID = containerImageID
	return s.updateInternal(&function)
}

// SetFunctionUrlConfig sets the function URL configuration for a Lambda function.
func (s *FunctionStore) SetFunctionUrlConfig(functionName string, config *FunctionUrlConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var function Function
	if err := s.BaseStore.Get(functionName, &function); err != nil {
		return ErrFunctionNotFound
	}
	now := time.Now().UTC()
	if function.UrlConfig == nil {
		config.CreationTime = now
	}
	config.LastModifiedTime = now
	config.FunctionArn = function.FunctionArn
	if config.FunctionUrl == "" {
		urlId := uuid.New().String()[:8]
		config.FunctionUrl = fmt.Sprintf("https://%s.lambda-url.%s.on.aws", urlId, s.region)
	}
	function.UrlConfig = config
	return s.updateInternal(&function)
}

// GetFunctionUrlConfig retrieves the function URL configuration for a Lambda function.
func (s *FunctionStore) GetFunctionUrlConfig(functionName string) (*FunctionUrlConfig, error) {
	function, err := s.Get(functionName)
	if err != nil {
		return nil, err
	}
	if function.UrlConfig == nil {
		return nil, ErrUrlConfigNotFound
	}
	return function.UrlConfig, nil
}

// DeleteFunctionUrlConfig deletes the function URL configuration for a Lambda function.
func (s *FunctionStore) DeleteFunctionUrlConfig(functionName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var function Function
	if err := s.BaseStore.Get(functionName, &function); err != nil {
		return ErrFunctionNotFound
	}
	function.UrlConfig = nil
	return s.updateInternal(&function)
}

// ResolveQualifier resolves a Lambda function qualifier (version or alias) to its function, version, or alias.
// Returns the function, version, and alias based on the qualifier.
func (s *FunctionStore) ResolveQualifier(functionName, qualifier string) (*Function, *Version, *Alias, error) {
	function, err := s.Get(functionName)
	if err != nil {
		return nil, nil, nil, err
	}

	if qualifier == "" || qualifier == "$LATEST" {
		return function, nil, nil, nil
	}

	if v, ok := function.versionsByNum[qualifier]; ok {
		return function, v, nil, nil
	}

	if a, ok := function.aliasesByName[qualifier]; ok {
		return function, nil, a, nil
	}

	return nil, nil, nil, ErrVersionNotFound
}

// SetProvisionedConcurrency sets the provisioned concurrency for a Lambda function qualifier.
func (s *FunctionStore) SetProvisionedConcurrency(functionName, qualifier string, concurrentExecutions int32) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var function Function
	if err := s.BaseStore.Get(functionName, &function); err != nil {
		return ErrFunctionNotFound
	}

	now := time.Now().UTC()
	found := false
	for i, pc := range function.ProvisionedConcurrency {
		if pc.Qualifier == qualifier {
			function.ProvisionedConcurrency[i].AllocatedProvisionedConcurrentExecutions = concurrentExecutions
			function.ProvisionedConcurrency[i].AvailableProvisionedConcurrentExecutions = concurrentExecutions
			function.ProvisionedConcurrency[i].Status = "READY"
			function.ProvisionedConcurrency[i].LastModified = now
			found = true
			break
		}
	}

	if !found {
		function.ProvisionedConcurrency = append(function.ProvisionedConcurrency, ProvisionedConcurrencyConfig{
			FunctionName:                             functionName,
			Qualifier:                                qualifier,
			AllocatedProvisionedConcurrentExecutions: concurrentExecutions,
			AvailableProvisionedConcurrentExecutions: concurrentExecutions,
			Status:                                   "READY",
			LastModified:                             now,
		})
	}

	return s.updateInternal(&function)
}

// GetProvisionedConcurrency retrieves the provisioned concurrency configuration for a Lambda function qualifier.
func (s *FunctionStore) GetProvisionedConcurrency(functionName, qualifier string) (*ProvisionedConcurrencyConfig, error) {
	function, err := s.Get(functionName)
	if err != nil {
		return nil, err
	}

	for _, pc := range function.ProvisionedConcurrency {
		if pc.Qualifier == qualifier {
			return &pc, nil
		}
	}

	return nil, ErrProvisionedConcurrencyNotFound
}

// DeleteProvisionedConcurrency deletes the provisioned concurrency configuration for a Lambda function qualifier.
func (s *FunctionStore) DeleteProvisionedConcurrency(functionName, qualifier string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var function Function
	if err := s.BaseStore.Get(functionName, &function); err != nil {
		return ErrFunctionNotFound
	}

	for i, pc := range function.ProvisionedConcurrency {
		if pc.Qualifier == qualifier {
			function.ProvisionedConcurrency = append(
				function.ProvisionedConcurrency[:i],
				function.ProvisionedConcurrency[i+1:]...,
			)
			return s.updateInternal(&function)
		}
	}

	return ErrProvisionedConcurrencyNotFound
}

// ListProvisionedConcurrency lists all provisioned concurrency configurations for a Lambda function.
func (s *FunctionStore) ListProvisionedConcurrency(functionName string) ([]ProvisionedConcurrencyConfig, error) {
	function, err := s.Get(functionName)
	if err != nil {
		return nil, err
	}

	return function.ProvisionedConcurrency, nil
}

// SetEventInvokeConfig sets the event invoke configuration for a Lambda function qualifier.
func (s *FunctionStore) SetEventInvokeConfig(functionName, qualifier string, config *EventInvokeConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var function Function
	if err := s.BaseStore.Get(functionName, &function); err != nil {
		return ErrFunctionNotFound
	}

	now := time.Now().UTC()
	config.FunctionName = functionName
	config.Qualifier = qualifier
	config.LastModified = now

	found := false
	for i, eic := range function.EventInvokeConfigs {
		if eic.Qualifier == qualifier {
			function.EventInvokeConfigs[i] = *config
			found = true
			break
		}
	}

	if !found {
		function.EventInvokeConfigs = append(function.EventInvokeConfigs, *config)
	}

	return s.updateInternal(&function)
}

// GetEventInvokeConfig retrieves the event invoke configuration for a Lambda function qualifier.
func (s *FunctionStore) GetEventInvokeConfig(functionName, qualifier string) (*EventInvokeConfig, error) {
	function, err := s.Get(functionName)
	if err != nil {
		return nil, err
	}

	for _, eic := range function.EventInvokeConfigs {
		if eic.Qualifier == qualifier {
			return &eic, nil
		}
	}

	return nil, ErrEventInvokeConfigNotFound
}

// DeleteEventInvokeConfig deletes the event invoke configuration for a Lambda function qualifier.
func (s *FunctionStore) DeleteEventInvokeConfig(functionName, qualifier string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var function Function
	if err := s.BaseStore.Get(functionName, &function); err != nil {
		return ErrFunctionNotFound
	}

	for i, eic := range function.EventInvokeConfigs {
		if eic.Qualifier == qualifier {
			function.EventInvokeConfigs = append(
				function.EventInvokeConfigs[:i],
				function.EventInvokeConfigs[i+1:]...,
			)
			return s.updateInternal(&function)
		}
	}

	return ErrEventInvokeConfigNotFound
}

// ListEventInvokeConfigs lists all event invoke configurations for a Lambda function.
func (s *FunctionStore) ListEventInvokeConfigs(functionName string) ([]EventInvokeConfig, error) {
	function, err := s.Get(functionName)
	if err != nil {
		return nil, err
	}

	return function.EventInvokeConfigs, nil
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

func deepCopyArchitectures(src []string) []string {
	if src == nil {
		return nil
	}
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}
