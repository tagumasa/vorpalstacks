import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../src/generated/logs.pbgrpc.dart';
import '../src/generated/logs.pb.dart';
import 'api_providers.dart';

class LogGroup {
  final String logGroupName;
  final String arn;
  final DateTime? creationTime;
  final int? retentionInDays;
  final int? storedBytes;

  LogGroup({
    required this.logGroupName,
    this.arn = '',
    this.creationTime,
    this.retentionInDays,
    this.storedBytes,
  });
}

class MetricAlarm {
  final String alarmName;
  final String arn;
  final String state;
  final String? metricName;
  final String? namespace;
  final int? evaluationPeriods;
  final int? datapointsToAlarm;
  final DateTime? stateUpdatedTimestamp;

  MetricAlarm({
    required this.alarmName,
    this.arn = '',
    this.state = 'INSUFFICIENT_DATA',
    this.metricName,
    this.namespace,
    this.evaluationPeriods,
    this.datapointsToAlarm,
    this.stateUpdatedTimestamp,
  });
}

class CloudWatchState {
  final List<LogGroup> logGroups;
  final List<MetricAlarm> alarms;
  final bool isLoading;
  final String? error;

  const CloudWatchState({
    this.logGroups = const [],
    this.alarms = const [],
    this.isLoading = false,
    this.error,
  });

  CloudWatchState copyWith({
    List<LogGroup>? logGroups,
    List<MetricAlarm>? alarms,
    bool? isLoading,
    String? error,
    bool clearError = false,
  }) {
    return CloudWatchState(
      logGroups: logGroups ?? this.logGroups,
      alarms: alarms ?? this.alarms,
      isLoading: isLoading ?? this.isLoading,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

class CloudWatchListNotifier extends Notifier<CloudWatchState> {
  @override
  CloudWatchState build() => const CloudWatchState();

  Future<void> loadData() async {
    state = state.copyWith(isLoading: true, clearError: true);

    try {
      final client = ref.read(cloudWatchLogsProvider);
      final request = ListLogGroupsRequest();
      final response = await client.listLogGroups(request);

      final logGroups = response.loggroups.map((lg) {
        return LogGroup(logGroupName: lg.loggroupname, arn: lg.loggrouparn);
      });

      state = state.copyWith(logGroups: logGroups.toList(), isLoading: false);
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  Future<void> createLogGroup(String logGroupName) async {
    state = state.copyWith(error: 'CloudWatch Admin API not yet implemented');
    await Future.delayed(const Duration(milliseconds: 100));
    state = state.copyWith(isLoading: false);
  }

  Future<void> deleteLogGroup(String logGroupName) async {
    state = state.copyWith(error: 'CloudWatch Admin API not yet implemented');
    await Future.delayed(const Duration(milliseconds: 100));
    state = state.copyWith(isLoading: false);
  }
}

final cloudwatchListProvider =
    NotifierProvider<CloudWatchListNotifier, CloudWatchState>(() {
      return CloudWatchListNotifier();
    });
