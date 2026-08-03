package acm

import (
	"context"
	"fmt"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	acmstorelib "vorpalstacks/internal/store/aws/acm"
)

// SearchCertificates retrieves a list of certificates matching search criteria.
func (s *ACMService) SearchCertificates(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters

	nextToken := request.GetStringParam(params, "NextToken")
	if err := validateNextToken(nextToken); err != nil {
		return nil, err
	}
	maxResults := 100
	if _, ok := params["MaxResults"]; ok {
		mr := request.GetIntParam(params, "MaxResults")
		if mr < 1 || mr > 500 {
			return nil, awserrors.NewValidationException(fmt.Sprintf("MaxResults must be between 1 and 500, got %d", mr))
		}
		maxResults = mr
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	allCerts, err := stores.certificates.ListAll()
	if err != nil {
		return nil, err
	}

	// Apply filters from the FilterStatement.
	filtered, err := applyCertificateFilters(allCerts, params)
	if err != nil {
		return nil, err
	}

	// Sort results.
	sortBy := request.GetStringParam(params, "SortBy")
	sortOrder := request.GetStringParam(params, "SortOrder")
	if err := validateSearchSortBy(sortBy); err != nil {
		return nil, err
	}
	if sortOrder != "" {
		if err := validateSortOrder(sortOrder); err != nil {
			return nil, err
		}
	}
	sortCertificateSearchResults(filtered, sortBy, sortOrder)

	// Paginate using offset-based tokens.
	offset := 0
	if nextToken != "" {
		n, parseErr := parseIntToken(nextToken)
		if parseErr != nil {
			return nil, NewInvalidParameterError(fmt.Sprintf("Invalid NextToken: %s", nextToken))
		}
		if n < 0 {
			return nil, NewInvalidParameterError("Invalid NextToken: negative offset")
		}
		offset = n
	}
	if offset > len(filtered) {
		offset = len(filtered)
	}

	end := offset + maxResults
	if end > len(filtered) {
		end = len(filtered)
	}

	page := filtered[offset:end]

	results := make([]interface{}, 0, len(page))
	for _, c := range page {
		results = append(results, buildCertificateSearchResult(c))
	}

	resp := map[string]interface{}{
		"Results": results,
	}
	if end < len(filtered) {
		resp["NextToken"] = formatIntToken(end)
	}

	return resp, nil
}

// applyCertificateFilters filters certificates based on FilterStatement
// parameters. Supports the full CertificateFilterStatement union type
// including And/Or/Not compound operators and nested filter unions.
// Returns an error if the filter structure violates Smithy constraints
// (e.g. list length outside 1-15, missing REQUIRED fields on filter members).
func applyCertificateFilters(certs []*acmstorelib.Certificate, params map[string]interface{}) ([]*acmstorelib.Certificate, error) {
	filterStatement := getNestedValue(params, "FilterStatement")
	if filterStatement == nil {
		return certs, nil
	}
	if err := validateFilterStatement(filterStatement); err != nil {
		return nil, err
	}
	var result []*acmstorelib.Certificate
	for _, c := range certs {
		if evaluateFilterStatement(c, filterStatement) {
			result = append(result, c)
		}
	}
	return result, nil
}
