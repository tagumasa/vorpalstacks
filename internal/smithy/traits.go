// Package smithy provides Smithy model handling for vorpalstacks.
package smithy

// TraitDocs represents the documentation trait.
const TraitDocs = "docs"

// TraitDocumentation represents the documentation trait.
const TraitDocumentation = "documentation"

// TraitSmithyAPIInput represents the input trait.
const TraitSmithyAPIInput = "smithy.api#input"

// TraitSmithyAPIOutput represents the output trait.
const TraitSmithyAPIOutput = "smithy.api#output"

// TraitSmithyAPIErrors represents the errors trait.
const TraitSmithyAPIErrors = "smithy.api#errors"

// TraitSmithyAPIHttp represents the HTTP trait.
const TraitSmithyAPIHttp = "smithy.api#http"

// TraitSmithyAPIReadonly represents the readonly trait.
const TraitSmithyAPIReadonly = "smithy.api#readonly"

// TraitSmithyAPIError represents the error trait.
const TraitSmithyAPIError = "smithy.api#error"

// TraitSmithyAPIHttpError represents the HTTP error trait.
const TraitSmithyAPIHttpError = "smithy.api#httpError"

// TraitSmithyAPIHttpCode represents the HTTP code trait.
const TraitSmithyAPIHttpCode = "smithy.api#httpCode"

// TraitSmithyAPIHttpLabel represents the HTTP label trait.
const TraitSmithyAPIHttpLabel = "smithy.api#httpLabel"

// TraitSmithyAPIHttpQuery represents the HTTP query trait.
const TraitSmithyAPIHttpQuery = "smithy.api#httpQuery"

// TraitSmithyAPIHttpHeader represents the HTTP header trait.
const TraitSmithyAPIHttpHeader = "smithy.api#httpHeader"

// TraitSmithyAPIHttpPayload represents the HTTP payload trait.
const TraitSmithyAPIHttpPayload = "smithy.api#httpPayload"

// TraitSmithyAPIJsonName represents the JSON name trait.
const TraitSmithyAPIJsonName = "smithy.api#jsonName"

// TraitSmithyAPIRequired represents the required trait.
const TraitSmithyAPIRequired = "smithy.api#required"

// TraitSmithyAPIHttpPrefix represents the HTTP prefix trait.
const TraitSmithyAPIHttpPrefix = "smithy.api#httpPrefix"

// TraitSmithyAPIHttpLocationName represents the HTTP location name trait.
const TraitSmithyAPIHttpLocationName = "smithy.api#httpLocationName"

// TraitSmithyHttpDefaultPayloadMediaType represents the default payload media type trait.
const TraitSmithyHttpDefaultPayloadMediaType = "smithy.http#defaultPayloadMediaType"

// TraitSmithyAPIEnumValue represents the enum value trait.
const TraitSmithyAPIEnumValue = "smithy.api#enumValue"

// TraitAwsAPIErrorCode represents the AWS error code trait.
const TraitAwsAPIErrorCode = "aws.api#errorCode"

// TraitAwsAPIService represents the AWS service trait.
const TraitAwsAPIService = "aws.api#service"

// TraitAwsAPIVersion represents the AWS version trait.
const TraitAwsAPIVersion = "aws.api#version"

// TraitAwsAPIAwsQueryCompatible represents the AWS query compatible trait.
const TraitAwsAPIAwsQueryCompatible = "aws.api#awsQueryCompatible"

// ShapeSkipTraits defines traits that should be skipped when processing shapes.
var ShapeSkipTraits = map[string]bool{
	TraitDocs:               true,
	TraitDocumentation:      true,
	TraitSmithyAPIInput:     true,
	TraitSmithyAPIOutput:    true,
	TraitSmithyAPIErrors:    true,
	TraitSmithyAPIHttp:      true,
	TraitSmithyAPIReadonly:  true,
	TraitSmithyAPIError:     true,
	TraitSmithyAPIHttpError: true,
	TraitSmithyAPIHttpCode:  true,
	TraitAwsAPIErrorCode:    true,
	TraitAwsAPIService:      true,
	TraitAwsAPIVersion:      true,
}

// MemberSkipTraits defines traits that should be skipped when processing shape members.
var MemberSkipTraits = map[string]bool{
	TraitDocs:                              true,
	TraitDocumentation:                     true,
	TraitSmithyAPIHttpLabel:                true,
	TraitSmithyAPIHttpQuery:                true,
	TraitSmithyAPIHttpHeader:               true,
	TraitSmithyAPIHttpPayload:              true,
	TraitSmithyAPIJsonName:                 true,
	TraitSmithyAPIRequired:                 true,
	TraitSmithyAPIHttpPrefix:               true,
	TraitSmithyHttpDefaultPayloadMediaType: true,
	TraitSmithyAPIHttpLocationName:         true,
}
