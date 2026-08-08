package cloudfront

// This file previously held CloudFrontStoresInterface, CloudFrontStores,
// NewCloudFrontStores, and per-store interface types that were never
// consumed outside the store package. They have been removed as dead code.
// The service layer (internal/services/aws/cloudfront) uses its own private
// cloudfrontStores struct which wraps the concrete store types directly.
