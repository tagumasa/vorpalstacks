package cloudfront

import (
	"context"
	"fmt"
	"sync"
	"time"

	"vorpalstacks/internal/services/aws/common/errors"
	"vorpalstacks/internal/services/aws/common/protocol"
	"vorpalstacks/internal/services/aws/common/request"
)

var (
	invalidationMu          sync.Mutex
	invalidationSeq         int64
	invalidations           = map[string][]map[string]interface{}{}
	maxInvalidationsPerDist = 1000
)

func cleanupInvalidations(distID string) {
	if items, ok := invalidations[distID]; ok && len(items) > maxInvalidationsPerDist {
		invalidations[distID] = items[len(items)-maxInvalidationsPerDist:]
	}
}

func generateInvalidationID() string {
	invalidationMu.Lock()
	invalidationSeq++
	id := invalidationSeq
	invalidationMu.Unlock()
	return fmt.Sprintf("INV%d-%s", id, time.Now().UTC().Format("20060102150405"))
}

// CreateInvalidation creates a new cache invalidation for a CloudFront distribution.
func (s *CloudFrontService) CreateInvalidation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	distID := request.GetStringParam(req.Parameters, "Id")
	if distID == "" {
		return nil, errors.NewAWSError("InvalidArgument", "Required parameter Id is missing.", 400)
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := stores.distributions.Get(distID); err != nil {
		return nil, errors.NewAWSError("NoSuchDistribution", fmt.Sprintf("The specified distribution does not exist: %s", distID), 404)
	}

	batch := request.GetMapParam(req.Parameters, "InvalidationBatch")
	if batch == nil {
		batch = req.Parameters
	}

	callerRef, _ := batch["CallerReference"].(string)
	if callerRef == "" {
		return nil, errors.NewAWSError("InvalidArgument", "CallerReference is required.", 400)
	}

	var paths []string
	if pathsMap, ok := batch["Paths"].(map[string]interface{}); ok {
		if items, ok := pathsMap["Items"]; ok {
			switch iv := items.(type) {
			case []interface{}:
				for _, item := range iv {
					if p, ok := item.(string); ok {
						paths = append(paths, p)
					}
				}
			case map[string]interface{}:
				if pathItems, ok := iv["Path"]; ok {
					switch pv := pathItems.(type) {
					case []interface{}:
						for _, item := range pv {
							if p, ok := item.(string); ok {
								paths = append(paths, p)
							}
						}
					case string:
						paths = append(paths, pv)
					}
				}
			}
		}
		quantity := int(request.GetIntParam(pathsMap, "Quantity"))
		if quantity > 0 && len(paths) != quantity {
			paths = paths[:quantity]
		}
	}

	if len(paths) == 0 {
		return nil, errors.NewAWSError("InvalidArgument", "At least one path is required.", 400)
	}
	if len(paths) > 3000 {
		return nil, errors.NewAWSError("InvalidArgument", "Cannot invalidate more than 3000 paths in a single request.", 400)
	}

	invalidationID := generateInvalidationID()
	now := time.Now().UTC()

	pathItemsXML := protocol.XMLElements{ElementName: "Path", Items: make([]interface{}, len(paths))}
	for i, p := range paths {
		pathItemsXML.Items[i] = p
	}

	invalidation := map[string]interface{}{
		"Id":         invalidationID,
		"CreateTime": now.Format(time.RFC3339),
		"Status":     "COMPLETED",
		"InvalidationBatch": map[string]interface{}{
			"CallerReference": callerRef,
			"Paths": map[string]interface{}{
				"Quantity": len(paths),
				"Items":    pathItemsXML,
			},
		},
	}

	invalidationMu.Lock()
	invalidations[distID] = append(invalidations[distID], invalidation)
	cleanupInvalidations(distID)
	invalidationMu.Unlock()

	return map[string]interface{}{
		"Location": fmt.Sprintf("/2020-05-31/distribution/%s/invalidation/%s", distID, invalidationID),
		"Invalidation": map[string]interface{}{
			"Id":         invalidationID,
			"CreateTime": now.Format(time.RFC3339),
			"Status":     "COMPLETED",
			"InvalidationBatch": map[string]interface{}{
				"CallerReference": callerRef,
				"Paths": map[string]interface{}{
					"Quantity": len(paths),
					"Items":    pathItemsXML,
				},
			},
		},
	}, nil
}

// ListInvalidations lists invalidations for a CloudFront distribution.
func (s *CloudFrontService) ListInvalidations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	distID := request.GetStringParam(req.Parameters, "Id")
	if distID == "" {
		return nil, errors.NewAWSError("InvalidArgument", "Required parameter Id is missing.", 400)
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := stores.distributions.Get(distID); err != nil {
		return nil, errors.NewAWSError("NoSuchDistribution", fmt.Sprintf("The specified distribution does not exist: %s", distID), 404)
	}

	invalidationMu.Lock()
	distInvalidations := invalidations[distID]
	invalidationMu.Unlock()

	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := request.GetIntParam(req.Parameters, "MaxItems")
	if maxItems <= 0 || maxItems > 200 {
		maxItems = 100
	}

	startIdx := 0
	if marker != "" {
		for i, inv := range distInvalidations {
			if id, ok := inv["Id"].(string); ok && id == marker {
				startIdx = i + 1
				break
			}
		}
	}

	available := distInvalidations[startIdx:]
	isTruncated := len(available) > int(maxItems)
	if isTruncated {
		available = available[:maxItems]
	}

	items := make([]map[string]interface{}, len(available))
	copy(items, available)

	nextMarker := ""
	if isTruncated && len(items) > 0 {
		if id, ok := items[len(items)-1]["Id"].(string); ok {
			nextMarker = id
		}
	}

	return map[string]interface{}{
		"IsTruncated": isTruncated,
		"Quantity":    len(items),
		"MaxItems":    maxItems,
		"Marker":      marker,
		"NextMarker":  nextMarker,
		"InvalidationList": map[string]interface{}{
			"IsTruncated": isTruncated,
			"Quantity":    len(items),
			"MaxItems":    maxItems,
			"Marker":      marker,
			"NextMarker":  nextMarker,
			"Items":       protocol.XMLElements{ElementName: "InvalidationSummary", Items: makeSliceFromMaps(items)},
		},
	}, nil
}

// GetInvalidation retrieves the status and details of a CloudFront invalidation.
func (s *CloudFrontService) GetInvalidation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	distID := request.GetStringParam(req.Parameters, "Id")
	if distID == "" {
		return nil, errors.NewAWSError("InvalidArgument", "Required parameter Id is missing.", 400)
	}

	invalidationID := request.GetStringParam(req.Parameters, "invalidationId")
	if invalidationID == "" {
		return nil, errors.NewAWSError("InvalidArgument", "Required parameter invalidationId is missing.", 400)
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := stores.distributions.Get(distID); err != nil {
		return nil, errors.NewAWSError("NoSuchDistribution", fmt.Sprintf("The specified distribution does not exist: %s", distID), 404)
	}

	invalidationMu.Lock()
	distInvalidations := invalidations[distID]
	invalidationMu.Unlock()

	for _, inv := range distInvalidations {
		if id, ok := inv["Id"].(string); ok && id == invalidationID {
			return map[string]interface{}{
				"Invalidation": inv,
			}, nil
		}
	}

	return nil, errors.NewAWSError("NoSuchInvalidation", fmt.Sprintf("The specified invalidation does not exist: %s", invalidationID), 404)
}

func makeSliceFromMaps(maps []map[string]interface{}) []interface{} {
	result := make([]interface{}, len(maps))
	for i, m := range maps {
		result[i] = m
	}
	return result
}
