// Transport-agnostic Core functions for IAM SAML providers: validation and
// store operations shared by the AWS-compatible HTTP API handlers and any
// admin surface (the xxxCore pattern).
package iam

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"time"
	"unicode/utf8"

	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// CreateSAMLProviderInput holds the parameters for creating a SAML provider.
type CreateSAMLProviderInput struct {
	Name                    string
	SAMLMetadataDocument    string
	AssertionEncryptionMode string
	AddPrivateKey           string
	Tags                    []tags.Tag
}

// UpdateSAMLProviderInput holds the parameters for updating a SAML provider.
type UpdateSAMLProviderInput struct {
	ProviderArn             string
	SAMLMetadataDocument    string
	AssertionEncryptionMode string
	AddPrivateKey           string
	RemovePrivateKey        string
}

// extractValidUntilFromSAMLMetadata extracts the X.509 certificate expiry
// from a SAML metadata document, if one is present and parseable.
func extractValidUntilFromSAMLMetadata(metadata string) *time.Time {
	matches := x509CertDataPattern.FindStringSubmatch(metadata)
	if len(matches) < 2 {
		return nil
	}
	certData := whitespacePattern.ReplaceAllString(matches[1], "")

	derBytes, err := base64.StdEncoding.DecodeString(certData)
	if err != nil {
		pemBlock, _ := pem.Decode([]byte("-----BEGIN CERTIFICATE-----\n" + certData + "\n-----END CERTIFICATE-----"))
		if pemBlock == nil {
			return nil
		}
		derBytes = pemBlock.Bytes
	}

	cert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		return nil
	}

	notAfter := cert.NotAfter.UTC()
	return &notAfter
}

// createSAMLProviderCore validates input and creates a SAML provider,
// returning its ARN.
func (s *IAMService) createSAMLProviderCore(store *iamstore.IAMStore, input *CreateSAMLProviderInput) (string, error) {
	if input.Name == "" {
		return "", NewValidationError("Name")
	}
	if err := validateSAMLProviderName(input.Name); err != nil {
		return "", err
	}
	if input.SAMLMetadataDocument == "" {
		return "", NewValidationError("SAMLMetadataDocument")
	}
	// SAMLMetadataDocumentType @length(1000,10000000) counts Unicode
	// characters (no pattern — XML metadata may carry multibyte text).
	if n := utf8.RuneCountInString(input.SAMLMetadataDocument); n < 1000 || n > 10000000 {
		return "", NewInvalidInputError("SAMLMetadataDocument", "must be between 1000 and 10000000 characters")
	}

	if err := validateNewTags(input.Tags); err != nil {
		return "", err
	}

	if input.AssertionEncryptionMode != "" && !validateAssertionEncryptionMode(input.AssertionEncryptionMode) {
		return "", NewInvalidInputError("AssertionEncryptionMode", "must be 'Required' or 'Allowed'")
	}

	if input.AddPrivateKey != "" {
		// privateKeyType carries a Latin-1 pattern, so lengths count Unicode
		// characters.
		if utf8.RuneCountInString(input.AddPrivateKey) > maxPrivateKeyLength {
			return "", NewInvalidInputError("AddPrivateKey", fmt.Sprintf("must be 1 to %d characters", maxPrivateKeyLength))
		}
	}

	validUntil := extractValidUntilFromSAMLMetadata(input.SAMLMetadataDocument)
	provider, err := store.SAMLProviders().Create(input.Name, input.SAMLMetadataDocument, validUntil, input.AssertionEncryptionMode, input.AddPrivateKey, input.Tags)
	if err != nil {
		if errors.Is(err, iamstore.ErrSAMLProviderAlreadyExists) {
			return "", NewEntityAlreadyExistsError("SAML Provider " + input.Name)
		}
		return "", err
	}
	return provider.Arn, nil
}

// getSAMLProviderCore returns the SAML provider with the given ARN.
func (s *IAMService) getSAMLProviderCore(store *iamstore.IAMStore, providerArn string) (*iamstore.SAMLProvider, error) {
	if providerArn == "" {
		return nil, NewValidationError("SAMLProviderArn")
	}
	if err := validateARNParameter("SAMLProviderArn", providerArn); err != nil {
		return nil, err
	}
	provider, err := store.SAMLProviders().Get(providerArn)
	if err != nil {
		return nil, NewNoSuchEntityError("SAML provider", providerArn)
	}
	return provider, nil
}

// listSAMLProvidersCore lists the SAML providers in the account.
func (s *IAMService) listSAMLProvidersCore(store *iamstore.IAMStore) ([]*iamstore.SAMLProvider, error) {
	result, err := store.SAMLProviders().List()
	if err != nil {
		return nil, err
	}
	return result.SAMLProviders, nil
}

// updateSAMLProviderCore validates input and updates the metadata document
// of the given SAML provider, returning its ARN.
func (s *IAMService) updateSAMLProviderCore(store *iamstore.IAMStore, input *UpdateSAMLProviderInput) (string, error) {
	if input.ProviderArn == "" {
		return "", NewValidationError("SAMLProviderArn")
	}
	if err := validateARNParameter("SAMLProviderArn", input.ProviderArn); err != nil {
		return "", err
	}
	if input.SAMLMetadataDocument != "" {
		// SAMLMetadataDocumentType @length(1000,10000000) counts Unicode
		// characters.
		if n := utf8.RuneCountInString(input.SAMLMetadataDocument); n < 1000 || n > 10000000 {
			return "", NewInvalidInputError("SAMLMetadataDocument", "must be between 1000 and 10000000 characters")
		}
	}

	if input.AssertionEncryptionMode != "" && !validateAssertionEncryptionMode(input.AssertionEncryptionMode) {
		return "", NewInvalidInputError("AssertionEncryptionMode", "must be 'Required' or 'Allowed'")
	}

	if input.AddPrivateKey != "" && utf8.RuneCountInString(input.AddPrivateKey) > maxPrivateKeyLength {
		return "", NewInvalidInputError("AddPrivateKey", fmt.Sprintf("must be 1 to %d characters", maxPrivateKeyLength))
	}

	// RemovePrivateKey is a privateKeyIdType (Smithy: length [22,64],
	// pattern ^[A-Z0-9]+$) identifying the key to remove.
	if input.RemovePrivateKey != "" {
		if len(input.RemovePrivateKey) < 22 || len(input.RemovePrivateKey) > 64 {
			return "", NewInvalidInputError("RemovePrivateKey", "must be 22 to 64 characters")
		}
		if err := validateSAMLPrivateKeyId(input.RemovePrivateKey); err != nil {
			return "", err
		}
	}

	if !store.SAMLProviders().Exists(input.ProviderArn) {
		return "", NewNoSuchEntityError("SAML provider", input.ProviderArn)
	}

	validUntil := extractValidUntilFromSAMLMetadata(input.SAMLMetadataDocument)
	if err := store.SAMLProviders().Update(input.ProviderArn, input.SAMLMetadataDocument, validUntil, input.AssertionEncryptionMode, input.AddPrivateKey, input.RemovePrivateKey); err != nil {
		if errors.Is(err, iamstore.ErrSAMLPrivateKeyNotFound) {
			return "", NewNoSuchEntityError("SAML private key", input.RemovePrivateKey)
		}
		return "", err
	}
	return input.ProviderArn, nil
}

// deleteSAMLProviderCore deletes the SAML provider with the given ARN.
func (s *IAMService) deleteSAMLProviderCore(store *iamstore.IAMStore, providerArn string) error {
	if providerArn == "" {
		return NewValidationError("SAMLProviderArn")
	}
	if err := validateARNParameter("SAMLProviderArn", providerArn); err != nil {
		return err
	}
	if !store.SAMLProviders().Exists(providerArn) {
		return NewNoSuchEntityError("SAML provider", providerArn)
	}
	return store.SAMLProviders().Delete(providerArn)
}
