package iot

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---------------------------------------------------------------------------
// CertificateProvider Core. AWS allows at most one certificate provider per
// account. When a provider is registered for the CreateCertificateFromCsr
// operation, all subsequent calls to that API invoke the provider's Lambda
// function instead of the platform's internal CA. The Lambda receives the CSR
// and returns a signed certificate PEM.
// ---------------------------------------------------------------------------

// CreateCertificateProviderInput carries the fields for
// CreateCertificateProvider. AccountDefaultForOperations keeps the raw wire
// value (a list of operation names).
type CreateCertificateProviderInput struct {
	CertificateProviderName     string
	LambdaFunctionARN           string
	AccountDefaultForOperations interface{}
}

// CertificateProviderResult is the transport-agnostic result of the
// certificate-provider create and update operations.
type CertificateProviderResult struct {
	CertificateProviderName string
	CertificateProviderARN  string
}

// UpdateCertificateProviderInput carries the fields for
// UpdateCertificateProvider. The Lambda function is applied only when
// non-empty; the operations list is applied only when the caller supplied it.
type UpdateCertificateProviderInput struct {
	CertificateProviderName     string
	LambdaFunctionARN           string
	AccountDefaultForOperations interface{}
	OperationsProvided          bool
}

// ListCertificateProvidersResult is the transport-agnostic result of
// ListCertificateProviders.
type ListCertificateProvidersResult struct {
	CertificateProviders []map[string]interface{}
}

// validateCertProviderOperations enforces the accountDefaultForOperations
// constraint: exactly one entry whose value is CreateCertificateFromCsr.
func validateCertProviderOperations(ops interface{}) error {
	list, ok := ops.([]interface{})
	if !ok || len(list) != 1 {
		return iotstore.ErrInvalidRequest
	}
	if s, ok := list[0].(string); !ok || s != "CreateCertificateFromCsr" {
		return iotstore.ErrInvalidRequest
	}
	return nil
}

// createCertificateProviderCore registers the account's certificate provider.
func (s *IoTService) createCertificateProviderCore(store iotstore.IotStoreInterface, in CreateCertificateProviderInput) (*CertificateProviderResult, error) {
	if in.CertificateProviderName == "" {
		return nil, iotstore.ErrMissingParam
	}
	if in.LambdaFunctionARN == "" {
		return nil, iotstore.ErrMissingParam
	}
	if err := validateCertProviderOperations(in.AccountDefaultForOperations); err != nil {
		return nil, err
	}
	// Enforce 1-per-account limit per AWS spec.
	activeRec := map[string]interface{}{}
	activeExists, _ := store.GetGenericExists("config/activeCertProvider", &activeRec)
	if activeExists {
		return nil, iotstore.ErrCertificateProviderAlreadyExists
	}
	now := time.Now().UTC().Unix()
	rec := map[string]interface{}{
		"certificateProviderName":     in.CertificateProviderName,
		"certificateProviderArn":      iotstore.BuildCertificateProviderARN(store.GetAccountID(), store.GetRegion(), in.CertificateProviderName),
		"lambdaFunctionArn":           in.LambdaFunctionARN,
		"accountDefaultForOperations": in.AccountDefaultForOperations,
		"creationDate":                now,
		"lastModifiedDate":            now,
	}
	if err := store.PutGeneric("iot-cert-provider/"+in.CertificateProviderName, rec); err != nil {
		return nil, err
	}
	// Mark as the active provider if CreateCertificateFromCsr is listed.
	if hasOperation(in.AccountDefaultForOperations, "CreateCertificateFromCsr") {
		if err := store.PutGeneric("config/activeCertProvider", map[string]interface{}{
			"providerName":      in.CertificateProviderName,
			"lambdaFunctionArn": in.LambdaFunctionARN,
		}); err != nil {
			return nil, err
		}
	}
	return &CertificateProviderResult{
		CertificateProviderName: in.CertificateProviderName,
		CertificateProviderARN:  rec["certificateProviderArn"].(string),
	}, nil
}

// describeCertificateProviderCore retrieves a certificate provider record.
func (s *IoTService) describeCertificateProviderCore(store iotstore.IotStoreInterface, name string) (map[string]interface{}, error) {
	if name == "" {
		return nil, iotstore.ErrMissingParam
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

// updateCertificateProviderCore applies the supplied fields to an existing
// certificate provider and refreshes the active-provider pointer when the
// operations list includes CreateCertificateFromCsr.
func (s *IoTService) updateCertificateProviderCore(store iotstore.IotStoreInterface, in UpdateCertificateProviderInput) (*CertificateProviderResult, error) {
	if in.CertificateProviderName == "" {
		return nil, iotstore.ErrMissingParam
	}
	key := "iot-cert-provider/" + in.CertificateProviderName
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrCertificateProviderNotFound
	}
	if in.OperationsProvided {
		if err := validateCertProviderOperations(in.AccountDefaultForOperations); err != nil {
			return nil, err
		}
	}
	if in.LambdaFunctionARN != "" {
		rec["lambdaFunctionArn"] = in.LambdaFunctionARN
	}
	if in.OperationsProvided {
		rec["accountDefaultForOperations"] = in.AccountDefaultForOperations
	}
	rec["lastModifiedDate"] = time.Now().UTC().Unix()
	if err := store.PutGeneric(key, rec); err != nil {
		return nil, err
	}
	// Update active provider pointer if this provider handles CreateCertificateFromCsr.
	if hasOperation(in.AccountDefaultForOperations, "CreateCertificateFromCsr") {
		if lambdaFn, ok := rec["lambdaFunctionArn"].(string); ok {
			if err := store.PutGeneric("config/activeCertProvider", map[string]interface{}{
				"providerName":      in.CertificateProviderName,
				"lambdaFunctionArn": lambdaFn,
			}); err != nil {
				return nil, err
			}
		}
	}
	arn, _ := rec["certificateProviderArn"].(string)
	return &CertificateProviderResult{
		CertificateProviderName: in.CertificateProviderName,
		CertificateProviderARN:  arn,
	}, nil
}

// deleteCertificateProviderCore removes the certificate provider, its tags,
// and the active-provider pointer.
func (s *IoTService) deleteCertificateProviderCore(store iotstore.IotStoreInterface, name string) error {
	if name == "" {
		return iotstore.ErrMissingParam
	}
	exists, err := store.GetGenericExists("iot-cert-provider/"+name, &map[string]interface{}{})
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrCertificateProviderNotFound
	}
	arn := iotstore.BuildCertificateProviderARN(store.GetAccountID(), store.GetRegion(), name)
	_ = store.DeleteAllTags(arn)

	if err := store.DeleteGeneric("iot-cert-provider/" + name); err != nil {
		return err
	}
	// Clear active provider pointer.
	return store.DeleteGeneric("config/activeCertProvider")
}

// listCertificateProvidersCore lists provider name/ARN summaries.
func (s *IoTService) listCertificateProvidersCore(store iotstore.IotStoreInterface) (*ListCertificateProvidersResult, error) {
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
	return &ListCertificateProvidersResult{CertificateProviders: summaries}, nil
}

// invokeCertProviderCore checks whether a certificate provider is registered
// for the CreateCertificateFromCsr operation. If so, it invokes the Lambda
// and returns the certificate PEM from the response. Returns (certPEM, true)
// when a provider was invoked, or ("", false) when no provider is active.
func (s *IoTService) invokeCertProviderCore(ctx context.Context, store iotstore.IotStoreInterface, csrPEM string) (string, bool, error) {
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
