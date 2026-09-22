package cloudwatchlogs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"unicode/utf8"

	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// The delivery-destination family: destination and destination-policy
// cores, including the S3-existence admission and the reference guard.

// --- Delivery destination ---

// PutDeliveryDestinationInput holds the parsed parameters of
// PutDeliveryDestination.
type PutDeliveryDestinationInput struct {
	Name                    string
	OutputFormat            string
	DestinationResourceArn  string
	ExplicitDestinationType string
	Tags                    map[string]string
	Region                  string
}

// deliveryDestinationTypeFor derives the delivery destination type from
// the destination resource ARN's service, and validates the explicit
// member against it when the request carries one. The derived and
// explicit forms must agree; a destination resource of a service with no
// delivery substrate (Firehose, X-Ray) rejects here — the never-accept-
// and-drop rule.
func deliveryDestinationTypeFor(destinationResourceArn, explicitType string) (string, error) {
	parsed, err := svcarn.ParseARN(destinationResourceArn)
	if err != nil {
		return "", errDeliveryValidation("destinationResourceArn %q is not a valid ARN", destinationResourceArn)
	}
	derived := ""
	switch parsed.Service {
	case "logs":
		if svcarn.ExtractLogGroupNameFromARN(destinationResourceArn) == "" {
			return "", errDeliveryValidation("destinationResourceArn %q does not address a CloudWatch Logs log group", destinationResourceArn)
		}
		derived = "CWL"
	case "s3":
		derived = "S3"
	case "firehose":
		return "", errDeliveryValidation("Firehose delivery destinations are not supported until the platform Firehose service exists")
	default:
		return "", errDeliveryValidation("destinationResourceArn %q does not address a supported delivery destination resource (a CloudWatch Logs log group or an S3 bucket)", destinationResourceArn)
	}
	if explicitType != "" && explicitType != derived {
		// X-Ray trace destinations carry no destinationResourceArn; the
		// platform rejects them with the same identity the service form
		// rejects with.
		if explicitType == "XRAY" || explicitType == "FH" {
			return "", errDeliveryValidation("deliveryDestinationType %s is not supported on this platform", explicitType)
		}
		return "", errDeliveryValidation("deliveryDestinationType %s does not match the destination resource %s", explicitType, derived)
	}
	return derived, nil
}

func (s *LogsService) putDeliveryDestinationCore(input *PutDeliveryDestinationInput) (*logsstore.DeliveryDestination, error) {
	if input.Name == "" {
		return nil, errValidationMember("name")
	}
	if len(input.Name) > logsstore.DeliveryNameMax || !deliveryNameRe.MatchString(input.Name) {
		return nil, errDeliveryValidation("invalid delivery destination name %q: 1-%d characters, word characters and hyphens only",
			input.Name, logsstore.DeliveryNameMax)
	}
	if input.OutputFormat != "" {
		switch input.OutputFormat {
		case "json", "plain", "w3c", "raw":
		case "parquet":
			// Parquet output stays excluded as a scope boundary: the
			// columnar writer the platform already runs for S3 inventory
			// reports is not carried into the delivery family, and a
			// destination configured for it could never deliver, so it
			// rejects at Put.
			return nil, errDeliveryValidation("outputFormat parquet is not supported on this platform")
		default:
			return nil, errDeliveryValidation("invalid outputFormat %q: valid values are json, plain, w3c, raw", input.OutputFormat)
		}
	}
	if input.DestinationResourceArn == "" {
		return nil, errValidationMember("deliveryDestinationConfiguration.destinationResourceArn")
	}
	if err := validateDeliveryTags(input.Tags); err != nil {
		return nil, err
	}

	destType, err := deliveryDestinationTypeFor(input.DestinationResourceArn, input.ExplicitDestinationType)
	if err != nil {
		return nil, err
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return nil, err
	}
	if destType == "CWL" {
		// The delivery engine writes to the destination group in the
		// region the resource ARN addresses (falling back to the
		// delivery's home region when the ARN carries none), so the
		// admission check must address the same store — checking the
		// call's region would validate a different region's group than
		// the one every delivery pass writes to. The check's store is
		// scoped to the check alone: the destination record itself is a
		// resource of the call's region (its ARN and every read of the
		// family resolve there), so persisting through the check's store
		// would file the destination in the destination resource's
		// region where no Get/Describe of this region can find it.
		admissionStore := store
		if destRegion := deliveryResourceRegion(input.DestinationResourceArn, input.Region); destRegion != input.Region {
			admissionStore, err = s.getLogsStoreByRegion(destRegion)
			if err != nil {
				return nil, err
			}
		}
		groupName := resolveLogGroupIdentifier(input.DestinationResourceArn)
		if _, err := admissionStore.GetLogGroup(groupName); err != nil {
			return nil, NewLogsError("ResourceNotFoundException",
				"The specified resource does not exist: "+input.DestinationResourceArn, 400)
		}
	} else if destType == "S3" {
		if s.eventBus() == nil || s.eventBus().S3Invoker() == nil {
			// The bucket's existence is part of the family's admission
			// contract ("Create a delivery destination for an existing
			// bucket", vended-logs setup): admitting a destination whose
			// bucket the platform cannot verify would accept a delivery
			// that could never flow — the engine refuses S3 writes without
			// the invoker. The unwired substrate state rejects here with
			// the operation's declared service-side identity instead of
			// surfacing per pass in the engine.
			return nil, NewLogsError("ServiceUnavailableException",
				"The S3 delivery substrate is not available to verify the destination bucket", http.StatusServiceUnavailable)
		}
		bucket, ok := s3BucketFromArn(input.DestinationResourceArn)
		if !ok {
			return nil, errDeliveryValidation("destinationResourceArn %q does not address an S3 bucket", input.DestinationResourceArn)
		}
		exists, err := s.eventBus().S3Invoker().BucketExists(context.Background(), input.Region, bucket)
		if err != nil {
			// A failed existence check must not admit the bucket as
			// existing: propagate the backend failure instead of
			// swallowing it into a silent pass.
			return nil, fmt.Errorf("bucket existence check for %s failed: %w", bucket, err)
		}
		if !exists {
			return nil, NewLogsError("ResourceNotFoundException",
				"The specified resource does not exist: "+input.DestinationResourceArn, 400)
		}
	}

	destination := &logsstore.DeliveryDestination{
		Name:                    input.Name,
		Arn:                     store.ARNBuilder().CloudWatch().DeliveryDestination(input.Name),
		DeliveryDestinationType: destType,
		OutputFormat:            input.OutputFormat,
		DestinationResourceArn:  input.DestinationResourceArn,
		Tags:                    input.Tags,
	}
	if err := store.PutDeliveryDestination(destination); err != nil {
		return nil, mapStoreError(err)
	}
	return destination, nil
}

func (s *LogsService) deleteDeliveryDestinationCore(region, name string) error {
	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return err
	}
	dest, err := store.GetDeliveryDestination(name)
	if err != nil {
		return mapStoreError(err)
	}
	deliveries, err := store.ListDeliveries()
	if err != nil {
		return mapStoreError(err)
	}
	if err := refuseDeliveryFamilyIfReferenced(deliveries,
		func(d *logsstore.Delivery) bool { return d.DeliveryDestinationArn == dest.Arn },
		"Delivery destination", name); err != nil {
		return err
	}
	return mapStoreError(store.DeleteDeliveryDestination(name))
}

// --- Delivery destination policy ---

func (s *LogsService) putDeliveryDestinationPolicyCore(region, name, policyDocument string) (*logsstore.DeliveryDestinationPolicy, error) {
	if name == "" {
		return nil, errValidationMember("deliveryDestinationName")
	}
	if policyDocument == "" {
		return nil, errValidationMember("deliveryDestinationPolicy")
	}
	// The delivery destination policy rides the same PolicyDocument
	// trait the log plane's policy family validates against (1..51200
	// Unicode characters): the one constant, counted the one way.
	if utf8.RuneCountInString(policyDocument) > logsstore.MaxPolicyDocumentLength {
		return nil, errDeliveryValidation("deliveryDestinationPolicy exceeds the maximum length of %d characters", logsstore.MaxPolicyDocumentLength)
	}
	// The family's declared error list carries ValidationException (not
	// the log-plane InvalidParameterException the shared JSON validator
	// raises), so the document's JSON validity rejects here with the
	// family's own identity.
	if !json.Valid([]byte(policyDocument)) {
		return nil, errDeliveryValidation("deliveryDestinationPolicy must be valid JSON")
	}

	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetDeliveryDestination(name); err != nil {
		return nil, mapStoreError(err)
	}
	policy := &logsstore.DeliveryDestinationPolicy{PolicyDocument: policyDocument}
	if err := store.PutDeliveryDestinationPolicy(name, policy); err != nil {
		return nil, mapStoreError(err)
	}
	return policy, nil
}

func (s *LogsService) getDeliveryDestinationPolicyCore(region, name string) (*logsstore.DeliveryDestinationPolicy, error) {
	if name == "" {
		return nil, errValidationMember("deliveryDestinationName")
	}
	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetDeliveryDestination(name); err != nil {
		return nil, mapStoreError(err)
	}
	policy, err := store.GetDeliveryDestinationPolicy(name)
	if err != nil {
		return nil, mapStoreError(err)
	}
	return policy, nil
}

func (s *LogsService) deleteDeliveryDestinationPolicyCore(region, name string) error {
	if name == "" {
		return errValidationMember("deliveryDestinationName")
	}
	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return err
	}
	if _, err := store.GetDeliveryDestination(name); err != nil {
		return mapStoreError(err)
	}
	return mapStoreError(store.DeleteDeliveryDestinationPolicy(name))
}
