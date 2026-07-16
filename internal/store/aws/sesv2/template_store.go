package sesv2

import (
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/store/aws/common"
)

// EmailTemplateContent represents the content of an email template.
type EmailTemplateContent struct {
	Subject string `json:"subject,omitempty"`
	Html    string `json:"html,omitempty"`
	Text    string `json:"text,omitempty"`
}

// EmailTemplate represents an email template.
type EmailTemplate struct {
	TemplateName         string                `json:"templateName"`
	TemplateContent      *EmailTemplateContent `json:"templateContent,omitempty"`
	CreatedTimestamp     time.Time             `json:"createdTimestamp"`
	LastUpdatedTimestamp time.Time             `json:"lastUpdatedTimestamp"`
}

// NewEmailTemplate creates a new email template.
func NewEmailTemplate(name string) *EmailTemplate {
	now := time.Now().UTC()
	return &EmailTemplate{
		TemplateName:         name,
		TemplateContent:      &EmailTemplateContent{},
		CreatedTimestamp:     now,
		LastUpdatedTimestamp: now,
	}
}

// CreateEmailTemplate creates a new email template.
func (s *SESv2Store) CreateEmailTemplate(template *EmailTemplate) (*EmailTemplate, error) {
	if s.templateStore.Exists(template.TemplateName) {
		return nil, ErrTemplateAlreadyExists
	}

	template.CreatedTimestamp = time.Now().UTC()
	template.LastUpdatedTimestamp = template.CreatedTimestamp

	if err := s.templateStore.Put(template.TemplateName, template); err != nil {
		return nil, err
	}

	return template, nil
}

// GetEmailTemplate retrieves an email template.
func (s *SESv2Store) GetEmailTemplate(name string) (*EmailTemplate, error) {
	var template EmailTemplate
	if err := s.templateStore.Get(name, &template); err != nil {
		return nil, ErrTemplateNotFound
	}
	return &template, nil
}

// DeleteEmailTemplate deletes an email template.
func (s *SESv2Store) DeleteEmailTemplate(name string) error {
	if !s.templateStore.Exists(name) {
		return ErrTemplateNotFound
	}
	arn := s.BuildTemplateArn(name)
	if err := s.TagStore.Delete(arn); err != nil {
		logs.Error("Failed to delete tags for template", logs.String("name", name), logs.Err(err))
	}
	return s.templateStore.Delete(name)
}

// UpdateEmailTemplate updates an email template.
func (s *SESv2Store) UpdateEmailTemplate(template *EmailTemplate) error {
	if !s.templateStore.Exists(template.TemplateName) {
		return ErrTemplateNotFound
	}
	template.LastUpdatedTimestamp = time.Now().UTC()
	return s.templateStore.Put(template.TemplateName, template)
}

// ListEmailTemplates lists email templates.
func (s *SESv2Store) ListEmailTemplates(opts common.ListOptions) (*common.ListResult[EmailTemplate], error) {
	return common.List[EmailTemplate](s.templateStore, opts, nil)
}

// EmailTemplateExists checks if an email template exists.
func (s *SESv2Store) EmailTemplateExists(name string) bool {
	return s.templateStore.Exists(name)
}
