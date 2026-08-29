package sesv2

import (
	"context"
	"net/mail"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	sesv2store "vorpalstacks/internal/store/aws/sesv2"
)

// ---------------------------------------------------------------------------
// Input DTOs — send family
// ---------------------------------------------------------------------------

// SendEmailInput carries every SendEmail member in a wire-independent form:
// flat members as typed scalars, nested structures (Destination, Content,
// ListManagementOptions) as raw wire maps interpreted by the Core.
type SendEmailInput struct {
	FromEmailAddress                          string
	FromEmailAddressIdentityArn               string
	FeedbackForwardingEmailAddress            string
	FeedbackForwardingEmailAddressIdentityArn string
	ConfigurationSetName                      string
	ReplyToAddresses                          []string
	EmailTags                                 []MessageTag
	Destination                               map[string]interface{}
	Content                                   map[string]interface{}
	ListManagementOptions                     map[string]interface{}
}

// SendBulkEmailInput carries every SendBulkEmail member. BulkEmailEntries
// travels as the raw wire list; the Core parses and validates each entry.
type SendBulkEmailInput struct {
	FromEmailAddress                          string
	FromEmailAddressIdentityArn               string
	FeedbackForwardingEmailAddress            string
	FeedbackForwardingEmailAddressIdentityArn string
	ConfigurationSetName                      string
	ReplyToAddresses                          []string
	DefaultEmailTags                          []MessageTag
	DefaultContent                            map[string]interface{}
	ListManagementOptions                     map[string]interface{}
	BulkEmailEntries                          []interface{}
}

// ---------------------------------------------------------------------------
// Wire-shape helpers (nested structures of the send family)
// ---------------------------------------------------------------------------

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

func parseDestination(destMap map[string]interface{}) (*sesv2store.Destination, error) {
	if destMap == nil {
		return nil, nil
	}

	dest := &sesv2store.Destination{}
	if to, ok := destMap["ToAddresses"].([]interface{}); ok {
		dest.ToAddresses = toStringSlice(to)
		for _, a := range dest.ToAddresses {
			if !validateEmailAddress(a) {
				return nil, ErrBadRequest
			}
		}
	}
	if cc, ok := destMap["CcAddresses"].([]interface{}); ok {
		dest.CcAddresses = toStringSlice(cc)
		for _, a := range dest.CcAddresses {
			if !validateEmailAddress(a) {
				return nil, ErrBadRequest
			}
		}
	}
	if bcc, ok := destMap["BccAddresses"].([]interface{}); ok {
		dest.BccAddresses = toStringSlice(bcc)
		for _, a := range dest.BccAddresses {
			if !validateEmailAddress(a) {
				return nil, ErrBadRequest
			}
		}
	}

	// A Destination must contain at least one recipient across To/Cc/Bcc.
	if len(dest.ToAddresses) == 0 && len(dest.CcAddresses) == 0 && len(dest.BccAddresses) == 0 {
		return nil, ErrBadRequest
	}

	return dest, nil
}

func parseContent(contentMap map[string]interface{}) *sesv2store.EmailContent {
	if contentMap == nil {
		return nil
	}

	content := &sesv2store.EmailContent{}

	// Per AWS docs EmailContent is a union — at most one of Simple, Raw,
	// or Template may be specified. Accept the first and ignore the rest
	// to avoid ambiguous content states.
	if simple, ok := contentMap["Simple"].(map[string]interface{}); ok {
		content.Simple = parseMessage(simple)
		return content
	}

	if raw, ok := contentMap["Raw"].(map[string]interface{}); ok {
		content.Raw = parseRawMessage(raw)
		return content
	}

	if tmpl, ok := contentMap["Template"].(map[string]interface{}); ok {
		content.Template = parseTemplate(tmpl)
		return content
	}

	// No recognised content member — treat as absent.
	return nil
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

func parseBulkEmailEntries(entryList []interface{}) ([]sesv2store.BulkEmailEntry, error) {
	var entries []sesv2store.BulkEmailEntry

	for _, e := range entryList {
		entryMap, ok := e.(map[string]interface{})
		if !ok {
			return nil, ErrBadRequest
		}

		entry := sesv2store.BulkEmailEntry{}
		// Per Smithy, BulkEmailEntry.Destination is @required.
		dest, destOk := entryMap["Destination"].(map[string]interface{})
		if !destOk {
			return nil, ErrBadRequest
		}
		parsed, err := parseDestination(dest)
		if err != nil {
			return nil, err
		}
		entry.Destination = parsed

		if tags, ok := entryMap["ReplacementTags"].([]interface{}); ok {
			tagsRaw := ParseMessageTagsFromList(tags)
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
// handled correctly via net/mail.ParseAddressList.
func extractIdentityFromEmail(email string) (string, bool) {
	if email == "" {
		return "", false
	}
	// Use net/mail.ParseAddressList for robust RFC 5322 parsing
	// instead of fragile manual angle-bracket extraction. Handles all
	// display-name variants correctly without edge-case breakage.
	if addrs, err := mail.ParseAddressList(email); err == nil && len(addrs) > 0 {
		email = addrs[0].Address
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
//
// AWS SES treats identity names as case-insensitive. The lookup is
// normalised to lowercase to match identities stored in any case.
func identityExistsForEmail(store sesv2store.SESv2StoreInterface, email string) bool {
	email = strings.ToLower(email)
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

// ---------------------------------------------------------------------------
// Core functions — send family
// ---------------------------------------------------------------------------

// sendEmailCore is the single entry point for the single-send path:
// address validation, destination/content parsing and validation,
// configuration-set existence, From-identity authorisation, and record
// persistence.
func (s *SESv2Service) sendEmailCore(store sesv2store.SESv2StoreInterface, in SendEmailInput) (map[string]interface{}, error) {
	fromEmailAddress := in.FromEmailAddress
	if fromEmailAddress == "" {
		return nil, ErrMissingParameter
	}
	if !validateFromAddress(fromEmailAddress) {
		return nil, ErrBadRequest
	}

	listMgmtOpts := parseListManagementOptions(in.ListManagementOptions)

	destMap := in.Destination
	if destMap == nil {
		return nil, ErrBadRequest
	}
	destination, err := parseDestination(destMap)
	if err != nil {
		return nil, err
	}

	content := parseContent(in.Content)

	// Content is Smithy @required — reject when absent or empty.
	if content == nil {
		return nil, ErrMissingParameter
	}
	// Reject content that has no meaningful payload (e.g. Simple with
	// empty Body, Raw with empty Data, Template with empty TemplateName).
	if err := validateContent(content); err != nil {
		return nil, ErrMessageRejected
	}
	// Reject raw messages exceeding the 10 MB AWS SES size limit.
	if content.Raw != nil && len(content.Raw.Data) > maxRawMessageSize {
		return nil, ErrMessageRejected
	}

	// Verify configuration-set existence when supplied.
	if in.ConfigurationSetName != "" {
		if !store.ConfigurationSetExists(in.ConfigurationSetName) {
			return nil, ErrConfigurationSetNotFound
		}
	}

	if !identityExistsForEmail(store, fromEmailAddress) {
		return nil, ErrIdentityNotFound
	}

	email := sesv2store.NewEmailRecord()
	email.FromEmailAddress = fromEmailAddress
	email.FromEmailAddressIdentityArn = in.FromEmailAddressIdentityArn
	email.Destination = destination
	email.Content = content
	email.ConfigurationSetName = in.ConfigurationSetName
	email.FeedbackForwardingEmailAddress = in.FeedbackForwardingEmailAddress
	email.FeedbackForwardingEmailAddressIdentityArn = in.FeedbackForwardingEmailAddressIdentityArn
	email.ListManagementOptions = listMgmtOpts
	email.Status = "SENT"
	email.SentTimestamp = time.Now().UTC()

	// Validate ReplyToAddresses email format.
	if len(in.ReplyToAddresses) > 0 {
		if err := validateReplyToAddresses(in.ReplyToAddresses); err != nil {
			return nil, ErrBadRequest
		}
		email.ReplyToAddresses = in.ReplyToAddresses
	}

	// Validate EmailTag Name/Value charset and length.
	if len(in.EmailTags) > 0 {
		tags := make([]sesv2store.MessageTag, len(in.EmailTags))
		for i, t := range in.EmailTags {
			tags[i] = sesv2store.MessageTag{Name: t.Name, Value: t.Value}
		}
		if err := validateMessageTags(tags); err != nil {
			return nil, ErrBadRequest
		}
		email.EmailTags = tags
	}

	if err := store.SaveEmailRecord(email); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"MessageId": email.MessageId,
	}, nil
}

// sendBulkEmailCore is the single entry point for the bulk-send path. Per
// AWS, FromEmailAddress and at least one BulkEmailEntries item are required;
// BulkEmailEntries is capped at 50 items.
func (s *SESv2Service) sendBulkEmailCore(store sesv2store.SESv2StoreInterface, in SendBulkEmailInput) (map[string]interface{}, error) {
	fromEmailAddress := in.FromEmailAddress
	if fromEmailAddress == "" {
		return nil, ErrMissingParameter
	}
	if !validateFromAddress(fromEmailAddress) {
		return nil, ErrBadRequest
	}

	// DefaultEmailTags are applied to every entry in the bulk send,
	// merged with per-entry ReplacementTags.
	var defaultEmailTags []sesv2store.MessageTag
	if len(in.DefaultEmailTags) > 0 {
		defaultEmailTags = make([]sesv2store.MessageTag, len(in.DefaultEmailTags))
		for i, t := range in.DefaultEmailTags {
			defaultEmailTags[i] = sesv2store.MessageTag{Name: t.Name, Value: t.Value}
		}
	}

	defaultContent := parseContent(in.DefaultContent)

	// DefaultContent is Smithy @required — reject when absent or empty.
	if defaultContent == nil {
		return nil, ErrMissingParameter
	}
	if err := validateContent(defaultContent); err != nil {
		return nil, ErrMessageRejected
	}

	entries, err := parseBulkEmailEntries(in.BulkEmailEntries)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, ErrMissingParameter
	}

	// AWS limits BulkEmailEntries to 50 items.
	if len(entries) > 50 {
		return nil, ErrLimitExceeded
	}

	// Verify configuration-set existence when supplied.
	if in.ConfigurationSetName != "" {
		if !store.ConfigurationSetExists(in.ConfigurationSetName) {
			return nil, ErrConfigurationSetNotFound
		}
	}

	// Validate ReplyToAddresses format.
	if len(in.ReplyToAddresses) > 0 {
		if err := validateReplyToAddresses(in.ReplyToAddresses); err != nil {
			return nil, ErrBadRequest
		}
	}

	if !identityExistsForEmail(store, fromEmailAddress) {
		return nil, ErrIdentityNotFound
	}

	listMgmtOpts := parseListManagementOptions(in.ListManagementOptions)

	results := make([]map[string]interface{}, 0, len(entries))
	for _, entry := range entries {
		email := sesv2store.NewEmailRecord()
		email.FromEmailAddress = fromEmailAddress
		email.FromEmailAddressIdentityArn = in.FromEmailAddressIdentityArn
		email.Destination = entry.Destination
		email.ConfigurationSetName = in.ConfigurationSetName
		email.FeedbackForwardingEmailAddress = in.FeedbackForwardingEmailAddress
		email.FeedbackForwardingEmailAddressIdentityArn = in.FeedbackForwardingEmailAddressIdentityArn
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
		var mergedTags []sesv2store.MessageTag
		if len(defaultEmailTags) > 0 {
			mergedTags = append(mergedTags, defaultEmailTags...)
		}
		if len(entry.ReplacementTags) > 0 {
			mergedTags = append(mergedTags, entry.ReplacementTags...)
		}
		// Validate merged tag set (count + charset + length).
		if len(mergedTags) > 0 {
			if err := validateMessageTags(mergedTags); err != nil {
				return nil, ErrBadRequest
			}
			email.EmailTags = mergedTags
		}

		if len(in.ReplyToAddresses) > 0 {
			email.ReplyToAddresses = in.ReplyToAddresses
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

// ---------------------------------------------------------------------------
// HTTP handlers — parse → DTO → Core → serialise
// ---------------------------------------------------------------------------

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
