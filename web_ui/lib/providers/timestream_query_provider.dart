import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../src/generated/query.timestream.pbgrpc.dart';
import 'api_providers.dart';

class TimestreamScheduledQuery {
  final String name;
  final String arn;
  final String state;
  final String lastRunStatus;
  final String creationTime;

  TimestreamScheduledQuery({
    required this.name,
    required this.arn,
    required this.state,
    required this.lastRunStatus,
    required this.creationTime,
  });
}

class TimestreamQueryState {
  final List<TimestreamScheduledQuery> scheduledQueries;
  final bool isLoading;
  final String? error;

  const TimestreamQueryState({
    this.scheduledQueries = const [],
    this.isLoading = false,
    this.error,
  });

  TimestreamQueryState copyWith({
    List<TimestreamScheduledQuery>? scheduledQueries,
    bool? isLoading,
    String? error,
    bool clearError = false,
  }) {
    return TimestreamQueryState(
      scheduledQueries: scheduledQueries ?? this.scheduledQueries,
      isLoading: isLoading ?? this.isLoading,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

class TimestreamQueryListNotifier extends Notifier<TimestreamQueryState> {
  @override
  TimestreamQueryState build() => const TimestreamQueryState();

  Future<void> loadScheduledQueries() async {
    state = state.copyWith(isLoading: true, clearError: true);

    try {
      final client = ref.read(timestreamQueryProvider);
      final request = ListScheduledQueriesRequest();
      final response = await client.listScheduledQueries(request);

      final queries = response.scheduledqueries.map((q) {
        return TimestreamScheduledQuery(
          name: q.name,
          arn: q.arn,
          state: q.state.name.replaceAll('SCHEDULED_QUERY_STATE_', ''),
          lastRunStatus: q.lastrunstatus.name.replaceAll(
            'SCHEDULED_QUERY_RUN_STATUS_',
            '',
          ),
          creationTime: q.creationtime,
        );
      }).toList();

      state = state.copyWith(scheduledQueries: queries, isLoading: false);
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }
}

final timestreamQueryListProvider =
    NotifierProvider<TimestreamQueryListNotifier, TimestreamQueryState>(() {
      return TimestreamQueryListNotifier();
    });
