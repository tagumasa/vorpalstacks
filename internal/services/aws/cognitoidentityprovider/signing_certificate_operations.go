package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// Handler for the pool signing-certificate family; the certificate
// derivation lives in signing_certificate_core.go.

// GetSigningCertificate returns the user pool's signing certificate.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_GetSigningCertificate.html
func (s *CognitoService) GetSigningCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	certificate, err := s.getSigningCertificateCore(reqCtx.GetRegion(), req.GetParam("UserPoolId"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Certificate": certificate,
	}, nil
}
