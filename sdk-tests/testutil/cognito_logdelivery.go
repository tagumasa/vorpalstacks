package testutil

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	cognitotypes "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"vorpalstacks-sdk-tests/config"
)

// cognitoLogDeliveryTests exercises the user-pool log-delivery family: the
// S3 destination export of userAuthEvents records, the CloudWatch Logs
// export of userNotification message-delivery records, and the destination
// validation contract (exactly one destination per configuration, the log
// level bound to the event source, userNotification CloudWatch-only).
func (r *TestRunner) cognitoLogDeliveryTests(tc *cognitoIDPContext) []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "cognito",
			TestName: "LogDelivery_Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) { o.UsePathStyle = true })
	cwlClient := cloudwatchlogs.NewFromConfig(cfg)

	// logDeliverySignIn performs one password sign-in, whose SignIn event
	// drives the userAuthEvents export.
	logDeliverySignIn := func(username, password, clientID string) error {
		_, err := tc.client.InitiateAuth(tc.ctx, &cognitoidentityprovider.InitiateAuthInput{
			AuthFlow: cognitotypes.AuthFlowTypeUserPasswordAuth,
			AuthParameters: map[string]string{
				"USERNAME": username,
				"PASSWORD": password,
			},
			ClientId: aws.String(clientID),
		})
		return err
	}
	logDeliveryPasswordClient := func() (string, error) {
		name := tc.unique("logdelivery-client")
		out, err := tc.client.CreateUserPoolClient(tc.ctx, &cognitoidentityprovider.CreateUserPoolClientInput{
			UserPoolId: aws.String(tc.userPoolID),
			ClientName: aws.String(name),
			ExplicitAuthFlows: []cognitotypes.ExplicitAuthFlowsType{
				cognitotypes.ExplicitAuthFlowsTypeAllowUserPasswordAuth,
			},
		})
		if err != nil {
			return "", err
		}
		return *out.UserPoolClient.ClientId, nil
	}

	results = append(results, r.RunTest("cognito", "LogDelivery_S3AuthEventExport", func() error {
		bucket := tc.unique("cognito-auth-logs")
		if _, err := s3Client.CreateBucket(tc.ctx, &s3.CreateBucketInput{Bucket: aws.String(bucket)}); err != nil {
			return err
		}
		defer func() {
			out, err := s3Client.ListObjectsV2(tc.ctx, &s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
			if err == nil && len(out.Contents) > 0 {
				objects := make([]s3types.ObjectIdentifier, 0, len(out.Contents))
				for _, o := range out.Contents {
					objects = append(objects, s3types.ObjectIdentifier{Key: o.Key})
				}
				_, _ = s3Client.DeleteObjects(tc.ctx, &s3.DeleteObjectsInput{
					Bucket: aws.String(bucket),
					Delete: &s3types.Delete{Objects: objects},
				})
			}
			_, _ = s3Client.DeleteBucket(tc.ctx, &s3.DeleteBucketInput{Bucket: aws.String(bucket)})
		}()

		if _, err := tc.client.SetLogDeliveryConfiguration(tc.ctx, &cognitoidentityprovider.SetLogDeliveryConfigurationInput{
			UserPoolId: aws.String(tc.userPoolID),
			LogConfigurations: []cognitotypes.LogConfigurationType{{
				EventSource:     cognitotypes.EventSourceNameUserAuthEvents,
				LogLevel:        cognitotypes.LogLevelInfo,
				S3Configuration: &cognitotypes.S3ConfigurationType{BucketArn: aws.String("arn:aws:s3:::" + bucket)},
			}},
		}); err != nil {
			return err
		}

		clientID, err := logDeliveryPasswordClient()
		if err != nil {
			return err
		}
		username := tc.unique("logdelivery-signin")
		if _, err := tc.client.AdminCreateUser(tc.ctx, &cognitoidentityprovider.AdminCreateUserInput{
			UserPoolId: aws.String(tc.userPoolID),
			Username:   aws.String(username),
		}); err != nil {
			return err
		}
		if _, err := tc.client.AdminSetUserPassword(tc.ctx, &cognitoidentityprovider.AdminSetUserPasswordInput{
			UserPoolId: aws.String(tc.userPoolID),
			Username:   aws.String(username),
			Password:   aws.String("LogDelivery123!"),
			Permanent:  true,
		}); err != nil {
			return err
		}
		if err := logDeliverySignIn(username, "LogDelivery123!", clientID); err != nil {
			return err
		}

		var key string
		deadline := time.Now().Add(20 * time.Second)
		for time.Now().Before(deadline) {
			out, err := s3Client.ListObjectsV2(tc.ctx, &s3.ListObjectsV2Input{
				Bucket: aws.String(bucket),
				Prefix: aws.String("AWSLogs/"),
			})
			if err != nil {
				return err
			}
			if len(out.Contents) > 0 {
				key = *out.Contents[0].Key
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
		if key == "" {
			return fmt.Errorf("no exported auth-event object landed under AWSLogs/ in %s", bucket)
		}

		obj, err := s3Client.GetObject(tc.ctx, &s3.GetObjectInput{Bucket: aws.String(bucket), Key: aws.String(key)})
		if err != nil {
			return err
		}
		body, err := io.ReadAll(obj.Body)
		if err != nil {
			return err
		}
		if !strings.Contains(string(body), `"USER_ACTIVITY"`) {
			return fmt.Errorf("exported record is not a userAuthEvents record: %s", body)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "LogDelivery_UserNotificationCloudWatch", func() error {
		logGroup := "cognito-logdelivery-" + tc.ts
		if _, err := cwlClient.CreateLogGroup(tc.ctx, &cloudwatchlogs.CreateLogGroupInput{LogGroupName: aws.String(logGroup)}); err != nil {
			return err
		}
		defer func() {
			_, _ = cwlClient.DeleteLogGroup(tc.ctx, &cloudwatchlogs.DeleteLogGroupInput{LogGroupName: aws.String(logGroup)})
		}()

		if _, err := tc.client.SetLogDeliveryConfiguration(tc.ctx, &cognitoidentityprovider.SetLogDeliveryConfigurationInput{
			UserPoolId: aws.String(tc.userPoolID),
			LogConfigurations: []cognitotypes.LogConfigurationType{{
				EventSource: cognitotypes.EventSourceNameUserNotification,
				LogLevel:    cognitotypes.LogLevelError,
				CloudWatchLogsConfiguration: &cognitotypes.CloudWatchLogsConfigurationType{
					LogGroupArn: aws.String(fmt.Sprintf("arn:aws:logs:%s:000000000000:log-group:%s", r.region, logGroup)),
				},
			}},
		}); err != nil {
			return err
		}

		// A sign-up delivers a verification code, driving the
		// userNotification export into the log group.
		clientID, err := logDeliveryPasswordClient()
		if err != nil {
			return err
		}
		username := tc.unique("logdelivery-notify")
		if _, err := tc.client.SignUp(tc.ctx, &cognitoidentityprovider.SignUpInput{
			ClientId: aws.String(clientID),
			Username: aws.String(username),
			Password: aws.String("LogDelivery123!"),
		}); err != nil {
			return err
		}

		stream := "cognito-" + tc.userPoolID
		deadline := time.Now().Add(20 * time.Second)
		for time.Now().Before(deadline) {
			out, err := cwlClient.GetLogEvents(tc.ctx, &cloudwatchlogs.GetLogEventsInput{
				LogGroupName:  aws.String(logGroup),
				LogStreamName: aws.String(stream),
			})
			if err == nil {
				for _, ev := range out.Events {
					if strings.Contains(*ev.Message, `"USER_NOTIFICATION"`) {
						return nil
					}
				}
			}
			time.Sleep(500 * time.Millisecond)
		}
		return fmt.Errorf("no USER_NOTIFICATION record landed in log group %s stream %s", logGroup, stream)
	}))

	results = append(results, r.RunTest("cognito", "LogDelivery_DestinationValidation", func() error {
		reject := func(lc cognitotypes.LogConfigurationType) error {
			_, err := tc.client.SetLogDeliveryConfiguration(tc.ctx, &cognitoidentityprovider.SetLogDeliveryConfigurationInput{
				UserPoolId:        aws.String(tc.userPoolID),
				LogConfigurations: []cognitotypes.LogConfigurationType{lc},
			})
			if err == nil {
				return fmt.Errorf("configuration %+v was accepted", lc)
			}
			return nil
		}
		if err := reject(cognitotypes.LogConfigurationType{
			EventSource: cognitotypes.EventSourceNameUserAuthEvents,
			LogLevel:    cognitotypes.LogLevelInfo,
		}); err != nil {
			return err
		}
		if err := reject(cognitotypes.LogConfigurationType{
			EventSource:                 cognitotypes.EventSourceNameUserAuthEvents,
			LogLevel:                    cognitotypes.LogLevelInfo,
			S3Configuration:             &cognitotypes.S3ConfigurationType{BucketArn: aws.String("arn:aws:s3:::two-dest")},
			CloudWatchLogsConfiguration: &cognitotypes.CloudWatchLogsConfigurationType{LogGroupArn: aws.String("arn:aws:logs:us-east-1:000000000000:log-group:two-dest")},
		}); err != nil {
			return err
		}
		if err := reject(cognitotypes.LogConfigurationType{
			EventSource:     cognitotypes.EventSourceNameUserNotification,
			LogLevel:        cognitotypes.LogLevelError,
			S3Configuration: &cognitotypes.S3ConfigurationType{BucketArn: aws.String("arn:aws:s3:::notify-dest")},
		}); err != nil {
			return err
		}
		if err := reject(cognitotypes.LogConfigurationType{
			EventSource: cognitotypes.EventSourceNameUserAuthEvents,
			LogLevel:    cognitotypes.LogLevelError,
			CloudWatchLogsConfiguration: &cognitotypes.CloudWatchLogsConfigurationType{
				LogGroupArn: aws.String("arn:aws:logs:us-east-1:000000000000:log-group:level-mismatch"),
			},
		}); err != nil {
			return err
		}
		return nil
	}))

	return results
}
