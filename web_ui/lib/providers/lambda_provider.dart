import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../src/generated/lambda.pbgrpc.dart';
import '../src/generated/lambda.pb.dart';
import 'api_providers.dart';

class LambdaFunction {
  final String functionName;
  final String runtime;
  final String role;
  final String status;
  final int memorySize;
  final int timeout;
  final DateTime created;

  LambdaFunction({
    required this.functionName,
    required this.runtime,
    required this.role,
    this.status = 'Active',
    this.memorySize = 128,
    this.timeout = 3,
    required this.created,
  });
}

class InvocationResult {
  final bool isSuccessful;
  final int statusCode;
  final String body;
  final Duration duration;

  InvocationResult({
    required this.isSuccessful,
    required this.statusCode,
    required this.body,
    required this.duration,
  });
}

class LambdaState {
  final List<LambdaFunction> functions;
  final bool isLoading;
  final String? error;

  const LambdaState({
    this.functions = const [],
    this.isLoading = false,
    this.error,
  });

  LambdaState copyWith({
    List<LambdaFunction>? functions,
    bool? isLoading,
    String? error,
    bool clearError = false,
  }) {
    return LambdaState(
      functions: functions ?? this.functions,
      isLoading: isLoading ?? this.isLoading,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

class LambdaListNotifier extends Notifier<LambdaState> {
  @override
  LambdaState build() => const LambdaState();

  Future<void> loadFunctions() async {
    state = state.copyWith(isLoading: true, clearError: true);

    try {
      final client = ref.read(lambdaProvider);
      final request = ListFunctionsRequest();
      final response = await client.listFunctions(request);

      final functions = response.functions.map((f) {
        return LambdaFunction(
          functionName: f.functionname,
          runtime: f.runtime.name.replaceAll('RUNTIME_', ''),
          role: f.role,
          status: f.state.name.replaceAll('STATE_', ''),
          memorySize: f.memorysize,
          timeout: f.timeout,
          created: f.lastmodified.isNotEmpty
              ? DateTime.parse(f.lastmodified)
              : DateTime.now(),
        );
      });

      state = state.copyWith(functions: functions.toList(), isLoading: false);
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  Runtime _parseRuntime(String runtime) {
    final runtimeMap = {
      'python3.11': Runtime.RUNTIME_PYTHON311,
      'python3.10': Runtime.RUNTIME_PYTHON310,
      'python3.9': Runtime.RUNTIME_PYTHON39,
      'nodejs20.x': Runtime.RUNTIME_NODEJS20X,
      'nodejs18.x': Runtime.RUNTIME_NODEJS18X,
      'nodejs16.x': Runtime.RUNTIME_NODEJS16X,
      'java17': Runtime.RUNTIME_JAVA17,
      'java11': Runtime.RUNTIME_JAVA11,
      'go1.x': Runtime.RUNTIME_GO1X,
      'ruby3.2': Runtime.RUNTIME_RUBY32,
      'dotnet8': Runtime.RUNTIME_DOTNET8,
    };
    return runtimeMap[runtime] ?? Runtime.RUNTIME_PYTHON311;
  }

  Future<void> createFunction({
    required String name,
    required String runtime,
    required String role,
    int memorySize = 128,
    int timeout = 30,
  }) async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      final client = ref.read(lambdaProvider);
      final runtimeEnum = _parseRuntime(runtime);
      final request = CreateFunctionRequest(
        functionname: name,
        runtime: runtimeEnum,
        role: role,
        handler: 'index.handler',
        memorysize: memorySize,
        timeout: timeout,
      );
      await client.createFunction(request);
      await loadFunctions();
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  Future<void> deleteFunction(String name) async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      final client = ref.read(lambdaProvider);
      final request = DeleteFunctionRequest(functionname: name);
      await client.deleteFunction(request);
      await loadFunctions();
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  Future<void> updateFunction({
    required String name,
    int? memorySize,
    int? timeout,
    String? status,
  }) async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      final client = ref.read(lambdaProvider);
      final request = UpdateFunctionConfigurationRequest(
        functionname: name,
        memorysize: memorySize,
        timeout: timeout,
      );
      await client.updateFunctionConfiguration(request);
      await loadFunctions();
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  Future<InvocationResult> invokeFunction(
    String name, {
    String payload = '{}',
  }) async {
    return InvocationResult(
      isSuccessful: false,
      statusCode: 501,
      body: 'Invoke not yet implemented in gRPC service',
      duration: Duration.zero,
    );
  }
}

final lambdaListProvider = NotifierProvider<LambdaListNotifier, LambdaState>(
  () {
    return LambdaListNotifier();
  },
);
