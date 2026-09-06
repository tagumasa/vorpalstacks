package cloudfront

import (
	awserrors "vorpalstacks/internal/common/errors"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
)

// CreatePublicKeyInput carries the parsed public key configuration.
type CreatePublicKeyInput struct {
	Config *cloudfrontstore.PublicKeyConfig
}

// UpdatePublicKeyInput carries the parameters for updating a public key.
type UpdatePublicKeyInput struct {
	Id      string
	IfMatch string
	Config  *cloudfrontstore.PublicKeyConfig
}

func (s *CloudFrontService) createPublicKeyCore(stores *cloudfrontStores, in CreatePublicKeyInput) (*cloudfrontstore.PublicKey, error) {
	if in.Config == nil || in.Config.Name == "" {
		return nil, invalidArgument("Name is required")
	}
	if in.Config.EncodedKey == "" {
		return nil, invalidArgument("EncodedKey is required")
	}
	// PublicKeyConfig models CallerReference as a required member; a
	// request without it must be rejected rather than silently minting
	// an idempotency token on the caller's behalf.
	if in.Config.CallerReference == "" {
		return nil, invalidArgument("CallerReference is required")
	}
	existing, _ := stores.publicKeys.GetByName(in.Config.Name)
	if existing != nil {
		return nil, awserrors.NewAWSError("PublicKeyAlreadyExists", "Public key with this name already exists", 409)
	}
	return stores.publicKeys.Create(in.Config)
}

func (s *CloudFrontService) getPublicKeyCore(stores *cloudfrontStores, id string) (*cloudfrontstore.PublicKey, error) {
	if err := requireID(id); err != nil {
		return nil, err
	}
	pk, err := stores.publicKeys.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchPublicKey", "Public key not found: "+id, 404)
		}
		return nil, err
	}
	return pk, nil
}

func (s *CloudFrontService) updatePublicKeyCore(stores *cloudfrontStores, in UpdatePublicKeyInput) (*cloudfrontstore.PublicKey, error) {
	if err := requireID(in.Id); err != nil {
		return nil, err
	}
	if in.Config == nil || in.Config.Name == "" {
		return nil, invalidArgument("Name is required")
	}
	if in.Config.EncodedKey == "" {
		return nil, invalidArgument("EncodedKey is required")
	}
	existing, err := stores.publicKeys.Get(in.Id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchPublicKey", "Public key not found: "+in.Id, 404)
		}
		return nil, err
	}
	if err := verifyIfMatch(in.IfMatch, existing.ETag); err != nil {
		return nil, err
	}
	if in.Config.EncodedKey != existing.PublicKeyConfig.EncodedKey {
		return nil, awserrors.NewAWSError("CannotChangeImmutablePublicKeyFields",
			"You can't change the value of a public key.", 400)
	}
	if in.Config.CallerReference == "" {
		in.Config.CallerReference = existing.PublicKeyConfig.CallerReference
	}
	return stores.publicKeys.Update(in.Id, in.Config)
}

func (s *CloudFrontService) deletePublicKeyCore(stores *cloudfrontStores, id, ifMatch string) error {
	if err := requireID(id); err != nil {
		return err
	}
	existing, err := stores.publicKeys.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return awserrors.NewAWSError("NoSuchPublicKey", "Public key not found: "+id, 404)
		}
		return err
	}
	if err := verifyIfMatch(ifMatch, existing.ETag); err != nil {
		return err
	}
	referenced, err := isPublicKeyReferenced(stores, id)
	if err != nil {
		return awserrors.NewInternalFailureException("Failed to check public key usage: " + err.Error())
	}
	if referenced {
		return awserrors.NewAWSError("PublicKeyInUse",
			"The specified public key is in use by one or more key groups", 409)
	}
	if err := stores.publicKeys.Delete(id); err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return awserrors.NewAWSError("NoSuchPublicKey", "Public key not found: "+id, 404)
		}
		return err
	}
	return nil
}

// PublicKeysResult is the transport-agnostic list result for public keys.
type PublicKeysResult struct {
	Keys              []*cloudfrontstore.PublicKey
	EffectiveMaxItems int
	NextMarker        string
}

func (s *CloudFrontService) listPublicKeysCore(stores *cloudfrontStores, marker string, maxItems int) (*PublicKeysResult, error) {
	effective := resolveListMaxItems(maxItems)
	result, err := stores.publicKeys.List(marker, effective)
	if err != nil {
		return nil, err
	}
	return &PublicKeysResult{
		Keys:              result.PublicKeys,
		EffectiveMaxItems: effective,
		NextMarker:        result.NextMarker,
	}, nil
}
