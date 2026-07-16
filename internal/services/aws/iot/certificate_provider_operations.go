package iot

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---------------------------------------------------------------------------
// CertificateProvider operations. AWS allows at most one certificate
// provider per account. When a provider is registered for the
// CreateCertificateFromCsr operation, all subsequent calls to that API
// invoke the provider's Lambda function instead of the platform's internal
// CA. The Lambda receives the CSR and returns a signed certificate PEM.
// ---------------------------------------------------------------------------

func (s *IoTService) CreateCertificateProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "certificateProviderName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// Enforce 1-per-account limit per AWS spec.
	activeRec := map[string]interface{}{}
	activeExists, _ := store.GetGenericExists("config/activeCertProvider", &activeRec)
	if activeExists {
		return nil, iotstore.ErrCertificateProviderAlreadyExists
	}
	lambdaFn := request.GetParamCaseInsensitive(req.Parameters, "lambdaFunctionArn")
	if lambdaFn == "" {
		return nil, iotstore.ErrMissingParam
	}
	now := time.Now().UTC().Unix()
	rec := map[string]interface{}{
		"certificateProviderName":     name,
		"certificateProviderArn":      iotstore.BuildCertificateProviderARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), name),
		"lambdaFunctionArn":           lambdaFn,
		"accountDefaultForOperations": req.Parameters["accountDefaultForOperations"],
		"creationDate":                now,
		"lastModifiedDate":            now,
	}
	if err := store.PutGeneric("iot-cert-provider/"+name, rec); err != nil {
		return nil, err
	}
	// Mark as the active provider if CreateCertificateFromCsr is listed.
	if hasOperation(req.Parameters["accountDefaultForOperations"], "CreateCertificateFromCsr") {
		if err := store.PutGeneric("config/activeCertProvider", map[string]interface{}{
			"providerName":      name,
			"lambdaFunctionArn": lambdaFn,
		}); err != nil {
			return nil, err
		}
	}
	return map[string]interface{}{
		"certificateProviderName": name,
		"certificateProviderArn":  rec["certificateProviderArn"],
	}, nil
}

func (s *IoTService) DescribeCertificateProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "certificateProviderName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("iot-cert-provider/"+name, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrCertificateProviderNotFound
	}
	return rec, nil
}

func (s *IoTService) UpdateCertificateProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "certificateProviderName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key := "iot-cert-provider/" + name
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrCertificateProviderNotFound
	}
	if lambdaFn := request.GetParamCaseInsensitive(req.Parameters, "lambdaFunctionArn"); lambdaFn != "" {
		rec["lambdaFunctionArn"] = lambdaFn
	}
	if ops, ok := req.Parameters["accountDefaultForOperations"]; ok {
		rec["accountDefaultForOperations"] = ops
	}
	rec["lastModifiedDate"] = time.Now().UTC().Unix()
	if err := store.PutGeneric(key, rec); err != nil {
		return nil, err
	}
	// Update active provider pointer if this provider handles CreateCertificateFromCsr.
	if hasOperation(req.Parameters["accountDefaultForOperations"], "CreateCertificateFromCsr") {
		if lambdaFn, ok := rec["lambdaFunctionArn"].(string); ok {
			if err := store.PutGeneric("config/activeCertProvider", map[string]interface{}{
				"providerName":      name,
				"lambdaFunctionArn": lambdaFn,
			}); err != nil {
				return nil, err
			}
		}
	}
	return map[string]interface{}{
		"certificateProviderName": name,
		"certificateProviderArn":  rec["certificateProviderArn"],
	}, nil
}

func (s *IoTService) DeleteCertificateProvider(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "certificateProviderName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	exists, err := store.GetGenericExists("iot-cert-provider/"+name, &map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrCertificateProviderNotFound
	}
	arn := iotstore.BuildCertificateProviderARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), name)
	_ = store.DeleteAllTags(arn)

	if err := store.DeleteGeneric("iot-cert-provider/" + name); err != nil {
		return nil, err
	}
	// Clear active provider pointer.
	if err := store.DeleteGeneric("config/activeCertProvider"); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) ListCertificateProviders(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := store.ListGeneric("iot-cert-provider/")
	if err != nil {
		return nil, err
	}
	summaries := make([]map[string]interface{}, 0, len(items))
	for _, rec := range items {
		summaries = append(summaries, map[string]interface{}{
			"certificateProviderName": rec["certificateProviderName"],
			"certificateProviderArn":  rec["certificateProviderArn"],
		})
	}
	return paginatedMaps("certificateProviders", summaries, req.Parameters), nil
}

// tryInvokeCertProvider checks whether a certificate provider is registered
// for the CreateCertificateFromCsr operation. If so, it invokes the Lambda
// and returns the certificate PEM from the response. Returns (certPEM, true)
// when a provider was invoked, or ("", false) when no provider is active.
func (s *IoTService) tryInvokeCertProvider(ctx context.Context, reqCtx *request.RequestContext, store interface {
	GetGenericExists(key string, dest interface{}) (bool, error)
}, csrPEM string) (string, bool, error) {
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("config/activeCertProvider", &rec)
	if err != nil || !exists {
		return "", false, nil
	}
	lambdaFn, _ := rec["lambdaFunctionArn"].(string)
	if lambdaFn == "" {
		return "", false, nil
	}
	invoker := s.deps.EventBus.LambdaInvoker()
	if invoker == nil {
		return "", false, nil
	}
	payload, _ := json.Marshal(map[string]string{
		"certificateSigningRequest": csrPEM,
	})
	_, resp, err := invoker.InvokeForGateway(ctx, lambdaFn, payload)
	if err != nil {
		return "", true, fmt.Errorf("certificate provider Lambda invocation failed: %w", err)
	}
	var result struct {
		CertificatePem string `json:"certificatePem"`
	}
	if err := json.Unmarshal(resp, &result); err != nil || result.CertificatePem == "" {
		return "", true, fmt.Errorf("certificate provider Lambda returned invalid response")
	}
	return result.CertificatePem, true, nil
}

// hasOperation checks whether the given operations list contains the
// specified operation name. Accepts both []interface{} and []string.
func hasOperation(ops interface{}, target string) bool {
	switch v := ops.(type) {
	case []interface{}:
		for _, op := range v {
			if s, ok := op.(string); ok && s == target {
				return true
			}
		}
	case []string:
		for _, op := range v {
			if op == target {
				return true
			}
		}
	}
	return false
}
