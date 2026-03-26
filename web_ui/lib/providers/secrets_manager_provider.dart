import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../src/generated/secretsmanager.pbgrpc.dart';
import '../src/generated/secretsmanager.pb.dart';
import 'api_providers.dart';

class Secret {
  final String name;
  final String arn;
  final String? description;
  final String? kmsKeyId;
  final DateTime? createdDate;
  final DateTime? lastChangedDate;
  final DateTime? lastAccessedDate;
  final int versionIdsToStagesCount;
  final bool replicationEnabled;

  Secret({
    required this.name,
    this.arn = '',
    this.description,
    this.kmsKeyId,
    this.createdDate,
    this.lastChangedDate,
    this.lastAccessedDate,
    this.versionIdsToStagesCount = 0,
    this.replicationEnabled = false,
  });
}

class SecretsManagerState {
  final List<Secret> secrets;
  final bool isLoading;
  final String? error;

  const SecretsManagerState({
    this.secrets = const [],
    this.isLoading = false,
    this.error,
  });

  SecretsManagerState copyWith({
    List<Secret>? secrets,
    bool? isLoading,
    String? error,
    bool clearError = false,
  }) {
    return SecretsManagerState(
      secrets: secrets ?? this.secrets,
      isLoading: isLoading ?? this.isLoading,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

class SecretsManagerListNotifier extends Notifier<SecretsManagerState> {
  @override
  SecretsManagerState build() => const SecretsManagerState();

  Future<void> loadSecrets() async {
    state = state.copyWith(isLoading: true, clearError: true);

    try {
      final client = ref.read(secretsManagerProvider);
      final request = ListSecretsRequest();
      final response = await client.listSecrets(request);

      final secrets = response.secretlist.map((s) {
        return Secret(name: s.name, arn: s.arn, description: s.description);
      });

      state = state.copyWith(secrets: secrets.toList(), isLoading: false);
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  Future<void> createSecret({
    required String name,
    String? description,
    String? secretString,
  }) async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      final client = ref.read(secretsManagerProvider);
      final request = CreateSecretRequest(
        name: name,
        description: description,
        secretstring: secretString,
      );
      await client.createSecret(request);
      await loadSecrets();
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  Future<void> deleteSecret(String name, {bool forceDelete = false}) async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      final client = ref.read(secretsManagerProvider);
      final request = DeleteSecretRequest(
        secretid: name,
        forcedeletewithoutrecovery: forceDelete,
      );
      await client.deleteSecret(request);
      await loadSecrets();
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  Future<String?> getSecretValue(String name) async {
    try {
      final client = ref.read(secretsManagerProvider);
      final request = GetSecretValueRequest(secretid: name);
      final response = await client.getSecretValue(request);
      return response.secretstring;
    } catch (e) {
      return null;
    }
  }
}

final secretsManagerListProvider =
    NotifierProvider<SecretsManagerListNotifier, SecretsManagerState>(() {
      return SecretsManagerListNotifier();
    });
