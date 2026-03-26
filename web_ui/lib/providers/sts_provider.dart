import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../src/generated/sts.pbgrpc.dart';
import 'api_providers.dart';

class STSIdentity {
  final String userId;
  final String arn;
  final String account;

  STSIdentity({required this.userId, required this.arn, required this.account});
}

class STSState {
  final STSIdentity? identity;
  final bool isLoading;
  final String? error;

  const STSState({this.identity, this.isLoading = false, this.error});

  STSState copyWith({
    STSIdentity? identity,
    bool? isLoading,
    String? error,
    bool clearError = false,
  }) {
    return STSState(
      identity: identity ?? this.identity,
      isLoading: isLoading ?? this.isLoading,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

class STSNotifier extends Notifier<STSState> {
  @override
  STSState build() => const STSState();

  Future<void> loadCallerIdentity() async {
    state = state.copyWith(isLoading: true, clearError: true);

    try {
      final client = ref.read(stsProvider);
      final request = GetCallerIdentityRequest();
      final response = await client.getCallerIdentity(request);

      final identity = STSIdentity(
        userId: response.userid,
        arn: response.arn,
        account: response.account,
      );

      state = state.copyWith(identity: identity, isLoading: false);
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }
}

final stsListProvider = NotifierProvider<STSNotifier, STSState>(() {
  return STSNotifier();
});
