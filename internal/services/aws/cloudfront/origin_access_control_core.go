package cloudfront

import (
	awserrors "vorpalstacks/internal/common/errors"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
)

// validateOriginAccessControlConfig enforces the required fields and enum
// values of an origin access control configuration, shared by the create
// and update cores.
func validateOriginAccessControlConfig(config *cloudfrontstore.OriginAccessControlConfig) error {
	if config == nil || config.Name == "" {
		return invalidArgument("Name is required")
	}
	if config.OriginAccessControlOriginType == "" {
		return invalidArgument("OriginAccessControlOriginType is required")
	}
	if config.SigningBehavior == "" {
		return invalidArgument("SigningBehavior is required")
	}
	if config.SigningProtocol == "" {
		return invalidArgument("SigningProtocol is required")
	}
	if !isValidOriginAccessControlOriginType(config.OriginAccessControlOriginType) {
		return invalidArgument("Invalid OriginAccessControlOriginType. Must be one of: " + originAccessControlOriginTypeValues())
	}
	if !isValidSigningBehavior(config.SigningBehavior) {
		return invalidArgument("Invalid SigningBehavior. Must be one of: " + signingBehaviorValues())
	}
	if !isValidSigningProtocol(config.SigningProtocol) {
		return invalidArgument("Invalid SigningProtocol. Must be: " + signingProtocolValues())
	}
	return nil
}

// CreateOriginAccessControlInput carries the parsed origin access control
// configuration.
type CreateOriginAccessControlInput struct {
	Config *cloudfrontstore.OriginAccessControlConfig
}

// UpdateOriginAccessControlInput carries the parameters for updating an
// origin access control.
type UpdateOriginAccessControlInput struct {
	Id      string
	IfMatch string
	Config  *cloudfrontstore.OriginAccessControlConfig
}

func (s *CloudFrontService) createOriginAccessControlCore(stores *cloudfrontStores, in CreateOriginAccessControlInput) (*cloudfrontstore.OriginAccessControl, error) {
	if err := validateOriginAccessControlConfig(in.Config); err != nil {
		return nil, err
	}
	existing, _ := stores.originAccessControls.GetByName(in.Config.Name)
	if existing != nil {
		return nil, awserrors.NewAWSError("OriginAccessControlAlreadyExists", "Origin access control with this name already exists", 409)
	}
	return stores.originAccessControls.Create(in.Config)
}

func (s *CloudFrontService) getOriginAccessControlCore(stores *cloudfrontStores, id string) (*cloudfrontstore.OriginAccessControl, error) {
	if err := requireID(id); err != nil {
		return nil, err
	}
	oac, err := stores.originAccessControls.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchOriginAccessControl", "Origin access control not found", 404)
		}
		return nil, err
	}
	return oac, nil
}

func (s *CloudFrontService) updateOriginAccessControlCore(stores *cloudfrontStores, in UpdateOriginAccessControlInput) (*cloudfrontstore.OriginAccessControl, error) {
	if err := requireID(in.Id); err != nil {
		return nil, err
	}
	if err := validateOriginAccessControlConfig(in.Config); err != nil {
		return nil, err
	}
	existing, err := stores.originAccessControls.Get(in.Id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchOriginAccessControl", "Origin access control not found", 404)
		}
		return nil, err
	}
	if err := verifyIfMatch(in.IfMatch, existing.ETag); err != nil {
		return nil, err
	}
	if err := ensureNameAvailable(existing.Name, in.Config.Name,
		func(name string) bool { dup, _ := stores.originAccessControls.GetByName(name); return dup != nil },
		awserrors.NewAWSError("OriginAccessControlAlreadyExists", "Origin access control with this name already exists", 409)); err != nil {
		return nil, err
	}
	oac, err := stores.originAccessControls.Update(in.Id, in.Config)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchOriginAccessControl", "Origin access control not found", 404)
		}
		return nil, err
	}
	return oac, nil
}

func (s *CloudFrontService) deleteOriginAccessControlCore(stores *cloudfrontStores, id, ifMatch string) error {
	if err := requireID(id); err != nil {
		return err
	}
	existing, err := stores.originAccessControls.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return awserrors.NewAWSError("NoSuchOriginAccessControl", "Origin access control not found", 404)
		}
		return err
	}
	if err := verifyIfMatch(ifMatch, existing.ETag); err != nil {
		return err
	}
	attached, err := isOriginAccessControlAttached(stores, id)
	if err != nil {
		return awserrors.NewAWSError("InternalError", "Failed to check origin access control usage: "+err.Error(), 500)
	}
	if attached {
		return awserrors.NewAWSError("OriginAccessControlInUse",
			"Cannot delete this origin access control because it is attached to one or more distributions", 409)
	}
	if err := stores.originAccessControls.Delete(id); err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return awserrors.NewAWSError("NoSuchOriginAccessControl", "Origin access control not found", 404)
		}
		return err
	}
	return nil
}

// OriginAccessControlsResult is the transport-agnostic list result for
// origin access controls.
type OriginAccessControlsResult struct {
	Controls          []*cloudfrontstore.OriginAccessControl
	IsTruncated       bool
	EffectiveMaxItems int
	NextMarker        string
}

func (s *CloudFrontService) listOriginAccessControlsCore(stores *cloudfrontStores, marker string, maxItems int) (*OriginAccessControlsResult, error) {
	effective := resolveListMaxItems(maxItems)
	result, err := stores.originAccessControls.List(marker, effective)
	if err != nil {
		return nil, err
	}
	return &OriginAccessControlsResult{
		Controls:          result.OriginAccessControls,
		IsTruncated:       result.IsTruncated,
		EffectiveMaxItems: effective,
		NextMarker:        result.NextMarker,
	}, nil
}
