// Package sesv2 provides SES v2 service operations for vorpalstacks.
package sesv2

import (
	"sync"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	sesv2store "vorpalstacks/internal/store/aws/sesv2"
)

// SESv2Service provides SES v2 email service operations.
type SESv2Service struct {
	accountID string
	stores    sync.Map // region → sesv2store.SESv2StoreInterface
}

// NewSESv2Service creates a new SES v2 service instance.
func NewSESv2Service(accountID string) *SESv2Service {
	return &SESv2Service{
		accountID: accountID,
	}
}

func (s *SESv2Service) store(reqCtx *request.RequestContext) (sesv2store.SESv2StoreInterface, error) {
	if store := reqCtx.GetSESv2Store(); store != nil {
		return store, nil
	}
	region := reqCtx.GetRegion()
	if cached, ok := s.stores.Load(region); ok {
		if typed, ok := cached.(sesv2store.SESv2StoreInterface); ok {
			return typed, nil
		}
	}
	storage, err := reqCtx.GetStorage()
	if err != nil {
		return nil, err
	}
	store := sesv2store.NewSESv2Store(storage, s.accountID, region)
	if actual, loaded := s.stores.LoadOrStore(region, store); loaded {
		if typed, ok := actual.(sesv2store.SESv2StoreInterface); ok {
			return typed, nil
		}
	}
	return store, nil
}

// RegisterHandlers registers the SES v2 service handlers with the dispatcher.
func (s *SESv2Service) RegisterHandlers(d handler.Registrar) {
	// Email sending operations
	d.RegisterHandlerForService("email", "SendEmail", s.SendEmail)
	d.RegisterHandlerForService("email", "SendBulkEmail", s.SendBulkEmail)

	// Email identity operations
	d.RegisterHandlerForService("email", "CreateEmailIdentity", s.CreateEmailIdentity)
	d.RegisterHandlerForService("email", "GetEmailIdentity", s.GetEmailIdentity)
	d.RegisterHandlerForService("email", "DeleteEmailIdentity", s.DeleteEmailIdentity)
	d.RegisterHandlerForService("email", "ListEmailIdentities", s.ListEmailIdentities)
	d.RegisterHandlerForService("email", "PutEmailIdentityDkimAttributes", s.PutEmailIdentityDkimAttributes)
	d.RegisterHandlerForService("email", "PutEmailIdentityDkimSigningAttributes", s.PutEmailIdentityDkimSigningAttributes)
	d.RegisterHandlerForService("email", "PutEmailIdentityFeedbackAttributes", s.PutEmailIdentityFeedbackAttributes)
	d.RegisterHandlerForService("email", "PutEmailIdentityMailFromAttributes", s.PutEmailIdentityMailFromAttributes)
	d.RegisterHandlerForService("email", "PutEmailIdentityConfigurationSetAttributes", s.PutEmailIdentityConfigurationSetAttributes)
	d.RegisterHandlerForService("email", "GetEmailIdentityPolicies", s.GetEmailIdentityPolicies)
	d.RegisterHandlerForService("email", "CreateEmailIdentityPolicy", s.CreateEmailIdentityPolicy)
	d.RegisterHandlerForService("email", "UpdateEmailIdentityPolicy", s.UpdateEmailIdentityPolicy)
	d.RegisterHandlerForService("email", "DeleteEmailIdentityPolicy", s.DeleteEmailIdentityPolicy)

	// Configuration set operations
	d.RegisterHandlerForService("email", "CreateConfigurationSet", s.CreateConfigurationSet)
	d.RegisterHandlerForService("email", "GetConfigurationSet", s.GetConfigurationSet)
	d.RegisterHandlerForService("email", "DeleteConfigurationSet", s.DeleteConfigurationSet)
	d.RegisterHandlerForService("email", "ListConfigurationSets", s.ListConfigurationSets)
	d.RegisterHandlerForService("email", "CreateConfigurationSetEventDestination", s.CreateConfigurationSetEventDestination)
	d.RegisterHandlerForService("email", "GetConfigurationSetEventDestinations", s.GetConfigurationSetEventDestinations)
	d.RegisterHandlerForService("email", "UpdateConfigurationSetEventDestination", s.UpdateConfigurationSetEventDestination)
	d.RegisterHandlerForService("email", "DeleteConfigurationSetEventDestination", s.DeleteConfigurationSetEventDestination)
	d.RegisterHandlerForService("email", "PutConfigurationSetDeliveryOptions", s.PutConfigurationSetDeliveryOptions)
	d.RegisterHandlerForService("email", "PutConfigurationSetReputationOptions", s.PutConfigurationSetReputationOptions)
	d.RegisterHandlerForService("email", "PutConfigurationSetSendingOptions", s.PutConfigurationSetSendingOptions)
	d.RegisterHandlerForService("email", "PutConfigurationSetSuppressionOptions", s.PutConfigurationSetSuppressionOptions)
	d.RegisterHandlerForService("email", "PutConfigurationSetTrackingOptions", s.PutConfigurationSetTrackingOptions)
	d.RegisterHandlerForService("email", "PutConfigurationSetVdmOptions", s.PutConfigurationSetVdmOptions)
	d.RegisterHandlerForService("email", "PutConfigurationSetArchivingOptions", s.PutConfigurationSetArchivingOptions)

	// Template operations
	d.RegisterHandlerForService("email", "CreateEmailTemplate", s.CreateEmailTemplate)
	d.RegisterHandlerForService("email", "GetEmailTemplate", s.GetEmailTemplate)
	d.RegisterHandlerForService("email", "UpdateEmailTemplate", s.UpdateEmailTemplate)
	d.RegisterHandlerForService("email", "DeleteEmailTemplate", s.DeleteEmailTemplate)
	d.RegisterHandlerForService("email", "ListEmailTemplates", s.ListEmailTemplates)
	d.RegisterHandlerForService("email", "TestRenderEmailTemplate", s.TestRenderEmailTemplate)

	// Tag operations
	d.RegisterHandlerForService("email", "TagResource", s.TagResource)
	d.RegisterHandlerForService("email", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("email", "ListTagsForResource", s.ListTagsForResource)

	// Account operations
	d.RegisterHandlerForService("email", "GetAccount", s.GetAccount)
	d.RegisterHandlerForService("email", "PutAccountSendingAttributes", s.PutAccountSendingAttributes)
	d.RegisterHandlerForService("email", "PutAccountSuppressionAttributes", s.PutAccountSuppressionAttributes)
	d.RegisterHandlerForService("email", "PutAccountDetails", s.PutAccountDetails)
	d.RegisterHandlerForService("email", "PutAccountVdmAttributes", s.PutAccountVdmAttributes)
	d.RegisterHandlerForService("email", "PutAccountDedicatedIpWarmupAttributes", s.PutAccountDedicatedIpWarmupAttributes)

	// Dedicated IP pool operations
	d.RegisterHandlerForService("email", "CreateDedicatedIpPool", s.CreateDedicatedIpPool)
	d.RegisterHandlerForService("email", "DeleteDedicatedIpPool", s.DeleteDedicatedIpPool)
	d.RegisterHandlerForService("email", "GetDedicatedIpPool", s.GetDedicatedIpPool)
	d.RegisterHandlerForService("email", "ListDedicatedIpPools", s.ListDedicatedIpPools)

	// Suppression list operations
	d.RegisterHandlerForService("email", "GetSuppressedDestination", s.GetSuppressedDestination)
	d.RegisterHandlerForService("email", "PutSuppressedDestination", s.PutSuppressedDestination)
	d.RegisterHandlerForService("email", "DeleteSuppressedDestination", s.DeleteSuppressedDestination)
	d.RegisterHandlerForService("email", "ListSuppressedDestinations", s.ListSuppressedDestinations)

	// Contact list operations
	d.RegisterHandlerForService("email", "CreateContactList", s.CreateContactList)
	d.RegisterHandlerForService("email", "GetContactList", s.GetContactList)
	d.RegisterHandlerForService("email", "DeleteContactList", s.DeleteContactList)
	d.RegisterHandlerForService("email", "ListContactLists", s.ListContactLists)
	d.RegisterHandlerForService("email", "UpdateContactList", s.UpdateContactList)
	d.RegisterHandlerForService("email", "CreateContact", s.CreateContact)
	d.RegisterHandlerForService("email", "GetContact", s.GetContact)
	d.RegisterHandlerForService("email", "DeleteContact", s.DeleteContact)
	d.RegisterHandlerForService("email", "ListContacts", s.ListContacts)
	d.RegisterHandlerForService("email", "UpdateContact", s.UpdateContact)
}
