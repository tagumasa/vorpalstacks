package testutil

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"vorpalstacks-sdk-tests/config"
)

// cognitoImportJobTests exercises the CSV user import data plane end to
// end: job creation with the model-required CloudWatchLogsRoleArn, CSV
// upload to the returned URL, the asynchronous row-by-row import with
// per-line CloudWatch Logs outcomes, password-hash import with immediate
// sign-in, stopping a running job, and the negative start/stop paths.
func (r *TestRunner) cognitoImportJobTests(ctx context.Context, client *cognitoidentityprovider.Client) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("cognito", "UserImportJob_ImportFlow", func() error {
		return r.runUserImportFlowTest(ctx, client)
	}))
	results = append(results, r.RunTest("cognito", "UserImportJob_PasswordHashImport", func() error {
		return r.runUserImportHashTest(ctx, client)
	}))
	results = append(results, r.RunTest("cognito", "UserImportJob_Stop", func() error {
		return r.runUserImportStopTest(ctx, client)
	}))
	results = append(results, r.RunTest("cognito", "UserImportJob_NegativePaths", func() error {
		return r.runUserImportNegativeTest(ctx, client)
	}))

	return results
}

// createImportTestPool creates a pool with email as the auto-verified
// attribute, a required email, and an optional custom rank attribute, plus
// an app client allowed to use USER_PASSWORD_AUTH.
func (r *TestRunner) createImportTestPool(ctx context.Context, client *cognitoidentityprovider.Client, poolName string) (string, string, error) {
	createPool, err := client.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
		PoolName:               aws.String(poolName),
		AutoVerifiedAttributes: []types.VerifiedAttributeType{types.VerifiedAttributeTypeEmail},
		Schema: []types.SchemaAttributeType{
			{Name: aws.String("email"), Required: aws.Bool(true)},
			{Name: aws.String("rank"), AttributeDataType: types.AttributeDataTypeString, Mutable: aws.Bool(true)},
		},
		Policies: &types.UserPoolPolicyType{
			PasswordPolicy: &types.PasswordPolicyType{
				MinimumLength:    aws.Int32(8),
				RequireUppercase: true,
				RequireLowercase: true,
				RequireNumbers:   true,
			},
		},
	})
	if err != nil {
		return "", "", err
	}
	poolID := aws.ToString(createPool.UserPool.Id)

	createClient, err := client.CreateUserPoolClient(ctx, &cognitoidentityprovider.CreateUserPoolClientInput{
		UserPoolId:        aws.String(poolID),
		ClientName:        aws.String(poolName + "-client"),
		ExplicitAuthFlows: []types.ExplicitAuthFlowsType{types.ExplicitAuthFlowsTypeUserPasswordAuth},
	})
	if err != nil {
		_, _ = client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{UserPoolId: aws.String(poolID)})
		return "", "", err
	}
	return poolID, aws.ToString(createClient.UserPoolClient.ClientId), nil
}

// importCSV builds a CSV file from the pool's own header: every column the
// GetCSVHeader response lists is emitted, values come from the row maps,
// and columns without a value stay blank.
func importCSV(ctx context.Context, client *cognitoidentityprovider.Client, poolID string, rows []map[string]string) (string, error) {
	headerResp, err := client.GetCSVHeader(ctx, &cognitoidentityprovider.GetCSVHeaderInput{UserPoolId: aws.String(poolID)})
	if err != nil {
		return "", err
	}
	if len(headerResp.CSVHeader) == 0 {
		return "", fmt.Errorf("GetCSVHeader returned an empty header")
	}
	// The AWS CLI example response carries the request's UserPoolId
	// alongside the header list.
	if aws.ToString(headerResp.UserPoolId) != poolID {
		return "", fmt.Errorf("GetCSVHeader UserPoolId = %q, want %q", aws.ToString(headerResp.UserPoolId), poolID)
	}
	var b strings.Builder
	b.WriteString(strings.Join(headerResp.CSVHeader, ","))
	b.WriteString("\n")
	for _, row := range rows {
		values := make([]string, 0, len(headerResp.CSVHeader))
		for _, column := range headerResp.CSVHeader {
			values = append(values, row[column])
		}
		b.WriteString(strings.Join(values, ","))
		b.WriteString("\n")
	}
	return b.String(), nil
}

// uploadImportCSV PUTs the CSV bytes to the job's upload URL, mirroring the
// documented curl -T upload against the presigned target.
func uploadImportCSV(url, csv string) error {
	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(csv))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/csv")
	resp, err := awshttp.NewBuildableClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("CSV upload returned HTTP %d", resp.StatusCode)
	}
	return nil
}

// waitImportJobTerminal polls DescribeUserImportJob until the job leaves
// the Created/Pending/InProgress/Stopping states or the deadline passes.
func waitImportJobTerminal(ctx context.Context, client *cognitoidentityprovider.Client, poolID, jobID string, deadline time.Duration) (*types.UserImportJobType, error) {
	deadlineAt := time.Now().Add(deadline)
	for {
		resp, err := client.DescribeUserImportJob(ctx, &cognitoidentityprovider.DescribeUserImportJobInput{
			UserPoolId: aws.String(poolID),
			JobId:      aws.String(jobID),
		})
		if err != nil {
			return nil, err
		}
		job := resp.UserImportJob
		switch job.Status {
		case types.UserImportJobStatusTypeSucceeded, types.UserImportJobStatusTypeFailed,
			types.UserImportJobStatusTypeStopped, types.UserImportJobStatusTypeExpired:
			return job, nil
		}
		if time.Now().After(deadlineAt) {
			return job, fmt.Errorf("import job %s did not reach a terminal state within %s (status %s)", jobID, deadline, job.Status)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func (r *TestRunner) runUserImportFlowTest(ctx context.Context, client *cognitoidentityprovider.Client) error {
	poolName := fmt.Sprintf("import-flow-pool-%d", time.Now().UnixNano())
	poolID, clientID, err := r.createImportTestPool(ctx, client, poolName)
	if err != nil {
		return fmt.Errorf("pool setup failed: %w", err)
	}
	defer func() {
		_, _ = client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{UserPoolId: aws.String(poolID)})
	}()

	csv, err := importCSV(ctx, client, poolID, []map[string]string{
		{"cognito:username": "imported-alice", "name": "Alice", "email": "alice.import@example.com", "email_verified": "true", "birthdate": "03/14/1990", "custom:rank": "gold"},
		{"cognito:username": "imported-bob", "email": "bob.import@example.com", "email_verified": "true"},
		{"cognito:username": "imported-bad", "email": "bad.import@example.com", "email_verified": "false"},
		{"cognito:username": "imported-alice", "email": "alice.import@example.com", "email_verified": "true"},
	})
	if err != nil {
		return err
	}

	createResp, err := client.CreateUserImportJob(ctx, &cognitoidentityprovider.CreateUserImportJobInput{
		JobName:               aws.String("import-flow-job"),
		UserPoolId:            aws.String(poolID),
		CloudWatchLogsRoleArn: aws.String("arn:aws:iam::123456789012:role/cognito-import-role"),
	})
	if err != nil {
		return err
	}
	job := createResp.UserImportJob
	if job.Status != types.UserImportJobStatusTypeCreated {
		return fmt.Errorf("initial status = %s, want Created", job.Status)
	}
	if aws.ToString(job.PreSignedUrl) == "" {
		return fmt.Errorf("PreSignedUrl is empty")
	}
	// AWS example responses always carry the counters, zero-valued on a
	// freshly created job.
	if job.ImportedUsers != 0 || job.SkippedUsers != 0 || job.FailedUsers != 0 {
		return fmt.Errorf("create response counters = %d/%d/%d, want 0/0/0",
			job.ImportedUsers, job.SkippedUsers, job.FailedUsers)
	}
	jobID := aws.ToString(job.JobId)

	if err := uploadImportCSV(aws.ToString(job.PreSignedUrl), csv); err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}

	if _, err := client.StartUserImportJob(ctx, &cognitoidentityprovider.StartUserImportJobInput{
		UserPoolId: aws.String(poolID),
		JobId:      aws.String(jobID),
	}); err != nil {
		return err
	}

	final, err := waitImportJobTerminal(ctx, client, poolID, jobID, 30*time.Second)
	if err != nil {
		return err
	}
	if final.Status != types.UserImportJobStatusTypeSucceeded {
		return fmt.Errorf("final status = %s (%s), want Succeeded", final.Status, aws.ToString(final.CompletionMessage))
	}
	if final.ImportedUsers != 2 {
		return fmt.Errorf("ImportedUsers = %d, want 2", final.ImportedUsers)
	}
	if final.SkippedUsers != 1 {
		return fmt.Errorf("SkippedUsers = %d, want 1", final.SkippedUsers)
	}
	if final.FailedUsers != 1 {
		return fmt.Errorf("FailedUsers = %d, want 1", final.FailedUsers)
	}

	userResp, err := client.AdminGetUser(ctx, &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(poolID),
		Username:   aws.String("imported-alice"),
	})
	if err != nil {
		return fmt.Errorf("imported user not found: %w", err)
	}
	if userResp.UserStatus != types.UserStatusTypeResetRequired {
		return fmt.Errorf("imported user status = %s, want RESET_REQUIRED", userResp.UserStatus)
	}
	attrs := map[string]string{}
	for _, a := range userResp.UserAttributes {
		attrs[aws.ToString(a.Name)] = aws.ToString(a.Value)
	}
	if attrs["custom:rank"] != "gold" {
		return fmt.Errorf("custom:rank = %q, want gold", attrs["custom:rank"])
	}
	if attrs["birthdate"] != "1990-03-14" {
		return fmt.Errorf("birthdate = %q, want converted 1990-03-14", attrs["birthdate"])
	}
	if attrs["email_verified"] != "true" {
		return fmt.Errorf("email_verified = %q, want true", attrs["email_verified"])
	}

	// An imported RESET_REQUIRED user signs in with any password and is
	// prompted to set a new one; completing the challenge confirms the
	// imported account is usable.
	initResp, err := client.InitiateAuth(ctx, &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		AuthParameters: map[string]string{
			"USERNAME": "imported-alice",
			"PASSWORD": "anypassword123",
		},
		ClientId: aws.String(clientID),
	})
	if err != nil {
		return fmt.Errorf("sign-in with any password failed: %w", err)
	}
	if initResp.ChallengeName != types.ChallengeNameTypeNewPasswordRequired {
		return fmt.Errorf("challenge = %s, want NEW_PASSWORD_REQUIRED", initResp.ChallengeName)
	}
	challengeResp, err := client.RespondToAuthChallenge(ctx, &cognitoidentityprovider.RespondToAuthChallengeInput{
		ClientId:      aws.String(clientID),
		ChallengeName: types.ChallengeNameTypeNewPasswordRequired,
		Session:       initResp.Session,
		ChallengeResponses: map[string]string{
			"USERNAME":     "imported-alice",
			"NEW_PASSWORD": "NewP@ssw0rd1",
		},
	})
	if err != nil {
		return fmt.Errorf("new-password challenge failed: %w", err)
	}
	if challengeResp.AuthenticationResult == nil || aws.ToString(challengeResp.AuthenticationResult.AccessToken) == "" {
		return fmt.Errorf("no tokens issued after completing the new-password challenge")
	}

	return r.assertImportLogs(ctx, poolID, poolName, jobID, "import-flow-job")
}

// assertImportLogs verifies the per-line import outcomes in CloudWatch
// Logs: group /aws/cognito/userpools/{POOL_ID}/{POOL_NAME}, stream
// {JOB_ID}/{JOB_NAME}.
func (r *TestRunner) assertImportLogs(ctx context.Context, poolID, poolName, jobID, jobName string) error {
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return fmt.Errorf("logs config failed: %w", err)
	}
	logsClient := cloudwatchlogs.NewFromConfig(cfg)
	logGroup := fmt.Sprintf("/aws/cognito/userpools/%s/%s", poolID, poolName)
	streams, err := logsClient.DescribeLogStreams(ctx, &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: aws.String(logGroup),
	})
	if err != nil {
		return fmt.Errorf("DescribeLogStreams failed: %w", err)
	}
	var streamName string
	for _, s := range streams.LogStreams {
		if strings.HasPrefix(aws.ToString(s.LogStreamName), jobID+"/") {
			streamName = aws.ToString(s.LogStreamName)
			break
		}
	}
	if streamName == "" {
		return fmt.Errorf("no log stream for job %s in group %s (streams: %d)", jobID, logGroup, len(streams.LogStreams))
	}
	events, err := logsClient.GetLogEvents(ctx, &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(logGroup),
		LogStreamName: aws.String(streamName),
	})
	if err != nil {
		return fmt.Errorf("GetLogEvents failed: %w", err)
	}
	var succeeded, skipped, failed bool
	for _, e := range events.Events {
		message := aws.ToString(e.Message)
		if strings.Contains(message, "[SUCCEEDED]") {
			succeeded = true
		}
		if strings.Contains(message, "[SKIPPED]") {
			skipped = true
		}
		if strings.Contains(message, "[FAILED]") {
			failed = true
		}
	}
	if !succeeded || !skipped || !failed {
		return fmt.Errorf("import log outcomes incomplete (succeeded=%v skipped=%v failed=%v)", succeeded, skipped, failed)
	}
	return nil
}

// importedHashBCrypt verifies "Import3d!Pass" (bcrypt cost 4).
const importedHashBCrypt = "$2a$04$YobgTqYzyZqdRBYyNUR6teChu./N4lVp2O0MWbNGwlvPEP6cYW/ke"

func (r *TestRunner) runUserImportHashTest(ctx context.Context, client *cognitoidentityprovider.Client) error {
	poolName := fmt.Sprintf("import-hash-pool-%d", time.Now().UnixNano())
	poolID, clientID, err := r.createImportTestPool(ctx, client, poolName)
	if err != nil {
		return fmt.Errorf("pool setup failed: %w", err)
	}
	defer func() {
		_, _ = client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{UserPoolId: aws.String(poolID)})
	}()

	// The GetCSVHeader output does not carry the password_hash column;
	// AWS documents it as an additional column of the file.
	headerResp, err := client.GetCSVHeader(ctx, &cognitoidentityprovider.GetCSVHeaderInput{UserPoolId: aws.String(poolID)})
	if err != nil {
		return err
	}
	hashRow := map[string]string{
		"cognito:username": "hashimport-carol",
		"email":            "carol.import@example.com",
		"email_verified":   "true",
	}
	columns := append(append([]string{}, headerResp.CSVHeader...), "password_hash")
	values := make([]string, len(columns))
	for i, column := range columns {
		if column == "password_hash" {
			values[i] = importedHashBCrypt
		} else {
			values[i] = hashRow[column]
		}
	}
	csv := strings.Join(columns, ",") + "\n" + strings.Join(values, ",") + "\n"

	createResp, err := client.CreateUserImportJob(ctx, &cognitoidentityprovider.CreateUserImportJobInput{
		JobName:                  aws.String("import-hash-job"),
		UserPoolId:               aws.String(poolID),
		CloudWatchLogsRoleArn:    aws.String("arn:aws:iam::123456789012:role/cognito-import-role"),
		PasswordHashingAlgorithm: types.PasswordHashingAlgorithmTypeBcrypt,
	})
	if err != nil {
		return err
	}
	jobID := aws.ToString(createResp.UserImportJob.JobId)
	if err := uploadImportCSV(aws.ToString(createResp.UserImportJob.PreSignedUrl), csv); err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}
	if _, err := client.StartUserImportJob(ctx, &cognitoidentityprovider.StartUserImportJobInput{
		UserPoolId: aws.String(poolID),
		JobId:      aws.String(jobID),
	}); err != nil {
		return err
	}
	final, err := waitImportJobTerminal(ctx, client, poolID, jobID, 30*time.Second)
	if err != nil {
		return err
	}
	if final.Status != types.UserImportJobStatusTypeSucceeded || final.ImportedUsers != 1 {
		return fmt.Errorf("status = %s, ImportedUsers = %d, want Succeeded/1 (%s)", final.Status, final.ImportedUsers, aws.ToString(final.CompletionMessage))
	}

	userResp, err := client.AdminGetUser(ctx, &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(poolID),
		Username:   aws.String("hashimport-carol"),
	})
	if err != nil {
		return fmt.Errorf("imported user not found: %w", err)
	}
	if userResp.UserStatus != types.UserStatusTypeConfirmed {
		return fmt.Errorf("hash-imported user status = %s, want CONFIRMED", userResp.UserStatus)
	}

	// The imported hash must verify at sign-in; the first successful
	// sign-in transparently migrates the credentials.
	authResp, err := client.InitiateAuth(ctx, &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		AuthParameters: map[string]string{
			"USERNAME": "hashimport-carol",
			"PASSWORD": "Import3d!Pass",
		},
		ClientId: aws.String(clientID),
	})
	if err != nil {
		return fmt.Errorf("sign-in with the imported password failed: %w", err)
	}
	if authResp.AuthenticationResult == nil || aws.ToString(authResp.AuthenticationResult.AccessToken) == "" {
		return fmt.Errorf("no tokens issued for the imported-hash user")
	}
	return nil
}

func (r *TestRunner) runUserImportStopTest(ctx context.Context, client *cognitoidentityprovider.Client) error {
	poolName := fmt.Sprintf("import-stop-pool-%d", time.Now().UnixNano())
	poolID, _, err := r.createImportTestPool(ctx, client, poolName)
	if err != nil {
		return fmt.Errorf("pool setup failed: %w", err)
	}
	defer func() {
		_, _ = client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{UserPoolId: aws.String(poolID)})
	}()

	// Enough rows that the worker is still processing when the stop
	// request lands.
	rows := make([]map[string]string, 0, 5000)
	for i := 0; i < 5000; i++ {
		rows = append(rows, map[string]string{
			"cognito:username": fmt.Sprintf("stopuser-%04d", i),
			"email":            fmt.Sprintf("stopuser-%04d@example.com", i),
			"email_verified":   "true",
		})
	}
	csv, err := importCSV(ctx, client, poolID, rows)
	if err != nil {
		return err
	}

	createResp, err := client.CreateUserImportJob(ctx, &cognitoidentityprovider.CreateUserImportJobInput{
		JobName:               aws.String("import-stop-job"),
		UserPoolId:            aws.String(poolID),
		CloudWatchLogsRoleArn: aws.String("arn:aws:iam::123456789012:role/cognito-import-role"),
	})
	if err != nil {
		return err
	}
	jobID := aws.ToString(createResp.UserImportJob.JobId)
	if err := uploadImportCSV(aws.ToString(createResp.UserImportJob.PreSignedUrl), csv); err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}
	if _, err := client.StartUserImportJob(ctx, &cognitoidentityprovider.StartUserImportJobInput{
		UserPoolId: aws.String(poolID),
		JobId:      aws.String(jobID),
	}); err != nil {
		return err
	}
	stopResp, err := client.StopUserImportJob(ctx, &cognitoidentityprovider.StopUserImportJobInput{
		UserPoolId: aws.String(poolID),
		JobId:      aws.String(jobID),
	})
	if err != nil {
		return fmt.Errorf("stop failed: %w", err)
	}
	// The AWS example Stop response reports the terminal Stopped status
	// with the completion date already set, not the transitional state.
	if stopResp.UserImportJob == nil || stopResp.UserImportJob.Status != types.UserImportJobStatusTypeStopped {
		return fmt.Errorf("stop response status = %v, want Stopped", stopResp.UserImportJob)
	}
	if stopResp.UserImportJob.CompletionDate == nil {
		return fmt.Errorf("stop response has no CompletionDate")
	}

	final, err := waitImportJobTerminal(ctx, client, poolID, jobID, 30*time.Second)
	if err != nil {
		return err
	}
	if final.Status != types.UserImportJobStatusTypeStopped {
		return fmt.Errorf("final status = %s, want Stopped", final.Status)
	}
	if aws.ToString(final.CompletionMessage) != "The Import Job was stopped by the developer." {
		return fmt.Errorf("completion message = %q, want the AWS stop wording", aws.ToString(final.CompletionMessage))
	}
	return nil
}

func (r *TestRunner) runUserImportNegativeTest(ctx context.Context, client *cognitoidentityprovider.Client) error {
	if _, err := client.StartUserImportJob(ctx, &cognitoidentityprovider.StartUserImportJobInput{
		UserPoolId: aws.String(r.region + "_nonexistentpool"),
		JobId:      aws.String("import-nonexistent"),
	}); err == nil {
		return fmt.Errorf("starting a non-existent job succeeded, want ResourceNotFoundException")
	} else if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
		return err
	}

	poolName := fmt.Sprintf("import-neg-pool-%d", time.Now().UnixNano())
	poolID, _, err := r.createImportTestPool(ctx, client, poolName)
	if err != nil {
		return fmt.Errorf("pool setup failed: %w", err)
	}
	defer func() {
		_, _ = client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{UserPoolId: aws.String(poolID)})
	}()

	csv, err := importCSV(ctx, client, poolID, []map[string]string{
		{"cognito:username": "neg-dave", "email": "dave.import@example.com", "email_verified": "true"},
	})
	if err != nil {
		return err
	}
	createResp, err := client.CreateUserImportJob(ctx, &cognitoidentityprovider.CreateUserImportJobInput{
		JobName:               aws.String("import-neg-job"),
		UserPoolId:            aws.String(poolID),
		CloudWatchLogsRoleArn: aws.String("arn:aws:iam::123456789012:role/cognito-import-role"),
	})
	if err != nil {
		return err
	}
	jobID := aws.ToString(createResp.UserImportJob.JobId)
	if err := uploadImportCSV(aws.ToString(createResp.UserImportJob.PreSignedUrl), csv); err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}
	if _, err := client.StartUserImportJob(ctx, &cognitoidentityprovider.StartUserImportJobInput{
		UserPoolId: aws.String(poolID),
		JobId:      aws.String(jobID),
	}); err != nil {
		return err
	}
	final, err := waitImportJobTerminal(ctx, client, poolID, jobID, 30*time.Second)
	if err != nil {
		return err
	}
	if final.Status != types.UserImportJobStatusTypeSucceeded {
		return fmt.Errorf("warm-up job status = %s, want Succeeded", final.Status)
	}

	if _, err := client.StartUserImportJob(ctx, &cognitoidentityprovider.StartUserImportJobInput{
		UserPoolId: aws.String(poolID),
		JobId:      aws.String(jobID),
	}); err == nil {
		return fmt.Errorf("restarting a finished job succeeded, want InvalidParameterException")
	} else if err := AssertErrorContains(err, "InvalidParameterException"); err != nil {
		return err
	}

	secondResp, err := client.CreateUserImportJob(ctx, &cognitoidentityprovider.CreateUserImportJobInput{
		JobName:               aws.String("import-neg-job-2"),
		UserPoolId:            aws.String(poolID),
		CloudWatchLogsRoleArn: aws.String("arn:aws:iam::123456789012:role/cognito-import-role"),
	})
	if err != nil {
		return err
	}
	if _, err := client.StopUserImportJob(ctx, &cognitoidentityprovider.StopUserImportJobInput{
		UserPoolId: aws.String(poolID),
		JobId:      aws.String(aws.ToString(secondResp.UserImportJob.JobId)),
	}); err == nil {
		return fmt.Errorf("stopping a Created job succeeded, want InvalidParameterException")
	} else if err := AssertErrorContains(err, "InvalidParameterException"); err != nil {
		return err
	}
	return nil
}
