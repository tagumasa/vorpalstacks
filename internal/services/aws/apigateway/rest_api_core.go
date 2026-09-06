package apigateway

import (
	"fmt"
	"slices"
	"strings"

	types "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/store/aws/apigateway"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// CreateRestApiInput is the transport-agnostic DTO for creating a REST API.
// Both the HTTP API handler and the admin gRPC handler build this struct and
// delegate to createRestApiCore, so validation and persistence follow a
// single code path on both planes — handlers must not re-implement either.
type CreateRestApiInput struct {
	Name                      string
	Description               string
	Version                   string
	BinaryMediaTypes          []string
	ApiKeySource              string
	Policy                    string
	SecurityPolicy            string
	EndpointAccessMode        string
	DisableExecuteApiEndpoint bool
	MinimumCompressionSize    *int32
	EndpointTypes             []string
	Tags                      []types.Tag
	CloneFrom                 string
}

func (s *APIGatewayService) createRestApiCore(
	stores *apiGatewayStores,
	in CreateRestApiInput,
) (*apigateway.RestApi, error) {
	api := &apigateway.RestApi{
		Name:                      in.Name,
		Description:               in.Description,
		Version:                   in.Version,
		BinaryMediaTypes:          in.BinaryMediaTypes,
		ApiKeySource:              in.ApiKeySource,
		Policy:                    in.Policy,
		SecurityPolicy:            in.SecurityPolicy,
		EndpointAccessMode:        in.EndpointAccessMode,
		DisableExecuteApiEndpoint: in.DisableExecuteApiEndpoint,
		MinimumCompressionSize:    in.MinimumCompressionSize,
		Tags:                      in.Tags,
	}
	if len(in.EndpointTypes) > 0 {
		api.EndpointConfiguration = &apigateway.EndpointConfiguration{Types: in.EndpointTypes}
	}
	cloneFrom := in.CloneFrom
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
		return nil, toApiGatewayError(err)
	}

	if cloneFrom != "" {
		if err := stores.restApis.CloneFromSource(created, cloneFrom); err != nil {
			// Roll back the just-persisted target so that a clone failure
			// does not leave an empty API lingering in Pebble.
			_ = stores.restApis.Delete(created.Id)
			return nil, toApiGatewayError(err)
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
		return toApiGatewayError(err)
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
		return nil, toApiGatewayError(err)
	}

	for _, po := range patches {
		handled, err := applyRestApiPatch(api, po)
		if err != nil {
			return nil, toApiGatewayError(err)
		}
		if !handled {
			return nil, unknownPatchPathError(po)
		}
	}

	if err := stores.restApis.Update(api); err != nil {
		return nil, toApiGatewayError(err)
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
		if err := requirePatchOp(po, opReplace); err != nil {
			return true, err
		}
		api.Name = po.Value
	case "/description":
		// The row documents add, replace and remove.
		if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
			return true, err
		}
		api.Description = po.Value
	case "/version":
		if err := requirePatchOp(po, opReplace); err != nil {
			return true, err
		}
		api.Version = po.Value
	case "/apiKeySource":
		if err := requirePatchOp(po, opReplace); err != nil {
			return true, err
		}
		if !validateApiKeySource(po.Value) {
			return true, NewBadRequestException("Invalid apiKeySource: must be HEADER or AUTHORIZER")
		}
		api.ApiKeySource = po.Value
	case "/policy":
		if err := requirePatchOp(po, opReplace); err != nil {
			return true, err
		}
		api.Policy = po.Value
	case "/disableExecuteApiEndpoint":
		if err := requirePatchOp(po, opReplace); err != nil {
			return true, err
		}
		api.DisableExecuteApiEndpoint = po.Value == "true"
	case "/securityPolicy":
		if err := requirePatchOp(po, opReplace); err != nil {
			return true, err
		}
		if !validateSecurityPolicy(po.Value) {
			return true, NewBadRequestException("Invalid securityPolicy: must be TLS_1_0, TLS_1_2, or start with SecurityPolicy_")
		}
		api.SecurityPolicy = po.Value
	case "/endpointAccessMode":
		if err := requirePatchOp(po, opReplace); err != nil {
			return true, err
		}
		if !validateEndpointAccessMode(po.Value) {
			return true, NewBadRequestException("Invalid endpointAccessMode: must be BASIC or STRICT")
		}
		api.EndpointAccessMode = po.Value
	case "/minimumCompressionSize":
		if err := requirePatchOp(po, opReplace); err != nil {
			return true, err
		}
		// The table footnote: "To disable compression, apply a replace
		// operation with the value property set to null or omit the value
		// property." A null or absent value arrives as the empty string.
		if po.Value == "" || po.Value == "null" {
			api.MinimumCompressionSize = nil
			return true, nil
		}
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
			// The whole-member row documents add, replace and remove; the
			// value-token element form follows the same operations, while
			// numeric index and append-marker tokens address no documented
			// row and reject.
			if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
				return true, err
			}
			if err := applyBinaryMediaTypePatch(api, po); err != nil {
				return true, err
			}
			return true, nil
		}
		if po.Path == "/endpointConfiguration/ipAddressType" {
			// The row documents replace only, "Only dualstack and ipv4 are
			// supported"; the model notes dualstack is the only value a
			// PRIVATE endpoint admits.
			if err := requirePatchOp(po, opReplace); err != nil {
				return true, err
			}
			if !validateIpAddressType(po.Value) {
				return true, NewBadRequestException("Invalid ipAddressType: must be ipv4 or dualstack")
			}
			if api.EndpointConfiguration != nil &&
				slices.Contains(api.EndpointConfiguration.Types, "PRIVATE") && po.Value != "dualstack" {
				return true, NewBadRequestException("ipAddressType for a PRIVATE endpoint must be dualstack")
			}
			if api.EndpointConfiguration == nil {
				api.EndpointConfiguration = &apigateway.EndpointConfiguration{}
			}
			api.EndpointConfiguration.IpAddressType = po.Value
			return true, nil
		}
		if strings.HasPrefix(po.Path, "/endpointConfiguration/types") {
			if err := applyEndpointConfigPatch(api, po); err != nil {
				return true, err
			}
			return true, nil
		}
		if po.Path == "/endpointConfiguration/vpcEndpointIds" {
			if err := applyVpcEndpointIdsPatch(api, po); err != nil {
				return true, err
			}
			return true, nil
		}
		return false, nil
	}
	return true, nil
}

// applyBinaryMediaTypePatch applies the whole-member and value-token
// /binaryMediaTypes patch forms: add and replace append the value, remove
// drops the addressed media type. Numeric index and JSON-Pointer append
// tokens address no documented row and reject.
func applyBinaryMediaTypePatch(api *apigateway.RestApi, po PatchOperation) error {
	if rest, ok := strings.CutPrefix(po.Path, "/binaryMediaTypes/"); ok && isIndexToken(rest) {
		return unknownPatchPathError(po)
	}
	switch po.Op {
	case "add", "replace":
		if !slices.Contains(api.BinaryMediaTypes, po.Value) {
			api.BinaryMediaTypes = append(api.BinaryMediaTypes, po.Value)
		}
	case "remove":
		target := resolveBinaryMediaTypeToRemove(po.Path, po.Value)
		api.BinaryMediaTypes = removeString(api.BinaryMediaTypes, target)
	}
	return nil
}

// applyEndpointConfigPatch applies the /endpointConfiguration/types/{type}
// patch: the row documents replace only, with the type constrained to
// REGIONAL, EDGE, or PRIVATE, and replace sets the endpoint types to the
// addressed type.
func applyEndpointConfigPatch(api *apigateway.RestApi, po PatchOperation) error {
	if err := requirePatchOp(po, opReplace); err != nil {
		return err
	}
	typeName := strings.TrimPrefix(po.Path, "/endpointConfiguration/types/")
	if typeName == "" || !slices.Contains([]string{"REGIONAL", "EDGE", "PRIVATE"}, typeName) {
		return NewBadRequestException("Invalid endpoint type: must be REGIONAL, EDGE, or PRIVATE")
	}
	if api.EndpointConfiguration == nil {
		api.EndpointConfiguration = &apigateway.EndpointConfiguration{}
	}
	api.EndpointConfiguration.Types = []string{typeName}
	return nil
}

// applyVpcEndpointIdsPatch applies the /endpointConfiguration/vpcEndpointIds
// patch: the row documents add and remove for PRIVATE endpoint types only —
// add appends the addressed endpoint id, remove drops the matching id.
func applyVpcEndpointIdsPatch(api *apigateway.RestApi, po PatchOperation) error {
	if err := requirePatchOp(po, opAdd|opRemove); err != nil {
		return err
	}
	if api.EndpointConfiguration == nil || !slices.Contains(api.EndpointConfiguration.Types, "PRIVATE") {
		return NewBadRequestException("vpcEndpointIds patches are supported only for PRIVATE endpoint types")
	}
	if po.Op == "remove" {
		api.EndpointConfiguration.VpcEndpointIds = removeString(api.EndpointConfiguration.VpcEndpointIds, po.Value)
	} else if po.Value != "" && !slices.Contains(api.EndpointConfiguration.VpcEndpointIds, po.Value) {
		api.EndpointConfiguration.VpcEndpointIds = append(api.EndpointConfiguration.VpcEndpointIds, po.Value)
	}
	return nil
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
		return nil, toApiGatewayError(err)
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
