package sesv2

import (
	"context"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
	sesv2store "vorpalstacks/internal/store/aws/sesv2"
)

// validateFromAddress validates that fromEmailAddress is a syntactically
// valid email address or domain.  This is a lightweight check (presence
// of "@" for email, non-empty for domain) — not a full RFC 5322 parser.
// The goal is to reject obviously malformed input early rather than
// relying on the identity-lookup to fail.
func validateFromAddress(addr string) bool {
	if addr == "" {
		return false
	}
	// Domain identity: no "@" — accept as long as non-empty.
	if !strings.Contains(addr, "@") {
		return true
	}
	// Email address: must have exactly one "@" with non-empty local
	// and domain parts.
	parts := strings.Split(addr, "@")
	if len(parts) != 2 {
		return false
	}
	return parts[0] != "" && parts[1] != ""
}

// rejectTenantName returns ErrBadRequest when the caller supplies
// TenantName.  Tenant management (48 ops) is not implemented in this
// edge/on-prem build; silently ignoring TenantName would return
// account-level results, which is misleading.  Fail-closed instead.
func rejectTenantName(params map[string]interface{}) error {
	if tn := request.GetStringParam(params, "TenantName"); tn != "" {
		return ErrBadRequest
	}
	return nil
}

// parseListManagementOptions extracts the ListManagementOptions from
// request parameters.  Returns nil when ContactListName is absent.
func parseListManagementOptions(params map[string]interface{}) *sesv2store.ListManagementOptions {
	lmm := request.GetMapParam(params, "ListManagementOptions")
	if lmm == nil {
		return nil
	}
	contactListName := request.GetStringParam(lmm, "ContactListName")
	if contactListName == "" {
		return nil
	}
	return &sesv2store.ListManagementOptions{
		ContactListName: contactListName,
		TopicName:       request.GetStringParam(lmm, "TopicName"),
	}
}

// SendEmail sends an email using the SESv2 service.
func (s *SESv2Service) SendEmail(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := rejectTenantName(req.Parameters); err != nil {
		return nil, err
	}

	fromEmailAddress := request.GetStringParam(req.Parameters, "FromEmailAddress")
	if fromEmailAddress == "" {
		return nil, ErrMissingParameter
	}
	if !validateFromAddress(fromEmailAddress) {
		return nil, ErrBadRequest
	}

	configurationSetName := request.GetStringParam(req.Parameters, "ConfigurationSetName")
	feedbackForwardingEmailAddress := request.GetStringParam(req.Parameters, "FeedbackForwardingEmailAddress")
	fromEmailAddressIdentityArn := request.GetStringParam(req.Parameters, "FromEmailAddressIdentityArn")
	feedbackForwardingEmailAddressIdentityArn := request.GetStringParam(req.Parameters, "FeedbackForwardingEmailAddressIdentityArn")
	listMgmtOpts := parseListManagementOptions(req.Parameters)

	destMap := request.GetMapParam(req.Parameters, "Destination")
	destination, err := parseDestination(destMap)
	if err != nil {
		return nil, err
	}

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
	email.FromEmailAddressIdentityArn = fromEmailAddressIdentityArn
	email.Destination = destination
	email.Content = content
	email.ConfigurationSetName = configurationSetName
	email.FeedbackForwardingEmailAddress = feedbackForwardingEmailAddress
	email.FeedbackForwardingEmailAddressIdentityArn = feedbackForwardingEmailAddressIdentityArn
	email.ListManagementOptions = listMgmtOpts
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
// Per AWS, FromEmailAddress and at least one BulkEmailEntries item are
// required. The previous implementation accepted an empty From and an
// empty Entries slice, returning a success response with zero entries.
func (s *SESv2Service) SendBulkEmail(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := rejectTenantName(req.Parameters); err != nil {
		return nil, err
	}

	fromEmailAddress := request.GetStringParam(req.Parameters, "FromEmailAddress")
	if fromEmailAddress == "" {
		return nil, ErrMissingParameter
	}
	if !validateFromAddress(fromEmailAddress) {
		return nil, ErrBadRequest
	}
	configurationSetName := request.GetStringParam(req.Parameters, "ConfigurationSetName")
	feedbackForwardingEmailAddress := request.GetStringParam(req.Parameters, "FeedbackForwardingEmailAddress")
	fromEmailAddressIdentityArn := request.GetStringParam(req.Parameters, "FromEmailAddressIdentityArn")
	feedbackForwardingEmailAddressIdentityArn := request.GetStringParam(req.Parameters, "FeedbackForwardingEmailAddressIdentityArn")
	listMgmtOpts := parseListManagementOptions(req.Parameters)

	// DefaultEmailTags are applied to every entry in the bulk send,
	// merged with per-entry ReplacementTags.
	var defaultEmailTags []sesv2store.MessageTag
	defaultTagsRaw := tagutil.ParseMessageTags(req.Parameters, "DefaultEmailTags")
	if len(defaultTagsRaw) > 0 {
		defaultEmailTags = make([]sesv2store.MessageTag, len(defaultTagsRaw))
		for i, t := range defaultTagsRaw {
			defaultEmailTags[i] = sesv2store.MessageTag{Name: t.Name, Value: t.Value}
		}
	}

	defaultContent := parseContent(request.GetMapParam(req.Parameters, "DefaultContent"))

	entries, err := parseBulkEmailEntries(req.Parameters)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if !identityExistsForEmail(store.Raw(), fromEmailAddress) {
		return nil, ErrIdentityNotFound
	}

	results := make([]map[string]interface{}, 0, len(entries))
	for _, entry := range entries {
		email := sesv2store.NewEmailRecord()
		email.FromEmailAddress = fromEmailAddress
		email.FromEmailAddressIdentityArn = fromEmailAddressIdentityArn
		email.Destination = entry.Destination
		email.ConfigurationSetName = configurationSetName
		email.FeedbackForwardingEmailAddress = feedbackForwardingEmailAddress
		email.FeedbackForwardingEmailAddressIdentityArn = feedbackForwardingEmailAddressIdentityArn
		email.ListManagementOptions = listMgmtOpts
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

		// Merge DefaultEmailTags with per-entry ReplacementTags.
		if len(defaultEmailTags) > 0 {
			email.EmailTags = append(email.EmailTags, defaultEmailTags...)
		}
		if len(entry.ReplacementTags) > 0 {
			email.EmailTags = append(email.EmailTags, entry.ReplacementTags...)
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

// isValidEmailAddress validates that addr is a syntactically valid email
// address (must contain exactly one "@" with non-empty local and domain
// parts).  Used for recipient addresses in Destination — unlike sender
// identities, recipients cannot be bare domains.
func isValidEmailAddress(addr string) bool {
	if addr == "" {
		return false
	}
	parts := strings.Split(addr, "@")
	if len(parts) != 2 {
		return false
	}
	return parts[0] != "" && parts[1] != ""
}

func parseDestination(destMap map[string]interface{}) (*sesv2store.Destination, error) {
	if destMap == nil {
		return nil, nil
	}

	dest := &sesv2store.Destination{}
	if to, ok := destMap["ToAddresses"].([]interface{}); ok {
		dest.ToAddresses = toStringSlice(to)
		for _, a := range dest.ToAddresses {
			if !isValidEmailAddress(a) {
				return nil, ErrBadRequest
			}
		}
	}
	if cc, ok := destMap["CcAddresses"].([]interface{}); ok {
		dest.CcAddresses = toStringSlice(cc)
		for _, a := range dest.CcAddresses {
			if !isValidEmailAddress(a) {
				return nil, ErrBadRequest
			}
		}
	}
	if bcc, ok := destMap["BccAddresses"].([]interface{}); ok {
		dest.BccAddresses = toStringSlice(bcc)
		for _, a := range dest.BccAddresses {
			if !isValidEmailAddress(a) {
				return nil, ErrBadRequest
			}
		}
	}

	return dest, nil
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
		msg.Subject = parseContentPart(subject)
	}

	if body, ok := msgMap["Body"].(map[string]interface{}); ok {
		msg.Body = parseBody(body)
	}

	return msg
}

func parseContentPart(contentMap map[string]interface{}) *sesv2store.Content {
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
		body.Text = parseContentPart(text)
	}

	if html, ok := bodyMap["Html"].(map[string]interface{}); ok {
		body.Html = parseContentPart(html)
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

func parseBulkEmailEntries(params map[string]interface{}) ([]sesv2store.BulkEmailEntry, error) {
	var entries []sesv2store.BulkEmailEntry

	entriesIf, ok := params["BulkEmailEntries"]
	if !ok {
		return entries, nil
	}

	entryList, ok := entriesIf.([]interface{})
	if !ok {
		return entries, nil
	}

	for _, e := range entryList {
		entryMap, ok := e.(map[string]interface{})
		if !ok {
			continue
		}

		entry := sesv2store.BulkEmailEntry{}
		if dest, ok := entryMap["Destination"].(map[string]interface{}); ok {
			parsed, err := parseDestination(dest)
			if err != nil {
				return nil, err
			}
			entry.Destination = parsed
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

	return entries, nil
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

// extractIdentityFromEmail returns the identity string that should be
// looked up to authorise the FromEmailAddress. Per AWS, identities are
// either an exact email address (user@example.com) or the bare domain
// (example.com). RFC 5322 display-name forms such as 'John Doe <u@d>' are
// not valid SES identities and are rejected upstream by the SDK before
// reaching the service; we still try to extract the addr-spec when
// possible so a slightly-lenient client does not get a confusing
// 'identity not found' error.
func extractIdentityFromEmail(email string) (string, bool) {
	if email == "" {
		return "", false
	}
	// Trim display-name form: 'Name <addr@domain>' -> 'addr@domain'
	if i := strings.LastIndex(email, "<"); i >= 0 {
		end := strings.Index(email[i:], ">")
		if end > 0 {
			email = email[i+1 : i+end]
		}
	}
	email = strings.TrimSpace(email)
	if !strings.Contains(email, "@") {
		// Bare domain (e.g. 'example.com') is a valid identity form.
		return email, true
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", false
	}
	return parts[len(parts)-1], true
}

// identityExistsForEmail reports whether a verified identity exists that
// authorises sending from the given FromEmailAddress. The lookup is
// performed in two steps: first the exact email address, then the
// extracted domain. Malformed inputs return false so the caller returns
// ErrIdentityNotFound rather than a confusing 'no identity matches'.
func identityExistsForEmail(store *sesv2store.SESv2Store, email string) bool {
	if _, err := store.GetEmailIdentity(email); err == nil {
		return true
	}
	domain, ok := extractIdentityFromEmail(email)
	if !ok || domain == email {
		return false
	}
	_, err := store.GetEmailIdentity(domain)
	return err == nil
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
