package sesv2

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
)

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
