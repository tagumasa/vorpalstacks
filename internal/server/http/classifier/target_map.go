package classifier

import "strings"

// lookupServiceByTarget maps an X-Amz-Target header value to an internal
// service name. It first checks a case-sensitive map, then handles the
// com.amazonaws.* multi-part format, then falls back to a lowercase map,
// and finally returns the lowercased prefix as a last resort.
func lookupServiceByTarget(target string) string {
	parts := strings.Split(target, ".")

	caseSensitive := map[string]string{
		"AmazonSNS":                         "sns",
		"AmazonSQS":                         "sqs",
		"AmazonSSM":                         "ssm",
		"AWSCognitoIdentityProviderService": "cognito-idp",
		"AWSCognitoIdentityService":         "cognito-identity",
		"AWSEvents":                         "eventbridge",
		"Logs_20140328":                     "logs",
		"DynamoDB_20120810":                 "dynamodb",
		"Kinesis_20131202":                  "kinesis",
		"KMS_20141211":                      "kms",
		"TrentService":                      "kms",
		"secretsmanager":                    "secretsmanager",
		"AWSStepFunctions":                  "states",
		"Timestream_20181101":               "timestream-write",
		"GraniteServiceVersion20100801":     "monitoring",
		"AmazonAthena":                      "athena",
		"CertificateManager":                "acm",
		"AWSWAF_20190729":                   "wafv2",
		"AWSWAF_20150824":                   "waf",
		"CloudTrail_20131101":               "cloudtrail",
		"AWSSTS":                            "sts",
	}

	lowercase := map[string]string{
		"amazonsns":                         "sns",
		"amazonsqs":                         "sqs",
		"amazonssm":                         "ssm",
		"amazonlambda":                      "lambda",
		"amazoncloudwatch":                  "cloudwatch",
		"amazondynamodb":                    "dynamodb",
		"amazonkinesis":                     "kinesis",
		"amazonevents":                      "eventbridge",
		"amazonathena":                      "athena",
		"awscognitoidentityproviderservice": "cognito-idp",
		"awscognitoidentityservice":         "cognito-identity",
		"awsstepfunctions":                  "states",
		"amazoncloudwatchlogs":              "logs",
		"timestream_20181101":               "timestream-write",
		"timestreamquery_20181101":          "timestream-query",
	}

	if len(parts) > 0 {
		if svc, ok := caseSensitive[parts[0]]; ok {
			return svc
		}
		if parts[0] == "com" && len(parts) >= 3 && parts[1] == "amazonaws" {
			switch parts[2] {
			case "cloudtrail":
				return "cloudtrail"
			}
		}

		lowerPrefix := strings.ToLower(parts[0])
		if svc, ok := lowercase[lowerPrefix]; ok {
			return svc
		}
		return lowerPrefix
	}

	return strings.ToLower(target)
}
