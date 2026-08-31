package route53

import (
	"context"
	"fmt"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
)

// generateCidrCollectionId generates a UUID-like ID for a CIDR collection.
func generateCidrCollectionId() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// CreateCidrCollection creates a new CIDR collection.
func (s *Route53Service) CreateCidrCollection(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	collection, err := s.createCidrCollectionCore(st, CreateCidrCollectionInput{
		Name:            request.GetStringParam(req.Parameters, "Name"),
		CallerReference: request.GetStringParam(req.Parameters, "CallerReference"),
	})
	if err != nil {
		return nil, err
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
	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := deleteCidrCollectionCore(st, DeleteCidrCollectionInput{
		Id: request.GetStringParam(req.Parameters, "Id"),
	}); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// ListCidrCollections lists CIDR collections with pagination.
func (s *Route53Service) ListCidrCollections(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := listCidrCollectionsCore(st, ListCidrCollectionsInput{
		NextToken:  request.GetStringParam(req.Parameters, "NextToken"),
		MaxResults: request.GetIntParam(req.Parameters, "MaxResults"),
	})
	if err != nil {
		return nil, err
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
		// The CIDR list shapes carry no xmlName trait, so restXml clients
		// read every entry from a <member> element.
		"CidrCollections": protocol.XMLElements{ElementName: "member", Items: summaries},
	}
	if result.NextToken != "" {
		resp["NextToken"] = result.NextToken
	}
	return resp, nil
}

// ListCidrLocations lists locations within a CIDR collection.
func (s *Route53Service) ListCidrLocations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := listCidrLocationsCore(st, ListCidrLocationsInput{
		CollectionId: request.GetStringParam(req.Parameters, "CollectionId"),
		NextToken:    request.GetStringParam(req.Parameters, "NextToken"),
		MaxResults:   request.GetIntParam(req.Parameters, "MaxResults"),
	})
	if err != nil {
		return nil, err
	}

	var locationItems []interface{}
	for _, locName := range result.LocationNames {
		locationItems = append(locationItems, map[string]interface{}{
			"LocationName": locName,
		})
	}

	resp := map[string]interface{}{
		// The CIDR list shapes carry no xmlName trait, so restXml clients
		// read every entry from a <member> element.
		"CidrLocations": protocol.XMLElements{ElementName: "member", Items: locationItems},
	}
	if result.NextToken != "" {
		resp["NextToken"] = result.NextToken
	}
	return resp, nil
}

// ListCidrBlocks lists CIDR blocks within a location of a CIDR collection.
func (s *Route53Service) ListCidrBlocks(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	locationName := request.GetStringParam(req.Parameters, "LocationName")
	if locationName == "" {
		// The model binds LocationName to the httpQuery key "location",
		// which differs from the member name by more than case.
		locationName = request.GetStringParam(req.Parameters, "location")
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := listCidrBlocksCore(st, ListCidrBlocksInput{
		CollectionId: request.GetStringParam(req.Parameters, "CollectionId"),
		LocationName: locationName,
		NextToken:    request.GetStringParam(req.Parameters, "NextToken"),
		MaxResults:   request.GetIntParam(req.Parameters, "MaxResults"),
	})
	if err != nil {
		return nil, err
	}

	var blockItems []interface{}
	for _, b := range result.Blocks {
		blockItems = append(blockItems, map[string]interface{}{
			"CidrBlock":    b.Cidr,
			"LocationName": b.LocationName,
		})
	}

	resp := map[string]interface{}{
		// The CIDR list shapes carry no xmlName trait, so restXml clients
		// read every entry from a <member> element.
		"CidrBlocks": protocol.XMLElements{ElementName: "member", Items: blockItems},
	}
	if result.NextToken != "" {
		resp["NextToken"] = result.NextToken
	}
	return resp, nil
}

// ChangeCidrCollection adds or removes CIDR blocks in a collection location.
func (s *Route53Service) ChangeCidrCollection(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var changesList []interface{}
	switch c := req.Parameters["Changes"].(type) {
	case []interface{}:
		changesList = c
	case map[string]interface{}:
		if changeArr, ok := c["CidrCollectionChange"].([]interface{}); ok {
			changesList = changeArr
		} else if changeMap, ok := c["CidrCollectionChange"].(map[string]interface{}); ok {
			changesList = []interface{}{changeMap}
		} else if memberArr, ok := c["member"].([]interface{}); ok {
			// The model carries no xmlName on CidrCollectionChange, so
			// restXml clients serialise the list with the default
			// <member> element.
			changesList = memberArr
		} else if memberMap, ok := c["member"].(map[string]interface{}); ok {
			changesList = []interface{}{memberMap}
		}
	default:
		return nil, awserrors.NewAWSError("InvalidInput", "Changes are required", 400)
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := changeCidrCollectionCore(st, ChangeCidrCollectionInput{
		Id:                request.GetStringParam(req.Parameters, "Id"),
		CollectionVersion: request.GetIntParam(req.Parameters, "CollectionVersion"),
		Changes:           changesList,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Id": result.ChangeId,
	}, nil
}
