import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../src/generated/email.pbgrpc.dart';
import 'api_providers.dart';

class SESIdentity {
  final String identityName;
  final String identityType;
  final String verificationStatus;
  final bool sendingEnabled;

  SESIdentity({
    required this.identityName,
    required this.identityType,
    required this.verificationStatus,
    required this.sendingEnabled,
  });
}

class SESState {
  final List<SESIdentity> identities;
  final bool isLoading;
  final String? error;

  const SESState({
    this.identities = const [],
    this.isLoading = false,
    this.error,
  });

  SESState copyWith({
    List<SESIdentity>? identities,
    bool? isLoading,
    String? error,
    bool clearError = false,
  }) {
    return SESState(
      identities: identities ?? this.identities,
      isLoading: isLoading ?? this.isLoading,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

class SESListNotifier extends Notifier<SESState> {
  @override
  SESState build() => const SESState();

  Future<void> loadIdentities() async {
    state = state.copyWith(isLoading: true, clearError: true);

    try {
      final client = ref.read(sesProvider);
      final request = ListEmailIdentitiesRequest();
      final response = await client.listEmailIdentities(request);

      final identities = response.emailidentities.map((i) {
        return SESIdentity(
          identityName: i.identityname,
          identityType: i.identitytype.name.replaceAll('IDENTITY_TYPE_', ''),
          verificationStatus: i.verificationstatus.name.replaceAll(
            'VERIFICATION_STATUS_',
            '',
          ),
          sendingEnabled: i.sendingenabled,
        );
      }).toList();

      state = state.copyWith(identities: identities, isLoading: false);
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }
}

final sesListProvider = NotifierProvider<SESListNotifier, SESState>(() {
  return SESListNotifier();
});
