package sesv2

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
)

// CreateEmailTemplate creates a new email template.
func (s *SESv2Service) CreateEmailTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	contentMap := request.GetMapParam(req.Parameters, "TemplateContent")
	if err := s.createEmailTemplateCore(store, CreateEmailTemplateInput{
		TemplateName: request.GetStringParam(req.Parameters, "TemplateName"),
		TemplateContent: TemplateContentInput{
			Map:      contentMap,
			Provided: contentMap != nil,
		},
		Tags: tags.ParseTags(req.Parameters, "Tags"),
	}); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// GetEmailTemplate retrieves the details of an email template.
func (s *SESv2Service) GetEmailTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.getEmailTemplateCore(store, request.GetStringParam(req.Parameters, "TemplateName"))
}

// UpdateEmailTemplate updates an existing email template.
func (s *SESv2Service) UpdateEmailTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	contentMap := request.GetMapParam(req.Parameters, "TemplateContent")
	if err := s.updateEmailTemplateCore(store, UpdateEmailTemplateInput{
		TemplateName: request.GetStringParam(req.Parameters, "TemplateName"),
		TemplateContent: TemplateContentInput{
			Map:      contentMap,
			Provided: contentMap != nil,
		},
	}); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// DeleteEmailTemplate deletes an email template.
func (s *SESv2Service) DeleteEmailTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteEmailTemplateCore(store, request.GetStringParam(req.Parameters, "TemplateName")); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// ListEmailTemplates returns a list of email templates.
func (s *SESv2Service) ListEmailTemplates(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.listEmailTemplatesCore(store,
		pagination.GetMaxItems(req.Parameters, 100, "PageSize"),
		pagination.GetMarker(req.Parameters, "NextToken"))
}

// TestRenderEmailTemplate renders an email template with the provided data.
func (s *SESv2Service) TestRenderEmailTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	_, templateDataProvided := req.Parameters["TemplateData"]
	return s.testRenderEmailTemplateCore(store, TestRenderEmailTemplateInput{
		TemplateName:         request.GetStringParam(req.Parameters, "TemplateName"),
		TemplateData:         request.GetStringParam(req.Parameters, "TemplateData"),
		TemplateDataProvided: templateDataProvided,
	})
}
