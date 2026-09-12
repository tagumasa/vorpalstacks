package cognitoidentityprovider

// Pagination bound of the limit-bearing list requests. Source: the Smithy
// model's list-limit shapes — QueryLimitType, PoolQueryLimitType and
// QueryLimit all carry the range maximum 60. This is the single definition;
// every other site, store defaults and the service layer's request parsing
// alike, references it.

// MaxListLimit is the upper bound of the limit-bearing Cognito list
// requests; it doubles as the default page size when a request omits the
// limit or supplies zero.
const MaxListLimit = 60
