import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../src/generated/ssm.pbgrpc.dart';
import '../src/generated/ssm.pb.dart';
import 'api_providers.dart';

class SSMParameter {
  final String name;
  final String arn;
  final String type;
  final String value;
  final int version;
  final DateTime? lastModifiedDate;
  final String? description;
  final String? dataType;

  SSMParameter({
    required this.name,
    this.arn = '',
    this.type = 'String',
    this.value = '',
    this.version = 1,
    this.lastModifiedDate,
    this.description,
    this.dataType,
  });
}

class SSMState {
  final List<SSMParameter> parameters;
  final bool isLoading;
  final String? error;

  const SSMState({
    this.parameters = const [],
    this.isLoading = false,
    this.error,
  });

  SSMState copyWith({
    List<SSMParameter>? parameters,
    bool? isLoading,
    String? error,
    bool clearError = false,
  }) {
    return SSMState(
      parameters: parameters ?? this.parameters,
      isLoading: isLoading ?? this.isLoading,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

class SSMListNotifier extends Notifier<SSMState> {
  @override
  SSMState build() => const SSMState();

  Future<void> loadParameters() async {
    state = state.copyWith(isLoading: true, clearError: true);

    try {
      final client = ref.read(ssmProvider);
      final request = GetParametersByPathRequest(path: '/');
      final response = await client.getParametersByPath(request);

      final parameters = response.parameters.map((p) {
        return SSMParameter(
          name: p.name,
          arn: p.arn,
          type: p.type.name.replaceAll('PARAMETER_TYPE_', ''),
          value: p.value,
          version: p.version.toInt(),
        );
      });

      state = state.copyWith(parameters: parameters.toList(), isLoading: false);
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  Future<void> createParameter({
    required String name,
    required String value,
    required String type,
    String? description,
  }) async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      final client = ref.read(ssmProvider);
      final typeEnum = _parseParameterType(type);
      final request = PutParameterRequest(
        name: name,
        value: value,
        type: typeEnum,
        description: description,
      );
      await client.putParameter(request);
      await loadParameters();
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  ParameterType _parseParameterType(String type) {
    switch (type.toLowerCase()) {
      case 'string':
        return ParameterType.PARAMETER_TYPE_STRING;
      case 'stringlist':
        return ParameterType.PARAMETER_TYPE_STRING_LIST;
      case 'securestring':
        return ParameterType.PARAMETER_TYPE_SECURE_STRING;
      default:
        return ParameterType.PARAMETER_TYPE_STRING;
    }
  }

  Future<void> deleteParameter(String name) async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      final client = ref.read(ssmProvider);
      final request = DeleteParameterRequest(name: name);
      await client.deleteParameter(request);
      await loadParameters();
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }
}

final ssmListProvider = NotifierProvider<SSMListNotifier, SSMState>(() {
  return SSMListNotifier();
});
