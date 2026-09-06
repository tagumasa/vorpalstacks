package apigateway

import "strings"

// EscapeResourcePath renders a resource path as the single RFC 6901 pointer
// token the method-settings key convention addresses it with: "~" becomes
// "~0" first, then every "/" becomes "~1" ("/pets/{petId}" →
// "~1pets~1{petId}", the root "/" → "~1").
func EscapeResourcePath(path string) string {
	escaped := strings.ReplaceAll(path, "~", "~0")
	return strings.ReplaceAll(escaped, "/", "~1")
}

// MethodSettingsKey derives the {resource_path}/{http_method} map key the
// stage methodSettings and usage-plan throttle maps share, for a concrete
// route. The official CLI update-stage example shows the stored keys as the
// addressed pointer token itself: the patch path
// "/~1resourceName/GET/logging/dataTrace" yields the key
// "~1resourceName/GET". The Stage methodSettings model documentation
// defines the keys as "{resource_path}/{http_method} ... or /*/* for
// overriding all methods in the stage". The control-plane patch appliers
// store exactly this form (they keep the addressed tokens verbatim), so the
// execution plane must derive its lookup candidates through this function
// to stay on the same contract.
func MethodSettingsKey(resourcePath, httpMethod string) string {
	return EscapeResourcePath(resourcePath) + "/" + httpMethod
}
