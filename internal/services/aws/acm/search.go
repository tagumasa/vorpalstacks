package acm

import (
	"context"

	"vorpalstacks/internal/common/request"
	acmstorelib "vorpalstacks/internal/store/aws/acm"
)

// SearchCertificates retrieves a list of certificates matching search criteria.
func (s *ACMService) SearchCertificates(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters

	in := SearchCertificatesInput{
		Params:    params,
		NextToken: request.GetStringParam(params, "NextToken"),
		SortBy:    request.GetStringParam(params, "SortBy"),
		SortOrder: request.GetStringParam(params, "SortOrder"),
	}
	if _, ok := params["MaxResults"]; ok {
		in.MaxResults = request.GetIntParam(params, "MaxResults")
		in.MaxResultsSet = true
	}

	result, err := s.searchCertificatesCore(reqCtx, in)
	if err != nil {
		return nil, err
	}

	results := make([]interface{}, 0, len(result.Page))
	for _, c := range result.Page {
		results = append(results, buildCertificateSearchResult(c))
	}

	resp := map[string]interface{}{
		"Results": results,
	}
	if result.NextToken != "" {
		resp["NextToken"] = result.NextToken
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
