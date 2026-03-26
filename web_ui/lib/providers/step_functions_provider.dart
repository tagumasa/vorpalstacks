import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../src/generated/states.pbgrpc.dart';
import '../src/generated/states.pb.dart';
import 'api_providers.dart';

class StateMachine {
  final String stateMachineArn;
  final String name;
  final String type;
  final String status;
  final DateTime? creationDate;
  final DateTime? lastModified;
  final String roleArn;
  final int executionsCount;

  StateMachine({
    required this.stateMachineArn,
    required this.name,
    this.type = '',
    this.status = '',
    this.creationDate,
    this.lastModified,
    this.roleArn = '',
    this.executionsCount = 0,
  });
}

class Execution {
  final String executionArn;
  final String name;
  final String status;
  final DateTime? startDate;
  final DateTime? stopDate;

  Execution({
    required this.executionArn,
    required this.name,
    this.status = '',
    this.startDate,
    this.stopDate,
  });
}

class StepFunctionsState {
  final List<StateMachine> stateMachines;
  final List<Execution> executions;
  final bool isLoading;
  final String? error;

  const StepFunctionsState({
    this.stateMachines = const [],
    this.executions = const [],
    this.isLoading = false,
    this.error,
  });

  StepFunctionsState copyWith({
    List<StateMachine>? stateMachines,
    List<Execution>? executions,
    bool? isLoading,
    String? error,
    bool clearError = false,
  }) {
    return StepFunctionsState(
      stateMachines: stateMachines ?? this.stateMachines,
      executions: executions ?? this.executions,
      isLoading: isLoading ?? this.isLoading,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

class StepFunctionsListNotifier extends Notifier<StepFunctionsState> {
  @override
  StepFunctionsState build() => const StepFunctionsState();

  Future<void> loadStateMachines() async {
    state = state.copyWith(isLoading: true, clearError: true);

    try {
      final client = ref.read(stepFunctionsProvider);
      final request = ListStateMachinesInput();
      final response = await client.listStateMachines(request);

      final stateMachines = response.statemachines.map((sm) {
        return StateMachine(
          stateMachineArn: sm.statemachinearn,
          name: sm.name,
          type: sm.type.name.replaceAll('STATE_MACHINE_TYPE_', ''),
          creationDate: sm.creationdate.isNotEmpty
              ? DateTime.tryParse(sm.creationdate)
              : null,
        );
      });

      state = state.copyWith(
        stateMachines: stateMachines.toList(),
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }
}

final stepFunctionsListProvider =
    NotifierProvider<StepFunctionsListNotifier, StepFunctionsState>(() {
      return StepFunctionsListNotifier();
    });
