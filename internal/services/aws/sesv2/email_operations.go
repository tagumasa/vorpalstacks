package sesv2

import (
	"context"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
	sesv2store "vorpalstacks/internal/store/aws/sesv2"
)

// SendEmail sends an email using the SESv2 service.
func (s *SESv2Service) SendEmail(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	fromEmailAddress := request.GetStringParam(req.Parameters, "FromEmailAddress")
	if fromEmailAddress == "" {
		return nil, ErrMissingParameter
	}
	configurationSetName := request.GetStringParam(req.Parameters, "ConfigurationSetName")
	feedbackForwardingEmailAddress := request.GetStringParam(req.Parameters, "FeedbackForwardingEmailAddress")

	destMap := request.GetMapParam(req.Parameters, "Destination")
	destination := parseDestination(destMap)

	contentMap := request.GetMapParam(req.Parameters, "Content")
	content := parseContent(contentMap)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !identityExistsForEmail(store.Raw(), fromEmailAddress) {
		return nil, ErrIdentityNotFound
	}

	email := sesv2store.NewEmailRecord()
	email.FromEmailAddress = fromEmailAddress
	email.Destination = destination
	email.Content = content
	email.ConfigurationSetName = configurationSetName
	email.FeedbackForwardingEmailAddress = feedbackForwardingEmailAddress
	email.Status = "SENT"
	email.SentTimestamp = time.Now().UTC()

	replyTo := request.GetStringList(req.Parameters, "ReplyToAddresses")
	if len(replyTo) > 0 {
		email.ReplyToAddresses = replyTo
	}

	emailTagsRaw := tagutil.ParseMessageTags(req.Parameters, "EmailTags")
	if len(emailTagsRaw) > 0 {
		email.EmailTags = make([]sesv2store.MessageTag, len(emailTagsRaw))
		for i, t := range emailTagsRaw {
			email.EmailTags[i] = sesv2store.MessageTag{Name: t.Name, Value: t.Value}
		}
	}

	if err := store.SaveEmailRecord(email); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"MessageId": email.MessageId,
	}, nil
}

// SendBulkEmail sends multiple emails in a single operation.
func (s *SESv2Service) SendBulkEmail(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	fromEmailAddress := request.GetStringParam(req.Parameters, "FromEmailAddress")
	configurationSetName := request.GetStringParam(req.Parameters, "ConfigurationSetName")
	feedbackForwardingEmailAddress := request.GetStringParam(req.Parameters, "FeedbackForwardingEmailAddress")

	defaultContent := parseContent(request.GetMapParam(req.Parameters, "DefaultContent"))

	entries := parseBulkEmailEntries(req.Parameters)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if fromEmailAddress != "" {
		if !identityExistsForEmail(store.Raw(), fromEmailAddress) {
			return nil, ErrIdentityNotFound
		}
	}

	results := make([]map[string]interface{}, 0, len(entries))
	for _, entry := range entries {
		email := sesv2store.NewEmailRecord()
		email.FromEmailAddress = fromEmailAddress
		email.Destination = entry.Destination
		email.ConfigurationSetName = configurationSetName
		email.FeedbackForwardingEmailAddress = feedbackForwardingEmailAddress
		email.Status = "SENT"
		email.SentTimestamp = time.Now().UTC()

		email.Content = defaultContent
		if entry.ReplacementEmailContent != nil && entry.ReplacementEmailContent.ReplacementTemplate != nil {
			replTemplate := entry.ReplacementEmailContent.ReplacementTemplate
			if email.Content == nil {
				email.Content = &sesv2store.EmailContent{}
			}
			if email.Content.Template == nil {
				email.Content.Template = &sesv2store.Template{}
			}
			if replTemplate.ReplacementTemplateData != "" {
				email.Content.Template.TemplateData = replTemplate.ReplacementTemplateData
			}
		}

		if len(entry.ReplacementTags) > 0 {
			email.EmailTags = entry.ReplacementTags
		}

		if len(entry.ReplacementHeaders) > 0 {
			email.ReplacementHeaders = entry.ReplacementHeaders
		}

		if err := store.SaveEmailRecord(email); err != nil {
			results = append(results, map[string]interface{}{
				"Status": "FAILED",
				"Error":  err.Error(),
			})
			continue
		}

		results = append(results, map[string]interface{}{
			"Status":    "SUCCESS",
			"MessageId": email.MessageId,
		})
	}

	return map[string]interface{}{
		"BulkEmailEntryResults": results,
	}, nil
}

func parseDestination(destMap map[string]interface{}) *sesv2store.Destination {
	if destMap == nil {
		return nil
	}

	dest := &sesv2store.Destination{}
	if to, ok := destMap["ToAddresses"].([]interface{}); ok {
		dest.ToAddresses = toStringSlice(to)
	}
	if cc, ok := destMap["CcAddresses"].([]interface{}); ok {
		dest.CcAddresses = toStringSlice(cc)
	}
	if bcc, ok := destMap["BccAddresses"].([]interface{}); ok {
		dest.BccAddresses = toStringSlice(bcc)
	}

	return dest
}

func parseContent(contentMap map[string]interface{}) *sesv2store.EmailContent {
	if contentMap == nil {
		return nil
	}

	content := &sesv2store.EmailContent{}

	if simple, ok := contentMap["Simple"].(map[string]interface{}); ok {
		content.Simple = parseMessage(simple)
	}

	if raw, ok := contentMap["Raw"].(map[string]interface{}); ok {
		content.Raw = parseRawMessage(raw)
	}

	if tmpl, ok := contentMap["Template"].(map[string]interface{}); ok {
		content.Template = parseTemplate(tmpl)
	}

	return content
}

func parseMessage(msgMap map[string]interface{}) *sesv2store.Message {
	if msgMap == nil {
		return nil
	}

	msg := &sesv2store.Message{}

	if subject, ok := msgMap["Subject"].(map[string]interface{}); ok {
		msg.Subject = parseContent_(subject)
	}

	if body, ok := msgMap["Body"].(map[string]interface{}); ok {
		msg.Body = parseBody(body)
	}

	return msg
}

func parseContent_(contentMap map[string]interface{}) *sesv2store.Content {
	if contentMap == nil {
		return nil
	}
	return &sesv2store.Content{
		Data:    request.GetStringParam(contentMap, "Data"),
		Charset: request.GetStringParam(contentMap, "Charset"),
	}
}

func parseBody(bodyMap map[string]interface{}) *sesv2store.Body {
	if bodyMap == nil {
		return nil
	}

	body := &sesv2store.Body{}

	if text, ok := bodyMap["Text"].(map[string]interface{}); ok {
		body.Text = parseContent_(text)
	}

	if html, ok := bodyMap["Html"].(map[string]interface{}); ok {
		body.Html = parseContent_(html)
	}

	return body
}

func parseRawMessage(rawMap map[string]interface{}) *sesv2store.RawMessage {
	if rawMap == nil {
		return nil
	}
	return &sesv2store.RawMessage{
		Data: []byte(request.GetStringParam(rawMap, "Data")),
	}
}

func parseTemplate(tmplMap map[string]interface{}) *sesv2store.Template {
	if tmplMap == nil {
		return nil
	}

	tmpl := &sesv2store.Template{
		TemplateName: request.GetStringParam(tmplMap, "TemplateName"),
		TemplateArn:  request.GetStringParam(tmplMap, "TemplateArn"),
		TemplateData: request.GetStringParam(tmplMap, "TemplateData"),
	}

	if content, ok := tmplMap["TemplateContent"].(map[string]interface{}); ok {
		tmpl.TemplateContent = &sesv2store.EmailTemplateContent{
			Subject: request.GetStringParam(content, "Subject"),
			Html:    request.GetStringParam(content, "Html"),
			Text:    request.GetStringParam(content, "Text"),
		}
	}

	return tmpl
}

func parseBulkEmailEntries(params map[string]interface{}) []sesv2store.BulkEmailEntry {
	var entries []sesv2store.BulkEmailEntry

	entriesIf, ok := params["BulkEmailEntries"]
	if !ok {
		return entries
	}

	entryList, ok := entriesIf.([]interface{})
	if !ok {
		return entries
	}

	for _, e := range entryList {
		entryMap, ok := e.(map[string]interface{})
		if !ok {
			continue
		}

		entry := sesv2store.BulkEmailEntry{}
		if dest, ok := entryMap["Destination"].(map[string]interface{}); ok {
			entry.Destination = parseDestination(dest)
		}

		if tags, ok := entryMap["ReplacementTags"].([]interface{}); ok {
			tagsRaw := tagutil.ParseMessageTagsFromList(tags)
			entry.ReplacementTags = make([]sesv2store.MessageTag, len(tagsRaw))
			for i, t := range tagsRaw {
				entry.ReplacementTags[i] = sesv2store.MessageTag{Name: t.Name, Value: t.Value}
			}
		}

		if replContent, ok := entryMap["ReplacementEmailContent"].(map[string]interface{}); ok {
			entry.ReplacementEmailContent = parseReplacementEmailContent(replContent)
		}

		if replHeaders, ok := entryMap["ReplacementHeaders"].([]interface{}); ok {
			entry.ReplacementHeaders = parseReplacementHeaders(replHeaders)
		}

		entries = append(entries, entry)
	}

	return entries
}

func toStringSlice(iface []interface{}) []string {
	result := make([]string, 0, len(iface))
	for _, v := range iface {
		if s, ok := v.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

func extractIdentityFromEmail(email string) string {
	if !strings.Contains(email, "@") {
		return email
	}
	parts := strings.Split(email, "@")
	return parts[len(parts)-1]
}

func identityExistsForEmail(store *sesv2store.SESv2Store, email string) bool {
	_, err := store.GetEmailIdentity(email)
	if err != nil {
		domain := extractIdentityFromEmail(email)
		if domain != email {
			_, err = store.GetEmailIdentity(domain)
			return err == nil
		}
		return false
	}
	return true
}

func parseReplacementEmailContent(replMap map[string]interface{}) *sesv2store.ReplacementEmailContent {
	if replMap == nil {
		return nil
	}

	repl := &sesv2store.ReplacementEmailContent{}
	if replTemplate, ok := replMap["ReplacementTemplate"].(map[string]interface{}); ok {
		repl.ReplacementTemplate = &sesv2store.ReplacementTemplate{
			ReplacementTemplateData: request.GetStringParam(replTemplate, "ReplacementTemplateData"),
		}
	}

	return repl
}

func parseReplacementHeaders(headers []interface{}) []sesv2store.MessageHeader {
	if len(headers) == 0 {
		return nil
	}

	result := make([]sesv2store.MessageHeader, 0, len(headers))
	for _, h := range headers {
		if headerMap, ok := h.(map[string]interface{}); ok {
			header := sesv2store.MessageHeader{
				Name:  request.GetStringParam(headerMap, "Name"),
				Value: request.GetStringParam(headerMap, "Value"),
			}
			result = append(result, header)
		}
	}

	return result
}
