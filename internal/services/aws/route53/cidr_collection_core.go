package route53

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	route53store "vorpalstacks/internal/store/aws/route53"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input / Result structs
// ---------------------------------------------------------------------------

// CreateCidrCollectionInput carries the parameters for CreateCidrCollection.
type CreateCidrCollectionInput struct {
	Name            string
	CallerReference string
}

// DeleteCidrCollectionInput carries the parameters for DeleteCidrCollection.
type DeleteCidrCollectionInput struct {
	Id string
}

// ListCidrCollectionsInput carries the pagination parameters for
// ListCidrCollections.
type ListCidrCollectionsInput struct {
	NextToken  string
	MaxResults int
}

// ListCidrLocationsInput carries the parameters for ListCidrLocations.
type ListCidrLocationsInput struct {
	CollectionId string
	NextToken    string
	MaxResults   int
}

// ListCidrLocationsResult holds one page of location names together with
// the continuation token when the page is truncated.
type ListCidrLocationsResult struct {
	LocationNames []string
	NextToken     string
}

// CidrBlockSummary pairs one CIDR block with the location it belongs to.
type CidrBlockSummary struct {
	Cidr         string
	LocationName string
}

// ListCidrBlocksInput carries the parameters for ListCidrBlocks.
type ListCidrBlocksInput struct {
	CollectionId string
	LocationName string
	NextToken    string
	MaxResults   int
}

// ListCidrBlocksResult holds one page of CIDR blocks together with the
// continuation token when the page is truncated.
type ListCidrBlocksResult struct {
	Blocks    []CidrBlockSummary
	NextToken string
}

// ChangeCidrCollectionInput carries the parameters for ChangeCidrCollection.
// Changes keeps the raw wire element list — Core validates each element in
// the same position the original handler did, after the collection lookup
// and version check.
type ChangeCidrCollectionInput struct {
	Id                string
	CollectionVersion int
	Changes           []interface{}
}

// ChangeCidrCollectionResult holds the updated collection together with the
// change id that GetChange resolves.
type ChangeCidrCollectionResult struct {
	Collection *route53store.CidrCollection
	ChangeId   string
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// cidrCollectionStore returns the CIDR collection store or the
// InternalFailure AWS error (the only 500 the Route 53 contract documents)
// when the substrate is unavailable.
func cidrCollectionStore(st *route53store.Route53Stores) (*route53store.CidrCollectionStore, error) {
	cidrStore := st.CidrCollections()
	if cidrStore == nil {
		return nil, awserrors.NewInternalFailureException("CIDR collection store not available")
	}
	return cidrStore, nil
}

// createCidrCollectionCore is the single entry point for creating a CIDR
// collection.
func (s *Route53Service) createCidrCollectionCore(st *route53store.Route53Stores, input CreateCidrCollectionInput) (*route53store.CidrCollection, error) {
	if input.Name == "" {
		return nil, awserrors.NewAWSError("InvalidInput", "Name is required", 400)
	}
	if err := validateCidrCallerReference(input.CallerReference); err != nil {
		return nil, err
	}

	cidrStore, err := cidrCollectionStore(st)
	if err != nil {
		return nil, err
	}

	// Collection names are unique per account: a retry carrying the same
	// CallerReference gets the existing collection back, while the same
	// name under a different CallerReference is rejected.
	if existing, err := cidrStore.GetByName(input.Name); err == nil && existing != nil {
		if existing.CallerReference != input.CallerReference {
			return nil, awserrors.NewAWSError("CidrCollectionAlreadyExistsException", fmt.Sprintf("A CIDR collection with name %q already exists", input.Name), 400)
		}
		return existing, nil
	}

	collection := &route53store.CidrCollection{
		ID:              generateCidrCollectionId(),
		Name:            input.Name,
		CallerReference: input.CallerReference,
		AccountID:       s.accountID,
	}

	if err := cidrStore.Create(collection); err != nil {
		if route53store.IsAlreadyExists(err) {
			return nil, awserrors.NewAWSError("CidrCollectionAlreadyExistsException", fmt.Sprintf("A CIDR collection with name %q already exists", input.Name), 400)
		}
		return nil, mapStoreError(err)
	}
	return collection, nil
}

// deleteCidrCollectionCore is the single entry point for deleting a CIDR
// collection.
func deleteCidrCollectionCore(st *route53store.Route53Stores, input DeleteCidrCollectionInput) error {
	cidrStore, err := cidrCollectionStore(st)
	if err != nil {
		return err
	}

	if err := cidrStore.Delete(input.Id); err != nil {
		if route53store.IsNotFound(err) {
			return awserrors.NewAWSError("NoSuchCidrCollectionException", fmt.Sprintf("No CIDR collection found with id: %s", input.Id), 404)
		}
		return mapStoreError(err)
	}
	return nil
}

// listCidrCollectionsCore is the single entry point for listing CIDR
// collections.
func listCidrCollectionsCore(st *route53store.Route53Stores, input ListCidrCollectionsInput) (*route53store.CidrCollectionListResult, error) {
	cidrStore, err := cidrCollectionStore(st)
	if err != nil {
		return nil, err
	}

	result, err := cidrStore.List(input.NextToken, input.MaxResults)
	if err != nil {
		return nil, mapStoreError(err)
	}
	return result, nil
}

// listCidrLocationsCore is the single entry point for listing the location
// names of a CIDR collection. Names are returned in lexical order so pages
// are stable; the continuation token wraps the last returned name and the
// next page resumes strictly after it.
func listCidrLocationsCore(st *route53store.Route53Stores, input ListCidrLocationsInput) (*ListCidrLocationsResult, error) {
	if input.CollectionId == "" {
		return nil, awserrors.NewAWSError("InvalidInput", "CollectionId is required", 400)
	}

	cidrStore, err := cidrCollectionStore(st)
	if err != nil {
		return nil, err
	}

	collection, err := cidrStore.Get(input.CollectionId)
	if err != nil {
		if route53store.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchCidrCollectionException", fmt.Sprintf("No CIDR collection found with id: %s", input.CollectionId), 404)
		}
		return nil, mapStoreError(err)
	}

	names := make([]string, 0, len(collection.Locations))
	for locName := range collection.Locations {
		names = append(names, locName)
	}
	sort.Strings(names)

	if input.NextToken != "" {
		startAfter, err := decodeCidrLocationsToken(input.NextToken)
		if err != nil {
			return nil, invalidCidrPaginationTokenError()
		}
		start := sort.SearchStrings(names, startAfter)
		for start < len(names) && names[start] == startAfter {
			start++
		}
		names = names[start:]
	}

	nextToken := ""
	if input.MaxResults > 0 && len(names) > input.MaxResults {
		names = names[:input.MaxResults]
		nextToken = encodeCidrLocationsToken(names[len(names)-1])
	}

	return &ListCidrLocationsResult{LocationNames: names, NextToken: nextToken}, nil
}

// listCidrBlocksCore is the single entry point for listing the CIDR blocks
// of one location (or of every location, in lexical location order, when no
// location name is given). Blocks are walked in lexical order within each
// location so the walk always matches the lexical cursor comparison,
// regardless of the order the blocks were PUT in. The continuation token is
// the base64 of the last returned "location\x00cidr" pair; the next page
// resumes strictly after it.
func listCidrBlocksCore(st *route53store.Route53Stores, input ListCidrBlocksInput) (*ListCidrBlocksResult, error) {
	if input.CollectionId == "" {
		return nil, awserrors.NewAWSError("InvalidInput", "CollectionId is required", 400)
	}

	cidrStore, err := cidrCollectionStore(st)
	if err != nil {
		return nil, err
	}

	collection, err := cidrStore.Get(input.CollectionId)
	if err != nil {
		if route53store.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchCidrCollectionException", fmt.Sprintf("No CIDR collection found with id: %s", input.CollectionId), 404)
		}
		return nil, mapStoreError(err)
	}

	if input.LocationName != "" {
		if _, ok := collection.Locations[input.LocationName]; !ok {
			return nil, awserrors.NewAWSError("NoSuchCidrLocationException",
				fmt.Sprintf("No CIDR collection location found with name: %s", input.LocationName), 404)
		}
	}

	cursorLoc, cursorCidr, err := decodeCidrBlocksToken(input.NextToken)
	if err != nil {
		return nil, invalidCidrPaginationTokenError()
	}

	var blocks []CidrBlockSummary
	locationNames := make([]string, 0, len(collection.Locations))
	for locName := range collection.Locations {
		if input.LocationName == "" || locName == input.LocationName {
			locationNames = append(locationNames, locName)
		}
	}
	sort.Strings(locationNames)
	for _, locName := range locationNames {
		// The cursor compares blocks lexically, so each location's blocks
		// are walked in lexical order; the store keeps them in insertion
		// order, and walking that order would strand blocks that sort
		// before the cursor on later pages.
		cidrs := append([]string(nil), collection.Locations[locName]...)
		sort.Strings(cidrs)
		for _, cidr := range cidrs {
			if locName < cursorLoc || (locName == cursorLoc && cidr <= cursorCidr) {
				continue
			}
			blocks = append(blocks, CidrBlockSummary{Cidr: cidr, LocationName: locName})
		}
	}

	nextToken := ""
	if input.MaxResults > 0 && len(blocks) > input.MaxResults {
		blocks = blocks[:input.MaxResults]
		last := blocks[len(blocks)-1]
		nextToken = base64.StdEncoding.EncodeToString([]byte(last.LocationName + "\x00" + last.Cidr))
	}

	return &ListCidrBlocksResult{Blocks: blocks, NextToken: nextToken}, nil
}

// invalidCidrPaginationTokenError is returned when a CIDR list continuation
// token cannot be decoded. Both CIDR list operations document InvalidInput
// (HTTP 400), and starting from the beginning is specified only for a
// request that provides no token at all, so an undecodable token is
// rejected instead of silently restarting the listing from the first page.
func invalidCidrPaginationTokenError() error {
	return awserrors.NewAWSError("InvalidInput", "The NextToken provided is not valid", 400)
}

// encodeCidrLocationsToken wraps a location name so that a continuation
// token this API never produced is detectable on the next request.
func encodeCidrLocationsToken(location string) string {
	return base64.StdEncoding.EncodeToString([]byte(location))
}

// decodeCidrLocationsToken unwraps a ListCidrLocations continuation token
// into the location name to resume after. An empty or undecodable payload
// is rejected.
func decodeCidrLocationsToken(token string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", err
	}
	if len(raw) == 0 {
		return "", fmt.Errorf("cidr locations token is empty")
	}
	return string(raw), nil
}

// decodeCidrBlocksToken decodes the opaque ListCidrBlocks continuation
// token into its location and cidr cursor. Tokens are produced only by this
// listing, so a token that fails base64 decoding or lacks the
// location\x00cidr split is rejected as invalid input.
func decodeCidrBlocksToken(token string) (location, cidr string, err error) {
	if token == "" {
		return "", "", nil
	}
	raw, derr := base64.StdEncoding.DecodeString(token)
	if derr != nil {
		return "", "", derr
	}
	loc, c, ok := strings.Cut(string(raw), "\x00")
	if !ok || loc == "" || c == "" {
		return "", "", fmt.Errorf("cidr blocks token lacks a location or cidr cursor")
	}
	return loc, c, nil
}

// maxCidrCollectionChangesPerRequest is the CidrCollectionChanges @length
// maximum from the API model: one change batch carries between 1 and 1000
// entries.
const maxCidrCollectionChangesPerRequest = 1000

// The CIDR collection change member constraints from the API model: a
// location name is 1-16 characters from [0-9A-Za-z_-], a CIDR entry is
// 1-50 characters and must not be blank, and each change's Cidr list
// carries between 1 and 1000 entries.
const (
	maxCidrLocationNameLen = 16
	maxCidrLen             = 50
	maxCidrListLen         = 1000
)

// cidrLocationNameRegex is the Smithy @pattern of the CIDR collection
// location name member.
var cidrLocationNameRegex = regexp.MustCompile(`^[0-9A-Za-z_\-]+$`)

// cidrNonWhitespaceRegex implements the Smithy @pattern \S on the Cidr
// member with Smithy search semantics: the value must contain at least
// one non-whitespace character.
var cidrNonWhitespaceRegex = regexp.MustCompile(`\S`)

// validateCidrLocationName enforces the location name member constraints
// (length 1-16, alphanumeric/underscore/hyphen only).
func validateCidrLocationName(name string) error {
	if len(name) > maxCidrLocationNameLen {
		return awserrors.NewAWSError("InvalidInput", fmt.Sprintf("LocationName must be 1 to %d characters long", maxCidrLocationNameLen), 400)
	}
	if !cidrLocationNameRegex.MatchString(name) {
		return awserrors.NewAWSError("InvalidInput", "LocationName may contain only alphanumeric characters, underscores, and hyphens", 400)
	}
	return nil
}

// validateCidrEntry enforces the Cidr member constraints (length 1-50, at
// least one non-whitespace character).
func validateCidrEntry(cidr string) error {
	if len(cidr) > maxCidrLen {
		return awserrors.NewAWSError("InvalidInput", fmt.Sprintf("Cidr must be 1 to %d characters long", maxCidrLen), 400)
	}
	if !cidrNonWhitespaceRegex.MatchString(cidr) {
		return awserrors.NewAWSError("InvalidInput", "Cidr must not be blank", 400)
	}
	return nil
}

// changeCidrCollectionCore is the single entry point applying CIDR location
// changes to a collection. It enforces the batch size contract and the
// optional collection version match, validates every change element,
// applies PUT / DELETE_IF_EXISTS semantics per location, and records the
// ChangeInfo whose id the response carries (resolvable via GetChange).
func changeCidrCollectionCore(st *route53store.Route53Stores, input ChangeCidrCollectionInput) (*ChangeCidrCollectionResult, error) {
	if input.Id == "" {
		return nil, awserrors.NewAWSError("InvalidInput", "Id is required", 400)
	}
	if len(input.Changes) == 0 {
		return nil, awserrors.NewAWSError("InvalidInput", "Changes are required", 400)
	}
	if len(input.Changes) > maxCidrCollectionChangesPerRequest {
		return nil, awserrors.NewAWSError("InvalidInput",
			fmt.Sprintf("Changes must not exceed %d entries", maxCidrCollectionChangesPerRequest), 400)
	}

	cidrStore, err := cidrCollectionStore(st)
	if err != nil {
		return nil, err
	}

	collection, err := cidrStore.Get(input.Id)
	if err != nil {
		if route53store.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchCidrCollectionException", fmt.Sprintf("No CIDR collection found with id: %s", input.Id), 404)
		}
		return nil, mapStoreError(err)
	}

	if input.CollectionVersion > 0 {
		if int64(input.CollectionVersion) != collection.Version {
			return nil, awserrors.NewAWSError("CidrCollectionVersionMismatchException",
				fmt.Sprintf("Collection version mismatch: expected %d, got %d", collection.Version, input.CollectionVersion), 409)
		}
	}

	for _, c := range input.Changes {
		changeMap, ok := c.(map[string]interface{})
		if !ok {
			return nil, awserrors.NewAWSError("InvalidInput", "Each change must be a map", 400)
		}

		locationName := request.GetStringParam(changeMap, "LocationName")
		if locationName == "" {
			return nil, awserrors.NewAWSError("InvalidInput", "LocationName is required for each change", 400)
		}
		if err := validateCidrLocationName(locationName); err != nil {
			return nil, err
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
		if len(cidrs) > maxCidrListLen {
			return nil, awserrors.NewAWSError("InvalidInput", fmt.Sprintf("CidrList must contain at most %d entries", maxCidrListLen), 400)
		}
		for _, cidr := range cidrs {
			if err := validateCidrEntry(cidr); err != nil {
				return nil, err
			}
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

	// The batch is applied synchronously, so the change is already
	// propagated when the response returns.
	changeId := generateChangeId()
	if err := st.Changes().Create(&route53store.ChangeInfo{
		ID:          changeId,
		Status:      "INSYNC",
		SubmittedAt: time.Now(),
	}); err != nil {
		return nil, mapStoreError(err)
	}

	return &ChangeCidrCollectionResult{
		Collection: collection,
		ChangeId:   changeId,
	}, nil
}
