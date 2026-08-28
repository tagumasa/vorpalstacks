package iam

import (
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// createAccountAliasCore validates input and sets the account alias.
func (s *IAMService) createAccountAliasCore(store *iamstore.IAMStore, alias string) error {
	if alias == "" {
		return NewValidationError("AccountAlias")
	}
	if err := validateAccountAlias(alias); err != nil {
		return err
	}
	return store.AccountAlias().Put(alias)
}

// deleteAccountAliasCore validates input and removes the account alias.
// The alias in the request must match the stored alias.
func (s *IAMService) deleteAccountAliasCore(store *iamstore.IAMStore, accountAlias string) error {
	if accountAlias == "" {
		return NewValidationError("AccountAlias")
	}
	existing, err := store.AccountAlias().Get()
	if err != nil {
		return err
	}
	if existing == nil || existing.AccountAlias != accountAlias {
		return NewNoSuchEntityError("Account Alias", accountAlias)
	}
	return store.AccountAlias().Delete()
}

// listAccountAliasesCore retrieves the account alias, or nil when none is
// set.
func (s *IAMService) listAccountAliasesCore(store *iamstore.IAMStore) (*iamstore.AccountAlias, error) {
	return store.AccountAlias().Get()
}
