package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// CreateProvisioningClaim issues a temporary provisioning claim bound to an
// existing provisioning template.
func (s *IoTService) CreateProvisioningClaim(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.createProvisioningClaimCore(store, request.GetParamCaseInsensitive(req.Parameters, "templateName"))
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"certificateId":  result.CertificateID,
		"certificatePem": result.CertificatePem,
		"keyPair": map[string]interface{}{
			"PublicKey":  result.PublicKeyPEM,
			"PrivateKey": result.PrivateKeyPEM,
		},
		"expiration": result.Expiration,
	}, nil
}

func (s *IoTService) GetRegistrationCode(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.getRegistrationCodeCore(store)
}

func (s *IoTService) DeleteRegistrationCode(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteRegistrationCodeCore(store); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
