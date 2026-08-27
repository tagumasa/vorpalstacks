package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// ---------------------------------------------------------------------------
// CertificateProvider operations. AWS allows at most one certificate
// provider per account. When a provider is registered for the
// CreateCertificateFromCsr operation, all subsequent calls to that API
// invoke the provider's Lambda function instead of the platform's internal
// CA. The Lambda receives the CSR and returns a signed certificate PEM.
// ---------------------------------------------------------------------------

func (s *IoTService) CreateCertificateProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.createCertificateProviderCore(store, CreateCertificateProviderInput{
		CertificateProviderName:     request.GetParamCaseInsensitive(req.Parameters, "certificateProviderName"),
		LambdaFunctionARN:           request.GetParamCaseInsensitive(req.Parameters, "lambdaFunctionArn"),
		AccountDefaultForOperations: req.Parameters["accountDefaultForOperations"],
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"certificateProviderName": result.CertificateProviderName,
		"certificateProviderArn":  result.CertificateProviderARN,
	}, nil
}

func (s *IoTService) DescribeCertificateProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	rec, err := s.describeCertificateProviderCore(store, request.GetParamCaseInsensitive(req.Parameters, "certificateProviderName"))
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func (s *IoTService) UpdateCertificateProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Presence is the raw member key (case-sensitive), matching the wire
	// member name exactly.
	_, operationsProvided := req.Parameters["accountDefaultForOperations"]

	result, err := s.updateCertificateProviderCore(store, UpdateCertificateProviderInput{
		CertificateProviderName:     request.GetParamCaseInsensitive(req.Parameters, "certificateProviderName"),
		LambdaFunctionARN:           request.GetParamCaseInsensitive(req.Parameters, "lambdaFunctionArn"),
		AccountDefaultForOperations: req.Parameters["accountDefaultForOperations"],
		OperationsProvided:          operationsProvided,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"certificateProviderName": result.CertificateProviderName,
		"certificateProviderArn":  result.CertificateProviderARN,
	}, nil
}

func (s *IoTService) DeleteCertificateProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteCertificateProviderCore(store, request.GetParamCaseInsensitive(req.Parameters, "certificateProviderName")); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) ListCertificateProviders(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.listCertificateProvidersCore(store)
	if err != nil {
		return nil, err
	}

	return paginatedMaps("certificateProviders", result.CertificateProviders, req.Parameters)
}
