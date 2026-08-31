package s3

import (
	"vorpalstacks/internal/common/request"
)

// PutBucketWebsiteInput contains the input parameters for the PutBucketWebsite operation.
type PutBucketWebsiteInput struct {
	Bucket               string
	WebsiteConfiguration *WebsiteConfigurationInput
}

// WebsiteConfigurationInput defines the website configuration for a bucket.
type WebsiteConfigurationInput struct {
	IndexDocument         *IndexDocumentInput         `xml:"IndexDocument,omitempty"`
	ErrorDocument         *ErrorDocumentInput         `xml:"ErrorDocument,omitempty"`
	RedirectAllRequestsTo *RedirectAllRequestsToInput `xml:"RedirectAllRequestsTo,omitempty"`
	RoutingRules          []RoutingRuleInput          `xml:"RoutingRules>RoutingRule,omitempty"`
}

// IndexDocumentInput defines the index document configuration.
type IndexDocumentInput struct {
	Suffix string `xml:"Suffix"`
}

// ErrorDocumentInput defines the error document configuration.
type ErrorDocumentInput struct {
	Key string `xml:"Key"`
}

// RedirectAllRequestsToInput defines redirect configuration for all requests.
type RedirectAllRequestsToInput struct {
	HostName string `xml:"HostName"`
	Protocol string `xml:"Protocol,omitempty"`
}

// RoutingRuleInput defines a routing rule for website configuration.
type RoutingRuleInput struct {
	Condition *RoutingRuleConditionInput `xml:"Condition,omitempty"`
	Redirect  *RedirectInput             `xml:"Redirect"`
}

// RoutingRuleConditionInput defines the condition for a routing rule.
type RoutingRuleConditionInput struct {
	HTTPErrorCodeReturnedEquals *string `xml:"HttpErrorCodeReturnedEquals,omitempty"`
	KeyPrefixEquals             *string `xml:"KeyPrefixEquals,omitempty"`
}

// RedirectInput defines the redirect configuration for a routing rule.
type RedirectInput struct {
	HostName             *string `xml:"HostName,omitempty"`
	HTTPRedirectCode     *string `xml:"HttpRedirectCode,omitempty"`
	Protocol             *string `xml:"Protocol,omitempty"`
	ReplaceKeyPrefixWith *string `xml:"ReplaceKeyPrefixWith,omitempty"`
	ReplaceKeyWith       *string `xml:"ReplaceKeyWith,omitempty"`
}

// PutBucketWebsite configures the website configuration for a bucket.
func (o *BucketOperations) PutBucketWebsite(ctx *request.RequestContext, input *PutBucketWebsiteInput) error {
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}
	return o.svc.putBucketWebsiteCore(store.buckets, input)
}

// GetBucketWebsiteInput contains the input parameters for the GetBucketWebsite operation.
type GetBucketWebsiteInput struct {
	Bucket string
}

// GetBucketWebsiteOutput contains the output from the GetBucketWebsite operation.
type GetBucketWebsiteOutput struct {
	IndexDocument         *IndexDocumentOutput         `xml:"IndexDocument,omitempty"`
	ErrorDocument         *ErrorDocumentOutput         `xml:"ErrorDocument,omitempty"`
	RedirectAllRequestsTo *RedirectAllRequestsToOutput `xml:"RedirectAllRequestsTo,omitempty"`
	RoutingRules          []RoutingRuleOutput          `xml:"RoutingRules>RoutingRule,omitempty"`
}

// IndexDocumentOutput defines the index document in the response.
type IndexDocumentOutput struct {
	Suffix string `xml:"Suffix"`
}

// ErrorDocumentOutput defines the error document in the response.
type ErrorDocumentOutput struct {
	Key string `xml:"Key"`
}

// RedirectAllRequestsToOutput defines redirect configuration in the response.
type RedirectAllRequestsToOutput struct {
	HostName string `xml:"HostName"`
	Protocol string `xml:"Protocol,omitempty"`
}

// RoutingRuleOutput defines a routing rule in the response.
type RoutingRuleOutput struct {
	Condition *RoutingRuleConditionOutput `xml:"Condition,omitempty"`
	Redirect  *RedirectOutput             `xml:"Redirect"`
}

// RoutingRuleConditionOutput defines the condition for a routing rule in the response.
type RoutingRuleConditionOutput struct {
	HTTPErrorCodeReturnedEquals *string `xml:"HttpErrorCodeReturnedEquals,omitempty"`
	KeyPrefixEquals             *string `xml:"KeyPrefixEquals,omitempty"`
}

// RedirectOutput defines the redirect configuration in the response.
type RedirectOutput struct {
	HostName             *string `xml:"HostName,omitempty"`
	HTTPRedirectCode     *string `xml:"HttpRedirectCode,omitempty"`
	Protocol             *string `xml:"Protocol,omitempty"`
	ReplaceKeyPrefixWith *string `xml:"ReplaceKeyPrefixWith,omitempty"`
	ReplaceKeyWith       *string `xml:"ReplaceKeyWith,omitempty"`
}

// GetBucketWebsite retrieves the website configuration for a bucket.
func (o *BucketOperations) GetBucketWebsite(ctx *request.RequestContext, input *GetBucketWebsiteInput) (*GetBucketWebsiteOutput, error) {
	store, err := o.svc.store(ctx)
	if err != nil {
		return nil, err
	}
	return o.svc.getBucketWebsiteCore(store.buckets, input)
}

// DeleteBucketWebsiteInput contains the input parameters for the DeleteBucketWebsite operation.
type DeleteBucketWebsiteInput struct {
	Bucket string
}

// DeleteBucketWebsite removes the website configuration from a bucket.
func (o *BucketOperations) DeleteBucketWebsite(ctx *request.RequestContext, input *DeleteBucketWebsiteInput) error {
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}
	return o.svc.deleteBucketWebsiteCore(store.buckets, input)
}
