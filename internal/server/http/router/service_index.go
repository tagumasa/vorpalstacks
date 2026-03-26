package router

import (
	"regexp"
	"sort"
	"strings"

	"vorpalstacks/internal/store/api"
)

// OperationRef references a service operation by name.
type OperationRef struct {
	ServiceName string
	OpName      string
}

// PathPattern represents a URL path pattern for matching HTTP operations.
type PathPattern struct {
	Method     string
	Regex      *regexp.Regexp
	ParamNames []string
	ServiceID  int64
	OpName     string
}

// ServiceIndex provides fast lookup of services and operations by various request attributes.
type ServiceIndex struct {
	endpointPrefixToName map[string]string
	sdkIDToName          map[string]string
	targetToOperation    map[string]*OperationRef
	actionToOperation    map[string]*OperationRef
	pathPatterns         []*PathPattern
	serviceCount         int
	serviceIDToName      map[int64]string

	serviceStore   *api.ServiceStore
	operationStore *api.OperationStore
}

// NewServiceIndex creates a new ServiceIndex from the given service and operation stores.
func NewServiceIndex(serviceStore *api.ServiceStore, operationStore *api.OperationStore) (*ServiceIndex, error) {
	idx := &ServiceIndex{
		endpointPrefixToName: make(map[string]string),
		sdkIDToName:          make(map[string]string),
		targetToOperation:    make(map[string]*OperationRef),
		actionToOperation:    make(map[string]*OperationRef),
		pathPatterns:         make([]*PathPattern, 0),
		serviceIDToName:      make(map[int64]string),
		serviceStore:         serviceStore,
		operationStore:       operationStore,
	}

	if err := idx.build(); err != nil {
		return nil, err
	}

	return idx, nil
}

func (idx *ServiceIndex) build() error {
	if err := idx.serviceStore.ForEach(func(svc *api.Service) error {
		idx.serviceCount++
		idx.serviceIDToName[svc.ID] = svc.Name
		if svc.EndpointPrefix != "" {
			idx.endpointPrefixToName[svc.EndpointPrefix] = svc.Name
		}
		if svc.SDKID != "" {
			idx.sdkIDToName[svc.SDKID] = svc.Name
		}
		return nil
	}); err != nil {
		return err
	}

	if err := idx.operationStore.ForEach(func(op *api.Operation) error {
		svcName := idx.getServiceNameByID(op.ServiceID)
		if svcName == "" {
			return nil
		}

		if op.XAmzTarget != nil && *op.XAmzTarget != "" {
			idx.targetToOperation[*op.XAmzTarget] = &OperationRef{
				ServiceName: svcName,
				OpName:      op.Name,
			}
		}

		idx.actionToOperation[svcName+":"+op.Name] = &OperationRef{
			ServiceName: svcName,
			OpName:      op.Name,
		}

		if op.HTTPMethod != nil {
			path := ""
			if op.HTTPURIPath != nil && *op.HTTPURIPath != "" {
				path = *op.HTTPURIPath
			} else if op.HTTPURI != nil && *op.HTTPURI != "" {
				path = extractPath(*op.HTTPURI)
			}
			if path != "" {
				pattern := idx.compilePathPattern(*op.HTTPMethod, path, op.ServiceID, op.Name)
				if pattern != nil {
					idx.pathPatterns = append(idx.pathPatterns, pattern)
				}
			}
		}

		return nil
	}); err != nil {
		return err
	}

	sort.SliceStable(idx.pathPatterns, func(i, j int) bool {
		pi, pj := idx.pathPatterns[i], idx.pathPatterns[j]
		if len(pi.ParamNames) != len(pj.ParamNames) {
			return len(pi.ParamNames) < len(pj.ParamNames)
		}
		return len(pi.Regex.String()) > len(pj.Regex.String())
	})

	return nil
}

func (idx *ServiceIndex) getServiceNameByID(serviceID int64) string {
	return idx.serviceIDToName[serviceID]
}

func (idx *ServiceIndex) compilePathPattern(method, path string, serviceID int64, opName string) *PathPattern {
	regexStr := "^"
	paramNames := make([]string, 0)

	segments := strings.Split(path, "/")
	for _, seg := range segments {
		if seg == "" {
			continue
		}
		regexStr += "/"

		if strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}") {
			paramName := seg[1 : len(seg)-1]
			if strings.HasSuffix(paramName, "+") {
				paramName = paramName[:len(paramName)-1]
				regexStr += "(.+)"
			} else {
				regexStr += "([^/]+)"
			}
			paramNames = append(paramNames, paramName)
		} else {
			regexStr += regexp.QuoteMeta(seg)
		}
	}

	regexStr += "/?$"

	regex, err := regexp.Compile(regexStr)
	if err != nil {
		return nil
	}

	return &PathPattern{
		Method:     method,
		Regex:      regex,
		ParamNames: paramNames,
		ServiceID:  serviceID,
		OpName:     opName,
	}
}

// GetServiceByEndpointPrefix returns the service name for the given endpoint prefix.
func (idx *ServiceIndex) GetServiceByEndpointPrefix(prefix string) string {
	return idx.endpointPrefixToName[prefix]
}

// GetServiceBySDKID returns the service name for the given SDK ID.
func (idx *ServiceIndex) GetServiceBySDKID(sdkID string) string {
	return idx.sdkIDToName[sdkID]
}

// GetOperationByTarget returns the operation reference for the given X-Amz-Target header value.
func (idx *ServiceIndex) GetOperationByTarget(target string) *OperationRef {
	return idx.targetToOperation[target]
}

// GetOperationByAction returns the operation reference for the given service name and action.
func (idx *ServiceIndex) GetOperationByAction(serviceName, action string) *OperationRef {
	return idx.actionToOperation[serviceName+":"+action]
}

// MatchOperationByPath matches an operation by HTTP method and request path.
func (idx *ServiceIndex) MatchOperationByPath(method, path string) *OperationRef {
	for _, pattern := range idx.pathPatterns {
		if pattern.Method == method && pattern.Regex.MatchString(path) {
			return &OperationRef{
				ServiceName: idx.getServiceNameByID(pattern.ServiceID),
				OpName:      pattern.OpName,
			}
		}
	}
	return nil
}

// GetServiceCount returns the number of services in the index.
func (idx *ServiceIndex) GetServiceCount() int {
	return idx.serviceCount
}

// GetOperationCount returns the total number of operations in the index.
func (idx *ServiceIndex) GetOperationCount() int {
	return len(idx.actionToOperation)
}

// GetPathPatternCount returns the number of path patterns in the index.
func (idx *ServiceIndex) GetPathPatternCount() int {
	return len(idx.pathPatterns)
}

func extractPath(uri string) string {
	for i, c := range uri {
		if c == '?' {
			return uri[:i]
		}
	}
	return uri
}
