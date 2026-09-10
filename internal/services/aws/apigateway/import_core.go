package apigateway

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"

	"vorpalstacks/internal/store/aws/apigateway"
)

// importParameters carries the operation-specific query parameters of the
// import operations: ignore selects what the import skips (only
// "documentation" is documented), endpointConfigurationTypes overrides the
// created API's endpoint type, and basepath selects how the document's
// basePath is applied to the declared paths.
type importParameters struct {
	Ignore                     string
	EndpointConfigurationTypes string
	BasePath                   string
}

// parseImportParameters reads the parameters map from flat query entries.
// Only the documented keys are recognised; ignore=documentation is
// accepted as the platform default (documentation parts are not part of
// this platform) and any other value under ignore is rejected.
func parseImportParameters(params map[string]interface{}) (importParameters, error) {
	var p importParameters
	if v, ok := params["ignore"].(string); ok {
		if v != "documentation" {
			return p, NewBadRequestException("Invalid parameter ignore: " + v + "; only ignore=documentation is supported")
		}
		p.Ignore = v
	}
	if v, ok := params["endpointConfigurationTypes"].(string); ok {
		if v != "REGIONAL" && v != "EDGE" && v != "PRIVATE" {
			return p, NewBadRequestException("Invalid endpointConfigurationTypes: " + v + "; must be REGIONAL, EDGE, or PRIVATE")
		}
		p.EndpointConfigurationTypes = v
	}
	if v, ok := params["basepath"].(string); ok {
		if v != "ignore" && v != "prepend" && v != "split" {
			return p, NewBadRequestException("Invalid basepath: " + v + "; must be ignore, prepend, or split")
		}
		p.BasePath = v
	}
	return p, nil
}

// apiDefinitionDocument is the parsed view of an imported OpenAPI/Swagger
// document: the API-level metadata, the declared resource paths with their
// operations, the schemas, the declared authorizers and the API-key
// security scheme names.
type apiDefinitionDocument struct {
	Name          string
	Description   string
	Version       string
	BasePath      string
	Paths         map[string][]*apigateway.ImportedMethod
	Models        map[string]*apigateway.Model
	Authorizers   map[string]*apigateway.Authorizer
	ApiKeySchemes map[string]bool
	Warnings      []string
}

// importRestApiCore creates a REST API from an external API definition
// document. The document's metadata names the API, the paths and schemas
// become the API's resources and models, and the basePath is applied
// according to the basepath parameter. With failOnWarnings set, any
// warning recorded during parsing rolls the import back; otherwise the
// warnings travel with the created API's record.
func (s *APIGatewayService) importRestApiCore(
	stores *apiGatewayStores,
	body []byte,
	params importParameters,
	failOnWarnings bool,
) (*apigateway.RestApi, []string, error) {
	doc, err := s.parseApiDefinition(body)
	if err != nil {
		return nil, nil, err
	}
	if failOnWarnings && len(doc.Warnings) > 0 {
		return nil, nil, NewBadRequestException(
			fmt.Sprintf("API import failed: %d warning(s): %s", len(doc.Warnings), strings.Join(doc.Warnings, "; ")))
	}

	createInput := CreateRestApiInput{
		Name:        doc.Name,
		Description: doc.Description,
		Version:     doc.Version,
	}
	if params.EndpointConfigurationTypes != "" {
		createInput.EndpointTypes = []string{params.EndpointConfigurationTypes}
	}
	created, err := s.createRestApiCore(stores, createInput)
	if err != nil {
		return nil, nil, err
	}

	def := &apigateway.ImportedDefinition{
		Resources:   applyBasePathMode(doc, params.BasePath),
		Models:      doc.Models,
		Authorizers: doc.Authorizers,
	}
	if err := stores.restApis.ImportDefinition(created.Id, def, false); err != nil {
		_ = stores.restApis.Delete(created.Id)
		return nil, nil, toApiGatewayError(err)
	}

	created, err = stores.restApis.Get(created.Id)
	if err != nil {
		return nil, nil, toApiGatewayError(err)
	}
	created.Warnings = doc.Warnings
	if err := stores.restApis.Update(created); err != nil {
		return nil, nil, toApiGatewayError(err)
	}
	return created, doc.Warnings, nil
}

// putRestApiCore updates an existing API from an external API definition
// document. Merge (the default) installs the document's paths, models and
// authorizers on top of the existing API; overwrite replaces them
// wholesale. The API's identity, stages and deployments survive both
// modes. With failOnWarnings set, any parsing warning rejects the update
// before the API is touched.
func (s *APIGatewayService) putRestApiCore(
	stores *apiGatewayStores,
	apiId, mode string,
	body []byte,
	params importParameters,
	failOnWarnings bool,
) (*apigateway.RestApi, []string, error) {
	if apiId == "" {
		return nil, nil, NewBadRequestException("restApiId is required")
	}
	if mode != "" && mode != "merge" && mode != "overwrite" {
		return nil, nil, NewBadRequestException("Invalid mode: " + mode + "; must be merge or overwrite")
	}
	if _, err := stores.restApis.Get(apiId); err != nil {
		return nil, nil, toApiGatewayError(err)
	}

	doc, err := s.parseApiDefinition(body)
	if err != nil {
		return nil, nil, err
	}
	if failOnWarnings && len(doc.Warnings) > 0 {
		return nil, nil, NewBadRequestException(
			fmt.Sprintf("API import failed: %d warning(s): %s", len(doc.Warnings), strings.Join(doc.Warnings, "; ")))
	}

	def := &apigateway.ImportedDefinition{
		Resources:   applyBasePathMode(doc, params.BasePath),
		Models:      doc.Models,
		Authorizers: doc.Authorizers,
	}
	if err := stores.restApis.ImportDefinition(apiId, def, mode == "overwrite"); err != nil {
		return nil, nil, toApiGatewayError(err)
	}

	updated, err := stores.restApis.Get(apiId)
	if err != nil {
		return nil, nil, toApiGatewayError(err)
	}
	updated.Warnings = doc.Warnings
	if err := stores.restApis.Update(updated); err != nil {
		return nil, nil, toApiGatewayError(err)
	}
	return updated, doc.Warnings, nil
}

// parseApiDefinition decodes an OpenAPI (3.0) or Swagger (2.0) document,
// JSON or YAML, into the store-layer definition shape. Unresolvable
// constructs are recorded as warnings rather than errors, matching the
// import contract where failonwarnings turns them fatal.
func (s *APIGatewayService) parseApiDefinition(body []byte) (*apiDefinitionDocument, error) {
	if len(body) == 0 {
		return nil, NewBadRequestException("body is required")
	}
	if len(body) > apigateway.ImportDefinitionMaxBytes {
		return nil, NewBadRequestException("Invalid API definition: the maximum size is 6MB")
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		if err := yaml.Unmarshal(body, &raw); err != nil {
			return nil, NewBadRequestException("Invalid API definition: body is neither JSON nor YAML")
		}
	}

	family := ""
	if _, ok := raw["swagger"]; ok {
		family = "swagger"
	} else if _, ok := raw["openapi"]; ok {
		family = "openapi"
	} else {
		return nil, NewBadRequestException("Invalid API definition: neither a swagger nor an openapi document")
	}

	doc := &apiDefinitionDocument{
		Paths:         make(map[string][]*apigateway.ImportedMethod),
		Models:        make(map[string]*apigateway.Model),
		Authorizers:   make(map[string]*apigateway.Authorizer),
		ApiKeySchemes: make(map[string]bool),
	}

	info, _ := raw["info"].(map[string]interface{})
	if info != nil {
		doc.Name, _ = info["title"].(string)
		doc.Version, _ = info["version"].(string)
		doc.Description, _ = info["description"].(string)
	}
	if doc.Name == "" {
		return nil, NewBadRequestException("Invalid API definition: info.title is required")
	}

	doc.BasePath = importBasePath(family, raw)
	if err := parseImportSecuritySchemes(family, raw, doc); err != nil {
		return nil, err
	}
	parseImportModels(family, raw, doc)
	parseImportPaths(raw, doc)

	return doc, nil
}

// importBasePath derives the base path of the document: the basePath
// property in Swagger 2.0, or the server URL path (the single basePath
// server variable when one is declared) in OpenAPI 3.0.
func importBasePath(family string, raw map[string]interface{}) string {
	if family == "swagger" {
		basePath, _ := raw["basePath"].(string)
		return normaliseImportBasePath(basePath)
	}
	servers, _ := raw["servers"].([]interface{})
	for _, entry := range servers {
		server, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}
		vars, _ := server["variables"].(map[string]interface{})
		if basePathVar, ok := vars["basePath"].(map[string]interface{}); ok {
			def, _ := basePathVar["default"].(string)
			return normaliseImportBasePath(def)
		}
		url, _ := server["url"].(string)
		path := url
		if schemeIdx := strings.Index(url, "://"); schemeIdx >= 0 {
			rest := url[schemeIdx+3:]
			if slash := strings.Index(rest, "/"); slash >= 0 {
				path = rest[slash:]
			} else {
				path = ""
			}
		}
		return normaliseImportBasePath(path)
	}
	return ""
}

// normaliseImportBasePath trims a declared base path to its pure path
// form: no trailing slash, empty or "/" both meaning no base path.
func normaliseImportBasePath(basePath string) string {
	basePath = strings.TrimSuffix(strings.TrimSpace(basePath), "/")
	if basePath == "" {
		return ""
	}
	if !strings.HasPrefix(basePath, "/") {
		basePath = "/" + basePath
	}
	return basePath
}

// applyBasePathMode maps the document's declared paths to the resource
// paths actually installed, per the basepath parameter: ignore (the
// default) keeps the declared paths, prepend prefixes every path with the
// base path, and split drops the base path's first segment before
// prefixing with the remainder.
func applyBasePathMode(doc *apiDefinitionDocument, mode string) map[string][]*apigateway.ImportedMethod {
	prefix := ""
	switch mode {
	case "prepend":
		prefix = doc.BasePath
	case "split":
		rest := strings.TrimPrefix(doc.BasePath, "/")
		if idx := strings.Index(rest, "/"); idx >= 0 {
			prefix = rest[idx:]
		}
	}
	if prefix == "" {
		return doc.Paths
	}

	remapped := make(map[string][]*apigateway.ImportedMethod, len(doc.Paths))
	for path, methods := range doc.Paths {
		remapped[prefix+path] = methods
	}
	return remapped
}

// parseImportSecuritySchemes reads the document's security definitions:
// securityDefinitions in Swagger 2.0, components.securitySchemes in
// OpenAPI 3.0. A scheme carrying the x-amazon-apigateway-authorizer
// extension — where AWS defines it, on the security definition or scheme
// itself — becomes an authorizer of that name; a plain apiKey scheme is
// recorded for the security-requirement resolution. The authorizer
// extension rides on an API-key scheme: a carrier whose type is present
// and is not apiKey is invalid for API Gateway, so the import rejects
// the document instead of building an authorizer AWS would refuse.
func parseImportSecuritySchemes(family string, raw map[string]interface{}, doc *apiDefinitionDocument) error {
	schemes, _ := securitySchemesOf(family, raw)
	for name, entry := range schemes {
		scheme, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}
		if config, ok := scheme["x-amazon-apigateway-authorizer"].(map[string]interface{}); ok {
			if t, _ := scheme["type"].(string); t != "" && t != "apiKey" {
				return NewBadRequestException(
					"Invalid security scheme " + name + ": the authorizer extension requires the apiKey scheme type")
			}
			doc.Authorizers[name] = importedAuthorizer(name, config)
			continue
		}
		if t, _ := scheme["type"].(string); t == "apiKey" {
			doc.ApiKeySchemes[name] = true
		}
	}
	return nil
}

// importedAuthorizer builds the stored authorizer from a security scheme's
// x-amazon-apigateway-authorizer extension.
func importedAuthorizer(name string, config map[string]interface{}) *apigateway.Authorizer {
	authorizer := &apigateway.Authorizer{Name: name}
	switch t, _ := config["type"].(string); t {
	case "token":
		authorizer.Type = "TOKEN"
	case "request":
		authorizer.Type = "REQUEST"
	case "cognito_user_pools":
		authorizer.Type = "COGNITO_USER_POOLS"
	}
	authorizer.AuthorizerUri, _ = config["authorizerUri"].(string)
	authorizer.AuthorizerCredentials, _ = config["authorizerCredentials"].(string)
	authorizer.IdentitySource, _ = config["identitySource"].(string)
	authorizer.IdentityValidationExpression, _ = config["identityValidationExpression"].(string)
	authorizer.AuthorizerResultTtlInSeconds = int32(importNumber(config["authorizerResultTtlInSeconds"]))
	// providerARNs is documented as an array of user pool ARN strings.
	if providerArns, ok := config["providerARNs"].([]interface{}); ok {
		for _, v := range providerArns {
			if arn, ok := v.(string); ok {
				authorizer.ProviderArns = append(authorizer.ProviderArns, arn)
			}
		}
	}
	return authorizer
}

func securitySchemesOf(family string, raw map[string]interface{}) (map[string]interface{}, bool) {
	if family == "swagger" {
		schemes, ok := raw["securityDefinitions"].(map[string]interface{})
		return schemes, ok
	}
	components, _ := raw["components"].(map[string]interface{})
	schemes, ok := components["securitySchemes"].(map[string]interface{})
	return schemes, ok
}

// importNumber reads a numeric document value as an int64 regardless of
// whether the decoder produced an int (YAML) or a float64 (JSON).
func importNumber(v interface{}) int64 {
	switch n := v.(type) {
	case int:
		return int64(n)
	case int64:
		return n
	case float64:
		return int64(n)
	}
	return 0
}

// parseImportModels reads the document's schemas: definitions in Swagger
// 2.0, components.schemas in OpenAPI 3.0. Each schema becomes a model
// whose content is the schema serialised back to canonical JSON.
func parseImportModels(family string, raw map[string]interface{}, doc *apiDefinitionDocument) {
	var schemas map[string]interface{}
	if family == "swagger" {
		schemas, _ = raw["definitions"].(map[string]interface{})
	} else {
		components, _ := raw["components"].(map[string]interface{})
		schemas, _ = components["schemas"].(map[string]interface{})
	}
	for name, schema := range schemas {
		encoded, err := json.Marshal(schema)
		if err != nil {
			continue
		}
		doc.Models[name] = &apigateway.Model{
			Name:        name,
			Schema:      string(encoded),
			ContentType: "application/json",
		}
	}
}

// importOperationVerbs is the set of OpenAPI operation keys, mapped to the
// uppercase HTTP verbs the platform stores.
var importOperationVerbs = map[string]string{
	"get":     "GET",
	"put":     "PUT",
	"post":    "POST",
	"delete":  "DELETE",
	"options": "OPTIONS",
	"head":    "HEAD",
	"patch":   "PATCH",
}

// parseImportPaths walks the document's paths and builds the per-verb
// method payloads: operationId, the API key requirement, the referenced
// authorizer, the method responses and the integration extension.
// Document-level security is deliberately not consulted: API Gateway
// ignores the OpenAPI root security declaration on import, so a method
// is protected only by its own operation-level security.
func parseImportPaths(raw map[string]interface{}, doc *apiDefinitionDocument) {
	paths, _ := raw["paths"].(map[string]interface{})
	for path, entry := range paths {
		item, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		for key, op := range item {
			verb, ok := importOperationVerbs[key]
			if !ok {
				continue
			}
			operation, ok := op.(map[string]interface{})
			if !ok {
				continue
			}
			doc.Paths[path] = append(doc.Paths[path], parseImportOperation(verb, operation, doc))
		}
	}
}

func parseImportOperation(verb string, operation map[string]interface{}, doc *apiDefinitionDocument) *apigateway.ImportedMethod {
	method := &apigateway.Method{}
	imported := &apigateway.ImportedMethod{Verb: verb, Method: method}

	if operationId, _ := operation["operationId"].(string); operationId != "" {
		method.OperationName = operationId
	}
	if required, _ := operation["x-amazon-apigateway-api-key-required"].(bool); required {
		method.ApiKeyRequired = true
	}
	method.AuthorizationType = "NONE"
	requirements, _ := operation["security"].([]interface{})
	parseImportSecurity(requirements, doc, imported)
	method.MethodResponses = parseImportResponses(operation)
	if integration, ok := operation["x-amazon-apigateway-integration"].(map[string]interface{}); ok {
		method.MethodIntegration = parseImportIntegration(integration)
	}
	return imported
}

// parseImportSecurity resolves an operation's own security requirements
// against the declared authorizers and API-key security schemes; the
// document-level declaration never reaches this resolution. A
// requirement naming a declared authorizer binds the method to it; one
// naming an API-key scheme marks the method as requiring an API key;
// anything else is recorded as a warning.
func parseImportSecurity(requirements []interface{}, doc *apiDefinitionDocument, imported *apigateway.ImportedMethod) {
	for _, entry := range requirements {
		requirement, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}
		for name := range requirement {
			if authorizer, ok := doc.Authorizers[name]; ok {
				imported.AuthorizerName = authorizer.Name
				imported.Method.AuthorizationType = authorizer.Type
				continue
			}
			if doc.ApiKeySchemes[name] {
				imported.Method.ApiKeyRequired = true
				continue
			}
			doc.Warnings = append(doc.Warnings,
				"security scheme "+name+" is neither a declared authorizer nor an API key scheme; ignored")
		}
	}
}

// parseImportResponses reads an operation's responses object into method
// response entries keyed by status code.
func parseImportResponses(operation map[string]interface{}) map[string]*apigateway.MethodResponse {
	responses, _ := operation["responses"].(map[string]interface{})
	if len(responses) == 0 {
		return nil
	}
	methodResponses := make(map[string]*apigateway.MethodResponse, len(responses))
	for code := range responses {
		if !statusCodePattern.MatchString(code) {
			continue
		}
		methodResponses[code] = &apigateway.MethodResponse{StatusCode: code}
	}
	return methodResponses
}

// parseImportIntegration reads the x-amazon-apigateway-integration
// extension into the stored integration shape.
func parseImportIntegration(integration map[string]interface{}) *apigateway.Integration {
	parsed := &apigateway.Integration{}
	parsed.Type, _ = integration["type"].(string)
	parsed.IntegrationHttpMethod, _ = integration["httpMethod"].(string)
	parsed.Uri, _ = integration["uri"].(string)
	parsed.Credentials, _ = integration["credentials"].(string)
	parsed.PassthroughBehavior, _ = integration["passthroughBehavior"].(string)
	if templates, ok := integration["requestTemplates"].(map[string]interface{}); ok {
		parsed.RequestTemplates = make(map[string]string, len(templates))
		for contentType, template := range templates {
			if text, ok := template.(string); ok {
				parsed.RequestTemplates[contentType] = text
			}
		}
	}
	if params, ok := integration["requestParameters"].(map[string]interface{}); ok {
		parsed.RequestParameters = make(map[string]string, len(params))
		for key, value := range params {
			if text, ok := value.(string); ok {
				parsed.RequestParameters[key] = text
			}
		}
	}
	return parsed
}

// importApiKeysCore imports API keys from a CSV document. The first record
// holds the column names (case-insensitive, any order); the Name and Key
// columns are required, Description, Enabled and UsagePlanIds are
// optional, and unknown columns are ignored. Rows failing validation are
// skipped with a warning; with failOnWarnings set, any warning — including
// a plan association that already exists — rejects the import before
// anything is persisted. Re-importing a key value overwrites the stored
// key.
func (s *APIGatewayService) importApiKeysCore(
	stores *apiGatewayStores,
	body []byte,
	format string,
	failOnWarnings bool,
) ([]string, []string, error) {
	if format == "" {
		return nil, nil, NewBadRequestException("format is required")
	}
	if format != "csv" {
		return nil, nil, NewBadRequestException("Invalid format: " + format + "; only csv is supported")
	}

	rows, err := parseApiKeyCsv(body)
	if err != nil {
		return nil, nil, err
	}

	// Phase one validates every row and gathers warnings without writing,
	// so failonwarnings can reject the import atomically.
	type preparedKey struct {
		key          *apigateway.ApiKey
		usagePlanIds []string
	}
	// A plan association that already exists — in the store, or earlier in
	// this file, since the upsert makes the first row's association real
	// for the second — is reported as a warning.
	type planAssociation struct{ value, planId string }
	seenAssociations := make(map[planAssociation]bool)
	var warnings []string
	prepared := make([]preparedKey, 0, len(rows))
	for _, row := range rows {
		if len(row.value) < apigateway.ApiKeyValueMinLength || len(row.value) > apigateway.ApiKeyValueMaxLength {
			warnings = append(warnings, fmt.Sprintf(
				"Invalid key value for %q: must be between %d and %d characters", row.name,
				apigateway.ApiKeyValueMinLength, apigateway.ApiKeyValueMaxLength))
			continue
		}
		if len(row.name) > apigateway.ApiKeyNameMaxLength {
			warnings = append(warnings, fmt.Sprintf(
				"Invalid key name: must not exceed %d characters", apigateway.ApiKeyNameMaxLength))
			continue
		}
		// Usage-plan-key records are keyed by API key identifier, so the
		// stored check resolves the row's value to its key id first. A
		// lookup that fails for a reason other than absence fails the
		// import: swallowing it would silently drop the duplicate
		// warning.
		var keyId string
		existing, err := stores.usage.GetApiKeyByValue(row.value)
		if err == nil {
			keyId = existing.Id
		} else if !errors.Is(err, apigateway.ErrApiKeyNotFound) {
			return nil, nil, err
		}
		for _, planId := range row.usagePlanIds {
			associated := seenAssociations[planAssociation{row.value, planId}]
			if !associated && keyId != "" {
				if _, err := stores.usage.GetUsagePlanKey(planId, keyId); err == nil {
					associated = true
				}
			}
			if associated {
				warnings = append(warnings, "API key "+row.value+" is already associated with usage plan "+planId)
			}
			seenAssociations[planAssociation{row.value, planId}] = true
		}
		prepared = append(prepared, preparedKey{
			key: &apigateway.ApiKey{
				Value:       row.value,
				Name:        row.name,
				Description: row.description,
				Enabled:     row.enabled,
			},
			usagePlanIds: row.usagePlanIds,
		})
	}
	if failOnWarnings && len(warnings) > 0 {
		return nil, nil, NewBadRequestException(
			fmt.Sprintf("API key import failed: %d warning(s): %s", len(warnings), strings.Join(warnings, "; ")))
	}

	// Phase two persists; nothing here adds warnings, so the import is
	// all-or-nothing with respect to the failonwarnings contract.
	var ids []string
	for _, entry := range prepared {
		if err := s.upsertImportedApiKey(stores, entry.key, entry.usagePlanIds); err != nil {
			return nil, nil, err
		}
		ids = append(ids, entry.key.Id)
	}
	return ids, warnings, nil
}

// upsertImportedApiKey persists one imported key: an existing key with the
// same value keeps its identifier, creation date and plan memberships and
// has its descriptive fields replaced. A lookup failure other than
// absence aborts the import — falling through to the create branch would
// mint a second record for the same value. Plan associations that
// already exist are skipped — the caller reports them as warnings up
// front.
func (s *APIGatewayService) upsertImportedApiKey(stores *apiGatewayStores, key *apigateway.ApiKey, usagePlanIds []string) error {
	existing, err := stores.usage.GetApiKeyByValue(key.Value)
	if err != nil && !errors.Is(err, apigateway.ErrApiKeyNotFound) {
		return err
	}
	if err == nil {
		key.Id = existing.Id
		key.CreatedDate = existing.CreatedDate
		key.StageKeys = existing.StageKeys
		key.Tags = existing.Tags
		if err := stores.usage.UpdateApiKey(key); err != nil {
			return err
		}
	} else {
		generateDistinctId := false
		created, err := s.createApiKeyCore(stores, &ApiKeyInput{
			Name:               key.Name,
			Description:        key.Description,
			Enabled:            &key.Enabled,
			Value:              key.Value,
			GenerateDistinctId: &generateDistinctId,
		})
		if err != nil {
			return err
		}
		key.Id = created.Id
	}

	for _, planId := range usagePlanIds {
		if _, err := stores.usage.CreateUsagePlanKey(planId, &apigateway.UsagePlanKey{Id: key.Id}); err != nil {
			if errors.Is(err, apigateway.ErrUsagePlanKeyAlreadyExists) {
				continue
			}
			return err
		}
	}
	return nil
}

// importedApiKeyRow is one parsed CSV data row.
type importedApiKeyRow struct {
	name         string
	value        string
	description  string
	enabled      bool
	usagePlanIds []string
}

// parseApiKeyCsv parses the CSV key file: a header record naming the
// columns (case-insensitive; Name and Key required) followed by one record
// per key. The UsagePlanIds column holds a single cell whose value is a
// comma-separated list of plan identifiers. An absent or empty Enabled
// column leaves the key enabled.
func parseApiKeyCsv(body []byte) ([]importedApiKeyRow, error) {
	reader := csv.NewReader(bytes.NewReader(body))
	reader.TrimLeadingSpace = true
	records, err := reader.ReadAll()
	if err != nil {
		return nil, NewBadRequestException("Invalid API key file: " + err.Error())
	}
	if len(records) == 0 {
		return nil, NewBadRequestException("Invalid API key file: the header record is missing")
	}

	columns := make(map[string]int)
	for idx, name := range records[0] {
		columns[strings.ToLower(strings.TrimSpace(name))] = idx
	}
	if _, ok := columns["name"]; !ok {
		return nil, NewBadRequestException("Invalid API key file: the Name column is required")
	}
	if _, ok := columns["key"]; !ok {
		return nil, NewBadRequestException("Invalid API key file: the Key column is required")
	}

	cell := func(record []string, column string) string {
		idx, ok := columns[column]
		if !ok || idx >= len(record) {
			return ""
		}
		return strings.TrimSpace(record[idx])
	}

	rows := make([]importedApiKeyRow, 0, len(records)-1)
	for _, record := range records[1:] {
		if len(record) == 0 {
			continue
		}
		row := importedApiKeyRow{
			name:        cell(record, "name"),
			value:       cell(record, "key"),
			description: cell(record, "description"),
			enabled:     true,
		}
		if strings.EqualFold(cell(record, "enabled"), "false") {
			row.enabled = false
		}
		if plans := cell(record, "usageplanids"); plans != "" {
			for _, planId := range strings.Split(plans, ",") {
				if planId = strings.TrimSpace(planId); planId != "" {
					row.usagePlanIds = append(row.usagePlanIds, planId)
				}
			}
		}
		rows = append(rows, row)
	}
	return rows, nil
}
