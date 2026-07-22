package sesv2

import (
	"bytes"
	"context"
	"encoding/json"
	"regexp"
	"strings"
	"text/template"

	awserrors "vorpalstacks/internal/common/errors"
	pagination "vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/store/aws/common"
	sesv2store "vorpalstacks/internal/store/aws/sesv2"
)

var (
	awsVarPattern     = regexp.MustCompile(`\{\{([^#/^}{]+)\}\}`)
	awsThisPattern    = regexp.MustCompile(`\{\{\s*this\s*\}\}`)
	awsThisFieldStart = regexp.MustCompile(`\{\{\s*this\.`)
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

// CreateEmailTemplate creates a new email template.
// Per AWS, a template must contain at least one of Subject, Html, or Text;
// we replicate that validation rather than persisting a contentless
// template that TestRenderEmailTemplate would later fail on.
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

	if tmpl.TemplateContent == nil ||
		(tmpl.TemplateContent.Subject == "" && tmpl.TemplateContent.Html == "" && tmpl.TemplateContent.Text == "") {
		return nil, ErrInvalidParameter
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
		if err := store.TagFromSlice(arn, parsedTags); err != nil {
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

	if tmpl.TemplateContent == nil ||
		(tmpl.TemplateContent.Subject == "" && tmpl.TemplateContent.Html == "" && tmpl.TemplateContent.Text == "") {
		return nil, ErrInvalidParameter
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
	pageSize := pagination.GetMaxItems(req.Parameters, 100, "PageSize")
	nextToken := pagination.GetMarker(req.Parameters, "NextToken")

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
		entry := map[string]interface{}{
			"TemplateName": tmpl.TemplateName,
		}
		// Per Smithy com.amazonaws.sesv2#EmailTemplateMetadata the list
		// response carries CreatedTimestamp alongside TemplateName. Older
		// records persisted before the field was added serialise as the
		// zero time; only emit the field when it carries useful data.
		if !tmpl.CreatedTimestamp.IsZero() {
			entry["CreatedTimestamp"] = float64(tmpl.CreatedTimestamp.Unix())
		}
		templates = append(templates, entry)
	}

	response := map[string]interface{}{
		"TemplatesMetadata": templates,
	}

	pagination.SetNextToken(response, "NextToken", result.NextMarker)

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
		return nil, awserrors.NewBadRequestException("Template has no content")
	}

	data := make(map[string]interface{})
	if tmplData != "" {
		if err := json.Unmarshal([]byte(tmplData), &data); err != nil {
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

	rendered := "Subject: " + renderedSubject + "\nContent-Type: " + contentType + "\n\n" + body

	return map[string]interface{}{
		"RenderedTemplate": rendered,
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
