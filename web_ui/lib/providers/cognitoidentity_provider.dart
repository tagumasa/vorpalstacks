import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../src/generated/cognito_identity.pbgrpc.dart';
import 'api_providers.dart';

class CognitoIdentityPool {
  final String identityPoolId;
  final String identityPoolName;

  CognitoIdentityPool({
    required this.identityPoolId,
    required this.identityPoolName,
  });
}

class CognitoIdentityState {
  final List<CognitoIdentityPool> identityPools;
  final bool isLoading;
  final String? error;

  const CognitoIdentityState({
    this.identityPools = const [],
    this.isLoading = false,
    this.error,
  });

  CognitoIdentityState copyWith({
    List<CognitoIdentityPool>? identityPools,
    bool? isLoading,
    String? error,
    bool clearError = false,
  }) {
    return CognitoIdentityState(
      identityPools: identityPools ?? this.identityPools,
      isLoading: isLoading ?? this.isLoading,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

class CognitoIdentityListNotifier extends Notifier<CognitoIdentityState> {
  @override
  CognitoIdentityState build() => const CognitoIdentityState();

  Future<void> loadIdentityPools() async {
    state = state.copyWith(isLoading: true, clearError: true);

    try {
      final client = ref.read(cognitoIdentityProvider);
      final request = ListIdentityPoolsInput(maxresults: 60);
      final response = await client.listIdentityPools(request);

      final pools = response.identitypools.map((p) {
        return CognitoIdentityPool(
          identityPoolId: p.identitypoolid,
          identityPoolName: p.identitypoolname,
        );
      }).toList();

      state = state.copyWith(identityPools: pools, isLoading: false);
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }
}

final cognitoIdentityListProvider =
    NotifierProvider<CognitoIdentityListNotifier, CognitoIdentityState>(() {
      return CognitoIdentityListNotifier();
    });
