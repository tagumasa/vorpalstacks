package apigateway

import (
	"sort"
	"strings"
)

// ImportedMethod pairs a parsed operation with the name of the authorizer
// it references, if any. The document identifies authorizers by name; the
// store resolves the name to the identifier of the authorizer installed
// from the same definition.
type ImportedMethod struct {
	Verb           string
	Method         *Method
	AuthorizerName string
}

// ImportedDefinition is the store-layer shape of a parsed API definition
// document: operations grouped by resource path (intermediate path
// segments become resources without methods), models by name, and
// authorizers by name.
type ImportedDefinition struct {
	Resources   map[string][]*ImportedMethod
	Models      map[string]*Model
	Authorizers map[string]*Authorizer
}

// ImportDefinition installs an imported definition into an existing API in
// one atomic write. Merge keeps resources, models and authorizers that the
// definition does not mention (a mentioned path gains or replaces only the
// operations the definition carries); overwrite drops everything the
// definition could replace first. The API identity, stages, deployments
// and gateway responses survive either mode — serving a changed definition
// always requires a new deployment.
func (s *RestApiStore) ImportDefinition(apiId string, def *ImportedDefinition, overwrite bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	api, err := s.Get(apiId)
	if err != nil {
		return err
	}
	ensureRestApiMaps(api)

	if overwrite {
		overwriteDefinition(api)
	}

	s.installAuthorizers(api, def)
	s.installModels(api, def)
	s.installResources(api, def)

	return s.updateLocked(api)
}

// overwriteDefinition empties the replaceable parts of the API: every
// resource except the root (whose methods are cleared), all models and all
// authorizers.
func overwriteDefinition(api *RestApi) {
	for id, resource := range api.Resources {
		if resource.Path == "/" {
			resource.ResourceMethods = make(map[string]*Method)
			continue
		}
		delete(api.Resources, id)
	}
	api.Models = make(map[string]*Model)
	api.Authorizers = make(map[string]*Authorizer)
}

// installAuthorizers upserts the definition's authorizers by name: an
// existing authorizer with the same name keeps its identifier and has its
// configuration replaced.
func (s *RestApiStore) installAuthorizers(api *RestApi, def *ImportedDefinition) {
	names := make([]string, 0, len(def.Authorizers))
	for name := range def.Authorizers {
		names = append(names, name)
	}
	sort.Strings(names)

	byName := make(map[string]*Authorizer)
	for _, existing := range api.Authorizers {
		byName[existing.Name] = existing
	}
	for _, name := range names {
		imported := def.Authorizers[name]
		if existing, ok := byName[name]; ok {
			imported.Id = existing.Id
		} else {
			imported.Id = s.arnBuilder.GenerateAuthorizerId()
		}
		imported.RestApiId = api.Id
		api.Authorizers[imported.Id] = imported
	}
}

// authorizerIdByName resolves an authorizer name to the identifier of the
// authorizer currently installed on the API, or empty when unknown.
func authorizerIdByName(api *RestApi, name string) string {
	for _, existing := range api.Authorizers {
		if existing.Name == name {
			return existing.Id
		}
	}
	return ""
}

// installModels upserts the definition's models by name: an existing model
// with the same name keeps its identifier and has its content replaced.
func (s *RestApiStore) installModels(api *RestApi, def *ImportedDefinition) {
	names := make([]string, 0, len(def.Models))
	for name := range def.Models {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		imported := def.Models[name]
		if existing, ok := api.Models[name]; ok {
			imported.Id = existing.Id
		} else {
			imported.Id = s.arnBuilder.GenerateModelId()
		}
		imported.RestApiId = api.Id
		api.Models[name] = imported
	}
}

// installResources installs the definition's operations. Intermediate
// path segments become resources without methods — the definition need not
// declare them; a path already present on the API keeps its identifier and
// gains or replaces exactly the verbs the definition carries.
func (s *RestApiStore) installResources(api *RestApi, def *ImportedDefinition) {
	byPath := make(map[string]*Resource, len(api.Resources))
	for _, resource := range api.Resources {
		byPath[resource.Path] = resource
	}

	paths := make([]string, 0, len(def.Resources))
	for path := range def.Resources {
		paths = append(paths, path)
	}
	// Parents precede children lexicographically for path-shaped keys, so
	// a plain sort installs every parent before its descendants.
	sort.Strings(paths)

	for _, path := range paths {
		resource := s.ensureResourcePath(api, byPath, path)
		if resource.ResourceMethods == nil {
			// An empty method map is omitted by the record's JSON
			// encoding, so a freshly decoded resource carries nil here.
			resource.ResourceMethods = make(map[string]*Method)
		}
		for _, imported := range def.Resources[path] {
			method := imported.Method
			method.RestApiId = api.Id
			method.ResourceId = resource.Id
			method.HttpMethod = imported.Verb
			if imported.AuthorizerName != "" {
				method.AuthorizerId = authorizerIdByName(api, imported.AuthorizerName)
			}
			resource.ResourceMethods[imported.Verb] = method
		}
	}
}

// ensureResourcePath returns the resource at the given path, creating the
// resource and any undeclared ancestors on the way down from the root,
// which always exists.
func (s *RestApiStore) ensureResourcePath(api *RestApi, byPath map[string]*Resource, path string) *Resource {
	if existing, ok := byPath[path]; ok {
		return existing
	}
	parent := s.ensureResourcePath(api, byPath, parentOf(path))
	resource := &Resource{
		Id:              s.arnBuilder.GenerateResourceId(),
		RestApiId:       api.Id,
		ParentId:        parent.Id,
		Path:            path,
		PathPart:        strings.TrimPrefix(path, parent.Path),
		ResourceMethods: make(map[string]*Method),
	}
	api.Resources[resource.Id] = resource
	byPath[path] = resource
	return resource
}

// parentOf returns the parent path of a resource path: the path minus its
// last non-root segment, or "/" when the parent is the root.
func parentOf(path string) string {
	trimmed := strings.TrimSuffix(path, "/")
	idx := strings.LastIndex(trimmed, "/")
	if idx <= 0 {
		return "/"
	}
	return trimmed[:idx]
}
