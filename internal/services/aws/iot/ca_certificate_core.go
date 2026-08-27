package iot

import (
	"encoding/json"
	"time"

	iotstore "vorpalstacks/internal/store/aws/iot"
	vcrypto "vorpalstacks/internal/utils/crypto"
)

// ---------------------------------------------------------------------------
// CA certificate Core — registration, description, transfer bookkeeping
// ---------------------------------------------------------------------------

// RegisterCACertificateInput carries the fields for RegisterCACertificate.
type RegisterCACertificateInput struct {
	CACertificatePEM           string
	VerificationCertificatePEM string
	SetAsActive                bool
	AllowAutoRegistration      bool
	RegistrationConfig         string // serialised RegistrationConfig structure
	Tags                       map[string]string
	CertificateMode            string
}

// RegisterCACertificateResult is the transport-agnostic result of
// RegisterCACertificate.
type RegisterCACertificateResult struct {
	CertificateARN string
	CertificateID  string
}

// ListCACertificatesResult is the transport-agnostic result of
// ListCACertificates.
type ListCACertificatesResult struct {
	Certificates []map[string]interface{}
}

// UpdateCACertificateInput carries the fields for UpdateCACertificate. An
// empty NewStatus/NewAutoRegistrationStatus/RegistrationConfig leaves the
// stored value untouched; RemoveAutoRegistration clears the auto-registration
// setting (applied after NewAutoRegistrationStatus).
type UpdateCACertificateInput struct {
	CertificateID             string
	NewStatus                 string
	NewAutoRegistrationStatus string
	RegistrationConfig        string // serialised RegistrationConfig structure
	RemoveAutoRegistration    bool
}

// ListCertificatesByCACertificateInput carries the fields for
// ListCertificatesByCA.
type ListCertificatesByCACertificateInput struct {
	CACertificateID string
	Marker          string
	MaxItems        int
}

// ListCertificatesByCACertificateResult is the transport-agnostic result of
// ListCertificatesByCA.
type ListCertificatesByCACertificateResult struct {
	Certificates []*iotstore.Certificate
	NextMarker   string
}

// RegisterCertificateWithoutCAInput carries the fields for
// RegisterCertificateWithoutCA. An empty Status defaults to INACTIVE.
type RegisterCertificateWithoutCAInput struct {
	CertificatePEM string
	Status         string
}

// RegisterCertificateWithoutCAResult is the transport-agnostic result of
// RegisterCertificateWithoutCA.
type RegisterCertificateWithoutCAResult struct {
	Certificate *iotstore.Certificate
}

// TransferCertificateInput carries the fields for TransferCertificate.
type TransferCertificateInput struct {
	CertificateID    string
	TargetAWSAccount string
	TransferMessage  string
}

// TransferCertificateResult is the transport-agnostic result of
// TransferCertificate.
type TransferCertificateResult struct {
	TransferredCertificateARN string
}

// ListOutgoingCertificatesResult is the transport-agnostic result of
// ListOutgoingCertificates.
type ListOutgoingCertificatesResult struct {
	OutgoingCertificates []map[string]interface{}
}

// registerCACertificateCore registers a CA certificate PEM with a
// deterministic ID derived from the certificate content hash (prevents
// duplicate registration of the same CA certificate, matching AWS behaviour).
// The initial status follows setAsActive and the auto-registration status
// follows allowAutoRegistration; both default to the inactive/disabled side
// per the model's documented defaults.
func (s *IoTService) registerCACertificateCore(store iotstore.IotStoreInterface, in RegisterCACertificateInput) (*RegisterCACertificateResult, error) {
	if in.CACertificatePEM == "" {
		return nil, iotstore.ErrMissingParam
	}
	if _, ok := caCertificateModes[in.CertificateMode]; !ok && in.CertificateMode != "" {
		return nil, iotstore.ErrInvalidRequest
	}
	// The documented pairing rule: SNI_ONLY forbids a verification
	// certificate while DEFAULT (or an omitted mode) requires one.
	if in.CertificateMode == "SNI_ONLY" && in.VerificationCertificatePEM != "" {
		return nil, iotstore.ErrInvalidRequest
	}
	if in.CertificateMode != "SNI_ONLY" && in.VerificationCertificatePEM == "" {
		return nil, iotstore.ErrInvalidRequest
	}
	status := "INACTIVE"
	if in.SetAsActive {
		status = "ACTIVE"
	}
	autoRegistration := "DISABLE"
	if in.AllowAutoRegistration {
		autoRegistration = "ENABLE"
	}
	now := time.Now()
	caCertID := vcrypto.FingerprintPEM(in.CACertificatePEM)
	rec := map[string]interface{}{
		"certificateId":              caCertID,
		"status":                     status,
		"autoRegistrationStatus":     autoRegistration,
		"caCertificatePem":           in.CACertificatePEM,
		"verificationCertificatePem": in.VerificationCertificatePEM,
		"creationDate":               now.Unix(),
	}
	if in.CertificateMode != "" {
		rec["certificateMode"] = in.CertificateMode
	}
	if in.RegistrationConfig != "" {
		rec["registrationConfig"] = in.RegistrationConfig
	}
	if err := store.PutGeneric("caCert/"+caCertID, rec); err != nil {
		return nil, err
	}
	arn := iotstore.BuildCACertificateARN(store.GetAccountID(), store.GetRegion(), caCertID)
	if len(in.Tags) > 0 {
		if err := store.TagResource(arn, in.Tags); err != nil {
			return nil, err
		}
	}
	return &RegisterCACertificateResult{
		CertificateARN: arn,
		CertificateID:  caCertID,
	}, nil
}

// The CACertificateStatus enum values (the model's only two members).
var caCertificateStatuses = map[string]struct{}{
	"ACTIVE": {}, "INACTIVE": {},
}

// The AutoRegistrationStatus enum values.
var caAutoRegistrationStatuses = map[string]struct{}{
	"ENABLE": {}, "DISABLE": {},
}

// The CertificateMode enum values.
var caCertificateModes = map[string]struct{}{
	"DEFAULT": {}, "SNI_ONLY": {},
}

// DescribeCACertificateResult carries the two top-level response members:
// the CACertificateDescription projection and the stored registration
// configuration (nil when unset).
type DescribeCACertificateResult struct {
	CertificateDescription map[string]interface{}
	RegistrationConfig     map[string]interface{}
}

// describeCACertificateCore retrieves a registered CA certificate and
// projects it onto the response shape: the certificateDescription members
// plus the top-level registrationConfig when one is stored.
func (s *IoTService) describeCACertificateCore(store iotstore.IotStoreInterface, caCertID string) (*DescribeCACertificateResult, error) {
	if caCertID == "" {
		return nil, iotstore.ErrMissingParam
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("caCert/"+caCertID, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrCertificateNotFound
	}
	desc := map[string]interface{}{
		"certificateArn": iotstore.BuildCACertificateARN(store.GetAccountID(), store.GetRegion(), caCertID),
		"certificateId":  caCertID,
		"status":         rec["status"],
		"ownedBy":        store.GetAccountID(),
	}
	if pem, ok := rec["caCertificatePem"].(string); ok && pem != "" {
		desc["certificatePem"] = pem
	}
	for _, member := range []string{"autoRegistrationStatus", "certificateMode"} {
		if v, ok := rec[member]; ok {
			desc[member] = v
		}
	}
	for _, member := range []string{"creationDate", "lastModifiedDate"} {
		switch v := rec[member].(type) {
		case int64:
			desc[member] = v
		case float64:
			desc[member] = int64(v)
		}
	}
	var registrationConfig map[string]interface{}
	if stored, ok := rec["registrationConfig"].(string); ok && stored != "" {
		_ = json.Unmarshal([]byte(stored), &registrationConfig)
	}
	return &DescribeCACertificateResult{
		CertificateDescription: desc,
		RegistrationConfig:     registrationConfig,
	}, nil
}

// listCACertificatesCore lists every registered CA certificate with its ARN.
func (s *IoTService) listCACertificatesCore(store iotstore.IotStoreInterface) (*ListCACertificatesResult, error) {
	items, err := store.ListGeneric("caCert/")
	if err != nil {
		return nil, err
	}
	certs := make([]map[string]interface{}, 0, len(items))
	for _, c := range items {
		id, _ := c["certificateId"].(string)
		entry := map[string]interface{}{
			"certificateArn": iotstore.BuildCACertificateARN(store.GetAccountID(), store.GetRegion(), id),
			"certificateId":  c["certificateId"],
			"status":         c["status"],
		}
		switch created := c["creationDate"].(type) {
		case int64:
			entry["creationDate"] = created
		case float64:
			entry["creationDate"] = int64(created)
		}
		certs = append(certs, entry)
	}
	return &ListCACertificatesResult{Certificates: certs}, nil
}

// updateCACertificateCore applies a new status, auto-registration status or
// registration config to a registered CA certificate. Empty members leave the
// record untouched; non-empty enum values must be model members.
func (s *IoTService) updateCACertificateCore(store iotstore.IotStoreInterface, in UpdateCACertificateInput) error {
	if in.CertificateID == "" {
		return iotstore.ErrMissingParam
	}
	if _, ok := caCertificateStatuses[in.NewStatus]; !ok && in.NewStatus != "" {
		return iotstore.ErrInvalidRequest
	}
	if _, ok := caAutoRegistrationStatuses[in.NewAutoRegistrationStatus]; !ok && in.NewAutoRegistrationStatus != "" {
		return iotstore.ErrInvalidRequest
	}
	rec := map[string]interface{}{}
	caExists, err := store.GetGenericExists("caCert/"+in.CertificateID, &rec)
	if err != nil {
		return err
	}
	if !caExists {
		return iotstore.ErrCertificateNotFound
	}
	changed := false
	if in.NewStatus != "" {
		rec["status"] = in.NewStatus
		changed = true
	}
	if in.NewAutoRegistrationStatus != "" {
		rec["autoRegistrationStatus"] = in.NewAutoRegistrationStatus
		changed = true
	}
	if in.RemoveAutoRegistration {
		rec["autoRegistrationStatus"] = "DISABLE"
		changed = true
	}
	if in.RegistrationConfig != "" {
		rec["registrationConfig"] = in.RegistrationConfig
		changed = true
	}
	if changed {
		rec["lastModifiedDate"] = time.Now().Unix()
		if err := store.PutGeneric("caCert/"+in.CertificateID, rec); err != nil {
			return err
		}
	}
	return nil
}

// deleteCACertificateCore removes a registered CA certificate. AWS requires
// the CA certificate to be INACTIVE before deletion.
func (s *IoTService) deleteCACertificateCore(store iotstore.IotStoreInterface, caCertID string) error {
	if caCertID == "" {
		return iotstore.ErrMissingParam
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("caCert/"+caCertID, &rec)
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrCertificateNotFound
	}
	if status, _ := rec["status"].(string); status == "ACTIVE" {
		return iotstore.ErrCertHasAttachments
	}
	arn := iotstore.BuildCACertificateARN(store.GetAccountID(), store.GetRegion(), caCertID)
	_ = store.DeleteAllTags(arn)
	return store.DeleteGeneric("caCert/" + caCertID)
}

// listCertificatesByCACore lists the certificates issued by one CA.
func (s *IoTService) listCertificatesByCACore(store iotstore.IotStoreInterface, in ListCertificatesByCACertificateInput) (*ListCertificatesByCACertificateResult, error) {
	if in.CACertificateID == "" {
		return nil, iotstore.ErrMissingParam
	}
	maxItems := in.MaxItems
	if maxItems <= 0 {
		maxItems = 100
	}
	certs, err := store.ListCertificates(iotstoreListOpts(maxItems, in.Marker))
	if err != nil {
		return nil, err
	}
	items := make([]*iotstore.Certificate, 0)
	for _, c := range certs.Items {
		if c.CaCertificateID == in.CACertificateID {
			items = append(items, c)
		}
	}
	return &ListCertificatesByCACertificateResult{
		Certificates: items,
		NextMarker:   certs.NextMarker,
	}, nil
}

// registerCertificateWithoutCACore registers a PEM certificate directly,
// defaulting its status to INACTIVE.
func (s *IoTService) registerCertificateWithoutCACore(store iotstore.IotStoreInterface, in RegisterCertificateWithoutCAInput) (*RegisterCertificateWithoutCAResult, error) {
	if in.CertificatePEM == "" {
		return nil, iotstore.ErrMissingParam
	}
	certID := vcrypto.FingerprintPEM(in.CertificatePEM)
	status := in.Status
	if status == "" {
		status = "INACTIVE"
	} else if err := validateRegistrationStatus(status); err != nil {
		return nil, err
	}
	cert := &iotstore.Certificate{
		CertificateID:    certID,
		CertificatePEM:   in.CertificatePEM,
		Status:           status,
		CertificateMode:  "DEFAULT",
		CreationDate:     time.Now().UTC(),
		LastModifiedDate: time.Now().UTC(),
	}
	created, err := store.CreateCertificate(cert)
	if err != nil {
		return nil, err
	}
	return &RegisterCertificateWithoutCAResult{Certificate: created}, nil
}

// transferCertificateCore moves a certificate into PENDING_TRANSFER and
// records the outgoing transfer. AWS rejects the transfer of a non-existent
// certificate with ResourceNotFoundException rather than silently creating a
// transfer record.
func (s *IoTService) transferCertificateCore(store iotstore.IotStoreInterface, in TransferCertificateInput) (*TransferCertificateResult, error) {
	if in.CertificateID == "" || in.TargetAWSAccount == "" {
		return nil, iotstore.ErrMissingParam
	}
	if _, err := store.GetCertificate(in.CertificateID); err != nil {
		return nil, err
	}
	// Move the certificate into PENDING_TRANSFER so DescribeCertificate
	// reflects the transfer, matching AWS behaviour.
	if _, err := store.UpdateCertificate(in.CertificateID, iotstore.CertificateUpdateOpts{NewStatus: "PENDING_TRANSFER"}); err != nil {
		return nil, err
	}
	rec := map[string]interface{}{
		"certificateId": in.CertificateID,
		"status":        "PENDING_ACCEPTANCE",
		"transferredTo": in.TargetAWSAccount,
		"transferDate":  time.Now().UTC().Unix(),
	}
	if in.TransferMessage != "" {
		rec["transferMessage"] = in.TransferMessage
	}
	if err := store.PutGeneric("certTransfer/"+in.CertificateID, rec); err != nil {
		return nil, err
	}
	return &TransferCertificateResult{
		TransferredCertificateARN: iotstore.BuildCertificateARN(store.GetAccountID(), store.GetRegion(), in.CertificateID),
	}, nil
}

// listOutgoingCertificatesCore lists the transfers awaiting acceptance by
// their target accounts.
func (s *IoTService) listOutgoingCertificatesCore(store iotstore.IotStoreInterface) (*ListOutgoingCertificatesResult, error) {
	items, err := store.ListGeneric("certTransfer/")
	if err != nil {
		return nil, err
	}
	certs := make([]map[string]interface{}, 0, len(items))
	for _, c := range items {
		status, _ := c["status"].(string)
		if status != "PENDING_ACCEPTANCE" {
			continue
		}
		certID, _ := c["certificateId"].(string)
		entry := map[string]interface{}{
			"certificateArn": iotstore.BuildCertificateARN(store.GetAccountID(), store.GetRegion(), certID),
			"certificateId":  certID,
		}
		for _, member := range []string{"transferredTo", "transferDate", "transferMessage"} {
			if v, ok := c[member]; ok {
				entry[member] = v
			}
		}
		if cert, err := store.GetCertificate(certID); err == nil && cert != nil {
			entry["creationDate"] = cert.CreationDate.Unix()
		}
		certs = append(certs, entry)
	}
	return &ListOutgoingCertificatesResult{OutgoingCertificates: certs}, nil
}

// completeCertTransferCore resolves a pending certificate transfer: it sets
// the certificate to the terminal status (ACTIVE or INACTIVE per the caller)
// and removes the transfer record. A missing transfer record yields
// ResourceNotFoundException, matching AWS.
func (s *IoTService) completeCertTransferCore(store iotstore.IotStoreInterface, certID, certStatus string) error {
	if certID == "" {
		return iotstore.ErrMissingParam
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("certTransfer/"+certID, &rec)
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrCertificateNotFound
	}
	// Transition the certificate status and clear the transfer record.
	if _, err := store.UpdateCertificate(certID, iotstore.CertificateUpdateOpts{NewStatus: certStatus}); err != nil {
		return err
	}
	return store.DeleteGeneric("certTransfer/" + certID)
}
