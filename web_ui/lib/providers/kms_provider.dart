import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../src/generated/kms.pbgrpc.dart';
import 'api_providers.dart';

class KMSKey {
  final String keyId;
  final String keyArn;
  final String state;
  final String description;
  final DateTime created;

  KMSKey({
    required this.keyId,
    required this.keyArn,
    this.state = 'Enabled',
    this.description = '',
    required this.created,
  });
}

class KMSState {
  final List<KMSKey> keys;
  final bool isLoading;
  final String? error;

  const KMSState({this.keys = const [], this.isLoading = false, this.error});

  KMSState copyWith({
    List<KMSKey>? keys,
    bool? isLoading,
    String? error,
    bool clearError = false,
  }) {
    return KMSState(
      keys: keys ?? this.keys,
      isLoading: isLoading ?? this.isLoading,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

class KMSListNotifier extends Notifier<KMSState> {
  @override
  KMSState build() => const KMSState();

  Future<void> loadKeys() async {
    state = state.copyWith(isLoading: true, clearError: true);

    try {
      final client = ref.read(kmsProvider);
      final request = ListKeysRequest();
      final response = await client.listKeys(request);

      final keys = response.keys.map((k) {
        return KMSKey(
          keyId: k.keyid,
          keyArn: k.keyarn,
          created: DateTime.now(),
        );
      });

      state = state.copyWith(keys: keys.toList(), isLoading: false);
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  Future<void> createKey({String? description}) async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      final client = ref.read(kmsProvider);
      final request = CreateKeyRequest(description: description);
      await client.createKey(request);
      await loadKeys();
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  Future<void> enableKey(String keyId) async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      final client = ref.read(kmsProvider);
      final request = EnableKeyRequest(keyid: keyId);
      await client.enableKey(request);
      await loadKeys();
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  Future<void> disableKey(String keyId) async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      final client = ref.read(kmsProvider);
      final request = DisableKeyRequest(keyid: keyId);
      await client.disableKey(request);
      await loadKeys();
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }
}

final kmsListProvider = NotifierProvider<KMSListNotifier, KMSState>(() {
  return KMSListNotifier();
});
