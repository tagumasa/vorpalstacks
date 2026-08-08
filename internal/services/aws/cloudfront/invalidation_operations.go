package cloudfront

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/resilience"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
)

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
		if quantity != len(paths) {
			return nil, errors.NewAWSError("InconsistentQuantities",
				fmt.Sprintf("The Quantity value (%d) does not match the number of path items (%d)", quantity, len(paths)), 400)
		}
	}

	if len(paths) == 0 {
		return nil, errors.NewAWSError("InvalidArgument", "At least one path is required.", 400)
	}
	if len(paths) > 3000 {
		return nil, errors.NewAWSError("InvalidArgument", "Cannot invalidate more than 3000 paths in a single request.", 400)
	}

	inv, err := stores.invalidations.Create(distID, callerRef, paths)
	if err != nil {
		return nil, errors.NewAWSError("InternalError", fmt.Sprintf("Failed to create invalidation: %v", err), 500)
	}

	go s.transitionInvalidation(stores, inv)

	return map[string]interface{}{
		"Location":     fmt.Sprintf("/2020-05-31/distribution/%s/invalidation/%s", distID, inv.ID),
		"Invalidation": formatInvalidation(inv),
	}, nil
}

// transitionInvalidation asynchronously transitions an invalidation from
// InProgress to Completed, simulating the real CloudFront edge propagation.
func (s *CloudFrontService) transitionInvalidation(stores *cloudfrontStores, inv *cloudfrontstore.Invalidation) {
	defer func() {
		if r := resilience.RecoverPanic("cloudfront invalidation status transition"); r != nil {
			slog.Error("panic during invalidation status transition",
				"invalidationId", inv.ID, "distributionId", inv.DistributionID, "panic", r)
		}
	}()

	time.Sleep(2 * time.Second)

	inv.Status = "Completed"
	if err := stores.invalidations.Update(inv); err != nil {
		slog.Error("failed to persist invalidation status transition",
			"invalidationId", inv.ID, "distributionId", inv.DistributionID, "error", err)
	}
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

	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := request.GetIntParam(req.Parameters, "MaxItems")
	if maxItems <= 0 || maxItems > 100 {
		maxItems = 100
	}

	result, err := stores.invalidations.List(distID, marker, maxItems)
	if err != nil {
		return nil, errors.NewAWSError("InternalError", fmt.Sprintf("Failed to list invalidations: %v", err), 500)
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

	return map[string]interface{}{
		"InvalidationList": map[string]interface{}{
			"IsTruncated": result.IsTruncated,
			"Quantity":    len(items),
			"MaxItems":    maxItems,
			"Marker":      marker,
			"NextMarker":  nextMarker,
			"Items":       protocol.XMLElements{ElementName: "InvalidationSummary", Items: items},
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

	inv, err := stores.invalidations.Get(distID, invalidationID)
	if err != nil {
		return nil, errors.NewAWSError("NoSuchInvalidation", fmt.Sprintf("The specified invalidation does not exist: %s", invalidationID), 404)
	}

	return map[string]interface{}{
		"Invalidation": formatInvalidation(inv),
	}, nil
}
