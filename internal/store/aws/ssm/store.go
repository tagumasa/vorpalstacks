// Package ssm provides AWS Systems Manager Parameter Store storage functionality for vorpalstacks.
package ssm

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// Store provides SSM parameter storage functionality.
type Store struct {
	*common.BaseStore
	historyStore *common.BaseStore
	*common.TagStore
	arnBuilder *svcarn.ARNBuilder
	accountID  string
	region     string
	mu         sync.Mutex
}

func parameterBucketName(region string) string {
	return "ssm-parameters-" + region
}

func historyBucketName(region string) string {
	return "ssm-history-" + region
}

// NewStore creates a new SSM store instance.
func NewStore(store storage.BasicStorage, accountID, region string) *Store {
	return &Store{
		BaseStore:    common.NewBaseStore(store.Bucket(parameterBucketName(region)), "ssm-parameters"),
		historyStore: common.NewBaseStore(store.Bucket(historyBucketName(region)), "ssm-history"),
		TagStore:     common.NewTagStoreWithRegion(store, "ssm", region),
		arnBuilder:   svcarn.NewARNBuilder(accountID, region),
		accountID:    accountID,
		region:       region,
	}
}

func (s *Store) paramKey(name string) string {
	return name
}

func (s *Store) historyKey(name string, version int64) string {
	return fmt.Sprintf("%s:v%d", name, version)
}

func (s *Store) buildARN(name string) string {
	if strings.HasPrefix(name, "/") {
		return s.arnBuilder.Build("ssm", "parameter"+name)
	}
	return s.arnBuilder.Build("ssm", "parameter/"+name)
}

// PutParameter stores a parameter in the parameter store.
func (s *Store) PutParameter(param *Parameter, overwrite bool) (int64, error) {
	if err := ValidateParameterName(param.Name); err != nil {
		return 0, err
	}

	key := s.paramKey(param.Name)

	s.mu.Lock()
	defer s.mu.Unlock()

	var existingParam Parameter
	err := s.BaseStore.Get(key, &existingParam)

	if err == nil {
		if !overwrite {
			return 0, ErrParameterAlreadyExists
		}
		param.Version = existingParam.Version + 1
	} else {
		param.Version = 1
	}

	param.LastModifiedDate = time.Now().UTC()
	param.ARN = s.buildARN(param.Name)

	if param.Tier == "" {
		param.Tier = ParameterTierStandard
	}
	if param.DataType == "" {
		param.DataType = "text"
	}
	if param.Tags == nil {
		param.Tags = make(map[string]string)
	}

	if err := s.Put(key, param); err != nil {
		return 0, err
	}

	if len(param.Tags) > 0 {
		if err := s.TagStore.Tag(param.Name, param.Tags); err != nil {
			return 0, err
		}
	}

	historyEntry := &ParameterVersion{
		ParameterName:    param.Name,
		Version:          param.Version,
		Value:            param.Value,
		Type:             param.Type,
		Tier:             param.Tier,
		LastModifiedDate: param.LastModifiedDate,
		DataType:         param.DataType,
		Description:      param.Description,
		KeyID:            param.KeyID,
		AllowedPattern:   param.AllowedPattern,
	}
	hKey := s.historyKey(param.Name, param.Version)
	if err := s.historyStore.Put(hKey, historyEntry); err != nil {
		return 0, fmt.Errorf("failed to store parameter history: %w", err)
	}

	return param.Version, nil
}

// GetParameter retrieves a parameter by name.
func (s *Store) GetParameter(name string, withDecryption bool) (*Parameter, error) {
	key := s.paramKey(name)
	var param Parameter
	if err := s.BaseStore.Get(key, &param); err != nil {
		return nil, ErrParameterNotFound
	}
	return &param, nil
}

// GetParameterByVersion retrieves a specific version of a parameter.
// The returned parameter includes the full ARN reconstructed from the
// parameter name, matching AWS behaviour for versioned lookups.
func (s *Store) GetParameterByVersion(name string, version int64) (*Parameter, error) {
	if _, err := s.GetParameter(name, false); err != nil {
		return nil, err
	}

	hKey := s.historyKey(name, version)
	var paramVersion ParameterVersion
	if err := s.historyStore.Get(hKey, &paramVersion); err != nil {
		return nil, ErrParameterVersionNotFound
	}

	return &Parameter{
		Name:             paramVersion.ParameterName,
		Value:            paramVersion.Value,
		Type:             paramVersion.Type,
		Version:          paramVersion.Version,
		LastModifiedDate: paramVersion.LastModifiedDate,
		DataType:         paramVersion.DataType,
		ARN:              s.buildARN(paramVersion.ParameterName),
	}, nil
}

// GetParameterByLabel retrieves a parameter by its label.
func (s *Store) GetParameterByLabel(name string, label string) (*Parameter, error) {
	param, err := s.GetParameter(name, false)
	if err != nil {
		return nil, err
	}

	var targetVersion int64
	for ver, labels := range param.VersionLabels {
		for _, l := range labels {
			if l == label {
				targetVersion = ver
				break
			}
		}
		if targetVersion != 0 {
			break
		}
	}

	if targetVersion == 0 {
		return nil, ErrParameterLabelNotFound
	}

	return s.GetParameterByVersion(name, targetVersion)
}

// GetParameters retrieves multiple parameters by names.
func (s *Store) GetParameters(names []string, withDecryption bool) ([]*Parameter, []string) {
	var parameters []*Parameter
	var invalidNames []string

	for _, name := range names {
		param, err := s.GetParameter(name, withDecryption)
		if err != nil {
			invalidNames = append(invalidNames, name)
			continue
		}
		parameters = append(parameters, param)
	}

	return parameters, invalidNames
}

// GetParameterHistory retrieves the version history of a parameter.
// Results are returned in reverse chronological order (newest first),
// matching the AWS API behaviour.
// Returns the history items along with pagination metadata (NextMarker, IsTruncated).
func (s *Store) GetParameterHistory(name string, maxResults int32, marker string) ([]*ParameterVersion, string, bool, error) {
	param, err := s.GetParameter(name, false)
	if err != nil {
		return nil, "", false, err
	}

	opts := common.ListOptions{MaxItems: int(maxResults), Marker: marker}
	if maxResults <= 0 {
		opts.MaxItems = 50
	}

	result, err := common.List[ParameterVersion](s.historyStore, opts, func(pv *ParameterVersion) bool {
		return pv.ParameterName == name
	})
	if err != nil {
		return nil, "", false, err
	}

	for _, pv := range result.Items {
		if labels, ok := param.VersionLabels[pv.Version]; ok {
			pv.Labels = labels
		} else {
			pv.Labels = []string{}
		}
	}

	// Reverse results to return newest version first, matching AWS behaviour.
	for i, j := 0, len(result.Items)-1; i < j; i, j = i+1, j-1 {
		result.Items[i], result.Items[j] = result.Items[j], result.Items[i]
	}

	return result.Items, result.NextMarker, result.IsTruncated, nil
}

// DeleteParameter deletes a parameter.
func (s *Store) DeleteParameter(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.paramKey(name)
	if !s.Exists(key) {
		return ErrParameterNotFound
	}

	opts := common.ListOptions{MaxItems: 10000}
	result, err := common.List[ParameterVersion](s.historyStore, opts, func(pv *ParameterVersion) bool {
		return pv.ParameterName == name
	})
	if err != nil {
		return fmt.Errorf("failed to list parameter history: %w", err)
	}
	for _, pv := range result.Items {
		hKey := s.historyKey(name, pv.Version)
		if err := s.historyStore.Delete(hKey); err != nil {
			return fmt.Errorf("failed to delete parameter history: %w", err)
		}
	}

	if err := s.TagStore.Delete(name); err != nil {
		return fmt.Errorf("failed to delete parameter tags: %w", err)
	}

	return s.BaseStore.Delete(key)
}

// DeleteParameters deletes multiple parameters.
func (s *Store) DeleteParameters(names []string) (deleted []string, invalid []string) {
	for _, name := range names {
		if err := s.DeleteParameter(name); err == nil {
			deleted = append(deleted, name)
		} else {
			invalid = append(invalid, name)
		}
	}
	return deleted, invalid
}

// DescribeParameters lists parameters with optional filters.
func (s *Store) DescribeParameters(filters map[string]string, maxResults int32, marker string) ([]*ParameterMetadata, string, error) {
	opts := common.ListOptions{
		MaxItems: int(maxResults),
		Marker:   marker,
	}
	if maxResults <= 0 {
		opts.MaxItems = 50
	}

	result, err := common.List[Parameter](s.BaseStore, opts, func(p *Parameter) bool {
		for k, v := range filters {
			switch k {
			case "Type":
				if string(p.Type) != v {
					return false
				}
			case "KeyId":
				if p.KeyID != v {
					return false
				}
			case "Tier":
				if string(p.Tier) != v {
					return false
				}
			case "Name":
				if !strings.HasPrefix(p.Name, v) {
					return false
				}
			case "Path":
				if !strings.HasPrefix(p.Name, v) {
					return false
				}
			}
		}
		return true
	})
	if err != nil {
		return nil, "", err
	}

	var metadata []*ParameterMetadata
	for _, p := range result.Items {
		metadata = append(metadata, NewParameterMetadata(p))
	}

	return metadata, result.NextMarker, nil
}

// GetParametersByPath retrieves parameters under a specific path.
func (s *Store) GetParametersByPath(path string, recursive bool, withDecryption bool, maxResults int32, marker string) ([]*Parameter, string, error) {
	opts := common.ListOptions{
		MaxItems: int(maxResults),
		Marker:   marker,
	}
	if maxResults <= 0 {
		opts.MaxItems = 10
	}

	result, err := common.List[Parameter](s.BaseStore, opts, func(p *Parameter) bool {
		if recursive {
			if !strings.HasPrefix(p.Name, path) {
				return false
			}
			if len(p.Name) > len(path) && p.Name[len(path)] != '/' {
				return false
			}
			return true
		}
		return strings.HasPrefix(p.Name, path) && !strings.Contains(strings.TrimPrefix(strings.TrimPrefix(p.Name, path), "/"), "/")
	})
	if err != nil {
		return nil, "", err
	}

	return result.Items, result.NextMarker, nil
}

// AddTagsToResource adds tags to a parameter.
func (s *Store) AddTagsToResource(name string, tags map[string]string) error {
	if _, err := s.GetParameter(name, false); err != nil {
		return err
	}
	return s.TagStore.Tag(name, tags)
}

// RemoveTagsFromResource removes tags from a parameter.
func (s *Store) RemoveTagsFromResource(name string, tagKeys []string) error {
	if _, err := s.GetParameter(name, false); err != nil {
		return err
	}
	return s.TagStore.Untag(name, tagKeys)
}

// ListTagsForResource retrieves tags for a parameter.
func (s *Store) ListTagsForResource(name string) (map[string]string, error) {
	if _, err := s.GetParameter(name, false); err != nil {
		return nil, err
	}
	return s.TagStore.List(name)
}

var parameterNameRegex = regexp.MustCompile(`^/([a-zA-Z0-9_.-]+/)*[a-zA-Z0-9_.-]+$|^[a-zA-Z0-9_.-]+$`)

var reservedPrefixes = []string{"aws", "ssm"}

// ValidateParameterName validates a parameter name.
func ValidateParameterName(name string) error {
	if name == "" {
		return ErrInvalidParameterName
	}
	if len(name) > 2048 {
		return ErrInvalidParameterName
	}
	if !parameterNameRegex.MatchString(name) {
		return ErrInvalidParameterName
	}
	trimmed := strings.TrimPrefix(name, "/")
	if idx := strings.Index(trimmed, "/"); idx != -1 {
		trimmed = trimmed[:idx]
	}
	lower := strings.ToLower(trimmed)
	for _, prefix := range reservedPrefixes {
		if lower == prefix {
			return ErrReservedParameterName
		}
	}
	return nil
}

// LabelParameterVersion adds labels to a specific version of a parameter.
func (s *Store) LabelParameterVersion(name string, parameterVersion int64, labels []string) error {
	key := s.paramKey(name)

	s.mu.Lock()
	defer s.mu.Unlock()

	var param Parameter
	if err := s.BaseStore.Get(key, &param); err != nil {
		return ErrParameterNotFound
	}

	hKey := s.historyKey(name, parameterVersion)
	var paramVersion ParameterVersion
	if err := s.historyStore.Get(hKey, &paramVersion); err != nil {
		return ErrParameterVersionNotFound
	}

	if param.VersionLabels == nil {
		param.VersionLabels = make(map[int64][]string)
	}

	for _, label := range labels {
		for ver, existingLabels := range param.VersionLabels {
			if ver != parameterVersion {
				for i, l := range existingLabels {
					if l == label {
						param.VersionLabels[ver] = append(existingLabels[:i], existingLabels[i+1:]...)
						break
					}
				}
			}
		}
	}

	currentLabels := param.VersionLabels[parameterVersion]
	for _, label := range labels {
		found := false
		for _, l := range currentLabels {
			if l == label {
				found = true
				break
			}
		}
		if !found {
			currentLabels = append(currentLabels, label)
		}
	}
	param.VersionLabels[parameterVersion] = currentLabels

	return s.Put(key, param)
}

// UnlabelParameterVersion removes labels from a specific version of a parameter.
func (s *Store) UnlabelParameterVersion(name string, parameterVersion int64, labels []string) ([]string, error) {
	key := s.paramKey(name)

	s.mu.Lock()
	defer s.mu.Unlock()

	var param Parameter
	if err := s.BaseStore.Get(key, &param); err != nil {
		return nil, ErrParameterNotFound
	}

	hKey := s.historyKey(name, parameterVersion)
	var paramVersion ParameterVersion
	if err := s.historyStore.Get(hKey, &paramVersion); err != nil {
		return nil, ErrParameterVersionNotFound
	}

	if param.VersionLabels == nil {
		return []string{}, nil
	}

	currentLabels, ok := param.VersionLabels[parameterVersion]
	if !ok {
		return []string{}, nil
	}

	var removedLabels []string
	remainingLabels := make([]string, 0, len(currentLabels))
	for _, existing := range currentLabels {
		found := false
		for _, toRemove := range labels {
			if existing == toRemove {
				found = true
				removedLabels = append(removedLabels, existing)
				break
			}
		}
		if !found {
			remainingLabels = append(remainingLabels, existing)
		}
	}

	if len(remainingLabels) == 0 {
		delete(param.VersionLabels, parameterVersion)
	} else {
		param.VersionLabels[parameterVersion] = remainingLabels
	}

	if err := s.Put(key, param); err != nil {
		return nil, err
	}

	return removedLabels, nil
}
