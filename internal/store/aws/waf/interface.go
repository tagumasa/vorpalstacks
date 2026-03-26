package waf

// WebACLStoreInterface defines operations for managing WAF Web ACLs.
type WebACLStoreInterface interface {
	Get(id string) (*WebACL, error)
	GetByARN(arn string) (*WebACL, error)
	Create(id, name, description, scope string, capacity int64, rules []*Rule, defaultAction *Action, visibilityConfig *VisibilityConfig) (*WebACL, error)
	Update(id, lockToken string, capacity int64, rules []*Rule, defaultAction *Action, visibilityConfig *VisibilityConfig) (*WebACL, error)
	Delete(id, lockToken string) error
	List(marker string, maxItems int) (*WebACLListResult, error)
	Raw() *WebACLStore
}

// RuleGroupStoreInterface defines operations for managing WAF Rule Groups.
type RuleGroupStoreInterface interface {
	Get(id string) (*RuleGroup, error)
	GetByARN(arn string) (*RuleGroup, error)
	Create(id, name, description string, capacity int64, rules []*Rule, visibilityConfig *VisibilityConfig) (*RuleGroup, error)
	Update(id, lockToken string, capacity int64, rules []*Rule, visibilityConfig *VisibilityConfig) (*RuleGroup, error)
	Delete(id, lockToken string) error
	List(marker string, maxItems int) (*RuleGroupListResult, error)
	Raw() *RuleGroupStore
}

// IPSetStoreInterface defines operations for managing WAF IP Sets.
type IPSetStoreInterface interface {
	Get(id string) (*IPSet, error)
	GetByARN(arn string) (*IPSet, error)
	Create(id, name, description, ipAddressVersion string, addresses []string) (*IPSet, error)
	Update(id, lockToken string, addresses []string) (*IPSet, error)
	Delete(id, lockToken string) error
	List(marker string, maxItems int) (*IPSetListResult, error)
	Raw() *IPSetStore
}

// RegexPatternSetStoreInterface defines operations for managing WAF Regex Pattern Sets.
type RegexPatternSetStoreInterface interface {
	Get(id string) (*RegexPatternSet, error)
	GetByARN(arn string) (*RegexPatternSet, error)
	Create(id, name, description string, regularPatterns []string) (*RegexPatternSet, error)
	Update(id, lockToken string, regularPatterns []string) (*RegexPatternSet, error)
	Delete(id, lockToken string) error
	List(marker string, maxItems int) (*RegexPatternSetListResult, error)
	Raw() *RegexPatternSetStore
}

// WebACLAssociationStoreInterface defines operations for managing WAF Web ACL associations.
type WebACLAssociationStoreInterface interface {
	Associate(webACLArn, resourceArn string) error
	Disassociate(webACLArn, resourceArn string) error
	GetByResourceArn(resourceArn string) (*WebACLAssociation, error)
	GetByWebACLArn(webACLArn string) ([]*WebACLAssociation, error)
	List() ([]*WebACLAssociation, error)
	Raw() *WebACLAssociationStore
}

// LoggingStoreInterface defines operations for managing WAF logging configurations.
type LoggingStoreInterface interface {
	Create(resourceArn string, logDestinationConfigs []string, logScope, logType string, loggingFilter *LoggingFilter, managedByFirewallManager bool, redactedFields []*FieldToMatch) (*LoggingConfiguration, error)
	Get(resourceArn string) (*LoggingConfiguration, error)
	Update(resourceArn string, logDestinationConfigs []string, logScope, logType string, loggingFilter *LoggingFilter, managedByFirewallManager bool, redactedFields []*FieldToMatch) (*LoggingConfiguration, error)
	Delete(resourceArn string) error
	List(scope string, marker string, maxItems int) (*LoggingConfigurationListResult, error)
	GetByResourceArn(resourceArn string) (*LoggingConfiguration, error)
	Raw() *LoggingStore
}

// WAFStoresInterface provides access to all WAF stores.
type WAFStoresInterface interface {
	WebACLStore() WebACLStoreInterface
	RuleGroupStore() RuleGroupStoreInterface
	IPSetStore() IPSetStoreInterface
	RegexPatternSetStore() RegexPatternSetStoreInterface
	AssociationStore() WebACLAssociationStoreInterface
	LoggingStore() LoggingStoreInterface
}

// WAFStores provides access to all WAF stores.
type WAFStores struct {
	webACLs          *WebACLStore
	ruleGroups       *RuleGroupStore
	ipSets           *IPSetStore
	regexPatternSets *RegexPatternSetStore
	associations     *WebACLAssociationStore
	loggingConfigs   *LoggingStore
}

// NewWAFStores creates a new WAF stores instance.
func NewWAFStores(webACLs *WebACLStore, ruleGroups *RuleGroupStore, ipSets *IPSetStore, regexPatternSets *RegexPatternSetStore, associations *WebACLAssociationStore, loggingConfigs *LoggingStore) *WAFStores {
	return &WAFStores{
		webACLs:          webACLs,
		ruleGroups:       ruleGroups,
		ipSets:           ipSets,
		regexPatternSets: regexPatternSets,
		associations:     associations,
		loggingConfigs:   loggingConfigs,
	}
}

// WebACLStore returns the Web ACL store.
func (s *WAFStores) WebACLStore() WebACLStoreInterface {
	if s.webACLs == nil {
		return nil
	}
	return s.webACLs
}

// RuleGroupStore returns the Rule Group store.
func (s *WAFStores) RuleGroupStore() RuleGroupStoreInterface {
	if s.ruleGroups == nil {
		return nil
	}
	return s.ruleGroups
}

// IPSetStore returns the IP Set store.
func (s *WAFStores) IPSetStore() IPSetStoreInterface {
	if s.ipSets == nil {
		return nil
	}
	return s.ipSets
}

// RegexPatternSetStore returns the Regex Pattern Set store.
func (s *WAFStores) RegexPatternSetStore() RegexPatternSetStoreInterface {
	if s.regexPatternSets == nil {
		return nil
	}
	return s.regexPatternSets
}

// AssociationStore returns the Web ACL association store.
func (s *WAFStores) AssociationStore() WebACLAssociationStoreInterface {
	if s.associations == nil {
		return nil
	}
	return s.associations
}

// LoggingStore returns the logging configuration store.
func (s *WAFStores) LoggingStore() LoggingStoreInterface {
	if s.loggingConfigs == nil {
		return nil
	}
	return s.loggingConfigs
}

// Raw returns the underlying Web ACL store.
func (s *WebACLStore) Raw() *WebACLStore { return s }

// Raw returns the underlying Rule Group store.
func (s *RuleGroupStore) Raw() *RuleGroupStore { return s }

// Raw returns the underlying IP Set store.
func (s *IPSetStore) Raw() *IPSetStore { return s }

// Raw returns the underlying Regex Pattern Set store.
func (s *RegexPatternSetStore) Raw() *RegexPatternSetStore {
	return s
}

// Raw returns the underlying Web ACL association store.
func (s *WebACLAssociationStore) Raw() *WebACLAssociationStore {
	return s
}

// Raw returns the underlying logging store.
func (s *LoggingStore) Raw() *LoggingStore { return s }
