package iot

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/services/aws/iot/ca"
	iotstore "vorpalstacks/internal/store/aws/iot"
	vcrypto "vorpalstacks/internal/utils/crypto"
)

// Provisioning-template engine shared by RegisterThing and the bulk
// thing-registration task worker. The template format is defined by the
// AWS IoT provisioning-template guide: a Parameters section declaring
// name/type/default triples and a Resources section whose entries carry a
// caller-chosen logical name, a Type (AWS::IoT::Thing, AWS::IoT::Certificate,
// AWS::IoT::Policy), Properties, and optional OverrideSettings. Resources
// are identified by their Type — the logical name is arbitrary (the guide's
// own examples use "thing", "certificate", "policy") — so the engine never
// keys validation or provisioning on it.

const (
	thingResourceType       = "AWS::IoT::Thing"
	certificateResourceType = "AWS::IoT::Certificate"
	policyResourceType      = "AWS::IoT::Policy"
)

// OverrideSettings actions. MERGE is valid only for a thing's
// AttributePayload and ThingGroups properties; FAIL fails the request with
// ResourceConflictsException; OverrideSettings are not available for policy
// resources at all.
const (
	overrideDoNothing = "DO_NOTHING"
	overrideReplace   = "REPLACE"
	overrideFail      = "FAIL"
	overrideMerge     = "MERGE"
)

// thingOverrideProperties are the thing properties OverrideSettings may
// name; a certificate may override only Status.
var thingOverrideProperties = map[string]bool{
	"AttributePayload": true,
	"ThingTypeName":    true,
	"ThingGroups":      true,
}

type templateParameter struct {
	Default string
}

type templateResource struct {
	LogicalName      string
	Type             string
	Properties       map[string]interface{}
	OverrideSettings map[string]string
}

type provisioningTemplate struct {
	parameters map[string]templateParameter
	resources  []templateResource // sorted by logical name for determinism
}

// parseProvisioningTemplate validates the template body (the TemplateBody
// contract: required, within the documented length bound, a JSON object)
// and returns its parsed form. Structural rules: at least one
// AWS::IoT::Thing resource; every resource Type must be one of the three
// documented types; a certificate declares exactly one of
// CertificateSigningRequest, CertificateId, or CertificatePem; a policy
// carries PolicyName or PolicyDocument but not both; OverrideSettings name
// only the documented properties with the documented actions.
func parseProvisioningTemplate(body string) (*provisioningTemplate, error) {
	if body == "" {
		return nil, iotstore.ErrMissingParam
	}
	if len(body) > MaxTemplateBodyLength {
		return nil, iotstore.ErrInvalidRequest
	}
	var doc struct {
		Parameters map[string]json.RawMessage        `json:"Parameters"`
		Resources  map[string]map[string]interface{} `json:"Resources"`
	}
	if err := json.Unmarshal([]byte(body), &doc); err != nil {
		return nil, iotstore.ErrInvalidRequest
	}
	if doc.Resources == nil {
		return nil, iotstore.ErrInvalidRequest
	}
	tpl := &provisioningTemplate{parameters: map[string]templateParameter{}}
	for name, raw := range doc.Parameters {
		var p struct {
			Default string `json:"Default"`
		}
		_ = json.Unmarshal(raw, &p)
		tpl.parameters[name] = templateParameter{Default: p.Default}
	}
	names := make([]string, 0, len(doc.Resources))
	for name := range doc.Resources {
		names = append(names, name)
	}
	sort.Strings(names)
	hasThing := false
	for _, name := range names {
		def := doc.Resources[name]
		if def == nil {
			return nil, iotstore.ErrInvalidRequest
		}
		rtype, _ := def["Type"].(string)
		props, _ := def["Properties"].(map[string]interface{})
		var overrides map[string]string
		if raw, ok := def["OverrideSettings"].(map[string]interface{}); ok {
			overrides = map[string]string{}
			for k, v := range raw {
				action, _ := v.(string)
				overrides[k] = action
			}
		}
		res := templateResource{LogicalName: name, Type: rtype, Properties: props, OverrideSettings: overrides}
		switch rtype {
		case thingResourceType:
			hasThing = true
			for prop, action := range overrides {
				if !thingOverrideProperties[prop] || !validOverrideAction(action, prop) {
					return nil, iotstore.ErrInvalidRequest
				}
			}
		case certificateResourceType:
			if countCertificateDeclarations(props) != 1 {
				return nil, iotstore.ErrInvalidRequest
			}
			for prop, action := range overrides {
				if prop != "Status" || !validOverrideAction(action, prop) {
					return nil, iotstore.ErrInvalidRequest
				}
			}
		case policyResourceType:
			if len(overrides) > 0 {
				return nil, iotstore.ErrInvalidRequest
			}
			_, hasName := props["PolicyName"]
			_, hasDoc := props["PolicyDocument"]
			if hasName == hasDoc {
				return nil, iotstore.ErrInvalidRequest
			}
		default:
			return nil, iotstore.ErrInvalidRequest
		}
		tpl.resources = append(tpl.resources, res)
	}
	if !hasThing {
		return nil, iotstore.ErrInvalidRequest
	}
	return tpl, nil
}

// countCertificateDeclarations counts how many of the three documented
// certificate declaration methods the properties use; the guide requires
// exactly one.
func countCertificateDeclarations(props map[string]interface{}) int {
	n := 0
	for _, key := range []string{"CertificateSigningRequest", "CertificateId", "CertificatePem"} {
		if _, ok := props[key]; ok {
			n++
		}
	}
	return n
}

func validOverrideAction(action, prop string) bool {
	switch action {
	case overrideDoNothing, overrideReplace, overrideFail:
		return true
	case overrideMerge:
		return prop == "AttributePayload" || prop == "ThingGroups"
	default:
		return false
	}
}

// overrideAction returns the configured action for a property. Without an
// OverrideSettings section every property defaults to FAIL — the
// conservative outcome the guide describes for a template whose resource
// already exists; with a section present, unlisted properties are left
// untouched (DO_NOTHING).
func (r *templateResource) overrideAction(property string) string {
	if r.OverrideSettings == nil {
		return overrideFail
	}
	if action, ok := r.OverrideSettings[property]; ok {
		return action
	}
	return overrideDoNothing
}

// resolveValue resolves a template property value to a string. A {"Ref":X}
// object consults the caller parameters (X may be wrapped in double braces
// per the guide's examples) and falls back to the parameter's declared
// default; every other value is a literal.
func (t *provisioningTemplate) resolveValue(v interface{}, params map[string]string) string {
	m, ok := v.(map[string]interface{})
	if !ok || len(m) != 1 {
		return scalarString(v)
	}
	raw, ok := m["Ref"].(string)
	if !ok {
		return scalarString(v)
	}
	name := strings.Trim(raw, "{}")
	if val, ok := params[name]; ok && val != "" {
		return val
	}
	if p, ok := t.parameters[name]; ok {
		return p.Default
	}
	return ""
}

func scalarString(v interface{}) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	default:
		b, err := json.Marshal(x)
		if err != nil {
			return ""
		}
		return string(b)
	}
}

// resolveAttributePayload converts the AttributePayload property (an object
// whose values may be literals or Ref objects) to the attribute map.
func (t *provisioningTemplate) resolveAttributePayload(v interface{}, params map[string]string) map[string]string {
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	attrs := make(map[string]string, len(m))
	for k, val := range m {
		attrs[k] = t.resolveValue(val, params)
	}
	return attrs
}

// resolveStringList converts a list property whose entries may be literals
// or Ref objects.
func (t *provisioningTemplate) resolveStringList(v interface{}, params map[string]string) []string {
	list, ok := v.([]interface{})
	if !ok {
		return nil
	}
	out := make([]string, 0, len(list))
	for _, item := range list {
		out = append(out, t.resolveValue(item, params))
	}
	return out
}

// provisioningOutcome is the result of provisioning one template: the ARNs
// of every generated or reused resource keyed by the template's logical
// names, and the PEM of the provisioned certificate when the template
// declares one.
type provisioningOutcome struct {
	resourceArns   map[string]string
	certificatePem string
}

// provisionFromTemplate provisions the template's resources in dependency
// order: things first (attachments need them), then certificates (each
// attached to the first thing, honouring ThingPrincipalType exclusivity),
// then policies (each attached to the first certificate's principal).
// Partial failures are not rolled back — AWS does not document rollback
// for RegisterThing.
func (s *IoTService) provisionFromTemplate(store iotstore.IotStoreInterface, authority *ca.CertificateAuthority, tpl *provisioningTemplate, params map[string]string) (*provisioningOutcome, error) {
	out := &provisioningOutcome{resourceArns: map[string]string{}}
	var thing *iotstore.Thing
	for i := range tpl.resources {
		res := &tpl.resources[i]
		if res.Type != thingResourceType {
			continue
		}
		provisioned, err := s.provisionThing(store, res, tpl, params)
		if err != nil {
			return nil, err
		}
		if thing == nil {
			thing = provisioned
		}
		out.resourceArns[res.LogicalName] = provisioned.ThingARN
	}
	var certPrincipal string
	for i := range tpl.resources {
		res := &tpl.resources[i]
		if res.Type != certificateResourceType {
			continue
		}
		certARN, certPEM, err := s.provisionCertificate(store, authority, res, tpl, params, thing.ThingName)
		if err != nil {
			return nil, err
		}
		if certPrincipal == "" {
			certPrincipal = certARN
			out.certificatePem = certPEM
		}
		out.resourceArns[res.LogicalName] = certARN
	}
	for i := range tpl.resources {
		res := &tpl.resources[i]
		if res.Type != policyResourceType {
			continue
		}
		policyARN, err := s.provisionPolicy(store, res, tpl, params, certPrincipal)
		if err != nil {
			return nil, err
		}
		out.resourceArns[res.LogicalName] = policyARN
	}
	return out, nil
}

// provisionThing creates the thing resource or reconciles with an existing
// thing per OverrideSettings: DO_NOTHING keeps the existing thing, REPLACE
// overwrites the overridden properties, MERGE unions attributes and group
// memberships, and FAIL (the default without an OverrideSettings section)
// rejects the registration with ResourceConflictsException.
func (s *IoTService) provisionThing(store iotstore.IotStoreInterface, res *templateResource, tpl *provisioningTemplate, params map[string]string) (*iotstore.Thing, error) {
	props := res.Properties
	thingName := tpl.resolveValue(props["ThingName"], params)
	if thingName == "" {
		thingName = "registered-" + uuid.New().String()[:8]
	}
	attributes := tpl.resolveAttributePayload(props["AttributePayload"], params)
	thingTypeName := tpl.resolveValue(props["ThingTypeName"], params)
	groups := tpl.resolveStringList(props["ThingGroups"], params)
	billingGroup := tpl.resolveValue(props["BillingGroup"], params)

	existing, err := store.GetThing(thingName)
	if err != nil && !errors.Is(err, iotstore.ErrThingNotFound) {
		return nil, err
	}
	if existing != nil {
		attrAction := res.overrideAction("AttributePayload")
		typeAction := res.overrideAction("ThingTypeName")
		groupAction := res.overrideAction("ThingGroups")
		if attrAction == overrideFail || typeAction == overrideFail || groupAction == overrideFail {
			return nil, iotstore.ErrResourceConflicts
		}
		opts := iotstore.ThingUpdateOpts{
			Attributes:      attributes,
			MergeAttributes: attrAction == overrideMerge,
			PayloadProvided: attributes != nil,
		}
		if typeAction == overrideReplace {
			opts.ThingTypeName = thingTypeName
		}
		updated, err := store.UpdateThing(thingName, opts)
		if err != nil {
			return nil, err
		}
		if err := reconcileThingGroups(store, thingName, groups, groupAction); err != nil {
			return nil, err
		}
		return updated, nil
	}
	now := time.Now().UTC()
	created, err := store.CreateThing(&iotstore.Thing{
		ThingName:        thingName,
		ThingTypeName:    thingTypeName,
		Attributes:       attributes,
		BillingGroupName: billingGroup,
		Version:          1,
		CreationDate:     now,
		LastModifiedDate: now,
	})
	if err != nil {
		return nil, err
	}
	for _, group := range groups {
		if err := store.AddThingToThingGroup(thingName, group); err != nil {
			return nil, err
		}
	}
	return created, nil
}

// reconcileThingGroups applies the ThingGroups override action to an
// existing thing: MERGE adds the template's memberships, REPLACE makes the
// membership exactly the template's set, DO_NOTHING leaves them alone.
func reconcileThingGroups(store iotstore.IotStoreInterface, thingName string, desired []string, action string) error {
	switch action {
	case overrideMerge, overrideReplace:
	default:
		return nil
	}
	current, err := store.ListGroupsForThing(thingName)
	if err != nil {
		return err
	}
	currentSet := map[string]bool{}
	for _, g := range current {
		currentSet[g] = true
	}
	desiredSet := map[string]bool{}
	for _, g := range desired {
		desiredSet[g] = true
	}
	for _, g := range desired {
		if !currentSet[g] {
			if err := store.AddThingToThingGroup(thingName, g); err != nil {
				return err
			}
		}
	}
	if action == overrideReplace {
		for _, g := range current {
			if !desiredSet[g] {
				if err := store.RemoveThingFromThingGroup(thingName, g); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// provisionCertificate resolves the certificate resource's single declared
// method: a CSR is signed by the regional CA (RegisterThing never invokes
// certificate providers — that substitution belongs to the
// CreateCertificateFromCsr operation alone), an existing CertificateId is
// reused, and a PEM declaration registers the certificate (linked to its CA
// when CACertificatePem is present). The Status property applies on
// creation (default ACTIVE) and on an existing certificate only through a
// Status REPLACE override; ThingPrincipalType EXCLUSIVE_THING marks the
// attachment exclusive.
func (s *IoTService) provisionCertificate(store iotstore.IotStoreInterface, authority *ca.CertificateAuthority, res *templateResource, tpl *provisioningTemplate, params map[string]string, thingName string) (string, string, error) {
	props := res.Properties
	status := strings.ToUpper(tpl.resolveValue(props["Status"], params))
	if status == "" {
		status = "ACTIVE"
	}
	if status != "ACTIVE" && status != "INACTIVE" {
		return "", "", iotstore.ErrInvalidRequest
	}
	principalType := strings.ToUpper(tpl.resolveValue(props["ThingPrincipalType"], params))
	switch principalType {
	case "", "NON_EXCLUSIVE_THING", "EXCLUSIVE_THING":
	default:
		return "", "", iotstore.ErrInvalidRequest
	}

	var cert *iotstore.Certificate
	if raw, ok := props["CertificateSigningRequest"]; ok {
		if authority == nil {
			return "", "", fmt.Errorf("iot: certificate authority not available for the request region")
		}
		csrPEM := tpl.resolveValue(raw, params)
		if csrPEM == "" {
			return "", "", iotstore.ErrResourceRegistrationFailure
		}
		certPEM, certID, err := authority.IssueCertificateFromCSR(csrPEM)
		if err != nil {
			return "", "", err
		}
		cert = buildCertificateRecord(certPEM, certID, status == "ACTIVE")
		created, err := store.CreateCertificate(cert)
		if err != nil {
			return "", "", err
		}
		cert = created
	} else if raw, ok := props["CertificateId"]; ok {
		certID := tpl.resolveValue(raw, params)
		existing, err := store.GetCertificate(certID)
		if err != nil {
			return "", "", iotstore.ErrResourceRegistrationFailure
		}
		if res.overrideAction("Status") == overrideReplace && existing.Status != status {
			updated, err := store.UpdateCertificate(certID, iotstore.CertificateUpdateOpts{NewStatus: status})
			if err != nil {
				return "", "", err
			}
			existing = updated
		}
		cert = existing
	} else {
		// A PEM declaration — parseProvisioningTemplate guarantees presence.
		certPEM := tpl.resolveValue(props["CertificatePem"], params)
		if certPEM == "" {
			return "", "", iotstore.ErrResourceRegistrationFailure
		}
		certID := vcrypto.FingerprintPEM(certPEM)
		existing, err := store.GetCertificate(certID)
		if err == nil {
			cert = existing
		} else {
			if !errors.Is(err, iotstore.ErrCertificateNotFound) {
				return "", "", err
			}
			cert = buildCertificateRecord(certPEM, certID, status == "ACTIVE")
			if caPEM := tpl.resolveValue(props["CACertificatePem"], params); caPEM != "" {
				cert.CaCertificateID = vcrypto.FingerprintPEM(caPEM)
				cert.CertificateMode = "SNI_ONLY"
			}
			created, err := store.CreateCertificate(cert)
			if err != nil {
				return "", "", err
			}
			cert = created
		}
	}

	var attachErr error
	if principalType == "EXCLUSIVE_THING" {
		attachErr = store.AttachThingPrincipalExclusive(thingName, cert.CertificateARN)
	} else {
		attachErr = store.AttachThingPrincipal(thingName, cert.CertificateARN)
	}
	if attachErr != nil {
		return "", "", attachErr
	}
	return cert.CertificateARN, cert.CertificatePEM, nil
}

// provisionPolicy resolves the policy resource's single declaration: an
// existing PolicyName is reused, while a PolicyDocument creates a policy
// whose name defaults to a hash of the document ("Defaults to a hash of
// the policy document"). The policy is attached to the certificate
// principal; a template without a certificate cannot attach policies, so
// the registration fails.
func (s *IoTService) provisionPolicy(store iotstore.IotStoreInterface, res *templateResource, tpl *provisioningTemplate, params map[string]string, principal string) (string, error) {
	if principal == "" {
		return "", iotstore.ErrResourceRegistrationFailure
	}
	props := res.Properties
	policyName := tpl.resolveValue(props["PolicyName"], params)
	if policyName == "" {
		doc := policyDocumentString(props["PolicyDocument"], tpl, params)
		if err := validatePolicyDocument(doc); err != nil {
			return "", iotstore.ErrInvalidRequest
		}
		policyName = hashPolicyName(doc)
	}
	policy, err := store.GetPolicy(policyName)
	if err != nil {
		if !errors.Is(err, iotstore.ErrPolicyNotFound) {
			return "", err
		}
		if _, hasDoc := props["PolicyDocument"]; !hasDoc {
			// A named policy that does not exist cannot be attached.
			return "", iotstore.ErrResourceRegistrationFailure
		}
		created, cerr := s.createPolicyCore(store, CreatePolicyInput{
			PolicyName:     policyName,
			PolicyDocument: policyDocumentString(props["PolicyDocument"], tpl, params),
		})
		if cerr != nil {
			return "", cerr
		}
		policy = created.Policy
	}
	if err := store.AttachPolicyToPrincipal(policy.PolicyName, principal); err != nil {
		return "", err
	}
	return policy.PolicyARN, nil
}

// policyDocumentString renders the PolicyDocument property: the guide
// allows either an escaped JSON string or a nested object; a lone Ref
// object resolves against the parameters.
func policyDocumentString(v interface{}, tpl *provisioningTemplate, params map[string]string) string {
	switch x := v.(type) {
	case string:
		return x
	case map[string]interface{}:
		if len(x) == 1 {
			if _, ok := x["Ref"]; ok {
				return tpl.resolveValue(x, params)
			}
		}
		b, err := json.Marshal(x)
		if err != nil {
			return ""
		}
		return string(b)
	default:
		return ""
	}
}

// hashPolicyName derives the default policy name for a PolicyDocument-only
// declaration — the guide specifies "a hash of the policy document"; the
// digest's leading hex digits satisfy the policy name pattern.
func hashPolicyName(doc string) string {
	sum := sha256.Sum256([]byte(doc))
	return fmt.Sprintf("%x", sum)[:32]
}
