package s3

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/common/defaults"
	"vorpalstacks/internal/common/request"
	types "vorpalstacks/internal/common/tags"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// Core functions for the bucket sub-resource configuration operations
// (accelerate, ACL, CORS, encryption, lifecycle, location, logging,
// notification, object lock, ownership controls, policy, public access
// block, request payment, tagging, versioning, website). The operation
// methods on BucketOperations are thin adapters that acquire the
// per-region store and delegate all validation and persistence here, so
// the HTTP plane and any other caller share a single behavioural path.

// putBucketAccelerateConfigurationCore validates and persists the transfer
// acceleration configuration for a bucket.
func (s *S3Service) putBucketAccelerateConfigurationCore(bucketStore s3store.BucketStoreInterface, in *PutBucketAccelerateConfigurationInput) error {
	if err := validateAccelerateStatus(in.AccelerateConfiguration.Status); err != nil {
		return err
	}

	config := &s3store.AccelerateConfiguration{
		Status: in.AccelerateConfiguration.Status,
	}

	return bucketStore.SetAccelerateConfiguration(in.Bucket, config)
}

// getBucketAccelerateConfigurationCore reads the transfer acceleration
// configuration for a bucket.
func (s *S3Service) getBucketAccelerateConfigurationCore(bucketStore s3store.BucketStoreInterface, in *GetBucketAccelerateConfigurationInput) (*GetBucketAccelerateConfigurationOutput, error) {
	config, err := bucketStore.GetAccelerateConfiguration(in.Bucket)
	if err != nil {
		return nil, err
	}

	if config == nil {
		return &GetBucketAccelerateConfigurationOutput{
			AccelerateConfiguration: &AccelerateConfigurationInput{},
		}, nil
	}

	return &GetBucketAccelerateConfigurationOutput{
		AccelerateConfiguration: &AccelerateConfigurationInput{
			Status: config.Status,
		},
	}, nil
}

// putBucketAclCore resolves the requested ACL (canned ACL, access control
// policy body, or grant headers), enforces BlockPublicAcls and the
// BucketOwnerEnforced ownership rule, and persists the bucket ACL.
func (s *S3Service) putBucketAclCore(ctx *request.RequestContext, store *s3Stores, in *PutBucketAclInput) error {
	owner := &s3store.ACLOwner{ID: s.accountID, DisplayName: s.accountID}

	var acp *s3store.AccessControlPolicy
	var err error

	if in.ACL != "" {
		acp, err = CannedACLToPolicy(in.ACL, owner)
		if err != nil {
			return err
		}
	} else if in.AccessControlPolicy != nil {
		acp = in.AccessControlPolicy
	} else {
		grants, err := ParseGrantHeaders(in.GrantFullControl, in.GrantRead, in.GrantReadACP, in.GrantWrite, in.GrantWriteACP)
		if err != nil {
			return NewInvalidArgumentError(err.Error())
		}
		if len(grants) > 0 {
			acp = &s3store.AccessControlPolicy{Owner: owner, Grants: grants}
		} else {
			return NewInvalidArgumentError("missing required ACL specification")
		}
	}

	publicAccessBlock, _ := store.buckets.GetPublicAccessBlock(in.Bucket)
	if publicAccessBlock != nil && publicAccessBlock.BlockPublicAcls {
		if isPublicCannedACL(in.ACL) {
			return NewInvalidArgumentError("bucket has BlockPublicAcls enabled")
		}
		if acpContainsPublicAccess(acp) {
			return NewInvalidArgumentError("bucket has BlockPublicAcls enabled")
		}
	}

	// With Object Ownership set to BucketOwnerEnforced, "requests to set or
	// update ACLs fail" with AccessControlListNotSupported.
	if aclsDisabled, _ := s.bucketACLsDisabled(ctx, store, in.Bucket); aclsDisabled {
		return ErrAccessControlListNotSupported
	}

	return store.buckets.SetACL(in.Bucket, acp)
}

// getBucketAclCore reads the ACL for a bucket, materialising the default
// owner full-control grant when no ACL has been set.
func (s *S3Service) getBucketAclCore(bucketStore s3store.BucketStoreInterface, bucket string) (*GetBucketAclOutput, error) {
	b, err := bucketStore.Get(bucket)
	if err != nil {
		return nil, err
	}

	owner := &s3store.ACLOwner{ID: s.accountID, DisplayName: s.accountID}

	if b.ACL == nil {
		return &GetBucketAclOutput{
			Owner: owner,
			Grants: []*s3store.Grant{
				{
					Grantee:    &s3store.Grantee{Type: s3store.GranteeTypeCanonicalUser, ID: s.accountID, DisplayName: s.accountID},
					Permission: s3store.PermissionFullControl,
				},
			},
		}, nil
	}

	return &GetBucketAclOutput{
		Owner:  b.ACL.Owner,
		Grants: b.ACL.Grants,
	}, nil
}

// bucketACLsDisabled reports whether the bucket's Object Ownership setting
// is BucketOwnerEnforced, under which ACLs are disabled and set/update ACL
// requests fail.
func (s *S3Service) bucketACLsDisabled(ctx context.Context, store *s3Stores, bucket string) (bool, error) {
	b, err := store.buckets.Get(bucket)
	if err != nil {
		return false, err
	}
	return b.OwnershipControls != nil &&
		len(b.OwnershipControls.Rules) == 1 &&
		b.OwnershipControls.Rules[0].ObjectOwnership == "BucketOwnerEnforced", nil
}

// putBucketCORSCore validates the CORS rules and persists the CORS
// configuration for a bucket.
func (s *S3Service) putBucketCORSCore(bucketStore s3store.BucketStoreInterface, in *PutBucketCORSInput) error {
	var rules []s3store.CORSRule
	for _, r := range in.CORSConfiguration.CORSRules {
		if len(r.AllowedMethods) == 0 {
			return NewInvalidArgumentError("CORS rule must have at least one AllowedMethod")
		}
		if len(r.AllowedOrigins) == 0 {
			return NewInvalidArgumentError("CORS rule must have at least one AllowedOrigin")
		}
		for _, method := range r.AllowedMethods {
			if !validCORSMethods[method] {
				return NewInvalidArgumentError(fmt.Sprintf("invalid CORS method: %s (must be GET, PUT, HEAD, POST, or DELETE)", method))
			}
		}
		rules = append(rules, s3store.CORSRule{
			AllowedHeaders: r.AllowedHeaders,
			AllowedMethods: r.AllowedMethods,
			AllowedOrigins: r.AllowedOrigins,
			ExposeHeaders:  r.ExposeHeaders,
			MaxAgeSeconds:  r.MaxAgeSeconds,
			ID:             r.ID,
		})
	}

	return bucketStore.SetCORS(in.Bucket, &s3store.CORSConfiguration{
		CORSRules: rules,
	})
}

// getBucketCORSCore reads the CORS configuration for a bucket and fails
// with ErrNoSuchCORS when none is set.
func (s *S3Service) getBucketCORSCore(bucketStore s3store.BucketStoreInterface, in *GetBucketCORSInput) (*GetBucketCORSOutput, error) {
	bucket, err := bucketStore.Get(in.Bucket)
	if err != nil {
		return nil, err
	}

	if bucket.CORSConfiguration == nil {
		return nil, ErrNoSuchCORS
	}

	var rules []CORSRuleInput
	for _, r := range bucket.CORSConfiguration.CORSRules {
		rules = append(rules, CORSRuleInput{
			AllowedHeaders: r.AllowedHeaders,
			AllowedMethods: r.AllowedMethods,
			AllowedOrigins: r.AllowedOrigins,
			ExposeHeaders:  r.ExposeHeaders,
			MaxAgeSeconds:  r.MaxAgeSeconds,
			ID:             r.ID,
		})
	}

	return &GetBucketCORSOutput{
		CORSConfiguration: &CORSConfigurationInput{
			CORSRules: rules,
		},
	}, nil
}

// deleteBucketCORSCore removes the CORS configuration from a bucket.
func (s *S3Service) deleteBucketCORSCore(bucketStore s3store.BucketStoreInterface, in *DeleteBucketCORSInput) error {
	return bucketStore.SetCORS(in.Bucket, nil)
}

// putBucketEncryptionCore validates the default encryption rule and
// persists the encryption configuration for a bucket.
func (s *S3Service) putBucketEncryptionCore(bucketStore s3store.BucketStoreInterface, in *PutBucketEncryptionInput) error {
	if in.ServerSideEncryptionConfiguration == nil || len(in.ServerSideEncryptionConfiguration.Rules) != 1 {
		return NewInvalidArgumentError("exactly one encryption rule is required")
	}

	rule := in.ServerSideEncryptionConfiguration.Rules[0]
	sseAlgorithm := rule.ApplyServerSideEncryptionByDefault.SSEAlgorithm
	if sseAlgorithm != "AES256" && sseAlgorithm != "aws:kms" && sseAlgorithm != "aws:kms:dsse" {
		return fmt.Errorf("invalid SSE algorithm: %s (must be AES256, aws:kms, or aws:kms:dsse)", sseAlgorithm)
	}

	if err := validateKMSMasterKeyID(rule.ApplyServerSideEncryptionByDefault.KMSMasterKeyID, sseAlgorithm); err != nil {
		return err
	}

	config := &s3store.EncryptionConfig{
		SSEAlgorithm:   rule.ApplyServerSideEncryptionByDefault.SSEAlgorithm,
		KMSMasterKeyID: rule.ApplyServerSideEncryptionByDefault.KMSMasterKeyID,
	}
	if rule.BucketKeyEnabled != nil {
		config.BucketKeyEnabled = rule.BucketKeyEnabled
	}

	return bucketStore.SetEncryption(in.Bucket, config)
}

// getBucketEncryptionCore reads the encryption configuration for a bucket
// and fails with ErrNoSuchEncryption when none is set.
func (s *S3Service) getBucketEncryptionCore(bucketStore s3store.BucketStoreInterface, in *GetBucketEncryptionInput) (*GetBucketEncryptionOutput, error) {
	bucket, err := bucketStore.Get(in.Bucket)
	if err != nil {
		return nil, err
	}

	if bucket.EncryptionConfig == nil {
		return nil, ErrNoSuchEncryption
	}

	return &GetBucketEncryptionOutput{
		ServerSideEncryptionConfiguration: &ServerSideEncryptionConfiguration{
			Rules: []ServerSideEncryptionRule{
				{
					ApplyServerSideEncryptionByDefault: ApplyServerSideEncryptionByDefault{
						SSEAlgorithm:   bucket.EncryptionConfig.SSEAlgorithm,
						KMSMasterKeyID: bucket.EncryptionConfig.KMSMasterKeyID,
					},
					BucketKeyEnabled: proto.Bool(bucket.EncryptionConfig.BucketKeyEnabled != nil && *bucket.EncryptionConfig.BucketKeyEnabled),
				},
			},
		},
	}, nil
}

// deleteBucketEncryptionCore removes the encryption configuration from a
// bucket.
func (s *S3Service) deleteBucketEncryptionCore(bucketStore s3store.BucketStoreInterface, in *DeleteBucketEncryptionInput) error {
	return bucketStore.SetEncryption(in.Bucket, nil)
}

// putBucketLifecycleConfigurationCore validates the lifecycle rules and
// persists the lifecycle configuration for a bucket.
func (s *S3Service) putBucketLifecycleConfigurationCore(bucketStore s3store.BucketStoreInterface, in *PutBucketLifecycleConfigurationInput) error {
	if in.LifecycleConfiguration == nil {
		return NewInvalidArgumentError("lifecycle configuration is required")
	}

	if err := validateLifecycleRules(in.LifecycleConfiguration.Rules); err != nil {
		return err
	}

	var rules []s3store.LifecycleRule
	for _, rule := range in.LifecycleConfiguration.Rules {
		lifecycleRule := s3store.LifecycleRule{
			ID:     rule.ID,
			Status: rule.Status,
		}

		if rule.Filter != nil {
			lifecycleRule.Filter = &s3store.LifecycleRuleFilter{
				Prefix:                rule.Filter.Prefix,
				ObjectSizeGreaterThan: rule.Filter.ObjectSizeGreaterThan,
				ObjectSizeLessThan:    rule.Filter.ObjectSizeLessThan,
			}
			if rule.Filter.Tag != nil {
				lifecycleRule.Filter.Tag = &types.Tag{Key: rule.Filter.Tag.Key, Value: rule.Filter.Tag.Value}
			}
			if rule.Filter.And != nil {
				lifecycleRule.Filter.And = &s3store.LifecycleRuleAndOperator{
					Prefix:                rule.Filter.And.Prefix,
					ObjectSizeGreaterThan: rule.Filter.And.ObjectSizeGreaterThan,
					ObjectSizeLessThan:    rule.Filter.And.ObjectSizeLessThan,
				}
				for _, t := range rule.Filter.And.Tags {
					lifecycleRule.Filter.And.Tags = append(lifecycleRule.Filter.And.Tags, types.Tag{Key: t.Key, Value: t.Value})
				}
			}
		}

		if rule.Expiration != nil {
			lifecycleRule.Expiration = &s3store.LifecycleExpiration{
				Date:                      rule.Expiration.Date,
				Days:                      rule.Expiration.Days,
				ExpiredObjectDeleteMarker: rule.Expiration.ExpiredObjectDeleteMarker,
			}
		}

		for _, t := range rule.Transitions {
			lifecycleRule.Transitions = append(lifecycleRule.Transitions, s3store.LifecycleTransition{
				Date:         t.Date,
				Days:         t.Days,
				StorageClass: s3store.ObjectStorageClass(t.StorageClass),
			})
		}

		if rule.NoncurrentVersionExpiration != nil {
			lifecycleRule.NoncurrentVersionExpiration = &s3store.NoncurrentVersionExpiration{
				NoncurrentDays:          rule.NoncurrentVersionExpiration.NoncurrentDays,
				NewerNoncurrentVersions: rule.NoncurrentVersionExpiration.NewerNoncurrentVersions,
			}
		}

		for _, t := range rule.NoncurrentVersionTransitions {
			lifecycleRule.NoncurrentVersionTransitions = append(lifecycleRule.NoncurrentVersionTransitions, s3store.NoncurrentVersionTransition{
				NoncurrentDays:          t.NoncurrentDays,
				NewerNoncurrentVersions: t.NewerNoncurrentVersions,
				StorageClass:            s3store.ObjectStorageClass(t.StorageClass),
			})
		}

		if rule.AbortIncompleteMultipartUpload != nil {
			lifecycleRule.AbortIncompleteMultipartUpload = &s3store.AbortIncompleteUpload{
				DaysAfterInitiation: rule.AbortIncompleteMultipartUpload.DaysAfterInitiation,
			}
		}

		rules = append(rules, lifecycleRule)
	}

	return bucketStore.SetLifecycleConfiguration(in.Bucket, &s3store.LifecycleConfiguration{Rules: rules})
}

// getBucketLifecycleConfigurationCore reads the lifecycle configuration
// for a bucket and fails with ErrNoSuchLifecycle when none is set.
func (s *S3Service) getBucketLifecycleConfigurationCore(bucketStore s3store.BucketStoreInterface, in *GetBucketLifecycleConfigurationInput) (*GetBucketLifecycleConfigurationOutput, error) {
	bucket, err := bucketStore.Get(in.Bucket)
	if err != nil {
		return nil, err
	}

	if bucket.LifecycleConfiguration == nil {
		return nil, ErrNoSuchLifecycle
	}

	var rules []LifecycleRuleOutput
	for _, rule := range bucket.LifecycleConfiguration.Rules {
		outputRule := LifecycleRuleOutput{
			ID:     rule.ID,
			Status: rule.Status,
		}

		if rule.Filter != nil {
			outputRule.Filter = &LifecycleRuleFilterOutput{
				Prefix:                rule.Filter.Prefix,
				ObjectSizeGreaterThan: rule.Filter.ObjectSizeGreaterThan,
				ObjectSizeLessThan:    rule.Filter.ObjectSizeLessThan,
			}
			if rule.Filter.Tag != nil {
				outputRule.Filter.Tag = &Tag{Key: rule.Filter.Tag.Key, Value: rule.Filter.Tag.Value}
			}
			if rule.Filter.And != nil {
				outputRule.Filter.And = &LifecycleRuleAndOperatorOutput{
					Prefix:                rule.Filter.And.Prefix,
					ObjectSizeGreaterThan: rule.Filter.And.ObjectSizeGreaterThan,
					ObjectSizeLessThan:    rule.Filter.And.ObjectSizeLessThan,
				}
				for _, t := range rule.Filter.And.Tags {
					outputRule.Filter.And.Tags = append(outputRule.Filter.And.Tags, Tag{Key: t.Key, Value: t.Value})
				}
			}
		}

		if rule.Expiration != nil {
			outputRule.Expiration = &LifecycleExpirationOutput{
				Date:                      rule.Expiration.Date,
				Days:                      rule.Expiration.Days,
				ExpiredObjectDeleteMarker: rule.Expiration.ExpiredObjectDeleteMarker,
			}
		}

		for _, t := range rule.Transitions {
			outputRule.Transitions = append(outputRule.Transitions, LifecycleTransitionOutput{
				Date:         t.Date,
				Days:         t.Days,
				StorageClass: string(t.StorageClass),
			})
		}

		if rule.NoncurrentVersionExpiration != nil {
			outputRule.NoncurrentVersionExpiration = &NoncurrentVersionExpirationOutput{
				NoncurrentDays:          rule.NoncurrentVersionExpiration.NoncurrentDays,
				NewerNoncurrentVersions: rule.NoncurrentVersionExpiration.NewerNoncurrentVersions,
			}
		}

		for _, t := range rule.NoncurrentVersionTransitions {
			outputRule.NoncurrentVersionTransitions = append(outputRule.NoncurrentVersionTransitions, NoncurrentVersionTransitionOutput{
				NoncurrentDays:          t.NoncurrentDays,
				NewerNoncurrentVersions: t.NewerNoncurrentVersions,
				StorageClass:            string(t.StorageClass),
			})
		}

		if rule.AbortIncompleteMultipartUpload != nil {
			outputRule.AbortIncompleteMultipartUpload = &AbortIncompleteUploadOutput{
				DaysAfterInitiation: rule.AbortIncompleteMultipartUpload.DaysAfterInitiation,
			}
		}

		rules = append(rules, outputRule)
	}

	return &GetBucketLifecycleConfigurationOutput{Rules: rules}, nil
}

// deleteBucketLifecycleConfigurationCore removes the lifecycle
// configuration from a bucket.
func (s *S3Service) deleteBucketLifecycleConfigurationCore(bucketStore s3store.BucketStoreInterface, in *DeleteBucketLifecycleConfigurationInput) error {
	return bucketStore.SetLifecycleConfiguration(in.Bucket, nil)
}

// getBucketLocationCore reads the region constraint of a bucket; buckets
// in the default region report an empty location constraint.
func (s *S3Service) getBucketLocationCore(bucketStore s3store.BucketStoreInterface, in *GetBucketLocationInput) (*GetBucketLocationOutput, error) {
	bucket, err := bucketStore.Get(in.Bucket)
	if err != nil {
		return nil, err
	}

	if bucket.Region == defaults.DefaultRegion {
		return &GetBucketLocationOutput{LocationConstraint: ""}, nil
	}

	return &GetBucketLocationOutput{LocationConstraint: bucket.Region}, nil
}

// putBucketLoggingCore validates the logging configuration and persists it
// for a bucket.
func (s *S3Service) putBucketLoggingCore(bucketStore s3store.BucketStoreInterface, in *PutBucketLoggingInput) error {
	var config *s3store.LoggingConfiguration
	if in.LoggingConfiguration != nil {
		if in.LoggingConfiguration.TargetBucket == "" {
			return NewInvalidArgumentError("TargetBucket is required for logging configuration")
		}
		if len(in.LoggingConfiguration.TargetGrants) > maxLogTargetGrants {
			return NewInvalidArgumentError(fmt.Sprintf("too many TargetGrants (maximum %d)", maxLogTargetGrants))
		}
		config = &s3store.LoggingConfiguration{
			TargetBucket: in.LoggingConfiguration.TargetBucket,
			TargetPrefix: in.LoggingConfiguration.TargetPrefix,
		}

		for _, tg := range in.LoggingConfiguration.TargetGrants {
			if err := validateLogPermission(tg.Permission); err != nil {
				return err
			}
			if tg.Grantee != nil {
				if err := validateGranteeType(tg.Grantee.Type); err != nil {
					return err
				}
			}
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

	return bucketStore.SetLoggingConfiguration(in.Bucket, config)
}

// getBucketLoggingCore reads the logging configuration for a bucket.
func (s *S3Service) getBucketLoggingCore(bucketStore s3store.BucketStoreInterface, in *GetBucketLoggingInput) (*GetBucketLoggingOutput, error) {
	config, err := bucketStore.GetLoggingConfiguration(in.Bucket)
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

// putBucketNotificationConfigurationCore validates the notification
// configurations (Id uniqueness, ARN shape, target existence via the event
// bus, event names, filter rules) and persists the notification
// configuration for a bucket.
func (s *S3Service) putBucketNotificationConfigurationCore(ctx *request.RequestContext, bucketStore s3store.BucketStoreInterface, in *PutBucketNotificationInput) error {
	// Track Id uniqueness across all configurations.
	seenIds := make(map[string]bool)

	validateId := func(id string) error {
		if id == "" {
			return nil
		}
		if seenIds[id] {
			return NewInvalidArgumentError(fmt.Sprintf("duplicate notification configuration Id: %s", id))
		}
		seenIds[id] = true
		return nil
	}

	validateArn := func(arn, expectedService string) error {
		if arn == "" {
			return NewInvalidArgumentError("notification ARN is required")
		}
		parts := strings.SplitN(arn, ":", 6)
		if len(parts) < 6 || parts[0] != "arn" {
			return NewInvalidArgumentError(fmt.Sprintf("invalid ARN format: %s", arn))
		}
		if parts[2] != expectedService {
			return NewInvalidArgumentError(fmt.Sprintf("expected ARN service %s, got %s in: %s", expectedService, parts[2], arn))
		}
		return nil
	}

	validateEvents := func(events []string) error {
		return validateS3EventNames(events)
	}

	config := &s3store.NotificationConfiguration{}

	for _, tc := range in.NotificationConfiguration.TopicConfigurations {
		if err := validateId(tc.Id); err != nil {
			return err
		}
		if err := validateArn(tc.TopicArn, "sns"); err != nil {
			return err
		}
		if err := s.validateNotificationTarget(ctx, tc.TopicArn, "sns"); err != nil {
			return err
		}
		if err := validateEvents(tc.Events); err != nil {
			return err
		}
		topicConfig := s3store.TopicNotificationConfiguration{
			Id:       tc.Id,
			TopicArn: tc.TopicArn,
			Events:   tc.Events,
		}
		if tc.Filter != nil && tc.Filter.S3Key != nil {
			topicConfig.Filter = &s3store.NotificationConfigurationFilter{
				Key: &s3store.S3KeyFilter{},
			}
			for _, fr := range tc.Filter.S3Key.FilterRules {
				if err := validateFilterRule(fr.Name, fr.Value); err != nil {
					return err
				}
				topicConfig.Filter.Key.FilterRules = append(topicConfig.Filter.Key.FilterRules, s3store.FilterRule{
					Name:  fr.Name,
					Value: fr.Value,
				})
			}
		}
		config.TopicConfigurations = append(config.TopicConfigurations, topicConfig)
	}

	for _, qc := range in.NotificationConfiguration.QueueConfigurations {
		if err := validateId(qc.Id); err != nil {
			return err
		}
		if err := validateArn(qc.QueueArn, "sqs"); err != nil {
			return err
		}
		if err := s.validateNotificationTarget(ctx, qc.QueueArn, "sqs"); err != nil {
			return err
		}
		if err := validateEvents(qc.Events); err != nil {
			return err
		}
		queueConfig := s3store.QueueNotificationConfiguration{
			Id:       qc.Id,
			QueueArn: qc.QueueArn,
			Events:   qc.Events,
		}
		if qc.Filter != nil && qc.Filter.S3Key != nil {
			queueConfig.Filter = &s3store.NotificationConfigurationFilter{
				Key: &s3store.S3KeyFilter{},
			}
			for _, fr := range qc.Filter.S3Key.FilterRules {
				if err := validateFilterRule(fr.Name, fr.Value); err != nil {
					return err
				}
				queueConfig.Filter.Key.FilterRules = append(queueConfig.Filter.Key.FilterRules, s3store.FilterRule{
					Name:  fr.Name,
					Value: fr.Value,
				})
			}
		}
		config.QueueConfigurations = append(config.QueueConfigurations, queueConfig)
	}

	for _, lc := range in.NotificationConfiguration.LambdaConfigurations {
		if err := validateId(lc.Id); err != nil {
			return err
		}
		if err := validateArn(lc.LambdaFunctionArn, "lambda"); err != nil {
			return err
		}
		if err := s.validateNotificationTarget(ctx, lc.LambdaFunctionArn, "lambda"); err != nil {
			return err
		}
		if err := validateEvents(lc.Events); err != nil {
			return err
		}
		lambdaConfig := s3store.LambdaNotificationConfiguration{
			Id:                lc.Id,
			LambdaFunctionArn: lc.LambdaFunctionArn,
			Events:            lc.Events,
		}
		if lc.Filter != nil && lc.Filter.S3Key != nil {
			lambdaConfig.Filter = &s3store.NotificationConfigurationFilter{
				Key: &s3store.S3KeyFilter{},
			}
			for _, fr := range lc.Filter.S3Key.FilterRules {
				if err := validateFilterRule(fr.Name, fr.Value); err != nil {
					return err
				}
				lambdaConfig.Filter.Key.FilterRules = append(lambdaConfig.Filter.Key.FilterRules, s3store.FilterRule{
					Name:  fr.Name,
					Value: fr.Value,
				})
			}
		}
		config.LambdaConfigurations = append(config.LambdaConfigurations, lambdaConfig)
	}

	return bucketStore.SetNotificationConfiguration(in.Bucket, config)
}

// getBucketNotificationConfigurationCore reads the notification
// configuration for a bucket; a bucket without notifications reports an
// empty configuration.
func (s *S3Service) getBucketNotificationConfigurationCore(bucketStore s3store.BucketStoreInterface, in *GetBucketNotificationInput) (*GetBucketNotificationOutput, error) {
	config, err := bucketStore.GetNotificationConfiguration(in.Bucket)
	if err != nil {
		return nil, err
	}

	output := &GetBucketNotificationOutput{
		NotificationConfiguration: &NotificationConfigurationOutput{},
	}

	if config == nil {
		return output, nil
	}

	for _, tc := range config.TopicConfigurations {
		topicOut := TopicConfigurationOutput{
			Id:       tc.Id,
			TopicArn: tc.TopicArn,
			Events:   tc.Events,
		}
		if tc.Filter != nil && tc.Filter.Key != nil {
			topicOut.Filter = &NotificationFilterOutput{S3Key: &S3KeyFilterOutput{}}
			for _, fr := range tc.Filter.Key.FilterRules {
				topicOut.Filter.S3Key.FilterRules = append(topicOut.Filter.S3Key.FilterRules, FilterRuleOutput{
					Name:  fr.Name,
					Value: fr.Value,
				})
			}
		}
		output.NotificationConfiguration.TopicConfigurations = append(output.NotificationConfiguration.TopicConfigurations, topicOut)
	}

	for _, qc := range config.QueueConfigurations {
		queueOut := QueueConfigurationOutput{
			Id:       qc.Id,
			QueueArn: qc.QueueArn,
			Events:   qc.Events,
		}
		if qc.Filter != nil && qc.Filter.Key != nil {
			queueOut.Filter = &NotificationFilterOutput{S3Key: &S3KeyFilterOutput{}}
			for _, fr := range qc.Filter.Key.FilterRules {
				queueOut.Filter.S3Key.FilterRules = append(queueOut.Filter.S3Key.FilterRules, FilterRuleOutput{
					Name:  fr.Name,
					Value: fr.Value,
				})
			}
		}
		output.NotificationConfiguration.QueueConfigurations = append(output.NotificationConfiguration.QueueConfigurations, queueOut)
	}

	for _, lc := range config.LambdaConfigurations {
		lambdaOut := LambdaConfigurationOutput{
			Id:                lc.Id,
			LambdaFunctionArn: lc.LambdaFunctionArn,
			Events:            lc.Events,
		}
		if lc.Filter != nil && lc.Filter.Key != nil {
			lambdaOut.Filter = &NotificationFilterOutput{S3Key: &S3KeyFilterOutput{}}
			for _, fr := range lc.Filter.Key.FilterRules {
				lambdaOut.Filter.S3Key.FilterRules = append(lambdaOut.Filter.S3Key.FilterRules, FilterRuleOutput{
					Name:  fr.Name,
					Value: fr.Value,
				})
			}
		}
		output.NotificationConfiguration.LambdaConfigurations = append(output.NotificationConfiguration.LambdaConfigurations, lambdaOut)
	}

	return output, nil
}

// putObjectLockConfigurationCore validates the object lock configuration
// (enabled flag, retention mode, Days/Years exclusivity and bounds) and
// persists it for a bucket that has object lock enabled.
func (s *S3Service) putObjectLockConfigurationCore(bucketStore s3store.BucketStoreInterface, in *PutObjectLockConfigurationInput) error {
	if in.ObjectLockConfiguration == nil {
		return NewInvalidArgumentError("ObjectLockConfiguration is required")
	}
	if err := validateObjectLockEnabled(in.ObjectLockConfiguration.ObjectLockEnabled); err != nil {
		return err
	}

	bucket, err := bucketStore.Get(in.Bucket)
	if err != nil {
		return err
	}

	if !bucket.ObjectLockEnabled {
		return ErrObjectLockNotEnabled
	}

	config := &s3store.ObjectLockConfiguration{
		ObjectLockEnabled: in.ObjectLockConfiguration.ObjectLockEnabled,
	}

	if in.ObjectLockConfiguration.Rule != nil && in.ObjectLockConfiguration.Rule.DefaultRetention != nil {
		dr := in.ObjectLockConfiguration.Rule.DefaultRetention

		// Validate Mode: must be GOVERNANCE or COMPLIANCE
		if dr.Mode != string(s3store.ObjectLockRetentionModeGovernance) && dr.Mode != string(s3store.ObjectLockRetentionModeCompliance) {
			return NewInvalidArgumentError(fmt.Sprintf("invalid retention mode: %s (must be GOVERNANCE or COMPLIANCE)", dr.Mode))
		}

		// Days/Years are mutually exclusive: exactly one must be specified
		hasDays := dr.Days != nil
		hasYears := dr.Years != nil
		if hasDays && hasYears {
			return NewInvalidArgumentError("Days and Years are mutually exclusive; specify exactly one")
		}
		if !hasDays && !hasYears {
			return NewInvalidArgumentError("either Days or Years must be specified in DefaultRetention")
		}
		if hasDays {
			if *dr.Days < 1 || *dr.Days > 3650 {
				return NewInvalidArgumentError(fmt.Sprintf("Days must be between 1 and 3650, got %d", *dr.Days))
			}
		}
		if hasYears {
			if *dr.Years < 1 || *dr.Years > 100 {
				return NewInvalidArgumentError(fmt.Sprintf("Years must be between 1 and 100, got %d", *dr.Years))
			}
		}

		config.Rule = &s3store.ObjectLockRule{
			DefaultRetention: &s3store.DefaultRetention{
				Mode:  s3store.ObjectLockRetentionMode(dr.Mode),
				Days:  dr.Days,
				Years: dr.Years,
			},
		}
	}

	return bucketStore.SetObjectLockConfiguration(in.Bucket, config)
}

// getObjectLockConfigurationCore reads the object lock configuration for a
// bucket; a bucket without object lock enabled fails with
// ErrNoSuchObjectLock.
func (s *S3Service) getObjectLockConfigurationCore(bucketStore s3store.BucketStoreInterface, in *GetObjectLockConfigurationInput) (*GetObjectLockConfigurationOutput, error) {
	bucket, err := bucketStore.Get(in.Bucket)
	if err != nil {
		return nil, err
	}

	if !bucket.ObjectLockEnabled {
		return nil, ErrNoSuchObjectLock
	}

	if bucket.ObjectLockConfig == nil {
		return &GetObjectLockConfigurationOutput{
			ObjectLockConfiguration: &ObjectLockConfigurationOutput{
				ObjectLockEnabled: "Enabled",
			},
		}, nil
	}

	output := &GetObjectLockConfigurationOutput{
		ObjectLockConfiguration: &ObjectLockConfigurationOutput{
			ObjectLockEnabled: bucket.ObjectLockConfig.ObjectLockEnabled,
		},
	}

	if bucket.ObjectLockConfig.Rule != nil && bucket.ObjectLockConfig.Rule.DefaultRetention != nil {
		output.ObjectLockConfiguration.Rule = &ObjectLockRuleOutput{
			DefaultRetention: &DefaultRetentionOutput{
				Mode:  string(bucket.ObjectLockConfig.Rule.DefaultRetention.Mode),
				Days:  bucket.ObjectLockConfig.Rule.DefaultRetention.Days,
				Years: bucket.ObjectLockConfig.Rule.DefaultRetention.Years,
			},
		}
	}

	return output, nil
}

// getBucketCore reads the full bucket record.
func (s *S3Service) getBucketCore(bucketStore s3store.BucketStoreInterface, in *GetBucketInput) (*s3store.Bucket, error) {
	return bucketStore.Get(in.Bucket)
}

// headBucketCore reads the bucket's region for the HeadBucket response.
func (s *S3Service) headBucketCore(bucketStore s3store.BucketStoreInterface, in *HeadBucketInput) (*HeadBucketOutput, error) {
	bucket, err := bucketStore.Get(in.Bucket)
	if err != nil {
		return nil, err
	}
	return &HeadBucketOutput{
		BucketRegion: bucket.Region,
	}, nil
}

// putBucketOwnershipControlsCore validates the ownership rules and
// persists the ownership controls for a bucket.
func (s *S3Service) putBucketOwnershipControlsCore(bucketStore s3store.BucketStoreInterface, in *PutBucketOwnershipControlsInput) error {
	if in.OwnershipControls == nil {
		return NewInvalidArgumentError("OwnershipControls is required")
	}
	if err := validateOwnershipControls(in.OwnershipControls.Rules); err != nil {
		return err
	}

	config := &s3store.OwnershipControls{}
	for _, rule := range in.OwnershipControls.Rules {
		config.Rules = append(config.Rules, s3store.OwnershipControlsRule{
			ObjectOwnership: rule.ObjectOwnership,
		})
	}

	return bucketStore.SetOwnershipControls(in.Bucket, config)
}

// getBucketOwnershipControlsCore reads the ownership controls for a bucket
// and fails with ErrNoSuchOwnershipCtrls when none are set.
func (s *S3Service) getBucketOwnershipControlsCore(bucketStore s3store.BucketStoreInterface, in *GetBucketOwnershipControlsInput) (*GetBucketOwnershipControlsOutput, error) {
	config, err := bucketStore.GetOwnershipControls(in.Bucket)
	if err != nil {
		return nil, err
	}

	if config == nil {
		return nil, ErrNoSuchOwnershipCtrls
	}

	result := &GetBucketOwnershipControlsOutput{
		OwnershipControls: &OwnershipControlsInput{},
	}
	for _, rule := range config.Rules {
		result.OwnershipControls.Rules = append(result.OwnershipControls.Rules, OwnershipControlsRuleInput{
			ObjectOwnership: rule.ObjectOwnership,
		})
	}

	return result, nil
}

// deleteBucketOwnershipControlsCore removes the ownership controls from a
// bucket.
func (s *S3Service) deleteBucketOwnershipControlsCore(bucketStore s3store.BucketStoreInterface, in *DeleteBucketOwnershipControlsInput) error {
	return bucketStore.SetOwnershipControls(in.Bucket, nil)
}

// putBucketPolicyCore validates the policy document, enforces
// BlockPublicPolicy against public policies, and persists the bucket
// policy.
func (s *S3Service) putBucketPolicyCore(bucketStore s3store.BucketStoreInterface, in *PutBucketPolicyInput) error {
	if err := validatePolicyDocument(in.Policy); err != nil {
		return err
	}

	publicAccessBlock, _ := bucketStore.GetPublicAccessBlock(in.Bucket)
	if publicAccessBlock != nil && publicAccessBlock.BlockPublicPolicy {
		if policyContainsPublicAccess(in.Policy) {
			return NewInvalidArgumentError("bucket has BlockPublicPolicy enabled")
		}
	}

	return bucketStore.SetPolicy(in.Bucket, in.Policy)
}

// getBucketPolicyCore reads the bucket policy and fails with
// ErrNoSuchBucketPolicy when none is set.
func (s *S3Service) getBucketPolicyCore(bucketStore s3store.BucketStoreInterface, in *GetBucketPolicyInput) (*GetBucketPolicyOutput, error) {
	bucket, err := bucketStore.Get(in.Bucket)
	if err != nil {
		return nil, err
	}

	if bucket.Policy == "" {
		return nil, ErrNoSuchBucketPolicy
	}

	return &GetBucketPolicyOutput{
		Policy: bucket.Policy,
	}, nil
}

// deleteBucketPolicyCore removes the policy from a bucket.
func (s *S3Service) deleteBucketPolicyCore(bucketStore s3store.BucketStoreInterface, in *DeleteBucketPolicyInput) error {
	return bucketStore.SetPolicy(in.Bucket, "")
}

// getBucketPolicyStatusCore reports whether the bucket is public. The
// classification is shared with the BlockPublicPolicy enforcement so the
// reported status and the rejected policies can never diverge. A bucket
// without a policy grants nothing and is therefore never public.
func (s *S3Service) getBucketPolicyStatusCore(bucketStore s3store.BucketStoreInterface, in *GetBucketPolicyStatusInput) (*GetBucketPolicyStatusOutput, error) {
	bucket, err := bucketStore.Get(in.Bucket)
	if err != nil {
		return nil, err
	}

	return &GetBucketPolicyStatusOutput{
		PolicyStatus: &PolicyStatus{
			IsPublic: bucket.Policy != "" && policyContainsPublicAccess(bucket.Policy),
		},
	}, nil
}

// putPublicAccessBlockCore persists the public access block configuration
// for a bucket.
func (s *S3Service) putPublicAccessBlockCore(bucketStore s3store.BucketStoreInterface, in *PutPublicAccessBlockInput) error {
	config := &s3store.PublicAccessBlockConfig{
		BlockPublicAcls:       in.PublicAccessBlockConfiguration.BlockPublicAcls,
		BlockPublicPolicy:     in.PublicAccessBlockConfiguration.BlockPublicPolicy,
		IgnorePublicAcls:      in.PublicAccessBlockConfiguration.IgnorePublicAcls,
		RestrictPublicBuckets: in.PublicAccessBlockConfiguration.RestrictPublicBuckets,
	}

	return bucketStore.SetPublicAccessBlock(in.Bucket, config)
}

// getPublicAccessBlockCore reads the public access block configuration for
// a bucket and fails with ErrNoSuchPublicAccessBlk when none is set.
func (s *S3Service) getPublicAccessBlockCore(bucketStore s3store.BucketStoreInterface, in *GetPublicAccessBlockInput) (*GetPublicAccessBlockOutput, error) {
	config, err := bucketStore.GetPublicAccessBlock(in.Bucket)
	if err != nil {
		return nil, err
	}

	if config == nil {
		return nil, ErrNoSuchPublicAccessBlk
	}

	return &GetPublicAccessBlockOutput{
		PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
			BlockPublicAcls:       config.BlockPublicAcls,
			BlockPublicPolicy:     config.BlockPublicPolicy,
			IgnorePublicAcls:      config.IgnorePublicAcls,
			RestrictPublicBuckets: config.RestrictPublicBuckets,
		},
	}, nil
}

// deletePublicAccessBlockCore removes the public access block
// configuration from a bucket.
func (s *S3Service) deletePublicAccessBlockCore(bucketStore s3store.BucketStoreInterface, in *DeletePublicAccessBlockInput) error {
	return bucketStore.SetPublicAccessBlock(in.Bucket, nil)
}

// putBucketRequestPaymentCore validates the payer and persists the request
// payment configuration for a bucket.
func (s *S3Service) putBucketRequestPaymentCore(bucketStore s3store.BucketStoreInterface, in *PutBucketRequestPaymentInput) error {
	if in.RequestPaymentConfiguration == nil {
		return NewInvalidArgumentError("RequestPaymentConfiguration is required")
	}
	if err := validatePayer(in.RequestPaymentConfiguration.Payer); err != nil {
		return err
	}

	config := &s3store.RequestPaymentConfiguration{
		Payer: in.RequestPaymentConfiguration.Payer,
	}

	return bucketStore.SetRequestPayment(in.Bucket, config)
}

// getBucketRequestPaymentCore reads the request payment configuration for
// a bucket; a bucket without configuration reports BucketOwner.
func (s *S3Service) getBucketRequestPaymentCore(bucketStore s3store.BucketStoreInterface, in *GetBucketRequestPaymentInput) (*GetBucketRequestPaymentOutput, error) {
	config, err := bucketStore.GetRequestPayment(in.Bucket)
	if err != nil {
		return nil, err
	}

	if config == nil {
		return &GetBucketRequestPaymentOutput{
			RequestPaymentConfiguration: &RequestPaymentConfigurationInput{
				Payer: "BucketOwner",
			},
		}, nil
	}

	return &GetBucketRequestPaymentOutput{
		RequestPaymentConfiguration: &RequestPaymentConfigurationInput{
			Payer: config.Payer,
		},
	}, nil
}

// putBucketTaggingCore validates the tag set and persists the bucket tags.
func (s *S3Service) putBucketTaggingCore(bucketStore s3store.BucketStoreInterface, in *PutBucketTaggingInput) error {
	if err := validateTags(in.Tags); err != nil {
		return err
	}
	return bucketStore.SetTags(in.Bucket, TagsToCommon(in.Tags))
}

// getBucketTaggingCore reads the tags attached to a bucket.
func (s *S3Service) getBucketTaggingCore(bucketStore s3store.BucketStoreInterface, in *GetBucketTaggingInput) (*GetBucketTaggingOutput, error) {
	bucket, err := bucketStore.Get(in.Bucket)
	if err != nil {
		return nil, err
	}

	return &GetBucketTaggingOutput{
		TagSet: CommonToTags(bucket.Tags),
	}, nil
}

// deleteBucketTaggingCore removes all tags from a bucket.
func (s *S3Service) deleteBucketTaggingCore(bucketStore s3store.BucketStoreInterface, in *DeleteBucketTaggingInput) error {
	return bucketStore.SetTags(in.Bucket, nil)
}

// putBucketVersioningCore validates the versioning status and MFA delete
// setting, requires the bucket to exist, and persists the versioning state.
func (s *S3Service) putBucketVersioningCore(bucketStore s3store.BucketStoreInterface, in *PutBucketVersioningInput) error {
	status := s3store.BucketVersioningStatus(in.Status)
	if status != s3store.BucketVersioningEnabled && status != s3store.BucketVersioningSuspended {
		return NewInvalidArgumentError("invalid versioning status")
	}
	if err := validateMFADelete(in.MFADelete); err != nil {
		return err
	}
	if !bucketStore.Exists(in.Bucket) {
		return ErrNoSuchBucket
	}
	return bucketStore.SetVersioning(in.Bucket, status, in.MFADelete)
}

// getBucketVersioningCore reads the versioning configuration of a bucket;
// an unversioned bucket reports an empty configuration.
func (s *S3Service) getBucketVersioningCore(bucketStore s3store.BucketStoreInterface, in *GetBucketVersioningInput) (*GetBucketVersioningOutput, error) {
	bucket, err := bucketStore.Get(in.Bucket)
	if err != nil {
		return nil, err
	}

	if bucket.VersioningStatus == "" && bucket.MFADelete == "" {
		return &GetBucketVersioningOutput{}, nil
	}

	result := &GetBucketVersioningOutput{
		Status: string(bucket.VersioningStatus),
	}
	if bucket.MFADelete != "" {
		result.MFADelete = bucket.MFADelete
	}
	return result, nil
}

// putBucketWebsiteCore validates the website configuration and persists it
// for a bucket.
func (s *S3Service) putBucketWebsiteCore(bucketStore s3store.BucketStoreInterface, in *PutBucketWebsiteInput) error {
	if err := validateWebsiteConfig(in.WebsiteConfiguration); err != nil {
		return err
	}

	config := &s3store.WebsiteConfiguration{}

	if in.WebsiteConfiguration.IndexDocument != nil {
		config.IndexDocument = in.WebsiteConfiguration.IndexDocument.Suffix
	}

	if in.WebsiteConfiguration.ErrorDocument != nil {
		config.ErrorDocument = in.WebsiteConfiguration.ErrorDocument.Key
	}

	if in.WebsiteConfiguration.RedirectAllRequestsTo != nil {
		config.RedirectAllRequestsTo = &s3store.RedirectAllRequestsTo{
			HostName: in.WebsiteConfiguration.RedirectAllRequestsTo.HostName,
			Protocol: in.WebsiteConfiguration.RedirectAllRequestsTo.Protocol,
		}
	}

	for _, rule := range in.WebsiteConfiguration.RoutingRules {
		routingRule := s3store.RoutingRule{}

		if rule.Condition != nil {
			routingRule.Condition = &s3store.RoutingRuleCondition{
				HTTPErrorCodeReturnedEquals: rule.Condition.HTTPErrorCodeReturnedEquals,
				KeyPrefixEquals:             rule.Condition.KeyPrefixEquals,
			}
		}

		if rule.Redirect != nil {
			routingRule.Redirect = &s3store.RoutingRuleRedirect{
				HostName:             rule.Redirect.HostName,
				HTTPRedirectCode:     rule.Redirect.HTTPRedirectCode,
				Protocol:             rule.Redirect.Protocol,
				ReplaceKeyPrefixWith: rule.Redirect.ReplaceKeyPrefixWith,
				ReplaceKeyWith:       rule.Redirect.ReplaceKeyWith,
			}
		}

		config.RoutingRules = append(config.RoutingRules, routingRule)
	}

	return bucketStore.SetWebsiteConfiguration(in.Bucket, config)
}

// getBucketWebsiteCore reads the website configuration for a bucket and
// fails with ErrNoSuchWebsite when none is set.
func (s *S3Service) getBucketWebsiteCore(bucketStore s3store.BucketStoreInterface, in *GetBucketWebsiteInput) (*GetBucketWebsiteOutput, error) {
	bucket, err := bucketStore.Get(in.Bucket)
	if err != nil {
		return nil, err
	}

	if bucket.WebsiteConfiguration == nil {
		return nil, ErrNoSuchWebsite
	}

	output := &GetBucketWebsiteOutput{}

	if bucket.WebsiteConfiguration.IndexDocument != "" {
		output.IndexDocument = &IndexDocumentOutput{
			Suffix: bucket.WebsiteConfiguration.IndexDocument,
		}
	}

	if bucket.WebsiteConfiguration.ErrorDocument != "" {
		output.ErrorDocument = &ErrorDocumentOutput{
			Key: bucket.WebsiteConfiguration.ErrorDocument,
		}
	}

	if bucket.WebsiteConfiguration.RedirectAllRequestsTo != nil {
		output.RedirectAllRequestsTo = &RedirectAllRequestsToOutput{
			HostName: bucket.WebsiteConfiguration.RedirectAllRequestsTo.HostName,
			Protocol: bucket.WebsiteConfiguration.RedirectAllRequestsTo.Protocol,
		}
	}

	for _, rule := range bucket.WebsiteConfiguration.RoutingRules {
		outputRule := RoutingRuleOutput{}

		if rule.Condition != nil {
			outputRule.Condition = &RoutingRuleConditionOutput{
				HTTPErrorCodeReturnedEquals: rule.Condition.HTTPErrorCodeReturnedEquals,
				KeyPrefixEquals:             rule.Condition.KeyPrefixEquals,
			}
		}

		if rule.Redirect != nil {
			outputRule.Redirect = &RedirectOutput{
				HostName:             rule.Redirect.HostName,
				HTTPRedirectCode:     rule.Redirect.HTTPRedirectCode,
				Protocol:             rule.Redirect.Protocol,
				ReplaceKeyPrefixWith: rule.Redirect.ReplaceKeyPrefixWith,
				ReplaceKeyWith:       rule.Redirect.ReplaceKeyWith,
			}
		}

		output.RoutingRules = append(output.RoutingRules, outputRule)
	}

	return output, nil
}

// deleteBucketWebsiteCore removes the website configuration from a bucket.
func (s *S3Service) deleteBucketWebsiteCore(bucketStore s3store.BucketStoreInterface, in *DeleteBucketWebsiteInput) error {
	return bucketStore.SetWebsiteConfiguration(in.Bucket, nil)
}
