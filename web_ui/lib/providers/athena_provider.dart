import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../src/generated/athena.pbgrpc.dart';
import '../src/generated/athena.pb.dart';
import 'api_providers.dart';

class Workgroup {
  final String name;
  final String arn;
  final String state;
  final String? description;
  final DateTime? creationTime;
  final bool enforceWorkGroupConfiguration;

  Workgroup({
    required this.name,
    this.arn = '',
    this.state = '',
    this.description,
    this.creationTime,
    this.enforceWorkGroupConfiguration = false,
  });
}

class NamedQuery {
  final String name;
  final String? description;
  final String database;
  final String queryString;
  final String workGroup;
  final DateTime? creationTime;

  NamedQuery({
    required this.name,
    this.description,
    this.database = '',
    this.queryString = '',
    this.workGroup = '',
    this.creationTime,
  });
}

class AthenaState {
  final List<Workgroup> workgroups;
  final List<NamedQuery> namedQueries;
  final bool isLoading;
  final String? error;

  const AthenaState({
    this.workgroups = const [],
    this.namedQueries = const [],
    this.isLoading = false,
    this.error,
  });

  AthenaState copyWith({
    List<Workgroup>? workgroups,
    List<NamedQuery>? namedQueries,
    bool? isLoading,
    String? error,
    bool clearError = false,
  }) {
    return AthenaState(
      workgroups: workgroups ?? this.workgroups,
      namedQueries: namedQueries ?? this.namedQueries,
      isLoading: isLoading ?? this.isLoading,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

class AthenaListNotifier extends Notifier<AthenaState> {
  @override
  AthenaState build() => const AthenaState();

  Future<void> loadData() async {
    state = state.copyWith(isLoading: true, clearError: true);

    try {
      final client = ref.read(athenaProvider);

      final wgRequest = ListWorkGroupsInput();
      final wgResponse = await client.listWorkGroups(wgRequest);

      final workgroups = wgResponse.workgroups.map((wg) {
        return Workgroup(
          name: wg.name,
          state: wg.state.name.replaceAll('WORK_GROUP_STATE_', ''),
          description: wg.description,
        );
      });

      state = state.copyWith(workgroups: workgroups.toList(), isLoading: false);
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }
}

final athenaListProvider = NotifierProvider<AthenaListNotifier, AthenaState>(
  () {
    return AthenaListNotifier();
  },
);
