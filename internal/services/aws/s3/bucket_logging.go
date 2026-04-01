package s3

import (
	"vorpalstacks/internal/services/aws/common/request"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// PutBucketLoggingInput contains the request parameters for the PutBucketLogging operation.
type PutBucketLoggingInput struct {
	Bucket               string
	LoggingConfiguration *LoggingConfigurationInput
}

// LoggingConfigurationInput defines the logging configuration for a bucket.
type LoggingConfigurationInput struct {
	TargetBucket string             `xml:"TargetBucket"`
	TargetPrefix string             `xml:"TargetPrefix,omitempty"`
	TargetGrants []TargetGrantInput `xml:"TargetGrants>Grant,omitempty"`
}

// TargetGrantInput specifies a grant for bucket logging.
type TargetGrantInput struct {
	Grantee    *TargetGranteeInput `xml:"Grantee,omitempty"`
	Permission string              `xml:"Permission"`
}

// TargetGranteeInput defines the grantee for a logging grant.
type TargetGranteeInput struct {
	Type         string `xml:"type,attr"`
	ID           string `xml:"ID,omitempty"`
	DisplayName  string `xml:"DisplayName,omitempty"`
	URI          string `xml:"URI,omitempty"`
	EmailAddress string `xml:"EmailAddress,omitempty"`
}

// PutBucketLogging sets the logging configuration for a bucket.
// This operation requires the bucket owner to have WRITE and READ_ACP permissions.
func (o *BucketOperations) PutBucketLogging(ctx *request.RequestContext, input *PutBucketLoggingInput) error {
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}

	var config *s3store.LoggingConfiguration
	if input.LoggingConfiguration != nil {
		config = &s3store.LoggingConfiguration{
			TargetBucket: input.LoggingConfiguration.TargetBucket,
			TargetPrefix: input.LoggingConfiguration.TargetPrefix,
		}

		for _, tg := range input.LoggingConfiguration.TargetGrants {
			grant := s3store.TargetGrant{
				Permission: s3store.Permission(tg.Permission),
			}
			if tg.Grantee != nil {
				grant.Grantee = &s3store.Grantee{
					Type:        s3store.GranteeType(tg.Grantee.Type),
					ID:          tg.Grantee.ID,
					DisplayName: tg.Grantee.DisplayName,
					URI:         tg.Grantee.URI,
					Email:       tg.Grantee.EmailAddress,
				}
			}
			config.TargetGrants = append(config.TargetGrants, grant)
		}
	}

	return store.buckets.SetLoggingConfiguration(input.Bucket, config)
}

// GetBucketLoggingInput contains the request parameters for the GetBucketLogging operation.
type GetBucketLoggingInput struct {
	Bucket string
}

// GetBucketLoggingOutput contains the result of the GetBucketLogging operation.
type GetBucketLoggingOutput struct {
	LoggingConfiguration *LoggingConfigurationOutput `xml:"LoggingEnabled"`
}

// LoggingConfigurationOutput defines the logging configuration for a bucket.
type LoggingConfigurationOutput struct {
	TargetBucket string              `xml:"TargetBucket"`
	TargetPrefix string              `xml:"TargetPrefix,omitempty"`
	TargetGrants []TargetGrantOutput `xml:"TargetGrants>Grant,omitempty"`
}

// TargetGrantOutput specifies a grant for bucket logging.
type TargetGrantOutput struct {
	Grantee    *TargetGranteeOutput `xml:"Grantee,omitempty"`
	Permission string               `xml:"Permission"`
}

// TargetGranteeOutput defines the grantee for a logging grant.
type TargetGranteeOutput struct {
	Type         string `xml:"xsi:type,attr"`
	ID           string `xml:"ID,omitempty"`
	DisplayName  string `xml:"DisplayName,omitempty"`
	URI          string `xml:"URI,omitempty"`
	EmailAddress string `xml:"EmailAddress,omitempty"`
}

// GetBucketLogging retrieves the logging configuration for a bucket.
func (o *BucketOperations) GetBucketLogging(ctx *request.RequestContext, input *GetBucketLoggingInput) (*GetBucketLoggingOutput, error) {
	store, err := o.svc.store(ctx)
	if err != nil {
		return nil, err
	}

	config, err := store.buckets.GetLoggingConfiguration(input.Bucket)
	if err != nil {
		return nil, err
	}

	if config == nil {
		return &GetBucketLoggingOutput{}, nil
	}

	output := &GetBucketLoggingOutput{
		LoggingConfiguration: &LoggingConfigurationOutput{
			TargetBucket: config.TargetBucket,
			TargetPrefix: config.TargetPrefix,
		},
	}

	for _, tg := range config.TargetGrants {
		grantOut := TargetGrantOutput{
			Permission: string(tg.Permission),
		}
		if tg.Grantee != nil {
			grantOut.Grantee = &TargetGranteeOutput{
				Type:         string(tg.Grantee.Type),
				ID:           tg.Grantee.ID,
				DisplayName:  tg.Grantee.DisplayName,
				URI:          tg.Grantee.URI,
				EmailAddress: tg.Grantee.Email,
			}
		}
		output.LoggingConfiguration.TargetGrants = append(output.LoggingConfiguration.TargetGrants, grantOut)
	}

	return output, nil
}
