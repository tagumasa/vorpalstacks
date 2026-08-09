package route53

import (
	"context"
	"fmt"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	route53store "vorpalstacks/internal/store/aws/route53"
)

// generateCidrCollectionId generates a UUID-like ID for a CIDR collection.
func generateCidrCollectionId() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// CreateCidrCollection creates a new CIDR collection.
func (s *Route53Service) CreateCidrCollection(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")
	if name == "" {
		return nil, awserrors.NewAWSError("InvalidInput", "Name is required", 400)
	}
	callerRef := request.GetStringParam(req.Parameters, "CallerReference")
	if callerRef == "" {
		return nil, awserrors.NewAWSError("InvalidInput", "CallerReference is required", 400)
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cidrStore := st.CidrCollections()
	if cidrStore == nil {
		return nil, awserrors.NewAWSError("InternalError", "CIDR collection store not available", 500)
	}

	collection := &route53store.CidrCollection{
		ID:              generateCidrCollectionId(),
		Name:            name,
		CallerReference: callerRef,
		AccountID:       s.accountID,
	}

	if err := cidrStore.Create(collection); err != nil {
		if route53store.IsAlreadyExists(err) {
			return nil, awserrors.NewAWSError("CidrCollectionAlreadyExists", fmt.Sprintf("A CIDR collection with name %q already exists", name), 400)
		}
		return nil, mapStoreError(err)
	}

	arn := fmt.Sprintf("arn:aws:route53:::cidrcollection/%s", collection.ID)

	return map[string]interface{}{
		"Collection": map[string]interface{}{
			"Arn":     arn,
			"Id":      collection.ID,
			"Name":    collection.Name,
			"Version": collection.Version,
		},
		"Location": "/cidrcollection/" + collection.ID,
	}, nil
}

// DeleteCidrCollection deletes a CIDR collection.
func (s *Route53Service) DeleteCidrCollection(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, awserrors.NewAWSError("InvalidInput", "Id is required", 400)
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cidrStore := st.CidrCollections()
	if cidrStore == nil {
		return nil, awserrors.NewAWSError("InternalError", "CIDR collection store not available", 500)
	}

	if err := cidrStore.Delete(id); err != nil {
		if route53store.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchCidrCollection", fmt.Sprintf("No CIDR collection found with id: %s", id), 404)
		}
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{}, nil
}

// ListCidrCollections lists CIDR collections with pagination.
func (s *Route53Service) ListCidrCollections(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	nextToken := request.GetStringParam(req.Parameters, "NextToken")
	maxResults := request.GetIntParam(req.Parameters, "MaxResults")

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cidrStore := st.CidrCollections()
	if cidrStore == nil {
		return nil, awserrors.NewAWSError("InternalError", "CIDR collection store not available", 500)
	}

	result, err := cidrStore.List(nextToken, maxResults)
	if err != nil {
		return nil, mapStoreError(err)
	}

	var summaries []interface{}
	for _, c := range result.Collections {
		arn := fmt.Sprintf("arn:aws:route53:::cidrcollection/%s", c.ID)
		summaries = append(summaries, map[string]interface{}{
			"Arn":     arn,
			"Id":      c.ID,
			"Name":    c.Name,
			"Version": c.Version,
		})
	}

	resp := map[string]interface{}{
		"CidrCollections": protocol.XMLElements{ElementName: "CidrCollection", Items: summaries},
	}
	if result.NextToken != "" {
		resp["NextToken"] = result.NextToken
	}
	return resp, nil
}

// ListCidrLocations lists locations within a CIDR collection.
func (s *Route53Service) ListCidrLocations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	collectionId := request.GetStringParam(req.Parameters, "CollectionId")
	if collectionId == "" {
		return nil, awserrors.NewAWSError("InvalidInput", "CollectionId is required", 400)
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cidrStore := st.CidrCollections()
	if cidrStore == nil {
		return nil, awserrors.NewAWSError("InternalError", "CIDR collection store not available", 500)
	}

	collection, err := cidrStore.Get(collectionId)
	if err != nil {
		if route53store.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchCidrCollection", fmt.Sprintf("No CIDR collection found with id: %s", collectionId), 404)
		}
		return nil, mapStoreError(err)
	}

	var locations []interface{}
	for locName := range collection.Locations {
		locations = append(locations, map[string]interface{}{
			"LocationName": locName,
		})
	}

	return map[string]interface{}{
		"CidrLocations": protocol.XMLElements{ElementName: "CidrLocation", Items: locations},
	}, nil
}

// ListCidrBlocks lists CIDR blocks within a location of a CIDR collection.
func (s *Route53Service) ListCidrBlocks(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	collectionId := request.GetStringParam(req.Parameters, "CollectionId")
	if collectionId == "" {
		return nil, awserrors.NewAWSError("InvalidInput", "CollectionId is required", 400)
	}
	locationName := request.GetStringParam(req.Parameters, "LocationName")

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cidrStore := st.CidrCollections()
	if cidrStore == nil {
		return nil, awserrors.NewAWSError("InternalError", "CIDR collection store not available", 500)
	}

	collection, err := cidrStore.Get(collectionId)
	if err != nil {
		if route53store.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchCidrCollection", fmt.Sprintf("No CIDR collection found with id: %s", collectionId), 404)
		}
		return nil, mapStoreError(err)
	}

	var blocks []interface{}
	if locationName != "" {
		for _, cidr := range collection.Locations[locationName] {
			blocks = append(blocks, map[string]interface{}{
				"CidrBlock":    cidr,
				"LocationName": locationName,
			})
		}
	} else {
		for locName, cidrs := range collection.Locations {
			for _, cidr := range cidrs {
				blocks = append(blocks, map[string]interface{}{
					"CidrBlock":    cidr,
					"LocationName": locName,
				})
			}
		}
	}

	return map[string]interface{}{
		"CidrBlocks": protocol.XMLElements{ElementName: "CidrBlock", Items: blocks},
	}, nil
}

// ChangeCidrCollection adds or removes CIDR blocks in a collection location.
func (s *Route53Service) ChangeCidrCollection(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, awserrors.NewAWSError("InvalidInput", "Id is required", 400)
	}

	var changesList []interface{}
	switch c := req.Parameters["Changes"].(type) {
	case []interface{}:
		changesList = c
	case map[string]interface{}:
		if changeArr, ok := c["CidrCollectionChange"].([]interface{}); ok {
			changesList = changeArr
		} else if changeMap, ok := c["CidrCollectionChange"].(map[string]interface{}); ok {
			changesList = []interface{}{changeMap}
		}
	default:
		return nil, awserrors.NewAWSError("InvalidInput", "Changes are required", 400)
	}

	if len(changesList) == 0 {
		return nil, awserrors.NewAWSError("InvalidInput", "Changes are required", 400)
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cidrStore := st.CidrCollections()
	if cidrStore == nil {
		return nil, awserrors.NewAWSError("InternalError", "CIDR collection store not available", 500)
	}

	collection, err := cidrStore.Get(id)
	if err != nil {
		if route53store.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchCidrCollection", fmt.Sprintf("No CIDR collection found with id: %s", id), 404)
		}
		return nil, mapStoreError(err)
	}

	if reqVersion := request.GetIntParam(req.Parameters, "CollectionVersion"); reqVersion > 0 {
		if int64(reqVersion) != collection.Version {
			return nil, awserrors.NewAWSError("CidrCollectionVersionMismatch",
				fmt.Sprintf("Collection version mismatch: expected %d, got %d", collection.Version, reqVersion), 409)
		}
	}

	for _, c := range changesList {
		changeMap, ok := c.(map[string]interface{})
		if !ok {
			return nil, awserrors.NewAWSError("InvalidInput", "Each change must be a map", 400)
		}

		locationName := request.GetStringParam(changeMap, "LocationName")
		if locationName == "" {
			return nil, awserrors.NewAWSError("InvalidInput", "LocationName is required for each change", 400)
		}
		action := request.GetStringParam(changeMap, "Action")
		if action != "PUT" && action != "DELETE_IF_EXISTS" {
			return nil, awserrors.NewAWSError("InvalidInput", fmt.Sprintf("Invalid action: %s. Must be PUT or DELETE_IF_EXISTS", action), 400)
		}

		var cidrs []string
		if cidrList, ok := changeMap["CidrList"].([]interface{}); ok {
			for _, c := range cidrList {
				if cidr, ok := c.(string); ok {
					cidrs = append(cidrs, cidr)
				}
			}
		} else if cidrMap, ok := changeMap["CidrList"].(map[string]interface{}); ok {
			if arr, ok := cidrMap["Cidr"].([]interface{}); ok {
				for _, c := range arr {
					if cidr, ok := c.(string); ok {
						cidrs = append(cidrs, cidr)
					}
				}
			} else if single, ok := cidrMap["Cidr"].(string); ok {
				cidrs = append(cidrs, single)
			}
		}

		if len(cidrs) == 0 {
			return nil, awserrors.NewAWSError("InvalidInput", "CidrList is required for each change", 400)
		}

		if collection.Locations == nil {
			collection.Locations = make(map[string][]string)
		}

		switch action {
		case "PUT":
			existing := collection.Locations[locationName]
			collection.Locations[locationName] = append(existing, cidrs...)
		case "DELETE_IF_EXISTS":
			existing := collection.Locations[locationName]
			toDelete := make(map[string]bool, len(cidrs))
			for _, cidr := range cidrs {
				toDelete[cidr] = true
			}
			var remaining []string
			for _, cidr := range existing {
				if !toDelete[cidr] {
					remaining = append(remaining, cidr)
				}
			}
			if len(remaining) > 0 {
				collection.Locations[locationName] = remaining
			} else {
				delete(collection.Locations, locationName)
			}
		}
	}

	if err := cidrStore.Update(collection); err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"Id": collection.ID,
	}, nil
}
