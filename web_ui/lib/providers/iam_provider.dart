import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../src/generated/iam.pbgrpc.dart';
import 'api_providers.dart' show iamProvider;

class IAMUser {
  final String arn;
  final String userName;
  final String userId;
  final DateTime created;
  final String status;
  final List<String> groups;

  IAMUser({
    required this.arn,
    required this.userName,
    required this.userId,
    required this.created,
    this.status = '',
    this.groups = const [],
  });
}

class IAMRole {
  final String arn;
  final String roleName;
  final String roleId;
  final DateTime created;
  final String path;

  IAMRole({
    required this.arn,
    required this.roleName,
    required this.roleId,
    required this.created,
    this.path = '',
  });
}

class IAMPolicy {
  final String arn;
  final String policyName;
  final String type;
  final int attachedCount;
  final String? description;

  IAMPolicy({
    required this.arn,
    required this.policyName,
    required this.type,
    this.attachedCount = 0,
    this.description,
  });
}

class IAMState {
  final List<IAMUser> users;
  final List<IAMRole> roles;
  final List<IAMPolicy> policies;
  final bool isLoading;
  final String? error;
  final int selectedTabIndex;

  const IAMState({
    this.users = const [],
    this.roles = const [],
    this.policies = const [],
    this.isLoading = false,
    this.error,
    this.selectedTabIndex = 0,
  });

  IAMState copyWith({
    List<IAMUser>? users,
    List<IAMRole>? roles,
    List<IAMPolicy>? policies,
    bool? isLoading,
    String? error,
    bool clearError = false,
    int? selectedTabIndex,
  }) {
    return IAMState(
      users: users ?? this.users,
      roles: roles ?? this.roles,
      policies: policies ?? this.policies,
      isLoading: isLoading ?? this.isLoading,
      error: clearError ? null : (error ?? this.error),
      selectedTabIndex: selectedTabIndex ?? this.selectedTabIndex,
    );
  }
}

final iamListProvider = NotifierProvider<IAMNotifier, IAMState>(() {
  return IAMNotifier();
});

class IAMNotifier extends Notifier<IAMState> {
  @override
  IAMState build() => const IAMState();

  IAMServiceClient get _client => ref.read(iamProvider);

  Future<void> loadData() async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      final usersResponse = await _client.listUsers(ListUsersRequest());
      final rolesResponse = await _client.listRoles(ListRolesRequest());

      final users = usersResponse.users
          .map(
            (u) => IAMUser(
              arn: u.arn,
              userName: u.username,
              userId: u.userid,
              created: DateTime.tryParse(u.createdate) ?? DateTime.now(),
              status: 'Active',
            ),
          )
          .toList();

      final roles = rolesResponse.roles
          .map(
            (r) => IAMRole(
              arn: r.arn,
              roleName: r.rolename,
              roleId: r.roleid,
              created: DateTime.tryParse(r.createdate) ?? DateTime.now(),
              path: r.path,
            ),
          )
          .toList();

      state = state.copyWith(users: users, roles: roles, isLoading: false);
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  Future<void> createUser(String userName) async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      await _client.createUser(CreateUserRequest(username: userName));
      await loadData();
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  Future<void> deleteUser(String userName) async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      await _client.deleteUser(DeleteUserRequest(username: userName));
      state = state.copyWith(
        users: state.users.where((u) => u.userName != userName).toList(),
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  Future<void> deleteRole(String roleName) async {
    state = state.copyWith(isLoading: true, clearError: true);
    try {
      await _client.deleteRole(DeleteRoleRequest(rolename: roleName));
      state = state.copyWith(
        roles: state.roles.where((r) => r.roleName != roleName).toList(),
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }

  void setSelectedTab(int index) {
    state = state.copyWith(selectedTabIndex: index);
  }
}
