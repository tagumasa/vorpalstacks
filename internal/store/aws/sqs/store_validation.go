package sqs

import (
	"regexp"
	"strings"
)

var (
	queueNameRegex      = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	fifoQueueNameRegex  = regexp.MustCompile(`^[a-zA-Z0-9_-]+\.fifo$`)
	batchEntryIdRegex   = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	validAttributeNames = map[string]bool{
		"All":                                   true,
		"QueueArn":                              true,
		"ApproximateNumberOfMessages":           true,
		"ApproximateNumberOfMessagesDelayed":    true,
		"ApproximateNumberOfMessagesNotVisible": true,
		"CreatedTimestamp":                      true,
		"LastModifiedTimestamp":                 true,
		"VisibilityTimeout":                     true,
		"MaximumMessageSize":                    true,
		"MessageRetentionPeriod":                true,
		"DelaySeconds":                          true,
		"ReceiveMessageWaitTimeSeconds":         true,
		"Policy":                                true,
		"RedrivePolicy":                         true,
		"FifoQueue":                             true,
		"ContentBasedDeduplication":             true,
		"KmsMasterKeyId":                        true,
		"KmsDataKeyReusePeriodSeconds":          true,
		"DeduplicationScope":                    true,
		"FifoThroughputLimit":                   true,
		"RedriveAllowPolicy":                    true,
		"SqsManagedSseEnabled":                  true,
	}
)

func validateQueueName(name string) error {
	if len(name) == 0 {
		return ErrInvalidQueueName
	}
	if len(name) > maxQueueNameLength {
		return ErrInvalidQueueName
	}
	if strings.HasSuffix(name, ".fifo") {
		if !fifoQueueNameRegex.MatchString(name) {
			return ErrInvalidQueueName
		}
	} else {
		if !queueNameRegex.MatchString(name) {
			return ErrInvalidQueueName
		}
	}
	return nil
}

// ValidateBatchEntryId validates a batch entry ID.
func ValidateBatchEntryId(id string) error {
	if len(id) == 0 || len(id) > maxBatchEntryIdLength {
		return ErrInvalidBatchEntryId
	}
	if !batchEntryIdRegex.MatchString(id) {
		return ErrInvalidBatchEntryId
	}
	return nil
}

// IsValidAttributeName checks if an attribute name is valid for SQS queues.
func IsValidAttributeName(name string) bool {
	return validAttributeNames[name]
}

func validateVisibilityTimeout(value int32) error {
	if value < minVisibilityTimeout || value > maxVisibilityTimeout {
		return ErrInvalidParameterValue
	}
	return nil
}

func validateDelaySeconds(value int32) error {
	if value < minDelaySeconds || value > maxDelaySeconds {
		return ErrInvalidParameterValue
	}
	return nil
}

func validateMessageRetentionPeriod(value int32) error {
	if value < minMessageRetentionPeriod || value > maxMessageRetentionPeriod {
		return ErrInvalidParameterValue
	}
	return nil
}

func validateMaximumMessageSize(value int32) error {
	if value < minMaximumMessageSize || value > maxMaximumMessageSize {
		return ErrInvalidParameterValue
	}
	return nil
}

func validateReceiveMessageWaitTimeSeconds(value int32) error {
	if value < minReceiveMessageWaitTime || value > maxReceiveMessageWaitTime {
		return ErrInvalidParameterValue
	}
	return nil
}

func validateMessageBodySize(body string) error {
	if len(body) > maxMessageBodySize {
		return ErrMessageTooLarge
	}
	return nil
}

func validateTags(tags map[string]string) error {
	if len(tags) > maxTagsPerQueue {
		return ErrTooManyTags
	}
	for key, value := range tags {
		if len(key) > maxTagKeyLength {
			return ErrInvalidTagKey
		}
		if len(value) > maxTagValueLength {
			return ErrInvalidTagValue
		}
	}
	return nil
}
