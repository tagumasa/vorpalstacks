package wafv2

import (
	types "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

// parseAddressList parses the raw Addresses member and validates each
// entry as a well-formed CIDR notation matching the specified
// IPAddressVersion. Addresses is a required member of CreateIPSet and
// UpdateIPSet: an omitted member is rejected, while the documented empty
// array ("Addresses": []) remains valid.
func parseAddressList(raw interface{}, ipAddressVersion string) ([]string, error) {
	if raw == nil {
		return nil, invalidParamError("Addresses is a required member")
	}
	arr, ok := raw.([]interface{})
	if !ok {
		return nil, invalidParamError("Addresses must be a list of CIDR strings")
	}
	var addresses []string
	for _, a := range arr {
		s, ok := a.(string)
		if !ok {
			return nil, invalidParamError("Addresses entries must be strings")
		}
		if err := validateIPAddress(s, ipAddressVersion); err != nil {
			return nil, err
		}
		addresses = append(addresses, s)
	}
	return addresses, nil
}

// IPSetCreateInput is the transport-agnostic input for creating an IP set.
type IPSetCreateInput struct {
	Name             string
	Scope            string
	IPAddressVersion string
	AddressesRaw     interface{}
	Description      string
	Tags             []types.Tag
}

// createIPSetCore is the single entry point for creating an IP set. The
// validation ladder order matches the original handler: name, scope,
// address version, addresses, description.
func (s *WAFv2Service) createIPSetCore(stores *wafv2Stores, in IPSetCreateInput) (*wafstore.IPSet, error) {
	if err := validateEntityName(in.Name); err != nil {
		return nil, err
	}
	if err := validateScope(in.Scope); err != nil {
		return nil, err
	}
	if err := validateIPAddressVersion(in.IPAddressVersion); err != nil {
		return nil, err
	}
	addresses, err := parseAddressList(in.AddressesRaw, in.IPAddressVersion)
	if err != nil {
		return nil, err
	}
	if err := validateEntityDescription(in.Description); err != nil {
		return nil, err
	}

	id, err := generateID()
	if err != nil {
		return nil, err
	}

	ipSet, err := stores.ipSets.Create(id, in.Name, in.Description, in.IPAddressVersion, addresses, in.Scope)
	if err != nil {
		if wafstore.IsAlreadyExists(err) {
			return nil, newAPIError("WAFDuplicateItemException", "AWS WAF couldn't perform the operation because some resource in your request is a duplicate of an existing one", 400)
		}
		return nil, err
	}

	if len(in.Tags) > 0 {
		if err := stores.tags.TagFromSlice(ipSet.ARN, in.Tags); err != nil {
			logs.Warn("failed to persist tags for IPSet", logs.String("id", ipSet.ID), logs.Err(err))
		}
	}

	return ipSet, nil
}

// getIPSetCore is the single entry point for retrieving an IP set.
func (s *WAFv2Service) getIPSetCore(stores *wafv2Stores, id string) (*wafstore.IPSet, error) {
	if id == "" {
		return nil, invalidParamError("Id is required")
	}

	ipSet, err := stores.ipSets.Get(id)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("IPSet")
		}
		return nil, err
	}

	return ipSet, nil
}

// IPSetListInput is the transport-agnostic input for listing IP sets.
type IPSetListInput struct {
	Scope      string
	Limit      int
	NextMarker string
}

// listIPSetsCore is the single entry point for listing IP sets.
func (s *WAFv2Service) listIPSetsCore(stores *wafv2Stores, in IPSetListInput) (*wafstore.IPSetListResult, error) {
	if err := validateScope(in.Scope); err != nil {
		return nil, err
	}

	return stores.ipSets.List(in.NextMarker, in.Limit, in.Scope)
}

// IPSetUpdateInput is the transport-agnostic input for updating an IP set.
// The raw addresses travel untyped because their CIDR validation must run
// against the stored IPAddressVersion, which is only known after the
// pre-update fetch.
type IPSetUpdateInput struct {
	Id           string
	LockToken    string
	AddressesRaw interface{}
	Description  string
}

// updateIPSetCore is the single entry point for updating an IP set.
func (s *WAFv2Service) updateIPSetCore(stores *wafv2Stores, in IPSetUpdateInput) (*wafstore.IPSet, error) {
	if in.Id == "" {
		return nil, invalidParamError("Id is required")
	}

	if in.LockToken == "" {
		return nil, invalidParamError("LockToken is required")
	}

	ipSet, err := stores.ipSets.Get(in.Id)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("IPSet")
		}
		return nil, err
	}

	addresses, err := parseAddressList(in.AddressesRaw, ipSet.IPAddressVersion)
	if err != nil {
		return nil, err
	}

	if err := validateEntityDescription(in.Description); err != nil {
		return nil, err
	}

	ipSet, err = stores.ipSets.Update(in.Id, in.LockToken, addresses, in.Description)
	if err != nil {
		if wafstore.IsLockTokenMismatch(err) {
			return nil, lockTokenError()
		}
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("IPSet")
		}
		return nil, err
	}

	return ipSet, nil
}

// deleteIPSetCore is the single entry point for deleting an IP set,
// including the tag cleanup on the deleted ARN.
func (s *WAFv2Service) deleteIPSetCore(stores *wafv2Stores, id, lockToken string) error {
	if id == "" {
		return invalidParamError("Id is required")
	}

	if lockToken == "" {
		return invalidParamError("LockToken is required")
	}

	deleted, err := stores.ipSets.Delete(id, lockToken)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return notFoundError("IPSet")
		}
		if wafstore.IsLockTokenMismatch(err) {
			return lockTokenError()
		}
		return err
	}

	if deleted.ARN != "" {
		if err := stores.tags.Delete(deleted.ARN); err != nil {
			logs.Warn("failed to clean up tags for deleted IPSet", logs.String("id", id), logs.Err(err))
		}
	}

	return nil
}
