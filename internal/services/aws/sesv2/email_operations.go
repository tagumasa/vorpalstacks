package sesv2

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// SendEmail sends an email using the SESv2 service.
func (s *SESv2Service) SendEmail(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := rejectTenantName(req.Parameters); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.sendEmailCore(store, SendEmailInput{
		FromEmailAddress:                          request.GetStringParam(req.Parameters, "FromEmailAddress"),
		FromEmailAddressIdentityArn:               request.GetStringParam(req.Parameters, "FromEmailAddressIdentityArn"),
		FeedbackForwardingEmailAddress:            request.GetStringParam(req.Parameters, "FeedbackForwardingEmailAddress"),
		FeedbackForwardingEmailAddressIdentityArn: request.GetStringParam(req.Parameters, "FeedbackForwardingEmailAddressIdentityArn"),
		ConfigurationSetName:                      request.GetStringParam(req.Parameters, "ConfigurationSetName"),
		ReplyToAddresses:                          request.GetStringList(req.Parameters, "ReplyToAddresses"),
		EmailTags:                                 ParseMessageTags(req.Parameters, "EmailTags"),
		Destination:                               request.GetMapParam(req.Parameters, "Destination"),
		Content:                                   request.GetMapParam(req.Parameters, "Content"),
		ListManagementOptions:                     request.GetMapParam(req.Parameters, "ListManagementOptions"),
	})
}

// SendBulkEmail sends multiple emails in a single operation.
func (s *SESv2Service) SendBulkEmail(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := rejectTenantName(req.Parameters); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	entriesIf, _ := req.Parameters["BulkEmailEntries"].([]interface{})

	return s.sendBulkEmailCore(store, SendBulkEmailInput{
		FromEmailAddress:                          request.GetStringParam(req.Parameters, "FromEmailAddress"),
		FromEmailAddressIdentityArn:               request.GetStringParam(req.Parameters, "FromEmailAddressIdentityArn"),
		FeedbackForwardingEmailAddress:            request.GetStringParam(req.Parameters, "FeedbackForwardingEmailAddress"),
		FeedbackForwardingEmailAddressIdentityArn: request.GetStringParam(req.Parameters, "FeedbackForwardingEmailAddressIdentityArn"),
		ConfigurationSetName:                      request.GetStringParam(req.Parameters, "ConfigurationSetName"),
		ReplyToAddresses:                          request.GetStringList(req.Parameters, "ReplyToAddresses"),
		DefaultEmailTags:                          ParseMessageTags(req.Parameters, "DefaultEmailTags"),
		DefaultContent:                            request.GetMapParam(req.Parameters, "DefaultContent"),
		ListManagementOptions:                     request.GetMapParam(req.Parameters, "ListManagementOptions"),
		BulkEmailEntries:                          entriesIf,
	})
}
