// Package waflimits holds the WAF request-inspection limits shared by the
// enforcement planes (CloudFront, API Gateway, AppSync, Cognito) and the
// WAF store itself. The values describe the cross-service inspection
// contract, so their single definitions live here rather than in the WAF
// store package, which the enforcement planes must not import.
package waflimits

// DefaultBodyInspectionLimit is the default maximum body size in bytes a
// protected-resource plane forwards for WAF inspection (AWS WAF Developer
// Guide, "Oversize web request components": the default is 16 KB on the
// protected-resource types this platform hosts, with an upper bound of
// 64 KB).
const DefaultBodyInspectionLimit = 16384

// AppSyncBodyInspectionLimit is the fixed maximum body size in bytes the
// AppSync plane forwards for WAF inspection. The JsonBody
// OversizeHandling documentation fixes the Application Load Balancer and
// AppSync limit at 8 KB (8,192 bytes), independent of the configurable
// default.
const AppSyncBodyInspectionLimit = 8192
