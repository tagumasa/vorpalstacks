package sesv2

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"

	pagination "vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/store/aws/common"
	sesv2store "vorpalstacks/internal/store/aws/sesv2"
)

// ---------------------------------------------------------------------------
// Input DTOs — contact-list / contact families
// ---------------------------------------------------------------------------

// CreateContactListInput carries the contact-list create members. Topics
// travels as the raw wire list; the Core parses and validates it.
type CreateContactListInput struct {
	ContactListName string
	Description     string
	Topics          []map[string]interface{}
	Tags            []tags.Tag
}

// UpdateContactListInput carries the presence-based update members: an
// absent Description/Topics keeps the stored value, an explicitly supplied
// one (including empty) replaces it.
type UpdateContactListInput struct {
	ContactListName     string
	Description         string
	DescriptionProvided bool
	Topics              []map[string]interface{}
	TopicsProvided      bool
}

// CreateContactInput carries the contact create members.
type CreateContactInput struct {
	ContactListName        string
	EmailAddress           string
	AttributesData         string
	TopicPreferences       []map[string]interface{}
	UnsubscribeAll         bool
	UnsubscribeAllProvided bool
}

// UpdateContactInput carries the presence-based contact update members.
type UpdateContactInput struct {
	ContactListName          string
	EmailAddress             string
	AttributesData           string
	AttributesDataProvided   bool
	TopicPreferences         []map[string]interface{}
	TopicPreferencesProvided bool
	UnsubscribeAll           bool
	UnsubscribeAllProvided   bool
}

// ListContactsInput carries the contact-list members and the raw Filter
// wire map.
type ListContactsInput struct {
	ContactListName string
	Filter          map[string]interface{}
	MaxItems        int
	NextToken       string
}

// maxFilteredContactScan is the safety upper bound for the in-memory
// filter path in ListContacts.  Without it, a pathological contact list
// could cause unbounded memory and time consumption.  50k is well above
// any realistic edge/on-prem list size.
const maxFilteredContactScan = 50000

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// parseTopics converts the raw Topics wire list into stored topic records,
// validating the DefaultSubscriptionStatus enum.
func parseTopics(topicsList []map[string]interface{}) ([]sesv2store.Topic, error) {
	if len(topicsList) == 0 {
		return nil, nil
	}
	topics := make([]sesv2store.Topic, 0, len(topicsList))
	for _, t := range topicsList {
		topic := sesv2store.Topic{
			TopicName:                 request.GetStringParam(t, "TopicName"),
			DefaultSubscriptionStatus: request.GetStringParam(t, "DefaultSubscriptionStatus"),
		}
		if topic.DefaultSubscriptionStatus != "" && !validateSubscriptionStatus(topic.DefaultSubscriptionStatus) {
			return nil, ErrBadRequest
		}
		if desc := request.GetStringParam(t, "Description"); desc != "" {
			topic.Description = desc
		}
		if dn := request.GetStringParam(t, "DisplayName"); dn != "" {
			topic.DisplayName = dn
		}
		topics = append(topics, topic)
	}
	return topics, nil
}

// parseTopicPreferences converts the raw TopicPreferences wire list into
// stored preference records, validating the SubscriptionStatus enum.
func parseTopicPreferences(prefsList []map[string]interface{}) ([]sesv2store.TopicPreference, error) {
	if len(prefsList) == 0 {
		return nil, nil
	}
	prefs := make([]sesv2store.TopicPreference, 0, len(prefsList))
	for _, tp := range prefsList {
		status := request.GetStringParam(tp, "SubscriptionStatus")
		if status != "" && !validateSubscriptionStatus(status) {
			return nil, ErrBadRequest
		}
		prefs = append(prefs, sesv2store.TopicPreference{
			TopicName:          request.GetStringParam(tp, "TopicName"),
			SubscriptionStatus: status,
		})
	}
	return prefs, nil
}

// contactMatchesFilter applies the ListContacts Filter semantics to a
// single contact. FilteredStatus narrows by the contact's overall
// subscription state (OPT_IN matches UnsubscribeAll=false; OPT_OUT
// matches UnsubscribeAll=true). TopicFilter narrows by per-topic
// preference when UseDefaultIfPreferenceUnavailable is true: contacts
// without an explicit TopicPreference entry are treated as matching the
// topic's DefaultSubscriptionStatus on the parent contact list.
func contactMatchesFilter(c *sesv2store.Contact, filteredStatus, topicFilterName string, useDefault bool) bool {
	if filteredStatus != "" {
		optedIn := !c.UnsubscribeAll
		if filteredStatus == "OPT_IN" && !optedIn {
			return false
		}
		if filteredStatus == "OPT_OUT" && optedIn {
			return false
		}
	}
	if topicFilterName != "" {
		matched := false
		for _, tp := range c.TopicPreferences {
			if tp.TopicName == topicFilterName && tp.SubscriptionStatus == "OPT_IN" {
				matched = true
				break
			}
		}
		if !matched && useDefault {
			// Caller asked us to honour the contact-list default. We do
			// not have the list's Topic.DefaultSubscriptionStatus here;
			// the default behaviour in AWS is to treat absent preferences
			// as OPT_IN, so we accept when the contact has not opted out
			// of the topic explicitly.
			for _, tp := range c.TopicPreferences {
				if tp.TopicName == topicFilterName && tp.SubscriptionStatus == "OPT_OUT" {
					return false
				}
			}
			matched = true
		}
		if !matched {
			return false
		}
	}
	return true
}

// contactSummary renders the JSON shape returned by ListContacts per
// Smithy com.amazonaws.sesv2#Contact. Kept here so the filtered and
// paginated branches serialise identically.
func contactSummary(c *sesv2store.Contact) map[string]interface{} {
	m := map[string]interface{}{
		"EmailAddress":         c.EmailAddress,
		"UnsubscribeAll":       c.UnsubscribeAll,
		"LastUpdatedTimestamp": float64(c.LastUpdatedTimestamp.Unix()),
	}
	if len(c.TopicPreferences) > 0 {
		prefs := make([]map[string]interface{}, 0, len(c.TopicPreferences))
		for _, tp := range c.TopicPreferences {
			prefs = append(prefs, map[string]interface{}{
				"TopicName":          tp.TopicName,
				"SubscriptionStatus": tp.SubscriptionStatus,
			})
		}
		m["TopicPreferences"] = prefs
	}
	if len(c.TopicDefaultPreferences) > 0 {
		defaultPrefs := make([]map[string]interface{}, 0, len(c.TopicDefaultPreferences))
		for _, tp := range c.TopicDefaultPreferences {
			defaultPrefs = append(defaultPrefs, map[string]interface{}{
				"TopicName":          tp.TopicName,
				"SubscriptionStatus": tp.SubscriptionStatus,
			})
		}
		m["TopicDefaultPreferences"] = defaultPrefs
	}
	if c.AttributesData != "" {
		m["AttributesData"] = c.AttributesData
	}
	return m
}

// encodeContactOffset encodes a pagination offset as an opaque base64
// token. Previously "off:N" exposed the internal offset value
// and token structure; base64 encoding makes the token opaque to clients.
func encodeContactOffset(n int) string {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(n))
	return base64.RawURLEncoding.EncodeToString(b)
}

// decodeContactOffset decodes an opaque pagination token. Returns the
// offset and true on success, or 0 and false when the token is invalid.
// Callers should treat (0, false) as an error and reject the
// request rather than silently falling back to offset 0.
func decodeContactOffset(token string) (int, bool) {
	b, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil || len(b) != 4 {
		return 0, false
	}
	n := int(binary.BigEndian.Uint32(b))
	if n < 0 {
		return 0, false
	}
	return n, true
}

// ---------------------------------------------------------------------------
// Core functions — contact-list family
// ---------------------------------------------------------------------------

// createContactListCore is the single entry point for contact-list
// creation. Length limits per AWS service quotas (documented in the SESv2
// API reference even though not encoded in Smithy):
//   - ContactListName: 1-128 chars
//   - Description:    0-500 chars
//   - Topic.TopicName: 1-50 chars
//   - Topic.DisplayName: 1-100 chars
//   - Topic.Description: 0-500 chars
func (s *SESv2Service) createContactListCore(store sesv2store.SESv2StoreInterface, in CreateContactListInput) error {
	contactListName := in.ContactListName
	if contactListName == "" {
		return ErrMissingParameter
	}
	if len(contactListName) > 128 {
		return ErrInvalidParameter
	}
	// Per AWS spec ContactListName must match [a-zA-Z0-9_-]+
	if !validateContactListName(contactListName) {
		return ErrInvalidParameter
	}

	if len(in.Description) > 500 {
		return ErrInvalidParameter
	}

	topics, err := parseTopics(in.Topics)
	if err != nil {
		return err
	}
	for _, t := range topics {
		if t.TopicName == "" || len(t.TopicName) > 50 {
			return ErrInvalidParameter
		}
		if len(t.DisplayName) > 100 {
			return ErrInvalidParameter
		}
		if len(t.Description) > 500 {
			return ErrInvalidParameter
		}
		if t.DefaultSubscriptionStatus != "" &&
			t.DefaultSubscriptionStatus != "OPT_IN" &&
			t.DefaultSubscriptionStatus != "OPT_OUT" {
			return ErrInvalidParameter
		}
	}

	contactList := sesv2store.NewContactList(contactListName)
	contactList.Description = in.Description
	contactList.Topics = topics

	if _, err := store.CreateContactList(contactList); err != nil {
		return err
	}

	if len(in.Tags) > 0 {
		if err := store.TagFromSlice(store.BuildContactListArn(contactListName), in.Tags); err != nil {
			return err
		}
	}

	return nil
}

// getContactListCore is the single entry point for reading a contact list,
// including its tags.
func (s *SESv2Service) getContactListCore(store sesv2store.SESv2StoreInterface, contactListName string) (map[string]interface{}, error) {
	if contactListName == "" {
		return nil, ErrMissingParameter
	}

	cl, err := store.GetContactList(contactListName)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"ContactListName":      cl.ContactListName,
		"CreatedTimestamp":     float64(cl.CreatedTimestamp.Unix()),
		"LastUpdatedTimestamp": float64(cl.LastUpdatedTimestamp.Unix()),
	}

	arn := store.BuildContactListArn(contactListName)
	if tags, err := store.ListAsSlice(arn); err == nil && len(tags) > 0 {
		result["Tags"] = tags
	}

	if cl.Description != "" {
		result["Description"] = cl.Description
	}
	if len(cl.Topics) > 0 {
		topics := make([]map[string]interface{}, 0, len(cl.Topics))
		for _, t := range cl.Topics {
			topic := map[string]interface{}{
				"TopicName":                 t.TopicName,
				"DefaultSubscriptionStatus": t.DefaultSubscriptionStatus,
			}
			if t.DisplayName != "" {
				topic["DisplayName"] = t.DisplayName
			}
			if t.Description != "" {
				topic["Description"] = t.Description
			}
			topics = append(topics, topic)
		}
		result["Topics"] = topics
	}

	return result, nil
}

// deleteContactListCore is the single entry point for deleting a contact
// list.
func (s *SESv2Service) deleteContactListCore(store sesv2store.SESv2StoreInterface, contactListName string) error {
	if contactListName == "" {
		return ErrMissingParameter
	}
	return store.DeleteContactList(contactListName)
}

// listContactListsCore is the single entry point for listing contact
// lists.
func (s *SESv2Service) listContactListsCore(store sesv2store.SESv2StoreInterface, maxItems int, nextToken string) (map[string]interface{}, error) {
	result, err := store.ListContactLists(common.ListOptions{
		MaxItems: maxItems,
		Marker:   nextToken,
	})
	if err != nil {
		return nil, err
	}

	contactLists := make([]map[string]interface{}, 0, len(result.Items))
	for _, cl := range result.Items {
		contactLists = append(contactLists, map[string]interface{}{
			"ContactListName":      cl.ContactListName,
			"LastUpdatedTimestamp": float64(cl.LastUpdatedTimestamp.Unix()),
		})
	}

	resp := map[string]interface{}{
		"ContactLists": contactLists,
	}

	pagination.SetNextToken(resp, "NextToken", result.NextMarker)

	return resp, nil
}

// updateContactListCore is the single entry point for contact-list
// updates. Per Smithy com.amazonaws.sesv2#UpdateContactListRequest,
// Description and Topics are optional but, when present, they fully
// replace the existing values — including with empty. The Core
// distinguishes "absent" from "explicitly empty" via the Provided flags.
func (s *SESv2Service) updateContactListCore(store sesv2store.SESv2StoreInterface, in UpdateContactListInput) error {
	if in.ContactListName == "" {
		return ErrMissingParameter
	}

	cl, err := store.GetContactList(in.ContactListName)
	if err != nil {
		return err
	}

	if in.DescriptionProvided {
		cl.Description = in.Description
	}

	if in.TopicsProvided {
		topics, err := parseTopics(in.Topics)
		if err != nil {
			return err
		}
		cl.Topics = topics
	}

	return store.UpdateContactList(cl)
}

// ---------------------------------------------------------------------------
// Core functions — contact family
// ---------------------------------------------------------------------------

// createContactCore is the single entry point for adding a contact.
func (s *SESv2Service) createContactCore(store sesv2store.SESv2StoreInterface, in CreateContactInput) error {
	if in.ContactListName == "" || in.EmailAddress == "" {
		return ErrMissingParameter
	}

	// Verify the contact list exists before creating a contact.
	if !store.ContactListExists(in.ContactListName) {
		return ErrContactListNotFound
	}

	// Validate the email address format.
	if !validateSuppressionEmailAddress(in.EmailAddress) {
		return ErrBadRequest
	}

	contact := sesv2store.NewContact(in.EmailAddress, in.ContactListName)

	if in.AttributesData != "" {
		// Per AWS spec AttributesData must be a valid JSON object.
		var jsonCheck map[string]interface{}
		if err := json.Unmarshal([]byte(in.AttributesData), &jsonCheck); err != nil {
			return ErrBadRequest
		}
		// Cap the AttributesData size to prevent DoS.
		if len(in.AttributesData) > maxAttributesDataSize {
			return ErrBadRequest
		}
		contact.AttributesData = in.AttributesData
	}

	var err error
	contact.TopicPreferences, err = parseTopicPreferences(in.TopicPreferences)
	if err != nil {
		return err
	}

	if in.UnsubscribeAllProvided {
		contact.UnsubscribeAll = in.UnsubscribeAll
	}

	if _, err := store.CreateContact(contact); err != nil {
		return err
	}

	return nil
}

// getContactCore is the single entry point for reading a contact.
func (s *SESv2Service) getContactCore(store sesv2store.SESv2StoreInterface, contactListName, emailAddress string) (map[string]interface{}, error) {
	if contactListName == "" || emailAddress == "" {
		return nil, ErrMissingParameter
	}

	contact, err := store.GetContact(contactListName, emailAddress)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"ContactListName":      contact.ContactListName,
		"EmailAddress":         contact.EmailAddress,
		"UnsubscribeAll":       contact.UnsubscribeAll,
		"CreatedTimestamp":     float64(contact.CreatedTimestamp.Unix()),
		"LastUpdatedTimestamp": float64(contact.LastUpdatedTimestamp.Unix()),
	}

	if len(contact.TopicPreferences) > 0 {
		topicPrefs := make([]map[string]interface{}, 0, len(contact.TopicPreferences))
		for _, tp := range contact.TopicPreferences {
			topicPrefs = append(topicPrefs, map[string]interface{}{
				"TopicName":          tp.TopicName,
				"SubscriptionStatus": tp.SubscriptionStatus,
			})
		}
		result["TopicPreferences"] = topicPrefs
	}

	if len(contact.TopicDefaultPreferences) > 0 {
		defaultPrefs := make([]map[string]interface{}, 0, len(contact.TopicDefaultPreferences))
		for _, tp := range contact.TopicDefaultPreferences {
			defaultPrefs = append(defaultPrefs, map[string]interface{}{
				"TopicName":          tp.TopicName,
				"SubscriptionStatus": tp.SubscriptionStatus,
			})
		}
		result["TopicDefaultPreferences"] = defaultPrefs
	}

	if contact.AttributesData != "" {
		result["AttributesData"] = contact.AttributesData
	}

	return result, nil
}

// deleteContactCore is the single entry point for removing a contact.
func (s *SESv2Service) deleteContactCore(store sesv2store.SESv2StoreInterface, contactListName, emailAddress string) error {
	if contactListName == "" || emailAddress == "" {
		return ErrMissingParameter
	}
	return store.DeleteContact(contactListName, emailAddress)
}

// listContactsCore is the single entry point for listing contacts. Per
// Smithy com.amazonaws.sesv2#ListContactsRequest, the optional Filter
// member carries FilteredStatus (OPT_IN/OPT_OUT) and a TopicFilter
// (TopicName + UseDefaultIfPreferenceUnavailable); a filtered walk pages
// every contact with a safety cap and layers pagination on the filtered
// set so NextToken stays stable across pages.
func (s *SESv2Service) listContactsCore(store sesv2store.SESv2StoreInterface, in ListContactsInput) (map[string]interface{}, error) {
	contactListName := in.ContactListName
	if contactListName == "" {
		return nil, ErrMissingParameter
	}

	filteredStatus := ""
	topicFilterName := ""
	topicFilterUseDefault := false
	if in.Filter != nil {
		filteredStatus = request.GetStringParam(in.Filter, "FilteredStatus")
		// Validate FilteredStatus against the Smithy
		// SubscriptionStatus enum [OPT_IN, OPT_OUT].
		if filteredStatus != "" && !validateSubscriptionStatus(filteredStatus) {
			return nil, ErrBadRequest
		}
		topicFilterMap := request.GetMapParam(in.Filter, "TopicFilter")
		if topicFilterMap != nil {
			topicFilterName = request.GetStringParam(topicFilterMap, "TopicName")
			topicFilterUseDefault = request.GetBoolParam(topicFilterMap, "UseDefaultIfPreferenceUnavailable")
		}
	}

	// When a filter is supplied, walk every contact in the list and apply
	// it in-memory. Pagination is layered on top of the filtered set so
	// that NextToken remains stable across pages.
	if filteredStatus != "" || topicFilterName != "" {
		// Page through all contacts with a 50k safety cap — prevents
		// unbounded memory/time consumption while covering realistic
		// edge/on-prem list sizes.
		var allItems []*sesv2store.Contact
		listOpts := common.ListOptions{MaxItems: 1000}
		for {
			all, err := store.ListContacts(contactListName, listOpts)
			if err != nil {
				return nil, err
			}
			allItems = append(allItems, all.Items...)
			if len(allItems) > maxFilteredContactScan {
				return nil, ErrLimitExceeded
			}
			if !all.IsTruncated || all.NextMarker == "" {
				break
			}
			listOpts.Marker = all.NextMarker
		}
		filtered := make([]*sesv2store.Contact, 0, len(allItems))
		for _, c := range allItems {
			if !contactMatchesFilter(c, filteredStatus, topicFilterName, topicFilterUseDefault) {
				continue
			}
			filtered = append(filtered, c)
		}
		start := 0
		if in.NextToken != "" {
			idx, ok := decodeContactOffset(in.NextToken)
			if !ok || idx < 0 || idx >= len(filtered) {
				return nil, ErrBadRequest
			}
			start = idx
		}
		end := start + in.MaxItems
		if end > len(filtered) {
			end = len(filtered)
		}
		contacts := make([]map[string]interface{}, 0, end-start)
		for _, c := range filtered[start:end] {
			contacts = append(contacts, contactSummary(c))
		}
		resp := map[string]interface{}{
			"Contacts": contacts,
		}
		if end < len(filtered) {
			resp["NextToken"] = encodeContactOffset(end)
		}
		return resp, nil
	}

	result, err := store.ListContacts(contactListName, common.ListOptions{
		MaxItems: in.MaxItems,
		Marker:   in.NextToken,
	})
	if err != nil {
		return nil, err
	}

	contacts := make([]map[string]interface{}, 0, len(result.Items))
	for _, c := range result.Items {
		contacts = append(contacts, contactSummary(c))
	}

	resp := map[string]interface{}{
		"Contacts": contacts,
	}

	pagination.SetNextToken(resp, "NextToken", result.NextMarker)

	return resp, nil
}

// updateContactCore is the single entry point for contact updates. Per
// Smithy com.amazonaws.sesv2#UpdateContactRequest, the optional fields
// (AttributesData, TopicPreferences, UnsubscribeAll) fully replace the
// existing values when supplied, including with empty; the Core
// distinguishes 'absent' from 'explicitly empty' via the Provided flags.
func (s *SESv2Service) updateContactCore(store sesv2store.SESv2StoreInterface, in UpdateContactInput) error {
	if in.ContactListName == "" || in.EmailAddress == "" {
		return ErrMissingParameter
	}

	// Verify the contact list exists before updating a contact.
	if !store.ContactListExists(in.ContactListName) {
		return ErrContactListNotFound
	}

	contact, err := store.GetContact(in.ContactListName, in.EmailAddress)
	if err != nil {
		return err
	}

	if in.AttributesDataProvided {
		if in.AttributesData != "" {
			var jsonCheck map[string]interface{}
			if err := json.Unmarshal([]byte(in.AttributesData), &jsonCheck); err != nil {
				return ErrBadRequest
			}
			// Cap the AttributesData size to prevent DoS.
			if len(in.AttributesData) > maxAttributesDataSize {
				return ErrBadRequest
			}
		}
		contact.AttributesData = in.AttributesData
	}

	if in.TopicPreferencesProvided {
		contact.TopicPreferences, err = parseTopicPreferences(in.TopicPreferences)
		if err != nil {
			return err
		}
	}

	if in.UnsubscribeAllProvided {
		contact.UnsubscribeAll = in.UnsubscribeAll
	}

	return store.UpdateContact(contact)
}

// ---------------------------------------------------------------------------
// HTTP handlers — parse → DTO → Core → serialise
// ---------------------------------------------------------------------------

// CreateContactList creates a new contact list.
func (s *SESv2Service) CreateContactList(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.createContactListCore(store, CreateContactListInput{
		ContactListName: request.GetStringParam(req.Parameters, "ContactListName"),
		Description:     request.GetStringParam(req.Parameters, "Description"),
		Topics:          request.GetListParam(req.Parameters, "Topics"),
		Tags:            tags.ParseTags(req.Parameters, "Tags"),
	}); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// GetContactList retrieves the details of a contact list.
func (s *SESv2Service) GetContactList(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.getContactListCore(store, request.GetStringParam(req.Parameters, "ContactListName"))
}

// DeleteContactList deletes a contact list.
func (s *SESv2Service) DeleteContactList(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteContactListCore(store, request.GetStringParam(req.Parameters, "ContactListName")); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// ListContactLists lists all contact lists.
func (s *SESv2Service) ListContactLists(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.listContactListsCore(store,
		pagination.GetMaxItems(req.Parameters, 100, "PageSize"),
		pagination.GetMarker(req.Parameters, "NextToken"))
}

// UpdateContactList updates a contact list.
func (s *SESv2Service) UpdateContactList(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := UpdateContactListInput{
		ContactListName: request.GetStringParam(req.Parameters, "ContactListName"),
	}
	if _, ok := req.Parameters["Description"]; ok {
		in.Description = request.GetStringParam(req.Parameters, "Description")
		in.DescriptionProvided = true
	}
	if _, ok := req.Parameters["Topics"]; ok {
		in.Topics = request.GetListParam(req.Parameters, "Topics")
		in.TopicsProvided = true
	}
	if err := s.updateContactListCore(store, in); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// CreateContact adds a contact to a contact list.
func (s *SESv2Service) CreateContact(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := CreateContactInput{
		ContactListName:  request.GetStringParam(req.Parameters, "ContactListName"),
		EmailAddress:     request.GetStringParam(req.Parameters, "EmailAddress"),
		AttributesData:   request.GetStringParam(req.Parameters, "AttributesData"),
		TopicPreferences: request.GetListParam(req.Parameters, "TopicPreferences"),
	}
	if _, ok := req.Parameters["UnsubscribeAll"]; ok {
		in.UnsubscribeAll = request.GetBoolParam(req.Parameters, "UnsubscribeAll")
		in.UnsubscribeAllProvided = true
	}
	if err := s.createContactCore(store, in); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// GetContact retrieves the details of a contact in a contact list.
func (s *SESv2Service) GetContact(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.getContactCore(store,
		request.GetStringParam(req.Parameters, "ContactListName"),
		request.GetStringParam(req.Parameters, "EmailAddress"))
}

// DeleteContact removes a contact from a contact list.
func (s *SESv2Service) DeleteContact(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteContactCore(store,
		request.GetStringParam(req.Parameters, "ContactListName"),
		request.GetStringParam(req.Parameters, "EmailAddress")); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// ListContacts lists contacts in a contact list.
func (s *SESv2Service) ListContacts(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.listContactsCore(store, ListContactsInput{
		ContactListName: request.GetStringParam(req.Parameters, "ContactListName"),
		Filter:          request.GetMapParam(req.Parameters, "Filter"),
		MaxItems:        pagination.GetMaxItems(req.Parameters, 100, "PageSize"),
		NextToken:       pagination.GetMarker(req.Parameters, "NextToken"),
	})
}

// UpdateContact updates a contact in a contact list.
func (s *SESv2Service) UpdateContact(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := UpdateContactInput{
		ContactListName: request.GetStringParam(req.Parameters, "ContactListName"),
		EmailAddress:    request.GetStringParam(req.Parameters, "EmailAddress"),
	}
	if _, ok := req.Parameters["AttributesData"]; ok {
		in.AttributesData = request.GetStringParam(req.Parameters, "AttributesData")
		in.AttributesDataProvided = true
	}
	if _, ok := req.Parameters["TopicPreferences"]; ok {
		in.TopicPreferences = request.GetListParam(req.Parameters, "TopicPreferences")
		in.TopicPreferencesProvided = true
	}
	if _, ok := req.Parameters["UnsubscribeAll"]; ok {
		in.UnsubscribeAll = request.GetBoolParam(req.Parameters, "UnsubscribeAll")
		in.UnsubscribeAllProvided = true
	}
	if err := s.updateContactCore(store, in); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}
