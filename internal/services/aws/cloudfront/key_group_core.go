package cloudfront

import (
	awserrors "vorpalstacks/internal/common/errors"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
)

// validateKeyGroupConfig enforces the required fields and the public key
// count limit of a key group configuration, shared by the create and
// update cores.
func validateKeyGroupConfig(config *cloudfrontstore.KeyGroupConfig) error {
	if config == nil || config.Name == "" {
		return invalidArgument("Name is required")
	}
	if len(config.Items) == 0 {
		return invalidArgument("At least one public key ID is required in Items")
	}
	if len(config.Items) > cloudfrontstore.MaxPublicKeysPerKeyGroup {
		return awserrors.NewAWSError("TooManyPublicKeysInKeyGroup", tooManyPublicKeysMsg(), 400)
	}
	return nil
}

// CreateKeyGroupInput carries the parsed key group configuration.
type CreateKeyGroupInput struct {
	Config *cloudfrontstore.KeyGroupConfig
}

// UpdateKeyGroupInput carries the parameters for updating a key group.
type UpdateKeyGroupInput struct {
	Id      string
	IfMatch string
	Config  *cloudfrontstore.KeyGroupConfig
}

// validateKeyGroupItems verifies that every public key identifier in a key
// group refers to an existing public key. The Developer Guide's signing
// workflow requires public keys to be uploaded before they are added to a
// key group, and unknown identifiers are rejected with InvalidArgument.
func (s *CloudFrontService) validateKeyGroupItems(stores *cloudfrontStores, items []string) error {
	for _, pkID := range items {
		if pkID == "" {
			continue
		}
		if _, err := stores.publicKeys.Get(pkID); err != nil {
			if cloudfrontstore.IsNotFound(err) {
				return awserrors.NewAWSError("InvalidArgument", "The specified public key does not exist: "+pkID, 400)
			}
			return err
		}
	}
	return nil
}

func (s *CloudFrontService) createKeyGroupCore(stores *cloudfrontStores, in CreateKeyGroupInput) (*cloudfrontstore.KeyGroup, error) {
	if err := validateKeyGroupConfig(in.Config); err != nil {
		return nil, err
	}
	if err := s.validateKeyGroupItems(stores, in.Config.Items); err != nil {
		return nil, err
	}
	existing, _ := stores.keyGroups.GetByName(in.Config.Name)
	if existing != nil {
		return nil, awserrors.NewAWSError("KeyGroupAlreadyExists", "Key group with this name already exists", 409)
	}
	return stores.keyGroups.Create(in.Config)
}

func (s *CloudFrontService) getKeyGroupCore(stores *cloudfrontStores, id string) (*cloudfrontstore.KeyGroup, error) {
	if err := requireID(id); err != nil {
		return nil, err
	}
	kg, err := stores.keyGroups.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchResource", "Key group not found: "+id, 404)
		}
		return nil, err
	}
	return kg, nil
}

func (s *CloudFrontService) updateKeyGroupCore(stores *cloudfrontStores, in UpdateKeyGroupInput) (*cloudfrontstore.KeyGroup, error) {
	if err := requireID(in.Id); err != nil {
		return nil, err
	}
	if err := validateKeyGroupConfig(in.Config); err != nil {
		return nil, err
	}
	if err := s.validateKeyGroupItems(stores, in.Config.Items); err != nil {
		return nil, err
	}
	existing, err := stores.keyGroups.Get(in.Id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchResource", "Key group not found: "+in.Id, 404)
		}
		return nil, err
	}
	if err := verifyIfMatch(in.IfMatch, existing.ETag); err != nil {
		return nil, err
	}
	if err := ensureNameAvailable(existing.KeyGroupConfig.Name, in.Config.Name,
		func(name string) bool { dup, _ := stores.keyGroups.GetByName(name); return dup != nil },
		awserrors.NewAWSError("KeyGroupAlreadyExists", "Key group with this name already exists", 409)); err != nil {
		return nil, err
	}
	return stores.keyGroups.Update(in.Id, in.Config)
}

func (s *CloudFrontService) deleteKeyGroupCore(stores *cloudfrontStores, id, ifMatch string) error {
	if err := requireID(id); err != nil {
		return err
	}
	existing, err := stores.keyGroups.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return awserrors.NewAWSError("NoSuchResource", "Key group not found: "+id, 404)
		}
		return err
	}
	if err := verifyIfMatch(ifMatch, existing.ETag); err != nil {
		return err
	}
	referenced, err := isKeyGroupReferenced(stores, id)
	if err != nil {
		return awserrors.NewAWSError("InternalError", "Failed to check key group usage: "+err.Error(), 500)
	}
	if referenced {
		return awserrors.NewAWSError("ResourceInUse",
			"Cannot delete this key group because it is referenced by one or more distributions", 409)
	}
	if err := stores.keyGroups.Delete(id); err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return awserrors.NewAWSError("NoSuchResource", "Key group not found: "+id, 404)
		}
		return err
	}
	return nil
}

// KeyGroupsResult is the transport-agnostic list result for key groups.
type KeyGroupsResult struct {
	Groups            []*cloudfrontstore.KeyGroup
	EffectiveMaxItems int
	NextMarker        string
}

func (s *CloudFrontService) listKeyGroupsCore(stores *cloudfrontStores, marker string, maxItems int) (*KeyGroupsResult, error) {
	effective := resolveListMaxItems(maxItems)
	result, err := stores.keyGroups.List(marker, effective)
	if err != nil {
		return nil, err
	}
	return &KeyGroupsResult{
		Groups:            result.KeyGroups,
		EffectiveMaxItems: effective,
		NextMarker:        result.NextMarker,
	}, nil
}
