package ssm

// SSMStoreInterface defines operations for managing SSM parameters.
type SSMStoreInterface interface {
	PutParameter(param *Parameter, overwrite bool) (int64, error)
	GetParameter(name string, withDecryption bool) (*Parameter, error)
	GetParameterByVersion(name string, version int64) (*Parameter, error)
	GetParameterByLabel(name string, label string) (*Parameter, error)
	GetParameters(names []string, withDecryption bool) ([]*Parameter, []string)
	GetParameterHistory(name string, maxResults int32, marker string) ([]*ParameterVersion, string, bool, error)
	DeleteParameter(name string) error
	DeleteParameters(names []string) (deleted []string, invalid []string)
	DescribeParameters(filters map[string]string, maxResults int32, marker string) ([]*ParameterMetadata, string, error)
	GetParametersByPath(path string, recursive bool, withDecryption bool, maxResults int32, marker string) ([]*Parameter, string, error)
	AddTagsToResource(name string, tags map[string]string) error
	RemoveTagsFromResource(name string, tagKeys []string) error
	ListTagsForResource(name string) (map[string]string, error)
	LabelParameterVersion(name string, parameterVersion int64, labels []string) error
	UnlabelParameterVersion(name string, parameterVersion int64, labels []string) ([]string, error)
}

var _ SSMStoreInterface = (*Store)(nil)
