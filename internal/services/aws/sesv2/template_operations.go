package sesv2

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"text/template"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	"vorpalstacks/internal/services/aws/common/tags"
	"vorpalstacks/internal/store/aws/common"
	sesv2store "vorpalstacks/internal/store/aws/sesv2"
)

// CreateEmailTemplate creates a new email template.
func (s *SESv2Service) CreateEmailTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	templateName := request.GetStringParam(req.Parameters, "TemplateName")
	if templateName == "" {
		return nil, ErrMissingParameter
	}

	parsedTags := tags.ParseTags(req.Parameters, "Tags")

	tmpl := sesv2store.NewEmailTemplate(templateName)

	if contentMap := request.GetMapParam(req.Parameters, "TemplateContent"); contentMap != nil {
		tmpl.TemplateContent = &sesv2store.EmailTemplateContent{
			Subject: request.GetStringParam(contentMap, "Subject"),
			Html:    request.GetStringParam(contentMap, "Html"),
			Text:    request.GetStringParam(contentMap, "Text"),
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	_, err = store.CreateEmailTemplate(tmpl)
	if err != nil {
		return nil, err
	}

	if len(parsedTags) > 0 {
		arn := store.BuildTemplateArn(templateName)
		if err := store.TagResourceFromSlice(arn, parsedTags); err != nil {
			return nil, err
		}
	}

	return response.EmptyResponse(), nil
}

// GetEmailTemplate retrieves the details of an email template.
func (s *SESv2Service) GetEmailTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	templateName := request.GetStringParam(req.Parameters, "TemplateName")
	if templateName == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	tmpl, err := store.GetEmailTemplate(templateName)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"TemplateName": tmpl.TemplateName,
	}

	if tmpl.TemplateContent != nil {
		response["TemplateContent"] = map[string]interface{}{
			"Subject": tmpl.TemplateContent.Subject,
			"Html":    tmpl.TemplateContent.Html,
			"Text":    tmpl.TemplateContent.Text,
		}
	}

	return response, nil
}

// UpdateEmailTemplate updates an existing email template.
func (s *SESv2Service) UpdateEmailTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	templateName := request.GetStringParam(req.Parameters, "TemplateName")
	if templateName == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	tmpl, err := store.GetEmailTemplate(templateName)
	if err != nil {
		return nil, err
	}

	if contentMap := request.GetMapParam(req.Parameters, "TemplateContent"); contentMap != nil {
		tmpl.TemplateContent = &sesv2store.EmailTemplateContent{
			Subject: request.GetStringParam(contentMap, "Subject"),
			Html:    request.GetStringParam(contentMap, "Html"),
			Text:    request.GetStringParam(contentMap, "Text"),
		}
	}

	if err := store.UpdateEmailTemplate(tmpl); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DeleteEmailTemplate deletes an email template.
func (s *SESv2Service) DeleteEmailTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	templateName := request.GetStringParam(req.Parameters, "TemplateName")
	if templateName == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.DeleteEmailTemplate(templateName); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListEmailTemplates returns a list of email templates.
func (s *SESv2Service) ListEmailTemplates(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	pageSize := request.GetIntParam(req.Parameters, "PageSize")
	if pageSize == 0 {
		pageSize = 100
	}
	nextToken := request.GetStringParam(req.Parameters, "NextToken")

	opts := common.ListOptions{
		MaxItems: pageSize,
		Marker:   nextToken,
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := store.ListEmailTemplates(opts)
	if err != nil {
		return nil, err
	}

	templates := make([]map[string]interface{}, 0, len(result.Items))
	for _, tmpl := range result.Items {
		templates = append(templates, map[string]interface{}{
			"TemplateName": tmpl.TemplateName,
		})
	}

	response := map[string]interface{}{
		"TemplatesMetadata": templates,
	}

	if result.IsTruncated {
		response["NextToken"] = result.NextMarker
	}

	return response, nil
}

// TestRenderEmailTemplate renders an email template with the provided data.
func (s *SESv2Service) TestRenderEmailTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	templateName := request.GetStringParam(req.Parameters, "TemplateName")
	if templateName == "" {
		return nil, ErrMissingParameter
	}

	tmplData := request.GetStringParam(req.Parameters, "TemplateData")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	tmpl, err := store.GetEmailTemplate(templateName)
	if err != nil {
		return nil, err
	}

	if tmpl.TemplateContent == nil {
		return nil, NewBadRequestException("Template has no content")
	}

	data := make(map[string]interface{})
	if tmplData != "" {
		if err := json.Unmarshal([]byte(tmplData), &data); err != nil {
			return nil, NewBadRequestException("Invalid TemplateData JSON")
		}
	}

	renderedSubject, err := renderTemplateContent(tmpl.TemplateContent.Subject, data)
	if err != nil {
		return nil, NewBadRequestException("Failed to render subject: " + err.Error())
	}

	renderedHtml, err := renderTemplateContent(tmpl.TemplateContent.Html, data)
	if err != nil {
		return nil, NewBadRequestException("Failed to render HTML: " + err.Error())
	}

	renderedText, err := renderTemplateContent(tmpl.TemplateContent.Text, data)
	if err != nil {
		return nil, NewBadRequestException("Failed to render text: " + err.Error())
	}

	return map[string]interface{}{
		"RenderedTemplate": map[string]interface{}{
			"Subject": renderedSubject,
			"Html":    renderedHtml,
			"Text":    renderedText,
		},
	}, nil
}

func renderTemplateContent(templateStr string, data map[string]interface{}) (string, error) {
	if templateStr == "" {
		return "", nil
	}

	converted := convertAWSTemplateSyntax(templateStr)

	tmpl, err := template.New("template").Parse(converted)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func convertAWSTemplateSyntax(s string) string {
	s = strings.ReplaceAll(s, "{{#if", "{{if")
	s = strings.ReplaceAll(s, "{{#each", "{{range")
	s = strings.ReplaceAll(s, "{{#with", "{{with")
	s = strings.ReplaceAll(s, "{{/if}}", "{{end}}")
	s = strings.ReplaceAll(s, "{{/each}}", "{{end}}")
	s = strings.ReplaceAll(s, "{{/with}}", "{{end}}")
	s = strings.ReplaceAll(s, "{{else}}", "{{else}}")
	s = strings.ReplaceAll(s, "{{^if", "{{if not ")
	s = strings.ReplaceAll(s, "{{^each", "{{range")
	return s
}
