package s3

import (
	"vorpalstacks/internal/common/request"
	s3store "vorpalstacks/internal/store/aws/s3"
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
	config := &s3store.WebsiteConfiguration{}

	if input.WebsiteConfiguration.IndexDocument != nil {
		config.IndexDocument = input.WebsiteConfiguration.IndexDocument.Suffix
	}

	if input.WebsiteConfiguration.ErrorDocument != nil {
		config.ErrorDocument = input.WebsiteConfiguration.ErrorDocument.Key
	}

	if input.WebsiteConfiguration.RedirectAllRequestsTo != nil {
		config.RedirectAllRequestsTo = &s3store.RedirectAllRequestsTo{
			HostName: input.WebsiteConfiguration.RedirectAllRequestsTo.HostName,
			Protocol: input.WebsiteConfiguration.RedirectAllRequestsTo.Protocol,
		}
	}

	for _, rule := range input.WebsiteConfiguration.RoutingRules {
		routingRule := s3store.RoutingRule{}

		if rule.Condition != nil {
			routingRule.Condition = &s3store.RoutingRuleCondition{
				HTTPErrorCodeReturnedEquals: rule.Condition.HTTPErrorCodeReturnedEquals,
				KeyPrefixEquals:             rule.Condition.KeyPrefixEquals,
			}
		}

		if rule.Redirect != nil {
			routingRule.Redirect = &s3store.RoutingRuleRedirect{
				HostName:             rule.Redirect.HostName,
				HTTPRedirectCode:     rule.Redirect.HTTPRedirectCode,
				Protocol:             rule.Redirect.Protocol,
				ReplaceKeyPrefixWith: rule.Redirect.ReplaceKeyPrefixWith,
				ReplaceKeyWith:       rule.Redirect.ReplaceKeyWith,
			}
		}

		config.RoutingRules = append(config.RoutingRules, routingRule)
	}

	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}

	return store.buckets.SetWebsiteConfiguration(input.Bucket, config)
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

	bucket, err := store.buckets.Get(input.Bucket)
	if err != nil {
		return nil, err
	}

	if bucket.WebsiteConfiguration == nil {
		return nil, ErrNoSuchWebsite
	}

	output := &GetBucketWebsiteOutput{}

	if bucket.WebsiteConfiguration.IndexDocument != "" {
		output.IndexDocument = &IndexDocumentOutput{
			Suffix: bucket.WebsiteConfiguration.IndexDocument,
		}
	}

	if bucket.WebsiteConfiguration.ErrorDocument != "" {
		output.ErrorDocument = &ErrorDocumentOutput{
			Key: bucket.WebsiteConfiguration.ErrorDocument,
		}
	}

	if bucket.WebsiteConfiguration.RedirectAllRequestsTo != nil {
		output.RedirectAllRequestsTo = &RedirectAllRequestsToOutput{
			HostName: bucket.WebsiteConfiguration.RedirectAllRequestsTo.HostName,
			Protocol: bucket.WebsiteConfiguration.RedirectAllRequestsTo.Protocol,
		}
	}

	for _, rule := range bucket.WebsiteConfiguration.RoutingRules {
		outputRule := RoutingRuleOutput{}

		if rule.Condition != nil {
			outputRule.Condition = &RoutingRuleConditionOutput{
				HTTPErrorCodeReturnedEquals: rule.Condition.HTTPErrorCodeReturnedEquals,
				KeyPrefixEquals:             rule.Condition.KeyPrefixEquals,
			}
		}

		if rule.Redirect != nil {
			outputRule.Redirect = &RedirectOutput{
				HostName:             rule.Redirect.HostName,
				HTTPRedirectCode:     rule.Redirect.HTTPRedirectCode,
				Protocol:             rule.Redirect.Protocol,
				ReplaceKeyPrefixWith: rule.Redirect.ReplaceKeyPrefixWith,
				ReplaceKeyWith:       rule.Redirect.ReplaceKeyWith,
			}
		}

		output.RoutingRules = append(output.RoutingRules, outputRule)
	}

	return output, nil
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
	return store.buckets.SetWebsiteConfiguration(input.Bucket, nil)
}
