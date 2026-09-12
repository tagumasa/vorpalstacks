package cognitoidentityprovider

import (
	"errors"
	"testing"
)

// A pool without a stored log-delivery configuration or UI customisation is
// not an error: the getters answer nil for absent records, while storage
// failures propagate to the caller.
func TestOptionalConfigurationsAbsentIsNilNotError(t *testing.T) {
	s := newUserPoolTestStore(t)
	cfg, err := s.GetLogDeliveryConfiguration("nosuch-pool")
	if err != nil || cfg != nil {
		t.Fatalf("absent log delivery configuration = (%v, %v), want (nil, nil)", cfg, err)
	}
	ui, err := s.GetUICustomization("nosuch-pool", "")
	if err != nil || ui != nil {
		t.Fatalf("absent UI customisation = (%v, %v), want (nil, nil)", ui, err)
	}
}

// An app client holds at most one branding style: the save guard refuses a
// second record for the same client even when the caller bypasses the
// service-layer pre-check.
func TestSaveManagedLoginBrandingOnePerClient(t *testing.T) {
	s := newUserPoolTestStore(t)
	first := &ManagedLoginBranding{ManagedLoginBrandingId: "branding-1", UserPoolID: "pool", ClientID: "client-a"}
	if err := s.SaveManagedLoginBranding(first); err != nil {
		t.Fatalf("first save: %v", err)
	}
	second := &ManagedLoginBranding{ManagedLoginBrandingId: "branding-2", UserPoolID: "pool", ClientID: "client-a"}
	if err := s.SaveManagedLoginBranding(second); !errors.Is(err, ErrManagedLoginBrandingExists) {
		t.Fatalf("second branding for one client = %v, want ErrManagedLoginBrandingExists", err)
	}
	// A re-save of the same record (the update path) still passes.
	if err := s.SaveManagedLoginBranding(first); err != nil {
		t.Fatalf("re-save of the existing record: %v", err)
	}
}

// A terms name is unique to the app client: the save guard refuses a second
// document of the same name for one client, while another client may hold
// the same name.
func TestSaveTermsNameUniquePerClient(t *testing.T) {
	s := newUserPoolTestStore(t)
	first := &Terms{TermsID: "terms-1", UserPoolID: "pool", ClientID: "client-a", TermsName: "default"}
	if err := s.SaveTerms(first); err != nil {
		t.Fatalf("first save: %v", err)
	}
	dup := &Terms{TermsID: "terms-2", UserPoolID: "pool", ClientID: "client-a", TermsName: "default"}
	if err := s.SaveTerms(dup); !errors.Is(err, ErrTermsExists) {
		t.Fatalf("duplicate name for one client = %v, want ErrTermsExists", err)
	}
	other := &Terms{TermsID: "terms-3", UserPoolID: "pool", ClientID: "client-b", TermsName: "default"}
	if err := s.SaveTerms(other); err != nil {
		t.Fatalf("same name for another client: %v", err)
	}
	if err := s.SaveTerms(first); err != nil {
		t.Fatalf("re-save of the existing record: %v", err)
	}
}
