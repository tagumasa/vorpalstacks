package cloudfront

import (
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
)

// scanDistributions traverses every distribution page by page until fn
// returns true (match found) or all pages are exhausted. Reference checks
// use this instead of a single bounded list call so they remain correct
// regardless of how many distributions exist.
func scanDistributions(stores *cloudfrontStores, fn func(*cloudfrontstore.Distribution) bool) (bool, error) {
	marker := ""
	for {
		result, err := stores.distributions.List(marker, cloudfrontstore.DefaultListMaxItems)
		if err != nil {
			return false, err
		}
		for _, dist := range result.Distributions {
			if fn(dist) {
				return true, nil
			}
		}
		if !result.IsTruncated || result.NextMarker == "" {
			return false, nil
		}
		marker = result.NextMarker
	}
}

// scanKeyGroups traverses every key group page by page until fn returns
// true (match found) or all pages are exhausted.
func scanKeyGroups(stores *cloudfrontStores, fn func(*cloudfrontstore.KeyGroup) bool) (bool, error) {
	marker := ""
	for {
		result, err := stores.keyGroups.List(marker, cloudfrontstore.DefaultListMaxItems)
		if err != nil {
			return false, err
		}
		for _, kg := range result.KeyGroups {
			if fn(kg) {
				return true, nil
			}
		}
		if !result.IsTruncated || result.NextMarker == "" {
			return false, nil
		}
		marker = result.NextMarker
	}
}

// resolveListMaxItems applies the AWS list pagination default: a missing
// or non-positive MaxItems falls back to the default page size. AWS
// publishes no upper bound for CloudFront list operations, so positive
// values pass through unchanged and responses echo the effective value.
func resolveListMaxItems(raw int) int {
	if raw <= 0 {
		return cloudfrontstore.DefaultListMaxItems
	}
	return raw
}
