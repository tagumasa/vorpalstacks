package testutil

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

// The vended-logs delivery family pins (E1): the configuration CRUD on
// both destination types the platform carries, the destination policy,
// the delivery lifecycle with its conflict guards, the configuration
// template catalog, and the delivery engine's CWL flow. The S3
// destination flow is the integration suite's cross-service pin.

// describeAllDeliverySources walks every page of the source listing:
// during full regression other services create resources in parallel, so
// a membership assertion cannot rely on the first page alone.
func (tc *cwlogsTestCtx) describeAllDeliverySources() ([]types.DeliverySource, error) {
	var all []types.DeliverySource
	var nextToken *string
	for {
		resp, err := tc.client.DescribeDeliverySources(tc.ctx, &cloudwatchlogs.DescribeDeliverySourcesInput{NextToken: nextToken})
		if err != nil {
			return nil, err
		}
		all = append(all, resp.DeliverySources...)
		if resp.NextToken == nil || *resp.NextToken == "" {
			return all, nil
		}
		nextToken = resp.NextToken
	}
}

// describeAllDeliveries walks every page of the delivery listing.
func (tc *cwlogsTestCtx) describeAllDeliveries() ([]types.Delivery, error) {
	var all []types.Delivery
	var nextToken *string
	for {
		resp, err := tc.client.DescribeDeliveries(tc.ctx, &cloudwatchlogs.DescribeDeliveriesInput{NextToken: nextToken})
		if err != nil {
			return nil, err
		}
		all = append(all, resp.Deliveries...)
		if resp.NextToken == nil || *resp.NextToken == "" {
			return all, nil
		}
		nextToken = resp.NextToken
	}
}

func (tc *cwlogsTestCtx) deliveryTests() []TestResult {
	var results []TestResult

	results = append(results, tc.runner.RunTest("logs", "PutDeliverySource_StatusLifecycle", func() error {
		srcName := tc.uniquePrefix("vended-src")
		group := srcName + "-group"
		if err := tc.createLogGroup(group); err != nil {
			return fmt.Errorf("create group: %v", err)
		}
		defer tc.deleteLogGroup(group)
		// A failure past the put must not leak the source into later
		// runs' listings.
		defer func() {
			_, _ = tc.client.DeleteDeliverySource(tc.ctx, &cloudwatchlogs.DeleteDeliverySourceInput{Name: aws.String(srcName)})
		}()

		groupArn, err := tc.findLogGroupARN(group)
		if err != nil {
			return err
		}
		put, err := tc.client.PutDeliverySource(tc.ctx, &cloudwatchlogs.PutDeliverySourceInput{
			Name:        aws.String(srcName),
			ResourceArn: groupArn,
			LogType:     aws.String("APPLICATION_LOGS"),
		})
		if err != nil {
			return fmt.Errorf("put delivery source: %v", err)
		}
		src := put.DeliverySource
		if src.Name == nil || *src.Name != srcName {
			return fmt.Errorf("name mismatch: %v", src.Name)
		}
		if src.Arn == nil || !strings.Contains(*src.Arn, "delivery-source:"+srcName) {
			return fmt.Errorf("delivery source ARN form: %v", src.Arn)
		}
		if len(src.ResourceArns) != 1 || src.ResourceArns[0] != *groupArn {
			return fmt.Errorf("resourceArns mismatch: %v", src.ResourceArns)
		}
		if src.Service == nil || *src.Service != "logs" {
			return fmt.Errorf("service: %v", src.Service)
		}
		if src.Status != types.DeliverySourceStatusActive {
			return fmt.Errorf("status with live group: %v", src.Status)
		}

		got, err := tc.client.GetDeliverySource(tc.ctx, &cloudwatchlogs.GetDeliverySourceInput{Name: aws.String(srcName)})
		if err != nil {
			return fmt.Errorf("get delivery source: %v", err)
		}
		if got.DeliverySource.Status != types.DeliverySourceStatusActive {
			return fmt.Errorf("get status: %v", got.DeliverySource.Status)
		}

		all, err := tc.describeAllDeliverySources()
		if err != nil {
			return fmt.Errorf("describe delivery sources: %v", err)
		}
		found := false
		for _, s := range all {
			if s.Name != nil && *s.Name == srcName {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("source missing from the listing (%d sources)", len(all))
		}

		// The status tracks the underlying resource: deleting the group
		// flips the source to INACTIVE with RESOURCE_DELETED.
		tc.deleteLogGroup(group)
		got, err = tc.client.GetDeliverySource(tc.ctx, &cloudwatchlogs.GetDeliverySourceInput{Name: aws.String(srcName)})
		if err != nil {
			return fmt.Errorf("get delivery source after group delete: %v", err)
		}
		if got.DeliverySource.Status != types.DeliverySourceStatusInactive {
			return fmt.Errorf("status after group delete: %v", got.DeliverySource.Status)
		}
		if got.DeliverySource.StatusReason != types.DeliverySourceStatusReasonResourceDeleted {
			return fmt.Errorf("statusReason after group delete: %v", got.DeliverySource.StatusReason)
		}

		_, err = tc.client.DeleteDeliverySource(tc.ctx, &cloudwatchlogs.DeleteDeliverySourceInput{Name: aws.String(srcName)})
		if err != nil {
			return fmt.Errorf("delete delivery source: %v", err)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "PutDeliverySource_RejectsUnsupportedSource", func() error {
		var ve *types.ValidationException
		_, err := tc.client.PutDeliverySource(tc.ctx, &cloudwatchlogs.PutDeliverySourceInput{
			Name:        aws.String(tc.uniquePrefix("vpc-src")),
			ResourceArn: aws.String(fmt.Sprintf("arn:aws:vpc:%s:000000000000:vpc/vpc-0123456789abcdef0", tc.region)),
			LogType:     aws.String("FLOW_LOGS"),
		})
		if !errors.As(err, &ve) {
			return fmt.Errorf("vpc resource: want ValidationException, got %v", err)
		}
		_, err = tc.client.PutDeliverySource(tc.ctx, &cloudwatchlogs.PutDeliverySourceInput{
			Name:        aws.String(tc.uniquePrefix("raw-src")),
			ResourceArn: aws.String("not-an-arn"),
			LogType:     aws.String("APPLICATION_LOGS"),
		})
		if !errors.As(err, &ve) {
			return fmt.Errorf("non-ARN resource: want ValidationException, got %v", err)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "PutDeliveryDestination_CWLTypeAndRejections", func() error {
		destName := tc.uniquePrefix("vended-dest")
		group := destName + "-group"
		if err := tc.createLogGroup(group); err != nil {
			return fmt.Errorf("create group: %v", err)
		}
		defer tc.deleteLogGroup(group)
		groupArn, err := tc.findLogGroupARN(group)
		if err != nil {
			return err
		}

		put, err := tc.client.PutDeliveryDestination(tc.ctx, &cloudwatchlogs.PutDeliveryDestinationInput{
			Name: aws.String(destName),
			DeliveryDestinationConfiguration: &types.DeliveryDestinationConfiguration{
				DestinationResourceArn: groupArn,
			},
		})
		if err != nil {
			return fmt.Errorf("put delivery destination: %v", err)
		}
		// The successful put must not leak the destination into later
		// runs' listings.
		defer tc.client.DeleteDeliveryDestination(tc.ctx, &cloudwatchlogs.DeleteDeliveryDestinationInput{
			Name: aws.String(destName),
		})
		dest := put.DeliveryDestination
		if dest.DeliveryDestinationType != types.DeliveryDestinationTypeCwl {
			return fmt.Errorf("derived destination type: %v", dest.DeliveryDestinationType)
		}
		if dest.Arn == nil || !strings.Contains(*dest.Arn, "delivery-destination:"+destName) {
			return fmt.Errorf("destination ARN form: %v", dest.Arn)
		}
		if dest.DeliveryDestinationConfiguration == nil || dest.DeliveryDestinationConfiguration.DestinationResourceArn == nil ||
			*dest.DeliveryDestinationConfiguration.DestinationResourceArn != *groupArn {
			return fmt.Errorf("destination configuration echo: %v", dest.DeliveryDestinationConfiguration)
		}

		var ve *types.ValidationException
		_, err = tc.client.PutDeliveryDestination(tc.ctx, &cloudwatchlogs.PutDeliveryDestinationInput{
			Name: aws.String(tc.uniquePrefix("parquet-dest")),
			DeliveryDestinationConfiguration: &types.DeliveryDestinationConfiguration{
				DestinationResourceArn: groupArn,
			},
			OutputFormat: types.OutputFormatParquet,
		})
		if !errors.As(err, &ve) {
			return fmt.Errorf("parquet: want ValidationException, got %v", err)
		}
		_, err = tc.client.PutDeliveryDestination(tc.ctx, &cloudwatchlogs.PutDeliveryDestinationInput{
			Name: aws.String(tc.uniquePrefix("fh-dest")),
			DeliveryDestinationConfiguration: &types.DeliveryDestinationConfiguration{
				DestinationResourceArn: aws.String(fmt.Sprintf("arn:aws:firehose:%s:000000000000:deliverystream/vended", tc.region)),
			},
		})
		if !errors.As(err, &ve) {
			return fmt.Errorf("firehose: want ValidationException, got %v", err)
		}
		_, err = tc.client.PutDeliveryDestination(tc.ctx, &cloudwatchlogs.PutDeliveryDestinationInput{
			Name:                    aws.String(tc.uniquePrefix("xray-dest")),
			DeliveryDestinationType: types.DeliveryDestinationTypeXray,
		})
		if !errors.As(err, &ve) {
			return fmt.Errorf("xray: want ValidationException, got %v", err)
		}
		var rnf *types.ResourceNotFoundException
		_, err = tc.client.PutDeliveryDestination(tc.ctx, &cloudwatchlogs.PutDeliveryDestinationInput{
			Name: aws.String(tc.uniquePrefix("gone-dest")),
			DeliveryDestinationConfiguration: &types.DeliveryDestinationConfiguration{
				DestinationResourceArn: aws.String(fmt.Sprintf("arn:aws:logs:%s:000000000000:log-group:no-such-group-at-all", tc.region)),
			},
		})
		if !errors.As(err, &rnf) {
			return fmt.Errorf("missing destination group: want ResourceNotFoundException, got %v", err)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "DeliveryDestinationPolicy_Roundtrip", func() error {
		destName := tc.uniquePrefix("policy-dest")
		group := destName + "-group"
		if err := tc.createLogGroup(group); err != nil {
			return fmt.Errorf("create group: %v", err)
		}
		defer tc.deleteLogGroup(group)
		groupArn, err := tc.findLogGroupARN(group)
		if err != nil {
			return err
		}
		if _, err := tc.client.PutDeliveryDestination(tc.ctx, &cloudwatchlogs.PutDeliveryDestinationInput{
			Name: aws.String(destName),
			DeliveryDestinationConfiguration: &types.DeliveryDestinationConfiguration{
				DestinationResourceArn: groupArn,
			},
		}); err != nil {
			return fmt.Errorf("put delivery destination: %v", err)
		}
		defer tc.client.DeleteDeliveryDestination(tc.ctx, &cloudwatchlogs.DeleteDeliveryDestinationInput{
			Name: aws.String(destName),
		})

		policyDoc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"logs:CreateDelivery","Resource":"*"}]}`
		put, err := tc.client.PutDeliveryDestinationPolicy(tc.ctx, &cloudwatchlogs.PutDeliveryDestinationPolicyInput{
			DeliveryDestinationName:   aws.String(destName),
			DeliveryDestinationPolicy: aws.String(policyDoc),
		})
		if err != nil {
			return fmt.Errorf("put policy: %v", err)
		}
		if put.Policy == nil || put.Policy.DeliveryDestinationPolicy == nil || *put.Policy.DeliveryDestinationPolicy != policyDoc {
			return fmt.Errorf("policy echo mismatch: %v", put.Policy)
		}
		got, err := tc.client.GetDeliveryDestinationPolicy(tc.ctx, &cloudwatchlogs.GetDeliveryDestinationPolicyInput{
			DeliveryDestinationName: aws.String(destName),
		})
		if err != nil {
			return fmt.Errorf("get policy: %v", err)
		}
		if got.Policy == nil || got.Policy.DeliveryDestinationPolicy == nil || *got.Policy.DeliveryDestinationPolicy != policyDoc {
			return fmt.Errorf("policy get mismatch: %v", got.Policy)
		}
		if _, err := tc.client.DeleteDeliveryDestinationPolicy(tc.ctx, &cloudwatchlogs.DeleteDeliveryDestinationPolicyInput{
			DeliveryDestinationName: aws.String(destName),
		}); err != nil {
			return fmt.Errorf("delete policy: %v", err)
		}
		var rnf *types.ResourceNotFoundException
		_, err = tc.client.GetDeliveryDestinationPolicy(tc.ctx, &cloudwatchlogs.GetDeliveryDestinationPolicyInput{
			DeliveryDestinationName: aws.String(destName),
		})
		if !errors.As(err, &rnf) {
			return fmt.Errorf("deleted policy: want ResourceNotFoundException, got %v", err)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "CreateDelivery_Lifecycle", func() error {
		srcName := tc.uniquePrefix("life-src")
		destName := tc.uniquePrefix("life-dest")
		srcGroup := srcName + "-group"
		destGroup := destName + "-group"
		if err := tc.createLogGroup(srcGroup); err != nil {
			return fmt.Errorf("create source group: %v", err)
		}
		defer tc.deleteLogGroup(srcGroup)
		if err := tc.createLogGroup(destGroup); err != nil {
			return fmt.Errorf("create destination group: %v", err)
		}
		defer tc.deleteLogGroup(destGroup)
		srcArn, err := tc.findLogGroupARN(srcGroup)
		if err != nil {
			return err
		}
		destGroupArn, err := tc.findLogGroupARN(destGroup)
		if err != nil {
			return err
		}
		if _, err := tc.client.PutDeliverySource(tc.ctx, &cloudwatchlogs.PutDeliverySourceInput{
			Name: aws.String(srcName), ResourceArn: srcArn, LogType: aws.String("APPLICATION_LOGS"),
		}); err != nil {
			return fmt.Errorf("put source: %v", err)
		}
		destPut, err := tc.client.PutDeliveryDestination(tc.ctx, &cloudwatchlogs.PutDeliveryDestinationInput{
			Name: aws.String(destName),
			DeliveryDestinationConfiguration: &types.DeliveryDestinationConfiguration{
				DestinationResourceArn: destGroupArn,
			},
		})
		if err != nil {
			return fmt.Errorf("put destination: %v", err)
		}

		created, err := tc.client.CreateDelivery(tc.ctx, &cloudwatchlogs.CreateDeliveryInput{
			DeliverySourceName:     aws.String(srcName),
			DeliveryDestinationArn: destPut.DeliveryDestination.Arn,
			RecordFields:           []string{"time", "message"},
		})
		if err != nil {
			return fmt.Errorf("create delivery: %v", err)
		}
		// A failure past the create must not leak the delivery and its
		// source and destination into later runs' listings; on the
		// success path the explicit deletions below have already run and
		// these deferred calls fail silently on already-deleted names.
		defer func() {
			_, _ = tc.client.DeleteDelivery(tc.ctx, &cloudwatchlogs.DeleteDeliveryInput{Id: created.Delivery.Id})
			_, _ = tc.client.DeleteDeliverySource(tc.ctx, &cloudwatchlogs.DeleteDeliverySourceInput{Name: aws.String(srcName)})
			_, _ = tc.client.DeleteDeliveryDestination(tc.ctx, &cloudwatchlogs.DeleteDeliveryDestinationInput{Name: aws.String(destName)})
		}()
		delivery := created.Delivery
		if delivery.Id == nil || *delivery.Id == "" {
			return fmt.Errorf("delivery id missing")
		}
		if delivery.DeliveryDestinationType != types.DeliveryDestinationTypeCwl {
			return fmt.Errorf("delivery destination type: %v", delivery.DeliveryDestinationType)
		}
		if len(delivery.RecordFields) != 2 || delivery.RecordFields[0] != "time" {
			return fmt.Errorf("recordFields: %v", delivery.RecordFields)
		}

		var conflict *types.ConflictException
		_, err = tc.client.CreateDelivery(tc.ctx, &cloudwatchlogs.CreateDeliveryInput{
			DeliverySourceName:     aws.String(srcName),
			DeliveryDestinationArn: destPut.DeliveryDestination.Arn,
		})
		if !errors.As(err, &conflict) {
			return fmt.Errorf("duplicate pairing: want ConflictException, got %v", err)
		}

		got, err := tc.client.GetDelivery(tc.ctx, &cloudwatchlogs.GetDeliveryInput{Id: delivery.Id})
		if err != nil {
			return fmt.Errorf("get delivery: %v", err)
		}
		if got.Delivery.Arn == nil || *got.Delivery.Arn != *delivery.Arn {
			return fmt.Errorf("delivery ARN echo: %v vs %v", got.Delivery.Arn, delivery.Arn)
		}

		all, err := tc.describeAllDeliveries()
		if err != nil {
			return fmt.Errorf("describe deliveries: %v", err)
		}
		found := false
		for _, d := range all {
			if d.Id != nil && *d.Id == *delivery.Id {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("delivery missing from the listing (%d deliveries)", len(all))
		}

		if _, err := tc.client.UpdateDeliveryConfiguration(tc.ctx, &cloudwatchlogs.UpdateDeliveryConfigurationInput{
			Id:           delivery.Id,
			RecordFields: []string{"message"},
		}); err != nil {
			return fmt.Errorf("update delivery configuration: %v", err)
		}
		got, err = tc.client.GetDelivery(tc.ctx, &cloudwatchlogs.GetDeliveryInput{Id: delivery.Id})
		if err != nil {
			return fmt.Errorf("get delivery after update: %v", err)
		}
		if len(got.Delivery.RecordFields) != 1 || got.Delivery.RecordFields[0] != "message" {
			return fmt.Errorf("updated recordFields: %v", got.Delivery.RecordFields)
		}

		// The referenced source and destination refuse deletion.
		_, err = tc.client.DeleteDeliverySource(tc.ctx, &cloudwatchlogs.DeleteDeliverySourceInput{Name: aws.String(srcName)})
		if !errors.As(err, &conflict) {
			return fmt.Errorf("delete referenced source: want ConflictException, got %v", err)
		}
		_, err = tc.client.DeleteDeliveryDestination(tc.ctx, &cloudwatchlogs.DeleteDeliveryDestinationInput{Name: aws.String(destName)})
		if !errors.As(err, &conflict) {
			return fmt.Errorf("delete referenced destination: want ConflictException, got %v", err)
		}

		if _, err := tc.client.DeleteDelivery(tc.ctx, &cloudwatchlogs.DeleteDeliveryInput{Id: delivery.Id}); err != nil {
			return fmt.Errorf("delete delivery: %v", err)
		}
		var rnf *types.ResourceNotFoundException
		_, err = tc.client.GetDelivery(tc.ctx, &cloudwatchlogs.GetDeliveryInput{Id: delivery.Id})
		if !errors.As(err, &rnf) {
			return fmt.Errorf("deleted delivery: want ResourceNotFoundException, got %v", err)
		}
		if _, err := tc.client.DeleteDeliverySource(tc.ctx, &cloudwatchlogs.DeleteDeliverySourceInput{Name: aws.String(srcName)}); err != nil {
			return fmt.Errorf("delete unreferenced source: %v", err)
		}
		if _, err := tc.client.DeleteDeliveryDestination(tc.ctx, &cloudwatchlogs.DeleteDeliveryDestinationInput{Name: aws.String(destName)}); err != nil {
			return fmt.Errorf("delete unreferenced destination: %v", err)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "DescribeConfigurationTemplates_Catalog", func() error {
		// Walk every page: the catalog listing is paginated like the
		// other describes.
		var templates []types.ConfigurationTemplate
		var nextToken *string
		for {
			resp, err := tc.client.DescribeConfigurationTemplates(tc.ctx, &cloudwatchlogs.DescribeConfigurationTemplatesInput{
				NextToken: nextToken,
			})
			if err != nil {
				return fmt.Errorf("describe configuration templates: %v", err)
			}
			templates = append(templates, resp.ConfigurationTemplates...)
			if resp.NextToken == nil || *resp.NextToken == "" {
				break
			}
			nextToken = resp.NextToken
		}
		typesSeen := map[string]bool{}
		for _, tmpl := range templates {
			if tmpl.Service == nil || *tmpl.Service != "logs" {
				return fmt.Errorf("template service: %v", tmpl.Service)
			}
			typesSeen[string(tmpl.DeliveryDestinationType)] = true
			if tmpl.DeliveryDestinationType == types.DeliveryDestinationTypeS3 {
				hasJSON, hasParquet := false, false
				for _, f := range tmpl.AllowedOutputFormats {
					if f == types.OutputFormatJson {
						hasJSON = true
					}
					if f == types.OutputFormatParquet {
						hasParquet = true
					}
				}
				if !hasJSON {
					return fmt.Errorf("S3 template omits json output format")
				}
				if hasParquet {
					return fmt.Errorf("S3 template offers parquet")
				}
				for _, f := range tmpl.AllowedFields {
					if f.Name == nil || (*f.Name != "time" && *f.Name != "stream" && *f.Name != "message") {
						return fmt.Errorf("S3 template allowed field %v outside the vocabulary", f.Name)
					}
					if f.Mandatory != nil && *f.Mandatory {
						return fmt.Errorf("S3 template marks %q mandatory; the platform's source vocabulary carries no mandatory field", *f.Name)
					}
				}
			}
		}
		if !typesSeen["CWL"] || !typesSeen["S3"] {
			return fmt.Errorf("template destination types: %v", typesSeen)
		}

		filtered, err := tc.client.DescribeConfigurationTemplates(tc.ctx, &cloudwatchlogs.DescribeConfigurationTemplatesInput{
			DeliveryDestinationTypes: []types.DeliveryDestinationType{types.DeliveryDestinationTypeS3},
		})
		if err != nil {
			return fmt.Errorf("filtered describe: %v", err)
		}
		if len(filtered.ConfigurationTemplates) != 1 {
			return fmt.Errorf("filtered rows: %d", len(filtered.ConfigurationTemplates))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "DeliveryEngine_CWLFlow", func() error {
		srcName := tc.uniquePrefix("flow-src")
		destName := tc.uniquePrefix("flow-dest")
		srcGroup := srcName + "-group"
		destGroup := destName + "-group"
		const stream = "engine/1"
		if err := tc.createLogGroup(srcGroup); err != nil {
			return fmt.Errorf("create source group: %v", err)
		}
		defer tc.deleteLogGroup(srcGroup)
		if err := tc.createLogGroup(destGroup); err != nil {
			return fmt.Errorf("create destination group: %v", err)
		}
		defer tc.deleteLogGroup(destGroup)
		if err := tc.createLogStream(srcGroup, stream); err != nil {
			return fmt.Errorf("create stream: %v", err)
		}
		now := time.Now().UnixMilli()
		if err := tc.putLogEvent(srcGroup, stream, "engine-first", now-2000); err != nil {
			return err
		}
		if err := tc.putLogEvent(srcGroup, stream, "engine-second", now-1000); err != nil {
			return err
		}

		srcArn, err := tc.findLogGroupARN(srcGroup)
		if err != nil {
			return err
		}
		destGroupArn, err := tc.findLogGroupARN(destGroup)
		if err != nil {
			return err
		}
		if _, err := tc.client.PutDeliverySource(tc.ctx, &cloudwatchlogs.PutDeliverySourceInput{
			Name: aws.String(srcName), ResourceArn: srcArn, LogType: aws.String("APPLICATION_LOGS"),
		}); err != nil {
			return fmt.Errorf("put source: %v", err)
		}
		destPut, err := tc.client.PutDeliveryDestination(tc.ctx, &cloudwatchlogs.PutDeliveryDestinationInput{
			Name: aws.String(destName),
			DeliveryDestinationConfiguration: &types.DeliveryDestinationConfiguration{
				DestinationResourceArn: destGroupArn,
			},
		})
		if err != nil {
			return fmt.Errorf("put destination: %v", err)
		}
		// A failure past the puts must not leak the source and destination
		// into later runs' listings; the delivery defer must unwind first —
		// a referenced source refuses deletion with ConflictException.
		defer func() {
			_, _ = tc.client.DeleteDeliverySource(tc.ctx, &cloudwatchlogs.DeleteDeliverySourceInput{Name: aws.String(srcName)})
			_, _ = tc.client.DeleteDeliveryDestination(tc.ctx, &cloudwatchlogs.DeleteDeliveryDestinationInput{Name: aws.String(destName)})
		}()
		created, err := tc.client.CreateDelivery(tc.ctx, &cloudwatchlogs.CreateDeliveryInput{
			DeliverySourceName:     aws.String(srcName),
			DeliveryDestinationArn: destPut.DeliveryDestination.Arn,
		})
		if err != nil {
			return fmt.Errorf("create delivery: %v", err)
		}
		defer func() {
			_, _ = tc.client.DeleteDelivery(tc.ctx, &cloudwatchlogs.DeleteDeliveryInput{Id: created.Delivery.Id})
		}()

		// The delivery engine runs on the cadence pass; poll the
		// destination group for the two delivered events.
		deadline := time.Now().Add(15 * time.Second)
		for time.Now().Before(deadline) {
			events, err := tc.client.FilterLogEvents(tc.ctx, &cloudwatchlogs.FilterLogEventsInput{
				LogGroupName: aws.String(destGroup),
			})
			if err == nil && len(events.Events) >= 2 {
				// The delivered form is shaped per the delivery's default
				// recordFields (time, stream and the message joined by the
				// default delimiter) — the raw message re-ingests under no
				// format.
				first := *events.Events[0].Message
				wantFirst := fmt.Sprintf("%d\t%s\tengine-first", now-2000, stream)
				if first != wantFirst {
					return fmt.Errorf("first delivered message: %q want %q", first, wantFirst)
				}
				// The events keep their originating stream.
				if events.Events[0].LogStreamName == nil || *events.Events[0].LogStreamName != stream {
					return fmt.Errorf("delivered stream name: %v", events.Events[0].LogStreamName)
				}
				_, err := tc.client.DeleteDelivery(tc.ctx, &cloudwatchlogs.DeleteDeliveryInput{Id: created.Delivery.Id})
				if err != nil {
					return fmt.Errorf("delete delivery: %v", err)
				}
				return nil
			}
			time.Sleep(500 * time.Millisecond)
		}
		return fmt.Errorf("delivered events never reached the destination group")
	}))

	return results
}
