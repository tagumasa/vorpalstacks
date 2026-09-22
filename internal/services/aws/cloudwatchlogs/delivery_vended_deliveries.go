package cloudwatchlogs

import (
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// The delivery family proper: creating, shaping, updating and deleting
// deliveries.

// --- Delivery ---

// CreateDeliveryInput holds the parsed parameters of CreateDelivery.
type CreateDeliveryInput struct {
	DeliverySourceName     string
	DeliveryDestinationArn string
	RecordFields           []string
	FieldDelimiter         string
	S3SuffixPath           string
	S3HiveCompatiblePath   bool
	S3ConfigPresent        bool
	Tags                   map[string]string
	Region                 string
}

// validateDeliveryShaping validates the record-shaping members the
// CreateDelivery and UpdateDeliveryConfiguration operations share.
func validateDeliveryShaping(recordFields []string, fieldDelimiter, suffixPath string, destType string, s3ConfigPresent bool) error {
	if len(recordFields) > logsstore.DeliveryRecordFieldsMax {
		return errDeliveryValidation("recordFields exceeds the maximum of %d entries", logsstore.DeliveryRecordFieldsMax)
	}
	for _, f := range recordFields {
		if f == "" || len(f) > logsstore.DeliveryRecordFieldLenMax {
			return errDeliveryValidation("invalid record field %q: 1-%d characters each", f, logsstore.DeliveryRecordFieldLenMax)
		}
		if !deliveryRecordFieldVocabulary[f] {
			return errDeliveryValidation("record field %q is not one of this source's allowed fields (time, stream, message)", f)
		}
	}
	if len(fieldDelimiter) > logsstore.DeliveryFieldDelimiterMax {
		return errDeliveryValidation("fieldDelimiter exceeds the maximum length of %d characters", logsstore.DeliveryFieldDelimiterMax)
	}
	if fieldDelimiter != "" && !deliveryFieldDelimiterAllowed(fieldDelimiter) {
		return errDeliveryValidation("fieldDelimiter %q is not one of the allowed field delimiters", fieldDelimiter)
	}
	if suffixPath != "" {
		if len(suffixPath) > logsstore.DeliverySuffixPathMax {
			return errDeliveryValidation("suffixPath exceeds the maximum length of %d characters", logsstore.DeliverySuffixPathMax)
		}
	}
	// The s3DeliveryConfiguration structure is "valid only when the
	// delivery's delivery destination is an S3 bucket" (member
	// documentation) — its presence alone rejects on another type.
	if s3ConfigPresent && destType != "S3" {
		return errDeliveryValidation("s3DeliveryConfiguration is valid only when the delivery destination is an S3 bucket")
	}
	return nil
}

func deliveryFieldDelimiterAllowed(d string) bool {
	for _, allowed := range deliveryAllowedFieldDelimiters {
		if d == allowed {
			return true
		}
	}
	return false
}

func (s *LogsService) createDeliveryCore(input *CreateDeliveryInput) (*logsstore.Delivery, error) {
	if input.DeliverySourceName == "" {
		return nil, errValidationMember("deliverySourceName")
	}
	if input.DeliveryDestinationArn == "" {
		return nil, errValidationMember("deliveryDestinationArn")
	}
	if err := validateDeliveryTags(input.Tags); err != nil {
		return nil, err
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetDeliverySource(input.DeliverySourceName); err != nil {
		return nil, mapStoreError(err)
	}
	destName := svcarn.ExtractDeliveryDestinationNameFromARN(input.DeliveryDestinationArn)
	if destName == "" {
		return nil, errDeliveryValidation("deliveryDestinationArn %q does not address a delivery destination", input.DeliveryDestinationArn)
	}
	dest, err := store.GetDeliveryDestination(destName)
	if err != nil {
		return nil, NewLogsError("ResourceNotFoundException",
			"The specified resource does not exist: "+input.DeliveryDestinationArn, 400)
	}

	if err := validateDeliveryShaping(input.RecordFields, input.FieldDelimiter, input.S3SuffixPath, dest.DeliveryDestinationType, input.S3ConfigPresent); err != nil {
		return nil, err
	}

	id, err := newDeliveryId()
	if err != nil {
		return nil, err
	}
	delivery := &logsstore.Delivery{
		Id:                 id,
		Arn:                store.ARNBuilder().CloudWatch().Delivery(id),
		DeliverySourceName: input.DeliverySourceName,
		// The record carries the destination's platform ARN, not the
		// caller's spelling of it: the deliveryDestinationArn a request
		// may cast is free to differ from the platform's own (a foreign
		// account ID, a region the platform did not mint), and the
		// family's identity comparisons — the delete-time reference
		// guard and the (source, destination) pair uniqueness — are
		// equality comparisons against the platform-cast ARN the
		// destination record carries.
		DeliveryDestinationArn:  dest.Arn,
		DeliveryDestinationType: dest.DeliveryDestinationType,
		RecordFields:            input.RecordFields,
		FieldDelimiter:          input.FieldDelimiter,
		S3SuffixPath:            input.S3SuffixPath,
		S3HiveCompatiblePath:    input.S3HiveCompatiblePath,
		Tags:                    input.Tags,
	}
	// The duplicate gate runs under the delivery record mutex with the
	// write: two concurrent creations of the same (source, destination)
	// pair admit exactly one.
	existed, err := store.PutDeliveryIfPairAbsent(delivery)
	if err != nil {
		return nil, mapStoreError(err)
	}
	if existed {
		return nil, errDeliveryConflict("A delivery between delivery source " + input.DeliverySourceName + " and this delivery destination already exists")
	}
	return delivery, nil
}

// newDeliveryId mints the delivery identifier: alphanumeric per the
// DeliveryId pattern, a 32-character random hex string.
func newDeliveryId() (string, error) {
	return newRandomHexID(16)
}

// UpdateDeliveryConfigurationInput holds the parsed parameters of
// UpdateDeliveryConfiguration.
type UpdateDeliveryConfigurationInput struct {
	Id                   string
	RecordFields         []string
	FieldDelimiter       string
	S3SuffixPath         string
	S3HiveCompatiblePath bool
	S3ConfigPresent      bool
	Region               string
}

func (s *LogsService) updateDeliveryConfigurationCore(input *UpdateDeliveryConfigurationInput) error {
	if input.Id == "" {
		return errValidationMember("id")
	}
	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return err
	}
	// The shaping validation and the field writes run inside the
	// delivery mutate seam: a concurrent delivery engine cursor advance
	// survives this update, and this update cannot revert a cursor that
	// committed while it was building.
	return mapStoreError(store.MutateDelivery(input.Id, func(delivery *logsstore.Delivery) error {
		if err := validateDeliveryShaping(input.RecordFields, input.FieldDelimiter, input.S3SuffixPath, delivery.DeliveryDestinationType, input.S3ConfigPresent); err != nil {
			return err
		}
		delivery.RecordFields = input.RecordFields
		delivery.FieldDelimiter = input.FieldDelimiter
		delivery.S3SuffixPath = input.S3SuffixPath
		delivery.S3HiveCompatiblePath = input.S3HiveCompatiblePath
		return nil
	}))
}

func (s *LogsService) deleteDeliveryCore(region, id string) error {
	if id == "" {
		return errValidationMember("id")
	}
	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return err
	}
	return mapStoreError(store.DeleteDelivery(id))
}
