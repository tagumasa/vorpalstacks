package testutil

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/aws/aws-sdk-go-v2/service/iot/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"vorpalstacks-sdk-tests/config"
)

// runIoTThingTests covers the Thing resource lifecycle (Create/Describe/Update/
// List/Delete) plus negative paths (Duplicate, NotFound, RemoveThingType) and
// the StartThingRegistrationTask bulk path. All resources use uniqueName so the
// suite is safe under parallel runs and re-runs.
func (r *TestRunner) runIoTThingTests(tc *iotTestContext) []TestResult {
	var results []TestResult
	thingName := uniqueName("thing")
	attrs := map[string]string{"version": "1.0", "location": tc.region}

	// Best-effort cleanup so a failed run never leaves the thing behind.
	defer tc.client.DeleteThing(tc.ctx, &iot.DeleteThingInput{ThingName: aws.String(thingName)})

	results = append(results, r.RunTest("iot", "Thing_CreateThing", func() error {
		out, err := tc.client.CreateThing(tc.ctx, &iot.CreateThingInput{
			ThingName:        aws.String(thingName),
			AttributePayload: &types.AttributePayload{Attributes: attrs},
		})
		if err != nil {
			return fmt.Errorf("CreateThing failed: %w", err)
		}
		if out.ThingName == nil || *out.ThingName != thingName {
			return fmt.Errorf("expected thingName=%s, got %v", thingName, out.ThingName)
		}
		if out.ThingArn == nil || *out.ThingArn == "" {
			return fmt.Errorf("expected non-empty thingArn")
		}
		if out.ThingId == nil || *out.ThingId == "" {
			return fmt.Errorf("expected non-empty thingId")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Thing_DescribeThing", func() error {
		out, err := tc.client.DescribeThing(tc.ctx, &iot.DescribeThingInput{ThingName: aws.String(thingName)})
		if err != nil {
			return fmt.Errorf("DescribeThing failed: %w", err)
		}
		if out.ThingName == nil || *out.ThingName != thingName {
			return fmt.Errorf("expected thingName=%s, got %v", thingName, out.ThingName)
		}
		if out.Attributes == nil {
			return fmt.Errorf("expected attributes to be non-nil")
		}
		if v, ok := out.Attributes["version"]; !ok || v != "1.0" {
			return fmt.Errorf("expected attribute version=1.0, got %v", out.Attributes["version"])
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Thing_UpdateThing", func() error {
		_, err := tc.client.UpdateThing(tc.ctx, &iot.UpdateThingInput{
			ThingName: aws.String(thingName),
			AttributePayload: &types.AttributePayload{
				Attributes: map[string]string{"version": "2.0"},
			},
		})
		if err != nil {
			return fmt.Errorf("UpdateThing failed: %w", err)
		}
		out, err := tc.client.DescribeThing(tc.ctx, &iot.DescribeThingInput{ThingName: aws.String(thingName)})
		if err != nil {
			return fmt.Errorf("DescribeThing after update failed: %w", err)
		}
		if v, ok := out.Attributes["version"]; !ok || v != "2.0" {
			return fmt.Errorf("expected attribute version=2.0 after update, got %v", out.Attributes["version"])
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Thing_ListThings_IncludesCreated", func() error {
		found, err := tc.thingExists(thingName)
		if err != nil {
			return fmt.Errorf("ListThings failed: %w", err)
		}
		if !found {
			return fmt.Errorf("%s not found in thing list", thingName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Thing_CreateThing_Duplicate_Conflict", func() error {
		_, err := tc.client.CreateThing(tc.ctx, &iot.CreateThingInput{ThingName: aws.String(thingName)})
		return expectConflict(err)
	}))

	results = append(results, r.RunTest("iot", "Thing_DescribeThing_NotFound", func() error {
		_, err := tc.client.DescribeThing(tc.ctx, &iot.DescribeThingInput{ThingName: aws.String(uniqueName("nope"))})
		return expectNotFound(err)
	}))

	// RemoveThingType: create a thing type, attach it to a fresh thing, then
	// remove it and verify DescribeThing no longer reports the type.
	ttName := uniqueName("thing-type")
	thingWithType := uniqueName("thing-tt")
	defer tc.client.DeleteThing(tc.ctx, &iot.DeleteThingInput{ThingName: aws.String(thingWithType)})
	defer tc.client.DeleteThingType(tc.ctx, &iot.DeleteThingTypeInput{ThingTypeName: aws.String(ttName)})

	results = append(results, r.RunTest("iot", "Thing_UpdateThing_RemoveThingType", func() error {
		if _, err := tc.client.CreateThingType(tc.ctx, &iot.CreateThingTypeInput{ThingTypeName: aws.String(ttName)}); err != nil {
			return fmt.Errorf("CreateThingType prerequisite failed: %w", err)
		}
		if _, err := tc.client.CreateThing(tc.ctx, &iot.CreateThingInput{
			ThingName:     aws.String(thingWithType),
			ThingTypeName: aws.String(ttName),
		}); err != nil {
			return fmt.Errorf("CreateThing with type failed: %w", err)
		}
		before, err := tc.client.DescribeThing(tc.ctx, &iot.DescribeThingInput{ThingName: aws.String(thingWithType)})
		if err != nil {
			return fmt.Errorf("DescribeThing before remove failed: %w", err)
		}
		if before.ThingTypeName == nil || *before.ThingTypeName != ttName {
			return fmt.Errorf("expected thingTypeName=%s before removal, got %v", ttName, before.ThingTypeName)
		}
		if _, err := tc.client.UpdateThing(tc.ctx, &iot.UpdateThingInput{
			ThingName:       aws.String(thingWithType),
			RemoveThingType: true,
		}); err != nil {
			return fmt.Errorf("UpdateThing RemoveThingType failed: %w", err)
		}
		after, err := tc.client.DescribeThing(tc.ctx, &iot.DescribeThingInput{ThingName: aws.String(thingWithType)})
		if err != nil {
			return fmt.Errorf("DescribeThing after remove failed: %w", err)
		}
		if after.ThingTypeName != nil && *after.ThingTypeName != "" {
			return fmt.Errorf("expected empty thingTypeName after removal, got %v", after.ThingTypeName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Thing_UpdateThing_NonexistentThingType_NotFound", func() error {
		_, err := tc.client.UpdateThing(tc.ctx, &iot.UpdateThingInput{
			ThingName:     aws.String(thingName),
			ThingTypeName: aws.String(uniqueName("nope-type")),
		})
		return expectNotFound(err)
	}))

	// Set a thing type on a previously typeless thing (RemoveThingType=false +
	// ThingTypeName path).
	thingForSet := uniqueName("thing-set")
	ttForSet := uniqueName("thing-type-set")
	defer tc.client.DeleteThing(tc.ctx, &iot.DeleteThingInput{ThingName: aws.String(thingForSet)})
	defer tc.client.DeleteThingType(tc.ctx, &iot.DeleteThingTypeInput{ThingTypeName: aws.String(ttForSet)})

	results = append(results, r.RunTest("iot", "Thing_UpdateThing_SetThingType", func() error {
		if _, err := tc.client.CreateThingType(tc.ctx, &iot.CreateThingTypeInput{ThingTypeName: aws.String(ttForSet)}); err != nil {
			return fmt.Errorf("CreateThingType prerequisite failed: %w", err)
		}
		if _, err := tc.client.CreateThing(tc.ctx, &iot.CreateThingInput{ThingName: aws.String(thingForSet)}); err != nil {
			return fmt.Errorf("CreateThing prerequisite failed: %w", err)
		}
		if _, err := tc.client.UpdateThing(tc.ctx, &iot.UpdateThingInput{
			ThingName:       aws.String(thingForSet),
			RemoveThingType: false,
			ThingTypeName:   aws.String(ttForSet),
		}); err != nil {
			return fmt.Errorf("UpdateThing set ThingTypeName failed: %w", err)
		}
		desc, err := tc.client.DescribeThing(tc.ctx, &iot.DescribeThingInput{ThingName: aws.String(thingForSet)})
		if err != nil {
			return fmt.Errorf("DescribeThing failed: %w", err)
		}
		if desc.ThingTypeName == nil || *desc.ThingTypeName != ttForSet {
			return fmt.Errorf("expected thingTypeName=%s, got %v", ttForSet, desc.ThingTypeName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Thing_StartThingRegistrationTask_Validation", func() error {
		// A minimal body must be rejected with a validation error, proving
		// the handler is registered and validates.
		_, err := tc.client.StartThingRegistrationTask(tc.ctx, &iot.StartThingRegistrationTaskInput{
			TemplateBody: aws.String("{}"),
		})
		if err == nil {
			return fmt.Errorf("expected the empty template to be rejected")
		}
		// The documented member bounds: templateBody shares the 10240-byte
		// TemplateBody bound; the input-file bucket carries its own pattern.
		bigBody := `{"Resources":{"thing":{"Type":"AWS::IoT::Thing"}},"padding":"` +
			strings.Repeat("x", 10240) + `"}`
		_, err = tc.client.StartThingRegistrationTask(tc.ctx, &iot.StartThingRegistrationTaskInput{
			TemplateBody:    aws.String(bigBody),
			InputFileBucket: aws.String("valid-bucket"),
			InputFileKey:    aws.String("devices.ndjson"),
			RoleArn:         aws.String(tc.iamRoleARN("bulk-provisioning")),
		})
		if err == nil {
			return fmt.Errorf("expected the over-length templateBody to be rejected")
		}
		if vErr := expectAWSErrorCode(err, "InvalidRequestException"); vErr != nil {
			return vErr
		}
		_, err = tc.client.StartThingRegistrationTask(tc.ctx, &iot.StartThingRegistrationTaskInput{
			TemplateBody:    aws.String(`{"Resources":{"thing":{"Type":"AWS::IoT::Thing"}}}`),
			InputFileBucket: aws.String("bad:bucket!"),
			InputFileKey:    aws.String("devices.ndjson"),
			RoleArn:         aws.String(tc.iamRoleARN("bulk-provisioning")),
		})
		return expectAWSErrorCode(err, "InvalidRequestException")
	}))

	// The bulk task engine: the S3 input file's newline-delimited JSON lines
	// provision one device each; Describe reports the real counts and echo
	// members, and the RESULTS/ERRORS reports are downloadable through the
	// presigned links from ListThingRegistrationTaskReports.
	results = append(results, r.RunTest("iot", "ThingRegistrationTask_BulkProvision", func() error {
		cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
			Endpoint: r.endpoint,
			Region:   r.region,
		})
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) { o.UsePathStyle = true })
		bucket := uniqueName("iot-bulkreg")
		if _, err := s3Client.CreateBucket(tc.ctx, &s3.CreateBucketInput{Bucket: aws.String(bucket)}); err != nil {
			return fmt.Errorf("create bucket: %w", err)
		}
		defer cleanupS3Bucket(tc.ctx, s3Client, bucket)

		names := []string{uniqueName("bulk-a"), uniqueName("bulk-b"), uniqueName("bulk-c")}
		for _, n := range names {
			defer tc.client.DeleteThing(tc.ctx, &iot.DeleteThingInput{ThingName: aws.String(n)})
		}
		var lines []string
		for _, n := range names {
			lines = append(lines, fmt.Sprintf(`{"ThingName":%q}`, n))
		}
		inputKey := "devices.ndjson"
		if _, err := s3Client.PutObject(tc.ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucket), Key: aws.String(inputKey),
			Body: strings.NewReader(strings.Join(lines, "\n") + "\n"),
		}); err != nil {
			return fmt.Errorf("put input file: %w", err)
		}
		templateBody := `{"Parameters":{"ThingName":{"Type":"String"}},"Resources":{"thing":{"Type":"AWS::IoT::Thing","Properties":{"ThingName":{"Ref":"{{ThingName}}"}}}}}`
		start, err := tc.client.StartThingRegistrationTask(tc.ctx, &iot.StartThingRegistrationTaskInput{
			TemplateBody:    aws.String(templateBody),
			InputFileBucket: aws.String(bucket),
			InputFileKey:    aws.String(inputKey),
			RoleArn:         aws.String(tc.iamRoleARN("bulk-provisioning")),
		})
		if err != nil {
			return fmt.Errorf("StartThingRegistrationTask failed: %w", err)
		}
		taskID := aws.ToString(start.TaskId)

		// The task completes asynchronously; poll Describe to a terminal state.
		var desc *iot.DescribeThingRegistrationTaskOutput
		deadline := time.Now().Add(30 * time.Second)
		for {
			desc, err = tc.client.DescribeThingRegistrationTask(tc.ctx, &iot.DescribeThingRegistrationTaskInput{TaskId: aws.String(taskID)})
			if err != nil {
				return fmt.Errorf("DescribeThingRegistrationTask failed: %w", err)
			}
			if desc.Status == types.StatusCompleted ||
				desc.Status == types.StatusFailed ||
				desc.Status == types.StatusCancelled {
				break
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("task did not reach a terminal state within 30s (status %s)", desc.Status)
			}
			time.Sleep(200 * time.Millisecond)
		}
		if desc.Status != types.StatusCompleted {
			return fmt.Errorf("expected Completed, got %s (message %q)", desc.Status, aws.ToString(desc.Message))
		}
		if desc.SuccessCount != int32(len(names)) {
			return fmt.Errorf("expected successCount=%d, got %d", len(names), desc.SuccessCount)
		}
		if desc.FailureCount != 0 {
			return fmt.Errorf("expected failureCount=0, got %d", desc.FailureCount)
		}
		if desc.PercentageProgress != 100 {
			return fmt.Errorf("expected percentageProgress=100, got %d", desc.PercentageProgress)
		}
		if aws.ToString(desc.TemplateBody) != templateBody || aws.ToString(desc.InputFileBucket) != bucket ||
			aws.ToString(desc.InputFileKey) != inputKey || aws.ToString(desc.RoleArn) != tc.iamRoleARN("bulk-provisioning") {
			return fmt.Errorf("describe echo mismatch: %+v", desc)
		}
		if desc.CreationDate == nil || desc.LastModifiedDate == nil {
			return fmt.Errorf("expected creationDate and lastModifiedDate echo")
		}
		for _, n := range names {
			if _, err := tc.client.DescribeThing(tc.ctx, &iot.DescribeThingInput{ThingName: aws.String(n)}); err != nil {
				return fmt.Errorf("DescribeThing on provisioned %s failed: %w", n, err)
			}
		}

		reports, err := tc.client.ListThingRegistrationTaskReports(tc.ctx, &iot.ListThingRegistrationTaskReportsInput{
			TaskId: aws.String(taskID), ReportType: types.ReportTypeResults,
		})
		if err != nil {
			return fmt.Errorf("ListThingRegistrationTaskReports failed: %w", err)
		}
		if len(reports.ResourceLinks) != 1 {
			return fmt.Errorf("expected one RESULTS report link, got %v", reports.ResourceLinks)
		}
		// The presigned link must carry the endpoint the client dialled,
		// not a server-side localhost fallback.
		linkURL, err := url.Parse(reports.ResourceLinks[0])
		if err != nil {
			return fmt.Errorf("parse report link: %w", err)
		}
		endpointURL, err := url.Parse(r.endpoint)
		if err != nil {
			return fmt.Errorf("parse runner endpoint: %w", err)
		}
		if linkURL.Scheme != endpointURL.Scheme || linkURL.Host != endpointURL.Host {
			return fmt.Errorf("report link must point at the dialled endpoint %s, got %s", r.endpoint, reports.ResourceLinks[0])
		}
		resp, err := http.Get(reports.ResourceLinks[0])
		if err != nil {
			return fmt.Errorf("fetch report link: %w", err)
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("report link returned %d: %s", resp.StatusCode, string(body))
		}
		if !strings.Contains(string(body), names[0]) || !strings.Contains(string(body), "Success") {
			return fmt.Errorf("results report missing provisioned rows: %s", string(body))
		}

		errReports, err := tc.client.ListThingRegistrationTaskReports(tc.ctx, &iot.ListThingRegistrationTaskReportsInput{
			TaskId: aws.String(taskID), ReportType: types.ReportTypeErrors,
		})
		if err != nil {
			return fmt.Errorf("ListThingRegistrationTaskReports(ERRORS) failed: %w", err)
		}
		if len(errReports.ResourceLinks) != 1 {
			return fmt.Errorf("expected one ERRORS report link, got %v", errReports.ResourceLinks)
		}
		return nil
	}))

	// A missing input file fails the task with a message, and stopping a
	// terminal task is rejected as invalid.
	results = append(results, r.RunTest("iot", "ThingRegistrationTask_MissingInputFileFails", func() error {
		templateBody := `{"Resources":{"thing":{"Type":"AWS::IoT::Thing"}}}`
		start, err := tc.client.StartThingRegistrationTask(tc.ctx, &iot.StartThingRegistrationTaskInput{
			TemplateBody:    aws.String(templateBody),
			InputFileBucket: aws.String(uniqueName("iot-bulkreg")),
			InputFileKey:    aws.String("does-not-exist.ndjson"),
			RoleArn:         aws.String(tc.iamRoleARN("bulk-provisioning")),
		})
		if err != nil {
			return fmt.Errorf("StartThingRegistrationTask failed: %w", err)
		}
		taskID := aws.ToString(start.TaskId)
		deadline := time.Now().Add(30 * time.Second)
		for {
			desc, dErr := tc.client.DescribeThingRegistrationTask(tc.ctx, &iot.DescribeThingRegistrationTaskInput{TaskId: aws.String(taskID)})
			if dErr != nil {
				return fmt.Errorf("DescribeThingRegistrationTask failed: %w", dErr)
			}
			if desc.Status == types.StatusFailed {
				if aws.ToString(desc.Message) == "" {
					return fmt.Errorf("expected a failure message on the failed task")
				}
				break
			}
			if desc.Status == types.StatusCompleted {
				return fmt.Errorf("expected Failed for a missing input file, got Completed")
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("task did not fail within 30s (status %s)", desc.Status)
			}
			time.Sleep(200 * time.Millisecond)
		}
		_, err = tc.client.StopThingRegistrationTask(tc.ctx, &iot.StopThingRegistrationTaskInput{TaskId: aws.String(taskID)})
		return expectAWSErrorCode(err, "InvalidRequestException")
	}))

	results = append(results, r.RunTest("iot", "Thing_AttachDetachThingPrincipal", func() error {
		// Create a certificate to use as the principal.
		cert, cleanup, err := tc.createCertificate(true)
		if err != nil {
			return fmt.Errorf("CreateKeysAndCertificate failed: %w", err)
		}
		certARN := cert.ARN
		defer tc.client.DetachThingPrincipal(tc.ctx, &iot.DetachThingPrincipalInput{
			ThingName: aws.String(thingName), Principal: aws.String(certARN),
		})
		defer cleanup()

		if _, err := tc.client.AttachThingPrincipal(tc.ctx, &iot.AttachThingPrincipalInput{
			ThingName: aws.String(thingName), Principal: aws.String(certARN),
		}); err != nil {
			return fmt.Errorf("AttachThingPrincipal failed: %w", err)
		}
		// The principal must appear in ListThingPrincipals.
		lp, err := tc.client.ListThingPrincipals(tc.ctx, &iot.ListThingPrincipalsInput{ThingName: aws.String(thingName)})
		if err != nil {
			return fmt.Errorf("ListThingPrincipals failed: %w", err)
		}
		found := false
		for _, p := range lp.Principals {
			if p == certARN {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("attached principal %s not found in ListThingPrincipals", certARN)
		}
		if _, err := tc.client.DetachThingPrincipal(tc.ctx, &iot.DetachThingPrincipalInput{
			ThingName: aws.String(thingName), Principal: aws.String(certARN),
		}); err != nil {
			return fmt.Errorf("DetachThingPrincipal failed: %w", err)
		}
		lp2, err := tc.client.ListThingPrincipals(tc.ctx, &iot.ListThingPrincipalsInput{ThingName: aws.String(thingName)})
		if err != nil {
			return fmt.Errorf("ListThingPrincipals after detach failed: %w", err)
		}
		for _, p := range lp2.Principals {
			if p == certARN {
				return fmt.Errorf("principal still attached after DetachThingPrincipal")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Thing_DeleteThing", func() error {
		_, err := tc.client.DeleteThing(tc.ctx, &iot.DeleteThingInput{ThingName: aws.String(thingName)})
		if err != nil {
			return fmt.Errorf("DeleteThing failed: %w", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Thing_DeleteThing_NotFound", func() error {
		_, err := tc.client.DeleteThing(tc.ctx, &iot.DeleteThingInput{ThingName: aws.String(uniqueName("nope"))})
		return expectNotFound(err)
	}))

	return results
}
