package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron/v3"

	"vorpalstacks/internal/common/endpoint"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
	svcarn "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/timeutils"
)

// schedulerTickerInterval is the interval between schedule checks.
// Reduced to 1 second in test mode for faster test execution.
var schedulerTickerInterval = 1 * time.Minute

func init() {
	if os.Getenv("TEST_MODE") == "true" {
		schedulerTickerInterval = 1 * time.Second
	}
}

// Engine manages scheduled task execution for EventBridge Scheduler.
type Engine struct {
	storageManager *storage.RegionStorageManager
	accountID      string
	bus            eventbus.Bus
	stores         sync.Map // region → *schedulerstore.SchedulerStore

	running   bool
	runningMu sync.RWMutex
	stopChan  chan struct{}
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewEngine creates a new scheduler engine with the given store dependencies.
func NewEngine(
	storageManager *storage.RegionStorageManager,
	accountID string,
) *Engine {
	return &Engine{
		storageManager: storageManager,
		accountID:      accountID,
		stopChan:       make(chan struct{}),
	}
}

// SetEventBus injects the event bus for publishing scheduler lifecycle events.
func (e *Engine) SetEventBus(bus eventbus.Bus) {
	e.bus = bus
}

// Start starts the scheduler engine.
func (e *Engine) Start() error {
	e.runningMu.Lock()
	defer e.runningMu.Unlock()
	if e.running {
		return nil
	}
	e.running = true
	e.stopChan = make(chan struct{})
	e.ctx, e.cancel = context.WithCancel(context.Background())

	e.wg.Add(1)
	go e.run()

	logs.Debug("Scheduler engine started")
	return nil
}

// Stop stops the scheduler engine.
func (e *Engine) Stop() error {
	e.runningMu.Lock()
	defer e.runningMu.Unlock()
	if !e.running {
		return nil
	}
	e.running = false
	if e.cancel != nil {
		e.cancel()
	}
	close(e.stopChan)
	e.wg.Wait()

	logs.Debug("Scheduler engine stopped")
	return nil
}

func (e *Engine) run() {
	defer e.wg.Done()
	defer func() { resilience.RecoverAndRestart("scheduler engine", &e.wg, e.run) }()

	ticker := time.NewTicker(schedulerTickerInterval)
	defer ticker.Stop()

	e.checkSchedules()

	for {
		select {
		case <-e.stopChan:
			return
		case <-ticker.C:
			e.checkSchedules()
		}
	}
}

func (e *Engine) checkSchedules() {
	if e.storageManager == nil {
		return
	}

	regions := e.storageManager.GetActiveRegions()
	for _, region := range regions {
		storage, err := e.storageManager.GetStorage(region)
		if err != nil {
			logs.Debug("Failed to get storage for region", logs.String("region", region), logs.String("error", err.Error()))
			continue
		}
		store := schedulerstore.NewSchedulerStore(storage, e.accountID, region)
		if actual, loaded := e.stores.LoadOrStore(region, store); loaded {
			store = actual.(*schedulerstore.SchedulerStore)
		}
		schedules, err := store.GetAllEnabledSchedules(e.ctx)
		if err != nil {
			logs.Debug("Failed to get enabled schedules", logs.String("region", region), logs.String("error", err.Error()))
			continue
		}

		now := time.Now().UTC()

		for _, schedule := range schedules {
			schedule.Region = region
			if e.shouldExecute(schedule, now) {
				e.wg.Add(1)
				go func(sch *schedulerstore.Schedule) {
					defer e.wg.Done()
					defer func() {
						if r := recover(); r != nil {
							logs.Error("scheduler: panic executing schedule", logs.String("name", sch.Name), logs.Any("panic", r))
						}
					}()
					select {
					case <-e.ctx.Done():
						return
					default:
						e.executeSchedule(e.ctx, sch)
					}
				}(schedule)
			}
		}
	}
}

func (e *Engine) shouldExecute(schedule *schedulerstore.Schedule, now time.Time) bool {
	if schedule.StartDate != nil && now.Before(*schedule.StartDate) {
		return false
	}
	if schedule.EndDate != nil && now.After(*schedule.EndDate) {
		return false
	}

	nextTime, err := e.getNextExecutionTime(schedule, now)
	if err != nil {
		logs.Debug("Failed to calculate next execution time",
			logs.String("schedule", schedule.Name),
			logs.String("error", err.Error()))
		return false
	}

	if schedule.FlexibleTimeWindow != nil && schedule.FlexibleTimeWindow.Mode == schedulerstore.FlexibleTimeWindowModeFlexible {
		maxWindow := 1
		if schedule.FlexibleTimeWindow.MaximumWindowInMinutes != nil {
			maxWindow = *schedule.FlexibleTimeWindow.MaximumWindowInMinutes
		}
		windowStart := nextTime.Add(-time.Duration(maxWindow) * time.Minute)
		return now.After(windowStart) && !now.After(nextTime.Add(time.Duration(maxWindow)*time.Minute))
	}

	diff := now.Sub(nextTime)
	return diff >= 0 && diff < time.Minute
}

func (e *Engine) getNextExecutionTime(schedule *schedulerstore.Schedule, now time.Time) (time.Time, error) {
	expr := schedule.ScheduleExpression

	if strings.HasPrefix(expr, "at(") {
		matches := regexp.MustCompile(`^at\((\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2})\)$`).FindStringSubmatch(expr)
		if len(matches) == 2 {
			t, err := time.Parse(timeutils.ISO8601NoZFormat, matches[1])
			if err != nil {
				return time.Time{}, err
			}
			return t, nil
		}
	}

	if strings.HasPrefix(expr, "rate(") {
		matches := regexp.MustCompile(`^rate\((\d+)\s+(minute|minutes|hour|hours|day|days)\)$`).FindStringSubmatch(expr)
		if len(matches) == 3 {
			value, _ := strconv.Atoi(matches[1])
			unit := matches[2]

			var duration time.Duration
			switch strings.TrimSuffix(unit, "s") {
			case "minute":
				duration = time.Duration(value) * time.Minute
			case "hour":
				duration = time.Duration(value) * time.Hour
			case "day":
				duration = time.Duration(value) * 24 * time.Hour
			}

			creationTime := schedule.CreationDate
			elapsed := now.Sub(creationTime)
			periods := int(elapsed / duration)
			nextTime := creationTime.Add(time.Duration(periods) * duration)
			return nextTime, nil
		}
	}

	if strings.HasPrefix(expr, "cron(") {
		return e.parseCronNextTime(expr, now)
	}

	return time.Time{}, fmt.Errorf("unsupported schedule expression: %s", expr)
}

func (e *Engine) parseCronNextTime(expr string, now time.Time) (time.Time, error) {
	matches := regexp.MustCompile(`^cron\((.+)\)$`).FindStringSubmatch(expr)
	if len(matches) != 2 {
		return time.Time{}, fmt.Errorf("invalid cron expression: %s", expr)
	}

	fields := strings.Fields(matches[1])
	if len(fields) != 6 {
		return time.Time{}, fmt.Errorf("AWS cron expression must have 6 fields, got %d", len(fields))
	}

	standardCron := fmt.Sprintf("%s %s %s %s %s",
		convertAWSCronField(fields[0]),
		convertAWSCronField(fields[1]),
		convertAWSCronField(fields[2]),
		convertAWSCronField(fields[3]),
		convertAWSCronField(fields[4]),
	)

	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(standardCron)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse cron expression: %w", err)
	}

	return schedule.Next(now.Add(-time.Minute)), nil
}

func convertAWSCronField(field string) string {
	field = strings.ReplaceAll(field, "?", "*")
	field = strings.ReplaceAll(field, "SUN", "0")
	field = strings.ReplaceAll(field, "MON", "1")
	field = strings.ReplaceAll(field, "TUE", "2")
	field = strings.ReplaceAll(field, "WED", "3")
	field = strings.ReplaceAll(field, "THU", "4")
	field = strings.ReplaceAll(field, "FRI", "5")
	field = strings.ReplaceAll(field, "SAT", "6")
	field = strings.ReplaceAll(field, "JAN", "1")
	field = strings.ReplaceAll(field, "FEB", "2")
	field = strings.ReplaceAll(field, "MAR", "3")
	field = strings.ReplaceAll(field, "APR", "4")
	field = strings.ReplaceAll(field, "MAY", "5")
	field = strings.ReplaceAll(field, "JUN", "6")
	field = strings.ReplaceAll(field, "JUL", "7")
	field = strings.ReplaceAll(field, "AUG", "8")
	field = strings.ReplaceAll(field, "SEP", "9")
	field = strings.ReplaceAll(field, "OCT", "10")
	field = strings.ReplaceAll(field, "NOV", "11")
	field = strings.ReplaceAll(field, "DEC", "12")
	return field
}

func (e *Engine) executeSchedule(ctx context.Context, schedule *schedulerstore.Schedule) {
	logs.Debug("Executing schedule",
		logs.String("name", schedule.Name),
		logs.String("group", schedule.GroupName))

	if schedule.Target == nil {
		logs.Debug("Schedule has no target", logs.String("schedule", schedule.Name))
		return
	}

	target := schedule.Target
	targetArn := target.Arn
	region := schedule.Region
	if region == "" {
		region = request.DefaultRegion
	}

	if e.bus != nil {
		input := target.Input
		if input == "" {
			input = "{}"
		}
		schedEvt := &eventbus.ScheduleFiredEvent{
			ScheduleName: schedule.Name,
			ScheduleArn:  "",
			TargetArn:    targetArn,
			Input:        input,
		}
		schedEvt.Region = region
		if err := e.bus.Publish(context.Background(), schedEvt); err != nil {
			logs.Warn("Failed to publish schedule fired event to event bus",
				logs.String("schedule", schedule.Name),
				logs.String("target", targetArn),
				logs.Err(err))
			return
		}
	} else {
		if strings.Contains(targetArn, ":lambda:") {
			e.invokeLambda(ctx, schedule, target)
		} else if strings.Contains(targetArn, ":sqs:") {
			e.sendToSQS(ctx, schedule, target)
		} else if strings.Contains(targetArn, ":sns:") {
			e.publishToSNS(ctx, schedule, target)
		} else if strings.Contains(targetArn, ":states:") {
			e.startStepFunctionExecution(ctx, schedule, target)
		} else if strings.Contains(targetArn, ":events:") {
			e.sendToEventBridge(ctx, schedule, target)
		} else {
			logs.Debug("Unsupported target type", logs.String("arn", targetArn))
		}
	}

	if schedule.ActionAfterCompletion == schedulerstore.ActionAfterCompletionDelete {
		region := schedule.Region
		if region == "" {
			region = request.DefaultRegion
		}
		storage, err := e.storageManager.GetStorage(region)
		if err != nil {
			logs.Debug("Failed to get storage for region",
				logs.String("schedule", schedule.Name),
				logs.String("region", region),
				logs.String("error", err.Error()))
			return
		}
		var store *schedulerstore.SchedulerStore
		if cached, ok := e.stores.Load(region); ok {
			store = cached.(*schedulerstore.SchedulerStore)
		} else {
			store = schedulerstore.NewSchedulerStore(storage, e.accountID, region)
			if actual, loaded := e.stores.LoadOrStore(region, store); loaded {
				store = actual.(*schedulerstore.SchedulerStore)
			}
		}
		if err := store.DeleteSchedule(ctx, schedule.GroupName, schedule.Name); err != nil {
			logs.Debug("Failed to delete schedule after completion",
				logs.String("schedule", schedule.Name),
				logs.String("error", err.Error()))
		}
	}
}

func (e *Engine) invokeLambda(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) {
	if e.bus == nil {
		logs.Debug("event bus not configured for Lambda invocation", logs.String("schedule", schedule.Name))
		return
	}

	input := target.Input
	if input == "" {
		input = "{}"
	}

	functionName := svcarn.ExtractFunctionNameFromARN(target.Arn)
	if functionName == "" {
		logs.Debug("Failed to extract function name from ARN",
			logs.String("schedule", schedule.Name),
			logs.String("arn", target.Arn))
		return
	}

	logs.Debug("Invoking Lambda for schedule",
		logs.String("schedule", schedule.Name),
		logs.String("function", functionName))

	statusCode, _, err := e.bus.LambdaInvoker().InvokeForGateway(ctx, target.Arn, []byte(input))
	if err != nil {
		logs.Debug("Failed to invoke Lambda",
			logs.String("schedule", schedule.Name),
			logs.String("function", functionName),
			logs.String("error", err.Error()))
		return
	}

	logs.Debug("Lambda invocation completed",
		logs.String("schedule", schedule.Name),
		logs.String("function", functionName),
		logs.Int("statusCode", int(statusCode)))
}

func (e *Engine) sendToSQS(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) {
	if e.bus == nil {
		logs.Debug("event bus not configured for SQS delivery", logs.String("schedule", schedule.Name))
		return
	}

	parts := strings.Split(target.Arn, ":")
	if len(parts) < 6 {
		logs.Debug("Invalid SQS ARN", logs.String("arn", target.Arn))
		return
	}
	queueName := parts[5]
	queueURL := endpoint.SQSQueueURL(e.accountID, queueName)

	messageBody := target.Input
	if messageBody == "" {
		messageBody = fmt.Sprintf(`{"schedule":"%s","timestamp":"%s"}`, schedule.Name, time.Now().UTC().Format(time.RFC3339))
	}

	logs.Debug("Sending to SQS for schedule",
		logs.String("schedule", schedule.Name),
		logs.String("queue", queueName))

	_, _, err := e.bus.SQSInvoker().SendMessage(ctx, queueURL, messageBody, 0, nil)
	if err != nil {
		logs.Debug("Failed to send to SQS",
			logs.String("schedule", schedule.Name),
			logs.String("queue", queueName),
			logs.String("error", err.Error()))
	}
}

func (e *Engine) publishToSNS(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) {
	if e.bus == nil {
		logs.Debug("event bus not configured for SNS delivery", logs.String("schedule", schedule.Name))
		return
	}

	parts := strings.Split(target.Arn, ":")
	if len(parts) < 6 {
		logs.Debug("Invalid SNS ARN", logs.String("arn", target.Arn))
		return
	}
	topicArn := target.Arn

	message := target.Input
	if message == "" {
		msgPayload := map[string]interface{}{
			"schedule":  schedule.Name,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		}
		msgBytes, err := json.Marshal(msgPayload)
		if err != nil {
			logs.Debug("Failed to marshal SNS message payload",
				logs.String("schedule", schedule.Name),
				logs.String("error", err.Error()))
			return
		}
		message = string(msgBytes)
	}

	result, err := e.bus.SNSInvoker().ListSubscriptionsByTopic(ctx, topicArn)
	if err != nil {
		logs.Debug("Failed to list SNS subscriptions",
			logs.String("schedule", schedule.Name),
			logs.String("topicArn", topicArn),
			logs.String("error", err.Error()))
		return
	}

	for _, sub := range result {
		if sub.PendingConfirmation {
			continue
		}
		switch sub.Protocol {
		case "sqs":
			e.deliverSNSToSQS(ctx, schedule, sub, message)
		case "lambda":
			e.deliverSNSToLambda(ctx, schedule, sub, message)
		default:
			logs.Debug("Unsupported SNS subscription protocol",
				logs.String("schedule", schedule.Name),
				logs.String("protocol", sub.Protocol))
		}
	}

	logs.Debug("SNS delivery completed",
		logs.String("schedule", schedule.Name),
		logs.String("topic", topicArn),
		logs.Int("subscriptions", len(result)))
}

func (e *Engine) deliverSNSToSQS(ctx context.Context, schedule *schedulerstore.Schedule, sub eventbus.SubscriptionInfo, message string) {
	if e.bus == nil {
		return
	}

	endpointParts := strings.Split(sub.Endpoint, ":")
	if len(endpointParts) < 6 {
		return
	}
	queueName := endpointParts[5]
	queueURL := endpoint.SQSQueueURL(e.accountID, queueName)

	_, _, err := e.bus.SQSInvoker().SendMessage(ctx, queueURL, message, 0, nil)
	if err != nil {
		logs.Debug("Failed to deliver SNS to SQS",
			logs.String("schedule", schedule.Name),
			logs.String("queue", queueName),
			logs.String("error", err.Error()))
	}
}

func (e *Engine) deliverSNSToLambda(ctx context.Context, schedule *schedulerstore.Schedule, sub eventbus.SubscriptionInfo, message string) {
	if e.bus == nil {
		return
	}

	functionName := svcarn.ExtractFunctionNameFromARN(sub.Endpoint)
	if functionName == "" {
		return
	}

	payload := map[string]interface{}{
		"Records": []map[string]interface{}{
			{
				"EventSource":          "aws:sns",
				"EventVersion":         "1.0",
				"EventSubscriptionArn": sub.SubscriptionARN,
				"Sns": map[string]interface{}{
					"Type":      "Notification",
					"MessageId": fmt.Sprintf("%d", time.Now().UnixNano()),
					"TopicArn":  sub.TopicARN,
					"Message":   message,
					"Timestamp": time.Now().UTC().Format(time.RFC3339),
				},
			},
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logs.Debug("Failed to marshal Lambda payload",
			logs.String("schedule", schedule.Name),
			logs.String("function", functionName),
			logs.String("error", err.Error()))
		return
	}

	_, _, err = e.bus.LambdaInvoker().InvokeForGateway(ctx, sub.Endpoint, payloadBytes)
	if err != nil {
		logs.Debug("Failed to deliver SNS to Lambda",
			logs.String("schedule", schedule.Name),
			logs.String("function", functionName),
			logs.String("error", err.Error()))
	}
}

func (e *Engine) sendToCloudWatchLogs(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) {
	if e.bus == nil {
		logs.Debug("event bus not configured, skipping CloudWatch Logs delivery",
			logs.String("schedule", schedule.Name))
		return
	}

	logGroup := svcarn.ExtractLogGroupNameFromARN(target.Arn)
	if logGroup == "" {
		logs.Debug("Failed to extract log group from CloudWatch Logs ARN",
			logs.String("schedule", schedule.Name),
			logs.String("arn", target.Arn))
		return
	}

	_, _, region, _, resource := svcarn.SplitARN(target.Arn)
	var logStream string
	if idx := strings.LastIndex(resource, ":log-stream:"); idx != -1 {
		logStream = resource[idx+12:]
	} else {
		logStream = fmt.Sprintf("scheduler-%s", schedule.Name)
	}

	message := target.Input
	if message == "" {
		msgPayload := map[string]interface{}{
			"schedule":  schedule.Name,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		}
		if msgBytes, err := json.Marshal(msgPayload); err == nil {
			message = string(msgBytes)
		}
	}

	evt := &eventbus.CloudWatchLogsPutEvent{
		LogGroup:  logGroup,
		LogStream: logStream,
		LogEvents: []eventbus.LogEntry{
			{Timestamp: time.Now().UnixMilli(), Message: message},
		},
	}
	evt.Region = region
	evt.AccountID = e.accountID

	if err := e.bus.Publish(ctx, evt); err != nil {
		logs.Debug("Failed to deliver schedule to CloudWatch Logs",
			logs.String("schedule", schedule.Name),
			logs.String("logGroup", logGroup),
			logs.String("error", err.Error()))
		return
	}

	logs.Debug("Schedule delivered to CloudWatch Logs",
		logs.String("schedule", schedule.Name),
		logs.String("logGroup", logGroup))
}

func (e *Engine) startStepFunctionExecution(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) {
	if e.bus == nil {
		logs.Debug("event bus not configured, skipping Step Functions delivery",
			logs.String("schedule", schedule.Name))
		return
	}

	_, _, smRegion, _, _ := svcarn.SplitARN(target.Arn)
	if smRegion == "" {
		smRegion = schedule.Region
	}

	input := target.Input
	if input == "" {
		input = "{}"
	}

	evt := &eventbus.StepFunctionsStartExecutionEvent{
		StateMachineArn: target.Arn,
		Input:           input,
	}
	evt.Region = smRegion
	evt.AccountID = e.accountID

	if err := e.bus.Publish(ctx, evt); err != nil {
		logs.Debug("Failed to start Step Functions execution from schedule",
			logs.String("schedule", schedule.Name),
			logs.String("stateMachineArn", target.Arn),
			logs.String("error", err.Error()))
		return
	}

	logs.Debug("Schedule delivered to Step Functions",
		logs.String("schedule", schedule.Name),
		logs.String("stateMachineArn", target.Arn))
}

func (e *Engine) sendToEventBridge(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) {
	if e.bus == nil {
		logs.Debug("event bus not configured, skipping EventBridge delivery",
			logs.String("schedule", schedule.Name))
		return
	}

	_, _, ebRegion, _, resource := svcarn.SplitARN(target.Arn)
	if ebRegion == "" {
		ebRegion = schedule.Region
	}

	eventBusName := "default"
	if idx := strings.LastIndex(resource, ":event-bus/"); idx != -1 {
		eventBusName = resource[idx+len(":event-bus/"):]
	}

	input := target.Input
	if input == "" {
		input = "{}"
	}

	evt := &eventbus.EventBridgePutEventsEvent{
		EventBusName: eventBusName,
		Input:        input,
	}
	evt.Region = ebRegion
	evt.AccountID = e.accountID

	if err := e.bus.Publish(ctx, evt); err != nil {
		logs.Debug("Failed to deliver schedule to EventBridge",
			logs.String("schedule", schedule.Name),
			logs.String("eventBus", eventBusName),
			logs.String("error", err.Error()))
		return
	}

	logs.Debug("Schedule delivered to EventBridge",
		logs.String("schedule", schedule.Name),
		logs.String("eventBus", eventBusName))
}
