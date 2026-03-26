import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../src/generated/ingest.timestream.pbgrpc.dart';
import 'api_providers.dart';

class TimestreamDatabase {
  final String databaseName;
  final String arn;
  final int tableCount;
  final String creationTime;

  TimestreamDatabase({
    required this.databaseName,
    required this.arn,
    required this.tableCount,
    required this.creationTime,
  });
}

class TimestreamWriteState {
  final List<TimestreamDatabase> databases;
  final bool isLoading;
  final String? error;

  const TimestreamWriteState({
    this.databases = const [],
    this.isLoading = false,
    this.error,
  });

  TimestreamWriteState copyWith({
    List<TimestreamDatabase>? databases,
    bool? isLoading,
    String? error,
    bool clearError = false,
  }) {
    return TimestreamWriteState(
      databases: databases ?? this.databases,
      isLoading: isLoading ?? this.isLoading,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

class TimestreamWriteListNotifier extends Notifier<TimestreamWriteState> {
  @override
  TimestreamWriteState build() => const TimestreamWriteState();

  Future<void> loadDatabases() async {
    state = state.copyWith(isLoading: true, clearError: true);

    try {
      final client = ref.read(timestreamWriteProvider);
      final request = ListDatabasesRequest();
      final response = await client.listDatabases(request);

      final databases = response.databases.map((d) {
        return TimestreamDatabase(
          databaseName: d.databasename,
          arn: d.arn,
          tableCount: d.tablecount.toInt(),
          creationTime: d.creationtime,
        );
      }).toList();

      state = state.copyWith(databases: databases, isLoading: false);
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }
}

final timestreamWriteListProvider =
    NotifierProvider<TimestreamWriteListNotifier, TimestreamWriteState>(() {
      return TimestreamWriteListNotifier();
    });
