package sesv2

import (
	"time"

	"vorpalstacks/internal/store/aws/common"

	"github.com/google/uuid"
)

// Destination specifies the recipients of an email.
type Destination struct {
	ToAddresses  []string `json:"toAddresses,omitempty"`
	CcAddresses  []string `json:"ccAddresses,omitempty"`
	BccAddresses []string `json:"bccAddresses,omitempty"`
}

// Content represents the content of an email body or subject.
type Content struct {
	Data    string `json:"data,omitempty"`
	Charset string `json:"charset,omitempty"`
}

// Body represents the body of an email with text and HTML parts.
type Body struct {
	Text *Content `json:"text,omitempty"`
	Html *Content `json:"html,omitempty"`
}

// Message represents an email message.
type Message struct {
	Subject *Content `json:"subject,omitempty"`
	Body    *Body    `json:"body,omitempty"`
}

// EmailContent represents the content of an email.
type EmailContent struct {
	Simple   *Message    `json:"simple,omitempty"`
	Raw      *RawMessage `json:"raw,omitempty"`
	Template *Template   `json:"template,omitempty"`
}

// RawMessage represents raw email content as bytes.
type RawMessage struct {
	Data []byte `json:"data,omitempty"`
}

// Template represents an email template.
type Template struct {
	TemplateName    string                `json:"templateName,omitempty"`
	TemplateArn     string                `json:"templateArn,omitempty"`
	TemplateData    string                `json:"templateData,omitempty"`
	TemplateContent *EmailTemplateContent `json:"templateContent,omitempty"`
}

// EmailRecord represents a sent email record.
type EmailRecord struct {
	MessageId                      string          `json:"messageId"`
	FromEmailAddress               string          `json:"fromEmailAddress,omitempty"`
	FromEmailAddressIdentityArn    string          `json:"fromEmailAddressIdentityArn,omitempty"`
	Destination                    *Destination    `json:"destination,omitempty"`
	Content                        *EmailContent   `json:"content,omitempty"`
	ReplyToAddresses               []string        `json:"replyToAddresses,omitempty"`
	FeedbackForwardingEmailAddress string          `json:"feedbackForwardingEmailAddress,omitempty"`
	EmailTags                      []MessageTag    `json:"emailTags,omitempty"`
	ConfigurationSetName           string          `json:"configurationSetName,omitempty"`
	ReplacementHeaders             []MessageHeader `json:"replacementHeaders,omitempty"`
	Status                         string          `json:"status"`
	CreatedTimestamp               time.Time       `json:"createdTimestamp"`
	SentTimestamp                  time.Time       `json:"sentTimestamp,omitempty"`
}

// MessageTag represents a tag for an email message.
type MessageTag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// BulkEmailEntry represents an entry in a bulk email send operation.
type BulkEmailEntry struct {
	Destination             *Destination             `json:"destination,omitempty"`
	ReplacementTags         []MessageTag             `json:"replacementTags,omitempty"`
	ReplacementEmailContent *ReplacementEmailContent `json:"replacementEmailContent,omitempty"`
	ReplacementHeaders      []MessageHeader          `json:"replacementHeaders,omitempty"`
}

// ReplacementEmailContent specifies replacement content for bulk email.
type ReplacementEmailContent struct {
	ReplacementTemplate *ReplacementTemplate `json:"replacementTemplate,omitempty"`
}

// ReplacementTemplate specifies template data for bulk email.
type ReplacementTemplate struct {
	ReplacementTemplateData string `json:"replacementTemplateData,omitempty"`
}

// MessageHeader represents a header in an email message.
type MessageHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// BulkEmailEntryResult represents the result of a bulk email entry.
type BulkEmailEntryResult struct {
	Status    string `json:"status"`
	Error     string `json:"error,omitempty"`
	MessageId string `json:"messageId,omitempty"`
}

// NewEmailRecord creates a new email record with default values.
func NewEmailRecord() *EmailRecord {
	now := time.Now().UTC()
	return &EmailRecord{
		MessageId:        generateMessageId(),
		Status:           "QUEUED",
		CreatedTimestamp: now,
	}
}

func generateMessageId() string {
	return uuid.New().String()
}

// SaveEmailRecord saves an email record.
func (s *SESv2Store) SaveEmailRecord(email *EmailRecord) error {
	return s.emailRecordStore.Put(email.MessageId, email)
}

// GetEmailRecord retrieves an email record by message ID.
func (s *SESv2Store) GetEmailRecord(messageId string) (*EmailRecord, error) {
	var email EmailRecord
	if err := s.emailRecordStore.Get(messageId, &email); err != nil {
		return nil, ErrEmailNotFound
	}
	return &email, nil
}

// ListEmailRecords lists email records.
func (s *SESv2Store) ListEmailRecords(opts common.ListOptions) (*common.ListResult[EmailRecord], error) {
	return common.List[EmailRecord](s.emailRecordStore, opts, nil)
}
