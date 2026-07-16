package iot

// Update option types for atomic read-modify-write operations in the store
// layer. Empty string fields mean "do not change". Pointer types are used for
// values where zero-value ambiguity exists (e.g., bool, int).

// ThingUpdateOpts specifies partial updates for UpdateThing.
type ThingUpdateOpts struct {
	Attributes      map[string]string // attributes to set or merge
	MergeAttributes bool              // true: merge into existing; false: replace
	PayloadProvided bool              // true if Attributes field was provided by the client
	ThingTypeName   string            // empty = no change
	RemoveThingType bool              // true to clear ThingTypeName
}

// ThingTypeUpdateOpts specifies partial updates for UpdateThingType.
type ThingTypeUpdateOpts struct {
	Description          string // empty = no change
	SearchableAttributes []string
}

// ThingGroupUpdateOpts specifies partial updates for UpdateThingGroup.
type ThingGroupUpdateOpts struct {
	Description     string            // empty = no change
	Attributes      map[string]string // attributes to set
	ExpectedVersion int64             // CAS check; 0 = skip
}

// BillingGroupUpdateOpts specifies partial updates for UpdateBillingGroup.
type BillingGroupUpdateOpts struct {
	Description     string // empty = no change
	ExpectedVersion int64  // CAS check; 0 = skip
}

// CertificateUpdateOpts specifies partial updates for UpdateCertificate.
type CertificateUpdateOpts struct {
	NewStatus string // must be a valid cert status
}

// JobUpdateOpts specifies partial updates for UpdateJob and CancelJob.
type JobUpdateOpts struct {
	Description string   // empty = no change
	Status      string   // empty = no change; "CANCELED" for CancelJob
	Targets     []string // nil = no change; set by AssociateTargetsWithJob
}

// AuthorizerUpdateOpts specifies partial updates for UpdateAuthorizer.
type AuthorizerUpdateOpts struct {
	FunctionARN    string
	TokenName      string
	TokenSignature string
	EnableCaching  *bool  // nil = no change
	Status         string // "" = no change, "ACTIVE"/"INACTIVE"
}

// ProvisioningTemplateUpdateOpts specifies partial updates for
// UpdateProvisioningTemplate.
type ProvisioningTemplateUpdateOpts struct {
	Description  string
	RoleARN      string
	Enabled      *bool   // nil = no change
	TemplateBody *string // nil = no change
}

// RoleAliasUpdateOpts specifies partial updates for UpdateRoleAlias.
type RoleAliasUpdateOpts struct {
	RoleARN         string
	DurationSeconds int64
}

// RuleUpdateOpts specifies partial updates for UpdateRule / ReplaceTopicRule.
type RuleUpdateOpts struct {
	SQL              string
	Description      string
	AwsIotSqlVersion string
	RuleDisabled     *bool                  // nil = no change; used by ReplaceTopicRule
	Actions          map[string]interface{} // nil = no change
	ErrorAction      map[string]interface{} // nil = no change
}
