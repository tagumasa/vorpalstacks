package dispatcher

import (
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
)

// accessDeniedCodes maps internal service names to their AWS error code
// for IAM policy denials. Services using REST-XML or Query protocol return
// "AccessDenied" (no Exception suffix), while JSON and REST-JSON services
// return "AccessDeniedException". A few services have bespoke codes.
//
// Reference: each service's Smithy model defines the error shape name.
// Services without an explicit shape in Smithy inherit the gateway-level
// code from the protocol (Query → AccessDenied, JSON → AccessDeniedException).
var accessDeniedCodes = map[string]string{
	// REST-XML / Query protocol — "AccessDenied" (no suffix)
	"s3":         "AccessDenied",
	"sqs":        "AccessDenied",
	"sns":        "AccessDenied",
	"iam":        "AccessDenied",
	"sts":        "AccessDenied",
	"monitoring": "AccessDenied", // CloudWatch
	"cloudfront": "AccessDenied", // restXml
	"route53":    "AccessDenied", // restXml
	"neptune":    "AccessDenied", // awsQuery (RDS Neptune sub-service)

	// Bespoke error codes
	"cognito-identity": "NotAuthorizedException",
}

// accessDeniedErrorForService returns the service-specific access denied
// error. Unknown services default to "AccessDeniedException", which is the
// correct code for the majority of JSON and REST-JSON protocol services.
func accessDeniedErrorForService(serviceName string) *awserrors.AWSError {
	if code, ok := accessDeniedCodes[serviceName]; ok {
		return awserrors.NewAWSError(code, "Access denied", http.StatusForbidden)
	}
	return awserrors.NewAWSError("AccessDeniedException", "Access denied", http.StatusForbidden)
}
