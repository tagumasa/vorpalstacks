package cloudfront

import (
	"context"
	"fmt"
	"time"

	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
)

// CreateInvalidation creates a new cache invalidation for a CloudFront distribution.
func (s *CloudFrontService) CreateInvalidation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	distID := request.GetStringParam(req.Parameters, "Id")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	batch := request.GetMapParam(req.Parameters, "InvalidationBatch")
	if batch == nil {
		batch = req.Parameters
	}

	callerRef, _ := batch["CallerReference"].(string)

	var paths []string
	var declaredQuantity int
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
		declaredQuantity = int(request.GetIntParam(pathsMap, "Quantity"))
	}

	inv, err := s.createInvalidationCore(stores, CreateInvalidationInput{
		Id:               distID,
		CallerReference:  callerRef,
		Paths:            paths,
		DeclaredQuantity: declaredQuantity,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Location":     fmt.Sprintf("/2020-05-31/distribution/%s/invalidation/%s", distID, inv.ID),
		"Invalidation": formatInvalidation(inv),
	}, nil
}

func formatInvalidation(inv *cloudfrontstore.Invalidation) map[string]interface{} {
	pathItemsXML := protocol.XMLElements{ElementName: "Path", Items: make([]interface{}, len(inv.Paths))}
	for i, p := range inv.Paths {
		pathItemsXML.Items[i] = p
	}

	return map[string]interface{}{
		"Id":         inv.ID,
		"CreateTime": inv.CreateTime.Format(time.RFC3339),
		"Status":     inv.Status,
		"InvalidationBatch": map[string]interface{}{
			"CallerReference": inv.CallerReference,
			"Paths": map[string]interface{}{
				"Quantity": len(inv.Paths),
				"Items":    pathItemsXML,
			},
		},
	}
}

// ListInvalidations lists invalidations for a CloudFront distribution.
func (s *CloudFrontService) ListInvalidations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	distID := request.GetStringParam(req.Parameters, "Id")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := resolveListMaxItems(request.GetIntParam(req.Parameters, "MaxItems"))

	result, err := s.listInvalidationsCore(stores, ListInvalidationsInput{
		Id:       distID,
		Marker:   marker,
		MaxItems: maxItems,
	})
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.Invalidations))
	for _, inv := range result.Invalidations {
		items = append(items, map[string]interface{}{
			"Id":         inv.ID,
			"CreateTime": inv.CreateTime.Format(time.RFC3339),
			"Status":     inv.Status,
		})
	}

	nextMarker := ""
	if result.IsTruncated && len(result.Invalidations) > 0 {
		nextMarker = result.Invalidations[len(result.Invalidations)-1].ID
	}

	invList := map[string]interface{}{
		"IsTruncated": result.IsTruncated,
		"Quantity":    len(items),
		"MaxItems":    maxItems,
		"Marker":      marker,
		"Items":       protocol.XMLElements{ElementName: "InvalidationSummary", Items: items},
	}
	if nextMarker != "" {
		invList["NextMarker"] = nextMarker
	}
	return map[string]interface{}{"InvalidationList": invList}, nil
}

// GetInvalidation retrieves the status and details of a CloudFront invalidation.
func (s *CloudFrontService) GetInvalidation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	invalidationID := request.GetStringParam(req.Parameters, "invalidationId")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	inv, err := s.getInvalidationCore(stores, GetInvalidationInput{
		Id:             request.GetStringParam(req.Parameters, "Id"),
		InvalidationId: invalidationID,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Invalidation": formatInvalidation(inv),
	}, nil
}
