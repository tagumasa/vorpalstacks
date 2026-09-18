package cloudtrail

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/utils/aws/arn"
)

// Destination validation for CreateTrail and UpdateTrail. AWS verifies the
// trail's delivery destinations before accepting the configuration: the S3
// bucket must exist and its policy must grant CloudTrail write access, and
// the SNS topic must exist and its policy must grant CloudTrail publish
// access. The failures are the model-declared destination errors:
// S3BucketDoesNotExistException (the bucket is gone), and
// InsufficientS3BucketPolicyException / InsufficientSnsTopicPolicyException
// (the policy does not allow CloudTrail's delivery — the same error AWS
// returns for a missing topic: "SNS Topic does not exist, or the topic
// policy is incorrect").

// cloudTrailServicePrincipal is the service principal the documented
// destination policies grant ("Principal": {"Service":
// "cloudtrail.amazonaws.com"}).
const cloudTrailServicePrincipal = "cloudtrail.amazonaws.com"

// snsTopicInsufficientMessage covers both failure modes AWS folds into
// InsufficientSnsTopicPolicyException (observed service message for a
// missing topic; the same error family covers the policy check).
const snsTopicInsufficientMessage = "SNS topic does not exist, or the topic policy does not allow CloudTrail to publish to it"

// validateTrailS3Destination verifies the trail's S3 destination: the bucket
// exists and its policy grants the CloudTrail service principal the
// documented write access — s3:GetBucketAcl on the bucket and s3:PutObject
// on the trail's AWSLogs delivery prefix. trailARN is the identity the
// aws:SourceArn condition keys are evaluated against.
func (s *CloudTrailService) validateTrailS3Destination(ctx context.Context, storeRegion, accountID, bucket, keyPrefix, trailARN string) error {
	invoker := s.s3Invoker()
	if invoker == nil {
		// No S3 service means the destination cannot be verified; the
		// configuration is rejected rather than accepted unverified.
		return newInsufficientS3BucketPolicyException(
			fmt.Sprintf("Cannot verify the S3 bucket policy for %s: the S3 service is unavailable", bucket))
	}
	exists, err := invoker.BucketExists(ctx, storeRegion, bucket)
	if err != nil {
		return newInsufficientS3BucketPolicyException(
			fmt.Sprintf("Cannot verify the S3 bucket %s: %s", bucket, err.Error()))
	}
	if !exists {
		return newS3BucketDoesNotExistException(
			fmt.Sprintf("The specified S3 bucket does not exist: %s", bucket))
	}
	policy, err := invoker.GetBucketPolicy(ctx, storeRegion, bucket)
	if err != nil {
		return newInsufficientS3BucketPolicyException(
			fmt.Sprintf("Cannot read the policy of S3 bucket %s: %s", bucket, err.Error()))
	}
	if !s3BucketPolicySufficient(policy, bucket, keyPrefix, accountID, storeRegion, trailARN) {
		return newInsufficientS3BucketPolicyException(
			fmt.Sprintf("The policy on the S3 bucket %s is not sufficient: CloudTrail requires s3:GetBucketAcl on the bucket and s3:PutObject on the AWSLogs prefix for the service principal %s", bucket, cloudTrailServicePrincipal))
	}
	return nil
}

// resolveTrailSnsTopic resolves the SnsTopicName member — which carries
// "the name or ARN of the Amazon SNS topic" — to the topic ARN, verifying
// that the topic exists and its policy grants the CloudTrail service
// principal SNS:Publish on it.
func (s *CloudTrailService) resolveTrailSnsTopic(ctx context.Context, storeRegion, accountID, nameOrARN, trailARN string) (string, error) {
	if err := validateSnsTopicReference(nameOrARN); err != nil {
		return "", err
	}
	topicARN := nameOrARN
	if !strings.HasPrefix(nameOrARN, "arn:") {
		topicARN = arn.NewARNBuilder(accountID, storeRegion).SNS().Topic(nameOrARN)
	}
	invoker := s.snsInvoker()
	if invoker == nil {
		// No SNS service means the topic cannot be resolved; the
		// configuration is rejected rather than accepted unverified.
		return "", newInsufficientSnsTopicPolicyException(snsTopicInsufficientMessage)
	}
	if _, err := invoker.GetTopic(ctx, topicARN); err != nil {
		return "", newInsufficientSnsTopicPolicyException(snsTopicInsufficientMessage)
	}
	policy, err := invoker.GetTopicPolicy(ctx, topicARN)
	if err != nil {
		return "", newInsufficientSnsTopicPolicyException(snsTopicInsufficientMessage)
	}
	if !snsTopicPolicySufficient(policy, topicARN, trailARN, accountID) {
		return "", newInsufficientSnsTopicPolicyException(snsTopicInsufficientMessage)
	}
	return topicARN, nil
}

// validateSnsTopicReference validates the SnsTopicName member's dual form:
// a bare topic name follows the topic-name rules, an ARN reference must be
// an SNS topic ARN.
func validateSnsTopicReference(nameOrARN string) error {
	if strings.HasPrefix(nameOrARN, "arn:") {
		parsed, err := arn.ParseARN(nameOrARN)
		if err != nil || parsed.Service != "sns" || parsed.Resource == "" {
			return newInvalidSnsTopicNameException(
				fmt.Sprintf("SnsTopicName is not a valid SNS topic ARN: %s", nameOrARN))
		}
		return nil
	}
	return validateSnsTopicName(nameOrARN)
}

// policyStatement is the subset of the IAM statement shape the sufficiency
// checks evaluate: Effect, Principal, Action, Resource and Condition.
// Principal/Action/Resource each accept the single-value and list forms.
type policyStatement struct {
	Effect    string
	Principal policyPrincipal
	Action    []string
	Resource  []string
	Condition map[string]map[string]jsonStringOrList
}

// policyPrincipal is the statement principal: "*" as AWS or Service, or an
// explicit service/account list.
type policyPrincipal struct {
	Service []string
	AWS     []string
	Any     bool
}

// jsonStringOrList decodes a JSON member that appears as either a scalar
// string or an array of strings. Policy condition values may drop the
// quotation marks on numeric and Boolean scalars ("Quotation marks are
// optional for numeric and Boolean values", IAM JSON policy grammar), so
// those decode as their JSON literal — dropping the statement over one
// would erase a Deny from the evaluation, the fail-open direction.
func jsonStringOrListDecode(data []byte) ([]string, error) {
	var list []string
	if err := json.Unmarshal(data, &list); err == nil {
		return list, nil
	}
	var raws []json.RawMessage
	if err := json.Unmarshal(data, &raws); err == nil {
		values := make([]string, 0, len(raws))
		for _, raw := range raws {
			s, err := scalarPolicyValue(raw)
			if err != nil {
				return nil, err
			}
			values = append(values, s)
		}
		return values, nil
	}
	s, err := scalarPolicyValue(data)
	if err != nil {
		return nil, err
	}
	return []string{s}, nil
}

// scalarPolicyValue renders one policy value — a string, or a Boolean or
// numeric scalar written without its quotation marks — as the literal it
// was written as, the form string-valued condition operators compare
// against.
func scalarPolicyValue(raw []byte) (string, error) {
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s, nil
	}
	var b bool
	if err := json.Unmarshal(raw, &b); err == nil {
		return strconv.FormatBool(b), nil
	}
	var n json.Number
	if err := json.Unmarshal(raw, &n); err == nil {
		return n.String(), nil
	}
	return "", fmt.Errorf("expected a string or array of strings")
}

type jsonStringOrList []string

func (v *jsonStringOrList) UnmarshalJSON(data []byte) error {
	values, err := jsonStringOrListDecode(data)
	if err != nil {
		return err
	}
	*v = values
	return nil
}

// policyDocument parses the policy JSON into statements; Statement appears
// as either one object or an array.
func policyDocument(doc string) []policyStatement {
	var envelope struct {
		Statement json.RawMessage `json:"Statement"`
	}
	if err := json.Unmarshal([]byte(doc), &envelope); err != nil {
		return nil
	}
	var statements []json.RawMessage
	if err := json.Unmarshal(envelope.Statement, &statements); err != nil {
		// Single-object form.
		statements = []json.RawMessage{envelope.Statement}
	}
	parsed := make([]policyStatement, 0, len(statements))
	for _, raw := range statements {
		var st struct {
			Effect    string                                 `json:"Effect"`
			Principal json.RawMessage                        `json:"Principal"`
			Action    jsonStringOrList                       `json:"Action"`
			Resource  jsonStringOrList                       `json:"Resource"`
			Condition map[string]map[string]jsonStringOrList `json:"Condition"`
		}
		if err := json.Unmarshal(raw, &st); err != nil {
			continue
		}
		stmt := policyStatement{
			Effect:    st.Effect,
			Action:    st.Action,
			Resource:  st.Resource,
			Condition: st.Condition,
		}
		if len(st.Principal) > 0 {
			var p string
			if err := json.Unmarshal(st.Principal, &p); err == nil {
				stmt.Principal.Any = p == "*"
			} else {
				var pr struct {
					Service jsonStringOrList `json:"Service"`
					AWS     jsonStringOrList `json:"AWS"`
				}
				if err := json.Unmarshal(st.Principal, &pr); err == nil {
					stmt.Principal.Service = pr.Service
					stmt.Principal.AWS = pr.AWS
				}
			}
		}
		parsed = append(parsed, stmt)
	}
	return parsed
}

// coversServicePrincipal reports whether the statement's principal grants
// the CloudTrail service principal (or everyone).
func (p policyPrincipal) coversServicePrincipal() bool {
	if p.Any {
		return true
	}
	for _, v := range append(append([]string{}, p.Service...), p.AWS...) {
		if v == "*" || strings.EqualFold(v, cloudTrailServicePrincipal) {
			return true
		}
	}
	return false
}

// actionCovers reports whether the statement's actions include one that
// covers the wanted action. Action values are patterns where "*" and "?"
// may appear anywhere in the action name and matching is case-insensitive
// ("The prefix and the action name are case insensitive"; "You can also
// use wildcards (* or ?) as part of the action name", IAM JSON policy
// elements: Action) — "s3:*", "s3:Put*" and "S3:PutObject" all cover
// s3:PutObject.
func actionCovers(actions []string, action string) bool {
	target := strings.ToLower(action)
	for _, a := range actions {
		if a == "*" || wildcardMatch(strings.ToLower(a), target) {
			return true
		}
	}
	return false
}

// wildcardMatch reports whether the pattern — an IAM resource pattern
// where "*" matches any sequence and "?" any single character — matches
// the target exactly.
func wildcardMatch(pattern, target string) bool {
	quoted := regexp.QuoteMeta(pattern)
	quoted = strings.ReplaceAll(quoted, `\*`, ".*")
	quoted = strings.ReplaceAll(quoted, `\?`, ".")
	matched, err := regexp.MatchString("^"+quoted+"$", target)
	return err == nil && matched
}

// resourceCovers reports whether any of the statement's resource patterns
// covers the target ARN.
func resourceCovers(resources []string, target string) bool {
	for _, r := range resources {
		if r == "*" || wildcardMatch(r, target) {
			return true
		}
	}
	return false
}

// conditionsHold evaluates the statement's condition block against
// CloudTrail's own delivery context: aws:SourceArn is the trail ARN,
// aws:SourceAccount the owning account, aws:SecureTransport the delivery
// transport (the S3 invocation is in-process — no plaintext hop — so the
// key resolves true), and s3:x-amz-acl the object ACL CloudTrail sets.
// Operators follow their IAM semantics: StringEquals is exact matching,
// StringLike/ArnEquals/ArnLike are wildcard patterns ("The values can
// include multi-character match wildcards (*) and single-character match
// wildcards (?) anywhere in the string"), Bool matches the boolean
// literal. Condition operators outside these families, and keys outside
// the delivery context, fail closed — "If the key that you specify in a
// policy condition is not present in the request context, the values do
// not match and the condition is false" (IAM condition operators) — so a
// policy whose conditions cannot be satisfied by CloudTrail's documented
// request context is not sufficient.
func conditionsHold(conditions map[string]map[string]jsonStringOrList, trailARN, accountID string) bool {
	for operator, block := range conditions {
		switch operator {
		case "StringEquals", "StringLike", "ArnEquals", "ArnLike", "Bool", "BoolIfExists":
		default:
			return false
		}
		for key, values := range block {
			satisfied := false
			switch strings.ToLower(key) {
			case "aws:sourcearn":
				for _, v := range values {
					switch operator {
					case "StringEquals":
						if v == trailARN {
							satisfied = true
						}
					case "StringLike", "ArnEquals", "ArnLike":
						// "The ArnEquals and ArnLike condition
						// operators behave identically" and both take
						// multi- and single-character wildcards (IAM
						// condition operators reference), so every one
						// of these operators is a pattern match.
						if wildcardMatch(v, trailARN) {
							satisfied = true
						}
					}
				}
			case "aws:sourceaccount":
				for _, v := range values {
					switch operator {
					case "StringEquals":
						if v == accountID {
							satisfied = true
						}
					case "StringLike", "ArnEquals", "ArnLike":
						if wildcardMatch(v, accountID) {
							satisfied = true
						}
					}
				}
			case "aws:securetransport":
				// The delivery invokes S3 in-process — there is no
				// plaintext hop — so the transport key resolves true and
				// only a demanded "true" holds ("Boolean matching", IAM
				// condition operators).
				if operator == "Bool" || operator == "BoolIfExists" {
					for _, v := range values {
						if v == "true" {
							satisfied = true
						}
					}
				}
			case "s3:x-amz-acl":
				// CloudTrail's delivery always carries
				// bucket-owner-full-control; any other demanded ACL
				// blocks it.
				for _, v := range values {
					if v == "bucket-owner-full-control" {
						satisfied = true
					}
				}
			default:
				return false
			}
			if !satisfied {
				return false
			}
		}
	}
	return true
}

// s3LogDeliveryPrefix builds the bucket-relative folder the trail's log
// files are delivered under: [prefix/]AWSLogs/<account>/.
func s3LogDeliveryPrefix(keyPrefix, accountID string) string {
	if keyPrefix == "" {
		return "AWSLogs/" + accountID + "/"
	}
	return strings.Trim(keyPrefix, "/") + "/AWSLogs/" + accountID + "/"
}

// representativeObjectARN is a delivered log file's ARN form, used as the
// s3:PutObject coverage target: the file names vary, so the representative
// carries the naming tail a real file has.
func representativeObjectARN(bucket, keyPrefix, accountID, region string) string {
	name := logFileName(accountID, "CloudTrail", region, time.Now().UTC(), "0123456789abcdef")
	return fmt.Sprintf("arn:aws:s3:::%s/%s%s", bucket, s3LogDeliveryPrefix(keyPrefix, accountID), name)
}

// s3BucketPolicySufficient reports whether the bucket policy grants the
// CloudTrail service principal the documented access: s3:GetBucketAcl on
// the bucket and s3:PutObject on the trail's delivery prefix, per the
// AWS-documented bucket policy (the AWSCloudTrailAclCheck and
// AWSCloudTrailWrite statements). An explicit Deny that covers the same
// access overrides the grant, exactly as IAM evaluation order does.
// homeRegion names the trail's home region in the representative
// delivered-object ARN.
func s3BucketPolicySufficient(policy, bucket, keyPrefix, accountID, homeRegion, trailARN string) bool {
	if strings.TrimSpace(policy) == "" {
		return false
	}
	statements := policyDocument(policy)
	bucketARN := fmt.Sprintf("arn:aws:s3:::%s", bucket)
	// The coverage target is a representative delivered object; a policy
	// resource covering it covers the delivery prefix family.
	objectARN := representativeObjectARN(bucket, keyPrefix, accountID, homeRegion)
	aclOK, putOK := false, false
	aclDenied, putDenied := false, false
	for _, st := range statements {
		if !st.Principal.coversServicePrincipal() {
			continue
		}
		if !conditionsHold(st.Condition, trailARN, accountID) {
			continue
		}
		deny := strings.EqualFold(st.Effect, "Deny")
		if !deny && !strings.EqualFold(st.Effect, "Allow") {
			continue
		}
		if actionCovers(st.Action, "s3:GetBucketAcl") && resourceCovers(st.Resource, bucketARN) {
			if deny {
				aclDenied = true
			} else {
				aclOK = true
			}
		}
		if actionCovers(st.Action, "s3:PutObject") && resourceCovers(st.Resource, objectARN) {
			if deny {
				putDenied = true
			} else {
				putOK = true
			}
		}
	}
	return aclOK && putOK && !aclDenied && !putDenied
}

// snsTopicPolicySufficient reports whether the topic policy grants the
// CloudTrail service principal SNS:Publish on the topic, per the
// AWS-documented topic policy statement (AWSCloudTrailSNSPolicy). An
// explicit Deny that covers the same access overrides the grant, exactly
// as IAM evaluation order does.
func snsTopicPolicySufficient(policy, topicARN, trailARN, accountID string) bool {
	if strings.TrimSpace(policy) == "" {
		return false
	}
	publishOK, publishDenied := false, false
	for _, st := range policyDocument(policy) {
		if !st.Principal.coversServicePrincipal() {
			continue
		}
		if !conditionsHold(st.Condition, trailARN, accountID) {
			continue
		}
		deny := strings.EqualFold(st.Effect, "Deny")
		if !deny && !strings.EqualFold(st.Effect, "Allow") {
			continue
		}
		if actionCovers(st.Action, "SNS:Publish") && resourceCovers(st.Resource, topicARN) {
			if deny {
				publishDenied = true
			} else {
				publishOK = true
			}
		}
	}
	return publishOK && !publishDenied
}

// snsInvoker resolves the SNS invoker through the registry, or nil when SNS
// is not initialised.
func (s *CloudTrailService) snsInvoker() invokers.SNSInvoker {
	if s.invokerRegistry == nil {
		return nil
	}
	return s.invokerRegistry.SNSInvoker()
}

// logsInvoker resolves the CloudWatch Logs invoker through the registry, or
// nil when the Logs service is not initialised.
func (s *CloudTrailService) logsInvoker() invokers.LogsInvoker {
	if s.invokerRegistry == nil {
		return nil
	}
	return s.invokerRegistry.LogsInvoker()
}
