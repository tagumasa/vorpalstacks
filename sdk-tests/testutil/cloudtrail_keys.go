package testutil

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
)

func (r *TestRunner) runCloudTrailKeysTests(tc *cloudTrailTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("cloudtrail", "ListPublicKeys", func() error {
		name := fmt.Sprintf("pk-trail-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

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
		if pk.Fingerprint == nil || *pk.Fingerprint == "" {
			return fmt.Errorf("expected non-empty Fingerprint")
		}
		if len(pk.Value) == 0 {
			return fmt.Errorf("expected non-empty Value (DER bytes)")
		}
		_, err = base64.StdEncoding.DecodeString(base64.StdEncoding.EncodeToString(pk.Value))
		if err != nil {
			return fmt.Errorf("value should be base64-decodable DER bytes")
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

	results = append(results, r.RunTest("cloudtrail", "ListPublicKeys_TimeFilter", func() error {
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
