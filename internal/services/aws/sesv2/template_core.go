package sesv2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"text/template"
	"unicode/utf8"

	awserrors "vorpalstacks/internal/common/errors"
	pagination "vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/store/aws/common"
	sesv2store "vorpalstacks/internal/store/aws/sesv2"
)

// ---------------------------------------------------------------------------
// Input DTOs — email-template family
// ---------------------------------------------------------------------------

// TemplateContentInput carries the raw TemplateContent wire map so the Core
// can distinguish an absent map from an explicitly empty one.
type TemplateContentInput struct {
	Map      map[string]interface{}
	Provided bool
}

// CreateEmailTemplateInput carries every create-time template member.
type CreateEmailTemplateInput struct {
	TemplateName    string
	TemplateContent TemplateContentInput
	Tags            []tags.Tag
}

// UpdateEmailTemplateInput carries the update members; TemplateContent is
// presence-based.
type UpdateEmailTemplateInput struct {
	TemplateName    string
	TemplateContent TemplateContentInput
}

// TestRenderEmailTemplateInput carries the render members. TemplateData is
// Smithy @required, so its presence travels on TemplateDataProvided.
type TestRenderEmailTemplateInput struct {
	TemplateName         string
	TemplateData         string
	TemplateDataProvided bool
}

var (
	awsVarPattern     = regexp.MustCompile(`\{\{([^#/^}{]+)\}\}`)
	awsThisPattern    = regexp.MustCompile(`\{\{\s*this\s*\}\}`)
	awsThisFieldStart = regexp.MustCompile(`\{\{\s*this\.`)

	// unsupportedHandlebarsPatterns detect constructs that
	// convertAWSTemplateSyntax cannot translate.  Matching any of these
	// causes renderTemplateContent to return a BadRequest so the caller
	// knows the template will not render correctly.
	unsupportedHandlebarsPatterns = []*regexp.Regexp{
		regexp.MustCompile(`\{\{>`),  // partials: {{> name}}
		regexp.MustCompile(`\{\{@`),  // @key, @index, etc.
		regexp.MustCompile(` as \|`), // block params: #each xs as |x|
	}
)

// convertAWSTemplateSyntax performs a best-effort translation from the
// Handlebars-style template syntax that AWS SES accepts to Go text/template
// syntax. The translation covers the most common constructs:
//
//   - {{var}}              -> {{.Var}}          (capitalised field access)
//   - {{this}}             -> {{.}}             (current element)
//   - {{this.field}}       -> {{.Field}}        (current element field)
//   - {{#if x}}...{{/if}}  -> {{if .X}}...{{end}}
//   - {{#each xs}}...{{/each}} -> {{range .Xs}}...{{end}}
//   - {{#with x}}...{{/with}}  -> {{with .X}}...{{end}}
//   - {{^if x}}...{{/if}}      -> {{if not .X}}...{{end}}
//   - {{else}}                  preserved
//
// This is not a complete Handlebars implementation. Constructs such as
// helpers ({{uppercase x}}), block parameters, partials and {{@key}} are
// not supported. AWS users relying on those constructs should switch to
// the simpler Go-template subset.
func convertAWSTemplateSyntax(s string) string {
	s = awsThisFieldStart.ReplaceAllString(s, "{{this.")
	s = awsThisPattern.ReplaceAllString(s, "{{.}}")
	s = strings.ReplaceAll(s, "{{#if", "{{if")
	s = strings.ReplaceAll(s, "{{#each", "{{range")
	s = strings.ReplaceAll(s, "{{#with", "{{with")
	s = strings.ReplaceAll(s, "{{/if}}", "{{end}}")
	s = strings.ReplaceAll(s, "{{/each}}", "{{end}}")
	s = strings.ReplaceAll(s, "{{/with}}", "{{end}}")
	s = strings.ReplaceAll(s, "{{else}}", "{{else}}")
	s = strings.ReplaceAll(s, "{{^if", "{{if not ")
	s = awsVarPattern.ReplaceAllStringFunc(s, func(match string) string {
		varName := awsVarPattern.FindStringSubmatch(match)[1]
		varName = strings.TrimSpace(varName)
		if varName == "" {
			return match
		}
		// Skip Go-template control keywords we just emitted above.
		switch varName {
		case "if", "else", "end", "range", "with":
			return match
		}
		if strings.HasPrefix(varName, "if ") || strings.HasPrefix(varName, "not ") {
			return match
		}
		if strings.HasPrefix(varName, ".") {
			return match
		}
		// Capitalise first letter: {{name}} -> {{.Name}}, matching AWS's
		// Handlebars convention of treating top-level vars as struct fields.
		return "{{." + strings.ToUpper(varName[:1]) + varName[1:] + "}}"
	})
	return s
}

// renderTemplateContent renders a single template part against the data map.
func renderTemplateContent(templateStr string, data map[string]interface{}) (string, error) {
	if templateStr == "" {
		return "", nil
	}

	for _, p := range unsupportedHandlebarsPatterns {
		if p.MatchString(templateStr) {
			return "", fmt.Errorf("template contains unsupported Handlebars syntax; use Go template syntax instead")
		}
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

// templateContentFromInput builds the stored content struct from a raw wire
// map. A nil map yields nil so the caller can distinguish "no content
// member" from "empty content member".
func templateContentFromInput(in TemplateContentInput) *sesv2store.EmailTemplateContent {
	if !in.Provided || in.Map == nil {
		return nil
	}
	return &sesv2store.EmailTemplateContent{
		Subject: request.GetStringParam(in.Map, "Subject"),
		Html:    request.GetStringParam(in.Map, "Html"),
		Text:    request.GetStringParam(in.Map, "Text"),
	}
}

// ---------------------------------------------------------------------------
// Core functions — email-template family
// ---------------------------------------------------------------------------

// createEmailTemplateCore is the single entry point for email-template
// creation. Per AWS, a template must contain at least one of Subject, Html,
// or Text; we replicate that validation rather than persisting a
// contentless template that TestRenderEmailTemplate would later fail on.
func (s *SESv2Service) createEmailTemplateCore(store sesv2store.SESv2StoreInterface, in CreateEmailTemplateInput) error {
	if in.TemplateName == "" {
		return ErrMissingParameter
	}
	if !validateEmailTemplateName(in.TemplateName) {
		return ErrBadRequest
	}

	tmpl := sesv2store.NewEmailTemplate(in.TemplateName)
	tmpl.TemplateContent = templateContentFromInput(in.TemplateContent)

	if tmpl.TemplateContent == nil ||
		(tmpl.TemplateContent.Subject == "" && tmpl.TemplateContent.Html == "" && tmpl.TemplateContent.Text == "") {
		return ErrInvalidParameter
	}

	if _, err := store.CreateEmailTemplate(tmpl); err != nil {
		return err
	}

	if len(in.Tags) > 0 {
		arn := store.BuildTemplateArn(in.TemplateName)
		if err := store.TagFromSlice(arn, in.Tags); err != nil {
			return err
		}
	}

	return nil
}

// getEmailTemplateCore is the single entry point for reading an email
// template, including its tags.
func (s *SESv2Service) getEmailTemplateCore(store sesv2store.SESv2StoreInterface, templateName string) (map[string]interface{}, error) {
	if templateName == "" {
		return nil, ErrMissingParameter
	}

	tmpl, err := store.GetEmailTemplate(templateName)
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"TemplateName": tmpl.TemplateName,
	}

	if tmpl.TemplateContent != nil {
		resp["TemplateContent"] = map[string]interface{}{
			"Subject": tmpl.TemplateContent.Subject,
			"Html":    tmpl.TemplateContent.Html,
			"Text":    tmpl.TemplateContent.Text,
		}
	}

	arn := store.BuildTemplateArn(templateName)
	if tags, err := store.ListAsSlice(arn); err == nil && len(tags) > 0 {
		resp["Tags"] = tags
	}

	return resp, nil
}

// updateEmailTemplateCore is the single entry point for updating an email
// template. The TemplateContent member is presence-based: an absent member
// keeps the stored content.
func (s *SESv2Service) updateEmailTemplateCore(store sesv2store.SESv2StoreInterface, in UpdateEmailTemplateInput) error {
	if in.TemplateName == "" {
		return ErrMissingParameter
	}

	tmpl, err := store.GetEmailTemplate(in.TemplateName)
	if err != nil {
		return err
	}

	if content := templateContentFromInput(in.TemplateContent); content != nil {
		tmpl.TemplateContent = content
	}

	if tmpl.TemplateContent == nil ||
		(tmpl.TemplateContent.Subject == "" && tmpl.TemplateContent.Html == "" && tmpl.TemplateContent.Text == "") {
		return ErrInvalidParameter
	}

	return store.UpdateEmailTemplate(tmpl)
}

// deleteEmailTemplateCore is the single entry point for deleting an email
// template.
func (s *SESv2Service) deleteEmailTemplateCore(store sesv2store.SESv2StoreInterface, templateName string) error {
	if templateName == "" {
		return ErrMissingParameter
	}
	return store.DeleteEmailTemplate(templateName)
}

// listEmailTemplatesCore is the single entry point for listing email
// templates.
func (s *SESv2Service) listEmailTemplatesCore(store sesv2store.SESv2StoreInterface, maxItems int, nextToken string) (map[string]interface{}, error) {
	opts := common.ListOptions{
		MaxItems: maxItems,
		Marker:   nextToken,
	}

	result, err := store.ListEmailTemplates(opts)
	if err != nil {
		return nil, err
	}

	templates := make([]map[string]interface{}, 0, len(result.Items))
	for _, tmpl := range result.Items {
		entry := map[string]interface{}{
			"TemplateName":     tmpl.TemplateName,
			"CreatedTimestamp": float64(tmpl.CreatedTimestamp.Unix()),
		}
		templates = append(templates, entry)
	}

	resp := map[string]interface{}{
		"TemplatesMetadata": templates,
	}

	pagination.SetNextToken(resp, "NextToken", result.NextMarker)

	return resp, nil
}

// testRenderEmailTemplateCore is the single entry point for rendering an
// email template with the provided data.
func (s *SESv2Service) testRenderEmailTemplateCore(store sesv2store.SESv2StoreInterface, in TestRenderEmailTemplateInput) (map[string]interface{}, error) {
	if in.TemplateName == "" {
		return nil, ErrMissingParameter
	}

	// TemplateData is Smithy @required.
	if !in.TemplateDataProvided {
		return nil, ErrMissingParameter
	}

	// Enforce the Smithy EmailTemplateData @length(max=262144) constraint,
	// counted in Unicode characters (the shape carries no pattern), to
	// prevent unbounded JSON payload DoS.
	if utf8.RuneCountInString(in.TemplateData) > maxTemplateDataSize {
		return nil, ErrBadRequest
	}

	tmpl, err := store.GetEmailTemplate(in.TemplateName)
	if err != nil {
		return nil, err
	}

	if tmpl.TemplateContent == nil {
		return nil, awserrors.NewBadRequestException("Template has no content")
	}

	data := make(map[string]interface{})
	if in.TemplateData != "" {
		if err := json.Unmarshal([]byte(in.TemplateData), &data); err != nil {
			return nil, awserrors.NewBadRequestException("Invalid TemplateData JSON")
		}
	}

	renderedSubject, err := renderTemplateContent(tmpl.TemplateContent.Subject, data)
	if err != nil {
		return nil, awserrors.NewBadRequestException("Failed to render subject: " + err.Error())
	}

	renderedHtml, err := renderTemplateContent(tmpl.TemplateContent.Html, data)
	if err != nil {
		return nil, awserrors.NewBadRequestException("Failed to render HTML: " + err.Error())
	}

	renderedText, err := renderTemplateContent(tmpl.TemplateContent.Text, data)
	if err != nil {
		return nil, awserrors.NewBadRequestException("Failed to render text: " + err.Error())
	}

	// Pick Content-Type based on what the template actually rendered.
	// AWS prefers HTML when available; text-only templates must declare
	// text/plain so a downstream MIME parser does not mis-handle the body.
	contentType := "text/plain"
	body := renderedText
	if renderedHtml != "" {
		contentType = "text/html"
		body = renderedHtml
	}

	// When the rendered subject is empty, omit the "Subject:" line
	// entirely to avoid producing a malformed MIME preamble
	// ("Subject: \nContent-Type: ...").
	var rendered string
	if renderedSubject != "" {
		rendered = "Subject: " + renderedSubject + "\nContent-Type: " + contentType + "\n\n" + body
	} else {
		rendered = "Content-Type: " + contentType + "\n\n" + body
	}

	return map[string]interface{}{
		"RenderedTemplate": rendered,
	}, nil
}

// ---------------------------------------------------------------------------
// HTTP handlers — parse → DTO → Core → serialise
// ---------------------------------------------------------------------------

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
