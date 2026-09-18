package testutil

import (
	"crypto/x509"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
)

func (r *TestRunner) runCloudTrailKeysTests(tc *cloudTrailTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("cloudtrail", "ListPublicKeys", func() error {
		name := tc.uniqueName("pk-trail")
		defer tc.deleteTrail(name)

		if err := tc.ensureTrailBucket("pk-bucket"); err != nil {
			return err
		}

		_, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:                    aws.String(name),
			S3BucketName:            aws.String("pk-bucket"),
			EnableLogFileValidation: aws.Bool(true),
		})
		if err != nil {
			return err
		}

		resp, err := tc.client.ListPublicKeys(tc.ctx, &cloudtrail.ListPublicKeysInput{})
		if err != nil {
			return fmt.Errorf("list public keys: %v", err)
		}
		if resp.PublicKeyList == nil {
			return fmt.Errorf("PublicKeyList is nil")
		}
		if len(resp.PublicKeyList) == 0 {
			return fmt.Errorf("expected at least 1 public key (from LogFileValidation trail)")
		}
		pk := resp.PublicKeyList[0]
		if pk.Fingerprint == nil {
			return fmt.Errorf("expected non-nil Fingerprint")
		}
		if len(*pk.Fingerprint) != 32 {
			return fmt.Errorf("Fingerprint = %q, want 32 hex characters", *pk.Fingerprint)
		}
		for _, c := range *pk.Fingerprint {
			if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f') {
				return fmt.Errorf("Fingerprint = %q, want lowercase hexadecimal", *pk.Fingerprint)
			}
		}
		if len(pk.Value) == 0 {
			return fmt.Errorf("expected non-empty Value (DER bytes)")
		}
		if _, err := x509.ParsePKCS1PublicKey(pk.Value); err != nil {
			return fmt.Errorf("Value is not parseable as a PKCS#1 DER public key: %v", err)
		}
		if pk.ValidityStartTime == nil || pk.ValidityStartTime.IsZero() {
			return fmt.Errorf("expected ValidityStartTime to be set")
		}
		if pk.ValidityEndTime == nil || pk.ValidityEndTime.IsZero() {
			return fmt.Errorf("expected ValidityEndTime to be set")
		}
		if pk.ValidityEndTime.Before(*pk.ValidityStartTime) {
			return fmt.Errorf("ValidityEndTime should be after ValidityStartTime")
		}
		return nil
	}))

	// Enabling log file validation through UpdateTrail provisions the key
	// material: the default listing (keys valid now) must gain a key whose
	// validity started after the update.
	results = append(results, r.RunTest("cloudtrail", "ListPublicKeys_UpdateEnablesValidation", func() error {
		name := tc.uniqueName("pk-update")
		defer tc.deleteTrail(name)

		if err := tc.ensureTrailBucket("pk-update-bucket"); err != nil {
			return err
		}

		if _, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("pk-update-bucket"),
		}); err != nil {
			return err
		}

		marker := time.Now().UTC().Add(-1 * time.Minute)
		if _, err := tc.client.UpdateTrail(tc.ctx, &cloudtrail.UpdateTrailInput{
			Name:                    aws.String(name),
			EnableLogFileValidation: aws.Bool(true),
		}); err != nil {
			return fmt.Errorf("update enabling validation: %v", err)
		}

		resp, err := tc.client.ListPublicKeys(tc.ctx, &cloudtrail.ListPublicKeysInput{})
		if err != nil {
			return fmt.Errorf("list public keys: %v", err)
		}
		for _, pk := range resp.PublicKeyList {
			if pk.ValidityStartTime != nil && pk.ValidityStartTime.After(marker) {
				return nil
			}
		}
		return fmt.Errorf("no public key became valid after the update-enable (%d keys listed)", len(resp.PublicKeyList))
	}))

	results = append(results, r.RunTest("cloudtrail", "ListPublicKeys_InvalidTimeRange", func() error {
		past := time.Now().UTC().Add(-1 * time.Hour)
		future := time.Now().UTC().Add(1 * time.Hour)
		_, err := tc.client.ListPublicKeys(tc.ctx, &cloudtrail.ListPublicKeysInput{
			StartTime: &future,
			EndTime:   &past,
		})
		return AssertErrorContains(err, "InvalidTimeRangeException")
	}))

	results = append(results, r.RunTest("cloudtrail", "ListPublicKeys_TimeFilter", func() error {
		trailName := tc.uniqueName("pk-time")
		if err := tc.ensureTrailBucket("pk-time-bucket"); err != nil {
			return err
		}
		_, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:                    aws.String(trailName),
			S3BucketName:            aws.String("pk-time-bucket"),
			EnableLogFileValidation: aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("create trail for time filter test: %v", err)
		}
		defer tc.deleteTrail(trailName)

		now := time.Now().UTC()
		past := now.Add(-1 * time.Hour)
		future := now.Add(1 * time.Hour)

		resp, err := tc.client.ListPublicKeys(tc.ctx, &cloudtrail.ListPublicKeysInput{
			StartTime: &past,
			EndTime:   &future,
		})
		if err != nil {
			return fmt.Errorf("list public keys with time filter: %v", err)
		}
		if resp.PublicKeyList == nil {
			return fmt.Errorf("PublicKeyList is nil")
		}
		if len(resp.PublicKeyList) == 0 {
			return fmt.Errorf("expected at least 1 public key in time range")
		}
		pk := resp.PublicKeyList[0]
		if pk.Fingerprint == nil || *pk.Fingerprint == "" {
			return fmt.Errorf("expected non-empty Fingerprint")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "ListPublicKeys_OutsideTimeRange", func() error {
		farFuture := time.Now().UTC().Add(10 * 24 * 365 * time.Hour)
		beyond := farFuture.Add(1 * time.Hour)

		resp, err := tc.client.ListPublicKeys(tc.ctx, &cloudtrail.ListPublicKeysInput{
			StartTime: &farFuture,
			EndTime:   &beyond,
		})
		if err != nil {
			return fmt.Errorf("list public keys outside range: %v", err)
		}
		if resp.PublicKeyList == nil {
			return fmt.Errorf("PublicKeyList is nil")
		}
		if len(resp.PublicKeyList) != 0 {
			return fmt.Errorf("expected 0 public keys outside validity range, got %d", len(resp.PublicKeyList))
		}
		return nil
	}))

	return results
}
