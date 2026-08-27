package iot

import (
	"time"

	storecommon "vorpalstacks/internal/store/aws/common"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input / Result structs
// ---------------------------------------------------------------------------

// CreateThingInput carries every field that CreateThing needs in a
// format independent of the wire protocol (HTTP REST-JSON vs gRPC-Web).
type CreateThingInput struct {
	ThingName        string
	ThingTypeName    string
	BillingGroupName string
	Attributes       map[string]string
}

// ListThingsInput holds the parameters for listing things.
type ListThingsInput struct {
	AttributeName  string
	AttributeValue string
	ThingTypeName  string
	NextToken      string
	MaxItems       int
}

// CreatePolicyInput carries the fields for CreatePolicy.
type CreatePolicyInput struct {
	PolicyName     string
	PolicyDocument string
}

// UpdateCertificateInput carries the fields for UpdateCertificate.
type UpdateCertificateInput struct {
	CertificateID string
	NewStatus     string
}

// UpdateThingInput carries the fields for UpdateThing.
type UpdateThingInput struct {
	ThingName       string
	Attributes      map[string]string
	MergeAttributes bool
	PayloadProvided bool
	ThingTypeName   string
	RemoveThingType bool
}

// CreateThingResult is the transport-agnostic result of CreateThing.
type CreateThingResult struct {
	Thing *iotstore.Thing
}

// UpdateThingResult is the transport-agnostic result of UpdateThing.
type UpdateThingResult struct {
	Thing *iotstore.Thing
}

// ListThingsResult is the transport-agnostic result of ListThings.
type ListThingsResult struct {
	Things    []*iotstore.Thing
	NextToken string
}

// CreatePolicyResult is the transport-agnostic result of CreatePolicy.
type CreatePolicyResult struct {
	Policy *iotstore.Policy
}

// ListPoliciesResult is the transport-agnostic result of ListPolicies.
type ListPoliciesResult struct {
	Policies  []*iotstore.Policy
	NextToken string
}

// ListCertificatesResult is the transport-agnostic result of ListCertificates.
type ListCertificatesResult struct {
	Certificates []*iotstore.Certificate
	NextToken    string
}

// DescribeThingResult is the transport-agnostic result of DescribeThing.
type DescribeThingResult struct {
	Thing *iotstore.Thing
	// BillingGroupName carries the thing's billing group, if any (at most
	// one per AWS constraints).
	BillingGroupName string
}

// GetPolicyResult is the transport-agnostic result of GetPolicy.
type GetPolicyResult struct {
	Policy *iotstore.Policy
}

// DescribeCertificateResult is the transport-agnostic result of DescribeCertificate.
type DescribeCertificateResult struct {
	Certificate *iotstore.Certificate
}

// ListTopicRulesResult is the transport-agnostic result of ListTopicRules.
type ListTopicRulesResult struct {
	Rules     []*iotstore.TopicRule
	NextToken string
}

// ---------------------------------------------------------------------------
// Core methods — single validation + persistence path
// ---------------------------------------------------------------------------

// createThingCore performs all validation and store operations for
// creating a thing. Both the HTTP API handler and the admin gRPC handler
// delegate to this method.
func (s *IoTService) createThingCore(store iotstore.IotStoreInterface, in CreateThingInput) (*CreateThingResult, error) {
	if in.ThingName == "" {
		return nil, iotstore.ErrMissingParam
	}
	if err := ValidateThingName(in.ThingName); err != nil {
		return nil, err
	}

	thing := &iotstore.Thing{
		ThingName:        in.ThingName,
		ThingTypeName:    in.ThingTypeName,
		Attributes:       in.Attributes,
		DefaultClientId:  in.ThingName,
		Version:          1,
		CreationDate:     time.Now().UTC(),
		LastModifiedDate: time.Now().UTC(),
	}

	created, err := store.CreateThing(thing)
	if err != nil {
		return nil, err
	}

	if in.BillingGroupName != "" {
		if err := store.AddThingToBillingGroup(in.ThingName, in.BillingGroupName); err != nil {
			return nil, err
		}
	}

	return &CreateThingResult{Thing: created}, nil
}

// listThingsCore lists things with optional filtering.
func (s *IoTService) listThingsCore(store iotstore.IotStoreInterface, in ListThingsInput) (*ListThingsResult, error) {
	maxItems := in.MaxItems
	if maxItems <= 0 {
		maxItems = 100
	}

	opts := iotstoreListOpts(maxItems, in.NextToken)

	if in.ThingTypeName != "" {
		result, err := store.ListThingsForThingType(in.ThingTypeName, opts)
		if err != nil {
			return nil, err
		}
		return &ListThingsResult{
			Things:    result.Items,
			NextToken: result.NextMarker,
		}, nil
	}

	result, err := store.ListThings(opts, in.AttributeName, in.AttributeValue)
	if err != nil {
		return nil, err
	}
	return &ListThingsResult{
		Things:    result.Items,
		NextToken: result.NextMarker,
	}, nil
}

// deleteThingCore removes a thing by name and cleans up its tags.
func (s *IoTService) deleteThingCore(store iotstore.IotStoreInterface, thingName, accountID, region string) error {
	if thingName == "" {
		return iotstore.ErrMissingParam
	}

	if err := store.DeleteThing(thingName); err != nil {
		return err
	}

	arn := iotstore.BuildThingARN(accountID, region, thingName)
	_ = store.DeleteAllTags(arn)
	return nil
}

// updateThingCore modifies a thing's attributes and optionally its type.
func (s *IoTService) updateThingCore(store iotstore.IotStoreInterface, in UpdateThingInput) (*UpdateThingResult, error) {
	if in.ThingName == "" {
		return nil, iotstore.ErrMissingParam
	}

	opts := iotstore.ThingUpdateOpts{
		Attributes:      in.Attributes,
		MergeAttributes: in.MergeAttributes,
		PayloadProvided: in.PayloadProvided,
		ThingTypeName:   in.ThingTypeName,
		RemoveThingType: in.RemoveThingType,
	}

	updated, err := store.UpdateThing(in.ThingName, opts)
	if err != nil {
		return nil, err
	}
	return &UpdateThingResult{Thing: updated}, nil
}

// createPolicyCore creates a policy with full validation.
func (s *IoTService) createPolicyCore(store iotstore.IotStoreInterface, in CreatePolicyInput) (*CreatePolicyResult, error) {
	if in.PolicyName == "" || in.PolicyDocument == "" {
		return nil, iotstore.ErrMissingParam
	}
	if !policyNamePattern.MatchString(in.PolicyName) {
		return nil, iotstore.ErrValidation
	}
	if err := validatePolicyDocument(in.PolicyDocument); err != nil {
		return nil, iotstore.ErrMalformedPolicy
	}

	p := &iotstore.Policy{
		PolicyName:       in.PolicyName,
		PolicyDocument:   in.PolicyDocument,
		CreationDate:     time.Now().UTC(),
		LastModifiedDate: time.Now().UTC(),
		Version:          1,
	}

	created, err := store.CreatePolicy(p)
	if err != nil {
		return nil, err
	}

	// The persisted version set is the single source of version identity:
	// version "1" (the initial default) is recorded alongside the policy so
	// list/get/default resolution never needs a record-side synthesis.
	v1 := map[string]interface{}{
		"versionId":        "1",
		"policyDocument":   created.PolicyDocument,
		"isDefaultVersion": true,
		"createDate":       created.CreationDate.Unix(),
	}
	if err := store.PutGeneric("policyVersions/"+in.PolicyName, map[string]interface{}{
		"versions": []interface{}{v1},
	}); err != nil {
		return nil, err
	}

	return &CreatePolicyResult{Policy: created}, nil
}

// listPoliciesCore lists policies with pagination.
func (s *IoTService) listPoliciesCore(store iotstore.IotStoreInterface, nextToken string, maxItems int) (*ListPoliciesResult, error) {
	if maxItems <= 0 {
		maxItems = 100
	}
	opts := iotstoreListOpts(maxItems, nextToken)
	result, err := store.ListPolicies(opts)
	if err != nil {
		return nil, err
	}
	return &ListPoliciesResult{
		Policies:  result.Items,
		NextToken: result.NextMarker,
	}, nil
}

// deletePolicyCore removes a policy by name and cleans up its tags.
func (s *IoTService) deletePolicyCore(store iotstore.IotStoreInterface, policyName, accountID, region string) error {
	if policyName == "" {
		return iotstore.ErrMissingParam
	}

	if err := store.DeletePolicy(policyName); err != nil {
		return err
	}

	arn := iotstore.BuildPolicyARN(accountID, region, policyName)
	_ = store.DeleteAllTags(arn)
	// A policy's versions cannot outlive the policy; a recreated policy
	// starts with only its default version.
	_ = store.DeleteGeneric("policyVersions/" + policyName)
	return nil
}

// listCertificatesCore lists certificates with pagination.
func (s *IoTService) listCertificatesCore(store iotstore.IotStoreInterface, nextToken string, maxItems int) (*ListCertificatesResult, error) {
	if maxItems <= 0 {
		maxItems = 100
	}
	opts := iotstoreListOpts(maxItems, nextToken)
	result, err := store.ListCertificates(opts)
	if err != nil {
		return nil, err
	}
	return &ListCertificatesResult{
		Certificates: result.Items,
		NextToken:    result.NextMarker,
	}, nil
}

// updateCertificateCore validates and applies a certificate status change.
func (s *IoTService) updateCertificateCore(store iotstore.IotStoreInterface, in UpdateCertificateInput) error {
	if in.CertificateID == "" || in.NewStatus == "" {
		return iotstore.ErrMissingParam
	}
	if !IsValidCertUpdateStatus(in.NewStatus) {
		return iotstore.ErrInvalidCertStatus
	}

	current, err := store.GetCertificate(in.CertificateID)
	if err != nil {
		return err
	}
	if !IsValidCertStateTransition(current.Status, in.NewStatus) {
		return iotstore.ErrInvalidCertStatus
	}

	opts := iotstore.CertificateUpdateOpts{NewStatus: in.NewStatus}
	_, err = store.UpdateCertificate(in.CertificateID, opts)
	return err
}

// deleteCertificateCore removes a certificate and cleans up tags.
func (s *IoTService) deleteCertificateCore(store iotstore.IotStoreInterface, certID, accountID, region string) error {
	if certID == "" {
		return iotstore.ErrMissingParam
	}

	if err := store.DeleteCertificate(certID); err != nil {
		return err
	}

	arn := iotstore.BuildCertificateARN(accountID, region, certID)
	_ = store.DeleteAllTags(arn)
	return nil
}

// iotstoreListOpts converts admin parameters to the store ListOptions.
func iotstoreListOpts(maxItems int, nextToken string) storecommon.ListOptions {
	return storecommon.ListOptions{
		MaxItems: maxItems,
		Marker:   nextToken,
	}
}

// describeThingCore retrieves a single thing by name.
func (s *IoTService) describeThingCore(store iotstore.IotStoreInterface, thingName string) (*DescribeThingResult, error) {
	if thingName == "" {
		return nil, iotstore.ErrMissingParam
	}
	thing, err := store.GetThing(thingName)
	if err != nil {
		return nil, err
	}
	// AWS: DescribeThing returns billingGroupName when the thing belongs to
	// a billing group (at most one per AWS constraints).
	result := &DescribeThingResult{Thing: thing}
	if groups, _ := store.ListBillingGroupsForThing(thingName); len(groups) > 0 {
		result.BillingGroupName = groups[0]
	}
	return result, nil
}

// getPolicyCore retrieves a single policy by name.
func (s *IoTService) getPolicyCore(store iotstore.IotStoreInterface, policyName string) (*GetPolicyResult, error) {
	if policyName == "" {
		return nil, iotstore.ErrMissingParam
	}
	policy, err := store.GetPolicy(policyName)
	if err != nil {
		return nil, err
	}
	return &GetPolicyResult{Policy: policy}, nil
}

// describeCertificateCore retrieves a single certificate by ID.
func (s *IoTService) describeCertificateCore(store iotstore.IotStoreInterface, certID string) (*DescribeCertificateResult, error) {
	if certID == "" {
		return nil, iotstore.ErrMissingParam
	}
	cert, err := store.GetCertificate(certID)
	if err != nil {
		return nil, err
	}
	return &DescribeCertificateResult{Certificate: cert}, nil
}

// listTopicRulesCore lists topic rules with pagination.
func (s *IoTService) listTopicRulesCore(store iotstore.IotStoreInterface, nextToken string, maxItems int) (*ListTopicRulesResult, error) {
	if maxItems <= 0 {
		maxItems = 100
	}
	opts := iotstoreListOpts(maxItems, nextToken)
	result, err := store.ListRules(opts)
	if err != nil {
		return nil, err
	}
	return &ListTopicRulesResult{
		Rules:     result.Items,
		NextToken: result.NextMarker,
	}, nil
}
