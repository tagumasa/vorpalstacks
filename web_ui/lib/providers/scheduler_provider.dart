import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../src/generated/scheduler.pbgrpc.dart';
import '../src/generated/scheduler.pb.dart';
import 'api_providers.dart';

class Schedule {
  final String name;
  final String arn;
  final String state;
  final String target;
  final DateTime created;
  final DateTime? lastModified;

  Schedule({
    required this.name,
    this.arn = '',
    this.state = 'ENABLED',
    this.target = '',
    required this.created,
    this.lastModified,
  });
}

class SchedulerState {
  final List<Schedule> schedules;
  final bool isLoading;
  final String? error;

  const SchedulerState({
    this.schedules = const [],
    this.isLoading = false,
    this.error,
  });

  SchedulerState copyWith({
    List<Schedule>? schedules,
    bool? isLoading,
    String? error,
    bool clearError = false,
  }) {
    return SchedulerState(
      schedules: schedules ?? this.schedules,
      isLoading: isLoading ?? this.isLoading,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

class SchedulerListNotifier extends Notifier<SchedulerState> {
  @override
  SchedulerState build() => const SchedulerState();

  Future<void> loadSchedules() async {
    state = state.copyWith(isLoading: true, clearError: true);

    try {
      final client = ref.read(schedulerProvider);
      final request = ListSchedulesInput();
      final response = await client.listSchedules(request);

      final schedules = response.schedules.map((s) {
        return Schedule(
          name: s.name,
          created: DateTime.tryParse(s.creationdate) ?? DateTime.now(),
          lastModified: s.lastmodificationdate.isNotEmpty
              ? DateTime.tryParse(s.lastmodificationdate)
              : null,
          target: s.target.arn,
        );
      });

      state = state.copyWith(schedules: schedules.toList(), isLoading: false);
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  Future<void> createSchedule({
    required String name,
    required String scheduleExpression,
    required String targetArn,
  }) async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      final client = ref.read(schedulerProvider);
      final request = CreateScheduleInput(
        name: name,
        flexibletimewindow: FlexibleTimeWindow(),
        scheduleexpression: scheduleExpression,
        target: Target(arn: targetArn),
      );
      await client.createSchedule(request);
      await loadSchedules();
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  Future<void> deleteSchedule(String name) async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      final client = ref.read(schedulerProvider);
      final request = DeleteScheduleInput(name: name);
      await client.deleteSchedule(request);
      await loadSchedules();
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }
}

final schedulerListProvider =
    NotifierProvider<SchedulerListNotifier, SchedulerState>(() {
      return SchedulerListNotifier();
    });
