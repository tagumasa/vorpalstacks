package apigateway

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"vorpalstacks/internal/store/aws/apigateway"
)

// sdkConfigurationProperty describes one configuration property of an SDK
// type, as the GetSdk documentation states it: which parameters each
// sdkType requires.
type sdkConfigurationProperty struct {
	Name         string
	FriendlyName string
	Description  string
	Required     bool
	DefaultValue string
}

// sdkType describes a static SDK type. The table's substance comes from the
// GetSdk documentation: the supported sdkType identifiers and, per type,
// the required generation parameters.
type sdkType struct {
	Id           string
	FriendlyName string
	Description  string
	Properties   []sdkConfigurationProperty
}

var sdkTypeTable = []*sdkType{
	{
		Id:           "android",
		FriendlyName: "Android",
		Description:  "SDK type for Android clients",
		Properties: []sdkConfigurationProperty{
			{Name: "groupId", FriendlyName: "Group ID", Description: "Maven group ID of the generated SDK", Required: true},
			{Name: "artifactId", FriendlyName: "Artifact ID", Description: "Maven artifact ID of the generated SDK", Required: true},
			{Name: "artifactVersion", FriendlyName: "Artifact Version", Description: "Maven artifact version of the generated SDK", Required: true},
			{Name: "invokerPackage", FriendlyName: "Invoker Package", Description: "Java package name of the generated SDK classes", Required: true},
		},
	},
	{
		Id:           "java",
		FriendlyName: "Java",
		Description:  "SDK type for Java clients",
		Properties: []sdkConfigurationProperty{
			{Name: "serviceName", FriendlyName: "Service Name", Description: "Service name of the generated SDK", Required: true},
			{Name: "javaPackageName", FriendlyName: "Java Package Name", Description: "Java package name of the generated SDK classes", Required: true},
		},
	},
	{
		Id:           "javascript",
		FriendlyName: "JavaScript",
		Description:  "SDK type for JavaScript clients",
	},
	{
		Id:           "objectivec",
		FriendlyName: "Objective-C",
		Description:  "SDK type for iOS clients in Objective-C",
		Properties: []sdkConfigurationProperty{
			{Name: "classPrefix", FriendlyName: "Class Prefix", Description: "Class prefix of the generated SDK types", Required: true},
		},
	},
	{
		Id:           "swift",
		FriendlyName: "Swift",
		Description:  "SDK type for iOS clients in Swift",
		Properties: []sdkConfigurationProperty{
			{Name: "classPrefix", FriendlyName: "Class Prefix", Description: "Class prefix of the generated SDK types", Required: true},
		},
	},
	{
		Id:           "ruby",
		FriendlyName: "Ruby",
		Description:  "SDK type for Ruby clients",
	},
}

// getSdkTypeCore resolves a static SDK type by identifier.
func (s *APIGatewayService) getSdkTypeCore(id string) (*sdkType, *ApiGatewayError) {
	if id == "" {
		return nil, NewBadRequestException("id is required")
	}
	for _, t := range sdkTypeTable {
		if t.Id == id {
			return t, nil
		}
	}
	return nil, NewNotFoundException("SDK type", id)
}

// listSdkTypesCore returns the static SDK type table, paginated.
func (s *APIGatewayService) listSdkTypesCore(limit int, position string) ([]*sdkType, string, error) {
	resolved, err := resolvePageLimit(limit)
	if err != nil {
		return nil, "", err
	}
	start := 0
	if position != "" {
		found := false
		for i, t := range sdkTypeTable {
			if t.Id == position {
				start = i + 1
				found = true
				break
			}
		}
		if !found {
			return nil, "", NewBadRequestException("Invalid position: " + position)
		}
	}
	end := len(sdkTypeTable)
	if start+resolved < end {
		end = start + resolved
	}
	page := sdkTypeTable[start:end]
	next := ""
	if end < len(sdkTypeTable) && len(page) > 0 {
		next = page[len(page)-1].Id
	}
	return page, next, nil
}

// getExportCore renders the API's export document for a stage. The export
// reflects the state the stage's deployment snapshotted. exportType
// selects the document family (swagger for Swagger 2.0, oas30 for
// OpenAPI 3.0); accepts selects application/json (default) or
// application/yaml.
func (s *APIGatewayService) getExportCore(
	stores *apiGatewayStores,
	apiId, stageName, exportType, accepts string,
) ([]byte, string, error) {
	if apiId == "" {
		return nil, "", NewBadRequestException("restApiId is required")
	}
	if stageName == "" {
		return nil, "", NewBadRequestException("stageName is required")
	}
	var documentFamily string
	switch exportType {
	case "swagger":
		documentFamily = "swagger"
	case "oas30":
		documentFamily = "openapi"
	default:
		return nil, "", NewBadRequestException(
			"Invalid exportType: " + exportType + "; must be swagger or oas30")
	}
	var contentType string
	mediaType := strings.ToLower(strings.TrimSpace(accepts))
	if idx := strings.Index(mediaType, ";"); idx >= 0 {
		mediaType = strings.TrimSpace(mediaType[:idx])
	}
	switch mediaType {
	case "", "application/json":
		contentType = "application/json"
	case "application/yaml", "text/yaml", "application/x-yaml":
		contentType = "application/yaml"
	default:
		return nil, "", NewBadRequestException(
			"Invalid accepts header: " + accepts + "; must be application/json or application/yaml")
	}

	api, err := stores.restApis.Get(apiId)
	if err != nil {
		return nil, "", toApiGatewayError(err)
	}
	stage, err := stores.restApis.GetStage(apiId, stageName)
	if err != nil {
		return nil, "", toApiGatewayError(err)
	}
	deployment, err := stores.restApis.GetDeployment(apiId, stage.DeploymentId)
	if err != nil {
		return nil, "", toApiGatewayError(err)
	}

	doc := buildExportDocument(api, deployment, documentFamily)
	if contentType == "application/yaml" {
		out, err := yaml.Marshal(doc)
		if err != nil {
			return nil, "", fmt.Errorf("render export document: %w", err)
		}
		return out, contentType, nil
	}
	out, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, "", fmt.Errorf("render export document: %w", err)
	}
	return out, contentType, nil
}

// buildExportDocument renders the export view of a deployment snapshot.
func buildExportDocument(api *apigateway.RestApi, deployment *apigateway.Deployment, family string) map[string]interface{} {
	version := api.Version
	if version == "" {
		version = deployment.CreatedDate.Format("2006-01-02")
	}
	info := map[string]interface{}{"title": api.Name, "version": version}
	if api.Description != "" {
		info["description"] = api.Description
	}

	doc := map[string]interface{}{
		family:  familySwaggerVersion(family),
		"info":  info,
		"paths": exportPaths(deployment),
	}
	if definitions := exportDefinitions(deployment); len(definitions) > 0 {
		doc["definitions"] = definitions
	}
	return doc
}

func familySwaggerVersion(family string) string {
	if family == "openapi" {
		return "3.0.1"
	}
	return "2.0"
}

// exportPaths walks the snapshot's resources in path order and renders each
// method with the settings the snapshot holds, including the API Gateway
// integration extension.
func exportPaths(deployment *apigateway.Deployment) map[string]interface{} {
	paths := make(map[string]interface{})
	if deployment.Snapshot == nil {
		return paths
	}
	resources := make([]*apigateway.Resource, 0, len(deployment.Snapshot.Resources))
	for _, resource := range deployment.Snapshot.Resources {
		resources = append(resources, resource)
	}
	sort.Slice(resources, func(i, j int) bool { return resources[i].Path < resources[j].Path })

	for _, resource := range resources {
		if len(resource.ResourceMethods) == 0 {
			continue
		}
		verbs := make([]string, 0, len(resource.ResourceMethods))
		for verb := range resource.ResourceMethods {
			verbs = append(verbs, verb)
		}
		sort.Strings(verbs)
		operations := make(map[string]interface{}, len(verbs))
		for _, verb := range verbs {
			operations[strings.ToLower(verb)] = exportOperation(resource.ResourceMethods[verb])
		}
		paths[resource.Path] = operations
	}
	return paths
}

func exportOperation(method *apigateway.Method) map[string]interface{} {
	op := map[string]interface{}{}
	if method.OperationName != "" {
		op["operationId"] = method.OperationName
	}
	if method.ApiKeyRequired {
		op["x-amazon-apigateway-api-key-required"] = true
	}
	responses := map[string]interface{}{}
	for code := range method.MethodResponses {
		responses[code] = map[string]interface{}{"description": "Success description"}
	}
	if len(responses) == 0 {
		responses["200"] = map[string]interface{}{"description": "Success description"}
	}
	op["responses"] = responses
	if method.MethodIntegration != nil {
		integration := map[string]interface{}{
			"type": method.MethodIntegration.Type,
		}
		if method.MethodIntegration.IntegrationHttpMethod != "" {
			integration["httpMethod"] = method.MethodIntegration.IntegrationHttpMethod
		}
		if method.MethodIntegration.Uri != "" {
			integration["uri"] = method.MethodIntegration.Uri
		}
		op["x-amazon-apigateway-integration"] = integration
	}
	return op
}

// exportDefinitions inlines the snapshot's model schemas as the document's
// definitions.
func exportDefinitions(deployment *apigateway.Deployment) map[string]interface{} {
	definitions := make(map[string]interface{})
	if deployment.Snapshot == nil {
		return definitions
	}
	for name, model := range deployment.Snapshot.Models {
		var schema interface{}
		if err := json.Unmarshal([]byte(model.Schema), &schema); err != nil {
			continue
		}
		definitions[name] = schema
	}
	return definitions
}

// getModelTemplateCore generates the sample mapping template of a model:
// a VTL skeleton that reads the payload root and echoes one quoted entry
// per object property.
func (s *APIGatewayService) getModelTemplateCore(stores *apiGatewayStores, apiId, modelName string) (string, error) {
	model, err := s.getModelCore(stores, apiId, modelName, false)
	if err != nil {
		return "", err
	}

	var schema struct {
		Properties map[string]interface{} `json:"properties"`
	}
	_ = json.Unmarshal([]byte(model.Schema), &schema)

	var b strings.Builder
	b.WriteString("#set($inputRoot = $input.path('$'))\n{")
	names := make([]string, 0, len(schema.Properties))
	for name := range schema.Properties {
		names = append(names, name)
	}
	sort.Strings(names)
	// The comma separates properties; the last one must not carry it, or
	// every render of the template ends in invalid JSON.
	for i, name := range names {
		if i > 0 {
			b.WriteString(",")
		}
		fmt.Fprintf(&b, "\n  \"%s\": \"$inputRoot.%s\"", name, name)
	}
	b.WriteString("\n}")
	return b.String(), nil
}
