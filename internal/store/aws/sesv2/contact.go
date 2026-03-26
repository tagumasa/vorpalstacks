package sesv2

import (
	"log"
	"time"

	"vorpalstacks/internal/store/aws/common"
)

// ContactListStore provides storage operations for SESv2 contact lists.
type ContactListStore struct {
	*common.BaseStore
}

// ContactList represents an SESv2 contact list with its topics and metadata.
type ContactList struct {
	ContactListName      string            `json:"contactListName"`
	Description          string            `json:"description,omitempty"`
	Topics               []Topic           `json:"topics,omitempty"`
	Tags                 map[string]string `json:"tags,omitempty"`
	CreatedTimestamp     time.Time         `json:"createdTimestamp"`
	LastUpdatedTimestamp time.Time         `json:"lastUpdatedTimestamp"`
}

// Topic represents a topic within an SESv2 contact list.
type Topic struct {
	TopicName                 string `json:"topicName"`
	DisplayName               string `json:"displayName,omitempty"`
	Description               string `json:"description,omitempty"`
	DefaultSubscriptionStatus string `json:"defaultSubscriptionStatus,omitempty"`
}

// ContactStore provides storage operations for SESv2 contacts.
type ContactStore struct {
	*common.BaseStore
}

// Contact represents an email address in an SESv2 contact list.
type Contact struct {
	EmailAddress            string            `json:"emailAddress"`
	ContactListName         string            `json:"contactListName"`
	TopicPreferences        []TopicPreference `json:"topicPreferences,omitempty"`
	TopicDefaultPreferences []TopicPreference `json:"topicDefaultPreferences,omitempty"`
	UnsubscribeAll          bool              `json:"unsubscribeAll"`
	AttributesData          string            `json:"attributesData,omitempty"`
	CreatedTimestamp        time.Time         `json:"createdTimestamp"`
	LastUpdatedTimestamp    time.Time         `json:"lastUpdatedTimestamp"`
}

// TopicPreference represents a contact's subscription preference for a topic.
type TopicPreference struct {
	TopicName          string `json:"topicName"`
	SubscriptionStatus string `json:"subscriptionStatus"`
}

// NewContactList creates a new ContactList with the given name and initialised timestamps.
func NewContactList(name string) *ContactList {
	now := time.Now().UTC()
	return &ContactList{
		ContactListName:      name,
		CreatedTimestamp:     now,
		LastUpdatedTimestamp: now,
	}
}

// CreateContactList creates a new SESv2 contact list.
func (s *SESv2Store) CreateContactList(cl *ContactList) (*ContactList, error) {
	if s.contactListStore.Exists(cl.ContactListName) {
		return nil, ErrContactListAlreadyExists
	}
	cl.CreatedTimestamp = time.Now().UTC()
	cl.LastUpdatedTimestamp = cl.CreatedTimestamp
	if err := s.contactListStore.Put(cl.ContactListName, cl); err != nil {
		return nil, err
	}
	return cl, nil
}

// GetContactList retrieves an SESv2 contact list by name.
func (s *SESv2Store) GetContactList(name string) (*ContactList, error) {
	var cl ContactList
	if err := s.contactListStore.Get(name, &cl); err != nil {
		return nil, ErrContactListNotFound
	}
	return &cl, nil
}

// DeleteContactList deletes an SESv2 contact list and all its associated contacts and tags.
func (s *SESv2Store) DeleteContactList(name string) error {
	if !s.contactListStore.Exists(name) {
		return ErrContactListNotFound
	}
	prefix := "contact#" + name + "#"
	result, err := common.List[Contact](s.contactStore, common.ListOptions{Prefix: prefix, MaxItems: 10000}, nil)
	if err == nil {
		for _, c := range result.Items {
			_ = s.contactStore.Delete("contact#" + c.ContactListName + "#" + c.EmailAddress)
		}
	}
	arn := s.arnBuilder.Build("ses", "contact-list/"+name)
	if err := s.TagStore.Delete(arn); err != nil {
		log.Printf("Failed to delete tags for contact list %s: %v", name, err)
	}
	return s.contactListStore.Delete(name)
}

// UpdateContactList updates an existing SESv2 contact list.
func (s *SESv2Store) UpdateContactList(cl *ContactList) error {
	if !s.contactListStore.Exists(cl.ContactListName) {
		return ErrContactListNotFound
	}
	cl.LastUpdatedTimestamp = time.Now().UTC()
	return s.contactListStore.Put(cl.ContactListName, cl)
}

// ListContactLists returns all SESv2 contact lists using the specified list options.
func (s *SESv2Store) ListContactLists(opts common.ListOptions) (*common.ListResult[ContactList], error) {
	return common.List[ContactList](s.contactListStore, opts, nil)
}

// NewContact creates a new Contact with the given email address and contact list name.
func NewContact(emailAddress, contactListName string) *Contact {
	now := time.Now().UTC()
	return &Contact{
		EmailAddress:         emailAddress,
		ContactListName:      contactListName,
		CreatedTimestamp:     now,
		LastUpdatedTimestamp: now,
	}
}

// CreateContact creates a new SESv2 contact in a contact list.
func (s *SESv2Store) CreateContact(c *Contact) (*Contact, error) {
	key := "contact#" + c.ContactListName + "#" + c.EmailAddress
	if s.contactStore.Exists(key) {
		return nil, ErrContactAlreadyExists
	}
	c.CreatedTimestamp = time.Now().UTC()
	c.LastUpdatedTimestamp = c.CreatedTimestamp
	if err := s.contactStore.Put(key, c); err != nil {
		return nil, err
	}
	return c, nil
}

// GetContact retrieves an SESv2 contact by contact list name and email address.
func (s *SESv2Store) GetContact(contactListName, emailAddress string) (*Contact, error) {
	key := "contact#" + contactListName + "#" + emailAddress
	var c Contact
	if err := s.contactStore.Get(key, &c); err != nil {
		return nil, ErrContactNotFound
	}
	return &c, nil
}

// DeleteContact removes an SESv2 contact from a contact list.
func (s *SESv2Store) DeleteContact(contactListName, emailAddress string) error {
	key := "contact#" + contactListName + "#" + emailAddress
	if !s.contactStore.Exists(key) {
		return ErrContactNotFound
	}
	return s.contactStore.Delete(key)
}

// UpdateContact updates an existing SESv2 contact's preferences and attributes.
func (s *SESv2Store) UpdateContact(c *Contact) error {
	key := "contact#" + c.ContactListName + "#" + c.EmailAddress
	if !s.contactStore.Exists(key) {
		return ErrContactNotFound
	}
	c.LastUpdatedTimestamp = time.Now().UTC()
	return s.contactStore.Put(key, c)
}

// ListContacts returns all contacts in an SESv2 contact list using the specified list options.
func (s *SESv2Store) ListContacts(contactListName string, opts common.ListOptions) (*common.ListResult[Contact], error) {
	prefix := "contact#" + contactListName + "#"
	return common.List[Contact](s.contactStore, common.ListOptions{
		Prefix:   prefix,
		Marker:   opts.Marker,
		MaxItems: opts.MaxItems,
	}, nil)
}
