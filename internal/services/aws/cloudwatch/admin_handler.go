package cloudwatch

import (
	"context"
	"net/http"
	"vorpalstacks/internal/common/defaults"

	svcerrors "vorpalstacks/internal/common/errors"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/cloudwatch"
	cloudwatchconnect "vorpalstacks/internal/pb/aws/cloudwatch/cloudwatchconnect"
	pbcommon "vorpalstacks/internal/pb/aws/common"
)

// AdminHandler implements the CloudWatch gRPC-Web admin console handler. It
// exposes list and describe operations for alarms and metrics for the Flutter
// management UI. It delegates to the shared CloudWatchService Core methods.
type AdminHandler struct {
	cloudwatchconnect.UnimplementedCloudWatchServiceHandler
	service *CloudWatchService
}

var _ cloudwatchconnect.CloudWatchServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new CloudWatch admin console handler.
func NewAdminHandler(svc *CloudWatchService) *AdminHandler {
	return &AdminHandler{service: svc}
}

func (h *AdminHandler) getStoreFromHeaders(headers http.Header) (*cloudwatchStores, error) {
	region := defaults.GetRegionFromHeader(headers)
	return h.service.GetStoreForRegion(region)
}

// ListMetrics lists the specified metrics within a namespace, optionally
// filtered by metric name and dimensions.
func (h *AdminHandler) ListMetrics(ctx context.Context, req *connect.Request[pb.ListMetricsInput]) (*connect.Response[pb.ListMetricsOutput], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	result, err := h.service.listMetricsCore(stores, &ListMetricsInput{
		Namespace:  req.Msg.Namespace,
		MetricName: req.Msg.Metricname,
		Dimensions: dimensionFiltersToStore(req.Msg.Dimensions),
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.ListMetricsOutput{
		Metrics: metricsToPb(result.Metrics),
	}), nil
}

// DescribeAlarms retrieves the alarms for the specified alarm name prefix.
func (h *AdminHandler) DescribeAlarms(ctx context.Context, req *connect.Request[pb.DescribeAlarmsInput]) (*connect.Response[pb.DescribeAlarmsOutput], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	alarms, _, err := h.service.describeAlarmsCore(stores, &DescribeAlarmsInput{
		AlarmNamePrefix: req.Msg.Alarmnameprefix,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	pbAlarms := make([]*pb.MetricAlarm, 0, len(alarms))
	for _, a := range alarms {
		pbAlarms = append(pbAlarms, toPbMetricAlarm(a))
	}

	return connect.NewResponse(&pb.DescribeAlarmsOutput{
		Metricalarms: pbAlarms,
	}), nil
}

// PutMetricAlarm creates or updates a CloudWatch metric alarm via the admin
// console gRPC-Web interface.
func (h *AdminHandler) PutMetricAlarm(ctx context.Context, req *connect.Request[pb.PutMetricAlarmInput]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	input := &PutMetricAlarmInput{
		AlarmName:               req.Msg.Alarmname,
		Namespace:               req.Msg.Namespace,
		MetricName:              req.Msg.Metricname,
		Dimensions:              dimensionsToStore(req.Msg.Dimensions),
		ComparisonOperator:      fromPbComparisonOperator(req.Msg.Comparisonoperator),
		Threshold:               req.Msg.Threshold,
		EvaluationPeriods:       req.Msg.GetEvaluationperiods(),
		Period:                  req.Msg.GetPeriod(),
		Statistic:               fromPbStatistic(req.Msg.Statistic),
		TreatMissingData:        req.Msg.Treatmissingdata,
		AlarmDescription:        req.Msg.Alarmdescription,
		ActionsEnabled:          true,
		AlarmActions:            req.Msg.Alarmactions,
		OKActions:               req.Msg.Okactions,
		InsufficientDataActions: req.Msg.Insufficientdataactions,
	}
	if req.Msg.Actionsenabled != nil {
		input.ActionsEnabled = *req.Msg.Actionsenabled
	}
	if input.EvaluationPeriods == 0 {
		input.EvaluationPeriods = 1
	}
	if input.Period == 0 {
		input.Period = 60
	}
	input.DatapointsToAlarm = req.Msg.GetDatapointstoalarm()
	if input.DatapointsToAlarm == 0 {
		input.DatapointsToAlarm = input.EvaluationPeriods
	}
	if input.TreatMissingData == "" {
		input.TreatMissingData = "missing"
	}

	_, err = h.service.putMetricAlarmCore(stores, input)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// DeleteAlarms deletes one or more CloudWatch alarms via the admin console
// gRPC-Web interface.
func (h *AdminHandler) DeleteAlarms(ctx context.Context, req *connect.Request[pb.DeleteAlarmsInput]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	if err := h.service.deleteAlarmsCore(stores, &DeleteAlarmsInput{AlarmNames: req.Msg.Alarmnames}); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the CloudWatch admin console.
func NewConnectHandler(svc *CloudWatchService) (string, http.Handler) {
	return cloudwatchconnect.NewCloudWatchServiceHandler(NewAdminHandler(svc))
}
