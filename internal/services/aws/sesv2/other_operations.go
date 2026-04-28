package sesv2

import (
	"context"
	"fmt"
	"strconv"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	pagination "vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/store/aws/common"
	sesv2store "vorpalstacks/internal/store/aws/sesv2"
)

func parseStoredTime(s string) float64 {
	if v, err := strconv.ParseInt(s, 10, 64); err == nil {
		return float64(v)
	}
	return float64(time.Now().UTC().Unix())
}

// CreateDedicatedIpPool creates a new dedicated IP pool.
func (s *SESv2Service) CreateDedicatedIpPool(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	poolName := request.GetStringParam(req.Parameters, "PoolName")
	if poolName == "" {
		return nil, ErrMissingParameter
	}

	scalingMode := request.GetStringParam(req.Parameters, "ScalingMode")
	if scalingMode == "" {
		scalingMode = "MANAGED"
	}

	pool := &sesv2store.DedicatedIpPool{
		PoolName:    poolName,
		ScalingMode: scalingMode,
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.CreateDedicatedIpPool(pool); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DeleteDedicatedIpPool deletes a dedicated IP pool.
func (s *SESv2Service) DeleteDedicatedIpPool(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	poolName := request.GetStringParam(req.Parameters, "PoolName")
	if poolName == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteDedicatedIpPool(poolName); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetDedicatedIpPool retrieves the details of a dedicated IP pool.
func (s *SESv2Service) GetDedicatedIpPool(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	poolName := request.GetStringParam(req.Parameters, "PoolName")
	if poolName == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	pool, err := store.GetDedicatedIpPool(poolName)
	if err != nil {
		return nil, awserrors.NewNotFoundException("Dedicated IP pool not found")
	}

	return map[string]interface{}{
		"DedicatedIpPool": map[string]interface{}{
			"PoolName":    pool.PoolName,
			"ScalingMode": pool.ScalingMode,
		},
	}, nil
}

// ListDedicatedIpPools returns a list of dedicated IP pools.
func (s *SESv2Service) ListDedicatedIpPools(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	result, err := store.ListDedicatedIpPools(opts)
	if err != nil {
		return nil, err
	}

	pools := make([]string, 0, len(result.Items))
	for _, pool := range result.Items {
		pools = append(pools, pool.PoolName)
	}

	response := map[string]interface{}{
		"DedicatedIpPools": pools,
	}

	pagination.SetNextToken(response, "NextToken", result.NextMarker)

	return response, nil
}

// GetSuppressedDestination retrieves details about a suppressed email destination.
func (s *SESv2Service) GetSuppressedDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	emailAddress := request.GetStringParam(req.Parameters, "EmailAddress")
	if emailAddress == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	dest, err := store.GetSuppressedDestination(emailAddress)
	if err != nil {
		if common.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return map[string]interface{}{
		"SuppressedDestination": map[string]interface{}{
			"EmailAddress":   dest.EmailAddress,
			"Reason":         dest.Reason,
			"LastUpdateTime": parseStoredTime(dest.LastUpdateTime),
		},
	}, nil
}

// PutSuppressedDestination adds or updates a suppressed destination.
func (s *SESv2Service) PutSuppressedDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	emailAddress := request.GetStringParam(req.Parameters, "EmailAddress")
	reason := request.GetStringParam(req.Parameters, "Reason")

	if emailAddress == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	dest := &sesv2store.SuppressedDestination{
		EmailAddress:      emailAddress,
		Reason:            reason,
		SuppressionReason: reason,
		LastUpdateTime:    fmt.Sprintf("%d", time.Now().UTC().Unix()),
	}

	if err := store.PutSuppressedDestination(dest); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DeleteSuppressedDestination removes an email address from the suppression list.
func (s *SESv2Service) DeleteSuppressedDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	emailAddress := request.GetStringParam(req.Parameters, "EmailAddress")
	if emailAddress == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteSuppressedDestination(emailAddress); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListSuppressedDestinations returns a list of suppressed destinations.
func (s *SESv2Service) ListSuppressedDestinations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	result, err := store.ListSuppressedDestinations(opts)
	if err != nil {
		return nil, err
	}

	destinations := make([]map[string]interface{}, 0, len(result.Items))
	for _, dest := range result.Items {
		destinations = append(destinations, map[string]interface{}{
			"EmailAddress":   dest.EmailAddress,
			"Reason":         dest.Reason,
			"LastUpdateTime": parseStoredTime(dest.LastUpdateTime),
		})
	}

	response := map[string]interface{}{
		"SuppressedDestinationSummaries": destinations,
	}

	pagination.SetNextToken(response, "NextToken", result.NextMarker)

	return response, nil
}

// CreateContactList creates a new contact list.
func (s *SESv2Service) CreateContactList(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	contactListName := request.GetStringParam(req.Parameters, "ContactListName")
	if contactListName == "" {
		return nil, ErrMissingParameter
	}

	description := request.GetStringParam(req.Parameters, "Description")
	parsedTags := tags.ParseTags(req.Parameters, "Tags")

	contactList := sesv2store.NewContactList(contactListName)
	contactList.Description = description

	if topics := request.GetStringList(req.Parameters, "Topics"); len(topics) > 0 {
		for _, t := range topics {
			contactList.Topics = append(contactList.Topics, sesv2store.Topic{
				TopicName:                 t,
				DefaultSubscriptionStatus: "OPT_IN",
			})
		}
	}

	cl, err := store.CreateContactList(contactList)
	if err != nil {
		return nil, err
	}

	if len(parsedTags) > 0 {
		if err := store.TagFromSlice("ses:contact-list/"+contactListName, parsedTags); err != nil {
			return nil, err
		}
	}

	_ = cl

	return response.EmptyResponse(), nil
}

// GetContactList retrieves the details of a contact list.
func (s *SESv2Service) GetContactList(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	contactListName := request.GetStringParam(req.Parameters, "ContactListName")
	if contactListName == "" {
		return nil, ErrMissingParameter
	}

	cl, err := store.GetContactList(contactListName)
	if err != nil {
		return nil, awserrors.NewNotFoundException("Contact list not found")
	}

	result := map[string]interface{}{
		"ContactListName":      cl.ContactListName,
		"LastUpdatedTimestamp": float64(cl.LastUpdatedTimestamp.Unix()),
	}

	if cl.Description != "" {
		result["Description"] = cl.Description
	}
	if len(cl.Topics) > 0 {
		topics := make([]map[string]interface{}, 0, len(cl.Topics))
		for _, t := range cl.Topics {
			topics = append(topics, map[string]interface{}{
				"TopicName":                 t.TopicName,
				"DefaultSubscriptionStatus": t.DefaultSubscriptionStatus,
			})
		}
		result["Topics"] = topics
	}

	return result, nil
}

// DeleteContactList deletes a contact list.
func (s *SESv2Service) DeleteContactList(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	contactListName := request.GetStringParam(req.Parameters, "ContactListName")
	if contactListName == "" {
		return nil, ErrMissingParameter
	}

	if err := store.DeleteContactList(contactListName); err != nil {
		return nil, awserrors.NewNotFoundException("Contact list not found")
	}

	return response.EmptyResponse(), nil
}

// ListContactLists lists all contact lists.
func (s *SESv2Service) ListContactLists(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	pageSize := pagination.GetMaxItems(req.Parameters, 100, "PageSize")
	nextToken := pagination.GetMarker(req.Parameters, "NextToken")

	opts := common.ListOptions{
		MaxItems: pageSize,
		Marker:   nextToken,
	}

	result, err := store.ListContactLists(opts)
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

// UpdateContactList updates a contact list.
func (s *SESv2Service) UpdateContactList(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	contactListName := request.GetStringParam(req.Parameters, "ContactListName")
	if contactListName == "" {
		return nil, ErrMissingParameter
	}

	cl, err := store.GetContactList(contactListName)
	if err != nil {
		return nil, awserrors.NewNotFoundException("Contact list not found")
	}

	if description := request.GetStringParam(req.Parameters, "Description"); description != "" {
		cl.Description = description
	}

	if topics := request.GetStringList(req.Parameters, "Topics"); len(topics) > 0 {
		newTopics := make([]sesv2store.Topic, 0, len(topics))
		for _, t := range topics {
			newTopics = append(newTopics, sesv2store.Topic{
				TopicName:                 t,
				DefaultSubscriptionStatus: "OPT_IN",
			})
		}
		cl.Topics = newTopics
	}

	if err := store.UpdateContactList(cl); err != nil {
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

	contactListName := request.GetStringParam(req.Parameters, "ContactListName")
	emailAddress := request.GetStringParam(req.Parameters, "EmailAddress")

	if contactListName == "" || emailAddress == "" {
		return nil, ErrMissingParameter
	}

	contact := sesv2store.NewContact(emailAddress, contactListName)

	if attrs := request.GetStringParam(req.Parameters, "AttributesData"); attrs != "" {
		contact.AttributesData = attrs
	}

	if topicPrefs := request.GetStringList(req.Parameters, "TopicPreferences"); len(topicPrefs) > 0 {
		for _, tp := range topicPrefs {
			contact.TopicPreferences = append(contact.TopicPreferences, sesv2store.TopicPreference{
				TopicName:          tp,
				SubscriptionStatus: "OPT_IN",
			})
		}
	}

	if unsSubAll := request.GetBoolParam(req.Parameters, "UnsubscribeAll"); unsSubAll {
		contact.UnsubscribeAll = true
	}

	if _, err := store.CreateContact(contact); err != nil {
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

	contactListName := request.GetStringParam(req.Parameters, "ContactListName")
	emailAddress := request.GetStringParam(req.Parameters, "EmailAddress")

	if contactListName == "" || emailAddress == "" {
		return nil, ErrMissingParameter
	}

	contact, err := store.GetContact(contactListName, emailAddress)
	if err != nil {
		return nil, awserrors.NewNotFoundException("Contact not found")
	}

	result := map[string]interface{}{
		"EmailAddress":         contact.EmailAddress,
		"UnsubscribeAll":       contact.UnsubscribeAll,
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

	if contact.AttributesData != "" {
		result["AttributesData"] = contact.AttributesData
	}

	return result, nil
}

// DeleteContact removes a contact from a contact list.
func (s *SESv2Service) DeleteContact(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	contactListName := request.GetStringParam(req.Parameters, "ContactListName")
	emailAddress := request.GetStringParam(req.Parameters, "EmailAddress")

	if contactListName == "" || emailAddress == "" {
		return nil, ErrMissingParameter
	}

	if err := store.DeleteContact(contactListName, emailAddress); err != nil {
		return nil, awserrors.NewNotFoundException("Contact not found")
	}

	return response.EmptyResponse(), nil
}

// ListContacts lists contacts in a contact list.
func (s *SESv2Service) ListContacts(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	contactListName := request.GetStringParam(req.Parameters, "ContactListName")
	if contactListName == "" {
		return nil, ErrMissingParameter
	}

	pageSize := pagination.GetMaxItems(req.Parameters, 100, "PageSize")
	nextToken := pagination.GetMarker(req.Parameters, "NextToken")

	opts := common.ListOptions{
		MaxItems: pageSize,
		Marker:   nextToken,
	}

	result, err := store.ListContacts(contactListName, opts)
	if err != nil {
		return nil, err
	}

	contacts := make([]map[string]interface{}, 0, len(result.Items))
	for _, c := range result.Items {
		contacts = append(contacts, map[string]interface{}{
			"EmailAddress":         c.EmailAddress,
			"UnsubscribeAll":       c.UnsubscribeAll,
			"LastUpdatedTimestamp": float64(c.LastUpdatedTimestamp.Unix()),
		})
	}

	resp := map[string]interface{}{
		"Contacts": contacts,
	}

	pagination.SetNextToken(resp, "NextToken", result.NextMarker)

	return resp, nil
}

// UpdateContact updates a contact in a contact list.
func (s *SESv2Service) UpdateContact(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	contactListName := request.GetStringParam(req.Parameters, "ContactListName")
	emailAddress := request.GetStringParam(req.Parameters, "EmailAddress")

	if contactListName == "" || emailAddress == "" {
		return nil, ErrMissingParameter
	}

	contact, err := store.GetContact(contactListName, emailAddress)
	if err != nil {
		return nil, awserrors.NewNotFoundException("Contact not found")
	}

	if attrs := request.GetStringParam(req.Parameters, "AttributesData"); attrs != "" {
		contact.AttributesData = attrs
	}

	if topicPrefs := request.GetStringList(req.Parameters, "TopicPreferences"); len(topicPrefs) > 0 {
		newPrefs := make([]sesv2store.TopicPreference, 0, len(topicPrefs))
		for _, tp := range topicPrefs {
			newPrefs = append(newPrefs, sesv2store.TopicPreference{
				TopicName:          tp,
				SubscriptionStatus: "OPT_IN",
			})
		}
		contact.TopicPreferences = newPrefs
	}

	if unsSubAll := request.GetBoolParam(req.Parameters, "UnsubscribeAll"); unsSubAll {
		contact.UnsubscribeAll = true
	}

	if err := store.UpdateContact(contact); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
