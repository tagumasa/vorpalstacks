package apigateway

import (
	"fmt"
	"strings"

	"vorpalstacks/internal/store/aws/apigateway"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// createRestApiCore is the shared business logic for creating a REST API.
// Both the data-plane HTTP handler (rest_api_operations.go) and the admin
// gRPC handler (admin_handler_rest_api.go) call this function after their
// transport-specific input parsing. Keeping the logic here eliminates the
// duplicated implementation that previously caused:
//   - missing validation in the admin handler (apiKeySource, securityPolicy,
//     endpointAccessMode);
//   - shallow-copy cloneFrom leaving the source and target sharing the
//     same nested pointer maps so writes to the clone mutated the source;
//   - cloneFrom that did not stamp RestApiId on child elements so
//     resources cloned from a source API still referenced the source id;
//   - source-not-found cases persisting an empty target API in Pebble
//     because Create ran before the source lookup;
//   - missing endpointConfiguration default (REGIONAL) in the admin handler.
func (s *APIGatewayService) createRestApiCore(
	stores *apiGatewayStores,
	api *apigateway.RestApi,
	cloneFrom string,
) (*apigateway.RestApi, error) {
	if api.Name == "" {
		return nil, NewBadRequestException("name is required")
	}
	if !validateApiKeySource(api.ApiKeySource) {
		return nil, NewBadRequestException("Invalid apiKeySource: must be HEADER or AUTHORIZER")
	}
	if !validateSecurityPolicy(api.SecurityPolicy) {
		return nil, NewBadRequestException("Invalid securityPolicy: must be TLS_1_0, TLS_1_2, or start with SecurityPolicy_")
	}
	if !validateEndpointAccessMode(api.EndpointAccessMode) {
		return nil, NewBadRequestException("Invalid endpointAccessMode: must be BASIC or STRICT")
	}
	if api.MinimumCompressionSize != nil {
		v := *api.MinimumCompressionSize
		if !validateMinimumCompressionSize(v) {
			return nil, NewBadRequestException("minimumCompressionSize must be between 0 and 10485760")
		}
	}
	if api.EndpointConfiguration == nil || len(api.EndpointConfiguration.Types) == 0 {
		api.EndpointConfiguration = &apigateway.EndpointConfiguration{
			Types: []string{"REGIONAL"},
		}
	}

	if cloneFrom != "" {
		if !stores.restApis.Exists(cloneFrom) {
			return nil, NewNotFoundException("RestApi", cloneFrom)
		}
	}

	created, err := stores.restApis.Create(api)
	if err != nil {
		return nil, err
	}

	if cloneFrom != "" {
		if err := stores.restApis.CloneFromSource(created, cloneFrom); err != nil {
			// Roll back the just-persisted target so that a clone failure
			// does not leave an empty API lingering in Pebble.
			_ = stores.restApis.Delete(created.Id)
			return nil, err
		}
	}

	return created, nil
}

// getRestApiCore retrieves a single REST API by ID. Returns
// NewBadRequestException when apiId is empty and NotFoundException
// (mapped from storeErr) when no API matches the id.
func (s *APIGatewayService) getRestApiCore(stores *apiGatewayStores, apiId string) (*apigateway.RestApi, error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}
	return stores.restApis.Get(apiId)
}

// deleteRestApiCore deletes a REST API and cleans up any base path
// mappings that referenced it. The base path mapping cleanup is best-effort:
// if it fails we still consider the API deletion successful because the API
// itself is gone and dangling mappings are inert.
func (s *APIGatewayService) deleteRestApiCore(stores *apiGatewayStores, apiId string) error {
	if apiId == "" {
		return NewBadRequestException("restApiId is required")
	}
	if err := stores.restApis.Delete(apiId); err != nil {
		return err
	}
	_ = stores.domains.RemoveBasePathMappingsForApi(apiId)
	return nil
}

// updateRestApiCore applies the given patch operations to a REST API. It is
// the single source of truth for the patch path grammar; both the
// data-plane and admin handlers translate their transport-specific input
// into []PatchOperation and delegate here.
func (s *APIGatewayService) updateRestApiCore(
	stores *apiGatewayStores,
	apiId string,
	patches []PatchOperation,
) (*apigateway.RestApi, error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	stores.keyLocker.Lock(apiId)
	defer stores.keyLocker.Unlock(apiId)

	api, err := stores.restApis.Get(apiId)
	if err != nil {
		return nil, err
	}

	for _, po := range patches {
		handled, err := applyRestApiPatch(api, po)
		if err != nil {
			return nil, err
		}
		if !handled {
			return nil, NewBadRequestException("unknown patch path: " + po.Path)
		}
	}

	if err := stores.restApis.Update(api); err != nil {
		return nil, err
	}
	return api, nil
}

// applyRestApiPatch applies a single patch operation to a RestApi. Returns
// (handled, err): handled=false signals an unrecognised path so the caller
// can produce an explicit "unknown patch path" error; err is non-nil when
// the path is recognised but the value fails validation, carrying a
// descriptive message instead of degrading to "unknown patch path".
func applyRestApiPatch(api *apigateway.RestApi, po PatchOperation) (bool, error) {
	switch po.Path {
	case "/name":
		api.Name = po.Value
	case "/description":
		api.Description = po.Value
	case "/version":
		api.Version = po.Value
	case "/apiKeySource":
		if !validateApiKeySource(po.Value) {
			return true, NewBadRequestException("Invalid apiKeySource: must be HEADER or AUTHORIZER")
		}
		api.ApiKeySource = po.Value
	case "/policy":
		api.Policy = po.Value
	case "/disableExecuteApiEndpoint":
		api.DisableExecuteApiEndpoint = po.Value == "true"
	case "/securityPolicy":
		if !validateSecurityPolicy(po.Value) {
			return true, NewBadRequestException("Invalid securityPolicy: must be TLS_1_0, TLS_1_2, or start with SecurityPolicy_")
		}
		api.SecurityPolicy = po.Value
	case "/endpointAccessMode":
		if !validateEndpointAccessMode(po.Value) {
			return true, NewBadRequestException("Invalid endpointAccessMode: must be BASIC or STRICT")
		}
		api.EndpointAccessMode = po.Value
	case "/minimumCompressionSize":
		v, err := parseInt32(po.Value)
		if err != nil {
			return true, NewBadRequestException("Invalid minimumCompressionSize: not a number")
		}
		if !validateMinimumCompressionSize(v) {
			return true, NewBadRequestException("minimumCompressionSize must be between 0 and 10485760")
		}
		api.MinimumCompressionSize = &v
	default:
		if strings.HasPrefix(po.Path, "/binaryMediaTypes") {
			applyBinaryMediaTypePatch(api, po)
			return true, nil
		}
		if strings.HasPrefix(po.Path, "/endpointConfiguration/types") {
			applyEndpointConfigPatch(api, po)
			return true, nil
		}
		return false, nil
	}
	return true, nil
}

func applyBinaryMediaTypePatch(api *apigateway.RestApi, po PatchOperation) {
	switch po.Op {
	case "add", "replace":
		if !sliceContains(api.BinaryMediaTypes, po.Value) {
			api.BinaryMediaTypes = append(api.BinaryMediaTypes, po.Value)
		}
	case "remove":
		target := resolveBinaryMediaTypeToRemove(po.Path, po.Value, api.BinaryMediaTypes)
		api.BinaryMediaTypes = removeString(api.BinaryMediaTypes, target)
	}
}

func applyEndpointConfigPatch(api *apigateway.RestApi, po PatchOperation) {
	if api.EndpointConfiguration == nil {
		api.EndpointConfiguration = &apigateway.EndpointConfiguration{}
	}
	typeName := strings.TrimPrefix(po.Path, "/endpointConfiguration/types/")
	switch po.Op {
	case "add", "replace":
		if !sliceContains(api.EndpointConfiguration.Types, typeName) {
			api.EndpointConfiguration.Types = append(api.EndpointConfiguration.Types, typeName)
		}
	case "remove":
		api.EndpointConfiguration.Types = removeString(api.EndpointConfiguration.Types, typeName)
	}
}

// listRestApisCore returns a page of REST APIs. limit=0 falls back to the
// default page size; values above MaxPaginationLimit are rejected.
func (s *APIGatewayService) listRestApisCore(
	stores *apiGatewayStores,
	limit int,
	marker string,
) (*storecommon.ListResult[apigateway.RestApi], error) {
	resolved, err := resolvePageLimit(limit)
	if err != nil {
		return nil, err
	}
	return stores.restApis.List(storecommon.ListOptions{
		Marker:   marker,
		MaxItems: resolved,
	})
}

// resolvePageLimit validates a pre-parsed page size. A non-positive value
// falls back to DefaultPaginationLimit; values exceeding MaxPaginationLimit
// produce a BadRequestException to match AWS API Gateway behaviour.
func resolvePageLimit(limit int) (int, error) {
	if limit <= 0 {
		return DefaultPaginationLimit, nil
	}
	if limit > MaxPaginationLimit {
		return 0, NewBadRequestException(fmt.Sprintf("limit must not exceed %d", MaxPaginationLimit))
	}
	return limit, nil
}
