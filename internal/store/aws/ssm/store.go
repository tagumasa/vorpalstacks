// Package ssm provides AWS Systems Manager Parameter Store storage functionality for vorpalstacks.
package ssm

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
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

// enforceAllowedPattern validates a parameter Value against an
// AllowedPattern. AWS returns ParameterPatternMismatchException when
// the value would violate the pattern. Empty pattern skips the check.
func enforceAllowedPattern(pattern, value string) error {
	if pattern == "" {
		return nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return ErrInvalidAllowedPattern
	}
	if !re.MatchString(value) {
		return ErrParameterPatternMismatch
	}
	return nil
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

	// Operation gate: if the parameter exists and the caller did not opt in
	// to Overwrite, reject the request before any input validation runs. AWS
	// returns ParameterAlreadyExists regardless of whether the new Value
	// would also violate an existing AllowedPattern.
	if err == nil && !overwrite {
		return 0, ErrParameterAlreadyExists
	}

	// Determine which AllowedPattern to validate against: the value in the
	// current request takes precedence, falling back to the previously-stored
	// pattern on Overwrite. A new parameter with no pattern specified skips
	// validation entirely.
	pattern := param.AllowedPattern
	if pattern == "" && err == nil {
		pattern = existingParam.AllowedPattern
	}
	if patternErr := enforceAllowedPattern(pattern, param.Value); patternErr != nil {
		return 0, patternErr
	}

	if err == nil {
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

	// Pre-validate the tags before any write so deterministic validation
	// failures cannot leave a partially-written parameter behind.
	if err := s.TagStore.ValidateTags(param.Tags); err != nil {
		return 0, err
	}

	// Snapshot the prior tag state for rollback; a storage-level failure
	// after the body is persisted must leave no partial state behind.
	priorTags, tagErr := s.TagStore.List(param.Name)
	if tagErr != nil {
		return 0, tagErr
	}
	rollback := func() {
		if len(param.Tags) > 0 {
			if err := s.TagStore.Untag(param.Name, tagKeySlice(param.Tags)); err != nil {
				logs.Warn("PutParameter rollback: failed to remove tags", logs.String("parameter", param.Name), logs.Err(err))
			}
			if len(priorTags) > 0 {
				if err := s.TagStore.Tag(param.Name, priorTags); err != nil {
					logs.Warn("PutParameter rollback: failed to restore prior tags", logs.String("parameter", param.Name), logs.Err(err))
				}
			}
		}
		if err == nil {
			// The parameter existed before this call: restore the prior body.
			if err := s.Put(key, &existingParam); err != nil {
				logs.Warn("PutParameter rollback: failed to restore prior body", logs.String("parameter", param.Name), logs.Err(err))
			}
		} else {
			if derr := s.BaseStore.Delete(key); derr != nil {
				logs.Warn("PutParameter rollback: failed to delete new parameter", logs.String("parameter", param.Name), logs.Err(derr))
			}
		}
	}

	if err := s.Put(key, param); err != nil {
		return 0, err
	}

	if len(param.Tags) > 0 {
		if err := s.TagStore.Tag(param.Name, param.Tags); err != nil {
			rollback()
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
		LastModifiedBy:   param.LastModifiedBy,
		DataType:         param.DataType,
		Description:      param.Description,
		KeyID:            param.KeyID,
		AllowedPattern:   param.AllowedPattern,
		Policies:         param.Policies,
	}
	hKey := s.historyKey(param.Name, param.Version)
	if err := s.historyStore.Put(hKey, historyEntry); err != nil {
		rollback()
		return 0, fmt.Errorf("failed to store parameter history: %w", err)
	}

	return param.Version, nil
}

// tagKeySlice extracts the keys of a tag map.
func tagKeySlice(tags map[string]string) []string {
	keys := make([]string, 0, len(tags))
	for k := range tags {
		keys = append(keys, k)
	}
	return keys
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
		Tier:             paramVersion.Tier,
		Description:      paramVersion.Description,
		KeyID:            paramVersion.KeyID,
		AllowedPattern:   paramVersion.AllowedPattern,
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
// Results are returned newest version first, matching the AWS API behaviour.
//
// History keys encode the version in decimal ("name:v%d"), so lexicographic
// key order does not match numeric version order (v10 sorts before v2) and
// reversing an ascending page breaks ordering across page boundaries. The
// full history for the parameter is therefore fetched, sorted numerically
// and paginated in memory. The marker is the decimal version number of the
// last entry returned; the next page contains entries with a strictly lower
// version. Returns the history items along with pagination metadata
// (NextMarker, IsTruncated); ErrInvalidNextToken is returned for markers
// that do not parse as a positive version number.
func (s *Store) GetParameterHistory(name string, maxResults int32, marker string) ([]*ParameterVersion, string, bool, error) {
	param, err := s.GetParameter(name, false)
	if err != nil {
		return nil, "", false, err
	}

	if maxResults <= 0 {
		maxResults = MaxPageResults
	}

	var versions []*ParameterVersion
	if err := common.ForEachAll[ParameterVersion](s.historyStore, name+":", func(pv *ParameterVersion) bool {
		return pv.ParameterName == name
	}, func(pv *ParameterVersion) error {
		versions = append(versions, pv)
		return nil
	}); err != nil {
		return nil, "", false, err
	}

	sort.Slice(versions, func(i, j int) bool {
		return versions[i].Version > versions[j].Version
	})

	start := 0
	if marker != "" {
		lastVersion, err := strconv.ParseInt(marker, 10, 64)
		if err != nil || lastVersion <= 0 {
			return nil, "", false, ErrInvalidNextToken
		}
		start = sort.Search(len(versions), func(i int) bool {
			return versions[i].Version < lastVersion
		})
	}

	end := start + int(maxResults)
	if end > len(versions) {
		end = len(versions)
	}

	page := versions[start:end]

	isTruncated := end < len(versions)
	nextMarker := ""
	if isTruncated && len(page) > 0 {
		nextMarker = strconv.FormatInt(page[len(page)-1].Version, 10)
	}

	for _, pv := range page {
		if labels, ok := param.VersionLabels[pv.Version]; ok {
			pv.Labels = labels
		} else {
			pv.Labels = []string{}
		}
	}

	return page, nextMarker, isTruncated, nil
}

// DeleteParameter deletes a parameter.
func (s *Store) DeleteParameter(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.paramKey(name)
	if !s.Exists(key) {
		return ErrParameterNotFound
	}

	err := common.ForEachAll[ParameterVersion](s.historyStore, name+":", func(pv *ParameterVersion) bool {
		return pv.ParameterName == name
	}, func(pv *ParameterVersion) error {
		hKey := s.historyKey(name, pv.Version)
		return s.historyStore.Delete(hKey)
	})
	if err != nil {
		return fmt.Errorf("failed to delete parameter history: %w", err)
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
func (s *Store) DescribeParameters(filters []ParameterFilter, maxResults int32, marker string) ([]*ParameterMetadata, string, error) {
	opts := common.ListOptions{
		MaxItems: int(maxResults),
		Marker:   marker,
	}

	result, err := common.List[Parameter](s.BaseStore, opts, func(p *Parameter) bool {
		return p.Matches(filters, false)
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
func (s *Store) GetParametersByPath(path string, recursive bool, withDecryption bool, filters []ParameterFilter, maxResults int32, marker string) ([]*Parameter, string, error) {
	opts := common.ListOptions{
		MaxItems: int(maxResults),
		Marker:   marker,
	}

	result, err := common.List[Parameter](s.BaseStore, opts, func(p *Parameter) bool {
		if !strings.HasPrefix(p.Name, path) {
			return false
		}
		if len(p.Name) == len(path) {
			return recursive && p.Matches(filters, true)
		}
		if recursive {
			return (strings.HasSuffix(path, "/") || p.Name[len(path)] == '/') && p.Matches(filters, true)
		}
		remainder := p.Name[len(path):]
		if len(remainder) > 0 && remainder[0] == '/' {
			remainder = remainder[1:]
		}
		return remainder != "" && !strings.Contains(remainder, "/") && p.Matches(filters, true)
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

// AWS specification limits for Parameter Store. These are the single source
// of truth; service handlers, validators and tests must reference them
// instead of inlining numeric literals.
const (
	// MaxParameterNameLength is the user-specifiable portion of the Smithy
	// PSParameterName length cap. The documented 2048-character maximum
	// includes 1037 characters reserved for internal use, leaving 1011
	// characters for the name a caller actually provides.
	MaxParameterNameLength = 1011

	// MaxHierarchyDepth is the maximum number of levels a parameter name
	// hierarchy may have ("/Engineering/Dev/ConnectionString" has 3 levels).
	MaxHierarchyDepth = 15

	// MaxParameterKeyIdLength is the Smithy ParameterKeyId length cap.
	MaxParameterKeyIdLength = 256

	// MaxParameterDescriptionLength is the Smithy ParameterDescription
	// length cap.
	MaxParameterDescriptionLength = 1024

	// MaxAllowedPatternLength is the Smithy AllowedPattern length cap.
	MaxAllowedPatternLength = 1024

	// MaxParameterDataTypeLength is the Smithy ParameterDataType length cap.
	// The valid value set is far shorter; the cap documents the trait.
	MaxParameterDataTypeLength = 128

	// MaxPageResults is the Smithy MaxResults range maximum shared by
	// DescribeParameters and GetParameterHistory (1-50, default 50).
	MaxPageResults = 50
)

// ValidateParameterName validates a parameter name against the AWS contract:
// regex shape, reserved top-level prefixes, the user-specifiable length cap
// and the hierarchy depth cap.
func ValidateParameterName(name string) error {
	if name == "" {
		return ErrInvalidParameterName
	}
	if len(name) > MaxParameterNameLength {
		return ErrInvalidParameterName
	}
	if !parameterNameRegex.MatchString(name) {
		return ErrInvalidParameterName
	}
	if strings.Count(name, "/") > MaxHierarchyDepth {
		return ErrHierarchyLevelLimitExceeded
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
// Returns the list of labels that could not be applied — AWS rejects only
// the labels that would push the version over the 10-label cap and accepts
// the rest.
func (s *Store) LabelParameterVersion(name string, parameterVersion int64, labels []string) ([]string, error) {
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
	var invalidLabels []string
	for _, label := range labels {
		found := false
		for _, l := range currentLabels {
			if l == label {
				found = true
				break
			}
		}
		if found {
			continue
		}
		if len(currentLabels) >= 10 {
			invalidLabels = append(invalidLabels, label)
			continue
		}
		currentLabels = append(currentLabels, label)
	}
	param.VersionLabels[parameterVersion] = currentLabels

	if err := s.Put(key, param); err != nil {
		return nil, err
	}
	return invalidLabels, nil
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
