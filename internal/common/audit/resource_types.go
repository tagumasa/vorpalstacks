package audit

import (
	"sort"
	"strings"

	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// arnResourceTypes maps an ARN's service to its resource-word vocabulary:
// the first resource segment (or the segment before a colon) selects the
// CloudTrail resource type identifier in the documented format
// AWS::{{aws-service-name}}::{{data-type-name}}. Services whose resource is
// a bare name (SQS queues, SNS topics) map through "", the word of a
// slash-free resource.
var arnResourceTypes = map[string]map[string]string{
	"cloudtrail": {
		"trail":          "AWS::CloudTrail::Trail",
		"eventdatastore": "AWS::CloudTrail::EventDataStore",
		"channel":        "AWS::CloudTrail::Channel",
	},
	"events": {
		"rule":      "AWS::Events::Rule",
		"event-bus": "AWS::Events::EventBus",
	},
	"dynamodb": {"table": "AWS::DynamoDB::Table"},
	"iam": {
		"user":   "AWS::IAM::User",
		"role":   "AWS::IAM::Role",
		"policy": "AWS::IAM::Policy",
		"group":  "AWS::IAM::Group",
	},
	"kms": {
		"key":   "AWS::KMS::Key",
		"alias": "AWS::KMS::Alias",
	},
	"lambda":         {"function": "AWS::Lambda::Function"},
	"logs":           {"log-group": "AWS::Logs::LogGroup"},
	"scheduler":      {"schedule": "AWS::Scheduler::Schedule"},
	"secretsmanager": {"secret": "AWS::SecretsManager::Secret"},
	"sns":            {"": "AWS::SNS::Topic"},
	"sqs":            {"": "AWS::SQS::Queue"},
	"states":         {"stateMachine": "AWS::StepFunctions::StateMachine"},
}

// resourceTypeFromARN classifies an ARN into the CloudTrail resource type
// identifier. S3 ARNs carry no type word (the resource is the bucket, or
// bucket and object), so the slash separates Bucket from Object; every
// other service classifies by resource word. A resource with neither
// separator is the bare name itself (an SQS queue, an SNS topic) and
// classifies through the empty word. An unmapped word returns an empty
// type: the record format permits a resources entry carrying only the ARN
// and accountId.
func resourceTypeFromARN(arn string) string {
	parsed, err := arnutil.ParseARN(arn)
	if err != nil {
		return ""
	}
	if parsed.Service == "s3" {
		if strings.Contains(parsed.Resource, "/") {
			return "AWS::S3::Object"
		}
		return "AWS::S3::Bucket"
	}
	words, ok := arnResourceTypes[parsed.Service]
	if !ok {
		return ""
	}
	resource := parsed.Resource
	if idx := strings.Index(resource, ":"); idx >= 0 {
		resource = resource[:idx]
	}
	if idx := strings.Index(resource, "/"); idx >= 0 {
		resource = resource[:idx]
	} else if !strings.Contains(parsed.Resource, ":") {
		resource = ""
	}
	return words[resource]
}

// collectResourceARNs walks the operation's request parameters and response
// and returns every distinct ARN value, ordered deterministically. The
// sweep is service-agnostic: any ARN the operation names is a resource the
// event references.
func collectResourceARNs(values ...map[string]interface{}) []string {
	seen := make(map[string]bool)
	var arns []string
	for _, value := range values {
		collectARNsFrom(value, seen, &arns)
	}
	sort.Strings(arns)
	return arns
}

func collectARNsFrom(value interface{}, seen map[string]bool, arns *[]string) {
	switch v := value.(type) {
	case map[string]interface{}:
		// Deterministic walk keeps the collection order stable.
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			collectARNsFrom(v[k], seen, arns)
		}
	case []interface{}:
		for _, item := range v {
			collectARNsFrom(item, seen, arns)
		}
	case string:
		if !seen[v] && strings.HasPrefix(v, "arn:") {
			if _, err := arnutil.ParseARN(v); err == nil {
				seen[v] = true
				*arns = append(*arns, v)
			}
		}
	}
}
