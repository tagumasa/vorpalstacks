import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

import '../../providers/api_providers.dart';
import '../generated/admin_auth.pb.dart';

class AuthState {
  final bool isAuthenticated;
  final bool isLoading;
  final String? username;
  final String? accessToken;
  final String? refreshToken;
  final String? error;

  const AuthState({
    required this.isAuthenticated,
    required this.isLoading,
    this.username,
    this.accessToken,
    this.refreshToken,
    this.error,
  });

  factory AuthState.initial() =>
      const AuthState(isAuthenticated: false, isLoading: true);

  AuthState copyWith({
    bool? isAuthenticated,
    bool? isLoading,
    String? username,
    String? accessToken,
    String? refreshToken,
    String? error,
  }) {
    return AuthState(
      isAuthenticated: isAuthenticated ?? this.isAuthenticated,
      isLoading: isLoading ?? this.isLoading,
      username: username ?? this.username,
      accessToken: accessToken ?? this.accessToken,
      refreshToken: refreshToken ?? this.refreshToken,
      error: error,
    );
  }
}

class AuthNotifier extends Notifier<AuthState> {
  late final FlutterSecureStorage _storage;

  @override
  AuthState build() {
    _storage = const FlutterSecureStorage(
      webOptions: WebOptions(
        dbName: 'VorpalStacks',
        publicKey: 'VorpalStacksSecureStorage',
      ),
    );
    _loadStoredTokens();
    return AuthState.initial();
  }

  Future<void> _loadStoredTokens() async {
    try {
      final accessToken = await _storage.read(key: 'access_token');
      final refreshToken = await _storage.read(key: 'refresh_token');

      if (accessToken != null && refreshToken != null) {
        try {
          final client = ref.read(adminAuthProvider);
          final response = await client.getCurrentUser(
            GetCurrentUserRequest(accessToken: accessToken),
          );
          state = AuthState(
            isAuthenticated: true,
            isLoading: false,
            username: response.username,
            accessToken: accessToken,
            refreshToken: refreshToken,
          );
          return;
        } catch (_) {
          await _refreshTokens(refreshToken);
          return;
        }
      }
    } catch (_) {}

    state = const AuthState(isAuthenticated: false, isLoading: false);
  }

  Future<bool> login(String username, String password) async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final client = ref.read(adminAuthProvider);
      final response = await client.login(
        LoginRequest(username: username, password: password),
      );
      await _storage.write(key: 'access_token', value: response.accessToken);
      await _storage.write(key: 'refresh_token', value: response.refreshToken);

      state = AuthState(
        isAuthenticated: true,
        isLoading: false,
        username: username,
        accessToken: response.accessToken,
        refreshToken: response.refreshToken,
      );
      return true;
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
      return false;
    }
  }

  Future<void> logout() async {
    try {
      if (state.accessToken != null) {
        final client = ref.read(adminAuthProvider);
        await client.logout(LogoutRequest(accessToken: state.accessToken!));
      }
    } catch (_) {}

    await _storage.deleteAll();
    state = const AuthState(isAuthenticated: false, isLoading: false);
  }

  Future<void> _refreshTokens(String refreshToken) async {
    try {
      final client = ref.read(adminAuthProvider);
      final response = await client.refreshToken(
        RefreshTokenRequest(refreshToken: refreshToken),
      );
      await _storage.write(key: 'access_token', value: response.accessToken);

      state = state.copyWith(
        isAuthenticated: true,
        isLoading: false,
        accessToken: response.accessToken,
      );
    } catch (_) {
      await logout();
    }
  }
}

final authProvider = NotifierProvider<AuthNotifier, AuthState>(() {
  return AuthNotifier();
});
