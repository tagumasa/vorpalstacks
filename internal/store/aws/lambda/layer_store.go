// Package lambda provides AWS Lambda store functionality for vorpalstacks.
package lambda

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"

	"github.com/google/uuid"
)

// LayerStore provides storage operations for AWS Lambda layers.
type LayerStore struct {
	*common.BaseStore
	arnBuilder *ArnBuilder
	accountId  string
	region     string
	mu         sync.Mutex
}

// NewLayerStore creates a new LayerStore instance.
func NewLayerStore(store storage.BasicStorage, accountId, region string) *LayerStore {
	bucket := store.Bucket("lambda-layers-" + region)
	return &LayerStore{
		BaseStore:  common.NewBaseStore(bucket, "lambda-layers"),
		arnBuilder: NewArnBuilder(accountId, region),
		accountId:  accountId,
		region:     region,
	}
}

// Create stores a new Lambda layer.
func (s *LayerStore) Create(layer *Layer) (*Layer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Exists(layer.LayerName) {
		return nil, ErrResourceConflict
	}

	layer.LayerArn = s.arnBuilder.LayerArn(layer.LayerName)
	layer.CreatedDate = time.Now().UTC()

	if err := s.Put(layer.LayerName, layer); err != nil {
		return nil, err
	}

	return layer, nil
}

// Get retrieves a Lambda layer by name.
func (s *LayerStore) Get(layerName string) (*Layer, error) {
	var layer Layer
	if err := s.BaseStore.Get(layerName, &layer); err != nil {
		return nil, ErrLayerNotFound
	}
	return &layer, nil
}

// Delete removes a Lambda layer by name.
func (s *LayerStore) Delete(layerName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.Exists(layerName) {
		return ErrLayerNotFound
	}
	return s.BaseStore.Delete(layerName)
}

// List returns a list of Lambda layers with pagination options.
func (s *LayerStore) List(opts common.ListOptions) (*common.ListResult[Layer], error) {
	return common.List[Layer](s.BaseStore, opts, nil)
}

// PublishVersion creates a new version of a Lambda layer.
func (s *LayerStore) PublishVersion(layer *Layer, version *LayerVersion) (*LayerVersion, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	versionNum := int64(0)
	for _, v := range layer.Versions {
		if v.Version > versionNum {
			versionNum = v.Version
		}
	}
	versionNum++

	versionCopy := deepCopyLayerVersion(version)
	versionCopy.Version = versionNum
	versionCopy.LayerVersionArn = s.arnBuilder.LayerVersionArn(layer.LayerName, versionNum)
	versionCopy.CreatedDate = time.Now().UTC()

	if versionCopy.CodeSize > 0 && versionCopy.CodeSha256 == "" {
		versionCopy.CodeSha256 = GenerateCodeHash([]byte(fmt.Sprintf("%s-%d-%d", layer.LayerName, versionCopy.Version, versionCopy.CodeSize)))
	}

	layer.Versions = append(layer.Versions, versionCopy)
	layer.LatestMatchingVersion = versionCopy

	if err := s.updateInternal(layer); err != nil {
		return nil, err
	}

	return versionCopy, nil
}

// GetVersion retrieves a specific version of a Lambda layer.
func (s *LayerStore) GetVersion(layerName string, versionNumber int64) (*LayerVersion, error) {
	layer, err := s.Get(layerName)
	if err != nil {
		return nil, err
	}

	for _, v := range layer.Versions {
		if v.Version == versionNumber {
			return v, nil
		}
	}

	return nil, ErrLayerVersionNotFound
}

// GetVersionByArn retrieves a Lambda layer version by its ARN.
func (s *LayerStore) GetVersionByArn(arn string) (*LayerVersion, error) {
	layerName := s.arnBuilder.ParseLayerNameFromArn(arn)
	if layerName == "" {
		return nil, ErrLayerVersionNotFound
	}

	versionNumber := s.arnBuilder.ParseLayerVersionFromArn(arn)

	layer, err := s.Get(layerName)
	if err != nil {
		return nil, err
	}

	if versionNumber > 0 {
		for _, v := range layer.Versions {
			if v.Version == versionNumber {
				return v, nil
			}
		}
		return nil, ErrLayerVersionNotFound
	}

	if layer.LatestMatchingVersion != nil {
		return layer.LatestMatchingVersion, nil
	}

	return nil, ErrLayerVersionNotFound
}

// DeleteVersion removes a specific version of a Lambda layer.
func (s *LayerStore) DeleteVersion(layerName string, versionNumber int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var layer Layer
	if err := s.BaseStore.Get(layerName, &layer); err != nil {
		return ErrLayerNotFound
	}

	for i, v := range layer.Versions {
		if v.Version == versionNumber {
			layer.Versions = append(layer.Versions[:i], layer.Versions[i+1:]...)
			if layer.LatestMatchingVersion != nil && layer.LatestMatchingVersion.Version == versionNumber {
				if len(layer.Versions) > 0 {
					layer.LatestMatchingVersion = layer.Versions[len(layer.Versions)-1]
				} else {
					layer.LatestMatchingVersion = nil
				}
			}
			return s.updateInternal(&layer)
		}
	}

	return ErrLayerVersionNotFound
}

// ListVersions returns all versions of a Lambda layer with pagination.
func (s *LayerStore) ListVersions(layerName string, opts common.ListOptions) (*common.ListResult[LayerVersion], error) {
	layer, err := s.Get(layerName)
	if err != nil {
		return nil, err
	}

	allVersions := make([]*LayerVersion, len(layer.Versions))
	for i, v := range layer.Versions {
		allVersions[i] = deepCopyLayerVersion(v)
	}

	startIdx := 0
	if opts.Marker != "" {
		for i, v := range allVersions {
			if fmt.Sprintf("%d", v.Version) == opts.Marker {
				startIdx = i + 1
				break
			}
		}
	}

	maxItems := opts.MaxItems
	if maxItems <= 0 {
		maxItems = 50
	}

	endIdx := startIdx + maxItems
	if endIdx > len(allVersions) {
		endIdx = len(allVersions)
	}

	var resultVersions []*LayerVersion
	if startIdx < len(allVersions) {
		resultVersions = allVersions[startIdx:endIdx]
	} else {
		resultVersions = []*LayerVersion{}
	}

	isTruncated := endIdx < len(allVersions)
	nextMarker := ""
	if isTruncated && len(resultVersions) > 0 {
		nextMarker = fmt.Sprintf("%d", resultVersions[len(resultVersions)-1].Version)
	}

	return &common.ListResult[LayerVersion]{
		Items:       resultVersions,
		NextMarker:  nextMarker,
		IsTruncated: isTruncated,
	}, nil
}

// AddPolicy adds a permission policy to a Lambda layer version.
func (s *LayerStore) AddPolicy(layer *Layer, versionNumber int64, policy *LayerPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if layer.LatestMatchingVersion == nil || layer.LatestMatchingVersion.Version != versionNumber {
		return ErrLayerVersionNotFound
	}

	if policy.Id == "" {
		policy.Id = uuid.New().String()
	}
	layer.LatestMatchingVersion.Policies = append(layer.LatestMatchingVersion.Policies, *policy)
	return s.updateInternal(layer)
}

// RemovePolicy removes a permission policy from a Lambda layer version.
func (s *LayerStore) RemovePolicy(layerName string, versionNumber int64, statementId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var layer Layer
	if err := s.BaseStore.Get(layerName, &layer); err != nil {
		return ErrLayerNotFound
	}

	if layer.LatestMatchingVersion == nil || layer.LatestMatchingVersion.Version != versionNumber {
		return ErrLayerVersionNotFound
	}

	for i, p := range layer.LatestMatchingVersion.Policies {
		if p.Id == statementId {
			layer.LatestMatchingVersion.Policies = append(
				layer.LatestMatchingVersion.Policies[:i],
				layer.LatestMatchingVersion.Policies[i+1:]...,
			)
			return s.updateInternal(&layer)
		}
	}

	return ErrPolicyNotFound
}

// Update modifies an existing Lambda layer.
func (s *LayerStore) Update(layer *Layer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.updateInternal(layer)
}

// updateInternal updates without mutex (caller must hold lock).
func (s *LayerStore) updateInternal(layer *Layer) error {
	if !s.Exists(layer.LayerName) {
		return ErrLayerNotFound
	}
	return s.Put(layer.LayerName, layer)
}

// ListByCompatibleRuntime returns Lambda layers compatible with a specific runtime.
func (s *LayerStore) ListByCompatibleRuntime(runtime Runtime) ([]*Layer, error) {
	var layers []*Layer
	err := s.ForEach(func(key string, value []byte) error {
		var layer Layer
		if err := decodeJSON(value, &layer); err != nil {
			return err
		}
		if layer.LatestMatchingVersion != nil {
			for _, r := range layer.LatestMatchingVersion.CompatibleRuntimes {
				if r == runtime {
					layers = append(layers, &layer)
					break
				}
			}
		}
		return nil
	})
	return layers, err
}

func decodeJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func deepCopyLayerVersion(src *LayerVersion) *LayerVersion {
	if src == nil {
		return nil
	}
	dst := *src
	if src.CompatibleRuntimes != nil {
		dst.CompatibleRuntimes = make([]Runtime, len(src.CompatibleRuntimes))
		copy(dst.CompatibleRuntimes, src.CompatibleRuntimes)
	}
	if src.CompatibleArchitectures != nil {
		dst.CompatibleArchitectures = make([]string, len(src.CompatibleArchitectures))
		copy(dst.CompatibleArchitectures, src.CompatibleArchitectures)
	}
	if src.Policies != nil {
		dst.Policies = make([]LayerPolicy, len(src.Policies))
		copy(dst.Policies, src.Policies)
	}
	return &dst
}
