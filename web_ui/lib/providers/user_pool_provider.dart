import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../src/generated/cognito_idp.pbgrpc.dart';
import 'api_providers.dart' show cognitoIdpProvider;

class UserPool {
  final String id;
  final String name;
  final String status;
  final int usersCount;
  final DateTime created;

  UserPool({
    required this.id,
    required this.name,
    this.status = '',
    this.usersCount = 0,
    required this.created,
  });
}

class CognitoUser {
  final String username;
  final String email;
  final String status;
  final DateTime created;
  final DateTime? lastModified;
  final List<String> groups;
  CognitoUser({
    required this.username,
    this.email = '',
    this.status = '',
    required this.created,
    this.lastModified,
    this.groups = const [],
  });
}

class CognitoState {
  final List<UserPool> userPools;
  final List<CognitoUser> users;
  final UserPool? selectedPool;
  final bool isLoading;
  final String? error;
  const CognitoState({
    this.userPools = const [],
    this.users = const [],
    this.selectedPool,
    this.isLoading = false,
    this.error,
  });
  CognitoState copyWith({
    List<UserPool>? userPools,
    List<CognitoUser>? users,
    UserPool? selectedPool,
    bool clearSelectedPool = false,
    bool? isLoading,
    String? error,
    bool clearError = false,
  }) {
    return CognitoState(
      userPools: userPools ?? this.userPools,
      users: users ?? this.users,
      selectedPool: clearSelectedPool
          ? null
          : (selectedPool ?? this.selectedPool),
      isLoading: isLoading ?? this.isLoading,
      error: clearError ? null : (error ?? this.error),
    );
  }
}

final cognitoListProvider = NotifierProvider<CognitoNotifier, CognitoState>(() {
  return CognitoNotifier();
});

class CognitoNotifier extends Notifier<CognitoState> {
  @override
  CognitoState build() => const CognitoState();
  CognitoIdentityProviderServiceClient get _client =>
      ref.read(cognitoIdpProvider);
  Future<void> loadUserPools() async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      final response = await _client.listUserPools(ListUserPoolsRequest());
      final pools = response.userpools
          .map(
            (p) => UserPool(
              id: p.id,
              name: p.name,
              status: p.status.name.replaceAll('STATUS_TYPE_', ''),
              created: DateTime.tryParse(p.creationdate) ?? DateTime.now(),
            ),
          )
          .toList();
      state = state.copyWith(userPools: pools, isLoading: false);
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  Future<void> createUserPool(String name) async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      await _client.createUserPool(CreateUserPoolRequest(poolname: name));
      await loadUserPools();
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  void selectPool(UserPool pool) {
    state = state.copyWith(selectedPool: pool, users: []);
  }
}
