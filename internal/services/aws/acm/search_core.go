package acm

import (
	"fmt"

	awserrors "vorpalstacks/internal/common/errors"
	acmstorelib "vorpalstacks/internal/store/aws/acm"
)

// Smithy SearchMaxResults: @range(1-500) with @default 100.
const (
	searchMaxResultsDefault = 100
	searchMaxResultsLimit   = 500
)

// SearchCertificatesInput carries the wire-extracted fields of
// SearchCertificates. Params is the raw parameter map: the FilterStatement
// union is parsed and validated in the Core.
type SearchCertificatesInput struct {
	Params        map[string]interface{}
	NextToken     string
	MaxResults    int
	MaxResultsSet bool
	SortBy        string
	SortOrder     string
}

// SearchCertificatesResult is the transport-agnostic search page; the handler
// serialises each certificate via buildCertificateSearchResult.
type SearchCertificatesResult struct {
	Page      []*acmstorelib.Certificate
	NextToken string
}

// searchCertificatesCore is the single validation + persistence path for
// SearchCertificates: token and MaxResults validation, filter evaluation,
// sort validation, and offset-based pagination.
func (s *ACMService) searchCertificatesCore(stores *acmStores, in SearchCertificatesInput) (*SearchCertificatesResult, error) {
	if err := validateNextToken(in.NextToken); err != nil {
		return nil, err
	}

	maxResults := searchMaxResultsDefault
	if in.MaxResultsSet {
		if in.MaxResults < 1 || in.MaxResults > searchMaxResultsLimit {
			return nil, awserrors.NewValidationException(fmt.Sprintf("MaxResults must be between 1 and %d, got %d", searchMaxResultsLimit, in.MaxResults))
		}
		maxResults = in.MaxResults
	}

	allCerts, err := stores.certificates.ListAll()
	if err != nil {
		return nil, err
	}

	// Apply filters from the FilterStatement.
	filtered, err := applyCertificateFilters(allCerts, in.Params)
	if err != nil {
		return nil, err
	}

	// Sort results.
	if err := validateSearchSortBy(in.SortBy); err != nil {
		return nil, err
	}
	if in.SortOrder != "" {
		if err := validateSortOrder(in.SortOrder); err != nil {
			return nil, err
		}
	}
	sortCertificateSearchResults(filtered, in.SortBy, in.SortOrder)

	// Paginate using offset-based tokens.
	offset := 0
	if in.NextToken != "" {
		n, parseErr := parseIntToken(in.NextToken)
		if parseErr != nil {
			return nil, awserrors.NewValidationException(fmt.Sprintf("Invalid NextToken: %s", in.NextToken))
		}
		if n < 0 {
			return nil, awserrors.NewValidationException("Invalid NextToken: negative offset")
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

	nextToken := ""
	if end < len(filtered) {
		nextToken = formatIntToken(end)
	}

	return &SearchCertificatesResult{
		Page:      filtered[offset:end],
		NextToken: nextToken,
	}, nil
}
